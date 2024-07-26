package kafka

import (
	"fmt"
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
	sdkClient kafka.APIClient
}

var locationToURL = map[string]string{
	"de/fra": "https://kafka.de-fra.ionos.com",
	"de/txl": "https://kafka.de-txl.ionos.com",
	"es/vit": "https://kafka.es-vit.ionos.com",
	"gb/lhr": "https://kafka.gb-lhr.ionos.com",
	"us/ewr": "https://kafka.us-ewr.ionos.com",
	"us/las": "https://kafka.us-las.ionos.com",
	"us/mci": "https://kafka.us-mci.ionos.com",
	"fr/par": "https://kafka.fr-par.ionos.com",
	"pre":    "https://pre.kafka.de-fra.ionos.com",
}

// NewClient returns a new Kafka API client
func NewClient(username, password, token, url, version, terraformVersion string) *Client {
	config := kafka.NewConfiguration(username, password, token, url)

	if os.Getenv(constant.IonosDebug) != "" {
		config.Debug = true
	}
	config.MaxRetries = constant.MaxRetries
	config.MaxWaitTime = constant.MaxWaitTime
	config.HTTPClient = &http.Client{Transport: utils.CreateTransport()}
	config.UserAgent = fmt.Sprintf(
		"terraform-provider/%s_ionos-cloud-sdk-go-kafka/%s_hashicorp-terraform/%s_terraform-plugin-sdk/%s_os/%s_arch/%s",
		version, kafka.Version, terraformVersion, meta.SDKVersionString(), runtime.GOOS, runtime.GOARCH, //nolint:staticcheck
	)

	return &Client{sdkClient: *kafka.NewAPIClient(config)}
}

func (c *Client) changeConfigURL(location string) {
	config := c.sdkClient.GetConfig()
	config.Servers = kafka.ServerConfigurations{
		{
			URL: locationToURL[location],
		},
	}
}
