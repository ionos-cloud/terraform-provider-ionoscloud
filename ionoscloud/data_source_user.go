package ionoscloud

import (
	"context"
	"errors"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	ionoscloud "github.com/ionos-cloud/sdk-go/v6"
	"log"
)

func dataSourceUser() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceUserRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"first_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"last_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"email": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"administrator": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"force_sec_auth": {
				Type:     schema.TypeBool,
				Computed: true,
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
				Computed: true,
			},
		},
		Timeouts: &resourceDefaultTimeouts,
	}
}

func dataSourceUserRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(SdkBundle).CloudApiClient

	id, idOk := d.GetOk("id")
	email, emailOk := d.GetOk("email")

	if idOk && emailOk {
		diags := diag.FromErr(errors.New("id and email cannot be both specified in the same time"))
		return diags
	}
	if !idOk && !emailOk {
		diags := diag.FromErr(errors.New("please provide either the user id or email"))
		return diags
	}
	var user ionoscloud.User
	var err error
	var apiResponse *ionoscloud.APIResponse

	if idOk {
		/* search by ID */
		user, apiResponse, err = client.UserManagementApi.UmUsersFindById(ctx, id.(string)).Execute()
		logApiRequestTime(apiResponse)
		if err != nil {
			diags := diag.FromErr(fmt.Errorf("an error occurred while fetching user with ID %s: %w", id.(string), err))
			return diags
		}
	} else {
		/* search by name */
		users, apiResponse, err := client.UserManagementApi.UmUsersGet(ctx).Depth(1).Filter("email", email.(string)).Execute()
		logApiRequestTime(apiResponse)
		if err != nil {
			diags := diag.FromErr(fmt.Errorf("an error occurred while fetching users: %s", err.Error()))
			return diags
		}

		if users.Items != nil && len(*users.Items) > 0 {
			user = (*users.Items)[len(*users.Items)-1]
			log.Printf("[INFO] %v users found matching the search criteria. Getting the latest user from the list %v", len(*users.Items), *user.Id)
		} else {
			return diag.FromErr(fmt.Errorf("no user found with the specified email %s", email.(string)))
		}
	}

	if err = setUserData(d, &user); err != nil {
		return diag.FromErr(err)
	}

	return nil
}
