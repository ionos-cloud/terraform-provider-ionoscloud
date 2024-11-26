package objectstoragemanagement

import (
	"fmt"
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
func NewClient(username, password, token, url, version, terraformVersion string, insecure bool) *Client {
	newObjectStorageManagementConfig := objectstoragemanagement.NewConfiguration(username, password, token, url)

	if os.Getenv(constant.IonosDebug) != "" {
		newObjectStorageManagementConfig.Debug = true
	}
	newObjectStorageManagementConfig.MaxRetries = constant.MaxRetries
	newObjectStorageManagementConfig.MaxWaitTime = constant.MaxWaitTime
	newObjectStorageManagementConfig.HTTPClient = &http.Client{Transport: utils.CreateTransport(insecure)}
	newObjectStorageManagementConfig.UserAgent = fmt.Sprintf(
		"terraform-provider/%s_ionos-cloud-sdk-go-object-storage-management/%s_hashicorp-terraform/%s_terraform-plugin-sdk/%s_os/%s_arch/%s",
		version, objectstoragemanagement.Version, terraformVersion, meta.SDKVersionString(), runtime.GOOS, runtime.GOARCH) //nolint:staticcheck

	return &Client{client: objectstoragemanagement.NewAPIClient(newObjectStorageManagementConfig)}
}
