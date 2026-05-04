package ionoscloud

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	dns "github.com/ionos-cloud/sdk-go-bundle/products/dns/v2"
	"github.com/ionos-cloud/sdk-go-bundle/shared"

	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/bundleclient"
	diagutil "github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils/diags"
)

func dataSourceDNSZone() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceZoneRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:             schema.TypeString,
				Description:      "The ID of your DNS Zone.",
				ValidateDiagFunc: validation.ToDiagFunc(validation.IsUUID),
				Optional:         true,
				Computed:         true,
			},
			"name": {
				Type:        schema.TypeString,
				Description: "The name of your DNS Zone.",
				Optional:    true,
				Computed:    true,
			},
			"partial_match": {
				Type:        schema.TypeBool,
				Description: "Whether partial matching is allowed or not when using name argument.",
				Default:     false,
				Optional:    true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"enabled": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"nameservers": {
				Type:        schema.TypeList,
				Description: "A list of available name servers.",
				Computed:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func dataSourceZoneRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(bundleclient.SdkBundle).DNSClient
	idValue, idOk := d.GetOk("id")
	nameValue, nameOk := d.GetOk("name")
	partialMatch := d.Get("partial_match").(bool)

	id := idValue.(string)
	name := nameValue.(string)

	if idOk && nameOk {
		return diagutil.ToDiags(d, fmt.Errorf("ID and name cannot be both specified at the same time"), nil)
	}
	if !idOk && !nameOk {
		return diagutil.ToDiags(d, fmt.Errorf("please provide either the DNS Zone ID or name"), nil)
	}
	if partialMatch && !nameOk {
		return diagutil.ToDiags(d, fmt.Errorf("partial_match can only be used together with the name attribute"), nil)
	}

	var zone dns.ZoneRead
	var apiResponse *shared.APIResponse
	var err error

	if idOk {
		zone, apiResponse, err = client.GetZoneById(ctx, id)
		if err != nil {
			return diagutil.ToDiags(d, fmt.Errorf("an error occurred while fetching the DNS Zone with ID: %s, error: %w", id, err), &diagutil.ErrorContext{StatusCode: apiResponse.SafeStatusCode()})
		}
	} else {
		var results []dns.ZoneRead
		tflog.Info(ctx, "searching DNS zone by name", map[string]interface{}{"name": name, "partial_match": partialMatch})

		if partialMatch {
			// By default, when providing the name as a filter, for the GET requests, partial match
			// is true.
			zones, apiResponse, err := client.ListZones(ctx, name)
			if err != nil {
				return diagutil.ToDiags(d, fmt.Errorf("an error occurred while fetching DNS Zones: %w", err), &diagutil.ErrorContext{StatusCode: apiResponse.SafeStatusCode()})
			}
			results = zones.Items
		} else {
			// In order to have an exact name match, we must retrieve all the DNS Zones and then
			// build a list of exact matches based on the response, there is no other way since using
			// filter.zoneName only does a partial match.
			zones, apiResponse, err := client.ListZones(ctx, "")
			if err != nil {
				return diagutil.ToDiags(d, fmt.Errorf("an error occurred while fetching DNS Zones: %w", err), &diagutil.ErrorContext{StatusCode: apiResponse.SafeStatusCode()})
			}
			for _, zoneItem := range zones.Items {
				// Since each zone has a unique name, there is no need to keep on searching if
				// we already found the required zone.
				if len(results) == 1 {
					break
				}
				if strings.EqualFold(zoneItem.Properties.ZoneName, name) {
					results = append(results, zoneItem)
				}
			}
		}
		if results == nil || len(results) == 0 {
			return diagutil.ToDiags(d, fmt.Errorf("no DNS Zone found with the specified name = %s", name), nil)
		} else if len(results) > 1 {
			return diagutil.ToDiags(d, fmt.Errorf("more than one DNS Zone found with the specified name = %s", name), nil)
		} else {
			zone = results[0]
		}
	}

	if err := client.SetZoneData(d, zone); err != nil {
		return diagutil.ToDiags(d, err, nil)
	}

	return nil
}
