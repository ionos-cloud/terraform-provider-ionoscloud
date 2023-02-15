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

func resourceDataplatformNodePool() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDataplatformNodePoolCreate,
		ReadContext:   resourceDataplatformNodePoolRead,
		UpdateContext: resourceDataplatformNodePoolUpdate,
		DeleteContext: resourceDataplatformNodePoolDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceDataplatformNodePoolImport,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:             schema.TypeString,
				Description:      "The name of your node pool. Must be 63 characters or less and must be empty or begin and end with an alphanumeric character ([a-z0-9A-Z]) with dashes (-), underscores (_), dots (.), and alphanumerics between.",
				Required:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.All(validation.StringLenBetween(0, 63), validation.StringMatch(regexp.MustCompile("^[A-Za-z0-9][-A-Za-z0-9_.]*[A-Za-z0-9]$"), ""))),
				ForceNew:         true,
			},
			"node_count": {
				Type:        schema.TypeInt,
				Description: "The number of nodes that make up the node pool.",
				Required:    true,
			},
			"cpu_family": {
				Type:        schema.TypeString,
				Description: "A valid CPU family name or `AUTO` if the platform shall choose the best fitting option. Available CPU architectures can be retrieved from the datacenter resource.",
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
			},
			"cores_count": {
				Type:             schema.TypeInt,
				Description:      "The number of CPU cores per node.",
				Optional:         true,
				Computed:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.IntAtLeast(1)),
				ForceNew:         true,
			},
			"ram_size": {
				Type:             schema.TypeInt,
				Description:      "The RAM size for one node in MB. Must be set in multiples of 1024 MB, with a minimum size is of 2048 MB.",
				Optional:         true,
				Computed:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.All(validation.IntAtLeast(2048), validation.IntDivisibleBy(1024))),
				ForceNew:         true,
			},
			"availability_zone": {
				Type:             schema.TypeString,
				Description:      "The availability zone of the virtual datacenter region where the node pool resources should be provisioned.",
				Optional:         true,
				Computed:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"AUTO", "ZONE_1", "ZONE_2"}, true)),
				ForceNew:         true,
			},
			"storage_type": {
				Type:             schema.TypeString,
				Description:      "The type of hardware for the volume.",
				Optional:         true,
				Computed:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"HDD", "SSD"}, true)),
				ForceNew:         true,
			},
			"storage_size": {
				Type:             schema.TypeInt,
				Description:      "The size of the volume in GB. The size must be greater than 10GB.",
				Optional:         true,
				Computed:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.IntAtLeast(10)),
				ForceNew:         true,
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
			"labels": {
				Type:        schema.TypeMap,
				Description: "Key-value pairs attached to the node pool resource as [Kubernetes labels](https://kubernetes.io/docs/concepts/overview/working-with-objects/labels/)",
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"annotations": {
				Type:        schema.TypeMap,
				Description: "Key-value pairs attached to node pool resource as [Kubernetes annotations](https://kubernetes.io/docs/concepts/overview/working-with-objects/annotations/)",
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"version": {
				Type:        schema.TypeString,
				Description: "The version of the Data Platform.",
				Computed:    true,
			},
			"datacenter_id": {
				Type:        schema.TypeString,
				Description: "The UUID of the virtual data center (VDC) the cluster is provisioned.",
				Computed:    true,
			},
			"cluster_id": {
				Type:             schema.TypeString,
				Required:         true,
				Description:      "The UUID of an existing Dataplatform cluster.",
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringMatch(regexp.MustCompile("^[A-Za-z0-9][-A-Za-z0-9_.]*[A-Za-z0-9]$"), "")),
			},
		},
		Timeouts: &resourceDefaultTimeouts,
	}
}

func resourceDataplatformNodePoolCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(SdkBundle).DataplatformClient

	clusterId := d.Get("cluster_id").(string)

	dataplatformNodePoolResponse, _, err := client.CreateNodePool(ctx, clusterId, d)
	if err != nil {
		diags := diag.FromErr(fmt.Errorf("an error occured while creating a Dataplatform NodePool: %w", err))
		return diags
	}

	d.SetId(*dataplatformNodePoolResponse.Id)
	err = utils.WaitForResourceToBeReady(ctx, d, client.IsNodePoolReady)
	if err != nil {
		diags := diag.FromErr(fmt.Errorf(" while dataplaform nodepool waiting to be ready: %w", err))
		return diags
	}

	return resourceDataplatformNodePoolRead(ctx, d, meta)
}

func resourceDataplatformNodePoolRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(SdkBundle).DataplatformClient

	clusterId := d.Get("cluster_id").(string)
	nodePoolId := d.Id()

	dataplatformNodePool, apiResponse, err := client.GetNodePool(ctx, clusterId, nodePoolId)

	if err != nil {
		if apiResponse != nil && apiResponse.Response != nil && apiResponse.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		diags := diag.FromErr(fmt.Errorf("error while fetching Dataplatform Node Pool %s: %w", d.Id(), err))
		return diags
	}

	log.Printf("[INFO] Successfully retreived Dataplatform Node Pool %s: %+v", d.Id(), dataplatformNodePool)

	if err := dataplatformService.SetDataplatformNodePoolData(d, dataplatformNodePool); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceDataplatformNodePoolUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(SdkBundle).DataplatformClient

	clusterId := d.Get("cluster_id").(string)
	nodePoolId := d.Id()

	dataplatformNodePoolResponse, _, err := client.UpdateNodePool(ctx, clusterId, nodePoolId, d)
	if err != nil {
		diags := diag.FromErr(fmt.Errorf("an error occured while updating a Dataplatform NodePool: %w", err))
		return diags
	}

	d.SetId(*dataplatformNodePoolResponse.Id)

	err = utils.WaitForResourceToBeReady(ctx, d, client.IsNodePoolReady)
	if err != nil {
		diag.FromErr(fmt.Errorf("waiting until ready %w", err))
	}
	return resourceDataplatformNodePoolRead(ctx, d, meta)
}

func resourceDataplatformNodePoolDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(SdkBundle).DataplatformClient

	clusterId := d.Get("cluster_id").(string)
	nodePoolId := d.Id()

	_, apiResponse, err := client.DeleteNodePool(ctx, clusterId, nodePoolId)
	if err != nil {
		if apiResponse.HttpNotFound() {
			d.SetId("")
			return nil
		}
		diags := diag.FromErr(fmt.Errorf("error while deleting Dataplatform Node Pool %s: %w", d.Id(), err))
		return diags
	}
	err = utils.WaitForResourceToBeDeleted(ctx, d, client.IsNodePoolDeleted)
	if err != nil {
		diag.FromErr(fmt.Errorf("deleting %w", err))
	}

	return nil
}

func resourceDataplatformNodePoolImport(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	client := meta.(SdkBundle).DataplatformClient

	clusterId := d.Get("cluster_id").(string)
	nodePoolId := d.Id()

	dataplatformNodePool, apiResponse, err := client.GetNodePool(ctx, clusterId, nodePoolId)

	if err != nil {
		if apiResponse.HttpNotFound() {
			d.SetId("")
			return nil, fmt.Errorf("dataplatform Node Pool does not exist %q", nodePoolId)
		}
		return nil, fmt.Errorf("an error occured while trying to fetch the import of Dataplatform Node Pool %q", nodePoolId)
	}

	log.Printf("[INFO] Dataplatform Node Pool found: %+v", dataplatformNodePool)

	if err := dataplatformService.SetDataplatformNodePoolData(d, dataplatformNodePool); err != nil {
		return nil, err
	}

	return []*schema.ResourceData{d}, nil
}
