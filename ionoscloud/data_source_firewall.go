package ionoscloud

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	ionoscloud "github.com/ionos-cloud/sdk-go/v5"
)

func dataSourceFirewall() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceFirewallRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"protocol": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"source_mac": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"source_ip": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"target_ip": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"port_range_start": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"port_range_end": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"icmp_type": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"icmp_code": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"datacenter_id": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.All(validation.StringIsNotWhiteSpace),
			},
			"server_id": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.All(validation.StringIsNotWhiteSpace),
			},
			"nic_id": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.All(validation.StringIsNotWhiteSpace),
			},
		},
		Timeouts: &resourceDefaultTimeouts,
	}
}

func dataSourceFirewallRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*ionoscloud.APIClient)

	datacenterId := d.Get("datacenter_id").(string)
	serverId := d.Get("server_id").(string)
	nicId := d.Get("nic_id").(string)

	id, idOk := d.GetOk("id")
	name, nameOk := d.GetOk("name")

	if idOk && nameOk {
		return diag.FromErr(fmt.Errorf("id and name cannot be both specified in the same time"))
	}
	if !idOk && !nameOk {
		return diag.FromErr(fmt.Errorf("please provide either the firewall rule id or name"))
	}
	var firewall ionoscloud.FirewallRule
	var err error
	var apiResponse *ionoscloud.APIResponse

	found := false

	if idOk {
		/* search by ID */
		firewall, apiResponse, err = client.NicApi.DatacentersServersNicsFirewallrulesFindById(ctx, datacenterId, serverId, nicId, id.(string)).Execute()
		logApiRequestTime(apiResponse)
		if err != nil {
			return diag.FromErr(fmt.Errorf("an error occurred while fetching the firewall rule %s: %s", id.(string), err))
		}
		found = true
	} else {
		/* search by name */
		var firewalls ionoscloud.FirewallRules

		firewalls, apiResponse, err := client.NicApi.DatacentersServersNicsFirewallrulesGet(ctx, datacenterId, serverId, nicId).Execute()
		logApiRequestTime(apiResponse)

		if err != nil {
			return diag.FromErr(fmt.Errorf("an error occurred while fetching backup unit: %s", err.Error()))
		}

		if firewalls.Items != nil {
			for _, fr := range *firewalls.Items {
				tmpFirewall, apiResponse, err := client.NicApi.DatacentersServersNicsFirewallrulesFindById(ctx, datacenterId, serverId, nicId, *fr.Id).Execute()
				logApiRequestTime(apiResponse)
				if err != nil {
					return diag.FromErr(fmt.Errorf("an error occurred while fetching firewall rule with ID %s: %s", *fr.Id, err.Error()))
				}
				if tmpFirewall.Properties.Name != nil && *tmpFirewall.Properties.Name == name.(string) {
					firewall = tmpFirewall
					found = true
					break
				}
			}
		}
	}

	if !found {
		return diag.FromErr(fmt.Errorf("firewall rule not found"))
	}

	if firewall.Id != nil {
		if err := d.Set("id", *firewall.Id); err != nil {
			return diag.FromErr(err)
		}
	}

	if diags := setFirewallData(d, &firewall); diags != nil {
		return diag.FromErr(err)
	}

	return nil
}
