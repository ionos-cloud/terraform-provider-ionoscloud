package pgsqlv2

import (
	"context"

	pgsqlv2 "github.com/ionos-cloud/sdk-go-bundle/products/dbaas/psql/v3"
	"github.com/ionos-cloud/sdk-go-bundle/shared"
)

// GetBackupLocation retrieves a backup location by its ID.
func (c *Client) GetBackupLocation(ctx context.Context, backupLocationID string) (pgsqlv2.BackupLocationRead, *shared.APIResponse, error) {
	backupLocation, apiResponse, err := c.sdkClient.BackupLocationsApi.BackuplocationsFindById(ctx, backupLocationID).Execute()
	apiResponse.LogInfo()
	return backupLocation, apiResponse, err
}

// ListBackupLocations retrieves all available backup locations.
func (c *Client) ListBackupLocations(ctx context.Context) (pgsqlv2.BackupLocationReadList, *shared.APIResponse, error) {
	backupLocations, apiResponse, err := c.sdkClient.BackupLocationsApi.BackuplocationsGet(ctx).Execute()
	apiResponse.LogInfo()
	return backupLocations, apiResponse, err
}
