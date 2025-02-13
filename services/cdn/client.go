package cdn

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/meta"
	cdn "github.com/ionos-cloud/sdk-go-bundle/products/cdn/v2"
	"github.com/ionos-cloud/sdk-go-bundle/shared"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils/bundle"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils/loadedconfig"
	"net/http"
	"runtime"

	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils/constant"
)

// Client is a struct that defines the CDN client
type Client struct {
	SdkClient *cdn.APIClient
}

//
//// NewCDNClient returns a new CDN client
//func NewCDNClient(username, password, token, url, version, terraformVersion string, insecure bool) *sdkClient {
//	newConfigCDN := shared.NewConfiguration(username, password, token, url)
//	newConfigCDN.MaxRetries = constant.MaxRetries
//	newConfigCDN.MaxWaitTime = constant.MaxWaitTime
//
//	newConfigCDN.HTTPClient = &http.sdkClient{Transport: utils.CreateTransport(insecure)}
//	newConfigCDN.UserAgent = fmt.Sprintf(
//		"terraform-provider/%s_ionos-cloud-sdk-go-cdn/%s_hashicorp-terraform/%s_terraform-plugin-sdk/%s_os/%s_arch/%s",
//		version, cdn.Version, terraformVersion, meta.SDKVersionString(), runtime.GOOS, runtime.GOARCH) //nolint:staticcheck
//
//	return &sdkClient{
//		SdkClient: cdn.NewAPIClient(newConfigCDN),
//	}
//}

// NewCDNClient returns a new CDN client
func NewCDNClient(options bundle.ClientOptions, sharedLoadedConfig *shared.LoadedConfig) *Client {
	loadedconfig.SetClientOptionsFromLoadedConfig(&options, sharedLoadedConfig, shared.CDN)

	newConfigCDN := shared.NewConfiguration(options.Credentials.Username, options.Credentials.Password, options.Credentials.Token, options.Endpoint)
	newConfigCDN.MaxRetries = constant.MaxRetries
	newConfigCDN.MaxWaitTime = constant.MaxWaitTime

	newConfigCDN.HTTPClient = &http.Client{Transport: utils.CreateTransport(options.SkipTLSVerify)}
	newConfigCDN.UserAgent = fmt.Sprintf(
		"terraform-provider/_ionos-cloud-sdk-go-cdn/%s_hashicorp-terraform/%s_terraform-plugin-sdk/%s_os/%s_arch/%s",
		cdn.Version, options.TerraformVersion, meta.SDKVersionString(), runtime.GOOS, runtime.GOARCH) //nolint:staticcheck

	return &Client{
		SdkClient: cdn.NewAPIClient(newConfigCDN),
	}
}
