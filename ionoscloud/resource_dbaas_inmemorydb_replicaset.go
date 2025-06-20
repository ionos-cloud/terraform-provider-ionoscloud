package ionoscloud

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/bundleclient"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/dbaas/inmemorydb"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils/constant"
)

var (
	persistenceModes = []string{"None", "AOF", "RDB", "RDB_AOF"}
	evictionPolicies = []string{"allkeys-lru", "allkeys-lfu", "allkeys-random",
		"volatile-lru", "volatile-lfu", "volatile-random", "volatile-ttl", "noeviction"}
)

func resourceDBaaSInMemoryDBReplicaSet() *schema.Resource {
	return &schema.Resource{
		CreateContext: replicaSetCreate,
		DeleteContext: replicaSetDelete,
		ReadContext:   replicaSetRead,
		UpdateContext: replicaSetUpdate,
		Importer: &schema.ResourceImporter{
			StateContext: replicaSetImport,
		},
		Schema: map[string]*schema.Schema{
			"display_name": {
				Type:        schema.TypeString,
				Description: "The human readable name of your replica set.",
				Required:    true,
			},
			"location": {
				Type:        schema.TypeString,
				Description: "The replica set location",
				Optional:    true,
				ForceNew:    true,
			},
			"version": {
				Type:        schema.TypeString,
				Description: "The InMemoryDB version of your replica set.",
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
				Type:             schema.TypeString,
				Description:      "Specifies How and If data is persisted.",
				Required:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice(persistenceModes, true)),
			},
			"eviction_policy": {
				Type:             schema.TypeString,
				Description:      "The eviction policy for the replica set.",
				Required:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice(evictionPolicies, true)),
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
							ForceNew:    true,
						},
						"lan_id": {
							Type:        schema.TypeString,
							Description: "The numeric LAN ID to connect your instance to.",
							Required:    true,
							ForceNew:    true,
						},
						"cidr": {
							Type:        schema.TypeString,
							Description: "The IP and subnet for your instance. Note the following unavailable IP ranges: 10.233.64.0/18, 10.233.0.0/18, 10.233.114.0/24",
							Required:    true,
							ForceNew:    true,
						},
					},
				},
			},
			"credentials": {
				Type:        schema.TypeList,
				Description: "Credentials for the InMemoryDB replicaset.",
				Required:    true,
				MaxItems:    1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"username": {
							Type:        schema.TypeString,
							Description: "The username for the initial InMemoryDB user. Some system usernames are restricted (e.g. 'admin', 'standby').",
							Required:    true,
							ForceNew:    true,
						},
						"plain_text_password": {
							Type:         schema.TypeString,
							Description:  "The password for a InMemoryDB user.",
							Optional:     true,
							Sensitive:    true,
							ForceNew:     true,
							ExactlyOneOf: []string{"credentials.0.plain_text_password", "credentials.0.hashed_password"},
						},
						"hashed_password": {
							Type:         schema.TypeList,
							Description:  "The hashed password for a InMemoryDB user.",
							Optional:     true,
							ExactlyOneOf: []string{"credentials.0.hashed_password", "credentials.0.plain_text_password"},
							MaxItems:     1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"algorithm": {
										Type:     schema.TypeString,
										Required: true,
										ForceNew: true,
									},
									"hash": {
										Type:     schema.TypeString,
										Required: true,
										ForceNew: true,
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
				ForceNew:    true,
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

func replicaSetCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(bundleclient.SdkBundle).InMemoryDBClient

	replicaSet := inmemorydb.GetReplicaSetDataCreate(d)
	response, _, err := client.CreateReplicaSet(ctx, replicaSet, d.Get("location").(string))
	if err != nil {
		return diag.FromErr(fmt.Errorf("an error occurred while creating an InMemoryDB replica set: %w", err))
	}
	replicaSetID := response.Id
	d.SetId(replicaSetID)
	err = utils.WaitForResourceToBeReady(ctx, d, client.IsReplicaSetReady)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error occurred while checking the status for InMemoryDB replica set with ID: %v, error: %w", replicaSetID, err))
	}
	// Call the read function to save the DNS name in the state (DNS name is not present in the creation response).
	return replicaSetRead(ctx, d, meta)
}

func replicaSetDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(bundleclient.SdkBundle).InMemoryDBClient
	replicaSetID := d.Id()
	apiResponse, err := client.DeleteReplicaSet(ctx, replicaSetID, d.Get("location").(string))
	if err != nil {
		if apiResponse.HttpNotFound() {
			d.SetId("")
			return nil
		}
		return diag.FromErr(fmt.Errorf("error while deleting InMemoryDB replica set with ID: %v, error: %w", replicaSetID, err))
	}
	err = utils.WaitForResourceToBeDeleted(ctx, d, client.IsReplicaSetDeleted)
	if err != nil {
		return diag.FromErr(fmt.Errorf("deletion check failed for InMemoryDB replica set with ID: %v, error: %w", replicaSetID, err))
	}

	// wait for the lan to be freed after the deletion of the replica set
	time.Sleep(constant.SleepInterval * 10)
	return nil
}

func replicaSetRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(bundleclient.SdkBundle).InMemoryDBClient
	replicaSetID := d.Id()
	replicaSet, apiResponse, err := client.GetReplicaSet(ctx, replicaSetID, d.Get("location").(string))
	if err != nil {
		if apiResponse.HttpNotFound() {
			d.SetId("")
			return nil
		}
		return diag.FromErr(fmt.Errorf("error while fetching InMemoryDB replica set with ID: %v, error: %w", replicaSetID, err))
	}
	log.Printf("[INFO] Successfully retrieved InMemoryDB replica set with ID: %v, replica set info: %+v", replicaSetID, replicaSet)
	if err := client.SetReplicaSetData(d, replicaSet); err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func replicaSetUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(bundleclient.SdkBundle).InMemoryDBClient
	replicaSetID := d.Id()
	replicaSet := inmemorydb.GetReplicaSetDataUpdate(d)
	response, _, err := client.UpdateReplicaSet(ctx, replicaSetID, d.Get("location").(string), replicaSet)
	if err != nil {
		return diag.FromErr(fmt.Errorf("an error occurred while updating InMemoryDB replica set with ID: %v, error: %w", replicaSetID, err))
	}
	err = utils.WaitForResourceToBeReady(ctx, d, client.IsReplicaSetReady)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error occurred while checking the status for InMemoryDB replica set after update, replica set ID: %v, error: %w", replicaSetID, err))
	}
	if err := client.SetReplicaSetData(d, response); err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func replicaSetImport(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	client := meta.(bundleclient.SdkBundle).InMemoryDBClient
	parts := strings.Split(d.Id(), ":")
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid import ID: %q, expected ID in the format '<location>:<replica_set_id>'", d.Id())
	}
	location := parts[0]
	replicaSetID := parts[1]
	replicaSet, apiResponse, err := client.GetReplicaSet(ctx, replicaSetID, location)
	if err != nil {
		if apiResponse.HttpNotFound() {
			d.SetId("")
			return nil, fmt.Errorf("InMemoryDB replica set does not exist, error: %w", err)
		}
		return nil, fmt.Errorf("an error occurred while trying to import InMemoryDB replica set with ID: %v, error: %w", replicaSetID, err)
	}
	log.Printf("[INFO] InMemoryDB replica set found: %+v", replicaSet)
	if err := d.Set("location", location); err != nil {
		return nil, utils.GenerateSetError("InMemoryDB replica set", "location", err)
	}
	if err := client.SetReplicaSetData(d, replicaSet); err != nil {
		return nil, err
	}
	return []*schema.ResourceData{d}, nil
}
