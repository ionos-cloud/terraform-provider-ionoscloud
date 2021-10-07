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
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.All(validation.StringIsNotWhiteSpace),
			},
			"password": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.All(validation.StringIsNotWhiteSpace),
			},
			"administrator": {
				Type:     schema.TypeBool,
				Required: true,
			},
			"force_sec_auth": {
				Type:     schema.TypeBool,
				Required: true,
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
	forceSecAuth := d.Get("force_sec_auth").(bool)
	request.Properties.Administrator = &administrator
	request.Properties.ForceSecAuth = &forceSecAuth

	active := d.Get("active").(bool)
	request.Properties.Active = &active

	if _, ok := d.GetOk("sec_auth_active"); ok {
		diags := diag.FromErr(fmt.Errorf("sec_auth_active attribute is not allowed in create requests"))
		return diags
	}

	rsp, apiResponse, err := client.UserManagementApi.UmUsersPost(ctx).User(request).Execute()

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
	client := meta.(*ionoscloud.APIClient)

	rsp, apiResponse, err := client.UserManagementApi.UmUsersFindById(ctx, d.Id()).Execute()

	if err != nil {
		if apiResponse != nil && apiResponse.Response != nil && apiResponse.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		diags := diag.FromErr(fmt.Errorf("an error occured while fetching a User ID %s %s", d.Id(), err))
		return diags
	}

	if err := d.Set("first_name", *rsp.Properties.Firstname); err != nil {
		diags := diag.FromErr(err)
		return diags
	}
	if err := d.Set("last_name", *rsp.Properties.Lastname); err != nil {
		diags := diag.FromErr(err)
		return diags
	}
	if err := d.Set("email", *rsp.Properties.Email); err != nil {
		diags := diag.FromErr(err)
		return diags
	}
	if err := d.Set("administrator", *rsp.Properties.Administrator); err != nil {
		diags := diag.FromErr(err)
		return diags
	}
	if err := d.Set("force_sec_auth", *rsp.Properties.ForceSecAuth); err != nil {
		diags := diag.FromErr(err)
		return diags
	}

	if rsp.Properties.SecAuthActive != nil {
		if err := d.Set("sec_auth_active", *rsp.Properties.SecAuthActive); err != nil {
			diags := diag.FromErr(err)
			return diags
		}
	}

	if rsp.Properties.S3CanonicalUserId != nil {
		if err := d.Set("s3_canonical_user_id", *rsp.Properties.S3CanonicalUserId); err != nil {
			diags := diag.FromErr(err)
			return diags
		}
	}

	if rsp.Properties.Active != nil {
		if err := d.Set("active", *rsp.Properties.Active); err != nil {
			diags := diag.FromErr(err)
			return diags
		}
	}

	return nil
}

func resourceUserUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*ionoscloud.APIClient)

	rsp, apiResponse, err := client.UserManagementApi.UmUsersFindById(ctx, d.Id()).Execute()

	if err != nil {
		diags := diag.FromErr(fmt.Errorf("an error occured while fetching a User ID %s %s", d.Id(), err))
		return diags
	}

	administrator := d.Get("administrator").(bool)
	forceSecAuth := d.Get("force_sec_auth").(bool)
	userReq := ionoscloud.UserPut{
		Properties: &ionoscloud.UserPropertiesPut{
			Administrator: &administrator,
			ForceSecAuth:  &forceSecAuth,
		},
	}

	if d.HasChange("first_name") {
		_, newValue := d.GetChange("first_name")
		firstName := newValue.(string)
		userReq.Properties.Firstname = &firstName

	} else {
		userReq.Properties.Firstname = rsp.Properties.Firstname
	}

	if d.HasChange("last_name") {
		_, newValue := d.GetChange("last_name")
		lastName := newValue.(string)
		userReq.Properties.Lastname = &lastName
	} else {
		userReq.Properties.Lastname = rsp.Properties.Lastname
	}

	if d.HasChange("email") {
		_, newValue := d.GetChange("email")
		email := newValue.(string)
		userReq.Properties.Email = &email
	} else {
		userReq.Properties.Email = rsp.Properties.Email
	}

	if d.HasChange("active") {
		_, newValue := d.GetChange("active")
		active := newValue.(bool)
		userReq.Properties.Active = &active
	} else {
		userReq.Properties.Active = rsp.Properties.Active
	}

	if d.HasChange("sec_auth_active") {
		_, newValue := d.GetChange("sec_auth_active")
		active := newValue.(bool)
		userReq.Properties.SecAuthActive = &active
	} else {
		userReq.Properties.SecAuthActive = rsp.Properties.SecAuthActive
	}

	rsp, apiResponse, err = client.UserManagementApi.UmUsersPut(ctx, d.Id()).User(userReq).Execute()
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

	apiResponse, err := client.UserManagementApi.UmUsersDelete(ctx, d.Id()).Execute()
	if apiResponse == nil || err != nil {
		/* //try again in 20 seconds
		time.Sleep(20 * time.Second)
		apiResponse, err := client.UserManagementApi.UmUsersDelete(ctx, d.Id()).Execute()
		if err != nil { */
		if apiResponse == nil || apiResponse.Response != nil && apiResponse.StatusCode != 404 {
			diags := diag.FromErr(fmt.Errorf("an error occured while deleting a user %s %s, %s", d.Id(), err, responseBody(apiResponse)))
			return diags
		}
		// }
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
