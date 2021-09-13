package ionoscloud

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	dbaas "github.com/ionos-cloud/sdk-go-autoscaling"
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
	client := meta.(SdkBundle).DbaasClient

	id, idOk := d.GetOk("cluster_id")

	var postgresVersions dbaas.PostgresVersionList
	var err error

	if idOk {
		/* search by ID */
		postgresVersions, _, err = client.ClustersApi.ClusterPostgresVersionsGet(ctx, id.(string)).Execute()
		if err != nil {
			diags := diag.FromErr(fmt.Errorf("an error occurred while fetching postgres versions for cluster with ID %s: %s", id.(string), err))
			return diags
		}
	} else {
		postgresVersions, _, err = client.ClustersApi.PostgresVersionsGet(ctx).Execute()
		if err != nil {
			diags := diag.FromErr(fmt.Errorf("an error occurred while fetching postgres versions: %s", err.Error()))
			return diags
		}
	}

	if postgresVersions.Data != nil {
		var versions []string
		for _, version := range *postgresVersions.Data {
			if version.Name != nil {
				versions = append(versions, *version.Name)
			}
		}
		err := d.Set("postgres_versions", versions)
		if err != nil {
			diags := diag.FromErr(fmt.Errorf("error while setting postgres_versions: %s", err))
			return diags
		}
	}

	resourceId := uuid.New()
	d.SetId(resourceId.String())

	return nil

}
