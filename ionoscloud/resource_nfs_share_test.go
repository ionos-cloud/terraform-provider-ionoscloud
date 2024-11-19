//go:build all || nfs || nfs_share

package ionoscloud

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services"
)

func TestAccNFSShare_basic(t *testing.T) {
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
		CheckDestroy:      testAccCheckNFSShareDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckNFSShareConfig_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNFSShareExists("ionoscloud_nfs_share.example"),
					resource.TestCheckResourceAttr("ionoscloud_nfs_share.example", "name", "example-share"),
					resource.TestCheckResourceAttr("ionoscloud_nfs_share.example", "quota", "512"),
					resource.TestCheckResourceAttr("ionoscloud_nfs_share.example", "gid", "512"),
					resource.TestCheckResourceAttr("ionoscloud_nfs_share.example", "uid", "512"),
					resource.TestCheckResourceAttr("ionoscloud_nfs_share.example", "client_groups.0.description", "Client Group 1"),
					resource.TestCheckResourceAttr("ionoscloud_nfs_share.example", "client_groups.0.ip_networks.0", "10.234.50.0/24"),
					resource.TestCheckResourceAttr("ionoscloud_nfs_share.example", "client_groups.0.hosts.0", "10.234.62.123"),
					resource.TestCheckResourceAttr("ionoscloud_nfs_share.example", "client_groups.0.nfs.0.squash", "all-anonymous"),
				),
			},
			{
				Config: testAccCheckNFSShareConfig_update,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNFSShareExists("ionoscloud_nfs_share.example"),
					resource.TestCheckResourceAttr("ionoscloud_nfs_share.example", "name", "example-share-updated"),
					resource.TestCheckResourceAttr("ionoscloud_nfs_share.example", "quota", "1024"),
					resource.TestCheckResourceAttr("ionoscloud_nfs_share.example", "gid", "1024"),
					resource.TestCheckResourceAttr("ionoscloud_nfs_share.example", "uid", "1024"),
					resource.TestCheckResourceAttr("ionoscloud_nfs_share.example", "client_groups.0.description", "Client Group 1 Updated"),
					resource.TestCheckResourceAttr("ionoscloud_nfs_share.example", "client_groups.0.ip_networks.0", "10.234.50.0/24"),
					resource.TestCheckResourceAttr("ionoscloud_nfs_share.example", "client_groups.0.hosts.0", "10.234.62.124"),
					resource.TestCheckResourceAttr("ionoscloud_nfs_share.example", "client_groups.0.nfs.0.squash", "root-anonymous"),
				),
			},
			{
				Config: testAccDataSourceNFSShareMatchId,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair("data.ionoscloud_nfs_share.share_data_example", "name", "ionoscloud_nfs_share.example", "name"),
					resource.TestCheckResourceAttrPair("data.ionoscloud_nfs_share.share_data_example", "quota", "ionoscloud_nfs_share.example", "quota"),
					resource.TestCheckResourceAttrPair("data.ionoscloud_nfs_share.share_data_example", "gid", "ionoscloud_nfs_share.example", "gid"),
					resource.TestCheckResourceAttrPair("data.ionoscloud_nfs_share.share_data_example", "uid", "ionoscloud_nfs_share.example", "uid"),
				),
			},
		},
	})
}

func testAccCheckNFSShareDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(services.SdkBundle).NFSClient

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "ionoscloud_nfs_share" {
			continue
		}

		_, resp, err := client.GetNFSShareByID(context.Background(), rs.Primary.Attributes["cluster_id"], rs.Primary.ID, rs.Primary.Attributes["location"])
		if resp != nil && resp.StatusCode != 404 {
			return fmt.Errorf("NFS Share still exists: %s", rs.Primary.ID)
		}
		if err != nil {
			return fmt.Errorf("error fetching NFS Share with ID %s: %v", rs.Primary.ID, err)
		}
	}

	return nil
}

func testAccCheckNFSShareExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(services.SdkBundle).NFSClient

		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no ID is set")
		}

		ctx, cancel := context.WithTimeout(context.Background(), *resourceDefaultTimeouts.Default)
		defer cancel()

		found, _, err := client.GetNFSShareByID(ctx, rs.Primary.Attributes["cluster_id"], rs.Primary.ID, rs.Primary.Attributes["location"])
		if err != nil {
			return fmt.Errorf("an error occurred while fetching NFS Share with ID: %v, error: %w", rs.Primary.ID, err)
		}
		if *found.Id != rs.Primary.ID {
			return fmt.Errorf("resource not found")
		}

		return nil
	}
}

const testAccCheckNFSShareConfig_basic = testAccCheckNFSClusterConfig_basic + `
resource "ionoscloud_nfs_share" "example" {
  location = ionoscloud_nfs_cluster.example.location
  cluster_id = ionoscloud_nfs_cluster.example.id

  name = "example-share"
  quota = 512
  gid = 512
  uid = 512

  client_groups {
    description = "Client Group 1"
    ip_networks = ["10.234.50.0/24"]
    hosts = ["10.234.62.123"]

    nfs {
      squash = "all-anonymous"
    }
  }
}
`

const testAccCheckNFSShareConfig_update = testAccCheckNFSClusterConfig_basic + `
resource "ionoscloud_nfs_share" "example" {
  location = ionoscloud_nfs_cluster.example.location
  cluster_id = ionoscloud_nfs_cluster.example.id

  name = "example-share-updated"
  quota = 1024
  gid = 1024
  uid = 1024

  client_groups {
    description = "Client Group 1 Updated"
    ip_networks = ["10.234.50.0/24"]
    hosts = ["10.234.62.124"]

    nfs {
      squash = "root-anonymous"
    }
  }
}
`

const testAccDataSourceNFSShareMatchId = testAccCheckNFSShareConfig_update + `
data "ionoscloud_nfs_share" "share_data_example" {
  location = ionoscloud_nfs_cluster.example.location
  cluster_id = ionoscloud_nfs_cluster.example.id
  id = ionoscloud_nfs_share.example.id
}
`
