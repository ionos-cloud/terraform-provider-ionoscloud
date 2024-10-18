package objectstorage

import (
	"context"
	"fmt"
	"net/http"

	"github.com/ionos-cloud/sdk-go-bundle/shared"

	objstorage "github.com/ionos-cloud/sdk-go-s3"
)

// DeleteRequest represents a request to delete an object from a general purpose bucket.
type DeleteRequest struct {
	Bucket       string
	Key          string
	VersionID    string
	ForceDestroy bool
}

// DeleteAllObjectVersions deletes all versions of a specified key from a general purpose bucket.
// If key is empty then all versions of all objects are deleted.
// Set `force` to `true` to override any object lock protections on object lock enabled buckets.
// Returns the number of objects deleted.
// Use `emptyBucket` to delete all versions of all objects in a bucket.
func DeleteAllObjectVersions(ctx context.Context, client *objstorage.APIClient, req *DeleteRequest) (int, error) {
	var (
		objCount int
		lastErr  error
	)

	pages := NewListObjectVersionsPaginator(client, &ListObjectVersionsInput{
		Bucket: req.Bucket,
		Prefix: req.Key,
	})

	for pages.HasMorePages() {
		page, err := pages.NextPage(ctx)
		if err != nil {
			return objCount, err
		}

		var count int
		count, lastErr = deleteVersionsPage(ctx, client, page.Versions, req.Bucket, req.Key, req.ForceDestroy)
		objCount += count
	}

	if lastErr != nil {
		return objCount, lastErr
	}

	pages = NewListObjectVersionsPaginator(client, &ListObjectVersionsInput{
		Bucket: req.Bucket,
		Prefix: req.Key,
	})

	for pages.HasMorePages() {
		page, err := pages.NextPage(ctx)
		if err != nil {
			return objCount, err
		}

		var count int
		count, lastErr = deleteMarkersPage(ctx, client, page.DeleteMarkers, req.Bucket, req.Key)
		objCount += count
	}

	if lastErr != nil {
		return objCount, fmt.Errorf("deleting at least one object delete marker, last error: %w", lastErr)
	}

	return objCount, nil
}

func deleteObject(ctx context.Context, client *objstorage.APIClient, req *DeleteRequest) (*objstorage.APIResponse, error) {
	r := client.ObjectsApi.DeleteObject(ctx, req.Bucket, req.Key)
	if req.VersionID != "" {
		r = r.VersionId(req.VersionID)
	}

	if req.ForceDestroy {
		r = r.XAmzBypassGovernanceRetention(true)
	}

	_, apiResponse, err := r.Execute()
	return apiResponse, err
}

func deleteVersionsPage(ctx context.Context, client *objstorage.APIClient, versions *[]objstorage.ObjectVersion, bucket, key string, force bool) (int, error) {
	var (
		objCount int
		lastErr  error
	)

	if versions == nil {
		return objCount, nil
	}

	for _, v := range *versions {
		if key != shared.ToValueDefault(v.Key) {
			continue
		}

		apiResponse, err := deleteObject(ctx, client, &DeleteRequest{
			Bucket:       bucket,
			Key:          key,
			VersionID:    shared.ToValueDefault(v.VersionId),
			ForceDestroy: force,
		})

		if err == nil {
			objCount++
			continue
		}

		// If the object is locked by Object Lock, we need to disable the legal hold before deleting it
		if httpForbidden(apiResponse) && force {
			success, err := tryDisableLegalHold(ctx, client, bucket, key, shared.ToValueDefault(v.VersionId))
			if err != nil {
				lastErr = err
				continue
			}

			if !success {
				lastErr = err
				continue
			}

			if _, err = deleteObject(ctx, client, &DeleteRequest{
				Bucket:       bucket,
				Key:          key,
				VersionID:    shared.ToValueDefault(v.VersionId),
				ForceDestroy: force,
			}); err != nil {
				lastErr = err
				continue
			}

			objCount++
		} else {
			lastErr = err
		}
	}

	return objCount, lastErr
}

func deleteMarkersPage(ctx context.Context, client *objstorage.APIClient, markers *[]objstorage.DeleteMarkerEntry, bucket, key string) (int, error) {
	var (
		objCount int
		lastErr  error
	)

	if markers == nil {
		return objCount, nil
	}

	for _, m := range *markers {
		if key != shared.ToValueDefault(m.Key) {
			continue
		}

		_, err := deleteObject(ctx, client, &DeleteRequest{
			Bucket:    bucket,
			Key:       key,
			VersionID: *m.VersionId,
		})

		if err == nil {
			objCount++
			continue
		}

		lastErr = err
	}

	return objCount, lastErr

}

func tryDisableLegalHold(ctx context.Context, client *objstorage.APIClient, bucket, key, versionID string) (bool, error) {
	output, _, err := client.ObjectLockApi.GetObjectLegalHold(ctx, bucket, key).VersionId(versionID).Execute()
	if err != nil {
		return false, err
	}

	if *output.Status == "OFF" {
		return false, nil
	}

	_, err = client.ObjectLockApi.PutObjectLegalHold(ctx, bucket, key).VersionId(versionID).
		ObjectLegalHoldConfiguration(objstorage.ObjectLegalHoldConfiguration{
			Status: shared.ToPtr("OFF"),
		}).Execute()

	if err != nil {
		return false, err
	}

	return true, nil
}

func httpForbidden(response *objstorage.APIResponse) bool {
	return response != nil && response.Response != nil && response.StatusCode == http.StatusForbidden
}
