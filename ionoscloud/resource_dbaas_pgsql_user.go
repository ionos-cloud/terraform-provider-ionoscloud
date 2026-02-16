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
		},
		Timeouts: &resourceDefaultTimeouts,
	}
}

func resourceDbaasPgSqlUserCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(bundleclient.SdkBundle).PsqlClient

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
		return utils.ToDiags(d, fmt.Sprintf("an error occurred while adding the user: %s to the PgSql cluster with ID: %s, error: %s", username, clusterId, err), &utils.DiagsOpts{StatusCode: apiResponse.StatusCode})
	}
	d.SetId(user.Id)
	// Wait for the cluster to be ready again (when creating/updating the user, the cluster enters
	// 'BUSY' state).
	err = utils.WaitForResourceToBeReady(ctx, d, client.IsClusterReady)
	if err != nil {
		return utils.ToDiags(d, fmt.Sprintf("error while waiting for PgSql cluster with ID: %s to be ready, error: %s", clusterId, err), nil)
	}
	return utils.ToDiags(d, dbaas.SetUserPgSqlData(d, &user).Error(), nil)
}

func resourceDbaasPgSqlUserUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(bundleclient.SdkBundle).PsqlClient

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
		return utils.ToDiags(d, fmt.Sprintf("an error occurred while updating a PgSql user, username: %s, cluster ID: %s, error: %s", username, clusterId, err), &utils.DiagsOpts{StatusCode: apiResponse.StatusCode})
	}
	// Wait for the cluster to be ready again (when creating/updating the user, the cluster enters
	// 'BUSY' state).
	err = utils.WaitForResourceToBeReady(ctx, d, client.IsClusterReady)
	if err != nil {
		return utils.ToDiags(d, fmt.Sprintf("error while waiting for PgSql cluster with ID: %s to be ready after user: %s update, error: %s", clusterId, username, err), nil)
	}
	return utils.ToDiags(d, dbaas.SetUserPgSqlData(d, &user).Error(), nil)
}

func resourceDbaasPgSqlUserRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(bundleclient.SdkBundle).PsqlClient
	clusterId := d.Get("cluster_id").(string)
	username := d.Get("username").(string)

	user, apiResponse, err := client.FindUserByUsername(ctx, clusterId, username)
	if err != nil {
		if apiResponse.HttpNotFound() {
			d.SetId("")
			return nil
		}
		return utils.ToDiags(d, fmt.Sprintf("an error occurred while fetching the PgSql user: %s", err), &utils.DiagsOpts{StatusCode: apiResponse.StatusCode})
	}

	if err := dbaas.SetUserPgSqlData(d, &user); err != nil {
		return utils.ToDiags(d, err.Error(), nil)
	}

	return nil
}

func resourceDbaaSPgSqlUserDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(bundleclient.SdkBundle).PsqlClient

	clusterId := d.Get("cluster_id").(string)
	username := d.Get("username").(string)
	apiResponse, err := client.DeleteUser(ctx, clusterId, username)
	if err != nil {
		return utils.ToDiags(d, err.Error(), &utils.DiagsOpts{StatusCode: apiResponse.StatusCode})
	}
	err = utils.WaitForResourceToBeDeleted(ctx, d, client.IsUserDeleted)
	if err != nil {
		return utils.ToDiags(d, fmt.Sprintf("error while waiting for the PgSql username: %s to be deleted, error: %s", username, err), &utils.DiagsOpts{Timeout: schema.TimeoutDelete})
	}
	return nil
}

func resourceDbaasPgSqlUserImporter(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	parts := strings.Split(d.Id(), "/")
	if len(parts) != 2 {
		return nil, utils.ToError(d, "invalid import format:, expecting the following format: {clusterID}/{username}", nil)
	}
	clusterId := parts[0]
	username := parts[1]
	client := meta.(bundleclient.SdkBundle).PsqlClient
	user, apiResponse, err := client.FindUserByUsername(ctx, clusterId, username)
	if err != nil {
		if apiResponse.HttpNotFound() {
			d.SetId("")
			return nil, utils.ToError(d, fmt.Sprintf("unable to find PgSql username: %s, cluster ID: %s", username, clusterId), &utils.DiagsOpts{StatusCode: apiResponse.StatusCode})
		}
		return nil, utils.ToError(d, fmt.Sprintf("error occurred while fetching PgSql username: %s, cluster ID: %s, error: %s", username, clusterId, err), &utils.DiagsOpts{StatusCode: apiResponse.StatusCode})
	}
	if err := dbaas.SetUserPgSqlData(d, &user); err != nil {
		return nil, utils.ToError(d, err.Error(), nil)
	}
	if err := d.Set("cluster_id", clusterId); err != nil {
		return nil, utils.GenerateSetError("PgSQL user", "cluster_id", err)
	}
	return []*schema.ResourceData{d}, nil
}
