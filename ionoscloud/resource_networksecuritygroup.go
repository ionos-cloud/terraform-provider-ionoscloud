package ionoscloud

import (
	"context"
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/cloudapi"

	ionoscloud "github.com/ionos-cloud/sdk-go/v6"
)

func resourceNetworkSecurityGroup() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceNetworkSecurityGroupCreate,
		ReadContext:   resourceNetworkSecurityGroupRead,
		UpdateContext: resourceNetworkSecurityGroupUpdate,
		DeleteContext: resourceNetworkSecurityGroupDelete,
		// Importer: &schema.ResourceImporter{
		// 	StateContext: resourceFirewallImport,
		// },
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Required: true,
			},
			"rule_ids": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"datacenter_id": {
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringIsNotWhiteSpace),
			},
		},
		Timeouts: &resourceDefaultTimeouts,
	}
}

func resourceNetworkSecurityGroupCreate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(services.SdkBundle).CloudApiClient

	datacenterID := d.Get("datacenter_id").(string)
	sgName := d.Get("name").(string)
	sgDescription := d.Get("description").(string)

	sg := ionoscloud.SecurityGroupRequest{
		Properties: &ionoscloud.SecurityGroupProperties{
			Name:        &sgName,
			Description: &sgDescription,
		},
	}

	securityGroup, apiResponse, err := client.SecurityGroupsApi.DatacentersSecuritygroupsPost(ctx, datacenterID).SecurityGroup(sg).Execute()
	apiResponse.LogInfo()
	if err != nil {
		return diag.FromErr(fmt.Errorf("an error occured while creating a Network Security Group for datacenter dcID: %s, %w", datacenterID, err))
	}
	if errState := cloudapi.WaitForStateChange(ctx, meta, d, apiResponse, schema.TimeoutCreate); errState != nil {
		return diag.FromErr(fmt.Errorf("an error occured while waiting for Network Security Group to be created for datacenter dcID: %s,  %w", datacenterID, err))
	}
	d.SetId(*securityGroup.Id)

	return nil
}

func resourceNetworkSecurityGroupRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(services.SdkBundle).CloudApiClient
	datacenterID := d.Get("datacenter_id").(string)

	securityGroup, apiResponse, err := client.SecurityGroupsApi.DatacentersSecuritygroupsFindById(ctx, datacenterID, d.Id()).Depth(2).Execute()
	apiResponse.LogInfo()
	if err != nil {
		return diag.FromErr(fmt.Errorf("an error occured while retrieving a network security group: %w", err))
	}
	if securityGroup.Properties != nil {
		if securityGroup.Properties.Name != nil {
			err := d.Set("name", *securityGroup.Properties.Name)
			if err != nil {
				return diag.FromErr(fmt.Errorf("error while setting Network Security Group name  %s: %w", d.Id(), err))
			}
		}
		if securityGroup.Properties.Description != nil {
			err := d.Set("description", *securityGroup.Properties.Description)
			if err != nil {
				return diag.FromErr(fmt.Errorf("error while setting Network Security Group description  %s: %w", d.Id(), err))
			}
		}
	}
	if securityGroup.Entities != nil {
		if securityGroup.Entities.Rules != nil && securityGroup.Entities.Rules.Items != nil {
			rule_ids := make([]string, 0)
			for _, rule := range *securityGroup.Entities.Rules.Items {
				rule_ids = append(rule_ids, *rule.Id)

			}
			if len(rule_ids) > 0 {
				if err := d.Set("rule_ids", rule_ids); err != nil {
					return diag.FromErr(fmt.Errorf("error while setting rule_ids property for NetworkSecurityGroup  %s: %w", d.Id(), err))
				}
			}
		}
	}
	return nil
}

func resourceNetworkSecurityGroupUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(services.SdkBundle).CloudApiClient

	datacenterID := d.Get("datacenter_id").(string)
	sgName := d.Get("name").(string)
	sgDescription := d.Get("description").(string)

	sg := ionoscloud.SecurityGroupRequest{
		Properties: &ionoscloud.SecurityGroupProperties{
			Name:        &sgName,
			Description: &sgDescription,
		},
	}

	_, apiResponse, err := client.SecurityGroupsApi.DatacentersSecuritygroupsPut(ctx, datacenterID, d.Id()).SecurityGroup(sg).Execute()
	apiResponse.LogInfo()
	if err != nil {
		diags := diag.FromErr(fmt.Errorf("an error occured while updating security group: %w", err))
		return diags
	}
	return nil
}

func resourceNetworkSecurityGroupDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(services.SdkBundle).CloudApiClient

	datacenterID := d.Get("datacenter_id").(string)

	apiResponse, err := client.SecurityGroupsApi.DatacentersSecuritygroupsDelete(ctx, datacenterID, d.Id()).Execute()
	apiResponse.LogInfo()
	if err != nil {
		diags := diag.FromErr(fmt.Errorf("an error occured while deleting a network security group: %w", err))
		return diags
	}
	// todo - if needed
	// if errState := cloudapi.WaitForStateChange(ctx, meta, d, apiResponse, schema.TimeoutCreate); errState != nil {
	//
	// }

	return nil
}

func getFirewallRulesData(d *schema.ResourceData, onUpdate bool) ([]ionoscloud.FirewallRule, diag.Diagnostics) {
	rules, ok := d.GetOk("firewall_rule")
	if !ok {
		return []ionoscloud.FirewallRule{}, diag.Diagnostics{}
	}
	var firewallRules []ionoscloud.FirewallRule
	for i := range rules.([]interface{}) {
		var ruleId *string
		rulePath := fmt.Sprintf("firewall_rule.%d.", i)
		idValue, idOk := d.GetOk(rulePath + "id")
		if onUpdate && idOk {
			idStr := idValue.(string)
			ruleId = &idStr
		}
		rule, err := getFirewallData(d, rulePath, onUpdate && idOk)
		if err != nil {
			return []ionoscloud.FirewallRule{}, diag.Diagnostics{}
		}
		rule.Id = ruleId
		firewallRules = append(firewallRules, rule)
	}
	return firewallRules, nil
}

func setFirewallRulesData(d *schema.ResourceData, rules []ionoscloud.FirewallRule) diag.Diagnostics {

	var rulesData []map[string]any
	for _, rule := range rules {
		ruleProperties := rule.Properties
		ruleData := make(map[string]any)
		if rule.Id != nil {
			ruleData["id"] = *rule.Id
		}
		if ruleProperties != nil {
			if ruleProperties.Name != nil {
				ruleData["name"] = *ruleProperties.Name
			}
			if ruleProperties.Protocol != nil {
				ruleData["protocol"] = *ruleProperties.Protocol
			}
			if ruleProperties.SourceMac != nil {
				ruleData["source_mac"] = *ruleProperties.SourceMac
			}
			if ruleProperties.SourceIp != nil {
				ruleData["source_ip"] = *ruleProperties.SourceIp
			}
			if ruleProperties.Name != nil {
				ruleData["target_ip"] = *ruleProperties.TargetIp
			}
			if ruleProperties.PortRangeStart != nil {
				ruleData["port_range_start"] = *ruleProperties.PortRangeStart
			}
			if ruleProperties.PortRangeEnd != nil {
				ruleData["port_range_end"] = *ruleProperties.PortRangeEnd
			}
			if ruleProperties.IcmpType != nil {
				ruleData["icmp_type"] = strconv.Itoa(int(*ruleProperties.IcmpType))
			}
			if ruleProperties.IcmpCode != nil {
				ruleData["icmp_code"] = strconv.Itoa(int(*ruleProperties.IcmpCode))
			}
			if ruleProperties.Type != nil {
				ruleData["type"] = *ruleProperties.Type
			}
		}
		rulesData = append(rulesData, ruleData)
	}
	if err := d.Set("firewall_rule", rulesData); err != nil {
		return diag.FromErr(fmt.Errorf("error while setting firewall rules for NetworkSecurityGroup %s: %w", d.Id(), err))
	}
	return diag.Diagnostics{}
}
