package userobjectstorage

import (
	"context"
	"fmt"
	"time"

	"github.com/cenkalti/backoff/v4"
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/types"
	objstorage "github.com/ionos-cloud/sdk-go-bundle/products/userobjectstorage/v2"

	tftags "github.com/ionos-cloud/terraform-provider-ionoscloud/v6/internal/tags"
)

// BucketModel is used to represent a bucket.
type BucketModel struct {
	Name              types.String   `tfsdk:"name"`
	Region            types.String   `tfsdk:"region"`
	ObjectLockEnabled types.Bool     `tfsdk:"object_lock_enabled"`
	ForceDestroy      types.Bool     `tfsdk:"force_destroy"`
	Timeouts          timeouts.Value `tfsdk:"timeouts"`
	Tags              types.Map      `tfsdk:"tags"`
	ID                types.String   `tfsdk:"id"`
}

// BucketDataSourceModel is used to represent a bucket for a data source.
type BucketDataSourceModel struct {
	Name   types.String `tfsdk:"name"`
	Region types.String `tfsdk:"region"`
}

// CreateBucket creates a new bucket.
func (c *Client) CreateBucket(ctx context.Context, name, location types.String, objectLock types.Bool, tags types.Map, timeout time.Duration) error {
	createBucketConfig := objstorage.CreateBucketConfiguration{
		LocationConstraint: location.ValueStringPointer(),
	}

	if _, err := c.client.BucketsApi.CreateBucket(ctx, name.ValueString()).
		CreateBucketConfiguration(createBucketConfig).
		XAmzBucketObjectLockEnabled(objectLock.ValueBool()).
		Execute(); err != nil {
		return fmt.Errorf("failed to create user bucket: %w", err)
	}

	// Wait for bucket creation
	err := backoff.Retry(func() error {
		return c.bucketExistsCheck(ctx, name.ValueString())
	}, backoff.NewExponentialBackOff(backoff.WithMaxElapsedTime(timeout)))
	if err != nil {
		return fmt.Errorf("failed to wait for user bucket creation: %w", err)
	}

	if err = c.CreateBucketTags(ctx, name.ValueString(), tftags.NewFromMap(tags)); err != nil {
		return fmt.Errorf("failed to create user bucket tags: %w", err)
	}

	return nil
}

// GetBucket retrieves a bucket.
func (c *Client) GetBucket(ctx context.Context, name types.String) (*BucketModel, bool, error) {
	resp, err := c.client.BucketsApi.HeadBucket(ctx, name.ValueString()).Execute()
	if resp.HttpNotFound() {
		return nil, false, nil
	}

	if err != nil {
		return nil, false, fmt.Errorf("failed to get bucket: %w", err)
	}

	// location, err := c.GetBucketLocation(ctx, name)
	// if err != nil {
	// 	return nil, true, err
	// }

	lock, err := c.GetObjectLockEnabled(ctx, name)
	if err != nil {
		return nil, true, err
	}

	tags, err := c.ListBucketTags(ctx, name.ValueString())
	if err != nil {
		return nil, true, err
	}

	tagsMap, err := tags.ToMap(ctx)
	if err != nil {
		return nil, true, err
	}

	return &BucketModel{
		Name: name,
		// Region:            location,
		ObjectLockEnabled: lock,
		Tags:              tagsMap,
	}, true, nil
}

// GetBucketForDataSource retrieves a bucket for a data source.
func (c *Client) GetBucketForDataSource(ctx context.Context, name types.String) (*BucketDataSourceModel, bool, error) {
	resp, err := c.client.BucketsApi.HeadBucket(ctx, name.ValueString()).Execute()
	if resp.HttpNotFound() {
		return nil, false, nil
	}

	if err != nil {
		return nil, false, fmt.Errorf("failed to get bucket: %w", err)
	}

	location, err := c.GetBucketLocation(ctx, name)
	if err != nil {
		return nil, true, err
	}

	return &BucketDataSourceModel{
		Name:   name,
		Region: location,
	}, true, nil
}

// DeleteBucket deletes a bucket.
func (c *Client) DeleteBucket(ctx context.Context, name types.String, objectLockEnabled, forceDestroy types.Bool, timeout time.Duration) error {
	apiResponse, err := c.client.BucketsApi.DeleteBucket(ctx, name.ValueString()).Execute()
	if apiResponse.HttpNotFound() {
		return nil
	}

	if isBucketNotEmptyError(err) && forceDestroy.ValueBool() {
		if _, err = c.EmptyBucket(ctx, name.ValueString(), objectLockEnabled.ValueBool()); err != nil {
			return fmt.Errorf("failed to empty bucket: %w", err)
		}
		return c.DeleteBucket(ctx, name, objectLockEnabled, forceDestroy, timeout)
	}

	if err != nil {
		return err
	}

	backOff := backoff.NewExponentialBackOff(backoff.WithMaxElapsedTime(timeout))
	err = backoff.Retry(func() error {
		return c.bucketDeletedCheck(ctx, name.ValueString())
	}, backOff)

	if err != nil {
		return fmt.Errorf("failed to wait for bucket deletion: %w", err)
	}

	return nil
}

// GetBucketLocation should retrieve the location of a bucket. Currently does not work
func (c *Client) GetBucketLocation(ctx context.Context, name types.String) (types.String, error) {
	output, _, err := c.client.BucketsApi.GetBucketLocation(ctx, name.ValueString()).Execute()
	if err != nil {
		return types.StringNull(), fmt.Errorf("failed to get bucket location: %w", err)
	}

	return types.StringPointerValue(output.LocationConstraint), nil
}

func (c *Client) bucketDeletedCheck(ctx context.Context, name string) error {
	apiResponse, err := c.client.BucketsApi.HeadBucket(ctx, name).Execute()
	if apiResponse.HttpNotFound() {
		return nil
	}

	if err != nil {
		return backoff.Permanent(fmt.Errorf("failed to check if bucket exists: %w", err))
	}

	return fmt.Errorf("bucket still exists")
}

func (c *Client) bucketExistsCheck(ctx context.Context, name string) error {
	apiResponse, err := c.client.BucketsApi.HeadBucket(ctx, name).Execute()
	if apiResponse.HttpNotFound() {
		return fmt.Errorf("bucket not found")
	}

	if err != nil {
		return backoff.Permanent(fmt.Errorf("failed to check if bucket exists: %w", err))
	}

	return nil
}
