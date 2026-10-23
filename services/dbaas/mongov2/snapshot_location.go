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
	locations, apiResponse, err := c.sdkClient.SnapshotLocationsApi.SnapshotlocationsGet(ctx).Execute()
	apiResponse.LogInfo()
	return locations, apiResponse, err
}
