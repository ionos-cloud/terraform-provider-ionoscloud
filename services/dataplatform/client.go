package dataplatform

import (
	"fmt"
	"net/http"
	"runtime"

	"github.com/ionos-cloud/sdk-go-bundle/shared"

	"github.com/hashicorp/terraform-plugin-sdk/v2/meta"
	dataplatform "github.com/ionos-cloud/sdk-go-bundle/products/dataplatform/v2"
	"github.com/ionos-cloud/sdk-go-bundle/shared/fileconfiguration"

	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/clientoptions"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils/constant"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils/loadedconfig"
)

// Client - wrapper over sdk client, to allow for service layer
type Client struct {
	sdkClient dataplatform.APIClient
}

// NewClient - creates a new dataplatform client
func NewClient(clientOptions clientoptions.TerraformClientOptions, fileConfig *fileconfiguration.FileConfig) *Client {
	loadedconfig.SetGlobalClientOptionsFromFileConfig(&clientOptions, fileConfig, fileconfiguration.Dataplatform)
	config := shared.NewConfiguration(clientOptions.Credentials.Username, clientOptions.Credentials.Password, clientOptions.Credentials.Token, clientOptions.Endpoint)
	config.MaxRetries = constant.MaxRetries
	config.MaxWaitTime = constant.MaxWaitTime
	config.UserAgent = fmt.Sprintf(
		"terraform-provider/%s_ionos-cloud-sdk-go-dataplatform/%s_hashicorp-terraform/%s_terraform-plugin-sdk/%s_os/%s_arch/%s",
		clientOptions.Version, dataplatform.Version, clientOptions.TerraformVersion,
		meta.SDKVersionString(), runtime.GOOS, runtime.GOARCH, //nolint:staticcheck
	)
	client := &Client{
		sdkClient: *dataplatform.NewAPIClient(config),
	}
	client.sdkClient.GetConfig().HTTPClient = &http.Client{Transport: shared.CreateTransport(clientOptions.SkipTLSVerify, clientOptions.Certificate)}
	return client
}
