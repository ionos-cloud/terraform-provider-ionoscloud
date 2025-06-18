package ionoscloud

import (
	"context"
	"errors"
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
				Description: "The address of the client's grafana instance",
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
	pipelineResponse, _, err := client.CreatePipeline(ctx, d)
	if err != nil {
		return diag.FromErr(fmt.Errorf("an error occurred while creating a Logging pipeline: %w", err))
	}
	d.SetId(*pipelineResponse.Id)
	err = utils.WaitForResourceToBeReady(ctx, d, client.IsPipelineAvailable)
	if err != nil {
		return diag.FromErr(fmt.Errorf("an error occurred while waiting for the pipeline with ID: %s to become available: %w", *pipelineResponse.Id, err))
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
		return diag.FromErr(fmt.Errorf("error while fetching Logging pipeline with ID: %s, err: %w", pipelineID, err))
	}

	log.Printf("[INFO] Successfully retrieved Logging pipeline with ID: %s: %+v", pipelineID, pipeline)
	if err := client.SetPipelineData(d, pipeline); err != nil {
		return diag.FromErr(err)
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
		return diag.FromErr(fmt.Errorf("error while deleting Logging pipeline with ID: %s, location %s, error: %w", pipelineID, location, err))
	}

	err = utils.WaitForResourceToBeDeleted(ctx, d, client.IsPipelineDeleted)
	if err != nil {
		return diag.FromErr(fmt.Errorf("an error occurred while wainting for Logging pipeline with ID: %s to be deleted, error: %w", pipelineID, err))
	}
	return nil
}

func pipelineUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(bundleclient.SdkBundle).LoggingClient
	pipelineID := d.Id()

	pipelineResponse, _, err := client.UpdatePipeline(ctx, pipelineID, d)
	if err != nil {
		return diag.FromErr(fmt.Errorf("an error occurred while updating the Logging pipeline with ID: %s, error: %w", pipelineID, err))
	}
	err = utils.WaitForResourceToBeReady(ctx, d, client.IsPipelineAvailable)
	if err != nil {
		return diag.FromErr(fmt.Errorf("an error occurred while waiting for the Logging pipeline with ID: %s to become available: %w", pipelineID, err))
	}
	if err := client.SetPipelineData(d, pipelineResponse); err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func pipelineImport(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	parts := strings.Split(d.Id(), ":")
	if len(parts) != 2 {
		return nil, fmt.Errorf("expected ID in the format location:id")
	}

	location := parts[0]
	if err := d.Set("location", location); err != nil {
		return nil, fmt.Errorf("failed to set location for Logging Pipeline import: %w", err)
	}
	d.SetId(parts[1])

	diags := pipelineRead(ctx, d, meta)
	if diags != nil && diags.HasError() {
		return nil, errors.New(diags[0].Summary)
	}
	return []*schema.ResourceData{d}, nil
}
