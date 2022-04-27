package ionoscloud

import (
	"context"
	"errors"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	ionoscloud "github.com/ionos-cloud/sdk-go/v6"
	"log"
	"strings"
)

func dataSourceGroup() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceGroupRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"partial_match": {
				Type:        schema.TypeBool,
				Description: "Whether partial matching is allowed or not when using name argument.",
				Default:     false,
				Optional:    true,
			},
			"create_datacenter": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"create_snapshot": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"reserve_ip": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"access_activity_log": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"create_pcc": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"s3_privilege": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"create_backup_unit": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"create_internet_access": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"create_k8s_cluster": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"users": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
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
							Computed: true,
						},
						"administrator": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"force_sec_auth": {
							Type:     schema.TypeBool,
							Computed: true,
						},
					},
				},
			},
		},
		Timeouts: &resourceDefaultTimeouts,
	}
}

func dataSourceGroupRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(SdkBundle).CloudApiClient

	idValue, idOk := d.GetOk("id")
	nameValue, nameOk := d.GetOk("name")

	id := idValue.(string)
	name := nameValue.(string)

	if idOk && nameOk {
		diags := diag.FromErr(errors.New("id and name cannot be both specified in the same time"))
		return diags
	}
	if !idOk && !nameOk {
		diags := diag.FromErr(errors.New("please provide either the group id or name"))
		return diags
	}
	var group ionoscloud.Group
	var err error

	if idOk {
		/* search by ID */
		log.Printf("[INFO] Using data source for group by id %s", id)

		var apiResponse *ionoscloud.APIResponse
		group, apiResponse, err = client.UserManagementApi.UmGroupsFindById(ctx, id).Execute()
		logApiRequestTime(apiResponse)
		if err != nil {
			diags := diag.FromErr(fmt.Errorf("an error occurred while fetching group with ID %s: %w", id, err))
			return diags
		}
	} else {
		/* search by name */
		var results []ionoscloud.Group

		partialMatch := d.Get("partial_match").(bool)

		log.Printf("[INFO] Using data source for group by name with partial_match %t and name: %s", partialMatch, name)

		if partialMatch {
			groups, apiResponse, err := client.UserManagementApi.UmGroupsGet(ctx).Depth(1).Filter("name", name).Execute()
			logApiRequestTime(apiResponse)
			if err != nil {
				diags := diag.FromErr(fmt.Errorf("an error occurred while fetching groups: %w", err))
				return diags
			}
			results = *groups.Items
		} else {
			groups, apiResponse, err := client.UserManagementApi.UmGroupsGet(ctx).Depth(1).Execute()
			logApiRequestTime(apiResponse)
			if err != nil {
				diags := diag.FromErr(fmt.Errorf("an error occurred while fetching groups: %w", err))
				return diags
			}

			if groups.Items != nil {
				for _, g := range *groups.Items {
					if g.Properties != nil && g.Properties.Name != nil && strings.EqualFold(*g.Properties.Name, name) {
						/* group found */
						group, apiResponse, err = client.UserManagementApi.UmGroupsFindById(ctx, *g.Id).Execute()
						logApiRequestTime(apiResponse)
						if err != nil {
							diags := diag.FromErr(fmt.Errorf("an error occurred while fetching group %s: %w", *g.Id, err))
							return diags
						}
						results = append(results, g)
					}
				}
			}
		}

		if results == nil || len(results) == 0 {
			return diag.FromErr(fmt.Errorf("no group found with the specified name = %s", name))
		} else if len(results) > 1 {
			return diag.FromErr(fmt.Errorf("more than one group found with the specified criteria name = %s", name))
		} else {
			group = results[0]
		}

	}

	if err = setGroupData(ctx, client, d, &group); err != nil {
		return diag.FromErr(err)
	}

	return nil
}
