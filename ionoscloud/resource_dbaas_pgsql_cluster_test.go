//go:build all || dbaas
// +build all dbaas

package ionoscloud

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	dbaas "github.com/ionos-cloud/sdk-go-dbaas-postgres"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils"
	"regexp"
	"testing"
)

func TestAccDBaaSPgSqlClusterBasic(t *testing.T) {
	var dbaasCluster dbaas.ClusterResponse

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckDbaasPgSqlClusterDestroyCheck,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckDbaasPgSqlClusterConfigBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDbaasPgSqlClusterExists(DBaaSClusterResource+"."+DBaaSClusterTestResource, &dbaasCluster),
					resource.TestCheckResourceAttr(DBaaSClusterResource+"."+DBaaSClusterTestResource, "postgres_version", "12"),
					resource.TestCheckResourceAttr(DBaaSClusterResource+"."+DBaaSClusterTestResource, "instances", "1"),
					resource.TestCheckResourceAttr(DBaaSClusterResource+"."+DBaaSClusterTestResource, "cores", "1"),
					resource.TestCheckResourceAttr(DBaaSClusterResource+"."+DBaaSClusterTestResource, "ram", "2048"),
					resource.TestCheckResourceAttr(DBaaSClusterResource+"."+DBaaSClusterTestResource, "storage_size", "2048"),
					resource.TestCheckResourceAttr(DBaaSClusterResource+"."+DBaaSClusterTestResource, "storage_type", "HDD"),
					resource.TestCheckResourceAttrPair(DBaaSClusterResource+"."+DBaaSClusterTestResource, "connections.0.datacenter_id", DatacenterResource+".datacenter_example", "id"),
					resource.TestCheckResourceAttrPair(DBaaSClusterResource+"."+DBaaSClusterTestResource, "connections.0.lan_id", LanResource+".lan_example", "id"),
					resource.TestCheckResourceAttr(DBaaSClusterResource+"."+DBaaSClusterTestResource, "connections.0.cidr", "192.168.1.100/24"),
					resource.TestCheckResourceAttrPair(DBaaSClusterResource+"."+DBaaSClusterTestResource, "location", DatacenterResource+".datacenter_example", "location"),
					resource.TestCheckResourceAttr(DBaaSClusterResource+"."+DBaaSClusterTestResource, "backup_location", "de"),
					resource.TestCheckResourceAttr(DBaaSClusterResource+"."+DBaaSClusterTestResource, "display_name", DBaaSClusterTestResource),
					resource.TestCheckResourceAttr(DBaaSClusterResource+"."+DBaaSClusterTestResource, "maintenance_window.0.time", "09:00:00"),
					resource.TestCheckResourceAttr(DBaaSClusterResource+"."+DBaaSClusterTestResource, "maintenance_window.0.day_of_the_week", "Sunday"),
					resource.TestCheckResourceAttr(DBaaSClusterResource+"."+DBaaSClusterTestResource, "credentials.0.username", "username"),
					resource.TestCheckResourceAttr(DBaaSClusterResource+"."+DBaaSClusterTestResource, "credentials.0.password", "password"),
					resource.TestCheckResourceAttr(DBaaSClusterResource+"."+DBaaSClusterTestResource, "synchronization_mode", "ASYNCHRONOUS"),
				),
			},
			{
				Config: testAccDataSourceDBaaSPgSqlClusterMatchId,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(DataSource+"."+DBaaSClusterResource+"."+DBaaSClusterTestDataSourceById, "display_name", DBaaSClusterResource+"."+DBaaSClusterTestResource, "display_name"),
					resource.TestCheckResourceAttrPair(DataSource+"."+DBaaSClusterResource+"."+DBaaSClusterTestDataSourceById, "instances", DBaaSClusterResource+"."+DBaaSClusterTestResource, "instances"),
					resource.TestCheckResourceAttrPair(DataSource+"."+DBaaSClusterResource+"."+DBaaSClusterTestDataSourceById, "cores", DBaaSClusterResource+"."+DBaaSClusterTestResource, "cores"),
					resource.TestCheckResourceAttrPair(DataSource+"."+DBaaSClusterResource+"."+DBaaSClusterTestDataSourceById, "ram", DBaaSClusterResource+"."+DBaaSClusterTestResource, "ram"),
					resource.TestCheckResourceAttrPair(DataSource+"."+DBaaSClusterResource+"."+DBaaSClusterTestDataSourceById, "storage_size", DBaaSClusterResource+"."+DBaaSClusterTestResource, "storage_size"),
					resource.TestCheckResourceAttrPair(DataSource+"."+DBaaSClusterResource+"."+DBaaSClusterTestDataSourceById, "storage_type", DBaaSClusterResource+"."+DBaaSClusterTestResource, "storage_type"),
					resource.TestCheckResourceAttrPair(DataSource+"."+DBaaSClusterResource+"."+DBaaSClusterTestDataSourceById, "connections.datacenter_id", DBaaSClusterResource+"."+DBaaSClusterTestResource, "connections.datacenter_id"),
					resource.TestCheckResourceAttrPair(DataSource+"."+DBaaSClusterResource+"."+DBaaSClusterTestDataSourceById, "connections.lan_id", DBaaSClusterResource+"."+DBaaSClusterTestResource, "connections.lan_id"),
					resource.TestCheckResourceAttrPair(DataSource+"."+DBaaSClusterResource+"."+DBaaSClusterTestDataSourceById, "connections.cidr", DBaaSClusterResource+"."+DBaaSClusterTestResource, "connections.cidr"),
					resource.TestCheckResourceAttrPair(DataSource+"."+DBaaSClusterResource+"."+DBaaSClusterTestDataSourceById, "location", DBaaSClusterResource+"."+DBaaSClusterTestResource, "location"),
					resource.TestCheckResourceAttrPair(DataSource+"."+DBaaSClusterResource+"."+DBaaSClusterTestDataSourceById, "backup_location", DBaaSClusterResource+"."+DBaaSClusterTestResource, "backup_location"),
					resource.TestCheckResourceAttrPair(DataSource+"."+DBaaSClusterResource+"."+DBaaSClusterTestDataSourceById, "display_name", DBaaSClusterResource+"."+DBaaSClusterTestResource, "display_name"),
					resource.TestCheckResourceAttrPair(DataSource+"."+DBaaSClusterResource+"."+DBaaSClusterTestDataSourceById, "maintenance_window.day_of_the_week", DBaaSClusterResource+"."+DBaaSClusterTestResource, "maintenance_window.day_of_the_week"),
					resource.TestCheckResourceAttrPair(DataSource+"."+DBaaSClusterResource+"."+DBaaSClusterTestDataSourceById, "maintenance_window.time", DBaaSClusterResource+"."+DBaaSClusterTestResource, "maintenance_window.time"),
					resource.TestCheckResourceAttrPair(DataSource+"."+DBaaSClusterResource+"."+DBaaSClusterTestDataSourceById, "credentials.username", DBaaSClusterResource+"."+DBaaSClusterTestResource, "credentials.username"),
					resource.TestCheckResourceAttrPair(DataSource+"."+DBaaSClusterResource+"."+DBaaSClusterTestDataSourceById, "credentials.password", DBaaSClusterResource+"."+DBaaSClusterTestResource, "credentials.password"),
				),
			},
			{
				Config: testAccDataSourceDBaaSPgSqlClusterMatchName,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(DataSource+"."+DBaaSClusterResource+"."+DBaaSClusterTestDataSourceByName, "display_name", DBaaSClusterResource+"."+DBaaSClusterTestResource, "display_name"),
					resource.TestCheckResourceAttrPair(DataSource+"."+DBaaSClusterResource+"."+DBaaSClusterTestDataSourceByName, "instances", DBaaSClusterResource+"."+DBaaSClusterTestResource, "instances"),
					resource.TestCheckResourceAttrPair(DataSource+"."+DBaaSClusterResource+"."+DBaaSClusterTestDataSourceByName, "cores", DBaaSClusterResource+"."+DBaaSClusterTestResource, "cores"),
					resource.TestCheckResourceAttrPair(DataSource+"."+DBaaSClusterResource+"."+DBaaSClusterTestDataSourceByName, "ram", DBaaSClusterResource+"."+DBaaSClusterTestResource, "ram"),
					resource.TestCheckResourceAttrPair(DataSource+"."+DBaaSClusterResource+"."+DBaaSClusterTestDataSourceByName, "storage_size", DBaaSClusterResource+"."+DBaaSClusterTestResource, "storage_size"),
					resource.TestCheckResourceAttrPair(DataSource+"."+DBaaSClusterResource+"."+DBaaSClusterTestDataSourceByName, "storage_type", DBaaSClusterResource+"."+DBaaSClusterTestResource, "storage_type"),
					resource.TestCheckResourceAttrPair(DataSource+"."+DBaaSClusterResource+"."+DBaaSClusterTestDataSourceByName, "connections.datacenter_id", DBaaSClusterResource+"."+DBaaSClusterTestResource, "connections.datacenter_id"),
					resource.TestCheckResourceAttrPair(DataSource+"."+DBaaSClusterResource+"."+DBaaSClusterTestDataSourceByName, "connections.lan_id", DBaaSClusterResource+"."+DBaaSClusterTestResource, "connections.lan_id"),
					resource.TestCheckResourceAttrPair(DataSource+"."+DBaaSClusterResource+"."+DBaaSClusterTestDataSourceByName, "connections.cidr", DBaaSClusterResource+"."+DBaaSClusterTestResource, "connections.cidr"),
					resource.TestCheckResourceAttrPair(DataSource+"."+DBaaSClusterResource+"."+DBaaSClusterTestDataSourceByName, "location", DBaaSClusterResource+"."+DBaaSClusterTestResource, "location"),
					resource.TestCheckResourceAttrPair(DataSource+"."+DBaaSClusterResource+"."+DBaaSClusterTestDataSourceByName, "backup_location", DBaaSClusterResource+"."+DBaaSClusterTestResource, "backup_location"),
					resource.TestCheckResourceAttrPair(DataSource+"."+DBaaSClusterResource+"."+DBaaSClusterTestDataSourceByName, "display_name", DBaaSClusterResource+"."+DBaaSClusterTestResource, "display_name"),
					resource.TestCheckResourceAttrPair(DataSource+"."+DBaaSClusterResource+"."+DBaaSClusterTestDataSourceByName, "maintenance_window.day_of_the_week", DBaaSClusterResource+"."+DBaaSClusterTestResource, "maintenance_window.day_of_the_week"),
					resource.TestCheckResourceAttrPair(DataSource+"."+DBaaSClusterResource+"."+DBaaSClusterTestDataSourceByName, "maintenance_window.time", DBaaSClusterResource+"."+DBaaSClusterTestResource, "maintenance_window.time"),
					resource.TestCheckResourceAttrPair(DataSource+"."+DBaaSClusterResource+"."+DBaaSClusterTestDataSourceByName, "credentials.username", DBaaSClusterResource+"."+DBaaSClusterTestResource, "credentials.username"),
					resource.TestCheckResourceAttrPair(DataSource+"."+DBaaSClusterResource+"."+DBaaSClusterTestDataSourceByName, "credentials.password", DBaaSClusterResource+"."+DBaaSClusterTestResource, "credentials.password"),
				),
			},
			{
				Config:      testAccDataSourceDBaaSPgSqlClusterWrongNameError,
				ExpectError: regexp.MustCompile("no DBaaS cluster found with the specified name"),
			},
			{
				Config: testAccDataSourceDbaasPgSqlClusterBackups,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(DataSource+"."+DBaaSBackupsResource+"."+DBaaSBackupsTest, "cluster_backups.0.cluster_id", DataSource+"."+DBaaSBackupsResource+"."+DBaaSBackupsTest, "cluster_id"),
					utils.TestNotEmptySlice(DBaaSBackupsResource, "cluster_backups.#"),
				),
			},
			{
				Config: testAccDataSourceDbaasPgSqlVersionsByClusterId,
				Check: resource.ComposeTestCheckFunc(
					utils.TestNotEmptySlice(DBaaSVersionsResource, "postgres_versions.#"),
				),
			},
			{
				Config: testAccDataSourceDbaasPgSqlAllVersions,
				Check: resource.ComposeTestCheckFunc(
					utils.TestNotEmptySlice(DBaaSVersionsResource, "postgres_versions.#"),
				),
			},
			{
				Config: testAccCheckDbaasPgSqlClusterConfigUpdate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDbaasPgSqlClusterExists(DBaaSClusterResource+"."+DBaaSClusterTestResource, &dbaasCluster),
					resource.TestCheckResourceAttr(DBaaSClusterResource+"."+DBaaSClusterTestResource, "postgres_version", "12"),
					resource.TestCheckResourceAttr(DBaaSClusterResource+"."+DBaaSClusterTestResource, "instances", "2"),
					resource.TestCheckResourceAttr(DBaaSClusterResource+"."+DBaaSClusterTestResource, "cores", "2"),
					resource.TestCheckResourceAttr(DBaaSClusterResource+"."+DBaaSClusterTestResource, "ram", "3072"),
					resource.TestCheckResourceAttr(DBaaSClusterResource+"."+DBaaSClusterTestResource, "storage_size", "3072"),
					resource.TestCheckResourceAttr(DBaaSClusterResource+"."+DBaaSClusterTestResource, "storage_type", "HDD"),
					resource.TestCheckResourceAttrPair(DBaaSClusterResource+"."+DBaaSClusterTestResource, "connections.0.datacenter_id", DatacenterResource+".datacenter_example_update", "id"),
					resource.TestCheckResourceAttrPair(DBaaSClusterResource+"."+DBaaSClusterTestResource, "connections.0.lan_id", LanResource+".lan_example_update", "id"),
					resource.TestCheckResourceAttr(DBaaSClusterResource+"."+DBaaSClusterTestResource, "connections.0.cidr", "192.168.1.101/24"),
					resource.TestCheckResourceAttrPair(DBaaSClusterResource+"."+DBaaSClusterTestResource, "location", DatacenterResource+".datacenter_example_update", "location"),
					resource.TestCheckResourceAttr(DBaaSClusterResource+"."+DBaaSClusterTestResource, "backup_location", "de"),
					resource.TestCheckResourceAttr(DBaaSClusterResource+"."+DBaaSClusterTestResource, "display_name", UpdatedResources),
					resource.TestCheckResourceAttr(DBaaSClusterResource+"."+DBaaSClusterTestResource, "maintenance_window.0.time", "10:00:00"),
					resource.TestCheckResourceAttr(DBaaSClusterResource+"."+DBaaSClusterTestResource, "maintenance_window.0.day_of_the_week", "Saturday"),
					resource.TestCheckResourceAttr(DBaaSClusterResource+"."+DBaaSClusterTestResource, "credentials.0.username", "username"),
					resource.TestCheckResourceAttr(DBaaSClusterResource+"."+DBaaSClusterTestResource, "credentials.0.password", "password"),
					resource.TestCheckResourceAttr(DBaaSClusterResource+"."+DBaaSClusterTestResource, "synchronization_mode", "ASYNCHRONOUS"),
				),
			},
			{
				Config: testAccCheckDbaasPgSqlClusterConfigUpdateRemoveConnections,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDbaasPgSqlClusterExists(DBaaSClusterResource+"."+DBaaSClusterTestResource, &dbaasCluster),
					resource.TestCheckResourceAttr(DBaaSClusterResource+"."+DBaaSClusterTestResource, "postgres_version", "12"),
					resource.TestCheckResourceAttr(DBaaSClusterResource+"."+DBaaSClusterTestResource, "instances", "2"),
					resource.TestCheckResourceAttr(DBaaSClusterResource+"."+DBaaSClusterTestResource, "cores", "2"),
					resource.TestCheckResourceAttr(DBaaSClusterResource+"."+DBaaSClusterTestResource, "ram", "3072"),
					resource.TestCheckResourceAttr(DBaaSClusterResource+"."+DBaaSClusterTestResource, "storage_size", "3072"),
					resource.TestCheckResourceAttr(DBaaSClusterResource+"."+DBaaSClusterTestResource, "storage_type", "HDD"),
					resource.TestCheckNoResourceAttr(DBaaSClusterResource+"."+DBaaSClusterTestResource, "connections"),
					resource.TestCheckResourceAttrPair(DBaaSClusterResource+"."+DBaaSClusterTestResource, "location", DatacenterResource+".datacenter_example_update", "location"),
					resource.TestCheckResourceAttr(DBaaSClusterResource+"."+DBaaSClusterTestResource, "backup_location", "de"),
					resource.TestCheckResourceAttr(DBaaSClusterResource+"."+DBaaSClusterTestResource, "display_name", UpdatedResources),
					resource.TestCheckResourceAttr(DBaaSClusterResource+"."+DBaaSClusterTestResource, "maintenance_window.0.time", "10:00:00"),
					resource.TestCheckResourceAttr(DBaaSClusterResource+"."+DBaaSClusterTestResource, "maintenance_window.0.day_of_the_week", "Saturday"),
					resource.TestCheckResourceAttr(DBaaSClusterResource+"."+DBaaSClusterTestResource, "credentials.0.username", "username"),
					resource.TestCheckResourceAttr(DBaaSClusterResource+"."+DBaaSClusterTestResource, "credentials.0.password", "password"),
					resource.TestCheckResourceAttr(DBaaSClusterResource+"."+DBaaSClusterTestResource, "synchronization_mode", "ASYNCHRONOUS"),
				),
			},
			{
				//we need this as a separate test because the dbaas cluster needs to be deleted first
				//in order to be able to delete the associated lan after
				// otherwise we get 'Access Denied: Lan 1 is delete-protected by DBAAS'
				Config: testAccCheckDbaasPgSqlClusterConfigUpdateRemoveDBaaS,
			},
		},
	})
}

func TestAccDBaaSPgSqlClusterAdditionalParameters(t *testing.T) {
	var dbaasCluster dbaas.ClusterResponse
	t.Skip()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckDbaasPgSqlClusterDestroyCheck,
		Steps: []resource.TestStep{
			{
				Config: testAccFromBackup,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDbaasPgSqlClusterExists(DBaaSClusterResource+"."+DBaaSClusterTestResource, &dbaasCluster),
					resource.TestCheckResourceAttr(DBaaSClusterResource+"."+DBaaSClusterTestResource, "postgres_version", "12"),
					resource.TestCheckResourceAttr(DBaaSClusterResource+"."+DBaaSClusterTestResource, "instances", "1"),
					resource.TestCheckResourceAttr(DBaaSClusterResource+"."+DBaaSClusterTestResource, "cores", "1"),
					resource.TestCheckResourceAttr(DBaaSClusterResource+"."+DBaaSClusterTestResource, "ram", "2048"),
					resource.TestCheckResourceAttr(DBaaSClusterResource+"."+DBaaSClusterTestResource, "storage_size", "2048"),
					resource.TestCheckResourceAttr(DBaaSClusterResource+"."+DBaaSClusterTestResource, "storage_type", "HDD"),
					resource.TestCheckResourceAttrPair(DBaaSClusterResource+"."+DBaaSClusterTestResource, "connections.0.datacenter_id", DatacenterResource+".datacenter_example", "id"),
					resource.TestCheckResourceAttrPair(DBaaSClusterResource+"."+DBaaSClusterTestResource, "connections.0.lan_id", LanResource+".lan_example", "id"),
					resource.TestCheckResourceAttr(DBaaSClusterResource+"."+DBaaSClusterTestResource, "connections.0.cidr", "192.168.1.100/24"),
					resource.TestCheckResourceAttrPair(DBaaSClusterResource+"."+DBaaSClusterTestResource, "location", DatacenterResource+".datacenter_example", "location"),
					resource.TestCheckResourceAttr(DBaaSClusterResource+"."+DBaaSClusterTestResource, "display_name", DBaaSClusterTestResource),
					resource.TestCheckResourceAttr(DBaaSClusterResource+"."+DBaaSClusterTestResource, "maintenance_window.0.time", "09:00:00"),
					resource.TestCheckResourceAttr(DBaaSClusterResource+"."+DBaaSClusterTestResource, "maintenance_window.0.day_of_the_week", "Sunday"),
					resource.TestCheckResourceAttr(DBaaSClusterResource+"."+DBaaSClusterTestResource, "credentials.0.username", "username"),
					resource.TestCheckResourceAttr(DBaaSClusterResource+"."+DBaaSClusterTestResource, "credentials.0.password", "password"),
					resource.TestCheckResourceAttr(DBaaSClusterResource+"."+DBaaSClusterTestResource, "synchronization_mode", "ASYNCHRONOUS")),
			},
		},
	})
}

func testAccCheckDbaasPgSqlClusterDestroyCheck(s *terraform.State) error {
	client := testAccProvider.Meta().(SdkBundle).DbaasClient

	ctx, cancel := context.WithTimeout(context.Background(), *resourceDefaultTimeouts.Default)

	if cancel != nil {
		defer cancel()
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != DBaaSBackupsResource {
			continue
		}

		_, apiResponse, err := client.ClustersApi.ClustersFindById(ctx, rs.Primary.ID).Execute()

		if err != nil {
			if apiResponse == nil || apiResponse.StatusCode != 404 {
				return fmt.Errorf("an error occurred while checking the destruction of dbaas cluster %s: %s", rs.Primary.ID, err)
			}
		} else {
			return fmt.Errorf("k8s cluster %s still exists", rs.Primary.ID)
		}

	}

	return nil
}

func testAccCheckDbaasPgSqlClusterExists(n string, cluster *dbaas.ClusterResponse) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(SdkBundle).DbaasClient

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

const testAccCheckDbaasPgSqlClusterConfigBasic = `
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

resource ` + DBaaSClusterResource + ` ` + DBaaSClusterTestResource + ` {
  postgres_version   = 12
  instances          = 1
  cores              = 1
  ram                = 2048
  storage_size       = 2048
  storage_type       = "HDD"
  connections   {
	datacenter_id   =  ` + DatacenterResource + `.datacenter_example.id 
    lan_id          =  ` + LanResource + `.lan_example.id 
    cidr            =  "192.168.1.100/24"
  }
  location = ` + DatacenterResource + `.datacenter_example.location
  backup_location = "de"
  display_name = "` + DBaaSClusterTestResource + `"
  maintenance_window {
    day_of_the_week  = "Sunday"
    time             = "09:00:00"
  }
  credentials {
  	username = "username"
	password = "password"
  }
  synchronization_mode = "ASYNCHRONOUS"
}
`

const testAccCheckDbaasPgSqlClusterConfigUpdate = `
resource ` + DatacenterResource + ` "datacenter_example" {
  name        = "datacenter_example"
  location    = "de/txl"
  description = "Datacenter for testing dbaas cluster"
}

resource ` + DatacenterResource + ` "datacenter_example_update" {
  name        = "datacenter_example_update"
  location    = "de/txl"
  description = "Datacenter for testing dbaas cluster"
}

resource ` + LanResource + ` "lan_example" {
  datacenter_id = ` + DatacenterResource + `.datacenter_example.id 
  public        = false
  name          = "lan_example"
}

resource ` + LanResource + ` "lan_example_update" {
  datacenter_id = ` + DatacenterResource + `.datacenter_example_update.id 
  public        = false
  name          = "lan_example_update"
}


resource ` + DBaaSClusterResource + ` ` + DBaaSClusterTestResource + ` {
  postgres_version   = 12
  instances          = 2
  cores              = 2
  ram                = 3072
  storage_size       = 3072
  storage_type       = "HDD"
  connections   {
	datacenter_id   =  ` + DatacenterResource + `.datacenter_example_update.id 
    lan_id          =  ` + LanResource + `.lan_example_update.id     
    cidr            =  "192.168.1.101/24"
  }
  location = ` + DatacenterResource + `.datacenter_example_update.location
  backup_location = "de"
  display_name = "` + UpdatedResources + `"
  maintenance_window {
    day_of_the_week = "Saturday"
    time            = "10:00:00"
  }
  credentials {
  	username = "username"
	password = "password"
  }
  synchronization_mode = "ASYNCHRONOUS"
}
`

const testAccCheckDbaasPgSqlClusterConfigUpdateRemoveConnections = `
resource ` + DatacenterResource + ` "datacenter_example" {
  name        = "datacenter_example"
  location    = "de/txl"
  description = "Datacenter for testing dbaas cluster"
}

resource ` + DatacenterResource + ` "datacenter_example_update" {
  name        = "datacenter_example_update"
  location    = "de/txl"
  description = "Datacenter for testing dbaas cluster"
}

resource ` + LanResource + ` "lan_example" {
  datacenter_id = ` + DatacenterResource + `.datacenter_example.id 
  public        = false
  name          = "lan_example"
}

resource ` + LanResource + ` "lan_example_update" {
  datacenter_id = ` + DatacenterResource + `.datacenter_example_update.id 
  public        = false
  name          = "lan_example_update"
}


resource ` + DBaaSClusterResource + ` ` + DBaaSClusterTestResource + ` {
  postgres_version   = 12
  instances          = 2
  cores              = 2
  ram                = 3072
  storage_size       = 3072
  storage_type       = "HDD"
  location = ` + DatacenterResource + `.datacenter_example_update.location
  backup_location = "de"
  display_name = "` + UpdatedResources + `"
  maintenance_window {
    day_of_the_week = "Saturday"
    time            = "10:00:00"
  }
  credentials {
  	username = "username"
	password = "password"
  }
  synchronization_mode = "ASYNCHRONOUS"
}
`
const testAccCheckDbaasPgSqlClusterConfigUpdateRemoveDBaaS = `
resource ` + DatacenterResource + ` "datacenter_example" {
  name        = "datacenter_example"
  location    = "de/txl"
  description = "Datacenter for testing dbaas cluster"
}

resource ` + DatacenterResource + ` "datacenter_example_update" {
  name        = "datacenter_example_update"
  location    = "de/txl"
  description = "Datacenter for testing dbaas cluster"
}

resource ` + LanResource + ` "lan_example" {
  datacenter_id = ` + DatacenterResource + `.datacenter_example.id 
  public        = false
  name          = "lan_example"
}

resource ` + LanResource + ` "lan_example_update" {
  datacenter_id = ` + DatacenterResource + `.datacenter_example_update.id 
  public        = false
  name          = "lan_example_update"
}
`

const testAccFromBackup = `
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

resource ` + DBaaSClusterResource + ` ` + DBaaSClusterTestResource + ` {
  postgres_version   = 12
  instances          = 1
  cores              = 1
  ram                = 2048
  storage_size       = 2048
  storage_type       = "HDD"
  display_name = "` + DBaaSClusterTestResource + `"
  connections   {
	datacenter_id   =  ` + DatacenterResource + `.datacenter_example.id 
    lan_id          =  ` + LanResource + `.lan_example.id 
    cidr            =  "192.168.1.100/24"
  }
  location = ` + DatacenterResource + `.datacenter_example.location
  maintenance_window {
    day_of_the_week = "Sunday"
    time            = "09:00:00"
  }
  credentials {
  	username = "username"
	password = "password"
  }
  synchronization_mode = "ASYNCHRONOUS"
  from_backup {
	backup_id = "f767c6e5-747c-11ec-9bb6-4aa52b3d55f1-4oymiqu-12"
    recovery_target_time = "2022-01-13T16:27:42Z"
  }
}`

const testAccDataSourceDBaaSPgSqlClusterMatchId = testAccCheckDbaasPgSqlClusterConfigBasic + `
data ` + DBaaSClusterResource + ` ` + DBaaSClusterTestDataSourceById + ` {
  id	= ` + DBaaSClusterResource + `.` + DBaaSClusterTestResource + `.id
}
`

const testAccDataSourceDBaaSPgSqlClusterMatchName = testAccCheckDbaasPgSqlClusterConfigBasic + `
data ` + DBaaSClusterResource + ` ` + DBaaSClusterTestDataSourceByName + ` {
  display_name	= "` + DBaaSClusterTestResource + `"
}
`

const testAccDataSourceDBaaSPgSqlClusterWrongNameError = testAccCheckDbaasPgSqlClusterConfigBasic + `
data ` + DBaaSClusterResource + ` ` + DBaaSClusterTestDataSourceByName + ` {
  display_name	= "wrong_name"
}
`

const testAccDataSourceDbaasPgSqlClusterBackups = testAccCheckDbaasPgSqlClusterConfigBasic + `
data ` + DBaaSBackupsResource + ` ` + DBaaSBackupsTest + ` {
	cluster_id = ` + DBaaSClusterResource + `.` + DBaaSClusterTestResource + `.id
}
`

const testAccDataSourceDbaasPgSqlAllVersions = testAccCheckDbaasPgSqlClusterConfigBasic + `
data ` + DBaaSVersionsResource + ` ` + DBaaSVersionsTest + ` {
}
`

const testAccDataSourceDbaasPgSqlVersionsByClusterId = testAccCheckDbaasPgSqlClusterConfigBasic + `
data ` + DBaaSVersionsResource + ` ` + DBaaSVersionsTest + ` {
	cluster_id = ` + DBaaSClusterResource + `.` + DBaaSClusterTestResource + `.id
}
`
