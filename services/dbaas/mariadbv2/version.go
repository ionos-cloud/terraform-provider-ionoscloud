package mariadbv2

import (
	"context"

	mariadbv3 "github.com/ionos-cloud/sdk-go-bundle/products/dbaas/mariadb/v3"
	"github.com/ionos-cloud/sdk-go-bundle/shared"
)

// ListVersions retrieves a list of all supported MariaDB versions.
func (c *Client) ListVersions(ctx context.Context) (mariadbv3.MariadbVersionReadList, *shared.APIResponse, error) {
	versions, apiResponse, err := c.sdkClient.VersionsApi.VersionsGet(ctx).Execute()
	apiResponse.LogInfo()
	return versions, apiResponse, err
}
