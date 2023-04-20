package dnsaas

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/meta"
	dnsaas "github.com/ionos-cloud/sdk-go-dnsaas"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils"
	"net/http"
	"os"
	"runtime"
)

type Client struct {
	sdkClient dnsaas.APIClient
}

func NewClient(username, password, token, url, version, terraformVersion string) *Client {
	newConfigDNSaaS := dnsaas.NewConfiguration(username, password, token, url)

	if os.Getenv(utils.IonosDebug) != "" {
		newConfigDNSaaS.Debug = true
	}
	newConfigDNSaaS.MaxRetries = utils.MaxRetries
	newConfigDNSaaS.MaxWaitTime = utils.MaxWaitTime
	newConfigDNSaaS.HTTPClient = &http.Client{Transport: utils.CreateTransport()}
	newConfigDNSaaS.UserAgent = fmt.Sprintf(
		"terraform-provider/%s_ionos-cloud-sdk-go-dnsaas/%s_hashicorp-terraform/%s_terraform-plugin-sdk/%s_os/%s_arch/%s",
		version, dnsaas.Version, terraformVersion, meta.SDKVersionString(), runtime.GOOS, runtime.GOARCH)

	return &Client{sdkClient: *dnsaas.NewAPIClient(newConfigDNSaaS)}
}
