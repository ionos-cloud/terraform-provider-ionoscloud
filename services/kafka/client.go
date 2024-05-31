package kafka

import (
	"fmt"
	"net/http"
	"os"
	"runtime"

	kafka "github.com/ionos-cloud/sdk-go-kafka"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils/constant"

	"github.com/hashicorp/terraform-plugin-sdk/v2/meta"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils"
)

type Client struct {
	sdkClient kafka.APIClient
}

func NewClient(username, password, token, url, version, terraformVersion string) *Client {
	config := kafka.NewConfiguration(username, password, token, url)

	if os.Getenv(constant.IonosDebug) != "" {
		config.Debug = true
	}
	config.MaxRetries = constant.MaxRetries
	config.MaxWaitTime = constant.MaxWaitTime
	config.HTTPClient = &http.Client{Transport: utils.CreateTransport()}
	config.UserAgent = fmt.Sprintf(
		"terraform-provider/%s_ionos-cloud-sdk-go-kafka/%s_hashicorp-terraform/%s_terraform-plugin-sdk/%s_os/%s_arch/%s",
		version, kafka.Version, terraformVersion, meta.SDKVersionString(), runtime.GOOS, runtime.GOARCH)

	return &Client{sdkClient: *kafka.NewAPIClient(config)}
}
