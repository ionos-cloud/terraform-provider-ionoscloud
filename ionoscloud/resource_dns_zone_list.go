package ionoscloud

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/list"
	listschema "github.com/hashicorp/terraform-plugin-framework/list/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	dns "github.com/ionos-cloud/sdk-go-bundle/products/dns/v2"
	"github.com/ionos-cloud/sdk-go-bundle/shared"

	fwidentity "github.com/ionos-cloud/terraform-provider-ionoscloud/v6/internal/framework/identity"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/bundleclient"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils/constant"
)

// A list resource for ionoscloud_dns_zone, whose managed resource is still implemented
// with terraform-plugin-sdk/v2 and therefore lives on the other half of the mux. The
// protocol schemas the framework needs come from that resource via
// identity.SetRawV6Schemas.
//
// It lives in this package, next to the resource it lists, so it can call
// resourceDNSZone and setDNSZoneIdentity directly, and hand results to the same
// SetZoneData the DNS service client gives zoneRead. That is the whole design: results
// are produced by the resource's own state writer, so this file declares no model of
// the zone schema and there is nothing here to keep in sync when that schema changes.

var (
	_ list.ListResource                 = (*dnsZoneListResource)(nil)
	_ list.ListResourceWithConfigure    = (*dnsZoneListResource)(nil)
	_ list.ListResourceWithRawV6Schemas = (*dnsZoneListResource)(nil)
)

// dnsZoneListResource lists ionoscloud_dns_zone instances.
type dnsZoneListResource struct {
	bundle *bundleclient.SdkBundle

	// resourceSchema is the SDKv2 managed resource being listed: the source of the
	// protocol schemas returned by RawV6Schemas, and of the ResourceData every
	// result is written into.
	resourceSchema *schema.Resource
}

// NewDNSZoneListResource creates a new list resource for ionoscloud_dns_zone.
func NewDNSZoneListResource() list.ListResource {
	return &dnsZoneListResource{resourceSchema: resourceDNSZone()}
}

// RawV6Schemas hands the framework the protocol schemas of the SDKv2 managed resource.
// A framework-native list resource inherits them from the resource itself; this is only
// needed because ionoscloud_dns_zone lives on the SDKv2 side.
func (r *dnsZoneListResource) RawV6Schemas(ctx context.Context, _ list.RawV6SchemaRequest, resp *list.RawV6SchemaResponse) {
	fwidentity.SetRawV6Schemas(ctx, resp, constant.DNSZoneResource, r.resourceSchema)
}

// Metadata returns the type name of the managed resource being listed. It must match
// the SDKv2 resource exactly, otherwise terraform has no resource to attach the
// results to.
func (r *dnsZoneListResource) Metadata(_ context.Context, _ resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = constant.DNSZoneResource
}

// Configure stores the client bundle shared by the provider.
func (r *dnsZoneListResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	clientBundle, ok := req.ProviderData.(*bundleclient.SdkBundle)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected List Resource Configure Type",
			fmt.Sprintf("Expected *bundleclient.SdkBundle, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	r.bundle = clientBundle
}

// ListResourceConfigSchema returns the schema for the list resource config block.
func (r *dnsZoneListResource) ListResourceConfigSchema(_ context.Context, _ list.ListResourceSchemaRequest, resp *list.ListResourceSchemaResponse) {
	resp.Schema = listschema.Schema{
		Attributes: map[string]listschema.Attribute{
			fwidentity.FiltersKey: fwidentity.FilterAttribute("name", "description"),
		},
	}
}

// List fetches every DNS zone on the contract and streams the results. Cloud DNS is not
// a regional product - the zones of a contract come from a single endpoint - so there is
// nothing to fan out over here.
func (r *dnsZoneListResource) List(ctx context.Context, req list.ListRequest, stream *list.ListResultsStream) {
	fwidentity.StreamList(ctx, stream, req,
		func(ctx context.Context) ([]dns.ZoneRead, error) {
			// The same helper the ionoscloud_dns_zone data source lists with, called with
			// no server-side name filter: filtering stays client-side, in the mapper, so
			// that every filter field behaves the same way. The collection takes no depth
			// and the provider sends no limit, which leaves the page size to the DNS API.
			zones, apiResponse, err := r.bundle.DNSClient.ListZones(ctx, "")
			if apiResponse != nil {
				tflog.Debug(ctx, "listed dns zones", map[string]any{"status_code": apiResponse.SafeStatusCode()})
			}
			if err != nil {
				return nil, fmt.Errorf("failed to list dns zones: %w", err)
			}

			return zones.Items, nil
		},
		r.mapDNSZone,
	)
}

// mapDNSZone maps a DNS zone to an identity.MappedItem, or returns nil to skip it.
//
// The mapping itself is SetZoneData, the same state writer zoneRead uses, run against a
// ResourceData built from the live schema. Nothing here knows what attributes a zone has.
func (r *dnsZoneListResource) mapDNSZone(_ context.Context, includeResource bool, filters []fwidentity.Filter, zone dns.ZoneRead) (*fwidentity.MappedItem, diag.Diagnostics) {
	var diags diag.Diagnostics

	if zone.Id == "" {
		return nil, nil
	}

	if !fwidentity.MatchesFilters(map[string]string{
		"name":        zone.Properties.ZoneName,
		"description": shared.ToValueDefault(zone.Properties.Description),
	}, filters) {
		return nil, nil
	}

	data := r.resourceSchema.Data(&terraform.InstanceState{})
	if err := r.bundle.DNSClient.SetZoneData(data, zone); err != nil {
		diags.AddError("Failed to map the dns zone", err.Error())
		return nil, diags
	}

	// setDNSZoneIdentity reads the id back out of the ResourceData, so it has to run
	// after SetZoneData - which is also what sets that id in the first place.
	if err := setDNSZoneIdentity(data); err != nil {
		diags.AddError("Failed to map the dns zone identity", err.Error())
		return nil, diags
	}

	mapped, err := fwidentity.MappedItemFromResourceData(zone.Properties.ZoneName, data, includeResource)
	if err != nil {
		diags.AddError("Failed to convert the dns zone state", err.Error())
		return nil, diags
	}

	return mapped, diags
}
