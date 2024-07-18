package s3

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/cenkalti/backoff/v4"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils/constant"

	s3 "github.com/ionos-cloud/sdk-go-s3"
)

var (
	_ resource.ResourceWithImportState = (*bucketPolicyResource)(nil)
	_ resource.ResourceWithConfigure   = (*bucketPolicyResource)(nil)
)

var errBucketPolicyNotFound = errors.New("s3 bucket policy not found")

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

func (m bucketPolicyStatementModel) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"sid":       types.StringType,
		"effect":    types.StringType,
		"actions":   types.ListType{ElemType: types.StringType},
		"resources": types.ListType{ElemType: types.StringType},
		"principal": types.ListType{ElemType: types.StringType},
		"condition": types.ObjectType{AttrTypes: bucketPolicyStatementConditionModel{}.AttributeTypes()},
	}
}

type bucketPolicyStatementConditionModel struct {
	IPs             types.List   `tfsdk:"ip_addresses"`
	ExcludedIPs     types.List   `tfsdk:"excluded_ip_addresses"`
	DateGreaterThan types.String `tfsdk:"date_greater_than"`
	DateLessThan    types.String `tfsdk:"date_less_than"`
}

func (m bucketPolicyStatementConditionModel) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"ip_addresses":          types.ListType{ElemType: types.StringType},
		"excluded_ip_addresses": types.ListType{ElemType: types.StringType},
		"date_greater_than":     types.StringType,
		"date_less_than":        types.StringType,
	}
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
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
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
							Validators:  bucketPolicyStatementActionValidators(),
						},
						"effect": schema.StringAttribute{
							Description: "The outcome when the user requests a particular action.",
							Required:    true,
							Validators: []validator.String{
								stringvalidator.OneOf("Allow", "Deny"),
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
				Required: true,
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

	requestInput, diags := putBucketPolicyInput(ctx, &data)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}

	_, err := r.client.PolicyApi.PutBucketPolicy(ctx, data.BucketName.ValueString()).BucketPolicy(requestInput).Execute()
	if err != nil {
		resp.Diagnostics.AddError("failed to create bucket policy", err.Error())
		return
	}

	// Ensure policy is created
	err = backoff.Retry(func() error {
		if _, retryErr := getBucketPolicy(ctx, r.client, data.BucketName.ValueString()); retryErr != nil {
			if errors.Is(retryErr, errBucketPolicyNotFound) {
				return retryErr
			}
			return backoff.Permanent(fmt.Errorf("failed to check if bucket policy exists: %w", retryErr))
		}
		return nil
	}, backoff.NewExponentialBackOff(backoff.WithMaxElapsedTime(utils.DefaultTimeout)))

	if err != nil {
		resp.Diagnostics.AddError("failed to create bucket policy", fmt.Sprintf("error verifying bucket policy creation: %s", err))
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
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

	policy, err := getBucketPolicy(ctx, r.client, data.BucketName.ValueString())
	if err != nil {
		if errors.Is(err, errBucketPolicyNotFound) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Failed to read bucket policy", err.Error())
		return
	}

	resp.Diagnostics.Append(setBucketPolicyData(ctx, policy, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

}

// ImportState imports the state of the bucket policy.
func (r *bucketPolicyResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("bucket_policy"), req, resp)
}

// Update updates the bucket policy.
func (r *bucketPolicyResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data bucketPolicyModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	requestInput, diags := putBucketPolicyInput(ctx, &data)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}

	_, err := r.client.PolicyApi.PutBucketPolicy(ctx, data.BucketName.ValueString()).BucketPolicy(requestInput).Execute()
	if err != nil {
		resp.Diagnostics.AddError("failed to update bucket policy", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Delete deletes the bucket.
func (r *bucketPolicyResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	if r.client == nil {
		resp.Diagnostics.AddError("s3 api client not configured", "The provider client is not configured")
		return
	}

	var data bucketPolicyModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if apiResponse, err := r.client.PolicyApi.DeleteBucketPolicy(ctx, data.BucketName.ValueString()).Execute(); err != nil {
		if apiResponse.HttpNotFound() {
			return
		}

		resp.Diagnostics.AddError("failed to delete bucket policy", err.Error())
		return
	}

	err := backoff.Retry(func() error {
		if _, retryErr := getBucketPolicy(ctx, r.client, data.BucketName.ValueString()); retryErr != nil {
			if errors.Is(retryErr, errBucketPolicyNotFound) {
				return nil
			}
			return backoff.Permanent(fmt.Errorf("failed to check if bucket policy is deleted: %w", retryErr))
		}
		return nil
	}, backoff.NewExponentialBackOff(backoff.WithMaxElapsedTime(utils.DefaultTimeout)))

	if err != nil {
		resp.Diagnostics.AddError("failed to delete bucket policy", fmt.Sprintf("error verifying bucket policy deletion: %s", err))
		return
	}

}

func putBucketPolicyInput(ctx context.Context, policyModel *bucketPolicyModel) (s3.BucketPolicy, diag.Diagnostics) {
	var diags diag.Diagnostics

	policy := s3.BucketPolicy{
		Id:        policyModel.ID.ValueStringPointer(),
		Version:   policyModel.Version.ValueStringPointer(),
		Statement: []s3.BucketPolicyStatement{},
	}
	statementsModel := make([]bucketPolicyStatementModel, len(policyModel.Statements.Elements()))
	if diags = policyModel.Statements.ElementsAs(ctx, &statementsModel, false); diags.HasError() {
		return s3.BucketPolicy{}, diags
	}

	for _, statementModel := range statementsModel {
		statement := s3.BucketPolicyStatement{
			Sid:    statementModel.SID.ValueStringPointer(),
			Effect: statementModel.Effect.ValueString(),
		}

		statement.Action = make([]string, len(statementModel.Actions.Elements()))
		if diags = statementModel.Actions.ElementsAs(ctx, &statement.Action, false); diags.HasError() {
			return s3.BucketPolicy{}, diags
		}

		statement.Resource = make([]string, len(statementModel.Resources.Elements()))
		if diags = statementModel.Resources.ElementsAs(ctx, &statement.Resource, false); diags.HasError() {
			return s3.BucketPolicy{}, diags
		}

		principalList := make([]string, len(statementModel.Principal.Elements()))
		if diags = statementModel.Principal.ElementsAs(ctx, &principalList, false); diags.HasError() {
			return s3.BucketPolicy{}, diags
		}
		principal := s3.BucketPolicyStatementPrincipal{BucketPolicyStatementPrincipalAnyOf: s3.NewBucketPolicyStatementPrincipalAnyOf(principalList)}
		statement.Principal = &principal

		if !statementModel.Condition.IsNull() {
			if statement.Condition, diags = putBucketPolicyStatementConditionInput(ctx, statementModel); diags.HasError() {
				return s3.BucketPolicy{}, diags
			}
		}

		policy.Statement = append(policy.Statement, statement)
	}

	return policy, diags
}

func putBucketPolicyStatementConditionInput(ctx context.Context, statementModel bucketPolicyStatementModel) (*s3.BucketPolicyStatementCondition, diag.Diagnostics) {
	var diags diag.Diagnostics
	var err error

	condition := s3.NewBucketPolicyStatementCondition()
	conditionModel := bucketPolicyStatementConditionModel{}
	if diags = statementModel.Condition.As(ctx, &conditionModel, basetypes.ObjectAsOptions{}); diags.HasError() {
		return nil, diags
	}

	IPs := s3.BucketPolicyStatementConditionIpAddress{}
	IPs.AwsSourceIp = make([]string, len(conditionModel.IPs.Elements()))
	if diags = conditionModel.IPs.ElementsAs(ctx, &IPs.AwsSourceIp, false); diags.HasError() {
		return nil, diags
	}
	condition.IpAddress = &IPs

	excludedIPs := s3.BucketPolicyStatementConditionIpAddress{}
	excludedIPs.AwsSourceIp = make([]string, len(conditionModel.ExcludedIPs.Elements()))
	if diags = conditionModel.ExcludedIPs.ElementsAs(ctx, &excludedIPs.AwsSourceIp, false); diags.HasError() {
		return nil, diags
	}
	condition.NotIpAddress = &excludedIPs

	if !conditionModel.DateGreaterThan.IsNull() {
		var t *s3.IonosTime
		if t, err = convertToIonosTime(conditionModel.DateGreaterThan.ValueString()); err != nil {
			diags.AddError("Error converting policy condition 'greater than' date", err.Error())
			return nil, diags
		}
		dateTime := s3.BucketPolicyStatementConditionDateGreaterThan{BucketPolicyStatementConditionDateGreaterThanOneOf: s3.NewBucketPolicyStatementConditionDateGreaterThanOneOf()}
		dateTime.BucketPolicyStatementConditionDateGreaterThanOneOf.AwsCurrentTime = t
		condition.DateGreaterThan = &dateTime
	}
	if !conditionModel.DateLessThan.IsNull() {
		var t *s3.IonosTime
		if t, err = convertToIonosTime(conditionModel.DateLessThan.ValueString()); err != nil {
			diags.AddError("Error converting policy condition 'less than' date", err.Error())
			return nil, diags
		}
		dateTime := s3.BucketPolicyStatementConditionDateLessThan{BucketPolicyStatementConditionDateGreaterThanOneOf: s3.NewBucketPolicyStatementConditionDateGreaterThanOneOf()}
		dateTime.BucketPolicyStatementConditionDateGreaterThanOneOf.AwsCurrentTime = t
		condition.DateLessThan = &dateTime
	}

	return condition, diags
}

func getBucketPolicy(ctx context.Context, client *s3.APIClient, bucketName string) (*s3.BucketPolicy, error) {
	policy, apiResponse, err := client.PolicyApi.GetBucketPolicy(ctx, bucketName).Execute()
	if err != nil {
		if apiResponse.HttpNotFound() {
			return nil, errBucketPolicyNotFound
		}
		return nil, err
	}
	return policy, nil
}

func setBucketPolicyData(ctx context.Context, policy *s3.BucketPolicy, data *bucketPolicyModel) diag.Diagnostics {
	data.ID = types.StringPointerValue(policy.Id)
	data.Version = types.StringPointerValue(policy.Version)

	var statementsModel []bucketPolicyStatementModel
	var diags diag.Diagnostics
	for _, statement := range policy.Statement {
		statementModel := bucketPolicyStatementModel{}
		statementModel.SID = types.StringPointerValue(statement.Sid)
		statementModel.Effect = types.StringValue(statement.Effect)
		if statementModel.Actions, diags = types.ListValueFrom(ctx, types.StringType, statement.Action); diags.HasError() {
			return diags
		}
		if statementModel.Resources, diags = types.ListValueFrom(ctx, types.StringType, statement.Resource); diags.HasError() {
			return diags
		}
		if statement.Principal != nil {
			if statementModel.Principal, diags = types.ListValueFrom(ctx, types.StringType, statement.Principal.BucketPolicyStatementPrincipalAnyOf.AWS); diags.HasError() {
				return diags
			}
		}
		setBucketPolicyStatementConditionData(ctx, statement.Condition, &statementModel.Condition)

		statementsModel = append(statementsModel, statementModel)
	}

	data.Statements, diags = types.ListValueFrom(ctx, types.ObjectType{AttrTypes: bucketPolicyStatementModel{}.AttributeTypes()}, statementsModel)
	return diags
}

func setBucketPolicyStatementConditionData(ctx context.Context, condition *s3.BucketPolicyStatementCondition, model *types.Object) diag.Diagnostics {
	var conditionModel bucketPolicyStatementConditionModel
	var diags diag.Diagnostics
	if condition == nil {
		*model = types.ObjectNull(conditionModel.AttributeTypes())
		return diags
	}

	if condition.IpAddress != nil {
		if conditionModel.IPs, diags = types.ListValueFrom(ctx, types.StringType, condition.IpAddress.AwsSourceIp); diags.HasError() {
			return diags
		}
	}
	if condition.NotIpAddress != nil {
		if conditionModel.ExcludedIPs, diags = types.ListValueFrom(ctx, types.StringType, condition.NotIpAddress.AwsSourceIp); diags.HasError() {
			return diags
		}
	}

	if condition.DateGreaterThan != nil &&
		condition.DateGreaterThan.BucketPolicyStatementConditionDateGreaterThanOneOf != nil &&
		condition.DateGreaterThan.BucketPolicyStatementConditionDateGreaterThanOneOf.AwsCurrentTime != nil {
		conditionModel.DateGreaterThan = types.StringValue(condition.DateGreaterThan.BucketPolicyStatementConditionDateGreaterThanOneOf.AwsCurrentTime.Format(constant.DatetimeZLayout))
	}

	if condition.DateLessThan != nil &&
		condition.DateLessThan.BucketPolicyStatementConditionDateGreaterThanOneOf != nil &&
		condition.DateLessThan.BucketPolicyStatementConditionDateGreaterThanOneOf.AwsCurrentTime != nil {
		conditionModel.DateLessThan = types.StringValue(condition.DateLessThan.BucketPolicyStatementConditionDateGreaterThanOneOf.AwsCurrentTime.Format(constant.DatetimeZLayout))
	}

	*model, diags = types.ObjectValueFrom(ctx, conditionModel.AttributeTypes(), conditionModel)
	return diags
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
