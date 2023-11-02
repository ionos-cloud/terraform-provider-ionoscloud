package autoscaling

import (
	"fmt"
	"net/http"
	"os"
	"runtime"

	"github.com/hashicorp/terraform-plugin-sdk/v2/meta"
	logging "github.com/ionos-cloud/sdk-go-logging"
	autoscaling "github.com/ionos-cloud/sdk-go-vm-autoscaling"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils/constant"
)

type Client struct {
	sdkClient *autoscaling.APIClient
}

func NewClient(username, password, token, url, version, terraformVersion string) *Client {
	newConfigLogging := autoscaling.NewConfiguration(username, password, token, url)

	if os.Getenv(constant.IonosDebug) != "" {
		newConfigLogging.Debug = true
	}
	newConfigLogging.MaxRetries = constant.MaxRetries
	newConfigLogging.MaxWaitTime = constant.MaxWaitTime
	newConfigLogging.HTTPClient = &http.Client{Transport: utils.CreateTransport()}
	newConfigLogging.UserAgent = fmt.Sprintf(
		"terraform-provider/%s_ionos-cloud-sdk-go-logging/%s_hashicorp-terraform/%s_terraform-plugin-sdk/%s_os/%s_arch/%s",
		version, logging.Version, terraformVersion, meta.SDKVersionString(), runtime.GOOS, runtime.GOARCH)

	return &Client{sdkClient: autoscaling.NewAPIClient(newConfigLogging)}
}
