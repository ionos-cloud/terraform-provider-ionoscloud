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
	"github.com/ionos-cloud/sdk-go-bundle/shared"
	ionoscloud "github.com/ionos-cloud/sdk-go/v6"

	fwidentity "github.com/ionos-cloud/terraform-provider-ionoscloud/v6/internal/framework/identity"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/bundleclient"
)

// A list resource for ionoscloud_datacenter, whose managed resource is still
// implemented with terraform-plugin-sdk/v2 and therefore lives on the other half of
// the mux. The protocol schemas the framework needs come from that resource via
// identity.SetRawV6Schemas.
//
// It lives in this package, next to the resource it lists, so it can call
// resourceDatacenter, setDatacenterData and setDatacenterIdentity directly. That is
// the whole design: results are produced by the resource's own state writer, so this
// file declares no model of the datacenter schema and there is nothing here to keep
// in sync when that schema changes.

const datacenterResourceType = "ionoscloud_datacenter"

var (
	_ list.ListResource                 = (*datacenterListResource)(nil)
	_ list.ListResourceWithConfigure    = (*datacenterListResource)(nil)
	_ list.ListResourceWithRawV6Schemas = (*datacenterListResource)(nil)
)

// datacenterListResource lists ionoscloud_datacenter instances.
type datacenterListResource struct {
	bundle *bundleclient.SdkBundle

	// resourceSchema is the SDKv2 managed resource being listed: the source of the
	// protocol schemas returned by RawV6Schemas, and of the ResourceData every
	// result is written into.
	resourceSchema *schema.Resource
}

// NewDatacenterListResource creates a new list resource for ionoscloud_datacenter.
func NewDatacenterListResource() list.ListResource {
	return &datacenterListResource{resourceSchema: resourceDatacenter()}
}

// RawV6Schemas hands the framework the protocol schemas of the SDKv2 managed
// resource. A framework-native list resource inherits them from the resource
// itself; this is only needed because ionoscloud_datacenter lives on the SDKv2 side.
func (r *datacenterListResource) RawV6Schemas(ctx context.Context, _ list.RawV6SchemaRequest, resp *list.RawV6SchemaResponse) {
	fwidentity.SetRawV6Schemas(ctx, resp, datacenterResourceType, r.resourceSchema)
}

// Metadata returns the type name of the managed resource being listed. It must match
// the SDKv2 resource exactly, otherwise terraform has no resource to attach the
// results to.
func (r *datacenterListResource) Metadata(_ context.Context, _ resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = datacenterResourceType
}

// Configure stores the client bundle shared by the provider.
func (r *datacenterListResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
func (r *datacenterListResource) ListResourceConfigSchema(_ context.Context, _ list.ListResourceSchemaRequest, resp *list.ListResourceSchemaResponse) {
	resp.Schema = listschema.Schema{
		Attributes: map[string]listschema.Attribute{
			fwidentity.FiltersKey: fwidentity.FilterAttribute("name", "location"),
		},
	}
}

// List fetches every datacenter on the contract and streams the results. The Cloud
// API returns datacenters from all locations from a single collection, so unlike the
// regional products there is nothing to fan out over here.
func (r *datacenterListResource) List(ctx context.Context, req list.ListRequest, stream *list.ListResultsStream) {
	fwidentity.StreamList(ctx, stream, req,
		func(ctx context.Context) ([]ionoscloud.Datacenter, error) {
			// The datacenter collection is not location-scoped, so this uses the same
			// client every other global Cloud API listing uses.
			client, err := r.bundle.NewCloudAPIClientWithFailover(ctx)
			if err != nil {
				return nil, fmt.Errorf("failed to create the Cloud API client: %w", err)
			}

			// Depth(1) is what makes the API return the properties of every datacenter
			// instead of just its links. Filtering stays client-side, in the mapper.
			datacenters, apiResponse, err := client.DataCentersApi.DatacentersGet(ctx).Depth(1).Execute()
			if apiResponse != nil {
				tflog.Debug(ctx, "listed datacenters", map[string]any{"status_code": apiResponse.SafeStatusCode()})
			}
			if err != nil {
				return nil, fmt.Errorf("failed to list datacenters: %w", err)
			}
			if datacenters.Items == nil {
				return nil, nil
			}

			return *datacenters.Items, nil
		},
		r.mapDatacenter,
	)
}

// mapDatacenter maps a datacenter to an identity.MappedItem, or returns nil to skip it.
//
// The mapping itself is setDatacenterData, the same state writer resourceDatacenterRead
// uses, run against a ResourceData built from the live schema. Nothing here knows what
// attributes a datacenter has.
func (r *datacenterListResource) mapDatacenter(_ context.Context, includeResource bool, filters []fwidentity.Filter, dc ionoscloud.Datacenter) (*fwidentity.MappedItem, diag.Diagnostics) {
	var diags diag.Diagnostics

	if dc.Id == nil || dc.Properties == nil {
		return nil, nil
	}

	name := shared.ToValueDefault(dc.Properties.Name)
	location := shared.ToValueDefault(dc.Properties.Location)

	if !fwidentity.MatchesFilters(map[string]string{
		"name":     name,
		"location": location,
	}, filters) {
		return nil, nil
	}

	data := r.resourceSchema.Data(&terraform.InstanceState{})
	if err := setDatacenterData(data, &dc); err != nil {
		diags.AddError("Failed to map the datacenter", err.Error())
		return nil, diags
	}

	// setDatacenterIdentity reads id and location back out of the ResourceData, so it
	// has to run after setDatacenterData.
	if err := setDatacenterIdentity(data); err != nil {
		diags.AddError("Failed to map the datacenter identity", err.Error())
		return nil, diags
	}

	mapped, err := fwidentity.MappedItemFromResourceData(name, data, includeResource)
	if err != nil {
		diags.AddError("Failed to convert the datacenter state", err.Error())
		return nil, diags
	}

	return mapped, diags
}
