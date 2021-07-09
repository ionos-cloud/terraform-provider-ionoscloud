package ionoscloud

import (
	"context"
	"fmt"
	ionoscloud "github.com/ionos-cloud/sdk-go/v6"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccApplicationLoadBalancer_Basic(t *testing.T) {
	var applicationLoadBalancer ionoscloud.ApplicationLoadBalancer
	applicationaLoadbalancerName := "applicationLoadBalancer"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckApplicationLoadBalancerDestroyCheck,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccCheckApplicationLoadBalancerConfigBasic, applicationaLoadbalancerName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckApplicationLoadBalancerExists("ionoscloud_application_loadbalancer.alb", &applicationLoadBalancer),
					resource.TestCheckResourceAttr("ionoscloud_application_loadbalancer.alb", "name", applicationaLoadbalancerName),
				),
			},
			{
				Config: fmt.Sprintf(testAccCheckApplicationLoadBalancerConfigUpdate),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ionoscloud_application_loadbalancer.alb", "name", "updated"),
				),
			},
		},
	})
}

func testAccCheckApplicationLoadBalancerDestroyCheck(s *terraform.State) error {
	client := testAccProvider.Meta().(*ionoscloud.APIClient)
	ctx, cancel := context.WithTimeout(context.Background(), *resourceDefaultTimeouts.Delete)

	if cancel != nil {
		defer cancel()
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "ionoscloud_datacenter" {
			continue
		}

		dcId := rs.Primary.Attributes["datacenter_id"]
		albId := rs.Primary.ID

		_, apiResponse, err := client.ApplicationLoadBalancersApi.DatacentersApplicationloadbalancersFindByApplicationLoadBalancerId(ctx, dcId, albId).Execute()

		if err != nil {
			if apiResponse == nil || apiResponse.StatusCode != 404 {
				return fmt.Errorf("an error occured and checking deletion of application loadbalancer %s %s", rs.Primary.ID, responseBody(apiResponse))
			}
		} else {
			return fmt.Errorf("application loadbalancer still exists %s %s", rs.Primary.ID, err)
		}
	}

	return nil
}

func testAccCheckApplicationLoadBalancerExists(n string, alb *ionoscloud.ApplicationLoadBalancer) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*ionoscloud.APIClient)
		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("testAccCheckApplicationLoadBalancerExists: Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no Record ID is set")
		}

		ctx, cancel := context.WithTimeout(context.Background(), *resourceDefaultTimeouts.Delete)

		if cancel != nil {
			defer cancel()
		}

		dcId := rs.Primary.Attributes["datacenter_id"]
		albId := rs.Primary.ID

		foundNatGateway, _, err := client.ApplicationLoadBalancersApi.DatacentersApplicationloadbalancersFindByApplicationLoadBalancerId(ctx, dcId, albId).Execute()

		if err != nil {
			return fmt.Errorf("error occured while fetching NatGateway: %s", rs.Primary.ID)
		}
		if *foundNatGateway.Id != rs.Primary.ID {
			return fmt.Errorf("record not found")
		}

		alb = &foundNatGateway

		return nil
	}
}

const testAccCheckApplicationLoadBalancerConfigBasic = `
resource "ionoscloud_datacenter" "alb_datacenter" {
  name              = "test_alb"
  location          = "gb/lhr"
  description       = "datacenter for hosting "
}

resource "ionoscloud_lan" "alb_lan_1" {
  datacenter_id = ionoscloud_datacenter.alb_datacenter.id 
  public        = false
  name          = "test_alb_lan_1"
}

resource "ionoscloud_lan" "alb_lan_2" {
  datacenter_id = ionoscloud_datacenter.alb_datacenter.id 
  public        = false
  name          = "test_alb_lan_2"
}

resource "ionoscloud_application_loadbalancer" "alb" { 
  datacenter_id = ionoscloud_datacenter.alb_datacenter.id
  name          = "%s"
  listener_lan  = ionoscloud_lan.alb_lan_1.id
  ips           = [ "81.173.1.2",
                    "22.231.2.2",
                    "22.231.2.3"
                  ]
  target_lan    = ionoscloud_lan.alb_lan_2.id
  lb_private_ips= [ "81.173.1.5/24",
                    "22.231.2.5/24"
                  ]
}`

const testAccCheckApplicationLoadBalancerConfigUpdate = `
resource "ionoscloud_datacenter" "alb_datacenter" {
  name              = "test_alb"
  location          = "gb/lhr"
  description       = "datacenter for hosting "
}

resource "ionoscloud_lan" "alb_lan_1" {
  datacenter_id = ionoscloud_datacenter.alb_datacenter.id 
  public        = false
  name          = "test_alb_lan_1"
}

resource "ionoscloud_lan" "alb_lan_2" {
  datacenter_id = ionoscloud_datacenter.alb_datacenter.id 
  public        = false
  name          = "test_alb_lan_2"
}

resource "ionoscloud_application_loadbalancer" "alb" { 
  datacenter_id = ionoscloud_datacenter.alb_datacenter.id
  name          = "updated"
  listener_lan  = ionoscloud_lan.alb_lan_1.id
  ips           = [ "81.173.1.2",
                    "22.231.2.2",
                    "22.231.2.3"
                  ]
  target_lan    = ionoscloud_lan.alb_lan_2.id
  lb_private_ips= [ "81.173.1.5/24",
                    "22.231.2.5/24"
                  ]
}`
