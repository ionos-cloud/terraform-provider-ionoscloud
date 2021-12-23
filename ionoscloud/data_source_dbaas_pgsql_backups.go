package ionoscloud

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	dbaasService "github.com/ionos-cloud/terraform-provider-ionoscloud/services/dbaas"
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
	client := meta.(SdkBundle).DbaasClient

	id, idOk := d.GetOk("cluster_id")

	if !idOk {
		diags := diag.FromErr(fmt.Errorf("cluster_id has to be provided in order to search for backups"))
		return diags
	}

	/* search by ID */
	clusterBackups, _, err := client.GetClusterBackups(ctx, id.(string))

	if err != nil {
		diags := diag.FromErr(fmt.Errorf("an error occurred while fetching backup for cluster with ID %s: %s", id.(string), err))
		return diags
	}

	if diags := dbaasService.SetPgSqlClusterBackupData(d, &clusterBackups); diags != nil {
		return diags
	}

	return nil
}
