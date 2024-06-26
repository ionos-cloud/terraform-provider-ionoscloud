package ionoscloud

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceApiGateway() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceApiGatewayCreate,
		ReadContext:   resourceApiGatewayRead,
		UpdateContext: resourceApiGatewayUpdate,
		DeleteContext: resourceApiGatewayDelete,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Description: "The name of the API Gateway",
				Required:    true,
			},
			"logs": {
				Type:        schema.TypeBool,
				Description: "This field enables or disables the collection and reporting of logs for observability of this instance.",
				Optional:    true,
				Default:     false,
			},
			"metrics": {
				Type:        schema.TypeBool,
				Description: "This field enables or disables the collection and reporting of metrics for observability of this instance.",
				Optional:    true,
				Default:     false,
			},
			"custom_domains": {
				Type:        schema.TypeList,
				Description: "The custom domain that the API Gateway instance should listen on.",
				Required:    true,
				MaxItems:    5,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Description: "The domain name of the distribution.",
							Optional: true,
						},
						"certificate_id": {
							Type:        schema.TypeString,
							Description: "The ID of the certificate to use for the distribution.",
							Optional: true,
						},
					},
				},
			},
		},
	}
}

func resourceApiGatewayCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return nil
}

func resourceApiGatewayRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return nil
}

func resourceApiGatewayUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return nil
}

func resourceApiGatewayDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return nil
}
