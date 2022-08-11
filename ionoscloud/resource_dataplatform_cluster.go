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
				Type:         schema.TypeString,
				Description:  "The UUID of the virtual data center (VDC) the cluster is provisioned.",
				ValidateFunc: validation.All(validation.StringLenBetween(32, 63), validation.StringMatch(regexp.MustCompile("^[0-9a-f]{8}-?[0-9a-f]{4}-?[0-9a-f]{4}-?[0-9a-f]{4}-?[0-9a-f]{12}$"), "")),
				Required:     true,
			},
			"name": {
				Type:         schema.TypeString,
				Description:  "The name of your cluster. Must be 63 characters or less and must be empty or begin and end with an alphanumeric character ([a-z0-9A-Z]) with dashes (-), underscores (_), dots (.), and alphanumerics between.",
				ValidateFunc: validation.All(validation.StringLenBetween(0, 63), validation.StringMatch(regexp.MustCompile("^[A-Za-z0-9][-A-Za-z0-9_.]*[A-Za-z0-9]$"), "")),
				Required:     true,
			},
			"data_platform_version": {
				Type:         schema.TypeString,
				Description:  "The version of the DataPlatform.",
				ValidateFunc: validation.All(validation.StringLenBetween(0, 32)),
				Optional:     true,
				Computed:     true,
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

	dataplatformCluster := dataplatformService.GetDataplatformClusterDataCreate(d)
	dataplatformClusterResponse, _, err := client.CreateCluster(ctx, *dataplatformCluster)

	if err != nil {
		diags := diag.FromErr(fmt.Errorf("an error occured while creating a Dataplatform Cluster: %w", err))
		return diags
	}

	d.SetId(*dataplatformClusterResponse.Id)

	for {
		log.Printf("[INFO] Waiting for Dataplatform Cluster %s to be ready...", d.Id())

		clusterReady, rsErr := dataplatformClusterReady(ctx, client, d)

		if rsErr != nil {
			diags := diag.FromErr(fmt.Errorf("error while checking readiness status of Dataplatform Cluster %s: %w", d.Id(), rsErr))
			return diags
		}

		if clusterReady {
			log.Printf("[INFO] Dataplatform Cluster ready: %s", d.Id())
			break
		}

		select {
		case <-time.After(SleepInterval):
			log.Printf("[INFO] trying again ...")
		case <-ctx.Done():
			log.Printf("[INFO] create timed out")
			diags := diag.FromErr(fmt.Errorf("Dataplatform Cluster creation timed out! WARNING: your Dataplatform Cluster (%s) will still probably be created after some time but the terraform state wont reflect that; check your Ionos Cloud account for updates", d.Id()))
			return diags
		}

	}

	return resourceDataplatformClusterRead(ctx, d, meta)
}

func resourceDataplatformClusterRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	client := meta.(SdkBundle).DataplatformClient

	clusterId := d.Id()
	dataplatformCluster, apiResponse, err := client.GetCluster(ctx, clusterId)

	if err != nil {
		if apiResponse != nil && apiResponse.Response != nil && apiResponse.StatusCode == 404 {
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

	dataplatformCluster, diags := dataplatformService.GetDataplatformClusterDataUpdate(d)

	if diags != nil {
		return diags
	}

	dataplatformClusterResponse, _, err := client.UpdateCluster(ctx, clusterId, *dataplatformCluster)

	if err != nil {
		diags := diag.FromErr(fmt.Errorf("an error occured while updating a Dataplatform Cluster: %s", err))
		return diags
	}

	d.SetId(*dataplatformClusterResponse.Id)

	time.Sleep(SleepInterval)

	for {
		log.Printf("[INFO] Waiting for Cluster %s to be ready...", d.Id())

		clusterReady, rsErr := dataplatformClusterReady(ctx, client, d)

		if rsErr != nil {
			diags := diag.FromErr(fmt.Errorf("error while checking readiness status of Dataplatform Cluster %s: %w", d.Id(), rsErr))
			return diags
		}

		if clusterReady {
			log.Printf("[INFO] Dataplatform Cluster ready: %s", d.Id())
			break
		}

		select {
		case <-time.After(SleepInterval):
			log.Printf("[INFO] trying again ...")
		case <-ctx.Done():
			log.Printf("[INFO] create timed out")
			diags := diag.FromErr(fmt.Errorf("Dataplatform Cluster update timed out! WARNING: your Dataplatform Cluster (%s) will still probably be updated after some time but the terraform state wont reflect that; check your Ionos Cloud account for updates", d.Id()))
			return diags
		}

	}

	return resourceDataplatformClusterRead(ctx, d, meta)
}

func resourceDataplatformClusterDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(SdkBundle).DataplatformClient

	clusterId := d.Id()

	_, apiResponse, err := client.DeleteCluster(ctx, clusterId)

	if err != nil {
		if apiResponse != nil && apiResponse.Response != nil && apiResponse.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		diags := diag.FromErr(fmt.Errorf("error while deleting Dataplatform Cluster %s: %s", d.Id(), err))
		return diags
	}

	for {
		log.Printf("[INFO] Waiting for Dataplatform Cluster %s to be deleted...", d.Id())

		clusterdDeleted, dsErr := dataplatformClusterDeleted(ctx, client, d)

		if dsErr != nil {
			diags := diag.FromErr(fmt.Errorf("error while checking deletion status of Dataplatform Cluster %s: %s", d.Id(), dsErr))
			return diags
		}

		if clusterdDeleted {
			log.Printf("[INFO] Successfully deleted Dataplatform Cluster: %s", d.Id())
			break
		}

		select {
		case <-time.After(SleepInterval):
			log.Printf("[INFO] trying again ...")
		case <-ctx.Done():
			diags := diag.FromErr(fmt.Errorf("Dataplatform Cluster deletion timed out! WARNING: your Dataplatform Cluster (%s) will still probably be deleted after some time but the terraform state won't reflect that; check your Ionos Cloud account for updates", d.Id()))
			return diags
		}
	}

	// wait 15 seconds after the deletion of the Cluster, for the lan to be freed
	time.Sleep(SleepInterval * 3)

	return nil
}

func resourceDataplatformClusterImport(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	client := meta.(SdkBundle).DataplatformClient

	clusterId := d.Id()

	dataplatformCluster, apiResponse, err := client.GetCluster(ctx, clusterId)

	if err != nil {
		if apiResponse != nil && apiResponse.Response != nil && apiResponse.StatusCode == 404 {
			d.SetId("")
			return nil, fmt.Errorf("Dataplatform Cluster does not exist %q", clusterId)
		}
		return nil, fmt.Errorf("an error occured while trying to fetch the import of Dataplatform Cluster %q", clusterId)
	}

	log.Printf("[INFO] Dataplatform Cluster found: %+v", dataplatformCluster)

	if err := dataplatformService.SetDataplatformClusterData(d, dataplatformCluster); err != nil {
		return nil, err
	}

	return []*schema.ResourceData{d}, nil
}

func dataplatformClusterReady(ctx context.Context, client *dataplatformService.Client, d *schema.ResourceData) (bool, error) {
	clusterId := d.Id()

	subjectCluster, _, err := client.GetCluster(ctx, clusterId)

	if err != nil {
		return true, fmt.Errorf("error checking Dataplatform Cluster status: %s", err)
	}
	// ToDo: Removed this part since there are still problems with the clusters being unstable (failing for a short time and then recovering)
	//if *subjectCluster.LifecycleStatus == "FAILED" {
	//
	//	time.Sleep(time.Second * 3)
	//
	//	subjectCluster, _, err = client.GetCluster(ctx, d.Id())
	//
	//	if err != nil {
	//		return true, fmt.Errorf("error checking dbaas cluster status: %s", err)
	//	}
	//
	//	if *subjectCluster.LifecycleStatus == "FAILED" {
	//		return false, fmt.Errorf("dbaas cluster has failed. WARNING: your k8s cluster may still recover after some time but the terraform state won't reflect that; check your Ionos Cloud account for updates")
	//	}
	//}
	return *subjectCluster.Metadata.State == "AVAILABLE", nil
}

func dataplatformClusterDeleted(ctx context.Context, client *dataplatformService.Client, d *schema.ResourceData) (bool, error) {
	clusterId := d.Id()

	_, apiResponse, err := client.GetCluster(ctx, clusterId)

	if err != nil {
		if apiResponse != nil && apiResponse.Response != nil && apiResponse.StatusCode == 404 {
			return true, nil
		}
		return true, fmt.Errorf("error checking Dataplatform Cluster deletion status: %s", err)
	}
	return false, nil
}
