package inmemorydbv2

import (
	"context"

	inmemorydbv3 "github.com/ionos-cloud/sdk-go-bundle/products/dbaas/inmemorydb/v3"
	"github.com/ionos-cloud/sdk-go-bundle/shared"
)

// GetSnapshot retrieves a snapshot by ID.
func (c *Client) GetSnapshot(ctx context.Context, snapshotID string) (inmemorydbv3.SnapshotRead, *shared.APIResponse, error) {
	snapshot, apiResponse, err := c.sdkClient.SnapshotsApi.SnapshotsFindById(ctx, snapshotID).Execute()
	apiResponse.LogInfo()
	return snapshot, apiResponse, err
}

// ListSnapshots retrieves a list of snapshots with an optional cluster ID filter.
func (c *Client) ListSnapshots(ctx context.Context, filterClusterID string) (inmemorydbv3.SnapshotReadList, *shared.APIResponse, error) {
	request := c.sdkClient.SnapshotsApi.SnapshotsGet(ctx)
	if filterClusterID != "" {
		request = request.FilterClusterId(filterClusterID)
	}
	snapshots, apiResponse, err := request.Execute()
	apiResponse.LogInfo()
	return snapshots, apiResponse, err
}
