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

func resourceDbaasPgSqlUser() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDbaasPgSqlUserCreate,
		UpdateContext: resourceDbaasPgSqlUserUpdate,
		ReadContext:   resourceDbaasPgSqlUserRead,
		DeleteContext: resourceDbaaSPgSqlUserDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceDbaasPgSqlUserImporter,
		},
		Schema: map[string]*schema.Schema{
			"cluster_id": {
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.IsUUID),
			},
			"username": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"password": {
				Type:      schema.TypeString,
				Required:  true,
				Sensitive: true,
			},
			"is_system_user": {
				Type:        schema.TypeBool,
				Description: "Describes whether this user is a system user or not. A system user cannot be updated or deleted.",
				Computed:    true,
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

func resourceDbaasPgSqlUserCreate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client, err := meta.(bundleclient.SdkBundle).NewPsqlClient(ctx, d.Get("location").(string))
	if err != nil {
		return diag.FromErr(err)
	}

	clusterId := d.Get("cluster_id").(string)
	username := d.Get("username").(string)
	password := d.Get("password").(string)

	request := pgsql.User{
		Properties: pgsql.UserProperties{},
	}
	request.Properties.Username = username
	request.Properties.Password = &password

	user, apiResponse, err := client.CreateUser(ctx, clusterId, request)
	if err != nil {
		return diagutil.ToDiags(d, fmt.Errorf("an error occurred while adding the user: %s to the PgSql cluster with ID: %s, error: %w", username, clusterId, err), &diagutil.ErrorContext{StatusCode: apiResponse.SafeStatusCode()})
	}
	d.SetId(user.Id)
	// Wait for the cluster to be ready again (when creating/updating the user, the cluster enters
	// 'BUSY' state).
	err = utils.WaitForResourceToBeReady(ctx, d, client.IsClusterReady)
	if err != nil {
		return diagutil.ToDiags(d, fmt.Errorf("error while waiting for PgSql cluster with ID: %s to be ready, error: %w", clusterId, err), &diagutil.ErrorContext{Timeout: d.Timeout(schema.TimeoutCreate).String()})
	}
	return diagutil.ToDiags(d, dbaas.SetUserPgSqlData(d, &user), nil)
}

func resourceDbaasPgSqlUserUpdate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client, err := meta.(bundleclient.SdkBundle).NewPsqlClient(ctx, d.Get("location").(string))
	if err != nil {
		return diag.FromErr(err)
	}

	request := pgsql.UsersPatchRequest{
		Properties: *pgsql.NewPatchUserProperties(),
	}

	clusterId := d.Get("cluster_id").(string)
	username := d.Get("username").(string)

	if d.HasChange("password") {
		_, newPassword := d.GetChange("password")
		password := newPassword.(string)
		request.Properties.Password = &password
	}

	user, apiResponse, err := client.UpdateUser(ctx, clusterId, username, request)
	if err != nil {
		return diagutil.ToDiags(d, fmt.Errorf("an error occurred while updating a PgSql user, username: %s, cluster ID: %s, error: %w", username, clusterId, err), &diagutil.ErrorContext{StatusCode: apiResponse.SafeStatusCode()})
	}
	// Wait for the cluster to be ready again (when creating/updating the user, the cluster enters
	// 'BUSY' state).
	err = utils.WaitForResourceToBeReady(ctx, d, client.IsClusterReady)
	if err != nil {
		return diagutil.ToDiags(d, fmt.Errorf("error while waiting for PgSql cluster with ID: %s to be ready after user: %s update, error: %w", clusterId, username, err), &diagutil.ErrorContext{Timeout: d.Timeout(schema.TimeoutUpdate).String()})
	}
	return diagutil.ToDiags(d, dbaas.SetUserPgSqlData(d, &user), nil)
}

func resourceDbaasPgSqlUserRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client, err := meta.(bundleclient.SdkBundle).NewPsqlClient(ctx, d.Get("location").(string))
	if err != nil {
		return diag.FromErr(err)
	}
	clusterId := d.Get("cluster_id").(string)
	username := d.Get("username").(string)

	user, apiResponse, err := client.FindUserByUsername(ctx, clusterId, username)
	if err != nil {
		if apiResponse.HttpNotFound() {
			d.SetId("")
			return nil
		}
		return diagutil.ToDiags(d, fmt.Errorf("an error occurred while fetching the PgSql user: %w", err), &diagutil.ErrorContext{StatusCode: apiResponse.SafeStatusCode()})
	}

	if err := dbaas.SetUserPgSqlData(d, &user); err != nil {
		return diagutil.ToDiags(d, err, nil)
	}

	return nil
}

func resourceDbaaSPgSqlUserDelete(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client, err := meta.(bundleclient.SdkBundle).NewPsqlClient(ctx, d.Get("location").(string))
	if err != nil {
		return diag.FromErr(err)
	}

	clusterId := d.Get("cluster_id").(string)
	username := d.Get("username").(string)
	apiResponse, err := client.DeleteUser(ctx, clusterId, username)
	if err != nil {
		return diagutil.ToDiags(d, err, &diagutil.ErrorContext{StatusCode: apiResponse.SafeStatusCode()})
	}
	err = utils.WaitForResourceToBeDeleted(ctx, d, client.IsUserDeleted)
	if err != nil {
		return diagutil.ToDiags(d, fmt.Errorf("error while waiting for the PgSql username: %s to be deleted, error: %w", username, err), &diagutil.ErrorContext{Timeout: d.Timeout(schema.TimeoutDelete).String()})
	}
	return nil
}

func resourceDbaasPgSqlUserImporter(ctx context.Context, d *schema.ResourceData, meta any) ([]*schema.ResourceData, error) {
	importID := d.Id()
	location, parts := splitImportID(importID, "/")
	if len(parts) != 2 {
		return nil, diagutil.ToError(d, fmt.Errorf("invalid import identifier: expected one of <location>:<cluster-id>/<username> or <cluster-id>/<username>, got: %s", importID), nil)
	}

	if err := validateImportIDParts(parts); err != nil {
		return nil, diagutil.ToError(d, fmt.Errorf("failed validating import identifier %q: %w", importID, err), nil)
	}

	client, err := meta.(bundleclient.SdkBundle).NewPsqlClient(ctx, location)
	if err != nil {
		return nil, err
	}

	clusterId := parts[0]
	username := parts[1]

	user, apiResponse, err := client.FindUserByUsername(ctx, clusterId, username)
	if err != nil {
		if apiResponse.HttpNotFound() {
			d.SetId("")
			return nil, diagutil.ToError(d, fmt.Errorf("unable to find PgSql username: %s, cluster ID: %s", username, clusterId), &diagutil.ErrorContext{StatusCode: apiResponse.SafeStatusCode()})
		}
		return nil, diagutil.ToError(d, fmt.Errorf("error occurred while fetching PgSql username: %s, cluster ID: %s, error: %w", username, clusterId, err), &diagutil.ErrorContext{StatusCode: apiResponse.SafeStatusCode()})
	}
	if err := dbaas.SetUserPgSqlData(d, &user); err != nil {
		return nil, diagutil.ToError(d, err, nil)
	}
	if err := d.Set("cluster_id", clusterId); err != nil {
		return nil, utils.GenerateSetError("PgSQL user", "cluster_id", err)
	}
	if err := d.Set("location", location); err != nil {
		return nil, utils.GenerateSetError("PgSQL user", "location", err)
	}
	return []*schema.ResourceData{d}, nil
}
