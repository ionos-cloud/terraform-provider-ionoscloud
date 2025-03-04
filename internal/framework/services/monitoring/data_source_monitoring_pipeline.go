package monitoring

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-validators/datasourcevalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"

	monitoringSDK "github.com/ionos-cloud/sdk-go-bundle/products/monitoring/v2"

	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/bundleclient"
	monitoringService "github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/monitoring"
)

var _ datasource.DataSourceWithConfigure = (*pipelineDataSource)(nil)
var _ datasource.DataSourceWithConfigValidators = (*pipelineDataSource)(nil)

type pipelineDataSource struct {
	client *monitoringService.Client
}

type pipelineDataSourceModel struct {
	ID              types.String `tfsdk:"id"`
	Name            types.String `tfsdk:"name"`
	GrafanaEndpoint types.String `tfsdk:"grafana_endpoint"`
	HTTPEndpoint    types.String `tfsdk:"http_endpoint"`
	Location        types.String `tfsdk:"location"`
}

// NewPipelineDataSource creates a new pipeline data source.
func NewPipelineDataSource() datasource.DataSource {
	return &pipelineDataSource{}
}

// Metadata returns the metadata for the data source.
func (d *pipelineDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_monitoring_pipeline"
}

// Configure configures the data source.
func (d *pipelineDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	// Add a nil check when handling ProviderData because Terraform
	// sets that data after it calls the ConfigureProvider RPC.
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

	d.client = clientBundle.MonitoringClient
}

func (d *pipelineDataSource) ConfigValidators(ctx context.Context) []datasource.ConfigValidator {
	return []datasource.ConfigValidator{
		datasourcevalidator.ExactlyOneOf(
			path.MatchRoot("name"),
			path.MatchRoot("id"),
		),
	}
}

// Schema returns the schema for the data source.
func (d *pipelineDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The ID of the pipeline",
				Optional:    true,
				Computed:    true,
			},
			"name": schema.StringAttribute{
				Description: "The name of the pipeline",
				Optional:    true,
				Computed:    true,
			},
			"location": schema.StringAttribute{
				Description: "The location of the pipeline",
				Optional:    true,
			},
			"grafana_endpoint": schema.StringAttribute{
				Computed:    true,
				Description: "The endpoint of the Grafana instance",
			},
			"http_endpoint": schema.StringAttribute{
				Computed:    true,
				Description: "The HTTP endpoint of the monitoring instance",
			},
		},
	}
}

func (d *pipelineDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	if d.client == nil {
		resp.Diagnostics.AddError("Unconfigured Monitoring API client", "Expected configured Monitoring client. Please report this issue to the provider developers.")
		return
	}

	var data pipelineDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	pipelineID := data.ID.ValueString()
	pipelineName := data.Name.ValueString()
	location := data.Location.ValueString()

	var pipeline monitoringSDK.PipelineRead
	var pipelines []monitoringSDK.PipelineRead
	var err error

	// Search using the ID
	if pipelineID != "" {
		pipeline, _, err = d.client.GetPipelineByID(ctx, pipelineID, location)

		if err != nil {
			resp.Diagnostics.AddError("failed to get Monitoring pipeline", err.Error())
			return
		}
	}

	// Search using the name
	if pipelineName != "" {
		// Retrieve ALL pipelines.
		retrievedPipelines, _, err := d.client.GetPipelines(ctx, location)
		if err != nil {
			resp.Diagnostics.AddError(fmt.Sprintf("failed to get Monitoring pipelines from location: %s", location), err.Error())
			return
		}

		// Based on the provided name, build a list of pipelines.
		for _, p := range retrievedPipelines {
			if pipelineName == p.Properties.Name {
				pipelines = append(pipelines, p)
			}
		}

		if len(pipelines) > 1 {
			resp.Diagnostics.AddError("multiple Monitoring pipelines found with the same name", "Please search using the pipeline ID instead")
			return
		}

		if len(pipelines) == 0 {
			resp.Diagnostics.AddError(fmt.Sprintf("no Monitoring pipeline found with the specified name: %s in location: %s ", pipelineName, location), "Please make sure that the name and location are correct, or search using the pipeline ID instead")
			return
		}

		pipeline = pipelines[0]
	}

	data.ID = types.StringValue(pipeline.Id)
	data.Name = types.StringValue(pipeline.Properties.Name)
	data.GrafanaEndpoint = types.StringValue(pipeline.Metadata.GrafanaEndpoint)
	data.HTTPEndpoint = types.StringValue(pipeline.Metadata.HttpEndpoint)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
