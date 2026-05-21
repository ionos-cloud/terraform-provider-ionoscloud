package identity

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/list"
)

// StreamList fetches items using fetch and streams them into the results stream.
// If the fetch fails, the error is emitted as a diagnostic and iteration stops.
// The mapper returns a *MappedItem to populate each result with; nil skips the item.
func StreamList[T any](
	ctx context.Context,
	stream *list.ListResultsStream,
	req list.ListRequest,
	fetch func(context.Context) ([]T, error),
	mapper func(context.Context, bool, T) (*MappedItem, diag.Diagnostics),
) {
	items, err := fetch(ctx)
	if err != nil {
		var diags diag.Diagnostics
		diags.AddError("Failed to list resources", err.Error())
		stream.Results = list.ListResultsStreamDiagnostics(diags)
		return
	}
	stream.Results = func(push func(list.ListResult) bool) {
		for _, item := range items {
			result := req.NewListResult(ctx)
			mapped, diags := mapper(ctx, req.IncludeResource, item)
			if diags.HasError() {
				result.Diagnostics.Append(diags...)
				push(result)
				return
			}
			if mapped == nil {
				continue
			}
			result.DisplayName = mapped.DisplayName
			result.Diagnostics.Append(result.Identity.Set(ctx, mapped.Identity)...)
			if mapped.Resource != nil {
				result.Diagnostics.Append(result.Resource.Set(ctx, mapped.Resource)...)
			}
			if !push(result) {
				return
			}
		}
	}
}
