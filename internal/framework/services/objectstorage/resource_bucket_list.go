package objectstorage

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/list"
	listschema "github.com/hashicorp/terraform-plugin-framework/list/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	objstorage "github.com/ionos-cloud/sdk-go-bundle/products/objectstorage/v2"

	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/internal/framework/identity"
	diagutil "github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils/diags"
)

var (
	_ list.ListResource              = (*bucketResource)(nil)
	_ list.ListResourceWithConfigure = (*bucketResource)(nil)
)

// NewBucketListResource creates a new list resource for the bucket resource.
func NewBucketListResource() list.ListResource {
	return &bucketResource{}
}

// ListResourceConfigSchema returns the schema for the configuration of the bucket list resource.
func (r *bucketResource) ListResourceConfigSchema(_ context.Context, _ list.ListResourceSchemaRequest, resp *list.ListResourceSchemaResponse) {
	resp.Schema = listschema.Schema{
		Attributes: map[string]listschema.Attribute{
			"filters": identity.FilterAttribute(),
		},
	}
}

// List lists all bucket resources.
func (r *bucketResource) List(ctx context.Context, req list.ListRequest, stream *list.ListResultsStream) {
	identity.StreamList(ctx, stream, req, r.client.ListBuckets, r.Map)
}

// Map returns a MappedItem describing the bucket, or nil to skip it.
func (r *bucketResource) Map(ctx context.Context, includeResource bool, filters []identity.Filter, b objstorage.Bucket) (*identity.MappedItem, diag.Diagnostics) {
	var diags diag.Diagnostics
	if b.Name == nil {
		diags.AddError("unexpected API response", "bucket returned by API has no name")
		return nil, diags
	}

	region, err := r.client.GetBucketLocation(ctx, types.StringValue(*b.Name))
	if err != nil {
		diags.AddError("failed to get bucket location", diagutil.WrapError(err, &diagutil.ErrorContext{ResourceName: *b.Name}).Error())
		return nil, diags
	}

	if !identity.MatchesFilters(map[string]string{
		"name":   *b.Name,
		"region": region.ValueString(),
	}, filters) {
		return nil, diags
	}

	mapped := &identity.MappedItem{
		DisplayName: *b.Name,
		Identity:    &bucketIdentityModel{ID: types.StringValue(*b.Name), Region: region},
	}

	if !includeResource {
		return mapped, diags
	}

	bucketModel, found, err := r.client.GetBucket(ctx, types.StringValue(*b.Name), region)
	if err != nil {
		diags.AddError("failed to get bucket details", diagutil.WrapError(err, &diagutil.ErrorContext{ResourceName: *b.Name}).Error())
		return mapped, diags
	}
	if !found {
		diags.AddError("bucket not found during detailed fetch", fmt.Sprintf("Bucket %s was not found", *b.Name))
		return mapped, diags
	}

	timeoutsAttrTypes := map[string]attr.Type{
		"create": types.StringType,
		"read":   types.StringType,
		"update": types.StringType,
		"delete": types.StringType,
	}

	mapped.Resource = &bucketResourceModel{
		ID:                bucketModel.Name,
		Name:              bucketModel.Name,
		Region:            bucketModel.Region,
		ObjectLockEnabled: bucketModel.ObjectLockEnabled,
		ForceDestroy:      types.BoolValue(false),
		Tags:              bucketModel.Tags,
		Timeouts:          timeouts.Value{Object: types.ObjectNull(timeoutsAttrTypes)},
	}
	return mapped, diags
}
