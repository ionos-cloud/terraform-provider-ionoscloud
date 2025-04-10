package objectstorage

import (
	"context"
	"errors"
	"fmt"

	"github.com/ionos-cloud/sdk-go-bundle/shared"

	objstorage "github.com/ionos-cloud/sdk-go-bundle/products/objectstorage/v2"
)

const errAccessDenied = "AccessDenied"

// EmptyBucket deletes all objects and delete markers in the bucket.
// If `force` is `true` then Object Lock governance mode restrictions are bypassed and
// an attempt is made to remove any Object Lock legal holds.
// Returns the number of object versions and delete markers deleted.
func (c *Client) EmptyBucket(ctx context.Context, bucket string, forceDestroy bool) (int64, error) {
	objCount, err := c.forEachObjectVersionsPage(ctx, bucket, func(ctx context.Context, conn *objstorage.APIClient, bucket string, page *objstorage.ListObjectVersionsOutput) (int64, error) {
		return deletePageOfObjectVersions(ctx, conn, bucket, forceDestroy, page)
	})

	if err != nil {
		return objCount, err
	}

	n, err := c.forEachObjectVersionsPage(ctx, bucket, deletePageOfDeleteMarkers)
	objCount += n

	return objCount, err
}

func (c *Client) forEachObjectVersionsPage(ctx context.Context, bucket string, fn func(ctx context.Context, conn *objstorage.APIClient, bucket string, page *objstorage.ListObjectVersionsOutput) (int64, error)) (int64, error) {
	var objCount int64

	input := &ListObjectVersionsInput{
		Bucket: bucket,
	}
	var lastErr error

	pages := NewListObjectVersionsPaginator(c.client, input)
	for pages.HasMorePages() {
		page, err := pages.NextPage(ctx)

		if err != nil {
			return objCount, fmt.Errorf("listing bucket (%s) object versions: %w", bucket, err)
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

func getObjectsToDelete(page *objstorage.ListObjectVersionsOutput) []objstorage.ObjectIdentifier {
	if page.Versions == nil {
		return nil
	}

	toDelete := make([]objstorage.ObjectIdentifier, 0, len(page.Versions))
	for _, v := range page.Versions {
		toDelete = append(toDelete, objstorage.ObjectIdentifier{
			Key:       *v.Key,
			VersionId: v.VersionId,
		})
	}

	return toDelete
}

func getDeleteMarkersToDelete(page *objstorage.ListObjectVersionsOutput) []objstorage.ObjectIdentifier {
	if page.DeleteMarkers == nil {
		return nil
	}

	toDelete := make([]objstorage.ObjectIdentifier, 0, len(page.DeleteMarkers))
	for _, v := range page.DeleteMarkers {
		toDelete = append(toDelete, objstorage.ObjectIdentifier{
			Key:       *v.Key,
			VersionId: v.VersionId,
		})
	}

	return toDelete
}

func deletePageOfObjectVersions(ctx context.Context, conn *objstorage.APIClient, bucket string, force bool, page *objstorage.ListObjectVersionsOutput) (int64, error) {
	toDelete := getObjectsToDelete(page)
	var objCount int64
	if objCount = int64(len(toDelete)); objCount == 0 {
		return objCount, nil
	}

	req := conn.ObjectsApi.DeleteObjects(ctx, bucket).DeleteObjectsRequest(objstorage.DeleteObjectsRequest{
		Objects: toDelete,
		Quiet:   shared.ToPtr(true),
	})
	if force {
		req = req.XAmzBypassGovernanceRetention(true)
	}

	output, apiResponse, err := req.Execute()
	if apiResponse.HttpNotFound() {
		return objCount, nil
	}

	if err != nil {
		return objCount, fmt.Errorf("deleting bucket (%s) object versions: %w", bucket, err)
	}

	if output.Errors == nil {
		return objCount, nil
	}

	objCount -= int64(len(output.Errors))
	var errs []error
	for _, v := range output.Errors {
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
		return objCount, fmt.Errorf("deleting bucket (%s) object versions: %w", bucket, err)
	}

	return objCount, nil
}

func deletePageOfDeleteMarkers(ctx context.Context, conn *objstorage.APIClient, bucket string, page *objstorage.ListObjectVersionsOutput) (int64, error) {
	toDelete := getDeleteMarkersToDelete(page)
	var objCount int64
	if objCount = int64(len(toDelete)); objCount == 0 {
		return objCount, nil
	}

	output, apiResponse, err := conn.ObjectsApi.DeleteObjects(ctx, bucket).DeleteObjectsRequest(objstorage.DeleteObjectsRequest{
		Objects: toDelete,
		Quiet:   shared.ToPtr(true),
	}).Execute()
	if apiResponse.HttpNotFound() {
		return objCount, nil
	}

	if err != nil {
		return objCount, fmt.Errorf("deleting bucket (%s) object versions: %w", bucket, err)
	}

	if output.Errors == nil {
		return objCount, nil
	}

	objCount -= int64(len(output.Errors))
	errs := make([]error, 0, len(output.Errors))
	for _, v := range output.Errors {
		errs = append(errs, newDeleteObjectVersionError(v))
	}

	if err := errors.Join(errs...); err != nil {
		return objCount, fmt.Errorf("deleting bucket (%s) delete markers: %w", bucket, err)
	}

	return objCount, nil
}

func newDeleteObjectVersionError(err objstorage.DeletionError) error {
	sErr := fmt.Errorf("%s: %s", *err.Code, *err.Message)
	return fmt.Errorf("deleting: %w", newObjectVersionError(*err.Key, *err.VersionId, sErr))
}

func newObjectVersionError(key, versionID string, err error) error {
	if err == nil {
		return nil
	}

	if versionID == "" {
		return fmt.Errorf("object (%s): %w", key, err)
	}

	return fmt.Errorf("object (%s) version (%s): %w", key, versionID, err)
}
