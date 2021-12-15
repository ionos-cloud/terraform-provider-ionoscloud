package ionoscloud

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	ionoscloud "github.com/ionos-cloud/sdk-go/v5"
)

func resourceUser() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceUserCreate,
		ReadContext:   resourceUserRead,
		UpdateContext: resourceUserUpdate,
		DeleteContext: resourceUserDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceUserImporter,
		},
		Schema: map[string]*schema.Schema{
			"first_name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.All(validation.StringIsNotWhiteSpace),
			},
			"last_name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.All(validation.StringIsNotWhiteSpace),
			},
			"email": {
				Type:             schema.TypeString,
				Required:         true,
				ValidateFunc:     validation.All(validation.StringIsNotWhiteSpace),
				DiffSuppressFunc: DiffToLower,
			},
			"password": {
				Type:         schema.TypeString,
				Required:     true,
				Sensitive:    true,
				ValidateFunc: validation.All(validation.StringIsNotWhiteSpace),
			},
			"administrator": {
				Type:     schema.TypeBool,
				Required: true,
			},
			"force_sec_auth": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"sec_auth_active": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"s3_canonical_user_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"active": {
				Type:     schema.TypeBool,
				Optional: true,
			},
		},
		Timeouts: &resourceDefaultTimeouts,
	}
}

func resourceUserCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*ionoscloud.APIClient)

	request := ionoscloud.UserPost{
		Properties: &ionoscloud.UserPropertiesPost{},
	}

	if d.Get("first_name") != nil {
		firstName := d.Get("first_name").(string)
		request.Properties.Firstname = &firstName
	}
	if d.Get("last_name") != nil {
		lastName := d.Get("last_name").(string)
		request.Properties.Lastname = &lastName
	}
	if d.Get("email") != nil {
		email := d.Get("email").(string)
		request.Properties.Email = &email
	}
	if d.Get("password") != nil {
		password := d.Get("password").(string)
		request.Properties.Password = &password
	}

	administrator := d.Get("administrator").(bool)
	request.Properties.Administrator = &administrator

	forceSecAuth := d.Get("force_sec_auth").(bool)
	request.Properties.ForceSecAuth = &forceSecAuth

	active := d.Get("active").(bool)
	request.Properties.Active = &active

	rsp, apiResponse, err := client.UserManagementApi.UmUsersPost(ctx).User(request).Execute()
	logApiRequestTime(apiResponse)

	if err != nil {
		diags := diag.FromErr(fmt.Errorf("an error occured while creating a user: %s", err))
		return diags
	}
	if rsp.Id != nil {
		log.Printf("[DEBUG] USER ID: %s", *rsp.Id)
		d.SetId(*rsp.Id)
	}

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
	return resourceUserRead(ctx, d, meta)
}

func resourceUserRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*ionoscloud.APIClient)

	user, apiResponse, err := client.UserManagementApi.UmUsersFindById(ctx, d.Id()).Execute()
	logApiRequestTime(apiResponse)

	if err != nil {
		if apiResponse != nil && apiResponse.Response != nil && apiResponse.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		diags := diag.FromErr(fmt.Errorf("an error occured while fetching a User ID %s %s", d.Id(), err))
		return diags
	}

	if err := setUserData(d, &user); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceUserUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*ionoscloud.APIClient)
	foundUser, apiResponse, err := client.UserManagementApi.UmUsersFindById(ctx, d.Id()).Execute()
	logApiRequestTime(apiResponse)
	if err != nil {
		diags := diag.FromErr(fmt.Errorf("an error occured while fetching a User ID %s %s", d.Id(), err))
		return diags
	}
	userReq := ionoscloud.UserPut{
		Properties: &ionoscloud.UserPropertiesPut{},
	}
	//this is a PUT, so we first fill everything, then we change what has been modified
	if foundUser.Properties != nil {
		userReq.Properties.Firstname = foundUser.Properties.Firstname
		userReq.Properties.Lastname = foundUser.Properties.Lastname
		userReq.Properties.Email = foundUser.Properties.Email
		userReq.Properties.Active = foundUser.Properties.Active
		userReq.Properties.Administrator = foundUser.Properties.Administrator
		userReq.Properties.ForceSecAuth = foundUser.Properties.ForceSecAuth
	}
	if d.Get("first_name") != nil {
		firstName := d.Get("first_name").(string)
		userReq.Properties.Firstname = &firstName
	}
	if d.Get("last_name") != nil {
		lastName := d.Get("last_name").(string)
		userReq.Properties.Lastname = &lastName
	}
	if d.Get("email") != nil {
		email := d.Get("email").(string)
		userReq.Properties.Email = &email
	}

	administrator := d.Get("administrator").(bool)
	userReq.Properties.Administrator = &administrator

	forceSecAuth := d.Get("force_sec_auth").(bool)
	userReq.Properties.ForceSecAuth = &forceSecAuth

	active := d.Get("active").(bool)
	userReq.Properties.Active = &active

	if d.HasChange("password") {
		password := d.Get("password").(string)
		userReq.Properties.Password = &password
	}

	_, apiResponse, err = client.UserManagementApi.UmUsersPut(ctx, d.Id()).User(userReq).Execute()
	logApiRequestTime(apiResponse)
	if err != nil {
		diags := diag.FromErr(fmt.Errorf("an error occured while patching a user ID %s %s", d.Id(), err))
		return diags
	}

	// Wait, catching any errors
	_, errState := getStateChangeConf(meta, d, apiResponse.Header.Get("Location"), schema.TimeoutUpdate).WaitForStateContext(ctx)
	if errState != nil {
		diags := diag.FromErr(errState)
		return diags
	}

	return resourceUserRead(ctx, d, meta)
}

func resourceUserDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*ionoscloud.APIClient)

	_, apiResponse, err := client.UserManagementApi.UmUsersDelete(ctx, d.Id()).Execute()
	logApiRequestTime(apiResponse)
	if err != nil {
		//try again in 20 seconds
		time.Sleep(20 * time.Second)
		_, apiResponse, err = client.UserManagementApi.UmUsersDelete(ctx, d.Id()).Execute()
		logApiRequestTime(apiResponse)
		if err != nil {
			diags := diag.FromErr(fmt.Errorf("an error occured while deleting a user %s %s", d.Id(), err))
			return diags
		}
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

func resourceUserImporter(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	client := meta.(*ionoscloud.APIClient)

	userId := d.Id()

	user, apiResponse, err := client.UserManagementApi.UmUsersFindById(ctx, userId).Execute()
	logApiRequestTime(apiResponse)

	if err != nil {
		if apiResponse != nil && apiResponse.Response != nil && apiResponse.StatusCode == 404 {
			d.SetId("")
			return nil, fmt.Errorf("an error occured while trying to fetch the user %q", userId)
		}
		return nil, fmt.Errorf("user does not exist%q", userId)
	}

	if err := setUserData(d, &user); err != nil {
		return nil, err
	}

	log.Printf("[INFO] user found: %+v", user)

	return []*schema.ResourceData{d}, nil
}

func setUserData(d *schema.ResourceData, user *ionoscloud.User) error {
	if user.Id != nil {
		d.SetId(*user.Id)
	}
	if user.Properties != nil {
		if user.Properties.Firstname != nil {
			if err := d.Set("first_name", *user.Properties.Firstname); err != nil {
				return err
			}
		}

		if user.Properties.Lastname != nil {
			if err := d.Set("last_name", *user.Properties.Lastname); err != nil {
				return err
			}
		}
		if user.Properties.Email != nil {
			if err := d.Set("email", *user.Properties.Email); err != nil {
				return err
			}
		}
		if user.Properties.Administrator != nil {
			if err := d.Set("administrator", *user.Properties.Administrator); err != nil {
				return err
			}
		}
		if user.Properties.ForceSecAuth != nil {
			if err := d.Set("force_sec_auth", *user.Properties.ForceSecAuth); err != nil {
				return err
			}
		}

		if user.Properties.SecAuthActive != nil {
			if err := d.Set("sec_auth_active", *user.Properties.SecAuthActive); err != nil {
				return err
			}
		}

		if user.Properties.S3CanonicalUserId != nil {
			if err := d.Set("s3_canonical_user_id", *user.Properties.S3CanonicalUserId); err != nil {
				return err
			}
		}

		if user.Properties.Active != nil {
			if err := d.Set("active", *user.Properties.Active); err != nil {
				return err
			}
		}
	}

	return nil
}
