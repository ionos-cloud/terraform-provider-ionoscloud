package compute

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/list"
	listschema "github.com/hashicorp/terraform-plugin-framework/list/schema"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	ionoscloud "github.com/ionos-cloud/sdk-go/v6"

	dcschema "github.com/ionos-cloud/terraform-provider-ionoscloud/v6/internal/compute/datacenter"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/internal/framework/identity"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/bundleclient"
)

// datacenterSchemaResource returns a schema-only *schema.Resource (no CRUD,
// built from the same schema/identity schema the ionoscloud_datacenter SDKv2
// resource uses) for use in building raw V5 values and schemas for the list
// resource.
func datacenterSchemaResource() *schema.Resource {
	return &schema.Resource{
		Identity: &schema.ResourceIdentity{
			SchemaFunc: dcschema.IdentitySchema,
		},
		Schema: dcschema.Schema(),
	}
}

// datacenterToTfValues converts a Cloud API Datacenter into the identity and
// resource tftypes.Value pair expected by a list.ListResult, by round-tripping
// it through an SDKv2 ResourceData instance (the same field-mapping code the
// ionoscloud_datacenter resource and data source use).
func datacenterToTfValues(dc ionoscloud.Datacenter) (identity *tftypes.Value, resource *tftypes.Value, err error) {
	rd := datacenterSchemaResource().Data(nil)

	if err := dcschema.PopulateResourceData(rd, &dc); err != nil {
		return nil, nil, fmt.Errorf("failed to populate datacenter resource data: %w", err)
	}

	if dc.Id == nil {
		return nil, nil, fmt.Errorf("datacenter has no id")
	}

	identityData, err := rd.Identity()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get identity data: %w", err)
	}
	if err := identityData.Set("id", *dc.Id); err != nil {
		return nil, nil, fmt.Errorf("failed to set identity id: %w", err)
	}

	identityVal, err := rd.TfTypeIdentityState()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to convert identity state: %w", err)
	}

	resourceVal, err := rd.TfTypeResourceState()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to convert resource state: %w", err)
	}

	return identityVal, resourceVal, nil
}

var (
	_ list.ListResource                 = (*datacenterListResource)(nil)
	_ list.ListResourceWithConfigure    = (*datacenterListResource)(nil)
	_ list.ListResourceWithRawV5Schemas = (*datacenterListResource)(nil)
)

type datacenterListResource struct {
	bundle *bundleclient.SdkBundle
}

// NewDatacenterListResource creates a new list resource for ionoscloud_datacenter.
func NewDatacenterListResource() list.ListResource {
	return &datacenterListResource{}
}

// Metadata returns the full name of the list resource. This must match the
// ionoscloud_datacenter managed resource's full name exactly.
func (r *datacenterListResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_datacenter"
}

// Configure receives the shared SDK bundle. list.ListResourceWithConfigure's
// Configure method is declared using resource.ConfigureRequest/resource.ConfigureResponse
// (from the "resource" package, not "list"), matching resource.ResourceWithConfigure's
// signature so a single implementation can satisfy both interfaces if needed.
func (r *datacenterListResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	bundle, ok := req.ProviderData.(*bundleclient.SdkBundle)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Provider Data Type",
			fmt.Sprintf("Expected *bundleclient.SdkBundle, got: %T", req.ProviderData),
		)
		return
	}

	r.bundle = bundle
}

// ListResourceConfigSchema returns the schema for the list block's config.
func (r *datacenterListResource) ListResourceConfigSchema(_ context.Context, _ list.ListResourceSchemaRequest, resp *list.ListResourceSchemaResponse) {
	resp.Schema = listschema.Schema{
		Attributes: map[string]listschema.Attribute{
			identity.FiltersKey: identity.FilterAttribute("name", "location"),
		},
	}
}

// RawV5Schemas provides the ProtoV5 representations of the ionoscloud_datacenter
// resource and identity schemas, since it is an SDKv2-defined resource rather
// than a framework-native one.
func (r *datacenterListResource) RawV5Schemas(ctx context.Context, _ list.RawV5SchemaRequest, resp *list.RawV5SchemaResponse) {
	dcResource := datacenterSchemaResource()
	resp.ProtoV5Schema = dcResource.ProtoSchema(ctx)()
	resp.ProtoV5IdentitySchema = dcResource.ProtoIdentitySchema(ctx)()
}

// List streams datacenters matching the given filters.
func (r *datacenterListResource) List(ctx context.Context, req list.ListRequest, stream *list.ListResultsStream) {
	var filters []identity.Filter
	diags := req.Config.GetAttribute(ctx, path.Root(identity.FiltersKey), &filters)
	if diags.HasError() {
		stream.Results = list.ListResultsStreamDiagnostics(diags)
		return
	}

	nameFilter := identity.FilterValue(filters, "name")
	locationFilter := identity.FilterValue(filters, "location")

	client, err := r.bundle.NewCloudAPIClientWithFailover(ctx)
	if err != nil {
		var d diag.Diagnostics
		d.AddError("Failed to create Cloud API client", err.Error())
		stream.Results = list.ListResultsStreamDiagnostics(d)
		return
	}

	request := client.DataCentersApi.DatacentersGet(ctx).Depth(1)
	if nameFilter != "" {
		request = request.Filter("name", nameFilter)
	}
	if locationFilter != "" {
		request = request.Filter("location", locationFilter)
	}

	datacenters, _, err := request.Execute()
	if err != nil {
		var d diag.Diagnostics
		d.AddError("Failed to list datacenters", err.Error())
		stream.Results = list.ListResultsStreamDiagnostics(d)
		return
	}

	items := datacenters.Items
	stream.Results = func(push func(list.ListResult) bool) {
		if items == nil {
			return
		}
		for _, dc := range *items {
			result := req.NewListResult(ctx)

			if dc.Properties != nil && dc.Properties.Name != nil {
				result.DisplayName = *dc.Properties.Name
			}

			identityVal, resourceVal, err := datacenterToTfValues(dc)
			if err != nil {
				result.Diagnostics.AddError("Failed to convert datacenter", err.Error())
				push(result)
				return
			}

			result.Identity.Raw = *identityVal
			if req.IncludeResource {
				result.Resource.Raw = *resourceVal
			}

			if !push(result) {
				return
			}
		}
	}
}
