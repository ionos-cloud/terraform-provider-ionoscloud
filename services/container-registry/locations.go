package container_registry

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	cr "github.com/ionos-cloud/sdk-go-autoscaling"
)

type LocationsService interface {
	GetAllLocations(ctx context.Context) (cr.LocationsResponse, *cr.APIResponse, error)
}

func (c *Client) GetAllLocations(ctx context.Context) (cr.LocationsResponse, *cr.APIResponse, error) {
	versions, apiResponse, err := c.LocationsApi.LocationsGet(ctx).Execute()
	if apiResponse != nil {
		return versions, apiResponse, err

	}
	return versions, nil, err
}

func SetCRLocationsData(d *schema.ResourceData, locations *cr.LocationsResponse) diag.Diagnostics {

	resourceId := uuid.New()
	d.SetId(resourceId.String())

	if locations.Items != nil {
		var locationList []string
		for _, location := range *locations.Items {
			locationList = append(locationList, *location.Id)
		}
		err := d.Set("locations", locationList)
		if err != nil {
			diags := diag.FromErr(fmt.Errorf("error while setting locations: %s", err))
			return diags
		}

	}
	return nil
}
