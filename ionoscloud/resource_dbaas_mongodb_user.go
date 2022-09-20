package ionoscloud

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	mongo "github.com/ionos-cloud/sdk-go-dbaas-mongo"
	"log"
)

func resourceDbaasMongoUser() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDbaasMongoDBUserCreate,
		ReadContext:   resourceDbaasMongoDBUserRead,
		//UpdateContext: resourceUserUpdate,
		DeleteContext: resourceDbaasMongoDBUserDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceDbaasMongoDBUserImporter,
		},
		Schema: map[string]*schema.Schema{
			"cluster_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"username": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				//ValidateFunc:     validation.All(validation.StringIsNotWhiteSpace),
				//DiffSuppressFunc: DiffToLower,
			},
			"database": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				//ValidateFunc: validation.All(validation.StringIsNotWhiteSpace),
			},
			"password": {
				Type:     schema.TypeString,
				Required: true,
				//ValidateFunc: validation.All(validation.StringIsNotWhiteSpace),
				Sensitive: true,
				ForceNew:  true,
			},
			"roles": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				//ValidateFunc: validation.All(validation.StringIsNotWhiteSpace),
			},
		},
		Timeouts: &resourceDefaultTimeouts,
	}
}

func resourceDbaasMongoDBUserCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(SdkBundle).MongoClient
	request := mongo.User{
		Properties: &mongo.UserProperties{},
	}

	var clusterId string
	if d.Get("cluster_id") != nil {
		clusterId = d.Get("cluster_id").(string)
	}

	if d.Get("username") != nil {
		username := d.Get("username").(string)
		request.Properties.Username = &username
	}
	if d.Get("database") != nil {
		database := d.Get("database").(string)
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
				roleVal := role.(string)
				mongoRole := mongo.UserRoles{
					Role:     &roleVal,
					Database: request.Properties.Database,
				}
				roles = append(roles, mongoRole)
			}
		}

		request.Properties.Roles = &roles
	}

	//rsp, apiResponse, err := client2.UserManagementApi.UmUsersPost(ctx).User(request).Execute()
	_, apiResponse, err := client.UsersApi.ClustersUsersPost(ctx, clusterId).User(request).Execute()
	//logApiRequestTime(apiResponse)

	if err != nil {
		diags := diag.FromErr(fmt.Errorf("an error occurred while adding a user to mongoDB: %w", err))
		return diags
	}

	//d.SetId(*rsp.Id)

	// Wait, catching any errors
	_, errState := getStateChangeConf(meta, d, apiResponse.Header.Get("Location"), schema.TimeoutCreate).WaitForStateContext(ctx)
	if errState != nil {
		if IsRequestFailed(err) {
			// Request failed, so resource was not created, delete resource from state file
			d.SetId("")
		}
		diags := diag.FromErr(errState)
		return diags
	}
	return resourceUserMongoRead(ctx, d, meta)
}

func resourceDbaasMongoDBUserRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	//client := meta.(SdkBundle).CloudApiClient
	client := meta.(SdkBundle).MongoClient

	clusterId := d.Get("cluster_id").(string)
	database := d.Get("database").(string)
	username := d.Get("username").(string)
	//user, apiResponse, err := client.UserManagementApi.UmUsersFindById(ctx, d.Id()).Execute()
	user, apiResponse, err := client.UsersApi.ClustersUsersFindById(ctx, clusterId, database, username).Execute()
	//logApiRequestTime(apiResponse)

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

func resourceDbaasMongoDBUserDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	//client := meta.(SdkBundle).CloudApiClient
	client := meta.(SdkBundle).MongoClient

	clusterId := d.Get("cluster_id").(string)
	database := d.Get("database").(string)
	username := d.Get("username").(string)
	//apiResponse, err := client.UserManagementApi.UmUsersDelete(ctx, d.Id()).Execute()
	_, apiResponse, err := client.UsersApi.ClustersUsersDelete(ctx, clusterId, database, username).Execute()
	//logApiRequestTime(apiResponse)
	if err != nil {
		diags := diag.FromErr(err)
		return diags
	}

	// Wait, catching any errors
	_, errState := getStateChangeConf(meta, d, apiResponse.Header.Get("Location"), schema.TimeoutDelete).WaitForStateContext(ctx)
	if errState != nil {
		diags := diag.FromErr(errState)
		return diags
	}

	d.SetId("")
	return nil

}

func resourceDbaasMongoDBUserImporter(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	//client := meta.(SdkBundle).CloudApiClient
	client := meta.(SdkBundle).MongoClient

	userId := d.Id()

	clusterId := d.Get("cluster_id").(string)
	database := d.Get("database").(string)
	username := d.Get("username").(string)
	//user, apiResponse, err := client.UserManagementApi.UmUsersFindById(ctx, userId).Execute()
	user, apiResponse, err := client.UsersApi.ClustersUsersFindById(ctx, clusterId, database, username).Execute()
	//logApiRequestTime(apiResponse)

	if err != nil {
		if apiResponse != nil && apiResponse.Response != nil && apiResponse.StatusCode == 404 {
			d.SetId("")
			return nil, fmt.Errorf("an error occured while trying to fetch the user %q", userId)
		}
		return nil, fmt.Errorf("user does not exist%q", userId)
	}

	if err := setUserMongoData(d, &user); err != nil {
		return nil, err
	}

	log.Printf("[INFO] user found: %+v", user)

	return []*schema.ResourceData{d}, nil
}

func setUserMongoData(d *schema.ResourceData, user *mongo.User) error {
	//d.SetId(*user.Id)

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
		if user.Properties.Password != nil {
			if err := d.Set("password", *user.Properties.Password); err != nil {
				return err
			}
		}
		if user.Properties.Roles != nil {
			if err := d.Set("roles", *user.Properties.Roles); err != nil {
				return err
			}
		}
	}

	return nil
}

func resourceUserMongoRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	//client := meta.(SdkBundle).CloudApiClient
	client := meta.(SdkBundle).MongoClient

	clusterId := d.Get("cluster_id").(string)
	database := d.Get("database").(string)
	username := d.Get("username").(string)
	//user, apiResponse, err := client.UserManagementApi.UmUsersFindById(ctx, d.Id()).Execute()
	user, apiResponse, err := client.UsersApi.ClustersUsersFindById(ctx, clusterId, database, username).Execute()
	//logApiRequestTime(apiResponse)

	if err != nil {
		if apiResponse != nil && apiResponse.Response != nil && apiResponse.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		diags := diag.FromErr(fmt.Errorf("an error occured while fetching a mongo User ID %s %w", d.Id(), err))
		return diags
	}

	if err := setUserMongoData(d, &user); err != nil {
		return diag.FromErr(err)
	}

	return nil
}
