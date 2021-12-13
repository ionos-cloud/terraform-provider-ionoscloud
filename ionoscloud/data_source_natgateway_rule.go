package ionoscloud

import (
	"context"
	"errors"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	ionoscloud "github.com/ionos-cloud/sdk-go/v6"
	"strings"
)

func dataSourceNatGatewayRule() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceNatGatewayRuleRead,
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

func dataSourceNatGatewayRuleRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(SdkBundle).CloudApiClient

	datacenterId, dcIdOk := d.GetOk("datacenter_id")
	if !dcIdOk {
		return errors.New("no datacenter_id was specified")
	}

	natgatewayId, ngIdOk := d.GetOk("natgateway_id")
	if !ngIdOk {
		return errors.New("no natgateway_id was specified")
	}

	id, idOk := d.GetOk("id")
	name, nameOk := d.GetOk("name")

	if idOk && nameOk {
		return errors.New("id and name cannot be both specified in the same time")
	}
	if !idOk && !nameOk {
		return errors.New("please provide either the lan id or name")
	}

	var natGatewayRule ionoscloud.NatGatewayRule
	var err error
	var apiResponse *ionoscloud.APIResponse

	ctx, cancel := context.WithTimeout(context.Background(), *resourceDefaultTimeouts.Default)

	if cancel != nil {
		defer cancel()
	}

	if idOk {
		/* search by ID */
		natGatewayRule, apiResponse, err = client.NATGatewaysApi.DatacentersNatgatewaysRulesFindByNatGatewayRuleId(ctx, datacenterId.(string), natgatewayId.(string), id.(string)).Execute()
		logApiRequestTime(apiResponse)
		if err != nil {
			return fmt.Errorf("an error occurred while fetching the nat gateway rule %s: %s", id.(string), err)
		}
	} else {
		/* search by name */
		var natGatewayRules ionoscloud.NatGatewayRules

		ctx, cancel := context.WithTimeout(context.Background(), *resourceDefaultTimeouts.Default)

		if cancel != nil {
			defer cancel()
		}

		natGatewayRules, apiResponse, err := client.NATGatewaysApi.DatacentersNatgatewaysRulesGet(ctx, datacenterId.(string), natgatewayId.(string)).Execute()
		logApiRequestTime(apiResponse)
		if err != nil {
			return fmt.Errorf("an error occurred while fetching nat gateway rules: %s", err.Error())
		}

		if natGatewayRules.Items != nil {
			for _, c := range *natGatewayRules.Items {
				tmpNatGatewayRule, apiResponse, err := client.NATGatewaysApi.DatacentersNatgatewaysRulesFindByNatGatewayRuleId(ctx, datacenterId.(string), natgatewayId.(string), *c.Id).Execute()
				logApiRequestTime(apiResponse)
				if err != nil {
					return fmt.Errorf("an error occurred while fetching nat gateway rule with ID %s: %s", *c.Id, err.Error())
				}
				if tmpNatGatewayRule.Properties.Name != nil {
					if strings.Contains(*tmpNatGatewayRule.Properties.Name, name.(string)) {
						natGatewayRule = tmpNatGatewayRule
						break
					}
				}

			}
		}

	}

	if &natGatewayRule == nil {
		return errors.New("nat gateway rule not found")
	}

	if natGatewayRule.Id != nil {
		if err := d.Set("id", *natGatewayRule.Id); err != nil {
			return err
		}
	}

	if err = setNatGatewayRuleData(d, &natGatewayRule); err != nil {
		return err
	}

	return nil
}
