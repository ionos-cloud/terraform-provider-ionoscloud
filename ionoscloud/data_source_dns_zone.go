package ionoscloud

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	dns "github.com/ionos-cloud/sdk-go-bundle/products/dns/v2"
	"github.com/ionos-cloud/sdk-go-bundle/shared"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/bundleclient"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils"
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

func dataSourceZoneRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(bundleclient.SdkBundle).DNSClient
	idValue, idOk := d.GetOk("id")
	nameValue, nameOk := d.GetOk("name")
	partialMatch := d.Get("partial_match").(bool)

	id := idValue.(string)
	name := nameValue.(string)

	if idOk && nameOk {
		return utils.ToDiags(d, "ID and name cannot be both specified at the same time", nil)
	}
	if !idOk && !nameOk {
		return utils.ToDiags(d, "please provide either the DNS Zone ID or name", nil)
	}
	if partialMatch && !nameOk {
		return utils.ToDiags(d, "partial_match can only be used together with the name attribute", nil)
	}

	var zone dns.ZoneRead
	var apiResponse *shared.APIResponse
	var err error

	if idOk {
		zone, apiResponse, err = client.GetZoneById(ctx, id)
		if err != nil {
			return utils.ToDiags(d, fmt.Sprintf("an error occurred while fetching the DNS Zone with ID: %s, error: %s", id, err), &utils.DiagsOpts{StatusCode: apiResponse.StatusCode})
		}
	} else {
		var results []dns.ZoneRead
		log.Printf("[INFO] Populating data source for DNS Zone using name %s and partial_match %t", name, partialMatch)

		if partialMatch {
			// By default, when providing the name as a filter, for the GET requests, partial match
			// is true.
			zones, apiResponse, err := client.ListZones(ctx, name)
			if err != nil {
				return utils.ToDiags(d, fmt.Sprintf("an error occurred while fetching DNS Zones: %s", err), &utils.DiagsOpts{StatusCode: apiResponse.StatusCode})
			}
			results = zones.Items
		} else {
			// In order to have an exact name match, we must retrieve all the DNS Zones and then
			// build a list of exact matches based on the response, there is no other way since using
			// filter.zoneName only does a partial match.
			zones, apiResponse, err := client.ListZones(ctx, "")
			if err != nil {
				return utils.ToDiags(d, fmt.Sprintf("an error occurred while fetching DNS Zones: %s", err), &utils.DiagsOpts{StatusCode: apiResponse.StatusCode})
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
			return utils.ToDiags(d, fmt.Sprintf("no DNS Zone found with the specified name = %s", name), nil)
		} else if len(results) > 1 {
			return utils.ToDiags(d, fmt.Sprintf("more than one DNS Zone found with the specified name = %s", name), nil)
		} else {
			zone = results[0]
		}
	}

	if err := client.SetZoneData(d, zone); err != nil {
		return utils.ToDiags(d, err.Error(), nil)
	}

	return nil
}
