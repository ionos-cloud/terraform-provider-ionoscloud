package ionoscloud

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	ionoscloud "github.com/ionos-cloud/sdk-go/v6"

	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/bundleclient"
	diagutil "github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils/diags"
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
	client, err := meta.(bundleclient.SdkBundle).NewCloudAPIClientWithFailover(ctx)
	if err != nil {
		return diag.FromErr(err)
	}

	var results []ionoscloud.Resource

	resourceType := d.Get("resource_type").(string)
	resourceId := d.Get("resource_id").(string)

	ctx, cancel := context.WithTimeout(ctx, *resourceDefaultTimeouts.Default)
	if cancel != nil {
		defer cancel()
	}

	if resourceType != "" && resourceId != "" {
		result, apiResponse, err := client.UserManagementApi.UmResourcesFindByTypeAndId(ctx, resourceType, resourceId).Execute()
		logApiRequestTime(apiResponse)
		if err != nil {
			return diagutil.ToDiags(d, fmt.Errorf("an error occurred while fetching resource by type %w", err), &diagutil.ErrorContext{StatusCode: apiResponse.SafeStatusCode()})
		}
		results = append(results, result)

		err = d.Set("resource_type", result.Type)
		if err != nil {
			return diagutil.ToDiags(d, err, nil)
		}
		err = d.Set("resource_id", result.Id)
		if err != nil {
			return diagutil.ToDiags(d, err, nil)
		}
	} else if resourceType != "" {
		// items, err := client.ListResourcesByType(resource_type)
		items, apiResponse, err := client.UserManagementApi.UmResourcesFindByType(ctx, resourceType).Execute()
		logApiRequestTime(apiResponse)
		if err != nil {
			return diagutil.ToDiags(d, fmt.Errorf("an error occurred while fetching resources by type %w", err), &diagutil.ErrorContext{StatusCode: apiResponse.SafeStatusCode()})
		}
		results = *items.Items
		if len(results) > 0 && results[0].Type != nil {
			err = d.Set("resource_type", results[0].Type)
			if err != nil {
				return diagutil.ToDiags(d, err, nil)
			}
		}

	} else {
		// items, err := client.ListResources()
		items, apiResponse, err := client.UserManagementApi.UmResourcesGet(ctx).Depth(1).Execute()
		logApiRequestTime(apiResponse)
		if err != nil {
			return diagutil.ToDiags(d, fmt.Errorf("an error occurred while fetching resources %w", err), &diagutil.ErrorContext{StatusCode: apiResponse.SafeStatusCode()})
		}
		results = *items.Items
	}

	if len(results) == 0 {
		return diagutil.ToDiags(d, fmt.Errorf("there are no resources that match the search criteria"), nil)
	}

	d.SetId(*results[0].Id)

	return nil
}
