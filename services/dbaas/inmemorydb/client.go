package inmemorydb

import (
	"fmt"
	"github.com/ionos-cloud/sdk-go-bundle/shared"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils/bundle"
	"net/http"
	"os"
	"runtime"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/meta"
	inMemoryDB "github.com/ionos-cloud/sdk-go-dbaas-in-memory-db"

	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils/constant"
)

//nolint:golint
type Client struct {
	sdkClient    *inMemoryDB.APIClient
	loadedConfig *shared.LoadedConfig
}

func (c *Client) GetLoadedConfig() *shared.LoadedConfig {
	return c.loadedConfig
}

func (c *Client) GetConfig() *inMemoryDB.Configuration {
	return c.sdkClient.GetConfig()
}

// NewClient creates a new in-memory db client. LoadedConfig is used to set/override the client options if it exists
func NewClient(clientOptions bundle.ClientOptions, sharedLoadedConfig *shared.LoadedConfig) *Client {
	newConfigDbaas := inMemoryDB.NewConfiguration(clientOptions.Credentials.Username, clientOptions.Credentials.Password,
		clientOptions.Credentials.Token, clientOptions.Endpoint)

	if os.Getenv(constant.IonosDebug) != "" {
		newConfigDbaas.Debug = true
	}
	newConfigDbaas.MaxRetries = constant.MaxRetries
	newConfigDbaas.MaxWaitTime = constant.MaxWaitTime

	newConfigDbaas.HTTPClient = &http.Client{Transport: utils.CreateTransport(clientOptions.SkipTLSVerify)}
	newConfigDbaas.UserAgent = fmt.Sprintf(
		"terraform-provider/ionos-cloud-sdk-go-dbaas-in-memory-db/%s_hashicorp-terraform/%s_terraform-plugin-sdk/%s_os/%s_arch/%s",
		inMemoryDB.Version, clientOptions.TerraformVersion, meta.SDKVersionString(), runtime.GOOS, runtime.GOARCH) //nolint:staticcheck

	return &Client{
		sdkClient:    inMemoryDB.NewAPIClient(newConfigDbaas),
		loadedConfig: sharedLoadedConfig,
	}
}

// todo cguran override on each request
var (
	locationToURL = map[string]string{
		"":       "https://in-memory-db.de-fra.ionos.com",
		"de/fra": "https://in-memory-db.de-fra.ionos.com",
		"de/txl": "https://in-memory-db.de-txl.ionos.com",
		"es/vit": "https://in-memory-db.es-vit.ionos.com",
		"gb/lhr": "https://in-memory-db.gb-lhr.ionos.com",
		"us/ewr": "https://in-memory-db.us-ewr.ionos.com",
		"us/las": "https://in-memory-db.us-las.ionos.com",
		"us/mci": "https://in-memory-db.us-mci.ionos.com",
		"fr/par": "https://in-memory-db.fr-par.ionos.com",
	}
	ionosAPIURLInMemoryDB = "IONOS_API_URL_INMEMORYDB"
)

// changeConfigURL modifies the URL inside the client configuration.
// This function is required in order to make requests to different endpoints based on location.
func (c *Client) changeConfigURL(location string) {
	clientConfig := c.sdkClient.GetConfig()
	if location == "" && os.Getenv(ionosAPIURLInMemoryDB) != "" {
		clientConfig.Servers = inMemoryDB.ServerConfigurations{
			{
				URL: utils.CleanURL(os.Getenv(ionosAPIURLInMemoryDB)),
			},
		}
		return
	}
	for _, server := range clientConfig.Servers {
		if strings.EqualFold(server.Description, shared.EndpointOverridden+location) || strings.EqualFold(server.URL, locationToURL[location]) {
			clientConfig.Servers = inMemoryDB.ServerConfigurations{
				{
					URL:         server.URL,
					Description: shared.EndpointOverridden + location,
				},
			}
			return
		}
	}
}
