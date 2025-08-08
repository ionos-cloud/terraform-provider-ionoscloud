package kafka

import (
	"fmt"
	"net/http"
	"runtime"

	"github.com/hashicorp/terraform-plugin-sdk/v2/meta"
	kafka "github.com/ionos-cloud/sdk-go-bundle/products/kafka/v2"
	"github.com/ionos-cloud/sdk-go-bundle/shared"
	"github.com/ionos-cloud/sdk-go-bundle/shared/fileconfiguration"

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
func (c *Client) GetConfig() *shared.Configuration {
	return c.sdkClient.GetConfig()
}

var (
	// AvailableLocations is a list of available locations for Kafka
	AvailableLocations = []string{"de/fra", "de/fra/2", "de/txl", "fr/par", "es/vit", "gb/lhr", "gb/bhx", "us/las", "us/mci", "us/ewr"}

	locationToURL = map[string]string{
		"":         "https://kafka.de-fra.ionos.com",
		"de/fra":   "https://kafka.de-fra.ionos.com",
		"de/fra/2": "https://kafka.de-fra.ionos.com",
		"de/txl":   "https://kafka.de-txl.ionos.com",
		"fr/par":   "https://kafka.fr-par.ionos.com",
		"es/vit":   "https://kafka.es-vit.ionos.com",
		"gb/lhr":   "https://kafka.gb-lhr.ionos.com",
		"gb/bhx":   "https://kafka.gb-bhx.ionos.com",
		"us/las":   "https://kafka.us-las.ionos.com",
		"us/mci":   "https://kafka.us-mci.ionos.com",
		"us/ewr":   "https://kafka.us-ewr.ionos.com",
	}
)

// NewClient creates a new Kafka client
func NewClient(clientOptions clientoptions.TerraformClientOptions, fileConfig *fileconfiguration.FileConfig) *Client {
	config := shared.NewConfiguration(clientOptions.Credentials.Username, clientOptions.Credentials.Password, clientOptions.Credentials.Token, clientOptions.Endpoint)

	config.MaxRetries = constant.MaxRetries
	config.MaxWaitTime = constant.MaxWaitTime
	config.UserAgent = fmt.Sprintf(
		"terraform-provider/%s_ionos-cloud-sdk-go-kafka/%s_hashicorp-terraform/%s_terraform-plugin-sdk/%s_os/%s_arch/%s",
		clientOptions.Version, kafka.Version, clientOptions.TerraformVersion,
		meta.SDKVersionString(), runtime.GOOS, runtime.GOARCH, //nolint:staticcheck
	)

	client := Client{
		sdkClient:  *kafka.NewAPIClient(config),
		fileConfig: fileConfig,
	}
	client.sdkClient.GetConfig().HTTPClient = &http.Client{Transport: shared.CreateTransport(clientOptions.SkipTLSVerify, clientOptions.Certificate)}

	return &client
}

// ChangeConfigURL changes the url in the config based on the location
func (c *Client) ChangeConfigURL(location string) {
	config := c.sdkClient.GetConfig()
	config.Servers = shared.ServerConfigurations{
		{
			URL: locationToURL[location],
		},
	}
}
