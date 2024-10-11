//go:build cdn || all || distribution

package ionoscloud

import (
	"context"
	"fmt"
	"regexp"
	"testing"

	ionoscloud_cdn "github.com/ionos-cloud/sdk-go-bundle/products/cdn/v2"

	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils/constant"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccDistributionBasic(t *testing.T) {
	var distribution ionoscloud_cdn.Distribution

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckCDNDistributionDestroyCheck,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckCDNDistributionConfigOnlyRequired,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCDNDistributionExists(constant.CDNDistributionResource+"."+constant.CDNDistributionTestResource, &distribution),
					resource.TestCheckResourceAttr(constant.CDNDistributionResource+"."+constant.CDNDistributionTestResource, "domain", "unique.test.example.com"),
					resource.TestCheckResourceAttr(constant.CDNDistributionResource+"."+constant.CDNDistributionTestResource, "routing_rules.0.scheme", "http"),
					resource.TestCheckResourceAttr(constant.CDNDistributionResource+"."+constant.CDNDistributionTestResource, "routing_rules.0.prefix", "/api"),
					resource.TestCheckResourceAttr(constant.CDNDistributionResource+"."+constant.CDNDistributionTestResource, "routing_rules.0.upstream.0.host", "server.example.com"),
					resource.TestCheckResourceAttr(constant.CDNDistributionResource+"."+constant.CDNDistributionTestResource, "routing_rules.0.upstream.0.caching", "true"),
					resource.TestCheckResourceAttr(constant.CDNDistributionResource+"."+constant.CDNDistributionTestResource, "routing_rules.0.upstream.0.waf", "true"),
					resource.TestCheckResourceAttr(constant.CDNDistributionResource+"."+constant.CDNDistributionTestResource, "routing_rules.0.upstream.0.sni_mode", "distribution"),
					resource.TestCheckResourceAttr(constant.CDNDistributionResource+"."+constant.CDNDistributionTestResource, "routing_rules.0.upstream.0.rate_limit_class", "R100"),
					resource.TestCheckResourceAttr(constant.CDNDistributionResource+"."+constant.CDNDistributionTestResource, "certificate_id", ""),
					resource.TestCheckNoResourceAttr(constant.CDNDistributionResource+"."+constant.CDNDistributionTestResource, "routing_rules.0.upstream.0.geo_restrictions.#"),
				),
			},
			{
				Config: testAccCheckCDNDistributionConfigBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCDNDistributionExists(constant.CDNDistributionResource+"."+constant.CDNDistributionTestResource, &distribution),
					resource.TestCheckResourceAttr(constant.CDNDistributionResource+"."+constant.CDNDistributionTestResource, "domain", "unique.test.example.com"),
					resource.TestCheckResourceAttr(constant.CDNDistributionResource+"."+constant.CDNDistributionTestResource, "routing_rules.0.scheme", "http"),
					resource.TestCheckResourceAttr(constant.CDNDistributionResource+"."+constant.CDNDistributionTestResource, "routing_rules.0.prefix", "/api"),
					resource.TestCheckResourceAttr(constant.CDNDistributionResource+"."+constant.CDNDistributionTestResource, "routing_rules.0.upstream.0.host", "server.example.com"),
					resource.TestCheckResourceAttr(constant.CDNDistributionResource+"."+constant.CDNDistributionTestResource, "routing_rules.0.upstream.0.caching", "true"),
					resource.TestCheckResourceAttr(constant.CDNDistributionResource+"."+constant.CDNDistributionTestResource, "routing_rules.0.upstream.0.sni_mode", "distribution"),
					resource.TestCheckResourceAttr(constant.CDNDistributionResource+"."+constant.CDNDistributionTestResource, "routing_rules.0.upstream.0.waf", "true"),
					resource.TestCheckResourceAttr(constant.CDNDistributionResource+"."+constant.CDNDistributionTestResource, "routing_rules.0.upstream.0.rate_limit_class", "R100"),
					resource.TestCheckResourceAttr(constant.CDNDistributionResource+"."+constant.CDNDistributionTestResource, "routing_rules.0.upstream.0.geo_restrictions.0.allow_list.0", "RO"),
				),
			},
			{
				Config: testAccDataSourceCDNDistributionMatchId,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(constant.CDNDistributionResource+"."+constant.CDNDistributionTestResource, "domain", "unique.test.example.com"),

					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.CDNDistributionResource+"."+constant.CDNDistributionDataSourceByID, "public_endpoint_v4", constant.CDNDistributionResource+"."+constant.CDNDistributionTestResource, "public_endpoint_v4"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.CDNDistributionResource+"."+constant.CDNDistributionDataSourceByID, "public_endpoint_v6", constant.CDNDistributionResource+"."+constant.CDNDistributionTestResource, "public_endpoint_v6"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.CDNDistributionResource+"."+constant.CDNDistributionDataSourceByID, "resource_urn", constant.CDNDistributionResource+"."+constant.CDNDistributionTestResource, "resource_urn"),

					resource.TestCheckResourceAttr(constant.CDNDistributionResource+"."+constant.CDNDistributionTestResource, "routing_rules.0.scheme", "http"),
					resource.TestCheckResourceAttr(constant.CDNDistributionResource+"."+constant.CDNDistributionTestResource, "routing_rules.0.prefix", "/api"),
					resource.TestCheckResourceAttr(constant.CDNDistributionResource+"."+constant.CDNDistributionTestResource, "routing_rules.0.upstream.0.host", "server.example.com"),
					resource.TestCheckResourceAttr(constant.CDNDistributionResource+"."+constant.CDNDistributionTestResource, "routing_rules.0.upstream.0.caching", "true"),
					resource.TestCheckResourceAttr(constant.CDNDistributionResource+"."+constant.CDNDistributionTestResource, "routing_rules.0.upstream.0.waf", "true"),
					resource.TestCheckResourceAttr(constant.CDNDistributionResource+"."+constant.CDNDistributionTestResource, "routing_rules.0.upstream.0.sni_mode", "distribution"),
					resource.TestCheckResourceAttr(constant.CDNDistributionResource+"."+constant.CDNDistributionTestResource, "routing_rules.0.upstream.0.rate_limit_class", "R100"),
					resource.TestCheckResourceAttr(constant.CDNDistributionResource+"."+constant.CDNDistributionTestResource, "routing_rules.0.upstream.0.geo_restrictions.0.allow_list.0", "RO"),
				),
			},
			{
				Config: testAccDataSourceCDNDistributionMatchDomain,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.CDNDistributionResource+"."+constant.CDNDistributionDataSourceByDomain, "public_endpoint_v4", constant.CDNDistributionResource+"."+constant.CDNDistributionTestResource, "public_endpoint_v4"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.CDNDistributionResource+"."+constant.CDNDistributionDataSourceByDomain, "public_endpoint_v6", constant.CDNDistributionResource+"."+constant.CDNDistributionTestResource, "public_endpoint_v6"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.CDNDistributionResource+"."+constant.CDNDistributionDataSourceByDomain, "resource_urn", constant.CDNDistributionResource+"."+constant.CDNDistributionTestResource, "resource_urn"),

					resource.TestCheckResourceAttr(constant.CDNDistributionResource+"."+constant.CDNDistributionTestResource, "domain", "unique.test.example.com"),
					resource.TestCheckResourceAttr(constant.CDNDistributionResource+"."+constant.CDNDistributionTestResource, "routing_rules.0.scheme", "http"),
					resource.TestCheckResourceAttr(constant.CDNDistributionResource+"."+constant.CDNDistributionTestResource, "routing_rules.0.prefix", "/api"),
					resource.TestCheckResourceAttr(constant.CDNDistributionResource+"."+constant.CDNDistributionTestResource, "routing_rules.0.upstream.0.host", "server.example.com"),
					resource.TestCheckResourceAttr(constant.CDNDistributionResource+"."+constant.CDNDistributionTestResource, "routing_rules.0.upstream.0.caching", "true"),
					resource.TestCheckResourceAttr(constant.CDNDistributionResource+"."+constant.CDNDistributionTestResource, "routing_rules.0.upstream.0.waf", "true"),
					resource.TestCheckResourceAttr(constant.CDNDistributionResource+"."+constant.CDNDistributionTestResource, "routing_rules.0.upstream.0.rate_limit_class", "R100"),
					resource.TestCheckResourceAttr(constant.CDNDistributionResource+"."+constant.CDNDistributionTestResource, "routing_rules.0.upstream.0.sni_mode", "distribution"),
					resource.TestCheckResourceAttr(constant.CDNDistributionResource+"."+constant.CDNDistributionTestResource, "routing_rules.0.upstream.0.geo_restrictions.0.allow_list.0", "RO"),
				),
			},
			{
				Config:      testAccDataSourceCDNDistributionWrongDomainError,
				ExpectError: regexp.MustCompile("no distribution found with the specified criteria"),
			},
			{
				Config: testAccCheckCDNDistributionConfigUpdate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCDNDistributionExists(constant.CDNDistributionResource+"."+constant.CDNDistributionTestResource, &distribution),
					resource.TestCheckResourceAttr(constant.CDNDistributionResource+"."+constant.CDNDistributionTestResource, "domain", "unique.test.example.com"),
					resource.TestCheckResourceAttr(constant.CDNDistributionResource+"."+constant.CDNDistributionTestResource, "routing_rules.0.scheme", "http/https"),
					resource.TestCheckResourceAttr(constant.CDNDistributionResource+"."+constant.CDNDistributionTestResource, "routing_rules.0.prefix", "/api2"),
					resource.TestCheckResourceAttr(constant.CDNDistributionResource+"."+constant.CDNDistributionTestResource, "routing_rules.0.upstream.0.host", "server.server.example.com"),
					resource.TestCheckResourceAttr(constant.CDNDistributionResource+"."+constant.CDNDistributionTestResource, "routing_rules.0.upstream.0.caching", "false"),
					resource.TestCheckResourceAttr(constant.CDNDistributionResource+"."+constant.CDNDistributionTestResource, "routing_rules.0.upstream.0.sni_mode", "origin"),
					resource.TestCheckResourceAttr(constant.CDNDistributionResource+"."+constant.CDNDistributionTestResource, "routing_rules.0.upstream.0.waf", "true"),
					resource.TestCheckResourceAttr(constant.CDNDistributionResource+"."+constant.CDNDistributionTestResource, "routing_rules.0.upstream.0.rate_limit_class", "R10"),
					resource.TestCheckResourceAttr(constant.CDNDistributionResource+"."+constant.CDNDistributionTestResource, "routing_rules.0.upstream.0.geo_restrictions.0.block_list.0", "RO"),
					resource.TestCheckResourceAttr(constant.CDNDistributionResource+"."+constant.CDNDistributionTestResource, "routing_rules.1.scheme", "https"),
					resource.TestCheckResourceAttr(constant.CDNDistributionResource+"."+constant.CDNDistributionTestResource, "routing_rules.1.prefix", "/api3"),
					resource.TestCheckResourceAttr(constant.CDNDistributionResource+"."+constant.CDNDistributionTestResource, "routing_rules.1.upstream.0.host", "server2.example.com"),
					resource.TestCheckResourceAttr(constant.CDNDistributionResource+"."+constant.CDNDistributionTestResource, "routing_rules.1.upstream.0.caching", "true"),
					resource.TestCheckResourceAttr(constant.CDNDistributionResource+"."+constant.CDNDistributionTestResource, "routing_rules.1.upstream.0.sni_mode", "origin"),
					resource.TestCheckResourceAttr(constant.CDNDistributionResource+"."+constant.CDNDistributionTestResource, "routing_rules.1.upstream.0.waf", "false"),
					resource.TestCheckResourceAttr(constant.CDNDistributionResource+"."+constant.CDNDistributionTestResource, "routing_rules.1.upstream.0.rate_limit_class", "R100"),
					resource.TestCheckResourceAttr(constant.CDNDistributionResource+"."+constant.CDNDistributionTestResource, "routing_rules.1.upstream.0.geo_restrictions.0.allow_list.0", "CN"),
					resource.TestCheckResourceAttr(constant.CDNDistributionResource+"."+constant.CDNDistributionTestResource, "routing_rules.1.upstream.0.geo_restrictions.0.allow_list.1", "RU"),
				),
			},
		},
	})
}

func testAccCheckCDNDistributionDestroyCheck(s *terraform.State) error {
	client := testAccProvider.Meta().(services.SdkBundle).CDNClient

	ctx, cancel := context.WithTimeout(context.Background(), *resourceDefaultTimeouts.Default)

	if cancel != nil {
		defer cancel()
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != constant.CDNDistributionResource {
			continue
		}

		_, apiResponse, err := client.SdkClient.DistributionsApi.DistributionsFindById(ctx, rs.Primary.ID).Execute()

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

func testAccCheckCDNDistributionExists(n string, distribution *ionoscloud_cdn.Distribution) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(services.SdkBundle).CDNClient

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

		foundDistribution, _, err := client.SdkClient.DistributionsApi.DistributionsFindById(ctx, rs.Primary.ID).Execute()

		if err != nil {
			return fmt.Errorf("error occurred while fetching distribution: %s", rs.Primary.ID)
		}
		if foundDistribution.Id != rs.Primary.ID {
			return fmt.Errorf("record not found")
		}
		distribution = &foundDistribution

		return nil
	}
}

const testAccDataSourceCDNDistributionMatchId = testAccCheckCDNDistributionConfigBasic + `
data ` + constant.CDNDistributionResource + ` ` + constant.CDNDistributionDataSourceByID + ` {
  id			= ` + constant.CDNDistributionResource + `.` + constant.CDNDistributionTestResource + `.id
}`

const testAccDataSourceCDNDistributionMatchDomain = testAccCheckCDNDistributionConfigBasic + `
data ` + constant.CDNDistributionResource + ` ` + constant.CDNDistributionDataSourceByDomain + ` {
    domain = ` + constant.CDNDistributionResource + `.` + constant.CDNDistributionTestResource + `.domain
}`

const testAccDataSourceCDNDistributionWrongDomainError = `
data ` + constant.CDNDistributionResource + ` ` + "test_distribution_matching" + ` {
    domain =  "wrong.domain.com"
}`

const testAccCheckCDNDistributionConfigUpdate = `
resource ` + constant.CertificateResource + ` ` + constant.TestCertName + ` {
	name        	  = "` + constant.TestCertName + `"
	certificate 	  = <<EOT
` + constant.TestCertificate + `
EOT
	certificate_chain = <<EOT
` + constant.TestCertificate + `
EOT
	private_key 	  = <<EOT
` + constant.PrivateKey + `
EOT
}
` + `resource ` + constant.CDNDistributionResource + ` ` + constant.CDNDistributionTestResource + ` {
	domain         = "unique.test.example.com"
	certificate_id = ` + constant.CertificateResource + `.` + constant.TestCertName + `.id` + `
	routing_rules {
		scheme = "http/https"
		prefix = "/api2"
		upstream {
			host             = "server.server.example.com"
			caching          = false
			waf              = true
			rate_limit_class = "R10"
			sni_mode 		 = "origin"
			geo_restrictions {
				block_list = [ "RO"]
			}
		}
	}
	routing_rules {
		scheme = "https"
		prefix = "/api3"
		upstream {
			host             = "server2.example.com"
			caching          = true
			waf              = false
			rate_limit_class = "R100"
			sni_mode 		 = "origin"
			geo_restrictions {
				allow_list = [ "CN", "RU"]
			}
		}
	}
}`
