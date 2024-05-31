package kafka

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	kafka "github.com/ionos-cloud/sdk-go-kafka"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils/constant"

	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils"
)

var ClusterResourceName = "Kafka"

func (c *Client) CreateCluster(ctx context.Context, d *schema.ResourceData) (kafka.KafkaPostPutResponseData, utils.ApiResponseInfo, error) {
	request := setClusterPostRequest(d)
	cluster, apiResponse, err := c.sdkClient.KafkaApi.ClustersPost(ctx).KafkaPostRequestBodyData(*request).Execute()
	apiResponse.LogInfo()
	return cluster, apiResponse, err
}

func (c *Client) IsClusterAvailable(ctx context.Context, d *schema.ResourceData) (bool, error) {
	ClusterId := d.Id()
	Cluster, _, err := c.GetClusterById(ctx, ClusterId)
	if err != nil {
		return false, err
	}
	if Cluster.Metadata == nil || Cluster.Metadata.State == nil {
		return false, fmt.Errorf("expected metadata, got empty for Cluster with ID: %s", ClusterId)
	}
	log.Printf("[DEBUG] Cluster status: %s", *Cluster.Metadata.State)
	return strings.EqualFold(*Cluster.Metadata.State, constant.Available), nil
}

func (c *Client) UpdateCluster(ctx context.Context, id string, d *schema.ResourceData) (kafka.KafkaPostPutResponseData, utils.ApiResponseInfo, error) {
	request := setClusterPatchRequest(d)
	ClusterResponse, apiResponse, err := c.sdkClient.KafkaApi.ClustersPut(ctx, id).KafkaPutRequestBodyData(*request).Execute()
	apiResponse.LogInfo()
	return ClusterResponse, apiResponse, err
}

func (c *Client) DeleteCluster(ctx context.Context, id string) (utils.ApiResponseInfo, error) {
	apiResponse, err := c.sdkClient.KafkaApi.ClustersDelete(ctx, id).Execute()
	apiResponse.LogInfo()
	return apiResponse, err
}

func (c *Client) IsClusterDeleted(ctx context.Context, d *schema.ResourceData) (bool, error) {
	_, apiResponse, err := c.sdkClient.KafkaApi.ClustersFindById(ctx, d.Id()).Execute()
	apiResponse.LogInfo()
	return apiResponse.HttpNotFound(), err
}

func (c *Client) GetClusterById(ctx context.Context, id string) (kafka.KafkaClusterByIdResponseData, utils.ApiResponseInfo, error) {
	Cluster, apiResponse, err := c.sdkClient.KafkaApi.ClustersFindById(ctx, id).Execute()
	apiResponse.LogInfo()
	return Cluster, apiResponse, err
}

func (c *Client) ListClusters(ctx context.Context) (kafka.KafkaListResponseData, *kafka.APIResponse, error) {
	Clusters, apiResponse, err := c.sdkClient.KafkaApi.ClustersGet(ctx).Execute()
	apiResponse.LogInfo()
	return Clusters, apiResponse, err
}

func setClusterPostRequest(d *schema.ResourceData) *kafka.KafkaPostRequestBodyData {
	clusterName := d.Get("name").(string)
	version := d.Get("version").(string)
	size := d.Get("size").(string)
	return kafka.NewKafkaPostRequestBodyData(*kafka.NewKafkaCluster(clusterName, version, size))
}

func setClusterPatchRequest(d *schema.ResourceData) *kafka.KafkaPutRequestBodyData {
	name := d.Get("name").(string)
	version := d.Get("version").(string)
	size := d.Get("size").(string)
	request := kafka.NewKafkaPutRequestBodyData(*kafka.NewKafkaCluster(name, version, size))

	return request
}

func (c *Client) SetKafkaClusterData(d *schema.ResourceData, cluster *kafka.KafkaClusterByIdResponseData) error {

	if cluster.Id != nil {
		d.SetId(*cluster.Id)
	}

	if cluster.Properties != nil {
		if cluster.Properties.Name != nil {
			if err := d.Set("name", *cluster.Properties.Name); err != nil {
				return err
			}
		}
		if cluster.Properties.Version != nil {
			if err := d.Set("version", *cluster.Properties.Version); err != nil {
				return err
			}
		}
		if cluster.Properties.Size != nil {
			if err := d.Set("size", *cluster.Properties.Size); err != nil {
				return err
			}
		}
	}

	return nil
}

func (c *Client) SetKafkaClusterMetadata(d *schema.ResourceData, cluster *kafka.KafkaClusterByIdResponseData) error {

	if cluster.Metadata == nil {
		return fmt.Errorf("expected metadata in the response for the kafka Cluster with ID: %s, but received 'nil' instead", *cluster.Id)
	}
	if cluster.Metadata.ServerAddress != nil {
		if err := d.Set("server_address", *cluster.Metadata.ServerAddress); err != nil {
			return fmt.Errorf("error setting server_address: %w", err)
		}
	}
	if cluster.Metadata.BrokerAddresses != nil {
		if err := d.Set("broker_addresses", *cluster.Metadata.BrokerAddresses); err != nil {
			return fmt.Errorf("error setting broker_addresses: %w", err)
		}

	}

	return nil
}

//
// func (c *Client) SetClusterData(d *schema.ResourceData, Cluster kafka.Cluster) error {
// 	d.SetId(*Cluster.Id)
//
// 	if Cluster.Properties == nil {
// 		return fmt.Errorf("expected properties in the response for the kafka Cluster with ID: %s, but received 'nil' instead", *Cluster.Id)
// 	}
//
// 	if Cluster.Metadata == nil {
// 		return fmt.Errorf("expected metadata in the response for the kafka Cluster with ID: %s, but received 'nil' instead", *Cluster.Id)
// 	}
//
// 	if Cluster.Properties.Name != nil {
// 		if err := d.Set("name", *Cluster.Properties.Name); err != nil {
// 			return utils.GenerateSetError(ClusterResourceName, "name", err)
// 		}
// 	}
//
// 	if Cluster.Properties.Logs != nil {
// 		logs := make([]interface{}, len(*Cluster.Properties.Logs))
// 		for i, logElem := range *Cluster.Properties.Logs {
// 			// Populate the logElem entry.
// 			logEntry := make(map[string]interface{})
// 			logEntry["source"] = *logElem.Source
// 			logEntry["tag"] = *logElem.Tag
// 			logEntry["protocol"] = *logElem.Protocol
// 			logEntry["public"] = *logElem.Public
//
// 			// Logic for destinations
// 			destinations := make([]interface{}, len(*logElem.Destinations))
// 			for i, destination := range *logElem.Destinations {
// 				destinationEntry := make(map[string]interface{})
// 				destinationEntry["type"] = *destination.Type
// 				destinationEntry["retention_in_days"] = *destination.RetentionInDays
// 				destinations[i] = destinationEntry
// 			}
// 			logEntry["destinations"] = destinations
// 			logs[i] = logEntry
// 		}
// 		if err := d.Set("log", logs); err != nil {
// 			return utils.GenerateSetError(ClusterResourceName, "log", err)
// 		}
// 	}
//
// 	return nil
// }
