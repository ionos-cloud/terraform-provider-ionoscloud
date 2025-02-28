//go:build compute || all || server

package ionoscloud

import (
	"context"
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"

	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/bundleclient"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils/constant"
)

const serverTestResource2 = constant.ServerTestResource + "2"

func TestAccDataSourceServersBasic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ExternalProviders:        randomProviderVersion343(),
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactoriesInternal(t, &testAccProvider),
		CheckDestroy:             testAccCheckServersDestroyCheck,
		Steps: []resource.TestStep{
			{
				Config: testAccCheck2ServersByNameAndCores,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(constant.DataSource+"."+constant.ServersDataSource+"."+constant.ServerDataSourceByName, "servers.#", "1"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.ServersDataSource+"."+constant.ServerDataSourceByName, "servers.0.name",
						constant.ServerResource+"."+serverTestResource2, "name"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.ServersDataSource+"."+constant.ServerDataSourceByName, "servers.0.cores",
						constant.ServerResource+"."+serverTestResource2, "cores"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.ServersDataSource+"."+constant.ServerDataSourceByName, "servers.0.ram",
						constant.ServerResource+"."+serverTestResource2, "ram"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.ServersDataSource+"."+constant.ServerDataSourceByName, "servers.0.type",
						constant.ServerResource+"."+serverTestResource2, "type"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.ServersDataSource+"."+constant.ServerDataSourceByName, "servers.0.type",
						constant.ServerResource+"."+serverTestResource2, "type"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.ServersDataSource+"."+constant.ServerDataSourceByName, "servers.0.volumes.0.name",
						constant.ServerResource+"."+serverTestResource2, "volume.0.name"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.ServersDataSource+"."+constant.ServerDataSourceByName, "servers.0.volumes.0.size",
						constant.ServerResource+"."+serverTestResource2, "volume.0.size"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.ServersDataSource+"."+constant.ServerDataSourceByName, "servers.0.volumes.0.disk_type",
						constant.ServerResource+"."+serverTestResource2, "volume.0.disk_type"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.ServersDataSource+"."+constant.ServerDataSourceByName, "servers.0.volumes.0.bus",
						constant.ServerResource+"."+serverTestResource2, "volume.0.bus"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.ServersDataSource+"."+constant.ServerDataSourceByName, "servers.0.nics.0.name",
						constant.ServerResource+"."+serverTestResource2, "nic.0.name"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.ServersDataSource+"."+constant.ServerDataSourceByName, "servers.0.nics.0.lan",
						constant.ServerResource+"."+serverTestResource2, "nic.0.lan"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.ServersDataSource+"."+constant.ServerDataSourceByName, "servers.0.nics.0.dhcp",
						constant.ServerResource+"."+serverTestResource2, "nic.0.dhcp"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.ServersDataSource+"."+constant.ServerDataSourceByName, "servers.0.nics.0.firewall_active",
						constant.ServerResource+"."+serverTestResource2, "nic.0.firewall_active"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.ServersDataSource+"."+constant.ServerDataSourceByName, "servers.0.nics.0.firewall_rules.0.name",
						constant.ServerResource+"."+serverTestResource2, "nic.0.firewall.0.name"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.ServersDataSource+"."+constant.ServerDataSourceByName, "servers.0.nics.0.firewall_rules.0.protocol",
						constant.ServerResource+"."+serverTestResource2, "nic.0.firewall.0.protocol"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.ServersDataSource+"."+constant.ServerDataSourceByName, "servers.0.nics.0.firewall_rules.0.port_range_start",
						constant.ServerResource+"."+serverTestResource2, "nic.0.firewall.0.port_range_start"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.ServersDataSource+"."+constant.ServerDataSourceByName, "servers.0.nics.0.firewall_rules.0.port_range_end",
						constant.ServerResource+"."+serverTestResource2, "nic.0.firewall.0.port_range_end"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.ServersDataSource+"."+constant.ServerDataSourceByName, "servers.0.nics.0.firewall_rules.0.source_mac",
						constant.ServerResource+"."+serverTestResource2, "nic.0.firewall.0.source_mac"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.ServersDataSource+"."+constant.ServerDataSourceByName, "servers.0.nics.0.firewall_rules.0.source_ip",
						constant.ServerResource+"."+serverTestResource2, "nic.0.firewall.0.source_ip"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.ServersDataSource+"."+constant.ServerDataSourceByName, "servers.0.nics.0.firewall_rules.0.source_ip",
						constant.ServerResource+"."+serverTestResource2, "nic.0.firewall.0.source_ip"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.ServersDataSource+"."+constant.ServerDataSourceByName, "servers.0.nics.0.firewall_rules.0.target_ip",
						constant.ServerResource+"."+serverTestResource2, "nic.0.firewall.0.target_ip"),
				),
			},
			{
				Config: testAccCheck2ServersByCpuFamily,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(constant.DataSource+"."+constant.ServersDataSource+"."+constant.ServerDataSourceByName, "servers.#", "2"),
					// Check server labels.
					resource.TestCheckResourceAttr(constant.DataSource+"."+constant.ServersDataSource+"."+constant.ServerDataSourceByName, "servers.0.labels.#", "2"),
					resource.TestCheckResourceAttr(constant.DataSource+"."+constant.ServersDataSource+"."+constant.ServerDataSourceByName, "servers.0.labels.0.key", "labelkey0"),
					resource.TestCheckResourceAttr(constant.DataSource+"."+constant.ServersDataSource+"."+constant.ServerDataSourceByName, "servers.0.labels.0.value", "labelvalue0"),
					resource.TestCheckResourceAttr(constant.DataSource+"."+constant.ServersDataSource+"."+constant.ServerDataSourceByName, "servers.0.labels.1.key", "labelkey1"),
					resource.TestCheckResourceAttr(constant.DataSource+"."+constant.ServersDataSource+"."+constant.ServerDataSourceByName, "servers.0.labels.1.value", "labelvalue1"),
					resource.TestCheckResourceAttr(constant.DataSource+"."+constant.ServersDataSource+"."+constant.ServerDataSourceByName, "servers.1.labels.0.key", "labelkey0"),
					resource.TestCheckResourceAttr(constant.DataSource+"."+constant.ServersDataSource+"."+constant.ServerDataSourceByName, "servers.1.labels.0.value", "labelvalue0"),
					resource.TestCheckResourceAttr(constant.DataSource+"."+constant.ServersDataSource+"."+constant.ServerDataSourceByName, "servers.1.labels.1.key", "labelkey1"),
					resource.TestCheckResourceAttr(constant.DataSource+"."+constant.ServersDataSource+"."+constant.ServerDataSourceByName, "servers.1.labels.1.value", "labelvalue1"),
				),
			},
			{
				Config:      testAccCheck2ServersBadFilter,
				ExpectError: regexp.MustCompile("no servers found for criteria, please check your filter configuration"),
			},
		},
	})
}

func testAccCheckServersDestroyCheck(s *terraform.State) error {
	client := testAccProvider.Meta().(bundleclient.SdkBundle).CloudApiClient

	ctx, cancel := context.WithTimeout(context.Background(), *resourceDefaultTimeouts.Default)

	if cancel != nil {
		defer cancel()
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != constant.ServerResource {
			continue
		}

		dcId := rs.Primary.Attributes["datacenter_id"]

		_, apiResponse, err := client.ServersApi.DatacentersServersFindById(ctx, dcId, rs.Primary.ID).Execute()
		logApiRequestTime(apiResponse)

		if err != nil {
			if !httpNotFound(apiResponse) {
				return fmt.Errorf("unable to fetch server %s: %w", rs.Primary.ID, err)
			}
		} else {
			return fmt.Errorf("server still exists %s", rs.Primary.ID)
		}
	}

	return nil
}

const testAccCheck2ServersByNameAndCores = `
resource ` + constant.DatacenterResource + ` ` + constant.DatacenterTestResource + ` {
	name       = "server-test"
	location = "us/las"
}

resource ` + constant.LanResource + ` ` + constant.LanTestResource + ` {
  datacenter_id = ` + constant.DatacenterResource + `.` + constant.DatacenterTestResource + `.id
  public = true
  name = "public"
}
resource ` + constant.ServerResource + ` ` + constant.ServerTestResource + ` {
  name = "` + constant.ServerTestResource + `"
  datacenter_id = ` + constant.DatacenterResource + `.` + constant.DatacenterTestResource + `.id
  cores = 2
  ram = 2048
  image_name ="ubuntu:latest"
  image_password = ` + constant.RandomPassword + `.server_image_password.result
  type = "ENTERPRISE"
  volume {
    name = "` + constant.VolumeTestResource + `"
    size = 6
    disk_type = "SSD Standard"
    user_data = "foo"
    bus = "IDE"
	}
  nic {
    lan = ` + constant.LanResource + `.` + constant.LanTestResource + `.id
    name = "` + constant.LanTestResource + `"
    dhcp = false
    firewall_active = false
  }
}

resource "ionoscloud_ipblock" "webserver_ipblock" {
  location = ` + constant.DatacenterResource + `.` + constant.DatacenterTestResource + `.location
  size = 2
  name = "webserver_ipblock"
}

resource ` + constant.ServerResource + ` ` + serverTestResource2 + ` {
  name = "` + serverTestResource2 + `"
  datacenter_id = ` + constant.DatacenterResource + `.` + constant.DatacenterTestResource + `.id
  cores = ` + noCoresTest + `
  ram = 2048
  image_name ="ubuntu:latest"
  image_password = ` + constant.RandomPassword + `.server2_image_password.result
  type = "ENTERPRISE"
  volume {
    name = "` + constant.VolumeTestResource + "2" + `"
    size = 6
    disk_type = "SSD Standard"
    user_data = "foo"
    bus = "IDE"
	}
  nic {
    lan = 1
    name = "` + constant.LanTestResource + "2" + `"
    dhcp = false
    firewall_active = false
    firewall {
      protocol = "TCP"
      name = "SSH"
      port_range_start = 22
      port_range_end = 22
	  source_mac = "00:0a:95:9d:68:17"
	  source_ip = ionoscloud_ipblock.webserver_ipblock.ips[0]
	  target_ip = ionoscloud_ipblock.webserver_ipblock.ips[1]
	  type = "EGRESS"
    }
  }
}

` + ServerImagePassword + `
resource ` + constant.RandomPassword + ` "server2_image_password" {
  length           = 16
  special          = false
}

data ` + constant.ServersDataSource + ` ` + constant.ServerDataSourceByName + ` {
 depends_on = [` + constant.ServerResource + `.` + serverTestResource2 + `,
	` + constant.ServerResource + `.` + constant.ServerTestResource + `]
  datacenter_id = ` + constant.DatacenterResource + `.` + constant.DatacenterTestResource + `.id
  filter {
   name = "name"
   value = "${ionoscloud_server.test_server2.name}" 
  }
  filter {
    name = "cores"
    value = "` + noCoresTest + `" 
  }
} `

const testAccCheck2ServersByCpuFamily = `
resource ` + constant.DatacenterResource + ` ` + constant.DatacenterTestResource + ` {
	name       = "server-test"
	location = "us/las"
}

resource ` + constant.LanResource + ` ` + constant.LanTestResource + ` {
  datacenter_id = ` + constant.DatacenterResource + `.` + constant.DatacenterTestResource + `.id
  public = true
  name = "public"
}

resource ` + constant.ServerResource + ` ` + constant.ServerTestResource + ` {
  name = "` + constant.ServerTestResource + `"
  datacenter_id = ` + constant.DatacenterResource + `.` + constant.DatacenterTestResource + `.id
  cores = 2
  ram = 2048
  image_name ="ubuntu:latest"
  image_password = ` + constant.RandomPassword + `.server_image_password.result
  type = "ENTERPRISE"
  volume {
    name = "` + constant.VolumeTestResource + `"
    size = 6
    disk_type = "SSD Standard"
    user_data = "foo"
    bus = "IDE"
	}
  nic {
    lan = ` + constant.LanResource + `.` + constant.LanTestResource + `.id
    name = "` + constant.LanTestResource + `"
    dhcp = false
    firewall_active = false
  }
  label {
    key = "labelkey0"
    value = "labelvalue0"
  }
  label {
    key = "labelkey1"
    value = "labelvalue1"
  }
}

resource "ionoscloud_ipblock" "webserver_ipblock" {
  location = ` + constant.DatacenterResource + `.` + constant.DatacenterTestResource + `.location
  size = 2
  name = "webserver_ipblock"
}

resource ` + constant.ServerResource + ` ` + serverTestResource2 + ` {
  name = "` + serverTestResource2 + `"
  datacenter_id = ` + constant.DatacenterResource + `.` + constant.DatacenterTestResource + `.id
  cores = ` + noCoresTest + `
  ram = 2048
  image_name ="ubuntu:latest"
  image_password = ` + constant.RandomPassword + `.server2_image_password.result
  type = "ENTERPRISE"
  volume {
    name = "` + constant.VolumeTestResource + "2" + `"
    size = 6
    disk_type = "SSD Standard"
    user_data = "foo"
    bus = "IDE"
	}
  nic {
    lan = 1
    name = "` + constant.LanTestResource + "2" + `"
    dhcp = false
    firewall_active = false
    firewall {
      protocol = "TCP"
      name = "SSH"
      port_range_start = 22
      port_range_end = 22
	  source_mac = "00:0a:95:9d:68:17"
	  source_ip = ionoscloud_ipblock.webserver_ipblock.ips[0]
	  target_ip = ionoscloud_ipblock.webserver_ipblock.ips[1]
	  type = "EGRESS"
    }
  }
  label {
    key = "labelkey0"
    value = "labelvalue0"
  }
  label {
    key = "labelkey1"
    value = "labelvalue1"
  }
}

` + ServerImagePassword + `
resource ` + constant.RandomPassword + ` "server2_image_password" {
  length           = 16
  special          = false
}

data ` + constant.ServersDataSource + ` ` + constant.ServerDataSourceByName + ` {
 depends_on = [` + constant.ServerResource + `.` + constant.ServerTestResource + "2" + `,
	` + constant.ServerResource + `.` + constant.ServerTestResource + `]
  datacenter_id = ` + constant.DatacenterResource + `.` + constant.DatacenterTestResource + `.id
  filter {
   name = "cpu_family"
   value = "` + cpuFamilyTest + `" 
  }
} `

const testAccCheck2ServersBadFilter = `
resource ` + constant.DatacenterResource + ` ` + constant.DatacenterTestResource + ` {
	name       = "server-test"
	location = "us/las"
}

resource ` + constant.LanResource + ` ` + constant.LanTestResource + ` {
  datacenter_id = ` + constant.DatacenterResource + `.` + constant.DatacenterTestResource + `.id
  public = true
  name = "public"
}
resource ` + constant.ServerResource + ` ` + constant.ServerTestResource + ` {
  name = "` + constant.ServerTestResource + `"
  datacenter_id = ` + constant.DatacenterResource + `.` + constant.DatacenterTestResource + `.id
  cores = 2
  ram = 2048
  image_name ="ubuntu:latest"
  image_password = ` + constant.RandomPassword + `.server_image_password.result
  type = "ENTERPRISE"
  volume {
    name = "` + constant.VolumeTestResource + `"
    size = 6
    disk_type = "SSD Standard"
    user_data = "foo"
    bus = "IDE"
	}
  nic {
    lan = ` + constant.LanResource + `.` + constant.LanTestResource + `.id
    name = "` + constant.LanTestResource + `"
    dhcp = false
    firewall_active = false
  }
}

resource ` + constant.ServerResource + ` ` + constant.ServerTestResource + "2" + ` {
  name = "` + constant.ServerTestResource + "2" + `"
  datacenter_id = ` + constant.DatacenterResource + `.` + constant.DatacenterTestResource + `.id
  cores = ` + noCoresTest + `
  ram = 2048
  image_name ="ubuntu:latest"
  image_password = ` + constant.RandomPassword + `.server2_image_password.result
  type = "ENTERPRISE"
  volume {
    name = "` + constant.VolumeTestResource + "2" + `"
    size = 6
    disk_type = "SSD Standard"
    user_data = "foo"
    bus = "IDE"
	}
  nic {
    lan = 1
    name = "` + constant.LanTestResource + "2" + `"
    dhcp = false
    firewall_active = false
  }
}

` + ServerImagePassword + `
resource ` + constant.RandomPassword + ` "server2_image_password" {
  length           = 16
  special          = false
}

data ` + constant.ServersDataSource + ` ` + constant.ServerDataSourceByName + ` {
 depends_on = [` + constant.ServerResource + `.` + constant.ServerTestResource + "2" + `,
	` + constant.ServerResource + `.` + constant.ServerTestResource + `]
  datacenter_id = ` + constant.DatacenterResource + `.` + constant.DatacenterTestResource + `.id
  filter {
   name = "cpu_family"
   value = "doesNotExist" 
  }
} `

const cpuFamilyTest = "INTEL_XEON"
const noCoresTest = "1"
