package ionoscloud

import (
	"context"
	"fmt"

	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils/constant"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	kafka "github.com/ionos-cloud/sdk-go-kafka"
)

func dataSourceKafkaCluster() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceKafkaClusterRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:             schema.TypeString,
				Description:      "The ID of the Kafka Cluster",
				ValidateDiagFunc: validation.ToDiagFunc(validation.IsUUID),
				Optional:         true,
				Computed:         true,
			},
			"name": {
				Type:        schema.TypeString,
				Description: "The name of the Kafka Cluster",
				Computed:    true,
				Optional:    true,
			},
			"location": {
				Type:             schema.TypeString,
				Description:      "The location of your Kafka Cluster. Supported locations: de/fra, de/txl, es/vit, gb/lhr, us/ewr, us/las, us/mci, fr/par",
				Required:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice(constant.MariaDBClusterLocations, false)),
			},
			"version": {
				Type:        schema.TypeString,
				Description: "The version of the Kafka Cluster",
				Computed:    true,
			},
			"size": {
				Type:        schema.TypeString,
				Description: "The size of the Kafka Cluster",
				Computed:    true,
			},
			"connections": {
				Type:        schema.TypeList,
				Description: "The network connection for your Kafka Cluster. Only one connection is allowed.",
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"datacenter_id": {
							Type:        schema.TypeString,
							Description: "The datacenter to connect your Kafka Cluster to.",
							Computed:    true,
						},
						"lan_id": {
							Type:        schema.TypeString,
							Description: "The numeric LAN ID to connect your Kafka Cluster to.",
							Computed:    true,
						},
						"cidr": {
							Type:        schema.TypeString,
							Description: "The IP and subnet for your Kafka Cluster.",
							Computed:    true,
						},
						"broker_addresses": {
							Type:        schema.TypeList,
							Description: "The broker addresses of the Kafka Cluster",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Computed: true,
						},
					},
				},
			},
			"bootstrap_address": {
				Type:        schema.TypeString,
				Description: "The bootstrap IP address and port.",
				Computed:    true,
			},
			"broker_addresses": {
				Type:        schema.TypeList,
				Description: "IP address and port of cluster brokers.",
				Computed:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"partial_match": {
				Type:        schema.TypeBool,
				Description: "Whether partial matching is allowed or not when using the name filter.",
				Default:     false,
				Optional:    true,
			},
		},
	}
}

func dataSourceKafkaClusterRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(services.SdkBundle).KafkaClient
	idValue, idOk := d.GetOk("id")
	nameValue, nameOk := d.GetOk("name")
	id := idValue.(string)
	name := nameValue.(string)
	location := d.Get("location").(string)

	if idOk && nameOk {
		return diag.FromErr(fmt.Errorf("ID and name cannot be both specified at the same time"))
	}
	if !idOk && !nameOk {
		return diag.FromErr(fmt.Errorf("please provide either the Kafka Cluster ID or name"))
	}

	partialMatch := d.Get("partial_match").(bool)
	var cluster kafka.ClusterRead
	var err error
	if idOk {
		cluster, _, err = client.GetClusterById(ctx, id, location)
		if err != nil {
			return diag.FromErr(fmt.Errorf("an error occured while fetching the Kafka Cluster with ID: %s, error: %w", id, err))
		}
	} else {
		var results []kafka.ClusterRead

		clusters, _, err := client.ListClusters(ctx, location)
		if err != nil {
			return diag.FromErr(fmt.Errorf("an error occured while fetching Kafka Cluster: %w", err))
		}

		for _, cluster := range *clusters.Items {
			if cluster.Properties != nil && cluster.Properties.Name != nil && utils.NameMatches(*cluster.Properties.Name, name, partialMatch) {
				results = append(results, cluster)
			}
		}

		if results == nil || len(results) == 0 {
			return diag.FromErr(fmt.Errorf("no Kafka Clusters found with the specified name: %s", name))
		} else if len(results) > 1 {
			return diag.FromErr(fmt.Errorf("more than one Kafka Cluster found with the specified name: %s", name))
		} else {
			cluster = results[0]
		}
	}

	if err := client.SetKafkaClusterData(d, &cluster); err != nil {
		return diag.FromErr(err)
	}

	return nil
}
