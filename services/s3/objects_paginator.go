package s3

import (
	"context"
	"fmt"

	s3 "github.com/ionos-cloud/sdk-go-s3"
)

// ListObjectsV2Input represents the input for the ListObjectsV2Paginator.
type ListObjectsV2Input struct {
	Bucket            string
	Prefix            *string
	Delimiter         *string
	EncodingType      *string
	ContinuationToken *string
	FetchOwner        bool
	StartAfter        *string
	MaxKeys           *int32
}

// ListObjectsV2Paginator is a paginator for listing objects from S3 bucket
type ListObjectsV2Paginator struct {
	client  *s3.APIClient
	input   *ListObjectsV2Input
	hasMore bool
}

// NewListObjectsV2Paginator creates a new ListObjectsV2Paginator.
func NewListObjectsV2Paginator(client *s3.APIClient, input *ListObjectsV2Input) *ListObjectsV2Paginator {
	if input.MaxKeys == nil {
		defaultValue := int32(1000)
		input.MaxKeys = &defaultValue
	}

	return &ListObjectsV2Paginator{
		client:  client,
		input:   input,
		hasMore: true,
	}
}

// HasMorePages returns true if there are more pages to retrieve.
func (p *ListObjectsV2Paginator) HasMorePages() bool {
	return p.hasMore
}

// NextPage retrieves the next page of objects.
func (p *ListObjectsV2Paginator) NextPage(ctx context.Context) (*s3.ListBucketResultV2, error) {
	if !p.hasMore {
		return nil, fmt.Errorf("no more pages")
	}

	req := p.client.ObjectsApi.ListObjectsV2(ctx, p.input.Bucket)
	if p.input.EncodingType != nil {
		req = req.EncodingType(*p.input.EncodingType)
	}

	if p.input.Prefix != nil {
		req = req.Prefix(*p.input.Prefix)
	}

	if p.input.Delimiter != nil {
		req = req.Delimiter(*p.input.Delimiter)
	}

	if p.input.ContinuationToken != nil {
		req = req.ContinuationToken(*p.input.ContinuationToken)
	}

	if p.input.FetchOwner {
		req = req.FetchOwner(p.input.FetchOwner)
	}

	if p.input.StartAfter != nil {
		req = req.StartAfter(*p.input.StartAfter)
	}

	if p.input.MaxKeys != nil {
		req = req.MaxKeys(*p.input.MaxKeys)
	}

	output, _, err := req.Execute()
	if err != nil {
		return nil, err
	}

	if output.IsTruncated != nil && *output.IsTruncated {
		p.input.ContinuationToken = output.NextContinuationToken
	} else {
		p.hasMore = false
	}

	return output, nil
}
