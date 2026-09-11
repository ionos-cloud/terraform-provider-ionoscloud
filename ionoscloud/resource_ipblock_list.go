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
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils/constant"
)

// A list resource for ionoscloud_ipblock, whose managed resource is still implemented
// with terraform-plugin-sdk/v2 and therefore lives on the other half of the mux. The
// protocol schemas the framework needs come from that resource via
// identity.SetRawV6Schemas.
//
// It lives in this package, next to the resource it lists, so it can call
// resourceIPBlock, IpBlockSetData and setIPBlockIdentity directly. That is the whole
// design: results are produced by the resource's own state writer, so this file
// declares no model of the ipblock schema and there is nothing here to keep in sync
// when that schema changes.

var (
	_ list.ListResource                 = (*ipBlockListResource)(nil)
	_ list.ListResourceWithConfigure    = (*ipBlockListResource)(nil)
	_ list.ListResourceWithRawV6Schemas = (*ipBlockListResource)(nil)
)

// ipBlockListResource lists ionoscloud_ipblock instances.
type ipBlockListResource struct {
	bundle *bundleclient.SdkBundle

	// resourceSchema is the SDKv2 managed resource being listed: the source of the
	// protocol schemas returned by RawV6Schemas, and of the ResourceData every
	// result is written into.
	resourceSchema *schema.Resource
}

// NewIPBlockListResource creates a new list resource for ionoscloud_ipblock.
func NewIPBlockListResource() list.ListResource {
	return &ipBlockListResource{resourceSchema: resourceIPBlock()}
}

// RawV6Schemas hands the framework the protocol schemas of the SDKv2 managed
// resource. A framework-native list resource inherits them from the resource itself;
// this is only needed because ionoscloud_ipblock lives on the SDKv2 side.
func (r *ipBlockListResource) RawV6Schemas(ctx context.Context, _ list.RawV6SchemaRequest, resp *list.RawV6SchemaResponse) {
	fwidentity.SetRawV6Schemas(ctx, resp, constant.IpBlockResource, r.resourceSchema)
}

// Metadata returns the type name of the managed resource being listed. It must match
// the SDKv2 resource exactly, otherwise terraform has no resource to attach the
// results to - hence the same constant that keys ResourcesMap in provider.go.
func (r *ipBlockListResource) Metadata(_ context.Context, _ resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = constant.IpBlockResource
}

// Configure stores the client bundle shared by the provider.
func (r *ipBlockListResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
func (r *ipBlockListResource) ListResourceConfigSchema(_ context.Context, _ list.ListResourceSchemaRequest, resp *list.ListResourceSchemaResponse) {
	resp.Schema = listschema.Schema{
		Attributes: map[string]listschema.Attribute{
			fwidentity.FiltersKey: fwidentity.FilterAttribute("name", "location"),
		},
	}
}

// List fetches every IP block on the contract and streams the results. The Cloud API
// returns IP blocks from all locations from a single collection, so unlike the
// regional products there is nothing to fan out over here.
func (r *ipBlockListResource) List(ctx context.Context, req list.ListRequest, stream *list.ListResultsStream) {
	fwidentity.StreamList(ctx, stream, req,
		func(ctx context.Context) ([]ionoscloud.IpBlock, error) {
			// The ipblock collection is not location-scoped, so this uses the same
			// client every other global Cloud API listing uses.
			client, err := r.bundle.NewCloudAPIClientWithFailover(ctx)
			if err != nil {
				return nil, fmt.Errorf("failed to create the Cloud API client: %w", err)
			}

			// Depth(1) is what makes the API return the properties of every IP block
			// instead of just its links. Filtering stays client-side, in the mapper.
			//
			// The explicit Limit matches ionoscloud/data_source_ipblock.go: the SDK
			// client falls back to limit=100 for /ipblocks (api_ip_blocks.go:489),
			// which is ten times lower than what the data source asks for, so without
			// it the list resource would see less than its own data source does.
			ipBlocks, apiResponse, err := client.IPBlocksApi.IpblocksGet(ctx).Depth(1).Limit(constant.IPBlockLimit).Execute()
			if apiResponse != nil {
				tflog.Debug(ctx, "listed ip blocks", map[string]any{"status_code": apiResponse.SafeStatusCode()})
			}
			if err != nil {
				return nil, fmt.Errorf("failed to list ip blocks: %w", err)
			}
			if ipBlocks.Items == nil {
				return nil, nil
			}

			return *ipBlocks.Items, nil
		},
		r.mapIPBlock,
	)
}

// mapIPBlock maps an IP block to an identity.MappedItem, or returns nil to skip it.
//
// The mapping itself is IpBlockSetData, the same state writer resourceIPBlockRead
// uses, run against a ResourceData built from the live schema. Nothing here knows what
// attributes an IP block has.
func (r *ipBlockListResource) mapIPBlock(_ context.Context, includeResource bool, filters []fwidentity.Filter, ipBlock ionoscloud.IpBlock) (*fwidentity.MappedItem, diag.Diagnostics) {
	var diags diag.Diagnostics

	if ipBlock.Id == nil || ipBlock.Properties == nil {
		return nil, nil
	}

	name := shared.ToValueDefault(ipBlock.Properties.Name)
	location := shared.ToValueDefault(ipBlock.Properties.Location)

	if !fwidentity.MatchesFilters(map[string]string{
		"name":     name,
		"location": location,
	}, filters) {
		return nil, nil
	}

	// `name` is Optional on ionoscloud_ipblock, so an unnamed IP block would otherwise
	// render as a blank row in the terraform query output.
	displayName := name
	if displayName == "" {
		displayName = *ipBlock.Id
	}

	data := r.resourceSchema.Data(&terraform.InstanceState{})
	if err := IpBlockSetData(data, &ipBlock); err != nil {
		diags.AddError("Failed to map the ip block", err.Error())
		return nil, diags
	}

	// setIPBlockIdentity reads id and location back out of the ResourceData, so it has
	// to run after IpBlockSetData.
	if err := setIPBlockIdentity(data); err != nil {
		diags.AddError("Failed to map the ip block identity", err.Error())
		return nil, diags
	}

	mapped, err := fwidentity.MappedItemFromResourceData(displayName, data, includeResource)
	if err != nil {
		diags.AddError("Failed to convert the ip block state", err.Error())
		return nil, diags
	}

	return mapped, diags
}
