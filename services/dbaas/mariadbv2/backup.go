package mariadbv2

import (
	"context"

	mariadbv3 "github.com/ionos-cloud/sdk-go-bundle/products/dbaas/mariadb/v3"
	"github.com/ionos-cloud/sdk-go-bundle/shared"
)

// GetBackup retrieves a backup by ID.
func (c *Client) GetBackup(ctx context.Context, backupID string) (mariadbv3.BackupRead, *shared.APIResponse, error) {
	backup, apiResponse, err := c.sdkClient.BackupsApi.BackupsFindById(ctx, backupID).Execute()
	apiResponse.LogInfo()
	return backup, apiResponse, err
}

// ListBackups retrieves a list of backups with an optional cluster ID filter.
func (c *Client) ListBackups(ctx context.Context, filterClusterID string) (mariadbv3.BackupReadList, *shared.APIResponse, error) {
	request := c.sdkClient.BackupsApi.BackupsGet(ctx)
	if filterClusterID != "" {
		request = request.FilterClusterId(filterClusterID)
	}
	backups, apiResponse, err := request.Execute()
	apiResponse.LogInfo()
	return backups, apiResponse, err
}
