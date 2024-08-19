package ionoscloud

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/kafka"
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
				Description:      fmt.Sprintf("The location of your Kafka Cluster Topic. Supported locations: %s", strings.Join(kafka.AvailableLocations, ", ")),
				Required:         true,
				ForceNew:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice(kafka.AvailableLocations, false)),
			},
			"replication_factor": {
				Type:        schema.TypeInt,
				Description: "The number of replicas of the topic. The replication factor determines how many copies of the topic are stored on different brokers. The replication factor must be less than or equal to the number of brokers in the Kafka Cluster.",
				Optional:    true,
				Default:     3,
				ForceNew:    true,
			},
			"number_of_partitions": {
				Type:        schema.TypeInt,
				Description: "The number of partitions of the topic. Partitions allow for parallel processing of messages. The partition count must be greater than or equal to the replication factor.",
				Optional:    true,
				Default:     3,
				ForceNew:    true,
			},
			"retention_time": {
				Type:        schema.TypeInt,
				Description: "This configuration controls the maximum time we will retain a log before we will discard old log segments to free up space. This represents an SLA on how soon consumers must read their data. If set to -1, no time limit is applied.",
				Optional:    true,
				Default:     604800000,
				ForceNew:    true,
			},
			"segment_bytes": {
				Type:        schema.TypeInt,
				Description: "This configuration controls the segment file size for the log. Retention and cleaning is always done a file at a time so a larger segment size means fewer files but less granular control over retention.",
				Optional:    true,
				Default:     1073741824,
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

	// Sleep for 5 second to avoid 500 error from the API
	time.Sleep(5 * time.Second)

	err = utils.WaitForResourceToBeReady(ctx, d, client.IsTopicAvailable)
	if err != nil {
		diags := diag.FromErr(fmt.Errorf("an error occurred  while Kafka Cluster Topic waiting to be ready: %w", err))
		return diags
	}

	return resourceKafkaTopicRead(ctx, d, meta)
}

func resourceKafkaTopicRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(services.SdkBundle).KafkaClient
	topicID := d.Id()
	clusterID := d.Get("cluster_id").(string)
	location := d.Get("location").(string)

	topic, apiResponse, err := client.GetTopicByID(ctx, clusterID, topicID, location)
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
	topicID := d.Id()
	clusterID := d.Get("cluster_id").(string)
	location := d.Get("location").(string)

	apiResponse, err := client.DeleteTopic(ctx, clusterID, topicID, location)
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
	if len(parts) != 3 {
		return nil, fmt.Errorf("expected ID in the format location:cluster_id:topic_id")
	}

	if err := d.Set("location", parts[0]); err != nil {
		return nil, fmt.Errorf("failed to set location for Kafka Cluster Topic import: %w", err)
	}
	if err := d.Set("cluster_id", parts[1]); err != nil {
		return nil, fmt.Errorf("failed to set cluster_id for Kafka Cluster Topic import: %w", err)
	}
	d.SetId(parts[2])

	diags := resourceKafkaTopicRead(ctx, d, meta)
	if diags != nil && diags.HasError() {
		return nil, fmt.Errorf(diags[0].Summary)
	}

	return []*schema.ResourceData{d}, nil
}
