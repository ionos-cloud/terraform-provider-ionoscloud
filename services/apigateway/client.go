package apigateway

import (
	"fmt"
	"net/http"
	"os"
	"runtime"

	"github.com/hashicorp/terraform-plugin-sdk/v2/meta"
	apigateway "github.com/ionos-cloud/sdk-go-apigateway"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils/constant"
)

type Client struct {
	sdkClient apigateway.APIClient
}

func NewClient(username, password, token, url, version, terraformVersion string) *Client {
	config := apigateway.NewConfiguration(username, password, token, url)

	if os.Getenv(constant.IonosDebug) != "" {
		config.Debug = true
	}
	config.MaxRetries = constant.MaxRetries
	config.MaxWaitTime = constant.MaxWaitTime
	config.HTTPClient = &http.Client{Transport: utils.CreateTransport()}
	config.UserAgent = fmt.Sprintf(
		"terraform-provider/%s_ionos-cloud-sdk-go-apigateway/%s_hashicorp-terraform/%s_terraform-plugin-sdk/%s_os/%s_arch/%s",
		version, apigateway.Version, terraformVersion, meta.SDKVersionString(), runtime.GOOS, runtime.GOARCH,
	)

	return &Client{sdkClient: *apigateway.NewAPIClient(config)}
}
