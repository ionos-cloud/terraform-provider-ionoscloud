package diags

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/url"
	"regexp"
	"sort"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var requestIDPattern = regexp.MustCompile(`/requests/([-A-Fa-f0-9]+)/`)

var contractNumberFunc func() string

// SetupContractNumberResolver configures the contract number resolver based on
// available credentials. If contractNumber is set explicitly, it is used directly.
// If a token is available, the contract number is extracted from the JWT.
// Otherwise, apiFallback is called lazily on the first error (and cached).
func SetupContractNumberResolver(contractNumber, token string, apiFallback func() string) {
	switch {
	case contractNumber != "":
		cn := contractNumber
		contractNumberFunc = func() string { return cn }
	case token != "":
		cn := ContractNumberFromToken(token)
		contractNumberFunc = func() string { return cn }
	default:
		called := false
		var cn string
		contractNumberFunc = func() string {
			if !called {
				called = true
				cn = apiFallback()
			}
			return cn
		}
	}
}

func getContractNumber() string {
	if contractNumberFunc != nil {
		return contractNumberFunc()
	}
	return ""
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

// buildContextString constructs context info from an ErrorContext.
func buildContextString(errCtx *ErrorContext) string {
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

	if cn := getContractNumber(); cn != "" {
		fmt.Fprintf(&sb, "\nContract number: %s", cn)
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

// WrapError wraps an error with context from an ErrorContext.
// Used by framework code: resp.Diagnostics.AddError("title", diags.WrapError(err, errCtx).Error())
func WrapError(err error, errCtx *ErrorContext) error {
	if err == nil {
		return nil
	}
	contextStr := buildContextString(errCtx)
	if contextStr != "" {
		return fmt.Errorf("%w%s", err, contextStr)
	}
	return err
}

// ToDiags wraps an error into a Terraform diagnostic with additional context.
// When errCtx is nil, resource ID and name are auto-filled from d.
func ToDiags(d *schema.ResourceData, err error, errCtx *ErrorContext) diag.Diagnostics {
	if err == nil {
		return nil
	}
	return diag.FromErr(ToError(d, err, errCtx))
}

// ToError wraps an error with resource context, similar to ToDiags but returns an error.
// When errCtx is nil, resource ID and name are auto-filled from d.
func ToError(d *schema.ResourceData, err error, errCtx *ErrorContext) error {
	if err == nil {
		return nil
	}

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

	return WrapError(err, errCtx)
}
