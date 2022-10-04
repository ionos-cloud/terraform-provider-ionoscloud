package ionoscloud

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	dbaasService "github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/dbaas"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils"
	"log"
	"time"
)

func resourceDbaasPgSqlCluster() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDbaasPgSqlClusterCreate,
		ReadContext:   resourceDbaasPgSqlClusterRead,
		UpdateContext: resourceDbaasPgSqlClusterUpdate,
		DeleteContext: resourceDbaasPgSqlClusterDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceDbaasPgSqlClusterImport,
		},
		CustomizeDiff: checkDBaaSClusterImmutableFields,
		Schema: map[string]*schema.Schema{
			"postgres_version": {
				Type:         schema.TypeString,
				Description:  "The PostgreSQL version of your cluster.",
				Required:     true,
				ValidateFunc: validation.All(validation.StringIsNotWhiteSpace),
			},
			"instances": {
				Type:         schema.TypeInt,
				Description:  "The total number of instances in the cluster (one master and n-1 standbys)",
				Required:     true,
				ValidateFunc: validation.All(validation.IntBetween(1, 5)),
			},
			"cores": {
				Type:        schema.TypeInt,
				Description: "The number of CPU cores per replica.",
				Required:    true,
			},
			"ram": {
				Type:         schema.TypeInt,
				Description:  "The amount of memory per instance in megabytes. Has to be a multiple of 1024.",
				Required:     true,
				ValidateFunc: validation.All(validation.IntAtLeast(2048), validation.IntDivisibleBy(1024)),
			},
			"storage_size": {
				Type:         schema.TypeInt,
				Description:  "The amount of storage per instance in megabytes. Has to be a multiple of 2048.",
				Required:     true,
				ValidateFunc: validation.All(validation.IntAtLeast(2048)),
			},
			"storage_type": {
				Type:         schema.TypeString,
				Description:  "The storage type used in your cluster.",
				Required:     true,
				ValidateFunc: validation.All(validation.StringInSlice([]string{"HDD", "SSD", "SSD Premium", "SSD Standard"}, true)),
			},
			"connections": {
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "Details about the network connection for your cluster.",
				Optional:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"datacenter_id": {
							Type:         schema.TypeString,
							Description:  "The datacenter to connect your cluster to.",
							Required:     true,
							ValidateFunc: validation.All(validation.IsUUID),
						},
						"lan_id": {
							Type:         schema.TypeString,
							Description:  "The LAN to connect your cluster to.",
							Required:     true,
							ValidateFunc: validation.All(validation.StringIsNotWhiteSpace),
						},
						"cidr": {
							Type:         schema.TypeString,
							Description:  "The IP and subnet for the database.\n          Note the following unavailable IP ranges:\n          10.233.64.0/18\n          10.233.0.0/18\n          10.233.114.0/24",
							Required:     true,
							ValidateFunc: VerifyUnavailableIPs,
						},
					},
				},
			},
			"location": {
				Type:         schema.TypeString,
				Description:  "The physical location where the cluster will be created. This will be where all of your instances live. Property cannot be modified after datacenter creation (disallowed in update requests)",
				Required:     true,
				ValidateFunc: validation.All(validation.StringIsNotWhiteSpace),
			},
			"backup_location": {
				Type:         schema.TypeString,
				Description:  "The S3 location where the backups will be stored.",
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.All(validation.StringInSlice([]string{"de", "eu-south-2", "eu-central-2"}, true)),
			},
			"display_name": {
				Type:        schema.TypeString,
				Description: "The friendly name of your cluster.",
				Required:    true,
			},
			"maintenance_window": {
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "a weekly 4 hour-long window, during which maintenance might occur",
				Optional:    true,
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"time": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.All(validation.StringIsNotWhiteSpace),
						},
						"day_of_the_week": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.All(validation.IsDayOfTheWeek(true)),
						},
					},
				},
			},
			"credentials": {
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "Credentials for the database user to be created.",
				Required:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"username": {
							Type:         schema.TypeString,
							Description:  "the username for the initial postgres user. some system usernames are restricted (e.g. \"postgres\", \"admin\", \"standby\")",
							Required:     true,
							ValidateFunc: validation.All(validation.StringIsNotWhiteSpace),
						},
						"password": {
							Type:         schema.TypeString,
							Required:     true,
							Sensitive:    true,
							ValidateFunc: validation.All(validation.StringIsNotWhiteSpace),
						},
					},
				},
			},
			"synchronization_mode": {
				Type:         schema.TypeString,
				Description:  "Represents different modes of replication.",
				Required:     true,
				ValidateFunc: validation.All(validation.StringInSlice([]string{"ASYNCHRONOUS", "SYNCHRONOUS", "STRICTLY_SYNCHRONOUS"}, false)),
			},
			"from_backup": {
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "The PostgreSQL version of your cluster.",
				Optional:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"backup_id": {
							Type:         schema.TypeString,
							Description:  "The unique ID of the backup you want to restore.",
							Required:     true,
							ValidateFunc: validation.All(validation.StringIsNotWhiteSpace),
						},
						"recovery_target_time": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: " If this value is supplied as ISO 8601 timestamp, the backup will be replayed up until the given timestamp. If empty, the backup will be applied completely.",
						},
					},
				},
			},
		},
		Timeouts: &resourceDefaultTimeouts,
	}
}
func checkDBaaSClusterImmutableFields(_ context.Context, diff *schema.ResourceDiff, _ interface{}) error {
	//we do not want to check in case of resource creation
	if diff.Id() == "" {
		return nil
	}
	if diff.HasChange("storage_type") {
		return fmt.Errorf("storage_type %s", ImmutableError)
	}
	if diff.HasChange("location") {
		return fmt.Errorf("location %s", ImmutableError)
	}
	if diff.HasChange("backup_location") {
		return fmt.Errorf("backup_location %s", ImmutableError)
	}
	if diff.HasChange("credentials") {
		return fmt.Errorf("credentials %s", ImmutableError)
	}
	if diff.HasChange("synchronization_mode") {
		return fmt.Errorf("synchronization_mode %s", ImmutableError)
	}
	if diff.HasChange("from_backup") {
		return fmt.Errorf("from_backup %s", ImmutableError)
	}
	return nil

}

func resourceDbaasPgSqlClusterCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(SdkBundle).PsqlClient

	dbaasCluster, err := dbaasService.GetDbaasPgSqlClusterDataCreate(d)

	if err != nil {
		return diag.FromErr(err)
	}
	dbaasClusterResponse, _, err := client.CreateCluster(ctx, *dbaasCluster)

	if err != nil {
		diags := diag.FromErr(fmt.Errorf("an error occured while creating a dbaas cluster: %w", err))
		return diags
	}

	d.SetId(*dbaasClusterResponse.Id)

	for {
		log.Printf("[INFO] Waiting for cluster %s to be ready...", d.Id())

		clusterReady, rsErr := dbaasClusterReady(ctx, client, d)

		if rsErr != nil {
			diags := diag.FromErr(fmt.Errorf("error while checking readiness status of dbaas cluster %s: %w", d.Id(), rsErr))
			return diags
		}

		if clusterReady {
			log.Printf("[INFO] dbaas cluster ready: %s", d.Id())
			break
		}

		select {
		case <-time.After(utils.SleepInterval):
			log.Printf("[INFO] trying again ...")
		case <-ctx.Done():
			log.Printf("[INFO] create timed out")
			diags := diag.FromErr(fmt.Errorf("dbaas cluster creation timed out! WARNING: your dbaas cluster (%s) will still probably be created after some time but the terraform state wont reflect that; check your Ionos Cloud account for updates", d.Id()))
			return diags
		}

	}

	return resourceDbaasPgSqlClusterRead(ctx, d, meta)
}

func resourceDbaasPgSqlClusterRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	client := meta.(SdkBundle).PsqlClient

	cluster, apiResponse, err := client.GetCluster(ctx, d.Id())

	if err != nil {
		if apiResponse != nil && apiResponse.Response != nil && apiResponse.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		diags := diag.FromErr(fmt.Errorf("error while fetching dbaas cluster %s: %w", d.Id(), err))
		return diags
	}

	log.Printf("[INFO] Successfully retreived cluster %s: %+v", d.Id(), cluster)

	if err := dbaasService.SetDbaasPgSqlClusterData(d, cluster); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceDbaasPgSqlClusterUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(SdkBundle).PsqlClient

	cluster, diags := dbaasService.GetDbaasPgSqlClusterDataUpdate(d)

	if diags != nil {
		return diags
	}

	dbaasClusterResponse, _, err := client.ClustersApi.ClustersPatch(ctx, d.Id()).PatchClusterRequest(*cluster).Execute()

	if err != nil {
		diags := diag.FromErr(fmt.Errorf("an error occured while updating a dbaas cluster: %s", err))
		return diags
	}

	d.SetId(*dbaasClusterResponse.Id)

	time.Sleep(utils.SleepInterval)

	for {
		log.Printf("[INFO] Waiting for cluster %s to be ready...", d.Id())

		clusterReady, rsErr := dbaasClusterReady(ctx, client, d)

		if rsErr != nil {
			diags := diag.FromErr(fmt.Errorf("error while checking readiness status of dbaas cluster %s: %w", d.Id(), rsErr))
			return diags
		}

		if clusterReady {
			log.Printf("[INFO] dbaas cluster ready: %s", d.Id())
			break
		}

		select {
		case <-time.After(utils.SleepInterval):
			log.Printf("[INFO] trying again ...")
		case <-ctx.Done():
			log.Printf("[INFO] create timed out")
			diags := diag.FromErr(fmt.Errorf("dbaas cluster update timed out! WARNING: your dbaas cluster (%s) will still probably be updated after some time but the terraform state wont reflect that; check your Ionos Cloud account for updates", d.Id()))
			return diags
		}

	}

	return resourceDbaasPgSqlClusterRead(ctx, d, meta)
}

func resourceDbaasPgSqlClusterDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(SdkBundle).PsqlClient

	_, apiResponse, err := client.DeleteCluster(ctx, d.Id())

	if err != nil {
		if apiResponse != nil && apiResponse.Response != nil && apiResponse.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		diags := diag.FromErr(fmt.Errorf("error while deleting dbaas cluster %s: %s", d.Id(), err))
		return diags
	}

	for {
		log.Printf("[INFO] Waiting for cluster %s to be deleted...", d.Id())

		clusterdDeleted, dsErr := dbaasClusterDeleted(ctx, client, d)

		if dsErr != nil {
			diags := diag.FromErr(fmt.Errorf("error while checking deletion status of dbaas cluster %s: %s", d.Id(), dsErr))
			return diags
		}

		if clusterdDeleted {
			log.Printf("[INFO] Successfully deleted dbaas cluster: %s", d.Id())
			break
		}

		select {
		case <-time.After(utils.SleepInterval):
			log.Printf("[INFO] trying again ...")
		case <-ctx.Done():
			diags := diag.FromErr(fmt.Errorf("dbaas cluster deletion timed out! WARNING: your k8s cluster (%s) will still probably be deleted after some time but the terraform state won't reflect that; check your Ionos Cloud account for updates", d.Id()))
			return diags
		}
	}

	// wait 15 seconds after the deletion of the cluster, for the lan to be freed
	time.Sleep(utils.SleepInterval * 3)

	return nil
}

func resourceDbaasPgSqlClusterImport(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	client := meta.(SdkBundle).PsqlClient

	clusterId := d.Id()

	dbaasCluster, apiResponse, err := client.GetCluster(ctx, clusterId)

	if err != nil {
		if apiResponse != nil && apiResponse.Response != nil && apiResponse.StatusCode == 404 {
			d.SetId("")
			return nil, fmt.Errorf("dbaas cluster does not exist %q", clusterId)
		}
		return nil, fmt.Errorf("an error occured while trying to fetch the import of dbaas cluster %q", clusterId)
	}

	log.Printf("[INFO] dbaas cluster found: %+v", dbaasCluster)

	if err := dbaasService.SetDbaasPgSqlClusterData(d, dbaasCluster); err != nil {
		return nil, err
	}

	return []*schema.ResourceData{d}, nil
}

func dbaasClusterReady(ctx context.Context, client *dbaasService.PsqlClient, d *schema.ResourceData) (bool, error) {
	subjectCluster, _, err := client.GetCluster(ctx, d.Id())

	if err != nil {
		return true, fmt.Errorf("error checking dbaas cluster status: %s", err)
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
	return *subjectCluster.Metadata.State == utils.Available, nil
}

func dbaasClusterDeleted(ctx context.Context, client *dbaasService.PsqlClient, d *schema.ResourceData) (bool, error) {

	_, apiResponse, err := client.GetCluster(ctx, d.Id())

	if err != nil {
		if apiResponse != nil && apiResponse.Response != nil && apiResponse.StatusCode == 404 {
			return true, nil
		}
		return true, fmt.Errorf("error checking dbaas cluster deletion status: %s", err)
	}
	return false, nil
}
