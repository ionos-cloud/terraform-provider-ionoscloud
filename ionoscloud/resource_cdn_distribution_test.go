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
		CheckDestroy:      testAccCheckCdnDistributionDestroyCheck,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckCdnDistributionConfigBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCdnDistributionExists(constant.CdnDistributionResource+"."+constant.CdnDistributionTestResource, &distribution),
					resource.TestCheckResourceAttr(constant.CdnDistributionResource+"."+constant.CdnDistributionTestResource, "domain", "example.com"),
					resource.TestCheckResourceAttr(constant.CdnDistributionResource+"."+constant.CdnDistributionTestResource, "routing_rules.0.scheme", "http"),
					resource.TestCheckResourceAttr(constant.CdnDistributionResource+"."+constant.CdnDistributionTestResource, "routing_rules.0.prefix", "/api"),
					resource.TestCheckResourceAttr(constant.CdnDistributionResource+"."+constant.CdnDistributionTestResource, "routing_rules.0.upstream.0.host", "server.example.com"),
					resource.TestCheckResourceAttr(constant.CdnDistributionResource+"."+constant.CdnDistributionTestResource, "routing_rules.0.upstream.0.caching", "true"),
					resource.TestCheckResourceAttr(constant.CdnDistributionResource+"."+constant.CdnDistributionTestResource, "routing_rules.0.upstream.0.waf", "true"),
					resource.TestCheckResourceAttr(constant.CdnDistributionResource+"."+constant.CdnDistributionTestResource, "routing_rules.0.upstream.0.rate_limit_class", "R100"),
					resource.TestCheckResourceAttr(constant.CdnDistributionResource+"."+constant.CdnDistributionTestResource, "routing_rules.0.upstream.0.geo_restrictions.0.allow_list.0", "RO"),
				),
			},
			{
				Config: testAccDataSourceCdnDistributionMatchId,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(constant.CdnDistributionResource+"."+constant.CdnDistributionTestResource, "domain", "example.com"),
					resource.TestCheckResourceAttr(constant.CdnDistributionResource+"."+constant.CdnDistributionTestResource, "routing_rules.0.scheme", "http"),
					resource.TestCheckResourceAttr(constant.CdnDistributionResource+"."+constant.CdnDistributionTestResource, "routing_rules.0.prefix", "/api"),
					resource.TestCheckResourceAttr(constant.CdnDistributionResource+"."+constant.CdnDistributionTestResource, "routing_rules.0.upstream.0.host", "server.example.com"),
					resource.TestCheckResourceAttr(constant.CdnDistributionResource+"."+constant.CdnDistributionTestResource, "routing_rules.0.upstream.0.caching", "true"),
					resource.TestCheckResourceAttr(constant.CdnDistributionResource+"."+constant.CdnDistributionTestResource, "routing_rules.0.upstream.0.waf", "true"),
					resource.TestCheckResourceAttr(constant.CdnDistributionResource+"."+constant.CdnDistributionTestResource, "routing_rules.0.upstream.0.rate_limit_class", "R100"),
					resource.TestCheckResourceAttr(constant.CdnDistributionResource+"."+constant.CdnDistributionTestResource, "routing_rules.0.upstream.0.geo_restrictions.0.allow_list.0", "RO"),
				),
			},
			{
				Config: testAccDataSourceCdnDistributionMatchDomain,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(constant.CdnDistributionResource+"."+constant.CdnDistributionTestResource, "domain", "example.com"),
					resource.TestCheckResourceAttr(constant.CdnDistributionResource+"."+constant.CdnDistributionTestResource, "routing_rules.0.scheme", "http"),
					resource.TestCheckResourceAttr(constant.CdnDistributionResource+"."+constant.CdnDistributionTestResource, "routing_rules.0.prefix", "/api"),
					resource.TestCheckResourceAttr(constant.CdnDistributionResource+"."+constant.CdnDistributionTestResource, "routing_rules.0.upstream.0.host", "server.example.com"),
					resource.TestCheckResourceAttr(constant.CdnDistributionResource+"."+constant.CdnDistributionTestResource, "routing_rules.0.upstream.0.caching", "true"),
					resource.TestCheckResourceAttr(constant.CdnDistributionResource+"."+constant.CdnDistributionTestResource, "routing_rules.0.upstream.0.waf", "true"),
					resource.TestCheckResourceAttr(constant.CdnDistributionResource+"."+constant.CdnDistributionTestResource, "routing_rules.0.upstream.0.rate_limit_class", "R100"),
					resource.TestCheckResourceAttr(constant.CdnDistributionResource+"."+constant.CdnDistributionTestResource, "routing_rules.0.upstream.0.geo_restrictions.0.allow_list.0", "RO"),
				),
			},
			{
				Config:      testAccDataSourceCdnDistributionMultipleResultsError,
				ExpectError: regexp.MustCompile("more than one registry found with the specified criteria"),
			},
			{
				Config:      testAccDataSourceCdnDistributionWrongDomainError,
				ExpectError: regexp.MustCompile("no distribution found with the specified criteria"),
			},
			{
				Config: testAccCheckCdnDistributionConfigUpdate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCdnDistributionExists(constant.CdnDistributionResource+"."+constant.CdnDistributionTestResource, &distribution),
					resource.TestCheckResourceAttr(constant.CdnDistributionResource+"."+constant.CdnDistributionTestResource, "domain", "example.example.com"),
					resource.TestCheckResourceAttr(constant.CdnDistributionResource+"."+constant.CdnDistributionTestResource, "routing_rules.0.scheme", "http/https"),
					resource.TestCheckResourceAttr(constant.CdnDistributionResource+"."+constant.CdnDistributionTestResource, "routing_rules.0.prefix", "/api2"),
					resource.TestCheckResourceAttr(constant.CdnDistributionResource+"."+constant.CdnDistributionTestResource, "routing_rules.0.upstream.0.host", "server.server.example.com"),
					resource.TestCheckResourceAttr(constant.CdnDistributionResource+"."+constant.CdnDistributionTestResource, "routing_rules.0.upstream.0.caching", "false"),
					resource.TestCheckResourceAttr(constant.CdnDistributionResource+"."+constant.CdnDistributionTestResource, "routing_rules.0.upstream.0.waf", "true"),
					resource.TestCheckResourceAttr(constant.CdnDistributionResource+"."+constant.CdnDistributionTestResource, "routing_rules.0.upstream.0.rate_limit_class", "R10"),
					resource.TestCheckResourceAttr(constant.CdnDistributionResource+"."+constant.CdnDistributionTestResource, "routing_rules.0.upstream.0.geo_restrictions.0.block_list.0", "RO"),
					resource.TestCheckResourceAttr(constant.CdnDistributionResource+"."+constant.CdnDistributionTestResource, "routing_rules.1.scheme", "https"),
					resource.TestCheckResourceAttr(constant.CdnDistributionResource+"."+constant.CdnDistributionTestResource, "routing_rules.1.prefix", "/api3"),
					resource.TestCheckResourceAttr(constant.CdnDistributionResource+"."+constant.CdnDistributionTestResource, "routing_rules.1.upstream.0.host", "server2.example.com"),
					resource.TestCheckResourceAttr(constant.CdnDistributionResource+"."+constant.CdnDistributionTestResource, "routing_rules.1.upstream.0.caching", "true"),
					resource.TestCheckResourceAttr(constant.CdnDistributionResource+"."+constant.CdnDistributionTestResource, "routing_rules.1.upstream.0.waf", "false"),
					resource.TestCheckResourceAttr(constant.CdnDistributionResource+"."+constant.CdnDistributionTestResource, "routing_rules.1.upstream.0.rate_limit_class", "R10000"),
					resource.TestCheckResourceAttr(constant.CdnDistributionResource+"."+constant.CdnDistributionTestResource, "routing_rules.1.upstream.0.geo_restrictions.0.allow_list.0", "CN"),
					resource.TestCheckResourceAttr(constant.CdnDistributionResource+"."+constant.CdnDistributionTestResource, "routing_rules.1.upstream.0.geo_restrictions.0.allow_list.1", "RU"),
				),
			},
		},
	})
}

func testAccCheckCdnDistributionDestroyCheck(s *terraform.State) error {
	client := testAccProvider.Meta().(services.SdkBundle).CdnClient

	ctx, cancel := context.WithTimeout(context.Background(), *resourceDefaultTimeouts.Default)

	if cancel != nil {
		defer cancel()
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != constant.CdnDistributionResource {
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

func testAccCheckCdnDistributionExists(n string, distribution *ionoscloud_cdn.Distribution) resource.TestCheckFunc {
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

const testAccDataSourceCdnDistributionMatchId = testAccCheckCdnDistributionConfigBasic + `
data ` + constant.CdnDistributionResource + ` ` + constant.CdnDistributionDataSourceById + ` {
  id			= ` + constant.CdnDistributionResource + `.` + constant.CdnDistributionTestResource + `.id
}`

const testAccDataSourceCdnDistributionMatchDomain = testAccCheckCdnDistributionConfigBasic + `
data ` + constant.CdnDistributionResource + ` ` + constant.CdnDistributionDataSourceByDomain + ` {
    domain = ` + constant.CdnDistributionResource + `.` + constant.CdnDistributionTestResource + `.domain
}`

const testAccDataSourceCdnDistributionMultipleResultsError = testAccCheckCdnDistributionConfigBasic + `
resource ` + constant.CdnDistributionResource + ` ` + constant.CdnDistributionTestResource + `_multiple_results {
	domain         = "example.com"
	routing_rules {
		scheme = "http"
		prefix = "/api2"
		upstream {
			host             = "server.server.example.com"
			caching          = false
			waf              = true
			rate_limit_class = "R10"
			geo_restrictions {
				block_list = [ "RO"]
			}
		}
	}
}

data ` + constant.CdnDistributionResource + ` ` + constant.CdnDistributionDataSourceMatching + ` {
    domain = ` + constant.CdnDistributionResource + `.` + constant.CdnDistributionTestResource + `.domain
}`

const testAccDataSourceCdnDistributionWrongDomainError = `
data ` + constant.CdnDistributionResource + ` ` + "test_distribution_matching" + ` {
    domain =  "wrong.domain.com"
}`

const testAccCheckCdnDistributionConfigUpdate = `
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
` + `resource ` + constant.CdnDistributionResource + ` ` + constant.CdnDistributionTestResource + ` {
	domain         = "example.example.com"
	certificate_id = ` + constant.CertificateResource + `.` + constant.TestCertName + `.id` + `
	routing_rules {
		scheme = "http/https"
		prefix = "/api2"
		upstream {
			host             = "server.server.example.com"
			caching          = false
			waf              = true
			rate_limit_class = "R10"
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
			geo_restrictions {
				allow_list = [ "CN", "RU"]
			}
		}
	}
}`
