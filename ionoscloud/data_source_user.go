package ionoscloud

import (
	"context"
	"errors"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	ionoscloud "github.com/ionos-cloud/sdk-go/v6"
	"log"
	"strings"
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
				Optional: true,
			},
			"last_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"email": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"administrator": {
				Type:     schema.TypeBool,
				Optional: true,
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
				Optional: true,
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

	id, idOk := d.GetOk("id")
	email, emailOk := d.GetOk("email")
	firstNameValue, firstNameOk := d.GetOk("first_name")
	lastNameValue, lastNameOk := d.GetOk("last_name")
	s3CanonicalIdValue, s3CanonicalIdOk := d.GetOk("s3_canonical_user_id")
	administratorValue, administratorOk := d.GetOk("administrator")
	// todo active flags also, but I don't know where are they

	firstName := firstNameValue.(string)
	lastName := lastNameValue.(string)
	s3CanonicalId := s3CanonicalIdValue.(string)
	administrator := administratorValue.(bool)

	if idOk && (emailOk || firstNameOk || lastNameOk || s3CanonicalIdOk || administratorOk) {
		diags := diag.FromErr(errors.New("id and other lookup parameter cannot be both specified in the same time"))
		return diags
	}
	if !idOk && !emailOk && !firstNameOk && !lastNameOk && !s3CanonicalIdOk && !administratorOk {
		diags := diag.FromErr(errors.New("please provide either the user id or other lookup parameter, like email or first_name"))
		log.Printf("ADMINISTRATOROK = %t", administratorOk)
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
		var users ionoscloud.Users

		users, apiResponse, err := client.UserManagementApi.UmUsersGet(ctx).Depth(1).Execute()
		logApiRequestTime(apiResponse)
		if err != nil {
			diags := diag.FromErr(fmt.Errorf("an error occurred while fetching users: %w", err))
			return diags
		}

		if users.Items == nil {
			diags := diag.FromErr(fmt.Errorf("no users found"))
			return diags
		}

		var results = *users.Items

		if emailOk {
			var emailResults []ionoscloud.User
			for _, u := range results {
				if u.Properties != nil && u.Properties.Email != nil && *u.Properties.Email == email.(string) {
					/* user found */
					user, apiResponse, err = client.UserManagementApi.UmUsersFindById(ctx, *u.Id).Execute()
					logApiRequestTime(apiResponse)
					if err != nil {
						diags := diag.FromErr(fmt.Errorf("an error occurred while fetching user %s: %w", *u.Id, err))
						return diags
					}
					emailResults = append(emailResults, user)
				}
			}
			if emailResults == nil {
				return diag.FromErr(fmt.Errorf("no user found with the specified criteria: email = %s", email))
			}
			results = emailResults

		}

		if firstNameOk && firstName != "" {
			var firstNameResults []ionoscloud.User
			if results != nil {
				for _, user := range results {
					if user.Properties != nil && user.Properties.Firstname != nil && strings.EqualFold(*user.Properties.Firstname, firstName) {
						u, apiResponse, err := client.UserManagementApi.UmUsersFindById(ctx, *user.Id).Execute()
						logApiRequestTime(apiResponse)
						if err != nil {
							diags := diag.FromErr(fmt.Errorf("an error occurred while fetching user %s: %w", *u.Id, err))
							return diags
						}
						firstNameResults = append(firstNameResults, u)
					}
				}
			}
			if firstNameResults == nil {
				return diag.FromErr(fmt.Errorf("no user found with the specified criteria: first name = %s", firstName))
			}
			results = firstNameResults

		}

		if lastNameOk && lastName != "" {
			var lastNameResults []ionoscloud.User
			if results != nil {
				for _, user := range results {
					if user.Properties != nil && user.Properties.Lastname != nil && strings.EqualFold(*user.Properties.Lastname, lastName) {
						lastNameResults = append(lastNameResults, user)
					}
				}
			}
			if lastNameResults == nil {
				return diag.FromErr(fmt.Errorf("no user found with the specified criteria: last name = %s", lastName))
			}
			results = lastNameResults
		}

		if s3CanonicalIdOk && s3CanonicalId != "" {
			var s3CanonicalIdResults []ionoscloud.User
			if results != nil {
				for _, user := range results {
					if user.Properties != nil && user.Properties.S3CanonicalUserId != nil && strings.EqualFold(*user.Properties.S3CanonicalUserId, s3CanonicalId) {
						s3CanonicalIdResults = append(s3CanonicalIdResults, user)
					}
				}
			}
			if s3CanonicalIdResults == nil {
				return diag.FromErr(fmt.Errorf("no user found with the specified criteria: s3 canonical id = %s", s3CanonicalId))
			}
			results = s3CanonicalIdResults
		}

		if administratorOk {
			var administratorResults []ionoscloud.User
			if results != nil {
				for _, user := range results {
					if user.Properties != nil && user.Properties.Administrator != nil && *user.Properties.Administrator == administrator {
						administratorResults = append(administratorResults, user)
					}
				}
			}
			if administratorResults == nil {
				return diag.FromErr(fmt.Errorf("no user found with the specified criteria: administrator = %t", administrator))
			}
			results = administratorResults
		}

		if results == nil || len(results) == 0 {
			return diag.FromErr(fmt.Errorf("no user found with the specified criteria: email = %s", email.(string)))
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
