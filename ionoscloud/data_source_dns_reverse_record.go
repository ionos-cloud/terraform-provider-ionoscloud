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
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/bundleclient"
)

func dataSourceDNSReverseRecord() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceReverseRecordRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:             schema.TypeString,
				Description:      "The ID of your DNS Reverse Record.",
				ValidateDiagFunc: validation.ToDiagFunc(validation.IsUUID),
				Optional:         true,
				Computed:         true,
			},
			"name": {
				Type:        schema.TypeString,
				Description: "The name of your DNS Reverse Record.",
				Optional:    true,
				Computed:    true,
			},
			"partial_match": {
				Type:        schema.TypeBool,
				Description: "Whether partial matching is allowed or not when using name argument.",
				Default:     false,
				Optional:    true,
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
			"priority": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"enabled": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"fqdn": {
				Type:        schema.TypeString,
				Description: "Fully qualified domain name",
				Computed:    true,
			},
		},
		Timeouts: &resourceDefaultTimeouts,
	}
}

func dataSourceReverseRecordRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(bundleclient.SdkBundle).DNSClient
	partialMatch := d.Get("partial_match").(bool)
	idValue, idOk := d.GetOk("id")
	nameValue, nameOk := d.GetOk("name")
	recordId := idValue.(string)
	recordName := nameValue.(string)

	if idOk && nameOk {
		return diag.FromErr(fmt.Errorf("ID and name cannot be both specified at the same time"))
	}
	if !idOk && !nameOk {
		return diag.FromErr(fmt.Errorf("please provide either the DNS Reverse Record ID or name"))
	}
	if partialMatch && !nameOk {
		return diag.FromErr(fmt.Errorf("partial_match can only be used together with the name attribute"))
	}

	var record dns.ReverseRecordRead
	var err error

	if idOk {
		record, _, err = client.GetReverseRecordById(ctx, recordId)
		if err != nil {
			return diag.FromErr(fmt.Errorf("an error occurred while fetching the DNS Reverse Record with ID: %s, error: %w", recordId, err))
		}
	} else {
		var results []dns.ReverseRecordRead
		log.Printf("[INFO] Populating data source for DNS Reverse Record using name: %s and partial_match: %t", recordName, partialMatch)
		if partialMatch {
			// In order to have an exact name match, we must retrieve all the DNS Reverse Records and then
			// build a list of partial matches based on the response
			records, _, err := client.ListReverseRecords(ctx, nil)
			if err != nil {
				return diag.FromErr(fmt.Errorf("an error occurred while fetching DNS Reverse Records: %w", err))
			}
			for _, recordItem := range records.Items {
				if len(results) == 1 {
					break
				}
				if strings.Contains(recordItem.Properties.Name, recordName) {
					results = append(results, recordItem)
				}
			}
			results = records.Items
		} else {
			// In order to have an exact name match, we must retrieve all the DNS Reverse Records and then
			// build a list of exact matches based on the response
			records, _, err := client.ListReverseRecords(ctx, nil)
			if err != nil {
				return diag.FromErr(fmt.Errorf("an error occurred while fetching DNS Reverse Records: %w", err))
			}
			for _, recordItem := range records.Items {
				if strings.EqualFold(recordItem.Properties.Name, recordName) {
					results = append(results, recordItem)
				}
			}
		}
		if results == nil || len(results) == 0 {
			return diag.FromErr(fmt.Errorf("no DNS Reverse Record found with the specified name = %s", recordName))
		} else if len(results) > 1 {
			return diag.FromErr(fmt.Errorf("more than one DNS Reverse Record found with the specified name = %s", recordName))
		} else {
			record = results[0]
		}
	}

	if err := client.SetReverseRecordData(d, record); err != nil {
		return diag.FromErr(err)
	}

	return nil
}
