package ionoscloud

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils"
)

func resourceKafkaTopic() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceKafkaTopicCreate,
		ReadContext:   resourceKafkaTopicRead,
		DeleteContext: resourceKafkaTopicDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceKafkaTopicImport,
		},

		Schema: map[string]*schema.Schema{
			"cluster_id": {
				Type:        schema.TypeString,
				Description: "The ID of the Kafka cluster to which the topic belongs.",
				Required:    true,
				ForceNew:    true,
			},
			"name": {
				Type:             schema.TypeString,
				Description:      "The name of your Kafka cluster topic. Must be 63 characters or less and must begin and end with an alphanumeric character (`[a-z0-9A-Z]`) with dashes (`-`), underscores (`_`), dots (`.`), and alphanumerics between.",
				Required:         true,
				ForceNew:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringIsNotWhiteSpace),
			},
			"replication_factor": {
				Type:        schema.TypeInt,
				Description: "The number of replicas of the topic. The replication factor determines how many copies of the topic are stored on different brokers. The replication factor must be less than or equal to the number of brokers in the Kafka cluster.",
				Required:    true,
				ForceNew:    true,
			},
			"number_of_partitions": {
				Type:        schema.TypeInt,
				Description: "The number of partitions of the topic. Partitions allow for parallel processing of messages. The partition count must be greater than or equal to the replication factor.",
				Required:    true,
				ForceNew:    true,
			},
			"retention_time": {
				Type:        schema.TypeInt,
				Description: "The time in milliseconds that a message is retained in the topic log. Messages older than the retention time are deleted. If value is `0`, messages are retained indefinitely unless other retention is set.",
				Optional:    true,
				Default:     0,
				ForceNew:    true,
			},
			"segment_bytes": {
				Type:        schema.TypeInt,
				Description: "The maximum size in bytes that the topic log can grow to. When the log reaches this size, the oldest messages are deleted. If value is `0`, messages are retained indefinitely unless other retention is set.",
				Optional:    true,
				Default:     0,
				ForceNew:    true,
			},
		},
		Timeouts: &resourceDefaultTimeouts,
	}
}
func resourceKafkaTopicCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(services.SdkBundle).KafkaClient

	createdTopic, _, err := client.CreateTopic(ctx, d)
	if err != nil {
		d.SetId("")
		diags := diag.FromErr(fmt.Errorf("error creating kafka cluster topic: %w", err))
		return diags
	}

	d.SetId(*createdTopic.Id)
	log.Printf("[INFO] Created kafka cluster topic: %s", d.Id())

	err = utils.WaitForResourceToBeReady(ctx, d, client.IsTopicAvailable)
	if err != nil {
		diags := diag.FromErr(fmt.Errorf("an error occured  while kafka cluster topic waiting to be ready: %w", err))
		return diags
	}

	return resourceKafkaClusterRead(ctx, d, meta)
}

func resourceKafkaTopicRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(services.SdkBundle).KafkaClient
	topicId := d.Id()
	clusterId := d.Get("cluster_id").(string)

	topic, apiResponse, err := client.GetTopicById(ctx, clusterId, topicId)
	if err != nil {
		if apiResponse.HttpNotFound() {
			d.SetId("")
			return nil
		}
		diags := diag.FromErr(fmt.Errorf("error while fetching kafka cluster topic %s: %w", d.Id(), err))
		return diags
	}

	log.Printf("[INFO] Successfully retreived cluster topic %s: %+v", d.Id(), topic)

	if err := client.SetKafkaTopicData(d, &topic); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceKafkaTopicDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(services.SdkBundle).KafkaClient
	topicId := d.Id()
	clusterId := d.Get("cluster_id").(string)

	apiResponse, err := client.DeleteTopic(ctx, clusterId, topicId)
	if err != nil {
		if apiResponse.HttpNotFound() {
			d.SetId("")
			return nil
		}
		diags := diag.FromErr(fmt.Errorf("error while deleting kafka cluster topic %s: %w", d.Id(), err))
		return diags
	}

	err = utils.WaitForResourceToBeDeleted(ctx, d, client.IsTopicDeleted)
	if err != nil {
		return diag.FromErr(fmt.Errorf("the check for cluster topic deletion failed with the following error: %w", err))
	}

	d.SetId("")

	return nil
}

func resourceKafkaTopicImport(ctx context.Context, d *schema.ResourceData, meta interface{}) (
	[]*schema.ResourceData, error,
) {
	resourceKafkaTopicRead(ctx, d, meta)
	return []*schema.ResourceData{d}, nil
}
