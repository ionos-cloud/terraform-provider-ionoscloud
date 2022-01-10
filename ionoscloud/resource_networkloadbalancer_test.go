package ionoscloud

import (
	"context"
	"fmt"
	ionoscloud "github.com/ionos-cloud/sdk-go/v6"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

const networkLoadBalancerResource = NetworkLoadBalancerResource + "." + NetworkLoadBalancerTestResource

func TestAccNetworkLoadBalancer_Basic(t *testing.T) {
	var networkLoadBalancer ionoscloud.NetworkLoadBalancer

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckNetworkLoadBalancerDestroyCheck,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckNetworkLoadBalancerConfigBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkLoadBalancerExists(networkLoadBalancerResource, &networkLoadBalancer),
					resource.TestCheckResourceAttr(networkLoadBalancerResource, "name", NetworkLoadBalancerTestResource),
					resource.TestCheckResourceAttr(networkLoadBalancerResource, "ips.0", "10.12.118.224"),
					resource.TestCheckResourceAttr(networkLoadBalancerResource, "lb_private_ips.0", "10.13.72.225/24"),
					resource.TestCheckResourceAttrPair(networkLoadBalancerResource, "listener_lan", LanResource+".nlb_lan_1", "id"),
					resource.TestCheckResourceAttrPair(networkLoadBalancerResource, "target_lan", LanResource+".nlb_lan_2", "id"),
				),
			},
			{
				Config: testAccCheckNetworkLoadBalancerConfigUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(networkLoadBalancerResource, "name", UpdatedResources),
					resource.TestCheckResourceAttr(networkLoadBalancerResource, "ips.0", "10.12.118.224"),
					resource.TestCheckResourceAttr(networkLoadBalancerResource, "ips.1", "10.12.119.224"),
					resource.TestCheckResourceAttr(networkLoadBalancerResource, "lb_private_ips.0", "10.13.72.225/24"),
					resource.TestCheckResourceAttr(networkLoadBalancerResource, "lb_private_ips.1", "10.13.73.225/24"),
					resource.TestCheckResourceAttrPair(networkLoadBalancerResource, "listener_lan", LanResource+".nlb_lan_3", "id"),
					resource.TestCheckResourceAttrPair(networkLoadBalancerResource, "target_lan", LanResource+".nlb_lan_4", "id"),
				),
			},
		},
	})
}

func testAccCheckNetworkLoadBalancerDestroyCheck(s *terraform.State) error {
	client := testAccProvider.Meta().(*ionoscloud.APIClient)

	ctx, cancel := context.WithTimeout(context.Background(), *resourceDefaultTimeouts.Delete)

	if cancel != nil {
		defer cancel()
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != NetworkLoadBalancerResource {
			continue
		}

		_, apiResponse, err := client.NetworkLoadBalancersApi.DatacentersNetworkloadbalancersFindByNetworkLoadBalancerId(ctx, rs.Primary.Attributes["datacenter_id"], rs.Primary.ID).Execute()
		logApiRequestTime(apiResponse)

		if err != nil {
			if apiResponse == nil || apiResponse.Response != nil && apiResponse.StatusCode != 404 {
				return fmt.Errorf("an error occured while checking deletion of network loadbalancer %s %s", rs.Primary.ID, responseBody(apiResponse))
			}
		} else {
			return fmt.Errorf("network loadbalancer still exists %s %s", rs.Primary.ID, err)
		}
	}

	return nil
}

func testAccCheckNetworkLoadBalancerExists(n string, networkLoadBalancer *ionoscloud.NetworkLoadBalancer) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*ionoscloud.APIClient)
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

const testAccCheckNetworkLoadBalancerConfigBasic = `
resource ` + DatacenterResource + ` "datacenter" {
  name              = "test_nbl"
  location          = "gb/lhr"
  description       = "datacenter for hosting "
}

resource ` + LanResource + ` "nlb_lan_1" {
  datacenter_id = ` + DatacenterResource + `.datacenter.id
  public        = false
  name          = "lan_1"
}

resource ` + LanResource + ` "nlb_lan_2" {
  datacenter_id = ` + DatacenterResource + `.datacenter.id
  public        = false
  name          = "lan_2"
}


resource ` + NetworkLoadBalancerResource + ` ` + NetworkLoadBalancerTestResource + ` {
  datacenter_id = ` + DatacenterResource + `.datacenter.id
  name          = "` + NetworkLoadBalancerTestResource + `"
  listener_lan  = ` + LanResource + `.nlb_lan_1.id
  target_lan    = ` + LanResource + `.nlb_lan_2.id
  ips           = ["10.12.118.224"]
  lb_private_ips = ["10.13.72.225/24"]
}
`

const testAccCheckNetworkLoadBalancerConfigUpdate = `
resource ` + DatacenterResource + ` "datacenter" {
  name              = "test_nbl"
  location          = "gb/lhr"
  description       = "datacenter for hosting "
}

resource ` + LanResource + ` "nlb_lan_1" {
  datacenter_id = ` + DatacenterResource + `.datacenter.id
  public        = false
  name          = "lan_1"
}

resource ` + LanResource + ` "nlb_lan_2" {
  datacenter_id = ` + DatacenterResource + `.datacenter.id
  public        = false
  name          = "lan_2"
}

resource ` + LanResource + ` "nlb_lan_3" {
  datacenter_id = ` + DatacenterResource + `.datacenter.id
  public        = false
  name          = "lan_3"
}

resource ` + LanResource + ` "nlb_lan_4" {
  datacenter_id = ` + DatacenterResource + `.datacenter.id
  public        = false
  name          = "lan_4"
}

resource ` + NetworkLoadBalancerResource + ` ` + NetworkLoadBalancerTestResource + ` {
  datacenter_id = ` + DatacenterResource + `.datacenter.id
  name          = "` + UpdatedResources + `"
  listener_lan  = ` + LanResource + `.nlb_lan_3.id
  target_lan    = ` + LanResource + `.nlb_lan_4.id
  ips           = ["10.12.118.224", "10.12.119.224"]
  lb_private_ips = ["10.13.72.225/24", "10.13.73.225/24"]
}
`
