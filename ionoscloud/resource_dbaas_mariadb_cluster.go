package ionoscloud

import (
	"context"
	"fmt"
	"log"
	"slices"
	"strings"
	"time"

	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils/constant"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/dbaas/mariadb"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils"
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
		Schema: map[string]*schema.Schema{
			"mariadb_version": {
				Type:        schema.TypeString,
				Description: "The MariaDB version of your cluster.",
				Required:    true,
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
			"dns_name": {
				Type:        schema.TypeString,
				Description: "The DNS name pointing to your cluster.",
				Computed:    true,
			},
		},
		Timeouts: &resourceDefaultTimeouts,
	}
}

func mariaDBClusterCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(services.SdkBundle).MariaDBClient

	cluster, err := mariadb.GetMariaDBClusterDataCreate(d)
	if err != nil {
		return diag.FromErr(err)
	}

	response, _, err := client.CreateCluster(ctx, *cluster, d.Get("location").(string))
	if err != nil {
		return diag.FromErr(fmt.Errorf("an error occurred while creating a DBaaS MariaDB cluster: %w", err))
	}
	clusterID := *response.Id
	d.SetId(clusterID)

	err = utils.WaitForResourceToBeReady(ctx, d, client.IsClusterReady)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error occurred while checking the status for MariaDB cluster with ID: %v, error: %w", clusterID, err))
	}
	if err := client.SetMariaDBClusterData(d, response); err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func mariaDBClusterDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(services.SdkBundle).MariaDBClient
	clusterID := d.Id()
	_, apiResponse, err := client.DeleteCluster(ctx, d.Id(), d.Get("location").(string))
	if err != nil {
		if apiResponse.HttpNotFound() {
			d.SetId("")
			return nil
		}
		return diag.FromErr(fmt.Errorf("error while deleting MariaDB cluster with ID: %v, error: %w", clusterID, err))
	}
	err = utils.WaitForResourceToBeDeleted(ctx, d, client.IsClusterDeleted)
	if err != nil {
		return diag.FromErr(fmt.Errorf("deletion check failed for MariaDB cluster with ID: %v, error: %w", clusterID, err))
	}

	// wait after the deletion of the cluster, for the lan to be freed
	time.Sleep(constant.SleepInterval * 10)

	return nil
}

func mariaDBClusterImport(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	client := meta.(services.SdkBundle).MariaDBClient
	parts := strings.Split(d.Id(), ":")
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid import ID: %q, expected ID in the format '<location>:<cluster_id>'", d.Id())
	}
	location := parts[0]
	if !slices.Contains(constant.Locations, location) {
		return nil, fmt.Errorf("invalid import ID: %q, location must be one of %v", d.Id(), constant.Locations)
	}
	clusterID := parts[1]

	cluster, apiResponse, err := client.GetCluster(ctx, clusterID, location)
	if err != nil {
		if apiResponse.HttpNotFound() {
			d.SetId("")
			return nil, fmt.Errorf("MariaDB cluster with ID: %v does not exist", clusterID)
		}
		return nil, fmt.Errorf("an error occurred while trying to import MariaDB cluster with ID: %v, error: %w", clusterID, err)
	}

	log.Printf("[INFO] MariaDB cluster found: %+v", cluster)

	if err := d.Set("location", location); err != nil {
		return nil, utils.GenerateSetError("MariaDB cluster", "location", err)
	}
	if err := client.SetMariaDBClusterData(d, cluster); err != nil {
		return nil, err
	}

	return []*schema.ResourceData{d}, nil
}

func mariaDBClusterRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(services.SdkBundle).MariaDBClient
	clusterID := d.Id()
	cluster, apiResponse, err := client.GetCluster(ctx, clusterID, d.Get("location").(string))
	if err != nil {
		if apiResponse.HttpNotFound() {
			d.SetId("")
			return nil
		}
		return diag.FromErr(fmt.Errorf("error while fetching MariaDB cluster with ID: %v, error: %w", clusterID, err))
	}
	log.Printf("[INFO] Successfully retrieved MariaDB cluster with ID: %v, cluster info: %+v", clusterID, cluster)

	if err := client.SetMariaDBClusterData(d, cluster); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func mariaDBClusterUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(services.SdkBundle).MariaDBClient

	clusterID := d.Id()
	cluster, err := mariadb.GetMariaDBClusterDataUpdate(d)
	if err != nil {
		return diag.FromErr(err)
	}

	response, _, err := client.UpdateCluster(ctx, *cluster, clusterID, d.Get("location").(string))
	if err != nil {
		return diag.FromErr(fmt.Errorf("an error occurred while updating DBaaS MariaDB cluster with ID: %v, error: %w", clusterID, err))
	}

	err = utils.WaitForResourceToBeReady(ctx, d, client.IsClusterReady)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error occurred while checking the status for MariaDB cluster with ID: %v, error: %w", clusterID, err))
	}
	if err := client.SetMariaDBClusterData(d, response); err != nil {
		return diag.FromErr(err)
	}
	return nil
}
