package ionoscloud

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	ionoscloud "github.com/ionos-cloud/sdk-go/v6"

	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/bundleclient"
	diagutil "github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils/diags"
)

func dataSourceUser() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceUserRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
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
							Type:             schema.TypeString,
							Required:         true,
							ValidateDiagFunc: validation.ToDiagFunc(validation.StringIsNotWhiteSpace),
						},
					},
				},
			},
		},
		Timeouts: &resourceDefaultTimeouts,
	}
}

func dataSourceUserRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, err := meta.(bundleclient.SdkBundle).NewCloudAPIClientWithFailover(ctx)
	if err != nil {
		return diag.FromErr(err)
	}

	idValue, idOk := d.GetOk("id")
	emailValue, emailOk := d.GetOk("email")

	id := idValue.(string)
	email := emailValue.(string)

	if idOk && emailOk {
		return diagutil.ToDiags(d, fmt.Errorf("id and email cannot be both specified in the same time"), nil)
	}

	if !idOk && !emailOk {
		config := client.GetConfig()
		email = config.Username
		if email == "" {
			return diagutil.ToDiags(d, fmt.Errorf("please provide either the user id or email"), nil)
		}
		tflog.Info(ctx, "email taken from provider configuration since none was provided")
	}
	var user ionoscloud.User
	var apiResponse *ionoscloud.APIResponse

	if idOk {
		/* search by ID */
		user, apiResponse, err = client.UserManagementApi.UmUsersFindById(ctx, id).Execute()
		logApiRequestTime(apiResponse)
		if err != nil {
			return diagutil.ToDiags(d, fmt.Errorf("an error occurred while fetching user with ID %s: %w", id, err), &diagutil.ErrorContext{StatusCode: apiResponse.SafeStatusCode()})
		}
	} else {
		/* search by email */
		users, apiResponse, err := client.UserManagementApi.UmUsersGet(ctx).Depth(1).Filter("email", email).Execute()
		logApiRequestTime(apiResponse)
		if err != nil {
			return diagutil.ToDiags(d, fmt.Errorf("an error occurred while fetching users: %w", err), &diagutil.ErrorContext{StatusCode: apiResponse.SafeStatusCode()})
		}
		if users.Items == nil || len(*users.Items) == 0 {
			return diagutil.ToDiags(d, fmt.Errorf("no user found with the specified criteria: email = %s", email), nil)
		} else if len(*users.Items) > 1 {
			return diagutil.ToDiags(d, fmt.Errorf("multiple users found with the specified criteria: email = %s", email), nil)
		}
		user = (*users.Items)[0]
	}
	if err = setUsersForGroup(ctx, d, &user, *client); err != nil {
		return diagutil.ToDiags(d, err, nil)
	}

	if err = setUserData(d, &user); err != nil {
		return diagutil.ToDiags(d, err, nil)
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
		return fmt.Errorf("an error occurred while executing UmUsersGroupsGet %s (%w)", *user.Id, err)
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
