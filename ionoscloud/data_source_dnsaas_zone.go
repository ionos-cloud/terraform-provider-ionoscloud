package ionoscloud

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	dnsaas "github.com/ionos-cloud/sdk-go-dnsaas"
	"log"
	"strings"
)

func dataSourceDNSaaSZone() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceZoneRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Description: "The ID of your DNS Zone.",
				Optional:    true,
			},
			"name": {
				Type:        schema.TypeString,
				Description: "The name of your DNS Zone.",
				Optional:    true,
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
		},
	}
}

func dataSourceZoneRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(SdkBundle).DNSaaSClient
	idValue, idOk := d.GetOk("id")
	nameValue, nameOk := d.GetOk("name")

	id := idValue.(string)
	name := nameValue.(string)

	if idOk && nameOk {
		return diag.FromErr(fmt.Errorf("ID and name cannot be both specified at the same time"))
	}
	if !idOk && !nameOk {
		return diag.FromErr(fmt.Errorf("please provide either the DNS Zone ID or name"))
	}

	var zone dnsaas.ZoneResponse
	var err error

	if idOk {
		zone, _, err = client.GetZoneById(ctx, id)
		if err != nil {
			return diag.FromErr(fmt.Errorf("an error occured while fetching the DNS Zone with ID: %s, error: %w", id, err))
		}
	} else {
		var results []dnsaas.ZoneResponse
		partialMatch := d.Get("partial_match").(bool)
		log.Printf("[INFO] Populating data source for DNS Zone using name %s and partial_match %t", name, partialMatch)

		if partialMatch {
			// By default, when providing the name as a filter, for the GET requests, partial match
			// is true.
			zones, _, err := client.ListZones(ctx, name)
			if err != nil {
				return diag.FromErr(fmt.Errorf("an error occured while fetching DNS Zones: %w", err))
			}
			if zones.Items != nil {
				results = *zones.Items
			} else {
				return diag.FromErr(fmt.Errorf("expected items representing DNS Zones, got 'nil' instead"))
			}
		} else {
			// In order to have an exact name match, we must retrieve all the DNS Zones and then
			// build a list of exact matches based on the response, there is no other way since using
			// filter.zoneName only does a partial match.
			zones, _, err := client.ListZones(ctx, "")
			if err != nil {
				return diag.FromErr(fmt.Errorf("an error occured while fetching DNS Zones: %w", err))
			}
			if zones.Items != nil {
				for _, zoneItem := range *zones.Items {
					// Since each zone has a unique name, there is no need to keep on searching if
					// we already found the required zone.
					if len(results) == 1 {
						break
					}
					if zoneItem.Properties != nil && zoneItem.Properties.ZoneName != nil && strings.EqualFold(*zoneItem.Properties.ZoneName, name) {
						results = append(results, zoneItem)
					}
				}
			} else {
				return diag.FromErr(fmt.Errorf("expected items representing DNS Zones, got 'nil' instead"))
			}
		}
		if results == nil || len(results) == 0 {
			return diag.FromErr(fmt.Errorf("no DNS Zone found with the specified name = %s", name))
		} else if len(results) > 1 {
			return diag.FromErr(fmt.Errorf("more than one DNS Zone found with the specified name = %s", name))
		} else {
			zone = results[0]
		}
	}

	if err := client.SetZoneData(d, zone); err != nil {
		return diag.FromErr(err)
	}

	return nil
}
