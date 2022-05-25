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

func TestServersBasic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckServersDestroyCheck,
		Steps: []resource.TestStep{
			{
				Config: testAccCheck2ServersByNameAndCores,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(DataSource+"."+ServersDataSource+"."+ServerDataSourceByName, "servers.#", "1"),
				),
			},
			{
				Config: testAccCheck2ServersByCpuFamily,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(DataSource+"."+ServersDataSource+"."+ServerDataSourceByName, "servers.#", "2"),
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
			if apiResponse == nil || apiResponse.Response != nil && apiResponse.StatusCode != 404 {
				return fmt.Errorf("unable to fetch server %s: %s", rs.Primary.ID, err)
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
  image_password = "K3tTj8G14a3EgKyNeeiYsasad"
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
  cpu_family = "AMD_OPTERON"
  image_name ="ubuntu:latest"
  image_password = "K3tTj8G14a3EgKyNeeiYsasad1"
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

data ` + ServersDataSource + ` ` + ServerDataSourceByName + ` {
 depends_on = [` + ServerResource + `.` + ServerTestResource + "2" + `,
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
  image_password = "K3tTj8G14a3EgKyNeeiYsasad"
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
  image_password = "K3tTj8G14a3EgKyNeeiYsasad1"
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
  image_password = "K3tTj8G14a3EgKyNeeiYsasad"
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
  image_password = "K3tTj8G14a3EgKyNeeiYsasad1"
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
