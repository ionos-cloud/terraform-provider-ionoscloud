//go:build compute || all || ipfailover

package ionoscloud

import (
	"context"
	"fmt"
	"testing"

	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/bundleclient"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils/constant"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccLanIPFailoverBasic(t *testing.T) {
	testDeleted := func(n string) resource.TestCheckFunc {
		return func(s *terraform.State) error {
			_, ok := s.RootModule().Resources[n]
			if ok {
				return fmt.Errorf("Expected IP failover group %s to be deleted.", n)
			}
			return nil
		}
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ExternalProviders:        randomProviderVersion343(),
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactoriesInternal(t, &testAccProvider),
		CheckDestroy:             testAccCheckLanIPFailoverDestroyCheck,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckLanIPFailoverConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLanIPFailoverGroupExists(constant.ResourceIpFailover+"."+constant.IpfailoverName),
					testAccCheckLanIPFailoverGroupExists(constant.ResourceIpFailover+"."+constant.SecondIpfailoverName),
					// We can only check that the IP and NIC UUID are set in the state, we can't
					// use values to compare them since both of them are computed.
					// Checks for the first IP failover group
					resource.TestCheckResourceAttrSet(constant.ResourceIpFailover+"."+constant.IpfailoverName, "ip"),
					resource.TestCheckResourceAttrSet(constant.ResourceIpFailover+"."+constant.IpfailoverName, "nicuuid"),
					// Checks for the second IP failover group
					resource.TestCheckResourceAttrSet(constant.ResourceIpFailover+"."+constant.SecondIpfailoverName, "ip"),
					resource.TestCheckResourceAttrSet(constant.ResourceIpFailover+"."+constant.SecondIpfailoverName, "nicuuid"),
				),
			},
			{
				Config: testAccDataSourceIpFailoverConfigBasic,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(constant.IpfailoverResourceFullName, "id"),
					resource.TestCheckResourceAttrPair(constant.IpfailoverResourceFullName, "id", constant.DataSource+"."+constant.ResourceIpFailover+"."+constant.IpfailoverName, "id"),
					resource.TestCheckResourceAttrPair(constant.IpfailoverResourceFullName, "nicuuid", constant.DataSource+"."+constant.ResourceIpFailover+"."+constant.IpfailoverName, "nicuuid"),
					resource.TestCheckResourceAttrPair(constant.IpfailoverResourceFullName, "lan_id", constant.DataSource+"."+constant.ResourceIpFailover+"."+constant.IpfailoverName, "lan_id"),
					resource.TestCheckResourceAttrPair(constant.IpfailoverResourceFullName, "datacenter_id", constant.DataSource+"."+constant.ResourceIpFailover+"."+constant.IpfailoverName, "datacenter_id"),
				),
			},
			{
				Config: testAccCheckLanIPFailoverGroupUpdateIp,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(constant.ResourceIpFailover+"."+constant.SecondIpfailoverName, "ip"),
					resource.TestCheckResourceAttrSet(constant.ResourceIpFailover+"."+constant.SecondIpfailoverName, "nicuuid")),
			},
			{
				Config: testAccCheckLanIPFailoverGroupDeletion,
				Check: resource.ComposeTestCheckFunc(
					testDeleted("ionoscloud_ipfailover." + constant.IpfailoverName),
				),
			},
		},
	})
}

func testAccCheckLanIPFailoverGroupExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(bundleclient.SdkBundle).CloudApiClient

		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no ID` is set")
		}

		dcId := rs.Primary.Attributes["datacenter_id"]
		lanId := rs.Primary.Attributes["lan_id"]
		ip := rs.Primary.Attributes["ip"]

		ctx, cancel := context.WithTimeout(context.Background(), *resourceDefaultTimeouts.Default)
		defer cancel()

		lan, apiResponse, err := client.LANsApi.DatacentersLansFindById(ctx, dcId, lanId).Execute()
		logApiRequestTime(apiResponse)
		if err != nil {
			return fmt.Errorf("LAN with ID: %s not found, datacenter ID: %s", lanId, dcId)
		}
		if lan.Properties.IpFailover == nil {
			return fmt.Errorf("LAN with ID: %s has no IP failover groups", lanId)
		}
		for _, failoverGroup := range *lan.Properties.IpFailover {
			if *failoverGroup.Ip == ip {
				return nil
			}
		}
		return fmt.Errorf("IP failover group with IP: %s was not found in LAN: %s, datacenter ID: %s", ip, lanId, dcId)
	}
}

func testAccCheckLanIPFailoverDestroyCheck(s *terraform.State) error {
	client := testAccProvider.Meta().(bundleclient.SdkBundle).CloudApiClient

	ctx, cancel := context.WithTimeout(context.Background(), *resourceDefaultTimeouts.Default)
	defer cancel()

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "ionoscloud_ipfailover" {
			continue
		}

		dcId := rs.Primary.Attributes["datacenter_id"]
		lanId := rs.Primary.Attributes["lan_id"]
		nicUuid := rs.Primary.Attributes["nicuuid"]
		ip := rs.Primary.Attributes["ip"]

		lan, apiResponse, err := client.LANsApi.DatacentersLansFindById(ctx, dcId, lanId).Execute()
		logApiRequestTime(apiResponse)

		if err != nil {
			if !httpNotFound(apiResponse) {
				return fmt.Errorf("an error occurred while fetching a Lan ID %s %s", rs.Primary.Attributes["lan_id"], err)
			}
		} else {
			found := false
			if lan.Properties.IpFailover != nil {
				for _, failoverGroup := range *lan.Properties.IpFailover {
					if *failoverGroup.Ip == ip {
						found = true
						break
					}
				}
				if found {
					return fmt.Errorf("IP failover group with IP: %s, NIC UUID: %s, LAN: %s, datacenter ID: %s still exists", ip, nicUuid, lanId, dcId)
				}
			}
		}
	}

	return nil
}

const testAccCheckLanIPFailoverConfig = `
resource "ionoscloud_datacenter" "foobar" {
	name       = "ipfailover-test"
	location = "us/las"
}

resource "ionoscloud_ipblock" "webserver_ip" {
  location = "us/las"
  size = 3
  name = "failover test"
}

resource "ionoscloud_lan" "webserver_lan1" {
  datacenter_id = ionoscloud_datacenter.foobar.id
  public = true
  name = "terraform test"
}

resource "ionoscloud_server" "webserver" {
  name = "server"
  datacenter_id = ionoscloud_datacenter.foobar.id
  cores = 1
  ram = 1024
  availability_zone = "ZONE_1"
  cpu_family = "INTEL_XEON"
  image_name = "ubuntu:latest"
  image_password = ` + constant.RandomPassword + `.server_image_password.result
  volume {
    name = "system"
    size = 15
    disk_type = "SSD"
  }
  nic {
    lan = ionoscloud_lan.webserver_lan1.id
    dhcp = true
    firewall_active = true
     ips =[ionoscloud_ipblock.webserver_ip.ips[0]]
  }
}

resource "ionoscloud_server" "secondwebserver" {
  depends_on = [ionoscloud_server.webserver]
  name = "secondserver"
  datacenter_id = ionoscloud_datacenter.foobar.id
  cores = 1
  ram = 1024
  availability_zone = "ZONE_1"
  cpu_family = "INTEL_XEON"
  image_name = "ubuntu:latest"
  image_password = ` + constant.RandomPassword + `.server_image_password.result
  volume {
    name = "system"
    size = 15
    disk_type = "SSD"
  }
  nic {
    lan = ionoscloud_lan.webserver_lan1.id
    dhcp = true
    firewall_active = true
     ips =[ionoscloud_ipblock.webserver_ip.ips[1]]
  }
}
 
resource "` + constant.ResourceIpFailover + `" "` + constant.IpfailoverName + `" {
  datacenter_id = "${ionoscloud_datacenter.foobar.id}"
  lan_id=ionoscloud_lan.webserver_lan1.id
  ip =ionoscloud_ipblock.webserver_ip.ips[0]
  nicuuid= ionoscloud_server.webserver.primary_nic
}

resource "` + constant.ResourceIpFailover + `" "` + constant.SecondIpfailoverName + `" {
  depends_on = [ ` + constant.ResourceIpFailover + `.` + constant.IpfailoverName + ` ]
  datacenter_id = ionoscloud_datacenter.foobar.id
  lan_id = ionoscloud_lan.webserver_lan1.id
  ip = ionoscloud_ipblock.webserver_ip.ips[1]
  nicuuid = ionoscloud_server.secondwebserver.primary_nic
}

` + ServerImagePassword

const testAccCheckLanIPFailoverGroupUpdateNic = `
resource "ionoscloud_datacenter" "foobar" {
	name       = "ipfailover-test"
	location = "us/las"
}

resource "ionoscloud_ipblock" "webserver_ip" {
  location = "us/las"
  size = 3
  name = "failover test"
}

resource "ionoscloud_lan" "webserver_lan1" {
  datacenter_id = ionoscloud_datacenter.foobar.id
  public = true
  name = "terraform test"
}

resource "ionoscloud_server" "webserver" {
  name = "server"
  datacenter_id = ionoscloud_datacenter.foobar.id
  cores = 1
  ram = 1024
  availability_zone = "ZONE_1"
  cpu_family = "INTEL_XEON"
  image_name = "ubuntu:latest"
  image_password = ` + constant.RandomPassword + `.server_image_password.result
  volume {
    name = "system"
    size = 15
    disk_type = "SSD"
  }
  nic {
    lan = ionoscloud_lan.webserver_lan1.id
    dhcp = true
    firewall_active = true
     ips = [ionoscloud_ipblock.webserver_ip.ips[0], ionoscloud_ipblock.webserver_ip.ips[1]]
  }
}

resource "ionoscloud_server" "secondwebserver" {
  depends_on = [ionoscloud_server.webserver]
  name = "secondserver"
  datacenter_id = ionoscloud_datacenter.foobar.id
  cores = 1
  ram = 1024
  availability_zone = "ZONE_1"
  cpu_family = "INTEL_XEON"
  image_name = "ubuntu:latest"
  image_password = ` + constant.RandomPassword + `.server_image_password.result
  volume {
    name = "system"
    size = 15
    disk_type = "SSD"
  }
  nic {
    lan = ionoscloud_lan.webserver_lan1.id
    dhcp = true
    firewall_active = true
     ips =[ionoscloud_ipblock.webserver_ip.ips[2]]
  }
}

resource "` + constant.ResourceIpFailover + `" "` + constant.IpfailoverName + `" {
  datacenter_id = ionoscloud_datacenter.foobar.id
  lan_id = ionoscloud_lan.webserver_lan1.id
  ip = ionoscloud_ipblock.webserver_ip.ips[0]
  nicuuid = ionoscloud_server.webserver.primary_nic
}

resource "` + constant.ResourceIpFailover + `" "` + constant.SecondIpfailoverName + `" {
  depends_on = [ ` + constant.ResourceIpFailover + `.` + constant.IpfailoverName + ` ]
  datacenter_id = ionoscloud_datacenter.foobar.id
  lan_id = ionoscloud_lan.webserver_lan1.id
  ip = ionoscloud_ipblock.webserver_ip.ips[1]
  nicuuid = ionoscloud_server.webserver.primary_nic
}

` + ServerImagePassword

const testAccCheckLanIPFailoverGroupUpdateIp = `
resource "ionoscloud_datacenter" "foobar" {
	name       = "ipfailover-test"
	location = "us/las"
}

resource "ionoscloud_ipblock" "webserver_ip" {
  location = "us/las"
  size = 3
  name = "failover test"
}

resource "ionoscloud_lan" "webserver_lan1" {
  datacenter_id = ionoscloud_datacenter.foobar.id
  public = true
  name = "terraform test"
}

resource "ionoscloud_server" "webserver" {
  name = "server"
  datacenter_id = ionoscloud_datacenter.foobar.id
  cores = 1
  ram = 1024
  availability_zone = "ZONE_1"
  cpu_family = "INTEL_XEON"
  image_name = "ubuntu:latest"
  image_password = ` + constant.RandomPassword + `.server_image_password.result
  volume {
    name = "system"
    size = 15
    disk_type = "SSD"
  }
  nic {
    lan = ionoscloud_lan.webserver_lan1.id
    dhcp = true
    firewall_active = true
     ips = [ionoscloud_ipblock.webserver_ip.ips[0]]
  }
}

resource "ionoscloud_server" "secondwebserver" {
  depends_on = [ionoscloud_server.webserver]
  name = "secondserver"
  datacenter_id = ionoscloud_datacenter.foobar.id
  cores = 1
  ram = 1024
  availability_zone = "ZONE_1"
  cpu_family = "INTEL_XEON"
  image_name = "ubuntu:latest"
  image_password = ` + constant.RandomPassword + `.server_image_password.result
  volume {
    name = "system"
    size = 15
    disk_type = "SSD"
  }
  nic {
    lan = ionoscloud_lan.webserver_lan1.id
    dhcp = true
    firewall_active = true
     ips =[ionoscloud_ipblock.webserver_ip.ips[2]]
  }
}

resource "` + constant.ResourceIpFailover + `" "` + constant.IpfailoverName + `" {
  datacenter_id = ionoscloud_datacenter.foobar.id
  lan_id = ionoscloud_lan.webserver_lan1.id
  ip = ionoscloud_ipblock.webserver_ip.ips[0]
  nicuuid = ionoscloud_server.webserver.primary_nic
}

resource "` + constant.ResourceIpFailover + `" "` + constant.SecondIpfailoverName + `" {
  depends_on = [ ` + constant.ResourceIpFailover + `.` + constant.IpfailoverName + ` ]
  datacenter_id = ionoscloud_datacenter.foobar.id
  lan_id = ionoscloud_lan.webserver_lan1.id
  ip = ionoscloud_ipblock.webserver_ip.ips[2]
  nicuuid = ionoscloud_server.secondwebserver.primary_nic
}

` + ServerImagePassword

const testAccCheckLanIPFailoverGroupDeletion = `
resource "ionoscloud_datacenter" "foobar" {
	name       = "ipfailover-test"
	location = "us/las"
}

resource "ionoscloud_ipblock" "webserver_ip" {
  location = "us/las"
  size = 3
  name = "failover test"
}

resource "ionoscloud_lan" "webserver_lan1" {
  datacenter_id = ionoscloud_datacenter.foobar.id
  public = true
  name = "terraform test"
}

resource "ionoscloud_server" "webserver" {
  name = "server"
  datacenter_id = ionoscloud_datacenter.foobar.id
  cores = 1
  ram = 1024
  availability_zone = "ZONE_1"
  cpu_family = "INTEL_XEON"
  image_name = "ubuntu:latest"
  image_password = ` + constant.RandomPassword + `.server_image_password.result
  volume {
    name = "system"
    size = 15
    disk_type = "SSD"
  }
  nic {
    lan = "1"
    dhcp = true
    firewall_active = true
     ips =[ionoscloud_ipblock.webserver_ip.ips[0]]
  }
}

resource "ionoscloud_server" "secondwebserver" {
  name = "secondserver"
  datacenter_id = ionoscloud_datacenter.foobar.id
  cores = 1
  ram = 1024
  availability_zone = "ZONE_1"
  cpu_family = "INTEL_XEON"
  image_name = "ubuntu:latest"
  image_password = ` + constant.RandomPassword + `.server_image_password.result
  volume {
    name = "system"
    size = 15
    disk_type = "SSD"
  }
  nic {
    lan = ionoscloud_lan.webserver_lan1.id
    dhcp = true
    firewall_active = true
     ips =[ionoscloud_ipblock.webserver_ip.ips[2]]
  }
}
` + ServerImagePassword

var testAccDataSourceIpFailoverConfigBasic = testAccCheckLanIPFailoverConfig + `
data ` + constant.ResourceIpFailover + " " + constant.IpfailoverName + `{
  datacenter_id = ionoscloud_datacenter.foobar.id
  lan_id = ` + constant.ResourceIpFailover + `.` + constant.IpfailoverName + `.lan_id
  ip = ` + constant.ResourceIpFailover + `.` + constant.IpfailoverName + `.ip
}
`
