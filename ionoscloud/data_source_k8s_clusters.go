package ionoscloud

import (
	"context"
	"encoding/base64"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/iancoleman/strcase"

	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/internal/uuidgen"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils"

	ionoscloud "github.com/ionos-cloud/sdk-go/v6"
)

func dataSourceK8sClusters() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceK8sReadClusters,
		Schema: map[string]*schema.Schema{
			"clusters": {
				Type:        schema.TypeList,
				Description: "List of clusters which match the filtering criteria.",
				Computed:    true,
				Elem:        &schema.Resource{Schema: dataSourceK8sClusterSchema()},
			},
			"entries": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"filter": dataSourceFiltersSchema(),
		},
		Timeouts: &resourceDefaultTimeouts,
	}
}

func dataSourceK8sReadClusters(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {

	// strcase.ToLowerCamel doesn't produce the correct format for these keys, provide them directly (https://github.com/iancoleman/strcase/issues/51)
	filterKeys := map[string]string{
		"k8s_version": "k8sVersion",
	}

	client := meta.(services.SdkBundle).CloudApiClient
	req := client.KubernetesApi.K8sGet(ctx).Depth(1)

	filters, filtersOk := d.GetOk("filter")
	if filtersOk {
		for _, v := range filters.(*schema.Set).List() {
			filter := v.(map[string]any)
			key := filter["name"].(string)
			value := filter["value"].(string)
			key = strcase.ToLowerCamel(key)
			if v, ok := filterKeys[key]; ok {
				key = v
			}
			req.Filter(key, value)
		}
	}

	clusters, apiResponse, err := req.Execute()
	logApiRequestTime(apiResponse)
	if err != nil {
		return diag.FromErr(fmt.Errorf("an error occurred while fetching k8s clusters: %w", err))
	}
	if clusters.Items != nil && len(*clusters.Items) == 0 {
		return diag.FromErr(fmt.Errorf("no clusters found"))
	}
	if err := setDataSourceK8sSetClusters(ctx, d, *clusters.Items, client); err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func setDataSourceK8sSetClusters(ctx context.Context, d *schema.ResourceData, clusters []ionoscloud.KubernetesCluster, client *ionoscloud.APIClient) error {

	if d.Id() == "" {
		d.SetId(uuidgen.ResourceUuid().String())
	}
	clusterList := make([]map[string]any, 0, len(clusters))
	for _, c := range clusters {
		clusterProperties, err := K8sClusterProperties(ctx, c, client)
		if err != nil {
			return err
		}
		clusterList = append(clusterList, clusterProperties)
	}

	if err := d.Set("clusters", clusterList); err != nil {
		return err
	}
	if err := d.Set("entries", len(clusterList)); err != nil {
		return err
	}

	return nil
}

// K8sClusterProperties returns a map equivalent of dataSourceK8sClusterSchema
func K8sClusterProperties(ctx context.Context, cluster ionoscloud.KubernetesCluster, client *ionoscloud.APIClient) (map[string]any, error) {
	if cluster.Properties == nil {
		clusterID := "nil"
		if cluster.Id != nil {
			clusterID = *cluster.Id
		}
		return nil, fmt.Errorf("cannot set data, Properties was nil for cluster: %s", clusterID)
	}
	clusterProperties := make(map[string]any)

	utils.SetPropWithNilCheck(clusterProperties, "name", cluster.Properties.Name)
	utils.SetPropWithNilCheck(clusterProperties, "k8s_version", cluster.Properties.K8sVersion)
	utils.SetPropWithNilCheck(clusterProperties, "public", cluster.Properties.Public)
	utils.SetPropWithNilCheck(clusterProperties, "location", cluster.Properties.Location)
	utils.SetPropWithNilCheck(clusterProperties, "nat_gateway_ip", cluster.Properties.NatGatewayIp)
	utils.SetPropWithNilCheck(clusterProperties, "node_subnet", cluster.Properties.NodeSubnet)
	utils.SetPropWithNilCheck(clusterProperties, "viable_node_pool_versions", cluster.Properties.ViableNodePoolVersions)

	if cluster.Properties.MaintenanceWindow != nil && cluster.Properties.MaintenanceWindow.Time != nil && cluster.Properties.MaintenanceWindow.DayOfTheWeek != nil {
		clusterProperties["maintenance_window"] = []map[string]any{{"time": *cluster.Properties.MaintenanceWindow.Time, "day_of_the_week": *cluster.Properties.MaintenanceWindow.DayOfTheWeek}}
	}
	if cluster.Properties.S3Buckets != nil {
		s3Buckets := make([]map[string]any, len(*cluster.Properties.S3Buckets))
		for i, s3Bucket := range *cluster.Properties.S3Buckets {
			s3Buckets[i] = map[string]any{"name": s3Bucket.Name}
		}
		clusterProperties["s3_buckets"] = s3Buckets
	}
	if cluster.Properties.AvailableUpgradeVersions != nil {
		availableUpgradeVersions := make([]any, len(*cluster.Properties.AvailableUpgradeVersions))
		for i, availableUpgradeVersion := range *cluster.Properties.AvailableUpgradeVersions {
			availableUpgradeVersions[i] = availableUpgradeVersion
		}
		clusterProperties["available_upgrade_versions"] = availableUpgradeVersions
	}
	if cluster.Properties.ApiSubnetAllowList != nil {
		apiSubnetAllowList := make([]any, len(*cluster.Properties.ApiSubnetAllowList))
		for i, subnet := range *cluster.Properties.ApiSubnetAllowList {
			apiSubnetAllowList[i] = subnet
		}
		clusterProperties["api_subnet_allow_list"] = apiSubnetAllowList
	}
	if cluster.Metadata != nil {
		utils.SetPropWithNilCheck(clusterProperties, "state", cluster.Metadata.State)
	}
	if cluster.Id != nil {

		// kubeconfig
		clusterConfig, apiResponse, err := client.KubernetesApi.K8sKubeconfigGet(ctx, *cluster.Id).Execute()
		logApiRequestTime(apiResponse)
		if err != nil {
			return nil, fmt.Errorf("an error occurred while fetching the kubernetes config for cluster with ID %s: %w", *cluster.Id, err)
		}
		clusterProperties["kube_config"] = clusterConfig
		clusterConfigProperties, err := setKubeConfigData(clusterConfig)
		if err != nil {
			return nil, err
		}

		// node pools
		clusterNodePools, apiResponse, err := client.KubernetesApi.K8sNodepoolsGet(ctx, *cluster.Id).Execute()
		logApiRequestTime(apiResponse)
		if err != nil {
			return nil, fmt.Errorf("an error occurred while fetching the kubernetes cluster node pools for cluster with ID %s: %w", *cluster.Id, err)
		}
		if clusterNodePools.Items != nil && len(*clusterNodePools.Items) > 0 {
			var nodePools []any
			for _, nodePool := range *clusterNodePools.Items {
				nodePools = append(nodePools, *nodePool.Id)
			}
			clusterProperties["node_pools"] = nodePools
		}
		clusterProperties = utils.MergeMaps(clusterProperties, clusterConfigProperties)
	}

	return clusterProperties, nil
}

func setKubeConfigData(clusterConfig string) (map[string]any, error) {

	err, kubeConfig := parseClusterKubeconfig(clusterConfig)
	if err != nil {
		return nil, err
	}

	clusterProperties := make(map[string]any)
	clusterConfigProperties := make(map[string]any)

	clusterConfigProperties["api_version"] = kubeConfig.ApiVersion
	clusterConfigProperties["current_context"] = kubeConfig.CurrentContext
	clusterConfigProperties["kind"] = kubeConfig.Kind

	// Managed K8s clusters each have their own unique kubeconfig so there is only 1 Clusters entry
	if len(kubeConfig.Clusters) != 0 {
		caData := kubeConfig.Clusters[0].Cluster.CaData
		decodedCrt := make([]byte, base64.StdEncoding.DecodedLen(len(caData)))
		if _, err := base64.StdEncoding.Decode(decodedCrt, []byte(caData)); err != nil {
			return nil, err
		}
		clusterProperties["ca_crt"] = string(decodedCrt)
		clusterProperties["server"] = kubeConfig.Clusters[0].Cluster.Server
		clusterConfigProperties["clusters"] = []map[string]any{
			{"name": kubeConfig.Clusters[0].Name, "cluster": map[string]string{"server": kubeConfig.Clusters[0].Cluster.Server, "certificate_authority_data": string(decodedCrt)}},
		}
	}

	contexts := make([]map[string]any, len(kubeConfig.Contexts))
	for i, contextVal := range kubeConfig.Contexts {
		contexts[i] = map[string]any{
			"name": contextVal.Name,
			"context": map[string]string{
				"cluster": contextVal.Context.Cluster,
				"user":    contextVal.Context.User,
			},
		}
	}
	clusterConfigProperties["contexts"] = contexts

	users := make([]map[string]any, len(kubeConfig.Users))
	userTokens := make(map[string]string)
	for i, user := range kubeConfig.Users {
		users[i] = map[string]any{
			"name": user.Name,
			"user": map[string]any{
				"token": user.User.Token,
			},
		}
		userTokens[user.Name] = user.User.Token
	}
	clusterConfigProperties["users"] = users

	clusterProperties["config"] = []map[string]any{clusterConfigProperties}
	clusterProperties["user_tokens"] = userTokens

	return clusterProperties, nil
}
