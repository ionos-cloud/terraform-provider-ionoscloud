package ionoscloud

import (
	"context"
	"errors"
	"fmt"
	ionoscloud "github.com/ionos-cloud/sdk-go/v5"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceK8sCluster() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceK8sReadCluster,
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
			"public": {
				Type: schema.TypeBool,
				Description: "The indicator if the cluster is public or private. Be aware that setting it to false is " +
					"currently in beta phase.",
				Optional: true,
				Computed: true,
			},
			"gateway_ip": {
				Type: schema.TypeString,
				Description: "The IP address of the gateway used by the cluster. This is mandatory when `public` is set " +
					"to `false` and should not be provided otherwise.",
				Optional: true,
			},
		},
		Timeouts: &resourceDefaultTimeouts,
	}
}

func dataSourceK8sReadCluster(d *schema.ResourceData, meta interface{}) error {
	client := meta.(SdkBundle).Client

	id, idOk := d.GetOk("id")
	name, nameOk := d.GetOk("name")

	if idOk && nameOk {
		return errors.New("id and name cannot be both specified in the same time")
	}
	if !idOk && !nameOk {
		return errors.New("please provide either the lan id or name")
	}
	var cluster ionoscloud.KubernetesCluster
	var err error

	ctx, cancel := context.WithTimeout(context.Background(), *resourceDefaultTimeouts.Default)

	if cancel != nil {
		defer cancel()
	}

	if idOk {
		/* search by ID */
		cluster, _, err = client.KubernetesApi.K8sFindByClusterId(ctx, id.(string)).Execute()
		if err != nil {
			return fmt.Errorf("an error occurred while fetching the k8s cluster with ID %s: %s", id.(string), err)
		}
	} else {
		/* search by name */
		var clusters ionoscloud.KubernetesClusters

		ctx, cancel := context.WithTimeout(context.Background(), *resourceDefaultTimeouts.Default)

		if cancel != nil {
			defer cancel()
		}

		clusters, _, err := client.KubernetesApi.K8sGet(ctx).Execute()
		if err != nil {
			return fmt.Errorf("an error occurred while fetching k8s clusters: %s", err.Error())
		}

		if clusters.Items != nil {
			for _, c := range *clusters.Items {
				tmpCluster, _, err := client.KubernetesApi.K8sFindByClusterId(ctx, *c.Id).Execute()
				if err != nil {
					return fmt.Errorf("an error occurred while fetching k8s cluster with ID %s: %s", *c.Id, err.Error())
				}
				if tmpCluster.Properties.Name != nil {
					if strings.Contains(*tmpCluster.Properties.Name, name.(string)) {
						/* lan found */
						cluster = tmpCluster
						break
					}
				}

			}
		}

	}

	if &cluster == nil {
		return errors.New("k8s cluster not found")
	}

	if err = setK8sClusterData(d, &cluster, client); err != nil {
		return err
	}

	return nil
}

func setK8sClusterData(d *schema.ResourceData, cluster *ionoscloud.KubernetesCluster, client *ionoscloud.APIClient) error {

	if cluster.Id != nil {
		d.SetId(*cluster.Id)
		if err := d.Set("id", *cluster.Id); err != nil {
			return err
		}
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

		if cluster.Properties.AvailableUpgradeVersions != nil {
			availableUpgradeVersions := make([]interface{}, len(*cluster.Properties.AvailableUpgradeVersions), len(*cluster.Properties.AvailableUpgradeVersions))
			for i, availableUpgradeVersion := range *cluster.Properties.AvailableUpgradeVersions {
				availableUpgradeVersions[i] = availableUpgradeVersion
			}
			if err := d.Set("available_upgrade_versions", availableUpgradeVersions); err != nil {
				return err
			}
		}

		if cluster.Properties.ViableNodePoolVersions != nil {
			viableNodePoolVersions := make([]interface{}, len(*cluster.Properties.ViableNodePoolVersions), len(*cluster.Properties.ViableNodePoolVersions))
			for i, viableNodePoolVersion := range *cluster.Properties.ViableNodePoolVersions {
				viableNodePoolVersions[i] = viableNodePoolVersion
			}
			if err := d.Set("viable_node_pool_versions", viableNodePoolVersions); err != nil {
				return err
			}
		}

		if cluster.Properties.Public != nil {
			err := d.Set("public", *cluster.Properties.Public)
			if err != nil {
				return fmt.Errorf("Error while setting public property for cluser %s: %s", d.Id(), err)
			}
		}

		if cluster.Properties.GatewayIp != nil {
			err := d.Set("gateway_ip", *cluster.Properties.GatewayIp)
			if err != nil {
				return fmt.Errorf("Error while setting gateway_ip property for cluser %s: %s", d.Id(), err)
			}
		}

	}

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

	/* get and set the kubeconfig */
	if cluster.Id != nil {
		kubeConfig, _, err := client.KubernetesApi.K8sKubeconfigGet(ctx, *cluster.Id).Execute()
		if err != nil {
			return fmt.Errorf("an error occurred while fetching the kubernetes config for cluster with ID %s: %s", *cluster.Id, err)
		}

		if kubeConfig.Properties.Kubeconfig != nil {
			if err := d.Set("kube_config", *kubeConfig.Properties.Kubeconfig); err != nil {
				return err
			}
		}

		/* getting node pools */
		clusterNodePools, _, err := client.KubernetesApi.K8sNodepoolsGet(ctx, *cluster.Id).Execute()
		if err != nil {
			return fmt.Errorf("an error occurred while fetching the kubernetes cluster node pools for cluster with ID %s: %s", *cluster.Id, err)
		}

		nodePools := make([]interface{}, 0)

		if clusterNodePools.Items != nil && len(*clusterNodePools.Items) > 0 {
			nodePools = make([]interface{}, len(*clusterNodePools.Items), len(*clusterNodePools.Items))
			for i, nodePool := range *clusterNodePools.Items {
				nodePools[i] = *nodePool.Id
			}
		}

		if err := d.Set("node_pools", nodePools); err != nil {
			return err
		}
	}

	return nil
}
