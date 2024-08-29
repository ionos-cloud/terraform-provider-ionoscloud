package s3

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	s3 "github.com/ionos-cloud/sdk-go-s3"
)

// BucketVersioningResourceModel defines the expected inputs for creating a new BucketVersioning.
type BucketVersioningResourceModel struct {
	Bucket                  types.String             `tfsdk:"bucket"`
	VersioningConfiguration *versioningConfiguration `tfsdk:"versioning_configuration"`
}

type versioningConfiguration struct {
	Status    types.String `tfsdk:"status"`
	MfaDelete types.String `tfsdk:"mfa_delete"`
}

func (v versioningConfiguration) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"status":     types.StringType,
		"mfa_delete": types.StringType,
	}
}

// CreateBucketVersioning creates a new BucketVersioning.
func (c *Client) CreateBucketVersioning(ctx context.Context, data *BucketVersioningResourceModel) error {
	_, err := c.client.VersioningApi.PutBucketVersioning(ctx, data.Bucket.ValueString()).PutBucketVersioningRequest(buildPutVersioningRequestFromModel(data)).Execute()
	return err
}

// GetBucketVersioning gets a BucketVersioning.
func (c *Client) GetBucketVersioning(ctx context.Context, bucketName types.String) (*BucketVersioningResourceModel, bool, error) {
	output, resp, err := c.client.VersioningApi.GetBucketVersioning(ctx, bucketName.ValueString()).Execute()
	if resp.HttpNotFound() {
		return nil, false, nil
	}

	if err != nil {
		return nil, false, err
	}

	builtModel := buildBucketVersioningModelFromAPIResponse(output, bucketName)
	return builtModel, true, nil
}

// UpdateBucketVersioning updates a BucketVersioning.
func (c *Client) UpdateBucketVersioning(ctx context.Context, data *BucketVersioningResourceModel) error {
	if err := c.CreateBucketVersioning(ctx, data); err != nil {
		return err
	}

	model, found, err := c.GetBucketVersioning(ctx, data.Bucket)
	if !found {
		return fmt.Errorf("bucket versioning not found")
	}

	if err != nil {
		return err
	}

	*data = *model
	return nil
}

// DeleteBucketVersioning deletes a BucketVersioning.
func (c *Client) DeleteBucketVersioning(ctx context.Context, data *BucketVersioningResourceModel) error {
	// Removing S3 bucket versioning for un-versioned bucket from state
	if data.VersioningConfiguration.Status.ValueString() == string(s3.BUCKETVERSIONINGSTATUS_SUSPENDED) {
		return nil
	}

	_, err := c.client.VersioningApi.PutBucketVersioning(ctx, data.Bucket.ValueString()).
		PutBucketVersioningRequest(s3.PutBucketVersioningRequest{
			Status: s3.BUCKETVERSIONINGSTATUS_SUSPENDED.Ptr(),
		}).Execute()
	if isInvalidStateBucketWithObjectLock(err) {
		return nil
	}

	return err
}

func buildPutVersioningRequestFromModel(data *BucketVersioningResourceModel) s3.PutBucketVersioningRequest {
	var request s3.PutBucketVersioningRequest
	if !data.VersioningConfiguration.Status.IsNull() {
		request.Status = s3.BucketVersioningStatus(data.VersioningConfiguration.Status.ValueString()).Ptr()
	}

	if !data.VersioningConfiguration.MfaDelete.IsNull() {
		request.MfaDelete = s3.MfaDeleteStatus(data.VersioningConfiguration.MfaDelete.ValueString()).Ptr()
	}
	return request
}

func buildBucketVersioningModelFromAPIResponse(output *s3.GetBucketVersioningOutput, bucket types.String) *BucketVersioningResourceModel {
	var vc versioningConfiguration
	if output.Status != nil {
		vc.Status = types.StringValue(string(*output.Status))
	}

	if output.MfaDelete != nil {
		vc.MfaDelete = types.StringValue(string(*output.MfaDelete))
	}

	built := BucketVersioningResourceModel{
		Bucket:                  bucket,
		VersioningConfiguration: &vc,
	}

	return &built
}
