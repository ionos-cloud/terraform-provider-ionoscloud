package ionoscloud

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services"

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
			"firewall_rules": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"protocol": {
							Type:             schema.TypeString,
							Required:         true,
							ValidateDiagFunc: validation.ToDiagFunc(validation.StringIsNotWhiteSpace),
						},
						"source_mac": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"source_ip": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"target_ip": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"port_range_start": {
							Type:     schema.TypeInt,
							Optional: true,
							ValidateDiagFunc: validation.ToDiagFunc(func(v any, k string) (ws []string, errors []error) {
								if v.(int) < 1 && v.(int) > 65534 {
									errors = append(errors, fmt.Errorf("port start range must be between 1 and 65534"))
								}
								return
							}),
						},
						"port_range_end": {
							Type:     schema.TypeInt,
							Optional: true,
							ValidateDiagFunc: validation.ToDiagFunc(func(v any, k string) (ws []string, errors []error) {
								if v.(int) < 1 && v.(int) > 65534 {
									errors = append(errors, fmt.Errorf("port end range must be between 1 and 65534"))
								}
								return
							}),
						},
						"icmp_type": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"icmp_code": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"type": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
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
	// todo - if needed
	// if errState := cloudapi.WaitForStateChange(ctx, meta, d, apiResponse, schema.TimeoutCreate); errState != nil {
	//
	// }
	sg.Entities = ionoscloud.NewSecurityGroupRequestEntitiesWithDefaults()

	rulesObjects := ionoscloud.NewFirewallRules()
	rules := d.Get("firewall_rules").([]any)
	for i := range rules {
		rule, diags := getFirewallData(d, fmt.Sprintf("firewall_rules.%d.", i), false)
		if diags != nil {
			return diags
		}

		*rulesObjects.Items = append(*rulesObjects.Items, rule)
	}

	sg.Entities.SetRules(*rulesObjects)

	securityGroup, apiResponse, err := client.SecurityGroupsApi.DatacentersSecuritygroupsPost(ctx, datacenterID).SecurityGroup(sg).Execute()
	apiResponse.LogInfo()
	if err != nil {
		diags := diag.FromErr(fmt.Errorf("an error occured while creating a security group: %w", err))
		return diags
	}
	d.SetId(*securityGroup.Id)

	return nil
}

func resourceNetworkSecurityGroupRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(services.SdkBundle).CloudApiClient

	datacenterID := d.Get("datacenter_id").(string)

	securityGroup, apiResponse, err := client.SecurityGroupsApi.DatacentersSecuritygroupsFindById(ctx, datacenterID, d.Id()).Depth(3).Execute()
	apiResponse.LogInfo()
	if err != nil {
		diags := diag.FromErr(fmt.Errorf("an error occured while retrieving a security group: %w", err))
		return diags
	}

	if securityGroup.Properties != nil {
		if securityGroup.Properties.Name != nil {
			err := d.Set("name", *securityGroup.Properties.Name)
			if err != nil {
				return diag.FromErr(fmt.Errorf("error while setting securityGroup property  %s: %w", d.Id(), err))
			}
		}
		if securityGroup.Properties.Description != nil {
			err := d.Set("description", *securityGroup.Properties.Description)
			if err != nil {
				return diag.FromErr(fmt.Errorf("error while setting securityGroup property  %s: %w", d.Id(), err))
			}
		}
	}

	if securityGroup.Entities != nil {
		if securityGroup.Entities.Rules != nil && securityGroup.Entities.Rules.Items != nil {
			rules := make([]interface{}, 0)
			for _, rule := range *securityGroup.Entities.Rules.Items {
				ruleProperties := rule.Properties
				newRule := make(map[string]interface{})
				if rule.Id != nil {
					newRule["id"] = rule.Id
				}
				if ruleProperties != nil {
					if ruleProperties.Name != nil {
						newRule["name"] = *ruleProperties.Name
					}
					if ruleProperties.Protocol != nil {
						newRule["protocol"] = *ruleProperties.Protocol
					}
					if ruleProperties.SourceMac != nil {
						newRule["source_mac"] = *ruleProperties.SourceMac
					}
					if ruleProperties.SourceIp != nil {
						newRule["source_ip"] = *ruleProperties.SourceIp
					}
					if ruleProperties.Name != nil {
						newRule["target_ip"] = *ruleProperties.Name
					}
					if ruleProperties.PortRangeStart != nil {
						newRule["port_range_start"] = *ruleProperties.PortRangeStart
					}
					if ruleProperties.PortRangeEnd != nil {
						newRule["port_range_end"] = *ruleProperties.PortRangeEnd
					}
					if ruleProperties.IcmpType != nil {
						newRule["icmp_type"] = *ruleProperties.IcmpType
					}
					if ruleProperties.IcmpCode != nil {
						newRule["icmp_code"] = *ruleProperties.IcmpCode
					}
					if ruleProperties.Type != nil {
						newRule["type"] = *ruleProperties.Type
					}
				}
				rules = append(rules, newRule)

			}
			if len(rules) > 0 {
				if err := d.Set("firewall_rules", rules); err != nil {
					return diag.FromErr(fmt.Errorf("error while setting firewall_rules property for NetworkSecurityGroup  %s: %w", d.Id(), err))
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

	sg.Entities = ionoscloud.NewSecurityGroupRequestEntitiesWithDefaults()

	rulesObjects := ionoscloud.NewFirewallRules()
	rules := d.Get("firewall_rules").([]any)
	for i := range rules {
		rule, diags := getFirewallData(d, fmt.Sprintf("firewall_rules.%d.", i), false)
		if diags != nil {
			return diags
		}

		*rulesObjects.Items = append(*rulesObjects.Items, rule)
	}

	sg.Entities.SetRules(*rulesObjects)

	_, apiResponse, err := client.SecurityGroupsApi.DatacentersSecuritygroupsPut(ctx, datacenterID, d.Id()).SecurityGroup(sg).Execute()
	apiResponse.LogInfo()
	if err != nil {
		diags := diag.FromErr(fmt.Errorf("an error occured while creating a security group: %w", err))
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
		diags := diag.FromErr(fmt.Errorf("an error occured while creating a security group: %w", err))
		return diags
	}
	// todo - if needed
	// if errState := cloudapi.WaitForStateChange(ctx, meta, d, apiResponse, schema.TimeoutCreate); errState != nil {
	//
	// }

	return nil
}
