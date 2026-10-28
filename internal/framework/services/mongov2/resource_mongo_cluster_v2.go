package mongov2

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/cenkalti/backoff/v4"
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/identityschema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	mongov3 "github.com/ionos-cloud/sdk-go-bundle/products/dbaas/mongo/v3"

	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/bundleclient"
	mongov2service "github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/dbaas/mongov2"
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
	BiConnectorEnabled  types.Bool                `tfsdk:"bi_connector_enabled"`
	ConnectionURI       types.String              `tfsdk:"connection_uri"`
	Connections         *connectionModel          `tfsdk:"connections"`
	Credentials         *credentialsModel         `tfsdk:"credentials"`
	Description         types.String              `tfsdk:"description"`
	Edition             types.String              `tfsdk:"edition"`
	ID                  types.String              `tfsdk:"id"`
	InstanceCount       types.Int32               `tfsdk:"instance_count"`
	InstanceTemplateID  types.String              `tfsdk:"instance_template_id"`
	Instances           *instancesModel           `tfsdk:"instances"`
	Location            types.String              `tfsdk:"location"`
	LogsEnabled         types.Bool                `tfsdk:"logs_enabled"`
	MaintenanceWindow   *maintenanceWindowModel   `tfsdk:"maintenance_window"`
	MetricsEnabled      types.Bool                `tfsdk:"metrics_enabled"`
	Name                types.String              `tfsdk:"name"`
	RestoreFromSnapshot *restoreFromSnapshotModel `tfsdk:"restore_from_snapshot"`
	Sharding            *shardingModel            `tfsdk:"sharding"`
	Snapshot            *snapshotConfigModel      `tfsdk:"snapshot"`
	Timeouts            timeouts.Value            `tfsdk:"timeouts"`
	Version             types.String              `tfsdk:"version"`
}

type instancesModel struct {
	Cores       types.Int32 `tfsdk:"cores"`
	Count       types.Int32 `tfsdk:"count"`
	RAM         types.Int32 `tfsdk:"ram"`
	StorageSize types.Int32 `tfsdk:"storage_size"`
}

type shardingModel struct {
	Enabled types.Bool  `tfsdk:"enabled"`
	Shards  types.Int32 `tfsdk:"shards"`
}

type connectionModel struct {
	DatacenterID      types.String `tfsdk:"datacenter_id"`
	LanID             types.String `tfsdk:"lan_id"`
	InstanceAddresses types.List   `tfsdk:"instance_addresses"`
}

type snapshotConfigModel struct {
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

type restoreFromSnapshotModel struct {
	RecoveryTargetDatetime types.String `tfsdk:"recovery_target_datetime"`
	SourceSnapshotID       types.String `tfsdk:"source_snapshot_id"`
}

// NewClusterResource creates a new resource for the MongoDB v2 cluster.
func NewClusterResource() resource.Resource {
	return &clusterResource{}
}

// Metadata returns the resource type name.
func (r *clusterResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_mongo_cluster_v2"
}

func (r *clusterResource) IdentitySchema(_ context.Context, _ resource.IdentitySchemaRequest, resp *resource.IdentitySchemaResponse) {
	resp.IdentitySchema = identityschema.Schema{
		Attributes: map[string]identityschema.Attribute{
			"id":       identityschema.StringAttribute{RequiredForImport: true},
			"location": identityschema.StringAttribute{RequiredForImport: true},
		},
	}
}

// Schema returns the schema for the MongoDB v2 cluster resource.
func (r *clusterResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages an IONOS CLOUD MongoDB V2 cluster.",
		Attributes: map[string]schema.Attribute{
			"bi_connector_enabled": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Enables or disables the MongoDB BI Connector for this cluster. Can only be enabled for the ENTERPRISE edition. Defaults to false.",
			},
			"connection_uri": schema.StringAttribute{
				Computed:    true,
				Description: "The connection string used to access the cluster.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"connections": schema.SingleNestedAttribute{
				Required:    true,
				Description: "Connection information of the MongoDB cluster.",
				Attributes: map[string]schema.Attribute{
					"datacenter_id": schema.StringAttribute{
						Required:    true,
						Description: "The datacenter to connect your instance to. Immutable: cannot be changed after cluster creation.",
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.RequiresReplace(),
						},
					},
					"lan_id": schema.StringAttribute{
						Required:    true,
						Description: "The numeric LAN ID to connect your instance to. Immutable: cannot be changed after cluster creation.",
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.RequiresReplace(),
						},
					},
					"instance_addresses": schema.ListAttribute{
						Required:    true,
						ElementType: types.StringType,
						Description: "A list of IP addresses (with netmask, in CIDR notation) assigned to the cluster instances. The number of addresses must match the number of instances in the cluster. The range 10.208.0.0/12 is not available.",
					},
				},
			},
			"credentials": schema.SingleNestedAttribute{
				Required:    true,
				Description: "Credentials for the initial database user to be created.",
				Attributes: map[string]schema.Attribute{
					"database": schema.StringAttribute{
						Required:    true,
						Description: "The name of the initial database to be created. Between 2 and 63 characters, containing only alphanumeric characters and underscores.",
					},
					"password": schema.StringAttribute{
						Required:    true,
						Sensitive:   true,
						Description: "The password for the initial MongoDB user. Between 8 and 256 characters. Required whenever credentials are supplied; never returned by the API, so it will be null in state after `terraform import`.",
					},
					"username": schema.StringAttribute{
						Required:    true,
						Description: "The username of the initial MongoDB user. Between 2 and 16 characters, containing only alphanumeric characters and underscores. Some usernames are reserved for platform use.",
					},
				},
			},
			"description": schema.StringAttribute{
				Optional:    true,
				Description: "Human-readable description for the cluster.",
			},
			"edition": schema.StringAttribute{
				Required:    true,
				Description: "The edition of the MongoDB cluster. One of `ENTERPRISE` or `BUSINESS`. Changing this forces re-creation.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "The ID (UUID) of the cluster.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"instance_count": schema.Int32Attribute{
				Optional:    true,
				Description: "The total number of instances in the cluster, one primary and n-1 secondary (between 1 and 9). Required for the BUSINESS edition.",
			},
			"instance_template_id": schema.StringAttribute{
				Optional:    true,
				Description: "The ID (UUID) of the instance template defining cores, RAM and storage. See the `ionoscloud_mongo_templates_v2` data source. Required for the BUSINESS edition. Immutable: changing it forces re-creation.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"instances": schema.SingleNestedAttribute{
				Optional:    true,
				Description: "Compute and storage configuration for each instance. Required for the ENTERPRISE edition.",
				Attributes: map[string]schema.Attribute{
					"cores": schema.Int32Attribute{
						Required:    true,
						Description: "The number of CPU cores per instance (between 1 and 62).",
					},
					"count": schema.Int32Attribute{
						Required:    true,
						Description: "The total number of instances in the cluster, one primary and n-1 secondary (between 1 and 9).",
					},
					"ram": schema.Int32Attribute{
						Required:    true,
						Description: "The amount of memory per instance in gigabytes (GB) (between 6 and 240).",
					},
					"storage_size": schema.Int32Attribute{
						Required:    true,
						Description: "The amount of storage per instance in gigabytes (GB) (between 10 and 4096).",
					},
				},
			},
			"location": schema.StringAttribute{
				Required:    true,
				Description: "The location of the cluster. Changing this forces re-creation. Available locations: " + mongov2service.AvailableLocationsString() + ".",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"logs_enabled": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Allows or disallows the collection and reporting of logs for this cluster's observability. If the observability service is not activated on the contract, this setting is accepted but has no effect. Defaults to false.",
			},
			"maintenance_window": schema.SingleNestedAttribute{
				Required:    true,
				Description: "A weekly 4 hour-long window, during which maintenance might occur.",
				Attributes: map[string]schema.Attribute{
					"day_of_the_week": schema.StringAttribute{
						Required:    true,
						Description: "The name of the week day. One of `Monday`, `Tuesday`, `Wednesday`, `Thursday`, `Friday`, `Saturday`, `Sunday`.",
						Validators: []validator.String{
							stringvalidator.OneOf("Sunday", "Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday"),
						},
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
				Description: "Allows or disallows the collection and reporting of metrics for this cluster's observability. If the observability service is not activated on the contract, this setting is accepted but has no effect. Defaults to false.",
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "The name of the MongoDB cluster. Between 2 and 63 characters; must begin and end with an alphanumeric character, with dashes, underscores, dots and alphanumerics in between.",
			},
			"restore_from_snapshot": schema.SingleNestedAttribute{
				Optional:    true,
				Description: "Restores the cluster from a snapshot. On cluster creation, set source_snapshot_id (optionally recovery_target_datetime) to initialize the new cluster from an existing snapshot in the same region. On in-place update, set recovery_target_datetime only; the restore source is inferred from the cluster's own snapshots.",
				Attributes: map[string]schema.Attribute{
					"recovery_target_datetime": schema.StringAttribute{
						Optional:    true,
						Description: "An ISO 8601 timestamp for point-in-time recovery; the cluster is restored to this exact moment within the available recovery window. Optional on creation (if omitted, the source snapshot is applied in full); required for in-place restore.",
					},
					"source_snapshot_id": schema.StringAttribute{
						Optional:    true,
						Description: "UUID of the snapshot to restore from. Required for restore on cluster creation; not valid for in-place restore, where the source is inferred from the cluster's own snapshots.",
					},
				},
			},
			"sharding": schema.SingleNestedAttribute{
				Optional:    true,
				Description: "Sharding properties for the cluster. Can only be set for the ENTERPRISE edition.",
				Attributes: map[string]schema.Attribute{
					"enabled": schema.BoolAttribute{
						Required:    true,
						Description: "Activates or deactivates sharding for this cluster. Can only be set to true for the ENTERPRISE edition.",
					},
					"shards": schema.Int32Attribute{
						Optional:    true,
						Computed:    true,
						Description: "The number of shards to create in the cluster (between 2 and 5). When sharding is enabled, the instances count signifies the number of replicas per shard.",
					},
				},
			},
			"snapshot": schema.SingleNestedAttribute{
				Required:    true,
				Description: "Snapshot storage and retention configuration for the cluster.",
				Attributes: map[string]schema.Attribute{
					"location": schema.StringAttribute{
						Required:    true,
						Description: "The Object Storage location where the snapshots will be created. For added data safety, use a different location than the cluster.",
					},
					"retention_days": schema.Int32Attribute{
						Required:    true,
						Description: "Configures how many days cluster snapshots are retained before being automatically deleted (between 1 and 365). Reducing this value purges existing snapshots that fall outside the new retention window.",
					},
				},
			},
			"version": schema.StringAttribute{
				Required:    true,
				Description: "The MongoDB version for the cluster, e.g. `6.0`, `7.0`, `8.0`. Use the `ionoscloud_mongo_versions_v2` data source to retrieve the list of supported versions and their available upgrade paths.",
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

// Create creates a new MongoDB v2 cluster.
func (r *clusterResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan clusterResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(validateEdition(&plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if plan.RestoreFromSnapshot != nil && plan.RestoreFromSnapshot.SourceSnapshotID.IsNull() {
		resp.Diagnostics.AddError(
			"missing source_snapshot_id",
			"source_snapshot_id is required when restore_from_snapshot is set during cluster creation.",
		)
		return
	}

	location := plan.Location.ValueString()
	client, err := r.bundle.NewMongoDBV2Client(ctx, location)
	if err != nil {
		resp.Diagnostics.AddError("failed to create MongoDB V2 client", err.Error())
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

	clusterResponse, apiResponse, err := client.CreateCluster(ctx, mongov3.ClusterCreate{Properties: createProps})
	if err != nil {
		resp.Diagnostics.AddError("failed to create MongoDB V2 cluster", diagutil.WrapError(err, &diagutil.ErrorContext{
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
		resp.Diagnostics.AddError("error waiting for MongoDB V2 cluster to become available", diagutil.WrapError(err, &diagutil.ErrorContext{
			ResourceID:     clusterID,
			ResourceName:   plan.Name.ValueString(),
			AdditionalInfo: map[string]string{"location": location},
		}).Error())
		return
	}

	cluster, apiResponseGet, err := client.GetCluster(ctx, clusterID)
	if err != nil {
		resp.Diagnostics.AddError("error reading MongoDB V2 cluster after creation", diagutil.WrapError(err, &diagutil.ErrorContext{
			ResourceID:     clusterID,
			ResourceName:   plan.Name.ValueString(),
			StatusCode:     apiResponseGet.SafeStatusCode(),
			AdditionalInfo: map[string]string{"location": location},
		}).Error())
		return
	}

	resp.Diagnostics.Append(mapClusterResponseToModel(ctx, &cluster, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
	resp.Diagnostics.Append(resp.Identity.Set(ctx, &clusterIdentityModel{ID: plan.ID, Location: plan.Location})...)
}

// Read reads the MongoDB v2 cluster state.
func (r *clusterResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state clusterResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	clusterID := state.ID.ValueString()
	location := state.Location.ValueString()

	client, err := r.bundle.NewMongoDBV2Client(ctx, location)
	if err != nil {
		resp.Diagnostics.AddError("failed to create MongoDB V2 client", err.Error())
		return
	}

	cluster, apiResponse, err := client.GetCluster(ctx, clusterID)
	if err != nil {
		if apiResponse.HttpNotFound() {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("error reading MongoDB V2 cluster", diagutil.WrapError(err, &diagutil.ErrorContext{
			ResourceID:     clusterID,
			ResourceName:   state.Name.ValueString(),
			StatusCode:     apiResponse.SafeStatusCode(),
			AdditionalInfo: map[string]string{"location": location},
		}).Error())
		return
	}

	resp.Diagnostics.Append(mapClusterResponseToModel(ctx, &cluster, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
	resp.Diagnostics.Append(resp.Identity.Set(ctx, &clusterIdentityModel{ID: state.ID, Location: state.Location})...)
}

// Update updates the MongoDB v2 cluster using PUT semantics.
func (r *clusterResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state clusterResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(validateEdition(&plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Only request an in-place restore if restore_from_snapshot was newly set/changed vs
	// prior state; a block carried over unchanged must not re-trigger a restore on every update.
	restore, diags := buildUpdateRestore(&plan, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	clusterID := plan.ID.ValueString()
	location := plan.Location.ValueString()

	client, err := r.bundle.NewMongoDBV2Client(ctx, location)
	if err != nil {
		resp.Diagnostics.AddError("failed to create MongoDB V2 client", err.Error())
		return
	}

	updateTimeout, diags := plan.Timeouts.Update(ctx, utils.DefaultTimeout)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	updateProps, diags := buildClusterUpdateProperties(ctx, &plan, restore)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	updateReq := mongov3.ClusterEnsure{
		Id:         clusterID,
		Properties: updateProps,
	}

	_, apiResponseUpdate, err := client.UpdateCluster(ctx, updateReq, clusterID)
	if err != nil {
		resp.Diagnostics.AddError("error updating MongoDB V2 cluster", diagutil.WrapError(err, &diagutil.ErrorContext{
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
		resp.Diagnostics.AddError("error waiting for MongoDB V2 cluster to become available after update", diagutil.WrapError(err, &diagutil.ErrorContext{
			ResourceID:     clusterID,
			ResourceName:   state.Name.ValueString(),
			AdditionalInfo: map[string]string{"location": location},
		}).Error())
		return
	}

	cluster, apiResponseGetAfterUpdate, err := client.GetCluster(ctx, clusterID)
	if err != nil {
		resp.Diagnostics.AddError("error reading MongoDB V2 cluster after update", diagutil.WrapError(err, &diagutil.ErrorContext{
			ResourceID:     clusterID,
			ResourceName:   state.Name.ValueString(),
			StatusCode:     apiResponseGetAfterUpdate.SafeStatusCode(),
			AdditionalInfo: map[string]string{"location": location},
		}).Error())
		return
	}

	resp.Diagnostics.Append(mapClusterResponseToModel(ctx, &cluster, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
	resp.Diagnostics.Append(resp.Identity.Set(ctx, &clusterIdentityModel{ID: plan.ID, Location: plan.Location})...)
}

// Delete deletes the MongoDB v2 cluster.
func (r *clusterResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state clusterResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	clusterID := state.ID.ValueString()
	location := state.Location.ValueString()

	client, err := r.bundle.NewMongoDBV2Client(ctx, location)
	if err != nil {
		resp.Diagnostics.AddError("failed to create MongoDB V2 client", err.Error())
		return
	}

	deleteTimeout, diags := state.Timeouts.Delete(ctx, utils.DefaultTimeout)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	apiResponseDelete, err := client.DeleteCluster(ctx, clusterID)
	if err != nil {
		resp.Diagnostics.AddError("error deleting MongoDB V2 cluster", diagutil.WrapError(err, &diagutil.ErrorContext{
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
		resp.Diagnostics.AddError("error waiting for MongoDB V2 cluster to be deleted", diagutil.WrapError(err, &diagutil.ErrorContext{
			ResourceID:     clusterID,
			ResourceName:   state.Name.ValueString(),
			AdditionalInfo: map[string]string{"location": location},
		}).Error())
		return
	}
}

// ImportState imports a MongoDB v2 cluster. Supports identity-based import and legacy string import using "location:cluster_id" format.
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

// validateEdition checks that the edition-specific attributes match the chosen edition.
func validateEdition(plan *clusterResourceModel) diag.Diagnostics {
	var diagnostics diag.Diagnostics

	switch plan.Edition.ValueString() {
	case string(mongov3.MONGOEDITION_ENTERPRISE):
		if plan.Instances == nil {
			diagnostics.AddError("missing instances", "instances is required for the ENTERPRISE edition.")
		}
		if !plan.InstanceTemplateID.IsNull() || !plan.InstanceCount.IsNull() {
			diagnostics.AddError("invalid attributes for ENTERPRISE edition", "instance_template_id and instance_count are only valid for the BUSINESS edition.")
		}
	case string(mongov3.MONGOEDITION_BUSINESS):
		if plan.InstanceTemplateID.IsNull() || plan.InstanceCount.IsNull() {
			diagnostics.AddError("missing business attributes", "instance_template_id and instance_count are required for the BUSINESS edition.")
		}
		if plan.Instances != nil || plan.Sharding != nil || !plan.BiConnectorEnabled.IsUnknown() {
			diagnostics.AddError("invalid attributes for BUSINESS edition", "instances, sharding and bi_connector_enabled are only valid for the ENTERPRISE edition.")
		}
	default:
		diagnostics.AddError("invalid edition", fmt.Sprintf("edition must be one of %q or %q.", mongov3.MONGOEDITION_ENTERPRISE, mongov3.MONGOEDITION_BUSINESS))
	}

	return diagnostics
}

func buildClusterCreateProperties(ctx context.Context, plan *clusterResourceModel) (mongov3.ClusterCreateProperties, diag.Diagnostics) {
	var diagnostics diag.Diagnostics
	var props mongov3.ClusterCreateProperties

	connection, diags := buildConnection(ctx, plan)
	diagnostics.Append(diags...)
	restore, diags := buildRestore(plan)
	diagnostics.Append(diags...)
	if diagnostics.HasError() {
		return props, diagnostics
	}

	credentials := mongov3.MongoUser{
		Username: plan.Credentials.Username.ValueString(),
		Password: plan.Credentials.Password.ValueStringPointer(),
		Database: plan.Credentials.Database.ValueString(),
	}
	maintenanceWindow := mongov3.MaintenanceWindow{
		DayOfTheWeek: mongov3.DayOfTheWeek(plan.MaintenanceWindow.DayOfTheWeek.ValueString()),
		Time:         plan.MaintenanceWindow.Time.ValueString(),
	}
	snapshot := mongov3.ClusterSnapshot{
		Location:      plan.Snapshot.Location.ValueString(),
		RetentionDays: plan.Snapshot.RetentionDays.ValueInt32(),
	}

	switch plan.Edition.ValueString() {
	case string(mongov3.MONGOEDITION_ENTERPRISE):
		enterprise := mongov3.EnterpriseClusterCreate{
			Name:              plan.Name.ValueString(),
			Description:       plan.Description.ValueStringPointer(),
			Version:           plan.Version.ValueString(),
			Connection:        connection,
			MaintenanceWindow: maintenanceWindow,
			Credentials:       credentials,
			Edition:           mongov3.MONGOEDITION_ENTERPRISE,
			Instances: mongov3.InstanceConfiguration{
				Cores:       plan.Instances.Cores.ValueInt32(),
				Count:       plan.Instances.Count.ValueInt32(),
				Ram:         plan.Instances.RAM.ValueInt32(),
				StorageSize: plan.Instances.StorageSize.ValueInt32(),
			},
			Snapshot:            snapshot,
			RestoreFromSnapshot: restore,
		}
		if !plan.BiConnectorEnabled.IsUnknown() {
			enterprise.BiConnectorEnabled = plan.BiConnectorEnabled.ValueBoolPointer()
		}
		if !plan.LogsEnabled.IsUnknown() {
			enterprise.LogsEnabled = plan.LogsEnabled.ValueBoolPointer()
		}
		if !plan.MetricsEnabled.IsUnknown() {
			enterprise.MetricsEnabled = plan.MetricsEnabled.ValueBoolPointer()
		}
		if plan.Sharding != nil {
			sharding := &mongov3.MongoSharding{
				Enabled: plan.Sharding.Enabled.ValueBoolPointer(),
			}
			if !plan.Sharding.Shards.IsNull() && !plan.Sharding.Shards.IsUnknown() {
				sharding.Shards = plan.Sharding.Shards.ValueInt32Pointer()
			}
			enterprise.Sharding = sharding
		}
		return mongov3.EnterpriseClusterCreateAsClusterCreateProperties(&enterprise), diagnostics
	case string(mongov3.MONGOEDITION_BUSINESS):
		business := mongov3.BusinessClusterCreate{
			Name:                plan.Name.ValueString(),
			Description:         plan.Description.ValueStringPointer(),
			Version:             plan.Version.ValueString(),
			Connection:          connection,
			MaintenanceWindow:   maintenanceWindow,
			Credentials:         credentials,
			Edition:             mongov3.MONGOEDITION_BUSINESS,
			InstanceTemplateId:  plan.InstanceTemplateID.ValueString(),
			InstanceCount:       plan.InstanceCount.ValueInt32(),
			Snapshot:            snapshot,
			RestoreFromSnapshot: restore,
		}
		if !plan.LogsEnabled.IsUnknown() {
			business.LogsEnabled = plan.LogsEnabled.ValueBoolPointer()
		}
		if !plan.MetricsEnabled.IsUnknown() {
			business.MetricsEnabled = plan.MetricsEnabled.ValueBoolPointer()
		}
		return mongov3.BusinessClusterCreateAsClusterCreateProperties(&business), diagnostics
	default:
		diagnostics.AddError("invalid edition", fmt.Sprintf("edition must be one of %q or %q.", mongov3.MONGOEDITION_ENTERPRISE, mongov3.MONGOEDITION_BUSINESS))
		return props, diagnostics
	}
}

func buildClusterUpdateProperties(ctx context.Context, plan *clusterResourceModel, restore *mongov3.MongoRestoreFromSnapshot) (mongov3.Cluster, diag.Diagnostics) {
	var diagnostics diag.Diagnostics
	var props mongov3.Cluster

	connection, diags := buildConnection(ctx, plan)
	diagnostics.Append(diags...)
	if diagnostics.HasError() {
		return props, diagnostics
	}

	credentials := &mongov3.MongoUser{
		Username: plan.Credentials.Username.ValueString(),
		Password: plan.Credentials.Password.ValueStringPointer(),
		Database: plan.Credentials.Database.ValueString(),
	}
	maintenanceWindow := mongov3.MaintenanceWindow{
		DayOfTheWeek: mongov3.DayOfTheWeek(plan.MaintenanceWindow.DayOfTheWeek.ValueString()),
		Time:         plan.MaintenanceWindow.Time.ValueString(),
	}
	snapshot := mongov3.ClusterSnapshot{
		Location:      plan.Snapshot.Location.ValueString(),
		RetentionDays: plan.Snapshot.RetentionDays.ValueInt32(),
	}

	switch plan.Edition.ValueString() {
	case string(mongov3.MONGOEDITION_ENTERPRISE):
		enterprise := mongov3.EnterpriseCluster{
			Name:              plan.Name.ValueString(),
			Description:       plan.Description.ValueStringPointer(),
			Version:           plan.Version.ValueString(),
			Connection:        connection,
			MaintenanceWindow: maintenanceWindow,
			Credentials:       credentials,
			Edition:           mongov3.MONGOEDITION_ENTERPRISE,
			Instances: mongov3.InstanceConfiguration{
				Cores:       plan.Instances.Cores.ValueInt32(),
				Count:       plan.Instances.Count.ValueInt32(),
				Ram:         plan.Instances.RAM.ValueInt32(),
				StorageSize: plan.Instances.StorageSize.ValueInt32(),
			},
			Snapshot:            snapshot,
			RestoreFromSnapshot: restore,
		}
		if !plan.BiConnectorEnabled.IsUnknown() {
			enterprise.BiConnectorEnabled = plan.BiConnectorEnabled.ValueBoolPointer()
		}
		if !plan.LogsEnabled.IsUnknown() {
			enterprise.LogsEnabled = plan.LogsEnabled.ValueBoolPointer()
		}
		if !plan.MetricsEnabled.IsUnknown() {
			enterprise.MetricsEnabled = plan.MetricsEnabled.ValueBoolPointer()
		}
		if plan.Sharding != nil {
			sharding := &mongov3.MongoSharding{
				Enabled: plan.Sharding.Enabled.ValueBoolPointer(),
			}
			if !plan.Sharding.Shards.IsNull() && !plan.Sharding.Shards.IsUnknown() {
				sharding.Shards = plan.Sharding.Shards.ValueInt32Pointer()
			}
			enterprise.Sharding = sharding
		}
		return mongov3.EnterpriseClusterAsCluster(&enterprise), diagnostics
	case string(mongov3.MONGOEDITION_BUSINESS):
		business := mongov3.BusinessCluster{
			Name:                plan.Name.ValueString(),
			Description:         plan.Description.ValueStringPointer(),
			Version:             plan.Version.ValueString(),
			Connection:          connection,
			MaintenanceWindow:   maintenanceWindow,
			Credentials:         credentials,
			Edition:             mongov3.MONGOEDITION_BUSINESS,
			InstanceTemplateId:  plan.InstanceTemplateID.ValueString(),
			InstanceCount:       plan.InstanceCount.ValueInt32(),
			Snapshot:            snapshot,
			RestoreFromSnapshot: restore,
		}
		if !plan.LogsEnabled.IsUnknown() {
			business.LogsEnabled = plan.LogsEnabled.ValueBoolPointer()
		}
		if !plan.MetricsEnabled.IsUnknown() {
			business.MetricsEnabled = plan.MetricsEnabled.ValueBoolPointer()
		}
		return mongov3.BusinessClusterAsCluster(&business), diagnostics
	default:
		diagnostics.AddError("invalid edition", fmt.Sprintf("edition must be one of %q or %q.", mongov3.MONGOEDITION_ENTERPRISE, mongov3.MONGOEDITION_BUSINESS))
		return props, diagnostics
	}
}

func buildConnection(ctx context.Context, plan *clusterResourceModel) (mongov3.MongoClusterConnection, diag.Diagnostics) {
	var instanceAddresses []string
	diags := plan.Connections.InstanceAddresses.ElementsAs(ctx, &instanceAddresses, false)
	return mongov3.MongoClusterConnection{
		DatacenterId:      plan.Connections.DatacenterID.ValueString(),
		LanId:             plan.Connections.LanID.ValueString(),
		InstanceAddresses: instanceAddresses,
	}, diags
}

func buildRestore(plan *clusterResourceModel) (*mongov3.MongoRestoreFromSnapshot, diag.Diagnostics) {
	var diagnostics diag.Diagnostics
	if plan.RestoreFromSnapshot == nil {
		return nil, diagnostics
	}
	restore := mongov3.MongoRestoreFromSnapshot{}
	if !plan.RestoreFromSnapshot.SourceSnapshotID.IsNull() {
		restore.SourceSnapshotId = plan.RestoreFromSnapshot.SourceSnapshotID.ValueStringPointer()
	}
	if !plan.RestoreFromSnapshot.RecoveryTargetDatetime.IsNull() {
		t, err := time.Parse(time.RFC3339, plan.RestoreFromSnapshot.RecoveryTargetDatetime.ValueString())
		if err != nil {
			diagnostics.AddError("invalid recovery_target_datetime", err.Error())
			return nil, diagnostics
		}
		restore.RecoveryTargetDatetime = &mongov3.IonosTime{Time: t}
	}
	return &restore, diagnostics
}

// restoreUnchanged reports whether restore_from_snapshot equals the last-applied value.
// The block is write-only (never read back from the API), so state holds the prior config.
func restoreUnchanged(plan, state *clusterResourceModel) bool {
	if plan.RestoreFromSnapshot == nil && state.RestoreFromSnapshot == nil {
		return true
	}
	if plan.RestoreFromSnapshot == nil || state.RestoreFromSnapshot == nil {
		return false
	}
	return plan.RestoreFromSnapshot.SourceSnapshotID.Equal(state.RestoreFromSnapshot.SourceSnapshotID) &&
		plan.RestoreFromSnapshot.RecoveryTargetDatetime.Equal(state.RestoreFromSnapshot.RecoveryTargetDatetime)
}

// buildUpdateRestore returns an in-place restore request only when restore_from_snapshot was
// newly set or changed vs prior state. An in-place restore uses recovery_target_datetime only;
// source_snapshot_id is create-only and intentionally omitted.
func buildUpdateRestore(plan, state *clusterResourceModel) (*mongov3.MongoRestoreFromSnapshot, diag.Diagnostics) {
	var diagnostics diag.Diagnostics
	if plan.RestoreFromSnapshot == nil || restoreUnchanged(plan, state) {
		return nil, diagnostics
	}
	if plan.RestoreFromSnapshot.RecoveryTargetDatetime.IsNull() {
		diagnostics.AddError(
			"missing recovery_target_datetime",
			"recovery_target_datetime is required to request an in-place restore during update.",
		)
		return nil, diagnostics
	}
	t, err := time.Parse(time.RFC3339, plan.RestoreFromSnapshot.RecoveryTargetDatetime.ValueString())
	if err != nil {
		diagnostics.AddError("invalid recovery_target_datetime", err.Error())
		return nil, diagnostics
	}
	return &mongov3.MongoRestoreFromSnapshot{RecoveryTargetDatetime: &mongov3.IonosTime{Time: t}}, diagnostics
}

func mapClusterResponseToModel(ctx context.Context, cluster *mongov3.ClusterRead, model *clusterResourceModel) diag.Diagnostics {
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
		diagnostics.Append(mapCommonFields(ctx, model, ent.Connection, ent.MaintenanceWindow, ent.Snapshot, ent.Credentials)...)
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
		// bi_connector_enabled is Optional+Computed and ENTERPRISE-only; set it known (null) for BUSINESS.
		model.BiConnectorEnabled = types.BoolNull()
		diagnostics.Append(mapCommonFields(ctx, model, bus.Connection, bus.MaintenanceWindow, bus.Snapshot, bus.Credentials)...)
	default:
		diagnostics.AddError("unknown cluster edition", "the API returned a cluster with neither ENTERPRISE nor BUSINESS properties set.")
	}

	return diagnostics
}

func mapCommonFields(ctx context.Context, model *clusterResourceModel, connection mongov3.MongoClusterConnection, maintenanceWindow mongov3.MaintenanceWindow, snapshot mongov3.ClusterSnapshot, credentials *mongov3.MongoUser) diag.Diagnostics {
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

	// The API never returns the password; preserve it from state/plan.
	if credentials != nil {
		var existingPassword types.String
		if model.Credentials != nil {
			existingPassword = model.Credentials.Password
		}
		model.Credentials = &credentialsModel{
			Username: types.StringValue(credentials.Username),
			Password: existingPassword,
			Database: types.StringValue(credentials.Database),
		}
	}

	return diagnostics
}
