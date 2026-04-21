package pgsqlv2

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	pgsqlv2 "github.com/ionos-cloud/sdk-go-bundle/products/dbaas/psql/v3"

	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/bundleclient"
	pgsqlv2Service "github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/dbaas/pgsqlv2"
)

var _ datasource.DataSourceWithConfigure = (*backupsDataSource)(nil)

type backupsDataSource struct {
	bundle *bundleclient.SdkBundle
}

type backupsDataSourceModel struct {
	Location  types.String  `tfsdk:"location"`
	ClusterID types.String  `tfsdk:"cluster_id"`
	Backups   []backupModel `tfsdk:"backups"`
}

type backupModel struct {
	ID                         types.String `tfsdk:"id"`
	ClusterID                  types.String `tfsdk:"cluster_id"`
	PostgresClusterVersion     types.String `tfsdk:"postgres_cluster_version"`
	IsActive                   types.Bool   `tfsdk:"is_active"`
	EarliestRecoveryTargetTime types.String `tfsdk:"earliest_recovery_target_time"`
	LatestRecoveryTargetTime   types.String `tfsdk:"latest_recovery_target_time"`
	Location                   types.String `tfsdk:"location"`
}

// NewBackupsDataSource creates a new data source for listing PgSQL v2 backups.
func NewBackupsDataSource() datasource.DataSource {
	return &backupsDataSource{}
}

// Metadata returns the metadata for the backups data source.
func (d *backupsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_pg_backups_v2"
}

// Configure configures the data source.
func (d *backupsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

// Schema returns the schema for the backups data source.
func (d *backupsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Lists PostgreSQL v2 backups.",
		Attributes: map[string]schema.Attribute{
			"backups": schema.ListNestedAttribute{
				Computed:    true,
				Description: "The list of backups.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"cluster_id": schema.StringAttribute{
							Computed:    true,
							Description: "The ID (UUID) of the cluster the backup belongs to.",
						},
						"earliest_recovery_target_time": schema.StringAttribute{
							Computed:    true,
							Description: "The earliest point in time to which the cluster can be restored.",
						},
						"id": schema.StringAttribute{
							Computed:    true,
							Description: "The ID (UUID) of the backup.",
						},
						"is_active": schema.BoolAttribute{
							Computed:    true,
							Description: "Whether the backup is active.",
						},
						"latest_recovery_target_time": schema.StringAttribute{
							Computed:    true,
							Description: "The latest point in time to which the cluster can be restored. If the backup can be restored up to the current time, this field will be null.",
						},
						"location": schema.StringAttribute{
							Computed:    true,
							Description: "The S3 location where the backup is stored.",
						},
						"postgres_cluster_version": schema.StringAttribute{
							Computed:    true,
							Description: "The PostgreSQL version of the cluster.",
						},
					},
				},
			},
			"cluster_id": schema.StringAttribute{
				Optional:    true,
				Description: "The ID (UUID) of the cluster to filter backups by.",
			},
			"location": schema.StringAttribute{
				Required:    true,
				Description: "The region in which to look up backups. Available locations: " + pgsqlv2Service.AvailableLocationsString() + ".",
			},
		},
	}
}

// Read reads the PgSQL v2 backups data source.
func (d *backupsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data backupsDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	location := data.Location.ValueString()
	clusterID := data.ClusterID.ValueString()

	client, err := d.bundle.NewPgSQLV2Client(location)
	if err != nil {
		resp.Diagnostics.AddError("failed to create PostgreSQL v2 client", err.Error())
		return
	}

	backupList, _, err := client.ListBackups(ctx, clusterID)
	if err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("failed to list PostgreSQL v2 backups for cluster %s in location: %s", clusterID, location), err.Error())
		return
	}

	var backups []backupModel
	for _, b := range backupList.Items {
		var item backupModel
		mapBackupResponseToModel(&b, &item)
		backups = append(backups, item)
	}

	data.Backups = backups

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// mapBackupResponseToModel maps the API backup response to the data source model.
func mapBackupResponseToModel(backup *pgsqlv2.BackupRead, model *backupModel) {
	model.ID = types.StringValue(backup.Id)
	props := &backup.Properties
	model.ClusterID = types.StringPointerValue(props.ClusterId)
	model.PostgresClusterVersion = types.StringPointerValue(props.PostgresClusterVersion)
	model.IsActive = types.BoolPointerValue(props.IsActive)
	if props.EarliestRecoveryTargetTime != nil {
		model.EarliestRecoveryTargetTime = types.StringValue(props.EarliestRecoveryTargetTime.String())
	}
	if props.LatestRecoveryTargetTime != nil {
		model.LatestRecoveryTargetTime = types.StringValue(props.LatestRecoveryTargetTime.String())
	}
	model.Location = types.StringPointerValue(props.Location)
}
