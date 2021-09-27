package ionoscloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceDbaasPgSqlClusterBackups(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceDbaasPgSqlClusterBackups,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair("data.ionoscloud_dbaas_pgsql_backups.test_ds_dbaas_backups", "cluster_backups.0.cluster_id", "data.ionoscloud_dbaas_pgsql_backups.test_ds_dbaas_backups", "cluster_id"),
					testNotEmptySlice("ionoscloud_dbaas_pgsql_backups", "cluster_backups.#"),
				),
			},
		},
	})
}

const testAccDataSourceDbaasPgSqlClusterBackups = `

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

data "ionoscloud_dbaas_pgsql_backups" "test_ds_dbaas_backups" {
	cluster_id = ionoscloud_dbaas_pgsql_cluster.test_dbaas_cluster.id
}
`
