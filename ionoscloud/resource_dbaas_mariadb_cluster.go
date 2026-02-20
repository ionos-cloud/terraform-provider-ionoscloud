package ionoscloud

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	semversion "github.com/hashicorp/go-version"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/bundleclient"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/dbaas/mariadb"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils/constant"
)

func resourceDBaaSMariaDBCluster() *schema.Resource {
	return &schema.Resource{
		CreateContext: mariaDBClusterCreate,
		DeleteContext: mariaDBClusterDelete,
		ReadContext:   mariaDBClusterRead,
		UpdateContext: mariaDBClusterUpdate,
		Importer: &schema.ResourceImporter{
			StateContext: mariaDBClusterImport,
		},
		CustomizeDiff: errorOnVersionDowngrade,
		Schema: map[string]*schema.Schema{
			"mariadb_version": {
				Type:             schema.TypeString,
				Description:      "The MariaDB version of your cluster. Cannot be downgraded.",
				Required:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringIsNotEmpty),
			},
			"location": {
				Type:        schema.TypeString,
				Description: "The cluster location",
				Optional:    true,
				ForceNew:    true,
			},
			"instances": {
				Type:             schema.TypeInt,
				Description:      "The total number of instances in the cluster (one primary and n-1 secondary).",
				Required:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.IntBetween(1, 5)),
			},
			"cores": {
				Type:             schema.TypeInt,
				Description:      "The number of CPU cores per instance.",
				Required:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.IntAtLeast(1)),
			},
			"ram": {
				Type:             schema.TypeInt,
				Description:      "The amount of memory per instance in gigabytes (GB).",
				Required:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.IntAtLeast(2)),
			},
			"storage_size": {
				Type:             schema.TypeInt,
				Description:      "The amount of storage per instance in gigabytes (GB).",
				Required:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.IntBetween(10, 2000)),
			},
			"connections": {
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "The network connection for your cluster. Only one connection is allowed.",
				Required:    true,
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
							Description:      "The numeric LAN ID to connect your cluster to.",
							Required:         true,
							ValidateDiagFunc: validation.ToDiagFunc(validation.StringIsNotWhiteSpace),
						},
						"cidr": {
							Type:             schema.TypeString,
							Description:      "The IP and subnet for your cluster.",
							Required:         true,
							ValidateDiagFunc: validation.ToDiagFunc(VerifyUnavailableIPs),
						},
					},
				},
			},
			"display_name": {
				Type:        schema.TypeString,
				Description: "The friendly name of your cluster.",
				Required:    true,
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
							Description:      "The username for the initial MariaDB user. Some system usernames are restricted (e.g 'mariadb', 'admin', 'standby').",
							Required:         true,
							ValidateDiagFunc: validation.ToDiagFunc(validation.All(validation.StringIsNotWhiteSpace, validation.StringLenBetween(1, 63))),
						},
						"password": {
							Type:             schema.TypeString,
							Description:      "The password for a MariaDB user.",
							Required:         true,
							Sensitive:        true,
							ValidateDiagFunc: validation.ToDiagFunc(validation.StringIsNotWhiteSpace),
						},
					},
				},
			},
			"maintenance_window": {
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "A weekly 4 hour-long window, during which maintenance might occur.",
				Optional:    true,
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"time": {
							Type:             schema.TypeString,
							Description:      "Start of the maintenance window in UTC time.",
							Required:         true,
							ValidateDiagFunc: validation.ToDiagFunc(validation.StringIsNotWhiteSpace),
						},
						"day_of_the_week": {
							Type:             schema.TypeString,
							Description:      "The name of the week day.",
							Required:         true,
							ValidateDiagFunc: validation.ToDiagFunc(validation.IsDayOfTheWeek(true)),
						},
					},
				},
			},
			"backup": {
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "Properties configuring the backup of the cluster. Immutable, change forces re-creation of the cluster.",
				Optional:    true,
				ForceNew:    true,
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"location": {
							Type:             schema.TypeString,
							Description:      "The IONOS Object Storage location where the backups will be stored.",
							Required:         true,
							ForceNew:         true,
							ValidateDiagFunc: validation.ToDiagFunc(validation.StringIsNotWhiteSpace),
						},
					},
				},
			},
			"dns_name": {
				Type:        schema.TypeString,
				Description: "The DNS name pointing to your cluster.",
				Computed:    true,
			},
		},
		Timeouts: &resourceDefaultTimeouts,
	}
}

func errorOnVersionDowngrade(_ context.Context, diff *schema.ResourceDiff, _ interface{}) error {
	// we do not want to check in case of resource creation
	if diff.Id() == "" {
		return nil
	}
	if diff.HasChange("mariadb_version") {
		oldValue, newValue := diff.GetChange("mariadb_version")
		oldVersionStr := oldValue.(string)
		newVersionStr := newValue.(string)
		oldVersion, err := semversion.NewVersion(oldVersionStr)
		if err != nil {
			return err
		}
		newVersion, err := semversion.NewVersion(newVersionStr)
		if err != nil {
			return err
		}
		if newVersion.LessThan(oldVersion) {
			return fmt.Errorf("downgrade is not supported from %s to %s", oldVersionStr, newVersionStr)
		}
	}
	return nil
}

func mariaDBClusterCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(bundleclient.SdkBundle).MariaDBClient

	cluster, err := mariadb.GetMariaDBClusterDataCreate(d)
	if err != nil {
		return utils.ToDiags(d, err.Error(), nil)
	}

	response, apiResponse, err := client.CreateCluster(ctx, *cluster, d.Get("location").(string))
	if err != nil {
		return utils.ToDiags(d, fmt.Sprintf("an error occurred while creating a DBaaS MariaDB cluster: %s", err), &utils.DiagsOpts{StatusCode: apiResponse.StatusCode})
	}
	clusterID := *response.Id
	d.SetId(clusterID)

	err = utils.WaitForResourceToBeReady(ctx, d, client.IsClusterReady)
	if err != nil {
		return utils.ToDiags(d, fmt.Sprintf("error occurred while checking the status for MariaDB cluster with ID: %v, error: %s", clusterID, err), nil)
	}
	if err := client.SetMariaDBClusterData(d, response); err != nil {
		return utils.ToDiags(d, err.Error(), nil)
	}
	return nil
}

func mariaDBClusterDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(bundleclient.SdkBundle).MariaDBClient
	clusterID := d.Id()
	_, apiResponse, err := client.DeleteCluster(ctx, d.Id(), d.Get("location").(string))
	if err != nil {
		if apiResponse.HttpNotFound() {
			d.SetId("")
			return nil
		}
		return utils.ToDiags(d, fmt.Sprintf("error while deleting MariaDB cluster with ID: %v, error: %s", clusterID, err), &utils.DiagsOpts{StatusCode: apiResponse.StatusCode})
	}
	err = utils.WaitForResourceToBeDeleted(ctx, d, client.IsClusterDeleted)
	if err != nil {
		return utils.ToDiags(d, fmt.Sprintf("deletion check failed for MariaDB cluster with ID: %v, error: %s", clusterID, err), &utils.DiagsOpts{Timeout: schema.TimeoutDelete})
	}

	// wait after the deletion of the cluster, for the lan to be freed
	time.Sleep(constant.SleepInterval * 10)

	return nil
}

func mariaDBClusterImport(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	client := meta.(bundleclient.SdkBundle).MariaDBClient
	parts := strings.Split(d.Id(), ":")
	if len(parts) != 2 {
		return nil, utils.ToError(d, "invalid import, expected ID in the format '<location>:<cluster_id>'", nil)
	}
	location := parts[0]
	clusterID := parts[1]

	cluster, apiResponse, err := client.GetCluster(ctx, clusterID, location)
	if err != nil {
		if apiResponse.HttpNotFound() {
			d.SetId("")
			return nil, utils.ToError(d, fmt.Sprintf("MariaDB cluster with ID: %v does not exist", clusterID), &utils.DiagsOpts{StatusCode: apiResponse.StatusCode})
		}
		return nil, utils.ToError(d, fmt.Sprintf("an error occurred while trying to import MariaDB cluster with ID: %v, error: %s", clusterID, err), &utils.DiagsOpts{StatusCode: apiResponse.StatusCode})
	}

	log.Printf("[INFO] MariaDB cluster found: %+v", cluster)

	if err := d.Set("location", location); err != nil {
		return nil, utils.GenerateSetError("MariaDB cluster", "location", err)
	}
	if err := client.SetMariaDBClusterData(d, cluster); err != nil {
		return nil, utils.ToError(d, err.Error(), nil)
	}

	return []*schema.ResourceData{d}, nil
}

func mariaDBClusterRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(bundleclient.SdkBundle).MariaDBClient
	clusterID := d.Id()
	cluster, apiResponse, err := client.GetCluster(ctx, clusterID, d.Get("location").(string))
	if err != nil {
		if apiResponse.HttpNotFound() {
			d.SetId("")
			return nil
		}
		return utils.ToDiags(d, fmt.Sprintf("error while fetching MariaDB cluster with ID: %v, error: %s", clusterID, err), &utils.DiagsOpts{StatusCode: apiResponse.StatusCode})
	}
	log.Printf("[INFO] Successfully retrieved MariaDB cluster with ID: %v, cluster info: %+v", clusterID, cluster)

	if err := client.SetMariaDBClusterData(d, cluster); err != nil {
		return utils.ToDiags(d, err.Error(), nil)
	}

	return nil
}

func mariaDBClusterUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(bundleclient.SdkBundle).MariaDBClient

	clusterID := d.Id()
	cluster, err := mariadb.GetMariaDBClusterDataUpdate(d)
	if err != nil {
		return utils.ToDiags(d, err.Error(), nil)
	}
	location := d.Get("location").(string)
	response, apiResponse, err := client.UpdateCluster(ctx, *cluster, clusterID, location)
	if err != nil {
		return utils.ToDiags(d, fmt.Sprintf("an error occurred while updating DBaaS MariaDB cluster with ID: %v in location %s, error: %s", clusterID, location, err), &utils.DiagsOpts{StatusCode: apiResponse.StatusCode})
	}

	err = utils.WaitForResourceToBeReady(ctx, d, client.IsClusterReady)
	if err != nil {
		return utils.ToDiags(d, fmt.Sprintf("error occurred while checking the status for MariaDB cluster with ID: %v in location %s, error: %s", clusterID, location, err), nil)
	}
	if err := client.SetMariaDBClusterData(d, response); err != nil {
		return utils.ToDiags(d, err.Error(), nil)
	}
	return nil
}
