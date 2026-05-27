package identity

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/list"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// StreamList fetches items using fetch and streams them into the results stream.
// If the fetch fails, the error is emitted as a diagnostic and iteration stops.
//
// The mapper contract:
//   - Return nil to skip the item; any diagnostics are silently dropped (per-item error).
//   - Return a non-nil *MappedItem to include the item; MappedItem.Identity must be
//     non-nil. StreamList sets DisplayName, Identity, and Resource on the result.
//   - Errors during result population (Identity.Set, Resource.Set) are fatal —
//     the error result is pushed and iteration stops.
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
			result.Diagnostics.Append(diags...)
			if mapped == nil {
				if diags.HasError() {
					tflog.Error(ctx, "skipping item due to mapper error", map[string]any{
						"error": diags[0].Detail(),
					})
				}
				continue
			}
			if mapped.Identity == nil {
				result.Diagnostics.AddError("mapper contract violation", "MappedItem.Identity must not be nil")
				push(result)
				return
			}
			result.DisplayName = mapped.DisplayName
			result.Diagnostics.Append(result.Identity.Set(ctx, mapped.Identity)...)
			if result.Diagnostics.HasError() {
				push(result)
				return
			}
			if mapped.Resource != nil {
				result.Diagnostics.Append(result.Resource.Set(ctx, mapped.Resource)...)
				if result.Diagnostics.HasError() {
					push(result)
					return
				}
			}
			if !push(result) {
				return
			}
		}
	}
}
