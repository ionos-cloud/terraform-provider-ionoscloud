package ionoscloud

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	ionoscloud "github.com/ionos-cloud/sdk-go/v5"
	"log"
	"strings"
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
		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Description:  "The desired name for the cluster",
				Required:     true,
				ValidateFunc: validation.All(validation.StringIsNotWhiteSpace),
			},
			"k8s_version": {
				Type:        schema.TypeString,
				Description: "The desired kubernetes version",
				Optional:    true,
				Computed:    true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					var oldMajor, oldMinor string
					if old != "" {
						oldSplit := strings.Split(old, ".")
						oldMajor = oldSplit[0]
						oldMinor = oldSplit[1]

						newSplit := strings.Split(new, ".")
						newMajor := newSplit[0]
						newMinor := newSplit[1]

						if oldMajor == newMajor && oldMinor == newMinor {
							return true
						}
					}
					return false
				},
			},
			"maintenance_window": {
				Type:        schema.TypeList,
				Description: "A maintenance window comprise of a day of the week and a time for maintenance to be allowed",
				Optional:    true,
				MaxItems:    1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"time": {
							Type:        schema.TypeString,
							Description: "A clock time in the day when maintenance is allowed",
							Required:    true,
						},
						"day_of_the_week": {
							Type:        schema.TypeString,
							Description: "Day of the week when maintenance is allowed",
							Required:    true,
						},
					},
				},
			},
			"available_upgrade_versions": {
				Type:        schema.TypeList,
				Description: "List of available versions for upgrading the cluster",
				Optional:    true,
				Computed:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
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
			"public": {
				Type: schema.TypeBool,
				Description: "The indicator if the cluster is public or private. Be aware that setting it to false is " +
					"currently in beta phase.",
				Optional: true,
				Default:  true,
			},
			"gateway_ip": {
				Type: schema.TypeString,
				Description: "The IP address of the gateway used by the cluster. This is mandatory when `public` is set " +
					"to `false` and should not be provided otherwise.",
				Optional: true,
			},
		},
		Timeouts: &resourceDefaultTimeouts,
	}
}

func resourcek8sClusterCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*ionoscloud.APIClient)

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

	public := d.Get("public").(bool)
	cluster.Properties.Public = &public

	if gatewayIp, gatewayIpOk := d.GetOk("gateway_ip"); gatewayIpOk {
		gatewayIp := gatewayIp.(string)
		cluster.Properties.GatewayIp = &gatewayIp
	}

	createdCluster, _, err := client.KubernetesApi.K8sPost(ctx).KubernetesCluster(cluster).Execute()

	if err != nil {
		d.SetId("")
		diags := diag.FromErr(fmt.Errorf("error creating k8s cluster: %s \n", err))
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
			diags := diag.FromErr(fmt.Errorf("k8s cluster creation timed out! WARNING: your k8s cluster will still probably be created " +
				"after some time but the terraform state wont reflect that; check your Ionos Cloud account for updates"))
			return diags
		}
	}

	return resourcek8sClusterRead(ctx, d, meta)
}

func resourcek8sClusterRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*ionoscloud.APIClient)

	cluster, apiResponse, err := client.KubernetesApi.K8sFindByClusterId(ctx, d.Id()).Execute()

	if err != nil {
		if apiResponse != nil && apiResponse.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		diags := diag.FromErr(fmt.Errorf("error while fetching k8s cluster %s: %s", d.Id(), err))
		return diags
	}

	log.Printf("[INFO] Successfully retreived cluster %s: %+v", d.Id(), cluster)

	if cluster.Properties.Name != nil {
		err := d.Set("name", *cluster.Properties.Name)
		if err != nil {
			diags := diag.FromErr(fmt.Errorf("error while setting name property for cluser %s: %s", d.Id(), err))
			return diags
		}
	}

	if cluster.Properties.K8sVersion != nil {
		err := d.Set("k8s_version", *cluster.Properties.K8sVersion)
		if err != nil {
			diags := diag.FromErr(fmt.Errorf("error while setting k8s_version property for cluser %s: %s", d.Id(), err))
			return diags
		}
	}

	if cluster.Properties.Name != nil {
		err := d.Set("name", *cluster.Properties.Name)
		if err != nil {
			diags := diag.FromErr(fmt.Errorf("error while setting name property for cluser %s: %s", d.Id(), err))
			return diags
		}
	}

	if cluster.Properties.AvailableUpgradeVersions != nil {
		availableUpgradeVersions := make([]interface{}, len(*cluster.Properties.AvailableUpgradeVersions), len(*cluster.Properties.AvailableUpgradeVersions))
		for i, availableUpgradeVersion := range *cluster.Properties.AvailableUpgradeVersions {
			availableUpgradeVersions[i] = availableUpgradeVersion
		}
		if err := d.Set("available_upgrade_versions", availableUpgradeVersions); err != nil {
			diags := diag.FromErr(err)
			return diags
		}
	}

	if cluster.Properties.ViableNodePoolVersions != nil {
		viableNodePoolVersions := make([]interface{}, len(*cluster.Properties.ViableNodePoolVersions), len(*cluster.Properties.ViableNodePoolVersions))
		for i, viableNodePoolVersion := range *cluster.Properties.ViableNodePoolVersions {
			viableNodePoolVersions[i] = viableNodePoolVersion
		}
		if err := d.Set("viable_node_pool_versions", viableNodePoolVersions); err != nil {
			diags := diag.FromErr(err)
			return diags
		}
	}

	if cluster.Properties.Public != nil {
		err := d.Set("public", *cluster.Properties.Public)
		if err != nil {
			diags := diag.FromErr(fmt.Errorf("error while setting public property for cluser %s: %s", d.Id(), err))
			return diags
		}
	}

	if cluster.Properties.GatewayIp != nil {
		err := d.Set("gateway_ip", *cluster.Properties.GatewayIp)
		if err != nil {
			diags := diag.FromErr(fmt.Errorf("error while setting gateway_ip property for cluser %s: %s", d.Id(), err))
			return diags
		}
	}

	return nil
}

func resourcek8sClusterUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*ionoscloud.APIClient)

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

	_, apiResponse, err := client.KubernetesApi.K8sPut(ctx, d.Id()).KubernetesCluster(request).Execute()

	if err != nil {
		if _, ok := err.(ionoscloud.GenericOpenAPIError); ok {
			if apiResponse != nil && apiResponse.StatusCode == 404 {
				d.SetId("")
				return nil
			}
			diags := diag.FromErr(fmt.Errorf("error while updating k8s cluster: %s", err))
			return diags
		}
		diags := diag.FromErr(fmt.Errorf("error while updating k8s cluster %s: %s", d.Id(), err))
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
			diags := diag.FromErr(fmt.Errorf("k8s cluster update timed out! WARNING: your k8s cluster will still probably be created " +
				"after some time but the terraform state won't reflect that; check your Ionos Cloud account for updates"))
			return diags
		}
	}

	return resourcek8sClusterRead(ctx, d, meta)
}

func resourcek8sClusterDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*ionoscloud.APIClient)

	_, apiResponse, err := client.KubernetesApi.K8sDelete(ctx, d.Id()).Execute()

	if err != nil {
		if apiResponse != nil && apiResponse.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		diags := diag.FromErr(fmt.Errorf("error while deleting k8s cluster: %s", err))
		return diags
	}

	for {
		log.Printf("[INFO] Waiting for cluster %s to be deleted...", d.Id())

		clusterdDeleted, dsErr := k8sClusterDeleted(ctx, client, d)

		if dsErr != nil {
			diags := diag.FromErr(fmt.Errorf("error while checking deletion status of k8s cluster %s: %s", d.Id(), dsErr))
			return diags
		}

		if clusterdDeleted {
			log.Printf("[INFO] Successfully deleted k8s cluster: %s", d.Id())
			break
		}

		select {
		case <-time.After(SleepInterval):
			log.Printf("[INFO] trying again ...")
		case <-ctx.Done():
			diags := diag.FromErr(fmt.Errorf("k8s cluster deletion timed out! WARNING: your k8s cluster will still probably be deleted " +
				"after some time but the terraform state won't reflect that; check your Ionos Cloud account for updates"))
			return diags
		}
	}

	return nil
}

func k8sClusterReady(ctx context.Context, client *ionoscloud.APIClient, d *schema.ResourceData) (bool, error) {

	subjectCluster, _, err := client.KubernetesApi.K8sFindByClusterId(ctx, d.Id()).Execute()

	if err != nil {
		return true, fmt.Errorf("error checking k8s cluster status: %s", err)
	}
	return *subjectCluster.Metadata.State == "ACTIVE", nil
}

func k8sClusterDeleted(ctx context.Context, client *ionoscloud.APIClient, d *schema.ResourceData) (bool, error) {

	_, apiResponse, err := client.KubernetesApi.K8sFindByClusterId(ctx, d.Id()).Execute()

	if err != nil {
		if _, ok := err.(ionoscloud.GenericOpenAPIError); ok {
			if apiResponse != nil && apiResponse.StatusCode == 404 {
				return true, nil
			}
			return true, fmt.Errorf("error checking k8s cluster deletion status: %s", err)
		}
	}
	return false, nil
}
