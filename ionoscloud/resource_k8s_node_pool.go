package ionoscloud

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/ionos-cloud/sdk-go-bundle/shared"
	"log"
	"reflect"
	"strings"
	"time"

	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils/constant"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	ionoscloud "github.com/ionos-cloud/sdk-go-bundle/products/cloud/v2"
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
		CustomizeDiff: checkNodePoolImmutableFields,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:             schema.TypeString,
				Description:      "The desired name for the node pool",
				Required:         true,
				ForceNew:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringIsNotWhiteSpace),
			},
			"k8s_version": {
				Type:             schema.TypeString,
				Description:      "The desired Kubernetes Version. For supported values, please check the API documentation. Downgrades are not supported. The provider will ignore downgrades of patch level.",
				Required:         true,
				DiffSuppressFunc: DiffBasedOnVersion,
			},
			"auto_scaling": {
				Type:        schema.TypeList,
				Description: "The range defining the minimum and maximum number of worker nodes that the managed node group can scale in",
				Optional:    true,
				MaxItems:    1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"min_node_count": {
							Type:             schema.TypeInt,
							Description:      "The minimum number of worker nodes the node pool can scale down to. Should be less than max_node_count",
							Required:         true,
							ValidateDiagFunc: validation.ToDiagFunc(validation.IntAtLeast(1)),
						},
						"max_node_count": {
							Type:             schema.TypeInt,
							Description:      "The maximum number of worker nodes that the node pool can scale to. Should be greater than min_node_count",
							Required:         true,
							ValidateDiagFunc: validation.ToDiagFunc(validation.IntAtLeast(1)),
						},
					},
				},
			},
			"lans": {
				Type:        schema.TypeSet,
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
							Type:        schema.TypeSet,
							Description: "An array of additional LANs attached to worker nodes",
							Optional:    true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"network": {
										Type:             schema.TypeString,
										Description:      "IPv4 or IPv6 CIDR to be routed via the interface",
										Required:         true,
										ValidateDiagFunc: validation.ToDiagFunc(validation.StringIsNotWhiteSpace),
									},
									"gateway_ip": {
										Type:             schema.TypeString,
										Description:      "IPv4 or IPv6 Gateway IP for the route",
										Required:         true,
										ValidateDiagFunc: validation.ToDiagFunc(validation.StringIsNotWhiteSpace),
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
				Computed:    true,
				MaxItems:    1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"time": {
							Type:             schema.TypeString,
							Description:      "A clock time in the day when maintenance is allowed",
							Required:         true,
							ValidateDiagFunc: validation.ToDiagFunc(validation.StringIsNotWhiteSpace),
						},
						"day_of_the_week": {
							Type:             schema.TypeString,
							Description:      "Day of the week when maintenance is allowed",
							Required:         true,
							ValidateDiagFunc: validation.ToDiagFunc(validation.StringIsNotWhiteSpace),
						},
					},
				},
			},
			"datacenter_id": {
				Type:             schema.TypeString,
				Description:      "The UUID of the VDC",
				Required:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.IsUUID),
			},
			"k8s_cluster_id": {
				Type:             schema.TypeString,
				Description:      "The UUID of an existing kubernetes cluster",
				Required:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.IsUUID),
			},
			"cpu_family": {
				Type:             schema.TypeString,
				Description:      "CPU Family",
				Required:         true,
				ForceNew:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringIsNotWhiteSpace),
			},
			"availability_zone": {
				Type:             schema.TypeString,
				Description:      "The compute availability zone in which the nodes should exist",
				Required:         true,
				ForceNew:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"AUTO", "ZONE_1", "ZONE_2"}, true)),
			},
			"storage_type": {
				Type:             schema.TypeString,
				Description:      "Storage type to use",
				Required:         true,
				ForceNew:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringIsNotWhiteSpace),
			},
			"node_count": {
				Type:        schema.TypeInt,
				Description: "The number of nodes in this node pool",
				Required:    true,
			},
			"cores_count": {
				Type:        schema.TypeInt,
				Description: "CPU cores count",
				ForceNew:    true,
				Required:    true,
			},
			"ram_size": {
				Type:        schema.TypeInt,
				Description: "The amount of RAM in MB",
				ForceNew:    true,
				Required:    true,
			},
			"storage_size": {
				Type:        schema.TypeInt,
				Description: "The total allocated storage capacity of a node in GB",
				ForceNew:    true,
				Required:    true,
			},
			"public_ips": {
				Type:        schema.TypeList,
				Description: "A list of fixed IPs. Cannot be set on private clusters.",
				Optional:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			// "gateway_ip": {
			//	Type:        schema.TypeString,
			//	Description: "Public IP address for the gateway performing source NAT for the node pool's nodes belonging to a private cluster. Required only if the node pool belongs to a private cluster.",
			//	ForceNew:    true,
			//	Optional:    true,
			//},
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
			"allow_replace": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "When set to true, allows the update of immutable fields by destroying and re-creating the node pool",
			},
		},
		Timeouts:      &constant.ResourceK8sNodePoolTimeout,
		SchemaVersion: 1,
		StateUpgraders: []schema.StateUpgrader{
			{
				Type:    resourceK8sNodePool0().CoreConfigSchema().ImpliedType(),
				Upgrade: resourceK8sNodePoolUpgradeV0,
				Version: 0,
			},
		},
	}
}

func checkNodePoolImmutableFields(_ context.Context, diff *schema.ResourceDiff, _ interface{}) error {

	allowReplace := diff.Get("allow_replace").(bool)
	// allows the immutable fields to be updated
	if allowReplace {
		return nil
	}
	// we do not want to check in case of resource creation
	if diff.Id() == "" {
		return nil
	}
	if diff.HasChange("name") {
		return fmt.Errorf("name %s", ImmutableError)
	}

	if diff.HasChange("cpu_family") {
		return fmt.Errorf("cpu_family %s", ImmutableError)
	}

	if diff.HasChange("availability_zone") {
		return fmt.Errorf("availability_zone %s", ImmutableError)
	}

	if diff.HasChange("cores_count") {
		return fmt.Errorf("cores_count %s", ImmutableError)
	}

	if diff.HasChange("ram_size") {
		return fmt.Errorf("ram_size %s", ImmutableError)
	}

	if diff.HasChange("storage_size") {
		return fmt.Errorf("storage_size %s", ImmutableError)
	}

	if diff.HasChange("storage_type") {
		return fmt.Errorf("storage_type %s", ImmutableError)
	}

	if diff.HasChange("gateway_ip") {
		return fmt.Errorf("gateway_ip %s", ImmutableError)
	}
	return nil

}

func resourceK8sNodePool0() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"lans": {
				Type:        schema.TypeList,
				Description: "A list of Local Area Networks the node pool should be part of",
				Optional:    true,
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
			},
		},
		Timeouts: &resourceDefaultTimeouts,
	}
}

// resourceK8sNodePoolUpgradeV0 handles the differences that arise on lans when migrating from v5.X.X to a v6.X.X stable
// release and ignores the upgrade from a v6.0.0-beta.X, since the structure of lans is the same
func resourceK8sNodePoolUpgradeV0(_ context.Context, state map[string]interface{}, _ interface{}) (map[string]interface{}, error) {
	oldState := state
	var oldData []interface{}
	if d, ok := oldState["lans"].([]interface{}); ok {
		oldData = d
	}

	var lans []interface{}
	var floatType float64

	for _, lanId := range oldData {
		// this condition is for handling the migration from a v5.X.X to v6.X.X  release, when the content of lans property
		// is a list of floats. The content is mapped to the new v6.X.X lans structure
		if reflect.TypeOf(lanId) == reflect.TypeOf(floatType) {

			lanEntry := make(map[string]interface{})

			lanEntry["id"] = lanId

			// default value for dhcp
			lanEntry["dhcp"] = true

			var nodePoolRoutes []interface{}

			// empty list for routes
			lanEntry["routes"] = nodePoolRoutes
			lans = append(lans, lanEntry)
		} else {
			// this condition is for the migration from a v6.X.X-beta.X to v6.X.X  release, when no handling is necessary since the structure of lans is the same
			return state, nil
		}
	}

	state["lans"] = lans

	return state, nil
}

func getLanResourceData(lansList *schema.Set) []ionoscloud.KubernetesNodePoolLan {
	lans := make([]ionoscloud.KubernetesNodePoolLan, 0)
	if lansList.List() != nil {
		for _, lanItem := range lansList.List() {
			lanContent := lanItem.(map[string]interface{})
			lan := ionoscloud.KubernetesNodePoolLan{}

			if lanID, lanIdOk := lanContent["id"].(int); lanIdOk {
				log.Printf("[INFO] Adding LAN %v to node pool...", lanID)
				lanID := int32(lanID)
				lan.Id = lanID
			}

			if lanDhcp, lanDhcpOk := lanContent["dhcp"].(bool); lanDhcpOk {
				log.Printf("[INFO] Adding dhcp %v to node pool...", lanDhcp)
				lan.Dhcp = &lanDhcp
			}

			routes := make([]ionoscloud.KubernetesNodePoolLanRoutes, 0)

			if lanRoutes, lanRoutesOk := lanContent["routes"].(*schema.Set); lanRoutesOk {
				log.Printf("[INFO] Adding routes %v to node pool...", lanRoutes)
				if lanRoutes.List() != nil {
					for _, routeItem := range lanRoutes.List() {
						routeContent := routeItem.(map[string]interface{})
						route := ionoscloud.KubernetesNodePoolLanRoutes{}

						if routeNetwork, routeNewtworkOk := routeContent["network"].(string); routeNewtworkOk {
							route.Network = &routeNetwork
						}

						if routeGatewayIp, routeGatewayIpOk := routeContent["gateway_ip"].(string); routeGatewayIpOk {
							route.GatewayIp = &routeGatewayIp
						}

						routes = append(routes, route)

					}
				}
				log.Printf("[INFO] k8s node pool LanRoutes set to %+v", routes)
			}

			lan.Routes = routes
			lans = append(lans, lan)
		}
	}
	return lans
}

func getAutoscalingData(d *schema.ResourceData) (*ionoscloud.KubernetesAutoScaling, error) {
	var autoscaling ionoscloud.KubernetesAutoScaling

	asmnVal, asmnOk := d.GetOk("auto_scaling.0.min_node_count")
	asmxnVal, asmxnOk := d.GetOk("auto_scaling.0.max_node_count")

	if asmnOk && asmxnOk {
		asmnVal := int32(asmnVal.(int))
		asmxnVal := int32(asmxnVal.(int))
		if asmnVal == asmxnVal {
			return &autoscaling, fmt.Errorf("error creating k8s node pool: max_node_count cannot be equal to min_node_count")
		}

		if asmxnVal < asmnVal {
			return &autoscaling, fmt.Errorf("error creating k8s node pool: max_node_count cannot be lower than min_node_count")
		}

		log.Printf("[INFO] Setting Autoscaling minimum node count to : %d", asmnVal)
		autoscaling.MinNodeCount = asmnVal
		log.Printf("[INFO] Setting Autoscaling maximum node count to : %d", asmxnVal)
		autoscaling.MaxNodeCount = asmxnVal
	}

	return &autoscaling, nil
}
func resourcek8sNodePoolCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(services.SdkBundle).CloudApiClient

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
		Properties: ionoscloud.KubernetesNodePoolPropertiesForPost{
			AvailabilityZone: availabilityZone,
			CoresCount:       coresCount,
			CpuFamily:        cpuFamily,
			DatacenterId:     datacenterId,
			K8sVersion:       &k8sVersion,
			Name:             name,
			NodeCount:        nodeCount,
			RamSize:          ramSize,
			StorageSize:      storageSize,
			StorageType:      storageType,
		},
	}

	if _, mwOk := d.GetOk("maintenance_window.0"); mwOk {
		k8sNodepool.Properties.MaintenanceWindow = &ionoscloud.KubernetesMaintenanceWindow{}
	}

	if mtVal, mtOk := d.GetOk("maintenance_window.0.time"); mtOk {
		log.Printf("[INFO] Setting Maintenance Window Time to : %s", mtVal.(string))
		mtVal := mtVal.(string)
		k8sNodepool.Properties.MaintenanceWindow.Time = mtVal
	}

	if mdVal, mdOk := d.GetOk("maintenance_window.0.day_of_the_week"); mdOk {
		mdVal := mdVal.(string)
		k8sNodepool.Properties.MaintenanceWindow.DayOfTheWeek = mdVal
	}

	if autoscaling, err := getAutoscalingData(d); err != nil {
		return diag.FromErr(err)
	} else {
		k8sNodepool.Properties.AutoScaling = autoscaling
	}

	if k8sNodepool.Properties.AutoScaling != nil && k8sNodepool.Properties.NodeCount < k8sNodepool.Properties.AutoScaling.MinNodeCount {
		d.SetId("")
		diags := diag.FromErr(fmt.Errorf("error creating k8s node pool: node_count cannot be lower than min_node_count"))
		return diags
	}

	if lansVal, lansOK := d.GetOk("lans"); lansOK {
		lansList := lansVal.(*schema.Set)
		lans := getLanResourceData(lansList)
		k8sNodepool.Properties.Lans = lans
	}

	publicIpsProp, ok := d.GetOk("public_ips")
	if ok {
		publicIps := publicIpsProp.([]interface{})

		/* number of public IPs needs to be at least NodeCount + 1 */
		if len(publicIps) > 0 && int32(len(publicIps)) < k8sNodepool.Properties.NodeCount+1 {
			diags := diag.FromErr(fmt.Errorf("the number of public IPs must be at least %d", k8sNodepool.Properties.NodeCount+1))
			return diags
		}

		var requestPublicIps []string
		for i := range publicIps {
			requestPublicIps = append(requestPublicIps, fmt.Sprint(publicIps[i]))
		}
		k8sNodepool.Properties.PublicIps = requestPublicIps
	}

	// if gatewayIp, gatewayIpOk := d.GetOk("gateway_ip"); gatewayIpOk {
	//	gatewayIp := gatewayIp.(string)
	//	k8sNodepool.Properties.GatewayIp = &gatewayIp
	//}

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
		diags := diag.FromErr(fmt.Errorf("error creating k8s node pool: %w", err))
		return diags
	}

	d.SetId(*createdNodepool.Id)

	log.Printf("[INFO] Successfully created k8s node pool: %s", d.Id())

	for {
		log.Printf("[INFO] Waiting for k8s node pool %s to be ready...", d.Id())

		nodepoolReady, rsErr := k8sNodepoolReady(ctx, client, d)
		if rsErr != nil {
			diags := diag.FromErr(fmt.Errorf("error while checking readiness status of k8s node pool %s: %w", d.Id(), rsErr))
			return diags
		}

		if nodepoolReady {
			log.Printf("[INFO] k8s node pool ready: %s", d.Id())
			break
		}

		select {
		case <-time.After(constant.SleepInterval):
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
	client := meta.(services.SdkBundle).CloudApiClient

	k8sNodepool, apiResponse, err := client.KubernetesApi.K8sNodepoolsFindById(ctx, d.Get("k8s_cluster_id").(string), d.Id()).Execute()
	logApiRequestTime(apiResponse)

	if err != nil {
		log.Printf("[INFO] Resource %s not found: %+v", d.Id(), err)
		if httpNotFound(apiResponse) {
			d.SetId("")
			return nil
		}
		diags := diag.FromErr(err)
		return diags
	}

	log.Printf("[INFO] Successfully retrieved k8s node pool %s: %+v", d.Id(), k8sNodepool)

	if err := setK8sNodePoolData(d, &k8sNodepool); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourcek8sNodePoolUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(services.SdkBundle).CloudApiClient

	request := ionoscloud.KubernetesNodePoolForPut{}

	nodeCount := int32(d.Get("node_count").(int))

	request.Properties = ionoscloud.KubernetesNodePoolPropertiesForPut{
		NodeCount: nodeCount,
	}

	if d.HasChange("node_count") {
		oldNc, newNc := d.GetChange("node_count")
		log.Printf("[INFO] k8s node pool node_count changed from %+v to %+v", oldNc, newNc)
	}

	k8sVersion := d.Get("k8s_version").(string)
	request.Properties.K8sVersion = &k8sVersion
	if d.HasChange("k8s_version") {
		oldk8sVersion, newk8sVersion := d.GetChange("k8s_version")
		log.Printf("[INFO] k8s pool k8s version changed from %+v to %+v", oldk8sVersion, newk8sVersion)
	}

	if autoscaling, err := getAutoscalingData(d); err != nil {
		return diag.FromErr(err)
	} else {
		request.Properties.AutoScaling = autoscaling
	}

	if request.Properties.AutoScaling != nil && request.Properties.NodeCount < request.Properties.AutoScaling.MinNodeCount {
		d.SetId("")
		diags := diag.FromErr(fmt.Errorf("error creating k8s node pool: node_count cannot be lower than min_node_count"))
		return diags
	}

	if d.HasChange("auto_scaling.0.min_node_count") {
		oldMinNodes, newMinNodes := d.GetChange("auto_scaling.0.min_node_count")
		log.Printf("[INFO] k8s node pool autoscaling min # of nodes changed from %+v to %+v", oldMinNodes, newMinNodes)
	}

	if d.HasChange("auto_scaling.0.max_node_count") {
		oldMaxNodes, newMaxNodes := d.GetChange("auto_scaling.0.max_node_count")
		log.Printf("[INFO] k8s node pool autoscaling max # of nodes changed from %+v to %+v", oldMaxNodes, newMaxNodes)
	}

	if d.HasChange("lans") {
		oldLANs, newLANs := d.GetChange("lans")
		lansList := newLANs.(*schema.Set)
		lans := getLanResourceData(lansList)
		log.Printf("[INFO] k8s node pool LANs changed from %+v to %+v", oldLANs, newLANs)
		request.Properties.Lans = lans
	}

	if d.HasChange("maintenance_window.0") {

		_, newMw := d.GetChange("maintenance_window.0")

		if newMw.(map[string]interface{}) != nil {

			updateMaintenanceWindow := false
			dayOfTheWeek := d.Get("maintenance_window.0.day_of_the_week").(string)
			timeS := d.Get("maintenance_window.0.time").(string)
			maintenanceWindow := &ionoscloud.KubernetesMaintenanceWindow{
				DayOfTheWeek: dayOfTheWeek,
				Time:         timeS,
			}

			if d.HasChange("maintenance_window.0.day_of_the_week") {

				oldMd, newMd := d.GetChange("maintenance_window.0.day_of_the_week")
				if newMd.(string) != "" {
					log.Printf("[INFO] k8s node pool maintenance window DOW changed from %+v to %+v", oldMd, newMd)
					updateMaintenanceWindow = true
					newMd := newMd.(string)
					maintenanceWindow.DayOfTheWeek = newMd
				}
			}

			if d.HasChange("maintenance_window.0.time") {
				oldMt, newMt := d.GetChange("maintenance_window.0.time")
				if newMt.(string) != "" {
					log.Printf("[INFO] k8s node pool maintenance window time changed from %+v to %+v", oldMt, newMt)
					updateMaintenanceWindow = true
					newMt := newMt.(string)
					maintenanceWindow.Time = newMt
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
			if len(publicIps) > 0 && int32(len(publicIps)) < request.Properties.NodeCount+1 {
				diags := diag.FromErr(fmt.Errorf("the number of public IPs must be at least %d", request.Properties.NodeCount+1))
				return diags
			}

			for _, ip := range publicIps {
				requestPublicIps = append(requestPublicIps, ip.(string))
			}

		}
		request.Properties.PublicIps = requestPublicIps

	}

	if d.HasChange("gateway_ip") {
		diags := diag.FromErr(fmt.Errorf("gateway_ip attribute is immutable, therefore not allowed in update requests"))
		return diags
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
		if httpNotFound(apiResponse) {
			d.SetId("")
			return nil
		}
		diags := diag.FromErr(fmt.Errorf("error while updating k8s node pool %s: %w", d.Id(), err))
		return diags
	}

	for {
		log.Printf("[INFO] Waiting for k8s node pool %s to be ready...", d.Id())

		nodepoolReady, rsErr := k8sNodepoolReady(ctx, client, d)

		if rsErr != nil {
			diags := diag.FromErr(fmt.Errorf("error while checking readiness status of k8s node pool %s: %w", d.Id(), rsErr))
			return diags
		}

		if nodepoolReady {
			log.Printf("[INFO] k8s node pool ready: %s", d.Id())
			break
		}

		select {
		case <-time.After(constant.SleepInterval):
			log.Printf("[DEBUG] retrying ...")
		case <-ctx.Done():
			diags := diag.FromErr(fmt.Errorf("k8s node pool update timed out! WARNING: your k8s node pool will still probably be updated after some time but the terraform state won't reflect that; check your Ionos Cloud account for updates"))
			return diags
		}
	}

	return resourcek8sNodePoolRead(ctx, d, meta)
}

func resourcek8sNodePoolDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(services.SdkBundle).CloudApiClient

	apiResponse, err := client.KubernetesApi.K8sNodepoolsDelete(ctx, d.Get("k8s_cluster_id").(string), d.Id()).Execute()
	logApiRequestTime(apiResponse)

	if err != nil {
		if httpNotFound(apiResponse) {
			d.SetId("")
			return nil
		}
		diags := diag.FromErr(fmt.Errorf("error while deleting k8s node pool %s: %w", d.Id(), err))
		return diags
	}

	for {
		log.Printf("[INFO] Waiting for k8s node pool %s to be deleted...", d.Id())

		nodepoolDeleted, dsErr := k8sNodepoolDeleted(ctx, client, d)
		if dsErr != nil {
			diags := diag.FromErr(fmt.Errorf("error while checking deletion status of k8s node pool %s: %w", d.Id(), dsErr))
			return diags
		}

		if nodepoolDeleted {
			log.Printf("[INFO] Successfully deleted k8s node pool: %s", d.Id())
			break
		}

		select {
		case <-time.After(constant.SleepInterval):
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

	client := meta.(services.SdkBundle).CloudApiClient
	k8sNodepool, apiResponse, err := client.KubernetesApi.K8sNodepoolsFindById(ctx, clusterId, npId).Execute()
	logApiRequestTime(apiResponse)

	if err != nil {
		if _, ok := err.(shared.GenericOpenAPIError); ok {
			if httpNotFound(apiResponse) {
				d.SetId("")
				return nil, fmt.Errorf("unable to find k8s node pool %q", npId)
			}
		}
		return nil, fmt.Errorf("unable to retrieve k8s node pool %q, error:%w", npId, err)
	}

	log.Printf("[INFO] K8s node pool found: %+v", k8sNodepool)

	if err := d.Set("k8s_cluster_id", clusterId); err != nil {
		return nil, fmt.Errorf("error while setting k8s_cluster_id property for k8s node pool %q: %q", npId, err)
	}

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

	if err := d.Set("name", nodePool.Properties.Name); err != nil {
		return err
	}

	if err := d.Set("datacenter_id", nodePool.Properties.DatacenterId); err != nil {
		return err
	}

	if err := d.Set("node_count", nodePool.Properties.NodeCount); err != nil {
		return err
	}

	if err := d.Set("cpu_family", nodePool.Properties.CpuFamily); err != nil {
		return err
	}

	if err := d.Set("cores_count", nodePool.Properties.CoresCount); err != nil {
		return err
	}

	if err := d.Set("ram_size", nodePool.Properties.RamSize); err != nil {
		return err
	}

	if err := d.Set("availability_zone", nodePool.Properties.AvailabilityZone); err != nil {
		return err
	}

	if err := d.Set("storage_type", nodePool.Properties.StorageType); err != nil {
		return err
	}

	if err := d.Set("storage_size", nodePool.Properties.StorageSize); err != nil {
		return err
	}

	if nodePool.Properties.K8sVersion != nil {
		if err := d.Set("k8s_version", *nodePool.Properties.K8sVersion); err != nil {
			return err
		}
	}

	if len(nodePool.Properties.PublicIps) > 0 {
		if err := d.Set("public_ips", nodePool.Properties.PublicIps); err != nil {
			return err
		}
	}

	if nodePool.Properties.MaintenanceWindow != nil {
		if err := d.Set("maintenance_window", []map[string]string{
			{
				"time":            nodePool.Properties.MaintenanceWindow.Time,
				"day_of_the_week": nodePool.Properties.MaintenanceWindow.DayOfTheWeek,
			},
		}); err != nil {
			return err
		}
	}

	if nodePool.Properties.AutoScaling != nil && (nodePool.Properties.AutoScaling.MinNodeCount != 0 &&
		nodePool.Properties.AutoScaling.MaxNodeCount != 0) {
		if err := d.Set("auto_scaling", []map[string]uint32{
			{
				"min_node_count": uint32(nodePool.Properties.AutoScaling.MinNodeCount),
				"max_node_count": uint32(nodePool.Properties.AutoScaling.MaxNodeCount),
			},
		}); err != nil {
			return err
		}
	}

	if len(nodePool.Properties.Lans) > 0 {

		nodePoolLans := getK8sNodePoolLans(nodePool.Properties.Lans)

		if err := d.Set("lans", nodePoolLans); err != nil {
			return fmt.Errorf("error while setting lans property for k8sNodepool %s: %w", d.Id(), err)
		}

	}

	// if nodePool.Properties.GatewayIp != nil {
	//	if err := d.Set("gateway_ip", *nodePool.Properties.GatewayIp); err != nil {
	//		return fmt.Errorf("error while setting gateway_ip property for nodepool %s: %w", d.Id(), err)
	//	}
	//}

	labels := make(map[string]interface{})
	if nodePool.Properties.Labels != nil && len(*nodePool.Properties.Labels) > 0 {
		for k, v := range *nodePool.Properties.Labels {
			labels[k] = v
		}
	}

	if err := d.Set("labels", labels); err != nil {
		return fmt.Errorf("error while setting the labels property for k8sNodepool %s: %w", d.Id(), err)

	}

	annotations := make(map[string]interface{})
	if nodePool.Properties.Annotations != nil && len(*nodePool.Properties.Annotations) > 0 {
		for k, v := range *nodePool.Properties.Annotations {
			annotations[k] = v
		}
	}

	if err := d.Set("annotations", annotations); err != nil {
		return fmt.Errorf("error while setting the annotations property for k8sNodepool %s: %w", d.Id(), err)
	}

	return nil
}
func k8sNodepoolReady(ctx context.Context, client *ionoscloud.APIClient, d *schema.ResourceData) (bool, error) {
	resource, apiResponse, err := client.KubernetesApi.K8sNodepoolsFindById(ctx, d.Get("k8s_cluster_id").(string), d.Id()).Execute()
	logApiRequestTime(apiResponse)
	if err != nil {
		return true, fmt.Errorf("error checking k8s node pool status: %w", err)
	}
	if resource.Metadata == nil || resource.Metadata.State == nil {
		return false, fmt.Errorf("error while checking k8s node pool status: state is nil")
	}
	if utils.IsStateFailed(*resource.Metadata.State) {
		return false, fmt.Errorf("error while checking if k8s nodepool is ready %s, state %s", *resource.Id, *resource.Metadata.State)
	}
	log.Printf("[INFO] k8s nodepool state: %s", *resource.Metadata.State)
	// k8s is the only resource that has a state of ACTIVE when it is ready
	return strings.EqualFold(*resource.Metadata.State, ionoscloud.Active), nil
}

func k8sNodepoolDeleted(ctx context.Context, client *ionoscloud.APIClient, d *schema.ResourceData) (bool, error) {
	resource, apiResponse, err := client.KubernetesApi.K8sNodepoolsFindById(ctx, d.Get("k8s_cluster_id").(string), d.Id()).Execute()
	logApiRequestTime(apiResponse)

	if err != nil {
		if httpNotFound(apiResponse) {
			return true, nil
		}
		return true, fmt.Errorf("error checking k8s node pool deletion status: %w", err)
	}
	if resource.Metadata != nil && resource.Metadata.State != nil {
		if utils.IsStateFailed(*resource.Metadata.State) {
			return false, fmt.Errorf("error while checking if k8s nodepool is properly deleted, nodepool ID: %s, state: %s", *resource.Id, *resource.Metadata.State)
		}
	}
	return false, nil
}

func getK8sNodePoolLans(lans []ionoscloud.KubernetesNodePoolLan) []interface{} {

	var nodePoolLans []interface{}
	for _, nodePoolLan := range lans {
		lanEntry := make(map[string]interface{})

		lanEntry["id"] = nodePoolLan.Id

		if nodePoolLan.Dhcp != nil {
			lanEntry["dhcp"] = *nodePoolLan.Dhcp
		}

		if len(nodePoolLan.Routes) > 0 {
			var nodePoolRoutes []interface{}
			for _, nodePoolRoute := range nodePoolLan.Routes {
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
