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

// buildContext constructs context information to append after an error message.
func buildContext(d *schema.ResourceData, opts *DiagsOpts) string {
	targetField := "name"
	if opts != nil && opts.ResourceNameString != "" {
		targetField = opts.ResourceNameString
	}

	var sb strings.Builder

	if opts != nil && opts.StatusCode >= 500 {
		sb.WriteString("\nThis is an API Error. Please contact API support.")
	}

	if d.Id() != "" {
		fmt.Fprintf(&sb, "\nResource ID: %s", d.Id())
	}

	if v, ok := d.GetOk(targetField); ok {
		if nameVal, isString := v.(string); isString {
			fmt.Fprintf(&sb, "\nResource Name: %s", nameVal)
		}
	}

	if opts != nil {
		if opts.Timeout != "" {
			fmt.Fprintf(&sb, "\nConfigured timeout: %s", opts.Timeout)
		}
		if opts.RequestLocation != nil {
			reqID := extractRequestId(opts.RequestLocation.String())
			if reqID != "" {
				fmt.Fprintf(&sb, "\nRequest ID: %s", reqID)
			}
		}
	}

	return sb.String()
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
