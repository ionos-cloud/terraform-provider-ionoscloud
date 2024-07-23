package nfs

import (
	"fmt"
	"net/http"
	"os"
	"runtime"

	"github.com/hashicorp/terraform-plugin-sdk/v2/meta"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils/constant"

	sdk "github.com/ionos-cloud/sdk-go-nfs"
)

type Client struct {
	sdkClient sdk.APIClient
}

func NewClient(username, password, token, url, version, terraformVersion string) *Client {
	config := sdk.NewConfiguration(username, password, token, url)

	if os.Getenv(constant.IonosDebug) != "" {
		config.Debug = true
	}
	config.MaxRetries = constant.MaxRetries
	config.MaxWaitTime = constant.MaxWaitTime
	config.HTTPClient = &http.Client{Transport: utils.CreateTransport()}
	config.UserAgent = fmt.Sprintf(
		"terraform-provider/%s_ionos-cloud-sdk-go-nfs/%s_hashicorp-terraform/%s_terraform-plugin-sdk/%s_os/%s_arch/%s",
		version, sdk.Version, terraformVersion, meta.SDKVersionString(), runtime.GOOS, runtime.GOARCH,
	)

	return &Client{sdkClient: *sdk.NewAPIClient(config)}
}

func (c *Client) Location(location string) *Client {
	var locationToURL = map[string]string{
		"de/fra":    "nfs.de-fra.ionos.com",
		"de/txl":    "nfs.de-txl.ionos.com",
		"qa/de/txl": "qa.nfs.de-txl.ionos.com",
	}

	c.sdkClient.GetConfig().Servers = sdk.ServerConfigurations{
		{
			URL: locationToURL[location],
		},
	}

	return c
}
