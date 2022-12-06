package dbaas

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	dbaas "github.com/ionos-cloud/sdk-go-dbaas-postgres"
)

type BackupService interface {
	GetClusterBackups(ctx context.Context, clusterId string) (dbaas.ClusterBackupList, *dbaas.APIResponse, error)
	GetAllBackups(ctx context.Context) (dbaas.ClusterBackupList, *dbaas.APIResponse, error)
}

func (c *PsqlClient) GetClusterBackups(ctx context.Context, clusterId string) (dbaas.ClusterBackupList, *dbaas.APIResponse, error) {
	backups, apiResponse, err := c.sdkClient.BackupsApi.ClusterBackupsGet(ctx, clusterId).Execute()
	if apiResponse != nil {
		return backups, apiResponse, err

	}
	return backups, nil, err
}

func (c *PsqlClient) GetAllBackups(ctx context.Context) (dbaas.ClusterBackupList, *dbaas.APIResponse, error) {
	backups, apiResponse, err := c.sdkClient.BackupsApi.ClustersBackupsGet(ctx).Execute()
	if apiResponse != nil {
		return backups, apiResponse, err
	}
	return backups, nil, err
}

func SetPgSqlClusterBackupData(d *schema.ResourceData, clusterBackups *dbaas.ClusterBackupList) diag.Diagnostics {

	resourceId := uuid.New()
	d.SetId(resourceId.String())

	if clusterBackups.Items != nil {
		var backups []interface{}
		for _, backup := range *clusterBackups.Items {

			backupEntry := make(map[string]interface{})
			if backup.Id != nil {
				backupEntry["id"] = *backup.Id
			}

			if backup.Properties == nil {
				return diag.FromErr(fmt.Errorf("backup properties do not exist."))
			}

			if backup.Properties.ClusterId != nil {
				backupEntry["cluster_id"] = *backup.Properties.ClusterId
			}

			if backup.Properties.Size != nil {
				backupEntry["size"] = *backup.Properties.Size
			}

			if backup.Properties.Location != nil {
				backupEntry["location"] = *backup.Properties.Location
			}

			if backup.Properties.Version != nil {
				backupEntry["version"] = *backup.Properties.Version
			}

			if backup.Properties.IsActive != nil {
				backupEntry["is_active"] = *backup.Properties.IsActive
			}

			if backup.Properties.EarliestRecoveryTargetTime != nil {
				backupEntry["earliest_recovery_target_time"] = (*backup.Properties.EarliestRecoveryTargetTime).String()
			}

			if backup.Type != nil {
				backupEntry["type"] = *backup.Type
			}

			if backup.Metadata != nil {
				var metadata []interface{}

				metadataEntry := make(map[string]interface{})

				if backup.Metadata.CreatedDate != nil {
					metadataEntry["created_date"] = (*backup.Metadata.CreatedDate).Time.Format("2006-01-02T15:04:05Z")
				}

				metadata = append(metadata, metadataEntry)
				backupEntry["metadata"] = metadata
			}

			backups = append(backups, backupEntry)

		}
		err := d.Set("cluster_backups", backups)
		if err != nil {
			diags := diag.FromErr(fmt.Errorf("error while setting cluster_backups: %s", err))
			return diags
		}
	}
	return nil
}
