package identity

import (
	"context"
	"errors"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/list"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
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

// MappedItemFromResourceData builds a MappedItem out of an SDKv2 ResourceData that
// the managed resource's own state writer has already filled.
//
// This is what lets a list resource for an SDKv2 managed resource avoid declaring a
// model of its own. TfTypeIdentityState and TfTypeResourceState render the
// ResourceData through the same schema the framework built the result type from, so
// the two cannot disagree, and tfsdk.Resource.Set / tfsdk.ResourceIdentity.Set take
// the resulting tftypes.Value directly instead of reflecting over a Go struct. A
// schema change therefore only has to be followed in the resource's own state
// writer, which lives beside the schema.
//
// The values are passed by value on purpose: Set type-asserts a tftypes.Value, and a
// *tftypes.Value would fall through to struct reflection and fail.
//
// The caller must have run a writer that calls d.SetId - both conversions go through
// ResourceData.State(), which returns nil while the ID is empty.
func MappedItemFromResourceData(displayName string, rd *schema.ResourceData, includeResource bool) (*MappedItem, error) {
	identityValue, err := rd.TfTypeIdentityState()
	if err != nil {
		return nil, fmt.Errorf("failed to convert the resource identity: %w", err)
	}
	if identityValue == nil {
		return nil, errors.New("the resource identity is empty; the identity writer has to run before mapping")
	}

	mapped := &MappedItem{DisplayName: displayName, Identity: *identityValue}
	if !includeResource {
		return mapped, nil
	}

	resourceValue, err := rd.TfTypeResourceState()
	if err != nil {
		return nil, fmt.Errorf("failed to convert the resource state: %w", err)
	}
	if resourceValue == nil {
		return nil, errors.New("the resource state is empty")
	}

	withoutTimeouts, err := nullTimeouts(*resourceValue)
	if err != nil {
		return nil, fmt.Errorf("failed to null the timeouts block: %w", err)
	}
	mapped.Resource = withoutTimeouts

	return mapped, nil
}

// nullTimeouts nulls the timeouts block of a converted resource state.
//
// The flatmap shims cannot tell a null single block from an empty one, so they always
// materialise timeouts as an object of null attributes. SDKv2 nulls it back out on
// every read of its own resources for that exact reason (helper/schema/grpc_provider.go,
// "we can't determine if a single block was null from the flatmapped values"). A listed
// resource has no timeouts either - they are configuration, not state - so doing the
// same here keeps a query result identical to what a refresh writes.
func nullTimeouts(value tftypes.Value) (tftypes.Value, error) {
	if value.IsNull() {
		return value, nil
	}

	valueType := value.Type()
	object, ok := valueType.(tftypes.Object)
	if !ok {
		return value, nil
	}
	timeoutsType, ok := object.AttributeTypes[schema.TimeoutsConfigKey]
	if !ok {
		return value, nil
	}

	var attributes map[string]tftypes.Value
	if err := value.As(&attributes); err != nil {
		return tftypes.Value{}, err
	}
	attributes[schema.TimeoutsConfigKey] = tftypes.NewValue(timeoutsType, nil)

	return tftypes.NewValue(valueType, attributes), nil
}
