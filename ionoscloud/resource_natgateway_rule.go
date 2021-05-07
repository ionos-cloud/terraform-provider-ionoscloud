package ionoscloud

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	ionoscloud "github.com/ionos-cloud/sdk-go/v5"
	"log"
	"strings"
)

func resourceNatGatewayRule() *schema.Resource {
	return &schema.Resource{
		Create: resourceNatGatewayRuleCreate,
		Read:   resourceNatGatewayRuleRead,
		Update: resourceNatGatewayRuleUpdate,
		Delete: resourceNatGatewayRuleDelete,
		Schema: map[string]*schema.Schema{

			"name": {
				Type:        schema.TypeString,
				Description: "Name of the NAT gateway rule",
				Required:    true,
			},
			"type": {
				Type:        schema.TypeString,
				Description: "Type of the NAT gateway rule.",
				Optional:    true,
				Computed:    true,
			},
			"protocol": {
				Type: schema.TypeString,
				Description: "Protocol of the NAT gateway rule. Defaults to ALL. If protocol is 'ICMP' then " +
					"targetPortRange start and end cannot be set.",
				Optional: true,
				Computed: true,
			},
			"source_subnet": {
				Type: schema.TypeString,
				Description: "Source subnet of the NAT gateway rule. For SNAT rules it specifies which packets this " +
					"translation rule applies to based on the packets source IP address.",
				Required: true,
			},
			"public_ip": {
				Type: schema.TypeString,
				Description: "Public IP address of the NAT gateway rule. Specifies the address used for masking outgoing " +
					"packets source address field. Should be one of the customer reserved IP address already " +
					"configured on the NAT gateway resource",
				Required: true,
			},
			"target_subnet": {
				Type: schema.TypeString,
				Description: "Target or destination subnet of the NAT gateway rule. For SNAT rules it specifies which " +
					"packets this translation rule applies to based on the packets destination IP address. If " +
					"none is provided, rule will match any address.",
				Optional: true,
				Computed: true,
			},
			"target_port_range": {
				Type:     schema.TypeList,
				MaxItems: 1,
				Description: "Target port range of the NAT gateway rule. For SNAT rules it specifies which packets this " +
					"translation rule applies to based on destination port. If none is provided, rule will " +
					"match any port",
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"start": {
							Type:        schema.TypeInt,
							Description: "Target port range start associated with the NAT gateway rule.",
							Optional:    true,
							Computed:    true,
						},
						"end": {
							Type:        schema.TypeInt,
							Description: "Target port range end associated with the NAT gateway rule.",
							Optional:    true,
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

func resourceNatGatewayRuleCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(SdkBundle).Client

	natGatewayRule := ionoscloud.NatGatewayRule{
		Properties: &ionoscloud.NatGatewayRuleProperties{},
	}

	if name, nameOk := d.GetOk("name"); nameOk {
		name := name.(string)
		natGatewayRule.Properties.Name = &name
	} else {
		return fmt.Errorf("Name must be provided for nat gateway rule")
	}

	if sourceSubnet, sourceSubnetOk := d.GetOk("source_subnet"); sourceSubnetOk {
		sourceSubnet := sourceSubnet.(string)
		natGatewayRule.Properties.SourceSubnet = &sourceSubnet
	} else {
		return fmt.Errorf("Source subnet must be provided for nat gateway rule")
	}

	if publicIp, publicIpOk := d.GetOk("public_ip"); publicIpOk {
		publicIp := publicIp.(string)
		natGatewayRule.Properties.PublicIp = &publicIp
	} else {
		return fmt.Errorf("Public Ip must be provided for nat gateway rule")
	}

	if ruleType, ruleTypeOk := d.GetOk("type"); ruleTypeOk {

		if strings.ToUpper(ruleType.(string)) != "SNAT" {
			return fmt.Errorf("Attribute value '%s' not allowed. Expected one of [SNAT]", ruleType.(string))
		}
		ruleType := ionoscloud.NatGatewayRuleType(strings.ToUpper(ruleType.(string)))
		natGatewayRule.Properties.Type = &ruleType
	}

	if protocol, protocolOk := d.GetOk("protocol"); protocolOk {
		if strings.ToUpper(protocol.(string)) != "TCP" && strings.ToUpper(protocol.(string)) != "UDP" &&
			strings.ToUpper(protocol.(string)) != "ICMP" && strings.ToUpper(protocol.(string)) != "ALL" {
			return fmt.Errorf("Attribute value '%s' not allowed. Expected one of [TCP, UDP, ICMP, ALL]", protocol.(string))
		}
		protocol := ionoscloud.NatGatewayRuleProtocol(strings.ToUpper(protocol.(string)))
		natGatewayRule.Properties.Protocol = &protocol
	}

	if targetSubnet, targetSubnetOk := d.GetOk("source_subnet"); targetSubnetOk {
		targetSubnet := targetSubnet.(string)
		natGatewayRule.Properties.TargetSubnet = &targetSubnet
	}

	if _, targetPortRangeOk := d.GetOk("target_port_range.0"); targetPortRangeOk {
		if *natGatewayRule.Properties.Protocol == "ICMP" || *natGatewayRule.Properties.Protocol == "ALL" {
			return fmt.Errorf("TargetPortRange start/end can not be set if protocol is ICMP or ALL or not provided.")
		}
		natGatewayRule.Properties.TargetPortRange = &ionoscloud.TargetPortRange{}
	}

	if start, startOk := d.GetOk("target_port_range.0.start"); startOk {
		start := int32(start.(int))
		natGatewayRule.Properties.TargetPortRange.Start = &start
	}

	if end, endOk := d.GetOk("target_port_range.0.end"); endOk {
		end := int32(end.(int))
		natGatewayRule.Properties.TargetPortRange.End = &end
	}

	ngId := d.Get("natgateway_id").(string)
	dcId := d.Get("datacenter_id").(string)

	ctx, cancel := context.WithTimeout(context.Background(), *resourceDefaultTimeouts.Create)
	if cancel != nil {
		defer cancel()
	}

	natGatewayRuleResp, apiResponse, err := client.NATGatewaysApi.DatacentersNatgatewaysRulesPost(ctx, dcId, ngId).NatGatewayRule(natGatewayRule).Execute()

	if err != nil {
		d.SetId("")
		return fmt.Errorf("Error creating nat gateway rule: %s \n ApiError %s", err, string(apiResponse.Payload))
	}

	d.SetId(*natGatewayRuleResp.Id)

	// Wait, catching any errors
	_, errState := getStateChangeConf(meta, d, apiResponse.Header.Get("Location"), schema.TimeoutCreate).WaitForState()
	if errState != nil {
		if IsRequestFailed(err) {
			// Request failed, so resource was not created, delete resource from state file
			d.SetId("")
		}
		return errState
	}

	return resourceNatGatewayRuleRead(d, meta)
}

func resourceNatGatewayRuleRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(SdkBundle).Client

	dcId := d.Get("datacenter_id").(string)
	ngId := d.Get("natgateway_id").(string)

	ctx, cancel := context.WithTimeout(context.Background(), *resourceDefaultTimeouts.Default)

	if cancel != nil {
		defer cancel()
	}

	natGatewayRule, apiResponse, err := client.NATGatewaysApi.DatacentersNatgatewaysRulesFindByNatGatewayRuleId(ctx, dcId, ngId, d.Id()).Execute()

	if err != nil {
		log.Printf("[INFO] Resource %s not found: %+v", d.Id(), err)
		if _, ok := err.(ionoscloud.GenericOpenAPIError); ok {
			if apiResponse.Response.StatusCode == 404 {
				d.SetId("")
				return nil
			}
		}
	}

	log.Printf("[INFO] Successfully retreived nat gateway rule %s: %+v", d.Id(), natGatewayRule)

	if natGatewayRule.Properties.Name != nil {
		err := d.Set("name", *natGatewayRule.Properties.Name)
		if err != nil {
			return fmt.Errorf("Error while setting name property for nat gateway %s: %s", d.Id(), err)
		}
	}

	if natGatewayRule.Properties.Type != nil {
		err := d.Set("type", *natGatewayRule.Properties.Type)
		if err != nil {
			return fmt.Errorf("Error while setting type property for nat gateway %s: %s", d.Id(), err)
		}
	}

	if natGatewayRule.Properties.Protocol != nil {
		err := d.Set("protocol", *natGatewayRule.Properties.Protocol)
		if err != nil {
			return fmt.Errorf("Error while setting protocol property for nat gateway %s: %s", d.Id(), err)
		}
	}

	if natGatewayRule.Properties.SourceSubnet != nil {
		err := d.Set("source_subnet", *natGatewayRule.Properties.SourceSubnet)
		if err != nil {
			return fmt.Errorf("Error while setting source_subnet property for nat gateway %s: %s", d.Id(), err)
		}
	}

	if natGatewayRule.Properties.PublicIp != nil {
		err := d.Set("public_ip", *natGatewayRule.Properties.PublicIp)
		if err != nil {
			return fmt.Errorf("Error while setting public_ip property for nat gateway %s: %s", d.Id(), err)
		}
	}

	if natGatewayRule.Properties.TargetSubnet != nil {
		err := d.Set("target_subnet", *natGatewayRule.Properties.TargetSubnet)
		if err != nil {
			return fmt.Errorf("Error while setting target_subnet property for nat gateway %s: %s", d.Id(), err)
		}
	}

	if natGatewayRule.Properties.TargetPortRange != nil && natGatewayRule.Properties.TargetPortRange.Start != nil &&
		natGatewayRule.Properties.TargetPortRange.End != nil {
		err := d.Set("target_port_range", []map[string]int32{{
			"start": *natGatewayRule.Properties.TargetPortRange.Start,
			"end":   *natGatewayRule.Properties.TargetPortRange.End,
		},
		})
		if err != nil {
			return fmt.Errorf("Error while setting target_port_range property for nat gateway %s: %s", d.Id(), err)
		}
	}

	return nil
}
func resourceNatGatewayRuleUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(SdkBundle).Client
	request := ionoscloud.NatGatewayRule{
		Properties: &ionoscloud.NatGatewayRuleProperties{},
	}

	dcId := d.Get("datacenter_id").(string)
	ngId := d.Get("natgateway_id").(string)

	if d.HasChange("name") {
		_, v := d.GetChange("name")
		vStr := v.(string)
		request.Properties.Name = &vStr
	}

	if d.HasChange("type") {
		_, v := d.GetChange("type")
		vStr := ionoscloud.NatGatewayRuleType(v.(string))
		request.Properties.Type = &vStr
	}

	if d.HasChange("protocol") {
		_, v := d.GetChange("protocol")
		vStr := ionoscloud.NatGatewayRuleProtocol(v.(string))
		request.Properties.Protocol = &vStr
	}

	if d.HasChange("source_subnet") {
		_, v := d.GetChange("source_subnet")
		vStr := v.(string)
		request.Properties.SourceSubnet = &vStr
	}

	if d.HasChange("public_ip") {
		_, v := d.GetChange("public_ip")
		vStr := v.(string)
		request.Properties.PublicIp = &vStr
	}

	if d.HasChange("target_subnet") {
		_, v := d.GetChange("target_subnet")
		vStr := v.(string)
		request.Properties.TargetSubnet = &vStr
	}

	if d.HasChange("target_port_range.0") {
		_, v := d.GetChange("target_port_range.0")
		if v.(map[string]interface{}) != nil {
			updateTargetPortRange := false
			start := int32(d.Get("target_port_range.0.start").(int))
			end := int32(d.Get("target_port_range.0.end").(int))
			targetPortRange := &ionoscloud.TargetPortRange{
				Start: &start,
				End:   &end,
			}

			if d.HasChange("target_port_range.0.start") {
				_, newStart := d.GetChange("target_port_range.0.start")
				if newStart != 0 {
					updateTargetPortRange = true
					newStart := int32(newStart.(int))
					targetPortRange.Start = &newStart
				}
			}

			if d.HasChange("target_port_range.0.end") {
				_, newEnd := d.GetChange("target_port_range.0.end")
				if newEnd != 0 {
					updateTargetPortRange = true
					newEnd := int32(newEnd.(int))
					targetPortRange.End = &newEnd
				}
			}

			if updateTargetPortRange == true {
				request.Properties.TargetPortRange = targetPortRange
			}
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), *resourceDefaultTimeouts.Update)

	if cancel != nil {
		defer cancel()
	}
	_, apiResponse, err := client.NATGatewaysApi.DatacentersNatgatewaysRulesPatch(ctx, dcId, ngId, d.Id()).NatGatewayRuleProperties(*request.Properties).Execute()

	if err != nil {
		return fmt.Errorf("An error occured while updating a nat gateway rule ID %s %s \n ApiError: %s", d.Id(), err, string(apiResponse.Payload))
	}

	_, errState := getStateChangeConf(meta, d, apiResponse.Header.Get("Location"), schema.TimeoutUpdate).WaitForState()
	if errState != nil {
		return errState
	}

	return resourceNatGatewayRuleRead(d, meta)
}

func resourceNatGatewayRuleDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(SdkBundle).Client

	dcId := d.Get("datacenter_id").(string)
	ngId := d.Get("natgateway_id").(string)

	ctx, cancel := context.WithTimeout(context.Background(), *resourceDefaultTimeouts.Delete)

	if cancel != nil {
		defer cancel()
	}

	_, apiResponse, err := client.NATGatewaysApi.DatacentersNatgatewaysRulesDelete(ctx, dcId, ngId, d.Id()).Execute()

	if err != nil {
		return fmt.Errorf("An error occured while deleting a nat gateway rule %s %s", d.Id(), err)
	}

	// Wait, catching any errors
	_, errState := getStateChangeConf(meta, d, apiResponse.Header.Get("Location"), schema.TimeoutDelete).WaitForState()
	if errState != nil {
		return errState
	}

	d.SetId("")

	return nil
}
