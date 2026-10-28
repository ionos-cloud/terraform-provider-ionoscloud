package mongov2

import (
	"context"

	mongov3 "github.com/ionos-cloud/sdk-go-bundle/products/dbaas/mongo/v3"
	"github.com/ionos-cloud/sdk-go-bundle/shared"
)

// ListVersions retrieves a list of all supported MongoDB versions.
func (c *Client) ListVersions(ctx context.Context) (mongov3.MongoVersionReadList, *shared.APIResponse, error) {
	const pageSize int32 = 100
	var (
		all         []mongov3.MongoVersionRead
		offset      int32
		apiResponse *shared.APIResponse
	)
	for {
		page, resp, err := c.sdkClient.VersionsApi.VersionsGet(ctx).Offset(offset).Limit(pageSize).Execute()
		resp.LogInfo()
		apiResponse = resp
		if err != nil {
			return mongov3.MongoVersionReadList{}, apiResponse, err
		}
		all = append(all, page.Items...)
		if len(page.Items) < int(pageSize) {
			page.Items = all
			return page, apiResponse, nil
		}
		offset += pageSize
	}
}
