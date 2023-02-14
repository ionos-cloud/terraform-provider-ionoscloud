package dataplatform

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/meta"
	dataplatform "github.com/ionos-cloud/sdk-go-dataplatform"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils"
	"net/http"
	"os"
	"runtime"
	"time"
)

// Client - wrapper over sdk client, to allow for service layer
type Client struct {
	sdkClient dataplatform.APIClient
}

func NewClient(username, password, token, url, version, terraformVersion string) *Client {
	newConfigDataplatform := dataplatform.NewConfiguration(username, password, token, url)

	if os.Getenv("IONOS_DEBUG") != "" {
		newConfigDataplatform.Debug = true
	}
	newConfigDataplatform.MaxRetries = 999
	newConfigDataplatform.MaxWaitTime = 4 * time.Second

	newConfigDataplatform.HTTPClient = &http.Client{Transport: utils.CreateTransport()}
	newConfigDataplatform.UserAgent = fmt.Sprintf(
		"terraform-provider/%s_ionos-cloud-sdk-go-dataplatform/%s_hashicorp-terraform/%s_terraform-plugin-sdk/%s_os/%s_arch/%s",
		version, dataplatform.Version, terraformVersion, meta.SDKVersionString(), runtime.GOOS, runtime.GOARCH)

	return &Client{sdkClient: *dataplatform.NewAPIClient(newConfigDataplatform)}
}
