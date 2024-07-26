//go:build nfs || all || nfs_cluster

package ionoscloud

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	ionoscloud "github.com/ionos-cloud/sdk-go-nfs"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/internal/tftest"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services"
)

var testAccCheckNFSClusterConfig_basic = `
resource "ionoscloud_datacenter" "nfs_dc" {
  name                = "NFS Datacenter"
  location            = "de/txl"
  description         = "Datacenter Description"
  sec_auth_protection = false
}

resource "ionoscloud_lan" "nfs_lan" {
  datacenter_id = ionoscloud_datacenter.nfs_dc.id
  public        = false
  name          = "Lan for NFS"
}

data "ionoscloud_image" "HDD_image" {
  image_alias = "ubuntu:20.04"
  type       = "HDD"
  cloud_init = "V1"
  location   = "de/txl"
}

resource "random_password" "password" {
  length  = 16
  special = false
}

resource "ionoscloud_server" "nfs_server" {
  name              = "Server for NFS"
  datacenter_id     = ionoscloud_datacenter.nfs_dc.id
  cores             = 1
  ram               = 2048
  availability_zone = "ZONE_1"
  cpu_family        = "INTEL_SKYLAKE"
  image_name        = data.ionoscloud_image.HDD_image.id
  image_password    = random_password.password.result
  volume {
    name      = "system"
    size      = 14
    disk_type = "SSD"
  }
  nic {
    name            = "NIC A"
    lan             = ionoscloud_lan.nfs_lan.id
    dhcp            = true
    firewall_active = true
  }
}

resource "ionoscloud_nfs_cluster" "example" {
  name = "example"
  location = "de/txl"
  size = 2
  nfs {
    min_version = "4.2"
  }
  connections {
    datacenter_id = ionoscloud_datacenter.nfs_dc.id
    ip_address    = format("%s/24", ionoscloud_server.nfs_server.nic[0].ips[0])
    lan           = ionoscloud_lan.nfs_lan.id
  }
}
`

var testAccCheckNFSClusterConfig_update = `
resource "ionoscloud_datacenter" "nfs_dc" {
  name                = "NFS Datacenter"
  location            = "de/txl"
  description         = "Datacenter Description"
  sec_auth_protection = false
}

resource "ionoscloud_lan" "nfs_lan" {
  datacenter_id = ionoscloud_datacenter.nfs_dc.id
  public        = false
  name          = "Lan for NFS"
}

data "ionoscloud_image" "HDD_image" {
  image_alias = "ubuntu:20.04"
  type       = "HDD"
  cloud_init = "V1"
  location   = "de/txl"
}

resource "random_password" "password" {
  length  = 16
  special = false
}

resource "ionoscloud_server" "nfs_server" {
  name              = "Server for NFS"
  datacenter_id     = ionoscloud_datacenter.nfs_dc.id
  cores             = 1
  ram               = 2048
  availability_zone = "ZONE_1"
  cpu_family        = "INTEL_SKYLAKE"
  image_name        = data.ionoscloud_image.HDD_image.id
  image_password    = random_password.password.result
  volume {
    name      = "system"
    size      = 14
    disk_type = "SSD"
  }
  nic {
    name            = "NIC A"
    lan             = ionoscloud_lan.nfs_lan.id
    dhcp            = true
    firewall_active = true
  }
}

resource "ionoscloud_nfs_cluster" "example" {
  name = "example_updated"
  location = "de/txl"
  size = 2
  nfs {
    min_version = "4.2"
  }
  connections {
    datacenter_id = ionoscloud_datacenter.nfs_dc.id
    ip_address    = format("%s/24", ionoscloud_server.nfs_server.nic[0].ips[0])
    lan           = ionoscloud_lan.nfs_lan.id
  }
}
`

const nfsDataSourceMatchName = `
data "ionoscloud_nfs_cluster" "data_with_name" {
  name = "example_updated"
  location = "de/txl"
}
`

const nfsDataSourcePartialMatchName = `
data "ionoscloud_nfs_cluster" "data_with_name" {
  name = "example_upd"
  location = "de/txl"
  partial_match = true
}
`

func TestAccNFSCluster_basic(t *testing.T) {
	var nfsCluster ionoscloud.ClusterRead

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {
				VersionConstraint: "3.4.3",
				Source:            "hashicorp/random",
			},
		},
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckNFSClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckNFSClusterConfig_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNFSClusterExists("ionoscloud_nfs_cluster.example", &nfsCluster),
					tftest.TestCheckExpectedAttrs("ionoscloud_nfs_cluster.example", map[string]string{
						"name":                        "example",
						"location":                    "de/txl",
						"size":                        "2",
						"nfs.0.min_version":           "4.2",
						"connections.0.datacenter_id": "dynamic:ionoscloud_datacenter.nfs_dc.id",
						"connections.0.ip_address":    "dynamic:ionoscloud_server.nfs_server.nic.0.ips.0",
						"connections.0.lan":           "dynamic:ionoscloud_lan.nfs_lan.id",
					}),
				),
			},
			{
				Config: testAccCheckNFSClusterConfig_update,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNFSClusterExists("ionoscloud_nfs_cluster.example", &nfsCluster),
					tftest.TestCheckExpectedAttrs("ionoscloud_nfs_cluster.example", map[string]string{
						"name":                        "example_updated",
						"location":                    "de/txl",
						"size":                        "2",
						"nfs.0.min_version":           "4.2",
						"connections.0.datacenter_id": "dynamic:ionoscloud_datacenter.nfs_dc.id",
						"connections.0.ip_address":    "dynamic:ionoscloud_server.nfs_server.nic.0.ips.0",
						"connections.0.lan":           "dynamic:ionoscloud_lan.nfs_lan.id",
					}),
				),
			},
			{
				Config: nfsDataSourceMatchName,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNFSClusterExists("data.ionoscloud_nfs_cluster.data_with_name", &nfsCluster),
					tftest.TestCheckExpectedAttrs("data.ionoscloud_nfs_cluster.data_with_name", map[string]string{
						"name":                        "example_updated",
						"location":                    "de/txl",
						"size":                        "2",
						"nfs.0.min_version":           "4.2",
						"connections.0.datacenter_id": "dynamic:ionoscloud_datacenter.nfs_dc.id",
						"connections.0.ip_address":    "dynamic:ionoscloud_server.nfs_server.nic.0.ips.0",
						"connections.0.lan":           "dynamic:ionoscloud_lan.nfs_lan.id",
					}),
				),
			},
			{
				Config: nfsDataSourcePartialMatchName,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNFSClusterExists("data.ionoscloud_nfs_cluster.data_with_name", &nfsCluster),
					tftest.TestCheckExpectedAttrs("data.ionoscloud_nfs_cluster.data_with_name", map[string]string{
						"name":                        "example_updated",
						"location":                    "de/txl",
						"size":                        "2",
						"nfs.0.min_version":           "4.2",
						"connections.0.datacenter_id": "dynamic:ionoscloud_datacenter.nfs_dc.id",
						"connections.0.ip_address":    "dynamic:ionoscloud_server.nfs_server.nic.0.ips.0",
						"connections.0.lan":           "dynamic:ionoscloud_lan.nfs_lan.id",
					}),
				),
			},
		},
	})
}

func testAccCheckNFSClusterDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(services.SdkBundle).NFSClient

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "ionoscloud_nfs_cluster" {
			continue
		}

		_, resp, err := client.GetNFSClusterByID(context.Background(), rs.Primary.ID, "de/txl")
		if resp != nil && resp.StatusCode != 404 {
			return fmt.Errorf("NFS Cluster still exists: %s", rs.Primary.ID)
		}
		if err != nil {
			return fmt.Errorf("error fetching NFS Cluster with ID %s: %v", rs.Primary.ID, err)
		}
	}

	return nil
}

func testAccCheckNFSClusterExists(n string, nfsCluster *ionoscloud.ClusterRead) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		client := testAccProvider.Meta().(services.SdkBundle).NFSClient
		found, _, err := client.GetNFSClusterByID(context.Background(), rs.Primary.ID, "de/txl")
		if err != nil {
			return fmt.Errorf("error fetching NFS Cluster with ID %s: %v", rs.Primary.ID, err)
		}

		if found.Id != nil && *found.Id != rs.Primary.ID {
			return fmt.Errorf("NFS Cluster not found")
		}

		*nfsCluster = found
		return nil
	}
}
