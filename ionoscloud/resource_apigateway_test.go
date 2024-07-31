//go:build apigateway || all || gateway

package ionoscloud

import (
	"context"
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	ionoscloud "github.com/ionos-cloud/sdk-go-apigateway"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services"
)

var testAccCheckApiGatewayConfig_basic = `
resource "ionoscloud_apigateway" "example" {
  name = "example"
  logs = false
  metrics = true
  custom_domains {
    name = "example.com"
    certificate_id = "00000000-0000-0000-0000-000000000000"
  }
}
`

var testAccCheckApiGatewayConfig_update = `
resource "ionoscloud_apigateway" "example" {
  name = "example_updated"
  logs = false
  metrics = false
  custom_domains {
    name = "example-updated.com"
    certificate_id = "00000000-0000-0000-0000-000000000000"
  }
}
`

var testAccDataSourceApiGatewayMatchId = testAccCheckApiGatewayConfig_basic + `
data "ionoscloud_apigateway" "example_by_id" {
  id = ionoscloud_apigateway.example.id
}
`

var testAccDataSourceApiGatewayMatchName = testAccCheckApiGatewayConfig_basic + `
data "ionoscloud_apigateway" "example_by_name" {
  name = ionoscloud_apigateway.example.name
}
`

var testAccDataSourceApiGatewayMatching = testAccCheckApiGatewayConfig_basic + `
data "ionoscloud_apigateway" "example_matching" {
  name = ionoscloud_apigateway.example.name
}
`

var testAccDataSourceApiGatewayMultipleResultsError = testAccCheckApiGatewayConfig_basic + `
resource "ionoscloud_apigateway" "example_multiple" {
  name = "example"
  logs = false
  metrics = true
  custom_domains {
    name = "example.com"
    certificate_id = "00000000-0000-0000-0000-000000000000"
  }
}

data "ionoscloud_apigateway" "example_matching" {
  name = ionoscloud_apigateway.example.name
}
`

var testAccDataSourceApiGatewayWrongNameError = testAccCheckApiGatewayConfig_basic + `
data "ionoscloud_apigateway" "example_wrong_name" {
  name = "wrong_name"
}
`

func TestAccApiGateway_basic(t *testing.T) {
	var apiGateway ionoscloud.GatewayRead

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactoriesInternal(t, &testAccProvider),
		CheckDestroy:             testAccCheckApiGatewayDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckApiGatewayConfig_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckApiGatewayExists("ionoscloud_apigateway.example", &apiGateway),
					resource.TestCheckResourceAttr("ionoscloud_apigateway.example", "name", "example"),
					resource.TestCheckResourceAttr("ionoscloud_apigateway.example", "logs", "false"),
					resource.TestCheckResourceAttr("ionoscloud_apigateway.example", "metrics", "true"),
					resource.TestCheckResourceAttr("ionoscloud_apigateway.example", "custom_domains.0.name", "example.com"),
					resource.TestCheckResourceAttr("ionoscloud_apigateway.example", "custom_domains.0.certificate_id", "00000000-0000-0000-0000-000000000000"),
				),
			},
			{
				Config: testAccCheckApiGatewayConfig_update,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckApiGatewayExists("ionoscloud_apigateway.example", &apiGateway),
					resource.TestCheckResourceAttr("ionoscloud_apigateway.example", "name", "example_updated"),
					resource.TestCheckResourceAttr("ionoscloud_apigateway.example", "logs", "false"),
					resource.TestCheckResourceAttr("ionoscloud_apigateway.example", "metrics", "false"),
					resource.TestCheckResourceAttr("ionoscloud_apigateway.example", "custom_domains.0.name", "example-updated.com"),
					resource.TestCheckResourceAttr("ionoscloud_apigateway.example", "custom_domains.0.certificate_id", "00000000-0000-0000-0000-000000000000"),
				),
			},
			{
				Config: testAccDataSourceApiGatewayMatchId,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair("data.ionoscloud_apigateway.example_by_id", "name", "ionoscloud_apigateway.example", "name"),
					resource.TestCheckResourceAttrPair("data.ionoscloud_apigateway.example_by_id", "logs", "ionoscloud_apigateway.example", "logs"),
					resource.TestCheckResourceAttrPair("data.ionoscloud_apigateway.example_by_id", "metrics", "ionoscloud_apigateway.example", "metrics"),
					resource.TestCheckResourceAttrPair("data.ionoscloud_apigateway.example_by_id", "custom_domains.0.name", "ionoscloud_apigateway.example", "custom_domains.0.name"),
					resource.TestCheckResourceAttrPair("data.ionoscloud_apigateway.example_by_id", "custom_domains.0.certificate_id", "ionoscloud_apigateway.example", "custom_domains.0.certificate_id"),
				),
			},
			{
				Config: testAccDataSourceApiGatewayMatchName,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair("data.ionoscloud_apigateway.example_by_name", "name", "ionoscloud_apigateway.example", "name"),
					resource.TestCheckResourceAttrPair("data.ionoscloud_apigateway.example_by_name", "logs", "ionoscloud_apigateway.example", "logs"),
					resource.TestCheckResourceAttrPair("data.ionoscloud_apigateway.example_by_name", "metrics", "ionoscloud_apigateway.example", "metrics"),
					resource.TestCheckResourceAttrPair("data.ionoscloud_apigateway.example_by_name", "custom_domains.0.name", "ionoscloud_apigateway.example", "custom_domains.0.name"),
					resource.TestCheckResourceAttrPair("data.ionoscloud_apigateway.example_by_name", "custom_domains.0.certificate_id", "ionoscloud_apigateway.example", "custom_domains.0.certificate_id"),
				),
			},
			{
				Config: testAccDataSourceApiGatewayMatching,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair("data.ionoscloud_apigateway.example_matching", "name", "ionoscloud_apigateway.example", "name"),
					resource.TestCheckResourceAttrPair("data.ionoscloud_apigateway.example_matching", "logs", "ionoscloud_apigateway.example", "logs"),
					resource.TestCheckResourceAttrPair("data.ionoscloud_apigateway.example_matching", "metrics", "ionoscloud_apigateway.example", "metrics"),
					resource.TestCheckResourceAttrPair("data.ionoscloud_apigateway.example_matching", "custom_domains.0.name", "ionoscloud_apigateway.example", "custom_domains.0.name"),
					resource.TestCheckResourceAttrPair("data.ionoscloud_apigateway.example_matching", "custom_domains.0.certificate_id", "ionoscloud_apigateway.example", "custom_domains.0.certificate_id"),
				),
			},
			{
				Config:      testAccDataSourceApiGatewayMultipleResultsError,
				ExpectError: regexp.MustCompile("more than one API Gateway found with the specified"),
			},
			{
				Config:      testAccDataSourceApiGatewayWrongNameError,
				ExpectError: regexp.MustCompile("no API Gateway found with the specified"),
			},
		},
	})
}

func testAccCheckApiGatewayDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(services.SdkBundle).ApiGatewayClient

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "ionoscloud_apigateway" {
			continue
		}

		_, resp, err := client.GetApiGatewayById(context.Background(), rs.Primary.ID)
		if err != nil {
			if resp == nil || resp.StatusCode != 404 {
				return fmt.Errorf("error fetching API Gateway with ID %s: %v", rs.Primary.ID, err)
			}
		} else {
			return fmt.Errorf("API Gateway still exists: %s", rs.Primary.ID)
		}
	}

	return nil
}

func testAccCheckApiGatewayExists(n string, apiGateway *ionoscloud.GatewayRead) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		client := testAccProvider.Meta().(services.SdkBundle).ApiGatewayClient
		found, _, err := client.GetApiGatewayById(context.Background(), rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("error fetching API Gateway with ID %s: %v", rs.Primary.ID, err)
		}

		if found.Id != nil && *found.Id != rs.Primary.ID {
			return fmt.Errorf("API Gateway not found")
		}

		*apiGateway = found
		return nil
	}
}
