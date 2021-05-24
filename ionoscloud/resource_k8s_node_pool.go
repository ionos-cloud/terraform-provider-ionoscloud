package ionoscloud

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	ionoscloud "github.com/ionos-cloud/sdk-go/v5"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
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
				Type:         schema.TypeString,
				Description:  "The desired name for the node pool",
				Required:     true,
				ValidateFunc: validation.All(validation.StringIsNotWhiteSpace),
			},
			"k8s_version": {
				Type:         schema.TypeString,
				Description:  "The desired kubernetes version",
				Required:     true,
				ValidateFunc: validation.All(validation.StringIsNotWhiteSpace),
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
				Type:         schema.TypeString,
				Description:  "The UUID of the VDC",
				Required:     true,
				ValidateFunc: validation.All(validation.StringIsNotWhiteSpace),
			},
			"k8s_cluster_id": {
				Type:         schema.TypeString,
				Description:  "The UUID of an existing kubernetes cluster",
				Required:     true,
				ValidateFunc: validation.All(validation.StringIsNotWhiteSpace),
			},
			"cpu_family": {
				Type:         schema.TypeString,
				Description:  "CPU Family",
				Required:     true,
				ValidateFunc: validation.All(validation.StringIsNotWhiteSpace),
			},
			"availability_zone": {
				Type:         schema.TypeString,
				Description:  "The compute availability zone in which the nodes should exist",
				Required:     true,
				ValidateFunc: validation.All(validation.StringIsNotWhiteSpace),
			},
			"storage_type": {
				Type:         schema.TypeString,
				Description:  "Storage type to use",
				Required:     true,
				ValidateFunc: validation.All(validation.StringIsNotWhiteSpace),
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
	client := meta.(*ionoscloud.APIClient)

	name := d.Get("name").(string)
	datacenterId := d.Get("datacenter_id").(string)
	k8sVersion := d.Get("k8s_version").(string)
	availabilityZone := d.Get("availability_zone").(string)
	cpuFamily := d.Get("cpu_family").(string)
	storageType := d.Get("storage_type").(string)
	nodeCount := int32(d.Get("node_count").(int))
	coresCount := int32(d.Get("cores_count").(int))
	storageSize := int32(d.Get("storage_size").(int))
	ramSize := int32(d.Get("ram_size").(int))

	k8sNodepool := ionoscloud.KubernetesNodePoolForPost{
		Properties: &ionoscloud.KubernetesNodePoolPropertiesForPost{
			Name:             &name,
			DatacenterId:     &datacenterId,
			K8sVersion:       &k8sVersion,
			AvailabilityZone: &availabilityZone,
			CpuFamily:        &cpuFamily,
			StorageType:      &storageType,
			NodeCount:        &nodeCount,
			CoresCount:       &coresCount,
			StorageSize:      &storageSize,
			RamSize:          &ramSize,
		},
	}

	if _, asOk := d.GetOk("auto_scaling.0"); asOk {
		k8sNodepool.Properties.AutoScaling = &ionoscloud.KubernetesAutoScaling{}
	}

	if asmnVal, asmnOk := d.GetOk("auto_scaling.0.min_node_count"); asmnOk {
		log.Printf("[INFO] Setting Autoscaling minimum node count to : %d", uint32(asmnVal.(int)))
		asmnVal := int32(asmnVal.(int))
		k8sNodepool.Properties.AutoScaling.MinNodeCount = &asmnVal
	}

	if asmxnVal, asmxnOk := d.GetOk("auto_scaling.0.max_node_count"); asmxnOk {
		log.Printf("[INFO] Setting Autoscaling maximum node count to : %d", uint32(asmxnVal.(int)))
		asmxnVal := int32(asmxnVal.(int))
		k8sNodepool.Properties.AutoScaling.MaxNodeCount = &asmxnVal
	}

	if k8sNodepool.Properties.AutoScaling != nil && k8sNodepool.Properties.AutoScaling.MinNodeCount != nil &&
		k8sNodepool.Properties.AutoScaling.MaxNodeCount != nil && *k8sNodepool.Properties.AutoScaling.MinNodeCount != 0 &&
		*k8sNodepool.Properties.AutoScaling.MaxNodeCount != 0 && *k8sNodepool.Properties.AutoScaling.MinNodeCount != *k8sNodepool.Properties.AutoScaling.MaxNodeCount {
		log.Printf("[INFO] Autoscaling is on, doing some extra checks for k8s node pool")

		if *k8sNodepool.Properties.NodeCount < *k8sNodepool.Properties.AutoScaling.MinNodeCount {
			d.SetId("")
			return fmt.Errorf("Error creating k8s node pool: node_count cannot be lower than min_node_count")
		}

		if *k8sNodepool.Properties.AutoScaling.MaxNodeCount < *k8sNodepool.Properties.AutoScaling.MinNodeCount {
			d.SetId("")
			return fmt.Errorf("Error creating k8s node pool: max_node_count cannot be lower than min_node_count")
		}
	}

	if _, mwOk := d.GetOk("maintenance_window.0"); mwOk {
		k8sNodepool.Properties.MaintenanceWindow = &ionoscloud.KubernetesMaintenanceWindow{}
	}

	if mtVal, mtOk := d.GetOk("maintenance_window.0.time"); mtOk {
		log.Printf("[INFO] Setting Maintenance Window Time to : %s", mtVal.(string))
		mtVal := mtVal.(string)
		k8sNodepool.Properties.MaintenanceWindow.Time = &mtVal
	}

	if mdVal, mdOk := d.GetOk("maintenance_window.0.day_of_the_week"); mdOk {
		mdVal := mdVal.(string)
		k8sNodepool.Properties.MaintenanceWindow.DayOfTheWeek = &mdVal
	}

	if lansVal, lansOK := d.GetOk("lans"); lansOK {
		if lansVal.([]interface{}) != nil {
			updateLans := false

			lans := []ionoscloud.KubernetesNodePoolLan{}

			for lanIndex := range lansVal.([]interface{}) {
				if lanID, lanIDOk := d.GetOk(fmt.Sprintf("lans.%d", lanIndex)); lanIDOk {
					log.Printf("[INFO] Adding k8s node pool to LAN %+v...", lanID)
					lanID := int32(lanID.(int))
					lans = append(lans, ionoscloud.KubernetesNodePoolLan{Id: &lanID})
				}
			}

			if len(lans) > 0 {
				updateLans = true
			}

			if updateLans == true {
				log.Printf("[INFO] k8s node pool LANs set to %+v", lans)
				k8sNodepool.Properties.Lans = &lans
			}
		}
	}

	publicIpsProp, ok := d.GetOk("public_ips")
	if ok {
		publicIps := publicIpsProp.([]interface{})

		/* number of public IPs needs to be at least NodeCount + 1 */
		if len(publicIps) > 0 && int32(len(publicIps)) < *k8sNodepool.Properties.NodeCount+1 {
			return fmt.Errorf("the number of public IPs must be at least %d", *k8sNodepool.Properties.NodeCount+1)
		}

		requestPublicIps := make([]string, len(publicIps), len(publicIps))
		for i := range publicIps {
			requestPublicIps[i] = fmt.Sprint(publicIps[i])
		}
		k8sNodepool.Properties.PublicIps = &requestPublicIps
	}

	ctx, cancel := context.WithTimeout(context.Background(), *resourceDefaultTimeouts.Create)

	if cancel != nil {
		defer cancel()
	}

	createdNodepool, apiResponse, err := client.KubernetesApi.K8sNodepoolsPost(ctx, d.Get("k8s_cluster_id").(string)).KubernetesNodePool(k8sNodepool).Execute()

	if err != nil {
		d.SetId("")
		payload := "<nil>"
		if apiResponse != nil {
			payload = fmt.Sprintf("API response: %s", string(apiResponse.Payload))
		}
		return fmt.Errorf("error creating k8s node pool: %s; %s", err, payload)
	}

	d.SetId(*createdNodepool.Id)

	log.Printf("[INFO] Successfully created k8s node pool: %s", d.Id())

	for {
		log.Printf("[INFO] Waiting for k8s node pool %s to be ready...", d.Id())
		time.Sleep(10 * time.Second)

		nodepoolReady, rsErr := k8sNodepoolReady(client, d)

		if rsErr != nil {
			return fmt.Errorf("Error while checking readiness status of k8s node pool %s: %s", d.Id(), rsErr)
		}

		if nodepoolReady {
			log.Printf("[INFO] k8s node pool ready: %s", d.Id())
			break
		}
	}

	return resourcek8sNodePoolRead(d, meta)
}

func resourcek8sNodePoolRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ionoscloud.APIClient)

	ctx, cancel := context.WithTimeout(context.Background(), *resourceDefaultTimeouts.Create)

	if cancel != nil {
		defer cancel()
	}

	k8sNodepool, apiResponse, err := client.KubernetesApi.K8sNodepoolsFindById(ctx, d.Get("k8s_cluster_id").(string), d.Id()).Execute()

	if err != nil {
		log.Printf("[INFO] Resource %s not found: %+v", d.Id(), err)
		if apiResponse != nil && apiResponse.Response.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		return err
	}

	log.Printf("[INFO] Successfully retreived k8s node pool %s: %+v", d.Id(), k8sNodepool)

	d.SetId(*k8sNodepool.Id)

	if k8sNodepool.Properties.Name != nil {
		err := d.Set("name", *k8sNodepool.Properties.Name)
		if err != nil {
			return fmt.Errorf("Error while setting name property for k8sNodepool %s: %s", d.Id(), err)
		}

	}

	if k8sNodepool.Properties.K8sVersion != nil {
		err := d.Set("k8s_version", *k8sNodepool.Properties.K8sVersion)
		if err != nil {
			return fmt.Errorf("Error while setting k8s_version property for k8sNodepool %s: %s", d.Id(), err)
		}

	}

	if k8sNodepool.Properties.DatacenterId != nil {
		err := d.Set("datacenter_id", *k8sNodepool.Properties.DatacenterId)
		if err != nil {
			return fmt.Errorf("Error while setting datacenter_id property for k8sNodepool %s: %s", d.Id(), err)
		}

	}

	if k8sNodepool.Properties.CpuFamily != nil {
		err := d.Set("cpu_family", *k8sNodepool.Properties.CpuFamily)
		if err != nil {
			return fmt.Errorf("Error while setting cpu_family property for k8sNodepool %s: %s", d.Id(), err)
		}

	}

	if k8sNodepool.Properties.AvailabilityZone != nil {
		err := d.Set("availability_zone", *k8sNodepool.Properties.AvailabilityZone)
		if err != nil {
			return fmt.Errorf("Error while setting availability_zone property for k8sNodepool %s: %s", d.Id(), err)
		}

	}

	if k8sNodepool.Properties.StorageType != nil {
		err := d.Set("storage_type", *k8sNodepool.Properties.StorageType)
		if err != nil {
			return fmt.Errorf("Error while setting storage_type property for k8sNodepool %s: %s", d.Id(), err)
		}

	}

	if k8sNodepool.Properties.NodeCount != nil {
		err := d.Set("node_count", *k8sNodepool.Properties.NodeCount)
		if err != nil {
			return fmt.Errorf("Error while setting node_count property for k8sNodepool %s: %s", d.Id(), err)
		}

	}

	if k8sNodepool.Properties.CoresCount != nil {
		err := d.Set("cores_count", *k8sNodepool.Properties.CoresCount)
		if err != nil {
			return fmt.Errorf("Error while setting cores_count property for k8sNodepool %s: %s", d.Id(), err)
		}

	}

	if k8sNodepool.Properties.RamSize != nil {
		err := d.Set("ram_size", *k8sNodepool.Properties.RamSize)
		if err != nil {
			return fmt.Errorf("Error while setting ram_size property for k8sNodepool %s: %s", d.Id(), err)
		}

	}

	if k8sNodepool.Properties.StorageSize != nil {
		err := d.Set("storage_size", *k8sNodepool.Properties.StorageSize)
		if err != nil {
			return fmt.Errorf("Error while setting storage_size property for k8sNodepool %s: %s", d.Id(), err)
		}

	}

	if k8sNodepool.Properties.PublicIps != nil {
		err := d.Set("public_ips", *k8sNodepool.Properties.PublicIps)
		if err != nil {
			return fmt.Errorf("Error while setting public_ips property for k8sNodepool %s: %s", d.Id(), err)
		}
	}

	if k8sNodepool.Properties.AutoScaling != nil && k8sNodepool.Properties.AutoScaling.MinNodeCount != nil &&
		k8sNodepool.Properties.AutoScaling.MaxNodeCount != nil && (*k8sNodepool.Properties.AutoScaling.MinNodeCount != 0 &&
		*k8sNodepool.Properties.AutoScaling.MaxNodeCount != 0) {
		err := d.Set("auto_scaling", []map[string]int32{
			{
				"min_node_count": *k8sNodepool.Properties.AutoScaling.MinNodeCount,
				"max_node_count": *k8sNodepool.Properties.AutoScaling.MaxNodeCount,
			},
		})
		if err != nil {
			return fmt.Errorf("Error while setting auto_scaling property for k8sNodepool %s: %s", d.Id(), err)
		}
		log.Printf("[INFO] Setting AutoScaling for k8s node pool %s to %+v...", d.Id(), *k8sNodepool.Properties.AutoScaling)
	}

	return nil
}

func resourcek8sNodePoolUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ionoscloud.APIClient)

	request := ionoscloud.KubernetesNodePoolForPut{}

	nodeCount := int32(d.Get("node_count").(int))
	request.Properties = &ionoscloud.KubernetesNodePoolPropertiesForPut{
		NodeCount: &nodeCount,
	}

	if d.HasChange("k8s_version") {
		oldk8sVersion, newk8sVersion := d.GetChange("k8s_version")
		log.Printf("[INFO] k8s pool k8s version changed from %+v to %+v", oldk8sVersion, newk8sVersion)
		if newk8sVersion != nil {
			newk8sVersion := newk8sVersion.(string)
			request.Properties.K8sVersion = &newk8sVersion
		}
	}

	if d.HasChange("auto_scaling.0") {
		_, newAs := d.GetChange("auto_scaling.0")
		if newAs.(map[string]interface{}) != nil {
			updateAutoscaling := false
			minNodeCount := int32(d.Get("auto_scaling.0.min_node_count").(int))
			maxNodeCount := int32(d.Get("auto_scaling.0.max_node_count").(int))
			autoScaling := &ionoscloud.KubernetesAutoScaling{
				MinNodeCount: &minNodeCount,
				MaxNodeCount: &maxNodeCount,
			}

			if d.HasChange("auto_scaling.0.min_node_count") {
				oldMinNodes, newMinNodes := d.GetChange("auto_scaling.0.min_node_count")
				if newMinNodes != 0 {
					log.Printf("[INFO] k8s node pool autoscaling min # of nodes changed from %+v to %+v", oldMinNodes, newMinNodes)
					updateAutoscaling = true
					newMinNodes := int32(newMinNodes.(int))
					autoScaling.MinNodeCount = &newMinNodes
				}
			}

			if d.HasChange("auto_scaling.0.max_node_count") {
				oldMaxNodes, newMaxNodes := d.GetChange("auto_scaling.0.max_node_count")
				if newMaxNodes != 0 {
					log.Printf("[INFO] k8s node pool autoscaling max # of nodes changed from %+v to %+v", oldMaxNodes, newMaxNodes)
					updateAutoscaling = true
					newMaxNodes := int32(newMaxNodes.(int))
					autoScaling.MaxNodeCount = &newMaxNodes
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

			ctx, cancel := context.WithTimeout(context.Background(), *resourceDefaultTimeouts.Create)

			if cancel != nil {
				defer cancel()
			}

			np, _, npErr := client.KubernetesApi.K8sNodepoolsFindById(ctx, d.Get("k8s_cluster_id").(string), d.Id()).Execute()
			if npErr != nil {
				return fmt.Errorf("Error retrieving k8s node pool %q: %s", d.Id(), npErr)
			}

			log.Printf("[INFO] Setting node_count for node pool %q from server from %d to %d instead of due to autoscaling %+v", d.Id(), uint32(d.Get("node_count").(int)), *np.Properties.NodeCount, d.Get("auto_scaling.0"))
			request.Properties.NodeCount = np.Properties.NodeCount
		}

		if updateNodeCount {
			oldNc, newNc := d.GetChange("node_count")
			log.Printf("[INFO] k8s node pool node_count changed from %+v to %+v", oldNc, newNc)
			if oldNc.(int) != newNc.(int) {
				newNc := int32(newNc.(int))
				request.Properties.NodeCount = &newNc
			}
		}
	}

	if d.HasChange("lans") {
		oldLANs, newLANs := d.GetChange("lans")
		if newLANs.([]interface{}) != nil {
			updateLans := false

			lans := []ionoscloud.KubernetesNodePoolLan{}

			for lanIndex := range newLANs.([]interface{}) {
				if lanID, lanIDOk := d.GetOk(fmt.Sprintf("lans.%d", lanIndex)); lanIDOk {
					log.Printf("[INFO] Adding k8s node pool to LAN %+v...", lanID)
					lanID := int32(lanID.(int))
					lans = append(lans, ionoscloud.KubernetesNodePoolLan{Id: &lanID})
				}
			}

			if len(lans) > 0 {
				updateLans = true
			}

			if updateLans == true {
				log.Printf("[INFO] k8s node pool LANs changed from %+v to %+v", oldLANs, newLANs)
				request.Properties.Lans = &lans
			}
		}
	}

	if d.HasChange("maintenance_window.0") {

		_, newMw := d.GetChange("maintenance_window.0")

		if newMw.(map[string]interface{}) != nil {

			updateMaintenanceWindow := false
			dayOfTheWeek := d.Get("maintenance_window.0.day_of_the_week").(string)
			timeS := d.Get("maintenance_window.0.time").(string)
			maintenanceWindow := &ionoscloud.KubernetesMaintenanceWindow{
				DayOfTheWeek: &dayOfTheWeek,
				Time:         &timeS,
			}

			if d.HasChange("maintenance_window.0.day_of_the_week") {

				oldMd, newMd := d.GetChange("maintenance_window.0.day_of_the_week")
				if newMd.(string) != "" {
					log.Printf("[INFO] k8s node pool maintenance window DOW changed from %+v to %+v", oldMd, newMd)
					updateMaintenanceWindow = true
					newMd := newMd.(string)
					maintenanceWindow.DayOfTheWeek = &newMd
				}
			}

			if d.HasChange("maintenance_window.0.time") {
				oldMt, newMt := d.GetChange("maintenance_window.0.time")
				if newMt.(string) != "" {
					log.Printf("[INFO] k8s node pool maintenance window time changed from %+v to %+v", oldMt, newMt)
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

	if d.HasChange("public_ips") {
		oldPublicIps, newPublicIps := d.GetChange("public_ips")
		log.Printf("[INFO] k8s pool public IPs changed from %+v to %+v", oldPublicIps, newPublicIps)
		if newPublicIps != nil {

			publicIps := newPublicIps.([]interface{})

			/* number of public IPs needs to be at least NodeCount + 1 */
			if len(publicIps) > 0 && int32(len(publicIps)) < *request.Properties.NodeCount+1 {
				return fmt.Errorf("the number of public IPs must be at least %d", *request.Properties.NodeCount+1)
			}

			requestPublicIps := make([]string, len(publicIps), len(publicIps))

			for i := range publicIps {
				requestPublicIps[i] = fmt.Sprint(publicIps[i])
			}

			request.Properties.PublicIps = &requestPublicIps
		}
	}

	b, jErr := json.Marshal(request)

	if jErr == nil {
		log.Printf("[INFO] Update req: %s", string(b))
	}

	ctx, cancel := context.WithTimeout(context.Background(), *resourceDefaultTimeouts.Update)

	if cancel != nil {
		defer cancel()
	}
	_, apiResponse, err := client.KubernetesApi.K8sNodepoolsPut(ctx, d.Get("k8s_cluster_id").(string), d.Id()).KubernetesNodePool(request).Execute()

	if err != nil {
		if apiResponse != nil && apiResponse.Response.StatusCode == 404 {
			d.SetId("")
			return nil
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

		if nodepoolReady {
			log.Printf("[INFO] k8s node pool ready: %s", d.Id())
			break
		}
	}

	return resourcek8sNodePoolRead(d, meta)
}

func resourcek8sNodePoolDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ionoscloud.APIClient)

	ctx, cancel := context.WithTimeout(context.Background(), *resourceDefaultTimeouts.Delete)

	if cancel != nil {
		defer cancel()
	}

	_, apiResponse, err := client.KubernetesApi.K8sNodepoolsDelete(ctx, d.Get("k8s_cluster_id").(string), d.Id()).Execute()

	if err != nil {
		if apiResponse != nil && apiResponse.Response.StatusCode == 404 {
			d.SetId("")
			return nil
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

		if nodepoolDeleted {
			log.Printf("[INFO] Successfully deleted k8s node pool: %s", d.Id())
			break
		}
	}

	d.SetId("")
	return nil
}

func k8sNodepoolReady(client *ionoscloud.APIClient, d *schema.ResourceData) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), *resourceDefaultTimeouts.Default)

	if cancel != nil {
		defer cancel()
	}

	subjectNodepool, _, err := client.KubernetesApi.K8sNodepoolsFindById(ctx, d.Get("k8s_cluster_id").(string), d.Id()).Execute()
	if err != nil {
		return true, fmt.Errorf("Error checking k8s node pool status: %s", err)
	}
	return *subjectNodepool.Metadata.State == "ACTIVE", nil
}

func k8sNodepoolDeleted(client *ionoscloud.APIClient, d *schema.ResourceData) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), *resourceDefaultTimeouts.Default)

	if cancel != nil {
		defer cancel()
	}
	_, apiResponse, err := client.KubernetesApi.K8sNodepoolsFindById(ctx, d.Get("k8s_cluster_id").(string), d.Id()).Execute()

	if err != nil {
		if _, ok := err.(ionoscloud.GenericOpenAPIError); ok {
			if apiResponse != nil && apiResponse.Response.StatusCode == 404 {
				return true, nil
			}
			return true, fmt.Errorf("Error checking k8s node pool deletion status: %s", err)
		}
	}
	return false, nil
}
