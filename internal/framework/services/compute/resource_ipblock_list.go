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
	"github.com/ionos-cloud/sdk-go-bundle/shared"
	ionoscloud "github.com/ionos-cloud/sdk-go/v6"

	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/internal/framework/identity"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/bundleclient"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils/constant"
)

// A list resource for ionoscloud_ipblock, whose managed resource is still
// implemented with terraform-plugin-sdk/v2 and therefore lives on the other half of
// the mux. The protocol schemas the framework needs come from that resource via
// identity.SetRawV6Schemas, and the Identity it requires is declared alongside it in
// ionoscloud/resource_ipblock.go.

const ipblockResourceType = "ionoscloud_ipblock"

var (
	_ list.ListResource                 = (*ipblockListResource)(nil)
	_ list.ListResourceWithConfigure    = (*ipblockListResource)(nil)
	_ list.ListResourceWithRawV6Schemas = (*ipblockListResource)(nil)
)

// ipblockListResource lists ionoscloud_ipblock instances.
type ipblockListResource struct {
	bundle *bundleclient.SdkBundle

	// resourceSchema is the SDKv2 managed resource being listed, the source of the
	// protocol schemas returned by RawV6Schemas.
	resourceSchema *schema.Resource
}

// ipblockIdentityModel mirrors the resource identity declared by the SDKv2
// ionoscloud_ipblock resource.
type ipblockIdentityModel struct {
	ID       types.String `tfsdk:"id"`
	Location types.String `tfsdk:"location"`
}

// ipblockResourceModel mirrors the SDKv2 schema of ionoscloud_ipblock. It has to
// cover every attribute and block in that schema, because the framework fills the
// whole resource object from it.
//
// The fields are plain Go pointers rather than the types.String / types.Int64 values
// a framework-native resource would use. The schema behind this model is converted
// from the SDKv2 protocol schema, where an SDKv2 TypeInt arrives as a protocol
// Number and becomes a NumberAttribute - assigning a types.Int64 to it fails the
// type check. Plain Go types reflect into whichever framework type the converted
// schema ended up with, and the pointer carries the null/absent distinction.
type ipblockResourceModel struct {
	ID          *string                   `tfsdk:"id"`
	Name        *string                   `tfsdk:"name"`
	Location    *string                   `tfsdk:"location"`
	Size        *int32                    `tfsdk:"size"`
	IPs         *[]string                 `tfsdk:"ips"`
	IPConsumers *[]ipblockIPConsumerModel `tfsdk:"ip_consumers"`
	Timeouts    *ipblockTimeoutsModel     `tfsdk:"timeouts"`
}

// ipblockIPConsumerModel mirrors an entry of the ip_consumers block, which reports
// what each IP of the block is currently attached to.
type ipblockIPConsumerModel struct {
	IP              *string `tfsdk:"ip"`
	Mac             *string `tfsdk:"mac"`
	NicID           *string `tfsdk:"nic_id"`
	ServerID        *string `tfsdk:"server_id"`
	ServerName      *string `tfsdk:"server_name"`
	DatacenterID    *string `tfsdk:"datacenter_id"`
	DatacenterName  *string `tfsdk:"datacenter_name"`
	K8sNodepoolUUID *string `tfsdk:"k8s_nodepool_uuid"`
	K8sClusterUUID  *string `tfsdk:"k8s_cluster_uuid"`
}

// ipblockTimeoutsModel mirrors the timeouts block that SDKv2 adds to the schema.
// A listed IP block has no timeouts, so this is always left null.
type ipblockTimeoutsModel struct {
	Create  *string `tfsdk:"create"`
	Default *string `tfsdk:"default"`
	Delete  *string `tfsdk:"delete"`
	Update  *string `tfsdk:"update"`
}

// NewIPBlockListResource creates a new list resource for ionoscloud_ipblock.
// ipblockResource is the SDKv2 managed resource being listed; it is the source of
// the protocol schemas the framework needs.
func NewIPBlockListResource(ipblockResource *schema.Resource) list.ListResource {
	return &ipblockListResource{resourceSchema: ipblockResource}
}

// RawV6Schemas hands the framework the protocol schemas of the SDKv2 managed
// resource. A framework-native list resource inherits them from the resource
// itself; this is only needed because ionoscloud_ipblock lives on the SDKv2 side.
func (r *ipblockListResource) RawV6Schemas(ctx context.Context, _ list.RawV6SchemaRequest, resp *list.RawV6SchemaResponse) {
	identity.SetRawV6Schemas(ctx, resp, ipblockResourceType, r.resourceSchema)
}

// Metadata returns the type name of the managed resource being listed. It must match
// the SDKv2 resource exactly, otherwise terraform has no resource to attach the
// results to.
func (r *ipblockListResource) Metadata(_ context.Context, _ resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = ipblockResourceType
}

// Configure stores the client bundle shared by the provider.
func (r *ipblockListResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
func (r *ipblockListResource) ListResourceConfigSchema(_ context.Context, _ list.ListResourceSchemaRequest, resp *list.ListResourceSchemaResponse) {
	resp.Schema = listschema.Schema{
		Attributes: map[string]listschema.Attribute{
			identity.FiltersKey: identity.FilterAttribute("name", "location"),
		},
	}
}

// List fetches every IP block on the contract and streams the results. The Cloud API
// returns IP blocks from all locations from a single collection, so unlike the
// regional products there is nothing to fan out over here.
func (r *ipblockListResource) List(ctx context.Context, req list.ListRequest, stream *list.ListResultsStream) {
	identity.StreamList(ctx, stream, req,
		func(ctx context.Context) ([]ionoscloud.IpBlock, error) {
			// The IP block collection is not location-scoped, so this uses the same
			// client every other global Cloud API listing uses.
			client, err := r.bundle.NewCloudAPIClientWithFailover(ctx)
			if err != nil {
				return nil, fmt.Errorf("failed to create the Cloud API client: %w", err)
			}

			// Depth(1) is what makes the API return the properties of every IP block
			// instead of just its links. Filtering stays client-side, in the mapper.
			//
			// The explicit Limit matches ionoscloud/data_source_ipblock.go. Without it
			// the SDK client falls back to limit=100 (sdk-go/v6/api_ip_blocks.go:487-489,
			// which is the client's own default, not the endpoint's), an order of
			// magnitude below the 1000 most Cloud API collections return, so this
			// listing would silently see fewer IP blocks than the data source reading
			// the same collection.
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
func (r *ipblockListResource) mapIPBlock(_ context.Context, includeResource bool, filters []identity.Filter, ipBlock ionoscloud.IpBlock) (*identity.MappedItem, diag.Diagnostics) {
	if ipBlock.Id == nil || ipBlock.Properties == nil {
		return nil, nil
	}

	name := shared.ToValueDefault(ipBlock.Properties.Name)
	location := shared.ToValueDefault(ipBlock.Properties.Location)

	if !identity.MatchesFilters(map[string]string{
		"name":     name,
		"location": location,
	}, filters) {
		return nil, nil
	}

	// name is Optional on ionoscloud_ipblock, so it can legitimately be absent. The
	// display name is the label `terraform query` prints for each result, and an
	// empty one renders as a blank row, so fall back to the ID.
	displayName := name
	if displayName == "" {
		displayName = *ipBlock.Id
	}

	mapped := &identity.MappedItem{
		DisplayName: displayName,
		Identity: &ipblockIdentityModel{
			ID:       types.StringValue(*ipBlock.Id),
			Location: types.StringValue(location),
		},
	}

	if !includeResource {
		return mapped, nil
	}

	mapped.Resource = &ipblockResourceModel{
		ID:          ipBlock.Id,
		Name:        ipBlock.Properties.Name,
		Location:    ipBlock.Properties.Location,
		Size:        ipBlock.Properties.Size,
		IPs:         ipBlock.Properties.Ips,
		IPConsumers: mapIPBlockConsumers(ipBlock.Properties.IpConsumers),
	}

	return mapped, nil
}

// mapIPBlockConsumers maps the consumption detail reported for each IP of a block.
func mapIPBlockConsumers(consumers *[]ionoscloud.IpConsumer) *[]ipblockIPConsumerModel {
	if consumers == nil {
		return nil
	}

	mapped := make([]ipblockIPConsumerModel, 0, len(*consumers))
	for _, consumer := range *consumers {
		mapped = append(mapped, ipblockIPConsumerModel{
			IP:              consumer.Ip,
			Mac:             consumer.Mac,
			NicID:           consumer.NicId,
			ServerID:        consumer.ServerId,
			ServerName:      consumer.ServerName,
			DatacenterID:    consumer.DatacenterId,
			DatacenterName:  consumer.DatacenterName,
			K8sNodepoolUUID: consumer.K8sNodePoolUuid,
			K8sClusterUUID:  consumer.K8sClusterUuid,
		})
	}

	return &mapped
}
