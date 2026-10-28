package mongov2

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/bundleclient"
	mongov2service "github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/dbaas/mongov2"
	diagutil "github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils/diags"
)

var _ datasource.DataSourceWithConfigure = (*snapshotLocationsDataSource)(nil)

type snapshotLocationDataSourceModel struct {
	ID       types.String `tfsdk:"id"`
	Location types.String `tfsdk:"location"`
}

func snapshotLocationDataSourceAttributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"id":       schema.StringAttribute{Computed: true, Description: "The ID (UUID) of the snapshot location."},
		"location": schema.StringAttribute{Computed: true, Description: "The Object Storage location name."},
	}
}

type snapshotLocationsDataSource struct {
	bundle *bundleclient.SdkBundle
}

type snapshotLocationsDataSourceModel struct {
	Items    []snapshotLocationDataSourceModel `tfsdk:"items"`
	Location types.String                      `tfsdk:"location"`
}

// NewSnapshotLocationsDataSource creates a new data source for listing MongoDB v2 snapshot locations.
func NewSnapshotLocationsDataSource() datasource.DataSource {
	return &snapshotLocationsDataSource{}
}

func (d *snapshotLocationsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_mongo_snapshot_locations_v2"
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
		Description: "Lists all supported IONOS CLOUD MongoDB V2 snapshot locations in a given region.",
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
				Description: "The location to query. Available locations: " + mongov2service.AvailableLocationsString() + ".",
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
	client, err := d.bundle.NewMongoDBV2Client(ctx, location)
	if err != nil {
		resp.Diagnostics.AddError("failed to create MongoDB V2 client", err.Error())
		return
	}

	list, apiResponse, err := client.ListSnapshotLocations(ctx)
	if err != nil {
		resp.Diagnostics.AddError("failed to list MongoDB V2 snapshot locations", diagutil.WrapError(err, &diagutil.ErrorContext{
			StatusCode:     apiResponse.SafeStatusCode(),
			AdditionalInfo: map[string]string{"location": location},
		}).Error())
		return
	}

	items := make([]snapshotLocationDataSourceModel, 0, len(list.Items))
	for _, sl := range list.Items {
		items = append(items, snapshotLocationDataSourceModel{
			ID:       types.StringValue(sl.Id),
			Location: types.StringPointerValue(sl.Properties.Location),
		})
	}
	data.Items = items

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
