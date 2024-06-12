package ionoscloud

import (
	"context"
	"fmt"
	"log"
	"time"

	ionoscloud "github.com/ionos-cloud/sdk-go-dbaas-mongo"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils/constant"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/dbaas"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils"
)

func resourceDbaasMongoDBCluster() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDbaasMongoClusterCreate,
		ReadContext:   resourceDbaasMongoClusterRead,
		UpdateContext: resourceDbaasMongoClusterUpdate,
		DeleteContext: resourceDbaasMongoClusterDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceDbaasMongoClusterImport,
		},
		Schema: map[string]*schema.Schema{
			"maintenance_window": {
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "A weekly 4 hour-long window, during which maintenance might occur",
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
			"mongodb_version": {
				Type:             schema.TypeString,
				Description:      "The MongoDB version of your cluster. Update forces cluster re-creation.",
				Required:         true,
				ForceNew:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringIsNotWhiteSpace),
			},
			"instances": {
				Type:             schema.TypeInt,
				Description:      "The total number of instances in the cluster (one master and n-1 standbys). Example: 1, 3, 5, 7. For enterprise edition at least 3.",
				Required:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.IntBetween(1, 7)),
			},
			"display_name": {
				Type:        schema.TypeString,
				Description: "The name of your cluster.",
				Required:    true,
			},
			"location": {
				Type: schema.TypeString,
				Description: "The physical location where the cluster will be created. This will be where all of your instances live. " +
					"Property cannot be modified after datacenter creation (disallowed in update requests). Available locations: de/txl, gb/lhr, es/vit. Update forces cluster re-creation.",
				Required:         true,
				ForceNew:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringIsNotWhiteSpace),
			},
			"connections": {
				Type:        schema.TypeList,
				MaxItems:    1,
				Required:    true,
				Description: "Details about the network connection for your cluster.",
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
						"cidr_list": {
							Type:     schema.TypeList,
							Required: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Description: "The list of IPs and subnet for your cluster. Note the following unavailable IP ranges:10.233.64.0/18, 10.233.0.0/18, 10.233.114.0/24. example: [192.168.1.100/24, 192.168.1.101/24]",
						},
						// add after api adds support
						// "whitelist": {
						//	Type:     schema.TypeList,
						//	Optional: true,
						//	Elem: &schema.Schema{
						//		Type: schema.TypeString,
						//	},
						//	Description: "List of whitelisted CIDRs",
						// },
					},
				},
			},
			"template_id": {
				Type: schema.TypeString,
				Description: "The unique ID of the template, which specifies the number of cores, storage size, and memory. " +
					"You cannot downgrade to a smaller template or minor edition (e.g. from business to playground). " +
					"To get a list of all templates to confirm the changes use the /templates endpoint.",
				Optional:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.IsUUID),
			},
			"connection_string": {
				Type:        schema.TypeString,
				Description: "The connection string for your cluster.",
				Computed:    true,
			},
			// enterprise edition below
			"type": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The cluster type, either `replicaset` or `sharded-cluster`",
			},
			"shards": {
				Type:             schema.TypeInt,
				Optional:         true,
				Description:      "The total number of shards in the cluster.",
				ValidateDiagFunc: validation.ToDiagFunc(validation.IntBetween(1, 100)),
			},
			"bi_connector": {
				Type:     schema.TypeList,
				MaxItems: 1,
				Optional: true,
				Computed: true,
				Description: "The MongoDB Connector for Business Intelligence allows you to " +
					"query a MongoDB database using SQL commands to aid in data analysis.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"enabled": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "Enable or disable the BiConnector.",
							Default:     false,
						},
						"host": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The host where this new BI Connector is installed.",
						},
						"port": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Port number used when connecting to this new BI Connector.",
						},
					},
				},
			},
			"ram": {
				Type:          schema.TypeInt,
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{"template_id"},
				Description:   "The amount of memory per instance in megabytes. Multiple of 1024",
				ValidateDiagFunc: validation.AllDiag(validation.ToDiagFunc(validation.IntAtLeast(2048)),
					validation.ToDiagFunc(validation.IntDivisibleBy(1024))),
			},
			"storage_size": {
				Type:          schema.TypeInt,
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{"template_id"},

				Description: "The amount of storage per instance in megabytes. At least 5120, at most 2097152",
				ForceNew:    true,
				ValidateDiagFunc: validation.AllDiag(validation.ToDiagFunc(validation.IntBetween(5120, 2097152)),
					validation.ToDiagFunc(validation.IntDivisibleBy(1024))),
			},
			"storage_type": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{"template_id"},
				Description:   "The storage type. One of : HDD, SSD, SSD Standard, SSD Premium",
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{string(ionoscloud.STORAGETYPE_HDD), string(ionoscloud.STORAGETYPE_SSD),
					string(ionoscloud.STORAGETYPE_SSD_STANDARD), string(ionoscloud.STORAGETYPE_SSD_PREMIUM)}, false)),
			},
			"cores": {
				Type:             schema.TypeInt,
				Optional:         true,
				Computed:         true,
				ConflictsWith:    []string{"template_id"},
				ValidateDiagFunc: validation.ToDiagFunc(validation.IntAtLeast(1)),
				Description:      "The number of CPU cores per instance.",
			},
			"edition": {
				Type:             schema.TypeString,
				Optional:         true,
				Computed:         true,
				Description:      "The cluster edition. Must be one of: playground, business, enterprise",
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"enterprise", "playground", "business"}, false)),
			},
			// to be added when there is api support
			// "from_backup": {
			//	Type:        schema.TypeList,
			//	MaxItems:    1,
			//	Description: "Creates the cluster based on the existing backup.",
			//	Optional:    true,
			//	Elem: &schema.Resource{
			//		Schema: map[string]*schema.Schema{
			//			"snapshot_id": {
			//				Type:        schema.TypeString,
			//				Description: "The unique ID of the snapshot you want to restore.",
			//				Computed:    true,
			//			},
			//			"recovery_target_time": {
			//				Type:        schema.TypeString,
			//				Computed:    true,
			//				Description: " If this value is supplied as ISO 8601 timestamp, the backup will be replayed up until the given timestamp.",
			//			},
			//		},
			//	},
			// },
			"backup": {
				Type:        schema.TypeList,
				MaxItems:    1,
				Optional:    true,
				Description: "Backup related properties.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"snapshot_interval_hours": {
							Type:             schema.TypeInt,
							Optional:         true,
							Description:      "Number of hours between snapshots.",
							Default:          24,
							ValidateDiagFunc: validation.ToDiagFunc(validation.IntInSlice([]int{6, 8, 12, 24})),
						},
						"point_in_time_window_hours": {
							Type:             schema.TypeInt,
							Optional:         true,
							Description:      "Number of hours in the past for which a point-in-time snapshot can be created.",
							ValidateDiagFunc: validation.ToDiagFunc(validation.IntBetween(1, 24)),
						},
						// will be added at a later date
						// "backup_retention": {
						//	Type:        schema.TypeList,
						//	MaxItems:    1,
						//	Optional:    true,
						//	Description: "Backup retention related properties.",
						//	Elem: &schema.Resource{
						//		Schema: map[string]*schema.Schema{
						//			"snapshot_retention_days": {
						//				Type:             schema.TypeInt,
						//				Optional:         true,
						//				Description:      "Number of days to keep recent snapshots.",
						//				ValidateDiagFunc: validation.ToDiagFunc(validation.IntBetween(2, 5)),
						//				Default:          2,
						//			},
						//			"daily_snapshot_retention_days": {
						//				Type:             schema.TypeInt,
						//				Optional:         true,
						//				Description:      "Number of days to retain daily snapshots.",
						//				ValidateDiagFunc: validation.ToDiagFunc(validation.IntBetween(1, 365)),
						//				Default:          0,
						//			},
						//			"weekly_snapshot_retention_weeks": {
						//				Type:             schema.TypeInt,
						//				Optional:         true,
						//				Description:      "Number of weeks to retain weekly snapshots.",
						//				ValidateDiagFunc: validation.ToDiagFunc(validation.IntBetween(1, 52)),
						//				Default:          2,
						//			},
						//			"monthly_snapshot_retention_months": {
						//				Type:             schema.TypeInt,
						//				Optional:         true,
						//				Description:      "Number of months to retain monthly snapshots.",
						//				ValidateDiagFunc: validation.ToDiagFunc(validation.IntBetween(1, 36)),
						//				Default:          1,
						//			},
						//		},
						//	},
						// },
						"location": {
							Type:     schema.TypeString,
							Optional: true,
							Description: "The location where the cluster backups will be stored. " +
								"If not set, the backup is stored in the nearest location of the cluster. Examples: de, eu-sounth-2, eu-central-2",
						},
					},
				},
			},
		},
		Timeouts: &resourceDefaultTimeouts,
	}
}

func resourceDbaasMongoClusterCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(services.SdkBundle).MongoClient

	if err := dbaas.MongoClusterCheckRequiredFieldsSet(d); err != nil {
		return diag.FromErr(err)
	}

	cluster, err := dbaas.SetMongoClusterCreateProperties(d)
	if err != nil {
		return diag.FromErr(err)
	}
	createdCluster, _, err := client.CreateCluster(ctx, *cluster)
	if err != nil {
		d.SetId("")
		diags := diag.FromErr(fmt.Errorf("create error for dbaas mongo cluster %s: %w", d.Id(), err))
		return diags
	}
	d.SetId(*createdCluster.Id)

	err = utils.WaitForResourceToBeReady(ctx, d, client.IsClusterReady)
	if err != nil {
		return diag.FromErr(fmt.Errorf("creating %w ", err))
	}

	return resourceDbaasMongoClusterRead(ctx, d, meta)
}

func resourceDbaasMongoClusterUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(services.SdkBundle).MongoClient
	clusterId := d.Id()
	patchRequest := dbaas.SetMongoClusterPatchProperties(d)

	_, apiResponse, err := client.UpdateCluster(ctx, clusterId, *patchRequest)
	if err != nil {
		if apiResponse.HttpNotFound() {
			d.SetId("")
			return nil
		}
		diags := diag.FromErr(fmt.Errorf("update error while fetching dbaas mongo cluster %s: %w", d.Id(), err))
		return diags
	}

	err = utils.WaitForResourceToBeReady(ctx, d, client.IsClusterReady)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed checking if ready %w", err))
	}

	return resourceDbaasMongoClusterRead(ctx, d, meta)
}

func resourceDbaasMongoClusterRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(services.SdkBundle).MongoClient

	cluster, apiResponse, err := client.GetCluster(ctx, d.Id())

	if err != nil {
		if apiResponse.HttpNotFound() {
			d.SetId("")
			return nil
		}
		diags := diag.FromErr(fmt.Errorf("read error while fetching dbaas mongo cluster %s: %w", d.Id(), err))
		return diags
	}

	log.Printf("[INFO] Successfully retreived cluster %s: %+v", d.Id(), cluster)

	if err := dbaas.SetMongoDBClusterData(d, cluster); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceDbaasMongoClusterDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(services.SdkBundle).MongoClient

	_, apiResponse, err := client.DeleteCluster(ctx, d.Id())

	if err != nil {
		if apiResponse.HttpNotFound() {
			d.SetId("")
			return nil
		}
		diags := diag.FromErr(fmt.Errorf("error while deleting mongo dbaas cluster %s: %w", d.Id(), err))
		return diags
	}

	err = utils.WaitForResourceToBeDeleted(ctx, d, client.IsClusterDeleted)
	if err != nil {
		return diag.FromErr(fmt.Errorf("The check for cluster deletion failed with the following error: %w", err))
	}
	// wait 15 seconds after the deletion of the cluster, for the lan to be freed
	time.Sleep(constant.SleepInterval * 3)

	return nil
}

func resourceDbaasMongoClusterImport(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	client := meta.(services.SdkBundle).MongoClient

	clusterId := d.Id()

	dbaasCluster, apiResponse, err := client.GetCluster(ctx, clusterId)

	if err != nil {
		if apiResponse.HttpNotFound() {
			d.SetId("")
			return nil, fmt.Errorf("dbaas cluster does not exist %q, error:%w", clusterId, err)
		}
		return nil, fmt.Errorf("an error occured while trying to fetch the import of dbaas cluster %q, error:%w", clusterId, err)
	}

	log.Printf("[INFO] dbaas cluster found: %+v", dbaasCluster)

	if err := dbaas.SetMongoDBClusterData(d, dbaasCluster); err != nil {
		return nil, err
	}

	return []*schema.ResourceData{d}, nil
}
