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
				Type:     schema.TypeInt,
				Optional: true,
				ValidateDiagFunc: validation.ToDiagFunc(func(v interface{}, k string) (ws []string, errors []error) {
					if v.(int) < 1 && v.(int) > 65534 {
						errors = append(errors, fmt.Errorf("port start range must be between 1 and 65534"))
					}
					return
				}),
			},
			"port_range_end": {
				Type:     schema.TypeInt,
				Optional: true,
				ValidateDiagFunc: validation.ToDiagFunc(func(v interface{}, k string) (ws []string, errors []error) {
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
	fw, apiResponse, err := client.SecurityGroupsApi.DatacentersSecuritygroupsFirewallrulesPost(ctx, d.Get("datacenter_id").(string), d.Get("nsg_id").(string)).FirewallRule(firewall).Execute()
	logApiRequestTime(apiResponse)

	if err != nil {
		diags := diag.FromErr(fmt.Errorf("an error occured while creating a nsg firewall rule: %w", err))
		return diags
	}
	d.SetId(*fw.Id)

	// Wait, catching any errors
	if errState := cloudapi.WaitForStateChange(ctx, meta, d, apiResponse, schema.TimeoutCreate); errState != nil {
		if cloudapi.IsRequestFailed(errState) {
			log.Printf("[DEBUG] firewall resource failed to be created")
			d.SetId("")
		}
		return diag.FromErr(fmt.Errorf("an error occured while creating a nsg firewall rule dcId: %s _id: %s %w", d.Get("datacenter_id").(string), d.Get("nsg_id").(string), errState))
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
		diags := diag.FromErr(fmt.Errorf("an error occured while fetching a firewall rule dcId: %s nsg_id: %s ID: %s %s",
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
	_, apiResponse, err := client.SecurityGroupsApi.DatacentersSecuritygroupsRulesPatch(ctx, d.Get("datacenter_id").(string), d.Get("nsg_id").(string), d.Id()).Rule(*firewall.Properties).Execute()
	logApiRequestTime(apiResponse)

	if err != nil {
		diags := diag.FromErr(fmt.Errorf("an error occured while updating a nsg firewall rule ID %s %w", d.Id(), err))
		return diags
	}

	if errState := cloudapi.WaitForStateChange(ctx, meta, d, apiResponse, schema.TimeoutUpdate); errState != nil {
		return diag.FromErr(fmt.Errorf("error getting state change for nsg firewall patch %w", errState))
	}

	return resourceNSGFirewallRead(ctx, d, meta)
}

func resourceNSGFirewallDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(services.SdkBundle).CloudApiClient

	apiResponse, err := client.SecurityGroupsApi.
		DatacentersSecuritygroupsFirewallrulesDelete(
			ctx, d.Get("datacenter_id").(string),
			d.Get("nsg_id").(string), d.Id()).
		Execute()

	if err != nil {
		diags := diag.FromErr(fmt.Errorf("an error occured while deleting a nsg firewall rule ID %s %w", d.Id(), err))
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

	dcId := parts[0]
	nsgId := parts[1]
	firewallId := parts[2]

	fw, apiResponse, err := client.SecurityGroupsApi.DatacentersSecuritygroupsRulesFindById(ctx, dcId,
		nsgId, firewallId).Execute()
	logApiRequestTime(apiResponse)

	if err != nil {
		if httpNotFound(apiResponse) {
			d.SetId("")
			return nil, fmt.Errorf("unable to find firewall rule %q", firewallId)
		}
		return nil, fmt.Errorf("an error occured while retrieving firewall rule %q: %q ", firewallId, err)
	}

	if err := d.Set("datacenter_id", dcId); err != nil {
		return nil, err
	}
	if err := d.Set("nsg_id", nsgId); err != nil {
		return nil, err
	}

	if err := setFirewallData(d, &fw); err != nil {
		return nil, err
	}

	return []*schema.ResourceData{d}, nil
}
