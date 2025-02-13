package kafka

import (
	"fmt"
	"github.com/ionos-cloud/sdk-go-bundle/shared"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils/bundle"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils/loadedconfig"
	"net/http"
	"os"
	"runtime"

	"github.com/hashicorp/terraform-plugin-sdk/v2/meta"

	kafka "github.com/ionos-cloud/sdk-go-kafka"

	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils/constant"
)

// Client is a wrapper for the Kafka API Client
type Client struct {
	sdkClient    kafka.APIClient
	loadedConfig *shared.LoadedConfig
}

func (c *Client) GetLoadedConfig() *shared.LoadedConfig {
	return c.loadedConfig
}
func (c *Client) GetConfig() *kafka.Configuration {
	return c.sdkClient.GetConfig()
}

var (
	// AvailableLocations is a list of available locations for Kafka
	AvailableLocations = []string{"de/fra", "de/txl" /*, "es/vit", "gb/lhr", "us/ewr", "us/las", "us/mci", "fr/par"*/}

	locationToURL = map[string]string{
		"":       "https://kafka.de-fra.ionos.com",
		"de/fra": "https://kafka.de-fra.ionos.com",
		"de/txl": "https://kafka.de-txl.ionos.com",
		// other locations not yet available. will be added in the future.
		// "es/vit": "https://kafka.es-vit.ionos.com",
		// "gb/lhr": "https://kafka.gb-lhr.ionos.com",
		// "us/ewr": "https://kafka.us-ewr.ionos.com",
		// "us/las": "https://kafka.us-las.ionos.com",
		// "us/mci": "https://kafka.us-mci.ionos.com",
		// "fr/par": "https://kafka.fr-par.ionos.com",
	}
)

func NewClient(clientOptions bundle.ClientOptions, sharedLoadedConfig *shared.LoadedConfig) *Client {
	loadedconfig.SetClientOptionsFromLoadedConfig(&clientOptions, sharedLoadedConfig, shared.Kafka)

	config := kafka.NewConfiguration(clientOptions.Credentials.Username, clientOptions.Credentials.Password, clientOptions.Credentials.Token, clientOptions.Endpoint)

	if os.Getenv(constant.IonosDebug) != "" {
		config.Debug = true
	}
	config.MaxRetries = constant.MaxRetries
	config.MaxWaitTime = constant.MaxWaitTime
	config.HTTPClient = &http.Client{Transport: utils.CreateTransport(clientOptions.SkipTLSVerify)}
	config.UserAgent = fmt.Sprintf(
		"terraform-provider/_ionos-cloud-sdk-go-kafka/%s_hashicorp-terraform/%s_terraform-plugin-sdk/%s_os/%s_arch/%s",
		kafka.Version, clientOptions.TerraformVersion, meta.SDKVersionString(), runtime.GOOS, runtime.GOARCH, //nolint:staticcheck
	)
	client := Client{sdkClient: *kafka.NewAPIClient(config),
		loadedConfig: sharedLoadedConfig,
	}
	return &client
}
func (c *Client) changeConfigURL(location string) {
	config := c.sdkClient.GetConfig()
	config.Servers = kafka.ServerConfigurations{
		{
			URL: locationToURL[location],
		},
	}
}
