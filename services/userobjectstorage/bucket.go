package userobjectstorage

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/cenkalti/backoff/v4"
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/ionos-cloud/sdk-go-bundle/products/userobjectstorage/v2"
	"github.com/ionos-cloud/sdk-go-bundle/shared"

	diagutil "github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils/diags"
)

// BucketModel is the resource model for an ionoscloud_user_object_storage_bucket.
type BucketModel struct {
	ForceDestroy      types.Bool     `tfsdk:"force_destroy"`
	ID                types.String   `tfsdk:"id"`
	Name              types.String   `tfsdk:"name"`
	ObjectLockEnabled types.Bool     `tfsdk:"object_lock_enabled"`
	Region            types.String   `tfsdk:"region"`
	Timeouts          timeouts.Value `tfsdk:"timeouts"`
}

// CreateBucket creates a new user-owned bucket.
func (c *Client) CreateBucket(ctx context.Context, data BucketModel, timeout time.Duration) error {
	client, err := c.apiClientForRegion(data.Region.ValueString())
	if err != nil {
		return err
	}

	createBucketConfig := userobjectstorage.CreateBucketRequestCreateBucketConfiguration{
		LocationConstraint: data.Region.ValueStringPointer(),
	}
	createReq := userobjectstorage.NewCreateBucketRequest()
	createReq.CreateBucketConfiguration = &createBucketConfig
	if _, err := client.BucketsApi.CreateBucket(ctx, data.Name.ValueString()).
		CreateBucketRequest(*createReq).
		XAmzBucketObjectLockEnabled(data.ObjectLockEnabled.ValueBool()).
		Execute(); err != nil {
		return fmt.Errorf("failed to create bucket %q: %w", data.Name.ValueString(), err)
	}

	backoffErr := backoff.Retry(func() error {
		return c.bucketExistsCheck(ctx, client, data.Name.ValueString())
	}, backoff.NewExponentialBackOff(backoff.WithMaxElapsedTime(timeout)))
	if backoffErr != nil {
		return fmt.Errorf("failed to wait for bucket creation: %w", backoffErr)
	}

	return nil
}

// GetBucket checks whether a bucket exists. Returns (false, nil) if the bucket is not found.
// Region is sourced from Terraform state — the SDK's GetBucketLocation response model has a
// generator-level XML tag issue that cannot be fixed without changing the generator templates.
func (c *Client) GetBucket(ctx context.Context, name, region types.String) (bool, error) {
	client, err := c.apiClientForRegion(region.ValueString())
	if err != nil {
		return false, err
	}

	apiResp, err := client.BucketsApi.HeadBucket(ctx, name.ValueString()).Execute()
	if apiResp.HttpNotFound() {
		return false, nil
	}
	if err != nil {
		return false, fmt.Errorf("failed to check bucket %q: %w", name.ValueString(), err)
	}

	return true, nil
}

// GetObjectLockEnabled returns true if Object Lock is enabled for the bucket.
// Returns false when the bucket has no Object Lock configuration (404) or when it is not enabled.
func (c *Client) GetObjectLockEnabled(ctx context.Context, name types.String) (types.Bool, error) {
	// Object lock configuration is not region-specific at the API level; use the default region.
	api, err := c.apiClientForRegion(DefaultRegion)
	if err != nil {
		return types.BoolNull(), err
	}

	output, apiResp, err := api.ObjectLockApi.GetObjectLockConfiguration(ctx, name.ValueString()).Execute()
	if apiResp.HttpNotFound() {
		return types.BoolValue(false), nil
	}
	if err != nil {
		return types.BoolNull(), fmt.Errorf("failed to get object lock configuration for %q: %w", name.ValueString(), err)
	}
	return types.BoolValue(output.GetObjectLockEnabled() == "Enabled"), nil
}

// DeleteBucket deletes a bucket. If forceDestroy is true and the bucket is not empty, all objects
// are deleted first.
func (c *Client) DeleteBucket(ctx context.Context, data BucketModel, timeout time.Duration) error {
	api, err := c.apiClientForRegion(data.Region.ValueString())
	if err != nil {
		return err
	}

	apiResp, err := api.BucketsApi.DeleteBucket(ctx, data.Name.ValueString()).Execute()
	if apiResp.HttpNotFound() {
		return nil
	}

	if isBucketNotEmptyError(err) && data.ForceDestroy.ValueBool() {
		if err = c.deleteAllObjects(ctx, api, data.Name.ValueString()); err != nil {
			return fmt.Errorf("failed to empty bucket %q: %w", data.Name.ValueString(), err)
		}
		return c.DeleteBucket(ctx, data, timeout)
	}

	if err != nil {
		return fmt.Errorf("failed to delete bucket %q: %w", data.Name.ValueString(), err)
	}

	backOff := backoff.NewExponentialBackOff(backoff.WithMaxElapsedTime(timeout))
	if err = backoff.Retry(func() error {
		return c.bucketDeletedCheck(ctx, api, data.Name.ValueString())
	}, backOff); err != nil {
		return fmt.Errorf("failed to wait for bucket deletion: %w", err)
	}

	return nil
}

// deleteAllObjects lists and deletes every object in the bucket (used for force_destroy).
func (c *Client) deleteAllObjects(ctx context.Context, api *userobjectstorage.APIClient, bucketName string) error {
	var continuationToken string
	for {
		req := api.ObjectsApi.ListObjectsV2(ctx, bucketName)
		if continuationToken != "" {
			req = req.ContinuationToken(continuationToken)
		}

		output, _, err := req.Execute()
		if err != nil {
			return fmt.Errorf("failed to list objects: %w", err)
		}

		for _, obj := range output.Contents {
			if obj.Key == nil {
				continue
			}
			if _, _, err := api.ObjectsApi.DeleteObject(ctx, bucketName, *obj.Key).Execute(); err != nil {
				return fmt.Errorf("failed to delete object %q: %w", *obj.Key, err)
			}
		}

		if !output.IsTruncated {
			break
		}
		if output.NextContinuationToken == nil {
			break
		}
		continuationToken = *output.NextContinuationToken
	}
	return nil
}

func (c *Client) bucketExistsCheck(ctx context.Context, api *userobjectstorage.APIClient, name string) error {
	apiResp, err := api.BucketsApi.HeadBucket(ctx, name).Execute()
	if apiResp.HttpNotFound() {
		return fmt.Errorf("bucket not found")
	}
	if err != nil {
		errCtx := &diagutil.ErrorContext{
			ResourceID: name,
			StatusCode: apiResp.SafeStatusCode(),
		}
		if apiResp != nil {
			loc, _ := apiResp.Location()
			if loc != nil {
				errCtx.RequestID = loc.String()
			}
		}
		return backoff.Permanent(diagutil.WrapError(fmt.Errorf("failed to check if bucket exists: %w", err),
			errCtx))
	}
	return nil
}

func (c *Client) bucketDeletedCheck(ctx context.Context, api *userobjectstorage.APIClient, name string) error {
	apiResp, err := api.BucketsApi.HeadBucket(ctx, name).Execute()
	if apiResp.HttpNotFound() {
		return nil
	}
	if err != nil {
		return backoff.Permanent(fmt.Errorf("failed to check if bucket was deleted: %w", err))
	}
	return fmt.Errorf("bucket still exists")
}

func isBucketNotEmptyError(err error) bool {
	var apiErr shared.GenericOpenAPIError
	return errors.As(err, &apiErr) && apiErr.StatusCode() == 409
}
