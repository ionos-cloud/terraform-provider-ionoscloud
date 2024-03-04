package ionoscloud

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils/constant"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	dbaasService "github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/dbaas"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils"
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
				Type:             schema.TypeString,
				Description:      "The PostgreSQL version of your cluster.",
				Required:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringIsNotWhiteSpace),
			},
			"instances": {
				Type:             schema.TypeInt,
				Description:      "The total number of instances in the cluster (one master and n-1 standbys)",
				Required:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.IntBetween(1, 5)),
			},
			"cores": {
				Type:        schema.TypeInt,
				Description: "The number of CPU cores per replica.",
				Required:    true,
			},
			"ram": {
				Type:             schema.TypeInt,
				Description:      "The amount of memory per instance in megabytes. Has to be a multiple of 1024.",
				Required:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.All(validation.IntAtLeast(2048), validation.IntDivisibleBy(1024))),
			},
			"storage_size": {
				Type:             schema.TypeInt,
				Description:      "The amount of storage per instance in megabytes. Has to be a multiple of 2048.",
				Required:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.IntAtLeast(2048)),
			},
			"storage_type": {
				Type:             schema.TypeString,
				Description:      "The storage type used in your cluster.",
				Required:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"HDD", "SSD", "SSD Premium", "SSD Standard"}, true)),
			},
			"connections": {
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "Details about the network connection for your cluster.",
				Optional:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"datacenter_id": {
							Type:             schema.TypeString,
							Description:      "The datacenter to connect your cluster to.",
							Required:         true,
							ValidateDiagFunc: validation.ToDiagFunc(validation.IsUUID),
						},
						"lan_id": {
							Type:             schema.TypeString,
							Description:      "The LAN to connect your cluster to.",
							Required:         true,
							ValidateDiagFunc: validation.ToDiagFunc(validation.StringIsNotWhiteSpace),
						},
						"cidr": {
							Type:             schema.TypeString,
							Description:      "The IP and subnet for the database.\n          Note the following unavailable IP ranges:\n          10.233.64.0/18\n          10.233.0.0/18\n          10.233.114.0/24",
							Required:         true,
							ValidateDiagFunc: validation.ToDiagFunc(VerifyUnavailableIPs),
						},
					},
				},
			},
			"location": {
				Type:             schema.TypeString,
				Description:      "The physical location where the cluster will be created. This will be where all of your instances live. Property cannot be modified after datacenter creation (disallowed in update requests)",
				Required:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringIsNotWhiteSpace),
			},
			"backup_location": {
				Type:             schema.TypeString,
				Description:      "The S3 location where the backups will be stored.",
				Optional:         true,
				Computed:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"de", "eu-south-2", "eu-central-2"}, true)),
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
							Type:             schema.TypeString,
							Required:         true,
							ValidateDiagFunc: validation.ToDiagFunc(validation.StringIsNotWhiteSpace),
						},
						"day_of_the_week": {
							Type:             schema.TypeString,
							Required:         true,
							ValidateDiagFunc: validation.ToDiagFunc(validation.IsDayOfTheWeek(true)),
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
							Type:             schema.TypeString,
							Description:      "the username for the initial postgres user. some system usernames are restricted (e.g. \"postgres\", \"admin\", \"standby\")",
							Required:         true,
							ValidateDiagFunc: validation.ToDiagFunc(validation.StringIsNotWhiteSpace),
						},
						"password": {
							Type:             schema.TypeString,
							Required:         true,
							Sensitive:        true,
							ValidateDiagFunc: validation.ToDiagFunc(validation.StringIsNotWhiteSpace),
						},
					},
				},
			},
			"synchronization_mode": {
				Type:             schema.TypeString,
				Description:      "Represents different modes of replication.",
				Required:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"ASYNCHRONOUS", "SYNCHRONOUS", "STRICTLY_SYNCHRONOUS"}, false)),
			},
			"from_backup": {
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "Creates the cluster based on the existing backup.",
				Optional:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"backup_id": {
							Type:             schema.TypeString,
							Description:      "The unique ID of the backup you want to restore.",
							Required:         true,
							ValidateDiagFunc: validation.ToDiagFunc(validation.StringIsNotWhiteSpace),
						},
						"recovery_target_time": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "If this value is supplied as ISO 8601 timestamp, the backup will be replayed up until the given timestamp. If empty, the backup will be applied completely.",
						},
					},
				},
			},
			"dns_name": {
				Type:        schema.TypeString,
				Description: "The DNS name pointing to your cluster",
				Computed:    true,
			},
		},
		Timeouts: &resourceDefaultTimeouts,
	}
}
func checkDBaaSClusterImmutableFields(_ context.Context, diff *schema.ResourceDiff, _ interface{}) error {
	// we do not want to check in case of resource creation
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
	client := meta.(services.SdkBundle).PsqlClient

	dbaasCluster, err := dbaasService.GetPgSqlClusterDataCreate(d)

	if err != nil {
		return diag.FromErr(err)
	}
	dbaasClusterResponse, _, err := client.CreateCluster(ctx, *dbaasCluster)

	if err != nil {
		diags := diag.FromErr(fmt.Errorf("an error occured while creating a DBaaS psql cluster: %w", err))
		return diags
	}

	d.SetId(*dbaasClusterResponse.Id)

	err = utils.WaitForResourceToBeReady(ctx, d, client.IsClusterReady)
	if err != nil {
		return diag.FromErr(fmt.Errorf("creating psql %w ", err))
	}

	return resourceDbaasPgSqlClusterRead(ctx, d, meta)
}

func resourceDbaasPgSqlClusterRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	client := meta.(services.SdkBundle).PsqlClient

	cluster, apiResponse, err := client.GetCluster(ctx, d.Id())

	if err != nil {
		if apiResponse.HttpNotFound() {
			d.SetId("")
			return nil
		}
		diags := diag.FromErr(fmt.Errorf("error while fetching dbaas cluster %s: %w", d.Id(), err))
		return diags
	}

	log.Printf("[INFO] Successfully retreived cluster %s: %+v", d.Id(), cluster)

	if err := dbaasService.SetPgSqlClusterData(d, cluster); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceDbaasPgSqlClusterUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(services.SdkBundle).PsqlClient

	cluster, diags := dbaasService.GetPgSqlClusterDataUpdate(d)

	if diags != nil {
		return diags
	}

	dbaasClusterResponse, _, err := client.UpdateCluster(ctx, d.Id(), *cluster)

	if err != nil {
		diags := diag.FromErr(fmt.Errorf("an error occured while updating a dbaas cluster: %w", err))
		return diags
	}

	d.SetId(*dbaasClusterResponse.Id)

	time.Sleep(constant.SleepInterval)

	err = utils.WaitForResourceToBeReady(ctx, d, client.IsClusterReady)
	if err != nil {
		return diag.FromErr(fmt.Errorf("creating psql %w ", err))
	}

	return resourceDbaasPgSqlClusterRead(ctx, d, meta)
}

func resourceDbaasPgSqlClusterDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(services.SdkBundle).PsqlClient

	_, apiResponse, err := client.DeleteCluster(ctx, d.Id())

	if err != nil {
		if apiResponse.HttpNotFound() {
			d.SetId("")
			return nil
		}
		diags := diag.FromErr(fmt.Errorf("error while deleting dbaas cluster %s: %w", d.Id(), err))
		return diags
	}

	err = utils.WaitForResourceToBeDeleted(ctx, d, client.IsClusterDeleted)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed checking if deleted %w", err))
	}

	// wait 15 seconds after the deletion of the cluster, for the lan to be freed
	time.Sleep(constant.SleepInterval * 3)

	return nil
}

func resourceDbaasPgSqlClusterImport(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	client := meta.(services.SdkBundle).PsqlClient

	clusterId := d.Id()

	dbaasCluster, apiResponse, err := client.GetCluster(ctx, clusterId)

	if err != nil {
		if apiResponse.HttpNotFound() {
			d.SetId("")
			return nil, fmt.Errorf("dbaas cluster does not exist %q", clusterId)
		}
		return nil, fmt.Errorf("an error occured while trying to fetch the import of dbaas cluster %q", clusterId)
	}

	log.Printf("[INFO] dbaas cluster found: %+v", dbaasCluster)

	if err := dbaasService.SetPgSqlClusterData(d, dbaasCluster); err != nil {
		return nil, err
	}

	return []*schema.ResourceData{d}, nil
}
