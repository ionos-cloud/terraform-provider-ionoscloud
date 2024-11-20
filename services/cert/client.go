package cert

import (
	"fmt"
	"net/http"
	"os"
	"runtime"

	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils/constant"

	"github.com/hashicorp/terraform-plugin-sdk/v2/meta"
	certmanager "github.com/ionos-cloud/sdk-go-cert-manager"

	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils"
)

type Client struct {
	sdkClient *certmanager.APIClient
}

func NewClient(username, password, token, url, version, terraformVersion string, insecure bool) *Client {
	certConfig := certmanager.NewConfiguration(username, password, token, url)

	if os.Getenv(constant.IonosDebug) != "" {
		certConfig.Debug = true
	}
	certConfig.MaxRetries = constant.MaxRetries
	certConfig.MaxWaitTime = constant.MaxWaitTime

	certConfig.HTTPClient = &http.Client{Transport: utils.CreateTransport(insecure)}
	certConfig.UserAgent = fmt.Sprintf(
		"terraform-provider/%s_ionos-cloud-sdk-go-cert-manager/%s_hashicorp-terraform/%s_terraform-plugin-sdk/%s_os/%s_arch/%s",
		version, certmanager.Version, terraformVersion, meta.SDKVersionString(), runtime.GOOS, runtime.GOARCH)

	return &Client{
		sdkClient: certmanager.NewAPIClient(certConfig),
	}
}
