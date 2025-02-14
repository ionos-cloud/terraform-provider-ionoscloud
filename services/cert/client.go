package cert

import (
	"fmt"
	"github.com/ionos-cloud/sdk-go-bundle/shared"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils/bundle"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils/loadedconfig"
	"net/http"
	"os"
	"runtime"

	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils/constant"

	"github.com/hashicorp/terraform-plugin-sdk/v2/meta"
	certmanager "github.com/ionos-cloud/sdk-go-cert-manager"

	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils"
)

type Client struct {
	sdkClient    *certmanager.APIClient
	loadedConfig *shared.LoadedConfig
}

func (c *Client) GetLoadedConfig() *shared.LoadedConfig {
	return c.loadedConfig
}

func (c *Client) GetConfig() *certmanager.Configuration {
	return c.sdkClient.GetConfig()
}

// todo cguran cert has both location(auto-cert) and no location on certificate. How do we override?
func NewClient(clientOptions bundle.ClientOptions, loadedConfig *shared.LoadedConfig) *Client {
	loadedconfig.SetClientOptionsFromFileConfig(&clientOptions, loadedConfig, shared.Cert)
	certConfig := certmanager.NewConfiguration(clientOptions.Credentials.Username, clientOptions.Credentials.Password, clientOptions.Credentials.Token, clientOptions.Endpoint)

	if os.Getenv(constant.IonosDebug) != "" {
		certConfig.Debug = true
	}
	certConfig.MaxRetries = constant.MaxRetries
	certConfig.MaxWaitTime = constant.MaxWaitTime

	certConfig.HTTPClient = &http.Client{Transport: utils.CreateTransport(clientOptions.SkipTLSVerify)}
	certConfig.UserAgent = fmt.Sprintf(
		"terraform-provider/_ionos-cloud-sdk-go-cert-manager/%s_hashicorp-terraform/%s_terraform-plugin-sdk/%s_os/%s_arch/%s",
		certmanager.Version, clientOptions.TerraformVersion, meta.SDKVersionString(), runtime.GOOS, runtime.GOARCH)

	return &Client{
		sdkClient:    certmanager.NewAPIClient(certConfig),
		loadedConfig: loadedConfig,
	}
}
