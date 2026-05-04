package ionoscloud

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/bundleclient"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils"
	diagutil "github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils/diags"
)

func resourceDNSReverseRecord() *schema.Resource {
	return &schema.Resource{
		CreateContext: reverseRecordCreate,
		ReadContext:   reverseRecordRead,
		UpdateContext: reverseRecordUpdate,
		DeleteContext: reverseRecordDelete,
		Importer: &schema.ResourceImporter{
			StateContext: reverseRecordImport,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"ip": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
		Timeouts: &resourceDefaultTimeouts,
	}
}

func reverseRecordCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(bundleclient.SdkBundle).DNSClient

	recordResponse, apiResponse, err := client.CreateReverseRecord(ctx, d)
	if err != nil {
		return diagutil.ToDiags(d, fmt.Errorf("an error occurred while creating reverse record, error: %w", err), &diagutil.ErrorContext{StatusCode: apiResponse.SafeStatusCode()})
	}

	d.SetId(recordResponse.Id)
	return reverseRecordRead(ctx, d, meta)
}

func reverseRecordRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(bundleclient.SdkBundle).DNSClient
	recordId := d.Id()

	record, apiResponse, err := client.GetReverseRecordById(ctx, recordId)
	if err != nil {
		if apiResponse.HttpNotFound() {
			d.SetId("")
			return nil
		}
		return diagutil.ToDiags(d, fmt.Errorf("error while fetching the DNS Reverse Record: %w", err), &diagutil.ErrorContext{StatusCode: apiResponse.SafeStatusCode()})
	}
	tflog.Info(ctx, "retrieved DNS reverse record", map[string]interface{}{"record_id": recordId})
	if err := client.SetReverseRecordData(d, record); err != nil {
		return diagutil.ToDiags(d, err, nil)
	}
	return nil
}

func reverseRecordUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(bundleclient.SdkBundle).DNSClient
	recordId := d.Id()

	_, apiResponse, err := client.UpdateReverseRecord(ctx, recordId, d)
	if err != nil {
		return diagutil.ToDiags(d, fmt.Errorf("an error occurred while updating the DNS Reverse Record: %w", err), &diagutil.ErrorContext{StatusCode: apiResponse.SafeStatusCode()})
	}
	return reverseRecordRead(ctx, d, meta)
}

func reverseRecordDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(bundleclient.SdkBundle).DNSClient
	recordId := d.Id()

	apiResponse, err := client.DeleteReverseRecord(ctx, recordId)
	if err != nil {
		if apiResponse.HttpNotFound() {
			d.SetId("")
			return nil
		}
		return diagutil.ToDiags(d, fmt.Errorf("error while deleting DNS Reverse Record: %w", err), &diagutil.ErrorContext{StatusCode: apiResponse.SafeStatusCode()})
	}
	err = utils.WaitForResourceToBeDeleted(ctx, d, client.IsReverseRecordDeleted)
	if err != nil {
		return diagutil.ToDiags(d, fmt.Errorf("an error occurred while waiting for the DNS Reverse Record to be deleted: %w", err), &diagutil.ErrorContext{Timeout: d.Timeout(schema.TimeoutDelete).String()})
	}
	return nil
}

func reverseRecordImport(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	client := meta.(bundleclient.SdkBundle).DNSClient

	recordId := d.Id()

	record, apiResponse, err := client.GetReverseRecordById(ctx, recordId)
	if err != nil {
		if apiResponse.HttpNotFound() {
			d.SetId("")
			return nil, diagutil.ToError(d, fmt.Errorf("DNS Reverse Record with ID: %s does not exist", recordId), &diagutil.ErrorContext{StatusCode: apiResponse.SafeStatusCode()})
		}
		return nil, diagutil.ToError(d, fmt.Errorf("an error occurred while trying to import the DNS Reverse Record: %w", err), &diagutil.ErrorContext{StatusCode: apiResponse.SafeStatusCode()})
	}
	tflog.Info(ctx, "DNS reverse record imported", map[string]interface{}{"record_id": recordId})
	if err := client.SetReverseRecordData(d, record); err != nil {
		return nil, diagutil.ToError(d, err, nil)
	}
	return []*schema.ResourceData{d}, nil
}
