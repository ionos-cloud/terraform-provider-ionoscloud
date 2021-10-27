package ionoscloud

import (
	"context"
	"errors"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	ionoscloud "github.com/ionos-cloud/sdk-go/v6"
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
			"user_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"users": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
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
	client := meta.(*ionoscloud.APIClient)

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
		group, _, err = client.UserManagementApi.UmGroupsFindById(ctx, id.(string)).Execute()
		if err != nil {
			diags := diag.FromErr(fmt.Errorf("an error occurred while fetching group with ID %s: %s", id.(string), err))
			return diags
		}
	} else {
		/* search by name */
		var groups ionoscloud.Groups

		groups, _, err := client.UserManagementApi.UmGroupsGet(ctx).Execute()
		if err != nil {
			diags := diag.FromErr(fmt.Errorf("an error occurred while fetching groups: %s", err.Error()))
			return diags
		}

		found := false
		if groups.Items != nil {
			for _, g := range *groups.Items {
				if g.Properties.Name != nil && *g.Properties.Name == name.(string) {
					/* group found */
					group, _, err = client.UserManagementApi.UmGroupsFindById(ctx, *g.Id).Execute()
					if err != nil {
						diags := diag.FromErr(fmt.Errorf("an error occurred while fetching group %s: %s", *g.Id, err))
						return diags
					}
					found = true
					break
				}
			}
		}

		if !found {
			diags := diag.FromErr(fmt.Errorf("group not found"))
			return diags
		}
	}

	if err = setGroupData(ctx, client, d, &group); err != nil {
		return diag.FromErr(err)
	}

	return nil
}
