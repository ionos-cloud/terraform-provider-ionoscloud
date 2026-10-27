package mongov2

import (
	"context"

	mongov3 "github.com/ionos-cloud/sdk-go-bundle/products/dbaas/mongo/v3"
	"github.com/ionos-cloud/sdk-go-bundle/shared"
)

// GetSnapshot retrieves a snapshot by ID.
func (c *Client) GetSnapshot(ctx context.Context, snapshotID string) (mongov3.SnapshotRead, *shared.APIResponse, error) {
	snapshot, apiResponse, err := c.sdkClient.SnapshotsApi.SnapshotsFindById(ctx, snapshotID).Execute()
	apiResponse.LogInfo()
	return snapshot, apiResponse, err
}

// ListSnapshots retrieves a list of snapshots with an optional cluster ID filter.
func (c *Client) ListSnapshots(ctx context.Context, filterClusterID string) (mongov3.SnapshotReadList, *shared.APIResponse, error) {
	const pageSize int32 = 100
	var (
		all         []mongov3.SnapshotRead
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
			return mongov3.SnapshotReadList{}, apiResponse, err
		}
		all = append(all, page.Items...)
		if len(page.Items) < int(pageSize) {
			page.Items = all
			return page, apiResponse, nil
		}
		offset += pageSize
	}
}
