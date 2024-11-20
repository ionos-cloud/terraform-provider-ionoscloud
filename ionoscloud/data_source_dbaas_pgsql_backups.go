package ionoscloud

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services"
	dbaasService "github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/dbaas"
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
						"type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"size": {
							Type:        schema.TypeInt,
							Description: "Size of all base backups including the wal size in MB.",
							Computed:    true,
						},
						"location": {
							Type:        schema.TypeString,
							Description: "The Object Storage location where the backups will be stored.",
							Computed:    true,
						},
						"version": {
							Type:        schema.TypeString,
							Description: "The PostgreSQL version this backup was created from.",
							Computed:    true,
						},
						"is_active": {
							Type:        schema.TypeBool,
							Description: "Whether a cluster currently backs up data to this backup.",
							Computed:    true,
						},
						"earliest_recovery_target_time": {
							Type:        schema.TypeString,
							Description: "The oldest available timestamp to which you can restore.",
							Computed:    true,
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
	client := meta.(services.SdkBundle).PsqlClient

	id, idOk := d.GetOk("cluster_id")
	idStr := id.(string)
	if !idOk {
		diags := diag.FromErr(fmt.Errorf("cluster_id has to be provided in order to search for backups"))
		return diags
	}

	/* search by ID */
	clusterBackups, resp, err := client.GetClusterBackups(ctx, idStr)
	if resp != nil {
		log.Printf("operation %s", resp.Operation)
		if resp.Response != nil {
			log.Printf("[DEBUG] response status code : %d\n", resp.StatusCode)
		}
	}

	if err != nil {
		diags := diag.FromErr(fmt.Errorf("an error occurred while fetching backup for cluster with ID %s: %w", idStr, err))
		return diags
	}
	if len(*clusterBackups.Items) == 0 {
		diags := diag.FromErr(fmt.Errorf("could not find backups for cluster with ID %s: %w", idStr, err))
		return diags
	}

	if diags := dbaasService.SetPgSqlClusterBackupData(d, &clusterBackups); diags != nil {
		return diags
	}

	return nil
}
