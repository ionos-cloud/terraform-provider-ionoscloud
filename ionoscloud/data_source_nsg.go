package ionoscloud

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	ionoscloud "github.com/ionos-cloud/sdk-go-bundle/products/compute/v2"
	"github.com/ionos-cloud/sdk-go-bundle/shared"

	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/bundleclient"
	diagutil "github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils/diags"
)

func dataSourceNSGFirewallRuleSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"id": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"name": {
			Type:     schema.TypeString,
			Computed: true,
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
			Type:     schema.TypeString,
			Computed: true,
		},
		"icmp_code": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"type": {
			Type:     schema.TypeString,
			Computed: true,
		},
	}
}

func dataSourceNSG() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceNSGRead,
		Schema: map[string]*schema.Schema{
			"datacenter_id": {
				Type:             schema.TypeString,
				Required:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.IsUUID),
			},
			"location": {
				Type:        schema.TypeString,
				Description: "The location of the resource. This field should be used only if you are also using a file configuration and should not be configured otherwise.",
				Optional:    true,
			},
			"id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"rule_ids": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"rules": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Resource{Schema: dataSourceNSGFirewallRuleSchema()},
			},
		},
		Timeouts: &resourceDefaultTimeouts,
	}
}

func dataSourceNSGRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	datacenterID := d.Get("datacenter_id").(string)
	id, idOk := d.GetOk("id")
	name, nameOk := d.GetOk("name")

	if idOk && nameOk {
		return diagutil.ToDiags(d, fmt.Errorf("id and name cannot be both specified at the same time"), nil)
	}
	if !idOk && !nameOk {
		return diagutil.ToDiags(d, fmt.Errorf("please provide either the network security group id or name"), nil)
	}

	location := d.Get("location").(string)

	client, err := meta.(bundleclient.SdkBundle).NewCloudAPIClient(ctx, location)
	if err != nil {
		return diag.FromErr(err)
	}

	if idOk {
		securityGroup, apiResponse, err := client.SecurityGroupsApi.DatacentersSecuritygroupsFindById(ctx, datacenterID, id.(string)).Depth(3).Execute()
		apiResponse.LogInfo()
		if err != nil {
			return diagutil.ToDiags(d, fmt.Errorf("an error occurred while retrieving network security group with ID: %s, %w", id.(string), err), &diagutil.ErrorContext{StatusCode: apiResponse.SafeStatusCode()})
		}
		return diagutil.ToDiags(d, setNSGDataSource(d, &securityGroup), nil)
	}

	securityGroups, apiResponse, err := client.SecurityGroupsApi.DatacentersSecuritygroupsGet(ctx, datacenterID).Depth(3).Execute()
	apiResponse.LogInfo()
	if err != nil {
		return diagutil.ToDiags(d, fmt.Errorf("an error occurred while retrieving network security groups: %w", err), &diagutil.ErrorContext{StatusCode: apiResponse.SafeStatusCode()})
	}
	var results []ionoscloud.SecurityGroup
	for _, sg := range securityGroups.Items {
		if strings.EqualFold(sg.Properties.Name, name.(string)) {
			results = append(results, sg)
		}
	}

	if len(results) == 0 {
		return diagutil.ToDiags(d, fmt.Errorf("no network security group found with the specified name = %s", name), nil)
	} else if len(results) > 1 {
		return diagutil.ToDiags(d, fmt.Errorf("more than one network security group found with the specified name = %s", name), nil)
	}
	securityGroup := results[0]
	if err := setNSGDataSource(d, &securityGroup); err != nil {
		return diagutil.ToDiags(d, err, nil)
	}

	return nil
}

func setNSGDataSource(d *schema.ResourceData, securityGroup *ionoscloud.SecurityGroup) error {
	if err := setNSGData(d, securityGroup); err != nil {
		return err
	}
	if securityGroup.Entities != nil && securityGroup.Entities.Rules != nil {
		rulesData := make([]map[string]any, 0, len(securityGroup.Entities.Rules.Items))
		for _, rule := range securityGroup.Entities.Rules.Items {
			ruleData := make(map[string]any)
			if rule.Id == nil {
				continue
			}
			ruleData["id"] = *rule.Id
			if rule.Properties.Name != nil {
				ruleData["name"] = *rule.Properties.Name
			}
			if rule.Properties.SourceMac.IsSet() {
				ruleData["source_mac"] = shared.ToValueDefault(rule.Properties.SourceMac.Get())
			}
			if rule.Properties.SourceIp.IsSet() {
				ruleData["source_ip"] = shared.ToValueDefault(rule.Properties.SourceIp.Get())
			}
			if rule.Properties.TargetIp.IsSet() {
				ruleData["target_ip"] = shared.ToValueDefault(rule.Properties.TargetIp.Get())
			}
			if rule.Properties.Protocol != nil {
				ruleData["protocol"] = *rule.Properties.Protocol
			}
			if rule.Properties.Type != nil {
				ruleData["type"] = *rule.Properties.Type
			}
			if rule.Properties.PortRangeStart != nil {
				ruleData["port_range_start"] = *rule.Properties.PortRangeStart
			}
			if rule.Properties.PortRangeEnd != nil {
				ruleData["port_range_end"] = *rule.Properties.PortRangeEnd
			}
			if rule.Properties.IcmpType.IsSet() && rule.Properties.IcmpType.Get() != nil {
				ruleData["icmp_type"] = strconv.Itoa(int(*rule.Properties.IcmpType.Get()))
			}
			if rule.Properties.IcmpCode.IsSet() && rule.Properties.IcmpCode.Get() != nil {
				ruleData["icmp_code"] = strconv.Itoa(int(*rule.Properties.IcmpCode.Get()))
			}
			rulesData = append(rulesData, ruleData)
		}
		if len(rulesData) > 0 {
			if err := d.Set("rules", rulesData); err != nil {
				return fmt.Errorf("error while setting rules property for NSG data source %s: %w", d.Id(), err)
			}
		}
	}
	return nil
}
