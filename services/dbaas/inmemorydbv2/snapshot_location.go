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

// ListSnapshotLocations retrieves the list of all snapshot locations, following pagination.
func (c *Client) ListSnapshotLocations(ctx context.Context) (inmemorydbv3.SnapshotLocationReadList, *shared.APIResponse, error) {
	const pageSize int32 = 100
	var (
		all         []inmemorydbv3.SnapshotLocationRead
		offset      int32
		apiResponse *shared.APIResponse
	)
	for {
		page, resp, err := c.sdkClient.SnapshotLocationsApi.SnapshotlocationsGet(ctx).Offset(offset).Limit(pageSize).Execute()
		resp.LogInfo()
		apiResponse = resp
		if err != nil {
			return inmemorydbv3.SnapshotLocationReadList{}, apiResponse, err
		}
		all = append(all, page.Items...)
		// A short (or empty) page marks the end of the result set. Links.HasNext() is
		// deliberately not used: it has been observed to stay populated past the last page.
		if len(page.Items) < int(pageSize) {
			page.Items = all
			return page, apiResponse, nil
		}
		offset += pageSize
	}
}
