package ionoscloud

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	psql "github.com/ionos-cloud/sdk-go-bundle/products/dbaas/psql/v2"
	"github.com/ionos-cloud/sdk-go-bundle/shared"

	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/bundleclient"
	dbaasService "github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/dbaas"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils"
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
	client := meta.(bundleclient.SdkBundle).PsqlClient

	id, idOk := d.GetOk("cluster_id")

	var postgresVersions psql.PostgresVersionList
	var apiResponse *shared.APIResponse
	var err error

	if idOk {
		/* search by ID */
		postgresVersions, apiResponse, err = client.GetClusterVersions(ctx, id.(string))
		if err != nil {
			return utils.ToDiags(d, fmt.Sprintf("an error occurred while fetching postgres versions for cluster with ID %s: %s", id.(string), err), &utils.DiagsOpts{StatusCode: apiResponse.StatusCode})
		}
	} else {
		postgresVersions, apiResponse, err = client.GetAllVersions(ctx)
		if err != nil {
			return utils.ToDiags(d, fmt.Sprintf("an error occurred while fetching postgres versions: %s", err), &utils.DiagsOpts{StatusCode: apiResponse.StatusCode})
		}
	}

	dbaasService.SetPgSqlVersionsData(d, postgresVersions)

	return nil

}
