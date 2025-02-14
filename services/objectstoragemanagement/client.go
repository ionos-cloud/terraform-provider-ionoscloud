package objectstoragemanagement

import (
	"fmt"
	"github.com/ionos-cloud/sdk-go-bundle/shared"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils/bundle"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils/loadedconfig"
	"net/http"
	"os"
	"runtime"

	"github.com/hashicorp/terraform-plugin-sdk/v2/meta"

	objectstoragemanagement "github.com/ionos-cloud/sdk-go-object-storage-management"

	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils/constant"
)

// Client is a wrapper around the S3 client.
type Client struct {
	client *objectstoragemanagement.APIClient
}

// GetBaseClient returns the base client.
func (c *Client) GetBaseClient() *objectstoragemanagement.APIClient {
	return c.client
}

// NewClient creates a new S3 client with the given credentials and region.
func NewClient(clientOptions bundle.ClientOptions, sharedLoadedConfig *shared.LoadedConfig) *Client {
	loadedconfig.SetClientOptionsFromFileConfig(&clientOptions, sharedLoadedConfig, shared.ObjectStorage) //todo cguran is this ObjectStorageManagement?
	newObjectStorageManagementConfig := objectstoragemanagement.NewConfiguration(clientOptions.Credentials.Username, clientOptions.Credentials.Password,
		clientOptions.Credentials.Token, clientOptions.Endpoint)

	if os.Getenv(constant.IonosDebug) != "" {
		newObjectStorageManagementConfig.Debug = true
	}
	newObjectStorageManagementConfig.MaxRetries = constant.MaxRetries
	newObjectStorageManagementConfig.MaxWaitTime = constant.MaxWaitTime
	newObjectStorageManagementConfig.HTTPClient = &http.Client{Transport: utils.CreateTransport(clientOptions.SkipTLSVerify)}
	newObjectStorageManagementConfig.UserAgent = fmt.Sprintf(
		"terraform-provider/ionos-cloud-sdk-go-object-storage-management/%s_hashicorp-terraform/%s_terraform-plugin-sdk/%s_os/%s_arch/%s",
		objectstoragemanagement.Version, clientOptions.TerraformVersion, meta.SDKVersionString(), runtime.GOOS, runtime.GOARCH) //nolint:staticcheck

	return &Client{client: objectstoragemanagement.NewAPIClient(newObjectStorageManagementConfig)}
}
