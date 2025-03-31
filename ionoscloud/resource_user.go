package ionoscloud

import (
	"context"
	"fmt"
	"github.com/ionos-cloud/sdk-go-bundle/shared"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	ionoscloud "github.com/ionos-cloud/sdk-go/v6"

	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/internal/tf/writeonly"
	bundleclient "github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/bundleclient"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/slice"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils"
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
				Type:             schema.TypeString,
				Required:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringIsNotWhiteSpace),
			},
			"last_name": {
				Type:             schema.TypeString,
				Required:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringIsNotWhiteSpace),
			},
			"email": {
				Type:             schema.TypeString,
				Required:         true,
				Description:      "Email address of the user",
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringIsNotWhiteSpace),
				DiffSuppressFunc: utils.DiffToLower,
			},
			"password": {
				Type:             schema.TypeString,
				Optional:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringIsNotWhiteSpace),
				Sensitive:        true,
				ConflictsWith:    []string{"password_wo"},
				ExactlyOneOf:     []string{"password", "password_wo"},
			},
			"password_wo": {
				Type:          schema.TypeString,
				Optional:      true,
				WriteOnly:     true,
				Sensitive:     true,
				Description:   "Write-only attribute. Password for the user. To modify, must change the password_wo_version attribute.",
				ConflictsWith: []string{"password"},
				ExactlyOneOf:  []string{"password", "password_wo"},
				RequiredWith:  []string{"password_wo_version"},
			},
			"password_wo_version": {
				Type:         schema.TypeInt,
				Optional:     true,
				RequiredWith: []string{"password_wo"},
				Description:  "Version of the password_wo attribute. Must be incremented to apply changes to the password_wo attribute.",
				ValidateFunc: validation.IntAtLeast(1),
			},
			"administrator": {
				Type:        schema.TypeBool,
				Description: "Indicates if the user has administrative rights. Administrators do not need to be managed in groups, as they automatically have access to all resources associated with the contract.",
				Optional:    true,
			},
			"force_sec_auth": {
				Type:        schema.TypeBool,
				Description: "Indicates if secure (two-factor) authentication is forced for the user",
				Optional:    true,
			},
			"sec_auth_active": {
				Type:        schema.TypeBool,
				Description: "Indicates if secure (two-factor) authentication is active for the user. It can not be used in create requests - can be used in update.",
				Computed:    true,
			},
			"s3_canonical_user_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"active": {
				Type:        schema.TypeBool,
				Description: "Indicates if the user is active",
				Optional:    true,
			},
			"group_ids": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type:             schema.TypeString,
					ValidateDiagFunc: validation.ToDiagFunc(validation.IsUUID),
				},
				Description: "Ids of the groups that the user is a member of",
			},
		},
		Timeouts: &resourceDefaultTimeouts,
	}
}

func resourceUserCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(bundleclient.SdkBundle).CloudApiClient
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
	if _, isSet := d.GetOk("password"); isSet {
		password := d.Get("password").(string)
		request.Properties.Password = &password
	} else {
		// get write-only value from configuration
		passwordWO, err := writeonly.GetStringValue(d, "password_wo")
		if err != nil {
			return diag.FromErr(err)
		}
		request.Properties.Password = shared.ToPtr(passwordWO)
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
		diags := diag.FromErr(fmt.Errorf("an error occurred while creating a user: %w", err))
		return diags
	}

	d.SetId(*rsp.Id)

	if errState := bundleclient.WaitForStateChange(ctx, meta, d, apiResponse, schema.TimeoutCreate); errState != nil {
		if bundleclient.IsRequestFailed(errState) {
			d.SetId("")
		}
		return diag.FromErr(errState)
	}

	// Add the user to the specified groups, if any.
	if groupsVal, groupsOk := d.GetOk("group_ids"); groupsOk {
		groupsList := groupsVal.(*schema.Set).List()
		log.Printf("[INFO] Adding group_ids %+v ", groupsList)
		if groupsList != nil {
			for _, groupsItem := range groupsList {
				groupId := groupsItem.(string)
				if err := addUserToGroup(d.Id(), groupId, ctx, d, meta); err != nil {
					return diag.FromErr(err)
				}
			}
		}
	}

	return resourceUserRead(ctx, d, meta)
}

func resourceUserRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(bundleclient.SdkBundle).CloudApiClient

	user, apiResponse, err := client.UserManagementApi.UmUsersFindById(ctx, d.Id()).Depth(1).Execute()
	logApiRequestTime(apiResponse)

	if err != nil {
		if apiResponse.HttpNotFound() {
			d.SetId("")
			return nil
		}
		diags := diag.FromErr(fmt.Errorf("an error occurred while fetching a User ID %s %w", d.Id(), err))
		return diags
	}

	if err = setUserData(d, &user); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceUserUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(bundleclient.SdkBundle).CloudApiClient

	foundUser, apiResponse, err := client.UserManagementApi.UmUsersFindById(ctx, d.Id()).Execute()
	logApiRequestTime(apiResponse)
	if err != nil {
		diags := diag.FromErr(fmt.Errorf("an error occurred while fetching a User ID %s %w", d.Id(), err))
		return diags
	}

	userReq := ionoscloud.UserPut{
		Properties: &ionoscloud.UserPropertiesPut{},
	}
	// this is a PUT, so we first fill everything, then we change what has been modified
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

	if d.HasChange("password") {
		password := d.Get("password").(string)
		userReq.Properties.Password = &password
	} else if d.HasChange("password_wo_version") {
		passwordWO, err := writeonly.GetStringValue(d, "password_wo")
		if err != nil {
			return diag.FromErr(err)
		}
		userReq.Properties.Password = shared.ToPtr(passwordWO)
	}

	if d.HasChange("group_ids") {
		oldValues, newValues := d.GetChange("group_ids")
		oldGroupsList := slice.AnyToString(oldValues.(*schema.Set).List())
		newGroupsList := slice.AnyToString(newValues.(*schema.Set).List())

		newGroups := utils.DiffSliceOneWay(newGroupsList, oldGroupsList)
		deletedGroups := utils.DiffSliceOneWay(oldGroupsList, newGroupsList)

		if newGroups != nil && len(newGroups) > 0 {
			log.Printf("[INFO] New groups to add: %+v", newGroups)
			for _, groupId := range newGroups {
				if err := addUserToGroup(d.Id(), groupId, ctx, d, meta); err != nil {
					return diag.FromErr(err)
				}
			}
		}

		if deletedGroups != nil && len(deletedGroups) > 0 {
			log.Printf("[INFO] Groups to delete: %+v", deletedGroups)
			for _, groupId := range deletedGroups {
				if err := deleteUserFromGroup(d.Id(), groupId, ctx, d, meta); err != nil {
					return diag.FromErr(err)
				}
			}
		}
	}

	_, apiResponse, err = client.UserManagementApi.UmUsersPut(ctx, d.Id()).User(userReq).Execute()
	logApiRequestTime(apiResponse)
	if err != nil {
		diags := diag.FromErr(fmt.Errorf("an error occurred while patching a user ID %s %w", d.Id(), err))
		return diags
	}

	if errState := bundleclient.WaitForStateChange(ctx, meta, d, apiResponse, schema.TimeoutUpdate); errState != nil {
		return diag.FromErr(errState)
	}

	return resourceUserRead(ctx, d, meta)
}

func resourceUserDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(bundleclient.SdkBundle).CloudApiClient

	apiResponse, err := client.UserManagementApi.UmUsersDelete(ctx, d.Id()).Execute()
	logApiRequestTime(apiResponse)
	if err != nil {
		diags := diag.FromErr(err)
		return diags
	}

	if errState := bundleclient.WaitForStateChange(ctx, meta, d, apiResponse, schema.TimeoutDelete); errState != nil {
		return diag.FromErr(errState)
	}

	d.SetId("")
	return nil

}

func resourceUserImporter(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	client := meta.(bundleclient.SdkBundle).CloudApiClient

	userId := d.Id()

	user, apiResponse, err := client.UserManagementApi.UmUsersFindById(ctx, userId).Depth(1).Execute()
	logApiRequestTime(apiResponse)

	if err != nil {
		if apiResponse.HttpNotFound() {
			d.SetId("")
			return nil, fmt.Errorf("user does not exist%q", userId)
		}
		return nil, fmt.Errorf("an error occurred while trying to fetch the user %q, error:%w", userId, err)

	}

	if err = setUserData(d, &user); err != nil {
		return nil, err
	}

	if err = d.Set("group_ids", getUserGroups(&user)); err != nil {
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

func getUserGroups(user *ionoscloud.User) []string {
	var groupIDs []string
	if user.Entities == nil {
		return groupIDs
	}

	if !user.Entities.HasGroups() {
		return groupIDs
	}

	if user.Entities.Groups.Items == nil {
		return groupIDs
	}

	groups := *user.Entities.Groups.Items
	for _, g := range groups {
		if g.Id != nil {
			groupIDs = append(groupIDs, *g.Id)
		}
	}

	return groupIDs
}
