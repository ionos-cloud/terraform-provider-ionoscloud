//go:build all || kafka
// +build all kafka

package ionoscloud

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"testing"

	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils/constant"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestKafkaClusterResource(t *testing.T) {
	resource.Test(
		t, resource.TestCase{
			PreCheck: func() {
				testAccPreCheck(t)
			},
			CheckDestroy:             testCheckKafkaClusterDestroy,
			ProtoV6ProviderFactories: testAccProtoV6ProviderFactoriesInternal(t, &testAccProvider),
			Steps: []resource.TestStep{
				{
					Config: configKafkaClusterBasic(clusterResourceName, clusterAttributeNameValue),
					Check:  checkKafkaClusterResourceAttributes(constant.KafkaClusterResource+"."+clusterResourceName, clusterAttributeNameValue),
				},
				{
					Config: configKafkaClusterDataSourceGetByID(clusterResourceName, clusterDataName, constant.KafkaClusterResource+"."+clusterResourceName+".id"),
					Check: checkKafkaClusterResourceAttributesComparative(
						constant.DataSource+"."+constant.KafkaClusterResource+"."+clusterDataName, constant.KafkaClusterResource+"."+clusterResourceName,
					),
				},
				{
					Config: configKafkaClusterDataSourceGetByName(
						clusterResourceName, clusterDataName, constant.KafkaClusterResource+"."+clusterResourceName+"."+clusterAttributeName, false,
					),
					Check: checkKafkaClusterResourceAttributesComparative(
						constant.DataSource+"."+constant.KafkaClusterResource+"."+clusterDataName, constant.KafkaClusterResource+"."+clusterResourceName,
					),
				},
				{
					Config: configKafkaClusterDataSourceGetByName(
						clusterResourceName, clusterDataName, fmt.Sprintf(`"%v"`, clusterAttributeNameValue[:len(clusterAttributeNameValue)-2]), true,
					),
					Check: checkKafkaClusterResourceAttributesComparative(
						constant.DataSource+"."+constant.KafkaClusterResource+"."+clusterDataName, constant.KafkaClusterResource+"."+clusterResourceName,
					),
				},
				{
					Config:      configKafkaClusterDataSourceGetByName(clusterResourceName, clusterDataName, `"willnotwork"`, false),
					ExpectError: regexp.MustCompile(`no Kafka Clusters found with the specified name:`),
				},
			},
		},
	)
}

func testCheckKafkaClusterDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(services.SdkBundle).KafkaClient

	ctx, cancel := context.WithTimeout(context.Background(), *resourceDefaultTimeouts.Default)

	if cancel != nil {
		defer cancel()
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != constant.KafkaClusterResource {
			continue
		}

		_, apiResponse, err := client.GetClusterByID(ctx, rs.Primary.ID, rs.Primary.Attributes["location"])
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

func testAccCheckKafkaClusterExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(services.SdkBundle).KafkaClient
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no ID is set")
		}

		ctx, cancel := context.WithTimeout(context.Background(), *resourceDefaultTimeouts.Default)
		defer cancel()

		foundCluster, _, err := client.GetClusterByID(ctx, rs.Primary.ID, rs.Primary.Attributes["location"])
		if err != nil {
			return fmt.Errorf("an error occurred while fetching Kafka Cluster with ID: %v, error: %w", rs.Primary.ID, err)
		}
		if foundCluster.Id != rs.Primary.ID {
			return fmt.Errorf("resource not found")
		}

		return nil
	}
}

func checkKafkaClusterResource(attributeNameReferenceValue string) resource.TestCheckFunc {
	return resource.ComposeTestCheckFunc(
		testAccCheckKafkaClusterExists(constant.KafkaClusterResource+"."+clusterResourceName),
		checkKafkaClusterResourceAttributes(constant.KafkaClusterResource+"."+clusterResourceName, attributeNameReferenceValue),
	)
}

func checkKafkaClusterResourceAttributes(fullResourceName, attributeNameReferenceValue string) resource.TestCheckFunc {
	return resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttr(fullResourceName, clusterAttributeName, attributeNameReferenceValue),
		resource.TestCheckResourceAttr(fullResourceName, clusterAttributeVersion, clusterAttributeVersionValue),
		resource.TestCheckResourceAttr(fullResourceName, clusterAttributeSize, clusterAttributeSizeValue),
		resource.TestCheckResourceAttrPair(fullResourceName, "connections.0.datacenter_id", "ionoscloud_datacenter.test_datacenter", "id"),
		resource.TestCheckResourceAttrPair(fullResourceName, "connections.0.lan_id", constant.LanResource+"."+constant.LanTestResource, "id"),
		resource.TestCheckResourceAttr(fullResourceName, clusterAttributeBrokerAddresses, "3"),
		resource.TestCheckResourceAttr(fullResourceName, clusterMetadataBrokerAddresses, "3"),
	)
}

func checkKafkaClusterResourceAttributesComparative(fullResourceName, fullReferenceResourceName string) resource.TestCheckFunc {
	return resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttrPair(fullResourceName, clusterAttributeName, fullReferenceResourceName, clusterAttributeName),
		resource.TestCheckResourceAttrPair(fullResourceName, clusterAttributeVersion, fullReferenceResourceName, clusterAttributeVersion),
		resource.TestCheckResourceAttrPair(fullResourceName, clusterAttributeSize, fullReferenceResourceName, clusterAttributeSize),
		resource.TestCheckResourceAttrPair(fullResourceName, "connections.0.datacenter_id", fullReferenceResourceName, "connections.0.datacenter_id"),
		resource.TestCheckResourceAttrPair(fullResourceName, "connections.0.lan_id", fullReferenceResourceName, "connections.0.lan_id"),
		resource.TestCheckResourceAttrPair(fullResourceName, clusterAttributeBrokerAddresses, fullReferenceResourceName, clusterAttributeBrokerAddresses),
		resource.TestCheckResourceAttrPair(fullResourceName, clusterMetadataBrokerAddresses, fullReferenceResourceName, clusterMetadataBrokerAddresses),
	)
}

func configKafkaClusterBasic(resourceName, attributeName string) string {
	clusterBasicConfig := fmt.Sprintf(
		templateKafkaClusterConfig, clusterResourceName, clusterAttributeNameValue, clusterAttributeVersionValue, clusterAttributeSizeValue, clusterAttributeLocationValue,
		clusterAttributeBrokerAddressesValue,
	)

	return strings.Join([]string{defaultKafkaClusterBaseConfig, clusterBasicConfig}, "\n")
}

func configKafkaClusterDataSourceGetByID(resourceName, dataSourceName, dataSourceAttributeID string) string {
	dataSourceBasicConfig := fmt.Sprintf(
		templateKafkaClusterDataIDConfig, dataSourceName, dataSourceAttributeID, constant.KafkaClusterResource+"."+resourceName+"."+clusterAttributeLocation,
	)
	baseConfig := configKafkaClusterBasic(resourceName, clusterAttributeNameValue)

	return strings.Join([]string{baseConfig, dataSourceBasicConfig}, "\n")
}

func configKafkaClusterDataSourceGetByName(resourceName, dataSourceName, dataSourceAttributeName string, partialMatching bool) string {
	dataSourceBasicConfig := fmt.Sprintf(
		templateKafkaClusterDataNameConfig, dataSourceName, dataSourceAttributeName, constant.KafkaClusterResource+"."+resourceName+"."+clusterAttributeLocation,
		partialMatching,
	)
	baseConfig := configKafkaClusterBasic(resourceName, clusterAttributeNameValue)

	return strings.Join([]string{baseConfig, dataSourceBasicConfig}, "\n")
}

const (
	clusterResourceName = "test_kafka_cluster"
	clusterDataName     = "test_ds_kafka_cluster"

	clusterAttributeName      = "name"
	clusterAttributeNameValue = "test_kafka_cluster"

	clusterAttributeVersion      = "version"
	clusterAttributeVersionValue = "3.7.0"

	clusterAttributeSize      = "size"
	clusterAttributeSizeValue = "S"

	clusterAttributeLocation      = "location"
	clusterAttributeLocationValue = "de/fra"

	clusterAttributeBrokerAddresses      = "connections.0.broker_addresses.#"
	clusterAttributeBrokerAddressesValue = `"192.168.1.101/24", "192.168.1.102/24", "192.168.1.103/24"`

	clusterMetadataBrokerAddresses = "broker_addresses.#"
)

const templateKafkaClusterConfig = `
resource "ionoscloud_kafka_cluster" "%v" {
	name = "%v"
	version = "%v"
	size = "%v"
	location = "%v"
	connections {
		datacenter_id = ionoscloud_datacenter.test_datacenter.id
		lan_id = ionoscloud_lan.test_lan.id
		broker_addresses = [
			%v
		]
	}
}`

const templateKafkaClusterDataIDConfig = `
data "ionoscloud_kafka_cluster" "%v" {
  	id = %v
	location = %v
}`

const templateKafkaClusterDataNameConfig = `
data "ionoscloud_kafka_cluster" "%v" {
	name = %v
	location = %v
	partial_match = %v
}`

const defaultKafkaClusterBaseConfig = `
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
`
