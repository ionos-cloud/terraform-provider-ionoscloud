package identity

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/list"
)

// StreamResults streams a slice of items of type T into the list results stream.
// The mapper receives a pre-created *list.ListResult to populate and returns whether
// to push the result (false = skip). If diagnostics contain errors, iteration stops.
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
			if !shouldPush {
				continue
			}
			result.Diagnostics.Append(diags...)
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
