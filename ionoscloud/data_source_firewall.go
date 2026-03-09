package ionoscloud

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	ionoscloud "github.com/ionos-cloud/sdk-go/v6"

	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/bundleclient"
	diagutil "github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils/diags"
)

func dataSourceFirewall() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceFirewallRead,
		Schema: map[string]*schema.Schema{
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
			"datacenter_id": {
				Type:             schema.TypeString,
				Required:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringIsNotWhiteSpace),
			},
			"server_id": {
				Type:             schema.TypeString,
				Required:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringIsNotWhiteSpace),
			},
			"nic_id": {
				Type:             schema.TypeString,
				Required:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringIsNotWhiteSpace),
			},
			"location": {
				Type:        schema.TypeString,
				Description: "The location of the resource. This field should be used only if you are also using a file configuration and should not be configured otherwise.",
				Optional:    true,
			},
		},
		Timeouts: &resourceDefaultTimeouts,
	}
}

func dataSourceFirewallRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	location := d.Get("location").(string)
	client, err := meta.(bundleclient.SdkBundle).NewCloudAPIClient(location)
	if err != nil {
		return diag.FromErr(err)
	}

	datacenterId := d.Get("datacenter_id").(string)
	serverId := d.Get("server_id").(string)
	nicId := d.Get("nic_id").(string)

	id, idOk := d.GetOk("id")
	name, nameOk := d.GetOk("name")

	if idOk && nameOk {
		return diagutil.ToDiags(d, fmt.Errorf("id and name cannot be both specified in the same time"), nil)
	}
	if !idOk && !nameOk {
		return diagutil.ToDiags(d, fmt.Errorf("please provide either the firewall rule id or name"), nil)
	}
	var firewall ionoscloud.FirewallRule
	var apiResponse *ionoscloud.APIResponse

	if idOk {
		/* search by ID */
		firewall, apiResponse, err = client.FirewallRulesApi.DatacentersServersNicsFirewallrulesFindById(ctx, datacenterId, serverId, nicId, id.(string)).Execute()
		logApiRequestTime(apiResponse)
		if err != nil {
			return diagutil.ToDiags(d, fmt.Errorf("an error occurred while fetching the firewall rule %s: %w", id.(string), err), &diagutil.DiagsOpts{StatusCode: apiResponse.StatusCode})
		}
	} else {
		/* search by name */
		var firewalls ionoscloud.FirewallRules

		firewalls, apiResponse, err := client.FirewallRulesApi.DatacentersServersNicsFirewallrulesGet(ctx, datacenterId, serverId, nicId).Depth(1).Execute()
		logApiRequestTime(apiResponse)

		if err != nil {
			return diagutil.ToDiags(d, fmt.Errorf("an error occurred while fetching backup unit: %w", err), &diagutil.DiagsOpts{StatusCode: apiResponse.StatusCode})
		}

		var results []ionoscloud.FirewallRule

		if firewalls.Items != nil {
			for _, fr := range *firewalls.Items {
				if fr.Properties != nil && fr.Properties.Name != nil && *fr.Properties.Name == name.(string) {
					tmpFirewall, apiResponse, err := client.FirewallRulesApi.DatacentersServersNicsFirewallrulesFindById(ctx, datacenterId, serverId, nicId, *fr.Id).Execute()
					logApiRequestTime(apiResponse)
					if err != nil {
						return diagutil.ToDiags(d, fmt.Errorf("an error occurred while fetching firewall rule with ID %s: %w", *fr.Id, err), &diagutil.DiagsOpts{StatusCode: apiResponse.StatusCode})
					}
					results = append(results, tmpFirewall)
				}
			}
		}

		if results == nil || len(results) == 0 {
			return diagutil.ToDiags(d, fmt.Errorf("no firewall rule found with the specified name = %s", name), nil)
		} else if len(results) > 1 {
			return diagutil.ToDiags(d, fmt.Errorf("more than one firewall rule found with the specified criteria name = %s", name), nil)
		} else {
			firewall = results[0]
		}

	}

	if err := setFirewallData(d, &firewall); err != nil {
		return diagutil.ToDiags(d, err, nil)
	}

	return nil
}
