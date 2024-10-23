package objectstorage

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/types"
	objstorage "github.com/ionos-cloud/sdk-go-s3"

	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils/convptr"
)

// BucketCorsConfigurationModel is used to create, update and delete a bucket cors configuration.
type BucketCorsConfigurationModel struct {
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

// CreateBucketCors creates a new bucket cors configuration.
func (c *Client) CreateBucketCors(ctx context.Context, data *BucketCorsConfigurationModel) error {
	_, err := c.client.CORSApi.PutBucketCors(ctx, data.Bucket.ValueString()).
		PutBucketCorsRequest(buildBucketCorsConfigurationFromModel(data)).Execute()
	return err
}

// GetBucketCors gets a bucket cors configuration.
func (c *Client) GetBucketCors(ctx context.Context, bucketName types.String) (*BucketCorsConfigurationModel, bool, error) {
	output, apiResponse, err := c.client.CORSApi.GetBucketCors(ctx, bucketName.ValueString()).Execute()
	if apiResponse.HttpNotFound() {
		return nil, false, nil
	}

	if err != nil {
		return nil, false, err
	}

	return buildBucketCorsConfigurationModelFromAPIResponse(output, &BucketCorsConfigurationModel{Bucket: bucketName}), true, nil
}

// UpdateBucketCors updates a bucket cors configuration.
func (c *Client) UpdateBucketCors(ctx context.Context, data *BucketCorsConfigurationModel) error {
	if err := c.CreateBucketCors(ctx, data); err != nil {
		return err
	}

	model, found, err := c.GetBucketCors(ctx, data.Bucket)
	if !found {
		return fmt.Errorf("bucket cors configuration not found")
	}

	if err != nil {
		return err
	}

	*data = *model

	return nil
}

// DeleteBucketCors deletes a bucket cors configuration.
func (c *Client) DeleteBucketCors(ctx context.Context, bucketName types.String) error {
	apiResponse, err := c.client.CORSApi.DeleteBucketCors(ctx, bucketName.ValueString()).Execute()
	if apiResponse.HttpNotFound() {
		return nil
	}
	return err
}

func buildBucketCorsConfigurationModelFromAPIResponse(output *objstorage.GetBucketCorsOutput, data *BucketCorsConfigurationModel) *BucketCorsConfigurationModel {
	data.Cors = buildCorsRulesFromAPIResponse(output.CORSRules)
	return data
}

func buildCorsRulesFromAPIResponse(rules *[]objstorage.CORSRule) []corsRule {
	if rules == nil {
		return nil
	}

	result := make([]corsRule, 0, len(*rules))
	for _, r := range *rules {
		result = append(result, corsRule{
			ID:             types.Int64PointerValue(convptr.Int32ToInt64(r.ID)),
			AllowedHeaders: toTFStrings(r.AllowedHeaders),
			AllowedMethods: toTFStrings(r.AllowedMethods),
			AllowedOrigins: toTFStrings(r.AllowedOrigins),
			ExposeHeaders:  toTFStrings(r.ExposeHeaders),
			MaxAgeSeconds:  types.Int64PointerValue(convptr.Int32ToInt64(r.MaxAgeSeconds)),
		})
	}

	return result
}
func buildBucketCorsConfigurationFromModel(data *BucketCorsConfigurationModel) objstorage.PutBucketCorsRequest {
	return objstorage.PutBucketCorsRequest{
		CORSRules: buildCorsRulesFromModel(data.Cors),
	}
}

func buildCorsRulesFromModel(rules []corsRule) *[]objstorage.CORSRule {
	result := make([]objstorage.CORSRule, 0, len(rules))
	for _, r := range rules {
		result = append(result, objstorage.CORSRule{
			ID:             convptr.Int64ToInt32(r.ID.ValueInt64Pointer()),
			AllowedHeaders: toStrings(r.AllowedHeaders),
			AllowedMethods: toStrings(r.AllowedMethods),
			AllowedOrigins: toStrings(r.AllowedOrigins),
			ExposeHeaders:  toStrings(r.ExposeHeaders),
			MaxAgeSeconds:  convptr.Int64ToInt32(r.MaxAgeSeconds.ValueInt64Pointer()),
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
