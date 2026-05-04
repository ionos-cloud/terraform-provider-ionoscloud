package ionoscloud

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/ionos-cloud/sdk-go-bundle/products/dbaas/mongo/v2"

	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/bundleclient"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/dbaas"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils"
	diagutil "github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils/diags"
)

func resourceDbaasMongoUser() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDbaasMongoUserCreate,
		UpdateContext: resourceDbaasMongoUserUpdate,
		ReadContext:   resourceDbaasMongoUserRead,
		DeleteContext: resourceDbaasMongoUserDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceDbaasMongoUserImporter,
		},
		Schema: map[string]*schema.Schema{
			"location": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "The location of the resource. This field should be used only if you are also using a file configuration and should not be configured otherwise.",
			},
			"cluster_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"username": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"password": {
				Type:             schema.TypeString,
				Required:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringLenBetween(10, 100)),
				Sensitive:        true,
			},
			"roles": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"role": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "A list of mongodb user roles. Examples: read, readWrite, readAnyDatabase",
						},
						"database": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
		},
		Timeouts: &resourceDefaultTimeouts,
	}
}

func resourceDbaasMongoUserCreate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client, err := meta.(bundleclient.SdkBundle).NewMongoClient(ctx, d.Get("location").(string))
	if err != nil {
		return diag.FromErr(err)
	}
	request := mongo.User{
		Properties: &mongo.UserProperties{},
	}

	var clusterId string
	if d.Get("cluster_id") != nil {
		clusterId = d.Get("cluster_id").(string)
	}
	username := ""
	if d.Get("username") != nil {
		username = d.Get("username").(string)
		request.Properties.Username = username
	}
	if d.Get("password") != nil {
		password := d.Get("password").(string)
		request.Properties.Password = password
	}
	if rolesValue, ok := d.GetOk("roles"); ok {
		roles := make([]mongo.UserRoles, 0)
		rolesValue := rolesValue.([]any)
		if rolesValue != nil {
			for _, role := range rolesValue {
				roleVal := role.(map[string]any)
				roleStr := roleVal["role"].(string)
				roleDb := roleVal["database"].(string)
				mongoRole := mongo.UserRoles{
					Role:     &roleStr,
					Database: &roleDb,
				}
				roles = append(roles, mongoRole)
			}
		}
		request.Properties.Roles = roles
	}

	user, apiResponse, err := client.CreateUser(ctx, clusterId, request)
	if err != nil {
		return diagutil.ToDiags(d, fmt.Errorf("an error occurred while adding a user to mongoDB: %w", err), &diagutil.ErrorContext{StatusCode: apiResponse.SafeStatusCode()})
	}

	if user.Properties != nil {
		d.SetId(clusterId + user.Properties.Username)
	}

	err = utils.WaitForResourceToBeReady(ctx, d, client.IsUserReady)
	if err != nil {
		return diagutil.ToDiags(d, fmt.Errorf("an error occurred while waiting for the Mongo user to become available, user ID: %v, error: %w", clusterId+user.Properties.Username, err), &diagutil.ErrorContext{Timeout: d.Timeout(schema.TimeoutCreate).String()})
	}

	return diagutil.ToDiags(d, dbaas.SetUserMongoData(d, &user), nil)
}

func resourceDbaasMongoUserUpdate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client, err := meta.(bundleclient.SdkBundle).NewMongoClient(ctx, d.Get("location").(string))
	if err != nil {
		return diag.FromErr(err)
	}
	request := mongo.PatchUserRequest{
		Properties: mongo.NewPatchUserProperties(),
	}

	clusterId := d.Get("cluster_id").(string)
	username := d.Get("username").(string)

	if d.HasChange("password") {
		_, password := d.GetChange("password")
		pwdStr := password.(string)
		request.Properties.Password = &pwdStr
	}

	if d.HasChange("roles") {
		_, rolesIntf := d.GetChange("roles")
		roles := make([]mongo.UserRoles, 0)
		rolesValue := rolesIntf.([]any)
		if rolesValue != nil {
			for _, role := range rolesValue {
				roleVal := role.(map[string]any)
				roleStr := roleVal["role"].(string)
				roleDb := roleVal["database"].(string)
				mongoRole := mongo.UserRoles{
					Role:     &roleStr,
					Database: &roleDb,
				}
				roles = append(roles, mongoRole)
			}
		}
		request.Properties.Roles = roles
	}

	user, apiResponse, err := client.UpdateUser(ctx, clusterId, username, request)
	if err != nil {
		return diagutil.ToDiags(d, fmt.Errorf("an error occurred while updating a Mongo user, username: %v, cluster ID: %v, error: %w", username, clusterId, err), &diagutil.ErrorContext{StatusCode: apiResponse.SafeStatusCode()})
	}

	err = utils.WaitForResourceToBeReady(ctx, d, client.IsUserReady)
	if err != nil {
		return diagutil.ToDiags(d, fmt.Errorf("an error occurred while waiting for a Mongo user to become available after an update, username: %v, cluster ID: %v, error: %w", username, clusterId, err), &diagutil.ErrorContext{Timeout: d.Timeout(schema.TimeoutUpdate).String()})
	}

	return diagutil.ToDiags(d, dbaas.SetUserMongoData(d, &user), nil)
}

func resourceDbaasMongoUserRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client, err := meta.(bundleclient.SdkBundle).NewMongoClient(ctx, d.Get("location").(string))
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
		return diagutil.ToDiags(d, fmt.Errorf("an error occurred while fetching a User: %w", err), &diagutil.ErrorContext{StatusCode: apiResponse.SafeStatusCode()})
	}

	if err := dbaas.SetUserMongoData(d, &user); err != nil {
		return diagutil.ToDiags(d, err, nil)
	}

	return nil
}

func resourceDbaasMongoUserDelete(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client, err := meta.(bundleclient.SdkBundle).NewMongoClient(ctx, d.Get("location").(string))
	if err != nil {
		return diag.FromErr(err)
	}

	clusterId := d.Get("cluster_id").(string)
	username := d.Get("username").(string)
	apiResponse, err := client.DeleteUser(ctx, clusterId, username)
	if err != nil {
		return diagutil.ToDiags(d, fmt.Errorf("an error occurred while deleting the Mongo user: %w", err), &diagutil.ErrorContext{StatusCode: apiResponse.SafeStatusCode()})
	}

	// Wait, catching any errors
	err = utils.WaitForResourceToBeDeleted(ctx, d, client.IsUserDeleted)
	if err != nil {
		return diagutil.ToDiags(d, fmt.Errorf("an error occurred while waiting for the Mongo user: to be deleted, error: %w", err), &diagutil.ErrorContext{Timeout: d.Timeout(schema.TimeoutDelete).String()})
	}

	d.SetId("")
	return nil

}

func resourceDbaasMongoUserImporter(ctx context.Context, d *schema.ResourceData, meta any) ([]*schema.ResourceData, error) {
	importID := d.Id()
	location, parts := splitImportID(importID, "/")
	if len(parts) != 2 {
		return nil, diagutil.ToError(d, fmt.Errorf("invalid import identifier: expected one of <location>:<cluster-id>/<username> or <cluster-id>/<username>, got: %s", importID), nil)
	}

	if err := validateImportIDParts(parts); err != nil {
		return nil, diagutil.ToError(d, fmt.Errorf("failed validating import identifier %q: %w", importID, err), nil)
	}

	client, err := meta.(bundleclient.SdkBundle).NewMongoClient(ctx, location)
	if err != nil {
		return nil, err
	}

	clusterID := parts[0]
	username := parts[1]

	user, apiResponse, err := client.FindUserByUsername(ctx, clusterID, username)
	if err != nil {
		if apiResponse.HttpNotFound() {
			d.SetId("")
			return nil, diagutil.ToError(d, fmt.Errorf("unable to find MongoDB user: %s, cluster ID: %s", username, clusterID), &diagutil.ErrorContext{StatusCode: apiResponse.SafeStatusCode()})
		}
		return nil, diagutil.ToError(d, fmt.Errorf("error occurred while fetching MongoDB user: %s, cluster ID: %s, error: %w", username, clusterID, err), &diagutil.ErrorContext{StatusCode: apiResponse.SafeStatusCode()})
	}
	if err := dbaas.SetUserMongoData(d, &user); err != nil {
		return nil, diagutil.ToError(d, err, nil)
	}
	if err := d.Set("cluster_id", clusterID); err != nil {
		return nil, utils.GenerateSetError("MongoDB user", "cluster_id", err)
	}
	if err := d.Set("location", location); err != nil {
		return nil, utils.GenerateSetError("MongoDB user", "location", err)
	}
	return []*schema.ResourceData{d}, nil
}
