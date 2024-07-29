package ionoscloud

import (
	"context"
	"errors"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	ionoscloud "github.com/ionos-cloud/sdk-go/v6"

	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services"
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
			"create_flow_log": {
				Type:        schema.TypeBool,
				Description: "Create Flow Logs privilege.",
				Computed:    true,
			},
			"access_and_manage_monitoring": {
				Type: schema.TypeBool,
				Description: "Privilege for a group to access and manage monitoring related functionality " +
					"(access metrics, CRUD on alarms, alarm-actions etc) using Monotoring-as-a-Service (MaaS).",
				Computed: true,
			},
			"access_and_manage_certificates": {
				Type:        schema.TypeBool,
				Description: "Privilege for a group to access and manage certificates.",
				Computed:    true,
			},
			"manage_dbaas": {
				Type:        schema.TypeBool,
				Description: "Privilege for a group to manage DBaaS related functionality",
				Computed:    true,
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
	client := meta.(services.SdkBundle).CloudApiClient

	id, idOk := d.GetOk("id")
	name, nameOk := d.GetOk("name")

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
		var apiResponse *ionoscloud.APIResponse
		group, apiResponse, err = client.UserManagementApi.UmGroupsFindById(ctx, id.(string)).Execute()
		logApiRequestTime(apiResponse)
		if err != nil {
			diags := diag.FromErr(fmt.Errorf("an error occurred while fetching group with ID %s: %w", id.(string), err))
			return diags
		}
	} else {
		/* search by name */
		var groups ionoscloud.Groups

		groups, apiResponse, err := client.UserManagementApi.UmGroupsGet(ctx).Depth(1).Execute()
		logApiRequestTime(apiResponse)
		if err != nil {
			diags := diag.FromErr(fmt.Errorf("an error occurred while fetching groups: %w", err))
			return diags
		}

		var results []ionoscloud.Group

		if groups.Items != nil {
			for _, g := range *groups.Items {
				if g.Properties != nil && g.Properties.Name != nil && *g.Properties.Name == name.(string) {
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

		if len(results) == 0 {
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
