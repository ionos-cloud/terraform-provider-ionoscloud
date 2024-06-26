package ionoscloud

import (
	"context"
	"fmt"
	"log"
	"strings"

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
		Importer: &schema.ResourceImporter{
			StateContext: resourceNetworkSecurityGroupImport,
		},
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

	return diag.FromErr(setNetworkSecurityGroupData(d, &securityGroup))
}

func resourceNetworkSecurityGroupRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(services.SdkBundle).CloudApiClient
	datacenterID := d.Get("datacenter_id").(string)

	securityGroup, apiResponse, err := client.SecurityGroupsApi.DatacentersSecuritygroupsFindById(ctx, datacenterID, d.Id()).Depth(2).Execute()
	apiResponse.LogInfo()
	if err != nil {
		return diag.FromErr(fmt.Errorf("an error occured while retrieving a network security group: %w", err))
	}

	if err := setNetworkSecurityGroupData(d, &securityGroup); err != nil {
		return diag.FromErr(err)
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

func resourceNetworkSecurityGroupImport(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	client := meta.(services.SdkBundle).CloudApiClient

	parts := strings.Split(d.Id(), "/")

	if len(parts) < 2 || parts[0] == "" || parts[1] == "" {
		return nil, fmt.Errorf("invalid import id %q. Expecting {datacenter UUID}/{nsg UUID}", d.Id())
	}

	datacenterId := parts[0]
	nsgId := parts[1]

	nsg, apiResponse, err := client.SecurityGroupsApi.DatacentersSecuritygroupsFindById(ctx, datacenterId, nsgId).Execute()
	logApiRequestTime(apiResponse)

	if err != nil {
		if httpNotFound(apiResponse) {
			d.SetId("")
			return nil, fmt.Errorf("unable to find Network Security Group %q", nsgId)
		}
		return nil, fmt.Errorf("an error occured while retrieving the Network Security Group %q, %q", d.Id(), err)
	}

	log.Printf("[INFO] Datacenter found: %+v", nsg)
	if err = d.Set("datacenter_id", datacenterId); err != nil {
		return nil, err
	}
	if err = setNetworkSecurityGroupData(d, &nsg); err != nil {
		return nil, err
	}

	return []*schema.ResourceData{d}, nil
}

func setNetworkSecurityGroupData(d *schema.ResourceData, securityGroup *ionoscloud.SecurityGroup) error {

	if securityGroup.Id != nil {
		d.SetId(*securityGroup.Id)
	}

	if securityGroup.Properties != nil {
		if securityGroup.Properties.Name != nil {
			err := d.Set("name", *securityGroup.Properties.Name)
			if err != nil {
				return fmt.Errorf("error while setting Network Security Group name  %s: %w", d.Id(), err)
			}
		}
		if securityGroup.Properties.Description != nil {
			err := d.Set("description", *securityGroup.Properties.Description)
			if err != nil {
				return fmt.Errorf("error while setting Network Security Group description  %s: %w", d.Id(), err)
			}
		}
	}
	var ruleIDs []string
	if securityGroup.Entities != nil {
		if securityGroup.Entities.Rules != nil && securityGroup.Entities.Rules.Items != nil {
			ruleIDs = make([]string, 0, len(*securityGroup.Entities.Rules.Items))
			for _, rule := range *securityGroup.Entities.Rules.Items {
				ruleIDs = append(ruleIDs, *rule.Id)
			}
		}
	}
	if err := d.Set("rule_ids", ruleIDs); err != nil {
		return fmt.Errorf("error while setting rule_ids property for NetworkSecurityGroup  %s: %w", d.Id(), err)
	}
	return nil
}
