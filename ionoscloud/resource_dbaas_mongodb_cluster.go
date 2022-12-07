package ionoscloud

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	mongo "github.com/ionos-cloud/sdk-go-dbaas-mongo"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/dbaas"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils"
	"log"
	"strings"
	"time"
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
			"mongodb_version": {
				Type:             schema.TypeString,
				Description:      "The MongoDB version of your cluster.",
				Required:         true,
				ForceNew:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringIsNotWhiteSpace),
			},
			"instances": {
				Type:             schema.TypeInt,
				Description:      "The total number of instances in the cluster (one master and n-1 standbys)",
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
					"Property cannot be modified after datacenter creation (disallowed in update requests). Available locations: de/txl, gb/lhr, es/vit",
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
							Description: "The list of IPs and subnet for your cluster.\n          Note the following unavailable IP ranges:\n          10.233.64.0/18\n          10.233.0.0/18\n          10.233.114.0/24 		\n example: [192.168.1.100/24, 192.168.1.101/24]",
						},
					},
				},
			},
			"template_id": {
				Type: schema.TypeString,
				Description: "The unique ID of the template, which specifies the number of cores, storage size, and memory. " +
					"You cannot downgrade to a smaller template or minor edition (e.g. from business to playground).",
				Required:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.IsUUID),
			},
			"connection_string": {
				Type:        schema.TypeString,
				Description: "The connection string for your cluster.",
				Computed:    true,
			},
			"credentials": {
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "Credentials for the database user to be created.",
				Required:    true,
				ForceNew:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"username": {
							Type:             schema.TypeString,
							Description:      "the username for the initial mongoDB user.",
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
		},
		Timeouts: &resourceDefaultTimeouts,
	}
}

func resourceDbaasMongoClusterCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(SdkBundle).MongoClient

	cluster := dbaas.SetMongoClusterCreateProperties(d)

	createdCluster, apiResponse, err := client.ClustersApi.ClustersPost(ctx).CreateClusterRequest(*cluster).Execute()
	if err != nil {
		if apiResponse.HttpNotFound() {
			d.SetId("")
			return nil
		}
		diags := diag.FromErr(fmt.Errorf("create error while fetching dbaas mongo cluster %s: %w", d.Id(), err))
		return diags
	}
	d.SetId(*createdCluster.Id)

	_, err = waitForClusterToBeReady(ctx, client, d.Id())
	if err != nil {
		return diag.FromErr(fmt.Errorf("updating %w ", err))
	}

	return resourceDbaasMongoClusterRead(ctx, d, meta)
}

func resourceDbaasMongoClusterUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(SdkBundle).MongoClient
	clusterId := d.Id()
	patchRequest := dbaas.SetMongoClusterPatchProperties(d)

	_, apiResponse, err := client.ClustersApi.ClustersPatch(ctx, clusterId).PatchClusterRequest(*patchRequest).Execute()
	if err != nil {
		if apiResponse.HttpNotFound() {
			d.SetId("")
			return nil
		}
		diags := diag.FromErr(fmt.Errorf("update error while fetching dbaas mongo cluster %s: %w", d.Id(), err))
		return diags
	}

	_, err = waitForClusterToBeReady(ctx, client, d.Id())
	if err != nil {
		return diag.FromErr(fmt.Errorf("creating %w", err))
	}

	return resourceDbaasMongoClusterRead(ctx, d, meta)
}

// waitForClusterToBeReady - keeps retrying until cluster is in 'available' state, or context deadline is reached
func waitForClusterToBeReady(ctx context.Context, client *dbaas.MongoClient, clusterId string) (mongo.ClusterResponse, error) {
	var clusterRequest = mongo.NewClusterResponse()
	err := resource.RetryContext(ctx, *resourceDefaultTimeouts.Update, func() *resource.RetryError {
		var err error
		var apiResponse *mongo.APIResponse
		*clusterRequest, apiResponse, err = client.ClustersApi.ClustersFindById(ctx, clusterId).Execute()
		apiResponse.LogInfo()
		if apiResponse.HttpNotFound() {
			log.Printf("[INFO] Could not find cluster %s retrying...", clusterId)
			return resource.RetryableError(fmt.Errorf("could not find cluster %s, %w, retrying", clusterId, err))
		}
		if err != nil {
			resource.NonRetryableError(err)
		}

		if clusterRequest != nil && clusterRequest.Metadata != nil && !strings.EqualFold(string(*clusterRequest.Metadata.State), utils.Available) {
			log.Printf("[INFO] mongo cluster %s is still in state %s", clusterId, *clusterRequest.Metadata.State)
			return resource.RetryableError(fmt.Errorf("clusterRequest is still in state %s", *clusterRequest.Metadata.State))
		}
		return nil
	})
	if clusterRequest == nil || clusterRequest.Properties == nil || *clusterRequest.Properties.DisplayName == "" {
		return *clusterRequest, fmt.Errorf("could not find cluster %s", clusterId)
	}
	return *clusterRequest, err
}

func resourceDbaasMongoClusterRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(SdkBundle).MongoClient

	cluster, apiResponse, err := client.ClustersApi.ClustersFindById(ctx, d.Id()).Execute()

	if err != nil {
		if apiResponse.HttpNotFound() {
			d.SetId("")
			return nil
		}
		diags := diag.FromErr(fmt.Errorf("read error while fetching dbaas mongo cluster %s: %w", d.Id(), err))
		return diags
	}

	log.Printf("[INFO] Successfully retreived cluster %s: %+v", d.Id(), cluster)

	if err := dbaas.SetDbaasMongoDBClusterData(d, cluster); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceDbaasMongoClusterDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(SdkBundle).MongoClient

	_, apiResponse, err := client.ClustersApi.ClustersDelete(ctx, d.Id()).Execute()

	if err != nil {
		if apiResponse.HttpNotFound() {
			d.SetId("")
			return nil
		}
		diags := diag.FromErr(fmt.Errorf("error while deleting mongo dbaas cluster %s: %w", d.Id(), err))
		return diags
	}

	for {
		log.Printf("[INFO] Waiting for cluster %s to be deleted...", d.Id())

		clusterdDeleted, dsErr := dbaasMongoClusterDeleted(ctx, client, d)

		if dsErr != nil {
			diags := diag.FromErr(fmt.Errorf("error while checking deletion status of mongo dbaas cluster %s: %s", d.Id(), dsErr))
			return diags
		}

		if clusterdDeleted {
			log.Printf("[INFO] Successfully deleted dbaas mongo cluster: %s", d.Id())
			break
		}

		select {
		case <-time.After(utils.SleepInterval):
			log.Printf("[INFO] trying again ...")
		case <-ctx.Done():
			diags := diag.FromErr(fmt.Errorf("dbaas mongo cluster deletion timed out! WARNING: your mongo cluster (%s) "+
				"will still probably be deleted after some time but the terraform state won't reflect that; check your Ionos Cloud account for updates", d.Id()))
			return diags
		}
	}

	// wait 15 seconds after the deletion of the cluster, for the lan to be freed
	time.Sleep(utils.SleepInterval * 3)

	return nil
}

func resourceDbaasMongoClusterImport(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	client := meta.(SdkBundle).MongoClient

	clusterId := d.Id()

	dbaasCluster, apiResponse, err := client.ClustersApi.ClustersFindById(ctx, clusterId).Execute()

	if err != nil {
		if apiResponse.HttpNotFound() {
			d.SetId("")
			return nil, fmt.Errorf("dbaas cluster does not exist %q", clusterId)
		}
		return nil, fmt.Errorf("an error occured while trying to fetch the import of dbaas cluster %q", clusterId)
	}

	log.Printf("[INFO] dbaas cluster found: %+v", dbaasCluster)

	if err := dbaas.SetDbaasMongoDBClusterData(d, dbaasCluster); err != nil {
		return nil, err
	}

	return []*schema.ResourceData{d}, nil
}

func dbaasMongoClusterDeleted(ctx context.Context, client *dbaas.MongoClient, d *schema.ResourceData) (bool, error) {

	_, apiResponse, err := client.ClustersApi.ClustersFindById(ctx, d.Id()).Execute()

	if err != nil {
		if apiResponse.HttpNotFound() {
			return true, nil
		}
		return true, fmt.Errorf("error checking dbaas mongo cluster deletion status: %w", err)
	}
	return false, nil
}
