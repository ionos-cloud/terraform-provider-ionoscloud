package dataplatform

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/meta"
	"github.com/ionos-cloud/sdk-go-bundle/shared"
	dataplatform "github.com/ionos-cloud/sdk-go-dataplatform"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils/bundle"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils/constant"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils/loadedconfig"
	"net/http"
	"os"
	"runtime"
)

// Client - wrapper over sdk client, to allow for service layer
type Client struct {
	sdkClient dataplatform.APIClient
}

func NewClient(clientOptions bundle.ClientOptions, loadedConfig *shared.LoadedConfig) *Client {
	loadedconfig.SetClientOptionsFromLoadedConfig(&clientOptions, loadedConfig, shared.Dataplatform)
	newConfigDataplatform := dataplatform.NewConfiguration(clientOptions.Credentials.Username, clientOptions.Credentials.Password, clientOptions.Credentials.Token, clientOptions.Endpoint)
	if os.Getenv(constant.IonosDebug) != "" {
		newConfigDataplatform.Debug = true
	}
	newConfigDataplatform.MaxRetries = constant.MaxRetries
	newConfigDataplatform.MaxWaitTime = constant.MaxWaitTime

	newConfigDataplatform.HTTPClient = &http.Client{Transport: utils.CreateTransport(clientOptions.SkipTLSVerify)}
	newConfigDataplatform.UserAgent = fmt.Sprintf(
		"terraform-provider/_ionos-cloud-sdk-go-dataplatform/%s_hashicorp-terraform/%s_terraform-plugin-sdk/%s_os/%s_arch/%s",
		dataplatform.Version, clientOptions.TerraformVersion, meta.SDKVersionString(), runtime.GOOS, runtime.GOARCH)

	return &Client{sdkClient: *dataplatform.NewAPIClient(newConfigDataplatform)}
}
