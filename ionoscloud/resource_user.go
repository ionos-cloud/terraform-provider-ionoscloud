package ionoscloud

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	ionoscloud "github.com/ionos-cloud/sdk-go/v6"
	"log"
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
				ValidateFunc: validation.All(validation.StringIsNotWhiteSpace),
			},
			"administrator": {
				Type:     schema.TypeBool,
				Optional: true,
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
	client := meta.(SdkBundle).CloudApiClient
	request := ionoscloud.UserPost{
		Properties: &ionoscloud.UserPropertiesPost{},
	}

	log.Printf("[DEBUG] NAME %s", d.Get("first_name"))

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

	if _, ok := d.GetOk("sec_auth_active"); ok {
		diags := diag.FromErr(fmt.Errorf("sec_auth_active attribute is not allowed in create requests"))
		return diags
	}

	rsp, apiResponse, err := client.UserManagementApi.UmUsersPost(ctx).User(request).Execute()
	logApiRequestTime(apiResponse)

	if err != nil {
		diags := diag.FromErr(fmt.Errorf("an error occured while creating a user: %s", err))
		return diags
	}

	d.SetId(*rsp.Id)

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
	client := meta.(SdkBundle).CloudApiClient

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
	client := meta.(SdkBundle).CloudApiClient

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

	if d.HasChange("first_name") {
		_, newValue := d.GetChange("first_name")
		firstName := newValue.(string)
		userReq.Properties.Firstname = &firstName
	}

	if d.HasChange("last_name") {
		_, newValue := d.GetChange("last_name")
		lastName := newValue.(string)
		userReq.Properties.Lastname = &lastName
	}
	if d.HasChange("email") {
		_, newValue := d.GetChange("email")
		email := newValue.(string)
		userReq.Properties.Email = &email
	}

	if d.HasChange("active") {
		active := d.Get("active").(bool)
		userReq.Properties.Active = &active
	}

	if d.HasChange("administrator") {
		administrator := d.Get("administrator").(bool)
		userReq.Properties.Administrator = &administrator
	}

	if d.HasChange("force_sec_auth") {
		forceSecAuth := d.Get("force_sec_auth").(bool)
		userReq.Properties.ForceSecAuth = &forceSecAuth
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
	client := meta.(SdkBundle).CloudApiClient

	apiResponse, err := client.UserManagementApi.UmUsersDelete(ctx, d.Id()).Execute()
	logApiRequestTime(apiResponse)
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

func resourceUserImporter(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	client := meta.(SdkBundle).CloudApiClient

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
	d.SetId(*user.Id)

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
