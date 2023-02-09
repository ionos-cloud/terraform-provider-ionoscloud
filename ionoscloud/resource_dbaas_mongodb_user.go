package ionoscloud

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	mongo "github.com/ionos-cloud/sdk-go-dbaas-mongo"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/dbaas"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils"
	"log"
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
	client := meta.(SdkBundle).MongoClient
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
		request.Properties.Username = &username
	}
	if d.Get("password") != nil {
		password := d.Get("password").(string)
		request.Properties.Password = &password
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
		request.Properties.Roles = &roles
	}

	user, _, err := client.CreateUser(ctx, clusterId, request)
	if err != nil {
		diags := diag.FromErr(fmt.Errorf("an error occurred while adding a user to mongoDB: %w", err))
		return diags
	}

	if user.Properties != nil && user.Properties.Username != nil {
		d.SetId(clusterId + *user.Properties.Username)
	}

	err = utils.WaitForResourceToBeReady(ctx, d, client.IsUserReady)
	if err != nil {
		return diag.FromErr(fmt.Errorf("updating %w", err))
	}

	return diag.FromErr(dbaas.SetUserMongoData(d, &user))
}

func resourceDbaasMongoUserUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(SdkBundle).MongoClient
	request := mongo.PatchUserRequest{
		Properties: mongo.NewPatchUserProperties(),
	}

	var clusterId string
	if d.Get("cluster_id") != nil {
		clusterId = d.Get("cluster_id").(string)
	}
	username := ""
	if d.Get("username") != nil {
		username = d.Get("username").(string)
	}
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
		request.Properties.Roles = &roles
	}

	_, _, err := client.UsersApi.ClustersUsersPatch(ctx, clusterId, defaultMongoDatabase, username).PatchUserRequest(request).Execute()
	if err != nil {
		diags := diag.FromErr(fmt.Errorf("while updating a user to mongoDB cluster %s: %w", clusterId, err))
		return diags
	}

	err = utils.WaitForResourceToBeReady(ctx, d, client.IsUserReady)
	if err != nil {
		return diag.FromErr(fmt.Errorf("updating %w", err))
	}

	return diag.FromErr(setUserMongoData(d, &user))
}

// waitForUserToBeReady - keeps retrying until user is in 'available' state, or context deadline is reached
func waitForUserToBeReady(ctx context.Context, client *dbaas.MongoClient, clusterId, database, username string) (mongo.User, error) {
	var user = mongo.NewUser()
	err := resource.RetryContext(ctx, *resourceDefaultTimeouts.Update, func() *resource.RetryError {

		var err error
		var apiResponse *mongo.APIResponse
		*user, apiResponse, err = client.UsersApi.ClustersUsersFindById(ctx, clusterId, database, username).Execute()
		if apiResponse.HttpNotFound() {
			log.Printf("[INFO] Could not find username %s with database %s in cluster %s, retrying...", username, database, clusterId)
			return resource.RetryableError(fmt.Errorf("could not find username %s with database %s in cluster %s, %w", username, database, clusterId, err))
		}
		if err != nil {
			resource.NonRetryableError(err)
		}

		if user != nil && user.Metadata != nil && user.Metadata.State != nil && !strings.EqualFold(*user.Metadata.State, utils.Available) {
			log.Printf("[INFO] mongo user %s is still in state %s", username, *user.Metadata.State)
			return resource.RetryableError(fmt.Errorf("user is still in state %s", *user.Metadata.State))
		}
		return nil
	})
	if user == nil || user.Properties == nil || *user.Properties.Username == "" {
		return *user, fmt.Errorf("could not find username %s with database %s in cluster %s", username, database, clusterId)
	}
	return *user, err
}

func resourceDbaasMongoUserRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(SdkBundle).MongoClient

	clusterId := d.Get("cluster_id").(string)
	username := d.Get("username").(string)

	user, apiResponse, err := client.UsersApi.ClustersUsersFindById(ctx, clusterId, defaultMongoDatabase, username).Execute()

	if err != nil {
		if apiResponse.HttpNotFound() {
			d.SetId("")
			return nil
		}
		diags := diag.FromErr(fmt.Errorf("an error occured while fetching a User ID %s %w", d.Id(), err))
		return diags
	}

	if err := dbaas.SetUserMongoData(d, &user); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceDbaasMongoUserDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(SdkBundle).MongoClient

	clusterId := d.Get("cluster_id").(string)
	username := d.Get("username").(string)
	_, err := client.DeleteUser(ctx, clusterId, username)
	if err != nil {
		diags := diag.FromErr(err)
		return diags
	}

	// Wait, catching any errors
	err = resource.RetryContext(ctx, *resourceDefaultTimeouts.Create, func() *resource.RetryError {
		var err error
		var apiResponse *mongo.APIResponse
		var user = mongo.User{}
		user, apiResponse, err = client.UsersApi.ClustersUsersFindById(ctx, clusterId, defaultMongoDatabase, username).Execute()
		if apiResponse.HttpNotFound() {
			log.Printf("[INFO] Deleted successfuly user %s with database %s in cluster %s", username, defaultMongoDatabase, clusterId)
			return nil
		}
		if err != nil {
			resource.NonRetryableError(err)
		}
		if user.Properties != nil && user.Properties.Username != nil && *user.Properties.Username == username {
			return resource.RetryableError(fmt.Errorf("user still found, retrying"))
		}
		return resource.NonRetryableError(fmt.Errorf("unexpected error"))
	})

	if err != nil {
		return diag.FromErr(fmt.Errorf("user deleted %w ", err))
	}

	d.SetId("")
	return nil

}

func resourceDbaasMongoUserImporter(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	client := meta.(SdkBundle).MongoClient

	userId := d.Id()

	clusterId := d.Get("cluster_id").(string)
	username := d.Get("username").(string)

	user, apiResponse, err := client.FindUserByUsername(ctx, clusterId, username)

	if err != nil {
		if apiResponse.HttpNotFound() {
			d.SetId("")
			return nil, fmt.Errorf("an error occured while trying to fetch the mongodb user %q", userId)
		}
		return nil, fmt.Errorf("mongodb user does not exist %q", userId)
	}

	if err := dbaas.SetUserMongoData(d, &user); err != nil {
		return nil, err
	}

	log.Printf("[INFO] user found: %+v", user)

	return []*schema.ResourceData{d}, nil
}
