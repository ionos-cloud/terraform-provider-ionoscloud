package ionoscloud

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services"
	crService "github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/containerregistry"
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
	client := meta.(services.SdkBundle).ContainerClient

	locations, _, err := client.GetAllLocations(ctx)
	if err != nil {
		diags := diag.FromErr(fmt.Errorf("an error occurred while fetching container registry locations: %w", err))
		return diags
	}

	crService.SetCRLocationsData(d, locations)

	return nil

}
