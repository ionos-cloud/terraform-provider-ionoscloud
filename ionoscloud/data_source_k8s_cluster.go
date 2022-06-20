package ionoscloud

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	ionoscloud "github.com/ionos-cloud/sdk-go/v6"
	"gopkg.in/yaml.v3"
)

type KubeConfig struct {
	ApiVersion string `yaml:"apiVersion"`
	Clusters   []struct {
		Name    string
		Cluster struct {
			CaData string `yaml:"certificate-authority-data"`
			Server string
		}
	}
	Contexts []struct {
		Name    string
		Context struct {
			Cluster string
			User    string
		}
	}
	CurrentContext string `yaml:"current-context"`
	Kind           string
	Users          []struct {
		Name string
		User struct {
			Token string
		}
	}
	// preferences - add it when its structure is clear
}

func dataSourceK8sCluster() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceK8sReadCluster,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"state": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"k8s_version": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"maintenance_window": {
				Type:        schema.TypeList,
				Description: "A maintenance window comprise of a day of the week and a time for maintenance to be allowed",
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"time": {
							Type:        schema.TypeString,
							Description: "A clock time in the day when maintenance is allowed",
							Computed:    true,
						},
						"day_of_the_week": {
							Type:        schema.TypeString,
							Description: "Day of the week when maintenance is allowed",
							Computed:    true,
						},
					},
				},
			},
			"config": {
				Type:      schema.TypeList,
				Computed:  true,
				Sensitive: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"api_version": {
							Type:      schema.TypeString,
							Computed:  true,
							Sensitive: true,
						},
						"current_context": {
							Type:      schema.TypeString,
							Computed:  true,
							Sensitive: true,
						},
						"kind": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"users": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:      schema.TypeString,
										Computed:  true,
										Sensitive: true,
									},
									"user": {
										Type:      schema.TypeMap,
										Computed:  true,
										Sensitive: true,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
								},
							},
						},
						"clusters": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:      schema.TypeString,
										Computed:  true,
										Sensitive: true,
									},
									"cluster": {
										Type:      schema.TypeMap,
										Computed:  true,
										Sensitive: true,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
								},
							},
						},
						"contexts": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:      schema.TypeString,
										Computed:  true,
										Sensitive: true,
									},
									"context": {
										Type:      schema.TypeMap,
										Computed:  true,
										Sensitive: true,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
								},
							},
						},
					},
				},
			},
			"user_tokens": {
				Type:      schema.TypeMap,
				Sensitive: true,
				Computed:  true,
				Elem: &schema.Schema{
					Type:      schema.TypeString,
					Sensitive: true,
				},
			},
			"ca_crt": {
				Type:      schema.TypeString,
				Sensitive: true,
				Computed:  true,
			},
			"server": {
				Type:      schema.TypeString,
				Sensitive: true,
				Computed:  true,
			},
			"available_upgrade_versions": {
				Type:        schema.TypeList,
				Description: "A list of available versions for upgrading the cluster",
				Computed:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"viable_node_pool_versions": {
				Type:        schema.TypeList,
				Description: "A list of versions that may be used for node pools under this cluster",
				Computed:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"node_pools": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"kube_config": {
				Type:     schema.TypeString,
				Computed: true,
			},
			//"public": {
			//	Type:        schema.TypeBool,
			//	Description: "The indicator if the cluster is public or private. Be aware that setting it to false is currently in beta phase.",
			//	Computed:    true,
			//},
			"api_subnet_allow_list": {
				Type: schema.TypeList,
				Description: "Access to the K8s API server is restricted to these CIDRs. Cluster-internal traffic is not " +
					"affected by this restriction. If no allowlist is specified, access is not restricted. If an IP " +
					"without subnet mask is provided, the default value will be used: 32 for IPv4 and 128 for IPv6.",
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"s3_buckets": {
				Type:        schema.TypeList,
				Description: "List of S3 bucket configured for K8s usage. For now it contains only an S3 bucket used to store K8s API audit logs.",
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Description: "Name of the S3 bucket",
							Required:    true,
						},
					},
				},
			},
		},
		Timeouts: &resourceDefaultTimeouts,
	}
}

func dataSourceK8sReadCluster(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(SdkBundle).CloudApiClient

	id, idOk := d.GetOk("id")
	name, nameOk := d.GetOk("name")

	if idOk && nameOk {
		return diag.FromErr(errors.New("id and name cannot be both specified in the same time"))
	}
	if !idOk && !nameOk {
		return diag.FromErr(errors.New("please provide either the k8s cluster id or name"))
	}
	var cluster ionoscloud.KubernetesCluster
	var err error
	var apiResponse *ionoscloud.APIResponse

	if idOk {
		/* search by ID */
		cluster, apiResponse, err = client.KubernetesApi.K8sFindByClusterId(ctx, id.(string)).Execute()
		logApiRequestTime(apiResponse)
		if err != nil {
			return diag.FromErr(fmt.Errorf("an error occurred while fetching the k8s cluster with ID %s: %s", id.(string), err))
		}
	} else {
		/* search by name */
		var clusters ionoscloud.KubernetesClusters

		clusters, apiResponse, err := client.KubernetesApi.K8sGet(ctx).Depth(1).Execute()
		logApiRequestTime(apiResponse)
		if err != nil {
			return diag.FromErr(fmt.Errorf("an error occurred while fetching k8s clusters: %w", err))
		}

		if clusters.Items != nil {
			var results []ionoscloud.KubernetesCluster

			for _, c := range *clusters.Items {
				if c.Properties != nil && c.Properties.Name != nil && *c.Properties.Name == name.(string) {
					tmpCluster, apiResponse, err := client.KubernetesApi.K8sFindByClusterId(ctx, *c.Id).Execute()
					logApiRequestTime(apiResponse)
					if err != nil {
						return diag.FromErr(fmt.Errorf("an error occurred while fetching k8s cluster with ID %s: %w", *c.Id, err))
					}
					results = append(results, tmpCluster)
					break
				}
			}

			if results == nil || len(results) == 0 {
				return diag.FromErr(fmt.Errorf("no cluster found with the specified name %s", name.(string)))
			} else if len(results) > 1 {
				return diag.FromErr(fmt.Errorf("more than one cluster found with the specified name %s", name.(string)))
			} else {
				cluster = results[0]
			}
		}

	}

	if err = setK8sClusterData(d, &cluster); err != nil {
		return diag.FromErr(err)
	}

	if err = setAdditionalK8sClusterData(d, &cluster, client); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func setK8sConfigData(d *schema.ResourceData, configStr string) error {

	var kubeConfig KubeConfig
	if err := yaml.Unmarshal([]byte(configStr), &kubeConfig); err != nil {
		return err
	}

	userTokens := map[string]string{}

	var server string
	var caCrt []byte

	configMap := make(map[string]interface{})

	configMap["api_version"] = kubeConfig.ApiVersion
	configMap["current_context"] = kubeConfig.CurrentContext
	configMap["kind"] = kubeConfig.Kind

	clustersList := make([]map[string]interface{}, len(kubeConfig.Clusters))
	for i, cluster := range kubeConfig.Clusters {

		/* decode ca */
		decodedCrt := make([]byte, base64.StdEncoding.DecodedLen(len(cluster.Cluster.CaData)))
		if _, err := base64.StdEncoding.Decode(decodedCrt, []byte(cluster.Cluster.CaData)); err != nil {
			return err
		}

		if len(caCrt) == 0 {
			caCrt = decodedCrt
		}

		clustersList[i] = map[string]interface{}{
			"name": cluster.Name,
			"cluster": map[string]string{
				"server":                     cluster.Cluster.Server,
				"certificate_authority_data": string(decodedCrt),
			},
		}
	}

	configMap["clusters"] = clustersList

	contextsList := make([]map[string]interface{}, len(kubeConfig.Contexts))
	for i, contextVal := range kubeConfig.Contexts {
		contextsList[i] = map[string]interface{}{
			"name": contextVal.Name,
			"context": map[string]string{
				"cluster": contextVal.Context.Cluster,
				"user":    contextVal.Context.User,
			},
		}
	}

	configMap["contexts"] = contextsList

	userList := make([]map[string]interface{}, len(kubeConfig.Users))
	for i, user := range kubeConfig.Users {
		userList[i] = map[string]interface{}{
			"name": user.Name,
			"user": map[string]interface{}{
				"token": user.User.Token,
			},
		}

		userTokens[user.Name] = user.User.Token
	}

	configMap["users"] = userList

	configList := []map[string]interface{}{configMap}

	if err := d.Set("config", configList); err != nil {
		return err
	}

	if err := d.Set("user_tokens", userTokens); err != nil {
		return err
	}

	if err := d.Set("server", server); err != nil {
		return err
	}

	if err := d.Set("ca_crt", string(caCrt)); err != nil {
		return err
	}

	return nil
}

func setAdditionalK8sClusterData(d *schema.ResourceData, cluster *ionoscloud.KubernetesCluster, client *ionoscloud.APIClient) error {

	if cluster.Metadata != nil {
		if cluster.Metadata.State != nil {
			if err := d.Set("state", *cluster.Metadata.State); err != nil {
				return err
			}
		}

	}

	ctx, cancel := context.WithTimeout(context.Background(), *resourceDefaultTimeouts.Default)

	if cancel != nil {
		defer cancel()
	}

	/* get and set the kubeconfig and nodepools*/
	if cluster.Id != nil {
		kubeConfig, apiResponse, err := client.KubernetesApi.K8sKubeconfigGet(ctx, *cluster.Id).Execute()
		logApiRequestTime(apiResponse)
		if err != nil {
			return fmt.Errorf("an error occurred while fetching the kubernetes config for cluster with ID %s: %s", *cluster.Id, err)
		}

		if err := d.Set("kube_config", kubeConfig); err != nil {
			return err
		}

		if err := setK8sConfigData(d, kubeConfig); err != nil {
			return err
		}

		/* getting node pools */
		clusterNodePools, apiResponse, err := client.KubernetesApi.K8sNodepoolsGet(ctx, *cluster.Id).Execute()
		logApiRequestTime(apiResponse)
		if err != nil {
			return fmt.Errorf("an error occurred while fetching the kubernetes cluster node pools for cluster with ID %s: %s", *cluster.Id, err)
		}

		if clusterNodePools.Items != nil && len(*clusterNodePools.Items) > 0 {
			var nodePools []interface{}
			for _, nodePool := range *clusterNodePools.Items {
				nodePools = append(nodePools, *nodePool.Id)
			}
			if err := d.Set("node_pools", nodePools); err != nil {
				return err
			}
		}

		if cluster.Properties != nil && cluster.Properties.AvailableUpgradeVersions != nil {
			availableUpgradeVersions := make([]interface{}, len(*cluster.Properties.AvailableUpgradeVersions), len(*cluster.Properties.AvailableUpgradeVersions))
			for i, availableUpgradeVersion := range *cluster.Properties.AvailableUpgradeVersions {
				availableUpgradeVersions[i] = availableUpgradeVersion
			}
			if err := d.Set("available_upgrade_versions", availableUpgradeVersions); err != nil {
				return err
			}
		}
	}

	return nil
}
