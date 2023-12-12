package ionoscloud

import (
	"context"
	"fmt"
	"regexp"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	dataplatform "github.com/ionos-cloud/sdk-go-dataplatform"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services"
	dataplatformService "github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/dataplatform"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils/constant"
)

func dataSourceDataplatformNodePools() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceNodePoolsRead,
		Schema: map[string]*schema.Schema{
			"cluster_id": {
				Type:             schema.TypeString,
				Required:         true,
				Description:      "The UUID of an existing Dataplatform cluster",
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringMatch(regexp.MustCompile(constant.DataPlatformNameRegexConstraint), constant.DataPlatformRegexNameError)),
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
			"node_pools": {
				Type:        schema.TypeList,
				Description: "list of node pools",
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:        schema.TypeString,
							Description: "The name of your node pool.",
							Computed:    true,
						},
						"version": {
							Type:        schema.TypeString,
							Description: "The version of the Data Platform.",
							Computed:    true,
						},
						"datacenter_id": {
							Type:        schema.TypeString,
							Description: "The UUID of the virtual data center (VDC) in which the node pool is provisioned",
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
							Description: "Key-value pairs attached to the node pool resource as kubernetes labels",
							Computed:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
						},
						"annotations": {
							Type:        schema.TypeMap,
							Description: "Key-value pairs attached to node pool resource as kubernetes annotations",
							Computed:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},
		},
		Timeouts: &resourceDefaultTimeouts,
	}
}

func dataSourceNodePoolsRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(services.SdkBundle).DataplatformClient

	clusterId := d.Get("cluster_id").(string)
	nameValue, nameOk := d.GetOk("name")
	name := nameValue.(string)

	var results []dataplatform.NodePoolResponseData
	var err diag.Diagnostics

	if nameOk {
		results, err = filterNodePools(ctx, d, client, name)
		if err != nil {
			return err
		}
	} else {
		nodePools, _, err := client.ListNodePools(ctx, clusterId)
		if err != nil {
			diags := diag.FromErr(fmt.Errorf("an error occurred while fetching Dataplatform NodePools: %w", err))
			return diags
		}
		if nodePools.Items != nil {
			results = *nodePools.Items
		}
	}

	if results == nil || len(results) == 0 {
		err := fmt.Errorf("no Dataplatform NodePool found under cluster %s", clusterId)
		if nameOk {
			err = fmt.Errorf("%w with the specified name = %s", err, name)
		}
		return diag.FromErr(err)
	}

	if err = dataplatformService.SetNodePoolsData(d, results); err != nil {
		return err
	}

	return nil
}
