package ionoscloud

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	ionoscloud "github.com/ionos-cloud/sdk-go/v6"

	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services"
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
				Type:     schema.TypeString,
				Computed: true,
			},
			"icmp_code": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"datacenter_id": {
				Type:             schema.TypeString,
				Required:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringIsNotWhiteSpace),
			},
			"server_id": {
				Type:             schema.TypeString,
				Required:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringIsNotWhiteSpace),
			},
			"nic_id": {
				Type:             schema.TypeString,
				Required:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringIsNotWhiteSpace),
			},
		},
		Timeouts: &resourceDefaultTimeouts,
	}
}

func dataSourceFirewallRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(services.SdkBundle).CloudAPIClient

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

	if idOk {
		/* search by ID */
		firewall, apiResponse, err = client.FirewallRulesApi.DatacentersServersNicsFirewallrulesFindById(ctx, datacenterId, serverId, nicId, id.(string)).Execute()
		logApiRequestTime(apiResponse)
		if err != nil {
			return diag.FromErr(fmt.Errorf("an error occurred while fetching the firewall rule %s: %w", id.(string), err))
		}
	} else {
		/* search by name */
		var firewalls ionoscloud.FirewallRules

		firewalls, apiResponse, err := client.FirewallRulesApi.DatacentersServersNicsFirewallrulesGet(ctx, datacenterId, serverId, nicId).Depth(1).Execute()
		logApiRequestTime(apiResponse)

		if err != nil {
			return diag.FromErr(fmt.Errorf("an error occurred while fetching backup unit: %w", err))
		}

		var results []ionoscloud.FirewallRule

		if firewalls.Items != nil {
			for _, fr := range *firewalls.Items {
				if fr.Properties != nil && fr.Properties.Name != nil && *fr.Properties.Name == name.(string) {
					tmpFirewall, apiResponse, err := client.FirewallRulesApi.DatacentersServersNicsFirewallrulesFindById(ctx, datacenterId, serverId, nicId, *fr.Id).Execute()
					logApiRequestTime(apiResponse)
					if err != nil {
						return diag.FromErr(fmt.Errorf("an error occurred while fetching firewall rule with ID %s: %w", *fr.Id, err))
					}
					results = append(results, tmpFirewall)
				}
			}
		}

		if results == nil || len(results) == 0 {
			return diag.FromErr(fmt.Errorf("no firewall rule found with the specified name = %s", name))
		} else if len(results) > 1 {
			return diag.FromErr(fmt.Errorf("more than one firewall rule found with the specified criteria name = %s", name))
		} else {
			firewall = results[0]
		}

	}

	if err := setFirewallData(d, &firewall); err != nil {
		return diag.FromErr(err)
	}

	return nil
}
