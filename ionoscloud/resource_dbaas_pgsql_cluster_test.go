//go:build all || dbaas
// +build all dbaas

package ionoscloud

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	psql "github.com/ionos-cloud/sdk-go-dbaas-postgres"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils"
	"regexp"
	"testing"
)

func TestAccDBaaSPgSqlClusterBasic(t *testing.T) {
	var dbaasCluster psql.ClusterResponse

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
					testAccCheckDbaasPgSqlClusterExists(PsqlClusterResource+"."+DBaaSClusterTestResource, &dbaasCluster),
					resource.TestCheckResourceAttr(PsqlClusterResource+"."+DBaaSClusterTestResource, "postgres_version", "12"),
					resource.TestCheckResourceAttr(PsqlClusterResource+"."+DBaaSClusterTestResource, "instances", "1"),
					resource.TestCheckResourceAttr(PsqlClusterResource+"."+DBaaSClusterTestResource, "cores", "1"),
					resource.TestCheckResourceAttr(PsqlClusterResource+"."+DBaaSClusterTestResource, "ram", "2048"),
					resource.TestCheckResourceAttr(PsqlClusterResource+"."+DBaaSClusterTestResource, "storage_size", "2048"),
					resource.TestCheckResourceAttr(PsqlClusterResource+"."+DBaaSClusterTestResource, "storage_type", "HDD"),
					resource.TestCheckResourceAttrPair(PsqlClusterResource+"."+DBaaSClusterTestResource, "connections.0.datacenter_id", DatacenterResource+".datacenter_example", "id"),
					resource.TestCheckResourceAttrPair(PsqlClusterResource+"."+DBaaSClusterTestResource, "connections.0.lan_id", LanResource+".lan_example", "id"),
					resource.TestCheckResourceAttr(PsqlClusterResource+"."+DBaaSClusterTestResource, "connections.0.cidr", "192.168.1.100/24"),
					resource.TestCheckResourceAttrPair(PsqlClusterResource+"."+DBaaSClusterTestResource, "location", DatacenterResource+".datacenter_example", "location"),
					resource.TestCheckResourceAttr(PsqlClusterResource+"."+DBaaSClusterTestResource, "backup_location", "de"),
					resource.TestCheckResourceAttr(PsqlClusterResource+"."+DBaaSClusterTestResource, "display_name", DBaaSClusterTestResource),
					resource.TestCheckResourceAttr(PsqlClusterResource+"."+DBaaSClusterTestResource, "maintenance_window.0.time", "09:00:00"),
					resource.TestCheckResourceAttr(PsqlClusterResource+"."+DBaaSClusterTestResource, "maintenance_window.0.day_of_the_week", "Sunday"),
					resource.TestCheckResourceAttr(PsqlClusterResource+"."+DBaaSClusterTestResource, "credentials.0.username", "username"),
					resource.TestCheckResourceAttr(PsqlClusterResource+"."+DBaaSClusterTestResource, "credentials.0.password", "password"),
					resource.TestCheckResourceAttr(PsqlClusterResource+"."+DBaaSClusterTestResource, "synchronization_mode", "ASYNCHRONOUS"),
				),
			},
			{
				Config: testAccDataSourceDBaaSPgSqlClusterMatchId,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(DataSource+"."+PsqlClusterResource+"."+DBaaSClusterTestDataSourceById, "display_name", PsqlClusterResource+"."+DBaaSClusterTestResource, "display_name"),
					resource.TestCheckResourceAttrPair(DataSource+"."+PsqlClusterResource+"."+DBaaSClusterTestDataSourceById, "instances", PsqlClusterResource+"."+DBaaSClusterTestResource, "instances"),
					resource.TestCheckResourceAttrPair(DataSource+"."+PsqlClusterResource+"."+DBaaSClusterTestDataSourceById, "cores", PsqlClusterResource+"."+DBaaSClusterTestResource, "cores"),
					resource.TestCheckResourceAttrPair(DataSource+"."+PsqlClusterResource+"."+DBaaSClusterTestDataSourceById, "ram", PsqlClusterResource+"."+DBaaSClusterTestResource, "ram"),
					resource.TestCheckResourceAttrPair(DataSource+"."+PsqlClusterResource+"."+DBaaSClusterTestDataSourceById, "storage_size", PsqlClusterResource+"."+DBaaSClusterTestResource, "storage_size"),
					resource.TestCheckResourceAttrPair(DataSource+"."+PsqlClusterResource+"."+DBaaSClusterTestDataSourceById, "storage_type", PsqlClusterResource+"."+DBaaSClusterTestResource, "storage_type"),
					resource.TestCheckResourceAttrPair(DataSource+"."+PsqlClusterResource+"."+DBaaSClusterTestDataSourceById, "connections.datacenter_id", PsqlClusterResource+"."+DBaaSClusterTestResource, "connections.datacenter_id"),
					resource.TestCheckResourceAttrPair(DataSource+"."+PsqlClusterResource+"."+DBaaSClusterTestDataSourceById, "connections.lan_id", PsqlClusterResource+"."+DBaaSClusterTestResource, "connections.lan_id"),
					resource.TestCheckResourceAttrPair(DataSource+"."+PsqlClusterResource+"."+DBaaSClusterTestDataSourceById, "connections.cidr", PsqlClusterResource+"."+DBaaSClusterTestResource, "connections.cidr"),
					resource.TestCheckResourceAttrPair(DataSource+"."+PsqlClusterResource+"."+DBaaSClusterTestDataSourceById, "location", PsqlClusterResource+"."+DBaaSClusterTestResource, "location"),
					resource.TestCheckResourceAttrPair(DataSource+"."+PsqlClusterResource+"."+DBaaSClusterTestDataSourceById, "backup_location", PsqlClusterResource+"."+DBaaSClusterTestResource, "backup_location"),
					resource.TestCheckResourceAttrPair(DataSource+"."+PsqlClusterResource+"."+DBaaSClusterTestDataSourceById, "display_name", PsqlClusterResource+"."+DBaaSClusterTestResource, "display_name"),
					resource.TestCheckResourceAttrPair(DataSource+"."+PsqlClusterResource+"."+DBaaSClusterTestDataSourceById, "maintenance_window.day_of_the_week", PsqlClusterResource+"."+DBaaSClusterTestResource, "maintenance_window.day_of_the_week"),
					resource.TestCheckResourceAttrPair(DataSource+"."+PsqlClusterResource+"."+DBaaSClusterTestDataSourceById, "maintenance_window.time", PsqlClusterResource+"."+DBaaSClusterTestResource, "maintenance_window.time"),
					resource.TestCheckResourceAttrPair(DataSource+"."+PsqlClusterResource+"."+DBaaSClusterTestDataSourceById, "credentials.username", PsqlClusterResource+"."+DBaaSClusterTestResource, "credentials.username"),
					resource.TestCheckResourceAttrPair(DataSource+"."+PsqlClusterResource+"."+DBaaSClusterTestDataSourceById, "credentials.password", PsqlClusterResource+"."+DBaaSClusterTestResource, "credentials.password"),
				),
			},
			{
				Config: testAccDataSourceDBaaSPgSqlClusterMatchName,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(DataSource+"."+PsqlClusterResource+"."+DBaaSClusterTestDataSourceByName, "display_name", PsqlClusterResource+"."+DBaaSClusterTestResource, "display_name"),
					resource.TestCheckResourceAttrPair(DataSource+"."+PsqlClusterResource+"."+DBaaSClusterTestDataSourceByName, "instances", PsqlClusterResource+"."+DBaaSClusterTestResource, "instances"),
					resource.TestCheckResourceAttrPair(DataSource+"."+PsqlClusterResource+"."+DBaaSClusterTestDataSourceByName, "cores", PsqlClusterResource+"."+DBaaSClusterTestResource, "cores"),
					resource.TestCheckResourceAttrPair(DataSource+"."+PsqlClusterResource+"."+DBaaSClusterTestDataSourceByName, "ram", PsqlClusterResource+"."+DBaaSClusterTestResource, "ram"),
					resource.TestCheckResourceAttrPair(DataSource+"."+PsqlClusterResource+"."+DBaaSClusterTestDataSourceByName, "storage_size", PsqlClusterResource+"."+DBaaSClusterTestResource, "storage_size"),
					resource.TestCheckResourceAttrPair(DataSource+"."+PsqlClusterResource+"."+DBaaSClusterTestDataSourceByName, "storage_type", PsqlClusterResource+"."+DBaaSClusterTestResource, "storage_type"),
					resource.TestCheckResourceAttrPair(DataSource+"."+PsqlClusterResource+"."+DBaaSClusterTestDataSourceByName, "connections.datacenter_id", PsqlClusterResource+"."+DBaaSClusterTestResource, "connections.datacenter_id"),
					resource.TestCheckResourceAttrPair(DataSource+"."+PsqlClusterResource+"."+DBaaSClusterTestDataSourceByName, "connections.lan_id", PsqlClusterResource+"."+DBaaSClusterTestResource, "connections.lan_id"),
					resource.TestCheckResourceAttrPair(DataSource+"."+PsqlClusterResource+"."+DBaaSClusterTestDataSourceByName, "connections.cidr", PsqlClusterResource+"."+DBaaSClusterTestResource, "connections.cidr"),
					resource.TestCheckResourceAttrPair(DataSource+"."+PsqlClusterResource+"."+DBaaSClusterTestDataSourceByName, "location", PsqlClusterResource+"."+DBaaSClusterTestResource, "location"),
					resource.TestCheckResourceAttrPair(DataSource+"."+PsqlClusterResource+"."+DBaaSClusterTestDataSourceByName, "backup_location", PsqlClusterResource+"."+DBaaSClusterTestResource, "backup_location"),
					resource.TestCheckResourceAttrPair(DataSource+"."+PsqlClusterResource+"."+DBaaSClusterTestDataSourceByName, "display_name", PsqlClusterResource+"."+DBaaSClusterTestResource, "display_name"),
					resource.TestCheckResourceAttrPair(DataSource+"."+PsqlClusterResource+"."+DBaaSClusterTestDataSourceByName, "maintenance_window.day_of_the_week", PsqlClusterResource+"."+DBaaSClusterTestResource, "maintenance_window.day_of_the_week"),
					resource.TestCheckResourceAttrPair(DataSource+"."+PsqlClusterResource+"."+DBaaSClusterTestDataSourceByName, "maintenance_window.time", PsqlClusterResource+"."+DBaaSClusterTestResource, "maintenance_window.time"),
					resource.TestCheckResourceAttrPair(DataSource+"."+PsqlClusterResource+"."+DBaaSClusterTestDataSourceByName, "credentials.username", PsqlClusterResource+"."+DBaaSClusterTestResource, "credentials.username"),
					resource.TestCheckResourceAttrPair(DataSource+"."+PsqlClusterResource+"."+DBaaSClusterTestDataSourceByName, "credentials.password", PsqlClusterResource+"."+DBaaSClusterTestResource, "credentials.password"),
				),
			},
			{
				Config:      testAccDataSourceDBaaSPgSqlClusterWrongNameError,
				ExpectError: regexp.MustCompile("no DBaaS cluster found with the specified name"),
			},
			{
				Config: testAccDataSourceDbaasPgSqlClusterBackups,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(DataSource+"."+PsqlBackupsResource+"."+PsqlBackupsTest, "cluster_backups.0.cluster_id", DataSource+"."+PsqlBackupsResource+"."+PsqlBackupsTest, "cluster_id"),
					utils.TestNotEmptySlice(PsqlBackupsResource, "cluster_backups.#"),
				),
			},
			{
				Config: testAccDataSourceDbaasPgSqlVersionsByClusterId,
				Check: resource.ComposeTestCheckFunc(
					utils.TestNotEmptySlice(PsqlVersionsResource, "postgres_versions.#"),
				),
			},
			{
				Config: testAccDataSourceDbaasPgSqlAllVersions,
				Check: resource.ComposeTestCheckFunc(
					utils.TestNotEmptySlice(PsqlVersionsResource, "postgres_versions.#"),
				),
			},
			{
				Config: testAccCheckDbaasPgSqlClusterConfigUpdate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDbaasPgSqlClusterExists(PsqlClusterResource+"."+DBaaSClusterTestResource, &dbaasCluster),
					resource.TestCheckResourceAttr(PsqlClusterResource+"."+DBaaSClusterTestResource, "postgres_version", "12"),
					resource.TestCheckResourceAttr(PsqlClusterResource+"."+DBaaSClusterTestResource, "instances", "2"),
					resource.TestCheckResourceAttr(PsqlClusterResource+"."+DBaaSClusterTestResource, "cores", "2"),
					resource.TestCheckResourceAttr(PsqlClusterResource+"."+DBaaSClusterTestResource, "ram", "3072"),
					resource.TestCheckResourceAttr(PsqlClusterResource+"."+DBaaSClusterTestResource, "storage_size", "3072"),
					resource.TestCheckResourceAttr(PsqlClusterResource+"."+DBaaSClusterTestResource, "storage_type", "HDD"),
					resource.TestCheckResourceAttrPair(PsqlClusterResource+"."+DBaaSClusterTestResource, "connections.0.datacenter_id", DatacenterResource+".datacenter_example_update", "id"),
					resource.TestCheckResourceAttrPair(PsqlClusterResource+"."+DBaaSClusterTestResource, "connections.0.lan_id", LanResource+".lan_example_update", "id"),
					resource.TestCheckResourceAttr(PsqlClusterResource+"."+DBaaSClusterTestResource, "connections.0.cidr", "192.168.1.101/24"),
					resource.TestCheckResourceAttrPair(PsqlClusterResource+"."+DBaaSClusterTestResource, "location", DatacenterResource+".datacenter_example_update", "location"),
					resource.TestCheckResourceAttr(PsqlClusterResource+"."+DBaaSClusterTestResource, "backup_location", "de"),
					resource.TestCheckResourceAttr(PsqlClusterResource+"."+DBaaSClusterTestResource, "display_name", UpdatedResources),
					resource.TestCheckResourceAttr(PsqlClusterResource+"."+DBaaSClusterTestResource, "maintenance_window.0.time", "10:00:00"),
					resource.TestCheckResourceAttr(PsqlClusterResource+"."+DBaaSClusterTestResource, "maintenance_window.0.day_of_the_week", "Saturday"),
					resource.TestCheckResourceAttr(PsqlClusterResource+"."+DBaaSClusterTestResource, "credentials.0.username", "username"),
					resource.TestCheckResourceAttr(PsqlClusterResource+"."+DBaaSClusterTestResource, "credentials.0.password", "password"),
					resource.TestCheckResourceAttr(PsqlClusterResource+"."+DBaaSClusterTestResource, "synchronization_mode", "ASYNCHRONOUS"),
				),
			},
			{
				Config: testAccCheckDbaasPgSqlClusterConfigUpdateRemoveConnections,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDbaasPgSqlClusterExists(PsqlClusterResource+"."+DBaaSClusterTestResource, &dbaasCluster),
					resource.TestCheckResourceAttr(PsqlClusterResource+"."+DBaaSClusterTestResource, "postgres_version", "12"),
					resource.TestCheckResourceAttr(PsqlClusterResource+"."+DBaaSClusterTestResource, "instances", "2"),
					resource.TestCheckResourceAttr(PsqlClusterResource+"."+DBaaSClusterTestResource, "cores", "2"),
					resource.TestCheckResourceAttr(PsqlClusterResource+"."+DBaaSClusterTestResource, "ram", "3072"),
					resource.TestCheckResourceAttr(PsqlClusterResource+"."+DBaaSClusterTestResource, "storage_size", "3072"),
					resource.TestCheckResourceAttr(PsqlClusterResource+"."+DBaaSClusterTestResource, "storage_type", "HDD"),
					resource.TestCheckNoResourceAttr(PsqlClusterResource+"."+DBaaSClusterTestResource, "connections"),
					resource.TestCheckResourceAttrPair(PsqlClusterResource+"."+DBaaSClusterTestResource, "location", DatacenterResource+".datacenter_example_update", "location"),
					resource.TestCheckResourceAttr(PsqlClusterResource+"."+DBaaSClusterTestResource, "backup_location", "de"),
					resource.TestCheckResourceAttr(PsqlClusterResource+"."+DBaaSClusterTestResource, "display_name", UpdatedResources),
					resource.TestCheckResourceAttr(PsqlClusterResource+"."+DBaaSClusterTestResource, "maintenance_window.0.time", "10:00:00"),
					resource.TestCheckResourceAttr(PsqlClusterResource+"."+DBaaSClusterTestResource, "maintenance_window.0.day_of_the_week", "Saturday"),
					resource.TestCheckResourceAttr(PsqlClusterResource+"."+DBaaSClusterTestResource, "credentials.0.username", "username"),
					resource.TestCheckResourceAttr(PsqlClusterResource+"."+DBaaSClusterTestResource, "credentials.0.password", "password"),
					resource.TestCheckResourceAttr(PsqlClusterResource+"."+DBaaSClusterTestResource, "synchronization_mode", "ASYNCHRONOUS"),
				),
			},
			{
				//we need this as a separate test because the psql cluster needs to be deleted first
				//in order to be able to delete the associated lan after
				// otherwise we get 'Access Denied: Lan 1 is delete-protected by DBAAS'
				Config: testAccCheckDbaasPgSqlClusterConfigUpdateRemoveDBaaS,
			},
		},
	})
}

func TestAccDBaaSPgSqlClusterAdditionalParameters(t *testing.T) {
	var dbaasCluster psql.ClusterResponse
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
					testAccCheckDbaasPgSqlClusterExists(PsqlClusterResource+"."+DBaaSClusterTestResource, &dbaasCluster),
					resource.TestCheckResourceAttr(PsqlClusterResource+"."+DBaaSClusterTestResource, "postgres_version", "12"),
					resource.TestCheckResourceAttr(PsqlClusterResource+"."+DBaaSClusterTestResource, "instances", "1"),
					resource.TestCheckResourceAttr(PsqlClusterResource+"."+DBaaSClusterTestResource, "cores", "1"),
					resource.TestCheckResourceAttr(PsqlClusterResource+"."+DBaaSClusterTestResource, "ram", "2048"),
					resource.TestCheckResourceAttr(PsqlClusterResource+"."+DBaaSClusterTestResource, "storage_size", "2048"),
					resource.TestCheckResourceAttr(PsqlClusterResource+"."+DBaaSClusterTestResource, "storage_type", "HDD"),
					resource.TestCheckResourceAttrPair(PsqlClusterResource+"."+DBaaSClusterTestResource, "connections.0.datacenter_id", DatacenterResource+".datacenter_example", "id"),
					resource.TestCheckResourceAttrPair(PsqlClusterResource+"."+DBaaSClusterTestResource, "connections.0.lan_id", LanResource+".lan_example", "id"),
					resource.TestCheckResourceAttr(PsqlClusterResource+"."+DBaaSClusterTestResource, "connections.0.cidr", "192.168.1.100/24"),
					resource.TestCheckResourceAttrPair(PsqlClusterResource+"."+DBaaSClusterTestResource, "location", DatacenterResource+".datacenter_example", "location"),
					resource.TestCheckResourceAttr(PsqlClusterResource+"."+DBaaSClusterTestResource, "display_name", DBaaSClusterTestResource),
					resource.TestCheckResourceAttr(PsqlClusterResource+"."+DBaaSClusterTestResource, "maintenance_window.0.time", "09:00:00"),
					resource.TestCheckResourceAttr(PsqlClusterResource+"."+DBaaSClusterTestResource, "maintenance_window.0.day_of_the_week", "Sunday"),
					resource.TestCheckResourceAttr(PsqlClusterResource+"."+DBaaSClusterTestResource, "credentials.0.username", "username"),
					resource.TestCheckResourceAttr(PsqlClusterResource+"."+DBaaSClusterTestResource, "credentials.0.password", "password"),
					resource.TestCheckResourceAttr(PsqlClusterResource+"."+DBaaSClusterTestResource, "synchronization_mode", "ASYNCHRONOUS")),
			},
		},
	})
}

func testAccCheckDbaasPgSqlClusterDestroyCheck(s *terraform.State) error {
	client := testAccProvider.Meta().(SdkBundle).PsqlClient

	ctx, cancel := context.WithTimeout(context.Background(), *resourceDefaultTimeouts.Default)

	if cancel != nil {
		defer cancel()
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != PsqlBackupsResource {
			continue
		}

		_, apiResponse, err := client.ClustersApi.ClustersFindById(ctx, rs.Primary.ID).Execute()

		if err != nil {
			if apiResponse == nil || apiResponse.StatusCode != 404 {
				return fmt.Errorf("an error occurred while checking the destruction of psql cluster %s: %s", rs.Primary.ID, err)
			}
		} else {
			return fmt.Errorf("k8s cluster %s still exists", rs.Primary.ID)
		}

	}

	return nil
}

func testAccCheckDbaasPgSqlClusterExists(n string, cluster *psql.ClusterResponse) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(SdkBundle).PsqlClient

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
  description = "Datacenter for testing psql cluster"
}

resource ` + LanResource + ` "lan_example" {
  datacenter_id = ` + DatacenterResource + `.datacenter_example.id 
  public        = false
  name          = "lan_example"
}

resource ` + PsqlClusterResource + ` ` + DBaaSClusterTestResource + ` {
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
  description = "Datacenter for testing psql cluster"
}

resource ` + DatacenterResource + ` "datacenter_example_update" {
  name        = "datacenter_example_update"
  location    = "de/txl"
  description = "Datacenter for testing psql cluster"
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


resource ` + PsqlClusterResource + ` ` + DBaaSClusterTestResource + ` {
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
  description = "Datacenter for testing psql cluster"
}

resource ` + DatacenterResource + ` "datacenter_example_update" {
  name        = "datacenter_example_update"
  location    = "de/txl"
  description = "Datacenter for testing psql cluster"
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


resource ` + PsqlClusterResource + ` ` + DBaaSClusterTestResource + ` {
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
  description = "Datacenter for testing psql cluster"
}

resource ` + DatacenterResource + ` "datacenter_example_update" {
  name        = "datacenter_example_update"
  location    = "de/txl"
  description = "Datacenter for testing psql cluster"
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
  description = "Datacenter for testing psql cluster"
}

resource ` + LanResource + ` "lan_example" {
  datacenter_id = ` + DatacenterResource + `.datacenter_example.id 
  public        = false
  name          = "lan_example"
}

resource ` + PsqlClusterResource + ` ` + DBaaSClusterTestResource + ` {
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
data ` + PsqlClusterResource + ` ` + DBaaSClusterTestDataSourceById + ` {
  id	= ` + PsqlClusterResource + `.` + DBaaSClusterTestResource + `.id
}
`

const testAccDataSourceDBaaSPgSqlClusterMatchName = testAccCheckDbaasPgSqlClusterConfigBasic + `
data ` + PsqlClusterResource + ` ` + DBaaSClusterTestDataSourceByName + ` {
  display_name	= "` + DBaaSClusterTestResource + `"
}
`

const testAccDataSourceDBaaSPgSqlClusterWrongNameError = testAccCheckDbaasPgSqlClusterConfigBasic + `
data ` + PsqlClusterResource + ` ` + DBaaSClusterTestDataSourceByName + ` {
  display_name	= "wrong_name"
}
`

const testAccDataSourceDbaasPgSqlClusterBackups = testAccCheckDbaasPgSqlClusterConfigBasic + `
data ` + PsqlBackupsResource + ` ` + PsqlBackupsTest + ` {
	cluster_id = ` + PsqlClusterResource + `.` + DBaaSClusterTestResource + `.id
}
`

const testAccDataSourceDbaasPgSqlAllVersions = testAccCheckDbaasPgSqlClusterConfigBasic + `
data ` + PsqlVersionsResource + ` ` + PsqlVersionsTest + ` {
}
`

const testAccDataSourceDbaasPgSqlVersionsByClusterId = testAccCheckDbaasPgSqlClusterConfigBasic + `
data ` + PsqlVersionsResource + ` ` + PsqlVersionsTest + ` {
	cluster_id = ` + PsqlClusterResource + `.` + DBaaSClusterTestResource + `.id
}
`
