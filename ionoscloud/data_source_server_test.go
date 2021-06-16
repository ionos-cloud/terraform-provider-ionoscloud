package ionoscloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceServer_matchId(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceServerCreateResources,
			},
			{
				Config: testAccDataSourceServerMatchId,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.ionoscloud_server.test_server", "name", "test_datasource_server"),
				),
			},
		},
	})
}

func TestAccDataSourceServer_matchName(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceServerCreateResources,
			},
			{
				Config: testAccDataSourceServerMatchName,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.ionoscloud_server.test_server", "name", "test_datasource_server"),
				),
			},
		},
	})

}

const testAccDataSourceServerCreateResources = `
resource "ionoscloud_datacenter" "test_datasource_server" {
  name              = "test_datasource_server"
  location          = "de/fra"
  description       = "datacenter for testing the server terraform data source"
}
resource "ionoscloud_server" "test_datasource_server" {
  depends_on        = [ionoscloud_datacenter.test_datasource_server]
  datacenter_id     = ionoscloud_datacenter.test_datasource_server.id
  name              = "test_datasource_server"
  cores             = 2
  ram               = 2048
  availability_zone = "ZONE_1"
  cpu_family        = "INTEL_SKYLAKE"

  image_name        = "ubuntu:20.04"
  image_password    = "foobar123456"

  volume {
    size            =   "40"
    disk_type       =   "HDD"
  }

  nic {
    lan             = 1
  }
}
`

const testAccDataSourceServerMatchId = `
resource "ionoscloud_datacenter" "test_datasource_server" {
  name              = "test_datasource_server"
  location          = "de/fra"
  description       = "datacenter for testing the server terraform data source"
}

resource "ionoscloud_server" "test_datasource_server" {
  depends_on        = [ionoscloud_datacenter.test_datasource_server]
  datacenter_id     = ionoscloud_datacenter.test_datasource_server.id
  name              = "test_datasource_server"
  cores             = 2
  ram               = 2048
  availability_zone = "ZONE_1"
  cpu_family        = "INTEL_SKYLAKE"

  image_name        = "ubuntu:20.04"
  image_password    = "foobar123456"

  volume {
    size            =   "40"
    disk_type       =   "HDD"
  }

  nic {
    lan             = 1
  }
}

data "ionoscloud_server" "test_server" {
  datacenter_id = ionoscloud_datacenter.test_datasource_server.id
  id			= ionoscloud_server.test_datasource_server.id
}
`

const testAccDataSourceServerMatchName = `
resource "ionoscloud_datacenter" "test_datasource_server" {
  name              = "test_datasource_server"
  location          = "de/fra"
  description       = "datacenter for testing the server terraform data source"
}

resource "ionoscloud_server" "test_datasource_server" {
  depends_on        = [ionoscloud_datacenter.test_datasource_server]
  datacenter_id     = ionoscloud_datacenter.test_datasource_server.id
  name              = "test_datasource_server"

  cores             = 2
  ram               = 2048
  availability_zone = "ZONE_1"
  cpu_family        = "INTEL_SKYLAKE"

  image_name        = "ubuntu:20.04"
  image_password    = "foobar123456"

  volume {
    size            =   "40"
    disk_type       =   "HDD"
  }

  nic {
    lan             = 1
  }
}

data "ionoscloud_server" "test_server" {
  datacenter_id = ionoscloud_datacenter.test_datasource_server.id
  name			= "test_datasource_server"
}
`
