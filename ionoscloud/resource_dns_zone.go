package ionoscloud

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-log/tflog"
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
		Identity: &schema.ResourceIdentity{
			Version: 0,
			SchemaFunc: func() map[string]*schema.Schema {
				return map[string]*schema.Schema{
					"id": {
						Type:              schema.TypeString,
						RequiredForImport: true,
						Description:       "The UUID of the DNS zone.",
					},
				}
			},
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

func zoneCreate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(bundleclient.SdkBundle).DNSClient
	zoneResponse, apiResponse, err := client.CreateZone(ctx, d)

	if err != nil {
		return diagutil.ToDiags(d, fmt.Errorf("an error occurred while creating a DNS Zone: %w", err), &diagutil.ErrorContext{StatusCode: apiResponse.SafeStatusCode()})
	}
	if zoneResponse.Metadata.State == dns.PROVISIONINGSTATE_FAILED {
		// This is a temporary error message since right now the API is not returning errors that we can work with.
		return diagutil.ToDiags(d, fmt.Errorf("zone creation has failed, this can happen if the data in the request is not correct, "+
			"please check again the values defined in the plan"), nil)
	}
	d.SetId(zoneResponse.Id)
	return zoneRead(ctx, d, meta)
}

func zoneRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(bundleclient.SdkBundle).DNSClient
	zoneID := d.Id()
	zone, apiResponse, err := client.GetZoneById(ctx, zoneID)

	if err != nil {
		if apiResponse.HttpNotFound() {
			tflog.Info(ctx, "DNS zone not found", map[string]any{"zone_id": zoneID})
			d.SetId("")
			return nil
		}
		return diagutil.ToDiags(d, fmt.Errorf("error while fetching DNS Zone with ID: %s, error: %w", zoneID, err), &diagutil.ErrorContext{StatusCode: apiResponse.SafeStatusCode()})
	}

	tflog.Info(ctx, "retrieved DNS zone", map[string]any{"zone_id": zoneID})

	if err := client.SetZoneData(d, zone); err != nil {
		return diagutil.ToDiags(d, err, nil)
	}

	// Must run after SetZoneData: the identity reads attributes that only the data
	// setter fills in (this matters most on an identity-based import).
	if err := setDNSZoneIdentity(d); err != nil {
		return diagutil.ToDiags(d, err, nil)
	}
	return nil
}

func zoneUpdate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(bundleclient.SdkBundle).DNSClient
	zoneID := d.Id()

	zoneResponse, apiResponse, err := client.UpdateZone(ctx, zoneID, d)
	if err != nil {
		return diagutil.ToDiags(d, fmt.Errorf("an error occurred while updating the DNS Zone with ID: %s, error: %w", zoneID, err), &diagutil.ErrorContext{StatusCode: apiResponse.SafeStatusCode()})
	}
	if zoneResponse.Metadata.State == dns.PROVISIONINGSTATE_FAILED {
		// This is a temporary error message since right now the API is not returning errors that we can work with.
		return diagutil.ToDiags(d, fmt.Errorf("zone update has failed, this can happen if the data in the request is not correct, "+
			"please check again the values defined in the plan"), nil)
	}

	// Unlike Create, Update does not need this to satisfy terraform: the SDK carries the
	// prior identity into the apply on its own, so an update that never touches the
	// identity still returns one. It is a safety net for state written before this
	// resource declared an identity, e.g. a refresh-free apply over an old state file.
	// The value is read back out of state, never out of an API response, so whatever this
	// writes equals the prior identity and the stability check cannot trip.
	if err := setDNSZoneIdentity(d); err != nil {
		return diagutil.ToDiags(d, err, nil)
	}
	return nil
}

func zoneDelete(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(bundleclient.SdkBundle).DNSClient
	zoneID := d.Id()

	apiResponse, err := client.DeleteZone(ctx, zoneID)
	if err != nil {
		if apiResponse.HttpNotFound() {
			d.SetId("")
			return nil
		}
		return diagutil.ToDiags(d, fmt.Errorf("error while deleting DNS Zone with ID: %s, error: %w", zoneID, err), &diagutil.ErrorContext{StatusCode: apiResponse.SafeStatusCode()})
	}

	err = utils.WaitForResourceToBeDeleted(ctx, d, client.IsZoneDeleted)
	if err != nil {
		return diagutil.ToDiags(d, fmt.Errorf("an error occurred while waiting for the DNS Zone with ID: %s to be deleted, error: %w", zoneID, err), &diagutil.ErrorContext{Timeout: d.Timeout(schema.TimeoutDelete).String()})
	}
	return nil
}

func zoneImport(ctx context.Context, d *schema.ResourceData, meta any) ([]*schema.ResourceData, error) {
	client := meta.(bundleclient.SdkBundle).DNSClient

	zoneID, err := dnsZoneImportParts(d)
	if err != nil {
		return nil, err
	}

	// Terraform sends an empty ID for an identity-based import, so the ID has to be set
	// here for the error diagnostics and the log line below to name the resource.
	d.SetId(zoneID)

	zone, apiResponse, err := client.GetZoneById(ctx, zoneID)
	if err != nil {
		if apiResponse.HttpNotFound() {
			d.SetId("")
			return nil, diagutil.ToError(d, fmt.Errorf("DNS Zone with ID: %s does not exist", zoneID), &diagutil.ErrorContext{StatusCode: apiResponse.SafeStatusCode()})
		}
		return nil, diagutil.ToError(d, fmt.Errorf("an error occurred while trying to import the DNS Zone with ID: %s, error: %w", zoneID, err), &diagutil.ErrorContext{StatusCode: apiResponse.SafeStatusCode()})
	}
	tflog.Info(ctx, "DNS zone imported", map[string]any{"zone_id": zoneID})

	if err := client.SetZoneData(d, zone); err != nil {
		return nil, diagutil.ToError(d, err, nil)
	}

	if err := setDNSZoneIdentity(d); err != nil {
		return nil, diagutil.ToError(d, err, nil)
	}

	return []*schema.ResourceData{d}, nil
}

// dnsZoneImportParts resolves the zone ID from either import mode: the `id` attribute of
// an identity-based import block, or the legacy plain-UUID import string. A DNS zone has
// no location, so the identity is a lone `id` and there is nothing to split.
func dnsZoneImportParts(d *schema.ResourceData) (string, error) {
	if identity, identityErr := d.Identity(); identityErr == nil {
		if id, ok := identity.GetOk("id"); ok {
			zoneID, _ := id.(string)
			return zoneID, nil
		}
	}

	zoneID := d.Id()
	if zoneID == "" {
		return "", fmt.Errorf("invalid import identifier: expected a DNS zone UUID, got an empty string")
	}

	return zoneID, nil
}

// setDNSZoneIdentity writes the resource identity. It reads the id back out of the
// ResourceData, so it must run after the state writer on every path that produces state.
func setDNSZoneIdentity(d *schema.ResourceData) error {
	identity, err := d.Identity()
	if err != nil {
		return err
	}

	if err := identity.Set("id", d.Id()); err != nil {
		return fmt.Errorf("error while setting id identity attribute for DNS zone %s: %w", d.Id(), err)
	}

	return nil
}
