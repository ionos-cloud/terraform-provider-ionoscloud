package ionoscloud

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	validation "github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils"
	"log"
)

func resourceLaaSPipeline() *schema.Resource {
	return &schema.Resource{
		CreateContext: pipelineCreate,
		ReadContext:   pipelineRead,
		UpdateContext: pipelineUpdate,
		DeleteContext: pipelineDelete,
		Importer: &schema.ResourceImporter{
			StateContext: pipelineImport,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"log": {
				Type:     schema.TypeSet,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"source": {
							Type:         schema.TypeString,
							Required:     true,
							Description:  "The source parser to be used",
							ValidateFunc: validation.StringInSlice([]string{"kubernetes", "docker", "systemd"}, false),
						},
						"tag": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The tag is used to distinguish different pipelines. Must be unique amongst the pipeline's array items.",
						},
						"protocol": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringInSlice([]string{"http", "tcp"}, false),
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
		Timeouts: &resourceDefaultTimeouts,
	}
}

func pipelineCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(SdkBundle).LaaSClient
	pipelineResponse, _, err := client.CreatePipeline(ctx, d)
	if err != nil {
		return diag.FromErr(fmt.Errorf("an error occured while creating a LaaS pipeline: %w", err))
	}
	d.SetId(*pipelineResponse.Id)
	err = utils.WaitForResourceToBeReady(ctx, d, client.IsPipelineAvailable)
	if err != nil {
		return diag.FromErr(fmt.Errorf("an error occured while waiting for the pipeline with ID: %s to become available: %w", *pipelineResponse.Id, err))
	}
	return pipelineRead(ctx, d, meta)
}

func pipelineRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(SdkBundle).LaaSClient
	pipelineId := d.Id()
	pipeline, apiResponse, err := client.GetPipelineById(ctx, pipelineId)

	if err != nil {
		// TODO -- leave it as it is for now but check if this can be done in a better way
		if apiResponse.HttpNotFound() {
			log.Printf("[INFO] Could not find LaaS pipeline with ID: %s", pipelineId)
			d.SetId("")
			return nil
		}
		return diag.FromErr(fmt.Errorf("error while fetching LaaS pipeline with ID: %s, err: %w", pipelineId, err))
	}

	log.Printf("[INFO] Successfully retrieved LaaS pipeline with ID: %s: %+v", pipelineId, pipeline)
	if err := client.SetPipelineData(d, pipeline); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func pipelineDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(SdkBundle).LaaSClient
	pipelineId := d.Id()

	apiResponse, err := client.DeletePipeline(ctx, pipelineId)
	if err != nil {
		if apiResponse.HttpNotFound() {
			d.SetId("")
			return nil
		}
		return diag.FromErr(fmt.Errorf("error while deleting LaaS pipeline with ID: %s, error: %w", pipelineId, err))
	}

	err = utils.WaitForResourceToBeDeleted(ctx, d, client.IsPipelineDeleted)
	if err != nil {
		return diag.FromErr(fmt.Errorf("an error occured while wainting for LaaS pipeline with ID: %s to be deleted, error: %w", pipelineId, err))
	}
	return nil
}

func pipelineUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(SdkBundle).LaaSClient
	pipelineId := d.Id()

	_, _, err := client.UpdatePipeline(ctx, pipelineId, d)
	if err != nil {
		return diag.FromErr(fmt.Errorf("an error occured while updating the LaaS pipeline with ID: %s, error: %w", pipelineId, err))
	}
	err = utils.WaitForResourceToBeReady(ctx, d, client.IsPipelineAvailable)
	if err != nil {
		return diag.FromErr(fmt.Errorf("an error occured while waiting for the LaaS pipeline with ID: %s to become available: %w", pipelineId, err))
	}
	return pipelineRead(ctx, d, meta)
}

func pipelineImport(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	pipelineRead(ctx, d, meta)
	return []*schema.ResourceData{d}, nil
}
