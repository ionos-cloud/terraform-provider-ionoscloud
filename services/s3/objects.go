package s3

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils/convptr"
)

// ObjectsDataSourceModel is used to fetch objects from a bucket.
type ObjectsDataSourceModel struct {
	Bucket         types.String   `tfsdk:"bucket"`
	Delimiter      types.String   `tfsdk:"delimiter"`
	EncodingType   types.String   `tfsdk:"encoding_type"`
	MaxKeys        types.Int64    `tfsdk:"max_keys"`
	Prefix         types.String   `tfsdk:"prefix"`
	FetchOwner     types.Bool     `tfsdk:"fetch_owner"`
	StartAfter     types.String   `tfsdk:"start_after"`
	CommonPrefixes []types.String `tfsdk:"common_prefixes"`
	Keys           []types.String `tfsdk:"keys"`
	Owners         []types.String `tfsdk:"owners"`
}

// ListObjects fetches objects from a bucket.
func (c *Client) ListObjects(ctx context.Context, data *ObjectsDataSourceModel) error {
	input := &ListObjectsV2Input{
		Bucket:       data.Bucket.ValueString(),
		Delimiter:    data.Delimiter.ValueStringPointer(),
		EncodingType: data.EncodingType.ValueStringPointer(),
		MaxKeys:      convptr.Int64ToInt32(data.MaxKeys.ValueInt64Pointer()),
		Prefix:       data.Prefix.ValueStringPointer(),
		FetchOwner:   data.FetchOwner.ValueBool(),
		StartAfter:   data.StartAfter.ValueStringPointer(),
	}

	var maxKeys, nKeys int64
	if data.MaxKeys.IsNull() {
		maxKeys = 1000
	} else {
		maxKeys = data.MaxKeys.ValueInt64()
	}

	keys := make([]types.String, 0)
	owners := make([]types.String, 0)
	pages := NewListObjectsV2Paginator(c.client, input)
pageLoop:
	for pages.HasMorePages() {
		page, err := pages.NextPage(ctx)
		if err != nil {
			return fmt.Errorf("error fetching page: %w", err)
		}

		if page.CommonPrefixes != nil {
			data.CommonPrefixes = make([]types.String, len(*page.CommonPrefixes))
			for i, prefix := range *page.CommonPrefixes {
				data.CommonPrefixes[i] = types.StringPointerValue(prefix.Prefix)
			}
		}

		if page.Contents != nil {
			for _, v := range *page.Contents {
				if nKeys >= maxKeys {
					// The break statement with label is used to break out of the keys loop and the page loop when
					//the number of keys fetched is equal to or greater than the max keys specified.
					break pageLoop
				}

				keys = append(keys, types.StringPointerValue(v.Key))
				if v.Owner != nil {
					owners = append(owners, types.StringPointerValue(v.Owner.DisplayName))
				}

				nKeys++
			}
		}
	}

	data.Keys = keys
	data.Owners = owners
	return nil
}
