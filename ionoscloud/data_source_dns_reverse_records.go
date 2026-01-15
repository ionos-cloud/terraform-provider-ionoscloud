package ionoscloud

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	dns "github.com/ionos-cloud/sdk-go-bundle/products/dns/v2"

	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/bundleclient"
)

func dataSourceDNSReverseRecords() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceReverseRecordReads,
		Schema: map[string]*schema.Schema{
			"ips": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "List of IPs to filter by.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"name": {
				Type:        schema.TypeString,
				Description: "The name of your DNS Reverse Record.",
				Optional:    true,
			},
			"partial_match": {
				Type:        schema.TypeBool,
				Description: "Whether partial matching is allowed or not when using name argument.",
				Default:     false,
				Optional:    true,
			},
			"reverse_records": {
				Type:        schema.TypeList,
				Description: "list of recerse records",
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Description: "The unique ID of the server.",
							Computed:    true,
						},
						"name": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"ip": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"description": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
		},
		Timeouts: &resourceDefaultTimeouts,
	}
}

func dataSourceReverseRecordReads(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(bundleclient.SdkBundle).DNSClient
	partialMatch := d.Get("partial_match").(bool)
	nameValue, nameOk := d.GetOk("name")
	ipsValue, ipsOk := d.GetOk("ips")
	recordName := nameValue.(string)
	var filterIps []string
	if ipsOk {
		rawIps := ipsValue.([]interface{})

		for _, item := range rawIps {
			filterIps = append(filterIps, item.(string))
		}
	}

	var err error

	var results []dns.ReverseRecordRead

	records, _, err := client.ListReverseRecords(ctx, filterIps)
	if err != nil {
		return diag.FromErr(fmt.Errorf("an error occurred while fetching DNS Reverse Records: %w", err))
	}
	if nameOk {
		log.Printf("[INFO] Filtering Reverse Records in data source for DNS Reverse Records using name: %s and partial_match: %t", recordName, partialMatch)
		if partialMatch {
			for _, recordItem := range records.Items {
				if strings.Contains(recordItem.Properties.Name, recordName) {
					results = append(results, recordItem)
				}
			}
		} else {
			for _, recordItem := range records.Items {
				if strings.EqualFold(recordItem.Properties.Name, recordName) {
					results = append(results, recordItem)
				}
			}
		}
		records.Items = results
	}

	d.SetId("dns_reverse_records")
	if err := d.Set("reverse_records", reverseRecordsObjToIntf(records.Items)); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

// reverseRecordsObjToIntf converts a list of reverse records api objects to an interface list
func reverseRecordsObjToIntf(reverseRecordsObjects []dns.ReverseRecordRead) []interface{} {
	reverseRecordList := make([]interface{}, len(reverseRecordsObjects))

	for i, record := range reverseRecordsObjects {
		item := make(map[string]interface{})

		item["id"] = record.Id
		item["name"] = record.Properties.Name
		item["ip"] = record.Properties.Ip
		item["description"] = *record.Properties.Description

		reverseRecordList[i] = item
	}

	return reverseRecordList
}
