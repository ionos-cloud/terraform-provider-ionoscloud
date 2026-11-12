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

// ListBackupLocations retrieves a list of all backup locations, following pagination.
func (c *Client) ListBackupLocations(ctx context.Context) (mariadbv3.BackupLocationReadList, *shared.APIResponse, error) {
	const pageSize int32 = 100
	var (
		all         []mariadbv3.BackupLocationRead
		offset      int32
		apiResponse *shared.APIResponse
	)
	for {
		page, resp, err := c.sdkClient.BackupLocationsApi.BackuplocationsGet(ctx).Offset(offset).Limit(pageSize).Execute()
		resp.LogInfo()
		apiResponse = resp
		if err != nil {
			return mariadbv3.BackupLocationReadList{}, apiResponse, err
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
