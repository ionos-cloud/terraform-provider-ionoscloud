package ionoscloud

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func dataSourceIpFailover() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIpFailoverRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.All(validation.StringIsNotWhiteSpace),
			},
			"ip": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.All(validation.IsIPAddress),
			},
			"nicuuid": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.All(validation.IsUUID),
			},
			"lan_id": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.All(validation.StringIsNotWhiteSpace),
			},
			"datacenter_id": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.All(validation.IsUUID),
			},
		},
		Timeouts: &resourceDefaultTimeouts,
	}
}

func dataSourceIpFailoverRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	id, idOk := d.GetOk("id")

	if !idOk {
		return diag.FromErr(fmt.Errorf("please provide the ip failover id"))
	}
	d.SetId(id.(string))

	if diags := resourceLanIPFailoverRead(ctx, d, meta); diags != nil {
		return diags
	}

	return nil
}
