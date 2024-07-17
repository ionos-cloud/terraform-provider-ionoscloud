package ionoscloud

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils/constant"
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
			"id": {
				Type:        schema.TypeString,
				Description: "The ID of the Kafka Cluster Topic.",
				Computed:    true,
			},
			"cluster_id": {
				Type:        schema.TypeString,
				Description: "The ID of the Kafka Cluster to which the topic belongs.",
				Required:    true,
				ForceNew:    true,
			},
			"name": {
				Type:             schema.TypeString,
				Description:      "The name of your Kafka Cluster Topic. Must be 63 characters or less and must begin and end with an alphanumeric character (`[a-z0-9A-Z]`) with dashes (`-`), underscores (`_`), dots (`.`), and alphanumerics between.",
				Required:         true,
				ForceNew:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringIsNotWhiteSpace),
			},
			"location": {
				Type:             schema.TypeString,
				Description:      "The location of your Kafka Cluster Topic. Supported locations: de/fra, de/txl, es/vit, gb/lhr, us/ewr, us/las, us/mci, fr/par",
				Required:         true,
				ForceNew:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice(constant.MariaDBClusterLocations, false)),
			},
			"replication_factor": {
				Type:        schema.TypeInt,
				Description: "The number of replicas of the topic. The replication factor determines how many copies of the topic are stored on different brokers. The replication factor must be less than or equal to the number of brokers in the Kafka Cluster.",
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
		diags := diag.FromErr(fmt.Errorf("error creating Kafka Cluster Topic: %w", err))
		return diags
	}

	d.SetId(*createdTopic.Id)
	log.Printf("[INFO] Created Kafka Cluster Topic: %s", d.Id())

	err = utils.WaitForResourceToBeReady(ctx, d, client.IsTopicAvailable)
	if err != nil {
		diags := diag.FromErr(fmt.Errorf("an error occured  while Kafka Cluster Topic waiting to be ready: %w", err))
		return diags
	}

	return resourceKafkaTopicRead(ctx, d, meta)
}

func resourceKafkaTopicRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(services.SdkBundle).KafkaClient
	topicId := d.Id()
	clusterId := d.Get("cluster_id").(string)
	location := d.Get("location").(string)

	topic, apiResponse, err := client.GetTopicById(ctx, clusterId, topicId, location)
	if err != nil {
		if apiResponse.HttpNotFound() {
			d.SetId("")
			return nil
		}
		diags := diag.FromErr(fmt.Errorf("error while fetching Kafka Cluster Topic %s: %w", d.Id(), err))
		return diags
	}

	log.Printf("[INFO] Successfully retreived Kafka Cluster Topic %s: %+v", d.Id(), topic)

	if err := client.SetKafkaTopicData(d, &topic); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceKafkaTopicDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(services.SdkBundle).KafkaClient
	topicId := d.Id()
	clusterId := d.Get("cluster_id").(string)
	location := d.Get("location").(string)

	apiResponse, err := client.DeleteTopic(ctx, clusterId, topicId, location)
	if err != nil {
		if apiResponse.HttpNotFound() {
			d.SetId("")
			return nil
		}
		diags := diag.FromErr(fmt.Errorf("error while deleting Kafka Cluster Topic %s: %w", d.Id(), err))
		return diags
	}

	err = utils.WaitForResourceToBeDeleted(ctx, d, client.IsTopicDeleted)
	if err != nil {
		return diag.FromErr(fmt.Errorf("the check for Kafka Cluster Topic deletion failed with the following error: %w", err))
	}

	d.SetId("")

	return nil
}

func resourceKafkaTopicImport(ctx context.Context, d *schema.ResourceData, meta interface{}) (
	[]*schema.ResourceData, error,
) {
	parts := strings.Split(d.Id(), ":")
	if len(parts) != 2 {
		return nil, fmt.Errorf("expected ID in the format cluster_id:topic_id")
	}

	_ = d.Set("cluster_id", parts[0])
	d.SetId(parts[1])

	diags := resourceKafkaTopicRead(ctx, d, meta)
	if diags != nil && diags.HasError() {
		return nil, fmt.Errorf(diags[0].Summary)
	}

	return []*schema.ResourceData{d}, nil
}
