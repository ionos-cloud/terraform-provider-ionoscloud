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
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/internal/framework/sdkv2schema"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/bundleclient"
)

// This is the reference example for adding a list resource to a resource that is
// still implemented with terraform-plugin-sdk/v2, such as ionoscloud_datacenter.
//
// List resources cannot be written in SDKv2 - helper/schema answers every
// ListResource RPC with "list resource type is not supported by this provider" -
// so the list resource is registered on the plugin-framework side of the mux while
// the managed resource stays where it is. tf6muxserver keeps managed resources and
// list resources in separate routing tables, so the two halves can be served by
// different providers under the same ionoscloud_datacenter type name.
//
// Two things follow from the managed resource living elsewhere:
//
//   - The framework has no schema for ionoscloud_datacenter, so this list resource
//     implements list.ListResourceWithRawV6Schemas and hands it the protocol
//     schemas that the SDKv2 provider itself advertises (see sdkv2schema).
//     Without them the framework refuses to register the list resource at all:
//     "ListResource Type Defined without a Matching Managed Resource Type".
//   - The SDKv2 resource has to declare a schema.ResourceIdentity, since terraform
//     identifies every listed instance by its identity. See the Identity block in
//     ionoscloud/resource_datacenter.go.
//
// Everything else - the filters attribute, the streaming, the mapper contract - is
// the same as for a framework-native list resource such as pgsqlv2's pg_cluster_v2.

const datacenterResourceType = "ionoscloud_datacenter"

var (
	_ list.ListResource                 = (*datacenterListResource)(nil)
	_ list.ListResourceWithConfigure    = (*datacenterListResource)(nil)
	_ list.ListResourceWithRawV6Schemas = (*datacenterListResource)(nil)
)

// datacenterListResource lists ionoscloud_datacenter instances.
type datacenterListResource struct {
	bundle *bundleclient.SdkBundle

	// sdkv2Provider is the SDKv2 half of the muxed provider, the source of the
	// protocol schemas returned by RawV6Schemas.
	sdkv2Provider *schema.Provider
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
// sdkv2Provider is the provider that defines the managed resource being listed.
func NewDatacenterListResource(sdkv2Provider *schema.Provider) list.ListResource {
	return &datacenterListResource{sdkv2Provider: sdkv2Provider}
}

// Metadata returns the type name of the managed resource being listed. It must match
// the SDKv2 resource exactly, otherwise terraform has no resource to attach the
// results to.
func (r *datacenterListResource) Metadata(_ context.Context, _ resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = datacenterResourceType
}

// RawV6Schemas hands the framework the protocol schemas of the SDKv2 managed
// resource. This is only needed because ionoscloud_datacenter is not a framework
// resource; a framework-native list resource inherits them from the resource itself.
func (r *datacenterListResource) RawV6Schemas(ctx context.Context, _ list.RawV6SchemaRequest, resp *list.RawV6SchemaResponse) {
	resourceSchema, identitySchema, err := sdkv2schema.Schemas(ctx, r.sdkv2Provider, datacenterResourceType)
	if err != nil {
		// RawV6Schemas cannot report diagnostics. Leaving the schemas unset makes the
		// framework reject the list resource with an actionable error of its own, so
		// log the cause here to explain why.
		tflog.Error(ctx, "failed to read the SDKv2 schemas for the datacenter list resource", map[string]any{
			"resource_type": datacenterResourceType,
			"error":         err.Error(),
		})
		return
	}

	resp.ProtoV6Schema = resourceSchema
	resp.ProtoV6IdentitySchema = identitySchema
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
