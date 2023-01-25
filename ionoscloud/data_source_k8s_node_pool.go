package ionoscloud

import (
	"context"
	"errors"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	ionoscloud "github.com/ionos-cloud/sdk-go/v6"
	"log"
	"strings"
)

func dataSourceK8sNodePool() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceK8sReadNodePool,
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
			"partial_match": {
				Type:        schema.TypeBool,
				Description: "Whether partial matching is allowed or not when using name argument.",
				Default:     false,
				Optional:    true,
			},
			"datacenter_id": {
				Type:        schema.TypeString,
				Optional:    true,
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
				Optional:    true,
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
			//"gateway_ip": {
			//	Type:        schema.TypeString,
			//	Description: "Public IP address for the gateway performing source NAT for the node pool's nodes belonging to a private cluster. Required only if the node pool belongs to a private cluster.",
			//	Computed:    true,
			//},
		},
		Timeouts: &resourceDefaultTimeouts,
	}
}

func dataSourceK8sReadNodePool(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(SdkBundle).CloudApiClient

	clusterId := d.Get("k8s_cluster_id").(string)
	idValue, idOk := d.GetOk("id")
	nameValue, nameOk := d.GetOk("name")
	dcIdValue, dcIdOk := d.GetOk("datacenter_id")
	avZoneValue, avZoneOk := d.GetOk("availability_zone")

	id := idValue.(string)
	name := nameValue.(string)
	dcId := dcIdValue.(string)
	avZone := avZoneValue.(string)

	if idOk && (nameOk || dcIdOk || avZoneOk) {
		return diag.FromErr(errors.New("id and name/datacenter_id/availability_zone cannot be both specified in the same time, choose between id or a combination of other parameters"))
	}
	if !idOk && !nameOk && !dcIdOk && !avZoneOk {
		return diag.FromErr(errors.New("please provide either the lan id or other parameter like name or datacenter_id"))
	}

	var nodePool ionoscloud.KubernetesNodePool
	var err error
	var apiResponse *ionoscloud.APIResponse

	if idOk {
		/* search by ID */
		log.Printf("[INFO] Using data source for k8s nodepool by id %s", id)
		nodePool, apiResponse, err = client.KubernetesApi.K8sNodepoolsFindById(ctx, clusterId, id).Execute()
		logApiRequestTime(apiResponse)
		if err != nil {
			return diag.FromErr(fmt.Errorf("an error occurred while fetching the k8s nodePool with ID %s: %w", id, err))
		}
	} else {
		/* search by name */
		var results []ionoscloud.KubernetesNodePool

		if nameOk {
			partialMatch := d.Get("partial_match").(bool)

			log.Printf("[INFO] Using data source for k8s nodepool by name with partial_match %t and name: %s", partialMatch, name)

			if partialMatch {
				nodePools, apiResponse, err := client.KubernetesApi.K8sNodepoolsGet(ctx, clusterId).Depth(1).Filter("name", name).Execute()
				logApiRequestTime(apiResponse)
				if err != nil {
					return diag.FromErr(fmt.Errorf("an error occurred while fetching k8s nodepools: %s", err.Error()))
				}
				if len(*nodePools.Items) == 0 {
					return diag.FromErr(fmt.Errorf("no result found with the specified criteria name with partial match: %s", name))
				}
				results = *nodePools.Items
			} else {
				nodePools, apiResponse, err := client.KubernetesApi.K8sNodepoolsGet(ctx, clusterId).Depth(1).Execute()
				logApiRequestTime(apiResponse)
				if err != nil {
					return diag.FromErr(fmt.Errorf("an error occurred while fetching k8s nodepools: %s", err.Error()))
				}

				if nodePools.Items != nil {
					var nameResults []ionoscloud.KubernetesNodePool
					for _, c := range *nodePools.Items {
						if c.Properties != nil && c.Properties.Name != nil && strings.EqualFold(*c.Properties.Name, name) {
							tmpNodePool, apiResponse, err := client.KubernetesApi.K8sNodepoolsFindById(ctx, clusterId, *c.Id).Execute()
							logApiRequestTime(apiResponse)
							if err != nil {
								return diag.FromErr(fmt.Errorf("an error occurred while fetching k8s nodepool with ID %s: %w", *c.Id, err))
							}
							/* lan found */
							nameResults = append(nameResults, tmpNodePool)
							break
						}
					}
					if len(nameResults) == 0 {
						return diag.FromErr(fmt.Errorf("no result found with the specified criteria name: %s", name))
					}
					results = nameResults
				}
			}
		} else {
			nodePools, apiResponse, err := client.KubernetesApi.K8sNodepoolsGet(ctx, clusterId).Depth(1).Execute()
			logApiRequestTime(apiResponse)
			if err != nil {
				return diag.FromErr(fmt.Errorf("an error occurred while fetching k8s nodepools: %s", err.Error()))
			}
			results = *nodePools.Items
		}

		if dcIdOk && dcId != "" {
			var dcIdResults []ionoscloud.KubernetesNodePool
			if results != nil {
				for _, k8sNodepool := range results {
					if k8sNodepool.Properties != nil && k8sNodepool.Properties.DatacenterId != nil && strings.EqualFold(*k8sNodepool.Properties.DatacenterId, dcId) {
						dcIdResults = append(dcIdResults, k8sNodepool)
					}
				}
			}
			if dcIdResults == nil || len(dcIdResults) == 0 {
				return diag.FromErr(fmt.Errorf("no result found with the specified criteria: datacenter_id = %s", dcId))
			}
			results = dcIdResults
		}

		if avZoneOk && avZone != "" {
			var avZoneResults []ionoscloud.KubernetesNodePool
			if results != nil {
				for _, k8sNodepool := range results {
					if k8sNodepool.Properties != nil && k8sNodepool.Properties.AvailabilityZone != nil && strings.EqualFold(*k8sNodepool.Properties.AvailabilityZone, avZone) {
						avZoneResults = append(avZoneResults, k8sNodepool)
					}
				}
			}
			if avZoneResults == nil || len(avZoneResults) == 0 {
				return diag.FromErr(fmt.Errorf("no result found with the specified criteria: availability_zone = %s", avZone))
			}
			results = avZoneResults
		}

		if len(results) > 1 {
			return diag.FromErr(fmt.Errorf("more than one nodepool found with the specified name %s", name))
		} else {
			nodePool = results[0]
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
