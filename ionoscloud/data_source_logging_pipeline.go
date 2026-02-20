package ionoscloud

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/ionos-cloud/sdk-go-bundle/products/logging/v2"
	"github.com/ionos-cloud/sdk-go-bundle/shared"

	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/bundleclient"
	loggingService "github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/logging"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils"
)

func dataSourceLoggingPipeline() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourcePipelineRead,
		Schema: map[string]*schema.Schema{
			"location": {
				Type:        schema.TypeString,
				Description: fmt.Sprintf("The location of your logging pipeline. Default: de/txl. Supported locations: %s", strings.Join(loggingService.AvailableLocations, ", ")),
				Optional:    true,
				Default:     "de/txl",
				// no diff in case it moves from "" to de/txl since it's an upgrade from when we had no location
				DiffSuppressFunc: func(_, old, new string, _ *schema.ResourceData) bool {
					if old == "" && new == "de/txl" {
						return true
					}
					return false
				},
			},
			"id": {
				Type:             schema.TypeString,
				Description:      "The ID of the Logging pipeline",
				ValidateDiagFunc: validation.ToDiagFunc(validation.IsUUID),
				Optional:         true,
				Computed:         true,
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
			"http_address": {
				Type:        schema.TypeString,
				Description: "The HTTP address of the Logging pipeline",
				Computed:    true,
			},
			"tcp_address": {
				Type:        schema.TypeString,
				Description: "The TCP address of the Logging pipeline",
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
	client := meta.(bundleclient.SdkBundle).LoggingClient
	location := d.Get("location").(string)
	idValue, idOk := d.GetOk("id")
	nameValue, nameOk := d.GetOk("name")
	id := idValue.(string)
	name := nameValue.(string)

	if idOk && nameOk {
		return utils.ToDiags(d, "ID and name cannot be both specified at the same time", nil)
	}
	if !idOk && !nameOk {
		return utils.ToDiags(d, "please provide either the Logging pipeline ID or name", nil)
	}

	var pipeline logging.PipelineRead
	var apiResponse *shared.APIResponse
	var err error
	if idOk {
		pipeline, apiResponse, err = client.GetPipelineByID(ctx, location, id)
		if err != nil {
			return utils.ToDiags(d, fmt.Sprintf("an error occurred while fetching the Logging pipeline with ID: %s, error: %s", id, err), &utils.DiagsOpts{StatusCode: apiResponse.StatusCode})
		}
	} else {
		var results []logging.PipelineRead
		pipelines, apiResponse, err := client.ListPipelines(ctx, location)
		if err != nil {
			return utils.ToDiags(d, fmt.Sprintf("an error occurred while fetching Logging pipelines: %s", err), &utils.DiagsOpts{StatusCode: apiResponse.StatusCode})
		}
		for _, pipelineItem := range pipelines.Items {
			if strings.EqualFold(pipelineItem.Properties.Name, name) {
				results = append(results, pipelineItem)
			}
		}
		if results == nil || len(results) == 0 {
			return utils.ToDiags(d, fmt.Sprintf("no Logging pipelines found with the specified name: %s", name), nil)
		} else if len(results) > 1 {
			return utils.ToDiags(d, fmt.Sprintf("more than one Logging pipeline found with the specified name: %s", name), nil)
		} else {
			pipeline = results[0]
		}
	}
	if err := client.SetPipelineData(d, pipeline); err != nil {
		return utils.ToDiags(d, err.Error(), nil)
	}
	return nil
}
