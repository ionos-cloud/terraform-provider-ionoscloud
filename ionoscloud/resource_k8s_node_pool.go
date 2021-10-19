package ionoscloud

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	ionoscloud "github.com/ionos-cloud/sdk-go/v5"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourcek8sNodePool() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourcek8sNodePoolCreate,
		ReadContext:   resourcek8sNodePoolRead,
		UpdateContext: resourcek8sNodePoolUpdate,
		DeleteContext: resourcek8sNodePoolDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceK8sNodepoolImport,
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
			"labels": {
				Type:     schema.TypeMap,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"annotations": {
				Type:     schema.TypeMap,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
		Timeouts: &resourceDefaultTimeouts,
	}
}

func resourcek8sNodePoolCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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
			diags := diag.FromErr(fmt.Errorf("error creating k8s node pool: node_count cannot be lower than min_node_count"))
			return diags
		}

		if *k8sNodepool.Properties.AutoScaling.MaxNodeCount < *k8sNodepool.Properties.AutoScaling.MinNodeCount {
			d.SetId("")
			diags := diag.FromErr(fmt.Errorf("error creating k8s node pool: max_node_count cannot be lower than min_node_count"))
			return diags
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

			var lans []ionoscloud.KubernetesNodePoolLan

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
			diags := diag.FromErr(fmt.Errorf("the number of public IPs must be at least %d", *k8sNodepool.Properties.NodeCount+1))
			return diags
		}

		var requestPublicIps []string
		for i := range publicIps {
			requestPublicIps = append(requestPublicIps, fmt.Sprint(publicIps[i]))
		}
		k8sNodepool.Properties.PublicIps = &requestPublicIps
	}

	labelsProp, ok := d.GetOk("labels")
	if ok {
		labels := make(map[string]string)
		for k, v := range labelsProp.(map[string]interface{}) {
			labels[k] = v.(string)
		}
		k8sNodepool.Properties.Labels = &labels
	}

	annotationsProp, ok := d.GetOk("annotations")
	if ok {
		annotations := make(map[string]string)
		for k, v := range annotationsProp.(map[string]interface{}) {
			annotations[k] = v.(string)
		}
		k8sNodepool.Properties.Annotations = &annotations
	}

	createdNodepool, _, err := client.KubernetesApi.K8sNodepoolsPost(ctx, d.Get("k8s_cluster_id").(string)).KubernetesNodePool(k8sNodepool).Execute()

	if err != nil {
		d.SetId("")
		diags := diag.FromErr(fmt.Errorf("error creating k8s node pool: %s", err))
		return diags
	}

	d.SetId(*createdNodepool.Id)
	log.Printf("[INFO] Successfully created k8s node pool: %s", d.Id())

	for {
		log.Printf("[INFO] Waiting for k8s node pool %s to be ready...", d.Id())

		nodepoolReady, rsErr := k8sNodepoolReady(client, d)

		if rsErr != nil {
			diags := diag.FromErr(fmt.Errorf("error while checking readiness status of k8s node pool %s: %s", d.Id(), rsErr))
			return diags
		}

		if nodepoolReady {
			log.Printf("[INFO] k8s node pool ready: %s", d.Id())
			break
		}

		select {
		case <-time.After(SleepInterval):
			log.Printf("[INFO] trying again ...")
		case <-ctx.Done():
			diags := diag.FromErr(fmt.Errorf("k8s node pool creation timed out! WARNING: your k8s nodepool will still probably be " +
				"created after some time but the terraform state won't reflect that; check your Ionos Cloud account for updates"))
			return diags
		}
	}

	return resourcek8sNodePoolRead(ctx, d, meta)
}

func resourcek8sNodePoolRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*ionoscloud.APIClient)

	k8sNodepool, apiResponse, err := client.KubernetesApi.K8sNodepoolsFindById(ctx, d.Get("k8s_cluster_id").(string), d.Id()).Execute()

	if err != nil {
		if apiResponse != nil && apiResponse.Response != nil && apiResponse.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		diags := diag.FromErr(fmt.Errorf("error while fetching k8s node pool %s: %s", d.Id(), err))
		return diags
	}

	log.Printf("[INFO] Successfully retreived k8s node pool %s: %+v", d.Id(), k8sNodepool)

	d.SetId(*k8sNodepool.Id)

	if k8sNodepool.Properties.Name != nil {
		err := d.Set("name", *k8sNodepool.Properties.Name)
		if err != nil {
			diags := diag.FromErr(fmt.Errorf("error while setting name property for k8sNodepool %s: %s", d.Id(), err))
			return diags
		}

	}

	if k8sNodepool.Properties.K8sVersion != nil {
		err := d.Set("k8s_version", *k8sNodepool.Properties.K8sVersion)
		if err != nil {
			diags := diag.FromErr(fmt.Errorf("error while setting k8s_version property for k8sNodepool %s: %s", d.Id(), err))
			return diags
		}

	}

	if k8sNodepool.Properties.DatacenterId != nil {
		err := d.Set("datacenter_id", *k8sNodepool.Properties.DatacenterId)
		if err != nil {
			diags := diag.FromErr(fmt.Errorf("error while setting datacenter_id property for k8sNodepool %s: %s", d.Id(), err))
			return diags
		}

	}

	if k8sNodepool.Properties.CpuFamily != nil {
		err := d.Set("cpu_family", *k8sNodepool.Properties.CpuFamily)
		if err != nil {
			diags := diag.FromErr(fmt.Errorf("error while setting cpu_family property for k8sNodepool %s: %s", d.Id(), err))
			return diags
		}

	}

	if k8sNodepool.Properties.AvailabilityZone != nil {
		err := d.Set("availability_zone", *k8sNodepool.Properties.AvailabilityZone)
		if err != nil {
			diags := diag.FromErr(fmt.Errorf("error while setting availability_zone property for k8sNodepool %s: %s", d.Id(), err))
			return diags
		}

	}

	if k8sNodepool.Properties.StorageType != nil {
		err := d.Set("storage_type", *k8sNodepool.Properties.StorageType)
		if err != nil {
			diags := diag.FromErr(fmt.Errorf("error while setting storage_type property for k8sNodepool %s: %s", d.Id(), err))
			return diags
		}

	}

	if k8sNodepool.Properties.NodeCount != nil {
		err := d.Set("node_count", *k8sNodepool.Properties.NodeCount)
		if err != nil {
			diags := diag.FromErr(fmt.Errorf("error while setting node_count property for k8sNodepool %s: %s", d.Id(), err))
			return diags
		}

	}

	if k8sNodepool.Properties.CoresCount != nil {
		err := d.Set("cores_count", *k8sNodepool.Properties.CoresCount)
		if err != nil {
			diags := diag.FromErr(fmt.Errorf("error while setting cores_count property for k8sNodepool %s: %s", d.Id(), err))
			return diags
		}

	}

	if k8sNodepool.Properties.RamSize != nil {
		err := d.Set("ram_size", *k8sNodepool.Properties.RamSize)
		if err != nil {
			diags := diag.FromErr(fmt.Errorf("error while setting ram_size property for k8sNodepool %s: %s", d.Id(), err))
			return diags
		}

	}

	if k8sNodepool.Properties.StorageSize != nil {
		err := d.Set("storage_size", *k8sNodepool.Properties.StorageSize)
		if err != nil {
			diags := diag.FromErr(fmt.Errorf("error while setting storage_size property for k8sNodepool %s: %s", d.Id(), err))
			return diags
		}

	}

	if k8sNodepool.Properties.PublicIps != nil && len(*k8sNodepool.Properties.PublicIps) > 0 {
		err := d.Set("public_ips", *k8sNodepool.Properties.PublicIps)
		if err != nil {
			diags := diag.FromErr(fmt.Errorf("error while setting public_ips property for k8sNodepool %s: %s", d.Id(), err))
			return diags
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
			diags := diag.FromErr(fmt.Errorf("error while setting auto_scaling property for k8sNodepool %s: %s", d.Id(), err))
			return diags
		}
		log.Printf("[INFO] Setting AutoScaling for k8s node pool %s to %+v...", d.Id(), *k8sNodepool.Properties.AutoScaling)
	}

	if k8sNodepool.Properties.Lans != nil && len(*k8sNodepool.Properties.Lans) > 0 {
		lans := make([]int32, 0, 0)
		for _, lan := range *k8sNodepool.Properties.Lans {
			if lan.Id != nil {
				lans = append(lans, *lan.Id)
			}
		}
		err := d.Set("lans", lans)
		if err != nil {
			diags := diag.FromErr(fmt.Errorf("error while setting lans property for k8sNodepool %s: %s", d.Id(), err))
			return diags
		}
	}

	labels := make(map[string]interface{})
	if k8sNodepool.Properties.Labels != nil && len(*k8sNodepool.Properties.Labels) > 0 {
		for k, v := range *k8sNodepool.Properties.Labels {
			labels[k] = v
		}
	}

	if err := d.Set("labels", labels); err != nil {
		diags := diag.FromErr(fmt.Errorf("error while setting the labels property for k8sNodepool %s: %s", d.Id(), err))
		return diags
	}

	annotations := make(map[string]interface{})
	if k8sNodepool.Properties.Annotations != nil && len(*k8sNodepool.Properties.Annotations) > 0 {
		for k, v := range *k8sNodepool.Properties.Annotations {
			annotations[k] = v
		}
	}

	if err := d.Set("annotations", annotations); err != nil {
		diags := diag.FromErr(fmt.Errorf("error while setting the annotations property for k8sNodepool %s: %s", d.Id(), err))
		return diags
	}

	return nil
}

func resourcek8sNodePoolUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*ionoscloud.APIClient)

	request := ionoscloud.KubernetesNodePoolForPut{}

	nodeCount := int32(d.Get("node_count").(int))
	request.Properties = &ionoscloud.KubernetesNodePoolPropertiesForPut{
		NodeCount: &nodeCount,
	}

	if d.HasChange("name") {
		diags := diag.FromErr(fmt.Errorf("name attribute is immutable, therefore not allowed in update requests"))
		return diags
	}

	if d.HasChange("cpu_family") {
		diags := diag.FromErr(fmt.Errorf("cpu_family attribute is immutable, therefore not allowed in update requests"))
		return diags
	}

	if d.HasChange("availability_zone") {
		diags := diag.FromErr(fmt.Errorf("availability_zone attribute is immutable, therefore not allowed in update requests"))
		return diags
	}

	if d.HasChange("cores_count") {
		diags := diag.FromErr(fmt.Errorf("cores_count attribute is immutable, therefore not allowed in update requests"))
		return diags
	}

	if d.HasChange("ram_size") {
		diags := diag.FromErr(fmt.Errorf("ram_size attribute is immutable, therefore not allowed in update requests"))
		return diags
	}

	if d.HasChange("storage_size") {
		diags := diag.FromErr(fmt.Errorf("storage_size attribute is immutable, therefore not allowed in update requests"))
		return diags
	}

	if d.HasChange("storage_type") {
		diags := diag.FromErr(fmt.Errorf("storage_size attribute is immutable, therefore not allowed in update requests"))
		return diags
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
				diags := diag.FromErr(fmt.Errorf("error retrieving k8s node pool %q: %s", d.Id(), npErr))
				return diags
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

			var lans []ionoscloud.KubernetesNodePoolLan

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
				diags := diag.FromErr(fmt.Errorf("the number of public IPs must be at least %d", *request.Properties.NodeCount+1))
				return diags
			}

			requestPublicIps := make([]string, len(publicIps), len(publicIps))

			for i := range publicIps {
				requestPublicIps[i] = fmt.Sprint(publicIps[i])
			}

			request.Properties.PublicIps = &requestPublicIps
		}
	}

	if d.HasChange("labels") {
		oldLabels, newLabels := d.GetChange("labels")
		log.Printf("[INFO] k8s pool labels changed from %+v to %+v", oldLabels, newLabels)
		labels := make(map[string]string)
		if newLabels != nil {
			for k, v := range newLabels.(map[string]interface{}) {
				labels[k] = v.(string)
			}
		}
		request.Properties.Labels = &labels
	}

	if d.HasChange("annotations") {
		oldAnnotations, newAnnotations := d.GetChange("annotations")
		log.Printf("[INFO] k8s pool annotations changed from %+v to %+v", oldAnnotations, newAnnotations)
		annotations := make(map[string]string)
		if newAnnotations != nil {
			for k, v := range newAnnotations.(map[string]interface{}) {
				annotations[k] = v.(string)
			}
		}
		request.Properties.Annotations = &annotations
	}

	b, jErr := json.Marshal(request)

	if jErr == nil {
		log.Printf("[INFO] Update req: %s", string(b))
	}

	_, apiResponse, err := client.KubernetesApi.K8sNodepoolsPut(ctx, d.Get("k8s_cluster_id").(string), d.Id()).KubernetesNodePool(request).Execute()

	if err != nil {
		if apiResponse != nil && apiResponse.Response != nil && apiResponse.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		diags := diag.FromErr(fmt.Errorf("error while updating k8s node pool %s: %s", d.Id(), err))
		return diags
	}

	for {
		log.Printf("[INFO] Waiting for k8s node pool %s to be ready...", d.Id())

		nodepoolReady, rsErr := k8sNodepoolReady(client, d)

		if rsErr != nil {
			diags := diag.FromErr(fmt.Errorf("error while checking readiness status of k8s node pool %s: %s", d.Id(), rsErr))
			return diags
		}

		if nodepoolReady {
			log.Printf("[INFO] k8s node pool ready: %s", d.Id())
			break
		}

		time.Sleep(SleepInterval * 3)
		select {
		case <-time.After(SleepInterval):
			log.Printf("[INFO] trying again ...")
		case <-ctx.Done():
			diags := diag.FromErr(fmt.Errorf("k8s node pool update timed out! WARNING: your k8s nodepool will still probably be " +
				"updated after some time but the terraform state won't reflect that; check your Ionos Cloud account for updates"))
			return diags
		}
	}

	return resourcek8sNodePoolRead(ctx, d, meta)
}

func resourcek8sNodePoolDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*ionoscloud.APIClient)

	_, apiResponse, err := client.KubernetesApi.K8sNodepoolsDelete(ctx, d.Get("k8s_cluster_id").(string), d.Id()).Execute()

	if err != nil {
		if apiResponse != nil && apiResponse.Response != nil && apiResponse.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		diags := diag.FromErr(fmt.Errorf("error while deleting k8s node pool %s: %s", d.Id(), err))
		return diags
	}

	for {
		log.Printf("[INFO] Waiting for k8s node pool %s to be deleted...", d.Id())

		nodepoolDeleted, dsErr := k8sNodepoolDeleted(client, d)

		if dsErr != nil {
			diags := diag.FromErr(fmt.Errorf("error while checking deletion status of k8s node pool %s: %s", d.Id(), dsErr))
			return diags
		}

		if nodepoolDeleted {
			log.Printf("[INFO] Successfully deleted k8s node pool: %s", d.Id())
			break
		}

		select {
		case <-time.After(SleepInterval):
			log.Printf("[INFO] trying again ...")
		case <-ctx.Done():
			diags := diag.FromErr(fmt.Errorf("k8s node pool deletion timed out! WARNING: your k8s nodepool will still probably be " +
				"deleted after some time but the terraform state won't reflect that; check your Ionos Cloud account for updates"))
			return diags
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
		return true, fmt.Errorf("error checking k8s node pool status: %s", err)
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
		if apiResponse != nil && apiResponse.Response != nil && apiResponse.StatusCode == 404 {
			return true, nil
		}
		return true, fmt.Errorf("error checking k8s node pool deletion status: %s", err)
	}
	return false, nil
}
