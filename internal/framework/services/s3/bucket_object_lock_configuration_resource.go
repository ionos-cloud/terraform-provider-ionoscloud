package s3

import (
	"context"
	"crypto/md5" // nolint:gosec
	"encoding/hex"
	"encoding/xml"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/resourcevalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	s3 "github.com/ionos-cloud/sdk-go-s3"
)

var (
	_ resource.ResourceWithImportState      = (*objectLockConfiguration)(nil)
	_ resource.ResourceWithConfigure        = (*objectLockConfiguration)(nil)
	_ resource.ResourceWithConfigValidators = (*objectLockConfiguration)(nil)
)

type objectLockConfiguration struct {
	client *s3.APIClient
}

type objectLockConfigurationModel struct {
	Bucket            types.String `tfsdk:"bucket"`
	ObjectLockEnabled types.String `tfsdk:"object_lock_enabled"`
	Rule              *rule        `tfsdk:"rule"`
}

type rule struct {
	DefaultRetention *defaultRetention `tfsdk:"default_retention"`
}

type defaultRetention struct {
	Mode  types.String `tfsdk:"mode"`
	Days  types.Int64  `tfsdk:"days"`
	Years types.Int64  `tfsdk:"years"`
}

// NewObjectLockConfigurationResource creates a new resource for the bucket object lock configuration resource.
func NewObjectLockConfigurationResource() resource.Resource {
	return &objectLockConfiguration{}
}

// Metadata returns the metadata for the bucket object lock configuration resource.
func (r *objectLockConfiguration) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_s3_bucket_object_lock_configuration"
}

// Schema returns the schema for the bucket object lock configuration resource.
func (r *objectLockConfiguration) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"bucket": schema.StringAttribute{
				Required:    true,
				Description: "The name of the bucket.",
				Validators:  []validator.String{stringvalidator.LengthBetween(3, 63)},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"object_lock_enabled": schema.StringAttribute{
				Description: "Specifies whether Object Lock is enabled for the bucket.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Default: stringdefault.StaticString("Enabled"),
			},
		},
		Blocks: map[string]schema.Block{
			"rule": schema.SingleNestedBlock{
				Blocks: map[string]schema.Block{
					"default_retention": schema.SingleNestedBlock{
						Attributes: map[string]schema.Attribute{
							"mode": schema.StringAttribute{
								Optional:   true,
								Validators: []validator.String{stringvalidator.OneOf("GOVERNANCE", "COMPLIANCE")},
							},
							"days": schema.Int64Attribute{
								Optional:   true,
								Validators: []validator.Int64{int64validator.AtLeast(1)},
							},
							"years": schema.Int64Attribute{
								Optional:   true,
								Validators: []validator.Int64{int64validator.AtLeast(1)},
							},
						},
					},
				},
			},
		},
	}
}

// ConfigValidators returns the config validators for the bucket object lock configuration resource.
func (r *objectLockConfiguration) ConfigValidators(_ context.Context) []resource.ConfigValidator {
	return []resource.ConfigValidator{
		resourcevalidator.Conflicting(
			path.MatchRoot("rule").AtName("default_retention").AtName("days"),
			path.MatchRoot("rule").AtName("default_retention").AtName("years"),
		),
	}
}

// Configure configures the bucket object lock configuration resource.
func (r *objectLockConfiguration) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

// Create creates the bucket object lock configuration resource.
func (r *objectLockConfiguration) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	if r.client == nil {
		resp.Diagnostics.AddError("API client not configured", "The provider client is not configured")
		return
	}

	var data *objectLockConfigurationModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	input := buildObjectLockConfigurationFromModel(data)
	md5Sum, err := generateMD5Sum(input)
	if err != nil {
		resp.Diagnostics.AddError("Failed to generate MD5 sum", err.Error())
		return
	}
	_, err = r.client.ObjectLockApi.PutObjectLockConfiguration(ctx, data.Bucket.ValueString()).PutObjectLockConfigurationRequest(input).ContentMD5(md5Sum).Execute()
	if err != nil {
		resp.Diagnostics.AddError("Failed to create resource", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Read reads the bucket object lock configuration resource.
func (r *objectLockConfiguration) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	if r.client == nil {
		resp.Diagnostics.AddError("API client not configured", "The provider client is not configured")
		return
	}

	var data *objectLockConfigurationModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	output, apiResponse, err := r.client.ObjectLockApi.GetObjectLockConfiguration(ctx, data.Bucket.ValueString()).Execute()
	if apiResponse.HttpNotFound() {
		resp.State.RemoveResource(ctx)
		return
	}
	if err != nil {
		resp.Diagnostics.AddError("Failed to read resource", err.Error())
		return
	}

	data = buildObjectLockConfigurationModelFromAPIResponse(output, data)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// ImportState imports the state for the bucket object lock configuration resource.
func (r *objectLockConfiguration) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("bucket"), req, resp)
}

// Update updates the bucket object lock configuration resource.
func (r *objectLockConfiguration) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	if r.client == nil {
		resp.Diagnostics.AddError("API client not configured", "The provider client is not configured")
		return
	}

	var data *objectLockConfigurationModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	input := buildObjectLockConfigurationFromModel(data)
	md5Sum, err := generateMD5Sum(input)
	if err != nil {
		resp.Diagnostics.AddError("Failed to generate MD5 sum", err.Error())
		return
	}
	_, err = r.client.ObjectLockApi.PutObjectLockConfiguration(ctx, data.Bucket.ValueString()).PutObjectLockConfigurationRequest(input).ContentMD5(md5Sum).Execute()
	if err != nil {
		resp.Diagnostics.AddError("Failed to create resource", err.Error())
		return
	}

	output, apiResponse, err := r.client.ObjectLockApi.GetObjectLockConfiguration(ctx, data.Bucket.ValueString()).Execute()
	if apiResponse.HttpNotFound() {
		resp.State.RemoveResource(ctx)
		return
	}

	if err != nil {
		resp.Diagnostics.AddError("Failed to read resource", err.Error())
		return
	}

	data = buildObjectLockConfigurationModelFromAPIResponse(output, data)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Delete deletes the bucket object lock configuration resource.
func (r *objectLockConfiguration) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	if r.client == nil {
		resp.Diagnostics.AddError("API client not configured", "The provider client is not configured")
		return
	}

	var data *objectLockConfigurationModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Custom delete logic
}

func buildObjectLockConfigurationModelFromAPIResponse(output *s3.GetObjectLockConfigurationOutput, data *objectLockConfigurationModel) *objectLockConfigurationModel {
	built := &objectLockConfigurationModel{
		Bucket:            data.Bucket,
		ObjectLockEnabled: types.StringPointerValue(output.ObjectLockEnabled),
	}
	if output.Rule != nil {
		built.Rule = &rule{
			DefaultRetention: &defaultRetention{
				Mode:  types.StringPointerValue(output.Rule.DefaultRetention.Mode),
				Days:  types.Int64PointerValue(toInt64(output.Rule.DefaultRetention.Days)),
				Years: types.Int64PointerValue(toInt64(output.Rule.DefaultRetention.Years)),
			},
		}
	}

	return built
}

func buildObjectLockConfigurationFromModel(data *objectLockConfigurationModel) s3.PutObjectLockConfigurationRequest {
	req := s3.PutObjectLockConfigurationRequest{
		ObjectLockEnabled: data.ObjectLockEnabled.ValueStringPointer(),
		Rule: &s3.PutObjectLockConfigurationRequestRule{
			DefaultRetention: &s3.DefaultRetention{
				Mode:  data.Rule.DefaultRetention.Mode.ValueStringPointer(),
				Days:  toInt32(data.Rule.DefaultRetention.Days.ValueInt64Pointer()),
				Years: toInt32(data.Rule.DefaultRetention.Years.ValueInt64Pointer()),
			},
		},
	}
	return req
}

func generateMD5Sum(data interface{}) (string, error) {
	// Marshal the struct to JSON
	jsonBytes, err := xml.Marshal(data)
	if err != nil {
		return "", err
	}

	// Create an MD5 hasher
	hasher := md5.New() // nolint:gosec

	// Write the JSON bytes to the hasher
	_, err = hasher.Write(jsonBytes)
	if err != nil {
		return "", err
	}

	// Compute the MD5 checksum
	md5sum := hasher.Sum(nil)

	// Convert the checksum to a hex string
	return hex.EncodeToString(md5sum), nil
}

func toInt32(i *int64) *int32 {
	if i == nil {
		return nil
	}

	v := int32(*i)
	return &v
}

func toInt64(i *int32) *int64 {
	if i == nil {
		return nil
	}

	v := int64(*i)
	return &v
}
