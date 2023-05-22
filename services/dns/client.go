package dns

import (
	"fmt"
	"net/http"
	"os"
	"runtime"

	"github.com/hashicorp/terraform-plugin-sdk/v2/meta"
	dns "github.com/ionos-cloud/sdk-go-dns"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils"
)

type Client struct {
	sdkClient dns.APIClient
}

func NewClient(username, password, token, url, version, terraformVersion string) *Client {
	newConfigDNS := dns.NewConfiguration(username, password, token, url)

	if os.Getenv(utils.IonosDebug) != "" {
		newConfigDNS.Debug = true
	}
	newConfigDNS.MaxRetries = utils.MaxRetries
	newConfigDNS.MaxWaitTime = utils.MaxWaitTime
	newConfigDNS.HTTPClient = &http.Client{Transport: utils.CreateTransport()}
	newConfigDNS.UserAgent = fmt.Sprintf(
		"terraform-provider/%s_ionos-cloud-sdk-go-dns/%s_hashicorp-terraform/%s_terraform-plugin-sdk/%s_os/%s_arch/%s",
		version, dns.Version, terraformVersion, meta.SDKVersionString(), runtime.GOOS, runtime.GOARCH)

	return &Client{sdkClient: *dns.NewAPIClient(newConfigDNS)}
}
