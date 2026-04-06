package pgsqlv2

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-validators/datasourcevalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"
	pgsqlv2 "github.com/ionos-cloud/sdk-go-bundle/products/dbaas/psql/v3"

	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/bundleclient"
	pgsqlv2Service "github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/dbaas/pgsqlv2"
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
	Name              types.String            `tfsdk:"name"`
	Description       types.String            `tfsdk:"description"`
	Version           types.String            `tfsdk:"version"`
	DNSName           types.String            `tfsdk:"dns_name"`
	Location          types.String            `tfsdk:"location"`
	BackupLocation    types.String            `tfsdk:"backup_location"`
	ReplicationMode   types.String            `tfsdk:"replication_mode"`
	ConnectionPooler  types.String            `tfsdk:"connection_pooler"`
	LogsEnabled       types.Bool              `tfsdk:"logs_enabled"`
	MetricsEnabled    types.Bool              `tfsdk:"metrics_enabled"`
	Instances         *instancesModel         `tfsdk:"instances"`
	Connections       *connectionModel        `tfsdk:"connections"`
	MaintenanceWindow *maintenanceWindowModel `tfsdk:"maintenance_window"`
}

// NewClusterDataSource creates a new data source for reading a single PgSQL v2 cluster.
func NewClusterDataSource() datasource.DataSource {
	return &clusterDataSource{}
}

// Metadata returns the metadata for the cluster data source.
func (d *clusterDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_pg_cluster_v2"
}

// Configure configures the data source.
func (d *clusterDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

// ConfigValidators returns the config validators for the data source.
func (d *clusterDataSource) ConfigValidators(_ context.Context) []datasource.ConfigValidator {
	return []datasource.ConfigValidator{
		datasourcevalidator.ExactlyOneOf(
			path.MatchRoot("name"),
			path.MatchRoot("id"),
		),
	}
}

// Schema returns the schema for the cluster data source.
func (d *clusterDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Reads a single PostgreSQL v2 cluster by ID or name.",
		Attributes: map[string]schema.Attribute{
			"backup_location": schema.StringAttribute{
				Computed:    true,
				Description: "The S3 location where the backups are stored.",
			},
			"connection_pooler": schema.StringAttribute{
				Computed:    true,
				Description: "How database connections are managed and reused.",
			},
			"connections": schema.SingleNestedAttribute{
				Computed:    true,
				Description: "Connection information of the PostgreSQL cluster.",
				Attributes: map[string]schema.Attribute{
					"datacenter_id": schema.StringAttribute{
						Computed:    true,
						Description: "The datacenter the cluster is connected to.",
					},
					"lan_id": schema.StringAttribute{
						Computed:    true,
						Description: "The numeric LAN ID the cluster is connected to.",
					},
					"primary_instance_address": schema.StringAttribute{
						Computed:    true,
						Description: "The IP and netmask assigned to the cluster primary instance.",
					},
				},
			},
			"description": schema.StringAttribute{
				Computed:    true,
				Description: "Human-readable description of the cluster.",
			},
			"dns_name": schema.StringAttribute{
				Computed:    true,
				Description: "The DNS name used to access the cluster.",
			},
			"id": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "The ID (UUID) of the cluster.",
			},
			"instances": schema.SingleNestedAttribute{
				Computed:    true,
				Description: "The instance configuration for the PostgreSQL cluster.",
				Attributes: map[string]schema.Attribute{
					"cores": schema.Int32Attribute{
						Computed:    true,
						Description: "The number of CPU cores per instance.",
					},
					"count": schema.Int32Attribute{
						Computed:    true,
						Description: "The total number of instances in the cluster (one primary and n-1 secondary).",
					},
					"ram": schema.Int32Attribute{
						Computed:    true,
						Description: "The amount of memory per instance in gigabytes (GB).",
					},
					"storage_size": schema.Int32Attribute{
						Computed:    true,
						Description: "The amount of storage per instance in gigabytes (GB).",
					},
				},
			},
			"location": schema.StringAttribute{
				Required:    true,
				Description: "The region in which to look up the cluster. Available locations: " + pgsqlv2Service.AvailableLocationsString() + ".",
			},
			"logs_enabled": schema.BoolAttribute{
				Computed:    true,
				Description: "Whether the collection and reporting of logs is enabled for this cluster.",
			},
			"maintenance_window": schema.SingleNestedAttribute{
				Computed:    true,
				Description: "A weekly 4 hour-long window, during which maintenance might occur.",
				Attributes: map[string]schema.Attribute{
					"day_of_the_week": schema.StringAttribute{
						Computed:    true,
						Description: "The name of the week day.",
					},
					"time": schema.StringAttribute{
						Computed:    true,
						Description: "Start of the maintenance window in UTC time.",
					},
				},
			},
			"metrics_enabled": schema.BoolAttribute{
				Computed:    true,
				Description: "Whether the collection and reporting of metrics is enabled for this cluster.",
			},
			"name": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "The name of the PostgreSQL cluster.",
			},
			"replication_mode": schema.StringAttribute{
				Computed:    true,
				Description: "Replication mode across the instances.",
			},
			"version": schema.StringAttribute{
				Computed:    true,
				Description: "The PostgreSQL version of the cluster.",
			},
		},
	}
}

// Read reads the PgSQL v2 cluster data source.
func (d *clusterDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data clusterDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	clusterID := data.ID.ValueString()
	clusterName := data.Name.ValueString()
	location := data.Location.ValueString()

	client, err := d.bundle.NewPgSQLV2Client(location)
	if err != nil {
		resp.Diagnostics.AddError("failed to create PostgreSQL v2 client", err.Error())
		return
	}

	var cluster pgsqlv2.ClusterRead

	if clusterID != "" {
		retrieved, _, err := client.GetCluster(ctx, clusterID)
		if err != nil {
			resp.Diagnostics.AddError("failed to get PostgreSQL v2 cluster", err.Error())
			return
		}
		cluster = retrieved
	}

	if clusterName != "" {
		clusterList, _, err := client.ListClusters(ctx, clusterName)
		if err != nil {
			resp.Diagnostics.AddError(fmt.Sprintf("failed to list PostgreSQL v2 clusters in location: %s", location), err.Error())
			return
		}

		found, diags := findClusterByName(clusterList.Items, clusterName, location)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
		cluster = found
	}

	mapClusterResponseToDataSourceModel(&cluster, &data)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// findClusterByName finds exactly one cluster with the given name from a list.
// Returns an error diagnostic if zero or multiple clusters match.
func findClusterByName(clusters []pgsqlv2.ClusterRead, name, location string) (pgsqlv2.ClusterRead, diag.Diagnostics) {
	var matched []pgsqlv2.ClusterRead
	for _, c := range clusters {
		if c.Properties.Name == name {
			matched = append(matched, c)
		}
	}

	if len(matched) > 1 {
		return pgsqlv2.ClusterRead{}, diag.Diagnostics{
			diag.NewErrorDiagnostic("multiple PostgreSQL v2 clusters found with the same name", "Please search using the cluster ID instead."),
		}
	}
	if len(matched) == 0 {
		return pgsqlv2.ClusterRead{}, diag.Diagnostics{
			diag.NewErrorDiagnostic(
				fmt.Sprintf("no PostgreSQL v2 cluster found with name: %s in location: %s", name, location),
				"Please make sure that the name and location are correct, or search using the cluster ID instead.",
			),
		}
	}
	return matched[0], nil
}

// mapClusterResponseToDataSourceModel maps the API response to the data source model.
func mapClusterResponseToDataSourceModel(cluster *pgsqlv2.ClusterRead, model *clusterDataSourceModel) {
	model.ID = types.StringValue(cluster.Id)

	model.DNSName = types.StringPointerValue(cluster.Metadata.DnsName)

	props := &cluster.Properties

	model.Name = types.StringValue(props.Name)
	model.Description = types.StringPointerValue(props.Description)
	model.Version = types.StringPointerValue(props.Version)
	model.ReplicationMode = types.StringValue(string(props.ReplicationMode))
	model.BackupLocation = types.StringValue(props.BackupLocation)
	model.ConnectionPooler = types.StringPointerValue(props.ConnectionPooler)
	model.LogsEnabled = types.BoolPointerValue(props.LogsEnabled)
	model.MetricsEnabled = types.BoolPointerValue(props.MetricsEnabled)

	model.Instances = &instancesModel{
		Count:       types.Int32Value(props.Instances.Count),
		Cores:       types.Int32Value(props.Instances.Cores),
		RAM:         types.Int32Value(props.Instances.Ram),
		StorageSize: types.Int32Value(props.Instances.StorageSize),
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
}
