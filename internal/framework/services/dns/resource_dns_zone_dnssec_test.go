//go:build all || dns

package dns_test

import (
	"context"
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
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
				// Forget the key without deleting it, so that it can be imported into the state used by the next steps.
				Config: dnssecZoneConfig(zoneName) + `
removed {
  from = ionoscloud_dns_zone_dnssec.test
  lifecycle {
    destroy = false
  }
}`,
			},
			{
				Config:             dnssecConfig(zoneName, 2048, 1024),
				ResourceName:       dnssecResource,
				ImportState:        true,
				ImportStateIdFunc:  zoneIDFunc("ionoscloud_dns_zone.test"),
				ImportStatePersist: true,
			},
			{
				// The imported state lacks the creation parameters: they are adopted from the configuration in place,
				// without recreating the key.
				Config: dnssecConfig(zoneName, 2048, 1024),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{plancheck.ExpectResourceAction(dnssecResource, plancheck.ResourceActionUpdate)},
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dnssecResource, "algorithm", "RSASHA256"),
					resource.TestCheckResourceAttr(dnssecResource, "nsec_mode", "NSEC3"),
					resource.TestCheckResourceAttr(dnssecResource, "ksk_bits", "2048"),
					resource.TestCheckResourceAttr(dnssecResource, "validity", "120"),
				),
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

func dnssecZoneConfig(zoneName string) string {
	return fmt.Sprintf(`
resource "ionoscloud_dns_zone" "test" {
  name = %q
}
`, zoneName)
}

func dnssecConfig(zoneName string, kskBits, zskBits int) string {
	return dnssecZoneConfig(zoneName) + fmt.Sprintf(`
resource "ionoscloud_dns_zone_dnssec" "test" {
  zone_id          = ionoscloud_dns_zone.test.id
  ksk_bits         = %d
  zsk_bits         = %d
  nsec_mode        = "NSEC3"
  nsec3_iterations = 10
  nsec3_salt_bits  = 64
  validity         = 120
}
`, kskBits, zskBits)
}

// zoneIDFunc returns the ID of the given zone resource, to import a resource that is no longer in the state.
func zoneIDFunc(zoneResource string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[zoneResource]
		if !ok {
			return "", fmt.Errorf("resource %s not found in state", zoneResource)
		}
		return rs.Primary.ID, nil
	}
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
