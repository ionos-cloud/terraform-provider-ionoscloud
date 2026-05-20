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
		Attributes: map[string]listschema.Attribute{},
	}
}

// List lists all bucket resources.
func (r *bucketResource) List(ctx context.Context, req list.ListRequest, stream *list.ListResultsStream) {
	if r.client == nil {
		identity.StreamError(stream, "object storage api client not configured", "The provider client is not configured")
		return
	}

	buckets, err := r.client.ListBuckets(ctx)
	if err != nil {
		identity.StreamError(stream, "failed to list buckets", err.Error())
		return
	}

	identity.StreamResults(ctx, stream, req, buckets, r.Map)
}

// Map implements identity.ListMapper for mapping an objstorage.Bucket to a list.ListResult.
func (r *bucketResource) Map(ctx context.Context, req list.ListRequest, b objstorage.Bucket, result *list.ListResult) (bool, diag.Diagnostics) {
	var diags diag.Diagnostics
	if b.Name == nil {
		return false, diags
	}
	result.DisplayName = *b.Name

	diags.Append(result.Identity.Set(ctx, &identity.Model{
		ID: types.StringValue(*b.Name),
	})...)

	if req.IncludeResource {
		region, err := r.client.GetBucketLocation(ctx, types.StringValue(*b.Name))
		if err != nil {
			diags.AddError("failed to get bucket location", err.Error())
			return true, diags
		}

		bucketModel, found, err := r.client.GetBucket(ctx, types.StringValue(*b.Name), region)
		if err != nil {
			diags.AddError("failed to get bucket details", err.Error())
			return true, diags
		}
		if !found {
			diags.AddError("bucket not found during detailed fetch", fmt.Sprintf("Bucket %s was not found", *b.Name))
			return true, diags
		}

		timeoutsAttrTypes := map[string]attr.Type{
			"create": types.StringType,
			"read":   types.StringType,
			"update": types.StringType,
			"delete": types.StringType,
		}

		resourceModel := bucketResourceModel{
			ID:                bucketModel.Name,
			Name:              bucketModel.Name,
			Region:            bucketModel.Region,
			ObjectLockEnabled: bucketModel.ObjectLockEnabled,
			ForceDestroy:      types.BoolNull(),
			Tags:              bucketModel.Tags,
			Timeouts:          timeouts.Value{Object: types.ObjectNull(timeoutsAttrTypes)},
		}

		diags.Append(result.Resource.Set(ctx, &resourceModel)...)
	}

	return true, diags
}
