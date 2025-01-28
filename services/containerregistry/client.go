package containerregistry

import (
	"fmt"
	"net/http"
	"runtime"

	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils/constant"

	"github.com/hashicorp/terraform-plugin-sdk/v2/meta"
	cr "github.com/ionos-cloud/sdk-go-bundle/products/containerregistry/v2"
	"github.com/ionos-cloud/sdk-go-bundle/shared"

	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils"
)

type Client struct {
	sdkClient cr.APIClient
}

func NewClient(username, password, token, url, version, terraformVersion string, insecure bool) *Client {
	newConfigRegistry := shared.NewConfiguration(username, password, token, url)

	newConfigRegistry.MaxRetries = constant.MaxRetries
	newConfigRegistry.MaxWaitTime = constant.MaxWaitTime
	newConfigRegistry.UserAgent = fmt.Sprintf(
		"terraform-provider/%s_ionos-cloud-sdk-go-container-cr/%s_hashicorp-terraform/%s_terraform-plugin-sdk/%s_os/%s_arch/%s",
		version, cr.Version, terraformVersion, meta.SDKVersionString(), runtime.GOOS, runtime.GOARCH) //nolint:staticcheck
	client := &Client{
		sdkClient: *cr.NewAPIClient(newConfigRegistry),
	}
	client.sdkClient.GetConfig().HTTPClient = &http.Client{Transport: utils.CreateTransport(insecure)}

	return client
}
