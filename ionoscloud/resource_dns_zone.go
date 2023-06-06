package ionoscloud

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	dns "github.com/ionos-cloud/sdk-go-dns"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils"
)

func resourceDNSZone() *schema.Resource {
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
				Computed: true,
			},
			"nameservers": {
				Type:        schema.TypeList,
				Description: "A list of available name servers.",
				Computed:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
		Timeouts: &resourceDefaultTimeouts,
	}
}

func zoneCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(SdkBundle).DNSClient
	zoneResponse, _, err := client.CreateZone(ctx, d)

	if err != nil {
		return diag.FromErr(fmt.Errorf("an error occured while creating a DNS Zone: %w", err))
	}
	if zoneResponse.Metadata.State != nil {
		if *zoneResponse.Metadata.State == dns.FAILED {
			// This is a temporary error message since right now the API is not returning errors that we can work with.
			return diag.FromErr(fmt.Errorf("zone creation has failed, this can happen if the data in the request is not correct, " +
				"please check again the values defined in the plan"))
		}
	}
	d.SetId(*zoneResponse.Id)
	return zoneRead(ctx, d, meta)
}

func zoneRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(SdkBundle).DNSClient
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
	client := meta.(SdkBundle).DNSClient
	zoneId := d.Id()

	zoneResponse, _, err := client.UpdateZone(ctx, zoneId, d)
	if err != nil {
		return diag.FromErr(fmt.Errorf("an error occured while updating the DNS Zone with ID: %s, error: %w", zoneId, err))
	}
	if zoneResponse.Metadata.State != nil {
		if *zoneResponse.Metadata.State == dns.FAILED {
			// This is a temporary error message since right now the API is not returning errors that we can work with.
			return diag.FromErr(fmt.Errorf("zone update has failed, this can happen if the data in the request is not correct, " +
				"please check again the values defined in the plan"))
		}
	}
	return nil
}

func zoneDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(SdkBundle).DNSClient
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
		return diag.FromErr(fmt.Errorf("an error occured while waiting for the DNS Zone with ID: %s to be deleted, error: %w", zoneId, err))
	}
	return nil
}

func zoneImport(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	client := meta.(SdkBundle).DNSClient
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
