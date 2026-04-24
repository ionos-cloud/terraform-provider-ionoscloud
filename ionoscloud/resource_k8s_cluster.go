package ionoscloud

import (
	"context"
	"fmt"
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
				Type:             schema.TypeString,
				Description:      "The desired name for the cluster",
				Required:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringIsNotWhiteSpace),
			},
			"k8s_version": {
				Type:             schema.TypeString,
				Description:      "The desired Kubernetes Version. For supported values, please check the API documentation. Downgrades are not supported. The provider will ignore downgrades of patch level.",
				Optional:         true,
				Computed:         true,
				DiffSuppressFunc: DiffBasedOnVersion,
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
			"viable_node_pool_versions": {
				Type:        schema.TypeList,
				Description: "List of versions that may be used for node pools under this cluster",
				Computed:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"public": {
				Type:        schema.TypeBool,
				Description: "The indicator if the cluster is public or private.",
				Optional:    true,
				Default:     true,
				ForceNew:    true,
			},
			"nat_gateway_ip": {
				Type:             schema.TypeString,
				Description:      "The NAT gateway IP of the cluster if the cluster is private. This attribute is immutable. Must be a reserved IP in the same location as the cluster's location. This attribute is mandatory if the cluster is private.",
				ValidateDiagFunc: validation.ToDiagFunc(validation.All(validation.IsIPv4Address, validation.IsIPv6Address)),
				Optional:         true,
				ForceNew:         true,
			},
			"node_subnet": {
				Type:             schema.TypeString,
				Description:      "The node subnet of the cluster, if the cluster is private. This attribute is optional and immutable. Must be a valid CIDR notation for an IPv4 network prefix of 16 bits length.",
				ValidateDiagFunc: validation.ToDiagFunc(validation.IsCIDR),
				Optional:         true,
				Computed:         true,
				ForceNew:         true,
			},
			"location": {
				Type:        schema.TypeString,
				Description: "This attribute is mandatory if the cluster is private. The location must be enabled for your contract, or you must have a data center at that location. This attribute is immutable.",
				Optional:    true,
				ForceNew:    true,
			},
			"allow_replace": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "When set to true, allows the update of immutable fields by destroying and re-creating the cluster.",
			},
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
				Description: "List of Object Storage bucket configured for K8s usage. For now it contains only an Object Storage bucket used to store K8s API audit logs.",
				Optional:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Description: "Name of the Object Storage bucket",
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

	allowReplace := diff.Get("allow_replace").(bool)
	if allowReplace {
		return nil
	}
	// we do not want to check in case of resource creation
	if diff.Id() == "" {
		return nil
	}
	if diff.HasChange("public") {
		return fmt.Errorf("public %s", ImmutableError)
	}
	if diff.HasChange("location") {
		return fmt.Errorf("location %s", ImmutableError)
	}
	if diff.HasChange("nat_gateway_ip") {
		return fmt.Errorf("nat_gateway_ip %s", ImmutableError)
	}
	if diff.HasChange("node_subnet") {
		return fmt.Errorf("node_subnet %s", ImmutableError)
	}
	return nil

}
func resourcek8sClusterCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	location := d.Get("location").(string)
	client, err := meta.(bundleclient.SdkBundle).NewCloudAPIClient(location)
	if err != nil {
		return diag.FromErr(err)
	}

	clusterName := d.Get("name").(string)
	cluster := ionoscloud.KubernetesClusterForPost{
		Properties: &ionoscloud.KubernetesClusterPropertiesForPost{
			Name: &clusterName,
		},
	}

	if k8svVal, k8svOk := d.GetOk("k8s_version"); k8svOk {
		tflog.Info(ctx, "setting k8s version", map[string]interface{}{"version": k8svVal.(string)})
		k8svVal := k8svVal.(string)
		cluster.Properties.K8sVersion = &k8svVal
	}

	if _, mwOk := d.GetOk("maintenance_window.0"); mwOk {
		cluster.Properties.MaintenanceWindow = &ionoscloud.KubernetesMaintenanceWindow{}
	}

	if mtVal, mtOk := d.GetOk("maintenance_window.0.time"); mtOk {
		tflog.Info(ctx, "setting maintenance window time", map[string]interface{}{"time": mtVal.(string)})
		mtVal := mtVal.(string)
		cluster.Properties.MaintenanceWindow.Time = &mtVal
	}

	if mdVal, mdOk := d.GetOk("maintenance_window.0.day_of_the_week"); mdOk {
		mdVal := mdVal.(string)
		cluster.Properties.MaintenanceWindow.DayOfTheWeek = &mdVal
	}

	if publicVal, publicOk := d.GetOkExists("public"); publicOk {
		publicVal := publicVal.(bool)
		cluster.Properties.Public = &publicVal
	}

	if locationVal, locationOk := d.GetOk("location"); locationOk {
		locationVal := locationVal.(string)
		cluster.Properties.Location = &locationVal
	}

	if natGatewayIpVal, natGatewayIpOk := d.GetOk("nat_gateway_ip"); natGatewayIpOk {
		natGatewayIpVal := natGatewayIpVal.(string)
		cluster.Properties.NatGatewayIp = &natGatewayIpVal
	}

	if nodeSubnetVal, nodeSubnetOk := d.GetOk("node_subnet"); nodeSubnetOk {
		nodeSubnetVal := nodeSubnetVal.(string)
		cluster.Properties.NodeSubnet = &nodeSubnetVal
	}

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
					return diagutil.ToDiags(d, fmt.Errorf("name must be provided for Object Storage bucket"), nil)
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
		return diagutil.ToDiags(d, fmt.Errorf("error creating k8s cluster: %w", err), &diagutil.ErrorContext{StatusCode: apiResponse.SafeStatusCode()})
	}

	d.SetId(*createdCluster.Id)
	tflog.Info(ctx, "created k8s cluster", map[string]interface{}{"cluster_id": d.Id()})

	for {
		tflog.Info(ctx, "waiting for k8s cluster to be ACTIVE", map[string]interface{}{"cluster_id": d.Id()})

		clusterReady, rsErr := k8sClusterReady(ctx, client, d)

		if rsErr != nil {
			return diagutil.ToDiags(d, fmt.Errorf("error while checking readiness status of k8s cluster: %w", rsErr), nil)
		}

		if clusterReady {
			tflog.Info(ctx, "k8s cluster ready", map[string]interface{}{"cluster_id": d.Id()})
			break
		}

		select {
		case <-time.After(constant.SleepInterval):
			tflog.Info(ctx, "k8s cluster not ready, retrying")
		case <-ctx.Done():
			tflog.Info(ctx, "k8s cluster creation timed out")
			return diagutil.ToDiags(d, fmt.Errorf("k8s cluster creation timed out! WARNING: your k8s cluster will still probably be created after some time but the terraform state wont reflect that; check your Ionos Cloud account for updates"), nil)
		}

	}

	return resourcek8sClusterRead(ctx, d, meta)
}

func resourcek8sClusterRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	location := d.Get("location").(string)
	client, err := meta.(bundleclient.SdkBundle).NewCloudAPIClient(location)
	if err != nil {
		return diag.FromErr(err)
	}

	cluster, apiResponse, err := client.KubernetesApi.K8sFindByClusterId(ctx, d.Id()).Execute()
	logApiRequestTime(apiResponse)

	if err != nil {
		if httpNotFound(apiResponse) {
			d.SetId("")
			return nil
		}
		return diagutil.ToDiags(d, fmt.Errorf("error while fetching k8s cluster: %w", err), &diagutil.ErrorContext{StatusCode: apiResponse.SafeStatusCode()})
	}

	tflog.Info(ctx, "retrieved k8s cluster", map[string]interface{}{"cluster_id": d.Id()})

	if err := setK8sClusterData(d, &cluster); err != nil {
		return diagutil.ToDiags(d, err, nil)
	}

	return nil
}

func resourcek8sClusterUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	location := d.Get("location").(string)
	client, err := meta.(bundleclient.SdkBundle).NewCloudAPIClient(location)
	if err != nil {
		return diag.FromErr(err)
	}

	request := ionoscloud.KubernetesClusterForPut{}

	clusterName := d.Get("name").(string)
	request.Properties = &ionoscloud.KubernetesClusterPropertiesForPut{
		Name: &clusterName,
	}

	if d.HasChange("name") {
		oldName, newName := d.GetChange("name")
		tflog.Info(ctx, "k8s cluster name changed", map[string]interface{}{"old": oldName, "new": newName})
		newNameStr := newName.(string)
		request.Properties.Name = &newNameStr
	}

	tflog.Info(ctx, "attempting k8s cluster update", map[string]interface{}{"cluster_id": d.Id()})

	if d.HasChange("k8s_version") {
		oldk8sVersion, newk8sVersion := d.GetChange("k8s_version")
		tflog.Info(ctx, "k8s version changed", map[string]interface{}{"old": oldk8sVersion, "new": newk8sVersion})
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
					tflog.Info(ctx, "k8s maintenance window DOW changed", map[string]interface{}{"old": oldMd, "new": newMd})
					updateMaintenanceWindow = true
					newMd := newMd.(string)
					maintenanceWindow.DayOfTheWeek = &newMd
				}
			}

			if d.HasChange("maintenance_window.0.time") {

				oldMt, newMt := d.GetChange("maintenance_window.0.time")
				if newMt.(string) != "" {
					tflog.Info(ctx, "k8s maintenance window time changed", map[string]interface{}{"old": oldMt, "new": newMt})
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
		if httpNotFound(apiResponse) {
			d.SetId("")
			return nil
		}
		return diagutil.ToDiags(d, fmt.Errorf("error while updating k8s cluster: %w", err), &diagutil.ErrorContext{StatusCode: apiResponse.SafeStatusCode()})
	}

	for {
		tflog.Info(ctx, "waiting for k8s cluster to be ready", map[string]interface{}{"cluster_id": d.Id()})

		clusterReady, rsErr := k8sClusterReady(ctx, client, d)

		if rsErr != nil {
			return diagutil.ToDiags(d, fmt.Errorf("error while checking readiness status of k8s cluster: %w", rsErr), nil)
		}

		if clusterReady {
			tflog.Info(ctx, "k8s cluster ready", map[string]interface{}{"cluster_id": d.Id()})
			break
		}

		select {
		case <-time.After(constant.SleepInterval):
			tflog.Info(ctx, "k8s cluster not ready, retrying")
		case <-ctx.Done():
			return diagutil.ToDiags(d, fmt.Errorf("k8s cluster update timed out! WARNING: your k8s cluster will still probably be created after some time but the terraform state won't reflect that; check your Ionos Cloud account for updates"), nil)
		}

	}

	return resourcek8sClusterRead(ctx, d, meta)
}

func resourcek8sClusterDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	location := d.Get("location").(string)
	client, err := meta.(bundleclient.SdkBundle).NewCloudAPIClient(location)
	if err != nil {
		return diag.FromErr(err)
	}

	apiResponse, err := client.KubernetesApi.K8sDelete(ctx, d.Id()).Execute()
	logApiRequestTime(apiResponse)

	if err != nil {
		if httpNotFound(apiResponse) {
			d.SetId("")
			return nil
		}
		return diagutil.ToDiags(d, fmt.Errorf("error while deleting k8s cluster: %w", err), &diagutil.ErrorContext{StatusCode: apiResponse.SafeStatusCode()})
	}

	for {
		tflog.Info(ctx, "waiting for k8s cluster to be deleted", map[string]interface{}{"cluster_id": d.Id()})

		clusterDeleted, dsErr := k8sClusterDeleted(ctx, client, d)

		if dsErr != nil {
			return diagutil.ToDiags(d, fmt.Errorf("error while checking deletion status of k8s cluster: %w", dsErr), nil)
		}

		if clusterDeleted {
			tflog.Info(ctx, "successfully deleted k8s cluster", map[string]interface{}{"cluster_id": d.Id()})
			break
		}

		select {
		case <-time.After(constant.SleepInterval):
			tflog.Info(ctx, "k8s cluster not yet deleted, retrying")
		case <-ctx.Done():
			return diagutil.ToDiags(d, fmt.Errorf("k8s cluster deletion timed out! WARNING: your k8s cluster will still probably be deleted after some time but the terraform state won't reflect that; check your Ionos Cloud account for updates"), nil)
		}
	}

	d.SetId("")

	return nil
}

func resourceK8sClusterImport(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	importID := d.Id()

	location, parts := splitImportID(importID, ":")
	if len(parts) != 1 {
		return nil, fmt.Errorf("invalid import identifier: expected one of <location>:<cluster-id> or <cluster-id>, got: %s", importID)
	}

	if err := validateImportIDParts(parts); err != nil {
		return nil, fmt.Errorf("failed validating import identifier %q: %w", importID, err)
	}

	clusterID := parts[0]

	client, err := meta.(bundleclient.SdkBundle).NewCloudAPIClient(location)
	if err != nil {
		return nil, err
	}

	cluster, apiResponse, err := client.KubernetesApi.K8sFindByClusterId(ctx, clusterID).Execute()
	logApiRequestTime(apiResponse)

	if err != nil {
		if httpNotFound(apiResponse) {
			d.SetId("")
			return nil, diagutil.ToError(d, fmt.Errorf("unable to find k8s cluster %q", clusterID), &diagutil.ErrorContext{StatusCode: apiResponse.SafeStatusCode()})
		}
		return nil, diagutil.ToError(d, fmt.Errorf("unable to retrieve k8s cluster %q, error:%w", clusterID, err), &diagutil.ErrorContext{StatusCode: apiResponse.SafeStatusCode()})
	}

	tflog.Info(ctx, "k8s cluster imported", map[string]interface{}{"cluster_id": clusterID})

	if err := setK8sClusterData(d, &cluster); err != nil {
		return nil, diagutil.ToError(d, err, nil)
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

		if cluster.Properties.Public != nil {
			if err := d.Set("public", *cluster.Properties.Public); err != nil {
				return utils.GenerateSetError(constant.K8sClusterResource, "public", err)
			}
		}

		if cluster.Properties.Location != nil {
			if err := d.Set("location", *cluster.Properties.Location); err != nil {
				return utils.GenerateSetError(constant.K8sClusterResource, "location", err)
			}
		}

		if cluster.Properties.NatGatewayIp != nil {
			if err := d.Set("nat_gateway_ip", *cluster.Properties.NatGatewayIp); err != nil {
				return utils.GenerateSetError(constant.K8sClusterResource, "nat_gateway_ip", err)
			}
		}

		if cluster.Properties.NodeSubnet != nil {
			if err := d.Set("node_subnet", *cluster.Properties.NodeSubnet); err != nil {
				return utils.GenerateSetError(constant.K8sClusterResource, "node_subnet", err)
			}
		}

		if cluster.Properties.ApiSubnetAllowList != nil {
			apiSubnetAllowLists := make([]interface{}, len(*cluster.Properties.ApiSubnetAllowList), len(*cluster.Properties.ApiSubnetAllowList))
			for i, apiSubnetAllowList := range *cluster.Properties.ApiSubnetAllowList {
				apiSubnetAllowLists[i] = apiSubnetAllowList
			}
			if err := d.Set("api_subnet_allow_list", apiSubnetAllowLists); err != nil {
				return fmt.Errorf("error while setting api_subnet_allow_list property for cluster with ID: %s, error: %w", d.Id(), err)
			}
		} else {
			var emptySlice []interface{}
			if err := d.Set("api_subnet_allow_list", emptySlice); err != nil {
				return fmt.Errorf("error while setting api_subnet_allow_list property for cluster with ID: %s, error: %w", d.Id(), err)
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
				return fmt.Errorf("error while setting s3_buckets property for cluser %s: %w", d.Id(), err)
			}
		}

	}

	return nil
}

func k8sClusterReady(ctx context.Context, client *ionoscloud.APIClient, d *schema.ResourceData) (bool, error) {
	resource, apiResponse, err := client.KubernetesApi.K8sFindByClusterId(ctx, d.Id()).Execute()
	logApiRequestTime(apiResponse)
	if err != nil {
		return true, fmt.Errorf("error checking k8s cluster status: %w", err)
	}
	if resource.Metadata == nil || resource.Metadata.State == nil {
		return false, fmt.Errorf("error while checking k8s cluster status: state is nil")
	}
	if utils.IsStateFailed(*resource.Metadata.State) {
		return false, fmt.Errorf("error while checking if k8s cluster is ready %s, state %s", *resource.Id, *resource.Metadata.State)
	}
	tflog.Info(ctx, "k8s cluster state", map[string]interface{}{"state": *resource.Metadata.State})
	// k8s is the only resource that has a state of ACTIVE when it is ready
	return strings.EqualFold(*resource.Metadata.State, ionoscloud.Active), nil
}

func k8sClusterDeleted(ctx context.Context, client *ionoscloud.APIClient, d *schema.ResourceData) (bool, error) {

	cluster, apiResponse, err := client.KubernetesApi.K8sFindByClusterId(ctx, d.Id()).Execute()
	logApiRequestTime(apiResponse)

	if err != nil {
		if httpNotFound(apiResponse) {
			return true, nil
		}
		return true, fmt.Errorf("error checking k8s cluster deletion status: %w", err)
	}
	if cluster.Metadata != nil && cluster.Metadata.State != nil {
		if utils.IsStateFailed(*cluster.Metadata.State) {
			return false, fmt.Errorf("error while checking if k8s cluster is deleted properly, cluster ID: %s, state: %s", *cluster.Id, *cluster.Metadata.State)
		}
	}

	return false, nil
}
