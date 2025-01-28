package containerregistry

import (
	"fmt"
	"net/http"
	"os"
	"runtime"

	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils/constant"

	"github.com/hashicorp/terraform-plugin-sdk/v2/meta"
	cr "github.com/ionos-cloud/sdk-go-container-registry"

	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils"
)

type Client struct {
	sdkClient *cr.APIClient
}

func NewClient(username, password, token, url, version, terraformVersion string, insecure bool) *Client {
	newConfigRegistry := cr.NewConfiguration(username, password, token, url)

	if os.Getenv(constant.IonosDebug) != "" {
		newConfigRegistry.Debug = true
	}
	newConfigRegistry.MaxRetries = constant.MaxRetries
	newConfigRegistry.MaxWaitTime = constant.MaxWaitTime

	newConfigRegistry.HTTPClient = &http.Client{Transport: utils.CreateTransport(insecure)}
	newConfigRegistry.UserAgent = fmt.Sprintf(
		"terraform-provider/%s_ionos-cloud-sdk-go-container-cr/%s_hashicorp-terraform/%s_terraform-plugin-sdk/%s_os/%s_arch/%s",
		version, cr.Version, terraformVersion, meta.SDKVersionString(), runtime.GOOS, runtime.GOARCH)

	return &Client{
		sdkClient: cr.NewAPIClient(newConfigRegistry),
	}
}
