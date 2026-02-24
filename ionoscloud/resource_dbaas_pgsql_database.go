package ionoscloud

import (
	"context"
	"fmt"

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
			"location": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "The location of the cluster this database belongs to.",
			},
		},
		Timeouts: &resourceDefaultTimeouts,
	}
}

func resourceDbaasPgSqlDatabaseCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, err := meta.(bundleclient.SdkBundle).NewPsqlClient(d.Get("location").(string))
	if err != nil {
		return diag.FromErr(err)
	}
	clusterId := d.Get("cluster_id").(string)
	name := d.Get("name").(string)
	owner := d.Get("owner").(string)
	request := pgsql.Database{
		Properties: pgsql.DatabaseProperties{},
	}
	request.Properties.Name = name
	request.Properties.Owner = owner

	database, _, err := client.CreateDatabase(ctx, clusterId, request)
	if err != nil {
		return diag.FromErr(fmt.Errorf("an error occurred while creating the PgSql database named: %s inside the cluster with ID: %s, error: %w", name, clusterId, err))
	}
	return diag.FromErr(dbaas.SetDatabasePgSqlData(d, &database))
}

func resourceDbaasPgSqlDatabaseRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, err := meta.(bundleclient.SdkBundle).NewPsqlClient(d.Get("location").(string))
	if err != nil {
		return diag.FromErr(err)
	}
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
	client, err := meta.(bundleclient.SdkBundle).NewPsqlClient(d.Get("location").(string))
	if err != nil {
		return diag.FromErr(err)
	}

	clusterId := d.Get("cluster_id").(string)
	name := d.Get("name").(string)
	_, err = client.DeleteDatabase(ctx, clusterId, name)
	if err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func resourceDbaasPgSqlDatabaseImporter(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	importID := d.Id()
	location, parts := splitImportID(importID, "/")
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid import identifier: expected one of <location>:<cluster-id>/<database-name> or <cluster-id>/<database-name>, got: %s", importID)
	}

	if err := validateImportIDParts(parts); err != nil {
		return nil, fmt.Errorf("failed validating import identifier %q: %w", importID, err)
	}

	client, err := meta.(bundleclient.SdkBundle).NewPsqlClient(location)
	if err != nil {
		return nil, err
	}

	clusterId := parts[0]
	name := parts[1]

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
	if err := d.Set("location", location); err != nil {
		return nil, utils.GenerateSetError("PgSQL database", "location", err)
	}
	return []*schema.ResourceData{d}, nil
}
