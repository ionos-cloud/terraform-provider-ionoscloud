package ionoscloud

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils"
	"log"
)

func resourceDNSaaSZone() *schema.Resource {
	return &schema.Resource{
		CreateContext: zoneCreate,
		ReadContext:   zoneRead,
		UpdateContext: zoneUpdate,
		DeleteContext: zoneDelete,
		Importer: &schema.ResourceImporter{
			StateContext: zoneImport,
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
			"enabled": {
				Type:     schema.TypeBool,
				Optional: true,
			},
		},
		Timeouts: &resourceDefaultTimeouts,
	}
}

func zoneCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(SdkBundle).DNSaaSClient
	id, _, err := client.CreateZone(ctx, d)

	if err != nil {
		return diag.FromErr(fmt.Errorf("an error occured while creating a DNS Zone: %w", err))
	}
	d.SetId(id)
	err = utils.WaitForResourceToBeReady(ctx, d, client.IsZoneCreated)
	if err != nil {
		return diag.FromErr(fmt.Errorf("an error occured while waiting for DNS Zone with ID: %s to be ready, error: %w", id, err))
	}

	return zoneRead(ctx, d, meta)
}

func zoneRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(SdkBundle).DNSaaSClient
	zoneId := d.Id()
	zone, apiResponse, err := client.GetZoneById(ctx, zoneId)

	if err != nil {
		if apiResponse.HttpNotFound() {
			log.Printf("[INFO] Could not find zone with ID: %s", zoneId)
			d.SetId("")
			return nil
		}
		return diag.FromErr(fmt.Errorf("error while fetching DNS Zone with ID: %s, error: %w", zoneId, err))
	}

	log.Printf("[INFO] Successfully retrieved DNS Zone with ID: %s: %+v", zoneId, zone)

	if err := client.SetZoneData(d, zone); err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func zoneUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(SdkBundle).DNSaaSClient
	zoneId := d.Id()

	_, err := client.UpdateZone(ctx, zoneId, d)
	if err != nil {
		return diag.FromErr(fmt.Errorf("an error occured while updating the DNS Zone with ID: %s, error: %w", err))
	}

	err = utils.WaitForResourceToBeReady(ctx, d, client.IsZoneCreated)
	if err != nil {
		return diag.FromErr(fmt.Errorf("an error occured while waiting for DNS Zone with ID: %s to be ready, error: %w", zoneId, err))
	}
	return nil
}

func zoneDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(SdkBundle).DNSaaSClient
	zoneId := d.Id()

	apiResponse, err := client.DeleteZone(ctx, zoneId)
	if err != nil {
		if apiResponse.HttpNotFound() {
			d.SetId("")
			return nil
		}
		return diag.FromErr(fmt.Errorf("error while deleting DNS Zone with ID: %s, error: %w", zoneId, err))
	}

	err = utils.WaitForResourceToBeDeleted(ctx, d, client.IsZoneDeleted)
	if err != nil {
		diag.FromErr(fmt.Errorf("an error occured while waiting for the DNS Zone with ID: %s to be deleted, error: %w", zoneId, err))
	}
	return nil
}

func zoneImport(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	client := meta.(SdkBundle).DNSaaSClient
	zoneId := d.Id()

	zone, apiResponse, err := client.GetZoneById(ctx, zoneId)
	if err != nil {
		if apiResponse.HttpNotFound() {
			d.SetId("")
			return nil, fmt.Errorf("DNS Zone with ID: %s does not exist", zoneId)
		}
		return nil, fmt.Errorf("an error occured while trying to import the DNS Zone with ID: %s, error: %w", zoneId, err)
	}
	log.Printf("[INFO DNS Zone with ID: %s found: %+v", zoneId, zone)

	if err := client.SetZoneData(d, zone); err != nil {
		return nil, err
	}

	return []*schema.ResourceData{d}, nil
}
