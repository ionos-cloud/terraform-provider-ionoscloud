package ionoscloud

import (
	"context"
	"errors"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	ionoscloud "github.com/ionos-cloud/sdk-go/v6"
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
			"type": {
				Type:        schema.TypeString,
				Description: "Type of the NAT gateway rule.",
				Computed:    true,
			},
			"protocol": {
				Type: schema.TypeString,
				Description: "Protocol of the NAT gateway rule. Defaults to ALL. If protocol is 'ICMP' then " +
					"targetPortRange start and end cannot be set.",
				Computed: true,
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
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.All(validation.StringIsNotWhiteSpace),
			},
			"natgateway_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.All(validation.StringIsNotWhiteSpace),
			},
		},
		Timeouts: &resourceDefaultTimeouts,
	}
}

func dataSourceNatGatewayRuleRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(SdkBundle).CloudApiClient

	datacenterId, dcIdOk := d.GetOk("datacenter_id")
	if !dcIdOk {
		return diag.FromErr(errors.New("no datacenter_id was specified"))
	}

	natgatewayId, ngIdOk := d.GetOk("natgateway_id")
	if !ngIdOk {
		return diag.FromErr(errors.New("no natgateway_id was specified"))
	}

	id, idOk := d.GetOk("id")
	name, nameOk := d.GetOk("name")

	if idOk && nameOk {
		return diag.FromErr(errors.New("id and name cannot be both specified in the same time"))
	}
	if !idOk && !nameOk {
		return diag.FromErr(errors.New("please provide either the lan id or name"))
	}

	var natGatewayRule ionoscloud.NatGatewayRule
	var err error
	var apiResponse *ionoscloud.APIResponse

	if idOk {
		/* search by ID */
		natGatewayRule, apiResponse, err = client.NATGatewaysApi.DatacentersNatgatewaysRulesFindByNatGatewayRuleId(ctx, datacenterId.(string), natgatewayId.(string), id.(string)).Execute()
		logApiRequestTime(apiResponse)
		if err != nil {
			return diag.FromErr(fmt.Errorf("an error occurred while fetching the nat gateway rule %s: %s", id.(string), err))
		}
	} else {
		/* search by name */
		var natGatewayRules ionoscloud.NatGatewayRules

		natGatewayRules, apiResponse, err := client.NATGatewaysApi.DatacentersNatgatewaysRulesGet(ctx, datacenterId.(string), natgatewayId.(string)).Execute()
		logApiRequestTime(apiResponse)
		if err != nil {
			return diag.FromErr(fmt.Errorf("an error occurred while fetching nat gateway rules: %s", err.Error()))
		}

		var results []ionoscloud.NatGatewayRule
		if natGatewayRules.Items != nil {
			for _, c := range *natGatewayRules.Items {
				tmpNatGatewayRule, apiResponse, err := client.NATGatewaysApi.DatacentersNatgatewaysRulesFindByNatGatewayRuleId(ctx, datacenterId.(string), natgatewayId.(string), *c.Id).Execute()
				logApiRequestTime(apiResponse)
				if err != nil {
					return diag.FromErr(fmt.Errorf("an error occurred while fetching nat gateway rule with ID %s: %s", *c.Id, err.Error()))
				}
				if tmpNatGatewayRule.Properties != nil && tmpNatGatewayRule.Properties.Name != nil && *tmpNatGatewayRule.Properties.Name == name.(string) {
					results = append(results, tmpNatGatewayRule)
				}

			}
		}

		if results == nil || len(results) == 0 {
			return diag.FromErr(fmt.Errorf("no nat gateway rule found with the specified criteria: name = %s", name.(string)))
		} else if len(results) > 1 {
			return diag.FromErr(fmt.Errorf("more than one nat gateway rule found with the specified criteria: name = %s", name.(string)))
		} else {
			natGatewayRule = results[0]
		}

	}

	if err = setNatGatewayRuleData(d, &natGatewayRule); err != nil {
		return diag.FromErr(err)
	}

	return nil
}
