package ionoscloud

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	dataplatformService "github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/dataplatform"
	"log"
	"regexp"
	"time"
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
				Type:         schema.TypeString,
				Description:  "The name of your node pool. Must be 63 characters or less and must be empty or begin and end with an alphanumeric character ([a-z0-9A-Z]) with dashes (-), underscores (_), dots (.), and alphanumerics between.",
				Required:     true,
				ValidateFunc: validation.All(validation.StringLenBetween(0, 63), validation.StringMatch(regexp.MustCompile("^[A-Za-z0-9][-A-Za-z0-9_.]*[A-Za-z0-9]$"), "")),
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
			},
			"cores_count": {
				Type:         schema.TypeInt,
				Description:  "The number of CPU cores per node.",
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.All(validation.IntAtLeast(1)),
			},
			"ram_size": {
				Type:         schema.TypeInt,
				Description:  "The RAM size for one node in MB. Must be set in multiples of 1024 MB, with a minimum size is of 2048 MB.",
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.All(validation.IntAtLeast(2048), validation.IntDivisibleBy(1024)),
			},
			"availability_zone": {
				Type:         schema.TypeString,
				Description:  "The availability zone of the virtual datacenter region where the node pool resources should be provisioned.",
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.All(validation.StringInSlice([]string{"AUTO", "ZONE_1", "ZONE_2"}, true)),
			},
			"storage_type": {
				Type:         schema.TypeString,
				Description:  "The type of hardware for the volume.",
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.All(validation.StringInSlice([]string{"HDD", "SSD"}, true)),
			},
			"storage_size": {
				Type:         schema.TypeInt,
				Description:  "The size of the volume in GB. The size must be greater than 10GB.",
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.All(validation.IntAtLeast(10)),
			},
			"maintenance_window": {
				Type:        schema.TypeList,
				Description: "Starting time of a weekly 4 hour-long window, during which maintenance might occur in hh:mm:ss format",
				Optional:    true,
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"time": {
							Type:         schema.TypeString,
							Description:  "Time at which the maintenance should start.",
							ValidateFunc: validation.All(validation.StringMatch(regexp.MustCompile("^(?:[01]\\d|2[0-3]):(?:[0-5]\\d):(?:[0-5]\\d)$"), "")),
							Required:     true,
						},
						"day_of_the_week": {
							Type:         schema.TypeString,
							ValidateFunc: validation.All(validation.IsDayOfTheWeek(true)),
							Required:     true,
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
			"data_platform_version": {
				Type:        schema.TypeString,
				Description: "The version of the DataPlatform.",
				Computed:    true,
			},
			"datacenter_id": {
				Type:        schema.TypeString,
				Description: "The UUID of the virtual data center (VDC) the cluster is provisioned.",
				Computed:    true,
			},
			"cluster_id": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "The UUID of an existing Dataplatform cluster.",
				ValidateFunc: validation.All(validation.StringMatch(regexp.MustCompile("^[A-Za-z0-9][-A-Za-z0-9_.]*[A-Za-z0-9]$"), "")),
			},
		},
		Timeouts:      &resourceDefaultTimeouts,
		CustomizeDiff: checkDataplatformNodePoolImmutableFields,
	}
}

func checkDataplatformNodePoolImmutableFields(_ context.Context, diff *schema.ResourceDiff, _ interface{}) error {
	if diff.Id() == "" {
		return nil
	}

	if diff.HasChange("name") {
		return fmt.Errorf("name %s", ImmutableError)
	}

	if diff.HasChange("cpu_family") {
		return fmt.Errorf("cpu_family %s", ImmutableError)
	}

	if diff.HasChange("cores_count") {
		return fmt.Errorf("cores_count %s", ImmutableError)
	}

	if diff.HasChange("ram_size") {
		return fmt.Errorf("ram_size %s", ImmutableError)
	}

	if diff.HasChange("availability_zone") {
		return fmt.Errorf("availability_zone %s", ImmutableError)
	}

	if diff.HasChange("storage_size") {
		return fmt.Errorf("storage_size %s", ImmutableError)
	}

	if diff.HasChange("storage_type") {
		return fmt.Errorf("storage_type %s", ImmutableError)
	}

	return nil
}

func resourceDataplatformNodePoolCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(SdkBundle).DataplatformClient

	clusterId := d.Get("cluster_id").(string)

	dataplatformNodePool := dataplatformService.GetDataplatformNodePoolDataCreate(d)
	dataplatformNodePoolResponse, _, err := client.CreateNodePool(ctx, clusterId, *dataplatformNodePool)

	if err != nil {
		diags := diag.FromErr(fmt.Errorf("an error occured while creating a Dataplatform NodePool: %w", err))
		return diags
	}

	d.SetId(*dataplatformNodePoolResponse.Id)

	for {
		log.Printf("[INFO] Waiting for Dataplatform Node Pool %s to be ready...", d.Id())

		nodePoolReady, rsErr := dataplatformNodePoolReady(ctx, client, d)

		if rsErr != nil {
			diags := diag.FromErr(fmt.Errorf("error while checking readiness status of Dataplatform Node Pool %s: %w", d.Id(), rsErr))
			return diags
		}

		if nodePoolReady {
			log.Printf("[INFO] Dataplatform Node Pool ready: %s", d.Id())
			break
		}

		select {
		case <-time.After(SleepInterval):
			log.Printf("[INFO] trying again ...")
		case <-ctx.Done():
			log.Printf("[INFO] create timed out")
			diags := diag.FromErr(fmt.Errorf("Dataplatform Node Pool creation timed out! WARNING: your Dataplatform Node Pool (%s) will still probably be created after some time but the terraform state wont reflect that; check your Ionos Cloud account for updates", d.Id()))
			return diags
		}

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

	dataplatformNodePool, diags := dataplatformService.GetDataplatformNodePoolDataUpdate(d)

	if diags != nil {
		return diags
	}

	dataplatformNodePoolResponse, _, err := client.UpdateNodePool(ctx, clusterId, nodePoolId, *dataplatformNodePool)

	if err != nil {
		diags := diag.FromErr(fmt.Errorf("an error occured while updating a Dataplatform NodePool: %s", err))
		return diags
	}

	d.SetId(*dataplatformNodePoolResponse.Id)

	time.Sleep(SleepInterval)

	for {
		log.Printf("[INFO] Waiting for Node Pool %s to be ready...", d.Id())

		nodePoolReady, rsErr := dataplatformNodePoolReady(ctx, client, d)

		if rsErr != nil {
			diags := diag.FromErr(fmt.Errorf("error while checking readiness status of Dataplatform Node Pool %s: %w", d.Id(), rsErr))
			return diags
		}

		if nodePoolReady {
			log.Printf("[INFO] Dataplatform Node Pool ready: %s", d.Id())
			break
		}

		select {
		case <-time.After(SleepInterval):
			log.Printf("[INFO] trying again ...")
		case <-ctx.Done():
			log.Printf("[INFO] create timed out")
			diags := diag.FromErr(fmt.Errorf("Dataplatform Node Pool update timed out! WARNING: your Dataplatform Node Pool (%s) will still probably be updated after some time but the terraform state wont reflect that; check your Ionos Cloud account for updates", d.Id()))
			return diags
		}

	}

	return resourceDataplatformNodePoolRead(ctx, d, meta)
}

func resourceDataplatformNodePoolDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(SdkBundle).DataplatformClient

	clusterId := d.Get("cluster_id").(string)
	nodePoolId := d.Id()

	_, apiResponse, err := client.DeleteNodePool(ctx, clusterId, nodePoolId)

	if err != nil {
		if apiResponse != nil && apiResponse.Response != nil && apiResponse.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		diags := diag.FromErr(fmt.Errorf("error while deleting Dataplatform Node Pool %s: %s", d.Id(), err))
		return diags
	}

	for {
		log.Printf("[INFO] Waiting for Dataplatform Node Pool %s to be deleted...", d.Id())

		nodePoolDeleted, dsErr := dataplatformNodePoolDeleted(ctx, client, d)

		if dsErr != nil {
			diags := diag.FromErr(fmt.Errorf("error while checking deletion status of Dataplatform Node Pool %s: %s", d.Id(), dsErr))
			return diags
		}

		if nodePoolDeleted {
			log.Printf("[INFO] Successfully deleted Dataplatform Node Pool: %s", d.Id())
			break
		}

		select {
		case <-time.After(SleepInterval):
			log.Printf("[INFO] trying again ...")
		case <-ctx.Done():
			diags := diag.FromErr(fmt.Errorf("Dataplatform Node Pool deletion timed out! WARNING: your Dataplatform Node Pool (%s) will still probably be deleted after some time but the terraform state won't reflect that; check your Ionos Cloud account for updates", d.Id()))
			return diags
		}
	}

	// wait 15 seconds after the deletion of the Node Pool, for the lan to be freed
	time.Sleep(SleepInterval * 3)

	return nil
}

func resourceDataplatformNodePoolImport(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	client := meta.(SdkBundle).DataplatformClient

	clusterId := d.Get("cluster_id").(string)
	nodePoolId := d.Id()

	dataplatformNodePool, apiResponse, err := client.GetNodePool(ctx, clusterId, nodePoolId)

	if err != nil {
		if apiResponse != nil && apiResponse.Response != nil && apiResponse.StatusCode == 404 {
			d.SetId("")
			return nil, fmt.Errorf("Dataplatform Node Pool does not exist %q", nodePoolId)
		}
		return nil, fmt.Errorf("an error occured while trying to fetch the import of Dataplatform Node Pool %q", nodePoolId)
	}

	log.Printf("[INFO] Dataplatform Node Pool found: %+v", dataplatformNodePool)

	if err := dataplatformService.SetDataplatformNodePoolData(d, dataplatformNodePool); err != nil {
		return nil, err
	}

	return []*schema.ResourceData{d}, nil
}

func dataplatformNodePoolReady(ctx context.Context, client *dataplatformService.Client, d *schema.ResourceData) (bool, error) {

	clusterId := d.Get("cluster_id").(string)
	nodePoolId := d.Id()

	subjectNodePool, _, err := client.GetNodePool(ctx, clusterId, nodePoolId)

	if err != nil {
		return true, fmt.Errorf("error checking Dataplatform Node Pool status: %s", err)
	}
	// ToDo: Removed this part since there are still problems with the nodePools being unstable (failing for a short time and then recovering)
	//if *subjectNodePool.LifecycleStatus == "FAILED" {
	//
	//	time.Sleep(time.Second * 3)
	//
	//	subjectNodePool, _, err = client.GetNodePool(ctx, d.Id())
	//
	//	if err != nil {
	//		return true, fmt.Errorf("error checking dbaas nodePool status: %s", err)
	//	}
	//
	//	if *subjectNodePool.LifecycleStatus == "FAILED" {
	//		return false, fmt.Errorf("dbaas nodePool has failed. WARNING: your k8s nodePool may still recover after some time but the terraform state won't reflect that; check your Ionos Cloud account for updates")
	//	}
	//}
	return *subjectNodePool.Metadata.State == "AVAILABLE", nil
}

func dataplatformNodePoolDeleted(ctx context.Context, client *dataplatformService.Client, d *schema.ResourceData) (bool, error) {

	clusterId := d.Get("cluster_id").(string)
	nodePoolId := d.Id()

	_, apiResponse, err := client.GetNodePool(ctx, clusterId, nodePoolId)

	if err != nil {
		if apiResponse != nil && apiResponse.Response != nil && apiResponse.StatusCode == 404 {
			return true, nil
		}
		return true, fmt.Errorf("error checking Dataplatform Node Pool deletion status: %s", err)
	}
	return false, nil
}
