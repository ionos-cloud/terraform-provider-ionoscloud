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

// DataSourceAutoscalingGroupServers defines the schema for the Autoscaling Group Servers data source
func DataSourceAutoscalingGroupServers() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceAutoscalingServersRead,
		Schema: map[string]*schema.Schema{
			"group_id": {
				Type:             schema.TypeString,
				Description:      "Unique identifier for the group",
				Required:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.IsUUID),
			},
			"servers": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
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
		return diag.FromErr(fmt.Errorf("autoscaling group_id has to be provided in order to search for its servers"))
	}

	groupServers, _, err := client.GetAllGroupServers(ctx, id.(string))
	if err != nil {
		return diag.FromErr(fmt.Errorf("an error occurred while fetching the servers for the group with ID %s: %w", id.(string), err))
	}

	return autoscalingService.SetAutoscalingServersData(d, groupServers)
}
