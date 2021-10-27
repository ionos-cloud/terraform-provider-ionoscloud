package ionoscloud

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func dataSourceIpBlock() *schema.Resource {
	return &schema.Resource{
		ReadContext: datasourceIpBlockRead,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"location": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.All(validation.StringIsNotWhiteSpace),
			},
			"size": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
			},
			"ips": {
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"ip_consumers": {
				Type:     schema.TypeList,
				Computed: true,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ip": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"mac": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"nic_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"server_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"server_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"datacenter_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"datacenter_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"k8s_nodepool_uuid": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"k8s_cluster_uuid": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
		Timeouts: &resourceDefaultTimeouts,
	}

}

func datasourceIpBlockRead(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	id, idOk := data.GetOk("id")

	if !idOk {
		return diag.FromErr(fmt.Errorf("please provide the ip block id"))
	}
	data.SetId(id.(string))

	if diags := resourceIPBlockRead(ctx, data, meta); diags != nil {
		return diags
	}

	return nil
}
