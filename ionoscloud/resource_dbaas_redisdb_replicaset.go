package ionoscloud

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
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
			// TODO -- Check if we need to add validation here
			// TODO -- Check what to do with the default values
			"persistence_mode": {
				Type: schema.TypeString,
				// TODO -- In documentation, add the full description from the swagger.
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
						// TODO -- Check if we need to add validation here
						"username": {
							Type:        schema.TypeString,
							Description: "The username for the initial RedisDB user. Some system usernames are restricted (e.g. 'admin', 'standby').",
							Required:    true,
						},
						"plain_text_password": {
							Type:         schema.TypeString,
							Description:  "The password for a RedisDB user.",
							Optional:     true,
							ExactlyOneOf: []string{"plain_text_password", "hashed_password"},
						},
						"hashed_password": {
							Type:         schema.TypeList,
							Description:  "The hashed password for a RedisDB user.",
							Optional:     true,
							ExactlyOneOf: []string{"hashed_password", "plain_text_password"},
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
			// TODO -- Check if this needs to be marked as 'Computed'
			"initial_snapshot_id": {
				Type:        schema.TypeString,
				Description: "The ID of a snapshot to restore the replica set from. If set, the replica set will be created from the snapshot.",
				Optional:    true,
			},
		},
		Timeouts: &resourceDefaultTimeouts,
	}
}

func redisDBReplicaSetCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}

func redisDBReplicaSetDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}

func redisDBReplicaSetRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}

func redisDBReplicaSetUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}

func redisDBReplicaSetImport(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	return nil, nil
}
