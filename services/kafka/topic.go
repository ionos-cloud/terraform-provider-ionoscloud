package kafka

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	kafka "github.com/ionos-cloud/sdk-go-kafka"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils/constant"
)

func (c *Client) CreateTopic(ctx context.Context, d *schema.ResourceData) (
	kafka.TopicRead, utils.ApiResponseInfo, error,
) {
	topic := setTopicPostRequest(d)
	clusterId := d.Get("cluster_id").(string)

	topicResponse, apiResponse, err := c.sdkClient.TopicsApi.ClustersTopicsPost(ctx, clusterId).TopicCreate(*topic).Execute()
	apiResponse.LogInfo()

	return topicResponse, apiResponse, err
}

func (c *Client) GetTopicById(ctx context.Context, clusterId string, topicId string) (
	kafka.TopicRead, utils.ApiResponseInfo, error,
) {
	topic, apiResponse, err := c.sdkClient.TopicsApi.ClustersTopicsFindById(ctx, clusterId, topicId).Execute()
	apiResponse.LogInfo()

	return topic, apiResponse, err
}

func (c *Client) ListTopics(ctx context.Context, clusterId string) (
	kafka.TopicReadList, utils.ApiResponseInfo, error,
) {
	topics, apiResponse, err := c.sdkClient.TopicsApi.ClustersTopicsGet(ctx, clusterId).Execute()
	apiResponse.LogInfo()

	return topics, apiResponse, err
}

func (c *Client) DeleteTopic(ctx context.Context, clusterId string, topicId string) (utils.ApiResponseInfo, error) {
	apiResponse, err := c.sdkClient.TopicsApi.ClustersTopicsDelete(ctx, clusterId, topicId).Execute()
	apiResponse.LogInfo()

	return apiResponse, err
}

func (c *Client) IsTopicAvailable(ctx context.Context, d *schema.ResourceData) (bool, error) {
	topicId := d.Id()
	clusterId := d.Get("cluster_id").(string)

	topic, _, err := c.GetTopicById(ctx, clusterId, topicId)
	if err != nil {
		return false, err
	}
	if topic.Metadata == nil || topic.Metadata.State == nil {
		return false, fmt.Errorf("expected metadata, got empty for Topic with ID: %s", clusterId)
	}
	log.Printf("[DEBUG] Topic status: %s", *topic.Metadata.State)
	return strings.EqualFold(*topic.Metadata.State, constant.Available), nil
}

func (c *Client) IsTopicDeleted(ctx context.Context, d *schema.ResourceData) (bool, error) {
	clusterId := d.Get("cluster_id").(string)
	topicId := d.Id()

	_, apiResponse, err := c.sdkClient.TopicsApi.ClustersTopicsFindById(ctx, topicId, clusterId).Execute()
	apiResponse.LogInfo()

	return apiResponse.HttpNotFound(), err
}

func (c *Client) SetKafkaTopicData(d *schema.ResourceData, topic *kafka.TopicRead) error {
	if topic.Id != nil {
		d.SetId(*topic.Id)
	}

	if topic.Properties == nil {
		return fmt.Errorf("expected properties in the response for the Kafka cluster topic with ID: %s, but received 'nil' instead", *topic.Id)
	}

	if topic.Properties.Name != nil {
		if err := d.Set("name", *topic.Properties.Name); err != nil {
			return err
		}
	}
	if topic.Properties.NumberOfPartitions != nil {
		if err := d.Set("number_of_partitions", *topic.Properties.NumberOfPartitions); err != nil {
			return err
		}
	}
	if topic.Properties.ReplicationFactor != nil {
		if err := d.Set("replication_factor", *topic.Properties.ReplicationFactor); err != nil {
			return err
		}
	}

	if topic.Properties.LogRetention != nil {
		if topic.Properties.LogRetention.RetentionTime != nil {
			if err := d.Set("retention_time", *topic.Properties.LogRetention.RetentionTime); err != nil {
				return err
			}
		}

		if topic.Properties.LogRetention.SegmentBytes != nil {
			if err := d.Set("segment_bytes", *topic.Properties.LogRetention.SegmentBytes); err != nil {
				return err
			}
		}
	}

	return nil
}

func setTopicPostRequest(d *schema.ResourceData) *kafka.TopicCreate {
	topicName := d.Get("name").(string)
	replicationFactor := d.Get("replication_factor").(int32)
	partitionCount := d.Get("number_of_partitions").(int32)
	retentionTime := d.Get("retention_time").(int32)
	segmentBytes := d.Get("segment_bytes").(int32)

	return kafka.NewTopicCreate(
		*kafka.NewTopic(
			topicName, replicationFactor, partitionCount, *kafka.NewTopicLogRetention(retentionTime, segmentBytes),
		),
	)
}
