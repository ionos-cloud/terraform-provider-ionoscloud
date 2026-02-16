package ionoscloud

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/bundleclient"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/logging"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils"
)

func resourceLoggingPipeline() *schema.Resource {
	return &schema.Resource{
		CreateContext: pipelineCreate,
		ReadContext:   pipelineRead,
		UpdateContext: pipelineUpdate,
		DeleteContext: pipelineDelete,
		Importer: &schema.ResourceImporter{
			StateContext: pipelineImport,
		},
		Schema: map[string]*schema.Schema{
			"location": {
				Type:        schema.TypeString,
				Description: fmt.Sprintf("The location of your logging pipeline. Default: de/txl. Supported locations: %s", strings.Join(logging.AvailableLocations, ", ")),
				Optional:    true,
				ForceNew:    true,
				// no diff in case it moves from "" to de/txl since it's an upgrade from when we had no location
				DiffSuppressFunc: func(_, old, new string, _ *schema.ResourceData) bool {
					if old == "" && new == logging.DefaultLocation {
						return true
					}
					return false
				},
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"grafana_address": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The Grafana address is where user can access their logs, create dashboards, and set up alerts",
			},
			"http_address": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The HTTP address of the pipeline. This is the address to which logs are sent using the HTTP protocol",
			},
			"tcp_address": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The TCP address of the pipeline. This is the address to which logs are sent using the TCP protocol",
			},
			"key": {
				Type:        schema.TypeString,
				Computed:    true,
				Sensitive:   true,
				Description: "The key is shared once and is used to authenticate the logs sent to the pipeline",
			},
			"log": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"source": {
							Type:             schema.TypeString,
							Required:         true,
							Description:      "The source parser to be used",
							ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"kubernetes", "docker", "systemd", "generic"}, false)),
						},
						"tag": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The tag is used to distinguish different pipelines. Must be unique amongst the pipeline's array items.",
						},
						"protocol": {
							Type:             schema.TypeString,
							Required:         true,
							Description:      "Protocol to use as intake. Possible values are: http, tcp.",
							ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"http", "tcp"}, false)),
						},
						"public": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"destinations": {
							Type:        schema.TypeList,
							Optional:    true,
							Computed:    true,
							Description: "The internal output stream to send logs to. Possible values are: loki.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"type": {
										Type:     schema.TypeString,
										Optional: true,
										Computed: true,
									},
									"retention_in_days": {
										Type:             schema.TypeInt,
										Optional:         true,
										Computed:         true,
										Description:      "Defines the number of days a log record should be kept in loki. Works with loki destination type only. Possible values are: 7, 14, 30.",
										ValidateDiagFunc: validation.ToDiagFunc(validation.IntInSlice([]int{7, 14, 30})),
									},
								},
							},
						},
					},
				},
			},
		},
		Timeouts: &resourceDefaultTimeouts,
	}
}

func pipelineCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(bundleclient.SdkBundle).LoggingClient
	pipelineResponse, apiResponse, err := client.CreatePipeline(ctx, d)
	if err != nil {
		return utils.ToDiags(d, fmt.Sprintf("an error occurred while creating a Logging pipeline: %s", err), &utils.DiagsOpts{StatusCode: apiResponse.StatusCode})
	}
	d.SetId(pipelineResponse.Id)
	err = utils.WaitForResourceToBeReady(ctx, d, client.IsPipelineAvailable)
	if err != nil {
		return utils.ToDiags(d, fmt.Sprintf("an error occurred while waiting for the pipeline with ID: %s to become available: %s", pipelineResponse.Id, err), nil)
	}
	// only received in the create response
	if pipelineResponse.Properties.Key != nil {
		err = d.Set("key", pipelineResponse.Properties.Key)
		if err != nil {
			return utils.ToDiags(d, fmt.Sprintf("an error occurred while setting the key for the Logging pipeline with ID: %s, error: %s", pipelineResponse.Id, err), nil)
		}
	}
	// Make another read and set the data in the state because 'grafanaAddress` is not returned in the create response
	return pipelineRead(ctx, d, meta)
}

func pipelineRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(bundleclient.SdkBundle).LoggingClient
	pipelineID := d.Id()
	location := ""
	if newLocation, ok := d.GetOk("location"); ok {
		location = newLocation.(string)
	}
	pipeline, apiResponse, err := client.GetPipelineByID(ctx, location, pipelineID)
	if err != nil {
		if apiResponse.HttpNotFound() {
			log.Printf("[INFO] Could not find Logging pipeline with ID: %s", pipelineID)
			d.SetId("")
			return nil
		}
		return utils.ToDiags(d, fmt.Sprintf("error while fetching Logging pipeline: %s", err), &utils.DiagsOpts{StatusCode: apiResponse.StatusCode})
	}

	log.Printf("[INFO] Successfully retrieved Logging pipeline with ID: %s: %+v", pipelineID, pipeline)
	if err := client.SetPipelineData(d, pipeline); err != nil {
		return utils.ToDiags(d, err.Error(), nil)
	}

	return nil
}

func pipelineDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(bundleclient.SdkBundle).LoggingClient
	pipelineID := d.Id()
	location := d.Get("location").(string)
	apiResponse, err := client.DeletePipeline(ctx, location, pipelineID)
	if err != nil {
		if apiResponse.HttpNotFound() {
			d.SetId("")
			return nil
		}
		return utils.ToDiags(d, fmt.Sprintf("error while deleting Logging pipeline, location %s, error: %s", location, err), &utils.DiagsOpts{StatusCode: apiResponse.StatusCode})
	}

	err = utils.WaitForResourceToBeDeleted(ctx, d, client.IsPipelineDeleted)
	if err != nil {
		return utils.ToDiags(d, fmt.Sprintf("an error occurred while waiting for Logging pipeline to be deleted: %s", err), &utils.DiagsOpts{Timeout: schema.TimeoutDelete})
	}
	return nil
}

func pipelineUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(bundleclient.SdkBundle).LoggingClient
	pipelineID := d.Id()

	pipelineResponse, apiResponse, err := client.UpdatePipeline(ctx, pipelineID, d)
	if err != nil {
		return utils.ToDiags(d, fmt.Sprintf("an error occurred while updating the Logging pipeline: %s", err), &utils.DiagsOpts{StatusCode: apiResponse.StatusCode})
	}
	err = utils.WaitForResourceToBeReady(ctx, d, client.IsPipelineAvailable)
	if err != nil {
		return utils.ToDiags(d, fmt.Sprintf("an error occurred while waiting for the Logging pipeline to become available: %s", err), nil)
	}
	if err := client.SetPipelineData(d, pipelineResponse); err != nil {
		return utils.ToDiags(d, err.Error(), nil)
	}
	return nil
}

func pipelineImport(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	parts := strings.Split(d.Id(), ":")
	if len(parts) != 2 {
		return nil, utils.ToError(d, "expected ID in the format location:id", nil)
	}

	location := parts[0]
	if err := d.Set("location", location); err != nil {
		return nil, utils.ToError(d, fmt.Sprintf("failed to set location for Logging Pipeline import: %s", err), nil)
	}
	d.SetId(parts[1])

	diags := pipelineRead(ctx, d, meta)
	if diags != nil && diags.HasError() {
		return nil, utils.ToError(d, diags[0].Summary, nil)
	}
	return []*schema.ResourceData{d}, nil
}
