package laas

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/meta"
	laas "github.com/ionos-cloud/sdk-go-laas"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils"
	"net/http"
	"os"
	"runtime"
)

type Client struct {
	sdkClient laas.APIClient
}

func NewClient(username, password, token, url, version, terraformVersion string) *Client {
	newConfigLaaS := laas.NewConfiguration(username, password, token, url)

	if os.Getenv(utils.IonosDebug) != "" {
		newConfigLaaS.Debug = true
	}
	newConfigLaaS.MaxRetries = utils.MaxRetries
	newConfigLaaS.MaxWaitTime = utils.MaxWaitTime
	newConfigLaaS.HTTPClient = &http.Client{Transport: utils.CreateTransport()}
	newConfigLaaS.UserAgent = fmt.Sprintf(
		"terraform-provider/%s_ionos-cloud-sdk-go-laas/%s_hashicorp-terraform/%s_terraform-plugin-sdk/%s_os/%s_arch/%s",
		version, laas.Version, terraformVersion, meta.SDKVersionString(), runtime.GOOS, runtime.GOARCH)

	return &Client{sdkClient: *laas.NewAPIClient(newConfigLaaS)}
}
