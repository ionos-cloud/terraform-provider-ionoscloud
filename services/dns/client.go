package dns

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/meta"
	"github.com/ionos-cloud/sdk-go-bundle/shared/fileconfiguration"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils/bundle"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils/constant"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils/loadedconfig"
	"net/http"
	"os"
	"runtime"

	dns "github.com/ionos-cloud/sdk-go-dns"

	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils"
)

type Client struct {
	sdkClient dns.APIClient
}

func NewClient(clientOptions bundle.ClientOptions, fileConfig *fileconfiguration.FileConfig) *Client {
	loadedconfig.SetGlobalClientOptionsFromFileConfig(&clientOptions, fileConfig, fileconfiguration.DNS)

	config := dns.NewConfiguration(clientOptions.Credentials.Username, clientOptions.Credentials.Password, clientOptions.Credentials.Token, clientOptions.Endpoint)

	if os.Getenv(constant.IonosDebug) != "" {
		config.Debug = true
	}
	config.MaxRetries = constant.MaxRetries
	config.MaxWaitTime = constant.MaxWaitTime
	config.HTTPClient = &http.Client{Transport: utils.CreateTransport(clientOptions.SkipTLSVerify)}
	fileconfiguration.AddCertsToClient(config.HTTPClient, clientOptions.Certificate)
	config.UserAgent = fmt.Sprintf(
		"terraform-provider/ionos-cloud-sdk-go-dns/%s_hashicorp-terraform/%s_terraform-plugin-sdk/%s_os/%s_arch/%s",
		dns.Version, clientOptions.TerraformVersion, meta.SDKVersionString(), runtime.GOOS, runtime.GOARCH) //nolint:staticcheck

	return &Client{sdkClient: *dns.NewAPIClient(config)}
}
