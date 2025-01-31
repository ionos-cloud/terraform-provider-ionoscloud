package ionoscloud

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	kafkaSdk "github.com/ionos-cloud/sdk-go-bundle/products/kafka/v2"

	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/kafka"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils"
)

func dataSourceKafkaTopic() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceKafkaTopicRead,
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
				Description: "The name of your Kafka Cluster Topic. Must be 63 characters or less and must begin and end with an alphanumeric character (`[a-z0-9A-Z]`) with dashes (`-`), underscores (`_`), dots (`.`), and alphanumerics between.",
				Optional:    true,
				Computed:    true,
			},
			"location": {
				Type:        schema.TypeString,
				Description: fmt.Sprintf("The location of your Kafka Cluster Topic. Supported locations: %s", strings.Join(kafka.AvailableLocations, ", ")),
				Required:    true,
				ForceNew:    true,
			},
			"cluster_id": {
				Type:        schema.TypeString,
				Description: "The ID of the Kafka Cluster that the topic belongs to",
				Required:    true,
			},
			"replication_factor": {
				Type:        schema.TypeInt,
				Description: "The number of replicas of the topic. The replication factor determines how many copies of the topic are stored on different brokers. The replication factor must be less than or equal to the number of brokers in the Kafka Cluster.",
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
	clusterID := d.Get("cluster_id").(string)
	location := d.Get("location").(string)

	id := idValue.(string)
	name := nameValue.(string)

	if idOk && nameOk {
		return diag.FromErr(fmt.Errorf("ID and name cannot be both specified at the same time"))
	}
	if !idOk && !nameOk {
		return diag.FromErr(fmt.Errorf("please provide either the Kafka Cluster Topic ID or name"))
	}

	partialMatch := d.Get("partial_match").(bool)

	var topic kafkaSdk.TopicRead
	var err error
	if idOk {
		topic, _, err = client.GetTopicByID(ctx, clusterID, id, location)
		if err != nil {
			return diag.FromErr(fmt.Errorf("an error occurred while fetching the Kafka Cluster Topic with ID: %s, error: %w", id, err))
		}
	} else {
		var results []kafkaSdk.TopicRead

		topics, _, err := client.ListTopics(ctx, clusterID, location)
		if err != nil {
			return diag.FromErr(fmt.Errorf("an error occurred while fetching Kafka Cluster Topics: %w", err))
		}

		for _, t := range *topics.Items {
			if t.Properties != nil && t.Properties.Name != nil && utils.NameMatches(*t.Properties.Name, name, partialMatch) {
				results = append(results, t)
			}
		}

		switch {
		case len(results) == 0:
			return diag.FromErr(fmt.Errorf("no Kafka Cluster Topic found with the specified name: %s", name))
		case len(results) > 1:
			return diag.FromErr(fmt.Errorf("more than one Kafka Cluster Topic found with the specified name: %s", name))
		default:
			topic = results[0]
		}
	}

	if err = client.SetKafkaTopicData(d, &topic); err != nil {
		return diag.FromErr(err)
	}

	return nil
}
