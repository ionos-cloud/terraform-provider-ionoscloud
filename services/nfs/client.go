package nfs

import (
	"fmt"
	"github.com/ionos-cloud/sdk-go-bundle/shared"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils/bundle"
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
	sdkClient    sdk.APIClient
	loadedConfig *shared.LoadedConfig
}

func NewClient(clientOptions bundle.ClientOptions, sharedLoadedConfig *shared.LoadedConfig) *Client {
	config := sdk.NewConfiguration(clientOptions.Credentials.Username, clientOptions.Credentials.Password, clientOptions.Credentials.Token, clientOptions.Endpoint)

	if os.Getenv(constant.IonosDebug) != "" {
		config.Debug = true
	}
	config.MaxRetries = constant.MaxRetries
	config.MaxWaitTime = constant.MaxWaitTime
	config.HTTPClient = &http.Client{Transport: utils.CreateTransport(clientOptions.SkipTLSVerify)}
	config.UserAgent = fmt.Sprintf("terraform-provider/ionos-cloud-sdk-go-nfs/%s_hashicorp-terraform/%s_terraform-plugin-sdk/%s_os/%s_arch/%s",
		sdk.Version, clientOptions.TerraformVersion, meta.SDKVersionString(), runtime.GOOS, runtime.GOARCH) // nolint:staticcheck

	return &Client{sdkClient: *sdk.NewAPIClient(config),
		loadedConfig: sharedLoadedConfig}
}

// changeConfigURL sets the location of the NFS client which modifies the Host URL:
//   - de/fra:    https://nfs.de-fra.ionos.com
//   - de/txl:    https://nfs.de-txl.ionos.com
func (c *Client) changeConfigURL(location string) *Client {
	// if there is no location set, return the client as is. allows to overwrite the url with IONOS_API_URL
	if location == "" && os.Getenv(ionosAPIURLNFS) != "" {
		c.sdkClient.GetConfig().Servers = sdk.ServerConfigurations{
			{
				URL: utils.CleanURL(os.Getenv(ionosAPIURLNFS)),
			},
		}
		return c
	}

	c.sdkClient.GetConfig().Servers = sdk.ServerConfigurations{
		{
			URL: locationToURL[location],
		},
	}

	return c
}

var (
	locationToURL = map[string]string{
		"":       "https://nfs.de-fra.ionos.com",
		"de/fra": "https://nfs.de-fra.ionos.com",
		"de/txl": "https://nfs.de-txl.ionos.com",
		"fr/par": "https://nfs.fr-par.ionos.com",
		"gb/lhr": "https://nfs.gb-lhr.ionos.com",
		"es/vit": "https://nfs.es-vit.ionos.com",
		"us/las": "https://nfs.us-las.ionos.com",
		"us/ewr": "https://nfs.us-ewr.ionos.com",
		"us/mci": "https://nfs.us-mci.ionos.com",
	}
	// ValidNFSLocations is a list of valid locations for the Network File Storage Cluster.
	ValidNFSLocations = []string{"de/fra", "de/txl", "fr-par", "gb-lhr", "es/vit", "us/las", "us/ewr", "us/mci"}
)

// ionosAPIURLNFS is the environment variable key for the NFS API URL.
const ionosAPIURLNFS = "IONOS_API_URL_NFS"
