package cdn

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"runtime"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/meta"
	cdn "github.com/ionos-cloud/sdk-go-cdn"

	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils/constant"
)

// Client is a struct that defines the CDN client
type Client struct {
	SdkClient *cdn.APIClient
}

// NewCdnClient returns a new CDN client
func NewCdnClient(username, password, token, url, version, terraformVersion string) *Client {
	newConfigCdn := cdn.NewConfiguration(username, password, token, url)

	if os.Getenv(constant.IonosDebug) != "" {
		newConfigCdn.Debug = true
	}
	newConfigCdn.MaxRetries = constant.MaxRetries
	newConfigCdn.MaxWaitTime = constant.MaxWaitTime

	newConfigCdn.HTTPClient = &http.Client{Transport: utils.CreateTransport()}
	newConfigCdn.UserAgent = fmt.Sprintf(
		"terraform-provider/%s_ionos-cloud-sdk-go-cdn/%s_hashicorp-terraform/%s_terraform-plugin-sdk/%s_os/%s_arch/%s",
		version, cdn.Version, terraformVersion, meta.SDKVersionString(), runtime.GOOS, runtime.GOARCH) //nolint:staticcheck

	return &Client{
		SdkClient: cdn.NewAPIClient(newConfigCdn),
	}
}

// IsDistributionReady checks if the distribution is ready
func (c *Client) IsDistributionReady(ctx context.Context, d *schema.ResourceData) (bool, error) {
	distributionID := d.Id()
	distribution, _, err := c.SdkClient.DistributionsApi.DistributionsFindById(ctx, distributionID).Execute()
	if err != nil {
		return true, fmt.Errorf("status check failed for MariaDB distribution with ID: %v, error: %w", distributionID, err)
	}

	if distribution.Metadata == nil || distribution.Metadata.State == nil {
		return false, fmt.Errorf("distribution metadata or state is empty for MariaDB distribution with ID: %v", distributionID)
	}

	log.Printf("[INFO] state of the MariaDB distribution with ID: %v is: %s ", distributionID, *distribution.Metadata.State)
	return strings.EqualFold(*distribution.Metadata.State, constant.Available), nil
}
