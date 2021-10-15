package ionoscloud

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	dbaas "github.com/ionos-cloud/sdk-go-autoscaling"
	"testing"
)

func TestAccDBaaSPgSqlCluster_Basic(t *testing.T) {
	var dbaasCluster dbaas.Cluster

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
					resource.TestCheckResourceAttr(DBaaSClusterResource+"."+DBaaSClusterTestResource, "replicas", "1"),
					resource.TestCheckResourceAttr(DBaaSClusterResource+"."+DBaaSClusterTestResource, "cpu_core_count", "4"),
					resource.TestCheckResourceAttr(DBaaSClusterResource+"."+DBaaSClusterTestResource, "ram_size", "3Gi"),
					resource.TestCheckResourceAttr(DBaaSClusterResource+"."+DBaaSClusterTestResource, "storage_size", "2Gi"),
					resource.TestCheckResourceAttr(DBaaSClusterResource+"."+DBaaSClusterTestResource, "storage_type", "HDD"),
					resource.TestCheckResourceAttrPair(DBaaSClusterResource+"."+DBaaSClusterTestResource, "vdc_connections.0.vdc_id", "ionoscloud_datacenter.test_dbaas_cluster", "id"),
					resource.TestCheckResourceAttrPair(DBaaSClusterResource+"."+DBaaSClusterTestResource, "vdc_connections.0.lan_id", "ionoscloud_lan.test_dbaas_cluster", "id"),
					resource.TestCheckResourceAttr(DBaaSClusterResource+"."+DBaaSClusterTestResource, "vdc_connections.0.ip_address", "192.168.1.100/24"),
					resource.TestCheckResourceAttrPair(DBaaSClusterResource+"."+DBaaSClusterTestResource, "location", "ionoscloud_datacenter.test_dbaas_cluster", "location"),
					resource.TestCheckResourceAttr(DBaaSClusterResource+"."+DBaaSClusterTestResource, "display_name", DBaaSClusterTestResource),
					resource.TestCheckResourceAttr(DBaaSClusterResource+"."+DBaaSClusterTestResource, "maintenance_window.0.time", "09:00:00"),
					resource.TestCheckResourceAttr(DBaaSClusterResource+"."+DBaaSClusterTestResource, "maintenance_window.0.weekday", "Sunday"),
					resource.TestCheckResourceAttr(DBaaSClusterResource+"."+DBaaSClusterTestResource, "credentials.0.username", "username"),
					resource.TestCheckResourceAttr(DBaaSClusterResource+"."+DBaaSClusterTestResource, "credentials.0.password", "password"),
				),
			},
			{
				Config: testAccCheckDbaasPgSqlClusterConfigUpdate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDbaasPgSqlClusterExists(DBaaSClusterResource+"."+DBaaSClusterTestResource, &dbaasCluster),
					resource.TestCheckResourceAttr(DBaaSClusterResource+"."+DBaaSClusterTestResource, "postgres_version", "13"),
					resource.TestCheckResourceAttr(DBaaSClusterResource+"."+DBaaSClusterTestResource, "replicas", "1"),
					resource.TestCheckResourceAttr(DBaaSClusterResource+"."+DBaaSClusterTestResource, "cpu_core_count", "4"),
					resource.TestCheckResourceAttr(DBaaSClusterResource+"."+DBaaSClusterTestResource, "ram_size", "3Gi"),
					resource.TestCheckResourceAttr(DBaaSClusterResource+"."+DBaaSClusterTestResource, "storage_size", "3Gi"),
					resource.TestCheckResourceAttr(DBaaSClusterResource+"."+DBaaSClusterTestResource, "storage_type", "HDD"),
					resource.TestCheckResourceAttrPair(DBaaSClusterResource+"."+DBaaSClusterTestResource, "vdc_connections.0.vdc_id", "ionoscloud_datacenter.test_dbaas_cluster", "id"),
					resource.TestCheckResourceAttrPair(DBaaSClusterResource+"."+DBaaSClusterTestResource, "vdc_connections.0.lan_id", "ionoscloud_lan.test_dbaas_cluster", "id"),
					resource.TestCheckResourceAttr(DBaaSClusterResource+"."+DBaaSClusterTestResource, "vdc_connections.0.ip_address", "192.168.1.100/24"),
					resource.TestCheckResourceAttrPair(DBaaSClusterResource+"."+DBaaSClusterTestResource, "location", "ionoscloud_datacenter.test_dbaas_cluster", "location"),
					resource.TestCheckResourceAttr(DBaaSClusterResource+"."+DBaaSClusterTestResource, "display_name", UpdatedResources),
					resource.TestCheckResourceAttr(DBaaSClusterResource+"."+DBaaSClusterTestResource, "maintenance_window.0.time", "10:00:00"),
					resource.TestCheckResourceAttr(DBaaSClusterResource+"."+DBaaSClusterTestResource, "maintenance_window.0.weekday", "Saturday"),
					resource.TestCheckResourceAttr(DBaaSClusterResource+"."+DBaaSClusterTestResource, "credentials.0.username", "username"),
					resource.TestCheckResourceAttr(DBaaSClusterResource+"."+DBaaSClusterTestResource, "credentials.0.password", "password"),
				),
			},
		},
	})
}

func TestAccDBaaSPgSqlClusterAdditionalParameters(t *testing.T) {
	// if you want to remove this line in order to test, please be sure you replace from_backup and from_recovery_target_time
	// arguments with valid values since now they are hardcoded
	//t.Skip()

	var dbaasCluster dbaas.Cluster
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
					resource.TestCheckResourceAttr(DBaaSClusterResource+"."+DBaaSClusterTestResource, "replicas", "1"),
					resource.TestCheckResourceAttr(DBaaSClusterResource+"."+DBaaSClusterTestResource, "cpu_core_count", "4"),
					resource.TestCheckResourceAttr(DBaaSClusterResource+"."+DBaaSClusterTestResource, "ram_size", "2Gi"),
					resource.TestCheckResourceAttr(DBaaSClusterResource+"."+DBaaSClusterTestResource, "storage_size", "2Gi"),
					resource.TestCheckResourceAttr(DBaaSClusterResource+"."+DBaaSClusterTestResource, "storage_type", "HDD"),
					resource.TestCheckResourceAttrPair(DBaaSClusterResource+"."+DBaaSClusterTestResource, "vdc_connections.0.vdc_id", "ionoscloud_datacenter.test_dbaas_cluster", "id"),
					resource.TestCheckResourceAttrPair(DBaaSClusterResource+"."+DBaaSClusterTestResource, "vdc_connections.0.lan_id", "ionoscloud_lan.test_dbaas_cluster", "id"),
					resource.TestCheckResourceAttr(DBaaSClusterResource+"."+DBaaSClusterTestResource, "vdc_connections.0.ip_address", "192.168.1.100/24"),
					resource.TestCheckResourceAttrPair(DBaaSClusterResource+"."+DBaaSClusterTestResource, "location", "ionoscloud_datacenter.test_dbaas_cluster", "location"),
					resource.TestCheckResourceAttr(DBaaSClusterResource+"."+DBaaSClusterTestResource, "display_name", "PostgreSQL_cluster_from_Backup"),
					resource.TestCheckResourceAttr(DBaaSClusterResource+"."+DBaaSClusterTestResource, "credentials.0.username", "username"),
					resource.TestCheckResourceAttr(DBaaSClusterResource+"."+DBaaSClusterTestResource, "credentials.0.password", "password"),
				),
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
		if rs.Type != "ionoscloud_dbaas_pgsql_cluster" {
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

func testAccCheckDbaasPgSqlClusterExists(n string, cluster *dbaas.Cluster) resource.TestCheckFunc {
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
  replicas           = 1
  cpu_core_count     = 4
  ram_size           = "3Gi"
  storage_size       = "2Gi"
  storage_type       = "HDD"
  vdc_connections   {
	vdc_id          =  ionoscloud_datacenter.` + DBaaSClusterTestResource + `.id 
    lan_id          =  ionoscloud_lan.` + DBaaSClusterTestResource + `.id 
    ip_address      =  "192.168.1.100/24"
  }
  location = ionoscloud_datacenter.` + DBaaSClusterTestResource + `.location
  display_name = "` + DBaaSClusterTestResource + `"
  maintenance_window {
    weekday = "Sunday"
    time            = "09:00:00"
  }
  credentials {
  	username = "username"
	password = "password"
  }
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
  replicas           = 1
  cpu_core_count     = 4
  ram_size           = "3Gi"
  storage_size       = "3Gi"
  storage_type       = "HDD"
  vdc_connections   {
	vdc_id          =  ionoscloud_datacenter.` + DBaaSClusterTestResource + `.id 
    lan_id          =  ionoscloud_lan.` + DBaaSClusterTestResource + `.id 
    ip_address      =  "192.168.1.100/24"
  }
  location = ionoscloud_datacenter.` + DBaaSClusterTestResource + `.location
  display_name = "` + UpdatedResources + `"
  maintenance_window {
    weekday = "Saturday"
    time            = "10:00:00"
  }
  credentials {
  	username = "username"
	password = "password"
  }
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
  replicas           = 1
  cpu_core_count     = 4
  ram_size           = "2Gi"
  storage_size       = "2Gi"
  storage_type       = "HDD"
  vdc_connections   {
	vdc_id          =  ionoscloud_datacenter.` + DBaaSClusterTestResource + `.id 
    lan_id          =  ionoscloud_lan.` + DBaaSClusterTestResource + `.id 
    ip_address      =  "192.168.1.100/24"
  }
  location = ionoscloud_datacenter.` + DBaaSClusterTestResource + `.location
  display_name = "` + DBaaSClusterTestResource + `"
  maintenance_window {
    weekday = "Sunday"
    time            = "09:00:00"
  }
  credentials {
  	username = "username"
	password = "password"
  }
  from_backup = "ad7ac139-2d0b-11ec-a2e3-92fbe7e27ed1-4oymiqu-12"
  from_recovery_target_time = "2021-10-14T19:36:19Z"
}`
