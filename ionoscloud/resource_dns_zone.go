package ionoscloud

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	dns "github.com/ionos-cloud/sdk-go-bundle/products/dns/v2"

	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/bundleclient"
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
	client := meta.(bundleclient.SdkBundle).DNSClient
	zoneResponse, apiResponse, err := client.CreateZone(ctx, d)

	if err != nil {
		return utils.ToDiags(d, fmt.Sprintf("an error occurred while creating a DNS Zone: %s", err), &utils.DiagsOpts{StatusCode: apiResponse.StatusCode})
	}
	if zoneResponse.Metadata.State == dns.PROVISIONINGSTATE_FAILED {
		// This is a temporary error message since right now the API is not returning errors that we can work with.
		return utils.ToDiags(d, fmt.Sprintf("zone creation has failed, this can happen if the data in the request is not correct, "+
			"please check again the values defined in the plan"), nil)
	}
	d.SetId(zoneResponse.Id)
	return zoneRead(ctx, d, meta)
}

func zoneRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(bundleclient.SdkBundle).DNSClient
	zoneId := d.Id()
	zone, apiResponse, err := client.GetZoneById(ctx, zoneId)

	if err != nil {
		if apiResponse.HttpNotFound() {
			log.Printf("[INFO] Could not find zone with ID: %s", zoneId)
			d.SetId("")
			return nil
		}
		return utils.ToDiags(d, fmt.Sprintf("error while fetching DNS Zone with ID: %s, error: %s", zoneId, err), &utils.DiagsOpts{StatusCode: apiResponse.StatusCode})
	}

	log.Printf("[INFO] Successfully retrieved DNS Zone with ID: %s: %+v", zoneId, zone)

	if err := client.SetZoneData(d, zone); err != nil {
		return utils.ToDiags(d, err.Error(), nil)
	}
	return nil
}

func zoneUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(bundleclient.SdkBundle).DNSClient
	zoneId := d.Id()

	zoneResponse, apiResponse, err := client.UpdateZone(ctx, zoneId, d)
	if err != nil {
		return utils.ToDiags(d, fmt.Sprintf("an error occurred while updating the DNS Zone with ID: %s, error: %s", zoneId, err), &utils.DiagsOpts{StatusCode: apiResponse.StatusCode})
	}
	if zoneResponse.Metadata.State == dns.PROVISIONINGSTATE_FAILED {
		// This is a temporary error message since right now the API is not returning errors that we can work with.
		return utils.ToDiags(d, fmt.Sprintf("zone update has failed, this can happen if the data in the request is not correct, "+
			"please check again the values defined in the plan"), nil)
	}
	return nil
}

func zoneDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(bundleclient.SdkBundle).DNSClient
	zoneId := d.Id()

	apiResponse, err := client.DeleteZone(ctx, zoneId)
	if err != nil {
		if apiResponse.HttpNotFound() {
			d.SetId("")
			return nil
		}
		return utils.ToDiags(d, fmt.Sprintf("error while deleting DNS Zone with ID: %s, error: %s", zoneId, err), &utils.DiagsOpts{StatusCode: apiResponse.StatusCode})
	}

	err = utils.WaitForResourceToBeDeleted(ctx, d, client.IsZoneDeleted)
	if err != nil {
		return utils.ToDiags(d, fmt.Sprintf("an error occurred while waiting for the DNS Zone with ID: %s to be deleted, error: %s", zoneId, err), &utils.DiagsOpts{Timeout: schema.TimeoutDelete})
	}
	return nil
}

func zoneImport(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	client := meta.(bundleclient.SdkBundle).DNSClient
	zoneId := d.Id()

	zone, apiResponse, err := client.GetZoneById(ctx, zoneId)
	if err != nil {
		if apiResponse.HttpNotFound() {
			d.SetId("")
			return nil, utils.ToError(d, fmt.Sprintf("DNS Zone with ID: %s does not exist", zoneId), &utils.DiagsOpts{StatusCode: apiResponse.StatusCode})
		}
		return nil, utils.ToError(d, fmt.Sprintf("an error occurred while trying to import the DNS Zone with ID: %s, error: %s", zoneId, err), &utils.DiagsOpts{StatusCode: apiResponse.StatusCode})
	}
	log.Printf("[INFO DNS Zone with ID: %s found: %+v", zoneId, zone)

	if err := client.SetZoneData(d, zone); err != nil {
		return nil, utils.ToError(d, err.Error(), nil)
	}

	return []*schema.ResourceData{d}, nil
}
