package s3

import (
	"context"
	"fmt"

	s3 "github.com/ionos-cloud/sdk-go-s3"

	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/internal/tags"
)

// CreateBucketTags creates tags for a bucket.
func CreateBucketTags(ctx context.Context, client *s3.APIClient, bucketName string, tags tags.KeyValueTags) error {
	if len(tags) == 0 {
		return nil
	}

	return UpdateBucketTags(ctx, client, bucketName, tags, nil)
}

// UpdateBucketTags updates tags for a bucket.
func UpdateBucketTags(ctx context.Context, client *s3.APIClient, bucketName string, new, old tags.KeyValueTags) error {
	allTags, err := ListBucketTags(ctx, client, bucketName)
	if err != nil {
		return err
	}

	ignoredTags := allTags.Ignore(old).Ignore(new)
	if len(new)+len(ignoredTags) > 0 {
		if _, err = client.TaggingApi.PutBucketTagging(ctx, bucketName).PutBucketTaggingRequest(
			s3.PutBucketTaggingRequest{
				TagSet: new.Merge(ignoredTags).ToListPointer(),
			}).Execute(); err != nil {
			return fmt.Errorf("failed to update bucket tags: %w", err)
		}
	} else if len(old) > 0 && len(ignoredTags) == 0 {
		if _, err = client.TaggingApi.DeleteBucketTagging(ctx, bucketName).Execute(); err != nil {
			return fmt.Errorf("failed to delete bucket tags: %w", err)
		}
	}

	return nil
}

// ListBucketTags lists tags for a bucket.
func ListBucketTags(ctx context.Context, client *s3.APIClient, bucketName string) (tags.KeyValueTags, error) {
	output, apiResponse, err := client.TaggingApi.GetBucketTagging(ctx, bucketName).Execute()
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
func ListObjectTags(ctx context.Context, client *s3.APIClient, bucketName, objectName string) (tags.KeyValueTags, error) {
	output, apiResponse, err := client.TaggingApi.GetObjectTagging(ctx, bucketName, objectName).Execute()
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
func UpdateObjectTags(ctx context.Context, client *s3.APIClient, bucketName, objectName string, new, old tags.KeyValueTags) error {
	allTags, err := ListObjectTags(ctx, client, bucketName, objectName)
	if err != nil {
		return err
	}

	ignoredTags := allTags.Ignore(old).Ignore(new)
	if len(new)+len(ignoredTags) > 0 {
		if _, _, err = client.TaggingApi.PutObjectTagging(ctx, bucketName, objectName).PutObjectTaggingRequest(
			s3.PutObjectTaggingRequest{
				TagSet: new.Merge(ignoredTags).ToListPointer(),
			}).Execute(); err != nil {
			return fmt.Errorf("failed to update object tags: %w", err)
		}
	} else if len(old) > 0 && len(ignoredTags) == 0 {
		if _, _, err = client.TaggingApi.DeleteObjectTagging(ctx, bucketName, objectName).Execute(); err != nil {
			return fmt.Errorf("failed to delete object tags: %w", err)
		}
	}

	return nil
}
