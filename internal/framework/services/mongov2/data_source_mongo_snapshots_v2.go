package mongov2

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	mongov3 "github.com/ionos-cloud/sdk-go-bundle/products/dbaas/mongo/v3"

	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/bundleclient"
	mongov2service "github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/dbaas/mongov2"
	diagutil "github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils/diags"
)

var _ datasource.DataSourceWithConfigure = (*snapshotsDataSource)(nil)

type snapshotDataSourceModel struct {
	ClusterID                  types.String `tfsdk:"cluster_id"`
	ClusterName                types.String `tfsdk:"cluster_name"`
	EarliestRecoveryTargetTime types.String `tfsdk:"earliest_recovery_target_time"`
	ID                         types.String `tfsdk:"id"`
	LatestRecoveryTargetTime   types.String `tfsdk:"latest_recovery_target_time"`
	Location                   types.String `tfsdk:"location"`
	MongoVersion               types.String `tfsdk:"mongo_version"`
	RequiredStorageSize        types.Int32  `tfsdk:"required_storage_size"`
	TotalSize                  types.Int32  `tfsdk:"total_size"`
}

func snapshotDataSourceAttributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"cluster_id":                    schema.StringAttribute{Computed: true, Description: "The ID of the cluster this snapshot belongs to."},
		"cluster_name":                  schema.StringAttribute{Computed: true, Description: "The name of the cluster this snapshot belongs to."},
		"earliest_recovery_target_time": schema.StringAttribute{Computed: true, Description: "The earliest point in time to which the cluster can be restored from this snapshot."},
		"id":                            schema.StringAttribute{Computed: true, Description: "The ID (UUID) of the snapshot."},
		"latest_recovery_target_time":   schema.StringAttribute{Computed: true, Description: "The latest point in time to which the cluster can be restored. Null if the snapshot can be restored up to the current time."},
		"location":                      schema.StringAttribute{Computed: true, Description: "The Object Storage location where the snapshot is stored."},
		"mongo_version":                 schema.StringAttribute{Computed: true, Description: "The MongoDB version of the cluster at snapshot time."},
		"required_storage_size":         schema.Int32Attribute{Computed: true, Description: "The minimum storage size required to restore this snapshot, in Mebibytes (MiB)."},
		"total_size":                    schema.Int32Attribute{Computed: true, Description: "The total size of the snapshot in Mebibytes (MiB)."},
	}
}

type snapshotsDataSource struct {
	bundle *bundleclient.SdkBundle
}

type snapshotsDataSourceModel struct {
	ClusterID types.String              `tfsdk:"cluster_id"`
	Items     []snapshotDataSourceModel `tfsdk:"items"`
	Location  types.String              `tfsdk:"location"`
}

// NewSnapshotsDataSource creates a new data source for listing MongoDB v2 snapshots.
func NewSnapshotsDataSource() datasource.DataSource {
	return &snapshotsDataSource{}
}

func (d *snapshotsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_mongo_snapshots_v2"
}

func (d *snapshotsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *snapshotsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Lists IONOS CLOUD MongoDB V2 snapshots, with optional filter by cluster ID.",
		Attributes: map[string]schema.Attribute{
			"cluster_id": schema.StringAttribute{
				Optional:    true,
				Description: "Filter snapshots by the cluster they belong to.",
			},
			"items": schema.ListNestedAttribute{
				Computed:    true,
				Description: "The list of snapshots.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: snapshotDataSourceAttributes(),
				},
			},
			"location": schema.StringAttribute{
				Required:    true,
				Description: "The location to query. Available locations: " + mongov2service.AvailableLocationsString() + ".",
			},
		},
	}
}

func (d *snapshotsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data snapshotsDataSourceModel
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

	list, apiResponse, err := client.ListSnapshots(ctx, data.ClusterID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("failed to list MongoDB V2 snapshots", diagutil.WrapError(err, &diagutil.ErrorContext{
			StatusCode:     apiResponse.SafeStatusCode(),
			AdditionalInfo: map[string]string{"location": location},
		}).Error())
		return
	}

	items := make([]snapshotDataSourceModel, 0, len(list.Items))
	for _, s := range list.Items {
		items = append(items, mapSnapshotToModel(s))
	}
	data.Items = items

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func mapSnapshotToModel(s mongov3.SnapshotRead) snapshotDataSourceModel {
	item := snapshotDataSourceModel{
		ID:                  types.StringValue(s.Id),
		ClusterID:           types.StringPointerValue(s.Properties.ClusterId),
		ClusterName:         types.StringPointerValue(s.Properties.ClusterName),
		MongoVersion:        types.StringPointerValue(s.Properties.MongoVersion),
		Location:            types.StringPointerValue(s.Properties.Location),
		TotalSize:           types.Int32PointerValue(s.Properties.TotalSize),
		RequiredStorageSize: types.Int32PointerValue(s.Properties.RequiredStorageSize),
	}
	if s.Properties.EarliestRecoveryTargetTime != nil {
		item.EarliestRecoveryTargetTime = types.StringValue(s.Properties.EarliestRecoveryTargetTime.Format(time.RFC3339))
	}
	if s.Properties.LatestRecoveryTargetTime != nil {
		if t := s.Properties.LatestRecoveryTargetTime.Get(); t != nil {
			item.LatestRecoveryTargetTime = types.StringValue(t.Format(time.RFC3339))
		}
	}
	return item
}
