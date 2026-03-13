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
	diagutil "github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils/diags"
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
				ForceNew: true,
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
		return diagutil.ToDiags(d, fmt.Errorf("an error occurred while creating a DNS Zone: %w", err), &diagutil.ErrorContext{StatusCode: apiResponse.StatusCode})
	}
	if zoneResponse.Metadata.State == dns.PROVISIONINGSTATE_FAILED {
		// This is a temporary error message since right now the API is not returning errors that we can work with.
		return diagutil.ToDiags(d, fmt.Errorf("zone creation has failed, this can happen if the data in the request is not correct, "+
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
		return diagutil.ToDiags(d, fmt.Errorf("error while fetching DNS Zone with ID: %s, error: %w", zoneId, err), &diagutil.ErrorContext{StatusCode: apiResponse.StatusCode})
	}

	log.Printf("[INFO] Successfully retrieved DNS Zone with ID: %s: %+v", zoneId, zone)

	if err := client.SetZoneData(d, zone); err != nil {
		return diagutil.ToDiags(d, err, nil)
	}
	return nil
}

func zoneUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(bundleclient.SdkBundle).DNSClient
	zoneId := d.Id()

	zoneResponse, apiResponse, err := client.UpdateZone(ctx, zoneId, d)
	if err != nil {
		return diagutil.ToDiags(d, fmt.Errorf("an error occurred while updating the DNS Zone with ID: %s, error: %w", zoneId, err), &diagutil.ErrorContext{StatusCode: apiResponse.StatusCode})
	}
	if zoneResponse.Metadata.State == dns.PROVISIONINGSTATE_FAILED {
		// This is a temporary error message since right now the API is not returning errors that we can work with.
		return diagutil.ToDiags(d, fmt.Errorf("zone update has failed, this can happen if the data in the request is not correct, "+
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
		return diagutil.ToDiags(d, fmt.Errorf("error while deleting DNS Zone with ID: %s, error: %w", zoneId, err), &diagutil.ErrorContext{StatusCode: apiResponse.StatusCode})
	}

	err = utils.WaitForResourceToBeDeleted(ctx, d, client.IsZoneDeleted)
	if err != nil {
		return diagutil.ToDiags(d, fmt.Errorf("an error occurred while waiting for the DNS Zone with ID: %s to be deleted, error: %w", zoneId, err), &diagutil.ErrorContext{Timeout: schema.TimeoutDelete})
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
			return nil, diagutil.ToError(d, fmt.Errorf("DNS Zone with ID: %s does not exist", zoneId), &diagutil.ErrorContext{StatusCode: apiResponse.StatusCode})
		}
		return nil, diagutil.ToError(d, fmt.Errorf("an error occurred while trying to import the DNS Zone with ID: %s, error: %w", zoneId, err), &diagutil.ErrorContext{StatusCode: apiResponse.StatusCode})
	}
	log.Printf("[INFO DNS Zone with ID: %s found: %+v", zoneId, zone)

	if err := client.SetZoneData(d, zone); err != nil {
		return nil, diagutil.ToError(d, err, nil)
	}

	return []*schema.ResourceData{d}, nil
}
