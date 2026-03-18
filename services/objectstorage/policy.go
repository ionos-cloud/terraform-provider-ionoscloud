package objectstorage

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"reflect"
	"time"

	"github.com/cenkalti/backoff/v4"
	"github.com/hashicorp/terraform-plugin-framework-jsontypes/jsontypes"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	objstorage "github.com/ionos-cloud/sdk-go-bundle/products/objectstorage/v2"

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
	Principal bucketPolicyPrincipal           `json:"Principal"`
	Condition *bucketPolicyStatementCondition `json:"Condition,omitempty"`
}

// bucketPolicyPrincipal is the canonical representation of the S3 Principal field.
// It always marshals to {"AWS": [...]} to match what the API stores and returns,
// and accepts all common input forms on unmarshal:
//   - {"AWS": ["arn:..."]} — API array form
//   - {"AWS": "arn:..."}   — API single-string form
//   - ["arn:..."]          — legacy flat array form
//   - "arn:..." / "*"      — bare string form
type bucketPolicyPrincipal struct {
	AWS []string
}

func (p *bucketPolicyPrincipal) MarshalJSON() ([]byte, error) {
	// Match the API's canonical form: single principal → string, multiple → array
	if len(p.AWS) == 1 {
		return json.Marshal(struct {
			AWS string `json:"AWS"`
		}{AWS: p.AWS[0]})
	}
	return json.Marshal(struct {
		AWS []string `json:"AWS"`
	}{AWS: p.AWS})
}

func (p *bucketPolicyPrincipal) UnmarshalJSON(data []byte) error {
	// object form: {"AWS": "..." | [...]}
	var obj struct {
		AWS json.RawMessage `json:"AWS"`
	}
	if err := json.Unmarshal(data, &obj); err == nil && obj.AWS != nil {
		var arr []string
		if err := json.Unmarshal(obj.AWS, &arr); err == nil {
			p.AWS = arr
			return nil
		}
		var str string
		if err := json.Unmarshal(obj.AWS, &str); err == nil {
			p.AWS = []string{str}
			return nil
		}
	}
	// flat array: ["arn:..."]
	var arr []string
	if err := json.Unmarshal(data, &arr); err == nil {
		p.AWS = arr
		return nil
	}
	// bare string: "*" or "arn:..."
	var s string
	if err := json.Unmarshal(data, &s); err == nil {
		p.AWS = []string{s}
		return nil
	}
	return fmt.Errorf("cannot unmarshal Principal: %s", string(data))
}

type bucketPolicyStatementCondition struct {
	IPs             []string `json:"IpAddress,omitempty"`
	ExcludedIPs     []string `json:"NotIpAddress,omitempty"`
	DateGreaterThan *string  `json:"DateGreaterThan,omitempty"`
	DateLessThan    *string  `json:"DateLessThan,omitempty"`
}

// PoliciesSemanticEqual returns true if two policy JSON strings represent the
// same policy, regardless of JSON key ordering or Principal format differences
// (e.g. ["arn:..."] vs {"AWS": "arn:..."}).
func PoliciesSemanticEqual(statePolicy, apiPolicy string) bool {
	normalize := func(s string) any {
		var data bucketPolicy
		if err := json.Unmarshal([]byte(s), &data); err != nil {
			log.Printf("[WARN] PoliciesSemanticEqual: failed to unmarshal policy: %v", err)
			var raw any
			json.Unmarshal([]byte(s), &raw)
			return raw
		}
		normalized, err := json.Marshal(data)
		if err != nil {
			log.Printf("[WARN] PoliciesSemanticEqual: failed to marshal normalized policy: %v", err)
			var raw any
			json.Unmarshal([]byte(s), &raw)
			return raw
		}
		var result any
		json.Unmarshal(normalized, &result)
		return result
	}
	return reflect.DeepEqual(normalize(statePolicy), normalize(apiPolicy))
}

// CreateBucketPolicy creates a new bucket policy.
func (c *Client) CreateBucketPolicy(ctx context.Context, data *BucketPolicyModel) error {
	region, err := c.GetBucketLocation(ctx, data.Bucket)
	if err != nil {
		return err
	}
	err = c.ChangeConfigURL(region.ValueString())
	if err != nil {
		return err
	}

	input, diags := buildBucketPolicyFromModel(data)
	if diags.HasError() {
		return fmt.Errorf("error building bucket policy: %v", diags)
	}

	_, err = c.client.PolicyApi.PutBucketPolicy(ctx, data.Bucket.ValueString()).BucketPolicy(input).Execute()
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
	region, err := c.GetBucketLocation(ctx, bucketName)
	if err != nil {
		return nil, false, err
	}
	err = c.ChangeConfigURL(region.ValueString())
	if err != nil {
		return nil, false, err
	}

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
	region, err := c.GetBucketLocation(ctx, data.Bucket)
	if err != nil {
		return err
	}
	err = c.ChangeConfigURL(region.ValueString())
	if err != nil {
		return err
	}

	if err := c.CreateBucketPolicy(ctx, data); err != nil {
		return err
	}

	// Verify the policy was persisted, but don't overwrite data with the
	// re-serialized API response — that would change the JSON structure
	// (e.g. key ordering, Principal format) and cause a plan/state mismatch.
	_, found, err := c.GetBucketPolicy(ctx, data.Bucket)
	if !found {
		return fmt.Errorf("bucket policy not found for bucket %s", data.Bucket.ValueString())
	}

	return err
}

// DeleteBucketPolicy deletes a bucket policy.
func (c *Client) DeleteBucketPolicy(ctx context.Context, bucketName types.String) error {
	region, err := c.GetBucketLocation(ctx, bucketName)
	if err != nil {
		return err
	}
	err = c.ChangeConfigURL(region.ValueString())
	if err != nil {
		return err
	}

	_, err = c.client.PolicyApi.DeleteBucketPolicy(ctx, bucketName.ValueString()).Execute()
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
		policyData.Statement = make([]bucketPolicyStatement, 0, len(policyResponse.Statement))
		for _, statementResponse := range policyResponse.Statement {
			statementData := bucketPolicyStatement{
				SID:       statementResponse.Sid,
				Effect:    statementResponse.Effect,
				Action:    statementResponse.Action,
				Resources: statementResponse.Resource,
			}
			if statementResponse.Principal != nil && statementResponse.Principal.AWS != nil {
				statementData.Principal = bucketPolicyPrincipal{AWS: statementResponse.Principal.AWS}
			}
			if statementResponse.Condition != nil {
				conditionData := bucketPolicyStatementCondition{}
				if statementResponse.Condition.IpAddress != nil && statementResponse.Condition.IpAddress.AwsSourceIp != nil {
					conditionData.IPs = statementResponse.Condition.IpAddress.AwsSourceIp
				}
				if statementResponse.Condition.NotIpAddress != nil && statementResponse.Condition.NotIpAddress.AwsSourceIp != nil {
					conditionData.ExcludedIPs = statementResponse.Condition.NotIpAddress.AwsSourceIp
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
		statementInput.Principal = objstorage.NewPrincipal(statementData.Principal.AWS)

		if statementData.Condition != nil {
			statementInput.Condition = objstorage.NewBucketPolicyCondition()
			if statementData.Condition.IPs != nil {
				statementInput.Condition.IpAddress = objstorage.NewBucketPolicyConditionIpAddress()
				ips := statementData.Condition.IPs
				statementInput.Condition.IpAddress.AwsSourceIp = ips
			}
			if statementData.Condition.ExcludedIPs != nil {
				statementInput.Condition.NotIpAddress = objstorage.NewBucketPolicyConditionIpAddress()
				excludedIPs := statementData.Condition.ExcludedIPs
				statementInput.Condition.NotIpAddress.AwsSourceIp = excludedIPs
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
	policyInput.Statement = statement

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
