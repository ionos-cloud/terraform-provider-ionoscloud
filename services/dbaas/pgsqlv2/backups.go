package pgsqlv2

import (
	"context"

	pgsqlv2 "github.com/ionos-cloud/pgsqlv2"
	"github.com/ionos-cloud/sdk-go-bundle/shared"
)

// GetBackup retrieves a backup by its ID.
func (c *Client) GetBackup(ctx context.Context, backupID string) (pgsqlv2.BackupRead, *shared.APIResponse, error) {
	backup, apiResponse, err := c.sdkClient.BackupsApi.BackupsFindById(ctx, backupID).Execute()
	apiResponse.LogInfo()
	return backup, apiResponse, err
}

// ListBackups retrieves a list of backups. An optional clusterID filter can be used.
func (c *Client) ListBackups(ctx context.Context, clusterID string) (pgsqlv2.BackupReadList, *shared.APIResponse, error) {
	request := c.sdkClient.BackupsApi.BackupsGet(ctx)
	if clusterID != "" {
		request = request.FilterClusterId(clusterID)
	}
	backups, apiResponse, err := request.Execute()
	apiResponse.LogInfo()
	return backups, apiResponse, err
}
