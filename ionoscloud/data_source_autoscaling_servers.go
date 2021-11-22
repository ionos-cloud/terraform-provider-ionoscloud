package ionoscloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func dataSourceAutoscalingServers() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceAutoscalingGroupRead,
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
