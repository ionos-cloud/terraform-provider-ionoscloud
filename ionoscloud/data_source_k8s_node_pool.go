package ionoscloud

import (
	"context"
	"errors"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	ionoscloud "github.com/ionos-cloud/sdk-go/v6"

	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services"
)

func dataSourceK8sNodePool() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceK8sReadNodePool,
		Schema: map[string]*schema.Schema{
			"k8s_cluster_id": {
				Type:             schema.TypeString,
				Required:         true,
				Description:      "The UUID of an existing kubernetes cluster",
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringIsNotWhiteSpace),
			},
			"id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"name": {
				Type:        schema.TypeString,
				Description: "The desired name for the node pool",
				Optional:    true,
			},
			"datacenter_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The UUID of the VDC",
			},
			"state": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"node_count": {
				Type:        schema.TypeInt,
				Description: "The number of nodes in this node pool",
				Computed:    true,
			},
			"cpu_family": {
				Type:        schema.TypeString,
				Description: "CPU Family",
				Computed:    true,
			},
			"cores_count": {
				Type:        schema.TypeInt,
				Description: "CPU cores count",
				Computed:    true,
			},
			"ram_size": {
				Type:        schema.TypeInt,
				Description: "The amount of RAM in MB",
				Computed:    true,
			},
			"availability_zone": {
				Type:        schema.TypeString,
				Description: "The compute availability zone in which the nodes should exist",
				Computed:    true,
			},
			"storage_type": {
				Type:        schema.TypeString,
				Description: "Storage type to use",
				Computed:    true,
			},
			"storage_size": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"k8s_version": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The kubernetes version",
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
			"auto_scaling": {
				Type:        schema.TypeList,
				Description: "The range defining the minimum and maximum number of worker nodes that the managed node group can scale in",
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"min_node_count": {
							Type:        schema.TypeInt,
							Description: "The minimum number of worker nodes the node pool can scale down to. Should be less than max_node_count",
							Computed:    true,
						},
						"max_node_count": {
							Type:        schema.TypeInt,
							Description: "The maximum number of worker nodes that the node pool can scale to. Should be greater than min_node_count",
							Computed:    true,
						},
					},
				},
			},
			"lans": {
				Type:        schema.TypeList,
				Description: "A list of Local Area Networks the node pool should be part of",
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeInt,
							Description: "The LAN ID of an existing LAN at the related datacenter",
							Computed:    true,
						},
						"dhcp": {
							Type:        schema.TypeBool,
							Description: "Indicates if the Kubernetes Node Pool LAN will reserve an IP using DHCP",
							Computed:    true,
						},
						"routes": {
							Type:        schema.TypeList,
							Description: "An array of additional LANs attached to worker nodes",
							Optional:    true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"network": {
										Type:        schema.TypeString,
										Description: "IPv4 or IPv6 CIDR to be routed via the interface",
										Computed:    true,
									},
									"gateway_ip": {
										Type:        schema.TypeString,
										Description: "IPv4 or IPv6 Gateway IP for the route",
										Computed:    true,
									},
								},
							},
						},
					},
				},
			},
			"labels": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"annotations": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"available_upgrade_versions": {
				Type:        schema.TypeList,
				Description: "A list of kubernetes versions available for upgrade",
				Computed:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"public_ips": {
				Type:        schema.TypeList,
				Description: "A list of fixed IPs",
				Computed:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			// "gateway_ip": {
			//	Type:        schema.TypeString,
			//	Description: "Public IP address for the gateway performing source NAT for the node pool's nodes belonging to a private cluster. Required only if the node pool belongs to a private cluster.",
			//	Computed:    true,
			// },
		},
		Timeouts: &resourceDefaultTimeouts,
	}
}

func dataSourceK8sReadNodePool(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(services.SdkBundle).CloudApiClient

	clusterId := d.Get("k8s_cluster_id")
	id, idOk := d.GetOk("id")
	name, nameOk := d.GetOk("name")

	if idOk && nameOk {
		return diag.FromErr(errors.New("id and name cannot be both specified in the same time"))
	}
	if !idOk && !nameOk {
		return diag.FromErr(errors.New("please provide either the lan id or name"))
	}
	var nodePool ionoscloud.KubernetesNodePool
	var err error
	var apiResponse *ionoscloud.APIResponse
	if idOk {
		/* search by ID */
		nodePool, apiResponse, err = client.KubernetesApi.K8sNodepoolsFindById(ctx, clusterId.(string), id.(string)).Execute()
		logApiRequestTime(apiResponse)
		if err != nil {
			return diag.FromErr(fmt.Errorf("an error occurred while fetching the k8s nodePool with ID %s: %w", id.(string), err))
		}
	} else {
		/* search by name */
		var nodePools ionoscloud.KubernetesNodePools

		nodePools, apiResponse, err := client.KubernetesApi.K8sNodepoolsGet(ctx, clusterId.(string)).Depth(1).Execute()
		logApiRequestTime(apiResponse)
		if err != nil {
			return diag.FromErr(fmt.Errorf("an error occurred while fetching k8s nodepools: %w", err))
		}

		if nodePools.Items != nil {
			var results []ionoscloud.KubernetesNodePool

			for _, c := range *nodePools.Items {
				if c.Properties != nil && c.Properties.Name != nil && *c.Properties.Name == name.(string) {
					tmpNodePool, apiResponse, err := client.KubernetesApi.K8sNodepoolsFindById(ctx, clusterId.(string), *c.Id).Execute()
					logApiRequestTime(apiResponse)
					if err != nil {
						return diag.FromErr(fmt.Errorf("an error occurred while fetching k8s nodePool with ID %s: %w", *c.Id, err))
					}
					/* lan found */
					results = append(results, tmpNodePool)
					break
				}
			}

			if len(results) == 0 {
				return diag.FromErr(fmt.Errorf("no nodepool found with the specified name %s", name.(string)))
			} else if len(results) > 1 {
				return diag.FromErr(fmt.Errorf("more than one nodepool found with the specified name %s", name.(string)))
			} else {
				nodePool = results[0]
			}
		}
	}

	if err = setK8sNodePoolData(d, &nodePool); err != nil {
		return diag.FromErr(err)
	}

	if nodePool.Metadata != nil && nodePool.Metadata.State != nil {
		if err := d.Set("state", *nodePool.Metadata.State); err != nil {
			return diag.FromErr(err)
		}
	}

	if nodePool.Properties.AvailableUpgradeVersions != nil && len(*nodePool.Properties.AvailableUpgradeVersions) > 0 {
		if err := d.Set("available_upgrade_versions", *nodePool.Properties.AvailableUpgradeVersions); err != nil {
			return diag.FromErr(err)
		}
	}

	return nil
}
