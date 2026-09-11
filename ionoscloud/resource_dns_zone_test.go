//go:build all || dns

package ionoscloud

import (
	"context"
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/querycheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"

	dns "github.com/ionos-cloud/sdk-go-bundle/products/dns/v2"

	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/bundleclient"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils/constant"
)

// TestAccDNSZoneQuery covers `terraform query` against the dns zone list resource, and
// the resource identity it streams. It uses its own zone name and its own resource label:
// querycheck.ExpectLength asserts a CONTRACT-WIDE total, so reusing the suite's shared
// fixture would make ExpectLength(1) flap as soon as another test creates a zone.
func TestAccDNSZoneQuery(t *testing.T) {
	const (
		queryZoneName        = "tf-test-query.com"
		queryZoneDescription = "zone for the query acceptance test"
		queryZoneUpdatedDesc = "zone for the query acceptance test, updated"
		queryZoneAddr        = constant.DNSZoneResource + ".test_dns_zone_query"
	)

	resource.Test(t, resource.TestCase{
		PreCheck: func() { testAccPreCheck(t) },
		// `terraform query` and list blocks were introduced in Terraform 1.14.
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_14_0),
		},
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactoriesInternal(t, &testAccProvider),
		CheckDestroy:             testAccDNSZoneDestroyCheck,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
resource %[1]q "test_dns_zone_query" {
  name        = %[2]q
  description = %[3]q
  enabled     = true
}`, constant.DNSZoneResource, queryZoneName, queryZoneDescription),
			},
			// List without filters: the zone must show up with its identity. A dns zone
			// has no location, so the identity is a lone id.
			{
				Query: true,
				Config: fmt.Sprintf(`list %[1]q "test_dns_zone_query" {
  provider = ionoscloud
}`, constant.DNSZoneResource),
				QueryResultChecks: []querycheck.QueryResultCheck{
					querycheck.ExpectIdentity(queryZoneAddr, map[string]knownvalue.Check{
						"id": knownvalue.NotNull(),
					}),
				},
			},
			// Filter by name: the unique zone name guarantees exactly one result.
			{
				Query: true,
				Config: fmt.Sprintf(`list %[1]q "test_dns_zone_query" {
  provider = ionoscloud
  config {
    filters = [
      { field_name = "name", field_value = %[2]q },
    ]
  }
}`, constant.DNSZoneResource, queryZoneName),
				QueryResultChecks: []querycheck.QueryResultCheck{
					querycheck.ExpectLength(queryZoneAddr, 1),
				},
			},
			// Same name, a description that does not belong to it. A dns zone has no
			// location, so `description` is the only other allow-listed field that can
			// discriminate, and this is what proves the filters are ANDed.
			{
				Query: true,
				Config: fmt.Sprintf(`list %[1]q "test_dns_zone_query" {
  provider = ionoscloud
  config {
    filters = [
      { field_name = "name",        field_value = %[2]q },
      { field_name = "description", field_value = %[3]q },
    ]
  }
}`, constant.DNSZoneResource, queryZoneName, queryZoneUpdatedDesc),
				QueryResultChecks: []querycheck.QueryResultCheck{
					querycheck.ExpectLength(queryZoneAddr, 0),
				},
			},
			// Description only. `name` is ForceNew on a dns zone, so this is the one step
			// that actually enters zoneUpdate - which writes the identity itself, because
			// it does not delegate to the read.
			{
				Config: fmt.Sprintf(`
resource %[1]q "test_dns_zone_query" {
  name        = %[2]q
  description = %[3]q
  enabled     = true
}`, constant.DNSZoneResource, queryZoneName, queryZoneUpdatedDesc),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(queryZoneAddr, "name", queryZoneName),
					resource.TestCheckResourceAttr(queryZoneAddr, "description", queryZoneUpdatedDesc),
				),
			},
			// The same AND filter as above, which matched nothing before the update and
			// must match exactly this zone after it - so the listing is reading the
			// updated state rather than a cached one.
			{
				Query: true,
				Config: fmt.Sprintf(`list %[1]q "test_dns_zone_query" {
  provider = ionoscloud
  config {
    filters = [
      { field_name = "name",        field_value = %[2]q },
      { field_name = "description", field_value = %[3]q },
    ]
  }
}`, constant.DNSZoneResource, queryZoneName, queryZoneUpdatedDesc),
				QueryResultChecks: []querycheck.QueryResultCheck{
					querycheck.ExpectLength(queryZoneAddr, 1),
				},
			},
			// Import through the resource identity that the list results carry. This kind
			// already checks that the import succeeds, that the plan it leaves behind is a
			// no-op and that the planned identity matches the one in state; ImportStateVerify
			// cannot be combined with it, only ImportCommandWithID reads that field.
			{
				ResourceName:    queryZoneAddr,
				ImportState:     true,
				ImportStateKind: resource.ImportBlockWithResourceIdentity,
			},
		},
	})
}

func TestAccDNSZone(t *testing.T) {
	var Zone dns.ZoneRead

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactoriesInternal(t, &testAccProvider),
		CheckDestroy:             testAccDNSZoneDestroyCheck,
		Steps: []resource.TestStep{
			{
				Config: DNSZoneConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccDNSZoneExistenceCheck(constant.DNSZoneResource+"."+constant.DNSZoneTestResourceName, &Zone),
					resource.TestCheckResourceAttr(constant.DNSZoneResource+"."+constant.DNSZoneTestResourceName, zoneNameAttribute, zoneNameValue),
					resource.TestCheckResourceAttr(constant.DNSZoneResource+"."+constant.DNSZoneTestResourceName, zoneDescriptionAttribute, zoneDescriptionValue),
					resource.TestCheckResourceAttr(constant.DNSZoneResource+"."+constant.DNSZoneTestResourceName, zoneEnabledAttribute, zoneEnabledValue),
				),
			},
			{
				Config: DNSZoneDataSourceMatchById,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.DNSZoneResource+"."+constant.DNSZoneTestDataSourceName, zoneNameAttribute, constant.DNSZoneResource+"."+constant.DNSZoneTestResourceName, zoneNameAttribute),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.DNSZoneResource+"."+constant.DNSZoneTestDataSourceName, zoneDescriptionAttribute, constant.DNSZoneResource+"."+constant.DNSZoneTestResourceName, zoneDescriptionAttribute),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.DNSZoneResource+"."+constant.DNSZoneTestDataSourceName, zoneEnabledAttribute, constant.DNSZoneResource+"."+constant.DNSZoneTestResourceName, zoneEnabledAttribute),
				),
			},
			{
				Config: DNSZoneDataSourceMatchByName,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.DNSZoneResource+"."+constant.DNSZoneTestDataSourceName, zoneNameAttribute, constant.DNSZoneResource+"."+constant.DNSZoneTestResourceName, zoneNameAttribute),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.DNSZoneResource+"."+constant.DNSZoneTestDataSourceName, zoneDescriptionAttribute, constant.DNSZoneResource+"."+constant.DNSZoneTestResourceName, zoneDescriptionAttribute),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.DNSZoneResource+"."+constant.DNSZoneTestDataSourceName, zoneEnabledAttribute, constant.DNSZoneResource+"."+constant.DNSZoneTestResourceName, zoneEnabledAttribute),
				),
			},
			{
				Config: DNSZoneDataSourceMatchByNamePartialMatch,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.DNSZoneResource+"."+constant.DNSZoneTestDataSourceName, zoneNameAttribute, constant.DNSZoneResource+"."+constant.DNSZoneTestResourceName, zoneNameAttribute),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.DNSZoneResource+"."+constant.DNSZoneTestDataSourceName, zoneDescriptionAttribute, constant.DNSZoneResource+"."+constant.DNSZoneTestResourceName, zoneDescriptionAttribute),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.DNSZoneResource+"."+constant.DNSZoneTestDataSourceName, zoneEnabledAttribute, constant.DNSZoneResource+"."+constant.DNSZoneTestResourceName, zoneEnabledAttribute),
				),
			},
			{
				Config:      DNSZoneDataSourceInvalidBothIDAndName,
				ExpectError: regexp.MustCompile("ID and name cannot be both specified at the same time"),
			},
			{
				Config:      DNSZoneDataSourceInvalidNoIDNoName,
				ExpectError: regexp.MustCompile("please provide either the DNS Zone ID or name"),
			},
			{
				Config:      DNSZoneDataSourceInvalidPartialMatchUsedWithID,
				ExpectError: regexp.MustCompile("partial_match can only be used together with the name attribute"),
			},
			{
				Config:      DNSZoneDataSourceWrongNameError,
				ExpectError: regexp.MustCompile("no DNS Zone found with the specified name"),
			},
			{
				Config:      DNSZoneDataSourceWrongPartialNameError,
				ExpectError: regexp.MustCompile("no DNS Zone found with the specified name"),
			},
			{
				Config: DNSZoneConfigUpdate,
				Check: resource.ComposeTestCheckFunc(
					testAccDNSZoneExistenceCheck(constant.DNSZoneResource+"."+constant.DNSZoneTestResourceName, &Zone),
					resource.TestCheckResourceAttr(constant.DNSZoneResource+"."+constant.DNSZoneTestResourceName, zoneNameAttribute, zoneNameValue),
					resource.TestCheckResourceAttr(constant.DNSZoneResource+"."+constant.DNSZoneTestResourceName, zoneDescriptionAttribute, zoneUpdatedDescriptionValue),
					resource.TestCheckResourceAttr(constant.DNSZoneResource+"."+constant.DNSZoneTestResourceName, zoneEnabledAttribute, zoneupdatedEnabledValue),
				),
			},
		},
	})
}

func testAccDNSZoneDestroyCheck(s *terraform.State) error {
	client := testAccProvider.Meta().(bundleclient.SdkBundle).DNSClient
	ctx, cancel := context.WithTimeout(context.Background(), *resourceDefaultTimeouts.Default)
	defer cancel()

	for _, rs := range s.RootModule().Resources {
		if rs.Type != constant.DNSZoneResource {
			continue
		}
		zoneID := rs.Primary.ID
		_, apiResponse, err := client.GetZoneById(ctx, zoneID)
		if err != nil {
			if !apiResponse.HttpNotFound() {
				return fmt.Errorf("an error occurred while checking the destruction of DNS Zone with ID: %s, error: %w", zoneID, err)
			}
		} else {
			return fmt.Errorf("DNS Zone with ID: %s still exists", zoneID)
		}
	}
	return nil
}

func testAccDNSZoneExistenceCheck(path string, zone *dns.ZoneRead) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(bundleclient.SdkBundle).DNSClient
		rs, ok := s.RootModule().Resources[path]

		if !ok {
			return fmt.Errorf("not found: %s", path)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no ID is set for the DNS Zone")
		}
		ctx, cancel := context.WithTimeout(context.Background(), *resourceDefaultTimeouts.Default)
		defer cancel()
		zoneID := rs.Primary.ID
		zoneResponse, _, err := client.GetZoneById(ctx, zoneID)

		if err != nil {
			return fmt.Errorf("an error occurred while fetching DNS Zone with ID: %s, error: %w", zoneID, err)
		}
		zone = &zoneResponse
		return nil
	}
}

const DNSZoneDataSourceMatchById = DNSZoneConfig + `
` + constant.DataSource + ` ` + constant.DNSZoneResource + ` ` + constant.DNSZoneTestDataSourceName + ` {
	id = ` + constant.DNSZoneResource + `.` + constant.DNSZoneTestResourceName + `.id
}
`

const DNSZoneDataSourceMatchByName = DNSZoneConfig + `
` + constant.DataSource + ` ` + constant.DNSZoneResource + ` ` + constant.DNSZoneTestDataSourceName + ` {
	name = ` + constant.DNSZoneResource + `.` + constant.DNSZoneTestResourceName + `.name
}
`

var DNSZoneDataSourceMatchByNamePartialMatch = DNSZoneConfig + `
` + constant.DataSource + ` ` + constant.DNSZoneResource + ` ` + constant.DNSZoneTestDataSourceName + ` {
	name = "` + zoneNameValue[:5] + `"
	partial_match = true
}
`

const DNSZoneDataSourceInvalidBothIDAndName = DNSZoneConfig + `
` + constant.DataSource + ` ` + constant.DNSZoneResource + ` ` + constant.DNSZoneTestDataSourceName + ` {
	name = ` + constant.DNSZoneResource + `.` + constant.DNSZoneTestResourceName + `.name
	id = ` + constant.DNSZoneResource + `.` + constant.DNSZoneTestResourceName + `.id
}
`

const DNSZoneDataSourceInvalidNoIDNoName = `
` + constant.DataSource + ` ` + constant.DNSZoneResource + ` ` + constant.DNSZoneTestDataSourceName + ` {
}
`

const DNSZoneDataSourceInvalidPartialMatchUsedWithID = DNSZoneConfig + `
` + constant.DataSource + ` ` + constant.DNSZoneResource + ` ` + constant.DNSZoneTestDataSourceName + ` {
	id = ` + constant.DNSZoneResource + `.` + constant.DNSZoneTestResourceName + `.id
	partial_match = true
}
`

const DNSZoneDataSourceWrongNameError = `
` + constant.DataSource + ` ` + constant.DNSZoneResource + ` ` + constant.DNSZoneTestDataSourceName + ` {
	name = "nonexistent"
}
`

const DNSZoneDataSourceWrongPartialNameError = `
` + constant.DataSource + ` ` + constant.DNSZoneResource + ` ` + constant.DNSZoneTestDataSourceName + ` {
	name = "nonexistent"
	partial_match = true
}
`

const DNSZoneConfigUpdate = `
resource ` + constant.DNSZoneResource + ` ` + constant.DNSZoneTestResourceName + ` {
	` + zoneNameAttribute + ` = "` + zoneNameValue + `"
	` + zoneDescriptionAttribute + ` = "` + zoneUpdatedDescriptionValue + `"
    ` + zoneEnabledAttribute + ` = ` + zoneupdatedEnabledValue + `
}
`
