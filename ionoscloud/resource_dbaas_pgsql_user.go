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
	client := meta.(services.SdkBundle).PsqlClient

	clusterId := d.Get("cluster_id").(string)
	username := d.Get("username").(string)
	password := d.Get("password").(string)

	request := pgsql.User{
		Properties: &pgsql.UserProperties{},
	}
	request.Properties.Username = &username
	request.Properties.Password = &password

	user, _, err := client.CreateUser(ctx, clusterId, request)
	if err != nil {
		return diag.FromErr(fmt.Errorf("an error occurred while adding the user: %s to the PgSql cluster with ID: %s, error: %w", username, clusterId, err))
	}
	d.SetId(*user.Id)
	// Wait for the cluster to be ready again (when creating/updating the user, the cluster enters
	// 'BUSY' state).
	err = utils.WaitForResourceToBeReady(ctx, d, client.IsClusterReady)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error while waiting for PgSql cluster with ID: %s to be ready, error: %w", clusterId, err))
	}
	return diag.FromErr(dbaas.SetUserPgSqlData(d, &user))
}

func resourceDbaasPgSqlUserUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(services.SdkBundle).PsqlClient

	request := pgsql.UsersPatchRequest{
		Properties: pgsql.NewPatchUserProperties(),
	}

	clusterId := d.Get("cluster_id").(string)
	username := d.Get("username").(string)

	if d.HasChange("password") {
		_, newPassword := d.GetChange("password")
		password := newPassword.(string)
		request.Properties.Password = &password
	}

	user, _, err := client.UpdateUser(ctx, clusterId, username, request)
	if err != nil {
		return diag.FromErr(fmt.Errorf("an error occurred while updating a PgSql user, username: %s, cluster ID: %s, error: %w", username, clusterId, err))
	}
	// Wait for the cluster to be ready again (when creating/updating the user, the cluster enters
	// 'BUSY' state).
	err = utils.WaitForResourceToBeReady(ctx, d, client.IsClusterReady)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error while waiting for PgSql cluster with ID: %s to be ready after user: %s update, error: %w", clusterId, username, err))
	}
	return diag.FromErr(dbaas.SetUserPgSqlData(d, &user))
}

func resourceDbaasPgSqlUserRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(services.SdkBundle).PsqlClient
	clusterId := d.Get("cluster_id").(string)
	username := d.Get("username").(string)

	user, apiResponse, err := client.FindUserByUsername(ctx, clusterId, username)
	if err != nil {
		if apiResponse.HttpNotFound() {
			d.SetId("")
			return nil
		}
		return diag.FromErr(fmt.Errorf("an error occurred while fetching the PgSql user with ID: %s, error: %w", d.Id(), err))
	}

	if err := dbaas.SetUserPgSqlData(d, &user); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceDbaaSPgSqlUserDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(services.SdkBundle).PsqlClient

	clusterId := d.Get("cluster_id").(string)
	username := d.Get("username").(string)
	_, err := client.DeleteUser(ctx, clusterId, username)
	if err != nil {
		diags := diag.FromErr(err)
		return diags
	}
	err = utils.WaitForResourceToBeDeleted(ctx, d, client.IsUserDeleted)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error while waiting for the PgSql username: %s to be deleted, error: %w", username, err))
	}
	return nil
}

func resourceDbaasPgSqlUserImporter(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	parts := strings.Split(d.Id(), "/")
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid import format: %s, expecting the following format: {clusterID}/{username}", d.Id())
	}
	clusterId := parts[0]
	username := parts[1]
	client := meta.(services.SdkBundle).PsqlClient
	user, apiResponse, err := client.FindUserByUsername(ctx, clusterId, username)
	if err != nil {
		if apiResponse.HttpNotFound() {
			d.SetId("")
			return nil, fmt.Errorf("unable to find PgSql username: %s, cluster ID: %s", username, clusterId)
		}
		return nil, fmt.Errorf("error occurred while fetching PgSql username: %s, cluster ID: %s, error: %w", username, clusterId, err)
	}
	if err := dbaas.SetUserPgSqlData(d, &user); err != nil {
		return nil, err
	}
	if err := d.Set("cluster_id", clusterId); err != nil {
		return nil, utils.GenerateSetError("PgSQL user", "cluster_id", err)
	}
	return []*schema.ResourceData{d}, nil
}
