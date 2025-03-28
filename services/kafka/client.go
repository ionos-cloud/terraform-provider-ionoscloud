package kafka

import (
	"fmt"
	"net/http"
	"os"
	"runtime"

	"github.com/ionos-cloud/sdk-go-bundle/shared"

	"github.com/hashicorp/terraform-plugin-sdk/v2/meta"
	"github.com/ionos-cloud/sdk-go-bundle/shared/fileconfiguration"
	kafka "github.com/ionos-cloud/sdk-go-kafka"

	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/clientoptions"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils/constant"
)

// Client is a wrapper for the Kafka API Client
type Client struct {
	sdkClient  kafka.APIClient
	fileConfig *fileconfiguration.FileConfig
}

// GetFileConfig returns configuration read from the file
func (c *Client) GetFileConfig() *fileconfiguration.FileConfig {
	return c.fileConfig
}

// GetConfig returns the configuration
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

// NewClient creates a new Kafka client
func NewClient(clientOptions clientoptions.TerraformClientOptions, fileConfig *fileconfiguration.FileConfig) *Client {
	config := kafka.NewConfiguration(clientOptions.Credentials.Username, clientOptions.Credentials.Password, clientOptions.Credentials.Token, clientOptions.Endpoint)

	if os.Getenv(constant.IonosDebug) != "" {
		config.Debug = true
	}
	config.MaxRetries = constant.MaxRetries
	config.MaxWaitTime = constant.MaxWaitTime
	config.HTTPClient = http.DefaultClient
	config.HTTPClient.Transport = shared.CreateTransport(clientOptions.SkipTLSVerify, clientOptions.Certificate)
	config.UserAgent = fmt.Sprintf(
		"terraform-provider/%s_ionos-cloud-sdk-go-kafka/%s_hashicorp-terraform/%s_terraform-plugin-sdk/%s_os/%s_arch/%s",
		clientOptions.Version, kafka.Version, clientOptions.TerraformVersion,
		meta.SDKVersionString(), runtime.GOOS, runtime.GOARCH, //nolint:staticcheck
	)
	client := Client{sdkClient: *kafka.NewAPIClient(config),
		fileConfig: fileConfig,
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
