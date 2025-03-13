//go:build all || kafka
// +build all kafka

package ionoscloud

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"testing"

	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/bundleclient"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils/constant"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestKafkaTopicResource(t *testing.T) {
	resource.Test(
		t, resource.TestCase{
			PreCheck: func() {
				testAccPreCheck(t)
			},
			CheckDestroy:             testCheckKafkaTopicDestroy,
			ProtoV6ProviderFactories: testAccProtoV6ProviderFactoriesInternal(t, &testAccProvider),
			Steps: []resource.TestStep{
				{
					Config: configKafkaTopicBasic(topicResourceName, topicAttributeNameValue),
					Check:  checkKafkaTopicResource(topicAttributeNameValue),
				},
				{
					Config: configKafkaTopicDataSourceGetByID(topicResourceName, topicDataName, constant.KafkaClusterTopicResource+"."+topicResourceName+".id"),
					Check: checkKafkaTopicResourceAttributesComparative(
						constant.DataSource+"."+constant.KafkaClusterTopicResource+"."+topicDataName, constant.KafkaClusterTopicResource+"."+topicResourceName,
					),
				},
				{
					Config: configKafkaTopicDataSourceGetByName(
						topicResourceName, topicDataName, constant.KafkaClusterTopicResource+"."+topicResourceName+"."+topicAttributeName, false,
					),
					Check: checkKafkaTopicResourceAttributesComparative(
						constant.DataSource+"."+constant.KafkaClusterTopicResource+"."+topicDataName, constant.KafkaClusterTopicResource+"."+topicResourceName,
					),
				},
				{
					Config: configKafkaTopicDataSourceGetByName(
						topicResourceName, topicDataName, fmt.Sprintf(`"%v"`, topicAttributeNameValue[:len(topicAttributeNameValue)-3]), true,
					),
					Check: checkKafkaTopicResourceAttributesComparative(
						constant.DataSource+"."+constant.KafkaClusterTopicResource+"."+topicDataName, constant.KafkaClusterTopicResource+"."+topicResourceName,
					),
				},
				{
					Config:      configKafkaTopicDataSourceGetByName(topicResourceName, topicDataName, `"willnotwork"`, false),
					ExpectError: regexp.MustCompile(`no Kafka Cluster Topic found with the specified name`),
				},
			},
		},
	)
}

func testCheckKafkaTopicDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(bundleclient.SdkBundle).KafkaClient

	ctx, cancel := context.WithTimeout(context.Background(), *resourceDefaultTimeouts.Default)

	if cancel != nil {
		defer cancel()
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != constant.KafkaClusterTopicResource {
			continue
		}

		_, apiResponse, err := client.GetTopicByID(ctx, rs.Primary.Attributes["cluster_id"], rs.Primary.ID, rs.Primary.Attributes["location"])
		apiResponse.LogInfo()
		if err != nil {
			if !apiResponse.HttpNotFound() {
				return fmt.Errorf("an error occurred while checking the destruction of resource %s: %w", rs.Primary.ID, err)
			}
		} else {
			return fmt.Errorf("resource %s in %s still exists", rs.Primary.ID, rs.Primary.Attributes["location"])
		}
	}
	return nil
}

func testAccCheckKafkaTopicExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(bundleclient.SdkBundle).KafkaClient
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no ID is set")
		}

		ctx, cancel := context.WithTimeout(context.Background(), *resourceDefaultTimeouts.Default)
		defer cancel()

		foundCluster, _, err := client.GetTopicByID(ctx, rs.Primary.Attributes["cluster_id"], rs.Primary.ID, rs.Primary.Attributes["location"])
		if err != nil {
			return fmt.Errorf("an error occurred while fetching Kafka Cluster Topic with ID: %v, error: %w", rs.Primary.ID, err)
		}
		if foundCluster.Id != rs.Primary.ID {
			return fmt.Errorf("resource not found")
		}

		return nil
	}
}

func checkKafkaTopicResource(attributeNameReferenceValue string) resource.TestCheckFunc {
	return resource.ComposeTestCheckFunc(
		testAccCheckKafkaTopicExists(constant.KafkaClusterTopicResource+"."+topicResourceName),
		checkKafkaTopicResourceAttributes(constant.KafkaClusterTopicResource+"."+topicResourceName, attributeNameReferenceValue),
	)
}

func checkKafkaTopicResourceAttributes(fullResourceName, attributeNameReferenceValue string) resource.TestCheckFunc {
	return resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttr(fullResourceName, topicAttributeName, attributeNameReferenceValue),
		resource.TestCheckResourceAttr(fullResourceName, topicAttributeReplicationFactor, topicAttributeReplicationFactorValue),
		resource.TestCheckResourceAttr(fullResourceName, topicAttributeNumberOfPartitions, topicAttributeNumberOfPartitionsValue),
		resource.TestCheckResourceAttr(fullResourceName, topicAttributeSegmentBytes, topicAttributeSegmentBytesValue),
		resource.TestCheckResourceAttr(fullResourceName, topicAttributeRetentionTime, topicAttributeRetentionTimeValue),
	)
}

func checkKafkaTopicResourceAttributesComparative(fullResourceName, fullReferenceResourceName string) resource.TestCheckFunc {
	return resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttrPair(fullResourceName, topicAttributeName, fullReferenceResourceName, topicAttributeName),
		resource.TestCheckResourceAttrPair(fullResourceName, topicAttributeReplicationFactor, fullReferenceResourceName, topicAttributeReplicationFactor),
		resource.TestCheckResourceAttrPair(fullResourceName, topicAttributeNumberOfPartitions, fullReferenceResourceName, topicAttributeNumberOfPartitions),
		resource.TestCheckResourceAttrPair(fullResourceName, topicAttributeSegmentBytes, fullReferenceResourceName, topicAttributeSegmentBytes),
		resource.TestCheckResourceAttrPair(fullResourceName, topicAttributeRetentionTime, fullReferenceResourceName, topicAttributeRetentionTime),
	)
}

func configKafkaTopicBasic(resourceName, attributeName string) string {
	clusterBasicConfig := fmt.Sprintf(
		templateKafkaTopicConfig, topicResourceName, topicAttributeNameValue, topicAttributeReplicationFactorValue, topicAttributeNumberOfPartitionsValue,
		topicAttributeRetentionTimeValue, topicAttributeSegmentBytesValue,
	)

	return strings.Join([]string{defaultKafkaTopicBaseConfig, clusterBasicConfig}, "\n")
}

func configKafkaTopicDataSourceGetByID(resourceName, dataSourceName, dataSourceAttributeID string) string {
	dataSourceBasicConfig := fmt.Sprintf(
		templateKafkaTopicDataIDConfig, dataSourceName, dataSourceAttributeID,
	)
	baseConfig := configKafkaTopicBasic(resourceName, topicAttributeNameValue)

	return strings.Join([]string{baseConfig, dataSourceBasicConfig}, "\n")
}

func configKafkaTopicDataSourceGetByName(resourceName, dataSourceName, dataSourceAttributeName string, partialMatching bool) string {
	dataSourceBasicConfig := fmt.Sprintf(
		templateKafkaTopicDataNameConfig, dataSourceName, dataSourceAttributeName, partialMatching,
	)
	baseConfig := configKafkaTopicBasic(resourceName, topicAttributeNameValue)

	return strings.Join([]string{baseConfig, dataSourceBasicConfig}, "\n")
}

const (
	topicResourceName = "test_kafka_topic"
	topicDataName     = "test_ds_kafka_topic"

	topicAttributeName      = "name"
	topicAttributeNameValue = "test_kafka_topic"

	topicAttributeReplicationFactor      = "replication_factor"
	topicAttributeReplicationFactorValue = "3"

	topicAttributeNumberOfPartitions      = "number_of_partitions"
	topicAttributeNumberOfPartitionsValue = "3"

	topicAttributeRetentionTime      = "retention_time"
	topicAttributeRetentionTimeValue = "86400000"

	topicAttributeSegmentBytes      = "segment_bytes"
	topicAttributeSegmentBytesValue = "1073741824"
)

const templateKafkaTopicConfig = `
resource "ionoscloud_kafka_cluster_topic" "%v" {
	name = "%v"
	cluster_id = ionoscloud_kafka_cluster.test_kafka_cluster.id
	location = ionoscloud_kafka_cluster.test_kafka_cluster.location
	replication_factor = %v
	number_of_partitions = %v
	retention_time = %v
	segment_bytes = %v
}`

const templateKafkaTopicDataIDConfig = `
data "ionoscloud_kafka_cluster_topic" "%v" {
  	id = %v
	cluster_id = ionoscloud_kafka_cluster.test_kafka_cluster.id
	location = ionoscloud_kafka_cluster.test_kafka_cluster.location
}`

const templateKafkaTopicDataNameConfig = `
data "ionoscloud_kafka_cluster_topic" "%v" {
	name = %v
	cluster_id = ionoscloud_kafka_cluster.test_kafka_cluster.id
	location = ionoscloud_kafka_cluster.test_kafka_cluster.location
	partial_match = %v
}`

const defaultKafkaTopicBaseConfig = `
resource "ionoscloud_datacenter" "test_datacenter" {
	  name = "test_datacenter"
	  location = "de/fra"
	  sec_auth_protection = false
}

resource "ionoscloud_lan" "test_lan" {
	  name = "test_lan"
	  public = false
	  datacenter_id = ionoscloud_datacenter.test_datacenter.id
}

resource "ionoscloud_kafka_cluster" "test_kafka_cluster" {
	name = "test_kafka_cluster"
	version = "3.7.0"
	size = "S"
	location = ionoscloud_datacenter.test_datacenter.location
	connections {
		datacenter_id = ionoscloud_datacenter.test_datacenter.id
		lan_id = ionoscloud_lan.test_lan.id
		broker_addresses = [
			"192.168.1.1/24",
			"192.168.1.2/24",
			"192.168.1.3/24"
		]
	}
}`
