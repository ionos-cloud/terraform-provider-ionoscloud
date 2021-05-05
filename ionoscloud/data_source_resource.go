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
	client := meta.(SdkBundle).Client

	var results []ionoscloud.Resource

	resource_type := d.Get("resource_type").(string)
	resource_id := d.Get("resource_id").(string)

	ctx, cancel := context.WithTimeout(context.Background(), *resourceDefaultTimeouts.Default)
	if cancel != nil {
		defer cancel()
	}

	if resource_type != "" && resource_id != "" {
		result, _, err := client.UserManagementApi.UmResourcesFindByTypeAndId(ctx, resource_type, resource_id).Execute()
		if err != nil {
			return fmt.Errorf("An error occured while fetching resource by type %s", err)
		}
		results = append(results, result)

		d.Set("resource_type", result.Type)
		d.Set("resource_id", result.Id)
	} else if resource_type != "" {
		//items, err := client.ListResourcesByType(resource_type)
		items, _, err := client.UserManagementApi.UmResourcesFindByType(ctx, resource_type).Execute()
		if err != nil {
			return fmt.Errorf("An error occured while fetching resources by type %s", err)
		}

		results = *items.Items
		d.Set("resource_type", results[0].Type)
	} else {
		//items, err := client.ListResources()
		items, _, err := client.UserManagementApi.UmResourcesGet(ctx).Execute()
		if err != nil {
			return fmt.Errorf("An error occured while fetching resources %s", err)
		}
		results = *items.Items
	}

	if len(results) > 1 {
		return fmt.Errorf("There is more than one resource that match the search criteria")
	}

	if len(results) == 0 {
		return fmt.Errorf("There are no resources that match the search criteria")
	}

	d.SetId(*results[0].Id)

	return nil
}
