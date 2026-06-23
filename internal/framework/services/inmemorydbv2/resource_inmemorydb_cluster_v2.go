package inmemorydbv2

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
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	inmemorydbv3 "github.com/ionos-cloud/sdk-go-bundle/products/dbaas/inmemorydb/v3"

	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/bundleclient"
	inmemorydbv2Service "github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/dbaas/inmemorydbv2"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils"
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
	ID              types.String   `tfsdk:"id"`
	Location        types.String   `tfsdk:"location"`
	Name            types.String   `tfsdk:"name"`
	Description     types.String   `tfsdk:"description"`
	Version         types.String   `tfsdk:"version"`
	PersistenceMode types.String   `tfsdk:"persistence_mode"`
	EvictionPolicy  types.String   `tfsdk:"eviction_policy"`
	LogsEnabled     types.Bool     `tfsdk:"logs_enabled"`
	MetricsEnabled  types.Bool     `tfsdk:"metrics_enabled"`
	DNSName         types.String   `tfsdk:"dns_name"`
	Timeouts        timeouts.Value `tfsdk:"timeouts"`

	Instances           *instancesModel           `tfsdk:"instances"`
	Connections         *connectionModel          `tfsdk:"connections"`
	Snapshot            *snapshotConfigModel      `tfsdk:"snapshot"`
	MaintenanceWindow   *maintenanceWindowModel   `tfsdk:"maintenance_window"`
	Credentials         *credentialsModel         `tfsdk:"credentials"`
	RestoreFromSnapshot *restoreFromSnapshotModel `tfsdk:"restore_from_snapshot"`
}

type instancesModel struct {
	Count types.Int32 `tfsdk:"count"`
	Cores types.Int32 `tfsdk:"cores"`
	RAM   types.Int32 `tfsdk:"ram"`
}

type connectionModel struct {
	DatacenterID           types.String `tfsdk:"datacenter_id"`
	LanID                  types.String `tfsdk:"lan_id"`
	PrimaryInstanceAddress types.String `tfsdk:"primary_instance_address"`
}

type snapshotConfigModel struct {
	Location      types.String `tfsdk:"location"`
	RetentionDays types.Int32  `tfsdk:"retention_days"`
	SnapshotHours types.List   `tfsdk:"snapshot_hours"`
}

type maintenanceWindowModel struct {
	Time         types.String `tfsdk:"time"`
	DayOfTheWeek types.String `tfsdk:"day_of_the_week"`
}

type credentialsModel struct {
	Username types.String   `tfsdk:"username"`
	Password *passwordModel `tfsdk:"password"`
}

type passwordModel struct {
	Algorithm types.String `tfsdk:"algorithm"`
	Hash      types.String `tfsdk:"hash"`
}

type restoreFromSnapshotModel struct {
	SourceSnapshotID       types.String `tfsdk:"source_snapshot_id"`
	RecoveryTargetDatetime types.String `tfsdk:"recovery_target_datetime"`
}

// NewClusterResource creates a new resource for the InMemoryDB v2 cluster.
func NewClusterResource() resource.Resource {
	return &clusterResource{}
}

// Metadata returns the resource type name.
func (r *clusterResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_inmemorydb_cluster_v2"
}

func (r *clusterResource) IdentitySchema(_ context.Context, _ resource.IdentitySchemaRequest, resp *resource.IdentitySchemaResponse) {
	resp.IdentitySchema = identityschema.Schema{
		Attributes: map[string]identityschema.Attribute{
			"id":       identityschema.StringAttribute{RequiredForImport: true},
			"location": identityschema.StringAttribute{RequiredForImport: true},
		},
	}
}

// Schema returns the schema for the InMemoryDB v2 cluster resource.
func (r *clusterResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages an IONOS Cloud In-Memory DB v2 cluster.",
		Attributes: map[string]schema.Attribute{
			"connections": schema.SingleNestedAttribute{
				Required:    true,
				Description: "Network connection configuration for the cluster.",
				Attributes: map[string]schema.Attribute{
					"datacenter_id": schema.StringAttribute{
						Required:    true,
						Description: "The ID of the Virtual Data Center to connect the cluster to.",
					},
					"lan_id": schema.StringAttribute{
						Required:    true,
						Description: "The numeric LAN ID within the data center.",
					},
					"primary_instance_address": schema.StringAttribute{
						Required:    true,
						Description: "The IP address and subnet mask assigned to the primary instance in CIDR notation.",
					},
				},
			},
			"credentials": schema.SingleNestedAttribute{
				Required:    true,
				Description: "Credentials for the user with access to the cluster.",
				Attributes: map[string]schema.Attribute{
					"username": schema.StringAttribute{
						Required:    true,
						Description: "The username for the In-Memory DB user.",
					},
					"password": schema.SingleNestedAttribute{
						Required:    true,
						Description: "A pre-hashed password for the user.",
						Attributes: map[string]schema.Attribute{
							"algorithm": schema.StringAttribute{
								Required:    true,
								Description: "The hashing algorithm used.",
							},
							"hash": schema.StringAttribute{
								Required:    true,
								Sensitive:   true,
								Description: "The hex-encoded hash of the password.",
							},
						},
					},
				},
			},
			"description": schema.StringAttribute{
				Optional:    true,
				Description: "Human-readable description for the cluster.",
			},
			"dns_name": schema.StringAttribute{
				Computed:    true,
				Description: "The DNS name used to connect to the cluster's primary instance.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"eviction_policy": schema.StringAttribute{
				Required:    true,
				Description: "Defines the key eviction strategy when the memory limit is reached.",
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
				Description: "Compute configuration for each instance in the cluster.",
				Attributes: map[string]schema.Attribute{
					"cores": schema.Int32Attribute{
						Required:    true,
						Description: "The number of dedicated CPU cores per instance.",
					},
					"count": schema.Int32Attribute{
						Required:    true,
						Description: "The total number of instances in the cluster.",
					},
					"ram": schema.Int32Attribute{
						Required:    true,
						Description: "The amount of RAM per instance in gigabytes.",
					},
				},
			},
			"location": schema.StringAttribute{
				Required:    true,
				Description: "The location of the cluster. Used to route to the correct regional API endpoint. Available locations: " + inmemorydbv2Service.AvailableLocationsString() + ".",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"logs_enabled": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Enables or disables log collection and reporting for observability.",
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
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
				Description: "Enables or disables metrics collection and reporting for observability.",
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "The name of the In-Memory DB cluster.",
			},
			"persistence_mode": schema.StringAttribute{
				Required:    true,
				Description: "Specifies how and whether data is persisted to disk.",
			},
			"restore_from_snapshot": schema.SingleNestedAttribute{
				Optional:    true,
				WriteOnly:   true,
				Description: "Restores the cluster data from a snapshot.",
				Attributes: map[string]schema.Attribute{
					"source_snapshot_id": schema.StringAttribute{
						Optional:    true,
						Description: "The UUID of the snapshot to restore from.",
					},
					"recovery_target_datetime": schema.StringAttribute{
						Optional:    true,
						Description: "ISO 8601 timestamp to restore from the most recent snapshot at or before that time.",
					},
				},
			},
			"snapshot": schema.SingleNestedAttribute{
				Required:    true,
				Description: "Snapshot storage and retention configuration for the cluster.",
				Attributes: map[string]schema.Attribute{
					"location": schema.StringAttribute{
						Required:    true,
						Description: "The Object Storage location where snapshots will be stored.",
					},
					"retention_days": schema.Int32Attribute{
						Required:    true,
						Description: "Number of days snapshots are retained before automatic deletion.",
					},
					"snapshot_hours": schema.ListAttribute{
						Required:    true,
						Description: "Hours of the day (UTC) at which snapshots are taken. At least one hour must be specified.",
						ElementType: types.Int32Type,
					},
				},
			},
			"version": schema.StringAttribute{
				Required:    true,
				Description: "The In-Memory DB version. Use GET /versions to see supported values. Upgrades only (no downgrades).",
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

// Create creates a new InMemoryDB v2 cluster.
func (r *clusterResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan clusterResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	location := plan.Location.ValueString()
	client, err := r.bundle.NewInMemoryDBV2Client(ctx, location)
	if err != nil {
		resp.Diagnostics.AddError("failed to create InMemoryDB v2 client", err.Error())
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

	clusterResponse, _, err := client.CreateCluster(ctx, inmemorydbv3.ClusterCreate{Properties: createProps})
	if err != nil {
		resp.Diagnostics.AddError("failed to create InMemoryDB v2 cluster", err.Error())
		return
	}

	clusterID := clusterResponse.Id

	err = backoff.Retry(func() error {
		return client.IsClusterReady(ctx, clusterID)
	}, backoff.NewExponentialBackOff(backoff.WithMaxElapsedTime(createTimeout)))
	if err != nil {
		resp.Diagnostics.AddError("error waiting for InMemoryDB v2 cluster to become available", err.Error())
		return
	}

	cluster, _, err := client.GetCluster(ctx, clusterID)
	if err != nil {
		resp.Diagnostics.AddError("error reading InMemoryDB v2 cluster after creation", err.Error())
		return
	}

	resp.Diagnostics.Append(mapClusterResponseToModel(ctx, &cluster, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
	resp.Diagnostics.Append(resp.Identity.Set(ctx, &clusterIdentityModel{ID: plan.ID, Location: plan.Location})...)
}

// Read reads the InMemoryDB v2 cluster state.
func (r *clusterResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state clusterResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	clusterID := state.ID.ValueString()
	location := state.Location.ValueString()

	client, err := r.bundle.NewInMemoryDBV2Client(ctx, location)
	if err != nil {
		resp.Diagnostics.AddError("failed to create InMemoryDB v2 client", err.Error())
		return
	}

	cluster, apiResponse, err := client.GetCluster(ctx, clusterID)
	if err != nil {
		if apiResponse != nil && apiResponse.HttpNotFound() {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("error reading InMemoryDB v2 cluster", fmt.Sprintf("cluster ID: %s, error: %v", clusterID, err))
		return
	}

	resp.Diagnostics.Append(mapClusterResponseToModel(ctx, &cluster, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
	resp.Diagnostics.Append(resp.Identity.Set(ctx, &clusterIdentityModel{ID: state.ID, Location: state.Location})...)
}

// Update updates the InMemoryDB v2 cluster using PUT semantics.
func (r *clusterResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state clusterResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	clusterID := state.ID.ValueString()
	location := plan.Location.ValueString()

	client, err := r.bundle.NewInMemoryDBV2Client(ctx, location)
	if err != nil {
		resp.Diagnostics.AddError("failed to create InMemoryDB v2 client", err.Error())
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

	updateReq := inmemorydbv3.ClusterEnsure{
		Id:         clusterID,
		Properties: updateProps,
	}

	_, _, err = client.UpdateCluster(ctx, updateReq, clusterID)
	if err != nil {
		resp.Diagnostics.AddError("error updating InMemoryDB v2 cluster", fmt.Sprintf("cluster ID: %s, error: %v", clusterID, err))
		return
	}

	err = backoff.Retry(func() error {
		return client.IsClusterReady(ctx, clusterID)
	}, backoff.NewExponentialBackOff(backoff.WithMaxElapsedTime(updateTimeout)))
	if err != nil {
		resp.Diagnostics.AddError("error waiting for InMemoryDB v2 cluster to become available after update", fmt.Sprintf("cluster ID: %s, error: %v", clusterID, err))
		return
	}

	cluster, _, err := client.GetCluster(ctx, clusterID)
	if err != nil {
		resp.Diagnostics.AddError("error reading InMemoryDB v2 cluster after update", err.Error())
		return
	}

	resp.Diagnostics.Append(mapClusterResponseToModel(ctx, &cluster, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
	resp.Diagnostics.Append(resp.Identity.Set(ctx, &clusterIdentityModel{ID: plan.ID, Location: plan.Location})...)
}

// Delete deletes the InMemoryDB v2 cluster.
func (r *clusterResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state clusterResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	clusterID := state.ID.ValueString()
	location := state.Location.ValueString()

	client, err := r.bundle.NewInMemoryDBV2Client(ctx, location)
	if err != nil {
		resp.Diagnostics.AddError("failed to create InMemoryDB v2 client", err.Error())
		return
	}

	deleteTimeout, diags := state.Timeouts.Delete(ctx, utils.DefaultTimeout)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	_, err = client.DeleteCluster(ctx, clusterID)
	if err != nil {
		resp.Diagnostics.AddError("error deleting InMemoryDB v2 cluster", fmt.Sprintf("cluster ID: %s, error: %v", clusterID, err))
		return
	}

	err = backoff.Retry(func() error {
		return client.IsClusterDeleted(ctx, clusterID)
	}, backoff.NewExponentialBackOff(backoff.WithMaxElapsedTime(deleteTimeout)))
	if err != nil {
		resp.Diagnostics.AddError("error waiting for InMemoryDB v2 cluster to be deleted", fmt.Sprintf("cluster ID: %s, error: %v", clusterID, err))
		return
	}
}

// ImportState imports an InMemoryDB v2 cluster. Supports identity-based import (from terraform query)
// and legacy string import using "location:cluster_id" format.
func (r *clusterResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	if req.Identity != nil {
		var id *clusterIdentityModel
		resp.Diagnostics.Append(req.Identity.Get(ctx, &id)...)
		if resp.Diagnostics.HasError() {
			return
		}
		resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("location"), id.Location)...)
		resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), id.ID)...)
		return
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

func buildClusterCreateProperties(ctx context.Context, plan *clusterResourceModel) (inmemorydbv3.ClusterCreateProperties, diag.Diagnostics) {
	var diagnostics diag.Diagnostics

	snapshotConfig, diags := buildSnapshotConfig(ctx, plan.Snapshot)
	diagnostics.Append(diags...)
	if diagnostics.HasError() {
		return inmemorydbv3.ClusterCreateProperties{}, diagnostics
	}

	props := inmemorydbv3.ClusterCreateProperties{
		Name:            plan.Name.ValueString(),
		Version:         plan.Version.ValueString(),
		PersistenceMode: inmemorydbv3.PersistenceMode(plan.PersistenceMode.ValueString()),
		EvictionPolicy:  inmemorydbv3.EvictionPolicy(plan.EvictionPolicy.ValueString()),
		Instances: inmemorydbv3.InstanceConfiguration{
			Count: plan.Instances.Count.ValueInt32(),
			Cores: plan.Instances.Cores.ValueInt32(),
			Ram:   plan.Instances.RAM.ValueInt32(),
		},
		Connection: inmemorydbv3.ClusterConnection{
			DatacenterId:           plan.Connections.DatacenterID.ValueString(),
			LanId:                  plan.Connections.LanID.ValueString(),
			PrimaryInstanceAddress: plan.Connections.PrimaryInstanceAddress.ValueString(),
		},
		Snapshot: snapshotConfig,
		MaintenanceWindow: inmemorydbv3.MaintenanceWindow{
			Time:         plan.MaintenanceWindow.Time.ValueString(),
			DayOfTheWeek: inmemorydbv3.DayOfTheWeek(plan.MaintenanceWindow.DayOfTheWeek.ValueString()),
		},
		Credentials: inmemorydbv3.ClusterCredentials{
			Username: plan.Credentials.Username.ValueString(),
			Password: inmemorydbv3.HashedPassword{
				Algorithm: plan.Credentials.Password.Algorithm.ValueString(),
				Hash:      plan.Credentials.Password.Hash.ValueString(),
			},
		},
	}

	props.Description = plan.Description.ValueStringPointer()

	if !plan.LogsEnabled.IsUnknown() {
		props.LogsEnabled = plan.LogsEnabled.ValueBoolPointer()
	}
	if !plan.MetricsEnabled.IsUnknown() {
		props.MetricsEnabled = plan.MetricsEnabled.ValueBoolPointer()
	}

	if plan.RestoreFromSnapshot != nil {
		restore := inmemorydbv3.NewRestoreClusterFromSnapshot(plan.RestoreFromSnapshot.SourceSnapshotID.ValueString())
		if !plan.RestoreFromSnapshot.RecoveryTargetDatetime.IsNull() && !plan.RestoreFromSnapshot.RecoveryTargetDatetime.IsUnknown() {
			t, err := time.Parse(time.RFC3339, plan.RestoreFromSnapshot.RecoveryTargetDatetime.ValueString())
			if err != nil {
				diagnostics.AddError("invalid recovery_target_datetime", err.Error())
				return inmemorydbv3.ClusterCreateProperties{}, diagnostics
			}
			restore.SetRecoveryTargetDatetime(t)
		}
		wrapped := inmemorydbv3.RestoreClusterFromSnapshotAsClusterRestoreFromSnapshot(restore)
		props.RestoreFromSnapshot = &wrapped
	}

	return props, diagnostics
}

func buildClusterUpdateProperties(ctx context.Context, plan *clusterResourceModel) (inmemorydbv3.Cluster, diag.Diagnostics) {
	var diagnostics diag.Diagnostics

	snapshotConfig, diags := buildSnapshotConfig(ctx, plan.Snapshot)
	diagnostics.Append(diags...)
	if diagnostics.HasError() {
		return inmemorydbv3.Cluster{}, diagnostics
	}

	props := inmemorydbv3.Cluster{
		Name:            plan.Name.ValueString(),
		Version:         plan.Version.ValueString(),
		PersistenceMode: inmemorydbv3.PersistenceMode(plan.PersistenceMode.ValueString()),
		EvictionPolicy:  inmemorydbv3.EvictionPolicy(plan.EvictionPolicy.ValueString()),
		Instances: inmemorydbv3.InstanceConfiguration{
			Count: plan.Instances.Count.ValueInt32(),
			Cores: plan.Instances.Cores.ValueInt32(),
			Ram:   plan.Instances.RAM.ValueInt32(),
		},
		Connection: inmemorydbv3.ClusterConnection{
			DatacenterId:           plan.Connections.DatacenterID.ValueString(),
			LanId:                  plan.Connections.LanID.ValueString(),
			PrimaryInstanceAddress: plan.Connections.PrimaryInstanceAddress.ValueString(),
		},
		Snapshot: snapshotConfig,
		MaintenanceWindow: inmemorydbv3.MaintenanceWindow{
			Time:         plan.MaintenanceWindow.Time.ValueString(),
			DayOfTheWeek: inmemorydbv3.DayOfTheWeek(plan.MaintenanceWindow.DayOfTheWeek.ValueString()),
		},
	}

	props.Description = plan.Description.ValueStringPointer()

	if !plan.LogsEnabled.IsUnknown() {
		props.LogsEnabled = plan.LogsEnabled.ValueBoolPointer()
	}
	if !plan.MetricsEnabled.IsUnknown() {
		props.MetricsEnabled = plan.MetricsEnabled.ValueBoolPointer()
	}

	if plan.Credentials != nil {
		props.Credentials = &inmemorydbv3.ClusterCredentials{
			Username: plan.Credentials.Username.ValueString(),
			Password: inmemorydbv3.HashedPassword{
				Algorithm: plan.Credentials.Password.Algorithm.ValueString(),
				Hash:      plan.Credentials.Password.Hash.ValueString(),
			},
		}
	}

	if plan.RestoreFromSnapshot != nil && !plan.RestoreFromSnapshot.RecoveryTargetDatetime.IsNull() && !plan.RestoreFromSnapshot.RecoveryTargetDatetime.IsUnknown() {
		t, err := time.Parse(time.RFC3339, plan.RestoreFromSnapshot.RecoveryTargetDatetime.ValueString())
		if err != nil {
			diagnostics.AddError("invalid recovery_target_datetime", err.Error())
			return inmemorydbv3.Cluster{}, diagnostics
		}
		restore := inmemorydbv3.NewInPlaceRestoreClusterFromSnapshot(t)
		wrapped := inmemorydbv3.InPlaceRestoreClusterFromSnapshotAsClusterRestoreFromSnapshot(restore)
		props.RestoreFromSnapshot = &wrapped
	}

	return props, diagnostics
}

func buildSnapshotConfig(ctx context.Context, m *snapshotConfigModel) (inmemorydbv3.SnapshotConfiguration, diag.Diagnostics) {
	var diagnostics diag.Diagnostics
	cfg := inmemorydbv3.SnapshotConfiguration{
		Location:      m.Location.ValueString(),
		RetentionDays: m.RetentionDays.ValueInt32(),
	}

	var hours []types.Int32
	diagnostics.Append(m.SnapshotHours.ElementsAs(ctx, &hours, false)...)
	if !diagnostics.HasError() {
		sdkHours := make([]int32, len(hours))
		for i, h := range hours {
			sdkHours[i] = h.ValueInt32()
		}
		cfg.SnapshotHours = sdkHours
	}

	return cfg, diagnostics
}

func mapClusterResponseToModel(ctx context.Context, cluster *inmemorydbv3.ClusterRead, model *clusterResourceModel) diag.Diagnostics {
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

	// The API never returns the password hash.
	// Preserve username and password from the existing model (state/plan).
	if props.Credentials != nil && model.Credentials != nil {
		existingPassword := model.Credentials.Password
		model.Credentials = &credentialsModel{
			Username: types.StringValue(props.Credentials.Username),
			Password: existingPassword,
		}
	}

	return diagnostics
}
