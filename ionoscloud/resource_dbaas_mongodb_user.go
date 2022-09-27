package ionoscloud

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	mongo "github.com/ionos-cloud/sdk-go-dbaas-mongo"
	"log"
)

func resourceDbaasMongoUser() *schema.Resource {
	return &schema.Resource{
		// no update defined, forcenew on all fields
		CreateContext: resourceDbaasMongoUserCreate,
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
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The user database uses for authentication",
			},
			"database": {
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				Description:      "The user database to use for authentication",
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringIsNotWhiteSpace),
			},
			"password": {
				Type:             schema.TypeString,
				Required:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringIsNotWhiteSpace),
				Sensitive:        true,
				ForceNew:         true,
			},
			"roles": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"role": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "A list of mongodb user roles. Examples: read, readWrite, readAnyDatabase",
							ForceNew:    true,
						},
						"database": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
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
	database := ""
	if d.Get("database") != nil {
		database = d.Get("database").(string)
		request.Properties.Database = &database
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

	_, _, err := client.UsersApi.ClustersUsersPost(ctx, clusterId).User(request).Execute()
	if err != nil {
		diags := diag.FromErr(fmt.Errorf("an error occurred while adding a user to mongoDB: %w", err))
		return diags
	}

	var user = &mongo.User{}
	err = resource.RetryContext(ctx, *resourceDefaultTimeouts.Create, func() *resource.RetryError {
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
		if user != nil && user.Metadata != nil && *user.Metadata.State != "AVAILABLE" {
			log.Printf("[INFO] user is still getting created, is in state %s", *user.Metadata.State)
			return resource.RetryableError(fmt.Errorf("user is still getting created, is in state %s", *user.Metadata.State))
		}
		return nil
	})

	if err != nil {
		return diag.FromErr(err)
	}

	if user == nil || user.Properties == nil || *user.Properties.Username == "" {
		return diag.FromErr(fmt.Errorf("could not find username %s with database %s in cluster %s after creation ", username, database, clusterId))
	}

	if user.Properties != nil && user.Properties.Username != nil && user.Properties.Database != nil {
		d.SetId(clusterId + *user.Properties.Username + *user.Properties.Database)
	}

	return diag.FromErr(setUserMongoData(d, user))
}

func resourceDbaasMongoUserRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(SdkBundle).MongoClient

	clusterId := d.Get("cluster_id").(string)
	database := d.Get("database").(string)
	username := d.Get("username").(string)

	user, apiResponse, err := client.UsersApi.ClustersUsersFindById(ctx, clusterId, database, username).Execute()

	if err != nil {
		if apiResponse != nil && apiResponse.Response != nil && apiResponse.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		diags := diag.FromErr(fmt.Errorf("an error occured while fetching a User ID %s %w", d.Id(), err))
		return diags
	}

	if err := setUserMongoData(d, &user); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceDbaasMongoUserDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(SdkBundle).MongoClient

	clusterId := d.Get("cluster_id").(string)
	database := d.Get("database").(string)
	username := d.Get("username").(string)
	_, _, err := client.UsersApi.ClustersUsersDelete(ctx, clusterId, database, username).Execute()
	if err != nil {
		diags := diag.FromErr(err)
		return diags
	}

	// Wait, catching any errors
	err = resource.RetryContext(ctx, *resourceDefaultTimeouts.Create, func() *resource.RetryError {
		var err error
		var apiResponse *mongo.APIResponse
		var user = mongo.User{}
		user, apiResponse, err = client.UsersApi.ClustersUsersFindById(ctx, clusterId, database, username).Execute()
		if apiResponse.HttpNotFound() {
			log.Printf("[INFO] Deleted successfuly user %s with database %s in cluster %s", username, database, clusterId)
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
		return diag.FromErr(err)
	}

	d.SetId("")
	return nil

}

func resourceDbaasMongoUserImporter(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	client := meta.(SdkBundle).MongoClient

	userId := d.Id()

	clusterId := d.Get("cluster_id").(string)
	database := d.Get("database").(string)
	username := d.Get("username").(string)

	user, apiResponse, err := client.UsersApi.ClustersUsersFindById(ctx, clusterId, database, username).Execute()

	if err != nil {
		if apiResponse != nil && apiResponse.Response != nil && apiResponse.StatusCode == 404 {
			d.SetId("")
			return nil, fmt.Errorf("an error occured while trying to fetch the mongodb user %q", userId)
		}
		return nil, fmt.Errorf("mongodb user does not exist %q", userId)
	}

	if err := setUserMongoData(d, &user); err != nil {
		return nil, err
	}

	log.Printf("[INFO] user found: %+v", user)

	return []*schema.ResourceData{d}, nil
}

func setUserMongoData(d *schema.ResourceData, user *mongo.User) error {
	if user.Properties != nil {
		if user.Properties.Username != nil {
			if err := d.Set("username", *user.Properties.Username); err != nil {
				return err
			}
		}
		if user.Properties.Database != nil {
			if err := d.Set("database", *user.Properties.Database); err != nil {
				return err
			}
		}

		if user.Properties.Roles != nil && len(*user.Properties.Roles) > 0 {
			userRoles := make([]interface{}, len(*user.Properties.Roles))
			for index, user := range *user.Properties.Roles {
				userEntry := make(map[string]interface{})

				if user.Role != nil {
					userEntry["role"] = *user.Role
				}

				if user.Database != nil {
					userEntry["database"] = user.Database
				}
				userRoles[index] = userEntry
			}

			if len(userRoles) > 0 {
				if err := d.Set("roles", userRoles); err != nil {
					return fmt.Errorf("error setting user roles for user (%w)", err)
				}
			}
		}
	}
	return nil

}
