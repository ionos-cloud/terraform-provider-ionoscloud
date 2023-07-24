package ionoscloud

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils"
	"log"
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

func dataSourceDbaasPgSqlReadDatabases(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(SdkBundle).PsqlClient
	clusterId := d.Get("cluster_id").(string)
	owner, ownerOk := d.GetOk("owner")
	resourceName := "PgSQL databases"

	retrievedDatabases, _, err := client.GetDatabases(ctx, clusterId)
	if err != nil {
		return diag.FromErr(fmt.Errorf("an error occured while fetching PgSql databases for the cluster with ID: %s, error: %w", clusterId, err))
	}
	if retrievedDatabases.Items == nil {
		return diag.FromErr(fmt.Errorf("expected a list of PgSql databases, but received 'nil' instead, cluster ID: %s", clusterId))
	}
	var databases []interface{}
	for _, retrievedDatabase := range *retrievedDatabases.Items {
		if retrievedDatabase.Properties == nil {
			log.Printf("[INFO] 'nil' values in the response for the database with ID: %s, cluster ID: %s, skipping this database since there is not enough information to set in the state", *retrievedDatabase.Id, clusterId)
			continue
		}
		database := make(map[string]interface{})
		utils.SetPropWithNilCheck(database, "name", retrievedDatabase.Properties.Name)
		utils.SetPropWithNilCheck(database, "owner", retrievedDatabase.Properties.Owner)
		utils.SetPropWithNilCheck(database, "id", retrievedDatabase.Id)
		// Filter using the owner
		if ownerOk {
			if *retrievedDatabase.Properties.Owner == owner.(string) {
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
		return diag.FromErr(utils.GenerateSetError(resourceName, "databases", err))
	}
	return nil
}
