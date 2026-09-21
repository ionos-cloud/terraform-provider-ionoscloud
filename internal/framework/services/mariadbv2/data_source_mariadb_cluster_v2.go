package mariadbv2

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-validators/datasourcevalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"
	mariadbv3 "github.com/ionos-cloud/sdk-go-bundle/products/dbaas/mariadb/v3"

	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/bundleclient"
	mariadbv2service "github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/dbaas/mariadbv2"
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
	Backup            *clusterBackupModel         `tfsdk:"backup"`
	Connections       *connectionModel            `tfsdk:"connections"`
	Credentials       *credentialsDataSourceModel `tfsdk:"credentials"`
	Description       types.String                `tfsdk:"description"`
	DNSName           types.String                `tfsdk:"dns_name"`
	ID                types.String                `tfsdk:"id"`
	Instances         *instancesModel             `tfsdk:"instances"`
	Location          types.String                `tfsdk:"location"`
	LogsEnabled       types.Bool                  `tfsdk:"logs_enabled"`
	MaintenanceWindow *maintenanceWindowModel     `tfsdk:"maintenance_window"`
	MetricsEnabled    types.Bool                  `tfsdk:"metrics_enabled"`
	Name              types.String                `tfsdk:"name"`
	Version           types.String                `tfsdk:"version"`
}

type credentialsDataSourceModel struct {
	Database types.String `tfsdk:"database"`
	Username types.String `tfsdk:"username"`
}

// NewClusterDataSource creates a new data source for reading a single MariaDB v2 cluster.
func NewClusterDataSource() datasource.DataSource {
	return &clusterDataSource{}
}

func (d *clusterDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_mariadb_cluster_v2"
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
		Description: "Reads a single IONOS CLOUD MariaDB V2 cluster by ID or name.",
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
	client, err := d.bundle.NewMariaDBV2Client(ctx, location)
	if err != nil {
		resp.Diagnostics.AddError("failed to create MariaDB V2 client", err.Error())
		return
	}

	var cluster mariadbv3.ClusterRead

	if id := data.ID.ValueString(); id != "" {
		retrieved, apiResponse, err := client.GetCluster(ctx, id)
		if err != nil {
			resp.Diagnostics.AddError("failed to get MariaDB V2 cluster", diagutil.WrapError(err, &diagutil.ErrorContext{
				ResourceID:     id,
				StatusCode:     apiResponse.SafeStatusCode(),
				AdditionalInfo: map[string]string{"location": location},
			}).Error())
			return
		}
		cluster = retrieved
	} else {
		name := data.Name.ValueString()
		list, apiResponse, err := client.ListClusters(ctx, name)
		if err != nil {
			resp.Diagnostics.AddError("failed to list MariaDB V2 clusters", diagutil.WrapError(err, &diagutil.ErrorContext{
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

	mapClusterResponseToDataSourceModel(ctx, &cluster, &data)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func findClusterByName(clusters []mariadbv3.ClusterRead, name, location string) (mariadbv3.ClusterRead, diag.Diagnostics) {
	var matched []mariadbv3.ClusterRead
	for _, c := range clusters {
		if c.Properties.Name == name {
			matched = append(matched, c)
		}
	}
	switch len(matched) {
	case 1:
		return matched[0], nil
	case 0:
		return mariadbv3.ClusterRead{}, diag.Diagnostics{
			diag.NewErrorDiagnostic(
				fmt.Sprintf("no MariaDB V2 cluster found with name %q in location %s", name, location),
				"Verify the name and location, or search by ID instead.",
			),
		}
	default:
		return mariadbv3.ClusterRead{}, diag.Diagnostics{
			diag.NewErrorDiagnostic(
				fmt.Sprintf("multiple MariaDB V2 clusters found with name %q in location %s", name, location),
				"Use the cluster ID to uniquely identify the cluster.",
			),
		}
	}
}

func mapClusterResponseToDataSourceModel(_ context.Context, cluster *mariadbv3.ClusterRead, model *clusterDataSourceModel) {
	model.ID = types.StringValue(cluster.Id)
	model.DNSName = types.StringPointerValue(cluster.Metadata.DnsName)

	props := &cluster.Properties
	model.Name = types.StringValue(props.Name)
	model.Description = types.StringPointerValue(props.Description)
	model.Version = types.StringValue(props.Version)
	model.LogsEnabled = types.BoolPointerValue(props.LogsEnabled)
	model.MetricsEnabled = types.BoolPointerValue(props.MetricsEnabled)

	model.Instances = &instancesModel{
		Count:       types.Int32Value(props.Instances.Count),
		RAM:         types.Int32Value(props.Instances.Ram),
		Cores:       types.Int32Value(props.Instances.Cores),
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

	model.Backup = &clusterBackupModel{
		Location:      types.StringValue(props.Backup.Location),
		RetentionDays: types.Int32Value(props.Backup.RetentionDays),
	}

	if props.Credentials != nil {
		model.Credentials = &credentialsDataSourceModel{
			Database: types.StringValue(props.Credentials.Database),
			Username: types.StringValue(props.Credentials.Username),
		}
	}
}

// clusterListItemAttributes returns the schema for a cluster item inside the clusters list data source.
// All attributes are Computed since items are fully returned by the API.
func clusterListItemAttributes() map[string]schema.Attribute {
	attrs := clusterDataSourceAttributes()
	attrs["id"] = schema.StringAttribute{Computed: true, Description: "The ID (UUID) of the cluster."}
	attrs["name"] = schema.StringAttribute{Computed: true, Description: "The name of the cluster."}
	attrs["location"] = schema.StringAttribute{Computed: true, Description: "The location of the cluster."}
	return attrs
}

func clusterDataSourceAttributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"backup": schema.SingleNestedAttribute{
			Computed:    true,
			Description: "Backup location and retention configuration.",
			Attributes: map[string]schema.Attribute{
				"location":       schema.StringAttribute{Computed: true, Description: "The Object Storage location where the backup is stored."},
				"retention_days": schema.Int32Attribute{Computed: true, Description: "The number of days cluster backups are retained."},
			},
		},
		"connections": schema.SingleNestedAttribute{
			Computed:    true,
			Description: "Connection information of the MariaDB cluster.",
			Attributes: map[string]schema.Attribute{
				"datacenter_id":            schema.StringAttribute{Computed: true, Description: "The ID of the Virtual Data Center the cluster is connected to."},
				"lan_id":                   schema.StringAttribute{Computed: true, Description: "The numeric LAN ID the cluster is connected to."},
				"primary_instance_address": schema.StringAttribute{Computed: true, Description: "The IP address and netmask of the cluster's primary instance, in CIDR notation."},
			},
		},
		"credentials": schema.SingleNestedAttribute{
			Computed:    true,
			Description: "Credentials for the initial database user.",
			Attributes: map[string]schema.Attribute{
				"database": schema.StringAttribute{Computed: true, Description: "The name of the initial database."},
				"username": schema.StringAttribute{Computed: true, Description: "The username of the initial MariaDB user."},
			},
		},
		"description": schema.StringAttribute{
			Computed:    true,
			Description: "Human-readable description for the cluster.",
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
			Description: "Compute and storage configuration for each instance in the cluster.",
			Attributes: map[string]schema.Attribute{
				"cores":        schema.Int32Attribute{Computed: true, Description: "The number of CPU cores per instance."},
				"count":        schema.Int32Attribute{Computed: true, Description: "The total number of instances in the cluster (one primary and n-1 secondary)."},
				"ram":          schema.Int32Attribute{Computed: true, Description: "The amount of memory per instance in gigabytes (GB)."},
				"storage_size": schema.Int32Attribute{Computed: true, Description: "The amount of storage per instance in gigabytes (GB)."},
			},
		},
		"location": schema.StringAttribute{
			Required:    true,
			Description: "The location of the cluster. Available locations: " + mariadbv2service.AvailableLocationsString() + ".",
		},
		"logs_enabled": schema.BoolAttribute{
			Computed:    true,
			Description: "Whether log collection and reporting is enabled for this cluster's observability.",
		},
		"maintenance_window": schema.SingleNestedAttribute{
			Computed:    true,
			Description: "A weekly 4 hour-long window, during which maintenance might occur.",
			Attributes: map[string]schema.Attribute{
				"day_of_the_week": schema.StringAttribute{Computed: true, Description: "The name of the week day."},
				"time":            schema.StringAttribute{Computed: true, Description: "Start of the maintenance window in UTC time."},
			},
		},
		"metrics_enabled": schema.BoolAttribute{
			Computed:    true,
			Description: "Whether metrics collection and reporting is enabled for this cluster's observability.",
		},
		"name": schema.StringAttribute{
			Optional:    true,
			Computed:    true,
			Description: "The name of the MariaDB cluster.",
		},
		"version": schema.StringAttribute{
			Computed:    true,
			Description: "The MariaDB version for the cluster.",
		},
	}
}
