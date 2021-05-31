package ionoscloud

import (
	"context"
	"fmt"
	ionoscloud "github.com/ionos-cloud/sdk-go/v6"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccNetworkLoadBalancer_Basic(t *testing.T) {
	var networkLoadBalancer ionoscloud.NetworkLoadBalancer
	networkLoadBalancerName := "networkLoadBalancer"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNetworkLoadBalancerDestroyCheck,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccCheckNetworkLoadBalancerConfig_basic, networkLoadBalancerName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkLoadBalancerExists("ionoscloud_networkloadbalancer.test_networkloadbalancer", &networkLoadBalancer),
					resource.TestCheckResourceAttr("ionoscloud_networkloadbalancer.test_networkloadbalancer", "name", networkLoadBalancerName),
				),
			},
			{
				Config: testAccCheckNetworkLoadBalancerConfig_update,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ionoscloud_networkloadbalancer.test_networkloadbalancer", "name", "updated"),
				),
			},
		},
	})
}

func testAccCheckNetworkLoadBalancerDestroyCheck(s *terraform.State) error {
	client := testAccProvider.Meta().(SdkBundle).Client
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "ionoscloud_datacenter" {
			continue
		}

		ctx, cancel := context.WithTimeout(context.Background(), *resourceDefaultTimeouts.Delete)

		if cancel != nil {
			defer cancel()
		}

		_, apiResponse, err := client.NetworkLoadBalancersApi.DatacentersNetworkloadbalancersFindByNetworkLoadBalancerId(ctx, rs.Primary.Attributes["datacenter_id"], rs.Primary.ID).Execute()

		if _, ok := err.(ionoscloud.GenericOpenAPIError); ok {
			if apiResponse.Response.StatusCode != 404 {
				return fmt.Errorf("Network loadbalancer still exists %s %s", rs.Primary.ID, string(apiResponse.Payload))
			}
		} else {
			return fmt.Errorf("unable to fetch network loadbalancer %s %s", rs.Primary.ID, err)
		}
	}

	return nil
}

func testAccCheckNetworkLoadBalancerExists(n string, networkLoadBalancer *ionoscloud.NetworkLoadBalancer) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(SdkBundle).Client
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

		foundNetworkLoadBalancer, _, err := client.NetworkLoadBalancersApi.DatacentersNetworkloadbalancersFindByNetworkLoadBalancerId(ctx, rs.Primary.Attributes["datacenter_id"], rs.Primary.ID).Execute()

		if err != nil {
			return fmt.Errorf("error occured while fetching NetworkLoadBalancer: %s", rs.Primary.ID)
		}
		if *foundNetworkLoadBalancer.Id != rs.Primary.ID {
			return fmt.Errorf("Record not found")
		}

		networkLoadBalancer = &foundNetworkLoadBalancer

		return nil
	}
}

const testAccCheckNetworkLoadBalancerConfig_basic = `
resource "ionoscloud_datacenter" "datacenter" {
  name              = "test_nbl"
  location          = "gb/lhr"
  description       = "datacenter for hosting "
}

resource "ionoscloud_lan" "nlb_lan_1" {
  datacenter_id = ionoscloud_datacenter.datacenter.id
  public        = false
  name          = "lan_1"
}

resource "ionoscloud_lan" "nlb_lan_2" {
  datacenter_id = ionoscloud_datacenter.datacenter.id
  public        = false
  name          = "lan_2"
}


resource "ionoscloud_networkloadbalancer" "test_networkloadbalancer" {
  datacenter_id = ionoscloud_datacenter.datacenter.id
  name          = "%s"
  listener_lan  = ionoscloud_lan.nlb_lan_1.id
  target_lan    = ionoscloud_lan.nlb_lan_2.id
  ips           = ["10.12.118.224"]
  lb_private_ips = ["10.13.72.225/24"]
}
`

const testAccCheckNetworkLoadBalancerConfig_update = `
resource "ionoscloud_datacenter" "datacenter" {
  name              = "test_nbl"
  location          = "gb/lhr"
  description       = "datacenter for hosting "
}

resource "ionoscloud_lan" "nlb_lan_1" {
  datacenter_id = ionoscloud_datacenter.datacenter.id
  public        = false
  name          = "lan_1"
}

resource "ionoscloud_lan" "nlb_lan_2" {
  datacenter_id = ionoscloud_datacenter.datacenter.id
  public        = false
  name          = "lan_2"
}


resource "ionoscloud_networkloadbalancer" "test_networkloadbalancer" {
  datacenter_id = ionoscloud_datacenter.datacenter.id
  name          = "updated"
  listener_lan  = ionoscloud_lan.nlb_lan_1.id
  target_lan    = ionoscloud_lan.nlb_lan_2.id
  ips           = ["10.12.118.224"]
  lb_private_ips = ["10.13.72.225/24"]
}`
