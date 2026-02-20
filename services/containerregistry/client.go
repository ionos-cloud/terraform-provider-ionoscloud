package containerregistry

import (
	"fmt"
	"net/http"
	"runtime"

	"github.com/ionos-cloud/sdk-go-bundle/shared"

	"github.com/hashicorp/terraform-plugin-sdk/v2/meta"
	cr "github.com/ionos-cloud/sdk-go-bundle/products/containerregistry/v2"
	"github.com/ionos-cloud/sdk-go-bundle/shared/fileconfiguration"

	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/clientoptions"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils/constant"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils/loadedconfig"
)

// Client is a wrapper for the sdk client
type Client struct {
	sdkClient cr.APIClient
}

// NewClient creates a new Container Registry client
func NewClient(clientOptions clientoptions.TerraformClientOptions, fileConfig *fileconfiguration.FileConfig) *Client {
	loadedconfig.SetGlobalClientOptionsFromFileConfig(&clientOptions, fileConfig, fileconfiguration.ContainerRegistry)
	config := shared.NewConfiguration(clientOptions.Credentials.Username, clientOptions.Credentials.Password, clientOptions.Credentials.Token, clientOptions.Endpoint)

	config.MaxRetries = constant.MaxRetries
	config.MaxWaitTime = constant.MaxWaitTime
	config.UserAgent = fmt.Sprintf(
		"terraform-provider/%s_ionos-cloud-sdk-go-container-cr/%s_hashicorp-terraform/%s_terraform-plugin-sdk/%s_os/%s_arch/%s",
		clientOptions.Version, cr.Version, clientOptions.TerraformVersion,
		meta.SDKVersionString(), runtime.GOOS, runtime.GOARCH, //nolint:staticcheck
	)
	client := &Client{
		sdkClient: *cr.NewAPIClient(config),
	}
	client.sdkClient.GetConfig().HTTPClient = &http.Client{Transport: shared.CreateTransport(clientOptions.SkipTLSVerify, clientOptions.Certificate)}

	return client
}

// NewClientFromConfig creates a *Client from an existing shared.Configuration
func NewClientFromConfig(config *shared.Configuration) *Client {
	return &Client{
		sdkClient: *cr.NewAPIClient(config),
	}
}
