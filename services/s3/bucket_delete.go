package s3

import (
	"context"
	"errors"
	"fmt"
	"github.com/ionos-cloud/sdk-go-bundle/shared"

	s3 "github.com/ionos-cloud/sdk-go-s3"
)

const errAccessDenied = "AccessDenied"

// EmptyBucket deletes all objects and delete markers in the bucket.
// If `force` is `true` then S3 Object Lock governance mode restrictions are bypassed and
// an attempt is made to remove any S3 Object Lock legal holds.
// Returns the number of object versions and delete markers deleted.
func (c *Client) EmptyBucket(ctx context.Context, bucket string, forceDestroy bool) (int64, error) {
	objCount, err := c.forEachObjectVersionsPage(ctx, bucket, func(ctx context.Context, conn *s3.APIClient, bucket string, page *s3.ListObjectVersionsOutput) (int64, error) {
		return deletePageOfObjectVersions(ctx, conn, bucket, forceDestroy, page)
	})

	if err != nil {
		return objCount, err
	}

	n, err := c.forEachObjectVersionsPage(ctx, bucket, deletePageOfDeleteMarkers)
	objCount += n

	return objCount, err
}

func (c *Client) forEachObjectVersionsPage(ctx context.Context, bucket string, fn func(ctx context.Context, conn *s3.APIClient, bucket string, page *s3.ListObjectVersionsOutput) (int64, error)) (int64, error) {
	var objCount int64

	input := &ListObjectVersionsInput{
		Bucket: bucket,
	}
	var lastErr error

	pages := NewListObjectVersionsPaginator(c.client, input)
	for pages.HasMorePages() {
		page, err := pages.NextPage(ctx)

		if err != nil {
			return objCount, fmt.Errorf("listing S3 bucket (%s) object versions: %w", bucket, err)
		}

		n, err := fn(ctx, c.client, bucket, page)
		objCount += n

		if err != nil {
			lastErr = err
			break
		}
	}

	if lastErr != nil {
		return objCount, lastErr
	}

	return objCount, nil
}

func boolPtr(b bool) *bool {
	return &b
}

func getObjectsToDelete(page *s3.ListObjectVersionsOutput) []s3.ObjectIdentifier {
	if page.Versions == nil {
		return nil
	}

	toDelete := make([]s3.ObjectIdentifier, 0, len(*page.Versions))
	for _, v := range *page.Versions {
		toDelete = append(toDelete, s3.ObjectIdentifier{
			Key:       v.Key,
			VersionId: v.VersionId,
		})
	}

	return toDelete
}

func getDeleteMarkersToDelete(page *s3.ListObjectVersionsOutput) []s3.ObjectIdentifier {
	if page.DeleteMarkers == nil {
		return nil
	}

	toDelete := make([]s3.ObjectIdentifier, 0, len(*page.DeleteMarkers))
	for _, v := range *page.DeleteMarkers {
		toDelete = append(toDelete, s3.ObjectIdentifier{
			Key:       v.Key,
			VersionId: v.VersionId,
		})
	}

	return toDelete
}

func deletePageOfObjectVersions(ctx context.Context, conn *s3.APIClient, bucket string, force bool, page *s3.ListObjectVersionsOutput) (int64, error) {
	toDelete := getObjectsToDelete(page)
	var objCount int64
	if objCount = int64(len(toDelete)); objCount == 0 {
		return objCount, nil
	}

	req := conn.ObjectsApi.DeleteObjects(ctx, bucket).DeleteObjectsRequest(s3.DeleteObjectsRequest{
		Objects: &toDelete,
		Quiet:   boolPtr(true),
	})
	if force {
		req = req.XAmzBypassGovernanceRetention(true)
	}

	output, apiResponse, err := req.Execute()
	if apiResponse.HttpNotFound() {
		return objCount, nil
	}

	if err != nil {
		return objCount, fmt.Errorf("deleting S3 bucket (%s) object versions: %w", bucket, err)
	}

	if output.Errors == nil {
		return objCount, nil
	}

	objCount -= int64(len(*output.Errors))
	var errs []error
	for _, v := range *output.Errors {
		if force && shared.ToValueDefault(v.Code) == errAccessDenied {
			_, err := tryDisableLegalHold(ctx, conn, bucket, shared.ToValueDefault(v.Key), shared.ToValueDefault(v.VersionId))
			if err != nil {
				errs = append(errs, []error{newDeleteObjectVersionError(v), fmt.Errorf("removing legal hold: %w",
					newObjectVersionError(shared.ToValueDefault(v.Key), shared.ToValueDefault(v.VersionId), err))}...)
			} else {
				if _, err = deleteObject(ctx, conn, &DeleteRequest{
					Bucket:       bucket,
					Key:          *v.Key,
					VersionID:    *v.VersionId,
					ForceDestroy: force,
				}); err != nil {
					errs = append(errs, newObjectVersionError(shared.ToValueDefault(v.Key), shared.ToValueDefault(v.VersionId), err))
				} else {
					objCount++
				}
			}
		} else {
			errs = append(errs, newDeleteObjectVersionError(v))
		}
	}
	if err := errors.Join(errs...); err != nil {
		return objCount, fmt.Errorf("deleting S3 bucket (%s) object versions: %w", bucket, err)
	}

	return objCount, nil
}

func toString(v *string) string {
	if v == nil {
		return ""
	}

	return *v
}

func deletePageOfDeleteMarkers(ctx context.Context, conn *s3.APIClient, bucket string, page *s3.ListObjectVersionsOutput) (int64, error) {
	toDelete := getDeleteMarkersToDelete(page)
	var objCount int64
	if objCount = int64(len(toDelete)); objCount == 0 {
		return objCount, nil
	}

	output, apiResponse, err := conn.ObjectsApi.DeleteObjects(ctx, bucket).DeleteObjectsRequest(s3.DeleteObjectsRequest{
		Objects: &toDelete,
		Quiet:   boolPtr(true),
	}).Execute()
	if apiResponse.HttpNotFound() {
		return objCount, nil
	}

	if err != nil {
		return objCount, fmt.Errorf("deleting S3 bucket (%s) object versions: %w", bucket, err)
	}

	if output.Errors == nil {
		return objCount, nil
	}

	objCount -= int64(len(*output.Errors))
	errs := make([]error, 0, len(*output.Errors))
	for _, v := range *output.Errors {
		errs = append(errs, newDeleteObjectVersionError(v))
	}

	if err := errors.Join(errs...); err != nil {
		return objCount, fmt.Errorf("deleting S3 bucket (%s) delete markers: %w", bucket, err)
	}

	return objCount, nil
}

func newDeleteObjectVersionError(err s3.DeletionError) error {
	s3Err := fmt.Errorf("%s: %s", *err.Code, *err.Message)

	return fmt.Errorf("deleting: %w", newObjectVersionError(*err.Key, *err.VersionId, s3Err))
}

func newObjectVersionError(key, versionID string, err error) error {
	if err == nil {
		return nil
	}

	if versionID == "" {
		return fmt.Errorf("S3 object (%s): %w", key, err)
	}

	return fmt.Errorf("S3 object (%s) version (%s): %w", key, versionID, err)
}
