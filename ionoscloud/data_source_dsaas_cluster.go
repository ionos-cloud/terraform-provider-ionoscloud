package ionoscloud

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	dsaas "github.com/ionos-cloud/sdk-go-autoscaling"
	dsaasService "github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/dsaas"
	"gopkg.in/yaml.v3"
	"log"
	"regexp"
	"strings"
)

func dataSourceDSaaSCluster() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDSaaSClusterRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:         schema.TypeString,
				Description:  "The id of your cluster.",
				Optional:     true,
				ValidateFunc: validation.All(validation.StringMatch(regexp.MustCompile("^[A-Za-z0-9][-A-Za-z0-9_.]*[A-Za-z0-9]$"), "")),
			},
			"name": {
				Type:        schema.TypeString,
				Description: "The name of your cluster.",
				Optional:    true,
			},
			"partial_match": {
				Type:        schema.TypeBool,
				Description: "Whether partial matching is allowed or not when using name argument.",
				Default:     false,
				Optional:    true,
			},
			"data_platform_version": {
				Type:        schema.TypeString,
				Description: "The version of the DataPlatform.",
				Computed:    true,
			},
			"datacenter_id": {
				Type:        schema.TypeString,
				Description: "The UUID of the virtual data center (VDC) the cluster is provisioned.",
				Computed:    true,
			},
			"maintenance_window": {
				Type:        schema.TypeList,
				Description: "Starting time of a weekly 4 hour-long window, during which maintenance might occur in hh:mm:ss format",
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"time": {
							Type:        schema.TypeString,
							Description: "Time at which the maintenance should start. Must conform to the 'HH:MM:SS' 24-hour format.",
							Computed:    true,
						},
						"day_of_the_week": {
							Type:     schema.TypeString,
							Computed: true,
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
			"kube_config": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
		Timeouts: &resourceDefaultTimeouts,
	}
}

func dataSourceDSaaSClusterRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(SdkBundle).DSaaSClient

	idValue, idOk := d.GetOk("id")
	nameValue, nameOk := d.GetOk("name")

	id := idValue.(string)
	name := nameValue.(string)

	if idOk && nameOk {
		diags := diag.FromErr(errors.New("id and display_name cannot be both specified in the same time"))
		return diags
	}
	if !idOk && !nameOk {
		diags := diag.FromErr(errors.New("please provide either the DSaaS Cluster id or display_name"))
		return diags
	}

	var cluster dsaas.ClusterResponseData
	var err error

	if idOk {
		/* search by ID */
		cluster, _, err = client.GetCluster(ctx, id)
		if err != nil {
			diags := diag.FromErr(fmt.Errorf("an error occurred while fetching the DSaaS Cluster with ID %s: %s", id, err))
			return diags
		}
	} else {
		var results []dsaas.ClusterResponseData

		partialMatch := d.Get("partial_match").(bool)

		log.Printf("[INFO] Using data source for DSaaS Cluster by name with partial_match %t and name: %s", partialMatch, name)

		if partialMatch {
			clusters, _, err := client.ListClusters(ctx, name)
			if err != nil {
				diags := diag.FromErr(fmt.Errorf("an error occurred while fetching DSaaS Clusters: %s", err.Error()))
				return diags
			}
			results = *clusters.Items
		} else {
			clusters, _, err := client.ListClusters(ctx, "")
			if err != nil {
				diags := diag.FromErr(fmt.Errorf("an error occurred while fetching DSaaS Clusters: %s", err.Error()))
				return diags
			}
			if clusters.Items != nil && len(*clusters.Items) > 0 {
				for _, clusterItem := range *clusters.Items {
					if clusterItem.Properties != nil && clusterItem.Properties.Name != nil && strings.EqualFold(*clusterItem.Properties.Name, name) {
						tmpCluster, _, err := client.GetCluster(ctx, *clusterItem.Id)
						if err != nil {
							return diag.FromErr(fmt.Errorf("an error occurred while fetching the DSaaS Cluster with ID: %s while searching by full name: %s: %w", *clusterItem.Id, name, err))
						}
						results = append(results, tmpCluster)
					}
				}
			}
		}

		if results == nil || len(results) == 0 {
			return diag.FromErr(fmt.Errorf("no DSaaS Cluster found with the specified name = %s", name))
		} else if len(results) > 1 {
			return diag.FromErr(fmt.Errorf("more than one DSaaS Cluster found with the specified criteria name = %s", name))
		} else {
			cluster = results[0]
		}

	}

	if err := dsaasService.SetDSaaSClusterData(d, cluster); err != nil {
		return diag.FromErr(err)
	}

	if err = setAdditionalDSaaSClusterData(d, &cluster, client); err != nil {
		return diag.FromErr(err)
	}

	return nil

}

func setAdditionalDSaaSClusterData(d *schema.ResourceData, cluster *dsaas.ClusterResponseData, client *dsaasService.Client) error {

	ctx, cancel := context.WithTimeout(context.Background(), *resourceDefaultTimeouts.Default)

	if cancel != nil {
		defer cancel()
	}

	/* get and set the kubeconfig*/
	if cluster.Id != nil {
		kubeConfig, _, err := client.GetClusterKubeConfig(ctx, *cluster.Id)
		if err != nil {
			return fmt.Errorf("an error occurred while fetching the kubernetes config for cluster with ID %s: %s", *cluster.Id, err)
		}

		if err := d.Set("kube_config", kubeConfig); err != nil {
			return err
		}
		fmt.Printf("KubeConfig %+v", kubeConfig)

		if err := setDSaaSConfigData(d, kubeConfig); err != nil {
			return err
		}
	}

	return nil
}

func setDSaaSConfigData(d *schema.ResourceData, configStr string) error {

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