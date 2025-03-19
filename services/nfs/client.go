package nfs

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"runtime"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/meta"

	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/clientoptions"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils/constant"

	"github.com/ionos-cloud/sdk-go-bundle/products/nfs/v2"
	"github.com/ionos-cloud/sdk-go-bundle/shared"
	"github.com/ionos-cloud/sdk-go-bundle/shared/fileconfiguration"
)

// Client is a wrapper for the NFS SDK
type Client struct {
	sdkClient  nfs.APIClient
	fileConfig *fileconfiguration.FileConfig
}

// GetFileConfig returns the file configuration
func (c *Client) GetFileConfig() *fileconfiguration.FileConfig {
	return c.fileConfig
}

// GetConfig returns the configuration
func (c *Client) GetConfig() *shared.Configuration {
	return c.sdkClient.GetConfig()
}

func NewClient(clientOptions clientoptions.TerraformClientOptions, fileConfig *fileconfiguration.FileConfig) *Client {
	config := shared.NewConfiguration(clientOptions.Credentials.Username, clientOptions.Credentials.Password, clientOptions.Credentials.Token, clientOptions.Endpoint)

	config.MaxRetries = constant.MaxRetries
	config.MaxWaitTime = constant.MaxWaitTime
	config.HTTPClient = &http.Client{Transport: utils.CreateTransport(clientOptions.SkipTLSVerify)}
	config.UserAgent = fmt.Sprintf("terraform-provider/ionos-cloud-sdk-go-nfs/%s_hashicorp-terraform/%s_terraform-plugin-sdk/%s_os/%s_arch/%s",
		nfs.Version, clientOptions.TerraformVersion, meta.SDKVersionString(), runtime.GOOS, runtime.GOARCH) // nolint:staticcheck

	return &Client{sdkClient: *nfs.NewAPIClient(config),
		fileConfig: fileConfig}
}

// changeConfigURL sets the location of the NFS client which modifies the Host URL:
//   - de/fra:    https://nfs.de-fra.ionos.com
//   - de/txl:    https://nfs.de-txl.ionos.com
func (c *Client) changeConfigURL(location string) {
	config := c.sdkClient.GetConfig()
	// if there is no location set, return the client as is. allows to overwrite the url with IONOS_API_URL
	if location == "" && os.Getenv(ionosAPIURLNFS) != "" {
		config.Servers = shared.ServerConfigurations{
			{
				URL: utils.CleanURL(os.Getenv(ionosAPIURLNFS)),
			},
		}
		return
	}

	for _, server := range config.Servers {
		if strings.EqualFold(server.Description, shared.EndpointOverridden+location) || strings.EqualFold(server.URL, locationToURL[location]) {
			config.Servers = shared.ServerConfigurations{
				{
					URL:         server.URL,
					Description: shared.EndpointOverridden + location,
				},
			}
			return
		}
	}
}

// overrideClientEndpoint todo - after move to bundle, replace with generic function from fileConfig
func (c *Client) overrideClientEndpoint(productName, location string) {
	// whatever is set, at the end we need to check if the IONOS_API_URL_productname is set and use override the endpoint if yes
	defer c.changeConfigURL(location)
	if os.Getenv(shared.IonosApiUrlEnvVar) != "" {
		fmt.Printf("[DEBUG] Using custom endpoint %s\n", os.Getenv(shared.IonosApiUrlEnvVar))
		return
	}
	fileConfig := c.GetFileConfig()
	if fileConfig == nil {
		return
	}
	config := c.GetConfig()
	if config == nil {
		return
	}
	endpoint := fileConfig.GetProductLocationOverrides(productName, location)
	if endpoint == nil {
		log.Printf("[WARN] Missing endpoint for %s in location %s", productName, location)
		return
	}
	config.Servers = shared.ServerConfigurations{
		{
			URL:         endpoint.Name,
			Description: shared.EndpointOverridden + location,
		},
	}
	config.HTTPClient.Transport = shared.CreateTransport(endpoint.SkipTLSVerify, endpoint.CertificateAuthData)
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
