package s3

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	s3 "github.com/ionos-cloud/sdk-go-s3"
)

var (
	_ resource.ResourceWithImportState = (*serverSideEncryptionConfiguration)(nil)
	_ resource.ResourceWithConfigure   = (*serverSideEncryptionConfiguration)(nil)
)

type serverSideEncryptionConfiguration struct {
	client *s3.APIClient
}

type serverSideEncryptionConfigurationModel struct {
	Bucket types.String `tfsdk:"bucket"`
	Rules  []sseRule    `tfsdk:"rule"`
}

type sseRule struct {
	ApplyServerSideEncryptionByDefault applyServerSideEncryptionByDefault `tfsdk:"apply_server_side_encryption_by_default"`
}

type applyServerSideEncryptionByDefault struct {
	SSEAlgorithm types.String `tfsdk:"sse_algorithm"`
}

// NewServerSideEncryptionConfigurationResource creates a new resource for the server side encryption configuration resource.
func NewServerSideEncryptionConfigurationResource() resource.Resource {
	return &serverSideEncryptionConfiguration{}
}

// Metadata returns the metadata for the server side encryption configuration resource.
func (r *serverSideEncryptionConfiguration) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_s3_bucket_server_side_encryption_configuration"
}

// Schema returns the schema for the server side encryption configuration resource.
func (r *serverSideEncryptionConfiguration) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"bucket": schema.StringAttribute{
				Description: "The name of the bucket",
				Required:    true,
				Validators:  []validator.String{stringvalidator.LengthBetween(3, 63)},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
		},
		Blocks: map[string]schema.Block{
			"rule": schema.SetNestedBlock{
				Description: "Specifies the default server-side encryption configuration.",
				NestedObject: schema.NestedBlockObject{
					Blocks: map[string]schema.Block{
						"apply_server_side_encryption_by_default": schema.SingleNestedBlock{
							Description: "Defines the default encryption settings.",
							Attributes: map[string]schema.Attribute{
								"sse_algorithm": schema.StringAttribute{
									Required:    true,
									Description: "Server-side encryption algorithm to use. Valid values are 'AES256'",
									Validators: []validator.String{
										stringvalidator.OneOf("AES256"),
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

// Configure configures the server side encryption configuration resource.
func (r *serverSideEncryptionConfiguration) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*s3.APIClient)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *s3.APIClient, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.client = client
}

// Create creates the server side encryption configuration.
func (r *serverSideEncryptionConfiguration) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	if r.client == nil {
		resp.Diagnostics.AddError("API client not configured", "The provider client is not configured")
		return
	}

	var data *serverSideEncryptionConfigurationModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	_, err := r.client.EncryptionApi.PutBucketEncryption(ctx, data.Bucket.ValueString()).PutBucketEncryptionRequest(buildServerSideEncryptionConfigurationFromModel(data)).Execute()
	if err != nil {
		resp.Diagnostics.AddError("Failed to create resource", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Read reads the server side encryption configuration.
func (r *serverSideEncryptionConfiguration) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	if r.client == nil {
		resp.Diagnostics.AddError("API client not configured", "The provider client is not configured")
		return
	}

	var data *serverSideEncryptionConfigurationModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	output, apiResponse, err := r.client.EncryptionApi.GetBucketEncryption(ctx, data.Bucket.ValueString()).Execute()
	if apiResponse.HttpNotFound() {
		resp.State.RemoveResource(ctx)
		return
	}
	if err != nil {
		resp.Diagnostics.AddError("Failed to read resource", err.Error())
		return
	}

	data = buildServerSideEncryptionConfigurationModelFromAPIResponse(output, data)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// ImportState imports the state for the server side encryption configuration.
func (r *serverSideEncryptionConfiguration) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("bucket"), req, resp)
}

// Update updates the server side encryption configuration.
func (r *serverSideEncryptionConfiguration) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	if r.client == nil {
		resp.Diagnostics.AddError("API client not configured", "The provider client is not configured")
		return
	}

	var data *serverSideEncryptionConfigurationModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	_, err := r.client.EncryptionApi.PutBucketEncryption(ctx, data.Bucket.ValueString()).PutBucketEncryptionRequest(buildServerSideEncryptionConfigurationFromModel(data)).Execute()
	if err != nil {
		resp.Diagnostics.AddError("Failed to create resource", err.Error())
		return
	}

	output, _, err := r.client.EncryptionApi.GetBucketEncryption(ctx, data.Bucket.ValueString()).Execute()
	if err != nil {
		resp.Diagnostics.AddError("Failed to read resource", err.Error())
		return
	}

	data = buildServerSideEncryptionConfigurationModelFromAPIResponse(output, data)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Delete deletes the server side encryption configuration.
func (r *serverSideEncryptionConfiguration) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	if r.client == nil {
		resp.Diagnostics.AddError("API client not configured", "The provider client is not configured")
		return
	}

	var data *serverSideEncryptionConfigurationModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	_, err := r.client.EncryptionApi.DeleteBucketEncryption(ctx, data.Bucket.ValueString()).Execute()
	if err != nil {
		resp.Diagnostics.AddError("Failed to delete resource", err.Error())
		return
	}
}

func buildServerSideEncryptionConfigurationModelFromAPIResponse(output *s3.ServerSideEncryptionConfiguration, data *serverSideEncryptionConfigurationModel) *serverSideEncryptionConfigurationModel {
	return &serverSideEncryptionConfigurationModel{
		Bucket: data.Bucket,
		Rules:  buildServerSideEncryptionRulesFromAPIResponse(output.Rules),
	}
}

func buildServerSideEncryptionRulesFromAPIResponse(data *[]s3.ServerSideEncryptionRule) []sseRule {
	if data == nil {
		return nil
	}

	rules := make([]sseRule, 0, len(*data))
	for _, r := range *data {
		if r.ApplyServerSideEncryptionByDefault == nil {
			continue
		}

		if r.ApplyServerSideEncryptionByDefault.SSEAlgorithm == nil {
			continue
		}

		rules = append(rules, sseRule{
			ApplyServerSideEncryptionByDefault: applyServerSideEncryptionByDefault{
				SSEAlgorithm: types.StringValue(string(*r.ApplyServerSideEncryptionByDefault.SSEAlgorithm)),
			},
		})
	}

	return rules
}

func buildServerSideEncryptionConfigurationFromModel(data *serverSideEncryptionConfigurationModel) s3.PutBucketEncryptionRequest {
	return s3.PutBucketEncryptionRequest{
		Rules: buildServerSideEncryptionRulesFromModel(data.Rules),
	}
}

func buildServerSideEncryptionRulesFromModel(data []sseRule) *[]s3.ServerSideEncryptionRule {
	rules := make([]s3.ServerSideEncryptionRule, 0, len(data))
	for _, r := range data {
		rules = append(rules, s3.ServerSideEncryptionRule{
			ApplyServerSideEncryptionByDefault: &s3.ServerSideEncryptionByDefault{
				SSEAlgorithm: s3.ServerSideEncryption(r.ApplyServerSideEncryptionByDefault.SSEAlgorithm.ValueString()).Ptr(),
			},
		})
	}

	return &rules
}
