package inmemorydb

import (
	"fmt"
	"net/http"
	"os"
	"runtime"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/meta"
	inMemoryDB "github.com/ionos-cloud/sdk-go-bundle/products/dbaas/inmemorydb/v2"
	"github.com/ionos-cloud/sdk-go-bundle/shared"
	"github.com/ionos-cloud/sdk-go-bundle/shared/fileconfiguration"

	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/clientoptions"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils/constant"
)

// Client is a wrapper for the in-memory db client that has file configuration
type Client struct {
	sdkClient  *inMemoryDB.APIClient
	fileConfig *fileconfiguration.FileConfig
}

// GetFileConfig - returns the configuration read from the config file
func (c *Client) GetFileConfig() *fileconfiguration.FileConfig {
	return c.fileConfig
}

// GetConfig - returns the configuration
func (c *Client) GetConfig() *shared.Configuration {
	return c.sdkClient.GetConfig()
}

// NewClient creates a new in-memory db client. fileConfig is used to set/override the client options if it exists
func NewClient(clientOptions clientoptions.TerraformClientOptions, fileConfig *fileconfiguration.FileConfig) *Client {
	newConfigDbaas := shared.NewConfiguration(clientOptions.Credentials.Username, clientOptions.Credentials.Password,
		clientOptions.Credentials.Token, clientOptions.Endpoint)

	newConfigDbaas.MaxRetries = constant.MaxRetries
	newConfigDbaas.MaxWaitTime = constant.MaxWaitTime

	newConfigDbaas.HTTPClient = &http.Client{Transport: shared.CreateTransport(clientOptions.SkipTLSVerify, clientOptions.Certificate)}
	newConfigDbaas.UserAgent = fmt.Sprintf(
		"terraform-provider/%s_ionos-cloud-sdk-go-dbaas-in-memory-db/%s_hashicorp-terraform/%s_terraform-plugin-sdk/%s_os/%s_arch/%s",
		clientOptions.Version, inMemoryDB.Version, clientOptions.TerraformVersion,
		meta.SDKVersionString(), runtime.GOOS, runtime.GOARCH, //nolint:staticcheck
	)

	return &Client{
		sdkClient:  inMemoryDB.NewAPIClient(newConfigDbaas),
		fileConfig: fileConfig,
	}
}

var (
	locationToURL = map[string]string{
		"":         "https://in-memory-db.de-fra.ionos.com",
		"de/fra":   "https://in-memory-db.de-fra.ionos.com",
		"de/fra/2": "https://in-memory-db.de-fra.ionos.com",
		"de/txl":   "https://in-memory-db.de-txl.ionos.com",
		"es/vit":   "https://in-memory-db.es-vit.ionos.com",
		"gb/lhr":   "https://in-memory-db.gb-lhr.ionos.com",
		"us/ewr":   "https://in-memory-db.us-ewr.ionos.com",
		"us/las":   "https://in-memory-db.us-las.ionos.com",
		"us/mci":   "https://in-memory-db.us-mci.ionos.com",
		"fr/par":   "https://in-memory-db.fr-par.ionos.com",
	}
	ionosAPIURLInMemoryDB = "IONOS_API_URL_INMEMORYDB"
)

// changeConfigURL modifies the URL inside the client configuration.
// This function is required in order to make requests to different endpoints based on location.
func (c *Client) changeConfigURL(location string) {
	clientConfig := c.sdkClient.GetConfig()
	if location == "" && os.Getenv(ionosAPIURLInMemoryDB) != "" {
		clientConfig.Servers = shared.ServerConfigurations{
			{
				URL: utils.CleanURL(os.Getenv(ionosAPIURLInMemoryDB)),
			},
		}
		return
	}
	for _, server := range clientConfig.Servers {
		if strings.EqualFold(server.Description, shared.EndpointOverridden+location) || strings.EqualFold(server.URL, locationToURL[location]) {
			clientConfig.Servers = shared.ServerConfigurations{
				{
					URL:         server.URL,
					Description: shared.EndpointOverridden + location,
				},
			}
			return
		}
	}
}
