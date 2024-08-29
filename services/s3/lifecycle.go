package s3

import (
	"context"
	"encoding/base64"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/types"
	s3 "github.com/ionos-cloud/sdk-go-s3"

	convptr "github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils/convptr"
	hash2 "github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils/hash"
)

// BucketLifecycleConfigurationModel is used to create, update and delete a bucket lifecycle configuration.
type BucketLifecycleConfigurationModel struct {
	Bucket types.String    `tfsdk:"bucket"`
	Rule   []lifecycleRule `tfsdk:"rule"`
}

type lifecycleRule struct {
	ID                             types.String                    `tfsdk:"id"`
	Prefix                         types.String                    `tfsdk:"prefix"`
	Status                         types.String                    `tfsdk:"status"`
	Expiration                     *expiration                     `tfsdk:"expiration"`
	NoncurrentVersionExpiration    *noncurrentVersionExpiration    `tfsdk:"noncurrent_version_expiration"`
	AbortIncompleteMultipartUpload *abortIncompleteMultipartUpload `tfsdk:"abort_incomplete_multipart_upload"`
}

type expiration struct {
	Days                      types.Int64  `tfsdk:"days"`
	Date                      types.String `tfsdk:"date"`
	ExpiredObjectDeleteMarker types.Bool   `tfsdk:"expired_object_delete_marker"`
}

type noncurrentVersionExpiration struct {
	NoncurrentDays types.Int64 `tfsdk:"noncurrent_days"`
}

type abortIncompleteMultipartUpload struct {
	DaysAfterInitiation types.Int64 `tfsdk:"days_after_initiation"`
}

// CreateBucketLifecycle creates a new bucket lifecycle configuration.
func (c *Client) CreateBucketLifecycle(ctx context.Context, data *BucketLifecycleConfigurationModel) error {
	body := buildBucketLifecycleConfigurationFromModel(data)
	hash, err := hash2.MD5(body)
	if err != nil {
		return fmt.Errorf("failed to generate MD5 sum: %w", err)
	}

	_, err = c.client.LifecycleApi.PutBucketLifecycle(ctx, data.Bucket.ValueString()).PutBucketLifecycleRequest(body).ContentMD5(base64.StdEncoding.EncodeToString([]byte(hash))).Execute()
	return err
}

// GetBucketLifecycle gets a bucket lifecycle configuration.
func (c *Client) GetBucketLifecycle(ctx context.Context, bucketName types.String) (*BucketLifecycleConfigurationModel, bool, error) {
	output, apiResponse, err := c.client.LifecycleApi.GetBucketLifecycle(ctx, bucketName.ValueString()).Execute()
	if apiResponse.HttpNotFound() {
		return nil, false, nil
	}

	if err != nil {
		return nil, false, err
	}

	return buildBucketLifecycleConfigurationModelFromAPIResponse(output, &BucketLifecycleConfigurationModel{Bucket: bucketName}), true, nil
}

// UpdateBucketLifecycle updates a bucket lifecycle configuration.
func (c *Client) UpdateBucketLifecycle(ctx context.Context, data *BucketLifecycleConfigurationModel) error {
	if err := c.CreateBucketLifecycle(ctx, data); err != nil {
		return err
	}

	model, found, err := c.GetBucketLifecycle(ctx, data.Bucket)
	if !found {
		return fmt.Errorf("bucket lifecycle configuration not found")
	}

	if err != nil {
		return err
	}

	*data = *model
	return nil
}

// DeleteBucketLifecycle deletes a bucket lifecycle configuration.
func (c *Client) DeleteBucketLifecycle(ctx context.Context, bucketName types.String) error {
	apiResponse, err := c.client.LifecycleApi.DeleteBucketLifecycle(ctx, bucketName.ValueString()).Execute()
	if apiResponse.HttpNotFound() {
		return nil
	}
	return err
}

func buildBucketLifecycleConfigurationModelFromAPIResponse(output *s3.GetBucketLifecycleOutput, data *BucketLifecycleConfigurationModel) *BucketLifecycleConfigurationModel {
	data.Rule = buildRulesFromAPIResponse(output.Rules)
	return data
}

func buildRulesFromAPIResponse(rules *[]s3.Rule) []lifecycleRule {
	if rules == nil {
		return nil
	}

	result := make([]lifecycleRule, 0, len(*rules))
	for _, r := range *rules {
		result = append(result, lifecycleRule{
			ID:                             types.StringPointerValue(r.ID),
			Prefix:                         types.StringPointerValue(r.Prefix),
			Status:                         types.StringValue(string(*r.Status)),
			Expiration:                     buildExpirationFromAPIResponse(r.Expiration),
			NoncurrentVersionExpiration:    buildNoncurrentVersionExpirationFromAPIResponse(r.NoncurrentVersionExpiration),
			AbortIncompleteMultipartUpload: buildAbortIncompleteMultipartUploadFromAPIResponse(r.AbortIncompleteMultipartUpload),
		})
	}

	return result
}

func buildExpirationFromAPIResponse(exp *s3.LifecycleExpiration) *expiration {
	if exp == nil {
		return nil
	}

	return &expiration{
		Days:                      types.Int64PointerValue(convptr.Int32ToInt64(exp.Days)),
		Date:                      types.StringPointerValue(exp.Date),
		ExpiredObjectDeleteMarker: types.BoolPointerValue(exp.ExpiredObjectDeleteMarker),
	}
}

func buildNoncurrentVersionExpirationFromAPIResponse(exp *s3.NoncurrentVersionExpiration) *noncurrentVersionExpiration {
	if exp == nil {
		return nil
	}

	return &noncurrentVersionExpiration{
		NoncurrentDays: types.Int64PointerValue(convptr.Int32ToInt64(exp.NoncurrentDays)),
	}
}

func buildAbortIncompleteMultipartUploadFromAPIResponse(abort *s3.AbortIncompleteMultipartUpload) *abortIncompleteMultipartUpload {
	if abort == nil {
		return nil
	}

	return &abortIncompleteMultipartUpload{
		DaysAfterInitiation: types.Int64PointerValue(convptr.Int32ToInt64(abort.DaysAfterInitiation)),
	}
}

func buildBucketLifecycleConfigurationFromModel(data *BucketLifecycleConfigurationModel) s3.PutBucketLifecycleRequest {
	return s3.PutBucketLifecycleRequest{
		Rules: buildRulesFromModel(data.Rule),
	}
}

func buildRulesFromModel(rules []lifecycleRule) *[]s3.Rule {
	if rules == nil {
		return nil
	}

	result := make([]s3.Rule, 0, len(rules))
	for _, r := range rules {
		result = append(result, s3.Rule{
			ID:                             r.ID.ValueStringPointer(),
			Prefix:                         r.Prefix.ValueStringPointer(),
			Status:                         s3.ExpirationStatus(r.Status.ValueString()).Ptr(),
			Expiration:                     buildExpirationFromModel(r.Expiration),
			NoncurrentVersionExpiration:    buildNoncurrentVersionExpirationFromModel(r.NoncurrentVersionExpiration),
			AbortIncompleteMultipartUpload: buildAbortIncompleteMultipartUploadFromModel(r.AbortIncompleteMultipartUpload),
		})
	}

	return &result
}

func buildExpirationFromModel(expiration *expiration) *s3.LifecycleExpiration {
	if expiration == nil {
		return nil
	}

	return &s3.LifecycleExpiration{
		Days:                      convptr.Int64ToInt32(expiration.Days.ValueInt64Pointer()),
		Date:                      expiration.Date.ValueStringPointer(),
		ExpiredObjectDeleteMarker: expiration.ExpiredObjectDeleteMarker.ValueBoolPointer(),
	}
}

func buildNoncurrentVersionExpirationFromModel(expiration *noncurrentVersionExpiration) *s3.NoncurrentVersionExpiration {
	if expiration == nil {
		return nil
	}

	return &s3.NoncurrentVersionExpiration{
		NoncurrentDays: convptr.Int64ToInt32(expiration.NoncurrentDays.ValueInt64Pointer()),
	}
}

func buildAbortIncompleteMultipartUploadFromModel(abort *abortIncompleteMultipartUpload) *s3.AbortIncompleteMultipartUpload {
	if abort == nil {
		return nil
	}

	return &s3.AbortIncompleteMultipartUpload{
		DaysAfterInitiation: convptr.Int64ToInt32(abort.DaysAfterInitiation.ValueInt64Pointer()),
	}
}
