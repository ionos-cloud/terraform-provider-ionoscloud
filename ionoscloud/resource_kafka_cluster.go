package ionoscloud

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils"
)

func resourceKafkaCluster() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceKafkaClusterCreate,
		ReadContext:   resourceKafkaClusterRead,
		UpdateContext: resourceKafkaClusterUpdate,
		DeleteContext: resourceKafkaClusterDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceKafkaClusterImport,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:             schema.TypeString,
				Description:      "The name of your Kafka cluster. Must be 63 characters or less and must begin and end with an alphanumeric character (`[a-z0-9A-Z]`) with dashes (`-`), underscores (`_`), dots (`.`), and alphanumerics between.",
				Required:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringIsNotWhiteSpace),
			},
			"version": {
				Type:        schema.TypeString,
				Description: "The desired Kafka Version. Supported version: 3.5.1",
				Required:    true,
			},
			"size": {
				Type:             schema.TypeString,
				Description:      "The size of your Kafka cluster. The size of the Kafka cluster is given in T-shirt sizes. Valid values are: S",
				Required:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringIsNotWhiteSpace),
			},
		},
		Timeouts: &resourceDefaultTimeouts,
	}
}
func resourceKafkaClusterCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(services.SdkBundle).KafkaClient

	createdCluster, _, err := client.CreateCluster(ctx, d)
	if err != nil {
		d.SetId("")
		diags := diag.FromErr(fmt.Errorf("error creating kafka cluster: %w", err))
		return diags
	}

	d.SetId(*createdCluster.Id)
	log.Printf("[INFO] Created kafka cluster: %s", d.Id())

	err = utils.WaitForResourceToBeReady(ctx, d, client.IsClusterAvailable)
	if err != nil {
		diags := diag.FromErr(fmt.Errorf("an error occured  while kafka cluster waiting to be ready: %w", err))
		return diags
	}

	return resourceKafkaClusterRead(ctx, d, meta)
}

func resourceKafkaClusterRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(services.SdkBundle).KafkaClient

	cluster, apiResponse, err := client.GetClusterById(ctx, d.Id())
	if err != nil {
		if apiResponse.HttpNotFound() {
			d.SetId("")
			return nil
		}
		diags := diag.FromErr(fmt.Errorf("error while fetching kafka cluster %s: %w", d.Id(), err))
		return diags
	}

	log.Printf("[INFO] Successfully retreived cluster %s: %+v", d.Id(), cluster)

	if err := client.SetKafkaClusterData(d, &cluster); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceKafkaClusterUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	client := meta.(services.SdkBundle).KafkaClient

	_, apiResponse, err := client.UpdateCluster(ctx, d.Id(), d)
	apiResponse.LogInfo()

	if err != nil {
		if apiResponse.HttpNotFound() {
			d.SetId("")
			return nil
		}
		diags := diag.FromErr(fmt.Errorf("error while updating kafka cluster: %w", err))
		return diags
	}

	err = utils.WaitForResourceToBeReady(ctx, d, client.IsClusterAvailable)
	if err != nil {
		diags := diag.FromErr(fmt.Errorf("an error occured while kafka cluster waiting to be ready: %w", err))
		return diags
	}

	return resourceKafkaClusterRead(ctx, d, meta)
}

func resourceKafkaClusterDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	client := meta.(services.SdkBundle).KafkaClient

	apiResponse, err := client.DeleteCluster(ctx, d.Id())
	if err != nil {
		if apiResponse.HttpNotFound() {
			d.SetId("")
			return nil
		}
		diags := diag.FromErr(fmt.Errorf("error while deleting kafka cluster %s: %w", d.Id(), err))
		return diags
	}

	err = utils.WaitForResourceToBeDeleted(ctx, d, client.IsClusterDeleted)
	if err != nil {
		return diag.FromErr(fmt.Errorf("the check for cluster deletion failed with the following error: %w", err))
	}

	d.SetId("")

	return nil
}

func resourceKafkaClusterImport(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	resourceKafkaClusterRead(ctx, d, meta)
	return []*schema.ResourceData{d}, nil
}
