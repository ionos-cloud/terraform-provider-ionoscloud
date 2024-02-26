package ionoscloud

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	mariadbSDK "github.com/ionos-cloud/sdk-go-dbaas-maria"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services"
	mariaDBService "github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/dbaas/mariadb"
)

func dataSourceDBaaSMariaDBBackups() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDBaaSMariaDBReadBackups,
		Schema: map[string]*schema.Schema{
			"cluster_id": {
				Type:        schema.TypeString,
				Description: "The unique ID of the cluster that was backed up",
				Optional:    true,
			},
			"backup_id": {
				Type:        schema.TypeString,
				Description: "The unique ID of the backup",
				Optional:    true,
			},
			"cluster_backups": {
				Type:        schema.TypeList,
				Description: "The list of backups for the specified cluster",
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Description: "The unique ID of the backup",
							Computed:    true,
						},
						"size": {
							Type:        schema.TypeInt,
							Description: "The size of the backup in Mebibytes (MiB). This is the size of the binary backup file that was stored",
							Computed:    true,
						},
						"created": {
							Type:        schema.TypeString,
							Description: "The ISO 8601 creation timestamp",
							Computed:    true,
						},
					},
				},
			},
		},
		Timeouts: &resourceDefaultTimeouts,
	}
}

func dataSourceDBaaSMariaDBReadBackups(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(services.SdkBundle).MariaDBClient

	clusterIdIntf, clusterIdOk := d.GetOk("cluster_id")
	clusterId := clusterIdIntf.(string)
	backupIdIntf, backupIdOk := d.GetOk("backup_id")
	backupId := backupIdIntf.(string)

	if !clusterIdOk && !backupIdOk {
		return diag.FromErr(fmt.Errorf("please provide either the 'cluster_id' or 'backup_id'"))
	}
	if clusterIdOk && backupIdOk {
		return diag.FromErr(fmt.Errorf("'cluster_id' and 'backup_id' cannot be specified at the same time"))
	}

	var backups mariadbSDK.BackupResponse
	var err error
	if clusterIdOk {
		backups, _, err = client.GetClusterBackups(ctx, clusterId)
		if err != nil {
			return diag.FromErr(fmt.Errorf("an error occurred while fetching backups for cluster with ID %s: %w", clusterId, err))
		}
	} else if backupIdOk {
		backups, _, err = client.FindBackupById(ctx, backupId)
		if err != nil {
			return diag.FromErr(fmt.Errorf("an error occurred while fetching backup with ID %s: %w", backupId, err))
		}
	}

	if backups.Properties == nil || backups.Properties.Items == nil {
		return diag.FromErr(fmt.Errorf("expected valid properties in the API response for cluster backups, but received 'nil' instead"))
	}

	if diags := mariaDBService.SetMariaDBClusterBackupsData(d, &backups); diags != nil {
		return diags
	}

	return nil
}
