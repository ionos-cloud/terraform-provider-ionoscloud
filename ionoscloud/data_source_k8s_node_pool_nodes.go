package ionoscloud

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	mongo "github.com/ionos-cloud/sdk-go-dbaas-mongo"
	ionoscloud "github.com/ionos-cloud/sdk-go/v6"

	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services"
)

func dataSourceK8sNodePoolNodes() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceK8sReadNodePoolNodes,
		Schema: map[string]*schema.Schema{
			"k8s_cluster_id": {
				Type:             schema.TypeString,
				Required:         true,
				Description:      "The UUID of an existing kubernetes cluster",
				ValidateDiagFunc: validation.ToDiagFunc(validation.IsUUID),
			},
			"node_pool_id": {
				Type:             schema.TypeString,
				Required:         true,
				Description:      "The UUID of an existing nodepool",
				ValidateDiagFunc: validation.ToDiagFunc(validation.IsUUID),
			},
			"nodes": {
				Type:        schema.TypeList,
				Description: "list of nodes in the nodepool",
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"name": {
							Type:        schema.TypeString,
							Description: "The kubernetes node name",
							Optional:    true,
						},
						"public_ip": {
							Type:        schema.TypeString,
							Description: "A valid public IP",
							Optional:    true,
						},
						"private_ip": {
							Type:        schema.TypeString,
							Description: "A valid private IP",
							Optional:    true,
						},
						"k8s_version": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The kubernetes version",
						},
					},
				},
			},
		},
		Timeouts: &resourceDefaultTimeouts,
	}
}

func dataSourceK8sReadNodePoolNodes(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(services.SdkBundle).CloudAPIClient
	clusterId := d.Get("k8s_cluster_id")
	nodePoolId := d.Get("node_pool_id")
	nodePoolIdStr := nodePoolId.(string)
	d.SetId(nodePoolIdStr)
	nodesList, apiResponse, err := client.KubernetesApi.K8sNodepoolsNodesGet(ctx, clusterId.(string), nodePoolIdStr).Depth(1).Execute()
	logApiRequestTime(apiResponse)
	if err != nil {
		return diag.FromErr(fmt.Errorf("an error occurred while fetching k8s nodes: %w", err))
	}
	if nodesList.Items == nil {
		return diag.FromErr(fmt.Errorf("no nodes found for nodepool with id %s", nodePoolIdStr))
	}
	if len(*nodesList.Items) == 0 {
		return diag.FromErr(fmt.Errorf("nodes list is empty for of nodepool with id %s", nodePoolIdStr))
	}
	if len(*nodesList.Items) > 0 {
		var nodes []interface{}
		for _, node := range *nodesList.Items {
			nodeMap := setK8sNodesDataToMap(node)
			nodes = append(nodes, nodeMap)
		}
		err := d.Set("nodes", nodes)
		if err != nil {
			diags := diag.FromErr(fmt.Errorf("error while setting nodes: %w", err))
			return diags
		}
	}
	return nil
}

func setK8sNodesDataToMap(node ionoscloud.KubernetesNode) map[string]interface{} {
	nodeEntry := make(map[string]interface{})
	if node.Id != nil {
		nodeEntry["id"] = mongo.ToValueDefault(node.Id)
	}
	if node.Properties != nil {
		nodeEntry["name"] = mongo.ToValueDefault(node.Properties.Name)
		nodeEntry["public_ip"] = mongo.ToValueDefault(node.Properties.PublicIP)
		nodeEntry["private_ip"] = mongo.ToValueDefault(node.Properties.PrivateIP)
		nodeEntry["k8s_version"] = mongo.ToValueDefault(node.Properties.K8sVersion)
	}
	return nodeEntry
}
