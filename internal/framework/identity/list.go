package identity

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/list"
)

// StreamResults streams a slice of items of type T into the list results stream.
// It sets stream.Results to an iterator function that processes each item with the mapper.
// If mapping returns diagnostics with errors, they are pushed to the stream and execution stops.
func StreamResults[T any](
	ctx context.Context,
	stream *list.ListResultsStream,
	req list.ListRequest,
	items []T,
	mapper func(context.Context, list.ListRequest, T, *list.ListResult) (bool, diag.Diagnostics),
) {
	stream.Results = func(push func(list.ListResult) bool) {
		for _, item := range items {
			result := req.NewListResult(ctx)
			shouldPush, diags := mapper(ctx, req, item, &result)
			if diags.HasError() {
				result.Diagnostics.Append(diags...)
				push(result)
				return
			}
			result.Diagnostics.Append(diags...)
			if !shouldPush {
				continue
			}
			if !push(result) {
				return
			}
		}
	}
}

// StreamError sets the results stream to emit a single diagnostic error.
func StreamError(stream *list.ListResultsStream, summary, detail string) {
	var diags diag.Diagnostics
	diags.AddError(summary, detail)
	stream.Results = list.ListResultsStreamDiagnostics(diags)
}
