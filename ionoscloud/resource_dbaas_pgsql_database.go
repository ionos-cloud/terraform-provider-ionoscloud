package ionoscloud

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	pgsql "github.com/ionos-cloud/sdk-go-dbaas-postgres"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/dbaas"
)

func resourceDbaasPgSqlDatabase() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDbaasPgSqlDatabaseCreate,
		UpdateContext: resourceDbaasPgSqlDatabaseUpdate,
		ReadContext:   resourceDbaasPgSqlDatabaseRead,
		DeleteContext: resourceDbaasPgSqlDatabaseDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceDbaasPgSqlDatabaseImporter,
		},
		Schema: map[string]*schema.Schema{
			"cluster_id": {
				Type:     schema.TypeString,
				Required: true,
				// TODO -- not sure if ForceNew makes sense here since there is no way to modify
				// the database, so the update method will be empty
				ForceNew:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.IsUUID),
			},
			"name": {
				Type:        schema.TypeString,
				Description: "The databasename of a given database.",
				Required:    true,
			},
			"owner": {
				Type:        schema.TypeString,
				Description: "The name of the role owning a given database.",
				Required:    true,
			},
		},
		Timeouts: &resourceDefaultTimeouts,
	}
}

func resourceDbaasPgSqlDatabaseCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(SdkBundle).PsqlClient
	clusterId := d.Get("cluster_id").(string)
	name := d.Get("name").(string)
	owner := d.Get("owner").(string)
	request := pgsql.Database{
		Properties: &pgsql.DatabaseProperties{},
	}
	request.Properties.Name = &name
	request.Properties.Owner = &owner

	database, _, err := client.CreateDatabase(ctx, clusterId, request)
	if err != nil {
		return diag.FromErr(fmt.Errorf("an error occured while creating the PgSql database named: %s inside the cluster with ID: %s, error: %w", name, clusterId, err))
	}
	d.SetId(*database.Id)
	return resourceDbaasPgSqlDatabaseRead(ctx, d, meta)
}

func resourceDbaasPgSqlDatabaseUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}

func resourceDbaasPgSqlDatabaseRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(SdkBundle).PsqlClient
	clusterId := d.Get("cluster_id").(string)
	name := d.Get("name").(string)

	database, apiResponse, err := client.FindDatabaseByName(ctx, clusterId, name)
	if err != nil {
		if apiResponse.HttpNotFound() {
			d.SetId("")
			return nil
		}
		return diag.FromErr(fmt.Errorf("an error occured while fetching the PgSql database with ID: %s, error: %w", d.Id(), err))
	}
	if err := dbaas.SetDatabasePgSqlData(d, &database); err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func resourceDbaasPgSqlDatabaseDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(SdkBundle).PsqlClient

	clusterId := d.Get("cluster_id").(string)
	name := d.Get("name").(string)
	_, err := client.DeleteDatabase(ctx, clusterId, name)
	if err != nil {
		diags := diag.FromErr(err)
		return diags
	}
	return nil
}

// TODO -- Finish this
func resourceDbaasPgSqlDatabaseImporter(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	return nil, nil
}
