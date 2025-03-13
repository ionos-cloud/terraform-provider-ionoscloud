package objectstorage

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"runtime"
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
	client     *objstorage.APIClient
	fileConfig *fileconfiguration.FileConfig
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
	cfg := shared.NewConfigurationFromOptions(clientOptions.ClientOptions)
	signer := awsv4.NewSigner(credentials.NewStaticCredentials(clientOptions.StorageOptions.AccessKey, clientOptions.StorageOptions.SecretKey, ""))
	cfg.MiddlewareWithError = func(r *http.Request) error {
		var reader io.ReadSeeker
		if r.Body != nil {
			bodyBytes, err := io.ReadAll(r.Body)
			if err != nil {
				return err
			}
			reader = bytes.NewReader(bodyBytes)
		}

		if clientOptions.StorageOptions.Region == "" {
			clientOptions.StorageOptions.Region = "eu-central-3"
		}
		_, err := signer.Sign(r, reader, "s3", clientOptions.StorageOptions.Region, time.Now())
		if errors.Is(err, credentials.ErrStaticCredentialsEmpty) {
			return errors.New("object storage credentials are missing. Please set s3_access_key and s3_secret_key provider attributes or environment variables IONOS_S3_ACCESS_KEY and IONOS_S3_SECRET_KEY")
		}
		return err
	}
	cfg.UserAgent = fmt.Sprintf(
		"terraform-provider/%s_ionos-cloud-sdk-go-object-storage/%s_hashicorp-terraform/%s_terraform-plugin-sdk/%s_os/%s_arch/%s",
		clientOptions.Version, "1.1.0", clientOptions.TerraformVersion,
		meta.SDKVersionString(), runtime.GOOS, runtime.GOARCH, //nolint:staticcheck
	)

	cfg.HTTPClient = &http.Client{Transport: shared.CreateTransport(clientOptions.SkipTLSVerify, certificateAuthData)}
	return &Client{
		client:     objstorage.NewAPIClient(cfg),
		fileConfig: config,
	}
}
