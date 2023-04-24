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

func dataSourceDNSaaSRecord() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceRecordRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Description: "The ID of your DNS Record.",
				Optional:    true,
			},
			"name": {
				Type:        schema.TypeString,
				Description: "The name of your DNS Record.",
				Optional:    true,
			},
			"partial_match": {
				Type:        schema.TypeBool,
				Description: "Whether partial matching is allowed or not when using name argument.",
				Default:     false,
				Optional:    true,
			},
			"zone_id": {
				Type:        schema.TypeString,
				Description: "The UUID of an existing DNS Zone",
				Required:    true,
			},
			"type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"content": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"ttl": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"enabled": {
				Type:     schema.TypeBool,
				Computed: true,
			},
		},
		Timeouts: &resourceDefaultTimeouts,
	}
}

func dataSourceRecordRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(SdkBundle).DNSaaSClient
	partialMatch := d.Get("partial_match").(bool)
	zoneId := d.Get("zone_id").(string)
	idValue, idOk := d.GetOk("id")
	nameValue, nameOk := d.GetOk("name")
	recordId := idValue.(string)
	recordName := nameValue.(string)

	if idOk && nameOk {
		return diag.FromErr(fmt.Errorf("ID and name cannot be both specified at the same time"))
	}
	if !idOk && !nameOk {
		return diag.FromErr(fmt.Errorf("please provide either the DNS Record ID or name"))
	}
	if partialMatch && !nameOk {
		return diag.FromErr(fmt.Errorf("partial_match can only be used together with the name attribute"))
	}

	var record dnsaas.RecordResponse
	var err error

	if idOk {
		record, _, err = client.GetRecordById(ctx, zoneId, recordId)
		if err != nil {
			return diag.FromErr(fmt.Errorf("an error occured while fetching the DNS Record with ID: %s, DNS Zone ID: %s, error: %w", recordId, zoneId, err))
		}
	} else {
		var results []dnsaas.RecordResponse
		log.Printf("[INFO] Populating data source for DNS Record using name: %s and partial_match: %t", recordName, partialMatch)
		if partialMatch {
			// By default, when providing the name as a filter, for the GET requests, partial match
			// is true.
			records, _, err := client.ListRecords(ctx, zoneId, recordName)
			if err != nil {
				return diag.FromErr(fmt.Errorf("an error occured while fetching DNS Records: %w", err))
			}
			if records.Items != nil {
				results = *records.Items
			} else {
				return diag.FromErr(fmt.Errorf("expected items representing DNS Records, got 'nil' instead"))
			}
		} else {
			// In order to have an exact name match, we must retrieve all the DNS Records and then
			// build a list of exact matches based on the response, there is no other way since using
			// filter.name only does a partial match.
			records, _, err := client.ListRecords(ctx, zoneId, "")
			if err != nil {
				return diag.FromErr(fmt.Errorf("an error occured while fetching DNS Records: %w", err))
			}
			if records.Items != nil {
				for _, recordItem := range *records.Items {
					// Since each record has a unique name, there is no need to keep on searching if
					// we already found the required record.
					if len(results) == 1 {
						break
					}
					if recordItem.Properties != nil && recordItem.Properties.Name != nil && strings.EqualFold(*recordItem.Properties.Name, recordName) {
						results = append(results, recordItem)
					}
				}
			} else {
				return diag.FromErr(fmt.Errorf("expected items representing DNS Records, got 'nil' instead"))
			}
		}
		if results == nil || len(results) == 0 {
			return diag.FromErr(fmt.Errorf("no DNS Record found with the specified name = %s", recordName))
		} else if len(results) > 1 {
			return diag.FromErr(fmt.Errorf("more than one DNS Record found with the specified name = %s", recordName))
		} else {
			record = results[0]
		}
	}

	if err := client.SetRecordData(d, record); err != nil {
		return diag.FromErr(err)
	}

	return nil
}
