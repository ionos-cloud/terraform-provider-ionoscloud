package mariadbv2

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/bundleclient"
	mariadbv2service "github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/dbaas/mariadbv2"
	diagutil "github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils/diags"
)

var _ datasource.DataSourceWithConfigure = (*backupLocationsDataSource)(nil)

type backupLocationDataSourceModel struct {
	ID       types.String `tfsdk:"id"`
	Location types.String `tfsdk:"location"`
}

func backupLocationDataSourceAttributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"id":       schema.StringAttribute{Computed: true, Description: "The ID (UUID) of the backup location."},
		"location": schema.StringAttribute{Computed: true, Description: "The Object Storage location name."},
	}
}

type backupLocationsDataSource struct {
	bundle *bundleclient.SdkBundle
}

type backupLocationsDataSourceModel struct {
	Items    []backupLocationDataSourceModel `tfsdk:"items"`
	Location types.String                    `tfsdk:"location"`
}

// NewBackupLocationsDataSource creates a new data source for listing MariaDB v2 backup locations.
func NewBackupLocationsDataSource() datasource.DataSource {
	return &backupLocationsDataSource{}
}

func (d *backupLocationsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_mariadb_backup_locations_v2"
}

func (d *backupLocationsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *backupLocationsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Lists all supported IONOS CLOUD MariaDB V2 backup locations in a given region.",
		Attributes: map[string]schema.Attribute{
			"items": schema.ListNestedAttribute{
				Computed:    true,
				Description: "The list of backup locations.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: backupLocationDataSourceAttributes(),
				},
			},
			"location": schema.StringAttribute{
				Required:    true,
				Description: "The location to query. Available locations: " + mariadbv2service.AvailableLocationsString() + ".",
			},
		},
	}
}

func (d *backupLocationsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data backupLocationsDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	location := data.Location.ValueString()
	client, err := d.bundle.NewMariaDBV2Client(ctx, location)
	if err != nil {
		resp.Diagnostics.AddError("failed to create MariaDB V2 client", err.Error())
		return
	}

	list, apiResponse, err := client.ListBackupLocations(ctx)
	if err != nil {
		resp.Diagnostics.AddError("failed to list MariaDB V2 backup locations", diagutil.WrapError(err, &diagutil.ErrorContext{
			StatusCode:     apiResponse.SafeStatusCode(),
			AdditionalInfo: map[string]string{"location": location},
		}).Error())
		return
	}

	items := make([]backupLocationDataSourceModel, 0, len(list.Items))
	for _, bl := range list.Items {
		item := backupLocationDataSourceModel{
			ID:       types.StringValue(bl.Id),
			Location: types.StringPointerValue(bl.Properties.Location),
		}
		items = append(items, item)
	}
	data.Items = items

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
