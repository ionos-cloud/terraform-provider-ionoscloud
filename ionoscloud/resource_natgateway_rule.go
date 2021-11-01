package ionoscloud

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	ionoscloud "github.com/ionos-cloud/sdk-go/v6"
	"log"
	"strings"
)

func resourceNatGatewayRule() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceNatGatewayRuleCreate,
		ReadContext:   resourceNatGatewayRuleRead,
		UpdateContext: resourceNatGatewayRuleUpdate,
		DeleteContext: resourceNatGatewayRuleDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceNatGatewayRuleImport,
		},
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

func resourceNatGatewayRuleCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	client := meta.(SdkBundle).CloudApiClient

	natGatewayRule := ionoscloud.NatGatewayRule{
		Properties: &ionoscloud.NatGatewayRuleProperties{},
	}

	if name, nameOk := d.GetOk("name"); nameOk {
		name := name.(string)
		natGatewayRule.Properties.Name = &name
	} else {
		diags := diag.FromErr(fmt.Errorf("name must be provided for nat gateway rule"))
		return diags
	}

	if sourceSubnet, sourceSubnetOk := d.GetOk("source_subnet"); sourceSubnetOk {
		sourceSubnet := sourceSubnet.(string)
		natGatewayRule.Properties.SourceSubnet = &sourceSubnet
	} else {
		diags := diag.FromErr(fmt.Errorf("source subnet must be provided for nat gateway rule"))
		return diags
	}

	if publicIp, publicIpOk := d.GetOk("public_ip"); publicIpOk {
		publicIp := publicIp.(string)
		natGatewayRule.Properties.PublicIp = &publicIp
	} else {
		diags := diag.FromErr(fmt.Errorf("public Ip must be provided for nat gateway rule"))
		return diags
	}

	if ruleType, ruleTypeOk := d.GetOk("type"); ruleTypeOk {

		if strings.ToUpper(ruleType.(string)) != "SNAT" {
			diags := diag.FromErr(fmt.Errorf("attribute value '%s' not allowed. Expected one of [SNAT]", ruleType.(string)))
			return diags
		}
		ruleType := ionoscloud.NatGatewayRuleType(strings.ToUpper(ruleType.(string)))
		natGatewayRule.Properties.Type = &ruleType
	}

	if protocol, protocolOk := d.GetOk("protocol"); protocolOk {
		if strings.ToUpper(protocol.(string)) != "TCP" && strings.ToUpper(protocol.(string)) != "UDP" &&
			strings.ToUpper(protocol.(string)) != "ICMP" && strings.ToUpper(protocol.(string)) != "ALL" {
			diags := diag.FromErr(fmt.Errorf("attribute value '%s' not allowed. Expected one of [TCP, UDP, ICMP, ALL]", protocol.(string)))
			return diags
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
			diags := diag.FromErr(fmt.Errorf("target_port_range start/end can not be set if protocol is ICMP or ALL or not provided"))
			return diags
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

	natGatewayRuleResp, apiResponse, err := client.NATGatewaysApi.DatacentersNatgatewaysRulesPost(ctx, dcId, ngId).NatGatewayRule(natGatewayRule).Execute()
	logApiRequestTime(apiResponse)

	if err != nil {
		d.SetId("")
		diags := diag.FromErr(fmt.Errorf("error creating nat gateway rule: %s \n ApiError %s", err, responseBody(apiResponse)))
		return diags
	}

	d.SetId(*natGatewayRuleResp.Id)

	// Wait, catching any errors
	_, errState := getStateChangeConf(meta, d, apiResponse.Header.Get("Location"), schema.TimeoutCreate).WaitForStateContext(ctx)
	if errState != nil {
		if IsRequestFailed(err) {
			// Request failed, so resource was not created, delete resource from state file
			d.SetId("")
		}
		diags := diag.FromErr(errState)
		return diags
	}

	return resourceNatGatewayRuleRead(ctx, d, meta)
}

func resourceNatGatewayRuleRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(SdkBundle).CloudApiClient

	dcId := d.Get("datacenter_id").(string)
	ngId := d.Get("natgateway_id").(string)

	natGatewayRule, apiResponse, err := client.NATGatewaysApi.DatacentersNatgatewaysRulesFindByNatGatewayRuleId(ctx, dcId, ngId, d.Id()).Execute()
	logApiRequestTime(apiResponse)

	if err != nil {
		log.Printf("[INFO] Resource %s not found: %+v", d.Id(), err)
		if apiResponse != nil && apiResponse.Response != nil && apiResponse.StatusCode == 404 {
			d.SetId("")
			return nil
		}
	}

	log.Printf("[INFO] Successfully retreived nat gateway rule %s: %+v", d.Id(), natGatewayRule)

	if err := setNatGatewayRuleData(d, &natGatewayRule); err != nil {
		return diag.FromErr(err)
	}

	return nil
}
func resourceNatGatewayRuleUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(SdkBundle).CloudApiClient
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

	_, apiResponse, err := client.NATGatewaysApi.
		DatacentersNatgatewaysRulesPatch(ctx, dcId, ngId, d.Id()).
		NatGatewayRuleProperties(*request.Properties).
		Execute()

	if err != nil {
		diags := diag.FromErr(fmt.Errorf("an error occured while updating a nat gateway rule ID %s %s \n ApiError: %s", d.Id(), err, responseBody(apiResponse)))
		return diags
	}

	_, errState := getStateChangeConf(meta, d, apiResponse.Header.Get("Location"), schema.TimeoutUpdate).WaitForStateContext(ctx)
	if errState != nil {
		diags := diag.FromErr(errState)
		return diags
	}

	return resourceNatGatewayRuleRead(ctx, d, meta)
}

func resourceNatGatewayRuleDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(SdkBundle).CloudApiClient

	dcId := d.Get("datacenter_id").(string)
	ngId := d.Get("natgateway_id").(string)

	apiResponse, err := client.NATGatewaysApi.DatacentersNatgatewaysRulesDelete(ctx, dcId, ngId, d.Id()).Execute()
	logApiRequestTime(apiResponse)

	if err != nil {
		diags := diag.FromErr(fmt.Errorf("an error occured while deleting a nat gateway rule %s %s", d.Id(), err))
		return diags
	}

	// Wait, catching any errors
	_, errState := getStateChangeConf(meta, d, apiResponse.Header.Get("Location"), schema.TimeoutDelete).WaitForStateContext(ctx)
	if errState != nil {
		diags := diag.FromErr(errState)
		return diags
	}

	d.SetId("")

	return nil
}

func resourceNatGatewayRuleImport(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	client := meta.(*ionoscloud.APIClient)

	parts := strings.Split(d.Id(), "/")
	if len(parts) != 3 || parts[0] == "" || parts[1] == "" || parts[2] == "" {
		return nil, fmt.Errorf("invalid import id %q. Expecting {datacenter}/{natgateway}/{natgateway_rule}", d.Id())
	}

	dcId := parts[0]
	natGatewayId := parts[1]
	natGatewayRuleId := parts[2]

	natGatewayRule, apiResponse, err := client.NATGatewaysApi.DatacentersNatgatewaysRulesFindByNatGatewayRuleId(ctx, dcId, natGatewayId, natGatewayRuleId).Execute()
	logApiRequestTime(apiResponse)

	if err != nil {
		log.Printf("[INFO] Resource %s not found: %+v", d.Id(), err)
		if apiResponse != nil && apiResponse.Response != nil && apiResponse.StatusCode == 404 {
			d.SetId("")
			return nil, fmt.Errorf("unable to find nat gateway rule %q", natGatewayRuleId)
		}
		return nil, fmt.Errorf("an error occured while retrieving nat gateway rule  %q: %q ", natGatewayRuleId, err)
	}

	if err := d.Set("datacenter_id", dcId); err != nil {
		return nil, err
	}
	if err := d.Set("natgateway_id", natGatewayId); err != nil {
		return nil, err
	}

	if err := setNatGatewayRuleData(d, &natGatewayRule); err != nil {
		return nil, err
	}

	return []*schema.ResourceData{d}, nil
}
