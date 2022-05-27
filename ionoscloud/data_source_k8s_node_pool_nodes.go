package ionoscloud

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	ionoscloud "github.com/ionos-cloud/sdk-go/v6"
)

func dataSourceK8sNodePoolNodes() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceK8sReadNodePoolNodes,
		Schema: map[string]*schema.Schema{
			"k8s_cluster_id": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "The UUID of an existing kubernetes cluster",
				ValidateFunc: validation.All(validation.IsUUID),
			},
			"node_pool_id": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "The UUID of an existing nodepool",
				ValidateFunc: validation.All(validation.IsUUID),
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
	client := meta.(SdkBundle).CloudApiClient
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
		return diag.FromErr(fmt.Errorf("no nodes found for nodepool with id %s ", nodePoolIdStr))
	}
	if len(*nodesList.Items) > 0 {
		var nodes []interface{}
		for _, node := range *nodesList.Items {
			nodeMap := setK8sNodesDataToMap(node)
			nodes = append(nodes, nodeMap)
		}
		err := d.Set("nodes", nodes)
		if err != nil {
			diags := diag.FromErr(fmt.Errorf("error while setting nodes: %s", err))
			return diags
		}
	}

	return nil
}

func setK8sNodesDataToMap(node ionoscloud.KubernetesNode) map[string]interface{} {
	nodeEntry := make(map[string]interface{})
	if node.Id != nil {
		nodeEntry["id"] = stringOrDefault(node.Id, "")
	}
	if node.Properties != nil {
		nodeEntry["name"] = stringOrDefault(node.Properties.Name, "")
		nodeEntry["public_ip"] = stringOrDefault(node.Properties.PublicIP, "")
		nodeEntry["private_ip"] = stringOrDefault(node.Properties.PrivateIP, "")
		nodeEntry["k8s_version"] = stringOrDefault(node.Properties.K8sVersion, "")
	}

	return nodeEntry
}
