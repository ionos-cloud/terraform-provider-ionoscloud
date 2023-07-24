package ionoscloud

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	ionoscloud "github.com/ionos-cloud/sdk-go/v6"
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
			"groups": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.All(validation.StringIsNotWhiteSpace),
						},
					},
				},
			},
		},
		Timeouts: &resourceDefaultTimeouts,
	}
}

func dataSourceUserRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(SdkBundle).CloudApiClient

	idValue, idOk := d.GetOk("id")
	emailValue, emailOk := d.GetOk("email")

	id := idValue.(string)
	email := emailValue.(string)

	if idOk && emailOk {
		diags := diag.FromErr(errors.New("id and email cannot be both specified in the same time"))
		return diags
	}

	if !idOk && !emailOk {
		config := client.GetConfig()
		email = config.Username
		if email == "" {
			diags := diag.FromErr(errors.New("please provide either the user id or email"))
			return diags
		}
		log.Printf("[INFO] email got from provider configuration since none was provided")
	}
	var user ionoscloud.User
	var err error
	var apiResponse *ionoscloud.APIResponse

	if idOk {
		/* search by ID */
		user, apiResponse, err = client.UserManagementApi.UmUsersFindById(ctx, id).Execute()
		logApiRequestTime(apiResponse)
		if err != nil {
			diags := diag.FromErr(fmt.Errorf("an error occurred while fetching user with ID %s: %w", id, err))
			return diags
		}
	} else {
		/* search by name */
		var users ionoscloud.Users

		users, apiResponse, err := client.UserManagementApi.UmUsersGet(ctx).Depth(1).Execute()
		logApiRequestTime(apiResponse)
		if err != nil {
			diags := diag.FromErr(fmt.Errorf("an error occurred while fetching users: %w", err))
			return diags
		}

		var results []ionoscloud.User
		if users.Items != nil {
			for _, u := range *users.Items {
				if u.Properties != nil && u.Properties.Email != nil && *u.Properties.Email == email {
					/* user found */
					user, apiResponse, err = client.UserManagementApi.UmUsersFindById(ctx, *u.Id).Execute()
					logApiRequestTime(apiResponse)
					if err != nil {
						diags := diag.FromErr(fmt.Errorf("an error occurred while fetching user %s: %w", *u.Id, err))
						return diags
					}
					results = append(results, user)
				}
			}
		}

		if results == nil || len(results) == 0 {
			return diag.FromErr(fmt.Errorf("no user found with the specified criteria: email = %s", email))
		} else {
			user = results[0]
		}
	}
	if err = setUsersForGroup(ctx, d, &user, *client); err != nil {
		return diag.FromErr(err)
	}

	if err = setUserData(d, &user); err != nil {
		return diag.FromErr(err)
	}

	return nil
}
func setUsersForGroup(ctx context.Context, d *schema.ResourceData, user *ionoscloud.User, client ionoscloud.APIClient) error {
	if user == nil {
		return fmt.Errorf("did not expect empty user")
	}

	groups, apiResponse, err := client.UserManagementApi.UmUsersGroupsGet(ctx, *user.Id).Depth(1).Execute()
	logApiRequestTime(apiResponse)
	if err != nil {
		return fmt.Errorf("an error occured while executing UmUsersGroupsGet %s (%w)", *user.Id, err)
	}

	groupEntries := make([]interface{}, 0)
	if groups.Items != nil && len(*groups.Items) > 0 {
		groupEntries = make([]interface{}, len(*groups.Items))
		for groupIndex, group := range *groups.Items {
			groupEntry := make(map[string]interface{})

			if group.Id != nil {
				groupEntry["id"] = *group.Id
			}

			if group.Properties != nil && group.Properties.Name != nil {
				groupEntry["name"] = group.Properties.Name
			}
			groupEntries[groupIndex] = groupEntry
		}

		if len(groupEntries) > 0 {
			if err := d.Set("groups", groupEntries); err != nil {
				return fmt.Errorf("error while setting groups for user (%w)", err)
			}
		}
	}

	return nil
}
