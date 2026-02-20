package ionoscloud

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	pgsql "github.com/ionos-cloud/sdk-go-bundle/products/dbaas/psql/v2"

	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/bundleclient"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/dbaas"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils"
)

func resourceDbaasPgSqlDatabase() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDbaasPgSqlDatabaseCreate,
		ReadContext:   resourceDbaasPgSqlDatabaseRead,
		DeleteContext: resourceDbaasPgSqlDatabaseDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceDbaasPgSqlDatabaseImporter,
		},
		Schema: map[string]*schema.Schema{
			"cluster_id": {
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.IsUUID),
			},
			"name": {
				Type:        schema.TypeString,
				Description: "The databasename of a given database.",
				Required:    true,
				ForceNew:    true,
			},
			"owner": {
				Type:        schema.TypeString,
				Description: "The name of the role owning a given database.",
				Required:    true,
				ForceNew:    true,
			},
		},
		Timeouts: &resourceDefaultTimeouts,
	}
}

func resourceDbaasPgSqlDatabaseCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(bundleclient.SdkBundle).PsqlClient
	clusterId := d.Get("cluster_id").(string)
	name := d.Get("name").(string)
	owner := d.Get("owner").(string)
	request := pgsql.Database{
		Properties: pgsql.DatabaseProperties{},
	}
	request.Properties.Name = name
	request.Properties.Owner = owner

	database, apiResponse, err := client.CreateDatabase(ctx, clusterId, request)
	if err != nil {
		return utils.ToDiags(d, fmt.Sprintf("an error occurred while creating the PgSql database named: %s inside the cluster with ID: %s, error: %s", name, clusterId, err), &utils.DiagsOpts{StatusCode: apiResponse.StatusCode})
	}
	return utils.ToDiags(d, dbaas.SetDatabasePgSqlData(d, &database).Error(), nil)
}

func resourceDbaasPgSqlDatabaseRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(bundleclient.SdkBundle).PsqlClient
	clusterId := d.Get("cluster_id").(string)
	name := d.Get("name").(string)

	database, apiResponse, err := client.FindDatabaseByName(ctx, clusterId, name)
	if err != nil {
		if apiResponse.HttpNotFound() {
			d.SetId("")
			return nil
		}
		return utils.ToDiags(d, fmt.Sprintf("an error occurred while fetching the PgSql database: %s", err), &utils.DiagsOpts{StatusCode: apiResponse.StatusCode})
	}
	if err := dbaas.SetDatabasePgSqlData(d, &database); err != nil {
		return utils.ToDiags(d, err.Error(), nil)
	}
	return nil
}

func resourceDbaasPgSqlDatabaseDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(bundleclient.SdkBundle).PsqlClient

	clusterId := d.Get("cluster_id").(string)
	name := d.Get("name").(string)
	apiResponse, err := client.DeleteDatabase(ctx, clusterId, name)
	if err != nil {
		return utils.ToDiags(d, err.Error(), &utils.DiagsOpts{StatusCode: apiResponse.StatusCode})
	}
	return nil
}

func resourceDbaasPgSqlDatabaseImporter(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	parts := strings.Split(d.Id(), "/")
	if len(parts) != 2 {
		return nil, utils.ToError(d, "invalid import format:, expecting the following format: {clusterID}/{databaseName}", nil)
	}
	clusterId := parts[0]
	name := parts[1]
	client := meta.(bundleclient.SdkBundle).PsqlClient
	database, apiResponse, err := client.FindDatabaseByName(ctx, clusterId, name)
	if err != nil {
		if apiResponse.HttpNotFound() {
			d.SetId("")
			return nil, utils.ToError(d, fmt.Sprintf("unable to find PgSql database: %s, cluster ID: %s", name, clusterId), &utils.DiagsOpts{StatusCode: apiResponse.StatusCode})
		}
		return nil, utils.ToError(d, fmt.Sprintf("error occurred while fetching PgSql database: %s, cluster ID: %s, error: %s", name, clusterId, err), &utils.DiagsOpts{StatusCode: apiResponse.StatusCode})
	}
	if err := dbaas.SetDatabasePgSqlData(d, &database); err != nil {
		return nil, utils.ToError(d, err.Error(), nil)
	}
	if err := d.Set("cluster_id", clusterId); err != nil {
		return nil, utils.GenerateSetError("PgSQL database", "cluster_id", err)
	}
	return []*schema.ResourceData{d}, nil
}
