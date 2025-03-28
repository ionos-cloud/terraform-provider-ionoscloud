package mariadb

import (
	"fmt"
	"net/http"
	"os"
	"runtime"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/meta"
	mariadb "github.com/ionos-cloud/sdk-go-bundle/products/dbaas/mariadb/v2"
	shared "github.com/ionos-cloud/sdk-go-bundle/shared"
	"github.com/ionos-cloud/sdk-go-bundle/shared/fileconfiguration"

	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/clientoptions"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils/constant"
)

type Client struct {
	sdkClient  *mariadb.APIClient
	fileConfig *fileconfiguration.FileConfig
}

// GetFileConfig returns the loaded configuration of the client
func (c *Client) GetFileConfig() *fileconfiguration.FileConfig {
	return c.fileConfig
}

// GetConfig returns the configuration of the client
func (c *Client) GetConfig() *shared.Configuration {
	return c.sdkClient.GetConfig()
}

func NewClient(clientOptions clientoptions.TerraformClientOptions, fileConfig *fileconfiguration.FileConfig) *Client {
	newConfig := shared.NewConfiguration(clientOptions.Credentials.Username, clientOptions.Credentials.Password, clientOptions.Credentials.Token, clientOptions.Endpoint)

	newConfig.MaxRetries = constant.MaxRetries
	newConfig.MaxWaitTime = constant.MaxWaitTime

	newConfig.UserAgent = fmt.Sprintf(
		"terraform-provider/%s_ionos-cloud-sdk-go-dbaas-mariadb/%s_hashicorp-terraform/%s_terraform-plugin-sdk/%s_os/%s_arch/%s",
		clientOptions.Version, mariadb.Version, clientOptions.TerraformVersion,
		meta.SDKVersionString(), runtime.GOOS, runtime.GOARCH, //nolint:staticcheck
	)

	client := &Client{
		sdkClient:  mariadb.NewAPIClient(newConfig),
		fileConfig: fileConfig,
	}

	client.sdkClient.GetConfig().HTTPClient = &http.Client{Transport: utils.CreateTransport(clientOptions.SkipTLSVerify)}
	return client
}

// ChangeConfigURL modifies the URL inside the client configuration.
// This function is required in order to make requests to different endpoints based on location.
func (c *Client) ChangeConfigURL(location string) {
	clientConfig := c.sdkClient.GetConfig()
	if location == "" && os.Getenv(ionosAPIURLMariaDB) != "" {
		clientConfig.Servers = shared.ServerConfigurations{
			{
				URL: utils.CleanURL(os.Getenv(ionosAPIURLMariaDB)),
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

var (
	locationToURL = map[string]string{
		"":       "https://mariadb.de-txl.ionos.com",
		"de/fra": "https://mariadb.de-fra.ionos.com",
		"de/txl": "https://mariadb.de-txl.ionos.com",
		"es/vit": "https://mariadb.es-vit.ionos.com",
		"fr/par": "https://mariadb.fr-par.ionos.com",
		"gb/lhr": "https://mariadb.gb-lhr.ionos.com",
		"us/ewr": "https://mariadb.us-ewr.ionos.com",
		"us/las": "https://mariadb.us-las.ionos.com",
		"us/mci": "https://mariadb.us-mci.ionos.com",
	}
	ionosAPIURLMariaDB = "IONOS_API_URL_MARIADB"
)
