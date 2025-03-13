package kafka

import (
	"context"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	kafka "github.com/ionos-cloud/sdk-go-bundle/products/kafka/v2"
	"github.com/ionos-cloud/sdk-go-bundle/shared"
	"github.com/ionos-cloud/sdk-go-bundle/shared/fileconfiguration"

	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils/constant"
)

// todo replace when sdk-bundle is used
func overrideClientFromFileConfig(client *Client, productName, location string) {
	defer client.changeConfigURL(location)
	fileConfig := client.GetFileConfig()
	if fileConfig == nil {
		return
	}
	config := client.GetConfig()
	if config == nil {
		return
	}
	endpoint := fileConfig.GetProductLocationOverrides(productName, location)
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
	config.HTTPClient.Transport = shared.CreateTransport(endpoint.SkipTLSVerify, endpoint.CertificateAuthData)

}

// CreateCluster creates a new Kafka Cluster
func (c *Client) CreateCluster(ctx context.Context, d *schema.ResourceData) (
	kafka.ClusterRead, utils.ApiResponseInfo,
	error,
) {
	location := d.Get("location").(string)
	overrideClientFromFileConfig(c, fileconfiguration.Kafka, location)

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
	log.Printf("[DEBUG] Cluster status: %s", cluster.Metadata.State)
	return strings.EqualFold(cluster.Metadata.State, constant.Available), nil
}

// DeleteCluster deletes a Kafka Cluster
func (c *Client) DeleteCluster(ctx context.Context, id string, location string) (utils.ApiResponseInfo, error) {
	overrideClientFromFileConfig(c, fileconfiguration.Kafka, location)

	apiResponse, err := c.sdkClient.ClustersApi.ClustersDelete(ctx, id).Execute()
	apiResponse.LogInfo()
	return apiResponse, err
}

// IsClusterDeleted checks if the Kafka Cluster has been deleted
func (c *Client) IsClusterDeleted(ctx context.Context, d *schema.ResourceData) (bool, error) {
	location := d.Get("location").(string)
	overrideClientFromFileConfig(c, fileconfiguration.Kafka, location)

	_, apiResponse, err := c.sdkClient.ClustersApi.ClustersFindById(ctx, d.Id()).Execute()
	apiResponse.LogInfo()
	return apiResponse.HttpNotFound(), err
}

// GetClusterByID retrieves a Kafka Cluster by its ID
func (c *Client) GetClusterByID(ctx context.Context, id string, location string) (
	kafka.ClusterRead, utils.ApiResponseInfo, error,
) {
	overrideClientFromFileConfig(c, fileconfiguration.Kafka, location)

	Cluster, apiResponse, err := c.sdkClient.ClustersApi.ClustersFindById(ctx, id).Execute()
	apiResponse.LogInfo()
	return Cluster, apiResponse, err
}

// ListClusters retrieves a list of Kafka Clusters
func (c *Client) ListClusters(ctx context.Context, location string) (kafka.ClusterReadList, *kafka.APIResponse, error) {
	overrideClientFromFileConfig(c, fileconfiguration.Kafka, location)

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
		DatacenterId:    datacenterID,
		LanId:           lanID,
		BrokerAddresses: brokerAddresses,
	}

	return kafka.NewClusterCreate(*kafka.NewCluster(clusterName, version, size, []kafka.KafkaClusterConnection{connection}))
}

// SetKafkaClusterData sets the Kafka Cluster data to the Terraform Resource Data
func (c *Client) SetKafkaClusterData(d *schema.ResourceData, cluster *kafka.ClusterRead) error {
	d.SetId(cluster.Id)

	if err := d.Set("name", cluster.Properties.Name); err != nil {
		return err
	}
	if err := d.Set("version", cluster.Properties.Version); err != nil {
		return err
	}
	if err := d.Set("size", cluster.Properties.Size); err != nil {
		return err
	}

	if len(cluster.Properties.Connections) > 0 {
		var connections []interface{}

		for _, connection := range cluster.Properties.Connections {
			connectionEntry := c.setConnectionProperties(connection)
			connections = append(connections, connectionEntry)
		}

		if err := d.Set("connections", connections); err != nil {
			return err
		}
	}

	if err := d.Set("broker_addresses", cluster.Metadata.BrokerAddresses); err != nil {
		return err
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
