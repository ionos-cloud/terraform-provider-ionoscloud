package cert

import (
	"fmt"
	"net/http"
	"runtime"

	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils"

	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils/constant"

	"github.com/hashicorp/terraform-plugin-sdk/v2/meta"
	certmanager "github.com/ionos-cloud/sdk-go-bundle/products/cert/v2"
	"github.com/ionos-cloud/sdk-go-bundle/shared"
)

type Client struct {
	sdkClient certmanager.APIClient
}

func NewClient(username, password, token, url, version, terraformVersion string, insecure bool) *Client {
	certConfig := shared.NewConfiguration(username, password, token, url)

	certConfig.MaxRetries = constant.MaxRetries
	certConfig.MaxWaitTime = constant.MaxWaitTime
	certConfig.UserAgent = fmt.Sprintf(
		"terraform-provider/%s_ionos-cloud-sdk-go-cert-manager/%s_hashicorp-terraform/%s_terraform-plugin-sdk/%s_os/%s_arch/%s",
		version, certmanager.Version, terraformVersion, meta.SDKVersionString(), runtime.GOOS, runtime.GOARCH)
	certClient := &Client{
		sdkClient: *certmanager.NewAPIClient(certConfig),
	}
	certClient.sdkClient.GetConfig().HTTPClient = &http.Client{Transport: utils.CreateTransport(insecure)}

	return certClient
}
