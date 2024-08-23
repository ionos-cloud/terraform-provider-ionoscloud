package s3

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
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
	_ resource.ResourceWithImportState = (*bucketCorsConfiguration)(nil)
	_ resource.ResourceWithConfigure   = (*bucketCorsConfiguration)(nil)
)

type bucketCorsConfiguration struct {
	client *s3.APIClient
}

type bucketCorsConfigurationModel struct {
	Bucket types.String `tfsdk:"bucket"`
	Cors   []corsRule   `tfsdk:"cors_rule"`
}

type corsRule struct {
	ID             types.Int64    `tfsdk:"id"`
	AllowedHeaders []types.String `tfsdk:"allowed_headers"`
	AllowedMethods []types.String `tfsdk:"allowed_methods"`
	AllowedOrigins []types.String `tfsdk:"allowed_origins"`
	ExposeHeaders  []types.String `tfsdk:"expose_headers"`
	MaxAgeSeconds  types.Int64    `tfsdk:"max_age_seconds"`
}

// NewBucketCorsConfigurationResource creates a new resource for the bucket CORS configuration resource.
func NewBucketCorsConfigurationResource() resource.Resource {
	return &bucketCorsConfiguration{}
}

// Metadata returns the metadata for the bucket CORS configuration resource.
func (r *bucketCorsConfiguration) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_s3_bucket_cors_configuration"
}

// Schema returns the schema for the bucket CORS configuration resource.
func (r *bucketCorsConfiguration) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
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
			"cors_rule": schema.ListNestedBlock{
				Description: "A configuration for Cross-Origin Resource Sharing (CORS).",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.Int64Attribute{
							Optional:    true,
							Description: "Container for the Contract Number of the owner.",
						},
						"allowed_headers": schema.SetAttribute{
							ElementType: types.StringType,
							Optional:    true,
							Description: "Specifies which headers are allowed in a preflight OPTIONS request through the Access-Control-Request-Headers header.",
						},
						"allowed_methods": schema.SetAttribute{
							ElementType: types.StringType,
							Required:    true,
							Description: "An HTTP method that you allow the origin to execute. Valid values are GET, PUT, HEAD, POST, DELETE.",
							Validators: []validator.Set{
								setvalidator.SizeAtLeast(1),
								setvalidator.ValueStringsAre(stringvalidator.OneOf("GET", "PUT", "HEAD", "POST", "DELETE")),
							},
						},
						"allowed_origins": schema.SetAttribute{
							ElementType: types.StringType,
							Required:    true,
							Description: "One or more origins you want customers to be able to access the bucket from.",
							Validators: []validator.Set{
								setvalidator.SizeAtLeast(1),
							},
						},
						"expose_headers": schema.SetAttribute{
							ElementType: types.StringType,
							Optional:    true,
							Description: "One or more headers in the response that you want customers to be able to access from their applications.",
						},
						"max_age_seconds": schema.Int64Attribute{
							Optional:    true,
							Description: "The time in seconds that your browser is to cache the preflight response for the specified resource.",
						},
					},
				},
			},
		},
	}
}

// Configure configures the bucket CORS configuration resource.
func (r *bucketCorsConfiguration) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

// Create creates the bucket CORS configuration.
func (r *bucketCorsConfiguration) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	if r.client == nil {
		resp.Diagnostics.AddError("API client not configured", "The provider client is not configured")
		return
	}

	var data *bucketCorsConfigurationModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	_, err := r.client.CORSApi.PutBucketCors(ctx, data.Bucket.ValueString()).PutBucketCorsRequest(buildBucketCorsConfigurationFromModel(data)).Execute()
	if err != nil {
		resp.Diagnostics.AddError("Failed to create resource", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Read reads the bucket CORS configuration.
func (r *bucketCorsConfiguration) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	if r.client == nil {
		resp.Diagnostics.AddError("API client not configured", "The provider client is not configured")
		return
	}

	var data *bucketCorsConfigurationModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	output, apiResponse, err := r.client.CORSApi.GetBucketCors(ctx, data.Bucket.ValueString()).Execute()
	if apiResponse.HttpNotFound() {
		resp.State.RemoveResource(ctx)
		return
	}

	if err != nil {
		resp.Diagnostics.AddError("Failed to read resource", err.Error())
		return
	}

	data = buildBucketCorsConfigurationModelFromAPIResponse(output, data)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// ImportState imports the state of the bucket CORS configuration.
func (r *bucketCorsConfiguration) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("bucket"), req, resp)
}

// Update updates the bucket CORS configuration.
func (r *bucketCorsConfiguration) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	if r.client == nil {
		resp.Diagnostics.AddError("API client not configured", "The provider client is not configured")
		return
	}

	var data *bucketCorsConfigurationModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	_, err := r.client.CORSApi.PutBucketCors(ctx, data.Bucket.ValueString()).PutBucketCorsRequest(buildBucketCorsConfigurationFromModel(data)).Execute()
	if err != nil {
		resp.Diagnostics.AddError("Failed to update resource", err.Error())
		return
	}

	output, _, err := r.client.CORSApi.GetBucketCors(ctx, data.Bucket.ValueString()).Execute()
	if err != nil {
		resp.Diagnostics.AddError("Failed to read resource", err.Error())
		return
	}

	data = buildBucketCorsConfigurationModelFromAPIResponse(output, data)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Delete deletes the bucket CORS configuration.
func (r *bucketCorsConfiguration) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	if r.client == nil {
		resp.Diagnostics.AddError("API client not configured", "The provider client is not configured")
		return
	}

	var data *bucketCorsConfigurationModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	_, err := r.client.CORSApi.DeleteBucketCors(ctx, data.Bucket.ValueString()).Execute()
	if err != nil {
		resp.Diagnostics.AddError("Failed to delete resource", err.Error())
		return
	}
}

func buildBucketCorsConfigurationModelFromAPIResponse(output *s3.GetBucketCorsOutput, data *bucketCorsConfigurationModel) *bucketCorsConfigurationModel {
	data.Cors = buildCorsRulesFromAPIResponse(output.CORSRules)
	return data
}

func buildCorsRulesFromAPIResponse(rules *[]s3.CORSRule) []corsRule {
	if rules == nil {
		return nil
	}

	result := make([]corsRule, 0, len(*rules))
	for _, r := range *rules {
		result = append(result, corsRule{
			ID:             types.Int64PointerValue(toInt64(r.ID)),
			AllowedHeaders: toTFStrings(r.AllowedHeaders),
			AllowedMethods: toTFStrings(r.AllowedMethods),
			AllowedOrigins: toTFStrings(r.AllowedOrigins),
			ExposeHeaders:  toTFStrings(r.ExposeHeaders),
			MaxAgeSeconds:  types.Int64PointerValue(toInt64(r.MaxAgeSeconds)),
		})
	}

	return result
}
func buildBucketCorsConfigurationFromModel(data *bucketCorsConfigurationModel) s3.PutBucketCorsRequest {
	return s3.PutBucketCorsRequest{
		CORSRules: buildCorsRulesFromModel(data.Cors),
	}
}

func buildCorsRulesFromModel(rules []corsRule) *[]s3.CORSRule {
	result := make([]s3.CORSRule, 0, len(rules))
	for _, r := range rules {
		result = append(result, s3.CORSRule{
			ID:             toInt32(r.ID.ValueInt64Pointer()),
			AllowedHeaders: toStrings(r.AllowedHeaders),
			AllowedMethods: toStrings(r.AllowedMethods),
			AllowedOrigins: toStrings(r.AllowedOrigins),
			ExposeHeaders:  toStrings(r.ExposeHeaders),
			MaxAgeSeconds:  toInt32(r.MaxAgeSeconds.ValueInt64Pointer()),
		})
	}

	return &result
}

func toStrings(s []types.String) *[]string {
	if len(s) == 0 {
		return nil
	}

	result := make([]string, 0, len(s))
	for _, v := range s {
		result = append(result, v.ValueString())
	}

	return &result
}

func toTFStrings(s *[]string) []types.String {
	if s == nil {
		return nil
	}

	result := make([]types.String, 0, len(*s))
	for _, v := range *s {
		result = append(result, types.StringValue(v))
	}

	return result
}
