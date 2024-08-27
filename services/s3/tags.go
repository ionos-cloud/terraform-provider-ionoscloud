package s3

import (
	"context"
	"fmt"

	s3 "github.com/ionos-cloud/sdk-go-s3"

	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/internal/tags"
)

// CreateBucketTags creates tags for a bucket.
func (c *Client) CreateBucketTags(ctx context.Context, bucketName string, tags tags.KeyValueTags) error {
	if len(tags) == 0 {
		return nil
	}

	return c.UpdateBucketTags(ctx, bucketName, tags, nil)
}

// UpdateBucketTags updates tags for a bucket.
func (c *Client) UpdateBucketTags(ctx context.Context, bucketName string, new, old tags.KeyValueTags) error {
	allTags, err := c.ListBucketTags(ctx, bucketName)
	if err != nil {
		return err
	}

	ignoredTags := allTags.Ignore(old).Ignore(new)
	if len(new)+len(ignoredTags) > 0 {
		if _, err = c.client.TaggingApi.PutBucketTagging(ctx, bucketName).PutBucketTaggingRequest(
			s3.PutBucketTaggingRequest{
				TagSet: new.Merge(ignoredTags).ToListPointer(),
			}).Execute(); err != nil {
			return fmt.Errorf("failed to update bucket tags: %w", err)
		}
	} else if len(old) > 0 && len(ignoredTags) == 0 {
		if _, err = c.client.TaggingApi.DeleteBucketTagging(ctx, bucketName).Execute(); err != nil {
			return fmt.Errorf("failed to delete bucket tags: %w", err)
		}
	}

	return nil
}

// ListBucketTags lists tags for a bucket.
func (c *Client) ListBucketTags(ctx context.Context, bucketName string) (tags.KeyValueTags, error) {
	output, apiResponse, err := c.client.TaggingApi.GetBucketTagging(ctx, bucketName).Execute()
	if apiResponse.HttpNotFound() {
		return tags.New(nil), nil
	}

	if err != nil {
		return tags.New(nil), fmt.Errorf("failed to get bucket tags: %w", err)
	}

	if output.TagSet == nil {
		return tags.New(nil), nil
	}

	return tags.New(*output.TagSet), nil
}

// ListObjectTags lists tags for an object.
func (c *Client) ListObjectTags(ctx context.Context, bucketName, objectName string) (tags.KeyValueTags, error) {
	output, apiResponse, err := c.client.TaggingApi.GetObjectTagging(ctx, bucketName, objectName).Execute()
	if apiResponse.HttpNotFound() {
		return tags.New(nil), nil
	}

	if err != nil {
		return tags.New(nil), fmt.Errorf("failed to get object tags: %w", err)
	}

	if output.TagSet == nil {
		return tags.New(nil), nil
	}

	return tags.New(*output.TagSet), nil
}

// UpdateObjectTags updates tags for an object.
func (c *Client) UpdateObjectTags(ctx context.Context, bucketName, objectName string, new, old tags.KeyValueTags) error {
	allTags, err := c.ListObjectTags(ctx, bucketName, objectName)
	if err != nil {
		return err
	}

	ignoredTags := allTags.Ignore(old).Ignore(new)
	if len(new)+len(ignoredTags) > 0 {
		if _, _, err = c.client.TaggingApi.PutObjectTagging(ctx, bucketName, objectName).PutObjectTaggingRequest(
			s3.PutObjectTaggingRequest{
				TagSet: new.Merge(ignoredTags).ToListPointer(),
			}).Execute(); err != nil {
			return fmt.Errorf("failed to update object tags: %w", err)
		}
	} else if len(old) > 0 && len(ignoredTags) == 0 {
		if _, _, err = c.client.TaggingApi.DeleteObjectTagging(ctx, bucketName, objectName).Execute(); err != nil {
			return fmt.Errorf("failed to delete object tags: %w", err)
		}
	}

	return nil
}
