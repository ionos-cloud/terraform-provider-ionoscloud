package autoscaling

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/meta"
	"github.com/ionos-cloud/sdk-go-bundle/shared/fileconfiguration"
	autoscaling "github.com/ionos-cloud/sdk-go-vm-autoscaling"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils/bundle"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils/constant"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils/loadedconfig"
	"net/http"
	"os"
	"runtime"
)

type Client struct {
	sdkClient *autoscaling.APIClient
}

func NewClient(clientOptions bundle.ClientOptions, fileConfig *fileconfiguration.FileConfig) *Client {
	loadedconfig.SetGlobalClientOptionsFromFileConfig(&clientOptions, fileConfig, fileconfiguration.Autoscaling)
	config := autoscaling.NewConfiguration(clientOptions.Credentials.Username, clientOptions.Credentials.Password, clientOptions.Credentials.Token, clientOptions.Endpoint)
	config.UserAgent = fmt.Sprintf(
		"terraform-provider/_ionos-cloud-sdk-go-vm-autoscaling/%s_hashicorp-terraform/%s_terraform-plugin-sdk/%s_os/%s_arch/%s",
		autoscaling.Version, clientOptions.TerraformVersion, meta.SDKVersionString(), runtime.GOOS, runtime.GOARCH) //nolint:staticcheck

	if os.Getenv(constant.IonosDebug) != "" {
		config.Debug = true
	}
	config.MaxRetries = constant.MaxRetries
	config.WaitTime = constant.MaxWaitTime
	config.HTTPClient = &http.Client{Transport: utils.CreateTransport(clientOptions.SkipTLSVerify)}
	fileconfiguration.AddCertsToClient(config.HTTPClient, clientOptions.Certificate)
	client := &Client{sdkClient: autoscaling.NewAPIClient(config)}
	return client
}
