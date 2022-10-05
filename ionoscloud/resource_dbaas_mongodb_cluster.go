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
	"strings"
	"time"
)

func resourceDbaasMongoDBCluster() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDbaasMongoClusterCreate,
		ReadContext:   resourceDbaasMongoClusterRead,
		//no update operation, forcenew on all fields
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
				ForceNew:    true,
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
				ForceNew:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.IntBetween(1, 5)),
			},
			"display_name": {
				Type:        schema.TypeString,
				Description: "The name of your cluster.",
				Required:    true,
				ForceNew:    true,
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
				Description: "Details about the network connection for your cluster.",
				Required:    true,
				ForceNew:    true,
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
							Type: schema.TypeList,
							Description: "The list of IPs and subnet for your cluster.\n          Note the following unavailable IP ranges:\n          10.233.64.0/18\n          10.233.0.0/18\n          10.233.114.0/24 		\n example: [192.168.1.100/24, 192.168.1.101/24]",
							Required: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
			},
			"template_id": {
				Type:             schema.TypeString,
				Description:      "The unique ID of the template, which specifies the number of cores, storage size, and memory.",
				Required:         true,
				ForceNew:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.IsUUID),
			},
			"connection_string": {
				Type:        schema.TypeString,
				Description: "The connection string for your cluster.",
				Computed:    true,
				ForceNew:    true,
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

	cluster := dbaasService.GetDbaasMongoClusterDataCreate(d)

	createdCluster, apiResponse, err := client.ClustersApi.ClustersPost(ctx).CreateClusterRequest(*cluster).Execute()
	if err != nil {
		if apiResponse.HttpNotFound() {
			d.SetId("")
			return nil
		}
		diags := diag.FromErr(fmt.Errorf("error while fetching dbaas mongo cluster %s: %w", d.Id(), err))
		return diags
	}
	d.SetId(*createdCluster.Id)

	for {
		log.Printf("[INFO] Waiting for mongo cluster %s to be ready...", d.Id())

		clusterReady, rsErr := dbaasMongoClusterReady(ctx, client, d)
		if rsErr != nil {
			diags := diag.FromErr(fmt.Errorf("error while checking readiness status of dbaas mongo cluster %s: %w", d.Id(), rsErr))
			return diags
		}

		if clusterReady {
			log.Printf("[INFO] dbaas mongo cluster ready: %s", d.Id())
			break
		}

		select {
		case <-time.After(utils.SleepInterval):
			log.Printf("[INFO] trying again ...")
		case <-ctx.Done():
			log.Printf("[INFO] create timed out")
			diags := diag.FromErr(fmt.Errorf("dbaas mongo cluster creation timed out! WARNING: your dbaas cluster (%s) will still probably be created after some time but the terraform state wont reflect that; check your Ionos Cloud account for updates", d.Id()))
			return diags
		}

	}

	return resourceDbaasMongoClusterRead(ctx, d, meta)
}

func resourceDbaasMongoClusterRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(SdkBundle).MongoClient

	cluster, apiResponse, err := client.ClustersApi.ClustersFindById(ctx, d.Id()).Execute()

	if err != nil {
		if apiResponse.HttpNotFound() {
			d.SetId("")
			return nil
		}
		diags := diag.FromErr(fmt.Errorf("error while fetching dbaas mongo cluster %s: %w", d.Id(), err))
		return diags
	}

	log.Printf("[INFO] Successfully retreived cluster %s: %+v", d.Id(), cluster)

	if err := dbaasService.SetDbaasMongoDBClusterData(d, cluster); err != nil {
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
		diags := diag.FromErr(fmt.Errorf("error while deleting mongo dbaas cluster %s: %s", d.Id(), err))
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
			diags := diag.FromErr(fmt.Errorf("dbaas mongo cluster deletion timed out! WARNING: your mongo cluster (%s) will still probably be deleted after some time but the terraform state won't reflect that; check your Ionos Cloud account for updates", d.Id()))
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

	if err := dbaasService.SetDbaasMongoDBClusterData(d, dbaasCluster); err != nil {
		return nil, err
	}

	return []*schema.ResourceData{d}, nil
}

func dbaasMongoClusterReady(ctx context.Context, client *dbaasService.MongoClient, d *schema.ResourceData) (bool, error) {
	subjectCluster, _, err := client.ClustersApi.ClustersFindById(ctx, d.Id()).Execute()

	if err != nil {
		return true, fmt.Errorf("error checking dbaas mongo cluster status: %w", err)
	}
	return strings.EqualFold(string(*subjectCluster.Metadata.State), utils.Available), nil
}

func dbaasMongoClusterDeleted(ctx context.Context, client *dbaasService.MongoClient, d *schema.ResourceData) (bool, error) {

	_, apiResponse, err := client.ClustersApi.ClustersFindById(ctx, d.Id()).Execute()

	if err != nil {
		if apiResponse.HttpNotFound() {
			return true, nil
		}
		return true, fmt.Errorf("error checking dbaas mongo cluster deletion status: %s", err)
	}
	return false, nil
}
