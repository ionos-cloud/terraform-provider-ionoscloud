package dataplatform

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/meta"
	"github.com/ionos-cloud/sdk-go-bundle/shared/fileconfiguration"
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

func NewClient(clientOptions bundle.ClientOptions, fileConfig *fileconfiguration.FileConfig) *Client {
	loadedconfig.SetGlobalClientOptionsFromFileConfig(&clientOptions, fileConfig, fileconfiguration.Dataplatform)
	config := dataplatform.NewConfiguration(clientOptions.Credentials.Username, clientOptions.Credentials.Password, clientOptions.Credentials.Token, clientOptions.Endpoint)
	if os.Getenv(constant.IonosDebug) != "" {
		config.Debug = true
	}
	config.MaxRetries = constant.MaxRetries
	config.MaxWaitTime = constant.MaxWaitTime
	config.HTTPClient = &http.Client{Transport: utils.CreateTransport(clientOptions.SkipTLSVerify)}
	fileconfiguration.AddCertsToClient(config.HTTPClient, clientOptions.Certificate)

	config.UserAgent = fmt.Sprintf(
		"terraform-provider/_ionos-cloud-sdk-go-dataplatform/%s_hashicorp-terraform/%s_terraform-plugin-sdk/%s_os/%s_arch/%s",
		dataplatform.Version, clientOptions.TerraformVersion, meta.SDKVersionString(), runtime.GOOS, runtime.GOARCH)

	return &Client{sdkClient: *dataplatform.NewAPIClient(config)}
}
