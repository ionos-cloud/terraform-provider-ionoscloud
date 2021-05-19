package ionoscloud

import (
	"errors"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/profitbricks/profitbricks-sdk-go/v5"
)

func dataSourceK8sNodePool() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceK8sReadNodePool,
		Schema: map[string]*schema.Schema{
			"k8s_cluster_id": {
				Type:     schema.TypeString,
				Required: true,
				Description: "The UUID of an existing kubernetes cluster",
			},
			"id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"name": {
				Type:     schema.TypeString,
				Description: "The desired name for the node pool",
				Optional: true,
			},
			"datacenter_id": {
				Type:     schema.TypeString,
				Computed: true,
				Description: "The UUID of the VDC",
			},
			"state": {
				Type: schema.TypeString,
				Computed: true,
			},
			"node_count": {
				Type: schema.TypeInt,
				Description: "The number of nodes in this node pool",
				Computed: true,
			},
			"cpu_family": {
				Type: schema.TypeString,
				Description: "CPU Family",
				Computed: true,
			},
			"cores_count": {
				Type: schema.TypeInt,
				Description: "CPU cores count",
				Computed: true,
			},
			"ram_size": {
				Type: schema.TypeInt,
				Description: "The amount of RAM in MB",
				Computed: true,
			},
			"availability_zone": {
				Type: schema.TypeString,
				Description: "The compute availability zone in which the nodes should exist",
				Computed: true,
			},
			"storage_type": {
				Type: schema.TypeString,
				Description: "Storage type to use",
				Computed: true,
			},
			"storage_size": {
				Type: schema.TypeInt,
				Computed: true,
			},
			"k8s_version": {
				Type: schema.TypeString,
				Computed: true,
				Description: "The kubernetes version",
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
			"auto_scaling": {
				Type:        schema.TypeList,
				Description: "The range defining the minimum and maximum number of worker nodes that the managed node group can scale in",
				Computed:    true,
				MaxItems:    1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"min_node_count": {
							Type:        schema.TypeInt,
							Description: "The minimum number of worker nodes the node pool can scale down to. Should be less than max_node_count",
							Required:    true,
						},
						"max_node_count": {
							Type:        schema.TypeInt,
							Description: "The maximum number of worker nodes that the node pool can scale to. Should be greater than min_node_count",
							Required:    true,
						},
					},
				},
			},
			"lans": {
				Type:        schema.TypeList,
				Description: "A list of Local Area Networks the node pool should be part of",
				Computed:    true,
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
			},
			"labels": {
				Type: schema.TypeMap,
				Computed: true,
				Elem: &schema.Schema{Type: schema.TypeString},
			},
			"annotations": {
				Type: schema.TypeMap,
				Computed: true,
				Elem: &schema.Schema{Type: schema.TypeString},
			},
			"available_upgrade_versions": {
				Type: 			schema.TypeList,
				Description:	"A list of kubernetes versions available for upgrade",
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"public_ips": {
				Type:		schema.TypeList,
				Description: "A list of fixed IPs",
				Optional:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},

		},
		Timeouts: &resourceDefaultTimeouts,
	}
}

func dataSourceK8sReadNodePool(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*profitbricks.Client)

	clusterId := d.Get("k8s_cluster_id")
	id, idOk := d.GetOk("id")
	name, nameOk := d.GetOk("name")

	if idOk && nameOk {
		return errors.New("id and name cannot be both specified in the same time")
	}
	if !idOk && !nameOk {
		return errors.New("please provide either the lan id or name")
	}
	var nodePool *profitbricks.KubernetesNodePool
	var err error

	if idOk {
		/* search by ID */
		nodePool, err = client.GetKubernetesNodePool(clusterId.(string), id.(string))
		if err != nil {
			return fmt.Errorf("an error occurred while fetching the k8s nodePool with ID %s: %s", id.(string), err)
		}
	} else {
		/* search by name */
		var clusters *profitbricks.KubernetesNodePools
		clusters, err := client.ListKubernetesNodePools(clusterId.(string))
		if err != nil {
			return fmt.Errorf("an error occurred while fetching k8s nodepools: %s", err.Error())
		}

		for _, c := range clusters.Items {
			tmpNodePool, err := client.GetKubernetesNodePool(clusterId.(string), c.ID)
			if err != nil {
				return fmt.Errorf("an error occurred while fetching k8s nodePool with ID %s: %s", c.ID, err.Error())
			}
			if tmpNodePool.Properties.Name == name.(string) {
				/* lan found */
				nodePool = tmpNodePool
				break
			}
		}
	}

	if nodePool == nil {
		return errors.New("k8s nodePool not found")
	}

	if err = setK8sNodePoolData(d, nodePool); err != nil {
		return err
	}

	return nil
}

func setK8sNodePoolData(d *schema.ResourceData, nodePool *profitbricks.KubernetesNodePool) error {
	d.SetId(nodePool.ID)
	if err := d.Set("id", nodePool.ID); err != nil {
		return err
	}

	if err := d.Set("name", nodePool.Properties.Name); err != nil {
		return err
	}
	if err := d.Set("state", nodePool.Metadata.State); err != nil {
		return err
	}
	if err := d.Set("k8s_version", nodePool.Properties.K8sVersion); err != nil {
		return err
	}
	if err := d.Set("state", nodePool.Metadata.State); err != nil {
		return err
	}
	if err := d.Set("datacenter_id", nodePool.Properties.DatacenterID); err != nil {
		return err
	}
	if err := d.Set("node_count", nodePool.Properties.NodeCount); err != nil {
		return err
	}
	if err := d.Set("cpu_family", nodePool.Properties.CPUFamily); err != nil {
		return err
	}
	if err := d.Set("cores_count", nodePool.Properties.CoresCount); err != nil {
		return err
	}
	if err := d.Set("ram_size", nodePool.Properties.RAMSize); err != nil {
		return err
	}
	if err := d.Set("availability_zone", nodePool.Properties.AvailabilityZone); err != nil {
		return err
	}
	if err := d.Set("storage_type", nodePool.Properties.StorageType); err != nil {
		return err
	}
	if err := d.Set("storage_size", nodePool.Properties.StorageSize); err != nil {
		return err
	}
	if err := d.Set("k8s_version", nodePool.Properties.K8sVersion); err != nil {
		return err
	}
	if err := d.Set("public_ips", nodePool.Properties.PublicIPs); err != nil {
		return err
	}

	if nodePool.Properties.MaintenanceWindow != nil {
		if err := d.Set("maintenance_window", []map[string]string{
			{
				"time": nodePool.Properties.MaintenanceWindow.Time,
				"day_of_the_week": nodePool.Properties.MaintenanceWindow.DayOfTheWeek,
			},
		}); err != nil {
			return err
		}
	}

	if nodePool.Properties.AutoScaling != nil && (nodePool.Properties.AutoScaling.MinNodeCount != 0 && nodePool.Properties.AutoScaling.MaxNodeCount != 0) {
		if err := d.Set("auto_scaling", []map[string]uint32{
			{
				"min_node_count": nodePool.Properties.AutoScaling.MinNodeCount,
				"max_node_count": nodePool.Properties.AutoScaling.MaxNodeCount,
			},
		}); err != nil {
			return err
		}
	}

	if err := d.Set("lans", nodePool.Properties.LANs); err != nil {
		return err
	}

	return nil
}

