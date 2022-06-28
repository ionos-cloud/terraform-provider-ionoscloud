package ionoscloud

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	dsaasService "github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/dsaas"
)

func dataSourceDSaaSVersions() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDSaaSReadVersions,
		Schema: map[string]*schema.Schema{
			"versions": {
				Type:        schema.TypeList,
				Description: "Managed Data Stack API versions",
				Computed:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
		Timeouts: &resourceDefaultTimeouts,
	}
}

func dataSourceDSaaSReadVersions(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(SdkBundle).DSaaSClient

	dsaasVersions, _, err := client.GetVersions(ctx)
	if err != nil {
		diags := diag.FromErr(fmt.Errorf("an error occurred while fetching Data Stack API versions: %w", err))
		return diags
	}

	dsaasService.SetVersionsData(d, dsaasVersions)

	return nil

}
