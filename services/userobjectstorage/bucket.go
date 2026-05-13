package userobjectstorage

import (
	"context"
	"crypto/md5" //nolint:gosec
	"encoding/base64"
	"encoding/xml"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/cenkalti/backoff/v4"
	"github.com/ionos-cloud/sdk-go-bundle/products/userobjectstorage/v2"
	"github.com/ionos-cloud/sdk-go-bundle/shared"

	diagutil "github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils/diags"
)

// CreateBucket creates a new user-owned bucket in the given region.
func (c *Client) CreateBucket(ctx context.Context, name, region string, objectLock bool, timeout time.Duration) error {
	client, err := c.apiClientForRegion(region)
	if err != nil {
		return err
	}

	createBucketConfig := userobjectstorage.CreateBucketRequestCreateBucketConfiguration{
		LocationConstraint: &region,
	}
	createReq := userobjectstorage.NewCreateBucketRequest()
	createReq.CreateBucketConfiguration = &createBucketConfig
	if _, err := client.BucketsApi.CreateBucket(ctx, name).
		CreateBucketRequest(*createReq).
		XAmzBucketObjectLockEnabled(objectLock).
		Execute(); err != nil {
		return fmt.Errorf("failed to create bucket %q: %w", name, err)
	}

	backoffErr := backoff.Retry(func() error {
		return c.bucketExistsCheck(ctx, name)
	}, backoff.NewExponentialBackOff(backoff.WithMaxElapsedTime(timeout)))
	if backoffErr != nil {
		return fmt.Errorf("failed to wait for bucket creation: %w", backoffErr)
	}

	return nil
}

// GetBucket checks whether a bucket exists. Returns (false, nil) if the bucket is not found.
// Region is sourced from Terraform state — the SDK's GetBucketLocation response model has a
// generator-level XML tag issue that cannot be fixed without changing the generator templates.
func (c *Client) GetBucket(ctx context.Context, name, region string) (bool, error) {
	client, err := c.apiClientForRegion(region)
	if err != nil {
		return false, err
	}

	apiResp, err := client.BucketsApi.HeadBucket(ctx, name).Execute()
	if apiResp.HttpNotFound() {
		return false, nil
	}
	if err != nil {
		return false, fmt.Errorf("failed to check bucket %q: %w", name, err)
	}

	return true, nil
}

// GetObjectLockEnabled returns true if Object Lock is enabled for the bucket.
// Returns false when the bucket has no Object Lock configuration (404) or when it is not enabled.
func (c *Client) GetObjectLockEnabled(ctx context.Context, name string) (bool, error) {
	// Object lock configuration is not region-specific at the API level; use the default region.
	api, err := c.apiClientForRegion(DefaultRegion)
	if err != nil {
		return false, err
	}

	output, apiResp, err := api.ObjectLockApi.GetObjectLockConfiguration(ctx, name).Execute()
	if apiResp.HttpNotFound() {
		return false, nil
	}
	if err != nil {
		return false, fmt.Errorf("failed to get object lock configuration for %q: %w", name, err)
	}
	return output.GetObjectLockEnabled() == "Enabled", nil
}

// GetBucketTags returns the tags for a bucket. Returns nil when the bucket has no tags.
func (c *Client) GetBucketTags(ctx context.Context, name, region string) (map[string]string, error) {
	api, err := c.apiClientForRegion(region)
	if err != nil {
		return nil, err
	}
	output, apiResp, err := api.TaggingApi.GetBucketTagging(ctx, name).Execute()
	if apiResp.HttpNotFound() {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get tags for bucket %q: %w", name, err)
	}

	tagSet := output.GetTagSet()
	sdkTags := tagSet.GetTag()
	tags := make(map[string]string, len(sdkTags))
	for _, tag := range sdkTags {
		tags[tag.Key] = tag.Value
	}
	return tags, nil
}

// PutBucketTags sets the complete tag set for a bucket, replacing any existing tags.
func (c *Client) PutBucketTags(ctx context.Context, name, region string, tags map[string]string) error {
	api, err := c.apiClientForRegion(region)
	if err != nil {
		return err
	}

	tagSet := make([]userobjectstorage.Tag, 0, len(tags))
	for k, v := range tags {
		tagSet = append(tagSet, *userobjectstorage.NewTag(k, v))
	}

	req := userobjectstorage.NewPutBucketTaggingRequest(tagSet)
	body, err := xml.Marshal(req)
	if err != nil {
		return fmt.Errorf("failed to marshal tags for bucket %q: %w", name, err)
	}
	//nolint:gosec
	sum := md5.Sum(body)
	contentMD5 := base64.StdEncoding.EncodeToString(sum[:])

	if _, err := api.TaggingApi.PutBucketTagging(ctx, name).
		PutBucketTaggingRequest(*req).
		ContentMD5(contentMD5).Execute(); err != nil {
		return fmt.Errorf("failed to put tags for bucket %q: %w", name, err)
	}
	return nil
}

// DeleteBucketTags removes all tags from a bucket.
func (c *Client) DeleteBucketTags(ctx context.Context, name, region string) error {
	api, err := c.apiClientForRegion(region)
	if err != nil {
		return err
	}
	if _, err := api.TaggingApi.DeleteBucketTagging(ctx, name).Execute(); err != nil {
		return fmt.Errorf("failed to delete tags for bucket %q: %w", name, err)
	}
	return nil
}

// DeleteBucket deletes a bucket. If forceDestroy is true and the bucket is not empty, all objects
// are deleted first.
func (c *Client) DeleteBucket(ctx context.Context, name, region string, forceDestroy bool, timeout time.Duration) error {
	api, err := c.apiClientForRegion(region)
	if err != nil {
		return err
	}

	apiResp, err := api.BucketsApi.DeleteBucket(ctx, name).Execute()
	if apiResp.HttpNotFound() {
		return nil
	}

	if isBucketNotEmptyError(err) && forceDestroy {
		if err = c.deleteAllObjects(ctx, name); err != nil {
			return fmt.Errorf("failed to empty bucket %q: %w", name, err)
		}
		return c.DeleteBucket(ctx, name, region, forceDestroy, timeout)
	}

	if err != nil {
		return fmt.Errorf("failed to delete bucket %q: %w", name, err)
	}

	backOff := backoff.NewExponentialBackOff(backoff.WithMaxElapsedTime(timeout))
	if err = backoff.Retry(func() error {
		return c.bucketDeletedCheck(ctx, name)
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
	return errors.As(err, &apiErr) && apiErr.StatusCode() == http.StatusConflict
}
