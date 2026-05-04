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

var _ datasource.DataSourceWithConfigure = (*clustersDataSource)(nil)

type clustersDataSource struct {
	bundle *bundleclient.SdkBundle
}

type clustersDataSourceModel struct {
	Location types.String             `tfsdk:"location"`
	Name     types.String             `tfsdk:"name"`
	Clusters []clusterDataSourceModel `tfsdk:"clusters"`
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
			"clusters": schema.ListNestedAttribute{
				Computed:    true,
				Description: "The list of PostgreSQL v2 clusters.",
				NestedObject: schema.NestedAttributeObject{
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
							Computed:    true,
							Description: "The location of the PostgreSQL cluster.",
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
				},
			},
			"location": schema.StringAttribute{
				Required:    true,
				Description: "The region in which to look up clusters. Available locations: " + pgsqlv2Service.AvailableLocationsString() + ".",
			},
			"name": schema.StringAttribute{
				Optional:    true,
				Description: "Filters clusters by name. Matches cluster names that contain the provided string.",
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

	client, err := d.bundle.NewPgSQLV2Client(ctx, location)
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
