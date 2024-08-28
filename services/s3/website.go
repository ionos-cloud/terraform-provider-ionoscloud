package s3

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/types"
	s3 "github.com/ionos-cloud/sdk-go-s3"
)

// BucketWebsiteConfigurationModel defines the expected inputs for creating a new BucketWebsiteConfiguration.
type BucketWebsiteConfigurationModel struct {
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

// CreateBucketWebsite creates a new BucketWebsiteConfiguration.
func (c *Client) CreateBucketWebsite(ctx context.Context, data *BucketWebsiteConfigurationModel) error {
	_, err := c.client.WebsiteApi.PutBucketWebsite(ctx, data.Bucket.ValueString()).
		PutBucketWebsiteRequest(buildBucketWebsiteConfigurationFromModel(data)).Execute()
	return err
}

// GetBucketWebsite gets a BucketWebsiteConfiguration.
func (c *Client) GetBucketWebsite(ctx context.Context, bucketName types.String) (*BucketWebsiteConfigurationModel, bool, error) {
	output, apiResponse, err := c.client.WebsiteApi.GetBucketWebsite(ctx, bucketName.ValueString()).Execute()
	if apiResponse.HttpNotFound() {
		return nil, false, nil
	}

	if err != nil {
		return nil, false, err
	}

	return buildBucketWebsiteConfigurationModelFromAPIResponse(output, &BucketWebsiteConfigurationModel{Bucket: bucketName}), true, nil
}

// UpdateBucketWebsite updates a BucketWebsiteConfiguration.
func (c *Client) UpdateBucketWebsite(ctx context.Context, data *BucketWebsiteConfigurationModel) error {
	if err := c.CreateBucketWebsite(ctx, data); err != nil {
		return err
	}

	model, found, err := c.GetBucketWebsite(ctx, data.Bucket)
	if !found {
		return fmt.Errorf("bucket website configuration not found for bucket %s", data.Bucket.ValueString())
	}

	if err != nil {
		return err
	}

	*data = *model
	return nil
}

// DeleteBucketWebsite deletes a BucketWebsiteConfiguration.
func (c *Client) DeleteBucketWebsite(ctx context.Context, bucketName types.String) error {
	_, err := c.client.WebsiteApi.DeleteBucketWebsite(ctx, bucketName.ValueString()).Execute()
	return err
}

func buildBucketWebsiteConfigurationModelFromAPIResponse(output *s3.GetBucketWebsiteOutput, data *BucketWebsiteConfigurationModel) *BucketWebsiteConfigurationModel {
	built := &BucketWebsiteConfigurationModel{
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

func buildBucketWebsiteConfigurationFromModel(data *BucketWebsiteConfigurationModel) s3.PutBucketWebsiteRequest {
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
