package ionoscloud

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	dataplatformService "github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/dataplatform"
)

func dataSourceDataplatformVersions() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDataplatformReadVersions,
		Schema: map[string]*schema.Schema{
			"versions": {
				Type:        schema.TypeList,
				Description: "Managed Dataplatform API versions",
				Computed:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
		Timeouts: &resourceDefaultTimeouts,
	}
}

func dataSourceDataplatformReadVersions(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(SdkBundle).DataplatformClient

	dataplatformVersions, _, err := client.GetVersions(ctx)
	if err != nil {
		diags := diag.FromErr(fmt.Errorf("an error occurred while fetching Dataplatform API versions: %w", err))
		return diags
	}

	dataplatformService.SetVersionsData(d, dataplatformVersions)

	return nil

}
