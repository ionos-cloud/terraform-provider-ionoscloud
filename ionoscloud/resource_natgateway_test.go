// +build natgateway

package ionoscloud

import (
	"context"
	"fmt"
	ionoscloud "github.com/ionos-cloud/sdk-go/v6"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccNatGateway_Basic(t *testing.T) {
	var natGateway ionoscloud.NatGateway
	natGatewayName := "natGateway"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckNatGatewayDestroyCheck,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccCheckNatGatewayConfigBasic, natGatewayName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNatGatewayExists("ionoscloud_natgateway.natgateway", &natGateway),
					resource.TestCheckResourceAttr("ionoscloud_natgateway.natgateway", "name", natGatewayName),
					resource.TestCheckResourceAttrPair("ionoscloud_natgateway.natgateway", "public_ips.0", "ionoscloud_ipblock.natgateway_ips", "ips.0"),
				),
			},
			{
				Config: fmt.Sprintf(testAccCheckNatGatewayConfigUpdate),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ionoscloud_natgateway.natgateway", "name", "updated"),
					resource.TestCheckResourceAttrPair("ionoscloud_natgateway.natgateway", "public_ips.0", "ionoscloud_ipblock.natgateway_ips", "ips.0"),
				),
			},
		},
	})
}

func testAccCheckNatGatewayDestroyCheck(s *terraform.State) error {
	client := testAccProvider.Meta().(*ionoscloud.APIClient)
	ctx, cancel := context.WithTimeout(context.Background(), *resourceDefaultTimeouts.Delete)

	if cancel != nil {
		defer cancel()
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "ionoscloud_datacenter" {
			continue
		}

		_, apiResponse, err := client.NATGatewaysApi.DatacentersNatgatewaysFindByNatGatewayId(ctx, rs.Primary.Attributes["datacenter_id"], rs.Primary.ID).Execute()

		if err != nil {
			if apiResponse == nil || apiResponse.Response != nil && apiResponse.StatusCode != 404 {
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
		client := testAccProvider.Meta().(*ionoscloud.APIClient)
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

		foundNatGateway, _, err := client.NATGatewaysApi.DatacentersNatgatewaysFindByNatGatewayId(ctx, rs.Primary.Attributes["datacenter_id"], rs.Primary.ID).Execute()

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
resource "ionoscloud_datacenter" "natgateway_datacenter" {
  name              = "test_natgateway"
  location          = "de/fra"
  description       = "datacenter for hosting "
}

resource "ionoscloud_ipblock" "natgateway_ips" {
  location = ionoscloud_datacenter.natgateway_datacenter.location
  size = 1
  name = "natgateway_ips"
}

resource "ionoscloud_lan" "natgateway_lan" {
  datacenter_id = ionoscloud_datacenter.natgateway_datacenter.id 
  public        = false
  name          = "test_natgateway_lan"
}

resource "ionoscloud_natgateway" "natgateway" { 
  datacenter_id = ionoscloud_datacenter.natgateway_datacenter.id
  name          = "%s"
  public_ips    = [ ionoscloud_ipblock.natgateway_ips.ips[0] ]
  lans {
     id          = ionoscloud_lan.natgateway_lan.id
  }
}`

const testAccCheckNatGatewayConfigUpdate = `
resource "ionoscloud_datacenter" "natgateway_datacenter" {
  name              = "test_natgateway"
  location          = "de/fra"
  description       = "datacenter for hosting "
}

resource "ionoscloud_ipblock" "natgateway_ips" {
  location = ionoscloud_datacenter.natgateway_datacenter.location
  size = 1
  name = "natgateway_ips"
}

resource "ionoscloud_lan" "natgateway_lan" {
  datacenter_id = ionoscloud_datacenter.natgateway_datacenter.id 
  public        = false
  name          = "test_natgateway_lan"
}

resource "ionoscloud_natgateway" "natgateway" { 
  datacenter_id = ionoscloud_datacenter.natgateway_datacenter.id
  name          = "updated" 
  public_ips    = [ ionoscloud_ipblock.natgateway_ips.ips[0] ]
  lans {
     id          = ionoscloud_lan.natgateway_lan.id
  }
}`
