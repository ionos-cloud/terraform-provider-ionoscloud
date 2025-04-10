package objectstorage

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/types"
	objstorage "github.com/ionos-cloud/sdk-go-bundle/products/objectstorage/v2"
)

// BucketPublicAccessBlockResourceModel defines the expected inputs for creating a new BucketPublicAccessBlock.
type BucketPublicAccessBlockResourceModel struct {
	Bucket                types.String `tfsdk:"bucket"`
	BlockPublicACLS       types.Bool   `tfsdk:"block_public_acls"`
	BlockPublicPolicy     types.Bool   `tfsdk:"block_public_policy"`
	IgnorePublicACLS      types.Bool   `tfsdk:"ignore_public_acls"`
	RestrictPublicBuckets types.Bool   `tfsdk:"restrict_public_buckets"`
}

// CreateBucketPublicAccessBlock creates a new BucketPublicAccessBlock.
func (c *Client) CreateBucketPublicAccessBlock(ctx context.Context, data *BucketPublicAccessBlockResourceModel) error {
	_, err := c.client.PublicAccessBlockApi.PutPublicAccessBlock(ctx, data.Bucket.ValueString()).BlockPublicAccessPayload(buildPublicAccessBlockFromModel(data)).Execute()
	return err
}

// GetBucketPublicAccessBlock gets a BucketPublicAccessBlock.
func (c *Client) GetBucketPublicAccessBlock(ctx context.Context, bucketName types.String) (*BucketPublicAccessBlockResourceModel, bool, error) {
	output, apiResponse, err := c.client.PublicAccessBlockApi.GetPublicAccessBlock(ctx, bucketName.ValueString()).Execute()
	if apiResponse.HttpNotFound() {
		return nil, false, nil
	}

	if err != nil {
		return nil, false, err
	}

	return buildPublicAccessBlockModelFromAPIResponse(output, &BucketPublicAccessBlockResourceModel{Bucket: bucketName}), true, nil
}

// UpdateBucketPublicAccessBlock updates a BucketPublicAccessBlock.
func (c *Client) UpdateBucketPublicAccessBlock(ctx context.Context, data *BucketPublicAccessBlockResourceModel) error {
	if err := c.CreateBucketPublicAccessBlock(ctx, data); err != nil {
		return err
	}

	model, found, err := c.GetBucketPublicAccessBlock(ctx, data.Bucket)
	if !found {
		return err
	}

	if err != nil {
		return err
	}

	*data = *model
	return nil
}

// DeleteBucketPublicAccessBlock deletes a BucketPublicAccessBlock.
func (c *Client) DeleteBucketPublicAccessBlock(ctx context.Context, bucketName types.String) error {
	apiResponse, err := c.client.PublicAccessBlockApi.DeletePublicAccessBlock(ctx, bucketName.ValueString()).Execute()
	if apiResponse.HttpNotFound() {
		return nil
	}

	return err
}

// GetBucketPublicAccessBlockCheck gets a BucketPublicAccessBlock.
func buildPublicAccessBlockFromModel(model *BucketPublicAccessBlockResourceModel) objstorage.BlockPublicAccessPayload {
	input := objstorage.BlockPublicAccessPayload{
		BlockPublicPolicy:     model.BlockPublicPolicy.ValueBoolPointer(),
		IgnorePublicAcls:      model.IgnorePublicACLS.ValueBoolPointer(),
		BlockPublicAcls:       model.BlockPublicACLS.ValueBoolPointer(),
		RestrictPublicBuckets: model.RestrictPublicBuckets.ValueBoolPointer(),
	}
	return input
}

func buildPublicAccessBlockModelFromAPIResponse(output *objstorage.BlockPublicAccessOutput, model *BucketPublicAccessBlockResourceModel) *BucketPublicAccessBlockResourceModel {
	built := &BucketPublicAccessBlockResourceModel{
		Bucket:                model.Bucket,
		BlockPublicACLS:       types.BoolPointerValue(output.BlockPublicAcls),
		BlockPublicPolicy:     types.BoolPointerValue(output.BlockPublicPolicy),
		IgnorePublicACLS:      types.BoolPointerValue(output.IgnorePublicAcls),
		RestrictPublicBuckets: types.BoolPointerValue(output.RestrictPublicBuckets),
	}
	return built
}
