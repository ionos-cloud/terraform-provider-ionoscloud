package ionoscloud

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	ionoscloud "github.com/ionos-cloud/sdk-go/v6"

	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/bundleclient"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils"
)

func dataSourceResource() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceResourceRead,
		Schema: map[string]*schema.Schema{
			"resource_type": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"resource_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
		Timeouts: &resourceDefaultTimeouts,
	}
}

func dataSourceResourceRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(bundleclient.SdkBundle).CloudApiClient

	var results []ionoscloud.Resource

	resourceType := d.Get("resource_type").(string)
	resourceId := d.Get("resource_id").(string)

	ctx, cancel := context.WithTimeout(context.Background(), *resourceDefaultTimeouts.Default)
	if cancel != nil {
		defer cancel()
	}

	if resourceType != "" && resourceId != "" {
		result, apiResponse, err := client.UserManagementApi.UmResourcesFindByTypeAndId(ctx, resourceType, resourceId).Execute()
		logApiRequestTime(apiResponse)
		if err != nil {
			return utils.ToDiags(d, fmt.Sprintf("an error occurred while fetching resource by type %s", err), &utils.DiagsOpts{StatusCode: apiResponse.StatusCode})
		}
		results = append(results, result)

		err = d.Set("resource_type", result.Type)
		if err != nil {
			return utils.ToDiags(d, err.Error(), nil)
		}
		err = d.Set("resource_id", result.Id)
		if err != nil {
			return utils.ToDiags(d, err.Error(), nil)
		}
	} else if resourceType != "" {
		// items, err := client.ListResourcesByType(resource_type)
		items, apiResponse, err := client.UserManagementApi.UmResourcesFindByType(ctx, resourceType).Execute()
		logApiRequestTime(apiResponse)
		if err != nil {
			return utils.ToDiags(d, fmt.Sprintf("an error occurred while fetching resources by type %s", err), &utils.DiagsOpts{StatusCode: apiResponse.StatusCode})
		}
		results = *items.Items
		if len(results) > 0 && results[0].Type != nil {
			err = d.Set("resource_type", results[0].Type)
			if err != nil {
				return utils.ToDiags(d, err.Error(), nil)
			}
		}

	} else {
		// items, err := client.ListResources()
		items, apiResponse, err := client.UserManagementApi.UmResourcesGet(ctx).Depth(1).Execute()
		logApiRequestTime(apiResponse)
		if err != nil {
			return utils.ToDiags(d, fmt.Sprintf("an error occurred while fetching resources %s", err), &utils.DiagsOpts{StatusCode: apiResponse.StatusCode})
		}
		results = *items.Items
	}

	if len(results) == 0 {
		return utils.ToDiags(d, "there are no resources that match the search criteria", nil)
	}

	d.SetId(*results[0].Id)

	return nil
}
