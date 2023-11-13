package ionoscloud

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services"
	autoscalingService "github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/autoscaling"
)

func dataSourceAutoscalingGroupServers() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceAutoscalingServersRead,
		Schema: map[string]*schema.Schema{
			"group_id": {
				Type:             schema.TypeString,
				Description:      "Unique identifier for the group",
				Required:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringIsNotWhiteSpace),
			},
			"servers": {
				Type:     schema.TypeList,
				Computed: true, Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Description: "Unique identifier for the resource",
							Computed:    true,
						},
					}}},
		},
	}
}

func dataSourceAutoscalingServersRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(services.SdkBundle).AutoscalingClient

	id, idOk := d.GetOk("group_id")

	if !idOk {
		diags := diag.FromErr(fmt.Errorf("autoscaling group_id has to be provided in order to search for its servers"))
		return diags
	}

	groupServers, _, err := client.GetAllGroupServers(ctx, id.(string))
	if err != nil {
		return diag.FromErr(fmt.Errorf("an error occurred while fetching group with ID %s: %w", id.(string), err))
	}

	return autoscalingService.SetAutoscalingServersData(d, groupServers)
}
