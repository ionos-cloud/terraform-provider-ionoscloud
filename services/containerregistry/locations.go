package containerregistry

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	cr "github.com/ionos-cloud/sdk-go-bundle/products/containerregistry"
	"github.com/ionos-cloud/sdk-go-bundle/shared"
)

func (c *Client) GetAllLocations(ctx context.Context) (cr.LocationsResponse, *shared.APIResponse, error) {
	versions, apiResponse, err := c.sdkClient.LocationsApi.LocationsGet(ctx).Execute()
	apiResponse.LogInfo()
	return versions, apiResponse, err
}

func SetCRLocationsData(d *schema.ResourceData, locations cr.LocationsResponse) diag.Diagnostics {

	resourceId := uuid.New()
	d.SetId(resourceId.String())

	if locations.Items != nil {
		var locationList []string
		for _, location := range *locations.Items {
			locationList = append(locationList, *location.Id)
		}
		err := d.Set("locations", locationList)
		if err != nil {
			diags := diag.FromErr(fmt.Errorf("error while setting locations: %w", err))
			return diags
		}

	}
	return nil
}
