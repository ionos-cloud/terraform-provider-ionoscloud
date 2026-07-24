package diags

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/url"
	"regexp"
	"sort"
	"strings"
	"sync"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var requestIDPattern = regexp.MustCompile(`/requests/([-A-Fa-f0-9]+)/`)

// Enricher carries per-provider-configuration context (currently the contract
// number) used to enrich error messages. Exactly one instance is created per
// provider Configure and stored on the SdkBundle, so concurrent provider
// configurations — multiple Crossplane/Upjet ProviderConfigs or provider
// aliases sharing a single process — no longer race on, or overwrite, shared
// package-level state. It replaces the package-level contractNumberFunc.
type Enricher struct {
	contractNumberOnce sync.Once
	contractNumber     string
	resolve            func() string // lazy API fallback; nil when the value is known eagerly
}

// NewEnricher builds an Enricher from the available credentials, mirroring the
// resolution order of the legacy package-level resolver: an explicit contract
// number wins, then the contract number embedded in the JWT token, otherwise a
// lazy API fallback resolved once on first use.
func NewEnricher(contractNumber, token string, apiFallback func() string) *Enricher {
	switch {
	case contractNumber != "":
		return &Enricher{contractNumber: contractNumber}
	case token != "":
		return &Enricher{contractNumber: ContractNumberFromToken(token)}
	default:
		return &Enricher{resolve: apiFallback}
	}
}

// ContractNumber returns the resolved contract number, running the lazy API
// fallback at most once. Safe to call on a nil Enricher, which returns "".
func (e *Enricher) ContractNumber() string {
	if e == nil {
		return ""
	}
	if e.resolve != nil {
		e.contractNumberOnce.Do(func() { e.contractNumber = e.resolve() })
	}
	return e.contractNumber
}

// WrapError wraps an error with context from an ErrorContext, including this
// configuration's contract number. Safe to call on a nil Enricher.
func (e *Enricher) WrapError(err error, errCtx *ErrorContext) error {
	if err == nil {
		return nil
	}
	contextStr := buildContextString(errCtx, e.ContractNumber())
	if contextStr != "" {
		return fmt.Errorf("%w%s", err, contextStr)
	}
	return err
}

// ToError wraps an error with resource context. When errCtx is nil, resource ID
// and name are auto-filled from d. Safe to call on a nil Enricher.
func (e *Enricher) ToError(d *schema.ResourceData, err error, errCtx *ErrorContext) error {
	if err == nil {
		return nil
	}
	return e.WrapError(err, applyResourceData(d, errCtx))
}

// ToDiags wraps an error into a Terraform diagnostic with additional context.
// When errCtx is nil, resource ID and name are auto-filled from d. Safe to call
// on a nil Enricher.
func (e *Enricher) ToDiags(d *schema.ResourceData, err error, errCtx *ErrorContext) diag.Diagnostics {
	if err == nil {
		return nil
	}
	return diag.FromErr(e.ToError(d, err, errCtx))
}

// ErrorContext holds context for error enrichment.
// Used by both SDKv2 (via ToDiags/ToError) and framework code (via WrapError).
type ErrorContext struct {
	ResourceID     string
	ResourceName   string
	Timeout        string
	StatusCode     int               // HTTP status code from API response; 0 means not available
	RequestID      string            // pre-extracted UUID from request location
	AdditionalInfo map[string]string // additional context, e.g. "Cluster ID" -> "abc"
}

// ExtractRequestID extracts the request UUID from a *url.URL location header.
func ExtractRequestID(loc *url.URL) string {
	if loc == nil {
		return ""
	}
	matches := requestIDPattern.FindStringSubmatch(loc.String())
	if len(matches) > 1 {
		return matches[1]
	}
	return ""
}

// ContractNumberFromToken extracts the contract number from a JWT token payload.
func ContractNumberFromToken(token string) string {
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return ""
	}
	payload, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return ""
	}
	var claims struct {
		Identity struct {
			ContractNumber json.Number `json:"contractNumber"`
		} `json:"identity"`
	}
	if err := json.Unmarshal(payload, &claims); err != nil {
		return ""
	}
	return claims.Identity.ContractNumber.String()
}

// buildContextString constructs context info from an ErrorContext and the
// resolved contract number for the current configuration.
func buildContextString(errCtx *ErrorContext, contractNumber string) string {
	if errCtx == nil {
		return ""
	}

	var sb strings.Builder

	if errCtx.StatusCode >= 500 {
		sb.WriteString("\nThis is an API Error. Please contact API support.")
	}

	if errCtx.ResourceID != "" {
		fmt.Fprintf(&sb, "\nResource ID: %s", errCtx.ResourceID)
	}

	if errCtx.ResourceName != "" {
		fmt.Fprintf(&sb, "\nResource Name: %s", errCtx.ResourceName)
	}

	if errCtx.Timeout != "" {
		fmt.Fprintf(&sb, "\nConfigured timeout: %s", errCtx.Timeout)
	}

	if errCtx.RequestID != "" {
		fmt.Fprintf(&sb, "\nRequest ID: %s", errCtx.RequestID)
	}

	if contractNumber != "" {
		fmt.Fprintf(&sb, "\nContract number: %s", contractNumber)
	}

	if len(errCtx.AdditionalInfo) > 0 {
		keys := make([]string, 0, len(errCtx.AdditionalInfo))
		for k := range errCtx.AdditionalInfo {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		for _, k := range keys {
			fmt.Fprintf(&sb, "\n%s: %s", k, errCtx.AdditionalInfo[k])
		}
	}

	return sb.String()
}

// applyResourceData fills resource ID and name from d when not already set on
// errCtx, allocating an ErrorContext if one was not provided.
func applyResourceData(d *schema.ResourceData, errCtx *ErrorContext) *ErrorContext {
	if errCtx == nil {
		errCtx = &ErrorContext{}
	}

	if errCtx.ResourceID == "" {
		errCtx.ResourceID = d.Id()
	}

	if errCtx.ResourceName == "" {
		if v, ok := d.GetOk("name"); ok {
			if nameVal, isString := v.(string); isString {
				errCtx.ResourceName = nameVal
			}
		} else if v, ok := d.GetOk("display_name"); ok {
			if nameVal, isString := v.(string); isString {
				errCtx.ResourceName = nameVal
			}
		}
	}

	return errCtx
}

// Note: error wrapping is now done exclusively through the per-configuration
// (*Enricher) methods (WrapError/ToError/ToDiags), reachable from the SdkBundle
// — for SDKv2 via bundleclient.ToDiags/ToError(meta, ...) and for framework
// code via the configured client's Enricher. There are intentionally no
// package-level wrappers, so the contract number can never be read from shared
// process-global state.
