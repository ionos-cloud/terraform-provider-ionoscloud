package bundleclient

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	diagutil "github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils/diags"
)

// enricherFromMeta returns the per-configuration error Enricher carried on the
// SdkBundle meta. It returns nil when meta is not an SdkBundle or has no
// Enricher set; that is safe because every (*diagutil.Enricher) method is
// nil-receiver tolerant and simply omits the contract number.
//
// This is the SDKv2 replacement for the old package-level contract-number
// resolver: the contract number now travels with the per-configuration bundle
// (set in the provider's Configure / upjet's TerraformSetupBuilder) instead of
// process-global state, so concurrent provider configurations cannot overwrite
// one another's contract number.
func enricherFromMeta(m any) *diagutil.Enricher {
	if b, ok := m.(SdkBundle); ok {
		return b.Diags
	}
	return nil
}

// ToDiags wraps an error into Terraform diagnostics, enriched with context from
// errCtx and the contract number scoped to this provider configuration (taken
// from the SdkBundle meta). When errCtx is nil, resource ID and name are
// auto-filled from d. It is the per-configuration replacement for
// diagutil.ToDiags.
func ToDiags(m any, d *schema.ResourceData, err error, errCtx *diagutil.ErrorContext) diag.Diagnostics {
	return enricherFromMeta(m).ToDiags(d, err, errCtx)
}

// ToError is the error-returning counterpart of ToDiags.
func ToError(m any, d *schema.ResourceData, err error, errCtx *diagutil.ErrorContext) error {
	return enricherFromMeta(m).ToError(d, err, errCtx)
}
