package ionoscloud

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	ionoscloud "github.com/ionos-cloud/sdk-go/v5"
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
	client := meta.(*ionoscloud.APIClient)

	var results []ionoscloud.Resource

	resourceType := d.Get("resource_type").(string)
	resourceId := d.Get("resource_id").(string)

	ctx, cancel := context.WithTimeout(context.Background(), *resourceDefaultTimeouts.Default)
	if cancel != nil {
		defer cancel()
	}

	if resourceType != "" && resourceId != "" {
		result, _, err := client.UserManagementApi.UmResourcesFindByTypeAndId(ctx, resourceType, resourceId).Execute()
		if err != nil {
			return fmt.Errorf("an error occured while fetching resource by type %s", err)
		}
		results = append(results, result)

		d.Set("resource_type", result.Type)
		d.Set("resource_id", result.Id)
	} else if resourceType != "" {
		//items, err := client.ListResourcesByType(resource_type)
		items, _, err := client.UserManagementApi.UmResourcesFindByType(ctx, resourceType).Execute()
		if err != nil {
			return fmt.Errorf("an error occured while fetching resources by type %s", err)
		}

		results = *items.Items
		d.Set("resource_type", results[0].Type)
	} else {
		//items, err := client.ListResources()
		items, _, err := client.UserManagementApi.UmResourcesGet(ctx).Execute()
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
