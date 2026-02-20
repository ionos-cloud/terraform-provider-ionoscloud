package ionoscloud

import (
	"context"
	"fmt"
	"log"
	"strings"

	bundleclient "github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/bundleclient"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils/constant"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceNSGFirewallRule() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceNSGFirewallCreate,
		ReadContext:   resourceNSGFirewallRead,
		UpdateContext: resourceNSGFirewallUpdate,
		DeleteContext: resourceNSGFirewallDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceNSGFirewallImport,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"protocol": {
				Type:             schema.TypeString,
				Required:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice(constant.FirewallProtocolEnum, false)),
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
				Type:             schema.TypeInt,
				Optional:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.IntBetween(1, 65534)),
			},
			"port_range_end": {
				Type:             schema.TypeInt,
				Optional:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.IntBetween(1, 65534)),
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
			"datacenter_id": {
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.IsUUID),
			},
			"nsg_id": {
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.IsUUID),
			},
		},
		Timeouts: &resourceDefaultTimeouts,
	}
}

func resourceNSGFirewallCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(bundleclient.SdkBundle).CloudApiClient

	firewall, diags := getFirewallData(d, "", false)
	if diags != nil {
		return diags
	}
	nsgID := d.Get("nsg_id").(string)
	dcID := d.Get("datacenter_id").(string)
	fw, apiResponse, err := client.SecurityGroupsApi.DatacentersSecuritygroupsFirewallrulesPost(ctx, dcID, nsgID).FirewallRule(firewall).Execute()
	logApiRequestTime(apiResponse)
	if err != nil {
		requestLocation, _ := apiResponse.Location()
		return utils.ToDiags(d, fmt.Sprintf("an error occurred while creating a nsg firewall rule nsg id %s dcid %s : %s", nsgID, dcID, err), &utils.DiagsOpts{RequestLocation: requestLocation, StatusCode: apiResponse.StatusCode})
	}
	d.SetId(*fw.Id)

	// Wait, catching any errors
	if errState := bundleclient.WaitForStateChange(ctx, meta, d, apiResponse, schema.TimeoutCreate); errState != nil {
		d.SetId("")
		return utils.ToDiags(d, fmt.Sprintf("an error occurred while creating a nsg firewall rule dcId: %s nsg_id: %s %s", d.Get("datacenter_id").(string), d.Get("nsg_id").(string), errState), &utils.DiagsOpts{Timeout: schema.TimeoutCreate})
	}

	return resourceNSGFirewallRead(ctx, d, meta)
}

func resourceNSGFirewallRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(bundleclient.SdkBundle).CloudApiClient

	fw, apiResponse, err := client.SecurityGroupsApi.DatacentersSecuritygroupsRulesFindById(ctx, d.Get("datacenter_id").(string),
		d.Get("nsg_id").(string), d.Id()).Execute()
	logApiRequestTime(apiResponse)

	if err != nil {
		if httpNotFound(apiResponse) {
			log.Printf("[DEBUG] could not find firewall rule datacenter_id = %s nsg_id = %s with id = %s", d.Get("datacenter_id").(string), d.Get("nsg_id").(string), d.Id())
			d.SetId("")
			return nil
		}
		return utils.ToDiags(d, fmt.Sprintf("an error occurred while fetching a nsg firewall rule dcId: %s nsg_id: %s ID: %s %s",
			d.Get("datacenter_id").(string), d.Get("nsg_id").(string), d.Id(), err), nil)
	}

	if err := setFirewallData(d, &fw); err != nil {
		return utils.ToDiags(d, err.Error(), nil)
	}

	return nil
}

func resourceNSGFirewallUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(bundleclient.SdkBundle).CloudApiClient

	firewall, diags := getFirewallData(d, "", true)
	if diags != nil {
		return diags
	}
	nsgID := d.Get("nsg_id").(string)
	dcID := d.Get("datacenter_id").(string)
	_, apiResponse, err := client.SecurityGroupsApi.DatacentersSecuritygroupsRulesPatch(ctx, dcID, nsgID, d.Id()).Rule(*firewall.Properties).Execute()
	logApiRequestTime(apiResponse)

	if err != nil {
		requestLocation, _ := apiResponse.Location()
		return utils.ToDiags(d, fmt.Sprintf("an error occurred while updating a nsg firewall rule: dcID %s nsgID %s %s", dcID, nsgID, err), &utils.DiagsOpts{RequestLocation: requestLocation, StatusCode: apiResponse.StatusCode})
	}

	if errState := bundleclient.WaitForStateChange(ctx, meta, d, apiResponse, schema.TimeoutUpdate); errState != nil {
		return utils.ToDiags(d, fmt.Sprintf("error getting state change for nsg firewall patch %s", errState), &utils.DiagsOpts{Timeout: schema.TimeoutUpdate})
	}

	return resourceNSGFirewallRead(ctx, d, meta)
}

func resourceNSGFirewallDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(bundleclient.SdkBundle).CloudApiClient
	dcID := d.Get("datacenter_id").(string)
	nsgID := d.Get("nsg_id").(string)
	apiResponse, err := client.SecurityGroupsApi.
		DatacentersSecuritygroupsFirewallrulesDelete(
			ctx, dcID,
			nsgID, d.Id()).
		Execute()

	if err != nil {
		requestLocation, _ := apiResponse.Location()
		return utils.ToDiags(d, fmt.Sprintf("an error occurred while deleting a nsg firewall rule ID %s nsgID %s %s", nsgID, dcID, err), &utils.DiagsOpts{RequestLocation: requestLocation, StatusCode: apiResponse.StatusCode})
	}

	if errState := bundleclient.WaitForStateChange(ctx, meta, d, apiResponse, schema.TimeoutDelete); errState != nil {
		return utils.ToDiags(d, fmt.Sprintf("error getting state change for firewall delete %s", errState), &utils.DiagsOpts{Timeout: schema.TimeoutDelete})
	}

	d.SetId("")

	return nil
}

func resourceNSGFirewallImport(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {

	client := meta.(bundleclient.SdkBundle).CloudApiClient

	parts := strings.Split(d.Id(), "/")
	if len(parts) != 3 || parts[0] == "" || parts[1] == "" || parts[2] == "" {
		return nil, utils.ToError(d, "invalid import. Expecting {datacenter}/{nsg}/{firewall}", nil)
	}

	dcID := parts[0]
	nsgID := parts[1]
	firewallID := parts[2]

	fw, apiResponse, err := client.SecurityGroupsApi.DatacentersSecuritygroupsRulesFindById(ctx, dcID,
		nsgID, firewallID).Execute()
	logApiRequestTime(apiResponse)

	if err != nil {
		if httpNotFound(apiResponse) {
			d.SetId("")
			return nil, utils.ToError(d, fmt.Sprintf("unable to find nsg firewall rule %q", firewallID), nil)
		}
		return nil, utils.ToError(d, fmt.Sprintf("an error occurred while nsg retrieving firewall rule %q: %s ", firewallID, err), nil)
	}

	if err := d.Set("datacenter_id", dcID); err != nil {
		return nil, utils.ToError(d, err.Error(), nil)
	}
	if err := d.Set("nsg_id", nsgID); err != nil {
		return nil, utils.ToError(d, err.Error(), nil)
	}

	if err := setFirewallData(d, &fw); err != nil {
		return nil, utils.ToError(d, err.Error(), nil)
	}

	return []*schema.ResourceData{d}, nil
}
