package ionoscloud

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	pgsql "github.com/ionos-cloud/sdk-go-dbaas-postgres"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/dbaas"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils"
	"strings"
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
				Optional:    true,
			},
		},
		Timeouts: &resourceDefaultTimeouts,
	}
}

func resourceDbaasPgSqlUserCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(SdkBundle).PsqlClient

	clusterId := d.Get("cluster_id").(string)
	username := d.Get("username").(string)
	password := d.Get("password").(string)

	request := pgsql.User{
		Properties: &pgsql.UserProperties{},
	}
	request.Properties.Username = &username
	request.Properties.Password = &password

	if isSystemUser, ok := d.GetOk("is_system_user"); ok {
		isSystemUserValue := isSystemUser.(bool)
		request.Properties.System = &isSystemUserValue
	}

	user, _, err := client.CreateUser(ctx, clusterId, request)
	if err != nil {
		return diag.FromErr(fmt.Errorf("an error occured while adding an user to the PgSql cluster with ID: %s, error: %w", clusterId, err))
	}
	d.SetId(*user.Id)
	// Wait for the cluster to be ready again (when creating/updating the user, the cluster enters
	// 'BUSY' state).
	err = utils.WaitForResourceToBeReady(ctx, d, client.IsClusterReady)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error while waiting for PgSql cluster with ID: %s to be ready, error: %w", clusterId, err))
	}
	return resourceDbaasPgSqlUserRead(ctx, d, meta)
}

func resourceDbaasPgSqlUserUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(SdkBundle).PsqlClient

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

	_, _, err := client.UpdateUser(ctx, clusterId, username, request)
	if err != nil {
		return diag.FromErr(fmt.Errorf("an error occured while updating a PgSql user, username: %s, cluster ID: %s, error: %w", username, clusterId, err))
	}
	// Wait for the cluster to be ready again (when creating/updating the user, the cluster enters
	// 'BUSY' state).
	err = utils.WaitForResourceToBeReady(ctx, d, client.IsClusterReady)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error while waiting for PgSql cluster with ID: %s to be ready, error: %w", clusterId, err))
	}
	return resourceDbaasPgSqlUserRead(ctx, d, meta)
}

func resourceDbaasPgSqlUserRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(SdkBundle).PsqlClient
	clusterId := d.Get("cluster_id").(string)
	username := d.Get("username").(string)

	user, apiResponse, err := client.FindUserByUsername(ctx, clusterId, username)
	if err != nil {
		if apiResponse.HttpNotFound() {
			d.SetId("")
			return nil
		}
		return diag.FromErr(fmt.Errorf("an error occured while fetching the PgSql user with ID: %s, error: %w", d.Id(), err))
	}

	if err := dbaas.SetUserPgSqlData(d, &user); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceDbaaSPgSqlUserDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(SdkBundle).PsqlClient

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

// TODO -- make this in the proper way, the current implementation is not working because of this error:
// Error: Cannot import non-existent remote object.
// While attempting to import an existing object to "ionoscloud_pg_user.exampleuser", the provider detected that no object exists with the given id.
// Only pre-existing objects can be imported; check
// │ that the id is correct and that it is associated with the provider's configured region or endpoint,
// or use "terraform apply" to create a new remote object for this resource.
// From the error it looks like I'm trying to import a user that doesn't exist, but the user exists. I also debugged a little bit
// and it seems like the user is retrieved properly. I think that Terraform is bothered by the fact that the name is used to retrieve
// the user, not the ID.
func resourceDbaasPgSqlUserImporter(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	parts := strings.Split(d.Id(), "/")
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid import format: %s, expecting the following format: {clusterID}/{username}", d.Id())
	}
	clusterId := parts[0]
	username := parts[1]
	client := meta.(SdkBundle).PsqlClient
	user, apiResponse, err := client.FindUserByUsername(ctx, clusterId, username)
	if err != nil {
		if apiResponse.HttpNotFound() {
			d.SetId("")
			return nil, fmt.Errorf("unable to find PgSql username: %s, cluster ID: %s", username, clusterId)
		}
		return nil, fmt.Errorf("error occured while fetching PgSql username: %s, cluster ID: %s, error: %w", username, clusterId, err)
	}
	if err := dbaas.SetUserPgSqlData(d, &user); err != nil {
		return nil, err
	}
	return []*schema.ResourceData{d}, nil
}
