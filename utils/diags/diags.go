package diags

import (
	"fmt"
	"net/url"
"regexp"
	"sort"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var requestIDPattern = regexp.MustCompile(`/requests/([-A-Fa-f0-9]+)/`)

// ErrorContext holds context for error enrichment.
// Used by both SDKv2 (via ToDiags/ToError) and framework code (via WrapError).
type ErrorContext struct {
	ResourceID   string
	ResourceName string
	Timeout      string
	StatusCode   int               // HTTP status code from API response; 0 means not available
	RequestID    string            // pre-extracted UUID from request location
	AdditionalInfo      map[string]string // additional context, e.g. "Cluster ID" -> "abc"
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
