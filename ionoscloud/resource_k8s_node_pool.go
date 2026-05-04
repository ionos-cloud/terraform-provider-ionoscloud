package ionoscloud

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/bundleclient"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils/constant"
	diagutil "github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils/diags"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	ionoscloud "github.com/ionos-cloud/sdk-go/v6"
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
				ForceNew:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.IsUUID),
			},
			"k8s_cluster_id": {
				Type:             schema.TypeString,
				Description:      "The UUID of an existing kubernetes cluster",
				Required:         true,
				ForceNew:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.IsUUID),
			},
			"location": {
				Type:        schema.TypeString,
				Description: "The location of the resource. This field should be used only if you are also using a file configuration and should not be configured otherwise.",
				Optional:    true,
				ForceNew:    true,
			},
			"cpu_family": {
				Type:             schema.TypeString,
				Description:      "CPU Family",
				Optional:         true,
				ForceNew:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringIsNotWhiteSpace),
			},
			"server_type": {
				Type:             schema.TypeString,
				Description:      "The server type for the compute engine",
				Optional:         true,
				Default:          "DedicatedCore",
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"DedicatedCore", "VCPU"}, false)),
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
			// },
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

func checkNodePoolImmutableFields(_ context.Context, diff *schema.ResourceDiff, _ any) error {

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
func resourceK8sNodePoolUpgradeV0(_ context.Context, state map[string]any, _ any) (map[string]any, error) {
	oldState := state
	var oldData []any
	if d, ok := oldState["lans"].([]any); ok {
		oldData = d
	}

	var lans []any

	for _, lanId := range oldData {
		// this condition is for handling the migration from a v5.X.X to v6.X.X  release, when the content of lans property
		// is a list of floats. The content is mapped to the new v6.X.X lans structure
		if reflect.TypeOf(lanId) == reflect.TypeFor[float64]() {

			lanEntry := make(map[string]any)

			lanEntry["id"] = lanId

			// default value for dhcp
			lanEntry["dhcp"] = true

			var nodePoolRoutes []any

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

func getLanResourceData(ctx context.Context, lansList *schema.Set) []ionoscloud.KubernetesNodePoolLan {
	lans := make([]ionoscloud.KubernetesNodePoolLan, 0)
	if lansList.List() != nil {
		for _, lanItem := range lansList.List() {
			lanContent := lanItem.(map[string]any)
			lan := ionoscloud.KubernetesNodePoolLan{}

			if lanID, lanIdOk := lanContent["id"].(int); lanIdOk {
				tflog.Info(ctx, "adding LAN to node pool", map[string]interface{}{"lan_id": lanID})
				lanID := int32(lanID)
				lan.Id = &lanID
			}

			if lanDhcp, lanDhcpOk := lanContent["dhcp"].(bool); lanDhcpOk {
				tflog.Info(ctx, "adding dhcp to node pool", map[string]interface{}{"dhcp": lanDhcp})
				lan.Dhcp = &lanDhcp
			}

			routes := make([]ionoscloud.KubernetesNodePoolLanRoutes, 0)

			if lanRoutes, lanRoutesOk := lanContent["routes"].(*schema.Set); lanRoutesOk {
				tflog.Info(ctx, "adding routes to node pool", map[string]interface{}{"route_count": lanRoutes.Len()})
				if lanRoutes.List() != nil {
					for _, routeItem := range lanRoutes.List() {
						routeContent := routeItem.(map[string]any)
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
				tflog.Info(ctx, "node pool LAN routes set", map[string]interface{}{"route_count": len(routes)})
			}

			lan.Routes = &routes
			lans = append(lans, lan)
		}
	}
	return lans
}

func getAutoscalingData(ctx context.Context, d *schema.ResourceData) (*ionoscloud.KubernetesAutoScaling, error) {
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

		tflog.Info(ctx, "setting autoscaling min node count", map[string]interface{}{"min_node_count": asmnVal})
		autoscaling.MinNodeCount = &asmnVal
		tflog.Info(ctx, "setting autoscaling max node count", map[string]interface{}{"max_node_count": asmxnVal})
		autoscaling.MaxNodeCount = &asmxnVal
	}

	return &autoscaling, nil
}

func resourcek8sNodePoolCreate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	location := d.Get("location").(string)
	client, err := meta.(bundleclient.SdkBundle).NewCloudAPIClient(ctx, location)
	if err != nil {
		return diag.FromErr(err)
	}

	name := d.Get("name").(string)
	datacenterId := d.Get("datacenter_id").(string)
	k8sVersion := d.Get("k8s_version").(string)
	availabilityZone := d.Get("availability_zone").(string)
	storageType := d.Get("storage_type").(string)
	nodeCount := int32(d.Get("node_count").(int))
	coresCount := int32(d.Get("cores_count").(int))
	storageSize := int32(d.Get("storage_size").(int))
	ramSize := int32(d.Get("ram_size").(int))

	k8sNodepool := ionoscloud.KubernetesNodePoolForPost{
		Properties: &ionoscloud.KubernetesNodePoolPropertiesForPost{
			AvailabilityZone: &availabilityZone,
			CoresCount:       &coresCount,
			DatacenterId:     &datacenterId,
			K8sVersion:       &k8sVersion,
			Name:             &name,
			NodeCount:        &nodeCount,
			RamSize:          &ramSize,
			StorageSize:      &storageSize,
			StorageType:      &storageType,
		},
	}

	if serverType, serverTypeOk := d.GetOk("server_type"); serverTypeOk {
		serverType := ionoscloud.KubernetesNodePoolServerType(serverType.(string))
		k8sNodepool.Properties.ServerType = &serverType
	}

	if cpuFamily, cpuFamilyOk := d.GetOk("cpu_family"); cpuFamilyOk {
		cpuFamily := cpuFamily.(string)
		k8sNodepool.Properties.CpuFamily = &cpuFamily
	}

	if _, mwOk := d.GetOk("maintenance_window.0"); mwOk {
		k8sNodepool.Properties.MaintenanceWindow = &ionoscloud.KubernetesMaintenanceWindow{}
	}

	if mtVal, mtOk := d.GetOk("maintenance_window.0.time"); mtOk {
		tflog.Info(ctx, "setting maintenance window time", map[string]interface{}{"time": mtVal.(string)})
		mtVal := mtVal.(string)
		k8sNodepool.Properties.MaintenanceWindow.Time = &mtVal
	}

	if mdVal, mdOk := d.GetOk("maintenance_window.0.day_of_the_week"); mdOk {
		mdVal := mdVal.(string)
		k8sNodepool.Properties.MaintenanceWindow.DayOfTheWeek = &mdVal
	}

	if autoscaling, err := getAutoscalingData(ctx, d); err != nil {
		return diagutil.ToDiags(d, err, nil)
	} else {
		k8sNodepool.Properties.AutoScaling = autoscaling
	}

	if k8sNodepool.Properties.AutoScaling != nil && k8sNodepool.Properties.AutoScaling.MinNodeCount != nil && *k8sNodepool.Properties.NodeCount < *k8sNodepool.Properties.AutoScaling.MinNodeCount {
		d.SetId("")
		return diagutil.ToDiags(d, fmt.Errorf("error creating k8s node pool: node_count cannot be lower than min_node_count"), nil)
	}

	if lansVal, lansOK := d.GetOk("lans"); lansOK {
		lansList := lansVal.(*schema.Set)
		lans := getLanResourceData(ctx, lansList)
		k8sNodepool.Properties.Lans = &lans
	}

	publicIpsProp, ok := d.GetOk("public_ips")
	if ok {
		publicIps := publicIpsProp.([]any)

		/* number of public IPs needs to be at least NodeCount + 1 */
		if len(publicIps) > 0 && int32(len(publicIps)) < *k8sNodepool.Properties.NodeCount+1 {
			return diagutil.ToDiags(d, fmt.Errorf("the number of public IPs must be at least %d", *k8sNodepool.Properties.NodeCount+1), nil)
		}

		var requestPublicIps []string
		for i := range publicIps {
			requestPublicIps = append(requestPublicIps, fmt.Sprint(publicIps[i]))
		}
		k8sNodepool.Properties.PublicIps = &requestPublicIps
	}

	// if gatewayIp, gatewayIpOk := d.GetOk("gateway_ip"); gatewayIpOk {
	//	gatewayIp := gatewayIp.(string)
	//	k8sNodepool.Properties.GatewayIp = &gatewayIp
	// }

	labelsProp, ok := d.GetOk("labels")
	if ok {
		labels := make(map[string]string)
		for k, v := range labelsProp.(map[string]any) {
			labels[k] = v.(string)
		}
		k8sNodepool.Properties.Labels = &labels
	}

	annotationsProp, ok := d.GetOk("annotations")
	if ok {
		annotations := make(map[string]string)
		for k, v := range annotationsProp.(map[string]any) {
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
		return diagutil.ToDiags(d, fmt.Errorf("error creating k8s node pool: %w", err), &diagutil.ErrorContext{StatusCode: apiResponse.SafeStatusCode()})
	}

	d.SetId(*createdNodepool.Id)

	tflog.Info(ctx, "created k8s node pool", map[string]interface{}{"node_pool_id": d.Id()})

	for {
		tflog.Info(ctx, "waiting for k8s node pool to be ready", map[string]interface{}{"node_pool_id": d.Id()})

		nodepoolReady, rsErr := k8sNodepoolReady(ctx, client, d)
		if rsErr != nil {
			return diagutil.ToDiags(d, fmt.Errorf("error while checking readiness status of k8s node pool: %w", rsErr), nil)
		}

		if nodepoolReady {
			tflog.Info(ctx, "k8s node pool ready", map[string]interface{}{"node_pool_id": d.Id()})
			break
		}

		select {
		case <-time.After(constant.SleepInterval):
			tflog.Info(ctx, "k8s node pool not ready, retrying")
		case <-ctx.Done():
			tflog.Info(ctx, "k8s node pool creation timed out")
			return diagutil.ToDiags(d, fmt.Errorf("k8s creation timed out"), nil)
		}
	}

	return resourcek8sNodePoolRead(ctx, d, meta)
}

func resourcek8sNodePoolRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	location := d.Get("location").(string)
	client, err := meta.(bundleclient.SdkBundle).NewCloudAPIClient(ctx, location)
	if err != nil {
		return diag.FromErr(err)
	}

	k8sNodepool, apiResponse, err := client.KubernetesApi.K8sNodepoolsFindById(ctx, d.Get("k8s_cluster_id").(string), d.Id()).Execute()
	logApiRequestTime(apiResponse)

	if err != nil {
		tflog.Info(ctx, "k8s node pool not found", map[string]interface{}{"node_pool_id": d.Id(), "error": err.Error()})
		if httpNotFound(apiResponse) {
			d.SetId("")
			return nil
		}
		return diagutil.ToDiags(d, err, &diagutil.ErrorContext{StatusCode: apiResponse.SafeStatusCode()})
	}

	tflog.Info(ctx, "retrieved k8s node pool", map[string]interface{}{"node_pool_id": d.Id()})

	if err := setK8sNodePoolData(d, &k8sNodepool); err != nil {
		return diagutil.ToDiags(d, err, nil)
	}

	return nil
}

func resourcek8sNodePoolUpdate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	location := d.Get("location").(string)
	client, err := meta.(bundleclient.SdkBundle).NewCloudAPIClient(ctx, location)
	if err != nil {
		return diag.FromErr(err)
	}

	request := ionoscloud.KubernetesNodePoolForPut{}

	nodeCount := int32(d.Get("node_count").(int))

	request.Properties = &ionoscloud.KubernetesNodePoolPropertiesForPut{
		NodeCount: &nodeCount,
	}

	if d.HasChange("node_count") {
		oldNc, newNc := d.GetChange("node_count")
		tflog.Info(ctx, "node pool node_count changed", map[string]interface{}{"old": oldNc, "new": newNc})
	}

	k8sVersion := d.Get("k8s_version").(string)
	request.Properties.K8sVersion = &k8sVersion
	if d.HasChange("k8s_version") {
		oldk8sVersion, newk8sVersion := d.GetChange("k8s_version")
		tflog.Info(ctx, "node pool k8s version changed", map[string]interface{}{"old": oldk8sVersion, "new": newk8sVersion})
	}

	if autoscaling, err := getAutoscalingData(ctx, d); err != nil {
		return diagutil.ToDiags(d, err, nil)
	} else {
		request.Properties.AutoScaling = autoscaling
	}

	if request.Properties.AutoScaling != nil && request.Properties.AutoScaling.MinNodeCount != nil && *request.Properties.NodeCount < *request.Properties.AutoScaling.MinNodeCount {
		d.SetId("")
		return diagutil.ToDiags(d, fmt.Errorf("error creating k8s node pool: node_count cannot be lower than min_node_count"), nil)
	}

	if d.HasChange("auto_scaling.0.min_node_count") {
		oldMinNodes, newMinNodes := d.GetChange("auto_scaling.0.min_node_count")
		tflog.Info(ctx, "node pool autoscaling min node count changed", map[string]interface{}{"old": oldMinNodes, "new": newMinNodes})
	}

	if d.HasChange("auto_scaling.0.max_node_count") {
		oldMaxNodes, newMaxNodes := d.GetChange("auto_scaling.0.max_node_count")
		tflog.Info(ctx, "node pool autoscaling max node count changed", map[string]interface{}{"old": oldMaxNodes, "new": newMaxNodes})
	}

	if d.HasChange("lans") {
		oldLANs, newLANs := d.GetChange("lans")
		lansList := newLANs.(*schema.Set)
		lans := getLanResourceData(ctx, lansList)
		tflog.Info(ctx, "node pool LANs changed", map[string]interface{}{"old": oldLANs, "new": newLANs})
		request.Properties.Lans = &lans
	}

	if d.HasChange("maintenance_window.0") {

		_, newMw := d.GetChange("maintenance_window.0")

		if newMw.(map[string]any) != nil {

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
					tflog.Info(ctx, "node pool maintenance window DOW changed", map[string]interface{}{"old": oldMd, "new": newMd})
					updateMaintenanceWindow = true
					newMd := newMd.(string)
					maintenanceWindow.DayOfTheWeek = &newMd
				}
			}

			if d.HasChange("maintenance_window.0.time") {
				oldMt, newMt := d.GetChange("maintenance_window.0.time")
				if newMt.(string) != "" {
					tflog.Info(ctx, "node pool maintenance window time changed", map[string]interface{}{"old": oldMt, "new": newMt})
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

	if d.HasChange("server_type") {
		oldServerType, newServerType := d.GetChange("server_type")
		tflog.Info(ctx, "node pool server type changed", map[string]interface{}{"old": oldServerType, "new": newServerType})

		serverType := ionoscloud.KubernetesNodePoolServerType(newServerType.(string))
		request.Properties.ServerType = &serverType
	}

	if d.HasChange("public_ips") {
		oldPublicIps, newPublicIps := d.GetChange("public_ips")
		tflog.Info(ctx, "node pool public IPs changed", map[string]interface{}{"old": oldPublicIps, "new": newPublicIps})
		requestPublicIps := make([]string, 0)

		if newPublicIps != nil {

			publicIps := newPublicIps.([]any)

			/* number of public IPs needs to be at least NodeCount + 1 */
			if len(publicIps) > 0 && int32(len(publicIps)) < *request.Properties.NodeCount+1 {
				return diagutil.ToDiags(d, fmt.Errorf("the number of public IPs must be at least %d", *request.Properties.NodeCount+1), nil)
			}

			for _, ip := range publicIps {
				requestPublicIps = append(requestPublicIps, ip.(string))
			}

		}
		request.Properties.PublicIps = &requestPublicIps

	}

	if d.HasChange("gateway_ip") {
		return diagutil.ToDiags(d, fmt.Errorf("gateway_ip attribute is immutable, therefore not allowed in update requests"), nil)
	}

	if d.HasChange("labels") {
		oldLabels, newLabels := d.GetChange("labels")
		tflog.Info(ctx, "node pool labels changed", map[string]interface{}{"old": oldLabels, "new": newLabels})
		labels := make(map[string]string)
		if newLabels != nil {
			for k, v := range newLabels.(map[string]any) {
				labels[k] = v.(string)
			}
		}
		request.Properties.Labels = &labels
	}

	if d.HasChange("annotations") {
		oldAnnotations, newAnnotations := d.GetChange("annotations")
		tflog.Info(ctx, "node pool annotations changed", map[string]interface{}{"old": oldAnnotations, "new": newAnnotations})
		annotations := make(map[string]string)
		if newAnnotations != nil {
			for k, v := range newAnnotations.(map[string]any) {
				annotations[k] = v.(string)
			}
		}
		request.Properties.Annotations = &annotations
	}

	b, jErr := json.Marshal(request)

	if jErr == nil {
		tflog.Info(ctx, "node pool update request", map[string]interface{}{"request": string(b)})
	}

	_, apiResponse, err := client.KubernetesApi.K8sNodepoolsPut(ctx, d.Get("k8s_cluster_id").(string), d.Id()).KubernetesNodePool(request).Execute()
	logApiRequestTime(apiResponse)

	if err != nil {
		if httpNotFound(apiResponse) {
			d.SetId("")
			return nil
		}
		return diagutil.ToDiags(d, fmt.Errorf("error while updating k8s node pool: %w", err), &diagutil.ErrorContext{StatusCode: apiResponse.SafeStatusCode()})
	}

	for {
		tflog.Info(ctx, "waiting for k8s node pool to be ready", map[string]interface{}{"node_pool_id": d.Id()})

		nodepoolReady, rsErr := k8sNodepoolReady(ctx, client, d)

		if rsErr != nil {
			return diagutil.ToDiags(d, fmt.Errorf("error while checking readiness status of k8s node pool: %w", rsErr), nil)
		}

		if nodepoolReady {
			tflog.Info(ctx, "k8s node pool ready", map[string]interface{}{"node_pool_id": d.Id()})
			break
		}

		select {
		case <-time.After(constant.SleepInterval):
			tflog.Debug(ctx, "k8s node pool not ready, retrying")
		case <-ctx.Done():
			return diagutil.ToDiags(d, fmt.Errorf("k8s node pool update timed out! WARNING: your k8s node pool will still probably be updated after some time but the terraform state won't reflect that; check your Ionos Cloud account for updates"), nil)
		}
	}

	return resourcek8sNodePoolRead(ctx, d, meta)
}

func resourcek8sNodePoolDelete(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	location := d.Get("location").(string)
	client, err := meta.(bundleclient.SdkBundle).NewCloudAPIClient(ctx, location)
	if err != nil {
		return diag.FromErr(err)
	}

	apiResponse, err := client.KubernetesApi.K8sNodepoolsDelete(ctx, d.Get("k8s_cluster_id").(string), d.Id()).Execute()
	logApiRequestTime(apiResponse)

	if err != nil {
		if httpNotFound(apiResponse) {
			d.SetId("")
			return nil
		}
		return diagutil.ToDiags(d, fmt.Errorf("error while deleting k8s node pool: %w", err), &diagutil.ErrorContext{StatusCode: apiResponse.SafeStatusCode()})
	}

	for {
		tflog.Info(ctx, "waiting for k8s node pool to be deleted", map[string]interface{}{"node_pool_id": d.Id()})

		nodepoolDeleted, dsErr := k8sNodepoolDeleted(ctx, client, d)
		if dsErr != nil {
			return diagutil.ToDiags(d, fmt.Errorf("error while checking deletion status of k8s node pool: %w", dsErr), nil)
		}

		if nodepoolDeleted {
			tflog.Info(ctx, "successfully deleted k8s node pool", map[string]interface{}{"node_pool_id": d.Id()})
			break
		}

		select {
		case <-time.After(constant.SleepInterval):
			tflog.Debug(ctx, "k8s node pool not yet deleted, retrying")
		case <-ctx.Done():
			return diagutil.ToDiags(d, fmt.Errorf("k8s node pool deletion timed out! WARNING: your k8s node pool will still probably be deleted after some time but the terraform state won't reflect that; check your Ionos Cloud account for updates"), nil)
		}
	}

	d.SetId("")
	return nil
}

func resourceK8sNodepoolImport(ctx context.Context, d *schema.ResourceData, meta any) ([]*schema.ResourceData, error) {
	importID := d.Id()
	location, parts := splitImportID(importID, "/")

	if len(parts) != 2 {
		return nil, fmt.Errorf(
			"invalid import identifier: expected one of <location>:<k8s-cluster-id>/<k8s-nodepool-id> "+
				"or <k8s-cluster-id>/<k8s-nodepool-id>, got: %s", importID,
		)
	}

	if err := validateImportIDParts(parts); err != nil {
		return nil, diagutil.ToError(d, fmt.Errorf("failed validating import identifier %q: %w", importID, err), nil)
	}

	clusterId := parts[0]
	npId := parts[1]

	client, err := meta.(bundleclient.SdkBundle).NewCloudAPIClient(ctx, location)
	if err != nil {
		return nil, err
	}
	k8sNodepool, apiResponse, err := client.KubernetesApi.K8sNodepoolsFindById(ctx, clusterId, npId).Execute()
	logApiRequestTime(apiResponse)

	if err != nil {
		if _, ok := err.(ionoscloud.GenericOpenAPIError); ok {
			if httpNotFound(apiResponse) {
				d.SetId("")
				return nil, diagutil.ToError(d, fmt.Errorf("unable to find k8s node pool %q", npId), &diagutil.ErrorContext{StatusCode: apiResponse.SafeStatusCode()})
			}
		}
		return nil, diagutil.ToError(d, fmt.Errorf("unable to retrieve k8s node pool %q, error:%w", npId, err), &diagutil.ErrorContext{StatusCode: apiResponse.SafeStatusCode()})
	}

	tflog.Info(ctx, "k8s node pool imported", map[string]interface{}{"node_pool_id": npId, "cluster_id": clusterId})

	if err := d.Set("k8s_cluster_id", clusterId); err != nil {
		return nil, diagutil.ToError(d, fmt.Errorf("error while setting k8s_cluster_id property for k8s node pool %q: %w", npId, err), nil)
	}
	if err := d.Set("location", location); err != nil {
		return nil, diagutil.ToError(d, fmt.Errorf("error while setting location property for k8s node pool %q: %w", npId, err), nil)
	}

	if err := setK8sNodePoolData(d, &k8sNodepool); err != nil {
		return nil, diagutil.ToError(d, err, nil)
	}
	tflog.Info(ctx, "importing k8s node pool", map[string]interface{}{"node_pool_id": d.Id()})

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

		if nodePool.Properties.ServerType != nil {
			if err := d.Set("server_type", *nodePool.Properties.ServerType); err != nil {
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
				return fmt.Errorf("error while setting lans property for k8sNodepool %s: %w", d.Id(), err)
			}

		}

		// if nodePool.Properties.GatewayIp != nil {
		//	if err := d.Set("gateway_ip", *nodePool.Properties.GatewayIp); err != nil {
		//		return fmt.Errorf("error while setting gateway_ip property for nodepool %s: %w", d.Id(), err)
		//	}
		// }

		labels := make(map[string]any)
		if nodePool.Properties.Labels != nil && len(*nodePool.Properties.Labels) > 0 {
			for k, v := range *nodePool.Properties.Labels {
				labels[k] = v
			}
		}

		if err := d.Set("labels", labels); err != nil {
			return fmt.Errorf("error while setting the labels property for k8sNodepool %s: %w", d.Id(), err)

		}

		annotations := make(map[string]any)
		if nodePool.Properties.Annotations != nil && len(*nodePool.Properties.Annotations) > 0 {
			for k, v := range *nodePool.Properties.Annotations {
				annotations[k] = v
			}
		}

		if err := d.Set("annotations", annotations); err != nil {
			return fmt.Errorf("error while setting the annotations property for k8sNodepool %s: %w", d.Id(), err)
		}

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
	tflog.Info(ctx, "k8s node pool state", map[string]interface{}{"state": *resource.Metadata.State})
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

func getK8sNodePoolLans(lans []ionoscloud.KubernetesNodePoolLan) []any {

	var nodePoolLans []any
	for _, nodePoolLan := range lans {
		lanEntry := make(map[string]any)

		if nodePoolLan.Id != nil {
			lanEntry["id"] = *nodePoolLan.Id
		}

		if nodePoolLan.Dhcp != nil {
			lanEntry["dhcp"] = *nodePoolLan.Dhcp
		}

		if nodePoolLan.Routes != nil && len(*nodePoolLan.Routes) > 0 {
			var nodePoolRoutes []any
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
