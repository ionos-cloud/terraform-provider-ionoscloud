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
	diagutil "github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils/diags"
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
				Description: "The location of the resource. This field should be used only if you are also using a file configuration and should not be configured otherwise.",
			},
		},
		Timeouts: &resourceDefaultTimeouts,
	}
}

func resourceDbaasPgSqlDatabaseCreate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client, err := meta.(bundleclient.SdkBundle).NewPsqlClient(ctx, d.Get("location").(string))
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

	database, apiResponse, err := client.CreateDatabase(ctx, clusterId, request)
	if err != nil {
		return diagutil.ToDiags(d, fmt.Errorf("an error occurred while creating the PgSql database named: %s inside the cluster with ID: %s, error: %w", name, clusterId, err), &diagutil.ErrorContext{StatusCode: apiResponse.SafeStatusCode()})
	}
	return diagutil.ToDiags(d, dbaas.SetDatabasePgSqlData(d, &database), nil)
}

func resourceDbaasPgSqlDatabaseRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client, err := meta.(bundleclient.SdkBundle).NewPsqlClient(ctx, d.Get("location").(string))
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
		return diagutil.ToDiags(d, fmt.Errorf("an error occurred while fetching the PgSql database: %w", err), &diagutil.ErrorContext{StatusCode: apiResponse.SafeStatusCode()})
	}
	if err := dbaas.SetDatabasePgSqlData(d, &database); err != nil {
		return diagutil.ToDiags(d, err, nil)
	}
	return nil
}

func resourceDbaasPgSqlDatabaseDelete(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client, err := meta.(bundleclient.SdkBundle).NewPsqlClient(ctx, d.Get("location").(string))
	if err != nil {
		return diag.FromErr(err)
	}

	clusterId := d.Get("cluster_id").(string)
	name := d.Get("name").(string)
	apiResponse, err := client.DeleteDatabase(ctx, clusterId, name)
	if err != nil {
		return diagutil.ToDiags(d, err, &diagutil.ErrorContext{StatusCode: apiResponse.SafeStatusCode()})
	}
	return nil
}

func resourceDbaasPgSqlDatabaseImporter(ctx context.Context, d *schema.ResourceData, meta any) ([]*schema.ResourceData, error) {
	importID := d.Id()
	location, parts := splitImportID(importID, "/")
	if len(parts) != 2 {
		return nil, diagutil.ToError(d, fmt.Errorf("invalid import identifier: expected one of <location>:<cluster-id>/<database-name> or <cluster-id>/<database-name>, got: %s", importID), nil)
	}

	if err := validateImportIDParts(parts); err != nil {
		return nil, diagutil.ToError(d, fmt.Errorf("failed validating import identifier %q: %w", importID, err), nil)
	}

	client, err := meta.(bundleclient.SdkBundle).NewPsqlClient(ctx, location)
	if err != nil {
		return nil, err
	}

	clusterId := parts[0]
	name := parts[1]

	database, apiResponse, err := client.FindDatabaseByName(ctx, clusterId, name)
	if err != nil {
		if apiResponse.HttpNotFound() {
			d.SetId("")
			return nil, diagutil.ToError(d, fmt.Errorf("unable to find PgSql database: %s, cluster ID: %s", name, clusterId), &diagutil.ErrorContext{StatusCode: apiResponse.SafeStatusCode()})
		}
		return nil, diagutil.ToError(d, fmt.Errorf("error occurred while fetching PgSql database: %s, cluster ID: %s, error: %w", name, clusterId, err), &diagutil.ErrorContext{StatusCode: apiResponse.SafeStatusCode()})
	}
	if err := dbaas.SetDatabasePgSqlData(d, &database); err != nil {
		return nil, diagutil.ToError(d, err, nil)
	}
	if err := d.Set("cluster_id", clusterId); err != nil {
		return nil, utils.GenerateSetError("PgSQL database", "cluster_id", err)
	}
	if err := d.Set("location", location); err != nil {
		return nil, utils.GenerateSetError("PgSQL database", "location", err)
	}
	return []*schema.ResourceData{d}, nil
}
