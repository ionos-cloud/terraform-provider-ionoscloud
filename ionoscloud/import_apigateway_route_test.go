//go:build apigateway || all || route

package ionoscloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils/constant"
)

func TestAccApiGatewayRoute_import(t *testing.T) {
	resource.Test(
		t, resource.TestCase{
			PreCheck:                 func() { testAccPreCheck(t) },
			ProtoV6ProviderFactories: testAccProtoV6ProviderFactoriesInternal(t, &testAccProvider),
			CheckDestroy:             testAccCheckAPIGatewayRouteDestroyCheck,
			Steps: []resource.TestStep{
				{
					Config: configAPIGatewayRouteBasic(routeResourceName, routeAttributeNameValue),
				},
				{
					ResourceName:      constant.ApiGatewayRouteResource + "." + routeResourceName,
					ImportStateIdFunc: testAccApiGatewayRouteImportStateId,
					ImportState:       true,
					ImportStateVerify: true,
				},
			},
		},
	)
}

func testAccApiGatewayRouteImportStateId(s *terraform.State) (string, error) {
	importID := ""

	for _, rs := range s.RootModule().Resources {
		if rs.Type != constant.ApiGatewayRouteResource {
			continue
		}

		gatewayId := rs.Primary.Attributes["gateway_id"]
		id := rs.Primary.Attributes["id"]
		importID = fmt.Sprintf("%s:%s", gatewayId, id)
	}

	return importID, nil
}