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

// ListSnapshots retrieves a list of snapshots, following pagination. An optional cluster ID filter can be used.
func (c *Client) ListSnapshots(ctx context.Context, filterClusterID string) (inmemorydbv3.SnapshotReadList, *shared.APIResponse, error) {
	const pageSize int32 = 100
	var (
		all         []inmemorydbv3.SnapshotRead
		offset      int32
		apiResponse *shared.APIResponse
	)
	for {
		request := c.sdkClient.SnapshotsApi.SnapshotsGet(ctx).Offset(offset).Limit(pageSize)
		if filterClusterID != "" {
			request = request.FilterClusterId(filterClusterID)
		}
		page, resp, err := request.Execute()
		resp.LogInfo()
		apiResponse = resp
		if err != nil {
			return inmemorydbv3.SnapshotReadList{}, apiResponse, err
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
