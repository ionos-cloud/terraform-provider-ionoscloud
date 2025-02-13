package monitoring

import (
	"fmt"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils/bundle"
	"net/http"
	"os"
	"runtime"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/meta"
	monitoring "github.com/ionos-cloud/sdk-go-bundle/products/monitoring/v2"
	"github.com/ionos-cloud/sdk-go-bundle/shared"

	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils/constant"
)

// Client is a wrapper for the Monitoring SDK
type Client struct { //nolint:golint
	sdkClient    monitoring.APIClient
	loadedConfig *shared.LoadedConfig
}

func (c *Client) GetConfig() *shared.Configuration {
	return c.sdkClient.GetConfig()
}

func (c *Client) GetLoadedConfig() *shared.LoadedConfig {
	return c.loadedConfig
}

func NewClient(clientOptions bundle.ClientOptions, sharedLoadedConfig *shared.LoadedConfig) *Client {
	config := shared.NewConfiguration(clientOptions.Credentials.Username, clientOptions.Credentials.Password,
		clientOptions.Credentials.Token, clientOptions.Endpoint)
	config.MaxRetries = constant.MaxRetries
	config.MaxWaitTime = constant.MaxWaitTime
	config.HTTPClient = &http.Client{Transport: utils.CreateTransport(clientOptions.SkipTLSVerify)}
	config.UserAgent = fmt.Sprintf(
		"terraform-provider/ionos-cloud-sdk-go-monitoring/%s_hashicorp-terraform/%s_terraform-plugin-sdk/%s_os/%s_arch/%s",
		monitoring.Version, clientOptions.TerraformVersion, meta.SDKVersionString(), runtime.GOOS, runtime.GOARCH) // nolint:staticcheck

	return &Client{sdkClient: *monitoring.NewAPIClient(config),
		loadedConfig: sharedLoadedConfig}
}

var (
	ionosAPIURLMonitoring = "IONOS_API_URL_MONITORING"
	locationToURL         = map[string]string{
		"":       "https://monitoring.de-fra.ionos.com",
		"de/fra": "https://monitoring.de-fra.ionos.com",
		"de/txl": "https://monitoring.de-txl.ionos.com",
		"es/vit": "https://monitoring.es-vit.ionos.com",
		"gb/lhr": "https://monitoring.gb-lhr.ionos.com",
		"fr/par": "https://monitoring.fr-par.ionos.com",
	}
)

// todo cguran - check for description and override only after
//func (c *Client) ChangeConfigURL(location string) {
//	config := c.sdkClient.GetConfig()
//	if location == "" && os.Getenv(ionosAPIURLMonitoring) != "" {
//		config.Servers = shared.ServerConfigurations{
//			{
//				URL: utils.CleanURL(os.Getenv(ionosAPIURLMonitoring)),
//			},
//		}
//		return
//	}
//	config.Servers = shared.ServerConfigurations{
//		{
//			URL: locationToURL[location],
//		},
//	}
//}

func (c *Client) ChangeConfigURL(location string) {
	config := c.sdkClient.GetConfig()
	if location == "" && os.Getenv(ionosAPIURLMonitoring) != "" {
		config.Servers = shared.ServerConfigurations{
			{
				URL: utils.CleanURL(os.Getenv(ionosAPIURLMonitoring)),
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
