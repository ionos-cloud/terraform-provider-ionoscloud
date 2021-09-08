package ionoscloud

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	dbaas "github.com/ionos-cloud/sdk-go-autoscaling"
	"testing"
)

func TestAccDbaasCluster_Basic(t *testing.T) {
	var dbaasCluster dbaas.Cluster

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckDbaasClusterDestroyCheck,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckDbaasClusterConfigBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDbaasClusterExists("ionoscloud_dbaas_cluster.test_dbaas_cluster", &dbaasCluster),
				),
			},
			{
				Config: testAccCheckDbaasClusterConfigUpdate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDbaasClusterExists("ionoscloud_dbaas_cluster.test_dbaas_cluster", &dbaasCluster),
					resource.TestCheckResourceAttr("ionoscloud_dbaas_cluster.test_dbaas_cluster", "postgres_version", "13"),
					resource.TestCheckResourceAttr("ionoscloud_dbaas_cluster.test_dbaas_cluster", "storage_size", "1.5Gi"),
					resource.TestCheckResourceAttr("ionoscloud_dbaas_cluster.test_dbaas_cluster", "display_name", "PostgreSQL_cluster_update"),
					resource.TestCheckResourceAttr("ionoscloud_dbaas_cluster.test_dbaas_cluster", "backup_enabled", "true"),
					resource.TestCheckResourceAttr("ionoscloud_dbaas_cluster.test_dbaas_cluster", "maintenance_window.0.time", "10:00:00"),
					resource.TestCheckResourceAttr("ionoscloud_dbaas_cluster.test_dbaas_cluster", "maintenance_window.0.weekday", "Saturday"),
				),
			},
		},
	})
}

func testAccCheckDbaasClusterDestroyCheck(s *terraform.State) error {
	client := testAccProvider.Meta().(SdkBundle).DbaasClient

	ctx, cancel := context.WithTimeout(context.Background(), *resourceDefaultTimeouts.Default)

	if cancel != nil {
		defer cancel()
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "ionoscloud_dbaas_cluster" {
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

func testAccCheckDbaasClusterExists(n string, cluster *dbaas.Cluster) resource.TestCheckFunc {
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

const testAccCheckDbaasClusterConfigBasic = `
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

resource "ionoscloud_dbaas_cluster" "test_dbaas_cluster" {
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
}`

const testAccCheckDbaasClusterConfigUpdate = `
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

resource "ionoscloud_ipblock" "test_dbaas_cluster" {
  location = ionoscloud_datacenter.test_dbaas_cluster.location
  size = 1
  name = "test_dbaas_cluster"
}


resource "ionoscloud_dbaas_cluster" "test_dbaas_cluster" {
  postgres_version   = 13
  replicas           = 2
  cpu_core_count     = 4
  ram_size           = "2Gi"
  storage_size       = "1.5Gi"
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
}`
