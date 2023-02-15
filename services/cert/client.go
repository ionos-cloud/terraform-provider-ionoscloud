package cert

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/meta"
	certmanager "github.com/ionos-cloud/sdk-go-bundle/products/cert"
	"github.com/ionos-cloud/sdk-go-bundle/shared"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils"
	"net/http"
	"runtime"
)

type Client struct {
	sdkClient *certmanager.APIClient
}

func NewClient(username, password, token, url, version, terraformVersion string) *Client {
	certConfig := shared.NewConfiguration(username, password, token, url)

	certConfig.MaxRetries = utils.MaxRetries
	certConfig.MaxWaitTime = utils.MaxWaitTime

	certConfig.HTTPClient = &http.Client{Transport: utils.CreateTransport()}
	certConfig.UserAgent = fmt.Sprintf(
		"terraform-provider/%s_ionos-cloud-sdk-go-cert-manager/%s_hashicorp-terraform/%s_terraform-plugin-sdk/%s_os/%s_arch/%s",
		version, certmanager.Version, terraformVersion, meta.SDKVersionString(), runtime.GOOS, runtime.GOARCH)

	return &Client{
		sdkClient: certmanager.NewAPIClient(certConfig),
	}
}
