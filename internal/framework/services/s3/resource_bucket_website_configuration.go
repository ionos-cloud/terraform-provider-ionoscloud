package s3

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-validators/resourcevalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	s3 "github.com/ionos-cloud/sdk-go-s3"
)

var (
	_ resource.ResourceWithImportState      = (*bucketWebsiteConfiguration)(nil)
	_ resource.ResourceWithConfigure        = (*bucketWebsiteConfiguration)(nil)
	_ resource.ResourceWithConfigValidators = (*bucketWebsiteConfiguration)(nil)
)

type bucketWebsiteConfiguration struct {
	client *s3.APIClient
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
	}
}

type bucketWebsiteConfigurationModel struct {
	Bucket                types.String           `tfsdk:"bucket"`
	IndexDocument         *indexDocument         `tfsdk:"index_document"`
	ErrorDocument         *errorDocument         `tfsdk:"error_document"`
	RedirectAllRequestsTo *redirectAllRequestsTo `tfsdk:"redirect_all_requests_to"`
	RoutingRule           []routingRule          `tfsdk:"routing_rule"`
}

type indexDocument struct {
	Suffix types.String `tfsdk:"suffix"`
}

type errorDocument struct {
	Key types.String `tfsdk:"key"`
}

type redirectAllRequestsTo struct {
	HostName types.String `tfsdk:"host_name"`
	Protocol types.String `tfsdk:"protocol"`
}

type routingRule struct {
	Condition *condition `tfsdk:"condition"`
	Redirect  *redirect  `tfsdk:"redirect"`
}

type condition struct {
	HTTPErrorCodeReturnedEquals types.String `tfsdk:"http_error_code_returned_equals"`
	KeyPrefixEquals             types.String `tfsdk:"key_prefix_equals"`
}

type redirect struct {
	HostName             types.String `tfsdk:"host_name"`
	HTTPRedirectCode     types.String `tfsdk:"http_redirect_code"`
	Protocol             types.String `tfsdk:"protocol"`
	ReplaceKeyPrefixWith types.String `tfsdk:"replace_key_prefix_with"`
	ReplaceKeyWith       types.String `tfsdk:"replace_key_with"`
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
						Required:    true,
						Validators:  []validator.String{stringvalidator.LengthAtLeast(1), stringvalidator.NoneOf("/")},
					},
				},
			},
			"error_document": schema.SingleNestedBlock{
				Description: "The object key name to use when a 4XX class error occurs. Replacement must be made for object keys containing special characters (such as carriage returns) when using XML requests.",
				Attributes: map[string]schema.Attribute{
					"key": schema.StringAttribute{
						Description: "The object key.",
						Required:    true,
						Validators:  []validator.String{stringvalidator.LengthBetween(1, 1024)},
					},
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

	client, ok := req.ProviderData.(*s3.APIClient)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *s3.APIClient, got: %T. Please report this issue to the provider developers.", req.ProviderData),
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

	var data *bucketWebsiteConfigurationModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	_, err := r.client.WebsiteApi.PutBucketWebsite(ctx, data.Bucket.ValueString()).PutBucketWebsiteRequest(buildBucketWebsiteConfigurationFromModel(data)).Execute()
	if err != nil {
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

	var data *bucketWebsiteConfigurationModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	output, apiResponse, err := r.client.WebsiteApi.GetBucketWebsite(ctx, data.Bucket.ValueString()).Execute()
	if apiResponse.HttpNotFound() {
		resp.State.RemoveResource(ctx)
		return
	}
	if err != nil {
		resp.Diagnostics.AddError("Failed to read resource", err.Error())
		return
	}

	data = buildBucketWebsiteConfigurationModelFromAPIResponse(output, data)
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

	var data *bucketWebsiteConfigurationModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	_, err := r.client.WebsiteApi.PutBucketWebsite(ctx, data.Bucket.ValueString()).PutBucketWebsiteRequest(buildBucketWebsiteConfigurationFromModel(data)).Execute()
	if err != nil {
		resp.Diagnostics.AddError("Failed to update resource", err.Error())
		return
	}

	output, _, err := r.client.WebsiteApi.GetBucketWebsite(ctx, data.Bucket.ValueString()).Execute()
	if err != nil {
		resp.Diagnostics.AddError("Failed to read resource", err.Error())
		return
	}

	data = buildBucketWebsiteConfigurationModelFromAPIResponse(output, data)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Delete deletes the bucket website configuration.
func (r *bucketWebsiteConfiguration) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	if r.client == nil {
		resp.Diagnostics.AddError("API client not configured", "The provider client is not configured")
		return
	}

	var data *bucketWebsiteConfigurationModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	_, err := r.client.WebsiteApi.DeleteBucketWebsite(ctx, data.Bucket.ValueString()).Execute()
	if err != nil {
		resp.Diagnostics.AddError("Failed to delete resource", err.Error())
		return
	}
}

func buildBucketWebsiteConfigurationModelFromAPIResponse(output *s3.GetBucketWebsiteOutput, data *bucketWebsiteConfigurationModel) *bucketWebsiteConfigurationModel {
	built := &bucketWebsiteConfigurationModel{
		Bucket: data.Bucket,
	}

	if output.IndexDocument != nil {
		built.IndexDocument = &indexDocument{
			Suffix: types.StringPointerValue(output.IndexDocument.Suffix),
		}
	}

	if output.ErrorDocument != nil {
		built.ErrorDocument = &errorDocument{
			Key: types.StringPointerValue(output.ErrorDocument.Key),
		}
	}

	if output.RedirectAllRequestsTo != nil {
		built.RedirectAllRequestsTo = &redirectAllRequestsTo{
			HostName: types.StringPointerValue(output.RedirectAllRequestsTo.HostName),
			Protocol: types.StringPointerValue(output.RedirectAllRequestsTo.Protocol),
		}
	}

	if output.RoutingRules != nil {
		built.RoutingRule = make([]routingRule, 0, len(*output.RoutingRules))
		for _, r := range *output.RoutingRules {
			var rl routingRule
			if r.Condition != nil {
				rl.Condition = &condition{
					HTTPErrorCodeReturnedEquals: types.StringPointerValue(r.Condition.HttpErrorCodeReturnedEquals),
				}
			}
			if r.Redirect != nil {
				rl.Redirect = &redirect{
					HostName:             types.StringPointerValue(r.Redirect.HostName),
					HTTPRedirectCode:     types.StringPointerValue(r.Redirect.HttpRedirectCode),
					Protocol:             types.StringPointerValue(r.Redirect.Protocol),
					ReplaceKeyPrefixWith: types.StringPointerValue(r.Redirect.ReplaceKeyPrefixWith),
					ReplaceKeyWith:       types.StringPointerValue(r.Redirect.ReplaceKeyWith),
				}
			}
			built.RoutingRule = append(built.RoutingRule, rl)
		}
	}

	return built
}

func buildBucketWebsiteConfigurationFromModel(data *bucketWebsiteConfigurationModel) s3.PutBucketWebsiteRequest {
	return s3.PutBucketWebsiteRequest{
		IndexDocument:         buildIndexDocumentFromModel(data.IndexDocument),
		ErrorDocument:         buildErrorDocumentFromModel(data.ErrorDocument),
		RedirectAllRequestsTo: buildRedirectAllRequestsToFromModel(data.RedirectAllRequestsTo),
		RoutingRules:          buildRoutingRulesFromModel(data.RoutingRule),
	}
}

func buildIndexDocumentFromModel(data *indexDocument) *s3.IndexDocument {
	if data == nil {
		return nil
	}

	return &s3.IndexDocument{
		Suffix: data.Suffix.ValueStringPointer(),
	}
}

func buildErrorDocumentFromModel(data *errorDocument) *s3.ErrorDocument {
	if data == nil {
		return nil
	}

	return &s3.ErrorDocument{
		Key: data.Key.ValueStringPointer(),
	}
}

func buildRedirectAllRequestsToFromModel(data *redirectAllRequestsTo) *s3.RedirectAllRequestsTo {
	if data == nil {
		return nil
	}

	return &s3.RedirectAllRequestsTo{
		HostName: data.HostName.ValueStringPointer(),
		Protocol: data.Protocol.ValueStringPointer(),
	}
}

func buildRoutingRulesFromModel(data []routingRule) *[]s3.RoutingRule {
	if len(data) == 0 {
		return nil
	}

	rules := make([]s3.RoutingRule, 0, len(data))
	for _, r := range data {
		var rl s3.RoutingRule
		if r.Condition != nil {
			rl.Condition = &s3.RoutingRuleCondition{
				HttpErrorCodeReturnedEquals: r.Condition.HTTPErrorCodeReturnedEquals.ValueStringPointer(),
				KeyPrefixEquals:             r.Condition.KeyPrefixEquals.ValueStringPointer(),
			}
		}
		if r.Redirect != nil {
			rl.Redirect = &s3.Redirect{
				HostName:             r.Redirect.HostName.ValueStringPointer(),
				HttpRedirectCode:     r.Redirect.HTTPRedirectCode.ValueStringPointer(),
				Protocol:             r.Redirect.Protocol.ValueStringPointer(),
				ReplaceKeyPrefixWith: r.Redirect.ReplaceKeyPrefixWith.ValueStringPointer(),
				ReplaceKeyWith:       r.Redirect.ReplaceKeyWith.ValueStringPointer(),
			}
		}
		rules = append(rules, rl)
	}

	return &rules
}
