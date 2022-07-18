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

func dataSourceNatGatewayRule() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceNatGatewayRuleRead,
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
			"type": {
				Type:        schema.TypeString,
				Description: "Type of the NAT gateway rule.",
				Computed:    true,
			},
			"protocol": {
				Type: schema.TypeString,
				Description: "Protocol of the NAT gateway rule. Defaults to ALL. If protocol is 'ICMP' then " +
					"targetPortRange start and end cannot be set.",
				Optional: true,
			},
			"source_subnet": {
				Type: schema.TypeString,
				Description: "Source subnet of the NAT gateway rule. For SNAT rules it specifies which packets this " +
					"translation rule applies to based on the packets source IP address.",
				Computed: true,
			},
			"public_ip": {
				Type: schema.TypeString,
				Description: "Public IP address of the NAT gateway rule. Specifies the address used for masking outgoing " +
					"packets source address field. Should be one of the customer reserved IP address already " +
					"configured on the NAT gateway resource",
				Computed: true,
			},
			"target_subnet": {
				Type: schema.TypeString,
				Description: "Target or destination subnet of the NAT gateway rule. For SNAT rules it specifies which " +
					"packets this translation rule applies to based on the packets destination IP address. If " +
					"none is provided, rule will match any address.",
				Computed: true,
			},
			"target_port_range": {
				Type: schema.TypeList,
				Description: "Target port range of the NAT gateway rule. For SNAT rules it specifies which packets this " +
					"translation rule applies to based on destination port. If none is provided, rule will " +
					"match any port",
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"start": {
							Type:        schema.TypeInt,
							Description: "Target port range start associated with the NAT gateway rule.",
							Computed:    true,
						},
						"end": {
							Type:        schema.TypeInt,
							Description: "Target port range end associated with the NAT gateway rule.",
							Computed:    true,
						},
					},
				},
			},
			"datacenter_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"natgateway_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
		Timeouts: &resourceDefaultTimeouts,
	}
}

func dataSourceNatGatewayRuleRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(SdkBundle).CloudApiClient

	datacenterId := d.Get("datacenter_id").(string)

	natgatewayId := d.Get("natgateway_id").(string)

	idValue, idOk := d.GetOk("id")
	nameValue, nameOk := d.GetOk("name")
	protocolValue, protocolOk := d.GetOk("protocol")

	id := idValue.(string)
	name := nameValue.(string)
	protocol := protocolValue.(string)

	if idOk && (nameOk || protocolOk) {
		return diag.FromErr(errors.New("id and name/protocol cannot be both specified in the same time"))
	}
	if !idOk && !nameOk && !protocolOk {
		return diag.FromErr(errors.New("please provide either the lan id or other parameter like name or protocol"))
	}

	var natGatewayRule ionoscloud.NatGatewayRule
	var err error
	var apiResponse *ionoscloud.APIResponse

	if idOk {
		/* search by ID */
		log.Printf("[INFO] Using data source for nat gateway rule by id %s", id)
		natGatewayRule, apiResponse, err = client.NATGatewaysApi.DatacentersNatgatewaysRulesFindByNatGatewayRuleId(ctx, datacenterId, natgatewayId, id).Execute()
		logApiRequestTime(apiResponse)
		if err != nil {
			return diag.FromErr(fmt.Errorf("an error occurred while fetching the nat gateway rule %s: %s", id, err))
		}
	} else {
		/* search by name */
		var results []ionoscloud.NatGatewayRule

		if nameOk {
			//natGatewayRules, apiResponse, err := client.NATGatewaysApi.DatacentersNatgatewaysRulesGet(ctx, datacenterId, natgatewayId).Depth(1).Execute()
			//logApiRequestTime(apiResponse)
			//if err != nil {
			//	return diag.FromErr(fmt.Errorf("an error occurred while fetching nat gateway rules: %s", err.Error()))
			//}
			//
			//var results = *natGatewayRules.Items

			partialMatch := d.Get("partial_match").(bool)

			log.Printf("[INFO] Using data source for nat gateway rule by name with partial_match %t and name: %s", partialMatch, name)

			if partialMatch {
				natGatewayRules, apiResponse, err := client.NATGatewaysApi.DatacentersNatgatewaysRulesGet(ctx, datacenterId, natgatewayId).Depth(1).Filter("name", name).Execute()
				logApiRequestTime(apiResponse)
				if err != nil {
					return diag.FromErr(fmt.Errorf("an error occurred while fetching nat gateway rules: %s", err.Error()))
				}
				results = *natGatewayRules.Items
			} else {
				natGatewayRules, apiResponse, err := client.NATGatewaysApi.DatacentersNatgatewaysRulesGet(ctx, datacenterId, natgatewayId).Depth(1).Execute()
				logApiRequestTime(apiResponse)
				if err != nil {
					return diag.FromErr(fmt.Errorf("an error occurred while fetching nat gateway rules: %s", err.Error()))
				}

				if natGatewayRules.Items != nil && nameOk {
					var resultsByName []ionoscloud.NatGatewayRule
					for _, ngr := range *natGatewayRules.Items {
						if ngr.Properties != nil && ngr.Properties.Name != nil && strings.EqualFold(*ngr.Properties.Name, name) {
							tmpNatGatewayRule, apiResponse, err := client.NATGatewaysApi.DatacentersNatgatewaysRulesFindByNatGatewayRuleId(ctx, datacenterId, natgatewayId, *ngr.Id).Execute()
							logApiRequestTime(apiResponse)
							if err != nil {
								return diag.FromErr(fmt.Errorf("an error occurred while fetching nat gateway rule with ID %s: %s", *ngr.Id, err.Error()))
							}
							resultsByName = append(resultsByName, tmpNatGatewayRule)
						}
					}
					results = resultsByName
				}
			}
		} else {
			natGatewayRules, apiResponse, err := client.NATGatewaysApi.DatacentersNatgatewaysRulesGet(ctx, datacenterId, natgatewayId).Depth(1).Execute()
			logApiRequestTime(apiResponse)
			if err != nil {
				return diag.FromErr(fmt.Errorf("an error occurred while fetching nat gateway rules: %s", err.Error()))
			}
			results = *natGatewayRules.Items
		}

		if protocolOk && protocol != "" {
			var protocolResults []ionoscloud.NatGatewayRule
			if results != nil {
				for _, natGateway := range results {
					if natGateway.Properties != nil && natGateway.Properties.Protocol != nil && strings.EqualFold(string(*natGateway.Properties.Protocol), protocol) {
						protocolResults = append(protocolResults, natGateway)
					}
				}
			}
			if protocolResults == nil {
				return diag.FromErr(fmt.Errorf("no natgateway rule found with the specified criteria: protocolResults = %s", protocol))
			}
			results = protocolResults
		}

		if results == nil || len(results) == 0 {
			return diag.FromErr(fmt.Errorf("no nat gateway rule found with the specified criteria: name = %s", name))
		} else if len(results) > 1 {
			return diag.FromErr(fmt.Errorf("more than one nat gateway rule found with the specified criteria: name = %s", name))
		} else {
			natGatewayRule = results[0]
		}

	}

	if err = setNatGatewayRuleData(d, &natGatewayRule); err != nil {
		return diag.FromErr(err)
	}

	return nil
}
