package s3

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-validators/objectvalidator"

	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/s3"

	"github.com/hashicorp/terraform-plugin-framework-validators/resourcevalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

var (
	_ resource.ResourceWithImportState      = (*bucketWebsiteConfiguration)(nil)
	_ resource.ResourceWithConfigure        = (*bucketWebsiteConfiguration)(nil)
	_ resource.ResourceWithConfigValidators = (*bucketWebsiteConfiguration)(nil)
)

type bucketWebsiteConfiguration struct {
	client *s3.Client
}

func (r *bucketWebsiteConfiguration) ConfigValidators(ctx context.Context) []resource.ConfigValidator {
	return []resource.ConfigValidator{
		resourcevalidator.Conflicting(
			path.MatchRoot("redirect_all_requests_to"),
			path.MatchRoot("index_document"),
		),
		resourcevalidator.Conflicting(
			path.MatchRoot("redirect_all_requests_to"),
			path.MatchRoot("error_document"),
		),
		resourcevalidator.Conflicting(
			path.MatchRoot("redirect_all_requests_to"),
			path.MatchRoot("routing_rule"),
		),
		resourcevalidator.ExactlyOneOf(path.MatchRoot("index_document"), path.MatchRoot("redirect_all_requests_to")),
	}
}

// NewBucketWebsiteConfigurationResource creates a new resource for the bucket website configuration resource.
func NewBucketWebsiteConfigurationResource() resource.Resource {
	return &bucketWebsiteConfiguration{}
}

// Metadata returns the metadata for the bucket website configuration resource.
func (r *bucketWebsiteConfiguration) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_s3_bucket_website_configuration"
}

// Schema returns the schema for the bucket website configuration resource.
func (r *bucketWebsiteConfiguration) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"bucket": schema.StringAttribute{
				Description: "The name of the bucket.",
				Required:    true,
			},
		},
		Blocks: map[string]schema.Block{
			"index_document": schema.SingleNestedBlock{
				Description: "Container for the Suffix element.",
				Attributes: map[string]schema.Attribute{
					"suffix": schema.StringAttribute{
						Description: "A suffix that is appended to a request that is for a directory on the website endpoint (for example, if the suffix is index.html and you make a request to samplebucket/images/ the data that is returned will be for the object with the key name images/index.html) The suffix must not be empty and must not include a slash character. Replacement must be made for object keys containing special characters (such as carriage returns) when using XML requests.",
						Optional:    true,
						Validators:  []validator.String{stringvalidator.LengthAtLeast(1), stringvalidator.NoneOf("/")},
					},
				},
				Validators: []validator.Object{
					objectvalidator.AlsoRequires(path.Expressions{path.MatchRelative().AtName("suffix")}...),
				},
			},
			"error_document": schema.SingleNestedBlock{
				Description: "The object key name to use when a 4XX class error occurs. Replacement must be made for object keys containing special characters (such as carriage returns) when using XML requests.",
				Attributes: map[string]schema.Attribute{
					"key": schema.StringAttribute{
						Description: "The object key.",
						Optional:    true,
						Validators:  []validator.String{stringvalidator.LengthBetween(1, 1024)},
					},
				},

				Validators: []validator.Object{
					objectvalidator.AlsoRequires(path.Expressions{path.MatchRelative().AtName("key")}...),
				},
			},
			"redirect_all_requests_to": schema.SingleNestedBlock{
				Description: "Container for redirect information. You can redirect requests to another host, to another page, or with another protocol. In the event of an error, you can specify a different error code to return.",
				Attributes: map[string]schema.Attribute{
					"host_name": schema.StringAttribute{
						Description: "The host name to use in the redirect request.",
						Optional:    true,
					},
					"protocol": schema.StringAttribute{
						Description: "Protocol to use when redirecting requests. The default is the protocol that is used in the original request.",
						Optional:    true,
						Validators:  []validator.String{stringvalidator.OneOf("http", "https")},
					},
				},
			},
			"routing_rule": schema.ListNestedBlock{
				Description: "Rules that define when a redirect is applied and the redirect behavior.",
				NestedObject: schema.NestedBlockObject{
					Blocks: map[string]schema.Block{
						"condition": schema.SingleNestedBlock{
							Description: "A container for describing a condition that must be met for the specified redirect to apply. For example, 1. If request is for pages in the /docs folder, redirect to the /documents folder. 2. If request results in HTTP error 4xx, redirect request to another host where you might process the error.",
							Attributes: map[string]schema.Attribute{
								"http_error_code_returned_equals": schema.StringAttribute{
									Description: "The HTTP error code when the redirect is applied. In the event of an error, if the error code equals this value, then the specified redirect is applied. Required when parent element Condition is specified and sibling KeyPrefixEquals is not specified. If both are specified, then both must be true for the redirect to be applied",
									Optional:    true,
								},
								"key_prefix_equals": schema.StringAttribute{
									Description: "The object key name prefix when the redirect is applied. For example, to redirect requests for `ExamplePage.html`, the key prefix will be `ExamplePage.html`. To redirect request for all pages with the prefix `docs/`, the key prefix will be `/docs`, which identifies all objects in the `docs/` folder. Required when the parent element `Condition` is specified and sibling `HTTPErrorCodeReturnedEquals` is not specified. If both conditions are specified, both must be true for the redirect to be applied. Replacement must be made for object keys containing special characters (such as carriage returns) when using XML requests.",
									Optional:    true,
								},
							},
						},
						"redirect": schema.SingleNestedBlock{
							Description: "Container for redirect information. You can redirect requests to another host, to another page, or with another protocol. In the event of an error, you can specify a different error code to return.",
							Attributes: map[string]schema.Attribute{
								"host_name": schema.StringAttribute{
									Description: "The host name to use in the redirect request.",
									Optional:    true,
								},
								"http_redirect_code": schema.StringAttribute{
									Description: "The HTTP redirect code to use on the response.",
									Optional:    true,
								},
								"protocol": schema.StringAttribute{
									Description: "The protocol to use in the redirect request.",
									Optional:    true,
									Validators:  []validator.String{stringvalidator.OneOf("http", "https")},
								},
								"replace_key_prefix_with": schema.StringAttribute{
									Description: "The object key prefix to use in the redirect request. For example, to redirect requests for all pages with prefix `docs/` (objects in the `docs/` folder) to `documents/`, you can set a condition block with `KeyPrefixEquals` set to `docs/` and in the Redirect set `ReplaceKeyPrefixWith` to `/documents`. Not required if one of the siblings is present. Can be present only if `ReplaceKeyWith` is not provided.",
									Optional:    true,
								},
								"replace_key_with": schema.StringAttribute{
									Description: "The specific object key to use in the redirect request. For example, redirect request to error.html. Not required if one of the siblings is present. Can be present only if ReplaceKeyPrefixWith is not provided. Replacement must be made for object keys containing special characters (such as carriage returns) when using XML requests.",
									Optional:    true,
								},
							},
						},
					},
				},
			},
		},
	}
}

// Configure configures the bucket website configuration resource.
func (r *bucketWebsiteConfiguration) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*s3.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *s3.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.client = client
}

// Create creates the bucket website configuration.
func (r *bucketWebsiteConfiguration) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	if r.client == nil {
		resp.Diagnostics.AddError("API client not configured", "The provider client is not configured")
		return
	}

	var data *s3.BucketWebsiteConfigurationModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.client.CreateBucketWebsite(ctx, data); err != nil {
		resp.Diagnostics.AddError("Failed to create resource", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Read reads the bucket website configuration.
func (r *bucketWebsiteConfiguration) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	if r.client == nil {
		resp.Diagnostics.AddError("API client not configured", "The provider client is not configured")
		return
	}

	var data *s3.BucketWebsiteConfigurationModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	result, found, err := r.client.GetBucketWebsite(ctx, data.Bucket)
	if err != nil {
		resp.Diagnostics.AddError("Failed to read resource", err.Error())
		return
	}

	if !found {
		resp.State.RemoveResource(ctx)
		return
	}

	data = result
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// ImportState imports the state of the bucket website configuration.
func (r *bucketWebsiteConfiguration) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("bucket"), req, resp)
}

// Update updates the bucket website configuration.
func (r *bucketWebsiteConfiguration) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	if r.client == nil {
		resp.Diagnostics.AddError("API client not configured", "The provider client is not configured")
		return
	}

	var data *s3.BucketWebsiteConfigurationModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.client.UpdateBucketWebsite(ctx, data); err != nil {
		resp.Diagnostics.AddError("Failed to update resource", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Delete deletes the bucket website configuration.
func (r *bucketWebsiteConfiguration) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	if r.client == nil {
		resp.Diagnostics.AddError("API client not configured", "The provider client is not configured")
		return
	}

	var data *s3.BucketWebsiteConfigurationModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.client.DeleteBucketWebsite(ctx, data.Bucket); err != nil {
		resp.Diagnostics.AddError("Failed to delete resource", err.Error())
		return
	}
}
