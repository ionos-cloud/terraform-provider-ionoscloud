package s3

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils/constant"

	s3 "github.com/ionos-cloud/sdk-go-s3"
)

var (
	_ resource.ResourceWithImportState = (*bucketPolicyResource)(nil)
	_ resource.ResourceWithConfigure   = (*bucketPolicyResource)(nil)
)

// NewBucketPolicyResource creates a new resource for the bucket resource.
func NewBucketPolicyResource() resource.Resource {
	return &bucketPolicyResource{}
}

type bucketPolicyResource struct {
	client *s3.APIClient
}

type bucketPolicyModel struct {
	BucketName types.String `tfsdk:"bucket_name"`
	ID         types.String `tfsdk:"id"`
	Version    types.String `tfsdk:"version"`
	Statements types.List   `tfsdk:"statements"`
}

type bucketPolicyStatementModel struct {
	SID       types.String `tfsdk:"sid"`
	Effect    types.String `tfsdk:"effect"`
	Actions   types.List   `tfsdk:"actions"`
	Resources types.List   `tfsdk:"resources"`
	Principal types.List   `tfsdk:"principal"`
	Condition types.Object `tfsdk:"condition"`
}

type bucketPolicyStatementConditionModel struct {
	IpAddresses         types.List   `tfsdk:"ip_addresses"`
	ExcludedIpAddresses types.List   `tfsdk:"excluded_ip_addresses"`
	DateGreaterThan     types.String `tfsdk:"date_greater_than"`
	DateLessThan        types.String `tfsdk:"date_less_than"`
}

// Metadata returns the metadata for the bucket policy resource.
func (r *bucketPolicyResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_bucket_policy" // todo: use constant here maybe
}

// Schema returns the schema for the bucket policy resource.
func (r *bucketPolicyResource) Schema(_ context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"bucket_name": schema.StringAttribute{
				Description: "Name of the S3 bucket to which this policy will be applied.",
				Required:    true,
			},
			"id": schema.StringAttribute{
				Description: "Optional identifier for the bucket policy.",
				Optional:    true,
			},
			"version": schema.StringAttribute{
				Description: "The version of the bucket policy.",
				Optional:    true,
			},
			"statements": schema.ListNestedAttribute{
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"sid": schema.StringAttribute{
							Description: "Optional identifier for the policy statement.",
							Optional:    true,
						},
						"actions": schema.ListAttribute{
							Description: "List of allowed or denied actions.",
							ElementType: types.StringType,
							Required:    true,
						},
						"effect": schema.StringAttribute{
							Description: "The outcome when the user requests a particular action.",
							Required:    true,
							Validators: []validator.String{
								stringvalidator.OneOfCaseInsensitive("allow", "deny"),
							},
						},
						"resources": schema.ListAttribute{
							Description: "The bucket or object that the policy applies to.",
							ElementType: types.StringType,
							Required:    true,
						},
						"condition": schema.SingleNestedAttribute{
							Description: "Conditions for when a policy is in effect.",
							Attributes: map[string]schema.Attribute{
								"ip_addresses": schema.ListAttribute{
									Description: "List of affected IP addresses.",
									ElementType: types.StringType,
									Optional:    true,
								},
								"excluded_ip_addresses": schema.ListAttribute{
									Description: "List of unaffected IP addresses.",
									ElementType: types.StringType,
									Optional:    true,
								},
								"date_greater_than": schema.StringAttribute{
									Description: "Minimum date time.",
									Optional:    true,
								},
								"date_less_than": schema.StringAttribute{
									Description: "Maximum date time.",
									Optional:    true,
								},
							},
							Optional: true,
						},
						"principal": schema.ListAttribute{
							Description: "Users to which the policy applies to.",
							ElementType: types.StringType,
							Optional:    true,
						},
					},
				},
				Optional: true,
			},
		},
	}
}

// Configure configures the bucket policy resource.
func (r *bucketPolicyResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*s3.APIClient)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *hashicups.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.client = client
}

// Create creates the bucket policy.
func (r *bucketPolicyResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	if r.client == nil {
		resp.Diagnostics.AddError("s3 api client not configured", "The provider client is not configured") // todo: const for this error maybe?
		return
	}

	var data bucketPolicyModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	createInput, diags := putBucketPolicyInput(ctx, &data)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}

	_, err := r.client.PolicyApi.PutBucketPolicy(ctx, data.BucketName.ValueString()).BucketPolicy(createInput).Execute()
	if err != nil {
		resp.Diagnostics.AddError("failed to create bucket", err.Error())
		return
	}

}

// Read reads the bucket policy.
func (r *bucketPolicyResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	if r.client == nil {
		resp.Diagnostics.AddError("s3 api client not configured", "The provider client is not configured")
		return
	}

	var data bucketPolicyModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

}

// ImportState imports the state of the bucket policy.
func (r *bucketPolicyResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("bucket_policy"), req, resp)
}

// Update updates the bucket policy.
func (r *bucketPolicyResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data bucketResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Nothing to update for now
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Delete deletes the bucket.
func (r *bucketPolicyResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	if r.client == nil {
		resp.Diagnostics.AddError("s3 api client not configured", "The provider client is not configured")
		return
	}

}

// read state data into request object
func putBucketPolicyInput(ctx context.Context, data *bucketPolicyModel) (s3.BucketPolicy, diag.Diagnostics) {
	policy := s3.BucketPolicy{
		Id:        data.ID.ValueStringPointer(),
		Version:   data.Version.ValueStringPointer(),
		Statement: []s3.BucketPolicyStatement{},
	}

	statements := make([]bucketPolicyStatementModel, len(data.Statements.Elements()))
	diags := data.Statements.ElementsAs(ctx, &statements, false)

	for _, s := range statements {
		statement := s3.BucketPolicyStatement{
			Sid:    s.SID.ValueStringPointer(),
			Effect: s.Effect.String(),
		}
		condition := s3.BucketPolicyStatementCondition{}
		statement.Condition = &condition

		statement.Action = make([]string, len(s.Actions.Elements()))
		diags.Append(s.Actions.ElementsAs(ctx, &statement.Action, false)...)

		statement.Resource = make([]string, len(s.Resources.Elements()))
		diags.Append(s.Resources.ElementsAs(ctx, &statement.Resource, false)...)

		c := bucketPolicyStatementConditionModel{}
		diags.Append(s.Condition.As(ctx, &c, basetypes.ObjectAsOptions{})...)

		t, err := convertToIonosTime(c.DateGreaterThan.String())
		if err == nil {
			statement.Condition.DateGreaterThan.BucketPolicyStatementConditionDateGreaterThanOneOf.AwsCurrentTime = t
		}
		t, err = convertToIonosTime(c.DateLessThan.String())
		if err == nil {
			statement.Condition.DateLessThan.BucketPolicyStatementConditionDateGreaterThanOneOf.AwsCurrentTime = t
		}

		ipAddresses := s3.BucketPolicyStatementConditionIpAddress{}
		ipAddresses.AwsSourceIp = make([]string, len(c.IpAddresses.Elements()))
		diags.Append(c.IpAddresses.ElementsAs(ctx, &ipAddresses.AwsSourceIp, false)...)
		statement.Condition.IpAddress = &ipAddresses

		excludedIpAddresses := s3.BucketPolicyStatementConditionIpAddress{}
		ipAddresses.AwsSourceIp = make([]string, len(c.IpAddresses.Elements()))
		diags.Append(c.ExcludedIpAddresses.ElementsAs(ctx, &ipAddresses.AwsSourceIp, false)...)
		statement.Condition.NotIpAddress = &excludedIpAddresses

		policy.Statement = append(policy.Statement, s3.BucketPolicyStatement{})
	}

	return policy, diags
}

// duplicated
func convertToIonosTime(targetTime string) (*s3.IonosTime, error) {
	var ionosTime s3.IonosTime
	var convertedTime time.Time
	var err error

	// targetTime might have time zone offset layout (+0000 UTC)
	if convertedTime, err = time.Parse(constant.DatetimeTZOffsetLayout, targetTime); err != nil {
		if convertedTime, err = time.Parse(constant.DatetimeZLayout, targetTime); err != nil {
			return nil, fmt.Errorf("an error occurred while converting from IonosTime string to time.Time: %w", err)
		}
	}
	ionosTime.Time = convertedTime
	return &ionosTime, nil
}
