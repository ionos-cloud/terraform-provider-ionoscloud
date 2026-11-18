package mongov2

import (
	"context"

	mongov3 "github.com/ionos-cloud/sdk-go-bundle/products/dbaas/mongo/v3"
	"github.com/ionos-cloud/sdk-go-bundle/shared"
)

// GetTemplate retrieves an instance template by ID.
func (c *Client) GetTemplate(ctx context.Context, templateID string) (mongov3.InstanceTemplateRead, *shared.APIResponse, error) {
	template, apiResponse, err := c.sdkClient.InstanceTemplatesApi.InstancetemplatesFindById(ctx, templateID).Execute()
	apiResponse.LogInfo()
	return template, apiResponse, err
}

// ListTemplates retrieves a list of instance templates with an optional edition filter.
func (c *Client) ListTemplates(ctx context.Context, filterEdition string) (mongov3.InstanceTemplateReadList, *shared.APIResponse, error) {
	request := c.sdkClient.InstanceTemplatesApi.InstancetemplatesGet(ctx)
	if filterEdition != "" {
		request = request.FilterEdition(mongov3.MongoEdition(filterEdition))
	}
	templates, apiResponse, err := request.Execute()
	apiResponse.LogInfo()
	return templates, apiResponse, err
}
