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

	securityGroup, apiResponse, err := client.SecurityGroupsApi.DatacentersSecuritygroupsPost(ctx, datacenterID).SecurityGroup(sg).Execute()
	apiResponse.LogInfo()
	if err != nil {
		diags := diag.FromErr(fmt.Errorf("an error occured while creating a security group: %w", err))
		return diags
	}
	// todo - if needed
	// if errState := cloudapi.WaitForStateChange(ctx, meta, d, apiResponse, schema.TimeoutCreate); errState != nil {
	//
	// }

	rules := d.Get("firewall_rules").([]any)
	for i := range rules {
		properties := ionoscloud.FirewallruleProperties{}
		rulePath := fmt.Sprintf("firewall_rules.%d.", i)
		if v, ok := d.GetOk(rulePath + "name"); ok {
			vStr := v.(string)
			properties.Name = &vStr
		}
		if v, ok := d.GetOk(rulePath + "protocol"); ok {
			vStr := v.(string)
			properties.Protocol = &vStr
		}

		rule := ionoscloud.FirewallRule{Properties: &properties}
		_, apiResponse, err := client.SecurityGroupsApi.DatacentersSecuritygroupsFirewallrulesPost(ctx, datacenterID, *securityGroup.Id).FirewallRule(rule).Execute()
		// if errState := cloudapi.WaitForStateChange(ctx, meta, d, apiResponse, schema.TimeoutCreate); errState != nil {
		//
		// }
		apiResponse.LogInfo()
		if err != nil {
			diags := diag.FromErr(fmt.Errorf("an error occured while creating a security group firewall rule: %w, sgID: %s", err, *securityGroup.Id))
			return diags
		}
	}

	return nil
}
func resourceNetworkSecurityGroupRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}
func resourceNetworkSecurityGroupUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}
func resourceNetworkSecurityGroupDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}
