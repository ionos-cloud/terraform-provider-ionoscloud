package ionoscloud

import (
	"context"
	"fmt"
	ionoscloud "github.com/ionos-cloud/sdk-go/v5"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourcek8sCluster() *schema.Resource {
	return &schema.Resource{
		Create: resourcek8sClusterCreate,
		Read:   resourcek8sClusterRead,
		Update: resourcek8sClusterUpdate,
		Delete: resourcek8sClusterDelete,
		Importer: &schema.ResourceImporter{
			State: resourceK8sClusterImport,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Description: "The desired name for the cluster",
				Required:    true,
			},
			"k8s_version": {
				Type:        schema.TypeString,
				Description: "The desired kubernetes version",
				Optional:    true,
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
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"viable_node_pool_versions": {
				Type:        schema.TypeList,
				Description: "List of versions that may be used for node pools under this cluster",
				Optional:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
		Timeouts: &resourceDefaultTimeouts,
	}
}

func resourcek8sClusterCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(SdkBundle).Client

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

	auvVal, ok := d.GetOk("available_upgrade_versions")
	if ok {
		auvVal := auvVal.([]interface{})

		requestAvailableUpgradeVersions := make([]string, len(auvVal), len(auvVal))

		for i := range auvVal {
			requestAvailableUpgradeVersions[i] = fmt.Sprint(auvVal[i])
		}
	}

	vnpvVal, ok := d.GetOk("viable_node_pool_versions")
	if ok {
		vnpvVal := vnpvVal.([]interface{})

		requestViableNodePoolVersions := make([]string, len(vnpvVal), len(vnpvVal))

		for i := range vnpvVal {
			requestViableNodePoolVersions[i] = fmt.Sprint(vnpvVal[i])
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), *resourceDefaultTimeouts.Create)

	if cancel != nil {
		defer cancel()
	}

	createdCluster, apiResponse, err := client.KubernetesApi.K8sPost(ctx).KubernetesCluster(cluster).Execute()

	if err != nil {
		d.SetId("")
		return fmt.Errorf("Error creating k8s cluster: %s \n ApiError: %s ", err, string(apiResponse.Payload))
	}

	d.SetId(*createdCluster.Id)
	log.Printf("[INFO] Created k8s cluster: %s", d.Id())

	for {
		log.Printf("[INFO] Waiting for cluster %s to be ready...", d.Id())
		time.Sleep(5 * time.Second)

		clusterReady, rsErr := k8sClusterReady(client, d)

		if rsErr != nil {
			return fmt.Errorf("Error while checking readiness status of k8s cluster %s: %s", d.Id(), rsErr)
		}

		if clusterReady {
			log.Printf("[INFO] k8s cluster ready: %s", d.Id())
			break
		}
	}

	return resourcek8sClusterRead(d, meta)
}

func resourcek8sClusterRead(d *schema.ResourceData, meta interface{}) error {

	client := meta.(SdkBundle).Client

	ctx, cancel := context.WithTimeout(context.Background(), *resourceDefaultTimeouts.Default)

	if cancel != nil {
		defer cancel()
	}

	cluster, apiResponse, err := client.KubernetesApi.K8sFindByClusterId(ctx, d.Id()).Execute()

	if err != nil {
		if _, ok := err.(ionoscloud.GenericOpenAPIError); ok {
			if apiResponse.Response.StatusCode == 404 {
				d.SetId("")
				return nil
			}
		}
		return fmt.Errorf("Error while fetching k8s cluster %s: %s", d.Id(), err)
	}

	log.Printf("[INFO] Successfully retreived cluster %s: %+v", d.Id(), cluster)

	return nil
}

func resourcek8sClusterUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(SdkBundle).Client
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

	if d.HasChange("available_upgrade_versions") {
		oldAvailableUpgradeVersions, newAvailableUpgradeVersions := d.GetChange("available_upgrade_versions")
		log.Printf("[INFO] k8s cluster available upgrade versions changed from %+v to %+v", oldAvailableUpgradeVersions, newAvailableUpgradeVersions)
		if newAvailableUpgradeVersions != nil {

			availableUpgradeVersions := newAvailableUpgradeVersions.([]interface{})

			requestAvailableUpgradeVersions := make([]string, len(availableUpgradeVersions), len(availableUpgradeVersions))

			for i := range availableUpgradeVersions {
				requestAvailableUpgradeVersions[i] = fmt.Sprint(availableUpgradeVersions[i])
			}
		}
	}

	if d.HasChange("viable_node_pool_versions") {
		oldViableNodePoolVersions, newViableNodePoolVersions := d.GetChange("viable_node_pool_versions")
		log.Printf("[INFO] k8s cluster viable node pool versions changed from %+v to %+v", oldViableNodePoolVersions, newViableNodePoolVersions)
		if newViableNodePoolVersions != nil {

			availableViableNodePoolVersions := newViableNodePoolVersions.([]interface{})

			requestViableNodePoolVersions := make([]string, len(availableViableNodePoolVersions), len(availableViableNodePoolVersions))

			for i := range availableViableNodePoolVersions {
				requestViableNodePoolVersions[i] = fmt.Sprint(availableViableNodePoolVersions[i])
			}
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), *resourceDefaultTimeouts.Update)

	if cancel != nil {
		defer cancel()
	}

	_, apiResponse, err := client.KubernetesApi.K8sPut(ctx, d.Id()).KubernetesCluster(request).Execute()

	if err != nil {
		if _, ok := err.(ionoscloud.GenericOpenAPIError); ok {
			if apiResponse.Response.StatusCode == 404 {
				d.SetId("")
				return nil
			}
			return fmt.Errorf("Error while updating k8s cluster: %s", err)
		}
		return fmt.Errorf("Error while updating k8s cluster %s: %s", d.Id(), err)
	}

	for {
		log.Printf("[INFO] Waiting for cluster %s to be ready...", d.Id())
		time.Sleep(5 * time.Second)

		clusterReady, rsErr := k8sClusterReady(client, d)

		if rsErr != nil {
			return fmt.Errorf("Error while checking readiness status of k8s cluster %s: %s", d.Id(), rsErr)
		}

		if clusterReady {
			log.Printf("[INFO] k8s cluster ready: %s", d.Id())
			break
		}
	}

	return resourcek8sClusterRead(d, meta)
}

func resourcek8sClusterDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(SdkBundle).Client

	ctx, cancel := context.WithTimeout(context.Background(), *resourceDefaultTimeouts.Delete)

	if cancel != nil {
		defer cancel()
	}

	_, apiResponse, err := client.KubernetesApi.K8sDelete(ctx, d.Id()).Execute()

	if err != nil {
		if _, ok := err.(ionoscloud.GenericOpenAPIError); ok {
			if apiResponse.Response.StatusCode == 404 {
				d.SetId("")
				return nil
			}
			return fmt.Errorf("Error while deleting k8s cluster: %s", err)
		}

		return fmt.Errorf("Error while deleting k8s cluster %s: %s", d.Id(), err)
	}

	for {
		log.Printf("[INFO] Waiting for cluster %s to be deleted...", d.Id())
		time.Sleep(5 * time.Second)

		clusterdDeleted, dsErr := k8sClusterDeleted(client, d)

		if dsErr != nil {
			return fmt.Errorf("Error while checking deletion status of k8s cluster %s: %s", d.Id(), dsErr)
		}

		if clusterdDeleted {
			log.Printf("[INFO] Successfully deleted k8s cluster: %s", d.Id())
			break
		}
	}

	return nil
}

func k8sClusterReady(client *ionoscloud.APIClient, d *schema.ResourceData) (bool, error) {

	ctx, cancel := context.WithTimeout(context.Background(), *resourceDefaultTimeouts.Default)

	if cancel != nil {
		defer cancel()
	}

	subjectCluster, _, err := client.KubernetesApi.K8sFindByClusterId(ctx, d.Id()).Execute()

	if err != nil {
		return true, fmt.Errorf("Error checking k8s cluster status: %s", err)
	}
	return *subjectCluster.Metadata.State == "ACTIVE", nil
}

func k8sClusterDeleted(client *ionoscloud.APIClient, d *schema.ResourceData) (bool, error) {

	ctx, cancel := context.WithTimeout(context.Background(), *resourceDefaultTimeouts.Default)

	if cancel != nil {
		defer cancel()
	}

	_, apiResponse, err := client.KubernetesApi.K8sFindByClusterId(ctx, d.Id()).Execute()

	if err != nil {
		if _, ok := err.(ionoscloud.GenericOpenAPIError); ok {
			if apiResponse.Response.StatusCode == 404 {
				return true, nil
			}
			return true, fmt.Errorf("Error checking k8s cluster deletion status: %s", err)
		}
	}
	return false, nil
}
