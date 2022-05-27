package ionoscloud

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	ionoscloud "github.com/ionos-cloud/sdk-go/v6"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourcek8sCluster() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourcek8sClusterCreate,
		ReadContext:   resourcek8sClusterRead,
		UpdateContext: resourcek8sClusterUpdate,
		DeleteContext: resourcek8sClusterDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceK8sClusterImport,
		},
		CustomizeDiff: checkClusterImmutableFields,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Description:  "The desired name for the cluster",
				Required:     true,
				ValidateFunc: validation.All(validation.StringIsNotWhiteSpace),
			},
			"k8s_version": {
				Type:             schema.TypeString,
				Description:      "The desired kubernetes version",
				Optional:         true,
				Computed:         true,
				DiffSuppressFunc: DiffBasedOnVersion,
			},
			"maintenance_window": {
				Type:        schema.TypeList,
				Description: "A maintenance window comprise of a day of the week and a time for maintenance to be allowed",
				Optional:    true,
				MaxItems:    1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"time": {
							Type:         schema.TypeString,
							Description:  "A clock time in the day when maintenance is allowed",
							Required:     true,
							ValidateFunc: validation.All(validation.StringIsNotWhiteSpace),
						},
						"day_of_the_week": {
							Type:         schema.TypeString,
							Description:  "Day of the week when maintenance is allowed",
							Required:     true,
							ValidateFunc: validation.All(validation.StringIsNotWhiteSpace),
						},
					},
				},
			},
			"viable_node_pool_versions": {
				Type:        schema.TypeList,
				Description: "List of versions that may be used for node pools under this cluster",
				Optional:    true,
				Computed:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			//"public": {
			//	Type:        schema.TypeBool,
			//	Description: "The indicator if the cluster is public or private. Be aware that setting it to false is currently in beta phase.",
			//	Optional:    true,
			//	Default:     true,
			//},
			"api_subnet_allow_list": {
				Type: schema.TypeList,
				Description: "Access to the K8s API server is restricted to these CIDRs. Cluster-internal traffic is not " +
					"affected by this restriction. If no allowlist is specified, access is not restricted. If an IP " +
					"without subnet mask is provided, the default value will be used: 32 for IPv4 and 128 for IPv6.",
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"s3_buckets": {
				Type:        schema.TypeList,
				Description: "List of S3 bucket configured for K8s usage. For now it contains only an S3 bucket used to store K8s API audit logs.",
				Optional:    true,
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Description: "Name of the S3 bucket",
							Optional:    true,
						},
					},
				},
			},
		},
		Timeouts: &resourceDefaultTimeouts,
	}
}
func checkClusterImmutableFields(_ context.Context, diff *schema.ResourceDiff, _ interface{}) error {

	//we do not want to check in case of resource creation
	if diff.Id() == "" {
		return nil
	}
	if diff.HasChange("public") {
		return fmt.Errorf("public %s", ImmutableError)
	}
	return nil

}
func resourcek8sClusterCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(SdkBundle).CloudApiClient

	clusterName := d.Get("name").(string)
	cluster := ionoscloud.KubernetesClusterForPost{
		Properties: &ionoscloud.KubernetesClusterPropertiesForPost{
			Name: &clusterName,
		},
	}

	if k8svVal, k8svOk := d.GetOk("k8s_version"); k8svOk {
		log.Printf("[INFO] Setting K8s version to : %s", k8svVal.(string))
		k8svVal := k8svVal.(string)
		cluster.Properties.K8sVersion = &k8svVal
	}

	if _, mwOk := d.GetOk("maintenance_window.0"); mwOk {
		cluster.Properties.MaintenanceWindow = &ionoscloud.KubernetesMaintenanceWindow{}
	}

	if mtVal, mtOk := d.GetOk("maintenance_window.0.time"); mtOk {
		log.Printf("[INFO] Setting Maintenance Window Time to : %s", mtVal.(string))
		mtVal := mtVal.(string)
		cluster.Properties.MaintenanceWindow.Time = &mtVal
	}

	if mdVal, mdOk := d.GetOk("maintenance_window.0.day_of_the_week"); mdOk {
		mdVal := mdVal.(string)
		cluster.Properties.MaintenanceWindow.DayOfTheWeek = &mdVal
	}

	//public := d.Get("public").(bool)
	//cluster.Properties.Public = &public

	if apiSubnet, apiSubnetOk := d.GetOk("api_subnet_allow_list"); apiSubnetOk {
		apiSubnet := apiSubnet.([]interface{})
		if apiSubnet != nil && len(apiSubnet) > 0 {
			apiSubnets := make([]string, 0)
			for _, value := range apiSubnet {
				valueS := value.(string)
				apiSubnets = append(apiSubnets, valueS)
			}
			if len(apiSubnets) > 0 {
				cluster.Properties.ApiSubnetAllowList = &apiSubnets
			}
		}
	}

	if s3Bucket, s3BucketOk := d.GetOk("s3_buckets"); s3BucketOk {
		s3BucketValues := s3Bucket.([]interface{})
		if s3BucketValues != nil && len(s3BucketValues) > 0 {
			var s3Buckets []ionoscloud.S3Bucket
			for index := range s3BucketValues {
				var s3Bucket ionoscloud.S3Bucket
				addBucket := false
				if name, nameOk := d.GetOk(fmt.Sprintf("s3_buckets.%d.name", index)); nameOk {
					name := name.(string)
					s3Bucket.Name = &name
					addBucket = true
				} else {
					diags := diag.FromErr(fmt.Errorf("name must be provided for s3 bucket"))
					return diags
				}
				if addBucket {
					s3Buckets = append(s3Buckets, s3Bucket)
				}
			}
			if len(s3Buckets) > 0 {
				cluster.Properties.S3Buckets = &s3Buckets
			}
		}
	}

	createdCluster, apiResponse, err := client.KubernetesApi.K8sPost(ctx).KubernetesCluster(cluster).Execute()
	logApiRequestTime(apiResponse)

	if err != nil {
		d.SetId("")
		diags := diag.FromErr(fmt.Errorf("error creating k8s cluster: %s", err))
		return diags
	}

	d.SetId(*createdCluster.Id)
	log.Printf("[INFO] Created k8s cluster: %s", d.Id())

	for {
		log.Printf("[INFO] Waiting for cluster %s to be ready...", d.Id())

		clusterReady, rsErr := k8sClusterReady(ctx, client, d)

		if rsErr != nil {
			diags := diag.FromErr(fmt.Errorf("error while checking readiness status of k8s cluster %s: %s", d.Id(), rsErr))
			return diags
		}

		if clusterReady {
			log.Printf("[INFO] k8s cluster ready: %s", d.Id())
			break
		}

		select {
		case <-time.After(SleepInterval):
			log.Printf("[INFO] trying again ...")
		case <-ctx.Done():
			log.Printf("[INFO] create timed out")
			diags := diag.FromErr(fmt.Errorf("k8s cluster creation timed out! WARNING: your k8s cluster will still probably be created after some time but the terraform state wont reflect that; check your Ionos Cloud account for updates"))
			return diags
		}

	}

	return resourcek8sClusterRead(ctx, d, meta)
}

func resourcek8sClusterRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(SdkBundle).CloudApiClient

	cluster, apiResponse, err := client.KubernetesApi.K8sFindByClusterId(ctx, d.Id()).Execute()
	logApiRequestTime(apiResponse)

	if err != nil {
		if apiResponse != nil && apiResponse.Response != nil && apiResponse.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		diags := diag.FromErr(fmt.Errorf("error while fetching k8s cluster %s: %s", d.Id(), err))
		return diags
	}

	log.Printf("[INFO] Successfully retreived cluster %s: %+v", d.Id(), cluster)

	if err := setK8sClusterData(d, &cluster); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourcek8sClusterUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	client := meta.(SdkBundle).CloudApiClient

	request := ionoscloud.KubernetesClusterForPut{}

	clusterName := d.Get("name").(string)
	request.Properties = &ionoscloud.KubernetesClusterPropertiesForPut{
		Name: &clusterName,
	}

	if d.HasChange("name") {
		oldName, newName := d.GetChange("name")
		log.Printf("[INFO] k8s cluster name changed from %+v to %+v", oldName, newName)
		newNameStr := newName.(string)
		request.Properties.Name = &newNameStr
	}

	log.Printf("[INFO] Attempting update cluster Id %s", d.Id())

	if d.HasChange("k8s_version") {
		oldk8sVersion, newk8sVersion := d.GetChange("k8s_version")
		log.Printf("[INFO] k8s version changed from %+v to %+v", oldk8sVersion, newk8sVersion)
		newk8sVersionStr := newk8sVersion.(string)
		if newk8sVersion != nil {
			request.Properties.K8sVersion = &newk8sVersionStr
		}
	}

	if d.HasChange("maintenance_window.0") {

		_, newMw := d.GetChange("maintenance_window.0")

		if newMw.(map[string]interface{}) != nil {

			updateMaintenanceWindow := false
			dayofTheWeek := d.Get("maintenance_window.0.day_of_the_week").(string)
			winTime := d.Get("maintenance_window.0.time").(string)
			maintenanceWindow := &ionoscloud.KubernetesMaintenanceWindow{
				DayOfTheWeek: &dayofTheWeek,
				Time:         &winTime,
			}

			if d.HasChange("maintenance_window.0.day_of_the_week") {
				oldMd, newMd := d.GetChange("maintenance_window.0.day_of_the_week")
				if newMd.(string) != "" {
					log.Printf("[INFO] k8s maintenance window DOW changed from %+v to %+v", oldMd, newMd)
					updateMaintenanceWindow = true
					newMd := newMd.(string)
					maintenanceWindow.DayOfTheWeek = &newMd
				}
			}

			if d.HasChange("maintenance_window.0.time") {

				oldMt, newMt := d.GetChange("maintenance_window.0.time")
				if newMt.(string) != "" {
					log.Printf("[INFO] k8s maintenance window time changed from %+v to %+v", oldMt, newMt)
					updateMaintenanceWindow = true
					newMt := newMt.(string)
					maintenanceWindow.Time = &newMt
				}
			}

			if updateMaintenanceWindow == true {
				request.Properties.MaintenanceWindow = maintenanceWindow
			}
		}
	}

	if d.HasChange("api_subnet_allow_list") {
		_, newApiSubnet := d.GetChange("api_subnet_allow_list")
		apiSubnet := newApiSubnet.([]interface{})
		apiSubnets := make([]string, 0)
		if apiSubnet != nil && len(apiSubnet) > 0 {
			for _, value := range apiSubnet {
				valueS := value.(string)
				apiSubnets = append(apiSubnets, valueS)
			}
		}
		request.Properties.ApiSubnetAllowList = &apiSubnets
	}

	if d.HasChange("s3_buckets") {
		_, newS3Buckets := d.GetChange("s3_buckets")
		s3BucketValues := newS3Buckets.([]interface{})
		s3Buckets := make([]ionoscloud.S3Bucket, 0)
		if s3BucketValues != nil && len(s3BucketValues) > 0 {
			for index := range s3BucketValues {
				var s3Bucket ionoscloud.S3Bucket
				addBucket := false
				if name, nameOk := d.GetOk(fmt.Sprintf("s3_buckets.%d.name", index)); nameOk {
					name := name.(string)
					s3Bucket.Name = &name
					addBucket = true
				}
				if addBucket {
					s3Buckets = append(s3Buckets, s3Bucket)
				}
			}
		}
		request.Properties.S3Buckets = &s3Buckets
	}

	_, apiResponse, err := client.KubernetesApi.K8sPut(ctx, d.Id()).KubernetesCluster(request).Execute()
	logApiRequestTime(apiResponse)

	if err != nil {
		if apiResponse != nil && apiResponse.Response != nil && apiResponse.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		diags := diag.FromErr(fmt.Errorf("error while updating k8s cluster: %s", err))
		return diags
	}

	for {
		log.Printf("[INFO] Waiting for cluster %s to be ready...", d.Id())

		clusterReady, rsErr := k8sClusterReady(ctx, client, d)

		if rsErr != nil {
			diags := diag.FromErr(fmt.Errorf("error while checking readiness status of k8s cluster %s: %s", d.Id(), rsErr))
			return diags
		}

		if clusterReady {
			log.Printf("[INFO] k8s cluster ready: %s", d.Id())
			break
		}

		select {
		case <-time.After(SleepInterval):
			log.Printf("[INFO] trying again ...")
		case <-ctx.Done():
			diags := diag.FromErr(fmt.Errorf("k8s cluster update timed out! WARNING: your k8s cluster will still probably be created after some time but the terraform state won't reflect that; check your Ionos Cloud account for updates"))
			return diags
		}

	}

	return resourcek8sClusterRead(ctx, d, meta)
}

func resourcek8sClusterDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	client := meta.(SdkBundle).CloudApiClient

	apiResponse, err := client.KubernetesApi.K8sDelete(ctx, d.Id()).Execute()
	logApiRequestTime(apiResponse)

	if err != nil {
		if apiResponse != nil && apiResponse.Response != nil && apiResponse.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		diags := diag.FromErr(fmt.Errorf("error while deleting k8s cluster %s: %s", d.Id(), err))
		return diags
	}

	for {
		log.Printf("[INFO] Waiting for cluster %s to be deleted...", d.Id())

		clusterDeleted, dsErr := k8sClusterDeleted(ctx, client, d)

		if dsErr != nil {
			diags := diag.FromErr(fmt.Errorf("error while checking deletion status of k8s cluster %s: %w", d.Id(), dsErr))
			return diags
		}

		if clusterDeleted {
			log.Printf("[INFO] Successfully deleted k8s cluster: %s", d.Id())
			break
		}

		select {
		case <-time.After(SleepInterval):
			log.Printf("[INFO] trying again ...")
		case <-ctx.Done():
			diags := diag.FromErr(fmt.Errorf("k8s cluster deletion timed out! WARNING: your k8s cluster will still probably be deleted after some time but the terraform state won't reflect that; check your Ionos Cloud account for updates"))
			return diags
		}
	}

	d.SetId("")

	return nil
}

func resourceK8sClusterImport(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	client := meta.(SdkBundle).CloudApiClient

	clusterId := d.Id()

	cluster, apiResponse, err := client.KubernetesApi.K8sFindByClusterId(ctx, clusterId).Execute()
	logApiRequestTime(apiResponse)

	if err != nil {
		if apiResponse != nil && apiResponse.Response != nil && apiResponse.StatusCode == 404 {
			d.SetId("")
			return nil, fmt.Errorf("unable to find k8s cluster %q", clusterId)
		}
		return nil, fmt.Errorf("unable to retreive k8s cluster %q", d.Id())
	}

	log.Printf("[INFO] K8s cluster found: %+v", cluster)

	if err := setK8sClusterData(d, &cluster); err != nil {
		return nil, err
	}

	return []*schema.ResourceData{d}, nil
}

func setK8sClusterData(d *schema.ResourceData, cluster *ionoscloud.KubernetesCluster) error {

	if cluster.Id != nil {
		d.SetId(*cluster.Id)
	}

	if cluster.Properties != nil {
		if cluster.Properties.Name != nil {
			if err := d.Set("name", *cluster.Properties.Name); err != nil {
				return err
			}
		}

		if cluster.Properties.K8sVersion != nil {
			if err := d.Set("k8s_version", *cluster.Properties.K8sVersion); err != nil {
				return err
			}

		}

		if cluster.Properties.MaintenanceWindow != nil && cluster.Properties.MaintenanceWindow.Time != nil && cluster.Properties.MaintenanceWindow.DayOfTheWeek != nil {
			if err := d.Set("maintenance_window", []map[string]string{
				{
					"time":            *cluster.Properties.MaintenanceWindow.Time,
					"day_of_the_week": *cluster.Properties.MaintenanceWindow.DayOfTheWeek,
				},
			}); err != nil {
				return err
			}
		}

		if cluster.Properties.ViableNodePoolVersions != nil && len(*cluster.Properties.ViableNodePoolVersions) > 0 {
			var viableNodePoolVersions []interface{}
			for _, viableNodePoolVersion := range *cluster.Properties.ViableNodePoolVersions {
				viableNodePoolVersions = append(viableNodePoolVersions, viableNodePoolVersion)
			}
			if err := d.Set("viable_node_pool_versions", viableNodePoolVersions); err != nil {
				return err
			}
		}

		//if cluster.Properties.Public != nil {
		//	err := d.Set("public", *cluster.Properties.Public)
		//	if err != nil {
		//		return fmt.Errorf("error while setting public property for cluser %s: %s", d.Id(), err)
		//	}
		//}

		if cluster.Properties.ApiSubnetAllowList != nil {
			apiSubnetAllowLists := make([]interface{}, len(*cluster.Properties.ApiSubnetAllowList), len(*cluster.Properties.ApiSubnetAllowList))
			for i, apiSubnetAllowList := range *cluster.Properties.ApiSubnetAllowList {
				apiSubnetAllowLists[i] = apiSubnetAllowList
			}
			if err := d.Set("api_subnet_allow_list", apiSubnetAllowLists); err != nil {
				return fmt.Errorf("error while setting api_subnet_allow_list property for cluser %s: %s", d.Id(), err)
			}
		}

		if cluster.Properties.S3Buckets != nil {
			s3Buckets := make([]interface{}, len(*cluster.Properties.S3Buckets), len(*cluster.Properties.S3Buckets))
			for i, s3Bucket := range *cluster.Properties.S3Buckets {
				s3BucketEntry := make(map[string]interface{})
				s3BucketEntry["name"] = *s3Bucket.Name
				s3Buckets[i] = s3BucketEntry
			}
			if err := d.Set("s3_buckets", s3Buckets); err != nil {
				return fmt.Errorf("error while setting s3_buckets property for cluser %s: %s", d.Id(), err)
			}
		}

	}

	return nil
}

func k8sClusterReady(ctx context.Context, client *ionoscloud.APIClient, d *schema.ResourceData) (bool, error) {
	subjectCluster, apiResponse, err := client.KubernetesApi.K8sFindByClusterId(ctx, d.Id()).Execute()
	logApiRequestTime(apiResponse)

	if err != nil {
		return true, fmt.Errorf("error checking k8s cluster status: %s", err)
	}
	return *subjectCluster.Metadata.State == "ACTIVE", nil
}

func k8sClusterDeleted(ctx context.Context, client *ionoscloud.APIClient, d *schema.ResourceData) (bool, error) {

	_, apiResponse, err := client.KubernetesApi.K8sFindByClusterId(ctx, d.Id()).Execute()
	logApiRequestTime(apiResponse)

	if err != nil {
		if apiResponse != nil && apiResponse.Response != nil && apiResponse.StatusCode == 404 {
			return true, nil
		}
		return true, fmt.Errorf("error checking k8s cluster deletion status: %w", err)
	}
	return false, nil
}
