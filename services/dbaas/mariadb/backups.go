package mariadb

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	mariadb "github.com/ionos-cloud/sdk-go-bundle/products/dbaas/mariadb/v2"
	"github.com/ionos-cloud/sdk-go-bundle/shared"
	"github.com/ionos-cloud/sdk-go-bundle/shared/fileconfiguration"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils/loadedconfig"
)

// GetClusterBackups retrieves a list of backups for a given cluster ID and the location in which the cluster is created.
func (c *Client) GetClusterBackups(ctx context.Context, clusterID, location string) (mariadb.BackupList, *shared.APIResponse, error) {
	loadedconfig.SetClientOptionsFromConfig(c, fileconfiguration.Mariadb, location)
	backups, apiResponse, err := c.sdkClient.BackupsApi.ClusterBackupsGet(ctx, clusterID).Execute()
	apiResponse.LogInfo()
	return backups, apiResponse, err
}

// FindBackupByID retrieves a backup by its ID and the location in which the cluster is created.
func (c *Client) FindBackupByID(ctx context.Context, backupID, location string) (mariadb.BackupResponse, *shared.APIResponse, error) {
	loadedconfig.SetClientOptionsFromConfig(c, fileconfiguration.Mariadb, location)
	backups, apiResponse, err := c.sdkClient.BackupsApi.BackupsFindById(ctx, backupID).Execute()
	apiResponse.LogInfo()
	return backups, apiResponse, err
}

// SetMariaDBClusterBackupsData sets the data for the backups attribute in the MariaDB backup data source.
func SetMariaDBClusterBackupsData(d *schema.ResourceData, retrievedBackups []mariadb.BackupResponse) diag.Diagnostics {
	resourceId := uuid.New()
	d.SetId(resourceId.String())

	var backupsToBeSet []interface{}
	for _, retrievedBackup := range retrievedBackups {
		if retrievedBackup.Properties == nil {
			return diag.FromErr(fmt.Errorf("expected valid properties in the API response for backup, but received 'nil' instead, backup ID: %s", *retrievedBackup.Id))
		}
		backupEntry := make(map[string]interface{})
		if retrievedBackup.Properties.ClusterId != nil {
			backupEntry["cluster_id"] = *retrievedBackup.Properties.ClusterId
		}
		if retrievedBackup.Properties.EarliestRecoveryTargetTime != nil {
			backupEntry["earliest_recovery_target_time"] = (*retrievedBackup.Properties.EarliestRecoveryTargetTime).String()
		}
		if retrievedBackup.Properties.Size != nil {
			backupEntry["size"] = *retrievedBackup.Properties.Size
		}
		var baseBackupsToBeSet []interface{}
		for _, baseBackup := range retrievedBackup.Properties.BaseBackups {
			baseBackupEntry := make(map[string]interface{})
			if baseBackup.Size != nil {
				baseBackupEntry["size"] = *baseBackup.Size
			}
			if baseBackup.Created != nil {
				baseBackupEntry["created"] = (*baseBackup.Created).String()
			}
			baseBackupsToBeSet = append(baseBackupsToBeSet, baseBackupEntry)
		}
		backupEntry["base_backups"] = baseBackupsToBeSet
		backupsToBeSet = append(backupsToBeSet, backupEntry)
	}
	err := d.Set("backups", backupsToBeSet)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error while setting 'cluster_backups': %w", err))
	}
	return nil
}
