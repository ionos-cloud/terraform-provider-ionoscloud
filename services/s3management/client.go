package s3management

import (
	"fmt"
	"net/http"
	"os"
	"runtime"

	"github.com/hashicorp/terraform-plugin-sdk/v2/meta"
	s3management "github.com/ionos-cloud/sdk-go-s3-management"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils/constant"
)

// Client is a wrapper around the S3 client.
type Client struct {
	client *s3management.APIClient
}

// GetBaseClient returns the base client.
func (c *Client) GetBaseClient() *s3management.APIClient {
	return c.client
}

// NewClient creates a new S3 client with the given credentials and region.
func NewClient(username, password, token, url, version, terraformVersion string) *Client {
	newS3managementConfig := s3management.NewConfiguration(username, password, token, url)

	if os.Getenv(constant.IonosDebug) != "" {
		newS3managementConfig.Debug = true
	}
	newS3managementConfig.MaxRetries = constant.MaxRetries
	newS3managementConfig.MaxWaitTime = constant.MaxWaitTime
	newS3managementConfig.HTTPClient = &http.Client{Transport: utils.CreateTransport()}
	newS3managementConfig.UserAgent = fmt.Sprintf(
		"terraform-provider/%s_ionos-cloud-sdk-go-s3-management/%s_hashicorp-terraform/%s_terraform-plugin-sdk/%s_os/%s_arch/%s",
		version, s3management.Version, terraformVersion, meta.SDKVersionString(), runtime.GOOS, runtime.GOARCH) //nolint:staticcheck

	return &Client{client: s3management.NewAPIClient(newS3managementConfig)}
}
