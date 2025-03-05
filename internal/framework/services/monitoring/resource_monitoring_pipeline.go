package monitoring

import (
	"context"
	"fmt"
	"strings"

	"github.com/cenkalti/backoff/v4"
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	monitoringSDK "github.com/ionos-cloud/sdk-go-bundle/products/monitoring/v2"

	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services"
	monitoringService "github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/monitoring"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils"
)

var (
	_ resource.ResourceWithImportState = (*pipelineResource)(nil)
	_ resource.ResourceWithConfigure   = (*pipelineResource)(nil)
)

type pipelineResource struct {
	client *monitoringService.Client
}

type pipelineResourceModel struct {
	ID              types.String   `tfsdk:"id"`
	Name            types.String   `tfsdk:"name"`
	GrafanaEndpoint types.String   `tfsdk:"grafana_endpoint"`
	HTTPEndpoint    types.String   `tfsdk:"http_endpoint"`
	Key             types.String   `tfsdk:"key"`
	Location        types.String   `tfsdk:"location"`
	Timeouts        timeouts.Value `tfsdk:"timeouts"`
}

// NewPipelineResource creates a new resource for the pipeline resource.
func NewPipelineResource() resource.Resource {
	return &pipelineResource{}
}

// Metadata returns the metadata for the pipeline resource.
func (r *pipelineResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_monitoring_pipeline"
}

// Schema returns the schema for the pipeline resource.
func (r *pipelineResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "The ID of the pipeline",
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "The name of the pipeline",
			},
			"grafana_endpoint": schema.StringAttribute{
				Computed:    true,
				Description: "The endpoint of the Grafana instance",
			},
			"http_endpoint": schema.StringAttribute{
				Computed:    true,
				Description: "The HTTP endpoint of the monitoring instance",
			},
			"key": schema.StringAttribute{
				Computed:    true,
				Sensitive:   true,
				Description: "The authentication key of the monitoring instance",
			},
			"location": schema.StringAttribute{
				Optional: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Description: "The location of the pipeline",
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

// Configure configures the pipeline resource.
func (r *pipelineResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Add a nil check when handling ProviderData because Terraform
	// sets that data after it calls the ConfigureProvider RPC.
	if req.ProviderData == nil {
		return
	}

	clientBundle, ok := req.ProviderData.(*services.SdkBundle)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *services.SdkBundle, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	r.client = clientBundle.MonitoringClient
}

func (r *pipelineResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data *pipelineResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	createReq := monitoringSDK.PipelineCreate{
		Properties: monitoringSDK.Pipeline{
			Name: data.Name.ValueString(),
		},
	}
	location := data.Location.ValueString()
	createTimeout, diags := data.Timeouts.Create(ctx, utils.DefaultTimeout)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	pipelineResponse, _, err := r.client.CreatePipeline(ctx, createReq, location)
	if err != nil {
		resp.Diagnostics.AddError("failed to create Monitoring pipeline", err.Error())
		return
	}
	pipelineID := pipelineResponse.Id
	key := pipelineResponse.Metadata.Key

	err = backoff.Retry(func() error {
		return r.client.IsPipelineReady(ctx, pipelineID, location)
	}, backoff.NewExponentialBackOff(backoff.WithMaxElapsedTime(createTimeout)))
	if err != nil {
		resp.Diagnostics.AddError("error occurred while waiting for the Monitoring pipeline to become available", err.Error())
	}

	// Make another `GET` request after the pipeline becomes 'AVAILABLE' in order to retrieve some
	// attributes that are not set in the `POST` response.
	retrievedPipeline, _, err := r.client.GetPipelineByID(ctx, pipelineID, location)
	if err != nil {
		resp.Diagnostics.AddError("error while fetching Monitoring pipeline after creation", (fmt.Errorf("pipeline ID: %v, error: %w", pipelineID, err)).Error())
		return
	}

	data.ID = types.StringValue(pipelineID)
	data.Key = types.StringValue(key)
	data.Name = types.StringValue(retrievedPipeline.Properties.Name)
	data.GrafanaEndpoint = types.StringValue(retrievedPipeline.Metadata.GrafanaEndpoint)
	data.HTTPEndpoint = types.StringValue(retrievedPipeline.Metadata.HttpEndpoint)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *pipelineResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	if r.client == nil {
		resp.Diagnostics.AddError("Unconfigured Monitoring API client", "Expected configured Monitoring client. Please report this issue to the provider developers.")
		return
	}

	var data pipelineResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	pipelineID := data.ID.ValueString()
	location := data.Location.ValueString()

	pipeline, apiResponse, err := r.client.GetPipelineByID(ctx, pipelineID, location)
	if err != nil {
		if apiResponse.HttpNotFound() {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("error while fetching Monitoring pipeline", (fmt.Errorf("pipeline ID: %v, error: %w", pipelineID, err)).Error())
		return
	}

	data.Name = types.StringValue(pipeline.Properties.Name)
	data.GrafanaEndpoint = types.StringValue(pipeline.Metadata.GrafanaEndpoint)
	data.HTTPEndpoint = types.StringValue(pipeline.Metadata.HttpEndpoint)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *pipelineResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data pipelineResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	pipelineID := data.ID.ValueString()
	location := data.Location.ValueString()

	deleteTimeout, diags := data.Timeouts.Delete(ctx, utils.DefaultTimeout)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	_, err := r.client.DeletePipeline(ctx, pipelineID, location)
	if err != nil {
		resp.Diagnostics.AddError("error occurred while deleting Monitoring pipeline", (fmt.Errorf("pipeline ID: %v, error: %w", pipelineID, err)).Error())
		return
	}

	err = backoff.Retry(func() error {
		return r.client.IsPipelineDeleted(ctx, pipelineID, location)
	}, backoff.NewExponentialBackOff(backoff.WithMaxElapsedTime(deleteTimeout)))
	if err != nil {
		resp.Diagnostics.AddError("error occurred while waiting for the Monitoring pipeline to be deleted", (fmt.Errorf("pipeline ID: %v, error: %w", pipelineID, err)).Error())
		return
	}
}

func (r *pipelineResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state *pipelineResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	pipelineID := state.ID.ValueString()
	location := plan.Location.ValueString()
	updateReq := monitoringSDK.PipelineEnsure{
		Properties: monitoringSDK.Pipeline{
			Name: plan.Name.ValueString(),
		},
	}

	updateTimeout, diags := plan.Timeouts.Update(ctx, utils.DefaultTimeout)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	pipelineResponse, _, err := r.client.UpdatePipeline(ctx, updateReq, pipelineID, location)
	if err != nil {
		resp.Diagnostics.AddError("error while updating Monitoring pipeline", (fmt.Errorf("pipeline ID: %v, error: %w", pipelineID, err)).Error())
		return
	}
	err = backoff.Retry(func() error {
		return r.client.IsPipelineReady(ctx, pipelineID, location)
	}, backoff.NewExponentialBackOff(backoff.WithMaxElapsedTime(updateTimeout)))
	if err != nil {
		resp.Diagnostics.AddError("error while waiting for the Monitoring pipeline to become AVAILABLE after update", (fmt.Errorf("pipeline ID: %v, error: %w", pipelineID, err)).Error())
		return
	}

	state.Name = types.StringValue(pipelineResponse.Properties.Name)

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *pipelineResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	parts := strings.Split(req.ID, ":")
	if len(parts) != 2 {
		resp.Diagnostics.AddError(
			"Unexpected Import Identifier",
			fmt.Sprintf("Expected import identifier with format: '<location>:<pipeline_id>. Got: %q", req.ID),
		)
		return
	}
	location := parts[0]
	pipelineID := parts[1]

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("location"), location)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), pipelineID)...)
}
