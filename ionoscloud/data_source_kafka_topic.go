package ionoscloud

import (
	"context"
	"fmt"

	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	kafka "github.com/ionos-cloud/sdk-go-kafka"
)

func dataSourceKafkaTopic() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceKafkaTopicRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:             schema.TypeString,
				Description:      "The ID of the Kafka cluster",
				ValidateDiagFunc: validation.ToDiagFunc(validation.IsUUID),
				Optional:         true,
				Computed:         true,
			},
			"name": {
				Type:        schema.TypeString,
				Description: "The name of your Kafka cluster topic. Must be 63 characters or less and must begin and end with an alphanumeric character (`[a-z0-9A-Z]`) with dashes (`-`), underscores (`_`), dots (`.`), and alphanumerics between.",
				Optional:    true,
				Computed:    true,
			},
			"cluster_id": {
				Type:        schema.TypeString,
				Description: "The ID of the Kafka cluster that the topic belongs to",
				Required:    true,
			},
			"replication_factor": {
				Type:        schema.TypeInt,
				Description: "The number of replicas of the topic. The replication factor determines how many copies of the topic are stored on different brokers. The replication factor must be less than or equal to the number of brokers in the Kafka cluster.",
				Computed:    true,
			},
			"number_of_partitions": {
				Type:        schema.TypeInt,
				Description: "The number of partitions of the topic. Partitions allow for parallel processing of messages. The partition count must be greater than or equal to the replication factor.",
				Computed:    true,
			},
			"retention_time": {
				Type:        schema.TypeInt,
				Description: "The time in milliseconds that a message is retained in the topic log. Messages older than the retention time are deleted. If value is `0`, messages are retained indefinitely unless other retention is set.",
				Computed:    true,
			},
			"segment_bytes": {
				Type:        schema.TypeInt,
				Description: "The maximum size in bytes that the topic log can grow to. When the log reaches this size, the oldest messages are deleted. If value is `0`, messages are retained indefinitely unless other retention is set.",
				Computed:    true,
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

func dataSourceKafkaTopicRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(services.SdkBundle).KafkaClient
	idValue, idOk := d.GetOk("id")
	nameValue, nameOk := d.GetOk("name")
	clusterId := d.Get("cluster_id").(string)

	id := idValue.(string)
	name := nameValue.(string)

	if idOk && nameOk {
		return diag.FromErr(fmt.Errorf("ID and name cannot be both specified at the same time"))
	}
	if !idOk && !nameOk {
		return diag.FromErr(fmt.Errorf("please provide either the kafka cluster topic ID or name"))
	}

	partialMatch := d.Get("partial_match").(bool)

	var topic kafka.TopicRead
	var err error
	if idOk {
		topic, _, err = client.GetTopicById(ctx, clusterId, id)
		if err != nil {
			return diag.FromErr(fmt.Errorf("an error occured while fetching the kafka cluster topic with ID: %s, error: %w", id, err))
		}
	} else {
		var results []kafka.TopicRead

		clusters, _, err := client.ListTopics(ctx, clusterId)
		if err != nil {
			return diag.FromErr(fmt.Errorf("an error occured while fetching kafka clusters topics: %w", err))
		}

		for _, cluster := range *clusters.Items {
			if cluster.Properties != nil && cluster.Properties.Name != nil && utils.NameMatches(*cluster.Properties.Name, name, partialMatch) {
				results = append(results, cluster)
			}
		}

		if results == nil || len(results) == 0 {
			return diag.FromErr(fmt.Errorf("no kafka cluster topic found with the specified name: %s", name))
		} else if len(results) > 1 {
			return diag.FromErr(fmt.Errorf("more than one kafka cluster topic found with the specified name: %s", name))
		} else {
			topic = results[0]
		}
	}

	if err = client.SetKafkaTopicData(d, &topic); err != nil {
		return diag.FromErr(err)
	}

	return nil
}
