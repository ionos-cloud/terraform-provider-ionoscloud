package cdn

import (
	"fmt"
	"net/http"
	"os"
	"runtime"

	"github.com/hashicorp/terraform-plugin-sdk/v2/meta"
	cdn "github.com/ionos-cloud/sdk-go-cdn"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils/constant"
)

type Client struct {
	sdkClient *cdn.APIClient
}

func NewClient(username, password, token, url, version, terraformVersion string) *Client {
	newCdnConfig := cdn.NewConfiguration(username, password, token, url)

	if os.Getenv(constant.IonosDebug) != "" {
		newCdnConfig.Debug = true
	}
	newCdnConfig.MaxRetries = constant.MaxRetries
	newCdnConfig.MaxWaitTime = constant.MaxWaitTime
	newCdnConfig.HTTPClient = &http.Client{Transport: utils.CreateTransport()}
	newCdnConfig.UserAgent = fmt.Sprintf(
		"terraform-provider/%s_ionos-cloud-sdk-go-logging/%s_hashicorp-terraform/%s_terraform-plugin-sdk/%s_os/%s_arch/%s",
		version, cdn.Version, terraformVersion, meta.SDKVersionString(), runtime.GOOS, runtime.GOARCH)

	return &Client{sdkClient: cdn.NewAPIClient(newCdnConfig)}
}
