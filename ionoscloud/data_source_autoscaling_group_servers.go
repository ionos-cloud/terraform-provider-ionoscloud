package ionoscloud

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	autoscalingService "github.com/ionos-cloud/terraform-provider-ionoscloud/services/autoscaling"
)

func dataSourceAutoscalingGroupServers() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceAutoscalingServersRead,
		Schema: map[string]*schema.Schema{
			"group_id": {
				Type:         schema.TypeString,
				Description:  "Unique identifier for the group",
				Required:     true,
				ValidateFunc: validation.All(validation.StringIsNotWhiteSpace),
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
	client := meta.(SdkBundle).AutoscalingClient

	id, idOk := d.GetOk("group_id")

	if !idOk {
		diags := diag.FromErr(fmt.Errorf("autoscaling group_id has to be provided in order to search for its servers"))
		return diags
	}

	/* search by ID */
	groupServers, _, err := client.GetAllGroupServers(ctx, id.(string))

	if err != nil {
		diags := diag.FromErr(fmt.Errorf("an error occurred while fetching backup for cluster with ID %s: %s", id.(string), err))
		return diags
	}

	if diags := autoscalingService.SetAutoscalingServersData(d, groupServers); diags != nil {
		return diags
	}

	return nil
}
