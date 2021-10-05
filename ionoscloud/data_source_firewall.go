package ionoscloud

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	ionoscloud "github.com/ionos-cloud/sdk-go/v6"
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

	if idOk {
		/* search by ID */
		firewall, _, err = client.FirewallRulesApi.DatacentersServersNicsFirewallrulesFindById(ctx, datacenterId, serverId, nicId, id.(string)).Execute()
		if err != nil {
			return diag.FromErr(fmt.Errorf("an error occurred while fetching the firewall rule %s: %s", id.(string), err))
		}
	} else {
		/* search by name */
		var firewalls ionoscloud.FirewallRules

		firewalls, _, err := client.FirewallRulesApi.DatacentersServersNicsFirewallrulesGet(ctx, datacenterId, serverId, nicId).Execute()

		if err != nil {
			return diag.FromErr(fmt.Errorf("an error occurred while fetching backup unit: %s", err.Error()))
		}

		if firewalls.Items != nil {
			for _, fr := range *firewalls.Items {
				tmpFirewall, _, err := client.FirewallRulesApi.DatacentersServersNicsFirewallrulesFindById(ctx, datacenterId, serverId, nicId, *fr.Id).Execute()
				if err != nil {
					return diag.FromErr(fmt.Errorf("an error occurred while fetching firewall rule with ID %s: %s", *fr.Id, err.Error()))
				}
				if tmpFirewall.Properties.Name != nil && *tmpFirewall.Properties.Name == name.(string) {
					firewall = tmpFirewall
					break
				}

			}
		}

	}

	if &firewall == nil {
		return diag.FromErr(fmt.Errorf("firewall rule not found"))
	}

	if err := d.Set("id", *firewall.Id); err != nil {
		return diag.FromErr(err)
	}

	if diags := setFirewallData(d, &firewall); diags != nil {
		return diags
	}

	return nil
}

func setFirewallData(d *schema.ResourceData, firewall *ionoscloud.FirewallRule) diag.Diagnostics {

	if firewall.Id != nil {
		d.SetId(*firewall.Id)
	}

	if firewall.Properties != nil {

		if firewall.Properties.Protocol != nil {
			err := d.Set("protocol", *firewall.Properties.Protocol)
			if err != nil {
				diags := diag.FromErr(fmt.Errorf("error while setting protocol property for firewall %s: %s", d.Id(), err))
				return diags
			}
		}

		if firewall.Properties.Name != nil {
			err := d.Set("name", *firewall.Properties.Name)
			if err != nil {
				diags := diag.FromErr(fmt.Errorf("error while setting name property for firewall %s: %s", d.Id(), err))
				return diags
			}
		}

		if firewall.Properties.SourceMac != nil {
			err := d.Set("source_mac", *firewall.Properties.SourceMac)
			if err != nil {
				diags := diag.FromErr(fmt.Errorf("error while setting source_mac property for firewall %s: %s", d.Id(), err))
				return diags
			}
		}

		if firewall.Properties.SourceIp != nil {
			err := d.Set("source_ip", *firewall.Properties.SourceIp)
			if err != nil {
				diags := diag.FromErr(fmt.Errorf("error while setting source_ip property for firewall %s: %s", d.Id(), err))
				return diags
			}
		}

		if firewall.Properties.TargetIp != nil {
			err := d.Set("target_ip", *firewall.Properties.TargetIp)
			if err != nil {
				diags := diag.FromErr(fmt.Errorf("error while setting target_ip property for firewall %s: %s", d.Id(), err))
				return diags
			}
		}

		if firewall.Properties.PortRangeStart != nil {
			err := d.Set("port_range_start", *firewall.Properties.PortRangeStart)
			if err != nil {
				diags := diag.FromErr(fmt.Errorf("error while setting port_range_start property for firewall %s: %s", d.Id(), err))
				return diags
			}
		}

		if firewall.Properties.PortRangeEnd != nil {
			err := d.Set("port_range_end", *firewall.Properties.PortRangeEnd)
			if err != nil {
				diags := diag.FromErr(fmt.Errorf("error while setting port_range_end property for firewall %s: %s", d.Id(), err))
				return diags
			}
		}

		if firewall.Properties.IcmpType != nil {
			err := d.Set("icmp_type", *firewall.Properties.IcmpType)
			if err != nil {
				diags := diag.FromErr(fmt.Errorf("error while setting icmp_type property for firewall %s: %s", d.Id(), err))
				return diags
			}
		}

		if firewall.Properties.IcmpCode != nil {
			err := d.Set("icmp_code", *firewall.Properties.IcmpCode)
			if err != nil {
				diags := diag.FromErr(fmt.Errorf("error while setting icmp_code property for firewall %s: %s", d.Id(), err))
				return diags
			}
		}
	}
	return nil
}
