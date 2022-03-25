//go:build all || alb

package ionoscloud

import (
	"context"
	"fmt"
	ionoscloud "github.com/ionos-cloud/sdk-go/v6"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

var resourceNameAlb = ALBResource + "." + ALBTestResource
var dataSourceNameAlbById = DataSource + "." + ALBResource + "." + ALBDataSourceById
var dataSourceNameAlbByName = DataSource + "." + ALBResource + "." + ALBDataSourceByName

func TestAccApplicationLoadBalancerBasic(t *testing.T) {
	var applicationLoadBalancer ionoscloud.ApplicationLoadBalancer

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckApplicationLoadBalancerDestroyCheck,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckApplicationLoadBalancerConfigBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckApplicationLoadBalancerExists(resourceNameAlb, &applicationLoadBalancer),
					resource.TestCheckResourceAttr(resourceNameAlb, "name", ALBTestResource),
					resource.TestCheckResourceAttrPair(resourceNameAlb, "listener_lan", LanResource+".alb_lan_1", "id"),
					resource.TestCheckResourceAttrPair(resourceNameAlb, "target_lan", LanResource+".alb_lan_2", "id"),
					utils.TestValueInSlice(ALBResource, "ips.#", "10.12.118.224"),
					utils.TestValueInSlice(ALBResource, "lb_private_ips.#", "10.13.72.225/24"),
				),
			},
			{
				Config: testAccDataSourceApplicationLoadBalancerMatchId,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(resourceNameAlb, "name", dataSourceNameAlbById, "name"),
					resource.TestCheckResourceAttrPair(resourceNameAlb, "listener_lan", dataSourceNameAlbById, "listener_lan"),
					resource.TestCheckResourceAttrPair(resourceNameAlb, "target_lan", dataSourceNameAlbById, "target_lan"),
					resource.TestCheckResourceAttrPair(resourceNameAlb, "ips.0", dataSourceNameAlbById, "ips.0"),
					resource.TestCheckResourceAttrPair(resourceNameAlb, "lb_private_ips.0", dataSourceNameAlbById, "lb_private_ips.0"),
				),
			},
			{
				Config: testAccDataSourceApplicationLoadBalancerMatchName,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(resourceNameAlb, "name", dataSourceNameAlbByName, "name"),
					resource.TestCheckResourceAttrPair(resourceNameAlb, "listener_lan", dataSourceNameAlbByName, "listener_lan"),
					resource.TestCheckResourceAttrPair(resourceNameAlb, "target_lan", dataSourceNameAlbByName, "target_lan"),
					resource.TestCheckResourceAttrPair(resourceNameAlb, "ips.0", dataSourceNameAlbByName, "ips.0"),
					resource.TestCheckResourceAttrPair(resourceNameAlb, "lb_private_ips.0", dataSourceNameAlbByName, "lb_private_ips.0"),
				),
			},
			{
				Config:      testAccDataSourceApplicationLoadBalancerWrongNameError,
				ExpectError: regexp.MustCompile("no application load balanacer found with the specified criteria"),
			},
			{
				Config: testAccCheckApplicationLoadBalancerConfigUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceNameAlb, "name", UpdatedResources),
					resource.TestCheckResourceAttrPair(resourceNameAlb, "listener_lan", LanResource+".alb_lan_3", "id"),
					resource.TestCheckResourceAttrPair(resourceNameAlb, "target_lan", LanResource+".alb_lan_4", "id"),
					utils.TestValueInSlice(ALBResource, "ips.#", "10.12.118.224"),
					utils.TestValueInSlice(ALBResource, "ips.#", "10.12.119.224"),
					utils.TestValueInSlice(ALBResource, "lb_private_ips.#", "10.13.72.225/24"),
					utils.TestValueInSlice(ALBResource, "lb_private_ips.#", "10.13.73.225/24"),
				),
			},
		},
	})
}

func testAccCheckApplicationLoadBalancerDestroyCheck(s *terraform.State) error {
	client := testAccProvider.Meta().(SdkBundle).CloudApiClient
	ctx, cancel := context.WithTimeout(context.Background(), *resourceDefaultTimeouts.Delete)

	if cancel != nil {
		defer cancel()
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != ALBResource {
			continue
		}

		dcId := rs.Primary.Attributes["datacenter_id"]
		albId := rs.Primary.ID

		_, apiResponse, err := client.ApplicationLoadBalancersApi.DatacentersApplicationloadbalancersFindByApplicationLoadBalancerId(ctx, dcId, albId).Execute()
		logApiRequestTime(apiResponse)

		if err != nil {
			if errorBesideNotFound(apiResponse) {
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
		client := testAccProvider.Meta().(SdkBundle).CloudApiClient
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

		foundNatGateway, apiResponse, err := client.ApplicationLoadBalancersApi.DatacentersApplicationloadbalancersFindByApplicationLoadBalancerId(ctx, dcId, albId).Execute()
		logApiRequestTime(apiResponse)

		if err != nil {
			return fmt.Errorf("error occured while fetching NatGateway: %s, %w", rs.Primary.ID, err)
		}
		if *foundNatGateway.Id != rs.Primary.ID {
			return fmt.Errorf("record not found")
		}

		alb = &foundNatGateway

		return nil
	}
}

const testAccCheckApplicationLoadBalancerConfigBasic = `
resource ` + DatacenterResource + ` "alb_datacenter" {
  name              = "test_alb"
  location          = "de/fkb"
  description       = "datacenter for hosting "
}

resource ` + LanResource + ` "alb_lan_1" {
  datacenter_id = ` + DatacenterResource + `.alb_datacenter.id 
  public        = false
  name          = "test_alb_lan_1"
}

resource ` + LanResource + ` "alb_lan_2" {
  datacenter_id = ` + DatacenterResource + `.alb_datacenter.id 
  public        = false
  name          = "test_alb_lan_2"
}

resource ` + ALBResource + ` ` + ALBTestResource + ` { 
  datacenter_id = ` + DatacenterResource + `.alb_datacenter.id
  name          = "` + ALBTestResource + `"
  listener_lan  = ` + LanResource + `.alb_lan_1.id
  ips           = [ "10.12.118.224"]
  target_lan    = ` + LanResource + `.alb_lan_2.id
  lb_private_ips= [ "10.13.72.225/24"]
}`

const testAccCheckApplicationLoadBalancerConfigUpdate = `
resource ` + DatacenterResource + ` "alb_datacenter" {
  name              = "test_alb"
  location          = "de/fkb"
  description       = "datacenter for hosting "
}

resource ` + LanResource + ` "alb_lan_1" {
  datacenter_id = ` + DatacenterResource + `.alb_datacenter.id 
  public        = false
  name          = "test_alb_lan_1"
}

resource ` + LanResource + ` "alb_lan_2" {
  datacenter_id = ` + DatacenterResource + `.alb_datacenter.id 
  public        = false
  name          = "test_alb_lan_2"
}

resource ` + LanResource + ` "alb_lan_3" {
  datacenter_id = ` + DatacenterResource + `.alb_datacenter.id 
  public        = false
  name          = "test_alb_lan_3"
}

resource ` + LanResource + ` "alb_lan_4" {
  datacenter_id = ` + DatacenterResource + `.alb_datacenter.id 
  public        = false
  name          = "test_alb_lan_4"
}

resource ` + ALBResource + ` ` + ALBTestResource + ` { 
  datacenter_id = ` + DatacenterResource + `.alb_datacenter.id
  name          = "` + UpdatedResources + `"
  listener_lan    = ` + LanResource + `.alb_lan_3.id
  ips           = [ "10.12.118.224", "10.12.119.224"]
  target_lan    = ` + LanResource + `.alb_lan_4.id
  lb_private_ips= [ "10.13.72.225/24", "10.13.73.225/24"]
}`

const testAccDataSourceApplicationLoadBalancerMatchId = testAccCheckApplicationLoadBalancerConfigBasic + `
data ` + ALBResource + ` ` + ALBDataSourceById + ` {
  datacenter_id = ` + DatacenterResource + `.alb_datacenter.id
  id			= ` + ALBResource + `.` + ALBTestResource + `.id
}
`

const testAccDataSourceApplicationLoadBalancerMatchName = testAccCheckApplicationLoadBalancerConfigBasic + `
data ` + ALBResource + ` ` + ALBDataSourceByName + ` {
  datacenter_id = ` + DatacenterResource + `.alb_datacenter.id
  name          = ` + ALBResource + `.` + ALBTestResource + `.name
}
`

const testAccDataSourceApplicationLoadBalancerWrongNameError = testAccCheckApplicationLoadBalancerConfigBasic + `
data ` + ALBResource + ` ` + ALBDataSourceByName + ` {
  datacenter_id = ` + DatacenterResource + `.alb_datacenter.id
  name          = "wrong_name"
}
`
