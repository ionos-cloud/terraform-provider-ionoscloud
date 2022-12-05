package ionoscloud

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	dataplatformService "github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/dataplatform"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils"
	"log"
	"regexp"
)

func resourceDataplatformCluster() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDataplatformClusterCreate,
		ReadContext:   resourceDataplatformClusterRead,
		UpdateContext: resourceDataplatformClusterUpdate,
		DeleteContext: resourceDataplatformClusterDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceDataplatformClusterImport,
		},
		Schema: map[string]*schema.Schema{
			"datacenter_id": {
				Type:        schema.TypeString,
				Description: "The UUID of the virtual data center (VDC) the cluster is provisioned.",
				ValidateDiagFunc: validation.ToDiagFunc(validation.All(validation.StringLenBetween(32, 63),
					validation.StringMatch(regexp.MustCompile("^[0-9a-f]{8}-?[0-9a-f]{4}-?[0-9a-f]{4}-?[0-9a-f]{4}-?[0-9a-f]{12}$"), ""))),
				Required: true,
			},
			"name": {
				Type:             schema.TypeString,
				Description:      "The name of your cluster. Must be 63 characters or less and must be empty or begin and end with an alphanumeric character ([a-z0-9A-Z]) with dashes (-), underscores (_), dots (.), and alphanumerics between.",
				ValidateDiagFunc: validation.ToDiagFunc(validation.All(validation.StringLenBetween(0, 63), validation.StringMatch(regexp.MustCompile("^[A-Za-z0-9][-A-Za-z0-9_.]*[A-Za-z0-9]$"), ""))),
				Required:         true,
			},
			"data_platform_version": {
				Type:             schema.TypeString,
				Description:      "The version of the Data Platform.",
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringLenBetween(0, 32)),
				Optional:         true,
				Computed:         true,
			},
			"maintenance_window": {
				Type:        schema.TypeList,
				Description: "Starting time of a weekly 4 hour-long window, during which maintenance might occur in hh:mm:ss format",
				Optional:    true,
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"time": {
							Type:             schema.TypeString,
							Description:      "Time at which the maintenance should start.",
							ValidateDiagFunc: validation.ToDiagFunc(validation.StringMatch(regexp.MustCompile("^(?:[01]\\d|2[0-3]):(?:[0-5]\\d):(?:[0-5]\\d)$"), "")),
							Required:         true,
						},
						"day_of_the_week": {
							Type:             schema.TypeString,
							ValidateDiagFunc: validation.ToDiagFunc(validation.IsDayOfTheWeek(true)),
							Required:         true,
						},
					},
				},
			},
		},
		CustomizeDiff: checkDataplatformClusterImmutableFields,
		Timeouts:      &resourceDefaultTimeouts,
	}
}

func checkDataplatformClusterImmutableFields(_ context.Context, diff *schema.ResourceDiff, _ interface{}) error {
	if diff.Id() == "" {
		return nil
	}

	if diff.HasChange("datacenter_id") {
		return fmt.Errorf("datacenter_id %s", ImmutableError)
	}

	return nil
}

func resourceDataplatformClusterCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(SdkBundle).DataplatformClient

	id, _, err := client.CreateResource(ctx, d)

	if err != nil {
		diags := diag.FromErr(fmt.Errorf("while creating a Dataplatform Cluster: %w", err))
		return diags
	}

	d.SetId(id)

	err = client.WaitForClusterToBeReady(ctx, id)
	if err != nil {
		return diag.FromErr(fmt.Errorf("waitforCluster %w", err))
	}

	return resourceDataplatformClusterRead(ctx, d, meta)
}

func resourceDataplatformClusterRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	client := meta.(SdkBundle).DataplatformClient

	clusterId := d.Id()
	dataplatformCluster, apiResponse, err := client.GetClusterById(ctx, clusterId)

	if err != nil {
		if apiResponse.HttpNotFound() {
			d.SetId("")
			return nil
		}
		diags := diag.FromErr(fmt.Errorf("error while fetching Dataplatform Cluster %s: %w", d.Id(), err))
		return diags
	}

	log.Printf("[INFO] Successfully retreived Dataplatform Cluster %s: %+v", d.Id(), dataplatformCluster)

	if err := dataplatformService.SetDataplatformClusterData(d, dataplatformCluster); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceDataplatformClusterUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(SdkBundle).DataplatformClient

	clusterId := d.Id()

	_, err := client.UpdateResource(ctx, clusterId, d)

	if err != nil {
		diags := diag.FromErr(fmt.Errorf("an error occured while updating a Dataplatform Cluster: %s", err))
		return diags
	}

	err = client.WaitForClusterToBeReady(ctx, clusterId)
	if err != nil {
		return diag.FromErr(fmt.Errorf("waitforCluster update %w", err))
	}

	return resourceDataplatformClusterRead(ctx, d, meta)
}

func resourceDataplatformClusterDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(SdkBundle).DataplatformClient

	clusterId := d.Id()

	apiResponse, err := client.DeleteResource(ctx, clusterId)

	if err != nil {
		if apiResponse.HttpNotFound() {
			d.SetId("")
			return nil
		}
		diags := diag.FromErr(fmt.Errorf("error while deleting Dataplatform Cluster %s: %s", d.Id(), err))
		return diags
	}

	err = utils.WaitForResourceToBeDeleted(ctx, clusterId, client.DoesResourceExist)
	if err != nil {
		diag.FromErr(fmt.Errorf("while deleting %w", err))
	}

	return nil
}

func resourceDataplatformClusterImport(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	client := meta.(SdkBundle).DataplatformClient

	clusterId := d.Id()

	dataplatformCluster, apiResponse, err := client.GetClusterById(ctx, clusterId)

	if err != nil {
		if apiResponse.HttpNotFound() {
			d.SetId("")
			return nil, fmt.Errorf("dataplatform Cluster does not exist %q", clusterId)
		}
		return nil, fmt.Errorf("an error occured while trying to fetch the import of Dataplatform Cluster %q", clusterId)
	}

	log.Printf("[INFO] Dataplatform Cluster found: %+v", dataplatformCluster)

	if err := dataplatformService.SetDataplatformClusterData(d, dataplatformCluster); err != nil {
		return nil, err
	}

	return []*schema.ResourceData{d}, nil
}
