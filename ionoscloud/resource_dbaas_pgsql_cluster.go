package ionoscloud

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	dbaas "github.com/ionos-cloud/sdk-go-autoscaling"
	"log"
	"net"
	"time"
)

func resourceDbaasCluster() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDbaasClusterCreate,
		ReadContext:   resourceDbaasClusterRead,
		UpdateContext: resourceDbaasClusterUpdate,
		DeleteContext: resourceDbaasClusterDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceDbaasClusterImport,
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

func resourceDbaasClusterCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(SdkBundle).DbaasClient

	dbaasCluster := dbaas.CreateClusterRequest{}

	if postgresVersion, ok := d.GetOk("postgres_version"); ok {
		postgresVersion := postgresVersion.(string)
		dbaasCluster.PostgresVersion = &postgresVersion
	}

	if replicas, ok := d.GetOk("replicas"); ok {
		replicas := float32(replicas.(int))
		dbaasCluster.Replicas = &replicas
	}

	if cpuCoreCount, ok := d.GetOk("cpu_core_count"); ok {
		cpuCoreCount := float32(cpuCoreCount.(int))
		dbaasCluster.CpuCoreCount = &cpuCoreCount
	}

	if ramSize, ok := d.GetOk("ram_size"); ok {
		ramSize := ramSize.(string)
		dbaasCluster.RamSize = &ramSize
	}

	if storageSize, ok := d.GetOk("storage_size"); ok {
		storageSize := storageSize.(string)
		dbaasCluster.StorageSize = &storageSize
	}

	if storageType, ok := d.GetOk("storage_type"); ok {
		storageType := dbaas.StorageType(storageType.(string))
		dbaasCluster.StorageType = &storageType
	}

	if vdcConnection, ok := d.GetOk("vdc_connections"); ok {
		if vdcConnection.([]interface{}) != nil {

			var vdcConnections []dbaas.VDCConnection

			for vdcIndex := range vdcConnection.([]interface{}) {

				connection := dbaas.VDCConnection{}

				if vdcId, ok := d.GetOk(fmt.Sprintf("vdc_connections.%d.vdc_id", vdcIndex)); ok {
					vdcId := vdcId.(string)
					connection.VdcId = &vdcId
				}

				if lanId, ok := d.GetOk(fmt.Sprintf("vdc_connections.%d.lan_id", vdcIndex)); ok {
					lanId := lanId.(string)
					connection.LanId = &lanId
				}

				if ipAddress, ok := d.GetOk(fmt.Sprintf("vdc_connections.%d.ip_address", vdcIndex)); ok {
					ipAddress := ipAddress.(string)
					connection.IpAddress = &ipAddress
				}

				vdcConnections = append(vdcConnections, connection)

			}

			if len(vdcConnections) > 0 {
				dbaasCluster.VdcConnections = &vdcConnections
			}
		}
	}

	if location, ok := d.GetOk("location"); ok {
		location := location.(string)
		dbaasCluster.Location = &location
	}

	if displayName, ok := d.GetOk("display_name"); ok {
		displayName := displayName.(string)
		dbaasCluster.DisplayName = &displayName
	}

	backupEnabled := d.Get("backup_enabled").(bool)
	dbaasCluster.BackupEnabled = &backupEnabled

	if _, ok := d.GetOk("maintenance_window.0"); ok {
		dbaasCluster.MaintenanceWindow = &dbaas.MaintenanceWindow{}

		if timeV, ok := d.GetOk("maintenance_window.0.time"); ok {
			timeV := timeV.(string)
			dbaasCluster.MaintenanceWindow.Time = &timeV
		}

		if weekday, ok := d.GetOk("maintenance_window.0.weekday"); ok {
			weekday := weekday.(string)
			dbaasCluster.MaintenanceWindow.Weekday = &weekday
		}
	}

	if _, ok := d.GetOk("credentials.0"); ok {
		dbaasCluster.Credentials = &dbaas.DBUser{}
		if username, ok := d.GetOk("credentials.0.username"); ok {
			username := username.(string)
			dbaasCluster.Credentials.Username = &username
		}

		if password, ok := d.GetOk("credentials.0.password"); ok {
			password := password.(string)
			dbaasCluster.Credentials.Password = &password
		}
	}

	dbaasClusterResponse, _, err := client.ClustersApi.ClustersPost(ctx).Cluster(dbaasCluster).Execute()

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

	return resourceDbaasClusterRead(ctx, d, meta)
}

func resourceDbaasClusterRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	client := meta.(SdkBundle).DbaasClient

	cluster, apiResponse, err := client.ClustersApi.ClustersFindById(ctx, d.Id()).Execute()

	if err != nil {
		if apiResponse != nil && apiResponse.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		diags := diag.FromErr(fmt.Errorf("error while fetching dbaas cluster %s: %s", d.Id(), err))
		return diags
	}

	log.Printf("[INFO] Successfully retreived cluster %s: %+v", d.Id(), cluster)

	setDbaasClusterData(d, &cluster)

	return nil
}

func resourceDbaasClusterUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	fmt.Printf("\n________________UPDATE________________\n")

	client := meta.(SdkBundle).DbaasClient

	cluster := dbaas.PatchClusterRequest{}

	if d.HasChange("postgres_version") {
		_, v := d.GetChange("postgres_version")
		vStr := v.(string)
		cluster.PostgresVersion = &vStr
	}

	if d.HasChange("replicas") {
		diags := diag.FromErr(fmt.Errorf("replicas parameter is immutable"))
		return diags
	}

	if d.HasChange("cpu_core_count") {
		diags := diag.FromErr(fmt.Errorf("cpu_core_count parameter is immutable"))
		return diags
	}

	if d.HasChange("ram_size") {
		diags := diag.FromErr(fmt.Errorf("ram_size parameter is immutable"))
		return diags
	}

	if d.HasChange("storage_size") {
		_, v := d.GetChange("storage_size")
		vStr := v.(string)
		cluster.StorageSize = &vStr
	}

	if d.HasChange("vdc_connections") {
		diags := diag.FromErr(fmt.Errorf("vdc_connections parameter is immutable"))
		return diags
	}

	if d.HasChange("display_name") {
		_, v := d.GetChange("display_name")
		vStr := v.(string)
		cluster.DisplayName = &vStr
	}

	if d.HasChange("backup_enabled") {
		vBool := d.Get("backup_enabled").(bool)
		cluster.BackupEnabled = &vBool
	}

	if d.HasChange("maintenance_window") {
		cluster.MaintenanceWindow = &dbaas.MaintenanceWindow{}

		if timeV, ok := d.GetOk("maintenance_window.0.time"); ok {
			timeV := timeV.(string)
			cluster.MaintenanceWindow.Time = &timeV
		}

		if weekday, ok := d.GetOk("maintenance_window.0.weekday"); ok {
			weekday := weekday.(string)
			cluster.MaintenanceWindow.Weekday = &weekday
		}
	}

	dbaasClusterResponse, _, err := client.ClustersApi.ClustersPatch(ctx, d.Id()).Cluster(cluster).Execute()

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

	return resourceDbaasClusterRead(ctx, d, meta)
}

func resourceDbaasClusterDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(SdkBundle).DbaasClient

	_, apiResponse, err := client.ClustersApi.ClustersDelete(ctx, d.Id()).Execute()

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

	d.SetId("")

	return nil
}

func resourceDbaasClusterImport(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	client := meta.(SdkBundle).DbaasClient

	clusterId := d.Id()

	dbaasCluster, apiResponse, err := client.ClustersApi.ClustersFindById(ctx, clusterId).Execute()

	if err != nil {
		if apiResponse != nil && apiResponse.StatusCode == 404 {
			d.SetId("")
			return nil, fmt.Errorf("an error occured while trying to fetch the import of dbaas cluster %q", clusterId)
		}
		return nil, fmt.Errorf("dbaas cluster does not exist %q", clusterId)
	}

	log.Printf("[INFO] dbaas cluster found: %+v", dbaasCluster)

	setDbaasClusterData(d, &dbaasCluster)

	return []*schema.ResourceData{d}, nil
}

func dbaasClusterReady(ctx context.Context, client *dbaas.APIClient, d *schema.ResourceData) (bool, error) {
	subjectCluster, _, err := client.ClustersApi.ClustersFindById(ctx, d.Id()).Execute()

	if err != nil {
		return true, fmt.Errorf("error checking dbaas cluster status: %s", err)
	}

	if *subjectCluster.LifecycleStatus == "FAILED" {
		return false, fmt.Errorf("dbaas cluster has failed")
	}
	return *subjectCluster.LifecycleStatus == "AVAILABLE", nil
}

func dbaasClusterDeleted(ctx context.Context, client *dbaas.APIClient, d *schema.ResourceData) (bool, error) {

	_, apiResponse, err := client.ClustersApi.ClustersFindById(ctx, d.Id()).Execute()

	if err != nil {
		if apiResponse != nil && apiResponse.StatusCode == 404 {
			return true, nil
		}
		return true, fmt.Errorf("error checking dbaas cluster deletion status: %s", err)
	}
	return false, nil
}
