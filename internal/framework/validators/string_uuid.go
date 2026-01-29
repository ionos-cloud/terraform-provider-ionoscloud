package validators

import (
	"context"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

var _ validator.String = UUIDValidator{}

type UUIDValidator struct {
}

// Description returns a plain text description of the validator's behavior.
func (v UUIDValidator) Description(_ context.Context) string {
	return "string must be a valid UUID"
}

// MarkdownDescription returns a markdown formatted description.
func (v UUIDValidator) MarkdownDescription(_ context.Context) string {
	return "string must be a valid UUID"
}

// ValidateString checks if the string it's a valid UUID.
func (v UUIDValidator) ValidateString(_ context.Context, req validator.StringRequest, resp *validator.StringResponse) {
	if req.ConfigValue.IsUnknown() || req.ConfigValue.IsNull() {
		return
	}

	err := uuid.Validate(req.ConfigValue.ValueString())
	if err != nil {
		resp.Diagnostics.AddAttributeError(req.Path, "Invalid String Format", "String must be a valid UUID")
	}
}
