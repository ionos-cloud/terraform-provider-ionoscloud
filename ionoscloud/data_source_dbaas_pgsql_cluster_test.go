package ionoscloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceDbaasPgSqlCluster_matchId(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceDbaasPgSqlClusterCreateResources,
			},
			{
				Config: testAccDataSourceDbaasPgSqlClusterMatchId,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.ionoscloud_dbaas_pgsql_cluster.test_ds_dbaas_cluster", "display_name", "PostgreSQL_cluster"),
				),
			},
		},
	})
}

func TestAccDataSourceDbaasPgSqlCluster_matchName(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceDbaasPgSqlClusterCreateResources,
			},
			{
				Config: testAccDataSourceDbaasPgSqlClusterMatchName,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.ionoscloud_dbaas_pgsql_cluster.test_ds_dbaas_cluster", "display_name", "PostgreSQL_cluster"),
				),
			},
		},
	})

}

const testAccDataSourceDbaasPgSqlClusterCreateResources = `
resource "ionoscloud_datacenter" "test_dbaas_cluster" {
  name        = "test_dbaas_cluster"
  location    = "de/txl"
  description = "Datacenter for testing dbaas cluster"
}

resource "ionoscloud_ipblock" "test_dbaas_cluster" {
  location = ionoscloud_datacenter.test_dbaas_cluster.location
  size = 1
  name = "test_dbaas_cluster"
}

resource "ionoscloud_lan" "test_dbaas_cluster" {
  datacenter_id = ionoscloud_datacenter.test_dbaas_cluster.id 
  public        = false
  name          = "test_dbaas_cluster"
}

resource "ionoscloud_dbaas_pgsql_cluster" "test_dbaas_cluster" {
  postgres_version   = 12
  replicas           = 2
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
  backup_enabled = true
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

const testAccDataSourceDbaasPgSqlClusterMatchId = `
resource "ionoscloud_datacenter" "test_dbaas_cluster" {
  name        = "test_dbaas_cluster"
  location    = "de/txl"
  description = "Datacenter for testing dbaas cluster"
}

resource "ionoscloud_ipblock" "test_dbaas_cluster" {
  location = ionoscloud_datacenter.test_dbaas_cluster.location
  size = 1
  name = "test_dbaas_cluster"
}

resource "ionoscloud_lan" "test_dbaas_cluster" {
  datacenter_id = ionoscloud_datacenter.test_dbaas_cluster.id 
  public        = false
  name          = "test_dbaas_cluster"
}

resource "ionoscloud_dbaas_pgsql_cluster" "test_dbaas_cluster" {
  postgres_version   = 12
  replicas           = 2
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
  backup_enabled = true
  maintenance_window {
    weekday = "Sunday"
    time            = "09:00:00"
  }
  credentials {
  	username = "username"
	password = "password"
  }
}

data "ionoscloud_dbaas_pgsql_cluster" "test_ds_dbaas_cluster" {
  id	= ionoscloud_dbaas_pgsql_cluster.test_dbaas_cluster.id
}
`

const testAccDataSourceDbaasPgSqlClusterMatchName = `
resource "ionoscloud_datacenter" "test_dbaas_cluster" {
  name        = "test_dbaas_cluster"
  location    = "de/txl"
  description = "Datacenter for testing dbaas cluster"
}

resource "ionoscloud_ipblock" "test_dbaas_cluster" {
  location = ionoscloud_datacenter.test_dbaas_cluster.location
  size = 1
  name = "test_dbaas_cluster"
}

resource "ionoscloud_lan" "test_dbaas_cluster" {
  datacenter_id = ionoscloud_datacenter.test_dbaas_cluster.id 
  public        = false
  name          = "test_dbaas_cluster"
}

resource "ionoscloud_dbaas_pgsql_cluster" "test_dbaas_cluster" {
  postgres_version   = 12
  replicas           = 2
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
  backup_enabled = true
  maintenance_window {
    weekday = "Sunday"
    time            = "09:00:00"
  }
  credentials {
  	username = "username"
	password = "password"
  }
}

data "ionoscloud_dbaas_pgsql_cluster" "test_ds_dbaas_cluster" {
  display_name	= "PostgreSQL_cluster"
}
`
