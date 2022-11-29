package ionoscloud

import (
	"context"
	"errors"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	dataplatform "github.com/ionos-cloud/sdk-go-dataplatform"
	dataplatformService "github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/dataplatform"
	"log"
	"regexp"
	"strings"
)

func dataSourceDataplatformNodePool() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDataplatformReadNodePool,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.All(validation.StringMatch(regexp.MustCompile("^[A-Za-z0-9][-A-Za-z0-9_.]*[A-Za-z0-9]$"), "")),
			},
			"name": {
				Type:        schema.TypeString,
				Description: "The name of your node pool.",
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
				Description: "The version of the Data Platform.",
				Computed:    true,
			},
			"datacenter_id": {
				Type:        schema.TypeString,
				Description: "The UUID of the virtual data center (VDC) the cluster is provisioned.",
				Computed:    true,
			},
			"node_count": {
				Type:        schema.TypeInt,
				Description: "The number of nodes that make up the node pool.",
				Computed:    true,
			},
			"cpu_family": {
				Type:        schema.TypeString,
				Description: "A valid CPU family name or `AUTO` if the platform shall choose the best fitting option. Available CPU architectures can be retrieved from the datacenter resource.",
				Computed:    true,
			},
			"cores_count": {
				Type:        schema.TypeInt,
				Description: "The number of CPU cores per node.",
				Computed:    true,
			},
			"ram_size": {
				Type:        schema.TypeInt,
				Description: "The RAM size for one node in MB. Must be set in multiples of 1024 MB, with a minimum size is of 2048 MB.",
				Computed:    true,
			},
			"availability_zone": {
				Type:        schema.TypeString,
				Description: "The availability zone of the virtual datacenter region where the node pool resources should be provisioned.",
				Computed:    true,
			},
			"storage_type": {
				Type:        schema.TypeString,
				Description: "The type of hardware for the volume.",
				Computed:    true,
			},
			"storage_size": {
				Type:        schema.TypeInt,
				Description: "The size of the volume in GB. The size must be greater than 10GB.",
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
			"labels": {
				Type:        schema.TypeMap,
				Description: "Key-value pairs attached to the node pool resource as [Kubernetes labels](https://kubernetes.io/docs/concepts/overview/working-with-objects/labels/)",
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"annotations": {
				Type:        schema.TypeMap,
				Description: "Key-value pairs attached to node pool resource as [Kubernetes annotations](https://kubernetes.io/docs/concepts/overview/working-with-objects/annotations/)",
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"cluster_id": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "The UUID of an existing Dataplatform cluster",
				ValidateFunc: validation.All(validation.StringMatch(regexp.MustCompile("^[A-Za-z0-9][-A-Za-z0-9_.]*[A-Za-z0-9]$"), "")),
			},
		},
		Timeouts: &resourceDefaultTimeouts,
	}
}

func dataSourceDataplatformReadNodePool(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(SdkBundle).DataplatformClient

	clusterId := d.Get("cluster_id").(string)
	idValue, idOk := d.GetOk("id")
	nameValue, nameOk := d.GetOk("name")

	id := idValue.(string)
	name := nameValue.(string)

	if idOk && nameOk {
		return diag.FromErr(errors.New("id and name cannot be both specified in the same time"))
	}
	if !idOk && !nameOk {
		return diag.FromErr(errors.New("please provide either the Dataplatform Node Pool id or name"))
	}

	var nodePool dataplatform.NodePoolResponseData
	var err error

	if idOk {
		/* search by ID */
		nodePool, _, err = client.GetNodePool(ctx, clusterId, id)
		if err != nil {
			return diag.FromErr(fmt.Errorf("an error occurred while fetching the Dataplatform Node Pool with ID %s: %s", id, err))
		}
	} else {
		/* search by name */
		results, err := filterNodePools(ctx, d, client, name)

		if err != nil {
			return err
		}

		if results == nil || len(results) == 0 {
			return diag.FromErr(fmt.Errorf("no Dataplatform NodePool found with the specified name = %s", name))
		} else if len(results) > 1 {
			return diag.FromErr(fmt.Errorf("more than one Dataplatform NodePool found with the specified criteria name = %s", name))
		} else {
			nodePool = results[0]
		}
	}

	if err = dataplatformService.SetDataplatformNodePoolData(d, nodePool); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func filterNodePools(ctx context.Context, d *schema.ResourceData, client *dataplatformService.Client, name string) ([]dataplatform.NodePoolResponseData, diag.Diagnostics) {
	clusterId := d.Get("cluster_id").(string)

	var results []dataplatform.NodePoolResponseData

	partialMatch := d.Get("partial_match").(bool)

	log.Printf("[INFO] Using data source for Dataplatform Node Pool by name with partial_match %t and name: %s", partialMatch, name)

	nodePools, _, err := client.ListNodePools(ctx, clusterId)
	if err != nil {
		diags := diag.FromErr(fmt.Errorf("an error occurred while fetching Dataplatform NodePools: %s", err.Error()))
		return nil, diags
	}
	if nodePools.Items != nil && len(*nodePools.Items) > 0 {
		for _, nodePoolItem := range *nodePools.Items {
			if nodePoolItem.Properties != nil && nodePoolItem.Properties.Name != nil && (partialMatch && strings.Contains(*nodePoolItem.Properties.Name, name) ||
				!partialMatch && strings.EqualFold(*nodePoolItem.Properties.Name, name)) {
				tmpNodePool, _, err := client.GetNodePool(ctx, clusterId, *nodePoolItem.Id)
				if err != nil {
					return nil, diag.FromErr(fmt.Errorf("an error occurred while fetching the Dataplatform NodePool with ID: %s while searching by full name: %s: %w", *nodePoolItem.Id, name, err))
				}
				results = append(results, tmpNodePool)
			}
		}
	}

	return results, nil
}
