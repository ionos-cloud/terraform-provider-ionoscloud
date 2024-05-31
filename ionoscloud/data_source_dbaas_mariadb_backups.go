package ionoscloud

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	mariadbSDK "github.com/ionos-cloud/sdk-go-dbaas-mariadb"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/dbaas/mariadb"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils/constant"
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
			"location": {
				Type:             schema.TypeString,
				Description:      "The cluster location",
				Optional:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice(constant.MariaDBClusterLocations, false)),
			},
			"backups": {
				Type:        schema.TypeList,
				Description: "The list of backups",
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"cluster_id": {
							Type:        schema.TypeString,
							Description: "The unique ID of the cluster that was backed up",
							Computed:    true,
						},
						"earliest_recovery_target_time": {
							Type:        schema.TypeString,
							Description: "The oldest available timestamp to which you can restore",
							Computed:    true,
						},
						"size": {
							Type:        schema.TypeInt,
							Description: "Size of all base backups in Mebibytes (MiB). This is at least the sum of all base backup sizes",
							Computed:    true,
						},
						"base_backups": {
							Type:        schema.TypeList,
							Description: "The list of backups for the specified cluster",
							Computed:    true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
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

	location := d.Get("location").(string)

	var backups []mariadbSDK.BackupResponse
	var err error
	if clusterIdOk {
		var backupsResponse mariadbSDK.BackupList
		backupsResponse, _, err = client.GetClusterBackups(ctx, clusterId, location)
		if err != nil {
			return diag.FromErr(fmt.Errorf("an error occurred while fetching backups for cluster with ID %s: %w", clusterId, err))
		}
		if backupsResponse.Items == nil {
			return diag.FromErr(fmt.Errorf("expected valid properties in the API response for cluster backups, but received 'nil' instead, cluster ID: %s", clusterId))
		}
		backups = *backupsResponse.Items
	} else {
		var backup mariadbSDK.BackupResponse
		backup, _, err = client.FindBackupByID(ctx, backupId, location)
		if err != nil {
			return diag.FromErr(fmt.Errorf("an error occurred while fetching backup with ID %s: %w", backupId, err))
		}
		if backup.Properties == nil {
			return diag.FromErr(fmt.Errorf("expected valid properties in the API response for backup, but received 'nil' instead, backup ID: %s", backupId))
		}
		backups = append(backups, backup)
	}

	return mariadb.SetMariaDBClusterBackupsData(d, backups)
}
