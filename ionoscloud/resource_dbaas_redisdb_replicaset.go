package ionoscloud

import (
	"context"
	"fmt"
	"log"
	"slices"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/dbaas/redisdb"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils/constant"
)

func resourceDBaaSRedisDBReplicaSet() *schema.Resource {
	return &schema.Resource{
		CreateContext: redisDBReplicaSetCreate,
		DeleteContext: redisDBReplicaSetDelete,
		ReadContext:   redisDBReplicaSetRead,
		UpdateContext: redisDBReplicaSetUpdate,
		Importer: &schema.ResourceImporter{
			StateContext: redisDBReplicaSetImport,
		},
		Schema: map[string]*schema.Schema{
			"display_name": {
				Type:        schema.TypeString,
				Description: "The human readable name of your replica set.",
				Required:    true,
			},
			"location": {
				Type:             schema.TypeString,
				Description:      "The replica set location",
				Required:         true,
				ForceNew:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice(constant.DBaaSClusterLocations, false)),
			},
			"redis_version": {
				Type:        schema.TypeString,
				Description: "The RedisDB version of your replica set.",
				Required:    true,
			},
			"replicas": {
				Type:        schema.TypeInt,
				Description: "The total number of replicas in the replica set (one active and n-1 passive). In case of a standalone instance, the value is 1. In all other cases, the value is > 1. The replicas will not be available as read replicas, they are only standby for a failure of the active instance.",
				Required:    true,
			},
			"resources": {
				Type:        schema.TypeList,
				Description: "The resources of the individual replicas.",
				Required:    true,
				MaxItems:    1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"cores": {
							Type:        schema.TypeInt,
							Description: "The number of CPU cores per instance.",
							Required:    true,
						},
						"ram": {
							Type:        schema.TypeInt,
							Description: "The amount of memory per instance in gigabytes (GB).",
							Required:    true,
						},
						"storage": {
							Type:        schema.TypeInt,
							Description: "The size of the storage in GB. The size is derived from the amount of RAM and the persistence mode and is not configurable.",
							Computed:    true,
						},
					},
				},
			},
			"persistence_mode": {
				Type:        schema.TypeString,
				Description: "Specifies How and If data is persisted.",
				Required:    true,
			},
			"eviction_policy": {
				Type:        schema.TypeString,
				Description: "The eviction policy for the replica set.",
				Required:    true,
			},
			"connections": {
				Type:        schema.TypeList,
				Description: "The network connection for your replica set. Only one connection is allowed.",
				Required:    true,
				MaxItems:    1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"datacenter_id": {
							Type:        schema.TypeString,
							Description: "The datacenter to connect your instance to.",
							Required:    true,
						},
						"lan_id": {
							Type:        schema.TypeString,
							Description: "The numeric LAN ID to connect your instance to.",
							Required:    true,
						},
						"cidr": {
							Type:        schema.TypeString,
							Description: "The IP and subnet for your instance. Note the following unavailable IP ranges: 10.233.64.0/18, 10.233.0.0/18, 10.233.114.0/24",
							Required:    true,
						},
					},
				},
			},
			"credentials": {
				Type:        schema.TypeList,
				Description: "Credentials for the Redis replicaset.",
				Required:    true,
				MaxItems:    1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"username": {
							Type:        schema.TypeString,
							Description: "The username for the initial RedisDB user. Some system usernames are restricted (e.g. 'admin', 'standby').",
							Required:    true,
						},
						"plain_text_password": {
							Type:         schema.TypeString,
							Description:  "The password for a RedisDB user.",
							Optional:     true,
							Sensitive:    true,
							ExactlyOneOf: []string{"credentials.0.plain_text_password", "credentials.0.hashed_password"},
						},
						"hashed_password": {
							Type:         schema.TypeList,
							Description:  "The hashed password for a RedisDB user.",
							Optional:     true,
							ExactlyOneOf: []string{"credentials.0.hashed_password", "credentials.0.plain_text_password"},
							MaxItems:     1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"algorithm": {
										Type:     schema.TypeString,
										Required: true,
									},
									"hash": {
										Type:     schema.TypeString,
										Required: true,
									},
								},
							},
						},
					},
				},
			},
			"maintenance_window": {
				Type:        schema.TypeList,
				Description: "A weekly 4 hour-long window, during which maintenance might occur.",
				Optional:    true,
				Computed:    true,
				MaxItems:    1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"time": {
							Type:        schema.TypeString,
							Description: "Start of the maintenance window in UTC time.",
							Required:    true,
						},
						"day_of_the_week": {
							Type:        schema.TypeString,
							Description: "The name of the week day.",
							Required:    true,
						},
					},
				},
			},
			"initial_snapshot_id": {
				Type:        schema.TypeString,
				Description: "The ID of a snapshot to restore the replica set from. If set, the replica set will be created from the snapshot.",
				Optional:    true,
			},
			"dns_name": {
				Type:        schema.TypeString,
				Description: "The DNS name pointing to your replica set. Will be used to connect to the active/standalone instance.",
				Computed:    true,
			},
		},
		Timeouts: &resourceDefaultTimeouts,
	}
}

func redisDBReplicaSetCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(services.SdkBundle).RedisDBClient

	replicaSet := redisdb.GetRedisDBReplicaSetDataCreate(d)
	response, _, err := client.CreateRedisDPReplicaSet(ctx, replicaSet, d.Get("location").(string))
	if err != nil {
		return diag.FromErr(fmt.Errorf("an error occurred while creating a DBaaS RedisDB replica set: %w", err))
	}
	replicaSetID := *response.Id
	d.SetId(replicaSetID)
	err = utils.WaitForResourceToBeReady(ctx, d, client.IsReplicaSetReady)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error occurred while checking the status for RedisDB replica set with ID: %v, error: %w", replicaSetID, err))
	}
	// Call the read function to save the DNS name in the state (DNS name is not present in the create response).
	return redisDBReplicaSetRead(ctx, d, meta)
}

func redisDBReplicaSetDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(services.SdkBundle).RedisDBClient
	replicaSetID := d.Id()
	apiResponse, err := client.DeleteRedisDBReplicaSet(ctx, replicaSetID, d.Get("location").(string))
	if err != nil {
		if apiResponse.HttpNotFound() {
			d.SetId("")
			return nil
		}
		return diag.FromErr(fmt.Errorf("error while deleting RedisDB replica set with ID: %v, error: %w", replicaSetID, err))
	}
	err = utils.WaitForResourceToBeDeleted(ctx, d, client.IsReplicaSetDeleted)
	if err != nil {
		return diag.FromErr(fmt.Errorf("deletion check failed for RedisDB replica set with ID: %v, error: %w", replicaSetID, err))
	}

	// wait for the lan to be freed after the deletion of the replica set
	time.Sleep(constant.SleepInterval * 10)
	return nil
}

func redisDBReplicaSetRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(services.SdkBundle).RedisDBClient
	replicaSetID := d.Id()
	replicaSet, apiResponse, err := client.GetReplicaSet(ctx, replicaSetID, d.Get("location").(string))
	if err != nil {
		if apiResponse.HttpNotFound() {
			d.SetId("")
			return nil
		}
		return diag.FromErr(fmt.Errorf("error while fetching RedisDB replica set with ID: %v, error: %w", replicaSetID, err))
	}
	log.Printf("[INFO] Successfully retrieved RedisDB replica set with ID: %v, replica set info: %+v", replicaSetID, replicaSet)
	if err := client.SetRedisDBReplicaSetData(d, replicaSet); err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func redisDBReplicaSetUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(services.SdkBundle).RedisDBClient
	replicaSetID := d.Id()
	replicaSet := redisdb.GetRedisDBReplicaSetDataUpdate(d)
	response, _, err := client.UpdateRedisDBReplicaSet(ctx, replicaSetID, d.Get("location").(string), replicaSet)
	if err != nil {
		return diag.FromErr(fmt.Errorf("an error occurred while updating RedisDB replica set with ID: %v, error: %w", replicaSetID, err))
	}
	err = utils.WaitForResourceToBeReady(ctx, d, client.IsReplicaSetReady)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error occurred while checking the status for RedisDB replica set after update, replica set ID: %v, error: %w", replicaSetID, err))
	}
	if err := client.SetRedisDBReplicaSetData(d, response); err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func redisDBReplicaSetImport(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	client := meta.(services.SdkBundle).RedisDBClient
	parts := strings.Split(d.Id(), ":")
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid import ID: %q, expected ID in the format '<location>:<replica_set_id>'", d.Id())
	}
	location := parts[0]
	if !slices.Contains(constant.DBaaSClusterLocations, location) {
		return nil, fmt.Errorf("invalid import ID: %q, location must be one of %v", d.Id(), constant.DBaaSClusterLocations)
	}
	replicaSetID := parts[1]
	replicaSet, apiResponse, err := client.GetReplicaSet(ctx, replicaSetID, location)
	if err != nil {
		if apiResponse.HttpNotFound() {
			d.SetId("")
			return nil, fmt.Errorf("RedisDB replica set does not exist, error: %w", err)
		}
		return nil, fmt.Errorf("an error occurred while trying to import RedisDB replica set with ID: %v, error: %w", replicaSetID, err)
	}
	log.Printf("[INFO] RedisDB replica set found: %+v", replicaSet)
	if err := d.Set("location", location); err != nil {
		return nil, utils.GenerateSetError("RedisDB replica set", "location", err)
	}
	if err := client.SetRedisDBReplicaSetData(d, replicaSet); err != nil {
		return nil, err
	}
	return []*schema.ResourceData{d}, nil
}
