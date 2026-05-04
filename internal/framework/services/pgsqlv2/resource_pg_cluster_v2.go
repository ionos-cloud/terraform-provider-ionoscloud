package pgsqlv2

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/cenkalti/backoff/v4"
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	pgsqlv2 "github.com/ionos-cloud/sdk-go-bundle/products/dbaas/psql/v3"

	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/bundleclient"
	pgsqlv2Service "github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/dbaas/pgsqlv2"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils"
)

var (
	_ resource.ResourceWithImportState = (*clusterResource)(nil)
	_ resource.ResourceWithConfigure   = (*clusterResource)(nil)
)

type clusterResource struct {
	bundle *bundleclient.SdkBundle
}

type clusterResourceModel struct {
	ID               types.String   `tfsdk:"id"`
	Name             types.String   `tfsdk:"name"`
	Description      types.String   `tfsdk:"description"`
	Version          types.String   `tfsdk:"version"`
	DNSName          types.String   `tfsdk:"dns_name"`
	Location         types.String   `tfsdk:"location"`
	BackupLocation   types.String   `tfsdk:"backup_location"`
	ReplicationMode  types.String   `tfsdk:"replication_mode"`
	ConnectionPooler types.String   `tfsdk:"connection_pooler"`
	LogsEnabled      types.Bool     `tfsdk:"logs_enabled"`
	MetricsEnabled   types.Bool     `tfsdk:"metrics_enabled"`
	Timeouts         timeouts.Value `tfsdk:"timeouts"`

	Instances         *instancesModel         `tfsdk:"instances"`
	Connections       *connectionModel        `tfsdk:"connections"`
	MaintenanceWindow *maintenanceWindowModel `tfsdk:"maintenance_window"`
	Credentials       *credentialsModel       `tfsdk:"credentials"`
	RestoreFromBackup *restoreFromBackupModel `tfsdk:"restore_from_backup"`
}

type instancesModel struct {
	Count       types.Int32 `tfsdk:"count"`
	Cores       types.Int32 `tfsdk:"cores"`
	RAM         types.Int32 `tfsdk:"ram"`
	StorageSize types.Int32 `tfsdk:"storage_size"`
}

type connectionModel struct {
	DatacenterID           types.String `tfsdk:"datacenter_id"`
	LanID                  types.String `tfsdk:"lan_id"`
	PrimaryInstanceAddress types.String `tfsdk:"primary_instance_address"`
}

type maintenanceWindowModel struct {
	Time         types.String `tfsdk:"time"`
	DayOfTheWeek types.String `tfsdk:"day_of_the_week"`
}

type credentialsModel struct {
	Username types.String `tfsdk:"username"`
	Password types.String `tfsdk:"password"`
	Database types.String `tfsdk:"database"`
}

type restoreFromBackupModel struct {
	SourceBackupID         types.String `tfsdk:"source_backup_id"`
	RecoveryTargetDateTime types.String `tfsdk:"recovery_target_datetime"`
}

// NewClusterResource creates a new resource for the PgSQL v2 cluster resource.
func NewClusterResource() resource.Resource {
	return &clusterResource{}
}

// Metadata returns the metadata for the cluster resource.
func (r *clusterResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_pg_cluster_v2"
}

// Schema returns the schema for the cluster resource.
func (r *clusterResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"backup_location": schema.StringAttribute{
				Required:    true,
				Description: "The S3 location where the backups will be created. Supported locations are provided by the backup locations endpoint.",
			},
			"connection_pooler": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Defines how database connections are managed and reused.",
			},
			"connections": schema.SingleNestedAttribute{
				Required:    true,
				Description: "Connection information of the PostgreSQL cluster.",
				Attributes: map[string]schema.Attribute{
					"datacenter_id": schema.StringAttribute{
						Required:    true,
						Description: "The datacenter to connect your instance to.",
					},
					"lan_id": schema.StringAttribute{
						Required:    true,
						Description: "The numeric LAN ID to connect your instance to.",
					},
					"primary_instance_address": schema.StringAttribute{
						Required:    true,
						Description: "The IP and netmask that will be assigned to the cluster primary instance.",
					},
				},
			},
			"credentials": schema.SingleNestedAttribute{
				Required:    true,
				Description: "Credentials for the master database user to be created.",
				Attributes: map[string]schema.Attribute{
					"database": schema.StringAttribute{
						Required:    true,
						Description: "The name of the initial database to be created.",
					},
					"password": schema.StringAttribute{
						Required:    true,
						Sensitive:   true,
						Description: "The password for the master database user.",
					},
					"username": schema.StringAttribute{
						Required:    true,
						Description: "The username of the master database user.",
					},
				},
			},
			"description": schema.StringAttribute{
				Optional:    true,
				Description: "Human-readable description for the cluster.",
			},
			"dns_name": schema.StringAttribute{
				Computed:    true,
				Description: "The DNS name used to access the cluster.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "The ID (UUID) of the cluster.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"instances": schema.SingleNestedAttribute{
				Required:    true,
				Description: "The instance configuration for the PostgreSQL cluster.",
				Attributes: map[string]schema.Attribute{
					"cores": schema.Int32Attribute{
						Required:    true,
						Description: "The number of CPU cores per instance.",
					},
					"count": schema.Int32Attribute{
						Required:    true,
						Description: "The total number of instances in the cluster (one primary and n-1 secondary).",
					},
					"ram": schema.Int32Attribute{
						Required:    true,
						Description: "The amount of memory per instance in gigabytes (GB).",
					},
					"storage_size": schema.Int32Attribute{
						Required:    true,
						Description: "The amount of storage per instance in gigabytes (GB).",
					},
				},
			},
			"location": schema.StringAttribute{
				Required:    true,
				Description: "The location of the PostgreSQL cluster. This is used for routing to the regional API endpoint. Available locations: " + pgsqlv2Service.AvailableLocationsString() + ".",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"logs_enabled": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Enables or disables the collection and reporting of logs for observability of this cluster.",
			},
			"maintenance_window": schema.SingleNestedAttribute{
				Required:    true,
				Description: "A weekly 4 hour-long window, during which maintenance might occur.",
				Attributes: map[string]schema.Attribute{
					"day_of_the_week": schema.StringAttribute{
						Required:    true,
						Description: "The name of the week day.",
					},
					"time": schema.StringAttribute{
						Required:    true,
						Description: "Start of the maintenance window in UTC time.",
					},
				},
			},
			"metrics_enabled": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Enables or disables the collection and reporting of metrics for observability of this cluster.",
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "The name of the PostgreSQL cluster.",
			},
			"replication_mode": schema.StringAttribute{
				Required:    true,
				Description: "Replication mode across the instances.",
			},
			"restore_from_backup": schema.SingleNestedAttribute{
				Optional:    true,
				Description: "Configures the cluster to be initialized with data from an existing backup.",
				Attributes: map[string]schema.Attribute{
					"recovery_target_datetime": schema.StringAttribute{
						Optional:    true,
						Description: "If supplied as ISO 8601 timestamp, the backup will be replayed up until the given timestamp. If empty, the backup will be applied completely.",
					},
					"source_backup_id": schema.StringAttribute{
						Optional:    true,
						Description: "The UUID of the backup to restore data from.",
					},
				},
			},
			"version": schema.StringAttribute{
				Required:    true,
				Description: "The PostgreSQL version of the cluster.",
			},
		},
		Blocks: map[string]schema.Block{
			"timeouts": timeouts.Block(ctx, timeouts.Opts{
				Create: true,
				Read:   true,
				Update: true,
				Delete: true,
			}),
		},
	}
}

// Configure configures the cluster resource.
func (r *clusterResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	clientBundle, ok := req.ProviderData.(*bundleclient.SdkBundle)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *services.SdkBundle, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	r.bundle = clientBundle
}

// Create creates a new PgSQL v2 cluster.
func (r *clusterResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan clusterResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	location := plan.Location.ValueString()
	client, err := r.bundle.NewPgSQLV2Client(ctx, location)
	if err != nil {
		resp.Diagnostics.AddError("failed to create PostgreSQL v2 client", err.Error())
		return
	}

	createTimeout, diags := plan.Timeouts.Create(ctx, utils.DefaultTimeout)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	createProps, diags := buildClusterCreateProperties(&plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	createReq := pgsqlv2.ClusterCreate{
		Properties: createProps,
	}

	clusterResponse, _, err := client.CreateCluster(ctx, createReq)
	if err != nil {
		resp.Diagnostics.AddError("failed to create PostgreSQL v2 cluster", err.Error())
		return
	}

	clusterID := clusterResponse.Id

	// Poll for PROVISIONED state.
	err = backoff.Retry(func() error {
		return client.IsClusterReady(ctx, clusterID)
	}, backoff.NewExponentialBackOff(backoff.WithMaxElapsedTime(createTimeout)))
	if err != nil {
		resp.Diagnostics.AddError("error waiting for PostgreSQL v2 cluster to become available", err.Error())
		return
	}

	// Get all computed fields.
	cluster, _, err := client.GetCluster(ctx, clusterID)
	if err != nil {
		resp.Diagnostics.AddError("error reading PostgreSQL v2 cluster after creation", err.Error())
		return
	}

	mapClusterResponseToModel(&cluster, &plan)

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

// Read reads the PgSQL v2 cluster state.
func (r *clusterResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state clusterResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	clusterID := state.ID.ValueString()
	location := state.Location.ValueString()

	client, err := r.bundle.NewPgSQLV2Client(ctx, location)
	if err != nil {
		resp.Diagnostics.AddError("failed to create PostgreSQL v2 client", err.Error())
		return
	}

	cluster, apiResponse, err := client.GetCluster(ctx, clusterID)
	if err != nil {
		if apiResponse != nil && apiResponse.HttpNotFound() {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("error reading PostgreSQL v2 cluster", fmt.Sprintf("cluster ID: %s, error: %v", clusterID, err))
		return
	}

	mapClusterResponseToModel(&cluster, &state)

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Update updates the PgSQL v2 cluster using PUT semantics (full replacement).
func (r *clusterResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state clusterResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	clusterID := state.ID.ValueString()
	location := plan.Location.ValueString()

	client, err := r.bundle.NewPgSQLV2Client(ctx, location)
	if err != nil {
		resp.Diagnostics.AddError("failed to create PostgreSQL v2 client", err.Error())
		return
	}

	updateTimeout, diags := plan.Timeouts.Update(ctx, utils.DefaultTimeout)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	updateProps, diags := buildClusterUpdateProperties(&plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	updateReq := pgsqlv2.ClusterEnsure{
		Id:         clusterID,
		Properties: updateProps,
	}

	_, _, err = client.UpdateCluster(ctx, updateReq, clusterID)
	if err != nil {
		resp.Diagnostics.AddError("error updating PostgreSQL v2 cluster", fmt.Sprintf("cluster ID: %s, error: %v", clusterID, err))
		return
	}

	// Poll for PROVISIONED state.
	err = backoff.Retry(func() error {
		return client.IsClusterReady(ctx, clusterID)
	}, backoff.NewExponentialBackOff(backoff.WithMaxElapsedTime(updateTimeout)))
	if err != nil {
		resp.Diagnostics.AddError("error waiting for PostgreSQL v2 cluster to become available after update", fmt.Sprintf("cluster ID: %s, error: %v", clusterID, err))
		return
	}

	cluster, _, err := client.GetCluster(ctx, clusterID)
	if err != nil {
		resp.Diagnostics.AddError("error reading PostgreSQL v2 cluster after update", err.Error())
		return
	}

	mapClusterResponseToModel(&cluster, &plan)

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

// Delete deletes the PgSQL v2 cluster.

func (r *clusterResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state clusterResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	clusterID := state.ID.ValueString()
	location := state.Location.ValueString()

	client, err := r.bundle.NewPgSQLV2Client(ctx, location)
	if err != nil {
		resp.Diagnostics.AddError("failed to create PostgreSQL v2 client", err.Error())
		return
	}

	deleteTimeout, diags := state.Timeouts.Delete(ctx, utils.DefaultTimeout)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	_, err = client.DeleteCluster(ctx, clusterID)
	if err != nil {
		resp.Diagnostics.AddError("error deleting PostgreSQL v2 cluster", fmt.Sprintf("cluster ID: %s, error: %v", clusterID, err))
		return
	}

	// Poll until the cluster is deleted (404).
	err = backoff.Retry(func() error {
		return client.IsClusterDeleted(ctx, clusterID)
	}, backoff.NewExponentialBackOff(backoff.WithMaxElapsedTime(deleteTimeout)))
	if err != nil {
		resp.Diagnostics.AddError("error waiting for PostgreSQL v2 cluster to be deleted", fmt.Sprintf("cluster ID: %s, error: %v", clusterID, err))
		return
	}
}

// ImportState imports a PgSQL v2 cluster using the format "location:cluster_id".
func (r *clusterResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	parts := strings.Split(req.ID, ":")
	if len(parts) != 2 {
		resp.Diagnostics.AddError(
			"Unexpected Import Identifier",
			fmt.Sprintf("Expected import identifier with format: '<location>:<cluster_id>'. Got: %q", req.ID),
		)
		return
	}
	location := parts[0]
	clusterID := parts[1]

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("location"), location)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), clusterID)...)
}

// buildClusterCreateProperties constructs the cluster create request from the plan model.
func buildClusterCreateProperties(plan *clusterResourceModel) (pgsqlv2.ClusterCreateProperties, diag.Diagnostics) {
	var diagnostics diag.Diagnostics
	props := pgsqlv2.ClusterCreateProperties{
		Name:            plan.Name.ValueString(),
		ReplicationMode: pgsqlv2.PostgresClusterReplicationMode(plan.ReplicationMode.ValueString()),
		BackupLocation:  plan.BackupLocation.ValueString(),
	}

	props.Version = plan.Version.ValueStringPointer()

	props.Instances = pgsqlv2.InstanceConfiguration{
		Count:       plan.Instances.Count.ValueInt32(),
		Cores:       plan.Instances.Cores.ValueInt32(),
		Ram:         plan.Instances.RAM.ValueInt32(),
		StorageSize: plan.Instances.StorageSize.ValueInt32(),
	}

	props.Connection = pgsqlv2.PostgresClusterConnection{
		DatacenterId:           plan.Connections.DatacenterID.ValueString(),
		LanId:                  plan.Connections.LanID.ValueString(),
		PrimaryInstanceAddress: plan.Connections.PrimaryInstanceAddress.ValueString(),
	}

	props.MaintenanceWindow = pgsqlv2.MaintenanceWindow{
		Time:         plan.MaintenanceWindow.Time.ValueString(),
		DayOfTheWeek: pgsqlv2.DayOfTheWeek(plan.MaintenanceWindow.DayOfTheWeek.ValueString()),
	}

	props.Credentials = pgsqlv2.PostgresUser{
		Username: plan.Credentials.Username.ValueString(),
		Password: plan.Credentials.Password.ValueString(),
		Database: plan.Credentials.Database.ValueString(),
	}

	props.Description = plan.Description.ValueStringPointer()

	if !plan.ConnectionPooler.IsUnknown() {
		props.ConnectionPooler = plan.ConnectionPooler.ValueStringPointer()
	}

	if !plan.LogsEnabled.IsUnknown() {
		props.LogsEnabled = plan.LogsEnabled.ValueBoolPointer()
	}

	if !plan.MetricsEnabled.IsUnknown() {
		props.MetricsEnabled = plan.MetricsEnabled.ValueBoolPointer()
	}

	if plan.RestoreFromBackup != nil {
		restore := &pgsqlv2.PostgresClusterFromBackup{}
		if !plan.RestoreFromBackup.SourceBackupID.IsNull() {
			restore.SourceBackupId = plan.RestoreFromBackup.SourceBackupID.ValueString()
		}
		if !plan.RestoreFromBackup.RecoveryTargetDateTime.IsNull() {
			t, err := time.Parse(time.RFC3339, plan.RestoreFromBackup.RecoveryTargetDateTime.ValueString())
			if err != nil {
				diagnostics.AddError("invalid recovery_target_datetime",
					fmt.Sprintf("expected RFC3339 format (e.g. 2020-12-10T13:37:50+01:00), got %q, error: %v", plan.RestoreFromBackup.RecoveryTargetDateTime.ValueString(), err))
				return props, diagnostics
			}
			restore.RecoveryTargetDatetime = &pgsqlv2.IonosTime{Time: t}
		}
		props.RestoreFromBackup = restore
	}

	return props, diagnostics
}

// buildClusterUpdateProperties constructs the Cluster properties for PUT update.
func buildClusterUpdateProperties(plan *clusterResourceModel) (pgsqlv2.Cluster, diag.Diagnostics) {
	var diagnostics diag.Diagnostics
	props := pgsqlv2.Cluster{
		Name:            plan.Name.ValueString(),
		ReplicationMode: pgsqlv2.PostgresClusterReplicationMode(plan.ReplicationMode.ValueString()),
		BackupLocation:  plan.BackupLocation.ValueString(),
	}

	props.Version = plan.Version.ValueStringPointer()

	props.Instances = pgsqlv2.InstanceConfiguration{
		Count:       plan.Instances.Count.ValueInt32(),
		Cores:       plan.Instances.Cores.ValueInt32(),
		Ram:         plan.Instances.RAM.ValueInt32(),
		StorageSize: plan.Instances.StorageSize.ValueInt32(),
	}

	props.Connection = pgsqlv2.PostgresClusterConnection{
		DatacenterId:           plan.Connections.DatacenterID.ValueString(),
		LanId:                  plan.Connections.LanID.ValueString(),
		PrimaryInstanceAddress: plan.Connections.PrimaryInstanceAddress.ValueString(),
	}

	props.MaintenanceWindow = pgsqlv2.MaintenanceWindow{
		Time:         plan.MaintenanceWindow.Time.ValueString(),
		DayOfTheWeek: pgsqlv2.DayOfTheWeek(plan.MaintenanceWindow.DayOfTheWeek.ValueString()),
	}

	props.Credentials = &pgsqlv2.PostgresUser{
		Username: plan.Credentials.Username.ValueString(),
		Password: plan.Credentials.Password.ValueString(),
		Database: plan.Credentials.Database.ValueString(),
	}

	props.Description = plan.Description.ValueStringPointer()

	if !plan.ConnectionPooler.IsUnknown() {
		props.ConnectionPooler = plan.ConnectionPooler.ValueStringPointer()
	}

	if !plan.LogsEnabled.IsUnknown() {
		props.LogsEnabled = plan.LogsEnabled.ValueBoolPointer()
	}

	if !plan.MetricsEnabled.IsUnknown() {
		props.MetricsEnabled = plan.MetricsEnabled.ValueBoolPointer()
	}

	if plan.RestoreFromBackup != nil {
		restore := &pgsqlv2.PostgresClusterFromBackup{}
		if !plan.RestoreFromBackup.SourceBackupID.IsNull() {
			restore.SourceBackupId = plan.RestoreFromBackup.SourceBackupID.ValueString()
		}
		if !plan.RestoreFromBackup.RecoveryTargetDateTime.IsNull() {
			t, err := time.Parse(time.RFC3339, plan.RestoreFromBackup.RecoveryTargetDateTime.ValueString())
			if err != nil {
				diagnostics.AddError("invalid recovery_target_datetime",
					fmt.Sprintf("expected RFC3339 format (e.g. 2020-12-10T13:37:50+01:00), got %q, error: %v", plan.RestoreFromBackup.RecoveryTargetDateTime.ValueString(), err))
				return props, diagnostics
			}
			restore.RecoveryTargetDatetime = &pgsqlv2.IonosTime{Time: t}
		}
		props.RestoreFromBackup = restore
	}

	return props, diagnostics
}

// mapClusterResponseToModel maps API response fields to the Terraform model.
func mapClusterResponseToModel(cluster *pgsqlv2.ClusterRead, model *clusterResourceModel) {
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

	// Credentials block: the API returns username and database but not the password.
	// Preserve the password from the existing model (state/plan) since the API never returns it.
	if props.Credentials != nil {
		var existingPassword types.String
		if model.Credentials != nil {
			existingPassword = model.Credentials.Password
		}
		model.Credentials = &credentialsModel{
			Username: types.StringValue(props.Credentials.Username),
			Password: existingPassword,
			Database: types.StringValue(props.Credentials.Database),
		}
	}
}
