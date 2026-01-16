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

	recordResponse, _, err := client.CreateReverseRecord(ctx, d)
	if err != nil {
		return diag.FromErr(fmt.Errorf("an error occurred while creating reverse record, error: %w", err))
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
		return diag.FromErr(fmt.Errorf("error while fetching the DNS Reverse Record with ID: %s, error: %w", recordId, err))
	}
	log.Printf("[INFO] Successfully retrieved DNS Reverse Record %s: %+v", recordId, record)
	if err := client.SetReverseRecordData(d, record); err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func reverseRecordUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(bundleclient.SdkBundle).DNSClient
	recordId := d.Id()

	_, _, err := client.UpdateReverseRecord(ctx, recordId, d)
	if err != nil {
		return diag.FromErr(fmt.Errorf("an error occurred while updating the DNS Reverse Record with ID: %s, error: %w", recordId, err))
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
		return diag.FromErr(fmt.Errorf("error while deleting DNS Reverse Record with ID: %s, error: %w", recordId, err))
	}
	err = utils.WaitForResourceToBeDeleted(ctx, d, client.IsReverseRecordDeleted)
	if err != nil {
		return diag.FromErr(fmt.Errorf("an error occurred while waiting for the DNS Reverse Record with ID: %s to be deleted, error: %w", recordId, err))
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
			return nil, fmt.Errorf("DNS Reverse Record with ID: %s does not exist", recordId)
		}
		return nil, fmt.Errorf("an error occurred while trying to import the DNS Reverse Record with ID: %s, error: %w", recordId, err)
	}
	log.Printf("[INFO] DNS Reverse Record found: %+v", record)
	if err := client.SetReverseRecordData(d, record); err != nil {
		return nil, err
	}
	return []*schema.ResourceData{d}, nil
}
