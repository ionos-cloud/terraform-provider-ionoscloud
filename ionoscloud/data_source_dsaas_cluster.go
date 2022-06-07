package ionoscloud

import (
	"context"
	"errors"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	dsaas "github.com/ionos-cloud/sdk-go-autoscaling"
	dsaasService "github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/dsaas"
	"log"
	"strings"
)

func dataSourceDSaaSCluster() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDSaaSClusterRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:         schema.TypeString,
				Description:  "The id of your cluster.",
				Optional:     true,
				ValidateFunc: validation.All(validation.IsUUID),
			},
			"name": {
				Type:        schema.TypeString,
				Description: "The name of your cluster.",
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
				Description: "The version of the DataPlatform.",
				Computed:    true,
			},
			"datacenter_id": {
				Type:        schema.TypeString,
				Description: "The UUID of the virtual data center (VDC) the cluster is provisioned.",
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
							Description: "Time at which the maintenance should start.",
							Computed:    true,
						},
						"day_of_the_week": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
		Timeouts: &resourceDefaultTimeouts,
	}
}

func dataSourceDSaaSClusterRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(SdkBundle).DSaaSClient

	idValue, idOk := d.GetOk("id")
	nameValue, nameOk := d.GetOk("name")

	id := idValue.(string)
	name := nameValue.(string)

	if idOk && nameOk {
		diags := diag.FromErr(errors.New("id and display_name cannot be both specified in the same time"))
		return diags
	}
	if !idOk && !nameOk {
		diags := diag.FromErr(errors.New("please provide either the dbaas cluster id or display_name"))
		return diags
	}

	var cluster dsaas.ClusterResponseData
	var err error

	if idOk {
		/* search by ID */
		cluster, _, err = client.GetCluster(ctx, id)
		if err != nil {
			diags := diag.FromErr(fmt.Errorf("an error occurred while fetching the DSaaS cluster with ID %s: %s", id, err))
			return diags
		}
	} else {
		var results []dsaas.ClusterResponseData

		partialMatch := d.Get("partial_match").(bool)

		log.Printf("[INFO] Using data source for DSaaS cluster by name with partial_match %t and name: %s", partialMatch, name)

		if partialMatch {
			clusters, _, err := client.ListClusters(ctx, name)
			if err != nil {
				diags := diag.FromErr(fmt.Errorf("an error occurred while fetching DSaaS clusters: %s", err.Error()))
				return diags
			}
			results = clusters
		} else {
			clusters, _, err := client.ListClusters(ctx, "")
			if err != nil {
				diags := diag.FromErr(fmt.Errorf("an error occurred while fetching DSaaS clusters: %s", err.Error()))
				return diags
			}
			if clusters != nil && len(clusters) > 0 {
				for _, clusterItem := range clusters {
					if clusterItem.Properties != nil && clusterItem.Properties.Name != nil && strings.EqualFold(*clusterItem.Properties.Name, name) {
						tmpBackupUnit, _, err := client.GetCluster(ctx, *clusterItem.Id)
						if err != nil {
							return diag.FromErr(fmt.Errorf("an error occurred while fetching the DSaaS cluster with ID: %s while searching by full name: %s: %w", *clusterItem.Id, name, err))
						}
						results = append(results, tmpBackupUnit)
					}
				}
			}
		}

		if results == nil || len(results) == 0 {
			return diag.FromErr(fmt.Errorf("no DSaaS cluster found with the specified name = %s", name))
		} else if len(results) > 1 {
			return diag.FromErr(fmt.Errorf("more than one DSaaS cluster found with the specified criteria name = %s", name))
		} else {
			cluster = results[0]
		}

	}

	if err := dsaasService.SetDSaaSClusterData(d, cluster); err != nil {
		return diag.FromErr(err)
	}

	return nil

}
