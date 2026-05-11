package userobjectstorage

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/cenkalti/backoff/v4"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/ionos-cloud/sdk-go-bundle/products/userobjectstorage/v2"
	"github.com/ionos-cloud/sdk-go-bundle/shared"
)

// CreateBucket creates a new user-owned bucket in the given region.
func (c *Client) CreateBucket(ctx context.Context, name, region types.String, objectLock types.Bool, timeout time.Duration) error {
	if err := c.ChangeRegion(region.ValueString()); err != nil {
		return err
	}

	createBucketConfig := userobjectstorage.CreateBucketRequestCreateBucketConfiguration{
		LocationConstraint: region.ValueStringPointer(),
	}
	createReq := userobjectstorage.NewCreateBucketRequest()
	createReq.CreateBucketConfiguration = &createBucketConfig

	if _, err := c.client.BucketsApi.CreateBucket(ctx, name.ValueString()).
		CreateBucketRequest(*createReq).
		XAmzBucketObjectLockEnabled(objectLock.ValueBool()).
		Execute(); err != nil {
		return fmt.Errorf("failed to create bucket %q: %w", name.ValueString(), err)
	}

	err := backoff.Retry(func() error {
		return c.bucketExistsCheck(ctx, name.ValueString())
	}, backoff.NewExponentialBackOff(backoff.WithMaxElapsedTime(timeout)))
	if err != nil {
		return fmt.Errorf("failed to wait for bucket creation: %w", err)
	}

	return nil
}

// GetBucket checks whether a bucket exists. Returns (false, nil) if the bucket is not found.
// Region is not refreshed from the API — it is sourced from Terraform state since region is
// RequiresReplace and the SDK's GetBucketLocation response model has a generator-level XML tag
// issue that cannot be fixed without changing the generator templates.
func (c *Client) GetBucket(ctx context.Context, name, region types.String) (bool, error) {
	if err := c.ChangeRegion(region.ValueString()); err != nil {
		return false, err
	}

	apiResp, err := c.client.BucketsApi.HeadBucket(ctx, name.ValueString()).Execute()
	if apiResp.HttpNotFound() {
		return false, nil
	}
	if err != nil {
		return false, fmt.Errorf("failed to check bucket %q: %w", name.ValueString(), err)
	}

	return true, nil
}

// DeleteBucket deletes a bucket. If forceDestroy is true and the bucket is not empty, all objects
// are deleted first.
func (c *Client) DeleteBucket(ctx context.Context, name types.String, forceDestroy types.Bool, region types.String, timeout time.Duration) error {
	if err := c.ChangeRegion(region.ValueString()); err != nil {
		return err
	}

	apiResp, err := c.client.BucketsApi.DeleteBucket(ctx, name.ValueString()).Execute()
	if apiResp.HttpNotFound() {
		return nil
	}

	if isBucketNotEmptyError(err) && forceDestroy.ValueBool() {
		if err = c.deleteAllObjects(ctx, name.ValueString()); err != nil {
			return fmt.Errorf("failed to empty bucket %q: %w", name.ValueString(), err)
		}
		return c.DeleteBucket(ctx, name, forceDestroy, region, timeout)
	}

	if err != nil {
		return fmt.Errorf("failed to delete bucket %q: %w", name.ValueString(), err)
	}

	backOff := backoff.NewExponentialBackOff(backoff.WithMaxElapsedTime(timeout))
	if err = backoff.Retry(func() error {
		return c.bucketDeletedCheck(ctx, name.ValueString())
	}, backOff); err != nil {
		return fmt.Errorf("failed to wait for bucket deletion: %w", err)
	}

	return nil
}

// deleteAllObjects lists and deletes every object in the bucket (used for force_destroy).
func (c *Client) deleteAllObjects(ctx context.Context, bucketName string) error {
	var continuationToken string
	for {
		req := c.client.ObjectsApi.ListObjectsV2(ctx, bucketName)
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
			if _, _, err := c.client.ObjectsApi.DeleteObject(ctx, bucketName, *obj.Key).Execute(); err != nil {
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

func (c *Client) bucketExistsCheck(ctx context.Context, name string) error {
	apiResp, err := c.client.BucketsApi.HeadBucket(ctx, name).Execute()
	if apiResp.HttpNotFound() {
		return fmt.Errorf("bucket not found")
	}
	if err != nil {
		return backoff.Permanent(fmt.Errorf("failed to check if bucket exists: %w", err))
	}
	return nil
}

func (c *Client) bucketDeletedCheck(ctx context.Context, name string) error {
	apiResp, err := c.client.BucketsApi.HeadBucket(ctx, name).Execute()
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
