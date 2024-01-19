package ionoscloud

import (
	"context"
	"encoding/base64"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	ionoscloud "github.com/ionos-cloud/sdk-go/v6"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/internal/uuidgen"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils"
	"gopkg.in/yaml.v3"
)

func dataSourceK8sClusters() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceK8sReadClusterList,
		Schema: map[string]*schema.Schema{
			"clusters": {
				Type:        schema.TypeList,
				Description: "List of clusters which match the filtering criteria.",
				Computed:    true,
				Elem:        &schema.Resource{Schema: dataSourceK8sClusterSchema},
			},
			"entries": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"location": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"k8s_version": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"public": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"states": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
		Timeouts: &resourceDefaultTimeouts,
	}
}

func dataSourceK8sReadClusterList(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(services.SdkBundle).CloudApiClient
	filters := clusterListFilter{states: map[string]struct{}{}}

	if name, nameOk := d.GetOk("name"); nameOk {
		filters.name = name.(string)
	}
	if location, locationOk := d.GetOk("location"); locationOk {
		filters.location = location.(string)
	}
	if public, publicOk := d.GetOkExists("public"); publicOk {
		filters.public.set = true
		filters.public.value = public.(bool)
	}
	if states, statesOk := d.GetOkExists("states"); statesOk {
		statesVal := states.([]interface{})
		for _, state := range statesVal {
			filters.states[state.(string)] = struct{}{}
		}
	}

	clusters, apiResponse, err := client.KubernetesApi.K8sGet(ctx).Depth(1).Execute()
	logApiRequestTime(apiResponse)
	if err != nil {
		return diag.FromErr(fmt.Errorf("an error occurred while fetching k8s clusters: %w", err))
	}

	if clusters.Items != nil {

		var results []*ionoscloud.KubernetesCluster
		for _, c := range *clusters.Items {
			if filters.filter(c) && true {
				c := c
				results = append(results, &c)
			}
		}
		if len(results) == 0 {
			return diag.FromErr(fmt.Errorf("no clusters match the specified filtering criteria"))
		}
		if err := setDataSourceK8sSetClusterList(d, results, client); err != nil {
			return diag.FromErr(err)
		}
	}
	return nil
}

func setDataSourceK8sSetClusterList(d *schema.ResourceData, clusters []*ionoscloud.KubernetesCluster, client *ionoscloud.APIClient) error {

	if d.Id() == "" {
		d.SetId(uuidgen.ResourceUuid().String())
	}
	clusterList := make([]map[string]any, 0)
	for _, c := range clusters {

		clusterProperties, err := K8sClusterProperties(c, client)
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

func K8sClusterProperties(cluster *ionoscloud.KubernetesCluster, client *ionoscloud.APIClient) (map[string]any, error) {
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
	if cluster.Metadata != nil {
		utils.SetPropWithNilCheck(clusterProperties, "state", cluster.Metadata.State)
	}
	if cluster.Id != nil {
		ctx, cancel := context.WithTimeout(context.Background(), *resourceDefaultTimeouts.Default)
		defer cancel()

		// kubeconfig
		clusterConfig, apiResponse, err := client.KubernetesApi.K8sKubeconfigGet(ctx, *cluster.Id).Execute()
		logApiRequestTime(apiResponse)
		if err != nil {
			return nil, fmt.Errorf("an error occurred while fetching the kubernetes config for cluster with ID %s: %w", *cluster.Id, err)
		}
		clusterProperties["kube_config"] = clusterConfig
		clusterConfigProperties, err := K8sClusterConfigProperties(clusterConfig)
		if err != nil {
			return nil, err
		}

		// node pools

		clusterProperties = mergeMaps(clusterProperties, clusterConfigProperties)
	}

	return clusterProperties, nil
}

func K8sClusterConfigProperties(clusterConfig string) (map[string]any, error) {
	kubeConfig := KubeConfig{}
	if err := yaml.Unmarshal([]byte(clusterConfig), &kubeConfig); err != nil {
		return nil, err
	}

	clusterProperties := make(map[string]any)
	clusterConfigProperties := make(map[string]any)

	clusterConfigProperties["api_version"] = kubeConfig.ApiVersion
	clusterConfigProperties["current_context"] = kubeConfig.CurrentContext
	clusterConfigProperties["kind"] = kubeConfig.Kind

	var decodedCert string
	clusters := make([]map[string]any, len(kubeConfig.Clusters))
	for i, c := range kubeConfig.Clusters {
		caCert := make([]byte, base64.StdEncoding.DecodedLen(len(c.Cluster.CaData)))
		if _, err := base64.StdEncoding.Decode(caCert, []byte(c.Cluster.CaData)); err != nil {
			return nil, err
		}
		decodedCert = string(caCert)
		clusters[i] = map[string]any{"name": c.Name, "cluster": map[string]string{"server": c.Cluster.Server, "certificate_authority_data": decodedCert}}
	}
	clusterConfigProperties["clusters"] = clusters

	contexts := make([]map[string]any, len(kubeConfig.Contexts))
	for i, contextVal := range kubeConfig.Contexts {
		contexts[i] = map[string]interface{}{
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
	clusterProperties["ca_crt"] = decodedCert
	clusterProperties["user_tokens"] = userTokens

	return clusterProperties, nil
}

func mergeMaps(maps ...map[string]any) map[string]any {
	merged := map[string]any{}
	for _, m := range maps {

		for k := range m {
			merged[k] = m[k]
		}
	}
	return merged
}

type clusterListFilter struct {
	name       string
	location   string
	k8sVersion string
	states     map[string]struct{}
	public     struct {
		set   bool
		value bool
	}
}

func (clf clusterListFilter) filter(cluster ionoscloud.KubernetesCluster) bool {

	type filterFn func(ionoscloud.KubernetesCluster) bool
	filterFns := []filterFn{
		clf.nameFieldFilter,
		clf.locationFieldFilter,
		clf.publicFieldFilter,
		clf.statesFieldFilter,
		clf.k8sVersionFieldFilter,
	}

	if cluster.Properties == nil {
		return false
	}

	for _, f := range filterFns {
		if !f(cluster) {
			return false
		}
	}

	return true
}
func (clf clusterListFilter) nameFieldFilter(cluster ionoscloud.KubernetesCluster) bool {
	if clf.name == "" || (cluster.Properties.Name != nil && strings.Contains(*cluster.Properties.Name, clf.name)) {
		return true
	}
	return false
}
func (clf clusterListFilter) locationFieldFilter(cluster ionoscloud.KubernetesCluster) bool {
	if clf.location == "" || (cluster.Properties.Location != nil && strings.EqualFold(*cluster.Properties.Location, clf.location)) {
		return true
	}
	return false
}
func (clf clusterListFilter) publicFieldFilter(cluster ionoscloud.KubernetesCluster) bool {
	if !clf.public.set || (cluster.Properties.Public != nil && *cluster.Properties.Public == clf.public.value) {
		return true
	}
	return false
}
func (clf clusterListFilter) statesFieldFilter(cluster ionoscloud.KubernetesCluster) bool {
	if len(clf.states) == 0 {
		return true
	}
	if cluster.Metadata != nil && cluster.Metadata.State != nil {
		if _, ok := clf.states[*cluster.Metadata.State]; ok {
			return true
		}
	}
	return false
}
func (clf clusterListFilter) k8sVersionFieldFilter(cluster ionoscloud.KubernetesCluster) bool {
	if clf.k8sVersion == "" || (cluster.Properties.K8sVersion != nil && strings.EqualFold(*cluster.Properties.K8sVersion, clf.k8sVersion)) {
		return true
	}
	return false
}
