package ionoscloud

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/cloudapi"

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
	client := meta.(services.SdkBundle).CloudApiClient

	firewall, diags := getFirewallData(d, "", false)
	if diags != nil {
		return diags
	}
	nsgID := d.Get("nsg_id").(string)
	dcID := d.Get("datacenter_id").(string)
	fw, apiResponse, err := client.SecurityGroupsApi.DatacentersSecuritygroupsFirewallrulesPost(ctx, dcID, nsgID).FirewallRule(firewall).Execute()
	logApiRequestTime(apiResponse)
	if err != nil {
		diags := diag.FromErr(fmt.Errorf("an error occurred while creating a nsg firewall rule nsg id %s dcid %s : %w", nsgID, dcID, err))
		return diags
	}
	d.SetId(*fw.Id)

	// Wait, catching any errors
	if errState := cloudapi.WaitForStateChange(ctx, meta, d, apiResponse, schema.TimeoutCreate); errState != nil {
		d.SetId("")
		return diag.FromErr(fmt.Errorf("an error occurred while creating a nsg firewall rule dcId: %s nsg_id: %s %w", d.Get("datacenter_id").(string), d.Get("nsg_id").(string), errState))
	}

	return resourceNSGFirewallRead(ctx, d, meta)
}

func resourceNSGFirewallRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(services.SdkBundle).CloudApiClient

	fw, apiResponse, err := client.SecurityGroupsApi.DatacentersSecuritygroupsRulesFindById(ctx, d.Get("datacenter_id").(string),
		d.Get("nsg_id").(string), d.Id()).Execute()
	logApiRequestTime(apiResponse)

	if err != nil {
		if httpNotFound(apiResponse) {
			log.Printf("[DEBUG] could not find firewall rule datacenter_id = %s nsg_id = %s with id = %s", d.Get("datacenter_id").(string), d.Get("nsg_id").(string), d.Id())
			d.SetId("")
			return nil
		}
		diags := diag.FromErr(fmt.Errorf("an error occurred while fetching a nsg firewall rule dcId: %s nsg_id: %s ID: %s %w",
			d.Get("datacenter_id").(string), d.Get("nsg_id").(string), d.Id(), err))
		return diags
	}

	if err := setFirewallData(d, &fw); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceNSGFirewallUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(services.SdkBundle).CloudApiClient

	firewall, diags := getFirewallData(d, "", true)
	if diags != nil {
		return diags
	}
	nsgID := d.Get("nsg_id").(string)
	dcID := d.Get("datacenter_id").(string)
	_, apiResponse, err := client.SecurityGroupsApi.DatacentersSecuritygroupsRulesPatch(ctx, dcID, nsgID, d.Id()).Rule(*firewall.Properties).Execute()
	logApiRequestTime(apiResponse)

	if err != nil {
		diags := diag.FromErr(fmt.Errorf("an error occurred while updating a nsg firewall rule ID %s dcID %s nsgID %s %w", d.Id(), dcID, nsgID, err))
		return diags
	}

	if errState := cloudapi.WaitForStateChange(ctx, meta, d, apiResponse, schema.TimeoutUpdate); errState != nil {
		return diag.FromErr(fmt.Errorf("error getting state change for nsg firewall patch %w", errState))
	}

	return resourceNSGFirewallRead(ctx, d, meta)
}

func resourceNSGFirewallDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(services.SdkBundle).CloudApiClient
	dcID := d.Get("datacenter_id").(string)
	nsgID := d.Get("nsg_id").(string)
	apiResponse, err := client.SecurityGroupsApi.
		DatacentersSecuritygroupsFirewallrulesDelete(
			ctx, dcID,
			nsgID, d.Id()).
		Execute()

	if err != nil {
		diags := diag.FromErr(fmt.Errorf("an error occurred while deleting a nsg firewall rule ID %s nsgID %s dcID %s %w", nsgID, dcID, d.Id(), err))
		return diags
	}

	if errState := cloudapi.WaitForStateChange(ctx, meta, d, apiResponse, schema.TimeoutDelete); errState != nil {
		return diag.FromErr(fmt.Errorf("error getting state change for firewall delete %w", errState))
	}

	d.SetId("")

	return nil
}

func resourceNSGFirewallImport(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {

	client := meta.(services.SdkBundle).CloudApiClient

	parts := strings.Split(d.Id(), "/")
	if len(parts) != 3 || parts[0] == "" || parts[1] == "" || parts[2] == "" {
		return nil, fmt.Errorf("invalid import id %q. Expecting {datacenter}/{nsg}/{firewall}", d.Id())
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
			return nil, fmt.Errorf("unable to find nsg firewall rule %q", firewallID)
		}
		return nil, fmt.Errorf("an error occurred while nsg retrieving firewall rule %q: %w ", firewallID, err)
	}

	if err := d.Set("datacenter_id", dcID); err != nil {
		return nil, err
	}
	if err := d.Set("nsg_id", nsgID); err != nil {
		return nil, err
	}

	if err := setFirewallData(d, &fw); err != nil {
		return nil, err
	}

	return []*schema.ResourceData{d}, nil
}
