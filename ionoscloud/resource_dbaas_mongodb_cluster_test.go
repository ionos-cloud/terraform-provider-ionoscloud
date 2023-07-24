//go:build all || dbaas || mongo
// +build all dbaas mongo

package ionoscloud

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	mongo "github.com/ionos-cloud/sdk-go-dbaas-mongo"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils/constant"
	"regexp"
	"testing"
)

func TestAccDBaaSMongoClusterBasic(t *testing.T) {
	var dbaasCluster mongo.ClusterResponse

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ExternalProviders: randomProviderVersion343(),
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckDbaasMongoClusterDestroyCheck,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckDbaasMongoClusterConfigBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDbaasMongoClusterExists(constant.DBaasMongoClusterResource+"."+constant.DBaaSClusterTestResource, &dbaasCluster),
					resource.TestCheckResourceAttr(constant.DBaasMongoClusterResource+"."+constant.DBaaSClusterTestResource, "maintenance_window.0.time", "09:00:00"),
					resource.TestCheckResourceAttr(constant.DBaasMongoClusterResource+"."+constant.DBaaSClusterTestResource, "maintenance_window.0.day_of_the_week", "Monday"),
					resource.TestCheckResourceAttr(constant.DBaasMongoClusterResource+"."+constant.DBaaSClusterTestResource, "mongodb_version", "5.0"),
					resource.TestCheckResourceAttr(constant.DBaasMongoClusterResource+"."+constant.DBaaSClusterTestResource, "instances", "1"),
					resource.TestCheckResourceAttr(constant.DBaasMongoClusterResource+"."+constant.DBaaSClusterTestResource, "display_name", constant.DBaaSClusterTestResource),
					resource.TestCheckResourceAttrPair(constant.DBaasMongoClusterResource+"."+constant.DBaaSClusterTestResource, "location", constant.DatacenterResource+".datacenter_example", "location"),
					resource.TestCheckResourceAttrPair(constant.DBaasMongoClusterResource+"."+constant.DBaaSClusterTestResource, "connections.0.datacenter_id", constant.DatacenterResource+".datacenter_example", "id"),
					resource.TestCheckResourceAttrPair(constant.DBaasMongoClusterResource+"."+constant.DBaaSClusterTestResource, "connections.0.lan_id", constant.LanResource+".lan_example", "id"),
					resource.TestCheckResourceAttr(constant.DBaasMongoClusterResource+"."+constant.DBaaSClusterTestResource, "connections.0.cidr_list.0", "192.168.1.108/24"),
					resource.TestCheckResourceAttr(constant.DBaasMongoClusterResource+"."+constant.DBaaSClusterTestResource, "template_id", "33457e53-1f8b-4ed2-8a12-2d42355aa759"),
					resource.TestCheckResourceAttr(constant.DBaasMongoClusterResource+"."+constant.DBaaSClusterTestResource, "credentials.0.username", "username"),
					resource.TestCheckResourceAttrPair(constant.DBaasMongoClusterResource+"."+constant.DBaaSClusterTestResource, "credentials.0.password", constant.RandomPassword+".dbaas_mongo_cluster_password", "result"),
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
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.DBaasMongoClusterResource+"."+constant.DBaaSClusterTestDataSourceById, "credentials.username", constant.DBaasMongoClusterResource+"."+constant.DBaaSClusterTestResource, "credentials.username"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.DBaasMongoClusterResource+"."+constant.DBaaSClusterTestDataSourceById, "credentials.password", constant.DBaasMongoClusterResource+"."+constant.DBaaSClusterTestResource, "credentials.password"),
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
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.DBaasMongoClusterResource+"."+constant.DBaaSClusterTestDataSourceByName, "credentials.username", constant.DBaasMongoClusterResource+"."+constant.DBaaSClusterTestResource, "credentials.username"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.DBaasMongoClusterResource+"."+constant.DBaaSClusterTestDataSourceByName, "credentials.password", constant.DBaasMongoClusterResource+"."+constant.DBaaSClusterTestResource, "credentials.password"),
				),
			},
			{
				Config: testAccCheckDbaasMongoClusterUpdated,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDbaasMongoClusterExists(constant.DBaasMongoClusterResource+"."+constant.DBaaSClusterTestResource, &dbaasCluster),
					resource.TestCheckResourceAttr(constant.DBaasMongoClusterResource+"."+constant.DBaaSClusterTestResource, "maintenance_window.0.time", "09:00:00"),
					resource.TestCheckResourceAttr(constant.DBaasMongoClusterResource+"."+constant.DBaaSClusterTestResource, "maintenance_window.0.day_of_the_week", "Sunday"),
					resource.TestCheckResourceAttr(constant.DBaasMongoClusterResource+"."+constant.DBaaSClusterTestResource, "mongodb_version", "5.0"),
					resource.TestCheckResourceAttr(constant.DBaasMongoClusterResource+"."+constant.DBaaSClusterTestResource, "instances", "3"),
					resource.TestCheckResourceAttr(constant.DBaasMongoClusterResource+"."+constant.DBaaSClusterTestResource, "display_name", constant.DBaaSClusterTestResource+"update"),
					resource.TestCheckResourceAttrPair(constant.DBaasMongoClusterResource+"."+constant.DBaaSClusterTestResource, "location", constant.DatacenterResource+".datacenter_example", "location"),
					resource.TestCheckResourceAttrPair(constant.DBaasMongoClusterResource+"."+constant.DBaaSClusterTestResource, "connections.0.datacenter_id", constant.DatacenterResource+".datacenter_example", "id"),
					resource.TestCheckResourceAttrPair(constant.DBaasMongoClusterResource+"."+constant.DBaaSClusterTestResource, "connections.0.lan_id", constant.LanResource+".lan_example", "id"),
					resource.TestCheckResourceAttr(constant.DBaasMongoClusterResource+"."+constant.DBaaSClusterTestResource, "connections.0.cidr_list.0", "192.168.1.108/24"),
					resource.TestCheckResourceAttr(constant.DBaasMongoClusterResource+"."+constant.DBaaSClusterTestResource, "template_id", "6b78ea06-ee0e-4689-998c-fc9c46e781f6"),
					resource.TestCheckResourceAttr(constant.DBaasMongoClusterResource+"."+constant.DBaaSClusterTestResource, "credentials.0.username", "username"),
					resource.TestCheckResourceAttrPair(constant.DBaasMongoClusterResource+"."+constant.DBaaSClusterTestResource, "credentials.0.password", constant.RandomPassword+".dbaas_mongo_cluster_password", "result"),
				),
			},
			{
				Config: testAccCheckDbaasMongoClusterUpdateTemplateAndInstances,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDbaasMongoClusterExists(constant.DBaasMongoClusterResource+"."+constant.DBaaSClusterTestResource, &dbaasCluster),
					resource.TestCheckResourceAttr(constant.DBaasMongoClusterResource+"."+constant.DBaaSClusterTestResource, "maintenance_window.0.time", "09:00:00"),
					resource.TestCheckResourceAttr(constant.DBaasMongoClusterResource+"."+constant.DBaaSClusterTestResource, "maintenance_window.0.day_of_the_week", "Sunday"),
					resource.TestCheckResourceAttr(constant.DBaasMongoClusterResource+"."+constant.DBaaSClusterTestResource, "mongodb_version", "5.0"),
					resource.TestCheckResourceAttr(constant.DBaasMongoClusterResource+"."+constant.DBaaSClusterTestResource, "instances", "3"),
					resource.TestCheckResourceAttr(constant.DBaasMongoClusterResource+"."+constant.DBaaSClusterTestResource, "display_name", constant.DBaaSClusterTestResource+"update"),
					resource.TestCheckResourceAttrPair(constant.DBaasMongoClusterResource+"."+constant.DBaaSClusterTestResource, "location", constant.DatacenterResource+".datacenter_example", "location"),
					resource.TestCheckResourceAttr(constant.DBaasMongoClusterResource+"."+constant.DBaaSClusterTestResource, "template_id", "6b78ea06-ee0e-4689-998c-fc9c46e781f6"),
					resource.TestCheckResourceAttr(constant.DBaasMongoClusterResource+"."+constant.DBaaSClusterTestResource, "credentials.0.username", "username"),
					resource.TestCheckResourceAttrPair(constant.DBaasMongoClusterResource+"."+constant.DBaaSClusterTestResource, "credentials.0.password", constant.RandomPassword+".dbaas_mongo_cluster_password", "result"),
				),
			},
			{
				Config:      testAccDataSourceDBaaSMongoClusterWrongNameError,
				ExpectError: regexp.MustCompile("no DBaaS mongo cluster found with the specified name"),
			},
		},
	})
}

func testAccCheckDbaasMongoClusterDestroyCheck(s *terraform.State) error {
	client := testAccProvider.Meta().(services.SdkBundle).MongoClient

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
		client := testAccProvider.Meta().(services.SdkBundle).MongoClient

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
			return fmt.Errorf("an error occured while fetching k8s Cluster %s: %w", rs.Primary.ID, err)
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
  location    = "de/txl"
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
  mongodb_version = "5.0"
  instances          = 1
  display_name = "` + constant.DBaaSClusterTestResource + `"
  location = ` + constant.DatacenterResource + `.datacenter_example.location
  connections   {
	datacenter_id   =  ` + constant.DatacenterResource + `.datacenter_example.id 
    lan_id          =  ` + constant.LanResource + `.lan_example.id 
    cidr_list            =  ["192.168.1.108/24"]
  }
  template_id = "33457e53-1f8b-4ed2-8a12-2d42355aa759"
  credentials {
  	username = "username"
	password = ` + constant.RandomPassword + `.dbaas_mongo_cluster_password.result
  }
}

resource ` + constant.RandomPassword + ` "dbaas_mongo_cluster_password" {
  length           = 16
  special          = true
  override_special = "!#$%&*()-_=+[]{}<>:?"
}
`

const testAccDataSourceDBaaSMongoClusterMatchId = testAccCheckDbaasMongoClusterConfigBasic + `
data ` + constant.DBaasMongoClusterResource + ` ` + constant.DBaaSClusterTestDataSourceById + ` {
  id	= ` + constant.DBaasMongoClusterResource + `.` + constant.DBaaSClusterTestResource + `.id
}
`

const testAccDataSourceDBaaSMongoClusterMatchName = testAccCheckDbaasMongoClusterConfigBasic + `
data ` + constant.DBaasMongoClusterResource + ` ` + constant.DBaaSClusterTestDataSourceByName + ` {
  display_name	= "` + constant.DBaaSClusterTestResource + `"
}
`

const testAccCheckDbaasMongoClusterUpdated = `
resource ` + constant.DatacenterResource + ` "datacenter_example" {
  name        = "datacenter_example"
  location    = "de/txl"
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
  mongodb_version = "5.0"
  instances          = 3
  display_name = "` + constant.DBaaSClusterTestResource + `update"
  location = ` + constant.DatacenterResource + `.datacenter_example.location
  connections   {
	datacenter_id   =  ` + constant.DatacenterResource + `.datacenter_example.id 
    lan_id          =  ` + constant.LanResource + `.lan_example.id 
    cidr_list            =  ["192.168.1.108/24", "192.168.1.109/24", "192.168.1.110/24"]
  }
  template_id = "6b78ea06-ee0e-4689-998c-fc9c46e781f6"
  
  credentials {
  	username = "username"
	password = ` + constant.RandomPassword + `.dbaas_mongo_cluster_password.result
  }
}
resource ` + constant.RandomPassword + ` "dbaas_mongo_cluster_password" {
  length           = 16
  special          = true
  override_special = "!#$%&*()-_=+[]{}<>:?"
}
`

const testAccCheckDbaasMongoClusterUpdateTemplateAndInstances = `
resource ` + constant.DatacenterResource + ` "datacenter_example" {
  name        = "datacenter_example"
  location    = "de/txl"
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
  mongodb_version = "5.0"
  instances          = 3
  display_name = "` + constant.DBaaSClusterTestResource + `update"
  location = ` + constant.DatacenterResource + `.datacenter_example.location
  connections   {
	datacenter_id   =  ` + constant.DatacenterResource + `.datacenter_example.id 
    lan_id          =  ` + constant.LanResource + `.lan_example.id 
    cidr_list            =  ["192.168.1.108/24", "192.168.1.109/24", "192.168.1.110/24"]
  }
  template_id = "6b78ea06-ee0e-4689-998c-fc9c46e781f6"
  
  credentials {
  	username = "username"
	password = ` + constant.RandomPassword + `.dbaas_mongo_cluster_password.result
  }
}
resource ` + constant.RandomPassword + ` "dbaas_mongo_cluster_password" {
  length           = 16
  special          = true
  override_special = "!#$%&*()-_=+[]{}<>:?"
}
`

const testAccDataSourceDBaaSMongoClusterWrongNameError = testAccCheckDbaasMongoClusterConfigBasic + `
data ` + constant.DBaasMongoClusterResource + ` ` + constant.DBaaSClusterTestDataSourceByName + ` {
  display_name	= "wrong_name"
}
`
