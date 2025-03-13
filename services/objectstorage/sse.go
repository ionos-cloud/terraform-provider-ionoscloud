package objectstorage

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/types"
	objstorage "github.com/ionos-cloud/sdk-go-bundle/products/objectstorage/v2"
)

// ServerSideEncryptionConfigurationModel defines the expected inputs for creating a new ServerSideEncryptionConfiguration.
type ServerSideEncryptionConfigurationModel struct {
	Bucket types.String `tfsdk:"bucket"`
	Rules  []sseRule    `tfsdk:"rule"`
}

type sseRule struct {
	ApplyServerSideEncryptionByDefault applyServerSideEncryptionByDefault `tfsdk:"apply_server_side_encryption_by_default"`
}

type applyServerSideEncryptionByDefault struct {
	SSEAlgorithm types.String `tfsdk:"sse_algorithm"`
}

// CreateBucketSSE creates a new ServerSideEncryptionConfiguration.
func (c *Client) CreateBucketSSE(ctx context.Context, data *ServerSideEncryptionConfigurationModel) error {
	_, err := c.client.EncryptionApi.PutBucketEncryption(ctx, data.Bucket.ValueString()).
		PutBucketEncryptionRequest(buildServerSideEncryptionConfigurationFromModel(data)).
		Execute()
	return err
}

// GetBucketSSE gets a ServerSideEncryptionConfiguration.
func (c *Client) GetBucketSSE(ctx context.Context, bucketName types.String) (*ServerSideEncryptionConfigurationModel, bool, error) {
	output, apiResponse, err := c.client.EncryptionApi.GetBucketEncryption(ctx, bucketName.ValueString()).Execute()
	if apiResponse.HttpNotFound() {
		return nil, false, nil
	}

	if err != nil {
		return nil, false, err
	}

	return buildServerSideEncryptionConfigurationModelFromAPIResponse(output, &ServerSideEncryptionConfigurationModel{Bucket: bucketName}), true, nil
}

// UpdateBucketSSE updates a ServerSideEncryptionConfiguration.
func (c *Client) UpdateBucketSSE(ctx context.Context, data *ServerSideEncryptionConfigurationModel) error {
	if err := c.CreateBucketSSE(ctx, data); err != nil {
		return err
	}

	model, found, err := c.GetBucketSSE(ctx, data.Bucket)
	if !found {
		return err
	}

	if err != nil {
		return err
	}

	*data = *model
	return nil
}

// DeleteBucketSSE deletes a ServerSideEncryptionConfiguration.
func (c *Client) DeleteBucketSSE(ctx context.Context, bucketName types.String) error {
	_, err := c.client.EncryptionApi.DeleteBucketEncryption(ctx, bucketName.ValueString()).Execute()
	return err
}

func buildServerSideEncryptionConfigurationModelFromAPIResponse(output *objstorage.ServerSideEncryptionConfiguration, data *ServerSideEncryptionConfigurationModel) *ServerSideEncryptionConfigurationModel {
	return &ServerSideEncryptionConfigurationModel{
		Bucket: data.Bucket,
		Rules:  buildServerSideEncryptionRulesFromAPIResponse(output.Rules),
	}
}

func buildServerSideEncryptionRulesFromAPIResponse(data *[]objstorage.ServerSideEncryptionRule) []sseRule {
	if data == nil {
		return nil
	}

	rules := make([]sseRule, 0, len(*data))
	for _, r := range *data {
		if r.ApplyServerSideEncryptionByDefault == nil {
			continue
		}

		if r.ApplyServerSideEncryptionByDefault.SSEAlgorithm == nil {
			continue
		}

		rules = append(rules, sseRule{
			ApplyServerSideEncryptionByDefault: applyServerSideEncryptionByDefault{
				SSEAlgorithm: types.StringValue(string(*r.ApplyServerSideEncryptionByDefault.SSEAlgorithm)),
			},
		})
	}

	return rules
}

func buildServerSideEncryptionConfigurationFromModel(data *ServerSideEncryptionConfigurationModel) objstorage.PutBucketEncryptionRequest {
	return objstorage.PutBucketEncryptionRequest{
		Rules: buildServerSideEncryptionRulesFromModel(data.Rules),
	}
}

func buildServerSideEncryptionRulesFromModel(data []sseRule) *[]objstorage.ServerSideEncryptionRule {
	rules := make([]objstorage.ServerSideEncryptionRule, 0, len(data))
	for _, r := range data {
		rules = append(rules, objstorage.ServerSideEncryptionRule{
			ApplyServerSideEncryptionByDefault: &objstorage.ServerSideEncryptionByDefault{
				SSEAlgorithm: objstorage.ServerSideEncryption(r.ApplyServerSideEncryptionByDefault.SSEAlgorithm.ValueString()).Ptr(),
			},
		})
	}

	return &rules
}
