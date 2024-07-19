package s3

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/cenkalti/backoff/v4"
	"github.com/hashicorp/terraform-plugin-framework-jsontypes/jsontypes"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"

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
	Bucket types.String         `tfsdk:"bucket"`
	Policy jsontypes.Normalized `tfsdk:"policy"`
}

// Intermediary structs for policy string serialization/deserialization
type bucketPolicy struct {
	ID        *string                 `json:"Id,omitempty"`
	Version   *string                 `json:"Version,omitempty"`
	Statement []bucketPolicyStatement `json:"Statement"`
}

type bucketPolicyStatement struct {
	SID       *string                         `json:"Sid,omitempty"`
	Effect    string                          `json:"Effect"`
	Action    []string                        `json:"Action"`
	Resources []string                        `json:"Resource"`
	Principal []string                        `json:"Principal"`
	Condition *bucketPolicyStatementCondition `json:"Condition,omitempty"`
}

type bucketPolicyStatementCondition struct {
	IPs             []string `json:"IpAddress,omitempty"`
	ExcludedIPs     []string `json:"NotIpAddress,omitempty"`
	DateGreaterThan *string  `json:"DateGreaterThan,omitempty"`
	DateLessThan    *string  `json:"DateLessThan,omitempty"`
}

// Metadata returns the metadata for the bucket policy resource.
func (r *bucketPolicyResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_bucket_policy" // todo: use constant here maybe
}

// Schema returns the schema for the bucket policy resource.
func (r *bucketPolicyResource) Schema(_ context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"bucket": schema.StringAttribute{
				Description: "Name of the S3 bucket to which this policy will be applied.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"policy": schema.StringAttribute{
				CustomType:  jsontypes.NormalizedType{},
				Description: "Text of the policy",
				Required:    true,
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

	requestInput, diags := putBucketPolicyInput(&data)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}

	_, err := r.client.PolicyApi.PutBucketPolicy(ctx, data.Bucket.ValueString()).BucketPolicy(requestInput).Execute()
	if err != nil {
		resp.Diagnostics.AddError("failed to create bucket policy", err.Error())
		return
	}

	// Ensure policy is created
	err = backoff.Retry(func() error {
		if _, retryErr := getBucketPolicy(ctx, r.client, data.Bucket.ValueString()); retryErr != nil {
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

	policy, err := getBucketPolicy(ctx, r.client, data.Bucket.ValueString())
	if err != nil {
		if errors.Is(err, errBucketPolicyNotFound) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Failed to read bucket policy", err.Error())
		return
	}

	resp.Diagnostics.Append(setBucketPolicyData(policy, &data)...)
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

	requestInput, diags := putBucketPolicyInput(&data)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}

	_, err := r.client.PolicyApi.PutBucketPolicy(ctx, data.Bucket.ValueString()).BucketPolicy(requestInput).Execute()
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

	if apiResponse, err := r.client.PolicyApi.DeleteBucketPolicy(ctx, data.Bucket.ValueString()).Execute(); err != nil {
		if apiResponse.HttpNotFound() {
			return
		}

		resp.Diagnostics.AddError("failed to delete bucket policy", err.Error())
		return
	}

	err := backoff.Retry(func() error {
		if _, retryErr := getBucketPolicy(ctx, r.client, data.Bucket.ValueString()); retryErr != nil {
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

func putBucketPolicyInput(policyModel *bucketPolicyModel) (s3.BucketPolicy, diag.Diagnostics) {
	var diags diag.Diagnostics
	policyInput := s3.BucketPolicy{}
	policyData := bucketPolicy{}

	// Can't unmarshal directly in the API object, need to use an intermediary
	if diags = policyModel.Policy.Unmarshal(&policyData); diags.HasError() {
		return s3.BucketPolicy{}, diags
	}

	policyInput.Id = policyData.ID
	policyInput.Version = policyData.Version
	policyInput.Statement = make([]s3.BucketPolicyStatement, 0, len(policyData.Statement))
	for _, statementData := range policyData.Statement {
		statementInput := s3.BucketPolicyStatement{}
		statementInput.Sid = statementData.SID
		statementInput.Effect = statementData.Effect
		statementInput.Resource = statementData.Resources
		statementInput.Action = statementData.Action

		principalInput := s3.BucketPolicyStatementPrincipal{BucketPolicyStatementPrincipalAnyOf: s3.NewBucketPolicyStatementPrincipalAnyOf(statementData.Principal)}
		statementInput.Principal = &principalInput

		if statementData.Condition != nil {
			statementInput.Condition = s3.NewBucketPolicyStatementCondition()

			statementInput.Condition.IpAddress = s3.NewBucketPolicyStatementConditionIpAddress()
			statementInput.Condition.IpAddress.AwsSourceIp = statementData.Condition.IPs
			statementInput.Condition.NotIpAddress = s3.NewBucketPolicyStatementConditionIpAddress()
			statementInput.Condition.NotIpAddress.AwsSourceIp = statementData.Condition.ExcludedIPs
			if statementData.Condition.DateGreaterThan != nil {
				var t *s3.IonosTime
				var err error
				if t, err = convertToIonosTime(*statementData.Condition.DateGreaterThan); err != nil {
					diags.AddError("Error converting policy condition 'greater than' date", err.Error())
					return s3.BucketPolicy{}, diags
				}
				dateGreater := s3.BucketPolicyStatementConditionDateGreaterThan{
					BucketPolicyStatementConditionDateGreaterThanOneOf: s3.NewBucketPolicyStatementConditionDateGreaterThanOneOf(),
				}
				dateGreater.BucketPolicyStatementConditionDateGreaterThanOneOf.AwsCurrentTime = t
				statementInput.Condition.DateGreaterThan = &dateGreater
			}
			if statementData.Condition.DateLessThan != nil {
				var t *s3.IonosTime
				var err error
				if t, err = convertToIonosTime(*statementData.Condition.DateLessThan); err != nil {
					diags.AddError("Error converting policy condition 'less than' date", err.Error())
					return s3.BucketPolicy{}, diags
				}
				dateLess := s3.BucketPolicyStatementConditionDateLessThan{
					BucketPolicyStatementConditionDateGreaterThanOneOf: s3.NewBucketPolicyStatementConditionDateGreaterThanOneOf(),
				}
				dateLess.BucketPolicyStatementConditionDateGreaterThanOneOf.AwsCurrentTime = t
				statementInput.Condition.DateLessThan = &dateLess
			}
		}
		policyInput.Statement = append(policyInput.Statement, statementInput)
	}

	return policyInput, diags
}

func setBucketPolicyData(policyResponse *s3.BucketPolicy, data *bucketPolicyModel) diag.Diagnostics {
	var diags diag.Diagnostics

	policyData := bucketPolicy{
		ID:      policyResponse.Id,
		Version: policyResponse.Version,
	}
	policyData.Statement = make([]bucketPolicyStatement, 0, len(policyResponse.Statement))

	for _, statementResponse := range policyResponse.Statement {
		statementData := bucketPolicyStatement{
			SID:       statementResponse.Sid,
			Effect:    statementResponse.Effect,
			Action:    statementResponse.Action,
			Resources: statementResponse.Resource,
		}
		if statementResponse.Principal != nil && statementResponse.Principal.BucketPolicyStatementPrincipalAnyOf != nil {
			statementData.Principal = statementResponse.Principal.BucketPolicyStatementPrincipalAnyOf.AWS
		}
		if statementResponse.Condition != nil {
			conditionData := bucketPolicyStatementCondition{
				IPs:         statementResponse.Condition.GetIpAddress().AwsSourceIp,
				ExcludedIPs: statementResponse.Condition.GetNotIpAddress().AwsSourceIp,
			}
			if statementResponse.Condition.DateGreaterThan != nil &&
				statementResponse.Condition.DateGreaterThan.BucketPolicyStatementConditionDateGreaterThanOneOf != nil &&
				statementResponse.Condition.DateGreaterThan.BucketPolicyStatementConditionDateGreaterThanOneOf.AwsCurrentTime != nil {
				dateString := statementResponse.Condition.DateGreaterThan.BucketPolicyStatementConditionDateGreaterThanOneOf.AwsCurrentTime.Format(constant.DatetimeZLayout)
				conditionData.DateGreaterThan = &dateString
			}
			if statementResponse.Condition.DateLessThan != nil &&
				statementResponse.Condition.DateLessThan.BucketPolicyStatementConditionDateGreaterThanOneOf != nil &&
				statementResponse.Condition.DateLessThan.BucketPolicyStatementConditionDateGreaterThanOneOf.AwsCurrentTime != nil {
				dateString := statementResponse.Condition.DateLessThan.BucketPolicyStatementConditionDateGreaterThanOneOf.AwsCurrentTime.Format(constant.DatetimeZLayout)
				conditionData.DateLessThan = &dateString
			}
			statementData.Condition = &conditionData
		}
		policyData.Statement = append(policyData.Statement, statementData)
	}
	policyString, err := json.Marshal(policyData)
	if err != nil {
		diags.AddError("Error serializing policy data", err.Error())
		return diags
	}
	data.Policy = jsontypes.NewNormalizedValue(string(policyString))
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
