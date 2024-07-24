package vpn

import (
	"fmt"
	"net/http"
	"runtime"

	"github.com/hashicorp/terraform-plugin-sdk/v2/meta"
	vpn "github.com/ionos-cloud/sdk-go-bundle/products/vpn/v2"
	"github.com/ionos-cloud/sdk-go-bundle/shared"

	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils/constant"
)

// Client is a wrapper for the VPN SDK client
type Client struct {
	sdkClient vpn.APIClient
}

// NewClient returns a new ionoscloud logging client
func NewClient(username, password, token, url, terraformVersion string) *Client {
	newConfigLogging := shared.NewConfiguration(username, password, token, url)
	newConfigLogging.MaxRetries = constant.MaxRetries
	newConfigLogging.MaxWaitTime = constant.MaxWaitTime
	newConfigLogging.HTTPClient = &http.Client{Transport: utils.CreateTransport()}
	newConfigLogging.UserAgent = fmt.Sprintf(
		"terraform-provider/%s_ionos-cloud-sdk-vpn/%s_hashicorp-terraform/%s_terraform-plugin-sdk/%s_os/%s_arch",
		vpn.Version, terraformVersion, meta.SDKVersionString(), runtime.GOOS, runtime.GOARCH) //nolint:staticcheck

	return &Client{sdkClient: *vpn.NewAPIClient(newConfigLogging)}
}
