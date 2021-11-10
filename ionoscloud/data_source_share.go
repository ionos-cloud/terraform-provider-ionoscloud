package ionoscloud

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func dataSourceShare() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceShareRead,
		Schema: map[string]*schema.Schema{
			"edit_privilege": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"share_privilege": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"group_id": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.All(validation.StringIsNotWhiteSpace),
			},
			"resource_id": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.All(validation.StringIsNotWhiteSpace),
			},
			"id": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
		Timeouts: &resourceDefaultTimeouts,
	}
}

func dataSourceShareRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	id, idOk := d.GetOk("id")

	if !idOk {
		return diag.FromErr(fmt.Errorf("please provide the share id"))
	}
	d.SetId(id.(string))

	if diags := resourceShareRead(ctx, d, meta); diags != nil {
		return diags
	}

	return nil
}
