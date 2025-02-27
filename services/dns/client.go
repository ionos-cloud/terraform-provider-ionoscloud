package dns

import (
	"fmt"
	"net/http"
	"runtime"

	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils"

	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils/constant"

	"github.com/hashicorp/terraform-plugin-sdk/v2/meta"
	dns "github.com/ionos-cloud/sdk-go-bundle/products/dns/v2"
	"github.com/ionos-cloud/sdk-go-bundle/shared"
)

type Client struct {
	sdkClient dns.APIClient
}

func NewClient(username, password, token, url, version, terraformVersion string, insecure bool) *Client {
	newConfigDNS := shared.NewConfiguration(username, password, token, url)

	newConfigDNS.MaxRetries = constant.MaxRetries
	newConfigDNS.MaxWaitTime = constant.MaxWaitTime
	newConfigDNS.UserAgent = fmt.Sprintf(
		"terraform-provider/%s_ionos-cloud-sdk-go-dns/%s_hashicorp-terraform/%s_terraform-plugin-sdk/%s_os/%s_arch/%s",
		version, dns.Version, terraformVersion, meta.SDKVersionString(), runtime.GOOS, runtime.GOARCH) //nolint:staticcheck
	client := &Client{
		sdkClient: *dns.NewAPIClient(newConfigDNS),
	}
	client.sdkClient.GetConfig().HTTPClient = &http.Client{Transport: utils.CreateTransport(insecure)}

	return client
}
