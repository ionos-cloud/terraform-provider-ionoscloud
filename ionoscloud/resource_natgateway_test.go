//go:build natgateway
// +build natgateway

package ionoscloud

import (
	"context"
	"fmt"
	"regexp"
	"testing"

	ionoscloud "github.com/ionos-cloud/sdk-go-bundle/products/cloud/v2"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils/constant"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

const resourceNatGatewayResource = constant.NatGatewayResource + "." + constant.NatGatewayTestResource
const dataSourceIdNatGatewayResource = constant.DataSource + "." + constant.NatGatewayResource + "." + constant.NatGatewayDataSourceById
const dataSourceNameNatGatewayResource = constant.DataSource + "." + constant.NatGatewayResource + "." + constant.NatGatewayDataSourceByName

func TestAccNatGatewayBasic(t *testing.T) {
	var natGateway ionoscloud.NatGateway

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactoriesInternal(t, &testAccProvider),
		CheckDestroy:             testAccCheckNatGatewayDestroyCheck,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccCheckNatGatewayConfigBasic, constant.NatGatewayTestResource),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNatGatewayExists(resourceNatGatewayResource, &natGateway),
					resource.TestCheckResourceAttr(resourceNatGatewayResource, "name", constant.NatGatewayTestResource),
					resource.TestCheckResourceAttrPair(resourceNatGatewayResource, "public_ips.0", constant.IpBlockResource+".natgateway_ips", "ips.0"),
					resource.TestCheckResourceAttrPair(resourceNatGatewayResource, "lans.0.id", constant.LanResource+".natgateway_lan", "id"),
					resource.TestCheckResourceAttr(resourceNatGatewayResource, "lans.0.gateway_ips.0", "10.11.2.5/24"),
				),
			},
			{
				Config: fmt.Sprintf(testAccDataSourceNatGatewayMatchId, constant.NatGatewayTestResource),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceIdNatGatewayResource, "name", resourceNatGatewayResource, "name"),
					resource.TestCheckResourceAttrPair(dataSourceIdNatGatewayResource, "public_ips.0", resourceNatGatewayResource, "public_ips.0"),
					resource.TestCheckResourceAttrPair(dataSourceIdNatGatewayResource, "lans.0.id", resourceNatGatewayResource, "lans.0.id"),
					resource.TestCheckResourceAttrPair(dataSourceIdNatGatewayResource, "lans.0.gateway_ips", resourceNatGatewayResource, "lans.0.gateway_ips"),
				),
			},
			{
				Config: fmt.Sprintf(testAccDataSourceNatGatewayMatchName, constant.NatGatewayTestResource),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceNameNatGatewayResource, "name", resourceNatGatewayResource, "name"),
					resource.TestCheckResourceAttrPair(dataSourceNameNatGatewayResource, "public_ips.0", resourceNatGatewayResource, "public_ips.0"),
					resource.TestCheckResourceAttrPair(dataSourceNameNatGatewayResource, "lans.0.id", resourceNatGatewayResource, "lans.0.id"),
					resource.TestCheckResourceAttrPair(dataSourceNameNatGatewayResource, "lans.0.gateway_ips", resourceNatGatewayResource, "lans.0.gateway_ips"),
				),
			},
			{
				Config:      testAccDataSourceNatGatewayWrongNameError,
				ExpectError: regexp.MustCompile(`no nat gateway found with the specified criteria`),
			},
			{
				Config: fmt.Sprintf(testAccCheckNatGatewayConfigUpdate, constant.UpdatedResources),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNatGatewayExists(resourceNatGatewayResource, &natGateway),
					resource.TestCheckResourceAttr(resourceNatGatewayResource, "name", constant.UpdatedResources),
					resource.TestCheckResourceAttr(resourceNatGatewayResource, "public_ips.#", "2"),
					resource.TestCheckResourceAttrPair(resourceNatGatewayResource, "lans.0.id", constant.LanResource+".natgateway_lan_updated", "id"),
					resource.TestCheckResourceAttr(resourceNatGatewayResource, "lans.0.gateway_ips.0", "10.11.2.6/24"),
				),
			},
		},
	})
}

func testAccCheckNatGatewayDestroyCheck(s *terraform.State) error {
	client := testAccProvider.Meta().(services.SdkBundle).CloudApiClient
	ctx, cancel := context.WithTimeout(context.Background(), *resourceDefaultTimeouts.Delete)

	if cancel != nil {
		defer cancel()
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != constant.NatGatewayResource {
			continue
		}

		_, apiResponse, err := client.NATGatewaysApi.DatacentersNatgatewaysFindByNatGatewayId(ctx, rs.Primary.Attributes["datacenter_id"], rs.Primary.ID).Execute()
		logApiRequestTime(apiResponse)

		if err != nil {
			if !httpNotFound(apiResponse) {
				return fmt.Errorf("an error occurred and checking deletion of nat gateway %s %s", rs.Primary.ID, responseBody(apiResponse))
			}
		} else {
			return fmt.Errorf("nat gateway still exists %s %s", rs.Primary.ID, err)
		}
	}

	return nil
}

func testAccCheckNatGatewayExists(n string, natGateway *ionoscloud.NatGateway) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(services.SdkBundle).CloudApiClient
		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("testAccCheckNatGatewayExists: Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no Record ID is set")
		}

		ctx, cancel := context.WithTimeout(context.Background(), *resourceDefaultTimeouts.Delete)

		if cancel != nil {
			defer cancel()
		}

		foundNatGateway, apiResponse, err := client.NATGatewaysApi.DatacentersNatgatewaysFindByNatGatewayId(ctx, rs.Primary.Attributes["datacenter_id"], rs.Primary.ID).Execute()
		logApiRequestTime(apiResponse)

		if err != nil {
			return fmt.Errorf("error occurred while fetching NatGateway: %s", rs.Primary.ID)
		}
		if *foundNatGateway.Id != rs.Primary.ID {
			return fmt.Errorf("record not found")
		}

		natGateway = &foundNatGateway

		return nil
	}
}

const testAccCheckNatGatewayConfigBasic = `
resource ` + constant.DatacenterResource + ` "natgateway_datacenter" {
  name              = "test_natgateway"
  location          = "de/fra"
  description       = "datacenter for hosting "
}

resource ` + constant.IpBlockResource + ` "natgateway_ips" {
  location = ` + constant.DatacenterResource + `.natgateway_datacenter.location
  size = 2
  name = "natgateway_ips"
}

resource ` + constant.LanResource + ` "natgateway_lan" {
  datacenter_id = ` + constant.DatacenterResource + `.natgateway_datacenter.id 
  public        = false
  name          = "test_natgateway_lan"
}

resource ` + constant.NatGatewayResource + ` ` + constant.NatGatewayTestResource + ` { 
  datacenter_id = ` + constant.DatacenterResource + `.natgateway_datacenter.id
  name          = "%s" 
  public_ips    = [ ` + constant.IpBlockResource + `.natgateway_ips.ips[0] ]
  lans {
     id          = ` + constant.LanResource + `.natgateway_lan.id
     gateway_ips = [ "10.11.2.5"] 
  }
}`

const testAccCheckNatGatewayConfigUpdate = `
resource ` + constant.DatacenterResource + ` "natgateway_datacenter" {
  name              = "test_natgateway"
  location          = "de/fra"
  description       = "datacenter for hosting "
}

resource ` + constant.IpBlockResource + ` "natgateway_ips" {
  location = ` + constant.DatacenterResource + `.natgateway_datacenter.location
  size = 2
  name = "natgateway_ips"
}

resource ` + constant.LanResource + ` "natgateway_lan" {
  datacenter_id = ` + constant.DatacenterResource + `.natgateway_datacenter.id 
  public        = false
  name          = "test_natgateway_lan"
}


resource ` + constant.LanResource + ` "natgateway_lan_updated" {
  datacenter_id = ` + constant.DatacenterResource + `.natgateway_datacenter.id 
  public        = false
  name          = "test_natgateway_lan"
}

resource ` + constant.NatGatewayResource + ` ` + constant.NatGatewayTestResource + ` { 
  datacenter_id = ` + constant.DatacenterResource + `.natgateway_datacenter.id
  name          = "%s" 
  public_ips    = [ ` + constant.IpBlockResource + `.natgateway_ips.ips[0], ` + constant.IpBlockResource + `.natgateway_ips.ips[1] ]
  lans {
     id          = ` + constant.LanResource + `.natgateway_lan_updated.id
     gateway_ips = [ "10.11.2.6/24"] 
  }
}`

const testAccDataSourceNatGatewayMatchId = testAccCheckNatGatewayConfigBasic + `
data ` + constant.NatGatewayResource + ` ` + constant.NatGatewayDataSourceById + ` {
  datacenter_id = ` + constant.DatacenterResource + `.natgateway_datacenter.id
  id			= ` + constant.NatGatewayResource + `.` + constant.NatGatewayTestResource + `.id
}
`

const testAccDataSourceNatGatewayMatchName = testAccCheckNatGatewayConfigBasic + `
data ` + constant.NatGatewayResource + ` ` + constant.NatGatewayDataSourceByName + `  {
  datacenter_id = ` + constant.DatacenterResource + `.natgateway_datacenter.id
  name			= ` + constant.NatGatewayResource + `.` + constant.NatGatewayTestResource + `.name
}
`

const testAccDataSourceNatGatewayWrongNameError = testAccCheckNatGatewayConfigBasic + `
data ` + constant.NatGatewayResource + ` ` + constant.NatGatewayDataSourceByName + `  {
  datacenter_id = ` + constant.DatacenterResource + `.natgateway_datacenter.id
  name			= "wrong_name"
}
`
