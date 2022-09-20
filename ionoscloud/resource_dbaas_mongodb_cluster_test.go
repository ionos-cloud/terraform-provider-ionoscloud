//go:build all || dbaas
// +build all dbaas

package ionoscloud

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	mongo "github.com/ionos-cloud/sdk-go-dbaas-mongo"
	"regexp"
	"testing"
)

func TestAccDBaaSMongoClusterBasic(t *testing.T) {
	var dbaasCluster mongo.ClusterResponse

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckDbaasMongoClusterDestroyCheck,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckDbaasMongoClusterConfigBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDbaasMongoClusterExists(DBaasMongoClusterResource+"."+DBaaSClusterTestResource, &dbaasCluster),
					//resource.TestCheckResourceAttr(DBaasMongoClusterResource+"."+DBaaSClusterTestResource, "mongo_cluster", K8sClusterResource+"."+K8sClusterTestResource+".id"),
					resource.TestCheckResourceAttr(DBaasMongoClusterResource+"."+DBaaSClusterTestResource, "maintenance_window.0.time", "09:00:00"),
					resource.TestCheckResourceAttr(DBaasMongoClusterResource+"."+DBaaSClusterTestResource, "maintenance_window.0.day_of_the_week", "Sunday"),
					resource.TestCheckResourceAttr(DBaasMongoClusterResource+"."+DBaaSClusterTestResource, "mongodb_version", "5.0"),
					resource.TestCheckResourceAttr(DBaasMongoClusterResource+"."+DBaaSClusterTestResource, "instances", "3"),
					resource.TestCheckResourceAttr(DBaasMongoClusterResource+"."+DBaaSClusterTestResource, "display_name", DBaaSClusterTestResource),
					resource.TestCheckResourceAttrPair(DBaasMongoClusterResource+"."+DBaaSClusterTestResource, "location", DatacenterResource+".datacenter_example", "location"),
					resource.TestCheckResourceAttrPair(DBaasMongoClusterResource+"."+DBaaSClusterTestResource, "connections.0.datacenter_id", DatacenterResource+".datacenter_example", "id"),
					resource.TestCheckResourceAttrPair(DBaasMongoClusterResource+"."+DBaaSClusterTestResource, "connections.0.lan_id", LanResource+".lan_example", "id"),
					resource.TestCheckResourceAttr(DBaasMongoClusterResource+"."+DBaaSClusterTestResource, "connections.0.cidr_list.0", "192.168.1.108/24"),
					resource.TestCheckResourceAttr(DBaasMongoClusterResource+"."+DBaaSClusterTestResource, "connections.0.cidr_list.1", "192.168.1.109/24"),
					resource.TestCheckResourceAttr(DBaasMongoClusterResource+"."+DBaaSClusterTestResource, "connections.0.cidr_list.2", "192.168.1.110/24"),
					resource.TestCheckResourceAttr(DBaasMongoClusterResource+"."+DBaaSClusterTestResource, "template_id", "6b78ea06-ee0e-4689-998c-fc9c46e781f6"),
					resource.TestCheckResourceAttr(DBaasMongoClusterResource+"."+DBaaSClusterTestResource, "credentials.0.username", "username"),
					resource.TestCheckResourceAttr(DBaasMongoClusterResource+"."+DBaaSClusterTestResource, "credentials.0.password", "password"),
				),
			},
			{
				Config: testAccDataSourceDBaaSMongoClusterMatchId,
				Check: resource.ComposeTestCheckFunc(
					//resource.TestCheckResourceAttrPair(DataSource+"."+DBaasMongoClusterResource+"."+DBaaSClusterTestDataSourceById, "mongo_cluster", DBaasMongoClusterResource+"."+DBaaSClusterTestResource, "mongo_cluster"),
					resource.TestCheckResourceAttrPair(DataSource+"."+DBaasMongoClusterResource+"."+DBaaSClusterTestDataSourceById, "maintenance_window.day_of_the_week", DBaasMongoClusterResource+"."+DBaaSClusterTestResource, "maintenance_window.day_of_the_week"),
					resource.TestCheckResourceAttrPair(DataSource+"."+DBaasMongoClusterResource+"."+DBaaSClusterTestDataSourceById, "maintenance_window.time", DBaasMongoClusterResource+"."+DBaaSClusterTestResource, "maintenance_window.time"),
					resource.TestCheckResourceAttrPair(DataSource+"."+DBaasMongoClusterResource+"."+DBaaSClusterTestDataSourceById, "mongodb_version", DBaasMongoClusterResource+"."+DBaaSClusterTestResource, "mongodb_version"),
					resource.TestCheckResourceAttrPair(DataSource+"."+DBaasMongoClusterResource+"."+DBaaSClusterTestDataSourceById, "instances", DBaasMongoClusterResource+"."+DBaaSClusterTestResource, "instances"),
					resource.TestCheckResourceAttrPair(DataSource+"."+DBaasMongoClusterResource+"."+DBaaSClusterTestDataSourceById, "display_name", DBaasMongoClusterResource+"."+DBaaSClusterTestResource, "display_name"),
					resource.TestCheckResourceAttrPair(DataSource+"."+DBaasMongoClusterResource+"."+DBaaSClusterTestDataSourceById, "location", DBaasMongoClusterResource+"."+DBaaSClusterTestResource, "location"),
					resource.TestCheckResourceAttrPair(DataSource+"."+DBaasMongoClusterResource+"."+DBaaSClusterTestDataSourceById, "connections.datacenter_id", DBaasMongoClusterResource+"."+DBaaSClusterTestResource, "connections.datacenter_id"),
					resource.TestCheckResourceAttrPair(DataSource+"."+DBaasMongoClusterResource+"."+DBaaSClusterTestDataSourceById, "connections.lan_id", DBaasMongoClusterResource+"."+DBaaSClusterTestResource, "connections.lan_id"),
					resource.TestCheckResourceAttrPair(DataSource+"."+DBaasMongoClusterResource+"."+DBaaSClusterTestDataSourceById, "connections.0.cidr_list.0", DBaasMongoClusterResource+"."+DBaaSClusterTestResource, "connections.0.cidr_list.0"),
					resource.TestCheckResourceAttrPair(DataSource+"."+DBaasMongoClusterResource+"."+DBaaSClusterTestDataSourceById, "connections.0.cidr_list.1", DBaasMongoClusterResource+"."+DBaaSClusterTestResource, "connections.0.cidr_list.1"),
					resource.TestCheckResourceAttrPair(DataSource+"."+DBaasMongoClusterResource+"."+DBaaSClusterTestDataSourceById, "connections.0.cidr_list.2", DBaasMongoClusterResource+"."+DBaaSClusterTestResource, "connections.0.cidr_list.2"),
					resource.TestCheckResourceAttrPair(DataSource+"."+DBaasMongoClusterResource+"."+DBaaSClusterTestDataSourceById, "template_id", DBaasMongoClusterResource+"."+DBaaSClusterTestResource, "template_id"),
					resource.TestCheckResourceAttrPair(DataSource+"."+DBaasMongoClusterResource+"."+DBaaSClusterTestDataSourceById, "connection_string", DBaasMongoClusterResource+"."+DBaaSClusterTestResource, "connection_string"),
					resource.TestCheckResourceAttrPair(DataSource+"."+DBaasMongoClusterResource+"."+DBaaSClusterTestDataSourceById, "credentials.username", DBaasMongoClusterResource+"."+DBaaSClusterTestResource, "credentials.username"),
					resource.TestCheckResourceAttrPair(DataSource+"."+DBaasMongoClusterResource+"."+DBaaSClusterTestDataSourceById, "credentials.password", DBaasMongoClusterResource+"."+DBaaSClusterTestResource, "credentials.password"),
				),
			},
			{
				Config: testAccDataSourceDBaaSMongoClusterMatchName,
				Check: resource.ComposeTestCheckFunc(
					//	resource.TestCheckResourceAttrPair(DataSource+"."+DBaasMongoClusterResource+"."+DBaaSClusterTestDataSourceById, "mongo_cluster", DBaasMongoClusterResource+"."+DBaaSClusterTestResource, "mongo_cluster"),
					resource.TestCheckResourceAttrPair(DataSource+"."+DBaasMongoClusterResource+"."+DBaaSClusterTestDataSourceByName, "maintenance_window.day_of_the_week", DBaasMongoClusterResource+"."+DBaaSClusterTestResource, "maintenance_window.day_of_the_week"),
					resource.TestCheckResourceAttrPair(DataSource+"."+DBaasMongoClusterResource+"."+DBaaSClusterTestDataSourceByName, "maintenance_window.time", DBaasMongoClusterResource+"."+DBaaSClusterTestResource, "maintenance_window.time"),
					resource.TestCheckResourceAttrPair(DataSource+"."+DBaasMongoClusterResource+"."+DBaaSClusterTestDataSourceByName, "mongodb_version", DBaasMongoClusterResource+"."+DBaaSClusterTestResource, "mongodb_version"),
					resource.TestCheckResourceAttrPair(DataSource+"."+DBaasMongoClusterResource+"."+DBaaSClusterTestDataSourceByName, "instances", DBaasMongoClusterResource+"."+DBaaSClusterTestResource, "instances"),
					resource.TestCheckResourceAttrPair(DataSource+"."+DBaasMongoClusterResource+"."+DBaaSClusterTestDataSourceByName, "display_name", DBaasMongoClusterResource+"."+DBaaSClusterTestResource, "display_name"),
					resource.TestCheckResourceAttrPair(DataSource+"."+DBaasMongoClusterResource+"."+DBaaSClusterTestDataSourceByName, "location", DBaasMongoClusterResource+"."+DBaaSClusterTestResource, "location"),
					resource.TestCheckResourceAttrPair(DataSource+"."+DBaasMongoClusterResource+"."+DBaaSClusterTestDataSourceByName, "connections.datacenter_id", DBaasMongoClusterResource+"."+DBaaSClusterTestResource, "connections.datacenter_id"),
					resource.TestCheckResourceAttrPair(DataSource+"."+DBaasMongoClusterResource+"."+DBaaSClusterTestDataSourceByName, "connections.lan_id", DBaasMongoClusterResource+"."+DBaaSClusterTestResource, "connections.lan_id"),
					resource.TestCheckResourceAttrPair(DataSource+"."+DBaasMongoClusterResource+"."+DBaaSClusterTestDataSourceByName, "connections.0.cidr_list", DBaasMongoClusterResource+"."+DBaaSClusterTestResource, "connections.0.cidr_list"),
					resource.TestCheckResourceAttrPair(DataSource+"."+DBaasMongoClusterResource+"."+DBaaSClusterTestDataSourceByName, "template_id", DBaasMongoClusterResource+"."+DBaaSClusterTestResource, "template_id"),
					resource.TestCheckResourceAttrPair(DataSource+"."+DBaasMongoClusterResource+"."+DBaaSClusterTestDataSourceByName, "connection_string", DBaasMongoClusterResource+"."+DBaaSClusterTestResource, "connection_string"),
					resource.TestCheckResourceAttrPair(DataSource+"."+DBaasMongoClusterResource+"."+DBaaSClusterTestDataSourceByName, "credentials.username", DBaasMongoClusterResource+"."+DBaaSClusterTestResource, "credentials.username"),
					resource.TestCheckResourceAttrPair(DataSource+"."+DBaasMongoClusterResource+"."+DBaaSClusterTestDataSourceByName, "credentials.password", DBaasMongoClusterResource+"."+DBaaSClusterTestResource, "credentials.password"),
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
	//client := testAccProvider.Meta().(SdkBundle).DbaasClient
	client := testAccProvider.Meta().(SdkBundle).MongoClient

	ctx, cancel := context.WithTimeout(context.Background(), *resourceDefaultTimeouts.Default)

	if cancel != nil {
		defer cancel()
	}

	for _, rs := range s.RootModule().Resources {
		//if rs.Type != DBaaSBackupsResource {
		//	continue
		//}

		_, apiResponse, err := client.ClustersApi.ClustersFindById(ctx, rs.Primary.ID).Execute()

		if err != nil {
			if apiResponse == nil || apiResponse.StatusCode != 404 {
				return fmt.Errorf("an error occurred while checking the destruction of dbaas mongo cluster %s: %s", rs.Primary.ID, err)
			}
		} else {
			return fmt.Errorf("k8s cluster %s still exists", rs.Primary.ID)
		}

	}

	return nil
}

func testAccCheckDbaasMongoClusterExists(n string, cluster *mongo.ClusterResponse) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(SdkBundle).MongoClient

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

		foundCluster, _, err := client.ClustersApi.ClustersFindById(ctx, rs.Primary.ID).Execute()

		if err != nil {
			return fmt.Errorf("an error occured while fetching k8s Cluster %s: %s", rs.Primary.ID, err)
		}
		if *foundCluster.Id != rs.Primary.ID {
			return fmt.Errorf("record not found")
		}
		cluster = &foundCluster

		return nil
	}
}

const testAccCheckDbaasMongoClusterConfigBasic = `
resource ` + DatacenterResource + ` "datacenter_example" {
  name        = "datacenter_example"
  location    = "de/txl"
  description = "Datacenter for testing dbaas cluster"
}

resource ` + LanResource + ` "lan_example" {
  datacenter_id = ` + DatacenterResource + `.datacenter_example.id 
  public        = false
  name          = "lan_example"
}

resource ` + DBaasMongoClusterResource + ` ` + DBaaSClusterTestResource + ` {
  maintenance_window {
    day_of_the_week  = "Sunday"
    time             = "09:00:00"
  }
  mongodb_version = "5.0"
  instances          = 3
  display_name = "` + DBaaSClusterTestResource + `"
  location = ` + DatacenterResource + `.datacenter_example.location
  connections   {
	datacenter_id   =  ` + DatacenterResource + `.datacenter_example.id 
    lan_id          =  ` + LanResource + `.lan_example.id 
    cidr_list            =  ["192.168.1.108/24", "192.168.1.109/24", "192.168.1.110/24"]
  }
  template_id = "6b78ea06-ee0e-4689-998c-fc9c46e781f6"
  
  credentials {
  	username = "username"
	password = "password"
  }
}
`

const testAccDataSourceDBaaSMongoClusterMatchId = testAccCheckDbaasMongoClusterConfigBasic + `
data ` + DBaasMongoClusterResource + ` ` + DBaaSClusterTestDataSourceById + ` {
  id	= ` + DBaasMongoClusterResource + `.` + DBaaSClusterTestResource + `.id
}
`

const testAccDataSourceDBaaSMongoClusterMatchName = testAccCheckDbaasMongoClusterConfigBasic + `
data ` + DBaasMongoClusterResource + ` ` + DBaaSClusterTestDataSourceByName + ` {
  display_name	= "` + DBaaSClusterTestResource + `"
}
`

const testAccDataSourceDBaaSMongoClusterWrongNameError = testAccCheckDbaasMongoClusterConfigBasic + `
data ` + DBaasMongoClusterResource + ` ` + DBaaSClusterTestDataSourceByName + ` {
  display_name	= "wrong_name"
}
`
