package ionoscloud

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/bundleclient"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils"
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
		return utils.ToDiags(d, fmt.Sprintf("an error occurred while creating reverse record, error: %s", err), &utils.DiagsOpts{StatusCode: apiResponse.StatusCode})
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
		return utils.ToDiags(d, fmt.Sprintf("error while fetching the DNS Reverse Record: %s", err), &utils.DiagsOpts{StatusCode: apiResponse.StatusCode})
	}
	log.Printf("[INFO] Successfully retrieved DNS Reverse Record %s: %+v", recordId, record)
	if err := client.SetReverseRecordData(d, record); err != nil {
		return utils.ToDiags(d, err.Error(), nil)
	}
	return nil
}

func reverseRecordUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(bundleclient.SdkBundle).DNSClient
	recordId := d.Id()

	_, apiResponse, err := client.UpdateReverseRecord(ctx, recordId, d)
	if err != nil {
		return utils.ToDiags(d, fmt.Sprintf("an error occurred while updating the DNS Reverse Record: %s", err), &utils.DiagsOpts{StatusCode: apiResponse.StatusCode})
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
		return utils.ToDiags(d, fmt.Sprintf("error while deleting DNS Reverse Record: %s", err), &utils.DiagsOpts{StatusCode: apiResponse.StatusCode})
	}
	err = utils.WaitForResourceToBeDeleted(ctx, d, client.IsReverseRecordDeleted)
	if err != nil {
		return utils.ToDiags(d, fmt.Sprintf("an error occurred while waiting for the DNS Reverse Record to be deleted: %s", err), &utils.DiagsOpts{Timeout: schema.TimeoutDelete})
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
			return nil, utils.ToError(d, fmt.Sprintf("DNS Reverse Record with ID: %s does not exist", recordId), &utils.DiagsOpts{StatusCode: apiResponse.StatusCode})
		}
		return nil, utils.ToError(d, fmt.Sprintf("an error occurred while trying to import the DNS Reverse Record: %s", err), &utils.DiagsOpts{StatusCode: apiResponse.StatusCode})
	}
	log.Printf("[INFO] DNS Reverse Record found: %+v", record)
	if err := client.SetReverseRecordData(d, record); err != nil {
		return nil, utils.ToError(d, err.Error(), nil)
	}
	return []*schema.ResourceData{d}, nil
}
