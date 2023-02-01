//go:build all || natgateway
// +build all natgateway

package ionoscloud

import (
	"context"
	"fmt"
	ionoscloud "github.com/ionos-cloud/sdk-go-bundle/products/compute"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

const resourceNatGatewayResource = NatGatewayResource + "." + NatGatewayTestResource
const dataSourceIdNatGatewayResource = DataSource + "." + NatGatewayResource + "." + NatGatewayDataSourceById
const dataSourceNameNatGatewayResource = DataSource + "." + NatGatewayResource + "." + NatGatewayDataSourceByName

func TestAccNatGatewayBasic(t *testing.T) {
	var natGateway ionoscloud.NatGateway

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckNatGatewayDestroyCheck,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccCheckNatGatewayConfigBasic, NatGatewayTestResource),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNatGatewayExists(resourceNatGatewayResource, &natGateway),
					resource.TestCheckResourceAttr(resourceNatGatewayResource, "name", NatGatewayTestResource),
					resource.TestCheckResourceAttrPair(resourceNatGatewayResource, "public_ips.0", IpBlockResource+".natgateway_ips", "ips.0"),
					resource.TestCheckResourceAttrPair(resourceNatGatewayResource, "lans.0.id", LanResource+".natgateway_lan", "id"),
					resource.TestCheckResourceAttr(resourceNatGatewayResource, "lans.0.gateway_ips.0", "10.11.2.5/24"),
				),
			},
			{
				Config: fmt.Sprintf(testAccDataSourceNatGatewayMatchId, NatGatewayTestResource),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceIdNatGatewayResource, "name", resourceNatGatewayResource, "name"),
					resource.TestCheckResourceAttrPair(dataSourceIdNatGatewayResource, "public_ips.0", resourceNatGatewayResource, "public_ips.0"),
					resource.TestCheckResourceAttrPair(dataSourceIdNatGatewayResource, "lans.0.id", resourceNatGatewayResource, "lans.0.id"),
					resource.TestCheckResourceAttrPair(dataSourceIdNatGatewayResource, "lans.0.gateway_ips", resourceNatGatewayResource, "lans.0.gateway_ips"),
				),
			},
			{
				Config: fmt.Sprintf(testAccDataSourceNatGatewayMatchName, NatGatewayTestResource),
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
				Config: fmt.Sprintf(testAccCheckNatGatewayConfigUpdate, UpdatedResources),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNatGatewayExists(resourceNatGatewayResource, &natGateway),
					resource.TestCheckResourceAttr(resourceNatGatewayResource, "name", UpdatedResources),
					resource.TestCheckResourceAttr(resourceNatGatewayResource, "public_ips.#", "2"),
					resource.TestCheckResourceAttrPair(resourceNatGatewayResource, "lans.0.id", LanResource+".natgateway_lan_updated", "id"),
					resource.TestCheckResourceAttr(resourceNatGatewayResource, "lans.0.gateway_ips.0", "10.11.2.6/24"),
				),
			},
		},
	})
}

func testAccCheckNatGatewayDestroyCheck(s *terraform.State) error {
	client := testAccProvider.Meta().(SdkBundle).CloudApiClient
	ctx, cancel := context.WithTimeout(context.Background(), *resourceDefaultTimeouts.Delete)

	if cancel != nil {
		defer cancel()
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != NatGatewayResource {
			continue
		}

		_, apiResponse, err := client.NATGatewaysApi.DatacentersNatgatewaysFindByNatGatewayId(ctx, rs.Primary.Attributes["datacenter_id"], rs.Primary.ID).Execute()
		logApiRequestTime(apiResponse)

		if err != nil {
			if !httpNotFound(apiResponse) {
				return fmt.Errorf("an error occured and checking deletion of nat gateway %s %s", rs.Primary.ID, responseBody(apiResponse))
			}
		} else {
			return fmt.Errorf("nat gateway still exists %s %s", rs.Primary.ID, err)
		}
	}

	return nil
}

func testAccCheckNatGatewayExists(n string, natGateway *ionoscloud.NatGateway) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(SdkBundle).CloudApiClient
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
			return fmt.Errorf("error occured while fetching NatGateway: %s", rs.Primary.ID)
		}
		if *foundNatGateway.Id != rs.Primary.ID {
			return fmt.Errorf("record not found")
		}

		natGateway = &foundNatGateway

		return nil
	}
}

const testAccCheckNatGatewayConfigBasic = `
resource ` + DatacenterResource + ` "natgateway_datacenter" {
  name              = "test_natgateway"
  location          = "de/fra"
  description       = "datacenter for hosting "
}

resource ` + IpBlockResource + ` "natgateway_ips" {
  location = ` + DatacenterResource + `.natgateway_datacenter.location
  size = 2
  name = "natgateway_ips"
}

resource ` + LanResource + ` "natgateway_lan" {
  datacenter_id = ` + DatacenterResource + `.natgateway_datacenter.id 
  public        = false
  name          = "test_natgateway_lan"
}

resource ` + NatGatewayResource + ` ` + NatGatewayTestResource + ` { 
  datacenter_id = ` + DatacenterResource + `.natgateway_datacenter.id
  name          = "%s" 
  public_ips    = [ ` + IpBlockResource + `.natgateway_ips.ips[0] ]
  lans {
     id          = ` + LanResource + `.natgateway_lan.id
     gateway_ips = [ "10.11.2.5"] 
  }
}`

const testAccCheckNatGatewayConfigUpdate = `
resource ` + DatacenterResource + ` "natgateway_datacenter" {
  name              = "test_natgateway"
  location          = "de/fra"
  description       = "datacenter for hosting "
}

resource ` + IpBlockResource + ` "natgateway_ips" {
  location = ` + DatacenterResource + `.natgateway_datacenter.location
  size = 2
  name = "natgateway_ips"
}

resource ` + LanResource + ` "natgateway_lan" {
  datacenter_id = ` + DatacenterResource + `.natgateway_datacenter.id 
  public        = false
  name          = "test_natgateway_lan"
}


resource ` + LanResource + ` "natgateway_lan_updated" {
  datacenter_id = ` + DatacenterResource + `.natgateway_datacenter.id 
  public        = false
  name          = "test_natgateway_lan"
}

resource ` + NatGatewayResource + ` ` + NatGatewayTestResource + ` { 
  datacenter_id = ` + DatacenterResource + `.natgateway_datacenter.id
  name          = "%s" 
  public_ips    = [ ` + IpBlockResource + `.natgateway_ips.ips[0], ` + IpBlockResource + `.natgateway_ips.ips[1] ]
  lans {
     id          = ` + LanResource + `.natgateway_lan_updated.id
     gateway_ips = [ "10.11.2.6/24"] 
  }
}`

const testAccDataSourceNatGatewayMatchId = testAccCheckNatGatewayConfigBasic + `
data ` + NatGatewayResource + ` ` + NatGatewayDataSourceById + ` {
  datacenter_id = ` + DatacenterResource + `.natgateway_datacenter.id
  id			= ` + NatGatewayResource + `.` + NatGatewayTestResource + `.id
}
`

const testAccDataSourceNatGatewayMatchName = testAccCheckNatGatewayConfigBasic + `
data ` + NatGatewayResource + ` ` + NatGatewayDataSourceByName + `  {
  datacenter_id = ` + DatacenterResource + `.natgateway_datacenter.id
  name			= ` + NatGatewayResource + `.` + NatGatewayTestResource + `.name
}
`

const testAccDataSourceNatGatewayWrongNameError = testAccCheckNatGatewayConfigBasic + `
data ` + NatGatewayResource + ` ` + NatGatewayDataSourceByName + `  {
  datacenter_id = ` + DatacenterResource + `.natgateway_datacenter.id
  name			= "wrong_name"
}
`
