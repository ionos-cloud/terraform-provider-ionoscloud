package ionoscloud

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	psql "github.com/ionos-cloud/sdk-go-dbaas-postgres"
	dbaasService "github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/dbaas"
)

func dataSourceDbaasPgSqlVersions() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDbaasPgSqlReadVersions,
		Schema: map[string]*schema.Schema{
			"cluster_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"postgres_versions": {
				Type:        schema.TypeList,
				Description: "list of PostgreSQL versions",
				Computed:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
		Timeouts: &resourceDefaultTimeouts,
	}
}

func dataSourceDbaasPgSqlReadVersions(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(SdkBundle).PsqlClient

	id, idOk := d.GetOk("cluster_id")

	var postgresVersions psql.PostgresVersionList
	var err error

	if idOk {
		/* search by ID */
		postgresVersions, _, err = client.GetClusterVersions(ctx, id.(string))
		if err != nil {
			diags := diag.FromErr(fmt.Errorf("an error occurred while fetching postgres versions for cluster with ID %s: %w", id.(string), err))
			return diags
		}
	} else {
		postgresVersions, _, err = client.GetAllVersions(ctx)
		if err != nil {
			diags := diag.FromErr(fmt.Errorf("an error occurred while fetching postgres versions: %w", err))
			return diags
		}
	}

	dbaasService.SetPgSqlVersionsData(d, postgresVersions)

	return nil

}
