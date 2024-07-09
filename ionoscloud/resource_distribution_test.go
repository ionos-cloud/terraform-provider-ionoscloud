//go:build cdn || all || distribution

package ionoscloud

import (
	"context"
	"fmt"
	"regexp"
	"testing"

	ionoscloud_cdn "github.com/ionos-cloud/sdk-go-cdn"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils/constant"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccDistributionBasic(t *testing.T) {
	var distribution ionoscloud_cdn.Distribution

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckDistributionDestroyCheck,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckDistributionConfigUpdate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDistributionExists(constant.DistributionResource+"."+constant.DistributionTestResource, &distribution),
					resource.TestCheckResourceAttr(constant.DistributionResource+"."+constant.DistributionTestResource, "domain", "example.com.com"),
				),
			},
			{
				Config:      testAccDataSourceDistributionWrongDomainError,
				ExpectError: regexp.MustCompile("no distribution found with the specified criteria"),
			},
		},
	})
}

func testAccCheckDistributionDestroyCheck(s *terraform.State) error {
	client := testAccProvider.Meta().(services.SdkBundle).CdnClient

	ctx, cancel := context.WithTimeout(context.Background(), *resourceDefaultTimeouts.Default)

	if cancel != nil {
		defer cancel()
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != constant.DistributionResource {
			continue
		}

		_, apiResponse, err := client.DistributionsApi.DistributionsFindById(ctx, rs.Primary.ID).Execute()

		if err != nil {
			if apiResponse.StatusCode != 404 {
				return fmt.Errorf("an error occurred while checking the destruction of distribution %s: %w", rs.Primary.ID, err)
			}
		} else {
			return fmt.Errorf("distribution %s still exists", rs.Primary.ID)
		}
	}

	return nil
}

func testAccCheckDistributionExists(n string, distribution *ionoscloud_cdn.Distribution) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(services.SdkBundle).CdnClient

		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no Record ID is set")
		}

		ctx, cancel := context.WithTimeout(context.Background(), *resourceDefaultTimeouts.Default)

		if cancel != nil {
			defer cancel()
		}

		foundDistribution, _, err := client.DistributionsApi.DistributionsFindById(ctx, rs.Primary.ID).Execute()

		if err != nil {
			return fmt.Errorf("error occurred while fetching distribution: %s", rs.Primary.ID)
		}
		if *foundDistribution.Id != rs.Primary.ID {
			return fmt.Errorf("record not found")
		}
		distribution = &foundDistribution

		return nil
	}
}

const testAccCheckDistributionConfigUpdate = `resource ` + constant.DistributionResource + ` ` + constant.DistributionTestResource + ` {
	domain         = "example.com.com"
	routing_rules {
		scheme = "https"
		prefix = "/api"
		upstream {
			host             = "server.example.com"
			caching          = true
			waf              = true
			rate_limit_class = "none"
			geo_restrictions {
				allow_list = [ "CN"]
			}
		}
	}
	routing_rules {
		scheme = "http/https"
		prefix = "/api2123"
		upstream {
			host             = "server2.example.com"
			caching          = false
			waf              = false
			rate_limit_class = "R100"
			geo_restrictions {
				block_list = [ "RO", "RU"]
			}
		}
	}
}`

const testAccDataSourceDistributionWrongDomainError = `
data ` + constant.DistributionResource + ` ` + "test_distribution_matching" + ` {
    domain =  "wrong.domain.com"
}`
