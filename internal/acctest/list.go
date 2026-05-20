package acctest

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/list"
	fwresource "github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

// TestAccCheckResourceInList returns a TestCheckFunc that verifies a resource appears in the list stream.
//
// If lookupInState is true, expectedVal is a resource address (e.g. "ionoscloud_s3_bucket.test")
// and the ID is retrieved from state. Otherwise expectedVal is used directly as the identity ID.
func TestAccCheckResourceInList(
	listRes list.ListResource,
	wrappedRes fwresource.Resource,
	expectedVal string,
	lookupInState bool,
) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		ctx := context.Background()
		client := NewTestBundleClientFromEnv()

		var expectedID string
		if lookupInState {
			rs, ok := s.RootModule().Resources[expectedVal]
			if !ok {
				return fmt.Errorf("resource not found in state: %s", expectedVal)
			}
			expectedID = rs.Primary.ID
		} else {
			expectedID = expectedVal
		}

		var resp fwresource.ConfigureResponse
		listRes.(fwresource.ResourceWithConfigure).Configure(ctx, fwresource.ConfigureRequest{
			ProviderData: client,
		}, &resp)
		if resp.Diagnostics.HasError() {
			return fmt.Errorf("failed to configure list resource: %v", resp.Diagnostics)
		}

		var schemaResp fwresource.SchemaResponse
		wrappedRes.Schema(ctx, fwresource.SchemaRequest{}, &schemaResp)

		var identitySchemaResp fwresource.IdentitySchemaResponse
		wrappedRes.(fwresource.ResourceWithIdentity).IdentitySchema(ctx, fwresource.IdentitySchemaRequest{}, &identitySchemaResp)

		req := list.ListRequest{
			IncludeResource:        true,
			ResourceSchema:         schemaResp.Schema,
			ResourceIdentitySchema: identitySchemaResp.IdentitySchema,
		}

		var stream list.ListResultsStream
		listRes.List(ctx, req, &stream)
		if stream.Results == nil {
			return fmt.Errorf("results stream is nil")
		}

		found := false
		var errs []error
		stream.Results(func(result list.ListResult) bool {
			if result.Diagnostics.HasError() {
				errs = append(errs, fmt.Errorf("diagnostic error in stream: %v", result.Diagnostics))
				return true
			}

			var idModel struct {
				ID string `tfsdk:"id"`
			}
			if diags := result.Identity.Get(ctx, &idModel); diags.HasError() {
				errs = append(errs, fmt.Errorf("failed to get identity: %v", diags))
				return true
			}
			if idModel.ID == expectedID {
				found = true
				return false
			}
			return true
		})

		if len(errs) > 0 {
			return fmt.Errorf("errors during list iteration: %v", errs)
		}
		if !found {
			return fmt.Errorf("expected resource identifier %q was not found in the list", expectedID)
		}
		return nil
	}
}
