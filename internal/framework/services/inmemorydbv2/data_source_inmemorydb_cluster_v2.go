package inmemorydbv2

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-validators/datasourcevalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"
	inmemorydbv3 "github.com/ionos-cloud/sdk-go-bundle/products/dbaas/inmemorydb/v3"

	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/bundleclient"
	inmemorydbv2Service "github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/dbaas/inmemorydbv2"
)

var (
	_ datasource.DataSourceWithConfigure        = (*clusterDataSource)(nil)
	_ datasource.DataSourceWithConfigValidators = (*clusterDataSource)(nil)
)

type clusterDataSource struct {
	bundle *bundleclient.SdkBundle
}

type clusterDataSourceModel struct {
	ID                types.String            `tfsdk:"id"`
	Location          types.String            `tfsdk:"location"`
	Name              types.String            `tfsdk:"name"`
	Description       types.String            `tfsdk:"description"`
	Version           types.String            `tfsdk:"version"`
	PersistenceMode   types.String            `tfsdk:"persistence_mode"`
	EvictionPolicy    types.String            `tfsdk:"eviction_policy"`
	LogsEnabled       types.Bool              `tfsdk:"logs_enabled"`
	MetricsEnabled    types.Bool              `tfsdk:"metrics_enabled"`
	DNSName           types.String            `tfsdk:"dns_name"`
	Instances         *instancesModel         `tfsdk:"instances"`
	Connections       *connectionModel        `tfsdk:"connections"`
	Snapshot          *snapshotConfigModel    `tfsdk:"snapshot"`
	MaintenanceWindow *maintenanceWindowModel `tfsdk:"maintenance_window"`
}

// NewClusterDataSource creates a new data source for reading a single InMemoryDB v2 cluster.
func NewClusterDataSource() datasource.DataSource {
	return &clusterDataSource{}
}

func (d *clusterDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_inmemorydb_cluster_v2"
}

func (d *clusterDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *clusterDataSource) ConfigValidators(_ context.Context) []datasource.ConfigValidator {
	return []datasource.ConfigValidator{
		datasourcevalidator.ExactlyOneOf(
			path.MatchRoot("id"),
			path.MatchRoot("name"),
		),
	}
}

func (d *clusterDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Reads a single InMemoryDB v2 cluster by ID or name.",
		Attributes:  clusterDataSourceAttributes(),
	}
}

func (d *clusterDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data clusterDataSourceModel
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

	var cluster inmemorydbv3.ClusterRead

	if id := data.ID.ValueString(); id != "" {
		retrieved, _, err := client.GetCluster(ctx, id)
		if err != nil {
			resp.Diagnostics.AddError("failed to get InMemoryDB v2 cluster", err.Error())
			return
		}
		cluster = retrieved
	} else {
		name := data.Name.ValueString()
		// Partial match, case-insensitive.
		list, _, err := client.ListClusters(ctx, name)
		if err != nil {
			resp.Diagnostics.AddError(fmt.Sprintf("failed to list InMemoryDB v2 clusters in location %s", location), err.Error())
			return
		}
		// Exact match.
		found, diags := findClusterByName(list.Items, name, location)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
		cluster = found
	}

	resp.Diagnostics.Append(mapClusterResponseToDataSourceModel(ctx, &cluster, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func findClusterByName(clusters []inmemorydbv3.ClusterRead, name, location string) (inmemorydbv3.ClusterRead, diag.Diagnostics) {
	var matched []inmemorydbv3.ClusterRead
	for _, c := range clusters {
		if c.Properties.Name == name {
			matched = append(matched, c)
		}
	}
	switch len(matched) {
	case 1:
		return matched[0], nil
	case 0:
		return inmemorydbv3.ClusterRead{}, diag.Diagnostics{
			diag.NewErrorDiagnostic(
				fmt.Sprintf("no InMemoryDB v2 cluster found with name %q in location %s", name, location),
				"Verify the name and location, or search by ID instead.",
			),
		}
	default:
		return inmemorydbv3.ClusterRead{}, diag.Diagnostics{
			diag.NewErrorDiagnostic(
				fmt.Sprintf("multiple InMemoryDB v2 clusters found with name %q in location %s", name, location),
				"Use the cluster ID to uniquely identify the cluster.",
			),
		}
	}
}

func mapClusterResponseToDataSourceModel(ctx context.Context, cluster *inmemorydbv3.ClusterRead, model *clusterDataSourceModel) diag.Diagnostics {
	var diagnostics diag.Diagnostics
	model.ID = types.StringValue(cluster.Id)
	model.DNSName = types.StringPointerValue(cluster.Metadata.DnsName)

	props := &cluster.Properties
	model.Name = types.StringValue(props.Name)
	model.Description = types.StringPointerValue(props.Description)
	model.Version = types.StringValue(props.Version)
	model.PersistenceMode = types.StringValue(string(props.PersistenceMode))
	model.EvictionPolicy = types.StringValue(string(props.EvictionPolicy))
	model.LogsEnabled = types.BoolPointerValue(props.LogsEnabled)
	model.MetricsEnabled = types.BoolPointerValue(props.MetricsEnabled)

	model.Instances = &instancesModel{
		Count: types.Int32Value(props.Instances.Count),
		Cores: types.Int32Value(props.Instances.Cores),
		RAM:   types.Int32Value(props.Instances.Ram),
	}

	model.Connections = &connectionModel{
		DatacenterID:           types.StringValue(props.Connection.DatacenterId),
		LanID:                  types.StringValue(props.Connection.LanId),
		PrimaryInstanceAddress: types.StringValue(props.Connection.PrimaryInstanceAddress),
	}

	model.MaintenanceWindow = &maintenanceWindowModel{
		Time:         types.StringValue(props.MaintenanceWindow.Time),
		DayOfTheWeek: types.StringValue(string(props.MaintenanceWindow.DayOfTheWeek)),
	}

	snapshotModel := &snapshotConfigModel{
		Location:      types.StringValue(props.Snapshot.Location),
		RetentionDays: types.Int32Value(props.Snapshot.RetentionDays),
	}
	hours := make([]types.Int32, len(props.Snapshot.SnapshotHours))
	for i, h := range props.Snapshot.SnapshotHours {
		hours[i] = types.Int32Value(h)
	}
	listVal, diags := types.ListValueFrom(ctx, types.Int32Type, hours)
	diagnostics.Append(diags...)
	if diagnostics.HasError() {
		return diagnostics
	}
	snapshotModel.SnapshotHours = listVal
	model.Snapshot = snapshotModel
	return diagnostics
}

func clusterDataSourceAttributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"connections": schema.SingleNestedAttribute{
			Computed:    true,
			Description: "Network connection configuration for the cluster.",
			Attributes: map[string]schema.Attribute{
				"datacenter_id":            schema.StringAttribute{Computed: true, Description: "The Virtual Data Center ID."},
				"lan_id":                   schema.StringAttribute{Computed: true, Description: "The numeric LAN ID."},
				"primary_instance_address": schema.StringAttribute{Computed: true, Description: "The IP address and subnet mask in CIDR notation."},
			},
		},
		"description": schema.StringAttribute{
			Computed:    true,
			Description: "Human-readable description of the cluster.",
		},
		"dns_name": schema.StringAttribute{
			Computed:    true,
			Description: "The DNS name used to connect to the cluster's primary instance.",
		},
		"eviction_policy": schema.StringAttribute{
			Computed:    true,
			Description: "The key eviction strategy.",
		},
		"id": schema.StringAttribute{
			Optional:    true,
			Computed:    true,
			Description: "The ID (UUID) of the cluster.",
		},
		"instances": schema.SingleNestedAttribute{
			Computed:    true,
			Description: "Compute configuration for each instance. Storage size is automatically derived from RAM and persistence mode.",
			Attributes: map[string]schema.Attribute{
				"cores": schema.Int32Attribute{Computed: true, Description: "CPU cores per instance."},
				"count": schema.Int32Attribute{Computed: true, Description: "Number of instances."},
				"ram":   schema.Int32Attribute{Computed: true, Description: "RAM per instance in GB."},
			},
		},
		"location": schema.StringAttribute{
			Required:    true,
			Description: "The location of the cluster. Available locations: " + inmemorydbv2Service.AvailableLocationsString() + ".",
		},
		"logs_enabled": schema.BoolAttribute{
			Computed:    true,
			Description: "Whether log collection is enabled.",
		},
		"maintenance_window": schema.SingleNestedAttribute{
			Computed:    true,
			Description: "A weekly 4-hour maintenance window.",
			Attributes: map[string]schema.Attribute{
				"day_of_the_week": schema.StringAttribute{Computed: true, Description: "Day of the week."},
				"time":            schema.StringAttribute{Computed: true, Description: "Start time in UTC (HH:MM:SS)."},
			},
		},
		"metrics_enabled": schema.BoolAttribute{
			Computed:    true,
			Description: "Whether metrics collection is enabled.",
		},
		"name": schema.StringAttribute{
			Optional:    true,
			Computed:    true,
			Description: "The name of the cluster.",
		},
		"persistence_mode": schema.StringAttribute{
			Computed:    true,
			Description: "The data persistence mode.",
		},
		"snapshot": schema.SingleNestedAttribute{
			Computed:    true,
			Description: "Snapshot storage and retention configuration.",
			Attributes: map[string]schema.Attribute{
				"location":       schema.StringAttribute{Computed: true, Description: "Object Storage location for snapshots."},
				"retention_days": schema.Int32Attribute{Computed: true, Description: "Number of days snapshots are retained."},
				"snapshot_hours": schema.ListAttribute{Computed: true, ElementType: types.Int32Type, Description: "UTC hours at which snapshots are taken."},
			},
		},
		"version": schema.StringAttribute{
			Computed:    true,
			Description: "The In-Memory DB version.",
		},
	}
}
