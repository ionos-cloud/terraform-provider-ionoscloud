package ionoscloud

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/bundleclient"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils"
	diagutil "github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils/diags"
)

func dataSourceDbaasPgSqlDatabases() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDbaasPgSqlReadDatabases,
		Schema: map[string]*schema.Schema{
			"cluster_id": {
				Type:             schema.TypeString,
				Required:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.IsUUID),
			},
			"owner": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Filter for this data source, using this you can retrieve all databases that belong to a specific user.",
			},
			"location": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The location of the resource. This field should be used only if you are also using a file configuration and should not be configured otherwise.",
			},
			"databases": {
				Type:        schema.TypeList,
				Description: "The list of databases",
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Computed: true,
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
				},
			},
		},
		Timeouts: &resourceDefaultTimeouts,
	}
}

func dataSourceDbaasPgSqlReadDatabases(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client, err := meta.(bundleclient.SdkBundle).NewPsqlClient(ctx, d.Get("location").(string))
	if err != nil {
		return diag.FromErr(err)
	}
	clusterId := d.Get("cluster_id").(string)
	owner, ownerOk := d.GetOk("owner")
	resourceName := "PgSQL databases"

	retrievedDatabases, apiResponse, err := client.GetDatabases(ctx, clusterId)
	if err != nil {
		return diagutil.ToDiags(d, fmt.Errorf("an error occurred while fetching PgSql databases for the cluster with ID: %s, error: %w", clusterId, err), &diagutil.ErrorContext{StatusCode: apiResponse.SafeStatusCode()})
	}
	if retrievedDatabases.Items == nil {
		return diagutil.ToDiags(d, fmt.Errorf("expected a list of PgSql databases, but received 'nil' instead, cluster ID: %s", clusterId), nil)
	}
	var databases []any
	for _, retrievedDatabase := range retrievedDatabases.Items {
		database := make(map[string]any)
		utils.SetPropWithNilCheck(database, "name", retrievedDatabase.Properties.Name)
		utils.SetPropWithNilCheck(database, "owner", retrievedDatabase.Properties.Owner)
		utils.SetPropWithNilCheck(database, "id", retrievedDatabase.Id)
		// Filter using the owner
		if ownerOk {
			if retrievedDatabase.Properties.Owner == owner.(string) {
				databases = append(databases, database)
			}
		} else {
			databases = append(databases, database)
		}
	}
	if d.Id() == "" {
		d.SetId(clusterId)
	}
	if err := d.Set("databases", databases); err != nil {
		return diagutil.ToDiags(d, utils.GenerateSetError(resourceName, "databases", err), nil)
	}
	return nil
}
