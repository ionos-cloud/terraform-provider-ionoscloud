package inmemorydbv2

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/bundleclient"
	inmemorydbv2service "github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/dbaas/inmemorydbv2"
)

var _ datasource.DataSourceWithConfigure = (*snapshotLocationsDataSource)(nil)

type snapshotLocationDataSourceModel struct {
	ID             types.String `tfsdk:"id"`
	SnapshotRegion types.String `tfsdk:"snapshot_region"`
}

func snapshotLocationDataSourceAttributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"id":              schema.StringAttribute{Computed: true, Description: "The ID (UUID) of the snapshot location."},
		"snapshot_region": schema.StringAttribute{Computed: true, Description: "The Object Storage region identifier (e.g. eu-central-3)."},
	}
}

type snapshotLocationsDataSource struct {
	bundle *bundleclient.SdkBundle
}

type snapshotLocationsDataSourceModel struct {
	Items    []snapshotLocationDataSourceModel `tfsdk:"items"`
	Location types.String                      `tfsdk:"location"`
}

// NewSnapshotLocationsDataSource creates a data source for listing snapshot locations.
func NewSnapshotLocationsDataSource() datasource.DataSource {
	return &snapshotLocationsDataSource{}
}

func (d *snapshotLocationsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_inmemorydb_snapshot_locations_v2"
}

func (d *snapshotLocationsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	clientBundle, ok := req.ProviderData.(*bundleclient.SdkBundle)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *bundleclient.SdkBundle, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}
	d.bundle = clientBundle
}

func (d *snapshotLocationsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Lists all InMemoryDB v2 snapshot locations for a given API endpoint location.",
		Attributes: map[string]schema.Attribute{
			"items": schema.ListNestedAttribute{
				Computed:    true,
				Description: "The list of snapshot locations.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: snapshotLocationDataSourceAttributes(),
				},
			},
			"location": schema.StringAttribute{
				Required:    true,
				Description: "The InMemoryDB API endpoint location to query. Available locations: " + inmemorydbv2service.AvailableLocationsString() + ".",
			},
		},
	}
}

func (d *snapshotLocationsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data snapshotLocationsDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	location := data.Location.ValueString()
	client, err := d.bundle.NewInMemoryDBV2Client(ctx, location)
	if err != nil {
		resp.Diagnostics.AddError("failed to create InMemoryDB v2 client", err.Error())
		return
	}

	list, _, err := client.ListSnapshotLocations(ctx)
	if err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("failed to list InMemoryDB v2 snapshot locations for API location %s", location), err.Error())
		return
	}

	items := make([]snapshotLocationDataSourceModel, 0, len(list.Items))
	for _, loc := range list.Items {
		item := snapshotLocationDataSourceModel{
			ID: types.StringValue(loc.Id),
		}
		item.SnapshotRegion = types.StringPointerValue(loc.Properties.Location)
		items = append(items, item)
	}
	data.Items = items

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
