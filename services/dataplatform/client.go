package dataplatform

import (
	"fmt"
	"net/http"
	"os"
	"runtime"

	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils/constant"

	"github.com/hashicorp/terraform-plugin-sdk/v2/meta"
	dataplatform "github.com/ionos-cloud/sdk-go-dataplatform"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils"
)

// Client - wrapper over sdk client, to allow for service layer
type Client struct {
	sdkClient dataplatform.APIClient
}

func NewClient(username, password, token, url, version, terraformVersion string) *Client {
	newConfigDataplatform := dataplatform.NewConfiguration(username, password, token, url)

	if os.Getenv(constant.IonosDebug) != "" {
		newConfigDataplatform.Debug = true
	}
	newConfigDataplatform.MaxRetries = constant.MaxRetries
	newConfigDataplatform.MaxWaitTime = constant.MaxWaitTime

	newConfigDataplatform.HTTPClient = &http.Client{Transport: utils.CreateTransport()}
	newConfigDataplatform.UserAgent = fmt.Sprintf(
		"terraform-provider/%s_ionos-cloud-sdk-go-dataplatform/%s_hashicorp-terraform/%s_terraform-plugin-sdk/%s_os/%s_arch/%s",
		version, dataplatform.Version, terraformVersion, meta.SDKVersionString(), runtime.GOOS, runtime.GOARCH)

	return &Client{sdkClient: *dataplatform.NewAPIClient(newConfigDataplatform)}
}
