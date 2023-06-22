package ionoscloud

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	laas "github.com/ionos-cloud/sdk-go-laas"
	"strings"
)

func dataSourceLaaSPipeline() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourcePipelineRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:             schema.TypeString,
				Description:      "The ID of the LaaS pipeline",
				ValidateDiagFunc: validation.ToDiagFunc(validation.IsUUID),
				Optional:         true,
			},
			"name": {
				Type:        schema.TypeString,
				Description: "The name of the LaaS pipeline",
				Optional:    true,
				Computed:    true,
			},
			"log": {
				Type:        schema.TypeSet,
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
							Type:        schema.TypeString,
							Description: "Protocol to use as intake. Possible values are: http, tcp.",
							Computed:    true,
						},
						"public": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"destinations": {
							Type:        schema.TypeSet,
							Computed:    true,
							Description: "The internal output stream to send logs to. Possible values are: loki.",
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
	idValue, idOk := d.GetOk("id")
	nameValue, nameOk := d.GetOk("name")
	id := idValue.(string)
	name := nameValue.(string)

	if idOk && nameOk {
		return diag.FromErr(fmt.Errorf("ID and name cannot be both specified at the same time"))
	}
	if !idOk && !nameOk {
		return diag.FromErr(fmt.Errorf("please provide either the LaaS pipeline ID or name"))
	}

	var pipeline laas.Pipeline
	var err error
	if idOk {
		pipeline, _, err = client.GetPipelineById(ctx, id)
		if err != nil {
			return diag.FromErr(fmt.Errorf("an error occured while fetching the LaaS pipeline with ID: %s, error: %w", id, err))
		}
	} else {
		var results []laas.Pipeline
		pipelines, _, err := client.ListPipelines(ctx)
		if err != nil {
			return diag.FromErr(fmt.Errorf("an error occured while fetching LaaS pipelines: %w", err))
		}
		for _, pipelineItem := range *pipelines.Items {
			if pipelineItem.Properties != nil && pipelineItem.Properties.Name != nil && strings.EqualFold(*pipelineItem.Properties.Name, name) {
				results = append(results, pipelineItem)
			}
		}
		if results == nil || len(results) == 0 {
			return diag.FromErr(fmt.Errorf("no LaaS pipelines found with the specified name: %s", name))
		} else if len(results) > 1 {
			return diag.FromErr(fmt.Errorf("more than one LaaS pipeline found with the specified name: %s", name))
		} else {
			pipeline = results[0]
		}
	}
	if err := client.SetPipelineData(d, pipeline); err != nil {
		return diag.FromErr(err)
	}
	return nil
}
