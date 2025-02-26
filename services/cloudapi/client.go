package cloudapi

import (
	"fmt"
	"github.com/ionos-cloud/sdk-go-bundle/shared"
	"net/http"
	"os"
	"runtime"

	"github.com/hashicorp/terraform-plugin-sdk/v2/meta"
	"github.com/ionos-cloud/sdk-go-bundle/shared/fileconfiguration"
	ionoscloud "github.com/ionos-cloud/sdk-go/v6"

	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils/bundle"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils/constant"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils/loadedconfig"
)

func NewClient(clientOptions bundle.ClientOptions, fileConfig *fileconfiguration.FileConfig) *ionoscloud.APIClient {
	loadedconfig.SetGlobalClientOptionsFromFileConfig(&clientOptions, fileConfig, fileconfiguration.Cloud)
	config := ionoscloud.NewConfiguration(clientOptions.Credentials.Username, clientOptions.Credentials.Password, clientOptions.Credentials.Token, clientOptions.Endpoint)
	config.UserAgent = fmt.Sprintf(
		"terraform-provider/_ionos-cloud-sdk-go/%s_hashicorp-terraform/%s_terraform-plugin-sdk/%s_os/%s_arch/%s",
		ionoscloud.Version, clientOptions.TerraformVersion, meta.SDKVersionString(), runtime.GOOS, runtime.GOARCH, //nolint:staticcheck
	)
	if os.Getenv(constant.IonosDebug) != "" {
		config.Debug = true
	}
	config.MaxRetries = constant.MaxRetries
	config.WaitTime = constant.MaxWaitTime
	config.HTTPClient = http.DefaultClient
	config.HTTPClient.Transport = shared.CreateTransport(clientOptions.SkipTLSVerify, clientOptions.Certificate)
	client := ionoscloud.NewAPIClient(config)
	return client
}
