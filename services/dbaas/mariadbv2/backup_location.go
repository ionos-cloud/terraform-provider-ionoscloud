package mariadbv2

import (
	"context"

	mariadbv3 "github.com/ionos-cloud/sdk-go-bundle/products/dbaas/mariadb/v3"
	"github.com/ionos-cloud/sdk-go-bundle/shared"
)

// GetBackupLocation retrieves a backup location by ID.
func (c *Client) GetBackupLocation(ctx context.Context, backupLocationID string) (mariadbv3.BackupLocationRead, *shared.APIResponse, error) {
	backupLocation, apiResponse, err := c.sdkClient.BackupLocationsApi.BackuplocationsFindById(ctx, backupLocationID).Execute()
	apiResponse.LogInfo()
	return backupLocation, apiResponse, err
}

// ListBackupLocations retrieves a list of all backup locations.
func (c *Client) ListBackupLocations(ctx context.Context) (mariadbv3.BackupLocationReadList, *shared.APIResponse, error) {
	backupLocations, apiResponse, err := c.sdkClient.BackupLocationsApi.BackuplocationsGet(ctx).Execute()
	apiResponse.LogInfo()
	return backupLocations, apiResponse, err
}
