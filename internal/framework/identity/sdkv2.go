package identity

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/list"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-mux/tf5to6server/translate"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// SetRawV6Schemas fills resp with the protocol schemas of resourceSchema, the SDKv2
// managed resource being listed. Call it from a list resource's RawV6Schemas method.
//
// ProtoSchema and ProtoIdentitySchema return protocol v5, and this provider is
// served over v6, hence the translation.
//
// A resource that declares no Identity cannot be listed, since terraform identifies
// every result by its identity. RawV6Schemas cannot report diagnostics, so that case
// is logged here and the framework rejects the list resource with its own error.
func SetRawV6Schemas(ctx context.Context, resp *list.RawV6SchemaResponse, typeName string, resourceSchema *schema.Resource) {
	if resourceSchema == nil {
		tflog.Error(ctx, "no SDKv2 managed resource was set for this list resource", map[string]any{
			"resource_type": typeName,
		})
		return
	}

	identitySchema := resourceSchema.ProtoIdentitySchema(ctx)
	if identitySchema == nil {
		tflog.Error(ctx, "the SDKv2 managed resource declares no Identity, so it cannot be listed", map[string]any{
			"resource_type": typeName,
		})
		return
	}

	resp.ProtoV6Schema = translate.Schema(resourceSchema.ProtoSchema(ctx)())
	resp.ProtoV6IdentitySchema = translate.ResourceIdentitySchema(identitySchema())
}
