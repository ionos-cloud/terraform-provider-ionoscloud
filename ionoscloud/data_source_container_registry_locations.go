package ionoscloud

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/bundleclient"
	crService "github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/containerregistry"
	diagutil "github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils/diags"
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

func dataSourceContainerRegistryLocationsRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client, err := meta.(bundleclient.SdkBundle).NewContainerRegistryClient("")
	if err != nil {
		return diag.FromErr(err)
	}

	locations, apiResponse, err := client.GetAllLocations(ctx)
	if err != nil {
		return diagutil.ToDiags(d, fmt.Errorf("an error occurred while fetching container registry locations: %w", err), &diagutil.ErrorContext{StatusCode: apiResponse.SafeStatusCode()})
	}

	crService.SetCRLocationsData(d, locations)

	return nil

}
