package ionoscloud

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	dbaas "github.com/ionos-cloud/sdk-go-autoscaling"
	"testing"
)

func TestAccDbaasPgSqlCluster_Basic(t *testing.T) {
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
					testAccCheckDbaasPgSqlClusterExists("ionoscloud_dbaas_pgsql_cluster.test_dbaas_cluster", &dbaasCluster),
					resource.TestCheckResourceAttr("ionoscloud_dbaas_pgsql_cluster.test_dbaas_cluster", "postgres_version", "12"),
					resource.TestCheckResourceAttr("ionoscloud_dbaas_pgsql_cluster.test_dbaas_cluster", "replicas", "1"),
					resource.TestCheckResourceAttr("ionoscloud_dbaas_pgsql_cluster.test_dbaas_cluster", "cpu_core_count", "4"),
					resource.TestCheckResourceAttr("ionoscloud_dbaas_pgsql_cluster.test_dbaas_cluster", "ram_size", "3Gi"),
					resource.TestCheckResourceAttr("ionoscloud_dbaas_pgsql_cluster.test_dbaas_cluster", "storage_size", "1Gi"),
					resource.TestCheckResourceAttr("ionoscloud_dbaas_pgsql_cluster.test_dbaas_cluster", "storage_type", "HDD"),
					resource.TestCheckResourceAttrPair("ionoscloud_dbaas_pgsql_cluster.test_dbaas_cluster", "vdc_connections.0.vdc_id", "ionoscloud_datacenter.test_dbaas_cluster", "id"),
					resource.TestCheckResourceAttrPair("ionoscloud_dbaas_pgsql_cluster.test_dbaas_cluster", "vdc_connections.0.lan_id", "ionoscloud_lan.test_dbaas_cluster", "id"),
					resource.TestCheckResourceAttr("ionoscloud_dbaas_pgsql_cluster.test_dbaas_cluster", "vdc_connections.0.ip_address", "192.168.1.100/24"),
					resource.TestCheckResourceAttrPair("ionoscloud_dbaas_pgsql_cluster.test_dbaas_cluster", "location", "ionoscloud_datacenter.test_dbaas_cluster", "location"),
					resource.TestCheckResourceAttr("ionoscloud_dbaas_pgsql_cluster.test_dbaas_cluster", "display_name", "PostgreSQL_cluster"),
					resource.TestCheckResourceAttr("ionoscloud_dbaas_pgsql_cluster.test_dbaas_cluster", "backup_enabled", "true"),
					resource.TestCheckResourceAttr("ionoscloud_dbaas_pgsql_cluster.test_dbaas_cluster", "maintenance_window.0.time", "09:00:00"),
					resource.TestCheckResourceAttr("ionoscloud_dbaas_pgsql_cluster.test_dbaas_cluster", "maintenance_window.0.weekday", "Sunday"),
					resource.TestCheckResourceAttr("ionoscloud_dbaas_pgsql_cluster.test_dbaas_cluster", "credentials.0.username", "username"),
					resource.TestCheckResourceAttr("ionoscloud_dbaas_pgsql_cluster.test_dbaas_cluster", "credentials.0.password", "password"),
				),
			},
			{
				Config: testAccCheckDbaasPgSqlClusterConfigUpdate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDbaasPgSqlClusterExists("ionoscloud_dbaas_pgsql_cluster.test_dbaas_cluster", &dbaasCluster),
					resource.TestCheckResourceAttr("ionoscloud_dbaas_pgsql_cluster.test_dbaas_cluster", "postgres_version", "13"),
					resource.TestCheckResourceAttr("ionoscloud_dbaas_pgsql_cluster.test_dbaas_cluster", "replicas", "1"),
					resource.TestCheckResourceAttr("ionoscloud_dbaas_pgsql_cluster.test_dbaas_cluster", "cpu_core_count", "4"),
					resource.TestCheckResourceAttr("ionoscloud_dbaas_pgsql_cluster.test_dbaas_cluster", "ram_size", "3Gi"),
					resource.TestCheckResourceAttr("ionoscloud_dbaas_pgsql_cluster.test_dbaas_cluster", "storage_size", "2Gi"),
					resource.TestCheckResourceAttr("ionoscloud_dbaas_pgsql_cluster.test_dbaas_cluster", "storage_type", "HDD"),
					resource.TestCheckResourceAttrPair("ionoscloud_dbaas_pgsql_cluster.test_dbaas_cluster", "vdc_connections.0.vdc_id", "ionoscloud_datacenter.test_dbaas_cluster", "id"),
					resource.TestCheckResourceAttrPair("ionoscloud_dbaas_pgsql_cluster.test_dbaas_cluster", "vdc_connections.0.lan_id", "ionoscloud_lan.test_dbaas_cluster", "id"),
					resource.TestCheckResourceAttr("ionoscloud_dbaas_pgsql_cluster.test_dbaas_cluster", "vdc_connections.0.ip_address", "192.168.1.100/24"),
					resource.TestCheckResourceAttrPair("ionoscloud_dbaas_pgsql_cluster.test_dbaas_cluster", "location", "ionoscloud_datacenter.test_dbaas_cluster", "location"),
					resource.TestCheckResourceAttr("ionoscloud_dbaas_pgsql_cluster.test_dbaas_cluster", "display_name", "PostgreSQL_cluster_update"),
					resource.TestCheckResourceAttr("ionoscloud_dbaas_pgsql_cluster.test_dbaas_cluster", "backup_enabled", "true"),
					resource.TestCheckResourceAttr("ionoscloud_dbaas_pgsql_cluster.test_dbaas_cluster", "maintenance_window.0.time", "10:00:00"),
					resource.TestCheckResourceAttr("ionoscloud_dbaas_pgsql_cluster.test_dbaas_cluster", "maintenance_window.0.weekday", "Saturday"),
					resource.TestCheckResourceAttr("ionoscloud_dbaas_pgsql_cluster.test_dbaas_cluster", "credentials.0.username", "username"),
					resource.TestCheckResourceAttr("ionoscloud_dbaas_pgsql_cluster.test_dbaas_cluster", "credentials.0.password", "password"),
				),
			},
		},
	})
}

func TestAccDbaasPgSqlCluster_FromBackup(t *testing.T) {
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
					testAccCheckDbaasPgSqlClusterExists("ionoscloud_dbaas_pgsql_cluster.from_backup", &dbaasCluster),
					resource.TestCheckResourceAttr("ionoscloud_dbaas_pgsql_cluster.from_backup", "postgres_version", "12"),
					resource.TestCheckResourceAttr("ionoscloud_dbaas_pgsql_cluster.from_backup", "replicas", "1"),
					resource.TestCheckResourceAttr("ionoscloud_dbaas_pgsql_cluster.from_backup", "cpu_core_count", "4"),
					resource.TestCheckResourceAttr("ionoscloud_dbaas_pgsql_cluster.from_backup", "ram_size", "3Gi"),
					resource.TestCheckResourceAttr("ionoscloud_dbaas_pgsql_cluster.from_backup", "storage_size", "1Gi"),
					resource.TestCheckResourceAttr("ionoscloud_dbaas_pgsql_cluster.from_backup", "storage_type", "HDD"),
					resource.TestCheckResourceAttrPair("ionoscloud_dbaas_pgsql_cluster.from_backup", "vdc_connections.0.vdc_id", "ionoscloud_datacenter.test_dbaas_cluster", "id"),
					resource.TestCheckResourceAttrPair("ionoscloud_dbaas_pgsql_cluster.from_backup", "vdc_connections.0.lan_id", "ionoscloud_lan.test_dbaas_cluster", "id"),
					resource.TestCheckResourceAttr("ionoscloud_dbaas_pgsql_cluster.from_backup", "vdc_connections.0.ip_address", "192.168.2.100/24"),
					resource.TestCheckResourceAttrPair("ionoscloud_dbaas_pgsql_cluster.from_backup", "location", "ionoscloud_datacenter.test_dbaas_cluster", "location"),
					resource.TestCheckResourceAttr("ionoscloud_dbaas_pgsql_cluster.from_backup", "display_name", "PostgreSQL_cluster_from_Backup"),
					resource.TestCheckResourceAttr("ionoscloud_dbaas_pgsql_cluster.from_backup", "credentials.0.username", "username"),
					resource.TestCheckResourceAttr("ionoscloud_dbaas_pgsql_cluster.from_backup", "credentials.0.password", "password"),
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
  ram_size           = "3Gi"
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

const testAccCheckDbaasPgSqlClusterConfigUpdate = `
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
  postgres_version   = 13
  replicas           = 1
  cpu_core_count     = 4
  ram_size           = "3Gi"
  storage_size       = "2Gi"
  storage_type       = "HDD"
  vdc_connections   {
	vdc_id          =  ionoscloud_datacenter.test_dbaas_cluster.id 
    lan_id          =  ionoscloud_lan.test_dbaas_cluster.id 
    ip_address      =  "192.168.1.100/24"
  }
  location = ionoscloud_datacenter.test_dbaas_cluster.location
  display_name = "PostgreSQL_cluster_update"
  backup_enabled = true
  maintenance_window {
    weekday = "Saturday"
    time            = "10:00:00"
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

const testAccFromBackup = `
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

resource "ionoscloud_dbaas_pgsql_cluster" "from_backup" {
  postgres_version   = 12
  replicas           = 1
  cpu_core_count     = 4
  ram_size           = "3Gi"
  storage_size       = "1Gi"
  storage_type       = "HDD"
  vdc_connections   {
	vdc_id          =  ionoscloud_datacenter.test_dbaas_cluster.id 
    lan_id          =  ionoscloud_lan.test_dbaas_cluster.id 
    ip_address      =  "192.168.2.100/24"
  }
  location = ionoscloud_datacenter.test_dbaas_cluster.location
  display_name = "PostgreSQL_cluster_from_Backup"
  credentials {
  	username = "username"
	password = "password"

  }
  from_backup = "3273774b-2116-11ec-bd55-d6a61c20e878-4oymiqu-12"
}`
