package inmemorydbv2

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/bundleclient"
	inmemorydbv2Service "github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/dbaas/inmemorydbv2"
)

var _ datasource.DataSourceWithConfigure = (*snapshotsDataSource)(nil)

type snapshotDataSourceModel struct {
	ID                         types.String  `tfsdk:"id"`
	Location                   types.String  `tfsdk:"location"`
	ClusterID                  types.String  `tfsdk:"cluster_id"`
	DatacenterID               types.String  `tfsdk:"datacenter_id"`
	EarliestRecoveryTargetTime types.String  `tfsdk:"earliest_recovery_target_time"`
	LatestRecoveryTargetTime   types.String  `tfsdk:"latest_recovery_target_time"`
	SnapshotLocation           types.String  `tfsdk:"snapshot_location"`
	ClusterVersion             types.String  `tfsdk:"cluster_version"`
	SnapshotSize               types.Float32 `tfsdk:"snapshot_size"`
	RequiredSizeForRestore     types.Float32 `tfsdk:"required_size_for_restore"`
}

func snapshotDataSourceAttributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"cluster_id":                    schema.StringAttribute{Computed: true, Description: "The ID of the cluster this snapshot belongs to."},
		"cluster_version":               schema.StringAttribute{Computed: true, Description: "The version for the cluster."},
		"datacenter_id":                 schema.StringAttribute{Computed: true, Description: "The ID of the data center where the snapshot was created."},
		"earliest_recovery_target_time": schema.StringAttribute{Computed: true, Description: "The earliest time for which a snapshot is available to restore from."},
		"id":                            schema.StringAttribute{Computed: true, Description: "The ID (UUID) of the snapshot."},
		"latest_recovery_target_time":   schema.StringAttribute{Computed: true, Description: "The most recent time for which a snapshot is available to restore from."},
		"location":                      schema.StringAttribute{Computed: true, Description: "The location of the snapshot."},
		"required_size_for_restore":     schema.Float32Attribute{Computed: true, Description: "The minimum storage size in GB required to restore from this snapshot."},
		"snapshot_location":             schema.StringAttribute{Computed: true, Description: "The Object Storage location where the snapshot is stored."},
		"snapshot_size":                 schema.Float32Attribute{Computed: true, Description: "The size of the snapshot in GB."},
	}
}

type snapshotsDataSource struct {
	bundle *bundleclient.SdkBundle
}

type snapshotsDataSourceModel struct {
	Location  types.String              `tfsdk:"location"`
	ClusterID types.String              `tfsdk:"cluster_id"`
	Items     []snapshotDataSourceModel `tfsdk:"items"`
}

// NewSnapshotsDataSource creates a new data source for listing InMemoryDB v2 snapshots.
func NewSnapshotsDataSource() datasource.DataSource {
	return &snapshotsDataSource{}
}

func (d *snapshotsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_inmemorydb_snapshots_v2"
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
		Description: "Lists InMemoryDB v2 snapshots, with optional filter by cluster ID.",
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
				Description: "The location to query. Available locations: " + inmemorydbv2Service.AvailableLocationsString() + ".",
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
	client, err := d.bundle.NewInMemoryDBV2Client(ctx, location)
	if err != nil {
		resp.Diagnostics.AddError("failed to create InMemoryDB v2 client", err.Error())
		return
	}

	list, _, err := client.ListSnapshots(ctx, data.ClusterID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("failed to list InMemoryDB v2 snapshots in location %s", location), err.Error())
		return
	}

	items := make([]snapshotDataSourceModel, 0, len(list.Items))
	for _, s := range list.Items {
		item := snapshotDataSourceModel{
			ID:       types.StringValue(s.Id),
			Location: data.Location,
		}
		props := &s.Properties
		item.ClusterID = types.StringPointerValue(props.ClusterId)
		item.DatacenterID = types.StringPointerValue(props.DatacenterId)
		if props.EarliestRecoveryTargetTime != nil {
			item.EarliestRecoveryTargetTime = types.StringValue(props.EarliestRecoveryTargetTime.Time.Format(time.RFC3339))
		}
		if props.LatestRecoveryTargetTime != nil {
			if t := props.LatestRecoveryTargetTime.Get(); t != nil {
				item.LatestRecoveryTargetTime = types.StringValue(t.Format(time.RFC3339))
			}
		}
		item.SnapshotLocation = types.StringPointerValue(props.Location)
		item.ClusterVersion = types.StringPointerValue(props.ClusterVersion)
		if props.SnapshotSize != nil {
			item.SnapshotSize = types.Float32Value(*props.SnapshotSize)
		}
		if props.RequiredSizeForRestore != nil {
			item.RequiredSizeForRestore = types.Float32Value(*props.RequiredSizeForRestore)
		}
		items = append(items, item)
	}
	data.Items = items

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
