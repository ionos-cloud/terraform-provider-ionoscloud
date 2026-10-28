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
	const pageSize int32 = 100
	var (
		all         []mongov3.InstanceTemplateRead
		offset      int32
		apiResponse *shared.APIResponse
	)
	for {
		request := c.sdkClient.InstanceTemplatesApi.InstancetemplatesGet(ctx).Offset(offset).Limit(pageSize)
		if filterEdition != "" {
			request = request.FilterEdition(mongov3.MongoEdition(filterEdition))
		}
		page, resp, err := request.Execute()
		resp.LogInfo()
		apiResponse = resp
		if err != nil {
			return mongov3.InstanceTemplateReadList{}, apiResponse, err
		}
		all = append(all, page.Items...)
		if len(page.Items) < int(pageSize) {
			page.Items = all
			return page, apiResponse, nil
		}
		offset += pageSize
	}
}
