package ionoscloud

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	dbaas "github.com/ionos-cloud/sdk-go-dbaas-postgres"
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
					resource.TestCheckResourceAttr(DBaaSClusterResource+"."+DBaaSClusterTestResource, "cores", "3"),
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
					resource.TestCheckResourceAttr(DBaaSClusterResource+"."+DBaaSClusterTestResource, "synchronization_mode", "ASYNCHRONOUS"),
				),
			},
			{
				Config: testAccCheckDbaasPgSqlClusterConfigUpdate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDbaasPgSqlClusterExists(DBaaSClusterResource+"."+DBaaSClusterTestResource, &dbaasCluster),
					resource.TestCheckResourceAttr(DBaaSClusterResource+"."+DBaaSClusterTestResource, "postgres_version", "13"),
					resource.TestCheckResourceAttr(DBaaSClusterResource+"."+DBaaSClusterTestResource, "instances", "2"),
					resource.TestCheckResourceAttr(DBaaSClusterResource+"."+DBaaSClusterTestResource, "cores", "4"),
					resource.TestCheckResourceAttr(DBaaSClusterResource+"."+DBaaSClusterTestResource, "ram", "3072"),
					resource.TestCheckResourceAttr(DBaaSClusterResource+"."+DBaaSClusterTestResource, "storage_size", "3072"),
					resource.TestCheckResourceAttr(DBaaSClusterResource+"."+DBaaSClusterTestResource, "storage_type", "HDD"),
					resource.TestCheckResourceAttrPair(DBaaSClusterResource+"."+DBaaSClusterTestResource, "connections.0.datacenter_id", DatacenterResource+".datacenter_example_update", "id"),
					resource.TestCheckResourceAttrPair(DBaaSClusterResource+"."+DBaaSClusterTestResource, "connections.0.lan_id", LanResource+".lan_example_update", "id"),
					resource.TestCheckResourceAttr(DBaaSClusterResource+"."+DBaaSClusterTestResource, "connections.0.cidr", "192.168.1.101/24"),
					resource.TestCheckResourceAttrPair(DBaaSClusterResource+"."+DBaaSClusterTestResource, "location", DatacenterResource+".datacenter_example_update", "location"),
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
					resource.TestCheckResourceAttr(DBaaSClusterResource+"."+DBaaSClusterTestResource, "postgres_version", "13"),
					resource.TestCheckResourceAttr(DBaaSClusterResource+"."+DBaaSClusterTestResource, "instances", "2"),
					resource.TestCheckResourceAttr(DBaaSClusterResource+"."+DBaaSClusterTestResource, "cores", "4"),
					resource.TestCheckResourceAttr(DBaaSClusterResource+"."+DBaaSClusterTestResource, "ram", "3072"),
					resource.TestCheckResourceAttr(DBaaSClusterResource+"."+DBaaSClusterTestResource, "storage_size", "3072"),
					resource.TestCheckResourceAttr(DBaaSClusterResource+"."+DBaaSClusterTestResource, "storage_type", "HDD"),
					resource.TestCheckNoResourceAttr(DBaaSClusterResource+"."+DBaaSClusterTestResource, "connections"),
					resource.TestCheckResourceAttrPair(DBaaSClusterResource+"."+DBaaSClusterTestResource, "location", DatacenterResource+".datacenter_example_update", "location"),
					resource.TestCheckResourceAttr(DBaaSClusterResource+"."+DBaaSClusterTestResource, "display_name", UpdatedResources),
					resource.TestCheckResourceAttr(DBaaSClusterResource+"."+DBaaSClusterTestResource, "maintenance_window.0.time", "10:00:00"),
					resource.TestCheckResourceAttr(DBaaSClusterResource+"."+DBaaSClusterTestResource, "maintenance_window.0.day_of_the_week", "Saturday"),
					resource.TestCheckResourceAttr(DBaaSClusterResource+"."+DBaaSClusterTestResource, "credentials.0.username", "username"),
					resource.TestCheckResourceAttr(DBaaSClusterResource+"."+DBaaSClusterTestResource, "credentials.0.password", "password"),
					resource.TestCheckResourceAttr(DBaaSClusterResource+"."+DBaaSClusterTestResource, "synchronization_mode", "ASYNCHRONOUS"),
				),
			},
		},
	})
}

func TestAccDBaaSPgSqlClusterAdditionalParameters(t *testing.T) {
	var dbaasCluster dbaas.ClusterResponse
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
					resource.TestCheckResourceAttr(DBaaSClusterResource+"."+DBaaSClusterTestResource, "cores", "3"),
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
  cores              = 3
  ram                = 2048
  storage_size       = 2048
  storage_type       = "HDD"
  connections   {
	datacenter_id   =  ` + DatacenterResource + `.datacenter_example.id 
    lan_id          =  ` + LanResource + `.lan_example.id 
    cidr            =  "192.168.1.100/24"
  }
  location = ` + DatacenterResource + `.datacenter_example.location
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
  postgres_version   = 13
  instances          = 2
  cores              = 4
  ram                = 3072
  storage_size       = 3072
  storage_type       = "HDD"
  connections   {
	datacenter_id   =  ` + DatacenterResource + `.datacenter_example_update.id 
    lan_id          =  ` + LanResource + `.lan_example_update.id     
    cidr            =  "192.168.1.101/24"
  }
  location = ` + DatacenterResource + `.datacenter_example_update.location
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
  postgres_version   = 13
  instances          = 2
  cores              = 4
  ram                = 3072
  storage_size       = 3072
  storage_type       = "HDD"
  location = ` + DatacenterResource + `.datacenter_example_update.location
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
  cores              = 3
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
	backup_id = "2feaf22e-63f9-11ec-8761-8e8cbbc03eb8-4oymiqu-12"
    recovery_target_time = "2021-12-23T16:23:08Z"
  }
}`
