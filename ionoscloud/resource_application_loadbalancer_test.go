//go:build alb

package ionoscloud

import (
	"context"
	"fmt"
	"regexp"
	"testing"

	ionoscloud "github.com/ionos-cloud/sdk-go/v6"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils/constant"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

var resourceNameAlb = constant.ALBResource + "." + constant.ALBTestResource
var dataSourceNameAlbById = constant.DataSource + "." + constant.ALBResource + "." + constant.ALBDataSourceById
var dataSourceNameAlbByName = constant.DataSource + "." + constant.ALBResource + "." + constant.ALBDataSourceByName

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
					resource.TestCheckResourceAttr(resourceNameAlb, "name", constant.ALBTestResource),
					resource.TestCheckResourceAttr(resourceNameAlb, "central_logging", true),
					resource.TestCheckResourceAttr(resourceNameAlb, "logging_format", "%{+Q}o %{-Q}ci - - [%trg] %r %ST %B "" "" %cp %ms %ft %b %s %TR %Tw %Tc %Tr %Ta %tsc %ac %fc %bc %sc %rc %sq %bq %CC %CS %hrl %hsl"),
					resource.TestCheckResourceAttrPair(resourceNameAlb, "listener_lan", constant.LanResource+".alb_lan_1", "id"),
					resource.TestCheckResourceAttrPair(resourceNameAlb, "target_lan", constant.LanResource+".alb_lan_2", "id"),
					utils.TestValueInSlice(constant.ALBResource, "ips.#", "10.12.118.224"),
					utils.TestValueInSlice(constant.ALBResource, "lb_private_ips.#", "10.13.72.225/24"),
					resource.TestCheckResourceAttr(resourceNameAlb, "flowlog.0.name", "test_flowlog"),
					resource.TestCheckResourceAttr(resourceNameAlb, "flowlog.0.action", "ALL"),
					resource.TestCheckResourceAttr(resourceNameAlb, "flowlog.0.direction", "INGRESS"),
					resource.TestCheckResourceAttr(resourceNameAlb, "flowlog.0.bucket", constant.FlowlogBucket),
				),
			},
			{
				Config: testAccDataSourceApplicationLoadBalancerMatchId,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(resourceNameAlb, "name", dataSourceNameAlbById, "name"),
					resource.TestCheckResourceAttrPair(resourceNameAlb, "central_logging", true, "central_logging"),
					resource.TestCheckResourceAttrPair(resourceNameAlb, "logging_format", "%{+Q}o %{-Q}ci - - [%trg] %r %ST %B "" "" %cp %ms %ft %b %s %TR %Tw %Tc %Tr %Ta %tsc %ac %fc %bc %sc %rc %sq %bq %CC %CS %hrl %hsl", "logging_format"),
					resource.TestCheckResourceAttrPair(resourceNameAlb, "listener_lan", dataSourceNameAlbById, "listener_lan"),
					resource.TestCheckResourceAttrPair(resourceNameAlb, "target_lan", dataSourceNameAlbById, "target_lan"),
					resource.TestCheckResourceAttrPair(resourceNameAlb, "ips.0", dataSourceNameAlbById, "ips.0"),
					resource.TestCheckResourceAttrPair(resourceNameAlb, "lb_private_ips.0", dataSourceNameAlbById, "lb_private_ips.0"),
					resource.TestCheckResourceAttrPair(resourceNameAlb, "flowlog.0.name", dataSourceNameAlbById, "flowlog.0.name"),
					resource.TestCheckResourceAttrPair(resourceNameAlb, "flowlog.0.action", dataSourceNameAlbById, "flowlog.0.action"),
					resource.TestCheckResourceAttrPair(resourceNameAlb, "flowlog.0.direction", dataSourceNameAlbById, "flowlog.0.direction"),
					resource.TestCheckResourceAttrPair(resourceNameAlb, "flowlog.0.direction", dataSourceNameAlbById, "flowlog.0.direction"),
				),
			},
			{
				Config: testAccDataSourceApplicationLoadBalancerPartialMatchName,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(resourceNameAlb, "name", dataSourceNameAlbByName, "name"),
					resource.TestCheckResourceAttrPair(resourceNameAlb, "listener_lan", dataSourceNameAlbByName, "listener_lan"),
					resource.TestCheckResourceAttrPair(resourceNameAlb, "target_lan", dataSourceNameAlbByName, "target_lan"),
					resource.TestCheckResourceAttrPair(resourceNameAlb, "ips.0", dataSourceNameAlbByName, "ips.0"),
					resource.TestCheckResourceAttrPair(resourceNameAlb, "lb_private_ips.0", dataSourceNameAlbByName, "lb_private_ips.0"),
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
				Config:      testAccDataSourceApplicationLoadBalancerWrongPartialNameError,
				ExpectError: regexp.MustCompile("no application load balanacer found with the specified criteria"),
			},
			{
				Config: testAccCheckApplicationLoadBalancerConfigUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceNameAlb, "name", constant.UpdatedResources),
					resource.TestCheckResourceAttrPair(resourceNameAlb, "listener_lan", constant.LanResource+".alb_lan_3", "id"),
					resource.TestCheckResourceAttrPair(resourceNameAlb, "target_lan", constant.LanResource+".alb_lan_4", "id"),
					resource.TestCheckResourceAttr(resourceNameAlb, "central_logging", false),
					resource.TestCheckResourceAttr(resourceNameAlb, "logging_format", "%{+Q}o %{-Q}ci - - [%trg]"),
					utils.TestValueInSlice(constant.ALBResource, "ips.#", "10.12.118.224"),
					utils.TestValueInSlice(constant.ALBResource, "ips.#", "10.12.119.224"),
					utils.TestValueInSlice(constant.ALBResource, "lb_private_ips.#", "10.13.72.225/24"),
					utils.TestValueInSlice(constant.ALBResource, "lb_private_ips.#", "10.13.73.225/24"),
					resource.TestCheckResourceAttr(resourceNameAlb, "flowlog.0.name", "test_flowlog_updated"),
					resource.TestCheckResourceAttr(resourceNameAlb, "flowlog.0.action", "REJECTED"),
					resource.TestCheckResourceAttr(resourceNameAlb, "flowlog.0.direction", "EGRESS"),
					resource.TestCheckResourceAttr(resourceNameAlb, "flowlog.0.bucket", constant.FlowlogBucketUpdated),
				),
			},
		},
	})
}

func testAccCheckApplicationLoadBalancerDestroyCheck(s *terraform.State) error {
	client := testAccProvider.Meta().(services.SdkBundle).CloudApiClient
	ctx, cancel := context.WithTimeout(context.Background(), *resourceDefaultTimeouts.Delete)

	if cancel != nil {
		defer cancel()
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != constant.ALBResource {
			continue
		}

		dcId := rs.Primary.Attributes["datacenter_id"]
		albId := rs.Primary.ID

		_, apiResponse, err := client.ApplicationLoadBalancersApi.DatacentersApplicationloadbalancersFindByApplicationLoadBalancerId(ctx, dcId, albId).Execute()
		logApiRequestTime(apiResponse)

		if err != nil {
			if !httpNotFound(apiResponse) {
				return fmt.Errorf("an error occurred and checking deletion of application loadbalancer %s %s", rs.Primary.ID, responseBody(apiResponse))
			}
		} else {
			return fmt.Errorf("application loadbalancer still exists %s %s", rs.Primary.ID, err)
		}
	}

	return nil
}

func testAccCheckApplicationLoadBalancerExists(n string, alb *ionoscloud.ApplicationLoadBalancer) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(services.SdkBundle).CloudApiClient
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
			return fmt.Errorf("error occurred while fetching NatGateway: %s, %w", rs.Primary.ID, err)
		}
		if *foundNatGateway.Id != rs.Primary.ID {
			return fmt.Errorf("record not found")
		}

		alb = &foundNatGateway

		return nil
	}
}

const testAccCheckApplicationLoadBalancerConfigBasic = `
resource ` + constant.DatacenterResource + ` "alb_datacenter" {
  name              = "test_alb"
  location          = "de/txl"
  description       = "datacenter for hosting "
  central_logging   = true
  logging_format	= "%{+Q}o %{-Q}ci - - [%trg] %r %ST %B "" "" %cp %ms %ft %b %s %TR %Tw %Tc %Tr %Ta %tsc %ac %fc %bc %sc %rc %sq %bq %CC %CS %hrl %hsl"
}

resource ` + constant.LanResource + ` "alb_lan_1" {
  datacenter_id = ` + constant.DatacenterResource + `.alb_datacenter.id 
  public        = false
  name          = "test_alb_lan_1"
}

resource ` + constant.LanResource + ` "alb_lan_2" {
  datacenter_id = ` + constant.DatacenterResource + `.alb_datacenter.id 
  public        = false
  name          = "test_alb_lan_2"
}

resource ` + constant.ALBResource + ` ` + constant.ALBTestResource + ` { 
  datacenter_id = ` + constant.DatacenterResource + `.alb_datacenter.id
  name          = "` + constant.ALBTestResource + `"
  listener_lan  = ` + constant.LanResource + `.alb_lan_1.id
  ips           = [ "10.12.118.224"]
  target_lan    = ` + constant.LanResource + `.alb_lan_2.id
  lb_private_ips= [ "10.13.72.225/24"]
  flowlog {
    name = "test_flowlog"
    action = "ALL"
    direction = "INGRESS"
    bucket = "` + constant.FlowlogBucket + `"
  }
}`

const testAccCheckApplicationLoadBalancerConfigUpdate = `
resource ` + constant.DatacenterResource + ` "alb_datacenter" {
  name              = "test_alb"
  location          = "de/txl"
  description       = "datacenter for hosting "
}

resource ` + constant.LanResource + ` "alb_lan_1" {
  datacenter_id = ` + constant.DatacenterResource + `.alb_datacenter.id 
  public        = false
  name          = "test_alb_lan_1"
}

resource ` + constant.LanResource + ` "alb_lan_2" {
  datacenter_id = ` + constant.DatacenterResource + `.alb_datacenter.id 
  public        = false
  name          = "test_alb_lan_2"
}

resource ` + constant.LanResource + ` "alb_lan_3" {
  datacenter_id = ` + constant.DatacenterResource + `.alb_datacenter.id 
  public        = false
  name          = "test_alb_lan_3"
}

resource ` + constant.LanResource + ` "alb_lan_4" {
  datacenter_id = ` + constant.DatacenterResource + `.alb_datacenter.id 
  public        = false
  name          = "test_alb_lan_4"
}

resource ` + constant.ALBResource + ` ` + constant.ALBTestResource + ` { 
  datacenter_id = ` + constant.DatacenterResource + `.alb_datacenter.id
  name          = "` + constant.UpdatedResources + `"
  listener_lan    = ` + constant.LanResource + `.alb_lan_3.id
  ips           = [ "10.12.118.224", "10.12.119.224"]
  target_lan    = ` + constant.LanResource + `.alb_lan_4.id
  lb_private_ips= [ "10.13.72.225/24", "10.13.73.225/24"]
  flowlog {
    name = "test_flowlog_updated"
    action = "REJECTED"
    direction = "EGRESS"
    bucket = "` + constant.FlowlogBucketUpdated + `"
  }
  central_logging   = false
  logging_format	= "%{+Q}o %{-Q}ci - - [%trg]"
}`

const testAccDataSourceApplicationLoadBalancerMatchId = testAccCheckApplicationLoadBalancerConfigBasic + `
data ` + constant.ALBResource + ` ` + constant.ALBDataSourceById + ` {
  datacenter_id = ` + constant.DatacenterResource + `.alb_datacenter.id
  id			= ` + constant.ALBResource + `.` + constant.ALBTestResource + `.id
}
`

const testAccDataSourceApplicationLoadBalancerMatchName = testAccCheckApplicationLoadBalancerConfigBasic + `
data ` + constant.ALBResource + ` ` + constant.ALBDataSourceByName + ` {
  datacenter_id = ` + constant.DatacenterResource + `.alb_datacenter.id
  name          = ` + constant.ALBResource + `.` + constant.ALBTestResource + `.name
}
`

const testAccDataSourceApplicationLoadBalancerPartialMatchName = testAccCheckApplicationLoadBalancerConfigBasic + `
data ` + constant.ALBResource + ` ` + constant.ALBDataSourceByName + ` {
  datacenter_id = ` + constant.DatacenterResource + `.alb_datacenter.id
  name          = "` + constant.DataSourcePartial + `"
  partial_match = true
}
`

const testAccDataSourceApplicationLoadBalancerWrongNameError = testAccCheckApplicationLoadBalancerConfigBasic + `
data ` + constant.ALBResource + ` ` + constant.ALBDataSourceByName + ` {
  datacenter_id = ` + constant.DatacenterResource + `.alb_datacenter.id
  name          = "wrong_name"
}
`

const testAccDataSourceApplicationLoadBalancerWrongPartialNameError = testAccCheckApplicationLoadBalancerConfigBasic + `
data ` + constant.ALBResource + ` ` + constant.ALBDataSourceByName + ` {
  datacenter_id = ` + constant.DatacenterResource + `.alb_datacenter.id
  name          = "wrong_name"
  partial_match = true
}
`
