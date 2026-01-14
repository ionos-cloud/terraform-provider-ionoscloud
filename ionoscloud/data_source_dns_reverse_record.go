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
			"ip": {
				Type:        schema.TypeString,
				Description: "The IP of your DNS Reverse Record.",
				Optional:    true,
				Computed:    true,
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
			"description": {
				Type:     schema.TypeString,
				Computed: true,
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
	ipValue, ipOk := d.GetOk("ip")
	recordId := idValue.(string)
	recordName := nameValue.(string)
	recordIp := ipValue.(string)

	count := 0
	if idOk {
		count++
	}
	if nameOk {
		count++
	}
	if ipOk {
		count++
	}

	if count > 1 {
		return diag.FromErr(fmt.Errorf("only one of [Id, name, ip] can be specified at the same time"))
	}

	if count == 0 {
		return diag.FromErr(fmt.Errorf("please provide either the DNS Record Id, name or IP"))
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
		if nameOk {
			log.Printf("[INFO] Populating data source for DNS Reverse Record using name: %s and partial_match: %t", recordName, partialMatch)
			if partialMatch {
				// In order to have an exact name match, we must retrieve all the DNS Reverse Records and then
				// build a list of partial matches based on the response
				records, _, err := client.ListReverseRecords(ctx, nil)
				if err != nil {
					return diag.FromErr(fmt.Errorf("an error occurred while fetching DNS Reverse Records: %w", err))
				}
				for _, recordItem := range records.Items {
					if strings.Contains(recordItem.Properties.Name, recordName) {
						results = append(results, recordItem)
					}
				}
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
		} else {
			records, _, err := client.ListReverseRecords(ctx, []string{recordIp})
			if err != nil {
				return diag.FromErr(fmt.Errorf("an error occurred while fetching DNS Reverse Records: %w", err))
			}
			results = records.Items

		}
		var usedFilter string
		if ipOk {
			usedFilter = recordId
		} else if nameOk {
			usedFilter = recordName
		}

		if len(results) == 0 {
			return diag.FromErr(fmt.Errorf("no DNS Reverse Record found with the specified filter = %s", usedFilter))
		} else if len(results) > 1 {
			return diag.FromErr(fmt.Errorf("more than one DNS Reverse Record found with the specified name = %s", usedFilter))
		} else {
			record = results[0]
		}
	}

	if err := client.SetReverseRecordData(d, record); err != nil {
		return diag.FromErr(err)
	}

	return nil
}
