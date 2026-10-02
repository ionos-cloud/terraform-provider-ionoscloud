//go:build all || dns

package dns_test

import (
	"context"
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"

	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/internal/acctest"
)

const (
	dnssecResource   = "ionoscloud_dns_zone_dnssec.test"
	dnssecDataSource = "data.ionoscloud_dns_zone_dnssec.test"
)

func TestAccZoneDNSSEC(t *testing.T) {
	zoneName := acctest.GenerateRandomResourceName("tf-dnssec-") + ".com"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		CheckDestroy:             checkZoneDestroy,
		Steps: []resource.TestStep{
			{
				Config: dnssecConfig(zoneName, 2048, 1024),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dnssecResource, "zone_id", "ionoscloud_dns_zone.test", "id"),
					resource.TestCheckResourceAttr(dnssecResource, "algorithm", "RSASHA256"),
					resource.TestCheckResourceAttr(dnssecResource, "ksk_bits", "2048"),
					resource.TestCheckResourceAttr(dnssecResource, "zsk_bits", "1024"),
					resource.TestCheckResourceAttr(dnssecResource, "nsec_mode", "NSEC3"),
					resource.TestCheckResourceAttr(dnssecResource, "algorithm_number", "8"),
					resource.TestCheckResourceAttrSet(dnssecResource, "key_tag"),
					resource.TestCheckResourceAttrSet(dnssecResource, "digest"),
					resource.TestCheckResourceAttrSet(dnssecResource, "digest_type"),
					resource.TestCheckResourceAttrSet(dnssecResource, "public_key"),
					resource.TestCheckResourceAttrSet(dnssecResource, "composed_key_data"),
					resource.TestMatchResourceAttr(dnssecResource, "ds_record", regexp.MustCompile(`^\d+ 8 \d [0-9A-Fa-f]+$`)),
				),
			},
			{
				// The API does not return the creation parameters (bits, iterations, salt, validity).
				ResourceName:            dnssecResource,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"ksk_bits", "zsk_bits", "nsec3_iterations", "nsec3_salt_bits", "validity", "timeouts"},
			},
			{
				Config: dnssecConfig(zoneName, 2048, 1024) + `
data "ionoscloud_dns_zone_dnssec" "test" {
  zone_id = ionoscloud_dns_zone_dnssec.test.zone_id
}`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dnssecDataSource, "ds_record", dnssecResource, "ds_record"),
					resource.TestCheckResourceAttrPair(dnssecDataSource, "key_tag", dnssecResource, "key_tag"),
					resource.TestCheckResourceAttrPair(dnssecDataSource, "digest", dnssecResource, "digest"),
				),
			},
		},
	})
}

func dnssecConfig(zoneName string, kskBits, zskBits int) string {
	return fmt.Sprintf(`
resource "ionoscloud_dns_zone" "test" {
  name = %q
}

resource "ionoscloud_dns_zone_dnssec" "test" {
  zone_id          = ionoscloud_dns_zone.test.id
  ksk_bits         = %d
  zsk_bits         = %d
  nsec_mode        = "NSEC3"
  nsec3_iterations = 10
  nsec3_salt_bits  = 64
  validity         = 120
}
`, zoneName, kskBits, zskBits)
}

func checkZoneDestroy(s *terraform.State) error {
	client := acctest.NewTestBundleClientFromEnv().DNSClient
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "ionoscloud_dns_zone" {
			continue
		}
		_, apiResponse, err := client.GetZoneById(context.Background(), rs.Primary.ID)
		if err == nil {
			return fmt.Errorf("DNS zone %s still exists", rs.Primary.ID)
		}
		if !apiResponse.HttpNotFound() {
			return fmt.Errorf("unexpected error checking DNS zone %s: %w", rs.Primary.ID, err)
		}
	}
	return nil
}
