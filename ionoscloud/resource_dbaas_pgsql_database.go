package ionoscloud

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	pgsql "github.com/ionos-cloud/sdk-go-dbaas-postgres"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/dbaas"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils"
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
				Type:             schema.TypeString,
				Required:         true,
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
	client := meta.(services.SdkBundle).PsqlClient
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
		return diag.FromErr(fmt.Errorf("an error occurred while creating the PgSql database named: %s inside the cluster with ID: %s, error: %w", name, clusterId, err))
	}
	return diag.FromErr(dbaas.SetDatabasePgSqlData(d, &database))
}

func resourceDbaasPgSqlDatabaseUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}

func resourceDbaasPgSqlDatabaseRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(services.SdkBundle).PsqlClient
	clusterId := d.Get("cluster_id").(string)
	name := d.Get("name").(string)

	database, apiResponse, err := client.FindDatabaseByName(ctx, clusterId, name)
	if err != nil {
		if apiResponse.HttpNotFound() {
			d.SetId("")
			return nil
		}
		return diag.FromErr(fmt.Errorf("an error occurred while fetching the PgSql database with ID: %s, error: %w", d.Id(), err))
	}
	if err := dbaas.SetDatabasePgSqlData(d, &database); err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func resourceDbaasPgSqlDatabaseDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(services.SdkBundle).PsqlClient

	clusterId := d.Get("cluster_id").(string)
	name := d.Get("name").(string)
	_, err := client.DeleteDatabase(ctx, clusterId, name)
	if err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func resourceDbaasPgSqlDatabaseImporter(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	parts := strings.Split(d.Id(), "/")
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid import format: %s, expecting the following format: {clusterID}/{databaseName}", d.Id())
	}
	clusterId := parts[0]
	name := parts[1]
	client := meta.(services.SdkBundle).PsqlClient
	database, apiResponse, err := client.FindDatabaseByName(ctx, clusterId, name)
	if err != nil {
		if apiResponse.HttpNotFound() {
			d.SetId("")
			return nil, fmt.Errorf("unable to find PgSql database: %s, cluster ID: %s", name, clusterId)
		}
		return nil, fmt.Errorf("error occurred while fetching PgSql database: %s, cluster ID: %s, error: %w", name, clusterId, err)
	}
	if err := dbaas.SetDatabasePgSqlData(d, &database); err != nil {
		return nil, err
	}
	if err := d.Set("cluster_id", clusterId); err != nil {
		return nil, utils.GenerateSetError("PgSQL database", "cluster_id", err)
	}
	return []*schema.ResourceData{d}, nil
}
