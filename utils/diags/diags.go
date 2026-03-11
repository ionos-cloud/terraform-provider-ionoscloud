package diags

import (
	"fmt"
	"net/url"
	"regexp"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var requestIDPattern = regexp.MustCompile(`/requests/([-A-Fa-f0-9]+)/`)

// DiagsOpts contains optional fields for enhancing error messages.
type DiagsOpts struct {
	Timeout            string
	ResourceNameString string // The schema key for the resource name (e.g., "name", "display_name")
	RequestLocation    *url.URL
	StatusCode         int // HTTP status code from API response; 0 means not available
}

// ErrorContext holds plain-string context for error enrichment.
// Used directly by framework code, and populated from *schema.ResourceData by SDKv2 code.
type ErrorContext struct {
	ResourceID      string
	ResourceName    string
	Timeout         string
	StatusCode      int    // 0 = not available
	RequestLocation string // raw Location header/URL string for request ID extraction
}

// buildContextString constructs context info from plain values (no SDKv2 dependency).
func buildContextString(ctx *ErrorContext) string {
	if ctx == nil {
		return ""
	}

	var sb strings.Builder

	if ctx.StatusCode >= 500 {
		sb.WriteString("\nThis is an API Error. Please contact API support.")
	}

	if ctx.ResourceID != "" {
		fmt.Fprintf(&sb, "\nResource ID: %s", ctx.ResourceID)
	}

	if ctx.ResourceName != "" {
		fmt.Fprintf(&sb, "\nResource Name: %s", ctx.ResourceName)
	}

	if ctx.Timeout != "" {
		fmt.Fprintf(&sb, "\nConfigured timeout: %s", ctx.Timeout)
	}

	if ctx.RequestLocation != "" {
		reqID := extractRequestId(ctx.RequestLocation)
		if reqID != "" {
			fmt.Fprintf(&sb, "\nRequest ID: %s", reqID)
		}
	}

	return sb.String()
}

// WrapError wraps an error with context from an ErrorContext.
// Used by framework code: resp.Diagnostics.AddError("title", diags.WrapError(err, ctx).Error())
func WrapError(err error, ctx *ErrorContext) error {
	if err == nil {
		return nil
	}
	context := buildContextString(ctx)
	if context != "" {
		return fmt.Errorf("%w%s", err, context)
	}
	return err
}

// buildContext constructs context information to append after an error message.
func buildContext(d *schema.ResourceData, opts *DiagsOpts) string {
	targetField := "name"
	if opts != nil && opts.ResourceNameString != "" {
		targetField = opts.ResourceNameString
	}

	ctx := &ErrorContext{}
	ctx.ResourceID = d.Id()

	if v, ok := d.GetOk(targetField); ok {
		if nameVal, isString := v.(string); isString {
			ctx.ResourceName = nameVal
		}
	}

	if opts != nil {
		ctx.Timeout = opts.Timeout
		ctx.StatusCode = opts.StatusCode
		if opts.RequestLocation != nil {
			ctx.RequestLocation = opts.RequestLocation.String()
		}
	}

	return buildContextString(ctx)
}

// ToDiags wraps an error into a Terraform diagnostic with additional context.
func ToDiags(d *schema.ResourceData, err error, opts *DiagsOpts) diag.Diagnostics {
	if err == nil {
		return nil
	}
	return diag.FromErr(ToError(d, err, opts))
}

// ToError wraps an error with resource context, similar to ToDiags but returns an error.
// This is used in import functions and helper functions that return (T, error) instead of diag.Diagnostics.
func ToError(d *schema.ResourceData, err error, opts *DiagsOpts) error {
	if err == nil {
		return nil
	}
	context := buildContext(d, opts)
	if context != "" {
		return fmt.Errorf("%w%s", err, context)
	}
	return err
}

// extractRequestId parses the HTTP header string to find the UUID request ID.
func extractRequestId(loc string) string {
	matches := requestIDPattern.FindStringSubmatch(loc)
	if len(matches) > 1 {
		return matches[1]
	}
	return ""
}
