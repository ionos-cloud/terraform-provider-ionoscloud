package identity

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/identityschema"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// ResourceWithIdentityWrapper wraps any standard resource and automatically
// exposes either a default "id"-based Identity Schema or a specified custom schema.
type ResourceWithIdentityWrapper struct {
	resource.Resource
	CustomSchema *identityschema.Schema
}

// IdentitySchema implements resource.ResourceWithIdentity.
func (w ResourceWithIdentityWrapper) IdentitySchema(ctx context.Context, req resource.IdentitySchemaRequest, resp *resource.IdentitySchemaResponse) {
	if w.CustomSchema != nil {
		resp.IdentitySchema = *w.CustomSchema
	} else {
		resp.IdentitySchema = identityschema.Schema{
			Attributes: map[string]identityschema.Attribute{
				"id": identityschema.StringAttribute{
					RequiredForImport: true,
				},
			},
		}
	}
}

// Create implements resource.Resource.
func (w ResourceWithIdentityWrapper) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	w.Resource.Create(ctx, req, resp)
	if !resp.Diagnostics.HasError() {
		w.populateIdentity(ctx, resp.State, resp.Identity)
	}
}

// Read implements resource.Resource.
func (w ResourceWithIdentityWrapper) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	w.Resource.Read(ctx, req, resp)
	if !resp.Diagnostics.HasError() {
		w.populateIdentity(ctx, resp.State, resp.Identity)
	}
}

// Update implements resource.Resource.
func (w ResourceWithIdentityWrapper) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	w.Resource.Update(ctx, req, resp)
	if !resp.Diagnostics.HasError() {
		w.populateIdentity(ctx, resp.State, resp.Identity)
	}
}

// Delete implements resource.Resource.
func (w ResourceWithIdentityWrapper) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	w.Resource.Delete(ctx, req, resp)
}

func (w ResourceWithIdentityWrapper) populateIdentity(ctx context.Context, state tfsdk.State, identity *tfsdk.ResourceIdentity) {
	if identity == nil {
		return
	}

	if w.CustomSchema == nil {
		var idVal types.String
		diags := state.GetAttribute(ctx, path.Root("id"), &idVal)
		if diags.HasError() {
			return
		}
		if !idVal.IsUnknown() && !idVal.IsNull() {
			identity.Set(ctx, &Model{
				ID: idVal,
			})
		}
	}
}

// Configure delegates to the wrapped resource if it implements ResourceWithConfigure.
func (w ResourceWithIdentityWrapper) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if c, ok := w.Resource.(resource.ResourceWithConfigure); ok {
		c.Configure(ctx, req, resp)
	}
}

// ImportState delegates to the wrapped resource if it implements ResourceWithImportState.
func (w ResourceWithIdentityWrapper) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	if i, ok := w.Resource.(resource.ResourceWithImportState); ok {
		i.ImportState(ctx, req, resp)
	}
}

// ModifyPlan delegates to the wrapped resource if it implements ResourceWithModifyPlan.
func (w ResourceWithIdentityWrapper) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	if m, ok := w.Resource.(resource.ResourceWithModifyPlan); ok {
		m.ModifyPlan(ctx, req, resp)
	}
}

// UpgradeState delegates to the wrapped resource if it implements ResourceWithUpgradeState.
func (w ResourceWithIdentityWrapper) UpgradeState(ctx context.Context, req resource.UpgradeStateRequest, resp *resource.UpgradeStateResponse) {
	if u, ok := w.Resource.(resource.ResourceWithUpgradeState); ok {
		u.UpgradeState(ctx)
	}
}

// ValidateConfig delegates to the wrapped resource if it implements ResourceWithValidateConfig.
func (w ResourceWithIdentityWrapper) ValidateConfig(ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse) {
	if v, ok := w.Resource.(resource.ResourceWithValidateConfig); ok {
		v.ValidateConfig(ctx, req, resp)
	}
}

// ConfigValidators delegates to the wrapped resource if it implements ResourceWithConfigValidators.
func (w ResourceWithIdentityWrapper) ConfigValidators(ctx context.Context) []resource.ConfigValidator {
	if v, ok := w.Resource.(resource.ResourceWithConfigValidators); ok {
		return v.ConfigValidators(ctx)
	}
	return nil
}

// WithIdentity is a helper function to wrap a resource constructor with the default identity wrapper.
func WithIdentity(fn func() resource.Resource) func() resource.Resource {
	return func() resource.Resource {
		return ResourceWithIdentityWrapper{Resource: fn()}
	}
}

// WithCustomIdentity is a helper function to wrap a resource constructor with a custom identity schema (e.g. for composite keys).
func WithCustomIdentity(fn func() resource.Resource, schema identityschema.Schema) func() resource.Resource {
	return func() resource.Resource {
		return ResourceWithIdentityWrapper{
			Resource:     fn(),
			CustomSchema: &schema,
		}
	}
}
