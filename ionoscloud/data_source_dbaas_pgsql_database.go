package ionoscloud

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/bundleclient"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/dbaas"
	diagutil "github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils/diags"
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
			"location": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The location of the resource. This field should be used only if you are also using a file configuration and should not be configured otherwise.",
			},
		},
		Timeouts: &resourceDefaultTimeouts,
	}
}

func dataSourceDbaasPgSqlReadDatabase(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client, err := meta.(bundleclient.SdkBundle).NewPsqlClient(ctx, d.Get("location").(string))
	if err != nil {
		return diag.FromErr(err)
	}
	clusterID := d.Get("cluster_id").(string)
	name := d.Get("name").(string)

	database, apiResponse, err := client.FindDatabaseByName(ctx, clusterID, name)
	if err != nil {
		if apiResponse.HttpNotFound() {
			return bundleclient.ToDiags(meta, d, fmt.Errorf("no PgSql database found with the specified name: %s and cluster ID: %s", name, clusterID), &diagutil.ErrorContext{StatusCode: apiResponse.SafeStatusCode()})
		}
		return bundleclient.ToDiags(meta, d, fmt.Errorf("an error occurred while fetching the PgSql database: %s, cluster ID: %s, err: %w", name, clusterID, err), &diagutil.ErrorContext{StatusCode: apiResponse.SafeStatusCode()})
	}
	if err := dbaas.SetDatabasePgSqlData(d, &database); err != nil {
		return bundleclient.ToDiags(meta, d, err, nil)
	}
	return nil
}
