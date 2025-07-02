package objectstorage

import (
	"fmt"
	"net/http"
	"runtime"

	"github.com/hashicorp/terraform-plugin-sdk/v2/meta"
	objstorage "github.com/ionos-cloud/sdk-go-bundle/products/userobjectstorage/v2"
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

// const ionosAPIURLObjectStorage = "IONOS_API_URL_USER_OBJECT_STORAGE"

// GetBaseClient returns the base client.
func (c *Client) GetBaseClient() *objstorage.APIClient {
	return c.client
}

// NewClient creates a new Object Storage client with the given credentials and region.
func NewClient(clientOptions clientoptions.TerraformClientOptions, config *fileconfiguration.FileConfig) *Client {
	// Set custom endpoint if provided
	// if envValue := os.Getenv(ionosAPIURLObjectStorage); envValue != "" {
	// 	clientOptions.Endpoint = envValue
	// }
	certificateAuthData := ""
	if clientOptions.Endpoint == "" {
		// todo change to fileconfiguration.UserObjectStorage
		if endpointOverrides := config.GetProductLocationOverrides(fileconfiguration.ObjectStorage, clientOptions.StorageOptions.Region); endpointOverrides != nil {
			clientOptions.Endpoint = endpointOverrides.Name
			if !clientOptions.SkipTLSVerify {
				clientOptions.SkipTLSVerify = endpointOverrides.SkipTLSVerify
			}
			certificateAuthData = endpointOverrides.CertificateAuthData
		}
	}
	cfg := shared.NewConfigurationFromOptions(clientOptions.ClientOptions)

	cfg.HTTPClient = &http.Client{Transport: shared.CreateTransport(clientOptions.SkipTLSVerify, certificateAuthData)}

	cfg.MiddlewareWithError = shared.SignerMiddleware("s3-eu-central-1", "s3", clientOptions.StorageOptions.AccessKey, clientOptions.StorageOptions.SecretKey)
	cfg.UserAgent = fmt.Sprintf(
		"terraform-provider/%s_ionos-cloud-sdk-go-user-object-storage/%s_hashicorp-terraform/%s_terraform-plugin-sdk/%s_os/%s_arch/%s",
		clientOptions.Version, objstorage.Version, clientOptions.TerraformVersion,
		meta.SDKVersionString(), runtime.GOOS, runtime.GOARCH, //nolint:staticcheck
	)

	return &Client{
		client:     objstorage.NewAPIClient(cfg),
		fileConfig: config,
	}
}
