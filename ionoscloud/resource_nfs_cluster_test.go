//go:build all || nfs || nfs_cluster

package ionoscloud

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"

	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/bundleclient"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils/constant"
)

func TestAccNFSClusterBasic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {
				VersionConstraint: "3.4.3",
				Source:            "hashicorp/random",
			},
			"time": {
				Source:            "hashicorp/time",
				VersionConstraint: "0.11.1",
			},
		},
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckNFSClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckNFSClusterConfigBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNFSClusterExists("ionoscloud_nfs_cluster.example"),
					resource.TestCheckResourceAttr("ionoscloud_nfs_cluster.example", "name", "example"),
					resource.TestCheckResourceAttr("ionoscloud_nfs_cluster.example", "location", "de/txl"),
					resource.TestCheckResourceAttr("ionoscloud_nfs_cluster.example", "size", "2"),
					resource.TestCheckResourceAttr("ionoscloud_nfs_cluster.example", "nfs.0.min_version", "4.2"),
				),
			},
			{
				Config: testAccCheckNFSClusterConfigUpdate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNFSClusterExists("ionoscloud_nfs_cluster.example"),
					resource.TestCheckResourceAttr("ionoscloud_nfs_cluster.example", "name", "example_updated"),
					resource.TestCheckResourceAttr("ionoscloud_nfs_cluster.example", "location", "de/txl"),
					resource.TestCheckResourceAttr("ionoscloud_nfs_cluster.example", "size", "2"),
					resource.TestCheckResourceAttr("ionoscloud_nfs_cluster.example", "nfs.0.min_version", "4.2"),
				),
			},
			{
				Config: testAccDataSourceNFSClusterMatchName,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair("data.ionoscloud_nfs_cluster.data_with_name", "name", "ionoscloud_nfs_cluster.example", "name"),
					resource.TestCheckResourceAttrPair("data.ionoscloud_nfs_cluster.data_with_name", "location", "ionoscloud_nfs_cluster.example", "location"),
				),
			},
			{
				Config: testAccDataSourceNFSClusterPartialMatchName,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair("data.ionoscloud_nfs_cluster.data_with_name", "name", "ionoscloud_nfs_cluster.example", "name"),
					resource.TestCheckResourceAttrPair("data.ionoscloud_nfs_cluster.data_with_name", "location", "ionoscloud_nfs_cluster.example", "location"),
				),
			},
		},
	})
}

func testAccCheckNFSClusterDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(bundleclient.SdkBundle).NFSClient

	for _, rs := range s.RootModule().Resources {
		if rs.Type != constant.NFSClusterResource {
			continue
		}

		_, resp, err := client.GetNFSClusterByID(context.Background(), rs.Primary.ID, rs.Primary.Attributes["location"])
		if resp != nil && resp.StatusCode != 404 {
			return fmt.Errorf("NFS Cluster still exists: %s", rs.Primary.ID)
		}
		if err != nil {
			return fmt.Errorf("error fetching NFS Cluster with ID %s: %v", rs.Primary.ID, err)
		}
	}

	return nil
}

func testAccCheckNFSClusterExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(bundleclient.SdkBundle).NFSClient

		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no ID is set")
		}

		ctx, cancel := context.WithTimeout(context.Background(), *resourceDefaultTimeouts.Default)
		defer cancel()

		found, _, err := client.GetNFSClusterByID(ctx, rs.Primary.ID, rs.Primary.Attributes["location"])
		if err != nil {
			return fmt.Errorf("an error occurred while fetching NFS Cluster with ID: %v, error: %w", rs.Primary.ID, err)
		}
		if *found.Id != rs.Primary.ID {
			return fmt.Errorf("resource not found")
		}

		return nil
	}
}

const testAccCheckNFSClusterConfig = `
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
`

const testAccCheckNFSClusterConfigBasic = testAccCheckNFSClusterConfig + `
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

const testAccCheckNFSClusterConfigUpdate = testAccCheckNFSClusterConfig + `
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

const testAccDataSourceNFSClusterMatchName = testAccCheckNFSClusterConfigUpdate + `
data "ionoscloud_nfs_cluster" "data_with_name" {
  location = ionoscloud_nfs_cluster.example.location
  name = "example_updated"
}
`

const testAccDataSourceNFSClusterPartialMatchName = testAccCheckNFSClusterConfigUpdate + `
data "ionoscloud_nfs_cluster" "data_with_name" {
  location = ionoscloud_nfs_cluster.example.location
  name = "example_"
  partial_match = true
}
`
