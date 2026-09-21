package mariadbv2

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
	"github.com/hashicorp/terraform-plugin-framework/resource/identityschema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	mariadbv3 "github.com/ionos-cloud/sdk-go-bundle/products/dbaas/mariadb/v3"

	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/bundleclient"
	mariadbv2service "github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/dbaas/mariadbv2"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils"
	diagutil "github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils/diags"
)

var (
	_ resource.ResourceWithImportState = (*clusterResource)(nil)
	_ resource.ResourceWithConfigure   = (*clusterResource)(nil)
	_ resource.ResourceWithIdentity    = (*clusterResource)(nil)
)

type clusterIdentityModel struct {
	ID       types.String `tfsdk:"id"`
	Location types.String `tfsdk:"location"`
}

type clusterResource struct {
	bundle *bundleclient.SdkBundle
}

type clusterResourceModel struct {
	Backup            *clusterBackupModel     `tfsdk:"backup"`
	Connections       *connectionModel        `tfsdk:"connections"`
	Credentials       *credentialsModel       `tfsdk:"credentials"`
	Description       types.String            `tfsdk:"description"`
	DNSName           types.String            `tfsdk:"dns_name"`
	ID                types.String            `tfsdk:"id"`
	Instances         *instancesModel         `tfsdk:"instances"`
	Location          types.String            `tfsdk:"location"`
	LogsEnabled       types.Bool              `tfsdk:"logs_enabled"`
	MaintenanceWindow *maintenanceWindowModel `tfsdk:"maintenance_window"`
	MetricsEnabled    types.Bool              `tfsdk:"metrics_enabled"`
	Name              types.String            `tfsdk:"name"`
	RestoreFromBackup *restoreFromBackupModel `tfsdk:"restore_from_backup"`
	Timeouts          timeouts.Value          `tfsdk:"timeouts"`
	Version           types.String            `tfsdk:"version"`
}

type instancesModel struct {
	Cores       types.Int32 `tfsdk:"cores"`
	Count       types.Int32 `tfsdk:"count"`
	RAM         types.Int32 `tfsdk:"ram"`
	StorageSize types.Int32 `tfsdk:"storage_size"`
}

type connectionModel struct {
	DatacenterID           types.String `tfsdk:"datacenter_id"`
	LanID                  types.String `tfsdk:"lan_id"`
	PrimaryInstanceAddress types.String `tfsdk:"primary_instance_address"`
}

type clusterBackupModel struct {
	Location      types.String `tfsdk:"location"`
	RetentionDays types.Int32  `tfsdk:"retention_days"`
}

type maintenanceWindowModel struct {
	DayOfTheWeek types.String `tfsdk:"day_of_the_week"`
	Time         types.String `tfsdk:"time"`
}

type credentialsModel struct {
	Database types.String `tfsdk:"database"`
	Password types.String `tfsdk:"password"`
	Username types.String `tfsdk:"username"`
}

type restoreFromBackupModel struct {
	RecoveryTargetDatetime types.String `tfsdk:"recovery_target_datetime"`
	SourceBackupID         types.String `tfsdk:"source_backup_id"`
}

// NewClusterResource creates a new resource for the MariaDB v2 cluster.
func NewClusterResource() resource.Resource {
	return &clusterResource{}
}

// Metadata returns the resource type name.
func (r *clusterResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_mariadb_cluster_v2"
}

func (r *clusterResource) IdentitySchema(_ context.Context, _ resource.IdentitySchemaRequest, resp *resource.IdentitySchemaResponse) {
	resp.IdentitySchema = identityschema.Schema{
		Attributes: map[string]identityschema.Attribute{
			"id":       identityschema.StringAttribute{RequiredForImport: true},
			"location": identityschema.StringAttribute{RequiredForImport: true},
		},
	}
}

// Schema returns the schema for the MariaDB v2 cluster resource.
func (r *clusterResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages an IONOS CLOUD MariaDB V2 cluster.",
		Attributes: map[string]schema.Attribute{
			"backup": schema.SingleNestedAttribute{
				Required:    true,
				Description: "Configures backup location and retention.",
				Attributes: map[string]schema.Attribute{
					"location": schema.StringAttribute{
						Required:    true,
						Description: "The Object Storage location where the backup will be created. Changing this forces re-creation.",
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.RequiresReplace(),
						},
					},
					"retention_days": schema.Int32Attribute{
						Required:    true,
						Description: "Configures how many days cluster backups are retained.",
					},
				},
			},
			"connections": schema.SingleNestedAttribute{
				Required:    true,
				Description: "Connection information of the MariaDB cluster.",
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
						Description: "The IP address and netmask of the cluster's primary instance, in CIDR notation.",
					},
				},
			},
			"credentials": schema.SingleNestedAttribute{
				Required:    true,
				Description: "Credentials for the initial database user to be created.",
				Attributes: map[string]schema.Attribute{
					"database": schema.StringAttribute{
						Required:    true,
						Description: "The name of the initial database to be created.",
					},
					"password": schema.StringAttribute{
						Required:    true,
						Sensitive:   true,
						Description: "The password for the initial MariaDB user.",
					},
					"username": schema.StringAttribute{
						Required:    true,
						Description: "The username of the initial MariaDB user.",
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
				Description: "Compute and storage configuration for each instance in the cluster.",
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
				Description: "The location of the cluster. Changing this forces re-creation. Available locations: " + mariadbv2service.AvailableLocationsString() + ".",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"logs_enabled": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Allows or disallows the collection and reporting of logs for this cluster's observability.",
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
				Description: "Allows or disallows the collection and reporting of metrics for this cluster's observability.",
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "The name of the MariaDB cluster.",
			},
			"restore_from_backup": schema.SingleNestedAttribute{
				Optional:    true,
				Description: "Restores the cluster from a backup. On creation, set source_backup_id (and optionally recovery_target_datetime). On in-place update, set recovery_target_datetime only.",
				Attributes: map[string]schema.Attribute{
					"recovery_target_datetime": schema.StringAttribute{
						Optional:    true,
						Description: "Providing this value as an ISO 8601 timestamp causes the system to replay the backups up to the specified time.",
					},
					"source_backup_id": schema.StringAttribute{
						Optional:    true,
						Description: "UUID of the backup to restore from. Required for restore on cluster creation; not valid for in-place restore.",
					},
				},
			},
			"version": schema.StringAttribute{
				Required:    true,
				Description: "The MariaDB version for the cluster. Use GET /versions to retrieve the list of supported versions and their available upgrade paths.",
			},
		},
		Blocks: map[string]schema.Block{
			"timeouts": timeouts.Block(ctx, timeouts.Opts{
				Create: true,
				Update: true,
				Delete: true,
			}),
		},
	}
}

// Configure wires the bundle client into the resource.
func (r *clusterResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	clientBundle, ok := req.ProviderData.(*bundleclient.SdkBundle)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *bundleclient.SdkBundle, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	r.bundle = clientBundle
}

// Create creates a new MariaDB v2 cluster.
func (r *clusterResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan clusterResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if plan.RestoreFromBackup != nil && plan.RestoreFromBackup.SourceBackupID.IsNull() {
		resp.Diagnostics.AddError(
			"missing source_backup_id",
			"source_backup_id is required when restore_from_backup is set during cluster creation.",
		)
		return
	}

	location := plan.Location.ValueString()
	client, err := r.bundle.NewMariaDBV2Client(ctx, location)
	if err != nil {
		resp.Diagnostics.AddError("failed to create MariaDB V2 client", err.Error())
		return
	}

	createTimeout, diags := plan.Timeouts.Create(ctx, utils.DefaultTimeout)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	createProps, diags := buildClusterCreateProperties(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	clusterResponse, apiResponse, err := client.CreateCluster(ctx, mariadbv3.ClusterCreate{Properties: createProps})
	if err != nil {
		resp.Diagnostics.AddError("failed to create MariaDB V2 cluster", diagutil.WrapError(err, &diagutil.ErrorContext{
			ResourceName:   plan.Name.ValueString(),
			StatusCode:     apiResponse.SafeStatusCode(),
			AdditionalInfo: map[string]string{"location": location},
		}).Error())
		return
	}

	clusterID := clusterResponse.Id

	err = backoff.Retry(func() error {
		return client.IsClusterReady(ctx, clusterID)
	}, backoff.NewExponentialBackOff(backoff.WithMaxElapsedTime(createTimeout)))
	if err != nil {
		resp.Diagnostics.AddError("error waiting for MariaDB V2 cluster to become available", diagutil.WrapError(err, &diagutil.ErrorContext{
			ResourceID:     clusterID,
			ResourceName:   plan.Name.ValueString(),
			AdditionalInfo: map[string]string{"location": location},
		}).Error())
		return
	}

	cluster, apiResponseGet, err := client.GetCluster(ctx, clusterID)
	if err != nil {
		resp.Diagnostics.AddError("error reading MariaDB V2 cluster after creation", diagutil.WrapError(err, &diagutil.ErrorContext{
			ResourceID:     clusterID,
			ResourceName:   plan.Name.ValueString(),
			StatusCode:     apiResponseGet.SafeStatusCode(),
			AdditionalInfo: map[string]string{"location": location},
		}).Error())
		return
	}

	mapClusterResponseToModel(ctx, &cluster, &plan)

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
	resp.Diagnostics.Append(resp.Identity.Set(ctx, &clusterIdentityModel{ID: plan.ID, Location: plan.Location})...)
}

// Read reads the MariaDB v2 cluster state.
func (r *clusterResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state clusterResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	clusterID := state.ID.ValueString()
	location := state.Location.ValueString()

	client, err := r.bundle.NewMariaDBV2Client(ctx, location)
	if err != nil {
		resp.Diagnostics.AddError("failed to create MariaDB V2 client", err.Error())
		return
	}

	cluster, apiResponse, err := client.GetCluster(ctx, clusterID)
	if err != nil {
		if apiResponse.HttpNotFound() {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("error reading MariaDB V2 cluster", diagutil.WrapError(err, &diagutil.ErrorContext{
			ResourceID:     clusterID,
			ResourceName:   state.Name.ValueString(),
			StatusCode:     apiResponse.SafeStatusCode(),
			AdditionalInfo: map[string]string{"location": location},
		}).Error())
		return
	}

	mapClusterResponseToModel(ctx, &cluster, &state)

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
	resp.Diagnostics.Append(resp.Identity.Set(ctx, &clusterIdentityModel{ID: state.ID, Location: state.Location})...)
}

// Update updates the MariaDB v2 cluster using PUT semantics.
func (r *clusterResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state clusterResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if plan.RestoreFromBackup != nil {
		if plan.RestoreFromBackup.RecoveryTargetDatetime.IsNull() {
			resp.Diagnostics.AddError(
				"missing recovery_target_datetime",
				"recovery_target_datetime is required when restore_from_backup is set during cluster update.",
			)
			return
		}
		if !plan.RestoreFromBackup.SourceBackupID.IsNull() {
			tflog.Warn(ctx, "source_backup_id is set in restore_from_backup during update but has no effect — the restore source is inferred automatically for in-place restore")
		}
	}

	clusterID := plan.ID.ValueString()
	location := plan.Location.ValueString()

	client, err := r.bundle.NewMariaDBV2Client(ctx, location)
	if err != nil {
		resp.Diagnostics.AddError("failed to create MariaDB V2 client", err.Error())
		return
	}

	updateTimeout, diags := plan.Timeouts.Update(ctx, utils.DefaultTimeout)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	updateProps, diags := buildClusterUpdateProperties(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	updateReq := mariadbv3.ClusterEnsure{
		Id:         clusterID,
		Properties: updateProps,
	}

	_, apiResponseUpdate, err := client.UpdateCluster(ctx, updateReq, clusterID)
	if err != nil {
		resp.Diagnostics.AddError("error updating MariaDB V2 cluster", diagutil.WrapError(err, &diagutil.ErrorContext{
			ResourceID:     clusterID,
			ResourceName:   state.Name.ValueString(),
			StatusCode:     apiResponseUpdate.SafeStatusCode(),
			AdditionalInfo: map[string]string{"location": location},
		}).Error())
		return
	}

	err = backoff.Retry(func() error {
		return client.IsClusterReady(ctx, clusterID)
	}, backoff.NewExponentialBackOff(backoff.WithMaxElapsedTime(updateTimeout)))
	if err != nil {
		resp.Diagnostics.AddError("error waiting for MariaDB V2 cluster to become available after update", diagutil.WrapError(err, &diagutil.ErrorContext{
			ResourceID:     clusterID,
			ResourceName:   state.Name.ValueString(),
			AdditionalInfo: map[string]string{"location": location},
		}).Error())
		return
	}

	cluster, apiResponseGetAfterUpdate, err := client.GetCluster(ctx, clusterID)
	if err != nil {
		resp.Diagnostics.AddError("error reading MariaDB V2 cluster after update", diagutil.WrapError(err, &diagutil.ErrorContext{
			ResourceID:     clusterID,
			ResourceName:   state.Name.ValueString(),
			StatusCode:     apiResponseGetAfterUpdate.SafeStatusCode(),
			AdditionalInfo: map[string]string{"location": location},
		}).Error())
		return
	}

	mapClusterResponseToModel(ctx, &cluster, &plan)

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
	resp.Diagnostics.Append(resp.Identity.Set(ctx, &clusterIdentityModel{ID: plan.ID, Location: plan.Location})...)
}

// Delete deletes the MariaDB v2 cluster.
func (r *clusterResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state clusterResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	clusterID := state.ID.ValueString()
	location := state.Location.ValueString()

	client, err := r.bundle.NewMariaDBV2Client(ctx, location)
	if err != nil {
		resp.Diagnostics.AddError("failed to create MariaDB V2 client", err.Error())
		return
	}

	deleteTimeout, diags := state.Timeouts.Delete(ctx, utils.DefaultTimeout)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	apiResponseDelete, err := client.DeleteCluster(ctx, clusterID)
	if err != nil {
		resp.Diagnostics.AddError("error deleting MariaDB V2 cluster", diagutil.WrapError(err, &diagutil.ErrorContext{
			ResourceID:     clusterID,
			ResourceName:   state.Name.ValueString(),
			StatusCode:     apiResponseDelete.SafeStatusCode(),
			AdditionalInfo: map[string]string{"location": location},
		}).Error())
		return
	}

	err = backoff.Retry(func() error {
		return client.IsClusterDeleted(ctx, clusterID)
	}, backoff.NewExponentialBackOff(backoff.WithMaxElapsedTime(deleteTimeout)))
	if err != nil {
		resp.Diagnostics.AddError("error waiting for MariaDB V2 cluster to be deleted", diagutil.WrapError(err, &diagutil.ErrorContext{
			ResourceID:     clusterID,
			ResourceName:   state.Name.ValueString(),
			AdditionalInfo: map[string]string{"location": location},
		}).Error())
		return
	}
}

// ImportState imports a MariaDB v2 cluster. Supports identity-based import and legacy string import using "location:cluster_id" format.
func (r *clusterResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	if req.Identity != nil {
		var id *clusterIdentityModel
		resp.Diagnostics.Append(req.Identity.Get(ctx, &id)...)
		if resp.Diagnostics.HasError() {
			return
		}
		if id != nil {
			resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("location"), id.Location)...)
			resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), id.ID)...)
			return
		}
	}
	parts := strings.Split(req.ID, ":")
	if len(parts) != 2 {
		resp.Diagnostics.AddError(
			"Unexpected Import Identifier",
			fmt.Sprintf("Expected import identifier with format: '<location>:<cluster_id>'. Got: %q", req.ID),
		)
		return
	}
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("location"), parts[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), parts[1])...)
}

func buildClusterCreateProperties(_ context.Context, plan *clusterResourceModel) (mariadbv3.ClusterCreateProperties, diag.Diagnostics) {
	var diagnostics diag.Diagnostics

	props := mariadbv3.ClusterCreateProperties{
		Backup: mariadbv3.ClusterBackup{
			Location:      plan.Backup.Location.ValueString(),
			RetentionDays: plan.Backup.RetentionDays.ValueInt32(),
		},
		Connection: mariadbv3.MariadbClusterConnection{
			DatacenterId:           plan.Connections.DatacenterID.ValueString(),
			LanId:                  plan.Connections.LanID.ValueString(),
			PrimaryInstanceAddress: plan.Connections.PrimaryInstanceAddress.ValueString(),
		},
		Credentials: mariadbv3.MariadbUser{
			Database: plan.Credentials.Database.ValueString(),
			Password: plan.Credentials.Password.ValueString(),
			Username: plan.Credentials.Username.ValueString(),
		},
		Instances: mariadbv3.InstanceConfiguration{
			Cores:       plan.Instances.Cores.ValueInt32(),
			Count:       plan.Instances.Count.ValueInt32(),
			Ram:         plan.Instances.RAM.ValueInt32(),
			StorageSize: plan.Instances.StorageSize.ValueInt32(),
		},
		MaintenanceWindow: mariadbv3.MaintenanceWindow{
			DayOfTheWeek: mariadbv3.DayOfTheWeek(plan.MaintenanceWindow.DayOfTheWeek.ValueString()),
			Time:         plan.MaintenanceWindow.Time.ValueString(),
		},
		Name:    plan.Name.ValueString(),
		Version: plan.Version.ValueString(),
	}

	props.Description = plan.Description.ValueStringPointer()

	if !plan.LogsEnabled.IsUnknown() {
		props.LogsEnabled = plan.LogsEnabled.ValueBoolPointer()
	}
	if !plan.MetricsEnabled.IsUnknown() {
		props.MetricsEnabled = plan.MetricsEnabled.ValueBoolPointer()
	}

	if plan.RestoreFromBackup != nil {
		restore := mariadbv3.MariadbRestoreClusterFromBackup{}
		restore.SourceBackupId = new(plan.RestoreFromBackup.SourceBackupID.ValueString())
		if !plan.RestoreFromBackup.RecoveryTargetDatetime.IsNull() {
			t, err := time.Parse(time.RFC3339, plan.RestoreFromBackup.RecoveryTargetDatetime.ValueString())
			if err != nil {
				diagnostics.AddError("invalid recovery_target_datetime", err.Error())
				return props, diagnostics
			}
			restore.RecoveryTargetDatetime = &mariadbv3.IonosTime{Time: t}
		}
		props.RestoreFromBackup = &restore
	}

	return props, diagnostics
}

func buildClusterUpdateProperties(_ context.Context, plan *clusterResourceModel) (mariadbv3.Cluster, diag.Diagnostics) {
	var diagnostics diag.Diagnostics

	props := mariadbv3.Cluster{
		Backup: mariadbv3.ClusterBackup{
			Location:      plan.Backup.Location.ValueString(),
			RetentionDays: plan.Backup.RetentionDays.ValueInt32(),
		},
		Connection: mariadbv3.MariadbClusterConnection{
			DatacenterId:           plan.Connections.DatacenterID.ValueString(),
			LanId:                  plan.Connections.LanID.ValueString(),
			PrimaryInstanceAddress: plan.Connections.PrimaryInstanceAddress.ValueString(),
		},
		Credentials: &mariadbv3.MariadbUser{
			Database: plan.Credentials.Database.ValueString(),
			Password: plan.Credentials.Password.ValueString(),
			Username: plan.Credentials.Username.ValueString(),
		},
		Instances: mariadbv3.InstanceConfiguration{
			Cores:       plan.Instances.Cores.ValueInt32(),
			Count:       plan.Instances.Count.ValueInt32(),
			Ram:         plan.Instances.RAM.ValueInt32(),
			StorageSize: plan.Instances.StorageSize.ValueInt32(),
		},
		MaintenanceWindow: mariadbv3.MaintenanceWindow{
			DayOfTheWeek: mariadbv3.DayOfTheWeek(plan.MaintenanceWindow.DayOfTheWeek.ValueString()),
			Time:         plan.MaintenanceWindow.Time.ValueString(),
		},
		Name:    plan.Name.ValueString(),
		Version: plan.Version.ValueString(),
	}

	props.Description = plan.Description.ValueStringPointer()

	if !plan.LogsEnabled.IsUnknown() {
		props.LogsEnabled = plan.LogsEnabled.ValueBoolPointer()
	}
	if !plan.MetricsEnabled.IsUnknown() {
		props.MetricsEnabled = plan.MetricsEnabled.ValueBoolPointer()
	}

	if plan.RestoreFromBackup != nil {
		t, err := time.Parse(time.RFC3339, plan.RestoreFromBackup.RecoveryTargetDatetime.ValueString())
		if err != nil {
			diagnostics.AddError("invalid recovery_target_datetime", err.Error())
			return props, diagnostics
		}
		props.RestoreFromBackup = &mariadbv3.MariadbRestoreClusterFromBackup{
			RecoveryTargetDatetime: &mariadbv3.IonosTime{Time: t},
		}
	}

	return props, diagnostics
}

func mapClusterResponseToModel(_ context.Context, cluster *mariadbv3.ClusterRead, model *clusterResourceModel) {
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

	// The API never returns the password; preserve it from state/plan.
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
