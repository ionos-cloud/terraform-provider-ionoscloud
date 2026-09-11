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
// It lives in this package, next to the resource it lists, so it can call resourceDNSZone
// and setDNSZoneIdentity directly. That is the whole design: results are produced by the
// resource's own state writer, so this file declares no model of the DNS zone schema and
// there is nothing here to keep in sync when that schema changes. The writer itself is
// the one piece that is not in this package - SetZoneData is a method on the DNS service
// client, which the mapper reaches through the bundle it already holds.

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

// RawV6Schemas hands the framework the protocol schemas of the SDKv2 managed
// resource. A framework-native list resource inherits them from the resource itself;
// this is only needed because ionoscloud_dns_zone lives on the SDKv2 side.
func (r *dnsZoneListResource) RawV6Schemas(ctx context.Context, _ list.RawV6SchemaRequest, resp *list.RawV6SchemaResponse) {
	fwidentity.SetRawV6Schemas(ctx, resp, constant.DNSZoneResource, r.resourceSchema)
}

// Metadata returns the type name of the managed resource being listed. It must match
// the SDKv2 resource exactly, otherwise terraform has no resource to attach the
// results to - hence the same constant that keys ResourcesMap in provider.go.
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
//
// `enabled` and `nameservers` are deliberately not filterable: MatchesFilters compares
// strings, so a bool and a list of strings have nothing it could match against.
func (r *dnsZoneListResource) ListResourceConfigSchema(_ context.Context, _ list.ListResourceSchemaRequest, resp *list.ListResourceSchemaResponse) {
	resp.Schema = listschema.Schema{
		Attributes: map[string]listschema.Attribute{
			fwidentity.FiltersKey: fwidentity.FilterAttribute("name", "description"),
		},
	}
}

// List fetches every DNS zone on the contract and streams the results. DNS is not a
// regional product - its client is built once and held on the bundle, and GET /zones
// returns the whole collection in one response - so unlike the regional database
// products there is nothing to fan out over here.
func (r *dnsZoneListResource) List(ctx context.Context, req list.ListRequest, stream *list.ListResultsStream) {
	fwidentity.StreamList(ctx, stream, req,
		func(ctx context.Context) ([]dns.ZoneRead, error) {
			// The same call ionoscloud/data_source_dns_zone.go makes, through the same
			// service helper, so the list resource and the data source cannot drift
			// apart in what they can see.
			//
			// The empty filter argument is deliberate. ListZones can push a zone name
			// down as filter.zoneName, but that would save no round trip here - this is
			// a single global request either way - and the mapper has to re-check every
			// filter regardless, so pushing it down would only add a second place for
			// the semantics to disagree.
			//
			// No depth and no limit are sent: ApiZonesGetRequest has no Depth parameter
			// at all (that is a Cloud API concept), and it sends no limit unless one is
			// asked for, which matches the data source. See the pagination note in
			// docs/list-resources/dns_zone.md for what that costs.
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
// ResourceData built from the live schema. It is a method on the DNS service client
// rather than a package-level function in this package, which costs nothing because the
// mapper already holds the bundle. Nothing here knows what attributes a DNS zone has.
func (r *dnsZoneListResource) mapDNSZone(_ context.Context, includeResource bool, filters []fwidentity.Filter, zone dns.ZoneRead) (*fwidentity.MappedItem, diag.Diagnostics) {
	var diags diag.Diagnostics

	// ZoneRead.Id and .Properties are values, not pointers, the way the sdk-go-bundle
	// generator emits them, so the only unusable item is one returned without an id.
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
	// SetZoneData takes its zone by value, and calls d.SetId itself.
	if err := r.bundle.DNSClient.SetZoneData(data, zone); err != nil {
		diags.AddError("Failed to map the dns zone", err.Error())
		return nil, diags
	}

	// setDNSZoneIdentity reads the id back out of the ResourceData, so it has to run
	// after SetZoneData.
	if err := setDNSZoneIdentity(data); err != nil {
		diags.AddError("Failed to map the dns zone identity", err.Error())
		return nil, diags
	}

	// `name` is Required on ionoscloud_dns_zone and ZoneName is a plain string rather
	// than a pointer, so there is always a display name - no fallback to the UUID of the
	// kind ionoscloud_ipblock needs.
	mapped, err := fwidentity.MappedItemFromResourceData(zone.Properties.ZoneName, data, includeResource)
	if err != nil {
		diags.AddError("Failed to convert the dns zone state", err.Error())
		return nil, diags
	}

	return mapped, diags
}
