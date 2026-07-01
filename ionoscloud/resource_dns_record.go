package ionoscloud

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	dns "github.com/ionos-cloud/sdk-go-bundle/products/dns/v2"

	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/bundleclient"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils"
	diagutil "github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils/diags"
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
				ForceNew: true,
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

func recordCreate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(bundleclient.SdkBundle).DNSClient
	zoneID := d.Get("zone_id").(string)

	recordResponse, apiResponse, err := client.CreateRecord(ctx, zoneID, d)
	if err != nil {
		return bundleclient.ToDiags(meta, d, fmt.Errorf("an error occurred while creating a record for the DNS zone with ID: %s, error: %w", zoneID, err), &diagutil.ErrorContext{StatusCode: apiResponse.SafeStatusCode()})
	}

	if recordResponse.Metadata.State == dns.PROVISIONINGSTATE_FAILED {
		// This is a temporary error message since right now the API is not returning errors that we can work with.
		return bundleclient.ToDiags(meta, d, fmt.Errorf("record creation has failed, this can happen if the data in the request is not correct, "+
			"please check again the values defined in the plan"), nil)
	}
	d.SetId(recordResponse.Id)
	return recordRead(ctx, d, meta)
}

func recordRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(bundleclient.SdkBundle).DNSClient
	zoneID := d.Get("zone_id").(string)
	recordID := d.Id()

	record, apiResponse, err := client.GetRecordById(ctx, zoneID, recordID)
	if err != nil {
		if apiResponse.HttpNotFound() {
			d.SetId("")
			return nil
		}
		return bundleclient.ToDiags(meta, d, fmt.Errorf("error while fetching the DNS Record, zone ID: %s, error: %w", zoneID, err), &diagutil.ErrorContext{StatusCode: apiResponse.SafeStatusCode()})
	}
	tflog.Info(ctx, "retrieved DNS record", map[string]any{"record_id": recordID, "zone_id": zoneID})
	if err := client.SetRecordData(d, record); err != nil {
		return bundleclient.ToDiags(meta, d, err, nil)
	}
	return nil
}

func recordUpdate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(bundleclient.SdkBundle).DNSClient
	zoneID := d.Get("zone_id").(string)
	recordID := d.Id()

	recordResponse, apiResponse, err := client.UpdateRecord(ctx, zoneID, recordID, d)
	if err != nil {
		return bundleclient.ToDiags(meta, d, fmt.Errorf("an error occurred while updating the DNS Record, zone ID: %s, error: %w", zoneID, err), &diagutil.ErrorContext{StatusCode: apiResponse.SafeStatusCode()})
	}
	if recordResponse.Metadata.State == dns.PROVISIONINGSTATE_FAILED {
		// This is a temporary error message since right now the API is not returning errors that we can work with.
		return bundleclient.ToDiags(meta, d, fmt.Errorf("record update has failed, this can happen if the data in the request is not correct, "+
			"please check again the values defined in the plan"), nil)
	}
	return recordRead(ctx, d, meta)
}

func recordDelete(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(bundleclient.SdkBundle).DNSClient
	zoneID := d.Get("zone_id").(string)
	recordID := d.Id()

	apiResponse, err := client.DeleteRecord(ctx, zoneID, recordID)
	if err != nil {
		if apiResponse.HttpNotFound() {
			d.SetId("")
			return nil
		}
		return bundleclient.ToDiags(meta, d, fmt.Errorf("error while deleting DNS Record, zone ID: %s, error: %w", zoneID, err), &diagutil.ErrorContext{StatusCode: apiResponse.SafeStatusCode()})
	}
	err = utils.WaitForResourceToBeDeleted(ctx, d, client.IsRecordDeleted)
	if err != nil {
		return bundleclient.ToDiags(meta, d, fmt.Errorf("an error occurred while waiting for the DNS Record to be deleted, error: %w", err), &diagutil.ErrorContext{Timeout: d.Timeout(schema.TimeoutDelete).String()})
	}
	return nil
}

func recordImport(ctx context.Context, d *schema.ResourceData, meta any) ([]*schema.ResourceData, error) {
	client := meta.(bundleclient.SdkBundle).DNSClient

	// Split the string provided in order to get the IDs for both zone and record.
	parts := strings.Split(d.Id(), "/")
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return nil, bundleclient.ToError(meta, d, fmt.Errorf("invalid import, expected {zone UUID}/{record UUID}"), nil)
	}
	zoneID := parts[0]
	recordID := parts[1]

	record, apiResponse, err := client.GetRecordById(ctx, zoneID, recordID)
	if err != nil {
		if apiResponse.HttpNotFound() {
			d.SetId("")
			return nil, bundleclient.ToError(meta, d, fmt.Errorf("DNS Record with ID: %s does not exist, zone ID: %s", recordID, zoneID), &diagutil.ErrorContext{StatusCode: apiResponse.SafeStatusCode()})
		}
		return nil, bundleclient.ToError(meta, d, fmt.Errorf("an error occurred while trying to import the DNS Record with ID: %s, zone ID: %s, error: %w", recordID, zoneID, err), &diagutil.ErrorContext{StatusCode: apiResponse.SafeStatusCode()})
	}
	tflog.Info(ctx, "DNS record imported", map[string]any{"record_id": recordID, "zone_id": zoneID})
	if err := client.SetRecordData(d, record); err != nil {
		return nil, bundleclient.ToError(meta, d, err, nil)
	}
	if err := d.Set("zone_id", zoneID); err != nil {
		return nil, bundleclient.ToError(meta, d, err, nil)
	}
	return []*schema.ResourceData{d}, nil
}
