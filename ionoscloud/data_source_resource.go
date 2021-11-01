package ionoscloud

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	ionoscloud "github.com/ionos-cloud/sdk-go/v6"
)

func dataSourceResource() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceResourceRead,
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

func dataSourceResourceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(SdkBundle).CloudApiClient

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
			return fmt.Errorf("an error occured while fetching resource by type %s", err)
		}
		results = append(results, result)

		err = d.Set("resource_type", result.Type)
		if err != nil {
			return err
		}
		err = d.Set("resource_id", result.Id)
		if err != nil {
			return err
		}
	} else if resourceType != "" {
		//items, err := client.ListResourcesByType(resource_type)
		items, apiResponse, err := client.UserManagementApi.UmResourcesFindByType(ctx, resourceType).Execute()
		logApiRequestTime(apiResponse)
		if err != nil {
			return fmt.Errorf("an error occured while fetching resources by type %s", err)
		}

		results = *items.Items
		err = d.Set("resource_type", results[0].Type)
		if err != nil {
			return err
		}
	} else {
		//items, err := client.ListResources()
		items, apiResponse, err := client.UserManagementApi.UmResourcesGet(ctx).Execute()
		logApiRequestTime(apiResponse)
		if err != nil {
			return fmt.Errorf("an error occured while fetching resources %s", err)
		}
		results = *items.Items
	}

	if len(results) == 0 {
		return fmt.Errorf("there are no resources that match the search criteria")
	}

	d.SetId(*results[0].Id)

	return nil
}
