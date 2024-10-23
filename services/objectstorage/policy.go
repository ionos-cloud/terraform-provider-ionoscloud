package objectstorage

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/cenkalti/backoff/v4"
	"github.com/hashicorp/terraform-plugin-framework-jsontypes/jsontypes"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	objstorage "github.com/ionos-cloud/sdk-go-s3"

	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils/constant"
)

// ErrBucketPolicyNotFound returned for 404
var ErrBucketPolicyNotFound = errors.New("bucket policy not found")

// BucketPolicyModel is used to create, update and delete a bucket policy.
type BucketPolicyModel struct {
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

// CreateBucketPolicy creates a new bucket policy.
func (c *Client) CreateBucketPolicy(ctx context.Context, data *BucketPolicyModel) error {
	input, diags := buildBucketPolicyFromModel(data)
	if diags.HasError() {
		return fmt.Errorf("error building bucket policy: %v", diags)
	}

	_, err := c.client.PolicyApi.PutBucketPolicy(ctx, data.Bucket.ValueString()).BucketPolicy(input).Execute()
	if err != nil {
		return err
	}

	err = backoff.Retry(func() error {
		if _, retryErr := c.GetBucketPolicyCheck(ctx, data.Bucket.ValueString()); retryErr != nil {
			if errors.Is(retryErr, ErrBucketPolicyNotFound) {
				return retryErr
			}

			return backoff.Permanent(fmt.Errorf("failed to check if bucket policy exists: %w", retryErr))
		}
		return nil
	}, backoff.NewExponentialBackOff(backoff.WithMaxElapsedTime(utils.DefaultTimeout)))

	return err
}

// GetBucketPolicy gets a bucket policy.
func (c *Client) GetBucketPolicy(ctx context.Context, bucketName types.String) (*BucketPolicyModel, bool, error) {
	output, apiResponse, err := c.client.PolicyApi.GetBucketPolicy(ctx, bucketName.ValueString()).Execute()
	if apiResponse.HttpNotFound() {
		return nil, false, nil
	}

	if err != nil {
		return nil, false, err
	}

	data := &BucketPolicyModel{
		Bucket: bucketName,
	}
	setBucketPolicyData(output, data)
	return data, true, nil
}

// UpdateBucketPolicy updates a bucket policy.
func (c *Client) UpdateBucketPolicy(ctx context.Context, data *BucketPolicyModel) error {
	if err := c.CreateBucketPolicy(ctx, data); err != nil {
		return err
	}

	model, found, err := c.GetBucketPolicy(ctx, data.Bucket)
	if !found {
		return fmt.Errorf("bucket policy not found for bucket %s", data.Bucket.ValueString())
	}

	if err != nil {
		return err
	}

	*data = *model
	return nil
}

// DeleteBucketPolicy deletes a bucket policy.
func (c *Client) DeleteBucketPolicy(ctx context.Context, bucketName types.String) error {
	_, err := c.client.PolicyApi.DeleteBucketPolicy(ctx, bucketName.ValueString()).Execute()
	if err != nil {
		return err
	}

	err = backoff.Retry(func() error {
		if _, retryErr := c.GetBucketPolicyCheck(ctx, bucketName.ValueString()); retryErr != nil {
			if errors.Is(retryErr, ErrBucketPolicyNotFound) {
				return nil
			}
			return backoff.Permanent(fmt.Errorf("failed to check if bucket policy is deleted: %w", retryErr))
		}
		return nil
	}, backoff.NewExponentialBackOff(backoff.WithMaxElapsedTime(utils.DefaultTimeout)))

	return err
}

// GetBucketPolicyCheck gets a bucket policy.
func (c *Client) GetBucketPolicyCheck(ctx context.Context, bucketName string) (*objstorage.BucketPolicy, error) {
	policy, apiResponse, err := c.client.PolicyApi.GetBucketPolicy(ctx, bucketName).Execute()
	if err != nil {
		if apiResponse.HttpNotFound() {
			return nil, ErrBucketPolicyNotFound
		}
		return nil, err
	}
	return policy, nil
}

func setBucketPolicyData(policyResponse *objstorage.BucketPolicy, data *BucketPolicyModel) diag.Diagnostics {
	var diags diag.Diagnostics

	policyData := bucketPolicy{
		ID:      policyResponse.Id,
		Version: policyResponse.Version,
	}

	if policyResponse.Statement != nil {
		policyData.Statement = make([]bucketPolicyStatement, 0, len(*policyResponse.Statement))
		for _, statementResponse := range *policyResponse.Statement {
			statementData := bucketPolicyStatement{
				SID:       statementResponse.Sid,
				Effect:    *statementResponse.Effect,
				Action:    *statementResponse.Action,
				Resources: *statementResponse.Resource,
			}
			if statementResponse.Principal != nil && statementResponse.Principal.AWS != nil {
				statementData.Principal = *statementResponse.Principal.AWS
			}
			if statementResponse.Condition != nil {
				conditionData := bucketPolicyStatementCondition{}
				if statementResponse.Condition.IpAddress != nil && statementResponse.Condition.IpAddress.AwsSourceIp != nil {
					conditionData.IPs = *statementResponse.Condition.IpAddress.AwsSourceIp
				}
				if statementResponse.Condition.NotIpAddress != nil && statementResponse.Condition.NotIpAddress.AwsSourceIp != nil {
					conditionData.ExcludedIPs = *statementResponse.Condition.NotIpAddress.AwsSourceIp
				}
				if statementResponse.Condition.DateGreaterThan != nil {
					dateString := statementResponse.Condition.DateGreaterThan.AwsCurrentTime.Format(constant.DatetimeZLayout)
					conditionData.DateGreaterThan = &dateString
				}
				if statementResponse.Condition.DateLessThan != nil {
					dateString := statementResponse.Condition.DateLessThan.AwsCurrentTime.Format(constant.DatetimeZLayout)
					conditionData.DateLessThan = &dateString
				}
				statementData.Condition = &conditionData
			}
			policyData.Statement = append(policyData.Statement, statementData)
		}
	}
	policyString, err := json.Marshal(policyData)
	if err != nil {
		diags.AddError("Error serializing policy data", err.Error())
		return diags
	}
	data.Policy = jsontypes.NewNormalizedValue(string(policyString))
	return diags
}

func buildBucketPolicyFromModel(policyModel *BucketPolicyModel) (objstorage.BucketPolicy, diag.Diagnostics) {
	var diags diag.Diagnostics
	policyInput := objstorage.BucketPolicy{}
	policyData := bucketPolicy{}

	// Can't unmarshal directly in the API object, need to use an intermediary
	if diags = policyModel.Policy.Unmarshal(&policyData); diags.HasError() {
		return objstorage.BucketPolicy{}, diags
	}

	policyInput.Id = policyData.ID
	policyInput.Version = policyData.Version

	statement := make([]objstorage.BucketPolicyStatement, 0, len(policyData.Statement))
	for _, statementData := range policyData.Statement {
		statementInput := objstorage.NewBucketPolicyStatement(statementData.Action, statementData.Effect, statementData.Resources)
		statementInput.Sid = statementData.SID
		statementInput.Principal = objstorage.NewPrincipal(statementData.Principal)

		if statementData.Condition != nil {
			statementInput.Condition = objstorage.NewBucketPolicyCondition()
			if statementData.Condition.IPs != nil {
				statementInput.Condition.IpAddress = objstorage.NewBucketPolicyConditionIpAddress()
				ips := statementData.Condition.IPs
				statementInput.Condition.IpAddress.AwsSourceIp = &ips
			}
			if statementData.Condition.ExcludedIPs != nil {
				statementInput.Condition.NotIpAddress = objstorage.NewBucketPolicyConditionIpAddress()
				excludedIPs := statementData.Condition.ExcludedIPs
				statementInput.Condition.NotIpAddress.AwsSourceIp = &excludedIPs
			}
			if statementData.Condition.DateGreaterThan != nil {
				var t *objstorage.IonosTime
				var err error
				if t, err = convertToIonosTime(*statementData.Condition.DateGreaterThan); err != nil {
					diags.AddError("Error converting policy condition 'greater than' date", err.Error())
					return objstorage.BucketPolicy{}, diags
				}
				dateGreater := objstorage.BucketPolicyConditionDate{AwsCurrentTime: t}
				statementInput.Condition.DateGreaterThan = &dateGreater
			}
			if statementData.Condition.DateLessThan != nil {
				var t *objstorage.IonosTime
				var err error
				if t, err = convertToIonosTime(*statementData.Condition.DateLessThan); err != nil {
					diags.AddError("Error converting policy condition 'less than' date", err.Error())
					return objstorage.BucketPolicy{}, diags
				}
				dateLess := objstorage.BucketPolicyConditionDate{AwsCurrentTime: t}
				statementInput.Condition.DateLessThan = &dateLess
			}
		}
		statement = append(statement, *statementInput)
	}
	policyInput.Statement = &statement

	return policyInput, diags
}

func convertToIonosTime(targetTime string) (*objstorage.IonosTime, error) {
	var ionosTime objstorage.IonosTime
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
