package objectstorage

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"runtime"
	"slices"
	"time"

	"github.com/aws/aws-sdk-go/aws/credentials"
	awsv4 "github.com/aws/aws-sdk-go/aws/signer/v4"
	"github.com/hashicorp/terraform-plugin-sdk/v2/meta"
	objstorage "github.com/ionos-cloud/sdk-go-bundle/products/objectstorage/v2"
	"github.com/ionos-cloud/sdk-go-bundle/shared"
	"github.com/ionos-cloud/sdk-go-bundle/shared/fileconfiguration"

	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/clientoptions"
)

// Client is a wrapper around the Object Storage client.
type Client struct {
	client       *objstorage.APIClient
	fileConfig   *fileconfiguration.FileConfig
	signer       *awsv4.Signer
	globalRegion string
}

// GetFileConfig returns the file configuration.
func (c *Client) GetFileConfig() *fileconfiguration.FileConfig {
	return c.fileConfig
}

const ionosAPIURLObjectStorage = "IONOS_API_URL_OBJECT_STORAGE"

// GetBaseClient returns the base client.
func (c *Client) GetBaseClient() *objstorage.APIClient {
	return c.client
}

// NewClient creates a new Object Storage client with the given credentials and region.
func NewClient(clientOptions clientoptions.TerraformClientOptions, config *fileconfiguration.FileConfig) *Client {
	var globalRegion string
	// Check if a global region was set using the "IONOS_S3_REGION" env var or using "s3_region" provider attribute.
	// If no global region was set, the requests will be routed according to the region set at the bucket level.
	if clientOptions.StorageOptions.Region != "" {
		globalRegion = clientOptions.StorageOptions.Region
	}
	// Set custom endpoint if provided
	if envValue := os.Getenv(ionosAPIURLObjectStorage); envValue != "" {
		clientOptions.Endpoint = envValue
	}
	certificateAuthData := ""
	if clientOptions.Endpoint == "" {
		if endpointOverrides := config.GetProductLocationOverrides(fileconfiguration.ObjectStorage, clientOptions.StorageOptions.Region); endpointOverrides != nil {
			clientOptions.Endpoint = endpointOverrides.Name
			if !clientOptions.SkipTLSVerify {
				clientOptions.SkipTLSVerify = endpointOverrides.SkipTLSVerify
			}
			certificateAuthData = endpointOverrides.CertificateAuthData
		}
	}

	if clientOptions.StorageOptions.Region == "" {
		clientOptions.StorageOptions.Region = "eu-central-3"
	}

	// TODO -- replace with WithObjectStorage from sdk-go-bundle
	cfg := shared.NewConfigurationFromOptions(clientOptions.ClientOptions)
	signer := awsv4.NewSigner(credentials.NewStaticCredentials(clientOptions.StorageOptions.AccessKey, clientOptions.StorageOptions.SecretKey, ""))
	cfg.MiddlewareWithError = signerMiddleware(clientOptions.StorageOptions.Region, signer)
	cfg.UserAgent = fmt.Sprintf(
		"terraform-provider/%s_ionos-cloud-sdk-go-object-storage/%s_hashicorp-terraform/%s_terraform-plugin-sdk/%s_os/%s_arch/%s",
		clientOptions.Version, objstorage.Version, clientOptions.TerraformVersion,
		meta.SDKVersionString(), runtime.GOOS, runtime.GOARCH, //nolint:staticcheck
	)

	cfg.HTTPClient = &http.Client{Transport: shared.CreateTransport(clientOptions.SkipTLSVerify, certificateAuthData)}

	return &Client{
		client:       objstorage.NewAPIClient(cfg),
		fileConfig:   config,
		signer:       signer,
		globalRegion: globalRegion,
	}
}

var (
	regionToURL = map[string]string{
		"":             "https://s3.eu-central-3.ionoscloud.com",
		"eu-central-3": "https://s3.eu-central-3.ionoscloud.com",
		"eu-central-4": "https://s3.eu-central-4.ionoscloud.com",
		"us-central-1": "https://s3.us-central-1.ionoscloud.com",
	}
	PublicS3Locations = []string{"eu-central-3", "eu-central-4", "us-central-1"}
)

func (c *Client) ChangeConfigURL(bucketRegion string) error {
	var region string

	if c.globalRegion != "" && c.globalRegion != bucketRegion {
		return fmt.Errorf("global region: %v is different from the bucket region: %v, if you want to use the bucket region, please unset the global region", c.globalRegion, bucketRegion)
	}
	region = bucketRegion
	if region != "" && !slices.Contains(PublicS3Locations, region) {
		// Custom region, don't overwrite the URL.
		return nil
	}

	cfg := c.client.GetConfig()
	cfg.Servers = shared.ServerConfigurations{
		{
			URL: regionToURL[region],
		},
	}
	cfg.MiddlewareWithError = signerMiddleware(region, c.signer)
	return nil
}

// signerMiddleware returns a middleware function that signs the request using the provided AWS v4 signer
// and region.
func signerMiddleware(region string, signer *awsv4.Signer) shared.MiddlewareFunctionWithError {
	return func(r *http.Request) error {
		var reader io.ReadSeeker
		if r.Body != nil {
			bodyBytes, err := io.ReadAll(r.Body)
			if err != nil {
				return err
			}
			reader = bytes.NewReader(bodyBytes)
		}

		_, err := signer.Sign(r, reader, "s3", region, time.Now())
		if errors.Is(err, credentials.ErrStaticCredentialsEmpty) {
			return errors.New("object storage credentials are missing. Please set s3_access_key and s3_secret_key provider attributes or environment variables IONOS_S3_ACCESS_KEY and IONOS_S3_SECRET_KEY")
		}
		return err
	}
}
