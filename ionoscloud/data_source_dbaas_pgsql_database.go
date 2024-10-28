package ionoscloud

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/dbaas"
)

func dataSourceDbaasPgSqlDatabase() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDbaasPgSqlReadDatabase,
		Schema: map[string]*schema.Schema{
			"cluster_id": {
				Type:             schema.TypeString,
				Required:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.IsUUID),
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"owner": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
		Timeouts: &resourceDefaultTimeouts,
	}
}

func dataSourceDbaasPgSqlReadDatabase(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(services.SdkBundle).PsqlClient
	clusterId := d.Get("cluster_id").(string)
	name := d.Get("name").(string)

	database, apiResponse, err := client.FindDatabaseByName(ctx, clusterId, name)
	if err != nil {
		if apiResponse.HttpNotFound() {
			return diag.FromErr(fmt.Errorf("no PgSql database found with the specified name: %s and cluster ID: %s", name, clusterId))
		}
		return diag.FromErr(fmt.Errorf("an error occurred while fetching the PgSql database: %s, cluster ID: %s, err: %w", name, clusterId, err))
	}
	if err := dbaas.SetDatabasePgSqlData(d, &database); err != nil {
		return diag.FromErr(err)
	}
	return nil
}
