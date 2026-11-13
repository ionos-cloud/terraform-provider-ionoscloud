package mongov2

import (
	"context"

	mongov3 "github.com/ionos-cloud/sdk-go-bundle/products/dbaas/mongo/v3"
	"github.com/ionos-cloud/sdk-go-bundle/shared"
)

// ListVersions retrieves a list of all supported MongoDB versions.
func (c *Client) ListVersions(ctx context.Context) (mongov3.MongoVersionReadList, *shared.APIResponse, error) {
	versions, apiResponse, err := c.sdkClient.VersionsApi.VersionsGet(ctx).Execute()
	apiResponse.LogInfo()
	return versions, apiResponse, err
}
