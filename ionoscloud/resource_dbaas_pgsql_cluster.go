package ionoscloud

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	dbaasService "github.com/ionos-cloud/terraform-provider-ionoscloud/services/dbaas"
	"log"
	"net"
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
		Schema: map[string]*schema.Schema{
			"postgres_version": {
				Type:         schema.TypeString,
				Description:  "The PostgreSQL version of your cluster.",
				Required:     true,
				ValidateFunc: validation.All(validation.StringIsNotWhiteSpace),
			},
			"replicas": {
				Type:        schema.TypeInt,
				Description: "The number of replicas in your cluster.",
				Required:    true,
				ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
					v := val.(int)
					if v < 1 || v > 5 {
						errs = append(errs, fmt.Errorf("%q should have a value between 1 and 5; got: %v", key, v))
					}
					return
				},
			},
			"cpu_core_count": {
				Type:        schema.TypeInt,
				Description: "The number of CPU cores per replica.",
				Required:    true,
			},
			"ram_size": {
				Type:         schema.TypeString,
				Description:  "The amount of memory per replica.",
				Required:     true,
				ValidateFunc: validation.All(validation.StringIsNotWhiteSpace),
			},
			"storage_size": {
				Type:         schema.TypeString,
				Description:  "The amount of storage per replica.",
				Required:     true,
				ValidateFunc: validation.All(validation.StringIsNotWhiteSpace),
			},
			"storage_type": {
				Type:         schema.TypeString,
				Description:  "The storage type used in your cluster.",
				Required:     true,
				ValidateFunc: validation.All(validation.StringIsNotWhiteSpace),
			},
			"vdc_connections": {
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "The VDC to connect to your cluster.",
				Required:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"vdc_id": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.All(validation.StringIsNotWhiteSpace),
						},
						"lan_id": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.All(validation.StringIsNotWhiteSpace),
						},
						"ip_address": {
							Type:        schema.TypeString,
							Description: "The IP and subnet for the database.\n          Note the following unavailable IP ranges:\n          10.233.64.0/18\n          10.233.0.0/18\n          10.233.114.0/24",
							Optional:    true,
							ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
								v := val.(string)
								unavailableNetworks := []string{"10.233.64.0/18", "10.233.0.0/18", "10.233.114.0/24"}

								ip, _, _ := net.ParseCIDR(v)

								for _, unavailableNetwork := range unavailableNetworks {
									if _, network, _ := net.ParseCIDR(unavailableNetwork); network.Contains(ip) {
										errs = append(errs, fmt.Errorf("for %q the following IP ranges are unavailable: 10.233.64.0/18, 10.233.0.0/18, 10.233.114.0/24; got: %v", key, v))
									}
								}
								return
							},
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
			"display_name": {
				Type:        schema.TypeString,
				Description: "The friendly name of your cluster.",
				Required:    true,
			},
			"backup_enabled": {
				Type:       schema.TypeBool,
				Deprecated: "Deprecated: backup is always enabled.\n      Enables automatic backups of your cluster.",
				Optional:   true,
				Default:    true,
			},
			"maintenance_window": {
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "a weekly 4 hour-long window, during which maintenance might occur",
				Optional:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"time": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.All(validation.StringIsNotWhiteSpace),
						},
						"weekday": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.All(validation.StringIsNotWhiteSpace),
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
							Description:  "the username for the initial postgres user. some system usernames\n          are restricted (e.g. \"postgres\", \"admin\", \"standby\")",
							Required:     true,
							ValidateFunc: validation.All(validation.StringIsNotWhiteSpace),
						},
						"password": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.All(validation.StringIsNotWhiteSpace),
						},
					},
				},
			},
		},
		Timeouts: &resourceDefaultTimeouts,
	}
}

func resourceDbaasPgSqlClusterCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(SdkBundle).DbaasClient

	dbaasCluster := dbaasService.GetDbaasPgSqlClusterDataCreate(d)

	dbaasClusterResponse, _, err := client.CreateCluster(ctx, *dbaasCluster)

	if err != nil {
		diags := diag.FromErr(fmt.Errorf("an error occured while creating a dbaas cluster: %s", err))
		return diags
	}

	d.SetId(*dbaasClusterResponse.Id)

	for {
		log.Printf("[INFO] Waiting for cluster %s to be ready...", d.Id())

		clusterReady, rsErr := dbaasClusterReady(ctx, client, d)

		if rsErr != nil {
			diags := diag.FromErr(fmt.Errorf("error while checking readiness status of dbaas cluster %s: %s", d.Id(), rsErr))
			return diags
		}

		if clusterReady {
			log.Printf("[INFO] dbaas cluster ready: %s", d.Id())
			break
		}

		select {
		case <-time.After(SleepInterval):
			log.Printf("[INFO] trying again ...")
		case <-ctx.Done():
			log.Printf("[INFO] create timed out")
			diags := diag.FromErr(fmt.Errorf("dbaas cluster creation timed out! WARNING: your dbaas cluster will still probably be created after some time but the terraform state wont reflect that; check your Ionos Cloud account for updates"))
			return diags
		}

	}

	return resourceDbaasPgSqlClusterRead(ctx, d, meta)
}

func resourceDbaasPgSqlClusterRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	client := meta.(SdkBundle).DbaasClient

	cluster, apiResponse, err := client.GetCluster(ctx, d.Id())

	if err != nil {
		if apiResponse != nil && apiResponse.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		diags := diag.FromErr(fmt.Errorf("error while fetching dbaas cluster %s: %s", d.Id(), err))
		return diags
	}

	log.Printf("[INFO] Successfully retreived cluster %s: %+v", d.Id(), cluster)

	dbaasService.SetDbaasPgSqlClusterData(d, cluster)

	return nil
}

func resourceDbaasPgSqlClusterUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(SdkBundle).DbaasClient

	cluster, diags := dbaasService.GetDbaasPgSqlClusterDataUpdate(d)

	if diags != nil {
		return diags
	}

	dbaasClusterResponse, _, err := client.ClustersApi.ClustersPatch(ctx, d.Id()).Cluster(*cluster).Execute()

	if err != nil {
		diags := diag.FromErr(fmt.Errorf("an error occured while updating a dbaas cluster: %s", err))
		return diags
	}

	d.SetId(*dbaasClusterResponse.Id)

	for {
		log.Printf("[INFO] Waiting for cluster %s to be ready...", d.Id())

		clusterReady, rsErr := dbaasClusterReady(ctx, client, d)

		if rsErr != nil {
			diags := diag.FromErr(fmt.Errorf("error while checking readiness status of dbaas cluster %s: %s", d.Id(), rsErr))
			return diags
		}

		if clusterReady {
			log.Printf("[INFO] dbaas cluster ready: %s", d.Id())
			break
		}

		select {
		case <-time.After(SleepInterval):
			log.Printf("[INFO] trying again ...")
		case <-ctx.Done():
			log.Printf("[INFO] create timed out")
			diags := diag.FromErr(fmt.Errorf("dbaas cluster update timed out! WARNING: your dbaas cluster will still probably be updated after some time but the terraform state wont reflect that; check your Ionos Cloud account for updates"))
			return diags
		}

	}

	return resourceDbaasPgSqlClusterRead(ctx, d, meta)
}

func resourceDbaasPgSqlClusterDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(SdkBundle).DbaasClient

	_, apiResponse, err := client.DeleteCluster(ctx, d.Id())

	if err != nil {
		if apiResponse != nil && apiResponse.StatusCode == 404 {
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
		case <-time.After(SleepInterval):
			log.Printf("[INFO] trying again ...")
		case <-ctx.Done():
			diags := diag.FromErr(fmt.Errorf("dbaas cluster deletion timed out! WARNING: your k8s cluster will still probably be deleted after some time but the terraform state won't reflect that; check your Ionos Cloud account for updates"))
			return diags
		}
	}

	return nil
}

func resourceDbaasPgSqlClusterImport(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	client := meta.(SdkBundle).DbaasClient

	clusterId := d.Id()

	dbaasCluster, apiResponse, err := client.GetCluster(ctx, clusterId)

	if err != nil {
		if apiResponse != nil && apiResponse.StatusCode == 404 {
			d.SetId("")
			return nil, fmt.Errorf("an error occured while trying to fetch the import of dbaas cluster %q", clusterId)
		}
		return nil, fmt.Errorf("dbaas cluster does not exist %q", clusterId)
	}

	log.Printf("[INFO] dbaas cluster found: %+v", dbaasCluster)

	if dbaasCluster.Id != nil {
		if err := d.Set("id", *dbaasCluster.Id); err != nil {
			return nil, err
		}
	}
	dbaasService.SetDbaasPgSqlClusterData(d, dbaasCluster)

	return []*schema.ResourceData{d}, nil
}

func dbaasClusterReady(ctx context.Context, client *dbaasService.Client, d *schema.ResourceData) (bool, error) {
	subjectCluster, _, err := client.GetCluster(ctx, d.Id())

	if err != nil {
		return true, fmt.Errorf("error checking dbaas cluster status: %s", err)
	}

	//if *subjectCluster.LifecycleStatus == "FAILED" {
	//	return false, fmt.Errorf("dbaas cluster has failed")
	//}
	return *subjectCluster.LifecycleStatus == "AVAILABLE", nil
}

func dbaasClusterDeleted(ctx context.Context, client *dbaasService.Client, d *schema.ResourceData) (bool, error) {

	_, apiResponse, err := client.GetCluster(ctx, d.Id())

	if err != nil {
		if apiResponse != nil && apiResponse.StatusCode == 404 {
			return true, nil
		}
		return true, fmt.Errorf("error checking dbaas cluster deletion status: %s", err)
	}
	return false, nil
}
