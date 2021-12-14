package ionoscloud

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	ionoscloud "github.com/ionos-cloud/sdk-go/v6"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceK8sNodePool() *schema.Resource {
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
				Type:        schema.TypeString,
				Description: "The desired kubernetes version",
				Optional:    true,
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
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeInt,
							Description: "The LAN ID of an existing LAN at the related datacenter",
							Required:    true,
						},
						"dhcp": {
							Type:        schema.TypeBool,
							Description: "Indicates if the Kubernetes Node Pool LAN will reserve an IP using DHCP",
							Optional:    true,
							Default:     true,
						},
						"routes": {
							Type:        schema.TypeList,
							Description: "An array of additional LANs attached to worker nodes",
							Optional:    true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"network": {
										Type:         schema.TypeString,
										Description:  "IPv4 or IPv6 CIDR to be routed via the interface",
										Required:     true,
										ValidateFunc: validation.All(validation.StringIsNotWhiteSpace),
									},
									"gateway_ip": {
										Type:         schema.TypeString,
										Description:  "IPv4 or IPv6 Gateway IP for the route",
										Required:     true,
										ValidateFunc: validation.All(validation.StringIsNotWhiteSpace),
									},
								},
							},
						},
					},
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
			"available_upgrade_versions": {
				Type:        schema.TypeList,
				Description: "A list of kubernetes versions available for upgrade",
				Computed:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
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
				lan := ionoscloud.KubernetesNodePoolLan{}
				addLan := false
				if lanID, lanIdOk := d.GetOk(fmt.Sprintf("lans.%d.id", lanIndex)); lanIdOk {
					log.Printf("[INFO] Adding k8s node pool to LAN %+v...", lanID)
					lanID := int32(lanID.(int))
					lan.Id = &lanID
					addLan = true
				}

				lanDhcp := d.Get(fmt.Sprintf("lans.%d.dhcp", lanIndex)).(bool)
				lan.Dhcp = &lanDhcp

				if lanRoutes, lanRoutesOk := d.GetOk(fmt.Sprintf("lans.%d.routes", lanIndex)); lanRoutesOk {
					if lanRoutes.([]interface{}) != nil {
						updateRoutes := false

						var routes []ionoscloud.KubernetesNodePoolLanRoutes

						for routeIndex := range lanRoutes.([]interface{}) {

							addRoute := false
							route := ionoscloud.KubernetesNodePoolLanRoutes{}
							if routeNetwork, routeNewtworkOk := d.GetOk(fmt.Sprintf("lans.%d.routes.%d.network", lanIndex, routeIndex)); routeNewtworkOk {
								routeNetwork := routeNetwork.(string)
								route.Network = &routeNetwork
								addRoute = true
							}

							if routeGatewayIp, routeGatewayIpOk := d.GetOk(fmt.Sprintf("lans.%d.routes.%d.gateway_ip", lanIndex, routeIndex)); routeGatewayIpOk {
								routeGatewayIp := routeGatewayIp.(string)
								route.GatewayIp = &routeGatewayIp
								addRoute = true
							}

							if addRoute {
								routes = append(routes, route)
							}
						}

						if len(routes) > 0 {
							updateRoutes = true
						}

						if updateRoutes == true {
							log.Printf("[INFO] k8s node pool LanRoutes set to %+v", routes)
							lan.Routes = &routes
						}
					}
				}
				if addLan {
					lans = append(lans, lan)
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

	createdNodepool, apiResponse, err := client.KubernetesApi.
		K8sNodepoolsPost(ctx, d.Get("k8s_cluster_id").(string)).
		KubernetesNodePool(k8sNodepool).
		Execute()
	logApiRequestTime(apiResponse)
	if err != nil {
		d.SetId("")
		diags := diag.FromErr(fmt.Errorf("error creating k8s node pool: %s", err))
		return diags
	}

	d.SetId(*createdNodepool.Id)

	log.Printf("[INFO] Successfully created k8s node pool: %s", d.Id())

	for {
		log.Printf("[INFO] Waiting for k8s node pool %s to be ready...", d.Id())

		nodepoolReady, rsErr := k8sNodepoolReady(ctx, client, d)

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
			log.Printf("[INFO] timed out")
			diags := diag.FromErr(fmt.Errorf("k8s creation timed out"))
			return diags
		}
	}

	return resourcek8sNodePoolRead(ctx, d, meta)
}

func resourcek8sNodePoolRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*ionoscloud.APIClient)

	k8sNodepool, apiResponse, err := client.KubernetesApi.K8sNodepoolsFindById(ctx, d.Get("k8s_cluster_id").(string), d.Id()).Execute()
	logApiRequestTime(apiResponse)

	if err != nil {
		log.Printf("[INFO] Resource %s not found: %+v", d.Id(), err)
		if apiResponse != nil && apiResponse.Response != nil && apiResponse.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		diags := diag.FromErr(err)
		return diags
	}

	log.Printf("[INFO] Successfully retreived k8s node pool %s: %+v", d.Id(), k8sNodepool)

	if err := setK8sNodePoolData(d, &k8sNodepool); err != nil {
		return diag.FromErr(err)
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

			if *autoScaling.MaxNodeCount < *autoScaling.MinNodeCount {
				d.SetId("")
				diags := diag.FromErr(fmt.Errorf("error updating k8s node pool: max_node_count cannot be lower than min_node_count"))
				return diags
			}

			if updateAutoscaling == true {
				request.Properties.AutoScaling = autoScaling
			}
		}
	}

	if d.HasChange("node_count") {
		oldNc, newNc := d.GetChange("node_count")
		nodeCount := int32(newNc.(int))

		if d.Get("auto_scaling.0").(map[string]interface{}) != nil && (d.Get("auto_scaling.0.min_node_count").(int) != 0 || d.Get("auto_scaling.0.max_node_count").(int) != 0) {

			if nodeCount < *request.Properties.AutoScaling.MinNodeCount {
				d.SetId("")
				diags := diag.FromErr(fmt.Errorf("error updating k8s node pool: node_count cannot be lower than min_node_count"))
				return diags
			}

			if nodeCount > *request.Properties.AutoScaling.MaxNodeCount {
				d.SetId("")
				diags := diag.FromErr(fmt.Errorf("error updating k8s node pool: node_count cannot be greater than max_node_count"))
				return diags
			}
		}

		log.Printf("[INFO] k8s node pool node_count changed from %+v to %+v", oldNc, newNc)
		request.Properties.NodeCount = &nodeCount
	}

	if d.HasChange("lans") {
		oldLANs, newLANs := d.GetChange("lans")
		lans := make([]ionoscloud.KubernetesNodePoolLan, 0)
		if newLANs.([]interface{}) != nil {
			for lanIndex := range newLANs.([]interface{}) {
				lan := ionoscloud.KubernetesNodePoolLan{}
				if lanID, lanIdOk := d.GetOk(fmt.Sprintf("lans.%d.id", lanIndex)); lanIdOk {
					log.Printf("[INFO] Adding k8s node pool to LAN %+v...", lanID)
					lanID := int32(lanID.(int))
					lan.Id = &lanID
				}

				lanDhcp := d.Get(fmt.Sprintf("lans.%d.dhcp", lanIndex)).(bool)
				lan.Dhcp = &lanDhcp
				routes := make([]ionoscloud.KubernetesNodePoolLanRoutes, 0)
				if lanRoutes, lanRoutesOk := d.GetOk(fmt.Sprintf("lans.%d.routes", lanIndex)); lanRoutesOk {
					if lanRoutes.([]interface{}) != nil {
						for routeIndex := range lanRoutes.([]interface{}) {

							route := ionoscloud.KubernetesNodePoolLanRoutes{}
							if routeNetwork, routeNewtworkOk := d.GetOk(fmt.Sprintf("lans.%d.routes.%d.network", lanIndex, routeIndex)); routeNewtworkOk {
								routeNetwork := routeNetwork.(string)
								route.Network = &routeNetwork
							}

							if routeGatewayIp, routeGatewayIpOk := d.GetOk(fmt.Sprintf("lans.%d.routes.%d.gateway_ip", lanIndex, routeIndex)); routeGatewayIpOk {
								routeGatewayIp := routeGatewayIp.(string)
								route.GatewayIp = &routeGatewayIp
							}

							routes = append(routes, route)

						}

						log.Printf("[INFO] k8s node pool LanRoutes set to %+v", routes)
					}
				}
				lan.Routes = &routes
				lans = append(lans, lan)
			}
		}
		log.Printf("[INFO] k8s node pool LANs changed from %+v to %+v", oldLANs, newLANs)

		request.Properties.Lans = &lans
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
		requestPublicIps := make([]string, 0)

		if newPublicIps != nil {

			publicIps := newPublicIps.([]interface{})

			/* number of public IPs needs to be at least NodeCount + 1 */
			if len(publicIps) > 0 && int32(len(publicIps)) < *request.Properties.NodeCount+1 {
				diags := diag.FromErr(fmt.Errorf("the number of public IPs must be at least %d", *request.Properties.NodeCount+1))
				return diags
			}

			for _, ip := range publicIps {
				requestPublicIps = append(requestPublicIps, ip.(string))
			}

		}
		request.Properties.PublicIps = &requestPublicIps

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
	logApiRequestTime(apiResponse)

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

		nodepoolReady, rsErr := k8sNodepoolReady(ctx, client, d)

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
			log.Printf("[DEBUG] retrying ...")
		case <-ctx.Done():
			diags := diag.FromErr(fmt.Errorf("k8s node pool update timed out! WARNING: your k8s node pool will still probably be updated after some time but the terraform state won't reflect that; check your Ionos Cloud account for updates"))
			return diags
		}
	}

	return resourcek8sNodePoolRead(ctx, d, meta)
}

func resourcek8sNodePoolDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*ionoscloud.APIClient)

	apiResponse, err := client.KubernetesApi.K8sNodepoolsDelete(ctx, d.Get("k8s_cluster_id").(string), d.Id()).Execute()
	logApiRequestTime(apiResponse)

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

		nodepoolDeleted, dsErr := k8sNodepoolDeleted(ctx, client, d)

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
			log.Printf("[DEBUG] retrying ...")
		case <-ctx.Done():
			diags := diag.FromErr(fmt.Errorf("k8s node pool deletion timed out! WARNING: your k8s node pool will still probably be deleted after some time but the terraform state won't reflect that; check your Ionos Cloud account for updates"))
			return diags
		}
	}

	d.SetId("")
	return nil
}

func resourceK8sNodepoolImport(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {

	parts := strings.Split(d.Id(), "/")

	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return nil, fmt.Errorf("invalid import id %q. Expecting {k8sClusterId}/{k8sNodePoolId}", d.Id())
	}

	clusterId := parts[0]
	npId := parts[1]

	client := meta.(*ionoscloud.APIClient)

	k8sNodepool, apiResponse, err := client.KubernetesApi.K8sNodepoolsFindById(ctx, clusterId, npId).Execute()
	logApiRequestTime(apiResponse)

	if err != nil {
		if _, ok := err.(ionoscloud.GenericOpenAPIError); ok {
			if apiResponse != nil && apiResponse.Response != nil && apiResponse.StatusCode == 404 {
				d.SetId("")
				return nil, fmt.Errorf("unable to find k8s node pool %q", npId)
			}
		}
		return nil, fmt.Errorf("unable to retreive k8s node pool %q", npId)
	}

	log.Printf("[INFO] K8s node pool found: %+v", k8sNodepool)

	if err := setK8sNodePoolData(d, &k8sNodepool); err != nil {
		return nil, err
	}
	log.Printf("[INFO] Importing k8s node pool %q...", d.Id())

	return []*schema.ResourceData{d}, nil
}

func setK8sNodePoolData(d *schema.ResourceData, nodePool *ionoscloud.KubernetesNodePool) error {

	if nodePool.Id != nil {
		d.SetId(*nodePool.Id)
	}

	if nodePool.Properties != nil {
		if nodePool.Properties.Name != nil {
			if err := d.Set("name", *nodePool.Properties.Name); err != nil {
				return err
			}
		}

		if nodePool.Properties.DatacenterId != nil {
			if err := d.Set("datacenter_id", *nodePool.Properties.DatacenterId); err != nil {
				return err
			}
		}

		if nodePool.Properties.NodeCount != nil {
			if err := d.Set("node_count", *nodePool.Properties.NodeCount); err != nil {
				return err
			}
		}

		if nodePool.Properties.CpuFamily != nil {
			if err := d.Set("cpu_family", *nodePool.Properties.CpuFamily); err != nil {
				return err
			}
		}

		if nodePool.Properties.CoresCount != nil {
			if err := d.Set("cores_count", *nodePool.Properties.CoresCount); err != nil {
				return err
			}
		}

		if nodePool.Properties.RamSize != nil {
			if err := d.Set("ram_size", *nodePool.Properties.RamSize); err != nil {
				return err
			}
		}

		if nodePool.Properties.AvailabilityZone != nil {
			if err := d.Set("availability_zone", *nodePool.Properties.AvailabilityZone); err != nil {
				return err
			}
		}

		if nodePool.Properties.StorageType != nil {
			if err := d.Set("storage_type", *nodePool.Properties.StorageType); err != nil {
				return err
			}
		}

		if nodePool.Properties.StorageSize != nil {
			if err := d.Set("storage_size", *nodePool.Properties.StorageSize); err != nil {
				return err
			}
		}

		if nodePool.Properties.K8sVersion != nil {
			if err := d.Set("k8s_version", *nodePool.Properties.K8sVersion); err != nil {
				return err
			}
		}

		if nodePool.Properties.PublicIps != nil && len(*nodePool.Properties.PublicIps) > 0 {
			if err := d.Set("public_ips", *nodePool.Properties.PublicIps); err != nil {
				return err
			}
		}

		if nodePool.Properties.MaintenanceWindow != nil && nodePool.Properties.MaintenanceWindow.Time != nil && nodePool.Properties.MaintenanceWindow.DayOfTheWeek != nil {
			if err := d.Set("maintenance_window", []map[string]string{
				{
					"time":            *nodePool.Properties.MaintenanceWindow.Time,
					"day_of_the_week": *nodePool.Properties.MaintenanceWindow.DayOfTheWeek,
				},
			}); err != nil {
				return err
			}
		}

		if nodePool.Properties.AutoScaling != nil && nodePool.Properties.AutoScaling.MinNodeCount != nil &&
			nodePool.Properties.AutoScaling.MaxNodeCount != nil && (*nodePool.Properties.AutoScaling.MinNodeCount != 0 &&
			*nodePool.Properties.AutoScaling.MaxNodeCount != 0) {
			if err := d.Set("auto_scaling", []map[string]uint32{
				{
					"min_node_count": uint32(*nodePool.Properties.AutoScaling.MinNodeCount),
					"max_node_count": uint32(*nodePool.Properties.AutoScaling.MaxNodeCount),
				},
			}); err != nil {
				return err
			}
		}

		if nodePool.Properties.Lans != nil && len(*nodePool.Properties.Lans) > 0 {

			nodePoolLans := getK8sNodePoolLans(*nodePool.Properties.Lans)

			if err := d.Set("lans", nodePoolLans); err != nil {
				return fmt.Errorf("error while setting lans property for k8sNodepool %s: %s", d.Id(), err)
			}

		}

		if nodePool.Properties.AvailableUpgradeVersions != nil && len(*nodePool.Properties.AvailableUpgradeVersions) > 0 {
			if err := d.Set("available_upgrade_versions", *nodePool.Properties.AvailableUpgradeVersions); err != nil {
				return err
			}
		}

		if nodePool.Properties.PublicIps != nil && len(*nodePool.Properties.PublicIps) > 0 {
			if err := d.Set("public_ips", *nodePool.Properties.PublicIps); err != nil {
				return err
			}
		}

		labels := make(map[string]interface{})
		if nodePool.Properties.Labels != nil && len(*nodePool.Properties.Labels) > 0 {
			for k, v := range *nodePool.Properties.Labels {
				labels[k] = v
			}
		}

		if err := d.Set("labels", labels); err != nil {
			return fmt.Errorf("error while setting the labels property for k8sNodepool %s: %s", d.Id(), err)

		}

		annotations := make(map[string]interface{})
		if nodePool.Properties.Annotations != nil && len(*nodePool.Properties.Annotations) > 0 {
			for k, v := range *nodePool.Properties.Annotations {
				annotations[k] = v
			}
		}

		if err := d.Set("annotations", annotations); err != nil {
			return fmt.Errorf("error while setting the annotations property for k8sNodepool %s: %s", d.Id(), err)
		}

	}

	return nil
}
func k8sNodepoolReady(ctx context.Context, client *ionoscloud.APIClient, d *schema.ResourceData) (bool, error) {
	subjectNodepool, apiResponse, err := client.KubernetesApi.K8sNodepoolsFindById(ctx, d.Get("k8s_cluster_id").(string), d.Id()).Execute()
	logApiRequestTime(apiResponse)
	if err != nil {
		return true, fmt.Errorf("error checking k8s node pool status: %s", err)
	}
	return *subjectNodepool.Metadata.State == "ACTIVE", nil
}

func k8sNodepoolDeleted(ctx context.Context, client *ionoscloud.APIClient, d *schema.ResourceData) (bool, error) {
	_, apiResponse, err := client.KubernetesApi.K8sNodepoolsFindById(ctx, d.Get("k8s_cluster_id").(string), d.Id()).Execute()
	logApiRequestTime(apiResponse)

	if err != nil {
		if apiResponse != nil && apiResponse.Response != nil && apiResponse.StatusCode == 404 {
			return true, nil
		}
		return true, fmt.Errorf("error checking k8s node pool deletion status: %s", err)
	}
	return false, nil
}

func getK8sNodePoolLans(lans []ionoscloud.KubernetesNodePoolLan) []interface{} {

	var nodePoolLans []interface{}
	for _, nodePoolLan := range lans {
		lanEntry := make(map[string]interface{})

		if nodePoolLan.Id != nil {
			lanEntry["id"] = *nodePoolLan.Id
		}

		if nodePoolLan.Dhcp != nil {
			lanEntry["dhcp"] = *nodePoolLan.Dhcp
		}

		if nodePoolLan.Routes != nil && len(*nodePoolLan.Routes) > 0 {
			var nodePoolRoutes []interface{}
			for _, nodePoolRoute := range *nodePoolLan.Routes {
				routeEntry := make(map[string]string)
				if nodePoolRoute.Network != nil {
					routeEntry["network"] = *nodePoolRoute.Network
				}
				if nodePoolRoute.GatewayIp != nil {
					routeEntry["gateway_ip"] = *nodePoolRoute.GatewayIp
				}
				nodePoolRoutes = append(nodePoolRoutes, routeEntry)
			}

			if len(nodePoolRoutes) > 0 {
				lanEntry["routes"] = nodePoolRoutes
			}
		}

		nodePoolLans = append(nodePoolLans, lanEntry)
	}

	return nodePoolLans

}
