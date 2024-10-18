package objectstorage

import (
	"context"
	"fmt"

	objstorage "github.com/ionos-cloud/sdk-go-s3"
)

// ListObjectVersionsInput represents the input for the ListObjectVersionsPaginator.
type ListObjectVersionsInput struct {
	Bucket  string
	Prefix  string
	MaxKeys int
}

// ListObjectVersionsPaginator is a paginator for listing object versions in a bucket.
type ListObjectVersionsPaginator struct {
	client          *objstorage.APIClient
	input           *ListObjectVersionsInput
	keyMarker       *string
	versionIDMarker *string
	hasMore         bool
}

// NewListObjectVersionsPaginator creates a new ListObjectVersionsPaginator.
func NewListObjectVersionsPaginator(client *objstorage.APIClient, input *ListObjectVersionsInput) *ListObjectVersionsPaginator {
	if input.MaxKeys == 0 {
		input.MaxKeys = 1000
	}

	return &ListObjectVersionsPaginator{
		client:  client,
		input:   input,
		hasMore: true,
	}
}

// HasMorePages returns true if there are more pages to retrieve.
func (p *ListObjectVersionsPaginator) HasMorePages() bool {
	return p.hasMore
}

// NextPage retrieves the next page of object versions.
func (p *ListObjectVersionsPaginator) NextPage(ctx context.Context) (*objstorage.ListObjectVersionsOutput, error) {
	if !p.hasMore {
		return nil, fmt.Errorf("no more pages")
	}

	req := p.client.VersionsApi.ListObjectVersions(ctx, p.input.Bucket).
		Prefix(p.input.Prefix).
		MaxKeys(int32(p.input.MaxKeys))

	if p.keyMarker != nil {
		req = req.KeyMarker(*p.keyMarker)
	}

	if p.versionIDMarker != nil {
		req = req.VersionIdMarker(*p.versionIDMarker)
	}

	output, _, err := req.Execute()
	if err != nil {
		return nil, err
	}

	if output.IsTruncated != nil && *output.IsTruncated {
		p.keyMarker = output.NextKeyMarker
		p.versionIDMarker = output.NextVersionIdMarker
	} else {
		p.hasMore = false
	}

	return output, nil
}
