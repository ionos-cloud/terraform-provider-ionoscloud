package ionoscloud

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services"
)

var testAccCheckApiGatewayConfig_basic = `
resource "ionoscloud_apigateway" "example" {
  name = "example"
  logs = true
  metrics = true
  custom_domains {
    name = "example.com"
    certificate_id = "example-certificate-id"
  }
}
`

func TestAccApiGateway_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckApiGatewayConfig_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckApiGatewayExists("ionoscloud_apigateway.example"),
					resource.TestCheckResourceAttr("ionoscloud_apigateway.example", "name", "example"),
					resource.TestCheckResourceAttr("ionoscloud_apigateway.example", "logs", "true"),
					resource.TestCheckResourceAttr("ionoscloud_apigateway.example", "metrics", "true"),
					resource.TestCheckResourceAttr("ionoscloud_apigateway.example", "custom_domains.0.name", "example.com"),
					resource.TestCheckResourceAttr("ionoscloud_apigateway.example", "custom_domains.0.certificate_id", "example-certificate-id"),
				),
			},
			{
				ResourceName:      "ionoscloud_apigateway.example",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

var testAccCheckApiGatewayConfig_update = `
resource "ionoscloud_apigateway" "example" {
  name = "example_updated"
  logs = false
  metrics = false
  custom_domains {
    name = "example_updated.com"
    certificate_id = "example-certificate-id-updated"
  }
}
`

func TestAccApiGateway_update(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckApiGatewayConfig_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckApiGatewayExists("ionoscloud_apigateway.example"),
				),
			},
			{
				Config: testAccCheckApiGatewayConfig_update,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckApiGatewayExists("ionoscloud_apigateway.example"),
					resource.TestCheckResourceAttr("ionoscloud_apigateway.example", "name", "example_updated"),
					resource.TestCheckResourceAttr("ionoscloud_apigateway.example", "logs", "false"),
					resource.TestCheckResourceAttr("ionoscloud_apigateway.example", "metrics", "false"),
					resource.TestCheckResourceAttr("ionoscloud_apigateway.example", "custom_domains.0.name", "example_updated.com"),
					resource.TestCheckResourceAttr("ionoscloud_apigateway.example", "custom_domains.0.certificate_id", "example-certificate-id-updated"),
				),
			},
		},
	})
}

func TestAccApiGateway_destroy(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckApiGatewayConfig_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckApiGatewayExists("ionoscloud_apigateway.example"),
				),
			},
		},
		CheckDestroy: testAccCheckApiGatewayDestroy,
	})
}

func TestAccApiGateway_import(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckApiGatewayConfig_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckApiGatewayExists("ionoscloud_apigateway.example"),
				),
			},
			{
				ResourceName:      "ionoscloud_apigateway.example",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateId:     "example_import_id",
			},
			{
				Config: testAccCheckApiGatewayConfig_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckApiGatewayExists("ionoscloud_apigateway.example"),
				),
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateId:     "example_import_id",
			},
		},
	})
}

func testAccCheckApiGatewayExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		client := testAccProvider.Meta().(services.SdkBundle).ApiGatewayClient
		_, _, err := client.GetApiGatewayById(context.Background(), rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("error fetching API Gateway with ID: %v, error: %w", rs.Primary.ID, err)
		}

		return nil
	}
}

func testAccCheckApiGatewayDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "ionoscloud_apigateway" {
			continue
		}

		client := testAccProvider.Meta().(services.SdkBundle).ApiGatewayClient
		_, resp, _ := client.GetApiGatewayById(context.Background(), rs.Primary.ID)
		if resp != nil && resp.StatusCode != 404 {
			return fmt.Errorf("API Gateway still exists: %s", rs.Primary.ID)
		}
	}

	return nil
}
