//go:build all || dbaas || mongo
// +build all dbaas mongo

package ionoscloud

import (
	"context"
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	mongo "github.com/ionos-cloud/sdk-go-dbaas-mongo"

	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/bundleclient"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils/constant"
)

func TestAccDBaaSMongoClusterBasic(t *testing.T) {
	var dbaasCluster mongo.ClusterResponse

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactoriesInternal(t, &testAccProvider),
		CheckDestroy:             testAccCheckDbaasMongoClusterDestroyCheck,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckDbaasMongoClusterConfigBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDbaasMongoClusterExists(constant.DBaasMongoClusterResource+"."+constant.DBaaSClusterTestResource, &dbaasCluster),
					resource.TestCheckResourceAttr(constant.DBaasMongoClusterResource+"."+constant.DBaaSClusterTestResource, "maintenance_window.0.time", "09:00:00"),
					resource.TestCheckResourceAttr(constant.DBaasMongoClusterResource+"."+constant.DBaaSClusterTestResource, "maintenance_window.0.day_of_the_week", "Monday"),
					resource.TestCheckResourceAttr(constant.DBaasMongoClusterResource+"."+constant.DBaaSClusterTestResource, "mongodb_version", mongoVersion),
					resource.TestCheckResourceAttr(constant.DBaasMongoClusterResource+"."+constant.DBaaSClusterTestResource, "instances", "1"),
					resource.TestCheckResourceAttr(constant.DBaasMongoClusterResource+"."+constant.DBaaSClusterTestResource, "display_name", constant.DBaaSClusterTestResource),
					resource.TestCheckResourceAttrPair(constant.DBaasMongoClusterResource+"."+constant.DBaaSClusterTestResource, "location", constant.DatacenterResource+".datacenter_example", "location"),
					resource.TestCheckResourceAttrPair(constant.DBaasMongoClusterResource+"."+constant.DBaaSClusterTestResource, "connections.0.datacenter_id", constant.DatacenterResource+".datacenter_example", "id"),
					resource.TestCheckResourceAttrPair(constant.DBaasMongoClusterResource+"."+constant.DBaaSClusterTestResource, "connections.0.lan_id", constant.LanResource+".lan_example", "id"),
					resource.TestCheckResourceAttr(constant.DBaasMongoClusterResource+"."+constant.DBaaSClusterTestResource, "connections.0.cidr_list.0", "192.168.1.108/24"),
					resource.TestCheckResourceAttr(constant.DBaasMongoClusterResource+"."+constant.DBaaSClusterTestResource, "template_id", "33457e53-1f8b-4ed2-8a12-2d42355aa759"),
				),
			},
			{
				Config: testAccDataSourceDBaaSMongoClusterMatchId,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.DBaasMongoClusterResource+"."+constant.DBaaSClusterTestDataSourceById, "maintenance_window.day_of_the_week", constant.DBaasMongoClusterResource+"."+constant.DBaaSClusterTestResource, "maintenance_window.day_of_the_week"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.DBaasMongoClusterResource+"."+constant.DBaaSClusterTestDataSourceById, "maintenance_window.time", constant.DBaasMongoClusterResource+"."+constant.DBaaSClusterTestResource, "maintenance_window.time"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.DBaasMongoClusterResource+"."+constant.DBaaSClusterTestDataSourceById, "mongodb_version", constant.DBaasMongoClusterResource+"."+constant.DBaaSClusterTestResource, "mongodb_version"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.DBaasMongoClusterResource+"."+constant.DBaaSClusterTestDataSourceById, "instances", constant.DBaasMongoClusterResource+"."+constant.DBaaSClusterTestResource, "instances"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.DBaasMongoClusterResource+"."+constant.DBaaSClusterTestDataSourceById, "display_name", constant.DBaasMongoClusterResource+"."+constant.DBaaSClusterTestResource, "display_name"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.DBaasMongoClusterResource+"."+constant.DBaaSClusterTestDataSourceById, "location", constant.DBaasMongoClusterResource+"."+constant.DBaaSClusterTestResource, "location"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.DBaasMongoClusterResource+"."+constant.DBaaSClusterTestDataSourceById, "connections.datacenter_id", constant.DBaasMongoClusterResource+"."+constant.DBaaSClusterTestResource, "connections.datacenter_id"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.DBaasMongoClusterResource+"."+constant.DBaaSClusterTestDataSourceById, "connections.lan_id", constant.DBaasMongoClusterResource+"."+constant.DBaaSClusterTestResource, "connections.lan_id"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.DBaasMongoClusterResource+"."+constant.DBaaSClusterTestDataSourceById, "connections.0.cidr_list.0", constant.DBaasMongoClusterResource+"."+constant.DBaaSClusterTestResource, "connections.0.cidr_list.0"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.DBaasMongoClusterResource+"."+constant.DBaaSClusterTestDataSourceById, "template_id", constant.DBaasMongoClusterResource+"."+constant.DBaaSClusterTestResource, "template_id"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.DBaasMongoClusterResource+"."+constant.DBaaSClusterTestDataSourceById, "connection_string", constant.DBaasMongoClusterResource+"."+constant.DBaaSClusterTestResource, "connection_string"),
				),
			},
			{
				Config: testAccDataSourceDBaaSMongoClusterMatchName,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.DBaasMongoClusterResource+"."+constant.DBaaSClusterTestDataSourceByName, "maintenance_window.day_of_the_week", constant.DBaasMongoClusterResource+"."+constant.DBaaSClusterTestResource, "maintenance_window.day_of_the_week"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.DBaasMongoClusterResource+"."+constant.DBaaSClusterTestDataSourceByName, "maintenance_window.time", constant.DBaasMongoClusterResource+"."+constant.DBaaSClusterTestResource, "maintenance_window.time"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.DBaasMongoClusterResource+"."+constant.DBaaSClusterTestDataSourceByName, "mongodb_version", constant.DBaasMongoClusterResource+"."+constant.DBaaSClusterTestResource, "mongodb_version"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.DBaasMongoClusterResource+"."+constant.DBaaSClusterTestDataSourceByName, "instances", constant.DBaasMongoClusterResource+"."+constant.DBaaSClusterTestResource, "instances"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.DBaasMongoClusterResource+"."+constant.DBaaSClusterTestDataSourceByName, "display_name", constant.DBaasMongoClusterResource+"."+constant.DBaaSClusterTestResource, "display_name"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.DBaasMongoClusterResource+"."+constant.DBaaSClusterTestDataSourceByName, "location", constant.DBaasMongoClusterResource+"."+constant.DBaaSClusterTestResource, "location"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.DBaasMongoClusterResource+"."+constant.DBaaSClusterTestDataSourceByName, "connections.datacenter_id", constant.DBaasMongoClusterResource+"."+constant.DBaaSClusterTestResource, "connections.datacenter_id"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.DBaasMongoClusterResource+"."+constant.DBaaSClusterTestDataSourceByName, "connections.lan_id", constant.DBaasMongoClusterResource+"."+constant.DBaaSClusterTestResource, "connections.lan_id"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.DBaasMongoClusterResource+"."+constant.DBaaSClusterTestDataSourceByName, "connections.0.cidr_list", constant.DBaasMongoClusterResource+"."+constant.DBaaSClusterTestResource, "connections.0.cidr_list"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.DBaasMongoClusterResource+"."+constant.DBaaSClusterTestDataSourceByName, "template_id", constant.DBaasMongoClusterResource+"."+constant.DBaaSClusterTestResource, "template_id"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.DBaasMongoClusterResource+"."+constant.DBaaSClusterTestDataSourceByName, "connection_string", constant.DBaasMongoClusterResource+"."+constant.DBaaSClusterTestResource, "connection_string"),
				),
			},
			{
				Config: testAccCheckDbaasMongoClusterUpdated,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDbaasMongoClusterExists(constant.DBaasMongoClusterResource+"."+constant.DBaaSClusterTestResource, &dbaasCluster),
					resource.TestCheckResourceAttr(constant.DBaasMongoClusterResource+"."+constant.DBaaSClusterTestResource, "maintenance_window.0.time", "09:00:00"),
					resource.TestCheckResourceAttr(constant.DBaasMongoClusterResource+"."+constant.DBaaSClusterTestResource, "maintenance_window.0.day_of_the_week", "Sunday"),
					resource.TestCheckResourceAttr(constant.DBaasMongoClusterResource+"."+constant.DBaaSClusterTestResource, "mongodb_version", mongoVersion),
					resource.TestCheckResourceAttr(constant.DBaasMongoClusterResource+"."+constant.DBaaSClusterTestResource, "instances", "3"),
					resource.TestCheckResourceAttr(constant.DBaasMongoClusterResource+"."+constant.DBaaSClusterTestResource, "display_name", constant.DBaaSClusterTestResource+"update"),
					resource.TestCheckResourceAttrPair(constant.DBaasMongoClusterResource+"."+constant.DBaaSClusterTestResource, "location", constant.DatacenterResource+".datacenter_example", "location"),
					resource.TestCheckResourceAttrPair(constant.DBaasMongoClusterResource+"."+constant.DBaaSClusterTestResource, "connections.0.datacenter_id", constant.DatacenterResource+".datacenter_example", "id"),
					resource.TestCheckResourceAttrPair(constant.DBaasMongoClusterResource+"."+constant.DBaaSClusterTestResource, "connections.0.lan_id", constant.LanResource+".lan_example", "id"),
					resource.TestCheckResourceAttr(constant.DBaasMongoClusterResource+"."+constant.DBaaSClusterTestResource, "connections.0.cidr_list.0", "192.168.1.108/24"),
					resource.TestCheckResourceAttr(constant.DBaasMongoClusterResource+"."+constant.DBaaSClusterTestResource, "template_id", "6b78ea06-ee0e-4689-998c-fc9c46e781f6"),
				),
			},
			{
				Config: testAccCheckDbaasMongoClusterUpdateTemplateAndInstances,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDbaasMongoClusterExists(constant.DBaasMongoClusterResource+"."+constant.DBaaSClusterTestResource, &dbaasCluster),
					resource.TestCheckResourceAttr(constant.DBaasMongoClusterResource+"."+constant.DBaaSClusterTestResource, "maintenance_window.0.time", "09:00:00"),
					resource.TestCheckResourceAttr(constant.DBaasMongoClusterResource+"."+constant.DBaaSClusterTestResource, "maintenance_window.0.day_of_the_week", "Sunday"),
					resource.TestCheckResourceAttr(constant.DBaasMongoClusterResource+"."+constant.DBaaSClusterTestResource, "mongodb_version", mongoVersion),
					resource.TestCheckResourceAttr(constant.DBaasMongoClusterResource+"."+constant.DBaaSClusterTestResource, "instances", "3"),
					resource.TestCheckResourceAttr(constant.DBaasMongoClusterResource+"."+constant.DBaaSClusterTestResource, "display_name", constant.DBaaSClusterTestResource+"update"),
					resource.TestCheckResourceAttrPair(constant.DBaasMongoClusterResource+"."+constant.DBaaSClusterTestResource, "location", constant.DatacenterResource+".datacenter_example", "location"),
					resource.TestCheckResourceAttr(constant.DBaasMongoClusterResource+"."+constant.DBaaSClusterTestResource, "template_id", "6b78ea06-ee0e-4689-998c-fc9c46e781f6"),
				),
			},
			{
				Config:      testAccDataSourceDBaaSMongoClusterWrongNameError,
				ExpectError: regexp.MustCompile("no DBaaS mongo cluster found with the specified name"),
			},
		},
	})
}

func TestAccMongoClusterEnterpriseEditionBasic(t *testing.T) {
	var dbaasCluster mongo.ClusterResponse

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactoriesInternal(t, &testAccProvider),
		CheckDestroy:             testAccCheckDbaasMongoClusterDestroyCheck,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckMongoClusterEnterpriseBasic,
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckDbaasMongoClusterExists(constant.DBaasMongoClusterResource+"."+constant.DBaaSClusterTestResource, &dbaasCluster),
					resource.TestCheckResourceAttr(constant.DBaasMongoClusterResource+"."+constant.DBaaSClusterTestResource, "maintenance_window.0.time", "09:00:00"),
					resource.TestCheckResourceAttr(constant.DBaasMongoClusterResource+"."+constant.DBaaSClusterTestResource, "maintenance_window.0.day_of_the_week", "Monday"),
					resource.TestCheckResourceAttr(constant.DBaasMongoClusterResource+"."+constant.DBaaSClusterTestResource, "mongodb_version", mongoVersion),
					resource.TestCheckResourceAttr(constant.DBaasMongoClusterResource+"."+constant.DBaaSClusterTestResource, "instances", "3"),
					resource.TestCheckResourceAttr(constant.DBaasMongoClusterResource+"."+constant.DBaaSClusterTestResource, "display_name", constant.DBaaSClusterTestResource),
					resource.TestCheckResourceAttrPair(constant.DBaasMongoClusterResource+"."+constant.DBaaSClusterTestResource, "location", constant.DatacenterResource+".datacenter_example", "location"),
					resource.TestCheckResourceAttrPair(constant.DBaasMongoClusterResource+"."+constant.DBaaSClusterTestResource, "connections.0.datacenter_id", constant.DatacenterResource+".datacenter_example", "id"),
					resource.TestCheckResourceAttrPair(constant.DBaasMongoClusterResource+"."+constant.DBaaSClusterTestResource, "connections.0.lan_id", constant.LanResource+".lan_example", "id"),
					resource.TestCheckResourceAttr(constant.DBaasMongoClusterResource+"."+constant.DBaaSClusterTestResource, "connections.0.cidr_list.0", "192.168.1.108/24"),
					// enterprise edition
					// add after api adds support
					// resource.TestCheckResourceAttr(constant.DBaasMongoClusterResource+"."+constant.DBaaSClusterTestResource, "connections.0.whitelist.0", "192.168.1.108/24"),
					resource.TestCheckResourceAttr(constant.DBaasMongoClusterResource+"."+constant.DBaaSClusterTestResource, "edition", "enterprise"),
					resource.TestCheckResourceAttr(constant.DBaasMongoClusterResource+"."+constant.DBaaSClusterTestResource, "type", "sharded-cluster"),
					resource.TestCheckResourceAttr(constant.DBaasMongoClusterResource+"."+constant.DBaaSClusterTestResource, "ram", "2048"),
					resource.TestCheckResourceAttr(constant.DBaasMongoClusterResource+"."+constant.DBaaSClusterTestResource, "shards", "2"),
					resource.TestCheckResourceAttr(constant.DBaasMongoClusterResource+"."+constant.DBaaSClusterTestResource, "storage_size", "5120"),
					resource.TestCheckResourceAttr(constant.DBaasMongoClusterResource+"."+constant.DBaaSClusterTestResource, "storage_type", "HDD"),
				),
			},
			{
				Config: testAccCheckDbaasMongoClusterEnterpriseUpdated,
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckDbaasMongoClusterExists(constant.DBaasMongoClusterResource+"."+constant.DBaaSClusterTestResource, &dbaasCluster),
					resource.TestCheckResourceAttr(constant.DBaasMongoClusterResource+"."+constant.DBaaSClusterTestResource, "maintenance_window.0.time", "09:00:00"),
					resource.TestCheckResourceAttr(constant.DBaasMongoClusterResource+"."+constant.DBaaSClusterTestResource, "maintenance_window.0.day_of_the_week", "Sunday"),
					resource.TestCheckResourceAttr(constant.DBaasMongoClusterResource+"."+constant.DBaaSClusterTestResource, "mongodb_version", mongoVersion),
					resource.TestCheckResourceAttr(constant.DBaasMongoClusterResource+"."+constant.DBaaSClusterTestResource, "instances", "3"),
					resource.TestCheckResourceAttrPair(constant.DBaasMongoClusterResource+"."+constant.DBaaSClusterTestResource, "location", constant.DatacenterResource+".datacenter_example", "location"),
					resource.TestCheckResourceAttrPair(constant.DBaasMongoClusterResource+"."+constant.DBaaSClusterTestResource, "connections.0.datacenter_id", constant.DatacenterResource+".datacenter_example", "id"),
					resource.TestCheckResourceAttrPair(constant.DBaasMongoClusterResource+"."+constant.DBaaSClusterTestResource, "connections.0.lan_id", constant.LanResource+".lan_example", "id"),
					resource.TestCheckResourceAttr(constant.DBaasMongoClusterResource+"."+constant.DBaaSClusterTestResource, "connections.0.cidr_list.0", "192.168.1.108/24"),
					// enterprise edition
					resource.TestCheckResourceAttr(constant.DBaasMongoClusterResource+"."+constant.DBaaSClusterTestResource, "edition", "enterprise"),
					resource.TestCheckResourceAttr(constant.DBaasMongoClusterResource+"."+constant.DBaaSClusterTestResource, "type", "sharded-cluster"),
					resource.TestCheckResourceAttr(constant.DBaasMongoClusterResource+"."+constant.DBaaSClusterTestResource, "shards", "3"),
					resource.TestCheckResourceAttr(constant.DBaasMongoClusterResource+"."+constant.DBaaSClusterTestResource, "ram", "4096"),
					resource.TestCheckResourceAttr(constant.DBaasMongoClusterResource+"."+constant.DBaaSClusterTestResource, "storage_size", "5120"),
					resource.TestCheckResourceAttr(constant.DBaasMongoClusterResource+"."+constant.DBaaSClusterTestResource, "storage_type", "HDD"),
					resource.TestCheckResourceAttr(constant.DBaasMongoClusterResource+"."+constant.DBaaSClusterTestResource, "bi_connector.0.enabled", "true"),
					resource.TestCheckResourceAttrSet(constant.DBaasMongoClusterResource+"."+constant.DBaaSClusterTestResource, "bi_connector.0.host"),
					resource.TestCheckResourceAttrSet(constant.DBaasMongoClusterResource+"."+constant.DBaaSClusterTestResource, "bi_connector.0.port"),
				),
			},
			{
				Config: testAccDataSourceDBaaSMongoClusterEnterpriseMatchId,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.DBaasMongoClusterResource+"."+constant.DBaaSClusterTestDataSourceById, "maintenance_window.day_of_the_week",
						constant.DBaasMongoClusterResource+"."+constant.DBaaSClusterTestResource, "maintenance_window.day_of_the_week"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.DBaasMongoClusterResource+"."+constant.DBaaSClusterTestDataSourceById, "maintenance_window.time",
						constant.DBaasMongoClusterResource+"."+constant.DBaaSClusterTestResource, "maintenance_window.time"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.DBaasMongoClusterResource+"."+constant.DBaaSClusterTestDataSourceById, "mongodb_version",
						constant.DBaasMongoClusterResource+"."+constant.DBaaSClusterTestResource, "mongodb_version"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.DBaasMongoClusterResource+"."+constant.DBaaSClusterTestDataSourceById, "instances",
						constant.DBaasMongoClusterResource+"."+constant.DBaaSClusterTestResource, "instances"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.DBaasMongoClusterResource+"."+constant.DBaaSClusterTestDataSourceById, "display_name",
						constant.DBaasMongoClusterResource+"."+constant.DBaaSClusterTestResource, "display_name"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.DBaasMongoClusterResource+"."+constant.DBaaSClusterTestDataSourceById, "location",
						constant.DBaasMongoClusterResource+"."+constant.DBaaSClusterTestResource, "location"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.DBaasMongoClusterResource+"."+constant.DBaaSClusterTestDataSourceById, "connections.datacenter_id",
						constant.DBaasMongoClusterResource+"."+constant.DBaaSClusterTestResource, "connections.datacenter_id"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.DBaasMongoClusterResource+"."+constant.DBaaSClusterTestDataSourceById, "connections.lan_id",
						constant.DBaasMongoClusterResource+"."+constant.DBaaSClusterTestResource, "connections.lan_id"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.DBaasMongoClusterResource+"."+constant.DBaaSClusterTestDataSourceById, "connections.0.cidr_list.0",
						constant.DBaasMongoClusterResource+"."+constant.DBaaSClusterTestResource, "connections.0.cidr_list.0"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.DBaasMongoClusterResource+"."+constant.DBaaSClusterTestDataSourceById, "template_id",
						constant.DBaasMongoClusterResource+"."+constant.DBaaSClusterTestResource, "template_id"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.DBaasMongoClusterResource+"."+constant.DBaaSClusterTestDataSourceById, "connection_string",
						constant.DBaasMongoClusterResource+"."+constant.DBaaSClusterTestResource, "connection_string"),
					// enterprise here
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.DBaasMongoClusterResource+"."+constant.DBaaSClusterTestDataSourceById, "edition",
						constant.DBaasMongoClusterResource+"."+constant.DBaaSClusterTestResource, "edition"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.DBaasMongoClusterResource+"."+constant.DBaaSClusterTestDataSourceById, "type",
						constant.DBaasMongoClusterResource+"."+constant.DBaaSClusterTestResource, "type"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.DBaasMongoClusterResource+"."+constant.DBaaSClusterTestDataSourceById, "ram",
						constant.DBaasMongoClusterResource+"."+constant.DBaaSClusterTestResource, "ram"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.DBaasMongoClusterResource+"."+constant.DBaaSClusterTestDataSourceById, "shards",
						constant.DBaasMongoClusterResource+"."+constant.DBaaSClusterTestResource, "shards"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.DBaasMongoClusterResource+"."+constant.DBaaSClusterTestDataSourceById, "storage_size",
						constant.DBaasMongoClusterResource+"."+constant.DBaaSClusterTestResource, "storage_size"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.DBaasMongoClusterResource+"."+constant.DBaaSClusterTestDataSourceById, "storage_type",
						constant.DBaasMongoClusterResource+"."+constant.DBaaSClusterTestResource, "storage_type"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.DBaasMongoClusterResource+"."+constant.DBaaSClusterTestDataSourceById, "bi_connector.0.enabled",
						constant.DBaasMongoClusterResource+"."+constant.DBaaSClusterTestResource, "bi_connector.0.enabled"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.DBaasMongoClusterResource+"."+constant.DBaaSClusterTestDataSourceById, "bi_connector.0.host",
						constant.DBaasMongoClusterResource+"."+constant.DBaaSClusterTestResource, "bi_connector.0.host"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.DBaasMongoClusterResource+"."+constant.DBaaSClusterTestDataSourceById, "bi_connector.0.port",
						constant.DBaasMongoClusterResource+"."+constant.DBaaSClusterTestResource, "bi_connector.0.port"),
				),
			},
			{
				Config: testAccDataSourceDBaaSMongoClusterEnterpriseMatchName,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.DBaasMongoClusterResource+"."+constant.DBaaSClusterTestDataSourceByName, "maintenance_window.day_of_the_week", constant.DBaasMongoClusterResource+"."+constant.DBaaSClusterTestResource, "maintenance_window.day_of_the_week"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.DBaasMongoClusterResource+"."+constant.DBaaSClusterTestDataSourceByName, "maintenance_window.time", constant.DBaasMongoClusterResource+"."+constant.DBaaSClusterTestResource, "maintenance_window.time"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.DBaasMongoClusterResource+"."+constant.DBaaSClusterTestDataSourceByName, "mongodb_version", constant.DBaasMongoClusterResource+"."+constant.DBaaSClusterTestResource, "mongodb_version"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.DBaasMongoClusterResource+"."+constant.DBaaSClusterTestDataSourceByName, "instances", constant.DBaasMongoClusterResource+"."+constant.DBaaSClusterTestResource, "instances"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.DBaasMongoClusterResource+"."+constant.DBaaSClusterTestDataSourceByName, "display_name", constant.DBaasMongoClusterResource+"."+constant.DBaaSClusterTestResource, "display_name"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.DBaasMongoClusterResource+"."+constant.DBaaSClusterTestDataSourceByName, "location", constant.DBaasMongoClusterResource+"."+constant.DBaaSClusterTestResource, "location"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.DBaasMongoClusterResource+"."+constant.DBaaSClusterTestDataSourceByName, "connections.datacenter_id", constant.DBaasMongoClusterResource+"."+constant.DBaaSClusterTestResource, "connections.datacenter_id"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.DBaasMongoClusterResource+"."+constant.DBaaSClusterTestDataSourceByName, "connections.lan_id", constant.DBaasMongoClusterResource+"."+constant.DBaaSClusterTestResource, "connections.lan_id"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.DBaasMongoClusterResource+"."+constant.DBaaSClusterTestDataSourceByName, "connections.0.cidr_list", constant.DBaasMongoClusterResource+"."+constant.DBaaSClusterTestResource, "connections.0.cidr_list"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.DBaasMongoClusterResource+"."+constant.DBaaSClusterTestDataSourceByName, "template_id", constant.DBaasMongoClusterResource+"."+constant.DBaaSClusterTestResource, "template_id"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.DBaasMongoClusterResource+"."+constant.DBaaSClusterTestDataSourceByName, "connection_string", constant.DBaasMongoClusterResource+"."+constant.DBaaSClusterTestResource, "connection_string"),
					// enterprise here
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.DBaasMongoClusterResource+"."+constant.DBaaSClusterTestDataSourceByName, "edition",
						constant.DBaasMongoClusterResource+"."+constant.DBaaSClusterTestResource, "edition"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.DBaasMongoClusterResource+"."+constant.DBaaSClusterTestDataSourceByName, "type",
						constant.DBaasMongoClusterResource+"."+constant.DBaaSClusterTestResource, "type"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.DBaasMongoClusterResource+"."+constant.DBaaSClusterTestDataSourceByName, "ram",
						constant.DBaasMongoClusterResource+"."+constant.DBaaSClusterTestResource, "ram"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.DBaasMongoClusterResource+"."+constant.DBaaSClusterTestDataSourceByName, "shards",
						constant.DBaasMongoClusterResource+"."+constant.DBaaSClusterTestResource, "shards"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.DBaasMongoClusterResource+"."+constant.DBaaSClusterTestDataSourceByName, "storage_size",
						constant.DBaasMongoClusterResource+"."+constant.DBaaSClusterTestResource, "storage_size"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.DBaasMongoClusterResource+"."+constant.DBaaSClusterTestDataSourceByName, "storage_type",
						constant.DBaasMongoClusterResource+"."+constant.DBaaSClusterTestResource, "storage_type"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.DBaasMongoClusterResource+"."+constant.DBaaSClusterTestDataSourceByName, "bi_connector.0.enabled",
						constant.DBaasMongoClusterResource+"."+constant.DBaaSClusterTestResource, "bi_connector.0.enabled"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.DBaasMongoClusterResource+"."+constant.DBaaSClusterTestDataSourceByName, "bi_connector.0.host",
						constant.DBaasMongoClusterResource+"."+constant.DBaaSClusterTestResource, "bi_connector.0.host"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.DBaasMongoClusterResource+"."+constant.DBaaSClusterTestDataSourceByName, "bi_connector.0.port",
						constant.DBaasMongoClusterResource+"."+constant.DBaaSClusterTestResource, "bi_connector.0.port"),
				),
			},
		},
	})
}

func testAccCheckDbaasMongoClusterDestroyCheck(s *terraform.State) error {
	client := testAccProvider.Meta().(bundleclient.SdkBundle).MongoClient

	ctx, cancel := context.WithTimeout(context.Background(), *resourceDefaultTimeouts.Default)

	if cancel != nil {
		defer cancel()
	}

	for _, rs := range s.RootModule().Resources {

		_, apiResponse, err := client.GetCluster(ctx, rs.Primary.ID)

		if err != nil {
			if !apiResponse.HttpNotFound() {
				return fmt.Errorf("an error occurred while checking the destruction of dbaas mongo cluster %s: %w", rs.Primary.ID, err)
			}
		} else {
			return fmt.Errorf("mongo cluster %s still exists", rs.Primary.ID)
		}

	}

	return nil
}

func testAccCheckDbaasMongoClusterExists(n string, cluster *mongo.ClusterResponse) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(bundleclient.SdkBundle).MongoClient

		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no Record ID is set")
		}

		ctx, cancel := context.WithTimeout(context.Background(), *resourceDefaultTimeouts.Default)

		if cancel != nil {
			defer cancel()
		}

		foundCluster, _, err := client.GetCluster(ctx, rs.Primary.ID)

		if err != nil {
			return fmt.Errorf("an error occurred while fetching k8s Cluster %s: %w", rs.Primary.ID, err)
		}
		if *foundCluster.Id != rs.Primary.ID {
			return fmt.Errorf("record not found")
		}
		cluster = &foundCluster

		return nil
	}
}

const testAccCheckDbaasMongoClusterConfigBasic = `
resource ` + constant.DatacenterResource + ` "datacenter_example" {
  name        = "datacenter_example"
  location    = "gb/lhr"
  description = "Datacenter for testing dbaas cluster"
}

resource ` + constant.LanResource + ` "lan_example" {
  datacenter_id = ` + constant.DatacenterResource + `.datacenter_example.id 
  public        = false
  name          = "lan_example"
}

resource ` + constant.DBaasMongoClusterResource + ` ` + constant.DBaaSClusterTestResource + ` {
  maintenance_window {
    day_of_the_week  = "Monday"
    time             = "09:00:00"
  }
  mongodb_version = "` + mongoVersion + `"
  instances          = 1
  display_name = "` + constant.DBaaSClusterTestResource + `"
  location = ` + constant.DatacenterResource + `.datacenter_example.location
  connections   {
    datacenter_id   =  ` + constant.DatacenterResource + `.datacenter_example.id 
    lan_id          =  ` + constant.LanResource + `.lan_example.id 
    cidr_list            =  ["192.168.1.108/24"]
  }
  template_id = "33457e53-1f8b-4ed2-8a12-2d42355aa759"
}
`

const testAccDataSourceDBaaSMongoClusterMatchId = testAccCheckDbaasMongoClusterConfigBasic + `
data ` + constant.DBaasMongoClusterResource + ` ` + constant.DBaaSClusterTestDataSourceById + ` {
  id = ` + constant.DBaasMongoClusterResource + `.` + constant.DBaaSClusterTestResource + `.id
}
`

const testAccDataSourceDBaaSMongoClusterMatchName = testAccCheckDbaasMongoClusterConfigBasic + `
data ` + constant.DBaasMongoClusterResource + ` ` + constant.DBaaSClusterTestDataSourceByName + ` {
  display_name = "` + constant.DBaaSClusterTestResource + `"
}
`

const testAccCheckDbaasMongoClusterUpdated = `
resource ` + constant.DatacenterResource + ` "datacenter_example" {
  name        = "datacenter_example"
  location    = "gb/lhr"
  description = "Datacenter for testing dbaas cluster"
}

resource ` + constant.LanResource + ` "lan_example" {
  datacenter_id = ` + constant.DatacenterResource + `.datacenter_example.id 
  public        = false
  name          = "lan_example"
}

resource ` + constant.DBaasMongoClusterResource + ` ` + constant.DBaaSClusterTestResource + ` {
  maintenance_window {
    day_of_the_week  = "Sunday"
    time             = "09:00:00"
  }
  mongodb_version = "` + mongoVersion + `"
  instances       = 3
  display_name    = "` + constant.DBaaSClusterTestResource + `update"
  location        = ` + constant.DatacenterResource + `.datacenter_example.location
  connections   {
    datacenter_id   =  ` + constant.DatacenterResource + `.datacenter_example.id 
    lan_id          =  ` + constant.LanResource + `.lan_example.id 
    cidr_list       =  ["192.168.1.108/24", "192.168.1.109/24", "192.168.1.110/24"]
  }
  template_id = "6b78ea06-ee0e-4689-998c-fc9c46e781f6" 
}
`

const testAccCheckDbaasMongoClusterUpdateTemplateAndInstances = `
resource ` + constant.DatacenterResource + ` "datacenter_example" {
  name        = "datacenter_example"
  location    = "gb/lhr"
  description = "Datacenter for testing dbaas cluster"
}

resource ` + constant.LanResource + ` "lan_example" {
  datacenter_id = ` + constant.DatacenterResource + `.datacenter_example.id 
  public        = false
  name          = "lan_example"
}

resource ` + constant.DBaasMongoClusterResource + ` ` + constant.DBaaSClusterTestResource + ` {
  maintenance_window {
    day_of_the_week  = "Sunday"
    time             = "09:00:00"
  }
  mongodb_version = "` + mongoVersion + `"
  instances          = 3
  display_name = "` + constant.DBaaSClusterTestResource + `update"
  location = ` + constant.DatacenterResource + `.datacenter_example.location
  connections   {
    datacenter_id   =  ` + constant.DatacenterResource + `.datacenter_example.id 
    lan_id          =  ` + constant.LanResource + `.lan_example.id 
    cidr_list       =  ["192.168.1.108/24", "192.168.1.109/24", "192.168.1.110/24"]
  }
  template_id = "6b78ea06-ee0e-4689-998c-fc9c46e781f6"
}
`

const testAccDataSourceDBaaSMongoClusterWrongNameError = testAccCheckDbaasMongoClusterConfigBasic + `
data ` + constant.DBaasMongoClusterResource + ` ` + constant.DBaaSClusterTestDataSourceByName + ` {
  display_name = "wrong_name"
}
`

const testAccCheckMongoClusterEnterpriseBasic = `
resource ` + constant.DatacenterResource + ` "datacenter_example" {
  name        = "datacenter_example"
  location    = "gb/lhr"
  description = "Datacenter for testing dbaas cluster"
}

resource ` + constant.LanResource + ` "lan_example" {
  datacenter_id = ` + constant.DatacenterResource + `.datacenter_example.id 
  public        = false
  name          = "lan_example"
}

resource ` + constant.DBaasMongoClusterResource + ` ` + constant.DBaaSClusterTestResource + ` {
  maintenance_window {
    day_of_the_week  = "Monday"
    time             = "09:00:00"
  }
  mongodb_version = "` + mongoVersion + `"
  instances       = 3
  display_name    = "` + constant.DBaaSClusterTestResource + `"
  location        = ` + constant.DatacenterResource + `.datacenter_example.location
  connections   {
	datacenter_id   =  ` + constant.DatacenterResource + `.datacenter_example.id 
    lan_id          =  ` + constant.LanResource + `.lan_example.id 
    cidr_list       =  ["192.168.1.108/24", "192.168.1.109/24", "192.168.1.110/24"]
  }
  type         = "sharded-cluster"
  shards       = 2
  edition      = "enterprise"
  ram          = 2048
  cores        = 1
  storage_size = 5120
  storage_type = "HDD"
  backup {
    location = "de"
  }
}
`
const testAccCheckDbaasMongoClusterEnterpriseUpdated = `
resource ` + constant.DatacenterResource + ` "datacenter_example" {
  name        = "datacenter_example"
  location    = "gb/lhr"
  description = "Datacenter for testing dbaas cluster"
}

resource ` + constant.LanResource + ` "lan_example" {
  datacenter_id = ` + constant.DatacenterResource + `.datacenter_example.id
  public        = false
  name          = "lan_example"
}

resource ` + constant.DBaasMongoClusterResource + ` ` + constant.DBaaSClusterTestResource + ` {
  maintenance_window {
    day_of_the_week  = "Sunday"
    time             = "09:00:00"
  }
  mongodb_version    = "` + mongoVersion + `"
  instances          = 3
  display_name       = "` + constant.DBaaSClusterTestResource + `update"
  location           = ` + constant.DatacenterResource + `.datacenter_example.location
  connections   {
    datacenter_id =  ` + constant.DatacenterResource + `.datacenter_example.id 
    lan_id        =  ` + constant.LanResource + `.lan_example.id 
    cidr_list     =  ["192.168.1.108/24", "192.168.1.109/24", "192.168.1.110/24"]
  }
  bi_connector {
    enabled = true
  }
  type         = "sharded-cluster"
  shards       = 3
  edition      = "enterprise"
  ram          = 4096
  cores        = 2
  storage_size = 5120
  storage_type = "HDD"
  backup {
    location = "de"
  }
}

`

const testAccDataSourceDBaaSMongoClusterEnterpriseMatchId = testAccCheckDbaasMongoClusterEnterpriseUpdated + `
data ` + constant.DBaasMongoClusterResource + ` ` + constant.DBaaSClusterTestDataSourceById + ` {
  id	= ` + constant.DBaasMongoClusterResource + `.` + constant.DBaaSClusterTestResource + `.id
}
`
const testAccDataSourceDBaaSMongoClusterEnterpriseMatchName = testAccCheckDbaasMongoClusterEnterpriseUpdated + `
data ` + constant.DBaasMongoClusterResource + ` ` + constant.DBaaSClusterTestDataSourceByName + ` {
  display_name	= "` + constant.DBaaSClusterTestResource + `update"
}
`
