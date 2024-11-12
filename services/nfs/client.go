package nfs

import (
	"fmt"
	"net/http"
	"os"
	"runtime"

	"github.com/hashicorp/terraform-plugin-sdk/v2/meta"

	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils/constant"

	sdk "github.com/ionos-cloud/sdk-go-nfs"
)

// Client is a wrapper for the NFS SDK
type Client struct {
	sdkClient sdk.APIClient
}

// NewClient returns a new NFS client
func NewClient(username, password, token, url, version, terraformVersion string) *Client {
	config := sdk.NewConfiguration(username, password, token, url)

	if os.Getenv(constant.IonosDebug) != "" {
		config.Debug = true
	}
	config.MaxRetries = constant.MaxRetries
	config.MaxWaitTime = constant.MaxWaitTime
	config.HTTPClient = &http.Client{Transport: utils.CreateTransport()}
	config.UserAgent = fmt.Sprintf("terraform-provider/%s_ionos-cloud-sdk-go-nfs/%s_hashicorp-terraform/%s_terraform-plugin-sdk/%s_os/%s_arch/%s", version, sdk.Version, terraformVersion, meta.SDKVersionString(), runtime.GOOS, runtime.GOARCH) // nolint:staticcheck

	return &Client{sdkClient: *sdk.NewAPIClient(config)}
}

// Location sets the location of the NFS client which modifies the Host URL:
//   - de/fra:    https://nfs.de-fra.ionos.com
//   - de/txl:    https://nfs.de-txl.ionos.com
func (c *Client) Location(location string) *Client {
	// if there is no location set, return the client as is. allows to overwrite the url with IONOS_API_URL
	if location == "" && os.Getenv(ionosAPIURLNFS) != "" {
		c.sdkClient.GetConfig().Servers = sdk.ServerConfigurations{
			{
				URL: utils.CleanURL(os.Getenv(ionosAPIURLNFS)),
			},
		}
		return c
	}
	var locationToURL = map[string]string{
		"de/fra": "https://nfs.de-fra.ionos.com",
		"de/txl": "https://nfs.de-txl.ionos.com",
	}
	c.sdkClient.GetConfig().Servers = sdk.ServerConfigurations{
		{
			URL: locationToURL[location],
		},
	}

	return c
}

// ValidNFSLocations is a list of valid locations for the Network File Storage Cluster.
var ValidNFSLocations = []string{"de/fra", "de/txl"}

// ionosAPIURLNFS is the environment variable key for the NFS API URL.
var ionosAPIURLNFS = "IONOS_API_URL_NFS"
