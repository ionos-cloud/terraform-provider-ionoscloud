package compute

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/list"
	listschema "github.com/hashicorp/terraform-plugin-framework/list/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	ionoscloud "github.com/ionos-cloud/sdk-go/v6"

	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/internal/framework/identity"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/bundleclient"
)

// A list resource for ionoscloud_datacenter, whose managed resource is still
// implemented with terraform-plugin-sdk/v2 and therefore lives on the other half of
// the mux. The protocol schemas the framework needs come from that resource via
// identity.SetRawV6Schemas, and the Identity it requires is declared alongside it in
// ionoscloud/resource_datacenter.go.

const datacenterResourceType = "ionoscloud_datacenter"

var (
	_ list.ListResource                 = (*datacenterListResource)(nil)
	_ list.ListResourceWithConfigure    = (*datacenterListResource)(nil)
	_ list.ListResourceWithRawV6Schemas = (*datacenterListResource)(nil)
)

// datacenterListResource lists ionoscloud_datacenter instances.
type datacenterListResource struct {
	bundle *bundleclient.SdkBundle

	// resourceSchema is the SDKv2 managed resource being listed, the source of the
	// protocol schemas returned by RawV6Schemas.
	resourceSchema *schema.Resource
}

// datacenterIdentityModel mirrors the resource identity declared by the SDKv2
// ionoscloud_datacenter resource.
type datacenterIdentityModel struct {
	ID       types.String `tfsdk:"id"`
	Location types.String `tfsdk:"location"`
}

// datacenterResourceModel mirrors the SDKv2 schema of ionoscloud_datacenter. It has
// to cover every attribute and block in that schema, because the framework fills the
// whole resource object from it.
//
// The fields are plain Go pointers rather than the types.String / types.Int64 values
// a framework-native resource would use. The schema behind this model is converted
// from the SDKv2 protocol schema, where an SDKv2 TypeInt arrives as a protocol
// Number and becomes a NumberAttribute - assigning a types.Int64 to it fails the
// type check. Plain Go types reflect into whichever framework type the converted
// schema ended up with, and the pointer carries the null/absent distinction.
type datacenterResourceModel struct {
	ID                *string                           `tfsdk:"id"`
	Name              *string                           `tfsdk:"name"`
	Location          *string                           `tfsdk:"location"`
	Description       *string                           `tfsdk:"description"`
	SecAuthProtection *bool                             `tfsdk:"sec_auth_protection"`
	Version           *int32                            `tfsdk:"version"`
	Features          *[]string                         `tfsdk:"features"`
	CPUArchitecture   *[]datacenterCPUArchitectureModel `tfsdk:"cpu_architecture"`
	IPv6CidrBlock     *string                           `tfsdk:"ipv6_cidr_block"`
	Timeouts          *datacenterTimeoutsModel          `tfsdk:"timeouts"`
}

// datacenterCPUArchitectureModel mirrors an entry of the cpu_architecture attribute.
type datacenterCPUArchitectureModel struct {
	CPUFamily *string `tfsdk:"cpu_family"`
	MaxCores  *int32  `tfsdk:"max_cores"`
	MaxRAM    *int32  `tfsdk:"max_ram"`
	Vendor    *string `tfsdk:"vendor"`
}

// datacenterTimeoutsModel mirrors the timeouts block that SDKv2 adds to the schema.
// A listed datacenter has no timeouts, so this is always left null.
type datacenterTimeoutsModel struct {
	Create  *string `tfsdk:"create"`
	Default *string `tfsdk:"default"`
	Delete  *string `tfsdk:"delete"`
	Update  *string `tfsdk:"update"`
}

// NewDatacenterListResource creates a new list resource for ionoscloud_datacenter.
// datacenterResource is the SDKv2 managed resource being listed; it is the source of
// the protocol schemas the framework needs.
func NewDatacenterListResource(datacenterResource *schema.Resource) list.ListResource {
	return &datacenterListResource{resourceSchema: datacenterResource}
}

// RawV6Schemas hands the framework the protocol schemas of the SDKv2 managed
// resource. A framework-native list resource inherits them from the resource
// itself; this is only needed because ionoscloud_datacenter lives on the SDKv2 side.
func (r *datacenterListResource) RawV6Schemas(ctx context.Context, _ list.RawV6SchemaRequest, resp *list.RawV6SchemaResponse) {
	identity.SetRawV6Schemas(ctx, resp, datacenterResourceType, r.resourceSchema)
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
			identity.FiltersKey: identity.FilterAttribute("name", "location"),
		},
	}
}

// List fetches every datacenter on the contract and streams the results. The Cloud
// API returns datacenters from all locations from a single collection, so unlike the
// regional products there is nothing to fan out over here.
func (r *datacenterListResource) List(ctx context.Context, req list.ListRequest, stream *list.ListResultsStream) {
	identity.StreamList(ctx, stream, req,
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
func (r *datacenterListResource) mapDatacenter(_ context.Context, includeResource bool, filters []identity.Filter, dc ionoscloud.Datacenter) (*identity.MappedItem, diag.Diagnostics) {
	if dc.Id == nil || dc.Properties == nil {
		return nil, nil
	}

	name := valueOrZero(dc.Properties.Name)
	location := valueOrZero(dc.Properties.Location)

	if !identity.MatchesFilters(map[string]string{
		"name":     name,
		"location": location,
	}, filters) {
		return nil, nil
	}

	mapped := &identity.MappedItem{
		DisplayName: name,
		Identity: &datacenterIdentityModel{
			ID:       types.StringValue(*dc.Id),
			Location: types.StringValue(location),
		},
	}

	if !includeResource {
		return mapped, nil
	}

	mapped.Resource = &datacenterResourceModel{
		ID:                dc.Id,
		Name:              dc.Properties.Name,
		Location:          dc.Properties.Location,
		Description:       dc.Properties.Description,
		SecAuthProtection: dc.Properties.SecAuthProtection,
		Version:           dc.Properties.Version,
		Features:          dc.Properties.Features,
		CPUArchitecture:   mapDatacenterCPUArchitecture(dc.Properties.CpuArchitecture),
		IPv6CidrBlock:     dc.Properties.Ipv6CidrBlock,
	}

	return mapped, nil
}

// mapDatacenterCPUArchitecture maps the CPU architectures reported for a datacenter.
func mapDatacenterCPUArchitecture(architectures *[]ionoscloud.CpuArchitectureProperties) *[]datacenterCPUArchitectureModel {
	if architectures == nil {
		return nil
	}

	mapped := make([]datacenterCPUArchitectureModel, 0, len(*architectures))
	for _, architecture := range *architectures {
		mapped = append(mapped, datacenterCPUArchitectureModel{
			CPUFamily: architecture.CpuFamily,
			MaxCores:  architecture.MaxCores,
			MaxRAM:    architecture.MaxRam,
			Vendor:    architecture.Vendor,
		})
	}

	return &mapped
}

// valueOrZero dereferences ptr, or returns the zero value if it is nil.
func valueOrZero[T any](ptr *T) T {
	if ptr == nil {
		var zero T
		return zero
	}

	return *ptr
}
