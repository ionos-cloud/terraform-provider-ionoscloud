package ionoscloud

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	ionoscloud "github.com/ionos-cloud/sdk-go/v6"

	bundleclient "github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/bundleclient"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils"
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
				Type:             schema.TypeString,
				Description:      "Name of the NAT gateway rule",
				Required:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringIsNotWhiteSpace),
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
				Required:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringIsNotWhiteSpace),
			},
			"public_ip": {
				Type: schema.TypeString,
				Description: "Public IP address of the NAT gateway rule. Specifies the address used for masking outgoing " +
					"packets source address field. Should be one of the customer reserved IP address already " +
					"configured on the NAT gateway resource",
				Required:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringIsNotWhiteSpace),
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
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringIsNotWhiteSpace),
			},
			"natgateway_id": {
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringIsNotWhiteSpace),
			},
		},
		Timeouts: &resourceDefaultTimeouts,
	}
}

func resourceNatGatewayRuleCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	client := meta.(bundleclient.SdkBundle).CloudApiClient

	natGatewayRule := ionoscloud.NatGatewayRule{
		Properties: &ionoscloud.NatGatewayRuleProperties{},
	}

	if name, nameOk := d.GetOk("name"); nameOk {
		name := name.(string)
		natGatewayRule.Properties.Name = &name
	} else {
		return utils.ToDiags(d, "name must be provided for nat gateway rule", nil)
	}

	if sourceSubnet, sourceSubnetOk := d.GetOk("source_subnet"); sourceSubnetOk {
		sourceSubnet := sourceSubnet.(string)
		natGatewayRule.Properties.SourceSubnet = &sourceSubnet
	} else {
		return utils.ToDiags(d, "source subnet must be provided for nat gateway rule", nil)
	}

	if publicIp, publicIpOk := d.GetOk("public_ip"); publicIpOk {
		publicIp := publicIp.(string)
		natGatewayRule.Properties.PublicIp = &publicIp
	} else {
		return utils.ToDiags(d, "public Ip must be provided for nat gateway rule", nil)
	}

	if ruleType, ruleTypeOk := d.GetOk("type"); ruleTypeOk {

		if strings.ToUpper(ruleType.(string)) != "SNAT" {
			return utils.ToDiags(d, fmt.Sprintf("attribute value '%s' not allowed. Expected one of [SNAT]", ruleType.(string)), nil)
		}
		ruleType := ionoscloud.NatGatewayRuleType(strings.ToUpper(ruleType.(string)))
		natGatewayRule.Properties.Type = &ruleType
	}

	if protocol, protocolOk := d.GetOk("protocol"); protocolOk {
		if strings.ToUpper(protocol.(string)) != "TCP" && strings.ToUpper(protocol.(string)) != "UDP" &&
			strings.ToUpper(protocol.(string)) != "ICMP" && strings.ToUpper(protocol.(string)) != "ALL" {
			return utils.ToDiags(d, fmt.Sprintf("attribute value '%s' not allowed. Expected one of [TCP, UDP, ICMP, ALL]", protocol.(string)), nil)
		}
		protocol := ionoscloud.NatGatewayRuleProtocol(strings.ToUpper(protocol.(string)))
		natGatewayRule.Properties.Protocol = &protocol
	}

	if targetSubnet, targetSubnetOk := d.GetOk("target_subnet"); targetSubnetOk {
		targetSubnet := targetSubnet.(string)
		natGatewayRule.Properties.TargetSubnet = &targetSubnet
	}

	if _, targetPortRangeOk := d.GetOk("target_port_range.0"); targetPortRangeOk {
		if *natGatewayRule.Properties.Protocol == "ICMP" || *natGatewayRule.Properties.Protocol == "ALL" {
			return utils.ToDiags(d, "target_port_range start/end can not be set if protocol is ICMP or ALL or not provided", nil)
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
		requestLocation, _ := apiResponse.Location()
		return utils.ToDiags(d, fmt.Sprintf("error creating nat gateway rule: %s \n ApiError %s", err, responseBody(apiResponse)), &utils.DiagsOpts{RequestLocation: requestLocation, StatusCode: apiResponse.StatusCode})
	}

	d.SetId(*natGatewayRuleResp.Id)

	if errState := bundleclient.WaitForStateChange(ctx, meta, d, apiResponse, schema.TimeoutCreate); errState != nil {
		if bundleclient.IsRequestFailed(errState) {
			d.SetId("")
		}
		return utils.ToDiags(d, errState.Error(), &utils.DiagsOpts{Timeout: schema.TimeoutCreate})
	}

	return resourceNatGatewayRuleRead(ctx, d, meta)
}

func resourceNatGatewayRuleRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(bundleclient.SdkBundle).CloudApiClient

	dcId := d.Get("datacenter_id").(string)
	ngId := d.Get("natgateway_id").(string)

	natGatewayRule, apiResponse, err := client.NATGatewaysApi.DatacentersNatgatewaysRulesFindByNatGatewayRuleId(ctx, dcId, ngId, d.Id()).Execute()
	logApiRequestTime(apiResponse)

	if err != nil {
		log.Printf("[INFO] Resource %s not found: %+v", d.Id(), err)
		if httpNotFound(apiResponse) {
			d.SetId("")
			return nil
		}
	}

	log.Printf("[INFO] Successfully retrieved nat gateway rule %s: %+v", d.Id(), natGatewayRule)

	if err := setNatGatewayRuleData(d, &natGatewayRule); err != nil {
		return utils.ToDiags(d, err.Error(), nil)
	}

	return nil
}
func resourceNatGatewayRuleUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(bundleclient.SdkBundle).CloudApiClient
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
		requestLocation, _ := apiResponse.Location()
		return utils.ToDiags(d, fmt.Sprintf("an error occurred while updating a nat gateway rule: %s \n ApiError: %s", err, responseBody(apiResponse)), &utils.DiagsOpts{RequestLocation: requestLocation, StatusCode: apiResponse.StatusCode})
	}

	if errState := bundleclient.WaitForStateChange(ctx, meta, d, apiResponse, schema.TimeoutUpdate); errState != nil {
		return utils.ToDiags(d, errState.Error(), &utils.DiagsOpts{Timeout: schema.TimeoutUpdate})
	}

	return resourceNatGatewayRuleRead(ctx, d, meta)
}

func resourceNatGatewayRuleDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(bundleclient.SdkBundle).CloudApiClient

	dcId := d.Get("datacenter_id").(string)
	ngId := d.Get("natgateway_id").(string)

	apiResponse, err := client.NATGatewaysApi.DatacentersNatgatewaysRulesDelete(ctx, dcId, ngId, d.Id()).Execute()
	logApiRequestTime(apiResponse)

	if err != nil {
		requestLocation, _ := apiResponse.Location()
		return utils.ToDiags(d, fmt.Sprintf("an error occurred while deleting a nat gateway rule: %s", err), &utils.DiagsOpts{RequestLocation: requestLocation, StatusCode: apiResponse.StatusCode})
	}

	if errState := bundleclient.WaitForStateChange(ctx, meta, d, apiResponse, schema.TimeoutDelete); errState != nil {
		return utils.ToDiags(d, errState.Error(), &utils.DiagsOpts{Timeout: schema.TimeoutDelete})
	}

	d.SetId("")

	return nil
}

func resourceNatGatewayRuleImport(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	client := meta.(bundleclient.SdkBundle).CloudApiClient

	parts := strings.Split(d.Id(), "/")
	if len(parts) != 3 || parts[0] == "" || parts[1] == "" || parts[2] == "" {
		return nil, utils.ToError(d, "invalid import. Expecting {datacenter}/{natgateway}/{natgateway_rule}", nil)
	}

	dcId := parts[0]
	natGatewayId := parts[1]
	natGatewayRuleId := parts[2]

	natGatewayRule, apiResponse, err := client.NATGatewaysApi.DatacentersNatgatewaysRulesFindByNatGatewayRuleId(ctx, dcId, natGatewayId, natGatewayRuleId).Execute()
	logApiRequestTime(apiResponse)

	if err != nil {
		log.Printf("[INFO] Resource %s not found: %+v", d.Id(), err)
		if httpNotFound(apiResponse) {
			d.SetId("")
			return nil, utils.ToError(d, fmt.Sprintf("unable to find nat gateway rule %q", natGatewayRuleId), nil)
		}
		return nil, utils.ToError(d, fmt.Sprintf("an error occurred while retrieving nat gateway rule  %q: %q ", natGatewayRuleId, err), nil)
	}

	if err := d.Set("datacenter_id", dcId); err != nil {
		return nil, utils.ToError(d, err.Error(), nil)
	}
	if err := d.Set("natgateway_id", natGatewayId); err != nil {
		return nil, utils.ToError(d, err.Error(), nil)
	}

	if err := setNatGatewayRuleData(d, &natGatewayRule); err != nil {
		return nil, utils.ToError(d, err.Error(), nil)
	}

	return []*schema.ResourceData{d}, nil
}

func setNatGatewayRuleData(d *schema.ResourceData, natGatewayRule *ionoscloud.NatGatewayRule) error {

	if natGatewayRule.Id != nil {
		d.SetId(*natGatewayRule.Id)
	}

	if natGatewayRule.Properties != nil {
		if natGatewayRule.Properties.Name != nil {
			err := d.Set("name", *natGatewayRule.Properties.Name)
			if err != nil {
				return fmt.Errorf("error while setting name property for nat gateway %s: %w", d.Id(), err)
			}
		}

		if natGatewayRule.Properties.Type != nil {
			err := d.Set("type", *natGatewayRule.Properties.Type)
			if err != nil {
				return fmt.Errorf("error while setting type property for nat gateway %s: %w", d.Id(), err)
			}
		}

		if natGatewayRule.Properties.Protocol != nil {
			err := d.Set("protocol", *natGatewayRule.Properties.Protocol)
			if err != nil {
				return fmt.Errorf("error while setting protocol property for nat gateway %s: %w", d.Id(), err)
			}
		}

		if natGatewayRule.Properties.SourceSubnet != nil {
			err := d.Set("source_subnet", *natGatewayRule.Properties.SourceSubnet)
			if err != nil {
				return fmt.Errorf("error while setting source_subnet property for nat gateway %s: %w", d.Id(), err)
			}
		}

		if natGatewayRule.Properties.PublicIp != nil {
			err := d.Set("public_ip", *natGatewayRule.Properties.PublicIp)
			if err != nil {
				return fmt.Errorf("error while setting public_ip property for nat gateway %s: %w", d.Id(), err)
			}
		}

		if natGatewayRule.Properties.TargetSubnet != nil {
			err := d.Set("target_subnet", *natGatewayRule.Properties.TargetSubnet)
			if err != nil {
				return fmt.Errorf("error while setting target_subnet property for nat gateway %s: %w", d.Id(), err)
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
				return fmt.Errorf("error while setting target_port_range property for nat gateway %s: %w", d.Id(), err)
			}
		}
	}
	return nil
}
