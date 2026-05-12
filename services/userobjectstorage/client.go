package userobjectstorage

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"runtime"
	"time"

	"github.com/aws/aws-sdk-go/aws/credentials"
	awsv4 "github.com/aws/aws-sdk-go/aws/signer/v4"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/meta"
	userobjectstorage "github.com/ionos-cloud/sdk-go-bundle/products/userobjectstorage/v2"
	"github.com/ionos-cloud/sdk-go-bundle/shared"

	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/clientoptions"
)

// Client holds the configuration needed to create per-region API clients.
// A new SDK client is instantiated for each operation so concurrent operations
// on different regions do not share mutable state.
type Client struct {
	clientOptions clientoptions.TerraformClientOptions
	httpClient    *http.Client
	signer        *awsv4.Signer
	userAgent     string
}

var regionToURL = map[string]string{
	"de":           "https://s3.eu-central-1.ionoscloud.com",
	"eu-central-2": "https://s3.eu-central-2.ionoscloud.com",
	"eu-south-2":   "https://s3.eu-south-2.ionoscloud.com",
}

// ValidRegions holds the accepted region values for schema validation.
var ValidRegions = []string{"de", "eu-central-2", "eu-south-2"}

// DefaultRegion is used when no region is specified (e.g. during import).
const DefaultRegion = "de"

// NewClient creates a new user object storage client factory.
func NewClient(ctx context.Context, clientOptions clientoptions.TerraformClientOptions) *Client {
	tflog.Debug(ctx, "User Object Storage: configuring client")

	signer := awsv4.NewSigner(credentials.NewStaticCredentials(
		clientOptions.StorageOptions.AccessKey,
		clientOptions.StorageOptions.SecretKey,
		"",
	))
	userAgent := fmt.Sprintf(
		"terraform-provider/%s_ionos-cloud-sdk-go-user-object-storage/%s_hashicorp-terraform/%s_terraform-plugin-sdk/%s_os/%s_arch/%s",
		clientOptions.Version, userobjectstorage.Version, clientOptions.TerraformVersion,
		meta.SDKVersionString(), runtime.GOOS, runtime.GOARCH, //nolint:staticcheck
	)

	return &Client{
		clientOptions: clientOptions,
		httpClient:    &http.Client{Transport: shared.CreateTransport(clientOptions.SkipTLSVerify, "")},
		signer:        signer,
		userAgent:     userAgent,
	}
}

// apiClientForRegion returns a new SDK API client configured for the given region.
// An empty region defaults to DefaultRegion.
func (c *Client) apiClientForRegion(region string) (*userobjectstorage.APIClient, error) {
	if region == "" {
		region = DefaultRegion
	}
	url, ok := regionToURL[region]
	if !ok {
		return nil, fmt.Errorf("unsupported region %q: must be one of %v", region, ValidRegions)
	}
	cfg := shared.NewConfigurationFromOptions(c.clientOptions.ClientOptions)
	cfg.Servers = shared.ServerConfigurations{{URL: url}}
	cfg.MiddlewareWithError = signerMiddleware(region, c.signer)
	cfg.UserAgent = c.userAgent
	cfg.HTTPClient = c.httpClient
	return userobjectstorage.NewAPIClient(cfg), nil
}

// GetBaseClient returns an SDK client for the default region. Used by acceptance tests.
func (c *Client) GetBaseClient() *userobjectstorage.APIClient {
	client, _ := c.apiClientForRegion(DefaultRegion)
	return client
}

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
			return errors.New("user object storage credentials are missing. Please set s3_access_key and s3_secret_key provider attributes or environment variables IONOS_S3_ACCESS_KEY and IONOS_S3_SECRET_KEY")
		}
		return err
	}
}
