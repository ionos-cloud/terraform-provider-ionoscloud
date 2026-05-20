package objectstoragemanagement

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/list"
	listschema "github.com/hashicorp/terraform-plugin-framework/list/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	objectstoragemanagementSDK "github.com/ionos-cloud/sdk-go-bundle/products/objectstoragemanagement/v2"

	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/internal/framework/identity"
	objectStorageManagementService "github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/objectstoragemanagement"
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
		Attributes: map[string]listschema.Attribute{},
	}
}

// List lists all accesskey resources.
func (r *accesskeyResource) List(ctx context.Context, req list.ListRequest, stream *list.ListResultsStream) {
	if r.client == nil {
		identity.StreamError(stream, "object storage management api client not configured", "The provider client is not configured")
		return
	}

	accessKeys, _, err := r.client.ListAccessKeys(ctx)
	if err != nil {
		identity.StreamError(stream, "failed to list access keys", err.Error())
		return
	}

	identity.StreamResults(ctx, stream, req, accessKeys.Items, r.Map)
}

// Map populates result and returns whether to push it (always true for access keys).
func (r *accesskeyResource) Map(ctx context.Context, req list.ListRequest, ak objectstoragemanagementSDK.AccessKeyRead, result *list.ListResult) (bool, diag.Diagnostics) {
	var diags diag.Diagnostics

	result.DisplayName = ak.Properties.AccessKey
	diags.Append(result.Identity.Set(ctx, &identity.Model{
		ID: types.StringValue(ak.Id),
	})...)

	if req.IncludeResource {
		timeoutsAttrTypes := map[string]attr.Type{
			"create": types.StringType,
			"read":   types.StringType,
			"delete": types.StringType,
		}

		resourceModel := objectStorageManagementService.AccesskeyResourceModel{
			AccessKey:       types.StringValue(ak.Properties.AccessKey),
			SecretKey:       types.StringNull(),
			CanonicalUserID: types.StringPointerValue(ak.Properties.CanonicalUserId),
			ContractUserID:  types.StringPointerValue(ak.Properties.ContractUserId),
			Description:     types.StringValue(ak.Properties.Description),
			ID:              types.StringValue(ak.Id),
			Timeouts:        timeouts.Value{Object: types.ObjectNull(timeoutsAttrTypes)},
		}
		diags.Append(result.Resource.Set(ctx, &resourceModel)...)
	}

	return true, diags
}
