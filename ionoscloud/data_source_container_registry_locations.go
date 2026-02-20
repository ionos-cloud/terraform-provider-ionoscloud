package ionoscloud

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/bundleclient"
	crService "github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/containerregistry"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils"
)

func dataSourceContainerRegistryLocations() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceContainerRegistryLocationsRead,
		Schema: map[string]*schema.Schema{
			"locations": {
				Type:        schema.TypeList,
				Description: "list of container registry locations",
				Computed:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
		Timeouts: &resourceDefaultTimeouts,
	}
}

func dataSourceContainerRegistryLocationsRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(bundleclient.SdkBundle).ContainerClient

	locations, apiResponse, err := client.GetAllLocations(ctx)
	if err != nil {
		return utils.ToDiags(d, fmt.Sprintf("an error occurred while fetching container registry locations: %s", err), &utils.DiagsOpts{StatusCode: apiResponse.StatusCode})
	}

	crService.SetCRLocationsData(d, locations)

	return nil

}
