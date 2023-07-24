package ionoscloud

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	mongo "github.com/ionos-cloud/sdk-go-dbaas-mongo"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services"
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
	client := meta.(services.SdkBundle).MongoClient
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
	client := meta.(services.SdkBundle).MongoClient
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

	user, _, err := client.UpdateUser(ctx, clusterId, username, request)
	if err != nil {
		diags := diag.FromErr(fmt.Errorf("an error occured while updating a user to mongoDB cluster %s: %w", clusterId, err))
		return diags
	}

	err = utils.WaitForResourceToBeReady(ctx, d, client.IsUserReady)
	if err != nil {
		return diag.FromErr(fmt.Errorf("updating %w", err))
	}

	return diag.FromErr(dbaas.SetUserMongoData(d, &user))
}

func resourceDbaasMongoUserRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(services.SdkBundle).MongoClient

	clusterId := d.Get("cluster_id").(string)
	username := d.Get("username").(string)

	user, apiResponse, err := client.FindUserByUsername(ctx, clusterId, username)

	if err != nil {
		if apiResponse.HttpNotFound() {
			d.SetId("")
			return nil
		}
		return diag.FromErr(fmt.Errorf("an error occured while fetching a User ID %s %w", d.Id(), err))
	}

	if err := dbaas.SetUserMongoData(d, &user); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceDbaasMongoUserDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(services.SdkBundle).MongoClient

	clusterId := d.Get("cluster_id").(string)
	username := d.Get("username").(string)
	_, err := client.DeleteUser(ctx, clusterId, username)
	if err != nil {
		diags := diag.FromErr(err)
		return diags
	}

	// Wait, catching any errors
	err = utils.WaitForResourceToBeDeleted(ctx, d, client.IsUserDeleted)
	if err != nil {
		return diag.FromErr(fmt.Errorf("user deleted %w", err))
	}

	if err != nil {
		return diag.FromErr(fmt.Errorf("user deleted %w", err))
	}

	d.SetId("")
	return nil

}

func resourceDbaasMongoUserImporter(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	client := meta.(services.SdkBundle).MongoClient

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
