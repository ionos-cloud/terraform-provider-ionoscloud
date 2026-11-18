package pgsqlv2

import (
	"context"

	pgsqlv2 "github.com/ionos-cloud/sdk-go-bundle/products/dbaas/psql/v3"
	"github.com/ionos-cloud/sdk-go-bundle/shared"
)

// GetBackup retrieves a backup by its ID.
func (c *Client) GetBackup(ctx context.Context, backupID string) (pgsqlv2.BackupRead, *shared.APIResponse, error) {
	backup, apiResponse, err := c.sdkClient.BackupsApi.BackupsFindById(ctx, backupID).Execute()
	apiResponse.LogInfo()
	return backup, apiResponse, err
}

// ListBackups retrieves a list of backups, following pagination. An optional clusterID filter can be used.
func (c *Client) ListBackups(ctx context.Context, clusterID string) (pgsqlv2.BackupReadList, *shared.APIResponse, error) {
	const pageSize int32 = 100
	var (
		all         []pgsqlv2.BackupRead
		offset      int32
		apiResponse *shared.APIResponse
	)
	for {
		request := c.sdkClient.BackupsApi.BackupsGet(ctx).Offset(offset).Limit(pageSize)
		if clusterID != "" {
			request = request.FilterClusterId(clusterID)
		}
		page, resp, err := request.Execute()
		resp.LogInfo()
		apiResponse = resp
		if err != nil {
			return pgsqlv2.BackupReadList{}, apiResponse, err
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
