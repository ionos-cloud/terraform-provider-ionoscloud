package inmemorydbv2

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/bundleclient"
	inmemorydbv2Service "github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/dbaas/inmemorydbv2"
)

var _ datasource.DataSourceWithConfigure = (*clustersDataSource)(nil)

type clustersDataSource struct {
	bundle *bundleclient.SdkBundle
}

type clustersDataSourceModel struct {
	Location types.String             `tfsdk:"location"`
	Name     types.String             `tfsdk:"name"`
	Items    []clusterDataSourceModel `tfsdk:"items"`
}

// NewClustersDataSource creates a new data source for listing InMemoryDB v2 clusters.
func NewClustersDataSource() datasource.DataSource {
	return &clustersDataSource{}
}

func (d *clustersDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_inmemorydb_clusters_v2"
}

func (d *clustersDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *clustersDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Lists all InMemoryDB v2 clusters in a given location, with optional name filter.",
		Attributes: map[string]schema.Attribute{
			"location": schema.StringAttribute{
				Required:    true,
				Description: "The location to query. Available locations: " + inmemorydbv2Service.AvailableLocationsString() + ".",
			},
			"name": schema.StringAttribute{
				Optional:    true,
				Description: "Filter clusters by name (partial match, case-insensitive).",
			},
			"items": schema.ListNestedAttribute{
				Computed:    true,
				Description: "The list of clusters.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: clusterDataSourceAttributes(),
				},
			},
		},
	}
}

func (d *clustersDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data clustersDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	location := data.Location.ValueString()
	client, err := d.bundle.NewInMemoryDBV2Client(location)
	if err != nil {
		resp.Diagnostics.AddError("failed to create InMemoryDB v2 client", err.Error())
		return
	}

	list, _, err := client.ListClusters(ctx, data.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("failed to list InMemoryDB v2 clusters in location %s", location), err.Error())
		return
	}

	items := make([]clusterDataSourceModel, 0, len(list.Items))
	for i := range list.Items {
		var item clusterDataSourceModel
		item.Location = data.Location
		mapClusterResponseToDataSourceModel(ctx, &list.Items[i], &item)
		items = append(items, item)
	}
	data.Items = items

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
