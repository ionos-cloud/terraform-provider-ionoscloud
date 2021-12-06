package ionoscloud

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	dbaas "github.com/ionos-cloud/sdk-go-autoscaling"
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
					resource.TestCheckResourceAttrPair(DBaaSClusterResource+"."+DBaaSClusterTestResource, "connections.0.datacenter_id", "ionoscloud_datacenter.test_dbaas_cluster", "id"),
					resource.TestCheckResourceAttrPair(DBaaSClusterResource+"."+DBaaSClusterTestResource, "connections.0.lan_id", "ionoscloud_lan.test_dbaas_cluster", "id"),
					resource.TestCheckResourceAttr(DBaaSClusterResource+"."+DBaaSClusterTestResource, "connections.0.cidr", "192.168.1.100/24"),
					resource.TestCheckResourceAttrPair(DBaaSClusterResource+"."+DBaaSClusterTestResource, "location", "ionoscloud_datacenter.test_dbaas_cluster", "location"),
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
					resource.TestCheckResourceAttr(DBaaSClusterResource+"."+DBaaSClusterTestResource, "ram", "2304"),
					resource.TestCheckResourceAttr(DBaaSClusterResource+"."+DBaaSClusterTestResource, "storage_size", "2304"),
					resource.TestCheckResourceAttr(DBaaSClusterResource+"."+DBaaSClusterTestResource, "storage_type", "HDD"),
					resource.TestCheckResourceAttrPair(DBaaSClusterResource+"."+DBaaSClusterTestResource, "connections.0.datacenter_id", "ionoscloud_datacenter.test_dbaas_cluster", "id"),
					resource.TestCheckResourceAttrPair(DBaaSClusterResource+"."+DBaaSClusterTestResource, "connections.0.lan_id", "ionoscloud_lan.test_dbaas_cluster", "id"),
					resource.TestCheckResourceAttr(DBaaSClusterResource+"."+DBaaSClusterTestResource, "connections.0.cidr", "192.168.1.100/24"),
					resource.TestCheckResourceAttrPair(DBaaSClusterResource+"."+DBaaSClusterTestResource, "location", "ionoscloud_datacenter.test_dbaas_cluster", "location"),
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
					resource.TestCheckResourceAttrPair(DBaaSClusterResource+"."+DBaaSClusterTestResource, "connections.0.datacenter_id", "ionoscloud_datacenter.test_dbaas_cluster", "id"),
					resource.TestCheckResourceAttrPair(DBaaSClusterResource+"."+DBaaSClusterTestResource, "connections.0.lan_id", "ionoscloud_lan.test_dbaas_cluster", "id"),
					resource.TestCheckResourceAttr(DBaaSClusterResource+"."+DBaaSClusterTestResource, "connections.0.cidr", "192.168.1.100/24"),
					resource.TestCheckResourceAttrPair(DBaaSClusterResource+"."+DBaaSClusterTestResource, "location", "ionoscloud_datacenter.test_dbaas_cluster", "location"),
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
resource "ionoscloud_datacenter" ` + DBaaSClusterTestResource + ` {
  name        = "` + DBaaSClusterTestResource + `"
  location    = "de/txl"
  description = "Datacenter for testing dbaas cluster"
}

resource "ionoscloud_lan" ` + DBaaSClusterTestResource + ` {
  datacenter_id = ionoscloud_datacenter.` + DBaaSClusterTestResource + `.id 
  public        = false
  name          = "` + DBaaSClusterTestResource + `"
}

resource ` + DBaaSClusterResource + ` ` + DBaaSClusterTestResource + ` {
  postgres_version   = 12
  instances          = 1
  cores              = 3
  ram                = 2048
  storage_size       = 2048
  storage_type       = "HDD"
  connections   {
	datacenter_id   =  ionoscloud_datacenter.` + DBaaSClusterTestResource + `.id 
    lan_id          =  ionoscloud_lan.` + DBaaSClusterTestResource + `.id 
    cidr            =  "192.168.1.100/24"
  }
  location = ionoscloud_datacenter.` + DBaaSClusterTestResource + `.location
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
resource "ionoscloud_datacenter" ` + DBaaSClusterTestResource + ` {
  name        = "` + DBaaSClusterTestResource + `"
  location    = "de/txl"
  description = "Datacenter for testing dbaas cluster"
}

resource "ionoscloud_lan" ` + DBaaSClusterTestResource + ` {
  datacenter_id = ionoscloud_datacenter.` + DBaaSClusterTestResource + `.id 
  public        = false
  name          = "` + DBaaSClusterTestResource + `"
}

resource ` + DBaaSClusterResource + ` ` + DBaaSClusterTestResource + ` {
  postgres_version   = 13
  instances          = 2
  cores              = 4
  ram                = 2304
  storage_size       = 2304
  storage_type       = "HDD"
  connections   {
	datacenter_id   =  ionoscloud_datacenter.` + DBaaSClusterTestResource + `.id 
    lan_id          =  ionoscloud_lan.` + DBaaSClusterTestResource + `.id 
    cidr            =  "192.168.1.100/24"
  }
  location = ionoscloud_datacenter.` + DBaaSClusterTestResource + `.location
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
resource "ionoscloud_datacenter" ` + DBaaSClusterTestResource + ` {
  name        = "` + DBaaSClusterTestResource + `"
  location    = "de/txl"
  description = "Datacenter for testing dbaas cluster"
}

resource "ionoscloud_lan" ` + DBaaSClusterTestResource + ` {
  datacenter_id = ionoscloud_datacenter.` + DBaaSClusterTestResource + `.id 
  public        = false
  name          = "` + DBaaSClusterTestResource + `"
}

resource ` + DBaaSClusterResource + ` ` + DBaaSClusterTestResource + ` {
  postgres_version   = 12
  instances          = 1
  cores              = 3
  ram                = 2048
  storage_size       = 2048
  storage_type       = "HDD"
  connections   {
	datacenter_id   =  ionoscloud_datacenter.` + DBaaSClusterTestResource + `.id 
    lan_id          =  ionoscloud_lan.` + DBaaSClusterTestResource + `.id 
    cidr            =  "192.168.1.100/24"
  }
  location = ionoscloud_datacenter.` + DBaaSClusterTestResource + `.location
  display_name = "` + DBaaSClusterTestResource + `"
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
	backup_id = "24e86a13-5688-11ec-903e-cebe353ec223-4oymiqu-12"
    recovery_target_time = "2021-12-06T13:54:08Z"
  }
}`
