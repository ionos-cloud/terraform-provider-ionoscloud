package autoscaling

import (
	"fmt"
	"net/http"
	"runtime"

	"github.com/hashicorp/terraform-plugin-sdk/v2/meta"
	autoscaling "github.com/ionos-cloud/sdk-go-bundle/products/vmautoscaling/v2"
	"github.com/ionos-cloud/sdk-go-bundle/shared"

	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils/constant"
)

type Client struct {
	sdkClient *autoscaling.APIClient
}

func NewClient(username, password, token, url, version, terraformVersion string, insecure bool) *Client {
	newAutoscalingConfig := shared.NewConfiguration(username, password, token, url)

	newAutoscalingConfig.MaxRetries = constant.MaxRetries
	newAutoscalingConfig.MaxWaitTime = constant.MaxWaitTime
	newAutoscalingConfig.UserAgent = fmt.Sprintf(
		"terraform-provider/%s_ionos-cloud-sdk-go-vm-autoscaling/%s_hashicorp-terraform/%s_terraform-plugin-sdk/%s_os/%s_arch/%s",
		version, autoscaling.Version, terraformVersion, meta.SDKVersionString(), runtime.GOOS, runtime.GOARCH) //nolint:staticcheck

	client := &Client{
		sdkClient: autoscaling.NewAPIClient(newAutoscalingConfig),
	}
	client.sdkClient.GetConfig().HTTPClient = &http.Client{Transport: utils.CreateTransport(insecure)}

	return client
}
