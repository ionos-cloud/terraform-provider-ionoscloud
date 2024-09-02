package logging

import (
	"fmt"
	"net/http"
	"runtime"

	"github.com/hashicorp/terraform-plugin-sdk/v2/meta"
	"github.com/ionos-cloud/sdk-go-bundle/products/logging/v2"
	"github.com/ionos-cloud/sdk-go-bundle/shared"

	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils/constant"
)

// Client is a wrapper for the Logging SDK client
type Client struct {
	sdkClient logging.APIClient
}

// NewClient returns a new Logging client
func NewClient(username, password, token, url, terraformVersion string) *Client {
	newConfigLogging := shared.NewConfiguration(username, password, token, url)

	newConfigLogging.MaxRetries = constant.MaxRetries
	newConfigLogging.MaxWaitTime = constant.MaxWaitTime
	newConfigLogging.HTTPClient = &http.Client{Transport: utils.CreateTransport()}
	newConfigLogging.UserAgent = fmt.Sprintf(
		"terraform-provider/%s_ionos-cloud-sdk-go-logging/%s_hashicorp-terraform/%s_terraform-plugin-sdk/%s_os/%s_arch",
		logging.Version, terraformVersion, meta.SDKVersionString(), runtime.GOOS, runtime.GOARCH) //nolint:staticcheck

	return &Client{sdkClient: *logging.NewAPIClient(newConfigLogging)}
}

func (c *Client) changeConfigURL(location string) {
	if location == "" {
		location = "de/txl"
	}
	config := c.sdkClient.GetConfig()
	config.Servers = shared.ServerConfigurations{
		{
			URL: locationToURL[location],
		},
	}
}

var (
	// AvailableLocations is a list of available locations
	AvailableLocations = []string{"de/fra", "de/txl", "es/vit", "gb/lhr", "fr/par"}

	locationToURL = map[string]string{
		"de/fra": "https://logging.de-fra.ionos.com",
		"de/txl": "https://logging.de-txl.ionos.com",
		"es/vit": "https://logging.es-vit.ionos.com",
		"gb/lhr": "https://logging.gb-lhr.ionos.com",
		"fr/par": "https://logging.fr-par.ionos.com",
	}
)
