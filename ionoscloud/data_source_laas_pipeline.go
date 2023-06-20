package ionoscloud

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func dataSourceLaaSPipeline() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourcePipelineRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:             schema.TypeString,
				Description:      "The ID of the LaaS pipeline",
				ValidateDiagFunc: validation.ToDiagFunc(validation.IsUUID),
				Required:         true,
			},
			"name": {
				Type:        schema.TypeString,
				Description: "The name of the LaaS pipeline",
				Computed:    true,
			},
			"log": {
				Type:        schema.TypeList,
				Description: "The logs for the LaaS pipeline",
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"source": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The source parser to be used",
						},
						"tag": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The tag is used to distinguish different pipelines. Must be unique amongst the pipeline's array items.",
						},
						"protocol": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"public": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"destinations": {
							Type:     schema.TypeSet,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"type": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func dataSourcePipelineRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(SdkBundle).LaaSClient
	idValue := d.Get("id")
	id := idValue.(string)

	pipeline, _, err := client.GetPipelineById(ctx, id)
	if err != nil {
		return diag.FromErr(fmt.Errorf("an error occured while fetching the LaaS pipeline with ID: %s, error: %w", id, err))
	}
	if err := client.SetPipelineData(d, pipeline); err != nil {
		return diag.FromErr(err)
	}

	return nil
}
