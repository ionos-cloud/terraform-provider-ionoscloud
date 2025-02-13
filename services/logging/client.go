package logging

import (
	"fmt"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils/bundle"
	"net/http"
	"os"
	"runtime"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/meta"
	"github.com/ionos-cloud/sdk-go-bundle/products/logging/v2"
	"github.com/ionos-cloud/sdk-go-bundle/shared"

	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils/constant"
)

// Client is a wrapper for the Logging SDK client
type Client struct {
	sdkClient    logging.APIClient
	loadedConfig *shared.LoadedConfig
}

func NewClient(clientOptions bundle.ClientOptions, loadedConfig *shared.LoadedConfig) *Client {
	config := shared.NewConfiguration(clientOptions.Credentials.Username, clientOptions.Credentials.Password, clientOptions.Credentials.Token, clientOptions.Endpoint)
	config.MaxRetries = constant.MaxRetries
	config.MaxWaitTime = constant.MaxWaitTime
	config.HTTPClient = &http.Client{Transport: utils.CreateTransport(clientOptions.SkipTLSVerify)}
	config.UserAgent = fmt.Sprintf(
		"terraform-provider/%s_ionos-cloud-sdk-go-logging/%s_hashicorp-terraform/%s_terraform-plugin-sdk/%s_os/%s_arch",
		logging.Version, clientOptions.TerraformVersion, meta.SDKVersionString(), runtime.GOOS, runtime.GOARCH) //nolint:staticcheck

	client := &Client{sdkClient: *logging.NewAPIClient(config), loadedConfig: loadedConfig}

	// override client with location from config file if it exists and no global endpoint it set
	// todo cguran - remove after testing
	//if clientOptions.Endpoint == "" && overridesCloudAPI != nil && len(overridesCloudAPI.Endpoints) > 0 {
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

func (c *Client) GetConfig() *shared.Configuration {
	return c.sdkClient.GetConfig()
}

func (c *Client) GetLoadedConfig() *shared.LoadedConfig {
	return c.loadedConfig
}

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
