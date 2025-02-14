package cert

import (
	"fmt"
	"github.com/ionos-cloud/sdk-go-bundle/shared/fileconfiguration"
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
	sdkClient  *certmanager.APIClient
	fileConfig *fileconfiguration.FileConfig
}

func (c *Client) GetFileConfig() *fileconfiguration.FileConfig {
	return c.fileConfig
}

func (c *Client) GetConfig() *certmanager.Configuration {
	return c.sdkClient.GetConfig()
}

// todo cguran cert has both location(auto-cert) and no location on certificate. How do we override?
func NewClient(clientOptions bundle.ClientOptions, fileConfig *fileconfiguration.FileConfig) *Client {
	loadedconfig.SetGlobalClientOptionsFromFileConfig(&clientOptions, fileConfig, fileconfiguration.Cert)
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
		sdkClient:  certmanager.NewAPIClient(certConfig),
		fileConfig: fileConfig,
	}
}
