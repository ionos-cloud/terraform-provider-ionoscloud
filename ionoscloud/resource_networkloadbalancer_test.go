//go:build nlb
// +build nlb

package ionoscloud

import (
	"context"
	"fmt"
	"regexp"
	"testing"

	ionoscloud "github.com/ionos-cloud/sdk-go/v6"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils/constant"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

const networkLoadBalancerResource = constant.NetworkLoadBalancerResource + "." + constant.NetworkLoadBalancerTestResource

const dataSourceNetworkLoadBalancerId = constant.DataSource + "." + constant.NetworkLoadBalancerResource + "." + constant.NetworkLoadBalancerDataSourceById
const dataSourceNetworkLoadBalancerName = constant.DataSource + "." + constant.NetworkLoadBalancerResource + "." + constant.NetworkLoadBalancerDataSourceByName

func TestAccNetworkLoadBalancerBasic(t *testing.T) {
	var networkLoadBalancer ionoscloud.NetworkLoadBalancer

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckNetworkLoadBalancerDestroyCheck,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckNetworkLoadBalancerConfigBasicWithoutPrivateIp,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkLoadBalancerExists(networkLoadBalancerResource, &networkLoadBalancer),
					resource.TestCheckResourceAttr(networkLoadBalancerResource, "name", constant.NetworkLoadBalancerTestResource),
					resource.TestCheckResourceAttr(networkLoadBalancerResource, "ips.0", "10.12.118.224"),
					resource.TestCheckResourceAttrSet(networkLoadBalancerResource, "lb_private_ips.0"),
					resource.TestCheckResourceAttrPair(networkLoadBalancerResource, "listener_lan", constant.LanResource+".nlb_lan_1", "id"),
					resource.TestCheckResourceAttrPair(networkLoadBalancerResource, "target_lan", constant.LanResource+".nlb_lan_2", "id"),
				),
			},
			{
				Config: testAccCheckNetworkLoadBalancerConfigBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkLoadBalancerExists(networkLoadBalancerResource, &networkLoadBalancer),
					resource.TestCheckResourceAttr(networkLoadBalancerResource, "name", constant.NetworkLoadBalancerTestResource),
					resource.TestCheckResourceAttr(networkLoadBalancerResource, "ips.0", "10.12.118.224"),
					resource.TestCheckResourceAttr(networkLoadBalancerResource, "lb_private_ips.0", "10.13.72.225/24"),
					resource.TestCheckResourceAttrPair(networkLoadBalancerResource, "listener_lan", constant.LanResource+".nlb_lan_1", "id"),
					resource.TestCheckResourceAttrPair(networkLoadBalancerResource, "target_lan", constant.LanResource+".nlb_lan_2", "id"),
				),
			},
			{
				Config: testAccDataSourceNetworkLoadBalancerMatchId,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(networkLoadBalancerResource, "name", dataSourceNetworkLoadBalancerId, "name"),
					resource.TestCheckResourceAttrPair(networkLoadBalancerResource, "listener_lan", dataSourceNetworkLoadBalancerId, "listener_lan"),
					resource.TestCheckResourceAttrPair(networkLoadBalancerResource, "ips", dataSourceNetworkLoadBalancerId, "ips"),
					resource.TestCheckResourceAttrPair(networkLoadBalancerResource, "target_lan", dataSourceNetworkLoadBalancerId, "target_lan"),
					resource.TestCheckResourceAttrPair(networkLoadBalancerResource, "lb_private_ips", dataSourceNetworkLoadBalancerId, "lb_private_ips"),
				),
			},
			{
				Config: testAccDataSourceNetworkLoadBalancerMatchName,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(networkLoadBalancerResource, "name", dataSourceNetworkLoadBalancerName, "name"),
					resource.TestCheckResourceAttrPair(networkLoadBalancerResource, "listener_lan", dataSourceNetworkLoadBalancerName, "listener_lan"),
					resource.TestCheckResourceAttrPair(networkLoadBalancerResource, "ips", dataSourceNetworkLoadBalancerName, "ips"),
					resource.TestCheckResourceAttrPair(networkLoadBalancerResource, "target_lan", dataSourceNetworkLoadBalancerName, "target_lan"),
					resource.TestCheckResourceAttrPair(networkLoadBalancerResource, "lb_private_ips", dataSourceNetworkLoadBalancerName, "lb_private_ips"),
				),
			},
			{
				Config:      testAccDataSourceNetworkLoadBalancerWrongNameError,
				ExpectError: regexp.MustCompile(`no network load balancer found with the specified criteria: name`),
			},
			{
				Config: testAccCheckNetworkLoadBalancerConfigUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(networkLoadBalancerResource, "name", constant.UpdatedResources),
					resource.TestCheckResourceAttr(networkLoadBalancerResource, "ips.0", "10.12.118.224"),
					resource.TestCheckResourceAttr(networkLoadBalancerResource, "ips.1", "10.12.119.224"),
					resource.TestCheckResourceAttr(networkLoadBalancerResource, "lb_private_ips.0", "10.13.72.225/24"),
					resource.TestCheckResourceAttr(networkLoadBalancerResource, "lb_private_ips.1", "10.13.73.225/24"),
					resource.TestCheckResourceAttrPair(networkLoadBalancerResource, "listener_lan", constant.LanResource+".nlb_lan_3", "id"),
					resource.TestCheckResourceAttrPair(networkLoadBalancerResource, "target_lan", constant.LanResource+".nlb_lan_4", "id"),
				),
			},
		},
	})
}

func testAccCheckNetworkLoadBalancerDestroyCheck(s *terraform.State) error {
	client := testAccProvider.Meta().(services.SdkBundle).CloudApiClient

	ctx, cancel := context.WithTimeout(context.Background(), *resourceDefaultTimeouts.Delete)

	if cancel != nil {
		defer cancel()
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != constant.NetworkLoadBalancerResource {
			continue
		}

		_, apiResponse, err := client.NetworkLoadBalancersApi.DatacentersNetworkloadbalancersFindByNetworkLoadBalancerId(ctx, rs.Primary.Attributes["datacenter_id"], rs.Primary.ID).Execute()
		logApiRequestTime(apiResponse)

		if err != nil {
			if !httpNotFound(apiResponse) {
				return fmt.Errorf("an error occured while checking deletion of network loadbalancer %s %s", rs.Primary.ID, responseBody(apiResponse))
			}
		} else {
			return fmt.Errorf("network loadbalancer still exists %s %w", rs.Primary.ID, err)
		}
	}

	return nil
}

func testAccCheckNetworkLoadBalancerExists(n string, networkLoadBalancer *ionoscloud.NetworkLoadBalancer) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(services.SdkBundle).CloudApiClient
		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("testAccCheckNetworkLoadBalancerExists: Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no Record ID is set")
		}

		ctx, cancel := context.WithTimeout(context.Background(), *resourceDefaultTimeouts.Delete)

		if cancel != nil {
			defer cancel()
		}

		foundNetworkLoadBalancer, apiResponse, err := client.NetworkLoadBalancersApi.DatacentersNetworkloadbalancersFindByNetworkLoadBalancerId(ctx, rs.Primary.Attributes["datacenter_id"], rs.Primary.ID).Execute()
		logApiRequestTime(apiResponse)

		if err != nil {
			return fmt.Errorf("error occured while fetching NetworkLoadBalancer: %s", rs.Primary.ID)
		}
		if *foundNetworkLoadBalancer.Id != rs.Primary.ID {
			return fmt.Errorf("record not found")
		}

		networkLoadBalancer = &foundNetworkLoadBalancer

		return nil
	}
}

const testAccCheckNetworkLoadBalancerConfigBasicWithoutPrivateIp = `
resource ` + constant.DatacenterResource + ` "datacenter" {
  name              = "test_nbl"
  location          = "gb/lhr"
  description       = "datacenter for hosting "
}

resource ` + constant.LanResource + ` "nlb_lan_1" {
  datacenter_id = ` + constant.DatacenterResource + `.datacenter.id
  public        = false
  name          = "lan_1"
}

resource ` + constant.LanResource + ` "nlb_lan_2" {
  datacenter_id = ` + constant.DatacenterResource + `.datacenter.id
  public        = false
  name          = "lan_2"
}


resource ` + constant.NetworkLoadBalancerResource + ` ` + constant.NetworkLoadBalancerTestResource + ` {
  datacenter_id = ` + constant.DatacenterResource + `.datacenter.id
  name          = "` + constant.NetworkLoadBalancerTestResource + `"
  listener_lan  = ` + constant.LanResource + `.nlb_lan_1.id
  target_lan    = ` + constant.LanResource + `.nlb_lan_2.id
  ips           = ["10.12.118.224"]
}
`

const testAccCheckNetworkLoadBalancerConfigBasic = `
resource ` + constant.DatacenterResource + ` "datacenter" {
  name              = "test_nbl"
  location          = "gb/lhr"
  description       = "datacenter for hosting "
}

resource ` + constant.LanResource + ` "nlb_lan_1" {
  datacenter_id = ` + constant.DatacenterResource + `.datacenter.id
  public        = false
  name          = "lan_1"
}

resource ` + constant.LanResource + ` "nlb_lan_2" {
  datacenter_id = ` + constant.DatacenterResource + `.datacenter.id
  public        = false
  name          = "lan_2"
}


resource ` + constant.NetworkLoadBalancerResource + ` ` + constant.NetworkLoadBalancerTestResource + ` {
  datacenter_id = ` + constant.DatacenterResource + `.datacenter.id
  name          = "` + constant.NetworkLoadBalancerTestResource + `"
  listener_lan  = ` + constant.LanResource + `.nlb_lan_1.id
  target_lan    = ` + constant.LanResource + `.nlb_lan_2.id
  ips           = ["10.12.118.224"]
  lb_private_ips = ["10.13.72.225/24"]
}
`

const testAccCheckNetworkLoadBalancerConfigUpdate = `
resource ` + constant.DatacenterResource + ` "datacenter" {
  name              = "test_nbl"
  location          = "gb/lhr"
  description       = "datacenter for hosting "
}

resource ` + constant.LanResource + ` "nlb_lan_1" {
  datacenter_id = ` + constant.DatacenterResource + `.datacenter.id
  public        = false
  name          = "lan_1"
}

resource ` + constant.LanResource + ` "nlb_lan_2" {
  datacenter_id = ` + constant.DatacenterResource + `.datacenter.id
  public        = false
  name          = "lan_2"
}

resource ` + constant.LanResource + ` "nlb_lan_3" {
  datacenter_id = ` + constant.DatacenterResource + `.datacenter.id
  public        = false
  name          = "lan_3"
}

resource ` + constant.LanResource + ` "nlb_lan_4" {
  datacenter_id = ` + constant.DatacenterResource + `.datacenter.id
  public        = false
  name          = "lan_4"
}

resource ` + constant.NetworkLoadBalancerResource + ` ` + constant.NetworkLoadBalancerTestResource + ` {
  datacenter_id = ` + constant.DatacenterResource + `.datacenter.id
  name          = "` + constant.UpdatedResources + `"
  listener_lan  = ` + constant.LanResource + `.nlb_lan_3.id
  target_lan    = ` + constant.LanResource + `.nlb_lan_4.id
  ips           = ["10.12.118.224", "10.12.119.224"]
  lb_private_ips = ["10.13.72.225/24", "10.13.73.225/24"]
}
`

const testAccDataSourceNetworkLoadBalancerMatchId = testAccCheckNetworkLoadBalancerConfigBasic + `
data ` + constant.NetworkLoadBalancerResource + ` ` + constant.NetworkLoadBalancerDataSourceById + ` {
  datacenter_id = ` + constant.NetworkLoadBalancerResource + `.` + constant.NetworkLoadBalancerTestResource + `.datacenter_id
  id			= ` + constant.NetworkLoadBalancerResource + `.` + constant.NetworkLoadBalancerTestResource + `.id
}
`

const testAccDataSourceNetworkLoadBalancerMatchName = testAccCheckNetworkLoadBalancerConfigBasic + `
data ` + constant.NetworkLoadBalancerResource + ` ` + constant.NetworkLoadBalancerDataSourceByName + ` {
  datacenter_id = ` + constant.NetworkLoadBalancerResource + `.` + constant.NetworkLoadBalancerTestResource + `.datacenter_id
  name			= ` + constant.NetworkLoadBalancerResource + `.` + constant.NetworkLoadBalancerTestResource + `.name
}
`

const testAccDataSourceNetworkLoadBalancerWrongNameError = testAccCheckNetworkLoadBalancerConfigBasic + `
data ` + constant.NetworkLoadBalancerResource + ` ` + constant.NetworkLoadBalancerDataSourceByName + ` {
  datacenter_id = ` + constant.NetworkLoadBalancerResource + `.` + constant.NetworkLoadBalancerTestResource + `.datacenter_id
  name			= "wrong_name"
}
`
