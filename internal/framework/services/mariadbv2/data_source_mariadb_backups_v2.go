package mariadbv2

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	mariadbv3 "github.com/ionos-cloud/sdk-go-bundle/products/dbaas/mariadb/v3"

	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/bundleclient"
	mariadbv2service "github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/dbaas/mariadbv2"
	diagutil "github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils/diags"
)

var _ datasource.DataSourceWithConfigure = (*backupsDataSource)(nil)

type backupDataSourceModel struct {
	ClusterID                  types.String `tfsdk:"cluster_id"`
	ClusterName                types.String `tfsdk:"cluster_name"`
	EarliestRecoveryTargetTime types.String `tfsdk:"earliest_recovery_target_time"`
	ID                         types.String `tfsdk:"id"`
	LatestRecoveryTargetTime   types.String `tfsdk:"latest_recovery_target_time"`
	Location                   types.String `tfsdk:"location"`
	MariadbClusterVersion      types.String `tfsdk:"mariadb_cluster_version"`
}

func backupDataSourceAttributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"cluster_id":                    schema.StringAttribute{Computed: true, Description: "The ID of the cluster this backup belongs to."},
		"cluster_name":                  schema.StringAttribute{Computed: true, Description: "The name of the cluster this backup belongs to."},
		"earliest_recovery_target_time": schema.StringAttribute{Computed: true, Description: "The earliest point in time to which the cluster can be restored from this backup."},
		"id":                            schema.StringAttribute{Computed: true, Description: "The ID (UUID) of the backup."},
		"latest_recovery_target_time":   schema.StringAttribute{Computed: true, Description: "The latest point in time to which the cluster can be restored. Null if the backup can be restored up to the current time."},
		"location":                      schema.StringAttribute{Computed: true, Description: "The Object Storage location where the backup is stored."},
		"mariadb_cluster_version":       schema.StringAttribute{Computed: true, Description: "The MariaDB version of the cluster at backup time."},
	}
}

type backupsDataSource struct {
	bundle *bundleclient.SdkBundle
}

type backupsDataSourceModel struct {
	ClusterID types.String            `tfsdk:"cluster_id"`
	Items     []backupDataSourceModel `tfsdk:"items"`
	Location  types.String            `tfsdk:"location"`
}

// NewBackupsDataSource creates a new data source for listing MariaDB v2 backups.
func NewBackupsDataSource() datasource.DataSource {
	return &backupsDataSource{}
}

func (d *backupsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_mariadb_backups_v2"
}

func (d *backupsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *backupsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Lists IONOS CLOUD MariaDB V2 backups, with optional filter by cluster ID.",
		Attributes: map[string]schema.Attribute{
			"cluster_id": schema.StringAttribute{
				Optional:    true,
				Description: "Filter backups by the cluster they belong to.",
			},
			"items": schema.ListNestedAttribute{
				Computed:    true,
				Description: "The list of backups.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: backupDataSourceAttributes(),
				},
			},
			"location": schema.StringAttribute{
				Required:    true,
				Description: "The location to query. Available locations: " + mariadbv2service.AvailableLocationsString() + ".",
			},
		},
	}
}

func (d *backupsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data backupsDataSourceModel
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

	list, apiResponse, err := client.ListBackups(ctx, data.ClusterID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("failed to list MariaDB V2 backups", diagutil.WrapError(err, &diagutil.ErrorContext{
			StatusCode:     apiResponse.SafeStatusCode(),
			AdditionalInfo: map[string]string{"location": location},
		}).Error())
		return
	}

	items := make([]backupDataSourceModel, 0, len(list.Items))
	for _, b := range list.Items {
		items = append(items, mapBackupToModel(b))
	}
	data.Items = items

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func mapBackupToModel(b mariadbv3.BackupRead) backupDataSourceModel {
	item := backupDataSourceModel{
		ID:                    types.StringValue(b.Id),
		ClusterID:             types.StringPointerValue(b.Properties.ClusterId),
		ClusterName:           types.StringPointerValue(b.Properties.ClusterName),
		MariadbClusterVersion: types.StringPointerValue(b.Properties.MariadbClusterVersion),
		Location:              types.StringPointerValue(b.Properties.Location),
	}
	if b.Properties.EarliestRecoveryTargetTime != nil {
		item.EarliestRecoveryTargetTime = types.StringValue(b.Properties.EarliestRecoveryTargetTime.Format(time.RFC3339))
	}
	if b.Properties.LatestRecoveryTargetTime != nil {
		if t := b.Properties.LatestRecoveryTargetTime.Get(); t != nil {
			item.LatestRecoveryTargetTime = types.StringValue(t.Format(time.RFC3339))
		}
	}
	return item
}
