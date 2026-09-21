package inmemorydbv2

import (
	"context"

	inmemorydbv3 "github.com/ionos-cloud/sdk-go-bundle/products/dbaas/inmemorydb/v3"
	"github.com/ionos-cloud/sdk-go-bundle/shared"
)

// GetSnapshotLocation retrieves a snapshot location by ID.
func (c *Client) GetSnapshotLocation(ctx context.Context, snapshotLocationID string) (inmemorydbv3.SnapshotLocationRead, *shared.APIResponse, error) {
	location, apiResponse, err := c.sdkClient.SnapshotLocationsApi.SnapshotlocationsFindById(ctx, snapshotLocationID).Execute()
	apiResponse.LogInfo()
	return location, apiResponse, err
}

// ListSnapshotLocations retrieves the list of all snapshot locations.
func (c *Client) ListSnapshotLocations(ctx context.Context) (inmemorydbv3.SnapshotLocationReadList, *shared.APIResponse, error) {
	locations, apiResponse, err := c.sdkClient.SnapshotLocationsApi.SnapshotlocationsGet(ctx).Execute()
	apiResponse.LogInfo()
	return locations, apiResponse, err
}
