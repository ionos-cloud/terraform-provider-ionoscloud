package ionoscloud

import (
	"context"
	"errors"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	ionoscloud "github.com/ionos-cloud/sdk-go/v6"
)

func dataSourceK8sNodePool() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceK8sReadNodePool,
		Schema: map[string]*schema.Schema{
			"k8s_cluster_id": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "The UUID of an existing kubernetes cluster",
				ValidateFunc: validation.All(validation.StringIsNotWhiteSpace),
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
							Required:    true,
						},
						"dhcp": {
							Type:        schema.TypeBool,
							Description: "Indicates if the Kubernetes Node Pool LAN will reserve an IP using DHCP",
							Optional:    true,
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
										Optional:    true,
									},
									"gateway_ip": {
										Type:        schema.TypeString,
										Description: "IPv4 or IPv6 Gateway IP for the route",
										Optional:    true,
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
		},
		Timeouts: &resourceDefaultTimeouts,
	}
}

func dataSourceK8sReadNodePool(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ionoscloud.APIClient)

	clusterId := d.Get("k8s_cluster_id")
	id, idOk := d.GetOk("id")
	name, nameOk := d.GetOk("name")

	if idOk && nameOk {
		return errors.New("id and name cannot be both specified in the same time")
	}
	if !idOk && !nameOk {
		return errors.New("please provide either the lan id or name")
	}
	var nodePool ionoscloud.KubernetesNodePool
	var err error

	if idOk {
		/* search by ID */
		ctx, cancel := context.WithTimeout(context.Background(), *resourceDefaultTimeouts.Default)

		if cancel != nil {
			defer cancel()
		}

		nodePool, _, err = client.KubernetesApi.K8sNodepoolsFindById(ctx, clusterId.(string), id.(string)).Execute()
		if err != nil {
			return fmt.Errorf("an error occurred while fetching the k8s nodePool with ID %s: %s", id.(string), err)
		}
	} else {
		/* search by name */
		var clusters ionoscloud.KubernetesNodePools

		ctx, cancel := context.WithTimeout(context.Background(), *resourceDefaultTimeouts.Default)

		if cancel != nil {
			defer cancel()
		}

		clusters, _, err := client.KubernetesApi.K8sNodepoolsGet(ctx, clusterId.(string)).Execute()
		if err != nil {
			return fmt.Errorf("an error occurred while fetching k8s nodepools: %s", err.Error())
		}

		found := false
		if clusters.Items != nil {
			for _, c := range *clusters.Items {
				tmpNodePool, _, err := client.KubernetesApi.K8sNodepoolsFindById(ctx, clusterId.(string), *c.Id).Execute()
				if err != nil {
					return fmt.Errorf("an error occurred while fetching k8s nodePool with ID %s: %s", *c.Id, err.Error())
				}
				if tmpNodePool.Properties.Name != nil && *tmpNodePool.Properties.Name == name.(string) {
					/* lan found */
					nodePool = tmpNodePool
					found = true
					break
				}
			}
		}

		if !found {
			return errors.New("k8s nodePool not found")
		}

	}

	if err = setK8sNodePoolData(d, &nodePool); err != nil {
		return err
	}

	return nil
}

func setK8sNodePoolData(d *schema.ResourceData, nodePool *ionoscloud.KubernetesNodePool) error {

	if nodePool.Id != nil {
		d.SetId(*nodePool.Id)
		if err := d.Set("id", *nodePool.Id); err != nil {
			return err
		}
	}

	if nodePool.Properties != nil {
		if nodePool.Properties.Name != nil {
			if err := d.Set("name", *nodePool.Properties.Name); err != nil {
				return err
			}
		}

		if nodePool.Properties.K8sVersion != nil {
			if err := d.Set("k8s_version", *nodePool.Properties.K8sVersion); err != nil {
				return err
			}
		}

		if nodePool.Properties.DatacenterId != nil {
			if err := d.Set("datacenter_id", *nodePool.Properties.DatacenterId); err != nil {
				return err
			}
		}

		if nodePool.Properties.NodeCount != nil {
			if err := d.Set("node_count", *nodePool.Properties.NodeCount); err != nil {
				return err
			}
		}

		if nodePool.Properties.CpuFamily != nil {
			if err := d.Set("cpu_family", *nodePool.Properties.CpuFamily); err != nil {
				return err
			}
		}

		if nodePool.Properties.CoresCount != nil {
			if err := d.Set("cores_count", *nodePool.Properties.CoresCount); err != nil {
				return err
			}
		}

		if nodePool.Properties.RamSize != nil {
			if err := d.Set("ram_size", *nodePool.Properties.RamSize); err != nil {
				return err
			}
		}

		if nodePool.Properties.AvailabilityZone != nil {
			if err := d.Set("availability_zone", *nodePool.Properties.AvailabilityZone); err != nil {
				return err
			}
		}

		if nodePool.Properties.StorageType != nil {
			if err := d.Set("storage_type", *nodePool.Properties.StorageType); err != nil {
				return err
			}
		}

		if nodePool.Properties.StorageSize != nil {
			if err := d.Set("storage_size", *nodePool.Properties.StorageSize); err != nil {
				return err
			}
		}

		if nodePool.Properties.K8sVersion != nil {
			if err := d.Set("k8s_version", *nodePool.Properties.K8sVersion); err != nil {
				return err
			}
		}

		if nodePool.Properties.PublicIps != nil {
			if err := d.Set("public_ips", *nodePool.Properties.PublicIps); err != nil {
				return err
			}
		}

		if nodePool.Properties.MaintenanceWindow != nil && nodePool.Properties.MaintenanceWindow.Time != nil && nodePool.Properties.MaintenanceWindow.DayOfTheWeek != nil {
			if err := d.Set("maintenance_window", []map[string]string{
				{
					"time":            *nodePool.Properties.MaintenanceWindow.Time,
					"day_of_the_week": *nodePool.Properties.MaintenanceWindow.DayOfTheWeek,
				},
			}); err != nil {
				return err
			}
		}

		if nodePool.Properties.AutoScaling != nil && nodePool.Properties.AutoScaling.MinNodeCount != nil &&
			nodePool.Properties.AutoScaling.MaxNodeCount != nil && (*nodePool.Properties.AutoScaling.MinNodeCount != 0 &&
			*nodePool.Properties.AutoScaling.MaxNodeCount != 0) {
			if err := d.Set("auto_scaling", []map[string]uint32{
				{
					"min_node_count": uint32(*nodePool.Properties.AutoScaling.MinNodeCount),
					"max_node_count": uint32(*nodePool.Properties.AutoScaling.MaxNodeCount),
				},
			}); err != nil {
				return err
			}
		}

		nodePoolLans := make([]interface{}, 0)
		if nodePool.Properties.Lans != nil && len(*nodePool.Properties.Lans) > 0 {
			nodePoolLans = make([]interface{}, len(*nodePool.Properties.Lans))
			for lanIndex, nodePoolLan := range *nodePool.Properties.Lans {
				lanEntry := make(map[string]interface{})

				if nodePoolLan.Id != nil {
					lanEntry["id"] = *nodePoolLan.Id
				}

				if nodePoolLan.Dhcp != nil {
					lanEntry["dhcp"] = *nodePoolLan.Dhcp
				}

				nodePoolRoutes := make([]interface{}, 0)
				if len(*nodePoolLan.Routes) > 0 {
					nodePoolRoutes = make([]interface{}, len(*nodePoolLan.Routes))
					for routeIndex, nodePoolRoute := range *nodePoolLan.Routes {
						routeEntry := make(map[string]string)
						if nodePoolRoute.Network != nil {
							routeEntry["network"] = *nodePoolRoute.Network
						}
						if nodePoolRoute.GatewayIp != nil {
							routeEntry["gateway_ip"] = *nodePoolRoute.GatewayIp
						}
						nodePoolRoutes[routeIndex] = routeEntry
					}
				}

				if len(nodePoolRoutes) > 0 {
					lanEntry["routes"] = nodePoolRoutes
				}

				nodePoolLans[lanIndex] = lanEntry
			}
		}

		if len(nodePoolLans) > 0 {
			if err := d.Set("lans", nodePoolLans); err != nil {
				return fmt.Errorf("error while setting lans property for k8sNodepool %s: %s", d.Id(), err)
			}
		}
	}

	if nodePool.Metadata != nil && nodePool.Metadata.State != nil {
		if err := d.Set("state", *nodePool.Metadata.State); err != nil {
			return err
		}
	}

	return nil
}
