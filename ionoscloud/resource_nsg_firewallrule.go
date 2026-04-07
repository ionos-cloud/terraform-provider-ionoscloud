package ionoscloud

import (
	"context"
	"fmt"
	"log"

	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/bundleclient"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils/constant"
	diagutil "github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils/diags"

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
			"location": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The location of the resource. This field should be used only if you are also using a file configuration and should not be configured otherwise.",
				ForceNew:    true,
			},
		},
		Timeouts: &resourceDefaultTimeouts,
	}
}

func resourceNSGFirewallCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	firewall, diags := getFirewallData(d, "", false)
	if diags != nil {
		return diags
	}
	nsgID := d.Get("nsg_id").(string)
	dcID := d.Get("datacenter_id").(string)
	location := d.Get("location").(string)
	client, err := meta.(bundleclient.SdkBundle).NewCloudAPIClient(location)
	if err != nil {
		return diag.FromErr(err)
	}

	fw, apiResponse, err := client.SecurityGroupsApi.DatacentersSecuritygroupsFirewallrulesPost(ctx, dcID, nsgID).FirewallRule(firewall).Execute()
	logApiRequestTime(apiResponse)
	if err != nil {
		requestLocation, _ := apiResponse.Location()
		return diagutil.ToDiags(d, fmt.Errorf("an error occurred while creating a nsg firewall rule nsg id %s dcid %s : %w", nsgID, dcID, err), &diagutil.ErrorContext{RequestID: diagutil.ExtractRequestID(requestLocation), StatusCode: apiResponse.SafeStatusCode()})
	}
	d.SetId(*fw.Id)

	// Wait, catching any errors
	if errState := bundleclient.WaitForStateChange(ctx, meta, d, apiResponse, schema.TimeoutCreate); errState != nil {
		d.SetId("")
		requestLocation, _ := apiResponse.Location()
		return diagutil.ToDiags(d, fmt.Errorf("an error occurred while creating a nsg firewall rule dcId: %s nsg_id: %s %w", d.Get("datacenter_id").(string), d.Get("nsg_id").(string), errState), &diagutil.ErrorContext{Timeout: d.Timeout(schema.TimeoutCreate).String(), RequestID: diagutil.ExtractRequestID(requestLocation)})
	}

	return resourceNSGFirewallRead(ctx, d, meta)
}

func resourceNSGFirewallRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	location := d.Get("location").(string)
	client, err := meta.(bundleclient.SdkBundle).NewCloudAPIClient(location)
	if err != nil {
		return diag.FromErr(err)
	}

	fw, apiResponse, err := client.SecurityGroupsApi.DatacentersSecuritygroupsRulesFindById(ctx, d.Get("datacenter_id").(string),
		d.Get("nsg_id").(string), d.Id()).Execute()
	logApiRequestTime(apiResponse)

	if err != nil {
		if httpNotFound(apiResponse) {
			log.Printf("[DEBUG] could not find firewall rule datacenter_id = %s nsg_id = %s with id = %s", d.Get("datacenter_id").(string), d.Get("nsg_id").(string), d.Id())
			d.SetId("")
			return nil
		}
		return diagutil.ToDiags(d, fmt.Errorf("an error occurred while fetching a nsg firewall rule dcId: %s nsg_id: %s ID: %s %w",
			d.Get("datacenter_id").(string), d.Get("nsg_id").(string), d.Id(), err), nil)
	}

	if err := setFirewallData(d, &fw); err != nil {
		return diagutil.ToDiags(d, err, nil)
	}

	return nil
}

func resourceNSGFirewallUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	firewall, diags := getFirewallData(d, "", true)
	if diags != nil {
		return diags
	}
	nsgID := d.Get("nsg_id").(string)
	dcID := d.Get("datacenter_id").(string)
	location := d.Get("location").(string)

	client, err := meta.(bundleclient.SdkBundle).NewCloudAPIClient(location)
	if err != nil {
		return diag.FromErr(err)
	}

	_, apiResponse, err := client.SecurityGroupsApi.DatacentersSecuritygroupsRulesPatch(ctx, dcID, nsgID, d.Id()).Rule(*firewall.Properties).Execute()
	logApiRequestTime(apiResponse)

	if err != nil {
		requestLocation, _ := apiResponse.Location()
		return diagutil.ToDiags(d, fmt.Errorf("an error occurred while updating a nsg firewall rule: dcID %s nsgID %s %w", dcID, nsgID, err), &diagutil.ErrorContext{RequestID: diagutil.ExtractRequestID(requestLocation), StatusCode: apiResponse.SafeStatusCode()})
	}

	if errState := bundleclient.WaitForStateChange(ctx, meta, d, apiResponse, schema.TimeoutUpdate); errState != nil {
		requestLocation, _ := apiResponse.Location()
		return diagutil.ToDiags(d, fmt.Errorf("error getting state change for nsg firewall patch %w", errState), &diagutil.ErrorContext{Timeout: d.Timeout(schema.TimeoutUpdate).String(), RequestID: diagutil.ExtractRequestID(requestLocation)})
	}

	return resourceNSGFirewallRead(ctx, d, meta)
}

func resourceNSGFirewallDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	dcID := d.Get("datacenter_id").(string)
	nsgID := d.Get("nsg_id").(string)
	location := d.Get("location").(string)

	client, err := meta.(bundleclient.SdkBundle).NewCloudAPIClient(location)
	if err != nil {
		return diag.FromErr(err)
	}

	apiResponse, err := client.SecurityGroupsApi.
		DatacentersSecuritygroupsFirewallrulesDelete(
			ctx, dcID,
			nsgID, d.Id()).
		Execute()

	if err != nil {
		requestLocation, _ := apiResponse.Location()
		return diagutil.ToDiags(d, fmt.Errorf("an error occurred while deleting a nsg firewall rule ID %s nsgID %s %w", nsgID, dcID, err), &diagutil.ErrorContext{RequestID: diagutil.ExtractRequestID(requestLocation), StatusCode: apiResponse.SafeStatusCode()})
	}

	if errState := bundleclient.WaitForStateChange(ctx, meta, d, apiResponse, schema.TimeoutDelete); errState != nil {
		requestLocation, _ := apiResponse.Location()
		return diagutil.ToDiags(d, fmt.Errorf("error getting state change for firewall delete %w", errState), &diagutil.ErrorContext{Timeout: d.Timeout(schema.TimeoutDelete).String(), RequestID: diagutil.ExtractRequestID(requestLocation)})
	}

	d.SetId("")

	return nil
}

func resourceNSGFirewallImport(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	importID := d.Id()

	location, parts := splitImportID(importID, "/")
	if len(parts) != 3 {
		return nil, fmt.Errorf(
			"invalid import identifier: expected one of <location>:<datacenter-id>/<nsg-id>/<firewall-id> "+
				"or <datacenter-id>/<nsg-id>/<firewall-id>, got: %s", importID,
		)
	}

	if err := validateImportIDParts(parts); err != nil {
		return nil, diagutil.ToError(d, fmt.Errorf("failed validating import identifier %q: %w", importID, err), nil)
	}

	dcID := parts[0]
	nsgID := parts[1]
	firewallID := parts[2]

	client, err := meta.(bundleclient.SdkBundle).NewCloudAPIClient(location)
	if err != nil {
		return nil, err
	}

	fw, apiResponse, err := client.SecurityGroupsApi.DatacentersSecuritygroupsRulesFindById(ctx, dcID,
		nsgID, firewallID).Execute()
	logApiRequestTime(apiResponse)

	if err != nil {
		if httpNotFound(apiResponse) {
			d.SetId("")
			return nil, diagutil.ToError(d, fmt.Errorf("unable to find nsg firewall rule %q", firewallID), nil)
		}
		return nil, diagutil.ToError(d, fmt.Errorf("an error occurred while nsg retrieving firewall rule %q: %w ", firewallID, err), nil)
	}

	if err := d.Set("datacenter_id", dcID); err != nil {
		return nil, diagutil.ToError(d, err, nil)
	}
	if err := d.Set("nsg_id", nsgID); err != nil {
		return nil, diagutil.ToError(d, err, nil)
	}
	if err := d.Set("location", location); err != nil {
		return nil, err
	}

	if err := setFirewallData(d, &fw); err != nil {
		return nil, diagutil.ToError(d, err, nil)
	}

	return []*schema.ResourceData{d}, nil
}
