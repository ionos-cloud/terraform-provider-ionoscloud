//go:build compute || all || servers

package ionoscloud

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"regexp"
	"testing"
)

const serverTestResource2 = ServerTestResource + "2"

func TestAccDataSourceServersBasic(t *testing.T) {
	t.Skip("problem with 500 error thrown by backend")
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ExternalProviders: randomProviderVersion343(),
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckServersDestroyCheck,
		Steps: []resource.TestStep{
			{
				Config: testAccCheck2ServersByNameAndCores,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(DataSource+"."+ServersDataSource+"."+ServerDataSourceByName, "servers.#", "1"),
					resource.TestCheckResourceAttrPair(DataSource+"."+ServersDataSource+"."+ServerDataSourceByName, "servers.0.name",
						ServerResource+"."+serverTestResource2, "name"),
					resource.TestCheckResourceAttrPair(DataSource+"."+ServersDataSource+"."+ServerDataSourceByName, "servers.0.cores",
						ServerResource+"."+serverTestResource2, "cores"),
					resource.TestCheckResourceAttrPair(DataSource+"."+ServersDataSource+"."+ServerDataSourceByName, "servers.0.ram",
						ServerResource+"."+serverTestResource2, "ram"),
					resource.TestCheckResourceAttrPair(DataSource+"."+ServersDataSource+"."+ServerDataSourceByName, "servers.0.availability_zone",
						ServerResource+"."+serverTestResource2, "availability_zone"),
					resource.TestCheckResourceAttrPair(DataSource+"."+ServersDataSource+"."+ServerDataSourceByName, "servers.0.cpu_family",
						ServerResource+"."+serverTestResource2, "cpu_family"),
					resource.TestCheckResourceAttrPair(DataSource+"."+ServersDataSource+"."+ServerDataSourceByName, "servers.0.type",
						ServerResource+"."+serverTestResource2, "type"),
					resource.TestCheckResourceAttrPair(DataSource+"."+ServersDataSource+"."+ServerDataSourceByName, "servers.0.type",
						ServerResource+"."+serverTestResource2, "type"),
					resource.TestCheckResourceAttrPair(DataSource+"."+ServersDataSource+"."+ServerDataSourceByName, "servers.0.volumes.0.name",
						ServerResource+"."+serverTestResource2, "volume.0.name"),
					resource.TestCheckResourceAttrPair(DataSource+"."+ServersDataSource+"."+ServerDataSourceByName, "servers.0.volumes.0.size",
						ServerResource+"."+serverTestResource2, "volume.0.size"),
					resource.TestCheckResourceAttrPair(DataSource+"."+ServersDataSource+"."+ServerDataSourceByName, "servers.0.volumes.0.disk_type",
						ServerResource+"."+serverTestResource2, "volume.0.disk_type"),
					resource.TestCheckResourceAttrPair(DataSource+"."+ServersDataSource+"."+ServerDataSourceByName, "servers.0.volumes.0.bus",
						ServerResource+"."+serverTestResource2, "volume.0.bus"),
					resource.TestCheckResourceAttrPair(DataSource+"."+ServersDataSource+"."+ServerDataSourceByName, "servers.0.volumes.0.availability_zone",
						ServerResource+"."+serverTestResource2, "volume.0.availability_zone"),
					resource.TestCheckResourceAttrPair(DataSource+"."+ServersDataSource+"."+ServerDataSourceByName, "servers.0.nics.0.name",
						ServerResource+"."+serverTestResource2, "nic.0.name"),
					resource.TestCheckResourceAttrPair(DataSource+"."+ServersDataSource+"."+ServerDataSourceByName, "servers.0.nics.0.lan",
						ServerResource+"."+serverTestResource2, "nic.0.lan"),
					resource.TestCheckResourceAttrPair(DataSource+"."+ServersDataSource+"."+ServerDataSourceByName, "servers.0.nics.0.dhcp",
						ServerResource+"."+serverTestResource2, "nic.0.dhcp"),
					resource.TestCheckResourceAttrPair(DataSource+"."+ServersDataSource+"."+ServerDataSourceByName, "servers.0.nics.0.firewall_active",
						ServerResource+"."+serverTestResource2, "nic.0.firewall_active"),
					resource.TestCheckResourceAttrPair(DataSource+"."+ServersDataSource+"."+ServerDataSourceByName, "servers.0.nics.0.firewall_rules.0.name",
						ServerResource+"."+serverTestResource2, "nic.0.firewall.0.name"),
					resource.TestCheckResourceAttrPair(DataSource+"."+ServersDataSource+"."+ServerDataSourceByName, "servers.0.nics.0.firewall_rules.0.protocol",
						ServerResource+"."+serverTestResource2, "nic.0.firewall.0.protocol"),
					resource.TestCheckResourceAttrPair(DataSource+"."+ServersDataSource+"."+ServerDataSourceByName, "servers.0.nics.0.firewall_rules.0.port_range_start",
						ServerResource+"."+serverTestResource2, "nic.0.firewall.0.port_range_start"),
					resource.TestCheckResourceAttrPair(DataSource+"."+ServersDataSource+"."+ServerDataSourceByName, "servers.0.nics.0.firewall_rules.0.port_range_end",
						ServerResource+"."+serverTestResource2, "nic.0.firewall.0.port_range_end"),
					resource.TestCheckResourceAttrPair(DataSource+"."+ServersDataSource+"."+ServerDataSourceByName, "servers.0.nics.0.firewall_rules.0.source_mac",
						ServerResource+"."+serverTestResource2, "nic.0.firewall.0.source_mac"),
					resource.TestCheckResourceAttrPair(DataSource+"."+ServersDataSource+"."+ServerDataSourceByName, "servers.0.nics.0.firewall_rules.0.source_ip",
						ServerResource+"."+serverTestResource2, "nic.0.firewall.0.source_ip"),
					resource.TestCheckResourceAttrPair(DataSource+"."+ServersDataSource+"."+ServerDataSourceByName, "servers.0.nics.0.firewall_rules.0.source_ip",
						ServerResource+"."+serverTestResource2, "nic.0.firewall.0.source_ip"),
					resource.TestCheckResourceAttrPair(DataSource+"."+ServersDataSource+"."+ServerDataSourceByName, "servers.0.nics.0.firewall_rules.0.target_ip",
						ServerResource+"."+serverTestResource2, "nic.0.firewall.0.target_ip"),
				),
			},
			{
				Config: testAccCheck2ServersByCpuFamily,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(DataSource+"."+ServersDataSource+"."+ServerDataSourceByName, "servers.#", "2"),
					// Check server labels.
					resource.TestCheckResourceAttr(DataSource+"."+ServersDataSource+"."+ServerDataSourceByName, "servers.0.labels.#", "2"),
					resource.TestCheckResourceAttr(DataSource+"."+ServersDataSource+"."+ServerDataSourceByName, "servers.0.labels.0.key", "labelkey0"),
					resource.TestCheckResourceAttr(DataSource+"."+ServersDataSource+"."+ServerDataSourceByName, "servers.0.labels.0.value", "labelvalue0"),
					resource.TestCheckResourceAttr(DataSource+"."+ServersDataSource+"."+ServerDataSourceByName, "servers.0.labels.1.key", "labelkey1"),
					resource.TestCheckResourceAttr(DataSource+"."+ServersDataSource+"."+ServerDataSourceByName, "servers.0.labels.1.value", "labelvalue1"),
					resource.TestCheckResourceAttr(DataSource+"."+ServersDataSource+"."+ServerDataSourceByName, "servers.1.labels.0.key", "labelkey0"),
					resource.TestCheckResourceAttr(DataSource+"."+ServersDataSource+"."+ServerDataSourceByName, "servers.1.labels.0.value", "labelvalue0"),
					resource.TestCheckResourceAttr(DataSource+"."+ServersDataSource+"."+ServerDataSourceByName, "servers.1.labels.1.key", "labelkey1"),
					resource.TestCheckResourceAttr(DataSource+"."+ServersDataSource+"."+ServerDataSourceByName, "servers.1.labels.1.value", "labelvalue1"),
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
	client := testAccProvider.Meta().(SdkBundle).CloudApiClient

	ctx, cancel := context.WithTimeout(context.Background(), *resourceDefaultTimeouts.Default)

	if cancel != nil {
		defer cancel()
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != ServerResource {
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
resource ` + DatacenterResource + ` ` + DatacenterTestResource + ` {
	name       = "server-test"
	location = "us/las"
}

resource ` + LanResource + ` ` + LanTestResource + ` {
  datacenter_id = ` + DatacenterResource + `.` + DatacenterTestResource + `.id
  public = true
  name = "public"
}
resource ` + ServerResource + ` ` + ServerTestResource + ` {
  name = "` + ServerTestResource + `"
  datacenter_id = ` + DatacenterResource + `.` + DatacenterTestResource + `.id
  cores = 2
  ram = 2048
  availability_zone = "ZONE_1"
  cpu_family = "AMD_OPTERON"
  image_name ="ubuntu:latest"
  image_password = ` + RandomPassword + `.server_image_password.result
  type = "ENTERPRISE"
  volume {
    name = "` + VolumeTestResource + `"
    size = 6
    disk_type = "SSD Standard"
    user_data = "foo"
    bus = "IDE"
    availability_zone = "ZONE_1"
	}
  nic {
    lan = ` + LanResource + `.` + LanTestResource + `.id
    name = "` + LanTestResource + `"
    dhcp = false
    firewall_active = false
  }
}

resource "ionoscloud_ipblock" "webserver_ipblock" {
  location = ` + DatacenterResource + `.` + DatacenterTestResource + `.location
  size = 2
  name = "webserver_ipblock"
}

resource ` + ServerResource + ` ` + serverTestResource2 + ` {
  name = "` + serverTestResource2 + `"
  datacenter_id = ` + DatacenterResource + `.` + DatacenterTestResource + `.id
  cores = ` + noCoresTest + `
  ram = 2048
  availability_zone = "ZONE_1"
  cpu_family = "AMD_OPTERON"
  image_name ="ubuntu:latest"
  image_password = ` + RandomPassword + `.server2_image_password.result
  type = "ENTERPRISE"
  volume {
    name = "` + VolumeTestResource + "2" + `"
    size = 6
    disk_type = "SSD Standard"
    user_data = "foo"
    bus = "IDE"
    availability_zone = "ZONE_1"
	}
  nic {
    lan = 1
    name = "` + LanTestResource + "2" + `"
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
resource ` + RandomPassword + ` "server2_image_password" {
  length           = 16
  special          = false
}

data ` + ServersDataSource + ` ` + ServerDataSourceByName + ` {
 depends_on = [` + ServerResource + `.` + serverTestResource2 + `,
	` + ServerResource + `.` + ServerTestResource + `]
  datacenter_id = ` + DatacenterResource + `.` + DatacenterTestResource + `.id
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
resource ` + DatacenterResource + ` ` + DatacenterTestResource + ` {
	name       = "server-test"
	location = "us/las"
}

resource ` + LanResource + ` ` + LanTestResource + ` {
  datacenter_id = ` + DatacenterResource + `.` + DatacenterTestResource + `.id
  public = true
  name = "public"
}

resource ` + ServerResource + ` ` + ServerTestResource + ` {
  name = "` + ServerTestResource + `"
  datacenter_id = ` + DatacenterResource + `.` + DatacenterTestResource + `.id
  cores = 2
  ram = 2048
  availability_zone = "ZONE_1"
  cpu_family = "` + cpuFamilyTest + `" 
  image_name ="ubuntu:latest"
  image_password = ` + RandomPassword + `.server_image_password.result
  type = "ENTERPRISE"
  volume {
    name = "` + VolumeTestResource + `"
    size = 6
    disk_type = "SSD Standard"
    user_data = "foo"
    bus = "IDE"
    availability_zone = "ZONE_1"
	}
  nic {
    lan = ` + LanResource + `.` + LanTestResource + `.id
    name = "` + LanTestResource + `"
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
  location = ` + DatacenterResource + `.` + DatacenterTestResource + `.location
  size = 2
  name = "webserver_ipblock"
}

resource ` + ServerResource + ` ` + serverTestResource2 + ` {
  name = "` + serverTestResource2 + `"
  datacenter_id = ` + DatacenterResource + `.` + DatacenterTestResource + `.id
  cores = ` + noCoresTest + `
  ram = 2048
  availability_zone = "ZONE_1"
  cpu_family = "` + cpuFamilyTest + `" 
  image_name ="ubuntu:latest"
  image_password = ` + RandomPassword + `.server2_image_password.result
  type = "ENTERPRISE"
  volume {
    name = "` + VolumeTestResource + "2" + `"
    size = 6
    disk_type = "SSD Standard"
    user_data = "foo"
    bus = "IDE"
    availability_zone = "ZONE_1"
	}
  nic {
    lan = 1
    name = "` + LanTestResource + "2" + `"
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
resource ` + RandomPassword + ` "server2_image_password" {
  length           = 16
  special          = false
}

data ` + ServersDataSource + ` ` + ServerDataSourceByName + ` {
 depends_on = [` + ServerResource + `.` + ServerTestResource + "2" + `,
	` + ServerResource + `.` + ServerTestResource + `]
  datacenter_id = ` + DatacenterResource + `.` + DatacenterTestResource + `.id
  filter {
   name = "cpu_family"
   value = "` + cpuFamilyTest + `" 
  }
} `

const testAccCheck2ServersBadFilter = `
resource ` + DatacenterResource + ` ` + DatacenterTestResource + ` {
	name       = "server-test"
	location = "us/las"
}

resource ` + LanResource + ` ` + LanTestResource + ` {
  datacenter_id = ` + DatacenterResource + `.` + DatacenterTestResource + `.id
  public = true
  name = "public"
}
resource ` + ServerResource + ` ` + ServerTestResource + ` {
  name = "` + ServerTestResource + `"
  datacenter_id = ` + DatacenterResource + `.` + DatacenterTestResource + `.id
  cores = 2
  ram = 2048
  availability_zone = "ZONE_1"
  cpu_family = "` + cpuFamilyTest + `"
  image_name ="ubuntu:latest"
  image_password = ` + RandomPassword + `.server_image_password.result
  type = "ENTERPRISE"
  volume {
    name = "` + VolumeTestResource + `"
    size = 6
    disk_type = "SSD Standard"
    user_data = "foo"
    bus = "IDE"
    availability_zone = "ZONE_1"
	}
  nic {
    lan = ` + LanResource + `.` + LanTestResource + `.id
    name = "` + LanTestResource + `"
    dhcp = false
    firewall_active = false
  }
}

resource ` + ServerResource + ` ` + ServerTestResource + "2" + ` {
  name = "` + ServerTestResource + "2" + `"
  datacenter_id = ` + DatacenterResource + `.` + DatacenterTestResource + `.id
  cores = ` + noCoresTest + `
  ram = 2048
  availability_zone = "ZONE_1"
  cpu_family = "` + cpuFamilyTest + `"
  image_name ="ubuntu:latest"
  image_password = ` + RandomPassword + `.server2_image_password.result
  type = "ENTERPRISE"
  volume {
    name = "` + VolumeTestResource + "2" + `"
    size = 6
    disk_type = "SSD Standard"
    user_data = "foo"
    bus = "IDE"
    availability_zone = "ZONE_1"
	}
  nic {
    lan = 1
    name = "` + LanTestResource + "2" + `"
    dhcp = false
    firewall_active = false
  }
}

` + ServerImagePassword + `
resource ` + RandomPassword + ` "server2_image_password" {
  length           = 16
  special          = false
}

data ` + ServersDataSource + ` ` + ServerDataSourceByName + ` {
 depends_on = [` + ServerResource + `.` + ServerTestResource + "2" + `,
	` + ServerResource + `.` + ServerTestResource + `]
  datacenter_id = ` + DatacenterResource + `.` + DatacenterTestResource + `.id
  filter {
   name = "cpu_family"
   value = "doesNotExist" 
  }
} `

const cpuFamilyTest = "AMD_OPTERON"
const noCoresTest = "1"
