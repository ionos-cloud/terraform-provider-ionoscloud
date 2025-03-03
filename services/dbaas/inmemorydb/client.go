package inmemorydb

import (
	"fmt"
	"net/http"
	"os"
	"runtime"

	"github.com/hashicorp/terraform-plugin-sdk/v2/meta"
	inMemoryDB "github.com/ionos-cloud/sdk-go-dbaas-in-memory-db"

	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils/constant"
)

//nolint:golint
type InMemoryDBClient struct {
	sdkClient *inMemoryDB.APIClient
}

//nolint:golint
func NewInMemoryDBClient(username, password, token, url, version, terraformVersion string, insecure bool) *InMemoryDBClient {
	newConfigDbaas := inMemoryDB.NewConfiguration(username, password, token, url)

	if os.Getenv(constant.IonosDebug) != "" {
		newConfigDbaas.Debug = true
	}
	newConfigDbaas.MaxRetries = constant.MaxRetries
	newConfigDbaas.MaxWaitTime = constant.MaxWaitTime

	newConfigDbaas.HTTPClient = &http.Client{Transport: utils.CreateTransport(insecure)}
	newConfigDbaas.UserAgent = fmt.Sprintf(
		"terraform-provider/%s_ionos-cloud-sdk-go-dbaas-in-memory-db/%s_hashicorp-terraform/%s_terraform-plugin-sdk/%s_os/%s_arch/%s",
		version, inMemoryDB.Version, terraformVersion, meta.SDKVersionString(), runtime.GOOS, runtime.GOARCH) //nolint:staticcheck

	return &Client{
		sdkClient:  inMemoryDB.NewAPIClient(newConfigDbaas),
		fileConfig: fileConfig,
	}
}

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
