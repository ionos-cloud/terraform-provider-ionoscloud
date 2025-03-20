package cert

import (
	"fmt"
	"net/http"
	"runtime"

	"github.com/ionos-cloud/sdk-go-bundle/shared"

	"github.com/hashicorp/terraform-plugin-sdk/v2/meta"
	certmanager "github.com/ionos-cloud/sdk-go-bundle/products/cert/v2"
	"github.com/ionos-cloud/sdk-go-bundle/shared/fileconfiguration"

	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/clientoptions"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils/constant"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils/loadedconfig"
)

type Client struct {
	sdkClient  certmanager.APIClient
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

// NewClient todo cert has both location(auto-cert) and no location on certificate. How do we override?
func NewClient(clientOptions clientoptions.TerraformClientOptions, fileConfig *fileconfiguration.FileConfig) *Client {
	loadedconfig.SetGlobalClientOptionsFromFileConfig(&clientOptions, fileConfig, fileconfiguration.Cert)
	config := shared.NewConfiguration(clientOptions.Credentials.Username, clientOptions.Credentials.Password, clientOptions.Credentials.Token, clientOptions.Endpoint)
	config.MaxRetries = constant.MaxRetries
	config.MaxWaitTime = constant.MaxWaitTime
	config.UserAgent = fmt.Sprintf(
		"terraform-provider/%s_ionos-cloud-sdk-go-cert-manager/%s_hashicorp-terraform/%s_terraform-plugin-sdk/%s_os/%s_arch/%s",
		clientOptions.Version, certmanager.Version, clientOptions.TerraformVersion,
		meta.SDKVersionString(), runtime.GOOS, runtime.GOARCH,
	)
	client := &Client{
		sdkClient: *certmanager.NewAPIClient(config),
	}
	client.sdkClient.GetConfig().HTTPClient = &http.Client{Transport: shared.CreateTransport(clientOptions.SkipTLSVerify, clientOptions.Certificate)}

	return client
}
