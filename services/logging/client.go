package logging

import (
	"fmt"
	"os"
	"runtime"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/meta"
	"github.com/ionos-cloud/sdk-go-bundle/products/logging/v2"
	"github.com/ionos-cloud/sdk-go-bundle/shared"
	"github.com/ionos-cloud/sdk-go-bundle/shared/fileconfiguration"

	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils/bundle"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils/constant"
)

// Client is a wrapper for the Logging SDK client
type Client struct {
	sdkClient  logging.APIClient
	fileConfig *fileconfiguration.FileConfig
}

// GetConfig returns the configuration
func (c *Client) GetConfig() *shared.Configuration {
	return c.sdkClient.GetConfig()
}

// GetFileConfig returns configuration read from the file
func (c *Client) GetFileConfig() *fileconfiguration.FileConfig {
	return c.fileConfig
}

// NewClient creates a new Logging client
func NewClient(clientOptions bundle.ClientOptions, fileConfig *fileconfiguration.FileConfig) *Client {
	config := shared.NewConfigurationFromOptions(clientOptions.ClientOptions)
	config.MaxRetries = constant.MaxRetries
	config.MaxWaitTime = constant.MaxWaitTime
	config.UserAgent = fmt.Sprintf(
		"terraform-provider/%s_ionos-cloud-sdk-go-logging/%s_hashicorp-terraform/%s_terraform-plugin-sdk/%s_os/%s_arch",
		logging.Version, clientOptions.TerraformVersion, meta.SDKVersionString(), runtime.GOOS, runtime.GOARCH) //nolint:staticcheck

	client := &Client{sdkClient: *logging.NewAPIClient(config), fileConfig: fileConfig}

	// override client with location from config file if it exists and no global endpoint it set
	// todo cguran - remove after testing
	// if clientOptions.Endpoint == "" && overridesCloudAPI != nil && len(overridesCloudAPI.Endpoints) > 0 {
	//	for _, endpoint := range overridesCloudAPI.Endpoints {
	//		endpointLocation := endpoint.ChangeConfigURL
	//		if sdkLocation, ok := location.TerraformToSDK[endpoint.ChangeConfigURL]; ok {
	//			endpointLocation = sdkLocation
	//		}
	//		replaceServers := len(overridesCloudAPI.Endpoints) == 1 && endpointLocation == ""
	//		shared.OverrideLocationFor(&client.sdkClient, endpointLocation, endpoint.Name, replaceServers)
	//	}
	//}
	return client
}

// ChangeConfigURL changes the configuration URL based on the location
func (c *Client) ChangeConfigURL(location string) {
	config := c.sdkClient.GetConfig()
	if location == "" && os.Getenv(ionosAPIURLLogging) != "" {
		c.GetConfig().Servers = shared.ServerConfigurations{
			{
				URL: utils.CleanURL(os.Getenv(ionosAPIURLLogging)),
			},
		}
		return
	}
	for _, server := range config.Servers {
		if strings.EqualFold(server.Description, shared.EndpointOverridden+location) || strings.EqualFold(server.URL, locationToURL[location]) {
			config.Servers = shared.ServerConfigurations{
				{
					URL:         server.URL,
					Description: shared.EndpointOverridden + location,
				},
			}
			return
		}
	}
}

var (
	ionosAPIURLLogging = "IONOS_API_URL_LOGGING"

	AvailableLocations = []string{"de/fra", "de/txl", "es/vit", "gb/lhr", "fr/par"}
	// DefaultLocation is the default logging pipeline location
	DefaultLocation = "de/txl"
	locationToURL   = map[string]string{
		"":       "https://logging.de-txl.ionos.com",
		"de/fra": "https://logging.de-fra.ionos.com",
		"de/txl": "https://logging.de-txl.ionos.com",
		"es/vit": "https://logging.es-vit.ionos.com",
		"gb/lhr": "https://logging.gb-lhr.ionos.com",
		"fr/par": "https://logging.fr-par.ionos.com",
	}
)
