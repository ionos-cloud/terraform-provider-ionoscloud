package pgsqlv2

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/bundleclient"
	pgsqlv2Service "github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/dbaas/pgsqlv2"
)

var _ datasource.DataSourceWithConfigure = (*backupLocationDataSource)(nil)

type backupLocationDataSource struct {
	bundle *bundleclient.SdkBundle
}

type backupLocationDataSourceModel struct {
	Location        types.String          `tfsdk:"location"`
	BackupLocations []backupLocationModel `tfsdk:"backup_locations"`
}

type backupLocationModel struct {
	ID       types.String `tfsdk:"id"`
	Location types.String `tfsdk:"location"`
}

// NewBackupLocationDataSource creates a new data source for listing PgSQL v2 backup locations.
func NewBackupLocationDataSource() datasource.DataSource {
	return &backupLocationDataSource{}
}

// Metadata returns the metadata for the backup location data source.
func (d *backupLocationDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_pg_backup_location"
}

// Configure configures the data source.
func (d *backupLocationDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	clientBundle, ok := req.ProviderData.(*bundleclient.SdkBundle)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *services.SdkBundle, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	d.bundle = clientBundle
}

// Schema returns the schema for the backup location data source.
func (d *backupLocationDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Lists PostgreSQL v2 backup locations.",
		Attributes: map[string]schema.Attribute{
			"location": schema.StringAttribute{
				Required:    true,
				Description: "The region in which to look up backup locations. Available locations: " + pgsqlv2Service.AvailableLocationsString() + ".",
			},
			"backup_locations": schema.ListNestedAttribute{
				Computed:    true,
				Description: "The list of available backup locations.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Computed:    true,
							Description: "The ID (UUID) of the backup location.",
						},
						"location": schema.StringAttribute{
							Computed:    true,
							Description: "The S3 location where the backup is stored.",
						},
					},
				},
			},
		},
	}
}

// Read reads the PgSQL v2 backup locations data source.
func (d *backupLocationDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data backupLocationDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	location := data.Location.ValueString()

	client, err := d.bundle.NewPgSQLV2Client(location)
	if err != nil {
		resp.Diagnostics.AddError("failed to create PostgreSQL v2 client", err.Error())
		return
	}

	backupLocationList, _, err := client.ListBackupLocations(ctx)
	if err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("failed to list PostgreSQL v2 backup locations in location: %s", location), err.Error())
		return
	}

	var backupLocations []backupLocationModel
	for _, bl := range backupLocationList.Items {
		item := backupLocationModel{
			ID: types.StringValue(bl.Id),
		}
		if bl.Properties.Location != nil {
			item.Location = types.StringValue(*bl.Properties.Location)
		}
		backupLocations = append(backupLocations, item)
	}

	data.BackupLocations = backupLocations

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
