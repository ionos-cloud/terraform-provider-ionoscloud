package ionoscloud

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	ionoscloud "github.com/ionos-cloud/sdk-go/v6"
	"log"
	"strings"
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
			"partial_match": {
				Type:        schema.TypeBool,
				Description: "Whether partial matching is allowed or not when using name argument.",
				Default:     false,
				Optional:    true,
			},
			"protocol": {
				Type:     schema.TypeString,
				Optional: true,
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
				Optional: true,
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
	client := meta.(SdkBundle).CloudApiClient

	datacenterId := d.Get("datacenter_id").(string)
	serverId := d.Get("server_id").(string)
	nicId := d.Get("nic_id").(string)

	idValue, idOk := d.GetOk("id")
	nameValue, nameOk := d.GetOk("name")
	firewallTypeValue, firewallTypeOk := d.GetOk("type")
	protocolValue, protocolOk := d.GetOk("protocol")

	id := idValue.(string)
	name := nameValue.(string)
	firewallType := firewallTypeValue.(string)
	protocol := protocolValue.(string)

	if idOk && (nameOk || firewallTypeOk || protocolOk) {
		return diag.FromErr(fmt.Errorf("id and name/type/protocol cannot be both specified in the same time, choose between id or a combination of other parameters"))
	}
	if !idOk && !nameOk && !firewallTypeOk && !protocolOk {
		return diag.FromErr(fmt.Errorf("please provide either the firewall rule id or other parameter like name, type or protocol"))
	}
	var firewall ionoscloud.FirewallRule
	var err error
	var apiResponse *ionoscloud.APIResponse

	if idOk {
		/* search by ID */
		log.Printf("[INFO] Using data source for firewall rule by id %s", id)
		firewall, apiResponse, err = client.FirewallRulesApi.DatacentersServersNicsFirewallrulesFindById(ctx, datacenterId, serverId, nicId, id).Execute()
		logApiRequestTime(apiResponse)
		if err != nil {
			return diag.FromErr(fmt.Errorf("an error occurred while fetching the firewall rule %s: %s", id, err))
		}
	} else {
		/* search by name */
		var results []ionoscloud.FirewallRule

		if nameOk {
			partialMatch := d.Get("partial_match").(bool)

			log.Printf("[INFO] Using data source for firewall by name with partial_match %t and name: %s", partialMatch, name)

			if partialMatch {
				firewalls, apiResponse, err := client.FirewallRulesApi.DatacentersServersNicsFirewallrulesGet(ctx, datacenterId, serverId, nicId).Depth(1).Filter("name", name).Execute()
				logApiRequestTime(apiResponse)

				if err != nil {
					return diag.FromErr(fmt.Errorf("an error occurred while fetching firewall rules while searching by partial name: %s, %w", name, err))
				}
				if len(*firewalls.Items) == 0 {
					return diag.FromErr(fmt.Errorf("no result found with the specified criteria: name with partial match: %s", name))
				}
				results = *firewalls.Items
			} else {
				firewalls, apiResponse, err := client.FirewallRulesApi.DatacentersServersNicsFirewallrulesGet(ctx, datacenterId, serverId, nicId).Depth(1).Execute()
				logApiRequestTime(apiResponse)

				if err != nil {
					return diag.FromErr(fmt.Errorf("an error occurred while fetching  firewall rule: %w", err))
				}

				if firewalls.Items != nil && nameOk {
					var nameResults []ionoscloud.FirewallRule
					for _, fr := range *firewalls.Items {
						if fr.Properties != nil && fr.Properties.Name != nil && strings.EqualFold(*fr.Properties.Name, name) {
							tmpFirewall, apiResponse, err := client.FirewallRulesApi.DatacentersServersNicsFirewallrulesFindById(ctx, datacenterId, serverId, nicId, *fr.Id).Execute()
							logApiRequestTime(apiResponse)
							if err != nil {
								return diag.FromErr(fmt.Errorf("an error occurred while fetching firewall rule with ID %s: %w", *fr.Id, err))
							}
							nameResults = append(nameResults, tmpFirewall)
						}
					}
					if len(nameResults) == 0 {
						return diag.FromErr(fmt.Errorf("no result found with the specified criteria: name %s", name))
					}
					results = nameResults
				}
			}

		} else {
			firewalls, apiResponse, err := client.FirewallRulesApi.DatacentersServersNicsFirewallrulesGet(ctx, datacenterId, serverId, nicId).Depth(1).Execute()
			logApiRequestTime(apiResponse)
			if err != nil {
				diags := diag.FromErr(fmt.Errorf("an error occurred while fetching firewalls: %s", err.Error()))
				return diags
			}
			results = *firewalls.Items
		}

		if firewallTypeOk && firewallType != "" {
			var firewallTypeResults []ionoscloud.FirewallRule
			if results != nil {
				for _, firewall := range results {
					if firewall.Properties != nil && firewall.Properties.Type != nil && strings.EqualFold(*firewall.Properties.Type, firewallType) {
						firewallTypeResults = append(firewallTypeResults, firewall)
					}
				}
			}
			if firewallTypeResults == nil || len(firewallTypeResults) == 0 {
				return diag.FromErr(fmt.Errorf("no result found with the specified criteria: type = %s", firewallType))
			}
			results = firewallTypeResults
		}

		if protocolOk && protocol != "" {
			var protocolResults []ionoscloud.FirewallRule
			if results != nil {
				for _, firewall := range results {
					if firewall.Properties != nil && firewall.Properties.Protocol != nil && strings.EqualFold(*firewall.Properties.Protocol, protocol) {
						protocolResults = append(protocolResults, firewall)
					}
				}
			}
			if protocolResults == nil || len(protocolResults) == 0 {
				return diag.FromErr(fmt.Errorf("no result found with the specified criteria: protocol = %s", protocol))
			}
			results = protocolResults
		}

		if len(results) > 1 {
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
