package objectstoragemanagement

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/list"
	listschema "github.com/hashicorp/terraform-plugin-framework/list/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	objstoragesdk "github.com/ionos-cloud/sdk-go-bundle/products/objectstoragemanagement/v2"

	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/internal/framework/identity"
	objstorageservice "github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/objectstoragemanagement"
)

var (
	_ list.ListResource              = (*accesskeyResource)(nil)
	_ list.ListResourceWithConfigure = (*accesskeyResource)(nil)
)

// NewAccesskeyListResource creates a new list resource for the accesskey resource.
func NewAccesskeyListResource() list.ListResource {
	return &accesskeyResource{}
}

// ListResourceConfigSchema returns the schema for the configuration of the accesskey list resource.
func (r *accesskeyResource) ListResourceConfigSchema(_ context.Context, _ list.ListResourceSchemaRequest, resp *list.ListResourceSchemaResponse) {
	resp.Schema = listschema.Schema{
		Attributes: map[string]listschema.Attribute{
			identity.FiltersKey: identity.FilterAttribute(),
		},
	}
}

// List lists all accesskey resources.
func (r *accesskeyResource) List(ctx context.Context, req list.ListRequest, stream *list.ListResultsStream) {
	identity.StreamList(ctx, stream, req,
		func(ctx context.Context) ([]objstoragesdk.AccessKeyRead, error) {
			result, _, err := r.client.ListAccessKeys(ctx)
			return result.Items, err
		},
		r.Map,
	)
}

// Map returns a MappedItem describing the access key.
func (r *accesskeyResource) Map(_ context.Context, includeResource bool, filters []identity.Filter, ak objstoragesdk.AccessKeyRead) (*identity.MappedItem, diag.Diagnostics) {
	var diags diag.Diagnostics

	if !identity.MatchesFilters(map[string]string{
		"id":          ak.Id,
		"description": ak.Properties.Description,
		"accesskey":   ak.Properties.AccessKey,
	}, filters) {
		return nil, diags
	}

	mapped := &identity.MappedItem{
		DisplayName: ak.Properties.AccessKey,
		Identity:    &identity.Model{ID: types.StringValue(ak.Id)},
	}

	if !includeResource {
		return mapped, diags
	}

	timeoutsAttrTypes := map[string]attr.Type{
		"create": types.StringType,
		"read":   types.StringType,
		"delete": types.StringType,
	}

	mapped.Resource = &objstorageservice.AccesskeyResourceModel{
		AccessKey:       types.StringValue(ak.Properties.AccessKey),
		SecretKey:       types.StringNull(),
		CanonicalUserID: types.StringPointerValue(ak.Properties.CanonicalUserId),
		ContractUserID:  types.StringPointerValue(ak.Properties.ContractUserId),
		Description:     types.StringValue(ak.Properties.Description),
		ID:              types.StringValue(ak.Id),
		Timeouts:        timeouts.Value{Object: types.ObjectNull(timeoutsAttrTypes)},
	}
	return mapped, diags
}
