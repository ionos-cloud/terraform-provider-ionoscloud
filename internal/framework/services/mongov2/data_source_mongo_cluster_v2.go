package mongov2

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-validators/datasourcevalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"
	mongov3 "github.com/ionos-cloud/sdk-go-bundle/products/dbaas/mongo/v3"

	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/bundleclient"
	mongov2service "github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/dbaas/mongov2"
	diagutil "github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils/diags"
)

var (
	_ datasource.DataSourceWithConfigure        = (*clusterDataSource)(nil)
	_ datasource.DataSourceWithConfigValidators = (*clusterDataSource)(nil)
)

type clusterDataSource struct {
	bundle *bundleclient.SdkBundle
}

type clusterDataSourceModel struct {
	BiConnectorEnabled types.Bool                  `tfsdk:"bi_connector_enabled"`
	ConnectionURI      types.String                `tfsdk:"connection_uri"`
	Connections        *connectionModel            `tfsdk:"connections"`
	Credentials        *credentialsDataSourceModel `tfsdk:"credentials"`
	Description        types.String                `tfsdk:"description"`
	Edition            types.String                `tfsdk:"edition"`
	ID                 types.String                `tfsdk:"id"`
	InstanceCount      types.Int32                 `tfsdk:"instance_count"`
	InstanceTemplateID types.String                `tfsdk:"instance_template_id"`
	Instances          *instancesModel             `tfsdk:"instances"`
	Location           types.String                `tfsdk:"location"`
	LogsEnabled        types.Bool                  `tfsdk:"logs_enabled"`
	MaintenanceWindow  *maintenanceWindowModel     `tfsdk:"maintenance_window"`
	MetricsEnabled     types.Bool                  `tfsdk:"metrics_enabled"`
	Name               types.String                `tfsdk:"name"`
	Sharding           *shardingModel              `tfsdk:"sharding"`
	Snapshot           *snapshotConfigModel        `tfsdk:"snapshot"`
	Version            types.String                `tfsdk:"version"`
}

type credentialsDataSourceModel struct {
	Database types.String `tfsdk:"database"`
	Username types.String `tfsdk:"username"`
}

// NewClusterDataSource creates a new data source for reading a single MongoDB v2 cluster.
func NewClusterDataSource() datasource.DataSource {
	return &clusterDataSource{}
}

func (d *clusterDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_mongo_cluster_v2"
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
		Description: "Reads a single IONOS CLOUD MongoDB V2 cluster by ID or name.",
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
	client, err := d.bundle.NewMongoDBV2Client(ctx, location)
	if err != nil {
		resp.Diagnostics.AddError("failed to create MongoDB V2 client", err.Error())
		return
	}

	var cluster mongov3.ClusterRead

	if id := data.ID.ValueString(); id != "" {
		retrieved, apiResponse, err := client.GetCluster(ctx, id)
		if err != nil {
			resp.Diagnostics.AddError("failed to get MongoDB V2 cluster", diagutil.WrapError(err, &diagutil.ErrorContext{
				ResourceID:     id,
				StatusCode:     apiResponse.SafeStatusCode(),
				AdditionalInfo: map[string]string{"location": location},
			}).Error())
			return
		}
		cluster = retrieved
	} else {
		name := data.Name.ValueString()
		list, apiResponse, err := client.ListClusters(ctx, name, "")
		if err != nil {
			resp.Diagnostics.AddError("failed to list MongoDB V2 clusters", diagutil.WrapError(err, &diagutil.ErrorContext{
				StatusCode:     apiResponse.SafeStatusCode(),
				AdditionalInfo: map[string]string{"location": location},
			}).Error())
			return
		}
		found, diags := findClusterByName(list.Items, name, location)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
		cluster = found
	}

	resp.Diagnostics.Append(mapClusterResponseToDataSourceModel(ctx, &cluster, &data)...)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func findClusterByName(clusters []mongov3.ClusterRead, name, location string) (mongov3.ClusterRead, diag.Diagnostics) {
	var matched []mongov3.ClusterRead
	for _, c := range clusters {
		if clusterName(c.Properties) == name {
			matched = append(matched, c)
		}
	}
	switch len(matched) {
	case 1:
		return matched[0], nil
	case 0:
		return mongov3.ClusterRead{}, diag.Diagnostics{
			diag.NewErrorDiagnostic(
				fmt.Sprintf("no MongoDB V2 cluster found with name %q in location %s", name, location),
				"Verify the name and location, or search by ID instead.",
			),
		}
	default:
		return mongov3.ClusterRead{}, diag.Diagnostics{
			diag.NewErrorDiagnostic(
				fmt.Sprintf("multiple MongoDB V2 clusters found with name %q in location %s", name, location),
				"Use the cluster ID to uniquely identify the cluster.",
			),
		}
	}
}

func mapClusterResponseToDataSourceModel(ctx context.Context, cluster *mongov3.ClusterRead, model *clusterDataSourceModel) diag.Diagnostics {
	var diagnostics diag.Diagnostics

	model.ID = types.StringValue(cluster.Id)
	model.ConnectionURI = types.StringPointerValue(cluster.Metadata.ConnectionUri)

	props := &cluster.Properties
	switch {
	case props.EnterpriseCluster != nil:
		ent := props.EnterpriseCluster
		model.Edition = types.StringValue(string(ent.Edition))
		model.Name = types.StringValue(ent.Name)
		model.Description = types.StringPointerValue(ent.Description)
		model.Version = types.StringValue(ent.Version)
		model.LogsEnabled = types.BoolPointerValue(ent.LogsEnabled)
		model.MetricsEnabled = types.BoolPointerValue(ent.MetricsEnabled)
		model.BiConnectorEnabled = types.BoolPointerValue(ent.BiConnectorEnabled)
		model.Instances = &instancesModel{
			Cores:       types.Int32Value(ent.Instances.Cores),
			Count:       types.Int32Value(ent.Instances.Count),
			RAM:         types.Int32Value(ent.Instances.Ram),
			StorageSize: types.Int32Value(ent.Instances.StorageSize),
		}
		if ent.Sharding != nil {
			model.Sharding = &shardingModel{
				Enabled: types.BoolPointerValue(ent.Sharding.Enabled),
				Shards:  types.Int32PointerValue(ent.Sharding.Shards),
			}
		}
		diagnostics.Append(mapCommonDataSourceFields(ctx, model, ent.Connection, ent.MaintenanceWindow, ent.Snapshot, ent.Credentials)...)
	case props.BusinessCluster != nil:
		bus := props.BusinessCluster
		model.Edition = types.StringValue(string(bus.Edition))
		model.Name = types.StringValue(bus.Name)
		model.Description = types.StringPointerValue(bus.Description)
		model.Version = types.StringValue(bus.Version)
		model.LogsEnabled = types.BoolPointerValue(bus.LogsEnabled)
		model.MetricsEnabled = types.BoolPointerValue(bus.MetricsEnabled)
		model.InstanceTemplateID = types.StringValue(bus.InstanceTemplateId)
		model.InstanceCount = types.Int32Value(bus.InstanceCount)
		diagnostics.Append(mapCommonDataSourceFields(ctx, model, bus.Connection, bus.MaintenanceWindow, bus.Snapshot, bus.Credentials)...)
	default:
		diagnostics.AddError("unknown cluster edition", "the API returned a cluster with neither ENTERPRISE nor BUSINESS properties set.")
	}

	return diagnostics
}

func mapCommonDataSourceFields(ctx context.Context, model *clusterDataSourceModel, connection mongov3.MongoClusterConnection, maintenanceWindow mongov3.MaintenanceWindow, snapshot mongov3.ClusterSnapshot, credentials *mongov3.MongoUser) diag.Diagnostics {
	var diagnostics diag.Diagnostics

	addresses, diags := types.ListValueFrom(ctx, types.StringType, connection.InstanceAddresses)
	diagnostics.Append(diags...)
	model.Connections = &connectionModel{
		DatacenterID:      types.StringValue(connection.DatacenterId),
		LanID:             types.StringValue(connection.LanId),
		InstanceAddresses: addresses,
	}

	model.MaintenanceWindow = &maintenanceWindowModel{
		DayOfTheWeek: types.StringValue(string(maintenanceWindow.DayOfTheWeek)),
		Time:         types.StringValue(maintenanceWindow.Time),
	}

	model.Snapshot = &snapshotConfigModel{
		Location:      types.StringValue(snapshot.Location),
		RetentionDays: types.Int32Value(snapshot.RetentionDays),
	}

	if credentials != nil {
		model.Credentials = &credentialsDataSourceModel{
			Database: types.StringValue(credentials.Database),
			Username: types.StringValue(credentials.Username),
		}
	}

	return diagnostics
}

// clusterListItemAttributes returns the schema for a cluster item inside the clusters list data source.
func clusterListItemAttributes() map[string]schema.Attribute {
	attrs := clusterDataSourceAttributes()
	attrs["id"] = schema.StringAttribute{Computed: true, Description: "The ID (UUID) of the cluster."}
	attrs["name"] = schema.StringAttribute{Computed: true, Description: "The name of the cluster."}
	attrs["location"] = schema.StringAttribute{Computed: true, Description: "The location of the cluster."}
	return attrs
}

func clusterDataSourceAttributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"bi_connector_enabled": schema.BoolAttribute{Computed: true, Description: "Whether the MongoDB BI Connector is enabled. Only set for the ENTERPRISE edition."},
		"connection_uri":       schema.StringAttribute{Computed: true, Description: "The connection string used to access the cluster."},
		"connections": schema.SingleNestedAttribute{
			Computed:    true,
			Description: "Connection information of the MongoDB cluster.",
			Attributes: map[string]schema.Attribute{
				"datacenter_id":      schema.StringAttribute{Computed: true, Description: "The ID of the Virtual Data Center the cluster is connected to."},
				"lan_id":             schema.StringAttribute{Computed: true, Description: "The numeric LAN ID the cluster is connected to."},
				"instance_addresses": schema.ListAttribute{Computed: true, ElementType: types.StringType, Description: "The IP addresses (with netmask, in CIDR notation) assigned to the cluster instances."},
			},
		},
		"credentials": schema.SingleNestedAttribute{
			Computed:    true,
			Description: "Credentials for the initial database user.",
			Attributes: map[string]schema.Attribute{
				"database": schema.StringAttribute{Computed: true, Description: "The name of the initial database."},
				"username": schema.StringAttribute{Computed: true, Description: "The username of the initial MongoDB user."},
			},
		},
		"description": schema.StringAttribute{Computed: true, Description: "Human-readable description for the cluster."},
		"edition":     schema.StringAttribute{Computed: true, Description: "The cluster edition (ENTERPRISE or BUSINESS)."},
		"id":          schema.StringAttribute{Optional: true, Computed: true, Description: "The ID (UUID) of the cluster."},
		"instance_count": schema.Int32Attribute{
			Computed:    true,
			Description: "The number of instances in the cluster. Only set for the BUSINESS edition.",
		},
		"instance_template_id": schema.StringAttribute{
			Computed:    true,
			Description: "The ID of the instance template. Only set for the BUSINESS edition.",
		},
		"instances": schema.SingleNestedAttribute{
			Computed:    true,
			Description: "Compute and storage configuration for each instance. Only set for the ENTERPRISE edition.",
			Attributes: map[string]schema.Attribute{
				"cores":        schema.Int32Attribute{Computed: true, Description: "The number of CPU cores per instance."},
				"count":        schema.Int32Attribute{Computed: true, Description: "The total number of instances in the cluster (one primary and n-1 secondary)."},
				"ram":          schema.Int32Attribute{Computed: true, Description: "The amount of memory per instance in gigabytes (GB)."},
				"storage_size": schema.Int32Attribute{Computed: true, Description: "The amount of storage per instance in gigabytes (GB)."},
			},
		},
		"location": schema.StringAttribute{
			Required:    true,
			Description: "The location of the cluster. Available locations: " + mongov2service.AvailableLocationsString() + ".",
		},
		"logs_enabled": schema.BoolAttribute{Computed: true, Description: "Whether log collection and reporting is enabled for this cluster's observability."},
		"maintenance_window": schema.SingleNestedAttribute{
			Computed:    true,
			Description: "A weekly 4 hour-long window, during which maintenance might occur.",
			Attributes: map[string]schema.Attribute{
				"day_of_the_week": schema.StringAttribute{Computed: true, Description: "The name of the week day."},
				"time":            schema.StringAttribute{Computed: true, Description: "Start of the maintenance window in UTC time."},
			},
		},
		"metrics_enabled": schema.BoolAttribute{Computed: true, Description: "Whether metrics collection and reporting is enabled for this cluster's observability."},
		"name":            schema.StringAttribute{Optional: true, Computed: true, Description: "The name of the MongoDB cluster."},
		"sharding": schema.SingleNestedAttribute{
			Computed:    true,
			Description: "Sharding configuration. Only set for the ENTERPRISE edition.",
			Attributes: map[string]schema.Attribute{
				"enabled": schema.BoolAttribute{Computed: true, Description: "Whether sharding is enabled."},
				"shards":  schema.Int32Attribute{Computed: true, Description: "The number of shards in the cluster."},
			},
		},
		"snapshot": schema.SingleNestedAttribute{
			Computed:    true,
			Description: "Snapshot location and retention configuration.",
			Attributes: map[string]schema.Attribute{
				"location":       schema.StringAttribute{Computed: true, Description: "The Object Storage location where the snapshot is stored."},
				"retention_days": schema.Int32Attribute{Computed: true, Description: "The number of days cluster snapshots are retained."},
			},
		},
		"version": schema.StringAttribute{Computed: true, Description: "The MongoDB version for the cluster."},
	}
}
