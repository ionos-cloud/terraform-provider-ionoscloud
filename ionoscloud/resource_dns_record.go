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
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils"
)

func resourceDNSRecord() *schema.Resource {
	return &schema.Resource{
		CreateContext: recordCreate,
		ReadContext:   recordRead,
		UpdateContext: recordUpdate,
		DeleteContext: recordDelete,
		Importer: &schema.ResourceImporter{
			StateContext: recordImport,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"type": {
				Type:     schema.TypeString,
				Required: true,
			},
			"content": {
				Type:     schema.TypeString,
				Required: true,
			},
			"ttl": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"priority": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"fqdn": {
				Type:        schema.TypeString,
				Description: "Fully qualified domain name",
				Computed:    true,
			},
			"zone_id": {
				Type: schema.TypeString,
				// This should be required, changing this would require adding extra checks in the
				// code where we rely on the fact that this is required.
				Required: true,
			},
		},
		Timeouts: &resourceDefaultTimeouts,
	}
}

func recordCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(bundleclient.SdkBundle).DNSClient
	zoneId := d.Get("zone_id").(string)

	recordResponse, apiResponse, err := client.CreateRecord(ctx, zoneId, d)
	if err != nil {
		return utils.ToDiags(d, fmt.Sprintf("an error occurred while creating a record for the DNS zone with ID: %s, error: %s", zoneId, err), &utils.DiagsOpts{StatusCode: apiResponse.StatusCode})
	}

	if recordResponse.Metadata.State == dns.PROVISIONINGSTATE_FAILED {
		// This is a temporary error message since right now the API is not returning errors that we can work with.
		return utils.ToDiags(d, fmt.Sprintf("record creation has failed, this can happen if the data in the request is not correct, "+
			"please check again the values defined in the plan"), nil)
	}
	d.SetId(recordResponse.Id)
	return recordRead(ctx, d, meta)
}

func recordRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(bundleclient.SdkBundle).DNSClient
	zoneId := d.Get("zone_id").(string)
	recordId := d.Id()

	record, apiResponse, err := client.GetRecordById(ctx, zoneId, recordId)
	if err != nil {
		if apiResponse.HttpNotFound() {
			d.SetId("")
			return nil
		}
		return utils.ToDiags(d, fmt.Sprintf("error while fetching the DNS Record, zone ID: %s, error: %s", zoneId, err), &utils.DiagsOpts{StatusCode: apiResponse.StatusCode})
	}
	log.Printf("[INFO] Successfully retrieved DNS Record %s: %+v", recordId, record)
	if err := client.SetRecordData(d, record); err != nil {
		return utils.ToDiags(d, err.Error(), nil)
	}
	return nil
}

func recordUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(bundleclient.SdkBundle).DNSClient
	zoneId := d.Get("zone_id").(string)
	recordId := d.Id()

	recordResponse, apiResponse, err := client.UpdateRecord(ctx, zoneId, recordId, d)
	if err != nil {
		return utils.ToDiags(d, fmt.Sprintf("an error occurred while updating the DNS Record, zone ID: %s, error: %s", zoneId, err), &utils.DiagsOpts{StatusCode: apiResponse.StatusCode})
	}
	if recordResponse.Metadata.State == dns.PROVISIONINGSTATE_FAILED {
		// This is a temporary error message since right now the API is not returning errors that we can work with.
		return utils.ToDiags(d, fmt.Sprintf("record update has failed, this can happen if the data in the request is not correct, "+
			"please check again the values defined in the plan"), nil)
	}
	return recordRead(ctx, d, meta)
}

func recordDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(bundleclient.SdkBundle).DNSClient
	zoneId := d.Get("zone_id").(string)
	recordId := d.Id()

	apiResponse, err := client.DeleteRecord(ctx, zoneId, recordId)
	if err != nil {
		if apiResponse.HttpNotFound() {
			d.SetId("")
			return nil
		}
		return utils.ToDiags(d, fmt.Sprintf("error while deleting DNS Record, zone ID: %s, error: %s", zoneId, err), &utils.DiagsOpts{StatusCode: apiResponse.StatusCode})
	}
	err = utils.WaitForResourceToBeDeleted(ctx, d, client.IsRecordDeleted)
	if err != nil {
		return utils.ToDiags(d, fmt.Sprintf("an error occurred while waiting for the DNS Record to be deleted, error: %s", err), &utils.DiagsOpts{Timeout: schema.TimeoutDelete})
	}
	return nil
}

func recordImport(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	client := meta.(bundleclient.SdkBundle).DNSClient

	// Split the string provided in order to get the IDs for both zone and record.
	parts := strings.Split(d.Id(), "/")
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return nil, utils.ToError(d, "invalid import, expected {zone UUID}/{record UUID}", nil)
	}
	zoneId := parts[0]
	recordId := parts[1]

	record, apiResponse, err := client.GetRecordById(ctx, zoneId, recordId)
	if err != nil {
		if apiResponse.HttpNotFound() {
			d.SetId("")
			return nil, utils.ToError(d, fmt.Sprintf("DNS Record with ID: %s does not exist, zone ID: %s", recordId, zoneId), &utils.DiagsOpts{StatusCode: apiResponse.StatusCode})
		}
		return nil, utils.ToError(d, fmt.Sprintf("an error occurred while trying to import the DNS Record with ID: %s, zone ID: %s, error: %s", recordId, zoneId, err), &utils.DiagsOpts{StatusCode: apiResponse.StatusCode})
	}
	log.Printf("[INFO] DNS Record found: %+v", record)
	if err := client.SetRecordData(d, record); err != nil {
		return nil, utils.ToError(d, err.Error(), nil)
	}
	if err := d.Set("zone_id", zoneId); err != nil {
		return nil, utils.ToError(d, err.Error(), nil)
	}
	return []*schema.ResourceData{d}, nil
}
