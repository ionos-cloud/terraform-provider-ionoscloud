package monitoring

import (
	"fmt"
	"os"
	"runtime"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/meta"
	"github.com/ionos-cloud/sdk-go-bundle/products/monitoring/v2"
	"github.com/ionos-cloud/sdk-go-bundle/shared"
	"github.com/ionos-cloud/sdk-go-bundle/shared/fileconfiguration"

	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/clientoptions"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils/constant"
)

// Client is a wrapper for the Monitoring SDK
type Client struct {
	sdkClient  monitoring.APIClient
	fileConfig *fileconfiguration.FileConfig
}

// GetConfig - returns the configuration
func (c *Client) GetConfig() *shared.Configuration {
	return c.sdkClient.GetConfig()
}

// GetFileConfig - returns the configuration read from the config file
func (c *Client) GetFileConfig() *fileconfiguration.FileConfig {
	return c.fileConfig
}

// NewClient creates a new Monitoring client
func NewClient(clientOptions clientoptions.TerraformClientOptions, fileConfig *fileconfiguration.FileConfig) *Client {
	config := shared.NewConfigurationFromOptions(clientOptions.ClientOptions)
	config.MaxRetries = constant.MaxRetries
	config.MaxWaitTime = constant.MaxWaitTime
	config.UserAgent = fmt.Sprintf(
		"terraform-provider/%s_ionos-cloud-sdk-go-monitoring/%s_hashicorp-terraform/%s_terraform-plugin-sdk/%s_os/%s_arch/%s",
		clientOptions.Version, monitoring.Version, clientOptions.TerraformVersion,
		meta.SDKVersionString(), runtime.GOOS, runtime.GOARCH, // nolint:staticcheck
	)

	return &Client{sdkClient: *monitoring.NewAPIClient(config),
		fileConfig: fileConfig}
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
