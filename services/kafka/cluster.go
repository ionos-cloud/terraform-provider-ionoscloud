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

func (c *Client) CreateCluster(ctx context.Context, d *schema.ResourceData) (
	kafka.ClusterRead, utils.ApiResponseInfo,
	error,
) {
	request := setClusterPostRequest(d)
	cluster, apiResponse, err := c.sdkClient.ClustersApi.ClustersPost(ctx).ClusterCreate(*request).Execute()
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

func (c *Client) DeleteCluster(ctx context.Context, id string) (utils.ApiResponseInfo, error) {
	apiResponse, err := c.sdkClient.ClustersApi.ClustersDelete(ctx, id).Execute()
	apiResponse.LogInfo()
	return apiResponse, err
}

func (c *Client) IsClusterDeleted(ctx context.Context, d *schema.ResourceData) (bool, error) {
	_, apiResponse, err := c.sdkClient.ClustersApi.ClustersFindById(ctx, d.Id()).Execute()
	apiResponse.LogInfo()
	return apiResponse.HttpNotFound(), err
}

func (c *Client) GetClusterById(ctx context.Context, id string) (
	kafka.ClusterRead, utils.ApiResponseInfo, error,
) {
	Cluster, apiResponse, err := c.sdkClient.ClustersApi.ClustersFindById(ctx, id).Execute()
	apiResponse.LogInfo()
	return Cluster, apiResponse, err
}

func (c *Client) ListClusters(ctx context.Context) (kafka.ClusterReadList, *kafka.APIResponse, error) {
	Clusters, apiResponse, err := c.sdkClient.ClustersApi.ClustersGet(ctx).Execute()
	apiResponse.LogInfo()
	return Clusters, apiResponse, err
}

func setClusterPostRequest(d *schema.ResourceData) *kafka.ClusterCreate {
	clusterName := d.Get("name").(string)
	version := d.Get("version").(string)
	size := d.Get("size").(string)
	datacenterId := d.Get("connections.0.datacenter_id").(string)
	lanId := d.Get("connections.0.lan_id").(string)
	cidr := d.Get("connections.0.cidr").(string)
	brokerAddressesRaw := d.Get("connections.0.broker_addresses").([]interface{})

	brokerAddresses := make([]string, 0)
	for _, v := range brokerAddressesRaw {
		brokerAddresses = append(brokerAddresses, v.(string))
	}

	connection := kafka.KafkaClusterConnection{
		DatacenterId:    &datacenterId,
		LanId:           &lanId,
		Cidr:            &cidr,
		BrokerAddresses: &brokerAddresses,
	}

	return kafka.NewClusterCreate(*kafka.NewCluster(clusterName, version, size, []kafka.KafkaClusterConnection{connection}))
}

func (c *Client) SetKafkaClusterData(d *schema.ResourceData, cluster *kafka.ClusterRead) error {

	if cluster.Id != nil {
		d.SetId(*cluster.Id)
	}

	if cluster.Properties == nil {
		return fmt.Errorf("expected properties in the response for the Kafka cluster with ID: %s, but received 'nil' instead", *cluster.Id)
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
			connectionEntry := c.SetConnectionProperties(connection)
			connections = append(connections, connectionEntry)
		}

		if err := d.Set("connections", connections); err != nil {
			return err
		}
	}

	return nil
}

func (c *Client) SetConnectionProperties(connection kafka.KafkaClusterConnection) map[string]interface{} {
	connectionMap := map[string]interface{}{}

	utils.SetPropWithNilCheck(connectionMap, "datacenter_id", connection.DatacenterId)
	utils.SetPropWithNilCheck(connectionMap, "lan_id", connection.LanId)
	utils.SetPropWithNilCheck(connectionMap, "cidr", connection.Cidr)
	utils.SetPropWithNilCheck(connectionMap, "broker_addresses", connection.BrokerAddresses)

	return connectionMap
}
