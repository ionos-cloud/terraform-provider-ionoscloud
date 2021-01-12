package ionoscloud

import (
	"errors"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/profitbricks/profitbricks-sdk-go/v5"
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
				Type: schema.TypeString,
				Computed: true,
			},
			"k8s_version": {
				Type: schema.TypeString,
				Computed: true,
			},
			"node_pools": {
				Type: schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"kube_config": {
				Type: schema.TypeString,
				Computed: true,
			},
		},
		Timeouts: &resourceDefaultTimeouts,
	}
}

func dataSourceK8sReadCluster(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*profitbricks.Client)

	id, idOk := d.GetOk("id")
	name, nameOk := d.GetOk("name")

	if idOk && nameOk {
		return errors.New("id and name cannot be both specified in the same time")
	}
	if !idOk && !nameOk {
		return errors.New("please provide either the lan id or name")
	}
	var cluster *profitbricks.KubernetesCluster
	var err error

	if idOk {
		/* search by ID */
		cluster, err = client.GetKubernetesCluster(id.(string))
		if err != nil {
			return fmt.Errorf("an error occurred while fetching the k8s cluster with ID %s: %s", id.(string), err)
		}
	} else {
		/* search by name */
		var clusters *profitbricks.KubernetesClusters
		clusters, err := client.ListKubernetesClusters()
		if err != nil {
			return fmt.Errorf("an error occurred while fetching k8s clusters: %s", err.Error())
		}

		for _, c := range clusters.Items {
			tmpCluster, err := client.GetKubernetesCluster(c.ID)
			if err != nil {
				return fmt.Errorf("an error occurred while fetching k8s cluster with ID %s: %s", c.ID, err.Error())
			}
			if strings.Contains(tmpCluster.Properties.Name, name.(string)) {
				/* lan found */
				cluster = tmpCluster
				break
			}
		}
	}

	if cluster == nil {
		return errors.New("k8s cluster not found")
	}

	if err = setK8sClusterData(d, cluster, client); err != nil {
		return err
	}

	return nil
}

func setK8sClusterData(d *schema.ResourceData, cluster *profitbricks.KubernetesCluster, client *profitbricks.Client) error {
	d.SetId(cluster.ID)
	if err := d.Set("id", cluster.ID); err != nil {
		return err
	}

	if err := d.Set("name", cluster.Properties.Name); err != nil {
		return err
	}
	if err := d.Set("state", cluster.Metadata.State); err != nil {
		return err
	}
	if err := d.Set("k8s_version", cluster.Properties.K8sVersion); err != nil {
		return err
	}
	if err := d.Set("state", cluster.Metadata.State); err != nil {
		return err
	}

	/* get and set the kubeconfig */
	kubeConfig, err := client.GetKubeconfig(cluster.ID)
	if err != nil {
		return fmt.Errorf("an error occurred while fetching the kubernetes config for cluster with ID %s: %s", cluster.ID, err)
	}

	if err := d.Set("kube_config", kubeConfig); err != nil {
		return err
	}

	/* getting node pools */
	clusterNodePools, err := client.ListKubernetesNodePools(cluster.ID)
	if err != nil {
		return fmt.Errorf("an error occurred while fetching the kubernetes cluster node pools for cluster with ID %s: %s", cluster.ID, err)
	}

	nodePools := make([]interface{}, 0)

	if clusterNodePools != nil && len(clusterNodePools.Items) > 0 {
		nodePools = make([]interface{}, len(clusterNodePools.Items), len(clusterNodePools.Items))
		for i, nodePool := range clusterNodePools.Items {
			nodePools[i] = nodePool.ID
		}
	}

	if err := d.Set("node_pools", nodePools); err != nil {
		return err
	}

	return nil
}

