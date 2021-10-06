package ionoscloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceDBaaSPgSqlCluster(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceDBaaSPgSqlClusterCreateResources,
			},
			{
				Config: testAccDataSourceDBaaSPgSqlClusterMatchId,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair("data.ionoscloud_dbaas_pgsql_cluster.test_ds_dbaas_cluster_id", "display_name", "ionoscloud_dbaas_pgsql_cluster.test_dbaas_cluster", "display_name"),
					resource.TestCheckResourceAttrPair("data.ionoscloud_dbaas_pgsql_cluster.test_ds_dbaas_cluster_id", "replicas", "ionoscloud_dbaas_pgsql_cluster.test_dbaas_cluster", "replicas"),
					resource.TestCheckResourceAttrPair("data.ionoscloud_dbaas_pgsql_cluster.test_ds_dbaas_cluster_id", "cpu_core_count", "ionoscloud_dbaas_pgsql_cluster.test_dbaas_cluster", "cpu_core_count"),
					resource.TestCheckResourceAttrPair("data.ionoscloud_dbaas_pgsql_cluster.test_ds_dbaas_cluster_id", "ram_size", "ionoscloud_dbaas_pgsql_cluster.test_dbaas_cluster", "ram_size"),
					resource.TestCheckResourceAttrPair("data.ionoscloud_dbaas_pgsql_cluster.test_ds_dbaas_cluster_id", "storage_size", "ionoscloud_dbaas_pgsql_cluster.test_dbaas_cluster", "storage_size"),
					resource.TestCheckResourceAttrPair("data.ionoscloud_dbaas_pgsql_cluster.test_ds_dbaas_cluster_id", "storage_type", "ionoscloud_dbaas_pgsql_cluster.test_dbaas_cluster", "storage_type"),
					resource.TestCheckResourceAttrPair("data.ionoscloud_dbaas_pgsql_cluster.test_ds_dbaas_cluster_id", "vdc_connections.vdc_id", "ionoscloud_dbaas_pgsql_cluster.test_dbaas_cluster", "vdc_connections.vdc_id"),
					resource.TestCheckResourceAttrPair("data.ionoscloud_dbaas_pgsql_cluster.test_ds_dbaas_cluster_id", "vdc_connections.lan_id", "ionoscloud_dbaas_pgsql_cluster.test_dbaas_cluster", "vdc_connections.lan_id"),
					resource.TestCheckResourceAttrPair("data.ionoscloud_dbaas_pgsql_cluster.test_ds_dbaas_cluster_id", "vdc_connections.ip_address", "ionoscloud_dbaas_pgsql_cluster.test_dbaas_cluster", "vdc_connections.ip_address"),
					resource.TestCheckResourceAttrPair("data.ionoscloud_dbaas_pgsql_cluster.test_ds_dbaas_cluster_id", "location", "ionoscloud_dbaas_pgsql_cluster.test_dbaas_cluster", "location"),
					resource.TestCheckResourceAttrPair("data.ionoscloud_dbaas_pgsql_cluster.test_ds_dbaas_cluster_id", "display_name", "ionoscloud_dbaas_pgsql_cluster.test_dbaas_cluster", "display_name"),
					resource.TestCheckResourceAttrPair("data.ionoscloud_dbaas_pgsql_cluster.test_ds_dbaas_cluster_id", "maintenance_window.weekday", "ionoscloud_dbaas_pgsql_cluster.test_dbaas_cluster", "maintenance_window.weekday"),
					resource.TestCheckResourceAttrPair("data.ionoscloud_dbaas_pgsql_cluster.test_ds_dbaas_cluster_id", "maintenance_window.time", "ionoscloud_dbaas_pgsql_cluster.test_dbaas_cluster", "maintenance_window.time"),
					resource.TestCheckResourceAttrPair("data.ionoscloud_dbaas_pgsql_cluster.test_ds_dbaas_cluster_id", "credentials.username", "ionoscloud_dbaas_pgsql_cluster.test_dbaas_cluster", "credentials.username"),
					resource.TestCheckResourceAttrPair("data.ionoscloud_dbaas_pgsql_cluster.test_ds_dbaas_cluster_id", "credentials.password", "ionoscloud_dbaas_pgsql_cluster.test_dbaas_cluster", "credentials.password"),
				),
			},
			{
				Config: testAccDataSourceDBaaSPgSqlClusterMatchName,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair("data.ionoscloud_dbaas_pgsql_cluster.test_ds_dbaas_cluster_name", "display_name", "ionoscloud_dbaas_pgsql_cluster.test_dbaas_cluster", "display_name"),
					resource.TestCheckResourceAttrPair("data.ionoscloud_dbaas_pgsql_cluster.test_ds_dbaas_cluster_name", "replicas", "ionoscloud_dbaas_pgsql_cluster.test_dbaas_cluster", "replicas"),
					resource.TestCheckResourceAttrPair("data.ionoscloud_dbaas_pgsql_cluster.test_ds_dbaas_cluster_name", "cpu_core_count", "ionoscloud_dbaas_pgsql_cluster.test_dbaas_cluster", "cpu_core_count"),
					resource.TestCheckResourceAttrPair("data.ionoscloud_dbaas_pgsql_cluster.test_ds_dbaas_cluster_name", "ram_size", "ionoscloud_dbaas_pgsql_cluster.test_dbaas_cluster", "ram_size"),
					resource.TestCheckResourceAttrPair("data.ionoscloud_dbaas_pgsql_cluster.test_ds_dbaas_cluster_name", "storage_size", "ionoscloud_dbaas_pgsql_cluster.test_dbaas_cluster", "storage_size"),
					resource.TestCheckResourceAttrPair("data.ionoscloud_dbaas_pgsql_cluster.test_ds_dbaas_cluster_name", "storage_type", "ionoscloud_dbaas_pgsql_cluster.test_dbaas_cluster", "storage_type"),
					resource.TestCheckResourceAttrPair("data.ionoscloud_dbaas_pgsql_cluster.test_ds_dbaas_cluster_name", "vdc_connections.vdc_id", "ionoscloud_dbaas_pgsql_cluster.test_dbaas_cluster", "vdc_connections.vdc_id"),
					resource.TestCheckResourceAttrPair("data.ionoscloud_dbaas_pgsql_cluster.test_ds_dbaas_cluster_name", "vdc_connections.lan_id", "ionoscloud_dbaas_pgsql_cluster.test_dbaas_cluster", "vdc_connections.lan_id"),
					resource.TestCheckResourceAttrPair("data.ionoscloud_dbaas_pgsql_cluster.test_ds_dbaas_cluster_name", "vdc_connections.ip_address", "ionoscloud_dbaas_pgsql_cluster.test_dbaas_cluster", "vdc_connections.ip_address"),
					resource.TestCheckResourceAttrPair("data.ionoscloud_dbaas_pgsql_cluster.test_ds_dbaas_cluster_name", "location", "ionoscloud_dbaas_pgsql_cluster.test_dbaas_cluster", "location"),
					resource.TestCheckResourceAttrPair("data.ionoscloud_dbaas_pgsql_cluster.test_ds_dbaas_cluster_name", "display_name", "ionoscloud_dbaas_pgsql_cluster.test_dbaas_cluster", "display_name"),
					resource.TestCheckResourceAttrPair("data.ionoscloud_dbaas_pgsql_cluster.test_ds_dbaas_cluster_name", "maintenance_window.weekday", "ionoscloud_dbaas_pgsql_cluster.test_dbaas_cluster", "maintenance_window.weekday"),
					resource.TestCheckResourceAttrPair("data.ionoscloud_dbaas_pgsql_cluster.test_ds_dbaas_cluster_name", "maintenance_window.time", "ionoscloud_dbaas_pgsql_cluster.test_dbaas_cluster", "maintenance_window.time"),
					resource.TestCheckResourceAttrPair("data.ionoscloud_dbaas_pgsql_cluster.test_ds_dbaas_cluster_name", "credentials.username", "ionoscloud_dbaas_pgsql_cluster.test_dbaas_cluster", "credentials.username"),
					resource.TestCheckResourceAttrPair("data.ionoscloud_dbaas_pgsql_cluster.test_ds_dbaas_cluster_name", "credentials.password", "ionoscloud_dbaas_pgsql_cluster.test_dbaas_cluster", "credentials.password"),
				),
			},
		},
	})
}

const testAccDataSourceDBaaSPgSqlClusterCreateResources = `
resource "ionoscloud_datacenter" "test_dbaas_cluster" {
  name        = "test_dbaas_cluster"
  location    = "de/txl"
  description = "Datacenter for testing dbaas cluster"
}

resource "ionoscloud_lan" "test_dbaas_cluster" {
  datacenter_id = ionoscloud_datacenter.test_dbaas_cluster.id 
  public        = false
  name          = "test_dbaas_cluster"
}

resource "ionoscloud_dbaas_pgsql_cluster" "test_dbaas_cluster" {
  postgres_version   = 12
  replicas           = 1
  cpu_core_count     = 4
  ram_size           = "2Gi"
  storage_size       = "1Gi"
  storage_type       = "HDD"
  vdc_connections   {
	vdc_id          =  ionoscloud_datacenter.test_dbaas_cluster.id 
    lan_id          =  ionoscloud_lan.test_dbaas_cluster.id 
    ip_address      =  "192.168.1.100/24"
  }
  location = ionoscloud_datacenter.test_dbaas_cluster.location
  display_name = "PostgreSQL_cluster"
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

const testAccDataSourceDBaaSPgSqlClusterMatchId = `
resource "ionoscloud_datacenter" "test_dbaas_cluster" {
  name        = "test_dbaas_cluster"
  location    = "de/txl"
  description = "Datacenter for testing dbaas cluster"
}

resource "ionoscloud_lan" "test_dbaas_cluster" {
  datacenter_id = ionoscloud_datacenter.test_dbaas_cluster.id 
  public        = false
  name          = "test_dbaas_cluster"
}

resource "ionoscloud_dbaas_pgsql_cluster" "test_dbaas_cluster" {
  postgres_version   = 12
  replicas           = 1
  cpu_core_count     = 4
  ram_size           = "2Gi"
  storage_size       = "1Gi"
  storage_type       = "HDD"
  vdc_connections   {
	vdc_id          =  ionoscloud_datacenter.test_dbaas_cluster.id 
    lan_id          =  ionoscloud_lan.test_dbaas_cluster.id 
    ip_address      =  "192.168.1.100/24"
  }
  location = ionoscloud_datacenter.test_dbaas_cluster.location
  display_name = "PostgreSQL_cluster"
  maintenance_window {
    weekday = "Sunday"
    time            = "09:00:00"
  }
  credentials {
  	username = "username"
	password = "password"
  }
}

data "ionoscloud_dbaas_pgsql_cluster" "test_ds_dbaas_cluster_id" {
  id	= ionoscloud_dbaas_pgsql_cluster.test_dbaas_cluster.id
}
`

const testAccDataSourceDBaaSPgSqlClusterMatchName = `
resource "ionoscloud_datacenter" "test_dbaas_cluster" {
  name        = "test_dbaas_cluster"
  location    = "de/txl"
  description = "Datacenter for testing dbaas cluster"
}

resource "ionoscloud_lan" "test_dbaas_cluster" {
  datacenter_id = ionoscloud_datacenter.test_dbaas_cluster.id 
  public        = false
  name          = "test_dbaas_cluster"
}

resource "ionoscloud_dbaas_pgsql_cluster" "test_dbaas_cluster" {
  postgres_version   = 12
  replicas           = 1
  cpu_core_count     = 4
  ram_size           = "2Gi"
  storage_size       = "1Gi"
  storage_type       = "HDD"
  vdc_connections   {
	vdc_id          =  ionoscloud_datacenter.test_dbaas_cluster.id 
    lan_id          =  ionoscloud_lan.test_dbaas_cluster.id 
    ip_address      =  "192.168.1.100/24"
  }
  location = ionoscloud_datacenter.test_dbaas_cluster.location
  display_name = "PostgreSQL_cluster"
  maintenance_window {
    weekday = "Sunday"
    time            = "09:00:00"
  }
  credentials {
  	username = "username"
	password = "password"
  }
}

data "ionoscloud_dbaas_pgsql_cluster" "test_ds_dbaas_cluster_name" {
  display_name	= "PostgreSQL_cluster"
}
`
