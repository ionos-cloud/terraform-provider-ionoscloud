package vpn

import (
	"fmt"
	"net/http"
	"runtime"

	"github.com/hashicorp/terraform-plugin-sdk/v2/meta"
	vpn "github.com/ionos-cloud/sdk-go-bundle/products/vpn/v2"
	"github.com/ionos-cloud/sdk-go-bundle/shared"

	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils/constant"
)

// Client is a wrapper for the VPN SDK client
type Client struct {
	sdkClient vpn.APIClient
}

var (
	// AvailableLocations is a list of supported locations for VPN
	AvailableLocations = []string{"de/fra", "de/txl", "es/vit", "gb/bhx", "gb/lhr", "us/ewr", "us/las", "us/mci", "fr/par"}

	locationToURL = map[string]string{
		"de/fra": "https://vpn.de-fra.ionos.com",
		"de/txl": "https://vpn.de-txl.ionos.com",
		"es/vit": "https://vpn.es-vit.ionos.com",
		"gb/bhx": "https://vpn.gb-bhx.ionos.com",
		"gb/lhr": "https://vpn.gb-lhr.ionos.com",
		"us/ewr": "https://vpn.us-ewr.ionos.com",
		"us/las": "https://vpn.us-las.ionos.com",
		"us/mci": "https://vpn.us-mci.ionos.com",
		"fr/par": "https://vpn.fr-par.ionos.com",
	}
)

// NewClient returns a new ionoscloud logging client
func NewClient(username, password, token, url, terraformVersion string) *Client {
	newConfig := shared.NewConfiguration(username, password, token, url)
	newConfig.MaxRetries = constant.MaxRetries
	newConfig.MaxWaitTime = constant.MaxWaitTime
	newConfig.HTTPClient = &http.Client{Transport: utils.CreateTransport()}
	newConfig.UserAgent = fmt.Sprintf(
		"terraform-provider/%s_ionos-cloud-sdk-vpn/%s_hashicorp-terraform/%s_terraform-plugin-sdk/%s_os/%s_arch",
		vpn.Version, terraformVersion, meta.SDKVersionString(), runtime.GOOS, runtime.GOARCH, //nolint:staticcheck
	)

	return &Client{sdkClient: *vpn.NewAPIClient(newConfig)}
}

func (c *Client) changeConfigURL(location string) {
	config := c.sdkClient.GetConfig()
	config.Servers = shared.ServerConfigurations{
		{
			URL: locationToURL[location],
		},
	}
}
