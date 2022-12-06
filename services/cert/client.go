package cert

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/meta"
	certmanager "github.com/ionos-cloud/sdk-go-cert-manager"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils"
	"net/http"
	"os"
	"runtime"
)

type Client struct {
	sdkClient *certmanager.APIClient
}

func NewClientService(username, password, token, url, version, terraformVersion string) *Client {
	certConfig := certmanager.NewConfiguration(username, password, token, url)

	if os.Getenv(utils.IonosDebug) != "" {
		certConfig.Debug = true
	}
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
