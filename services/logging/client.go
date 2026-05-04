package logging

import (
	"context"
	"fmt"
	"os"
	"runtime"
	"strings"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/meta"
	"github.com/ionos-cloud/sdk-go-bundle/products/logging/v2"
	"github.com/ionos-cloud/sdk-go-bundle/shared"
	"github.com/ionos-cloud/sdk-go-bundle/shared/fileconfiguration"

	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/clientoptions"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils/configlog"
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
func NewClient(clientOptions clientoptions.TerraformClientOptions, fileConfig *fileconfiguration.FileConfig) *Client {
	config := shared.NewConfigurationFromOptions(clientOptions.ClientOptions)
	config.MaxRetries = constant.MaxRetries
	config.MaxWaitTime = constant.MaxWaitTime
	config.UserAgent = fmt.Sprintf(
		"terraform-provider/%s_ionos-cloud-sdk-go-logging/%s_hashicorp-terraform/%s_terraform-plugin-sdk/%s_os/%s_arch/%s",
		clientOptions.Version, logging.Version, clientOptions.TerraformVersion,
		meta.SDKVersionString(), runtime.GOOS, runtime.GOARCH, //nolint:staticcheck
	)

	client := &Client{sdkClient: *logging.NewAPIClient(config), fileConfig: fileConfig}
	return client
}

// ChangeConfigURL changes the configuration URL based on the location
func (c *Client) ChangeConfigURL(ctx context.Context, location string) {
	config := c.sdkClient.GetConfig()
	if location == "" && os.Getenv(ionosAPIURLLogging) != "" {
		url := utils.CleanURL(os.Getenv(ionosAPIURLLogging))
		tflog.Debug(ctx, "Logging: endpoint from env", map[string]interface{}{"env": ionosAPIURLLogging, "url": url})
		c.GetConfig().Servers = shared.ServerConfigurations{
			{
				URL: url,
			},
		}
		return
	}
	for _, server := range config.Servers {
		if strings.EqualFold(server.Description, shared.EndpointOverridden+location) || strings.EqualFold(server.URL, locationToURL[location]) {
			tflog.Debug(ctx, "Logging: endpoint for location", map[string]interface{}{"location": configlog.FormatLocation(location), "url": server.URL})
			config.Servers = shared.ServerConfigurations{
				{
					URL:         server.URL,
					Description: shared.EndpointOverridden + location,
				},
			}
			return
		}
	}
	tflog.Debug(ctx, "Logging: endpoint for location", map[string]interface{}{"location": configlog.FormatLocation(location), "url": locationToURL[location]})
}

var (
	ionosAPIURLLogging = "IONOS_API_URL_LOGGING"

	AvailableLocations = []string{"de/fra", "de/fra/2", "de/txl", "es/vit", "gb/lhr", "fr/par"}
	// DefaultLocation is the default logging pipeline location
	DefaultLocation = "de/txl"
	locationToURL   = map[string]string{
		"":         "https://logging.de-txl.ionos.com",
		"de/fra":   "https://logging.de-fra.ionos.com",
		"de/fra/2": "https://logging.de-fra.ionos.com",
		"de/txl":   "https://logging.de-txl.ionos.com",
		"es/vit":   "https://logging.es-vit.ionos.com",
		"gb/lhr":   "https://logging.gb-lhr.ionos.com",
		"fr/par":   "https://logging.fr-par.ionos.com",
	}
)
