package containerregistry

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/meta"
	"github.com/ionos-cloud/sdk-go-bundle/shared"
	cr "github.com/ionos-cloud/sdk-go-container-registry"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils/bundle"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils/constant"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils/loadedconfig"
	"net/http"
	"os"
	"runtime"
)

type Client struct {
	sdkClient *cr.APIClient
}

func NewClient(clientOptions bundle.ClientOptions, loadedConfig *shared.LoadedConfig) *Client {
	loadedconfig.SetClientOptionsFromLoadedConfig(&clientOptions, loadedConfig, shared.ContainerRegistry)
	newConfig := cr.NewConfiguration(clientOptions.Credentials.Username, clientOptions.Credentials.Password, clientOptions.Credentials.Token, clientOptions.Endpoint)

	if os.Getenv(constant.IonosDebug) != "" {
		newConfig.Debug = true
	}
	newConfig.MaxRetries = constant.MaxRetries
	newConfig.MaxWaitTime = constant.MaxWaitTime

	newConfig.HTTPClient = &http.Client{Transport: utils.CreateTransport(clientOptions.SkipTLSVerify)}
	newConfig.UserAgent = fmt.Sprintf(
		"terraform-provider/_ionos-cloud-sdk-go-container-cr/%s_hashicorp-terraform/%s_terraform-plugin-sdk/%s_os/%s_arch/%s",
		cr.Version, clientOptions.TerraformVersion, meta.SDKVersionString(), runtime.GOOS, runtime.GOARCH)

	return &Client{
		sdkClient: cr.NewAPIClient(newConfig),
	}
}
