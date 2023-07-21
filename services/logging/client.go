package logging

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/meta"
	logging "github.com/ionos-cloud/sdk-go-logging"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils"
	"net/http"
	"os"
	"runtime"
)

type Client struct {
	sdkClient logging.APIClient
}

func NewClient(username, password, token, url, version, terraformVersion string) *Client {
	newConfigLogging := logging.NewConfiguration(username, password, token, url)

	if os.Getenv(utils.IonosDebug) != "" {
		newConfigLogging.Debug = true
	}
	newConfigLogging.MaxRetries = utils.MaxRetries
	newConfigLogging.MaxWaitTime = utils.MaxWaitTime
	newConfigLogging.HTTPClient = &http.Client{Transport: utils.CreateTransport()}
	newConfigLogging.UserAgent = fmt.Sprintf(
		"terraform-provider/%s_ionos-cloud-sdk-go-logging/%s_hashicorp-terraform/%s_terraform-plugin-sdk/%s_os/%s_arch/%s",
		version, logging.Version, terraformVersion, meta.SDKVersionString(), runtime.GOOS, runtime.GOARCH)

	return &Client{sdkClient: *logging.NewAPIClient(newConfigLogging)}
}
