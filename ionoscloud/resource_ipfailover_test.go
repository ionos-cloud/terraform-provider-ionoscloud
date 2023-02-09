//go:build compute || all || ipfailover

package ionoscloud

import (
	"context"
	"fmt"
	ionoscloud "github.com/ionos-cloud/sdk-go/v6"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccLanIPFailoverBasic(t *testing.T) {
	var lan ionoscloud.Lan
	var ipfailover ionoscloud.IPFailover

	testDeleted := func(n string) resource.TestCheckFunc {
		return func(s *terraform.State) error {
			_, ok := s.RootModule().Resources[n]
			if ok {
				return fmt.Errorf("Expected Failover group  %s to be deleted.", n)
			}
			return nil
		}
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ExternalProviders: randomProviderVersion343(),
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckLanIPFailoverDestroyCheck,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccCheckLanIPFailoverConfig),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLanIPFailoverGroupExists("ionoscloud_ipfailover.failover-test", &lan, &ipfailover),
				),
			},
			{
				Config: testAccDataSourceIpFailoverConfigBasic,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(ipfailoverResourceFullName, "id"),
					resource.TestCheckResourceAttrPair(ipfailoverResourceFullName, "id", DataSource+"."+ResourceIpFailover+"."+ipfailoverName, "id"),
					resource.TestCheckResourceAttrPair(ipfailoverResourceFullName, "nicuuid", DataSource+"."+ResourceIpFailover+"."+ipfailoverName, "nicuuid"),
					resource.TestCheckResourceAttrPair(ipfailoverResourceFullName, "lan_id", DataSource+"."+ResourceIpFailover+"."+ipfailoverName, "lan_id"),
					resource.TestCheckResourceAttrPair(ipfailoverResourceFullName, "datacenter_id", DataSource+"."+ResourceIpFailover+"."+ipfailoverName, "datacenter_id"),
				),
			},
			{
				Config: testAccCheckLanIPFailoverConfigUpdate,
				Check: resource.ComposeTestCheckFunc(
					testDeleted("ionoscloud_ipfailover.failover-test"),
				),
			},
			{
				Config: `/* */`,
			},
		},
	})
}

func testAccCheckLanIPFailoverGroupExists(n string, _ *ionoscloud.Lan, _ *ionoscloud.IPFailover) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(SdkBundle).CloudApiClient

		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no ID` is set")
		}

		dcId := rs.Primary.Attributes["datacenter_id"]
		lanId := rs.Primary.Attributes["lan_id"]
		nicUuid := rs.Primary.Attributes["nicuuid"]

		ctx, cancel := context.WithTimeout(context.Background(), *resourceDefaultTimeouts.Default)

		if cancel != nil {
			defer cancel()
		}

		lan, apiResponse, err := client.LANsApi.DatacentersLansFindById(ctx, dcId, lanId).Execute()
		logApiRequestTime(apiResponse)
		if err != nil {
			return fmt.Errorf("lan %s not found", lanId)
		}

		if lan.Properties.IpFailover == nil {
			return fmt.Errorf("lan %s has no failover groups", lanId)
		}
		found := false
		for _, fo := range *lan.Properties.IpFailover {
			if *fo.NicUuid == nicUuid {
				found = true
			}
		}
		if !found {
			return fmt.Errorf("expected NIC %s to be a part of a failover group, but not found in lans %s failover groups", nicUuid, lanId)
		}

		return nil
	}
}

func testAccCheckLanIPFailoverDestroyCheck(s *terraform.State) error {
	client := testAccProvider.Meta().(SdkBundle).CloudApiClient

	ctx, cancel := context.WithTimeout(context.Background(), *resourceDefaultTimeouts.Default)

	if cancel != nil {
		defer cancel()
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "ionoscloud_ipfailover" {
			continue
		}

		dcId := rs.Primary.Attributes["datacenter_id"]
		lanId := rs.Primary.Attributes["lan_id"]
		nicUuid := rs.Primary.Attributes["nicuuid"]

		lan, apiResponse, err := client.LANsApi.DatacentersLansFindById(ctx, dcId, lanId).Execute()
		logApiRequestTime(apiResponse)

		if err != nil {
			if !httpNotFound(apiResponse) {
				return fmt.Errorf("an error occured while fetching a Lan ID %s %s", rs.Primary.Attributes["lan_id"], err)
			}
		} else {
			found := false
			if lan.Properties.IpFailover != nil {
				for _, fo := range *lan.Properties.IpFailover {
					if *fo.NicUuid == nicUuid {
						found = true
					}
				}
				if found {
					apiResponse, err := client.DataCentersApi.DatacentersDelete(ctx, dcId).Execute()
					logApiRequestTime(apiResponse)
					if err != nil {
						return fmt.Errorf("IP failover group with nicId %s still exists %s %s, removing datacenter", nicUuid, rs.Primary.ID, err)
					}
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
  size = 1
  name = "failover test"
}

resource "ionoscloud_lan" "webserver_lan1" {
  datacenter_id = "${ionoscloud_datacenter.foobar.id}"
  public = true
  name = "terraform test"
}

resource "ionoscloud_server" "webserver" {
  name = "server"
  datacenter_id = "${ionoscloud_datacenter.foobar.id}"
  cores = 1
  ram = 1024
  availability_zone = "ZONE_1"
  cpu_family = "AMD_OPTERON"
  image_name = "ubuntu:latest"
  image_password = ` + RandomPassword + `.server_image_password.result
  volume {
    name = "system"
    size = 15
    disk_type = "SSD"
  }
  nic {
    lan = "${ionoscloud_lan.webserver_lan1.id}"
    dhcp = true
    firewall_active = true
     ips =["${ionoscloud_ipblock.webserver_ip.ips[0]}"]
  }
}
resource "ionoscloud_ipfailover" "failover-test" {
  depends_on = [ ionoscloud_lan.webserver_lan1 ]
  datacenter_id = "${ionoscloud_datacenter.foobar.id}"
  lan_id="${ionoscloud_lan.webserver_lan1.id}"
  ip ="${ionoscloud_ipblock.webserver_ip.ips[0]}"
  nicuuid= "${ionoscloud_server.webserver.primary_nic}"
}
` + ServerImagePassword

const testAccCheckLanIPFailoverConfigUpdate = `
resource "ionoscloud_datacenter" "foobar" {
	name       = "ipfailover-test"
	location = "us/las"
}

resource "ionoscloud_ipblock" "webserver_ip" {
  location = "us/las"
  size = 1
  name = "failover test"
}

resource "ionoscloud_lan" "webserver_lan1" {
  datacenter_id = "${ionoscloud_datacenter.foobar.id}"
  public = true
  name = "terraform test"
}

resource "ionoscloud_server" "webserver" {
  name = "server"
  datacenter_id = "${ionoscloud_datacenter.foobar.id}"
  cores = 1
  ram = 1024
  availability_zone = "ZONE_1"
  cpu_family = "AMD_OPTERON"
  image_name = "ubuntu:latest"
  image_password = ` + RandomPassword + `.server_image_password.result
  volume {
    name = "system"
    size = 15
    disk_type = "SSD"
  }
  nic {
    lan = "1"
    dhcp = true
    firewall_active = true
     ips =["${ionoscloud_ipblock.webserver_ip.ips[0]}"]
  }
}
` + ServerImagePassword

var testAccDataSourceIpFailoverConfigBasic = testAccCheckLanIPFailoverConfig + `
data ` + ResourceIpFailover + " " + ipfailoverName + `{
  datacenter_id = "${ionoscloud_datacenter.foobar.id}"
  id		    = ` + ipfailoverResourceFullName + `.id
}
`
