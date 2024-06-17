package ionoscloud

import (
	"context"
	"fmt"
	"strings"

	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/ionos-cloud/sdk-go-bundle/products/logging/v2"
)

func dataSourceLoggingPipeline() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourcePipelineRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:             schema.TypeString,
				Description:      "The ID of the Logging pipeline",
				ValidateDiagFunc: validation.ToDiagFunc(validation.IsUUID),
				Optional:         true,
			},
			"name": {
				Type:        schema.TypeString,
				Description: "The name of the Logging pipeline",
				Optional:    true,
				Computed:    true,
			},
			"grafana_address": {
				Type:        schema.TypeString,
				Description: "The address of the client's grafana instance",
				Computed:    true,
			},
			"log": {
				Type:        schema.TypeSet,
				Description: "The logs for the Logging pipeline",
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
									"retention_in_days": {
										Type:        schema.TypeInt,
										Optional:    true,
										Computed:    true,
										Description: "Defines the number of days a log record should be kept in loki. Works with loki destination type only. Possible values are: 7, 14, 30.",
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
	client := meta.(services.SdkBundle).LoggingClient
	idValue, idOk := d.GetOk("id")
	nameValue, nameOk := d.GetOk("name")
	id := idValue.(string)
	name := nameValue.(string)

	if idOk && nameOk {
		return diag.FromErr(fmt.Errorf("ID and name cannot be both specified at the same time"))
	}
	if !idOk && !nameOk {
		return diag.FromErr(fmt.Errorf("please provide either the Logging pipeline ID or name"))
	}

	var pipeline logging.Pipeline
	var err error
	if idOk {
		pipeline, _, err = client.GetPipelineById(ctx, id)
		if err != nil {
			return diag.FromErr(fmt.Errorf("an error occurred while fetching the Logging pipeline with ID: %s, error: %w", id, err))
		}
	} else {
		var results []logging.Pipeline
		pipelines, _, err := client.ListPipelines(ctx)
		if err != nil {
			return diag.FromErr(fmt.Errorf("an error occurred while fetching Logging pipelines: %w", err))
		}
		for _, pipelineItem := range pipelines.Items {
			if pipelineItem.Properties != nil && pipelineItem.Properties.Name != nil && strings.EqualFold(*pipelineItem.Properties.Name, name) {
				results = append(results, pipelineItem)
			}
		}
		if results == nil || len(results) == 0 {
			return diag.FromErr(fmt.Errorf("no Logging pipelines found with the specified name: %s", name))
		} else if len(results) > 1 {
			return diag.FromErr(fmt.Errorf("more than one Logging pipeline found with the specified name: %s", name))
		} else {
			pipeline = results[0]
		}
	}
	d.SetId(*pipeline.Id)
	if err := client.SetPipelineData(d, pipeline.Properties); err != nil {
		return diag.FromErr(err)
	}
	return nil
}
