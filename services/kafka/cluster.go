package kafka

import (
	"context"
	"fmt"
	"github.com/ionos-cloud/sdk-go-bundle/shared"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	kafka "github.com/ionos-cloud/sdk-go-kafka"

	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils/constant"
)

// todo replace when sdk-bundle is used
func overrideClientFromLoadedConfig(client *Client, productName, location string) {
	loadedConfig := client.GetLoadedConfig()
	if loadedConfig == nil {
		return
	}
	config := client.GetConfig()
	if config == nil {
		return
	}
	endpoint := loadedConfig.GetProductLocationOverrides(productName, location)
	if endpoint == nil {
		log.Printf("[WARN] Missing endpoint for %s in location %s", productName, location)
		return
	}
	config.Servers = kafka.ServerConfigurations{
		{
			URL:         endpoint.Name,
			Description: shared.EndpointOverridden + location,
		},
	}
	if endpoint.SkipTLSVerify {
		config.HTTPClient.Transport = utils.CreateTransport(true)
	}
}

// CreateCluster creates a new Kafka Cluster
func (c *Client) CreateCluster(ctx context.Context, d *schema.ResourceData) (
	kafka.ClusterRead, utils.ApiResponseInfo,
	error,
) {
	location := d.Get("location").(string)
	overrideClientFromLoadedConfig(c, shared.Kafka, location)
	c.changeConfigURL(location)

	request := setClusterPostRequest(d)
	cluster, apiResponse, err := c.sdkClient.ClustersApi.ClustersPost(ctx).ClusterCreate(*request).Execute()
	apiResponse.LogInfo()
	return cluster, apiResponse, err
}

// IsClusterAvailable checks if the Kafka Cluster is available
func (c *Client) IsClusterAvailable(ctx context.Context, d *schema.ResourceData) (bool, error) {
	clusterID := d.Id()
	cluster, _, err := c.GetClusterByID(ctx, clusterID, d.Get("location").(string))
	if err != nil {
		return false, err
	}
	if cluster.Metadata == nil || cluster.Metadata.State == nil {
		return false, fmt.Errorf("expected metadata, got empty for Cluster with ID: %s", clusterID)
	}
	log.Printf("[DEBUG] Cluster status: %s", *cluster.Metadata.State)
	return strings.EqualFold(*cluster.Metadata.State, constant.Available), nil
}

// DeleteCluster deletes a Kafka Cluster
func (c *Client) DeleteCluster(ctx context.Context, id string, location string) (utils.ApiResponseInfo, error) {
	overrideClientFromLoadedConfig(c, shared.Kafka, location)
	c.changeConfigURL(location)

	apiResponse, err := c.sdkClient.ClustersApi.ClustersDelete(ctx, id).Execute()
	apiResponse.LogInfo()
	return apiResponse, err
}

// IsClusterDeleted checks if the Kafka Cluster has been deleted
func (c *Client) IsClusterDeleted(ctx context.Context, d *schema.ResourceData) (bool, error) {
	location := d.Get("location").(string)
	overrideClientFromLoadedConfig(c, shared.Kafka, location)
	c.changeConfigURL(location)

	_, apiResponse, err := c.sdkClient.ClustersApi.ClustersFindById(ctx, d.Id()).Execute()
	apiResponse.LogInfo()
	return apiResponse.HttpNotFound(), err
}

// GetClusterByID retrieves a Kafka Cluster by its ID
func (c *Client) GetClusterByID(ctx context.Context, id string, location string) (
	kafka.ClusterRead, utils.ApiResponseInfo, error,
) {
	overrideClientFromLoadedConfig(c, shared.Kafka, location)
	c.changeConfigURL(location)

	Cluster, apiResponse, err := c.sdkClient.ClustersApi.ClustersFindById(ctx, id).Execute()
	apiResponse.LogInfo()
	return Cluster, apiResponse, err
}

// ListClusters retrieves a list of Kafka Clusters
func (c *Client) ListClusters(ctx context.Context, location string) (kafka.ClusterReadList, *kafka.APIResponse, error) {
	overrideClientFromLoadedConfig(c, shared.Kafka, location)
	c.changeConfigURL(location)

	Clusters, apiResponse, err := c.sdkClient.ClustersApi.ClustersGet(ctx).Execute()
	apiResponse.LogInfo()
	return Clusters, apiResponse, err
}

func setClusterPostRequest(d *schema.ResourceData) *kafka.ClusterCreate {
	clusterName := d.Get("name").(string)
	version := d.Get("version").(string)
	size := d.Get("size").(string)
	datacenterID := d.Get("connections.0.datacenter_id").(string)
	lanID := d.Get("connections.0.lan_id").(string)
	brokerAddressesRaw := d.Get("connections.0.broker_addresses").([]interface{})

	brokerAddresses := make([]string, 0)
	for _, v := range brokerAddressesRaw {
		brokerAddresses = append(brokerAddresses, v.(string))
	}

	connection := kafka.KafkaClusterConnection{
		DatacenterId:    &datacenterID,
		LanId:           &lanID,
		BrokerAddresses: &brokerAddresses,
	}

	return kafka.NewClusterCreate(*kafka.NewCluster(clusterName, version, size, []kafka.KafkaClusterConnection{connection}))
}

// SetKafkaClusterData sets the Kafka Cluster data to the Terraform Resource Data
func (c *Client) SetKafkaClusterData(d *schema.ResourceData, cluster *kafka.ClusterRead) error {
	if cluster.Id != nil {
		d.SetId(*cluster.Id)
	}

	if cluster.Properties == nil || cluster.Metadata == nil {
		return fmt.Errorf("expected properties and metadata in the response for the Kafka Cluster with ID: %s, but received 'nil' instead", *cluster.Id)
	}

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

	if cluster.Properties.Connections != nil && len(*cluster.Properties.Connections) > 0 {
		var connections []interface{}

		for _, connection := range *cluster.Properties.Connections {
			connectionEntry := c.setConnectionProperties(connection)
			connections = append(connections, connectionEntry)
		}

		if err := d.Set("connections", connections); err != nil {
			return err
		}
	}

	if cluster.Metadata.BrokerAddresses != nil {
		if err := d.Set("broker_addresses", *cluster.Metadata.BrokerAddresses); err != nil {
			return err
		}
	}

	return nil
}

func (c *Client) setConnectionProperties(connection kafka.KafkaClusterConnection) map[string]interface{} {
	connectionMap := map[string]interface{}{}

	utils.SetPropWithNilCheck(connectionMap, "datacenter_id", connection.DatacenterId)
	utils.SetPropWithNilCheck(connectionMap, "lan_id", connection.LanId)
	utils.SetPropWithNilCheck(connectionMap, "broker_addresses", connection.BrokerAddresses)

	return connectionMap
}
