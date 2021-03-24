package ionoscloud

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/profitbricks/profitbricks-sdk-go/v5"
)

func resourcek8sNodePool() *schema.Resource {
	return &schema.Resource{
		Create: resourcek8sNodePoolCreate,
		Read:   resourcek8sNodePoolRead,
		Update: resourcek8sNodePoolUpdate,
		Delete: resourcek8sNodePoolDelete,
		Importer: &schema.ResourceImporter{
			State: resourceK8sNodepoolImport,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Description: "The desired name for the node pool",
				Required:    true,
			},
			"k8s_version": {
				Type:        schema.TypeString,
				Description: "The desired kubernetes version",
				Required:    true,
			},
			"auto_scaling": {
				Type:        schema.TypeList,
				Description: "The range defining the minimum and maximum number of worker nodes that the managed node group can scale in",
				Optional:    true,
				MaxItems:    1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"min_node_count": {
							Type:        schema.TypeInt,
							Description: "The minimum number of worker nodes the node pool can scale down to. Should be less than max_node_count",
							Required:    true,
						},
						"max_node_count": {
							Type:        schema.TypeInt,
							Description: "The maximum number of worker nodes that the node pool can scale to. Should be greater than min_node_count",
							Required:    true,
						},
					},
				},
			},
			"lans": {
				Type:        schema.TypeList,
				Description: "A list of Local Area Networks the node pool should be part of",
				Optional:    true,
				Elem: &schema.Schema{
					Type: schema.TypeInt,
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
			"datacenter_id": {
				Type:        schema.TypeString,
				Description: "The UUID of the VDC",
				Required:    true,
			},
			"k8s_cluster_id": {
				Type:        schema.TypeString,
				Description: "The UUID of an existing kubernetes cluster",
				Required:    true,
			},
			"cpu_family": {
				Type:        schema.TypeString,
				Description: "CPU Family",
				Required:    true,
			},
			"availability_zone": {
				Type:        schema.TypeString,
				Description: "The compute availability zone in which the nodes should exist",
				Required:    true,
			},
			"storage_type": {
				Type:        schema.TypeString,
				Description: "Storage type to use",
				Required:    true,
			},
			"node_count": {
				Type:        schema.TypeInt,
				Description: "The number of nodes in this node pool",
				Required:    true,
			},
			"cores_count": {
				Type:        schema.TypeInt,
				Description: "CPU cores count",
				Required:    true,
			},
			"ram_size": {
				Type:        schema.TypeInt,
				Description: "The amount of RAM in MB",
				Required:    true,
			},
			"storage_size": {
				Type:        schema.TypeInt,
				Description: "The total allocated storage capacity of a node in GB",
				Required:    true,
			},
			"public_ips": {
				Type:        schema.TypeList,
				Description: "A list of fixed IPs",
				Optional:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
		Timeouts: &resourceDefaultTimeouts,
	}
}

func resourcek8sNodePoolCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(SdkBundle).LegacyClient

	k8sNodepool := profitbricks.KubernetesNodePool{
		Properties: &profitbricks.KubernetesNodePoolProperties{
			Name:             d.Get("name").(string),
			DatacenterID:     d.Get("datacenter_id").(string),
			K8sVersion:       d.Get("k8s_version").(string),
			AvailabilityZone: d.Get("availability_zone").(string),
			CPUFamily:        d.Get("cpu_family").(string),
			StorageType:      d.Get("storage_type").(string),
			NodeCount:        uint32(d.Get("node_count").(int)),
			CoresCount:       uint32(d.Get("cores_count").(int)),
			StorageSize:      uint32(d.Get("storage_size").(int)),
			RAMSize:          uint32(d.Get("ram_size").(int)),
		},
	}

	if _, asOk := d.GetOk("auto_scaling.0"); asOk {
		k8sNodepool.Properties.AutoScaling = &profitbricks.AutoScaling{}
	}

	if asmnVal, asmnOk := d.GetOk("auto_scaling.0.min_node_count"); asmnOk {
		log.Printf("[INFO] Setting Autoscaling minimum node count to : %d", uint32(asmnVal.(int)))
		k8sNodepool.Properties.AutoScaling.MinNodeCount = uint32(asmnVal.(int))
	}

	if asmxnVal, asmxnOk := d.GetOk("auto_scaling.0.max_node_count"); asmxnOk {
		log.Printf("[INFO] Setting Autoscaling maximum node count to : %d", uint32(asmxnVal.(int)))
		k8sNodepool.Properties.AutoScaling.MaxNodeCount = uint32(asmxnVal.(int))
	}

	if k8sNodepool.Properties.AutoScaling != nil && k8sNodepool.Properties.AutoScaling.MinNodeCount != 0 && k8sNodepool.Properties.AutoScaling.MaxNodeCount != 0 && k8sNodepool.Properties.AutoScaling.MinNodeCount != k8sNodepool.Properties.AutoScaling.MaxNodeCount {
		log.Printf("[INFO] Autoscaling is on, doing some extra checks for k8s node pool")

		if k8sNodepool.Properties.NodeCount < k8sNodepool.Properties.AutoScaling.MinNodeCount {
			d.SetId("")
			return fmt.Errorf("Error creating k8s node pool: node_count cannot be lower than min_node_count")
		}

		if k8sNodepool.Properties.AutoScaling.MaxNodeCount < k8sNodepool.Properties.AutoScaling.MinNodeCount {
			d.SetId("")
			return fmt.Errorf("Error creating k8s node pool: max_node_count cannot be lower than min_node_count")
		}
	}

	if _, mwOk := d.GetOk("maintenance_window.0"); mwOk {
		k8sNodepool.Properties.MaintenanceWindow = &profitbricks.MaintenanceWindow{}
	}

	if mtVal, mtOk := d.GetOk("maintenance_window.0.time"); mtOk {
		log.Printf("[INFO] Setting Maintenance Window Time to : %s", mtVal.(string))
		k8sNodepool.Properties.MaintenanceWindow.Time = mtVal.(string)
	}

	if mdVal, mdOk := d.GetOk("maintenance_window.0.day_of_the_week"); mdOk {
		k8sNodepool.Properties.MaintenanceWindow.DayOfTheWeek = mdVal.(string)
	}

	if lansVal, lansOK := d.GetOk("lans"); lansOK {
		if lansVal.([]interface{}) != nil {
			updateLans := false

			lans := []profitbricks.KubernetesNodePoolLAN{}

			for lanIndex := range lansVal.([]interface{}) {
				if lanID, lanIDOk := d.GetOk(fmt.Sprintf("lans.%d", lanIndex)); lanIDOk {
					log.Printf("[INFO] Adding k8s node pool to LAN %+v...", lanID)
					lans = append(lans, profitbricks.KubernetesNodePoolLAN{ID: uint32(lanID.(int))})
				}
			}

			if len(lans) > 0 {
				updateLans = true
			}

			if updateLans == true {
				log.Printf("[INFO] k8s node pool LANs set to %+v", lans)
				k8sNodepool.Properties.LANs = &lans
			}
		}
	}

	publicIpsProp, ok := d.GetOk("public_ips")
	if ok {
		publicIps := publicIpsProp.([]interface{})

		/* number of public IPs needs to be at least NodeCount + 1 */
		if len(publicIps) > 0 && uint32(len(publicIps)) < k8sNodepool.Properties.NodeCount + 1 {
			return fmt.Errorf("the number of public IPs must be at least %d", k8sNodepool.Properties.NodeCount + 1)
		}

		requestPublicIps := make([]string, len(publicIps), len(publicIps))
		for i := range publicIps {
			requestPublicIps[i] = fmt.Sprint(publicIps[i])
		}
		k8sNodepool.Properties.PublicIPs = &requestPublicIps
	}

	createdNodepool, err := client.CreateKubernetesNodePool(d.Get("k8s_cluster_id").(string), k8sNodepool)

	if err != nil {
		d.SetId("")
		return fmt.Errorf("Error creating k8s node pool: %s", err)
	}

	d.SetId(createdNodepool.ID)

	log.Printf("[INFO] Successfully created k8s node pool: %s", d.Id())

	for {
		log.Printf("[INFO] Waiting for k8s node pool %s to be ready...", d.Id())
		time.Sleep(10 * time.Second)

		nodepoolReady, rsErr := k8sNodepoolReady(client, d)

		if rsErr != nil {
			return fmt.Errorf("Error while checking readiness status of k8s node pool %s: %s", d.Id(), rsErr)
		}

		if nodepoolReady && rsErr == nil {
			log.Printf("[INFO] k8s node pool ready: %s", d.Id())
			break
		}
	}

	return resourcek8sNodePoolRead(d, meta)
}

func resourcek8sNodePoolRead(d *schema.ResourceData, meta interface{}) error {

	client := meta.(SdkBundle).LegacyClient
	k8sNodepool, err := client.GetKubernetesNodePool(d.Get("k8s_cluster_id").(string), d.Id())

	if err != nil {
		log.Printf("[INFO] Resource %s not found: %+v", d.Id(), err)
		if apiError, ok := err.(profitbricks.ApiError); ok {
			if apiError.HttpStatusCode() == 404 {
				d.SetId("")
				return nil
			}
		}
	}

	log.Printf("[INFO] Successfully retreived k8s node pool %s: %+v", d.Id(), k8sNodepool)

	d.SetId(k8sNodepool.ID)
	d.Set("name", k8sNodepool.Properties.Name)
	d.Set("k8s_version", k8sNodepool.Properties.K8sVersion)
	d.Set("datacenter_id", k8sNodepool.Properties.DatacenterID)
	d.Set("cpu_family", k8sNodepool.Properties.CPUFamily)
	d.Set("availability_zone", k8sNodepool.Properties.AvailabilityZone)
	d.Set("storage_type", k8sNodepool.Properties.StorageType)
	d.Set("node_count", k8sNodepool.Properties.NodeCount)
	d.Set("cores_count", k8sNodepool.Properties.CoresCount)
	d.Set("ram_size", k8sNodepool.Properties.RAMSize)
	d.Set("storage_size", k8sNodepool.Properties.StorageSize)

	if k8sNodepool.Properties.PublicIPs != nil {
		d.Set("public_ips", k8sNodepool.Properties.PublicIPs)
	}

	if k8sNodepool.Properties.AutoScaling != nil && (k8sNodepool.Properties.AutoScaling.MinNodeCount != 0 && k8sNodepool.Properties.AutoScaling.MaxNodeCount != 0) {
		d.Set("auto_scaling", []map[string]uint32{
			{
				"min_node_count": k8sNodepool.Properties.AutoScaling.MinNodeCount,
				"max_node_count": k8sNodepool.Properties.AutoScaling.MaxNodeCount,
			},
		})
		log.Printf("[INFO] Setting AutoScaling for k8s node pool %s to %+v...", d.Id(), k8sNodepool.Properties.AutoScaling)
	}

	return nil
}

func resourcek8sNodePoolUpdate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(SdkBundle).LegacyClient
	request := profitbricks.KubernetesNodePool{}

	request.Properties = &profitbricks.KubernetesNodePoolProperties{
		NodeCount: uint32(d.Get("node_count").(int)),
	}

	if d.HasChange("k8s_version") {
		oldk8sVersion, newk8sVersion := d.GetChange("k8s_version")
		log.Printf("[INFO] k8s pool k8s version changed from %+v to %+v", oldk8sVersion, newk8sVersion)
		if newk8sVersion != nil {
			request.Properties.K8sVersion = newk8sVersion.(string)
		}
	}

	if d.HasChange("auto_scaling.0") {
		_, newAs := d.GetChange("auto_scaling.0")
		if newAs.(map[string]interface{}) != nil {
			updateAutoscaling := false
			autoScaling := &profitbricks.AutoScaling{
				MinNodeCount: uint32(d.Get("auto_scaling.0.min_node_count").(int)),
				MaxNodeCount: uint32(d.Get("auto_scaling.0.max_node_count").(int)),
			}

			if d.HasChange("auto_scaling.0.min_node_count") {
				oldMinNodes, newMinNodes := d.GetChange("auto_scaling.0.min_node_count")
				if newMinNodes != 0 {
					log.Printf("[INFO] k8s node pool autoscaling min # of nodes changed from %+v to %+v", oldMinNodes, newMinNodes)
					updateAutoscaling = true
					autoScaling.MinNodeCount = uint32(newMinNodes.(int))
				}
			}

			if d.HasChange("auto_scaling.0.max_node_count") {
				oldMaxNodes, newMaxNodes := d.GetChange("auto_scaling.0.max_node_count")
				if newMaxNodes != 0 {
					log.Printf("[INFO] k8s node pool autoscaling max # of nodes changed from %+v to %+v", oldMaxNodes, newMaxNodes)
					updateAutoscaling = true
					autoScaling.MaxNodeCount = uint32(newMaxNodes.(int))
				}
			}

			if updateAutoscaling == true {
				request.Properties.AutoScaling = autoScaling
			}
		}
	}

	if d.HasChange("node_count") {
		updateNodeCount := true

		if d.Get("auto_scaling.0").(map[string]interface{}) != nil && (d.Get("auto_scaling.0.min_node_count").(int) != 0 || d.Get("auto_scaling.0.max_node_count").(int) != 0) {

			updateNodeCount = false
			np, npErr := client.GetKubernetesNodePool(d.Get("k8s_cluster_id").(string), d.Id())
			if npErr != nil {
				return fmt.Errorf("Error retrieving k8s node pool %q: %s", d.Id(), npErr)
			}

			log.Printf("[INFO] Setting node_count for node pool %q from server from %d to %d instead of due to autoscaling %+v", d.Id(), uint32(d.Get("node_count").(int)), np.Properties.NodeCount, d.Get("auto_scaling.0"))
			request.Properties.NodeCount = uint32(np.Properties.NodeCount)
		}

		if updateNodeCount {
			oldNc, newNc := d.GetChange("node_count")
			log.Printf("[INFO] k8s node pool node_count changed from %+v to %+v", oldNc, newNc)
			if oldNc.(int) != newNc.(int) {
				request.Properties.NodeCount = uint32(newNc.(int))
			}
		}
	}

	if d.HasChange("lans") {
		oldLANs, newLANs := d.GetChange("lans")
		if newLANs.([]interface{}) != nil {
			updateLans := false

			lans := []profitbricks.KubernetesNodePoolLAN{}

			for lanIndex := range newLANs.([]interface{}) {
				if lanID, lanIDOk := d.GetOk(fmt.Sprintf("lans.%d", lanIndex)); lanIDOk {
					log.Printf("[INFO] Adding k8s node pool to LAN %+v...", lanID)
					lans = append(lans, profitbricks.KubernetesNodePoolLAN{ID: uint32(lanID.(int))})
				}
			}

			if len(lans) > 0 {
				updateLans = true
			}

			if updateLans == true {
				log.Printf("[INFO] k8s node pool LANs changed from %+v to %+v", oldLANs, newLANs)
				request.Properties.LANs = &lans
			}
		}
	}

	if d.HasChange("maintenance_window.0") {

		_, newMw := d.GetChange("maintenance_window.0")

		if newMw.(map[string]interface{}) != nil {

			updateMaintenanceWindow := false
			maintenanceWindow := &profitbricks.MaintenanceWindow{
				DayOfTheWeek: d.Get("maintenance_window.0.day_of_the_week").(string),
				Time:         d.Get("maintenance_window.0.time").(string),
			}

			if d.HasChange("maintenance_window.0.day_of_the_week") {

				oldMd, newMd := d.GetChange("maintenance_window.0.day_of_the_week")
				if newMd.(string) != "" {
					log.Printf("[INFO] k8s node pool maintenance window DOW changed from %+v to %+v", oldMd, newMd)
					updateMaintenanceWindow = true
					maintenanceWindow.DayOfTheWeek = newMd.(string)
				}
			}

			if d.HasChange("maintenance_window.0.time") {
				oldMt, newMt := d.GetChange("maintenance_window.0.time")
				if newMt.(string) != "" {
					log.Printf("[INFO] k8s node pool maintenance window time changed from %+v to %+v", oldMt, newMt)
					updateMaintenanceWindow = true
					maintenanceWindow.Time = newMt.(string)
				}
			}

			if updateMaintenanceWindow == true {
				request.Properties.MaintenanceWindow = maintenanceWindow
			}
		}
	}

	if d.HasChange("public_ips") {
		oldPublicIps, newPublicIps := d.GetChange("public_ips")
		log.Printf("[INFO] k8s pool public IPs changed from %+v to %+v", oldPublicIps, newPublicIps)
		if newPublicIps != nil {

			publicIps := newPublicIps.([]interface{})

			/* number of public IPs needs to be at least NodeCount + 1 */
			if len(publicIps) > 0 && uint32(len(publicIps)) < request.Properties.NodeCount + 1 {
				return fmt.Errorf("the number of public IPs must be at least %d", request.Properties.NodeCount + 1)
			}

			requestPublicIps := make([]string, len(publicIps), len(publicIps))

			for i := range publicIps {
				requestPublicIps[i] = fmt.Sprint(publicIps[i])
			}

			request.Properties.PublicIPs = &requestPublicIps
		}
	}

	b, jErr := json.Marshal(request)

	if jErr == nil {
		log.Printf("[INFO] Update req: %s", string(b))
	}

	_, err := client.UpdateKubernetesNodePool(d.Get("k8s_cluster_id").(string), d.Id(), request)

	if err != nil {
		if apiError, ok := err.(profitbricks.ApiError); ok {
			if apiError.HttpStatusCode() == 404 {
				d.SetId("")
				return nil
			}
			return fmt.Errorf("Error while updating k8s node pool: %s", err)
		}
		return fmt.Errorf("Error while updating k8s node pool %s: %s", d.Id(), err)
	}

	for {
		log.Printf("[INFO] Waiting for k8s node pool %s to be ready...", d.Id())
		time.Sleep(10 * time.Second)

		nodepoolReady, rsErr := k8sNodepoolReady(client, d)

		if rsErr != nil {
			return fmt.Errorf("Error while checking readiness status of k8s node pool %s: %s", d.Id(), rsErr)
		}

		if nodepoolReady && rsErr == nil {
			log.Printf("[INFO] k8s node pool ready: %s", d.Id())
			break
		}
	}

	return resourcek8sNodePoolRead(d, meta)
}

func resourcek8sNodePoolDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(SdkBundle).LegacyClient

	_, err := client.DeleteKubernetesNodePool(d.Get("k8s_cluster_id").(string), d.Id())

	if err != nil {
		if apiError, ok := err.(profitbricks.ApiError); ok {
			if apiError.HttpStatusCode() == 404 {
				d.SetId("")
				return nil
			}
			return fmt.Errorf("Error while deleting k8s node pool: %s", err)
		}

		return fmt.Errorf("Error while deleting k8s node pool %s: %s", d.Id(), err)
	}

	for {
		log.Printf("[INFO] Waiting for k8s node pool %s to be deleted...", d.Id())
		time.Sleep(10 * time.Second)

		nodepoolDeleted, dsErr := k8sNodepoolDeleted(client, d)

		if dsErr != nil {
			return fmt.Errorf("Error while checking deletion status of k8s node pool %s: %s", d.Id(), dsErr)
		}

		if nodepoolDeleted && dsErr == nil {
			log.Printf("[INFO] Successfully deleted k8s node pool: %s", d.Id())
			break
		}
	}

	d.SetId("")
	return nil
}

func k8sNodepoolReady(client *profitbricks.Client, d *schema.ResourceData) (bool, error) {
	subjectNodepool, err := client.GetKubernetesNodePool(d.Get("k8s_cluster_id").(string), d.Id())

	if err != nil {
		return true, fmt.Errorf("Error checking k8s node pool status: %s", err)
	}
	return subjectNodepool.Metadata.State == "ACTIVE", nil
}

func k8sNodepoolDeleted(client *profitbricks.Client, d *schema.ResourceData) (bool, error) {
	_, err := client.GetKubernetesNodePool(d.Get("k8s_cluster_id").(string), d.Id())

	if err != nil {
		if apiError, ok := err.(profitbricks.ApiError); ok {
			if apiError.HttpStatusCode() == 404 {
				return true, nil
			}
			return true, fmt.Errorf("Error checking k8s node pool deletion status: %s", err)
		}
	}
	return false, nil
}
