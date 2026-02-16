package ionoscloud

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	mongo "github.com/ionos-cloud/sdk-go-bundle/products/dbaas/mongo/v2"

	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/bundleclient"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/dbaas"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils"
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

func resourceDbaasMongoUserCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(bundleclient.SdkBundle).MongoClient
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
		rolesValue := rolesValue.([]interface{})
		if rolesValue != nil {
			for _, role := range rolesValue {
				roleVal := role.(map[string]interface{})
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

	user, _, err := client.CreateUser(ctx, clusterId, request)
	if err != nil {
		return utils.ToDiags(d, fmt.Sprintf("an error occurred while adding a user to mongoDB: %s", err), nil)
	}

	if user.Properties != nil {
		d.SetId(clusterId + user.Properties.Username)
	}

	err = utils.WaitForResourceToBeReady(ctx, d, client.IsUserReady)
	if err != nil {
		return utils.ToDiags(d, fmt.Sprintf("an error occurred while waiting for the Mongo user to become available, user ID: %v, error: %s", clusterId+user.Properties.Username, err), nil)
	}

	return utils.ToDiags(d, dbaas.SetUserMongoData(d, &user).Error(), nil)
}

func resourceDbaasMongoUserUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(bundleclient.SdkBundle).MongoClient
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
		rolesValue := rolesIntf.([]interface{})
		if rolesValue != nil {
			for _, role := range rolesValue {
				roleVal := role.(map[string]interface{})
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
		return utils.ToDiags(d, fmt.Sprintf("an error occurred while updating a Mongo user, username: %v, cluster ID: %v, error: %s", username, clusterId, err), &utils.DiagsOpts{StatusCode: apiResponse.StatusCode})
	}

	err = utils.WaitForResourceToBeReady(ctx, d, client.IsUserReady)
	if err != nil {
		return utils.ToDiags(d, fmt.Sprintf("an error occurred while waiting for a Mongo user to become available after an update, username: %v, cluster ID: %v, error: %s", username, clusterId, err), nil)
	}

	return utils.ToDiags(d, dbaas.SetUserMongoData(d, &user).Error(), nil)
}

func resourceDbaasMongoUserRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(bundleclient.SdkBundle).MongoClient

	clusterId := d.Get("cluster_id").(string)
	username := d.Get("username").(string)

	user, apiResponse, err := client.FindUserByUsername(ctx, clusterId, username)

	if err != nil {
		if apiResponse.HttpNotFound() {
			d.SetId("")
			return nil
		}
		return utils.ToDiags(d, fmt.Sprintf("an error occurred while fetching a User: %s", err), &utils.DiagsOpts{StatusCode: apiResponse.StatusCode})
	}

	if err := dbaas.SetUserMongoData(d, &user); err != nil {
		return utils.ToDiags(d, err.Error(), nil)
	}

	return nil
}

func resourceDbaasMongoUserDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(bundleclient.SdkBundle).MongoClient

	clusterId := d.Get("cluster_id").(string)
	username := d.Get("username").(string)
	apiResponse, err := client.DeleteUser(ctx, clusterId, username)
	if err != nil {
		return utils.ToDiags(d, fmt.Sprintf("an error occurred while deleting the Mongo user: %s", err), &utils.DiagsOpts{StatusCode: apiResponse.StatusCode})
	}

	// Wait, catching any errors
	err = utils.WaitForResourceToBeDeleted(ctx, d, client.IsUserDeleted)
	if err != nil {
		return utils.ToDiags(d, fmt.Sprintf("an error occurred while waiting for the Mongo user: to be deleted, error: %s", err), &utils.DiagsOpts{Timeout: schema.TimeoutDelete})
	}

	d.SetId("")
	return nil

}

func resourceDbaasMongoUserImporter(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	client := meta.(bundleclient.SdkBundle).MongoClient

	parts := strings.Split(d.Id(), "/")
	if len(parts) != 2 {
		return nil, utils.ToError(d, "invalid import format:, expecting the following format: {clusterID}/{username}", nil)
	}
	clusterID := parts[0]
	username := parts[1]

	user, apiResponse, err := client.FindUserByUsername(ctx, clusterID, username)
	if err != nil {
		if apiResponse.HttpNotFound() {
			d.SetId("")
			return nil, utils.ToError(d, fmt.Sprintf("unable to find MongoDB user: %s, cluster ID: %s", username, clusterID), &utils.DiagsOpts{StatusCode: apiResponse.StatusCode})
		}
		return nil, utils.ToError(d, fmt.Sprintf("error occurred while fetching MongoDB user: %s, cluster ID: %s, error: %s", username, clusterID, err), &utils.DiagsOpts{StatusCode: apiResponse.StatusCode})
	}
	if err := dbaas.SetUserMongoData(d, &user); err != nil {
		return nil, utils.ToError(d, err.Error(), nil)
	}
	if err := d.Set("cluster_id", clusterID); err != nil {
		return nil, utils.GenerateSetError("MongoDB user", "cluster_id", err)
	}
	return []*schema.ResourceData{d}, nil
}
