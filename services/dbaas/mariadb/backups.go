package mariadb

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	mariadb "github.com/ionos-cloud/sdk-go-dbaas-maria"
)

func (c *MariaDBClient) GetClusterBackups(ctx context.Context, clusterId string) (mariadb.BackupResponse, *mariadb.APIResponse, error) {
	backups, apiResponse, err := c.sdkClient.BackupsApi.ClusterBackupsGet(ctx, clusterId).Execute()
	apiResponse.LogInfo()
	return backups, apiResponse, err
}

func (c *MariaDBClient) FindBackupById(ctx context.Context, backupId string) (mariadb.BackupResponse, *mariadb.APIResponse, error) {
	backups, apiResponse, err := c.sdkClient.BackupsApi.BackupsFindById(ctx, backupId).Execute()
	apiResponse.LogInfo()
	return backups, apiResponse, err
}

func SetMariaDBClusterBackupsData(d *schema.ResourceData, clusterBackups *mariadb.BackupResponse) diag.Diagnostics {
	resourceId := uuid.New()
	d.SetId(resourceId.String())

	var backups []interface{}
	for _, backup := range *clusterBackups.Properties.Items {
		backupEntry := make(map[string]interface{})
		backupEntry["id"] = *clusterBackups.Id

		if backup.Size != nil {
			backupEntry["size"] = *backup.Size
		}
		if backup.Created != nil {
			backupEntry["created"] = (*backup.Created).Time.Format("2006-01-02T15:04:05Z")
		}
		backups = append(backups, backupEntry)
	}
	err := d.Set("cluster_backups", backups)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error while setting 'cluster_backups': %w", err))
	}
	return nil
}
