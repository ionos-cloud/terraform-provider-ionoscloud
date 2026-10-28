package mongov2

import (
	"context"

	mongov3 "github.com/ionos-cloud/sdk-go-bundle/products/dbaas/mongo/v3"
	"github.com/ionos-cloud/sdk-go-bundle/shared"
)

// GetSnapshotLocation retrieves a snapshot location by ID.
func (c *Client) GetSnapshotLocation(ctx context.Context, locationID string) (mongov3.SnapshotLocationRead, *shared.APIResponse, error) {
	location, apiResponse, err := c.sdkClient.SnapshotLocationsApi.SnapshotlocationsFindById(ctx, locationID).Execute()
	apiResponse.LogInfo()
	return location, apiResponse, err
}

// ListSnapshotLocations retrieves a list of all snapshot locations.
func (c *Client) ListSnapshotLocations(ctx context.Context) (mongov3.SnapshotLocationReadList, *shared.APIResponse, error) {
	const pageSize int32 = 100
	var (
		all         []mongov3.SnapshotLocationRead
		offset      int32
		apiResponse *shared.APIResponse
	)
	for {
		page, resp, err := c.sdkClient.SnapshotLocationsApi.SnapshotlocationsGet(ctx).Offset(offset).Limit(pageSize).Execute()
		resp.LogInfo()
		apiResponse = resp
		if err != nil {
			return mongov3.SnapshotLocationReadList{}, apiResponse, err
		}
		all = append(all, page.Items...)
		if len(page.Items) < int(pageSize) {
			page.Items = all
			return page, apiResponse, nil
		}
		offset += pageSize
	}
}
