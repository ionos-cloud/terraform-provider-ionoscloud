package containerregistry

import (
	"fmt"
	"github.com/ionos-cloud/sdk-go-bundle/shared"
	"net/http"
	"os"
	"runtime"

	"github.com/hashicorp/terraform-plugin-sdk/v2/meta"
	"github.com/ionos-cloud/sdk-go-bundle/shared/fileconfiguration"
	cr "github.com/ionos-cloud/sdk-go-container-registry"

	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils/bundle"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils/constant"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils/loadedconfig"
)

// Client is a wrapper for the sdk client
type Client struct {
	sdkClient *cr.APIClient
}

// NewClient creates a new Container Registry client
func NewClient(clientOptions bundle.ClientOptions, fileConfig *fileconfiguration.FileConfig) *Client {
	loadedconfig.SetGlobalClientOptionsFromFileConfig(&clientOptions, fileConfig, fileconfiguration.ContainerRegistry)
	config := cr.NewConfiguration(clientOptions.Credentials.Username, clientOptions.Credentials.Password, clientOptions.Credentials.Token, clientOptions.Endpoint)

	if os.Getenv(constant.IonosDebug) != "" {
		config.Debug = true
	}
	config.MaxRetries = constant.MaxRetries
	config.MaxWaitTime = constant.MaxWaitTime
	config.HTTPClient = http.DefaultClient
	config.HTTPClient.Transport = shared.CreateTransport(clientOptions.SkipTLSVerify, clientOptions.Certificate)
	config.UserAgent = fmt.Sprintf(
		"terraform-provider/_ionos-cloud-sdk-go-container-cr/%s_hashicorp-terraform/%s_terraform-plugin-sdk/%s_os/%s_arch/%s",
		cr.Version, clientOptions.TerraformVersion, meta.SDKVersionString(), runtime.GOOS, runtime.GOARCH)

	return &Client{
		sdkClient: cr.NewAPIClient(config),
	}
}
