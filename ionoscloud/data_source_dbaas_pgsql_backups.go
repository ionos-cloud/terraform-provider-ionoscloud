package ionoscloud

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	dbaas "github.com/ionos-cloud/sdk-go-autoscaling"
)

func dataSourceDbaasPgSqlBackups() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDbaasPgSqlReadBackups,
		Schema: map[string]*schema.Schema{
			"cluster_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"cluster_backups": {
				Type:        schema.TypeList,
				Description: "list of backups",
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Description: "The unique ID of the resource.",
							Computed:    true,
						},
						"cluster_id": {
							Type:        schema.TypeString,
							Description: "The unique ID of the cluster",
							Computed:    true,
						},
						"display_name": {
							Type:        schema.TypeString,
							Description: "The friendly name of your cluster.",
							Computed:    true,
						},
						"type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"metadata": {
							Type:        schema.TypeList,
							Description: "Metadata of the resource",
							Computed:    true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"created_date": {
										Type:        schema.TypeString,
										Description: "The ISO 8601 creation timestamp.",
										Computed:    true,
									},
									"created_by": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"created_by_user_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"last_modified_date": {
										Type:        schema.TypeString,
										Description: "The ISO 8601 modified timestamp.",
										Computed:    true,
									},
									"last_modified_by": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"last_modified_by_user_id": {
										Type:     schema.TypeString,
										Computed: true,
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

func dataSourceDbaasPgSqlReadBackups(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(SdkBundle).DbaasClient

	id, idOk := d.GetOk("cluster_id")

	if !idOk {
		diags := diag.FromErr(fmt.Errorf("cluster_id has to be provided in order to search for backups"))
		return diags
	}

	/* search by ID */
	clusterBackups, _, err := client.BackupsApi.ClusterBackupsGet(ctx, id.(string)).Execute()
	if err != nil {
		diags := diag.FromErr(fmt.Errorf("an error occurred while fetching backup for cluster with ID %s: %s", id.(string), err))
		return diags
	}

	if diags := setPgSqlClusterBackupData(d, &clusterBackups); diags != nil {
		return diags
	}

	return nil
}

func setPgSqlClusterBackupData(d *schema.ResourceData, clusterBackups *dbaas.ClusterBackupList) diag.Diagnostics {

	resourceId := uuid.New()
	d.SetId(resourceId.String())

	if clusterBackups.Data != nil {
		var backups []interface{}
		for _, backup := range *clusterBackups.Data {

			backupEntry := make(map[string]interface{})
			if backup.Id != nil {
				backupEntry["id"] = *backup.Id
			}

			if backup.ClusterId != nil {
				backupEntry["cluster_id"] = *backup.ClusterId
			}

			if backup.DisplayName != nil {
				backupEntry["display_name"] = *backup.DisplayName
			}

			if backup.Type != nil {
				backupEntry["type"] = *backup.Type
			}

			if backup.Metadata != nil {
				var metadata []interface{}

				metadataEntry := make(map[string]interface{})

				if backup.Metadata.CreatedDate != nil {
					metadataEntry["created_date"] = *backup.Metadata.CreatedDate
				}

				if backup.Metadata.CreatedBy != nil {
					metadataEntry["created_by"] = *backup.Metadata.CreatedDate
				}

				if backup.Metadata.CreatedByUserId != nil {
					metadataEntry["created_by_user_id"] = *backup.Metadata.CreatedDate
				}

				if backup.Metadata.LastModifiedDate != nil {
					metadataEntry["last_modified_date"] = *backup.Metadata.CreatedDate
				}

				if backup.Metadata.LastModifiedBy != nil {
					metadataEntry["last_modified_by"] = *backup.Metadata.CreatedDate
				}

				if backup.Metadata.LastModifiedByUserId != nil {
					metadataEntry["last_modified_by_user_id"] = *backup.Metadata.CreatedDate
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
