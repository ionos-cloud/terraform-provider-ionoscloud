package containerregistry

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/meta"
	"github.com/ionos-cloud/sdk-go-bundle/common"
	cr "github.com/ionos-cloud/sdk-go-bundle/products/containerregistry"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils"
	"net/http"
	"runtime"
	"time"
)

type Client struct {
	sdkClient *cr.APIClient
}

func NewClient(username, password, token, url, version, terraformVersion string) *Client {
	newConfigRegistry := common.NewConfiguration(username, password, token, url)

	newConfigRegistry.MaxRetries = 999
	newConfigRegistry.MaxWaitTime = 4 * time.Second

	newConfigRegistry.HTTPClient = &http.Client{Transport: utils.CreateTransport()}
	newConfigRegistry.UserAgent = fmt.Sprintf(
		"terraform-provider/%s_ionos-cloud-sdk-go-container-cr/%s_hashicorp-terraform/%s_terraform-plugin-sdk/%s_os/%s_arch/%s",
		version, cr.Version, terraformVersion, meta.SDKVersionString(), runtime.GOOS, runtime.GOARCH)

	return &Client{
		sdkClient: cr.NewAPIClient(newConfigRegistry),
	}
}
