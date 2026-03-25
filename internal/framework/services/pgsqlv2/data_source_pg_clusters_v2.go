package pgsqlv2

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/bundleclient"
)

var _ datasource.DataSourceWithConfigure = (*clustersDataSource)(nil)

type clustersDataSource struct {
	bundle *bundleclient.SdkBundle
}

type clustersDataSourceModel struct {
	Location types.String               `tfsdk:"location"`
	Name     types.String               `tfsdk:"name"`
	Clusters []clusterDataSourceModel   `tfsdk:"clusters"`
}

// NewClustersDataSource creates a new data source for listing PgSQL v2 clusters.
func NewClustersDataSource() datasource.DataSource {
	return &clustersDataSource{}
}

// Metadata returns the metadata for the clusters data source.
func (d *clustersDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_pg_clusters_v2"
}

// Configure configures the data source.
func (d *clustersDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

// Schema returns the schema for the clusters data source.
func (d *clustersDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Lists PostgreSQL v2 clusters.",
		Attributes: map[string]schema.Attribute{
			"location": schema.StringAttribute{
				Required:    true,
				Description: "The region in which to look up clusters.",
			},
			"name": schema.StringAttribute{
				Optional:    true,
				Description: "Filters clusters by name. Matches cluster names that contain the provided string.",
			},
			"clusters": schema.ListNestedAttribute{
				Computed:    true,
				Description: "The list of PostgreSQL v2 clusters.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Computed:    true,
							Description: "The ID (UUID) of the cluster.",
						},
						"name": schema.StringAttribute{
							Computed:    true,
							Description: "The name of the PostgreSQL cluster.",
						},
						"description": schema.StringAttribute{
							Computed:    true,
							Description: "Human-readable description of the cluster.",
						},
						"version": schema.StringAttribute{
							Computed:    true,
							Description: "The PostgreSQL version of the cluster.",
						},
						"dns_name": schema.StringAttribute{
							Computed:    true,
							Description: "The DNS name used to access the cluster.",
						},
						"location": schema.StringAttribute{
							Computed:    true,
							Description: "The location of the PostgreSQL cluster.",
						},
						"backup_location": schema.StringAttribute{
							Computed:    true,
							Description: "The S3 location where the backups are stored.",
						},
						"replication_mode": schema.StringAttribute{
							Computed:    true,
							Description: "Replication mode across the instances.",
						},
						"connection_pooler": schema.StringAttribute{
							Computed:    true,
							Description: "How database connections are managed and reused.",
						},
						"logs_enabled": schema.BoolAttribute{
							Computed:    true,
							Description: "Whether the collection and reporting of logs is enabled for this cluster.",
						},
						"metrics_enabled": schema.BoolAttribute{
							Computed:    true,
							Description: "Whether the collection and reporting of metrics is enabled for this cluster.",
						},
						"instances": schema.SingleNestedAttribute{
							Computed:    true,
							Description: "The instance configuration for the PostgreSQL cluster.",
							Attributes: map[string]schema.Attribute{
								"count": schema.Int32Attribute{
									Computed:    true,
									Description: "The total number of instances in the cluster (one primary and n-1 secondary).",
								},
								"cores": schema.Int32Attribute{
									Computed:    true,
									Description: "The number of CPU cores per instance.",
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
						"maintenance_window": schema.SingleNestedAttribute{
							Computed:    true,
							Description: "A weekly 4 hour-long window, during which maintenance might occur.",
							Attributes: map[string]schema.Attribute{
								"time": schema.StringAttribute{
									Computed:    true,
									Description: "Start of the maintenance window in UTC time.",
								},
								"day_of_the_week": schema.StringAttribute{
									Computed:    true,
									Description: "The name of the week day.",
								},
							},
						},
					},
				},
			},
		},
	}
}

// Read reads the PgSQL v2 clusters data source.
func (d *clustersDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data clustersDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	location := data.Location.ValueString()
	filterName := data.Name.ValueString()

	client, err := d.bundle.NewPgSQLV2Client(location)
	if err != nil {
		resp.Diagnostics.AddError("failed to create PostgreSQL v2 client", err.Error())
		return
	}

	clusterList, _, err := client.ListClusters(ctx, filterName)
	if err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("failed to list PostgreSQL v2 clusters in location: %s", location), err.Error())
		return
	}

	var clusters []clusterDataSourceModel
	for _, c := range clusterList.Items {
		var clusterModel clusterDataSourceModel
		mapClusterResponseToDataSourceModel(&c, &clusterModel)
		clusterModel.Location = types.StringValue(location)
		clusters = append(clusters, clusterModel)
	}

	data.Clusters = clusters

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
