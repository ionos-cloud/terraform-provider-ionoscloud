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

// TestAccDNSZoneQuery exercises the ionoscloud_dns_zone list resource and the resource
// identity that listing depends on.
//
// The list resource is served by the plugin-framework half of the provider even though
// the DNS zone resource itself is implemented with SDKv2, so this also covers the mux
// serving the two halves under the same type name. See resource_dns_zone_list.go in this
// package.
func TestAccDNSZoneQuery(t *testing.T) {
	const (
		queryZoneName        = "tf-test-query-zone.com"
		queryZoneDescription = "the zone the query test looks for"
		queryZoneUpdatedDesc = "the zone the query test looks for, updated"
		queryZoneAddr        = constant.DNSZoneResource + ".test_dns_zone_query"
		otherDescription     = "a description no zone in this test carries"
	)

	// Its own name, label and create step rather than DNSZoneConfig: ExpectLength asserts
	// a contract-wide total, so a zone another test in this package creates under the same
	// name would make ExpectLength(1) flap.

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
}`, constant.DNSZoneResource, queryZoneName, queryZoneDescription),
			},
			// name is force-new, so description is what an update step can change: zoneUpdate
			// does not delegate to zoneRead, and this is the only step that enters it - where
			// the identity is written again.
			{
				Config: fmt.Sprintf(`
resource %[1]q "test_dns_zone_query" {
  name        = %[2]q
  description = %[3]q
}`, constant.DNSZoneResource, queryZoneName, queryZoneUpdatedDesc),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(queryZoneAddr, zoneDescriptionAttribute, queryZoneUpdatedDesc),
				),
			},
			// List without filters: the zone must show up with its identity.
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
			// Filter by name and description: the unique name guarantees exactly one result.
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
			// Same name, a description no zone carries: proves the description filter is
			// evaluated rather than ignored.
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
}`, constant.DNSZoneResource, queryZoneName, otherDescription),
				QueryResultChecks: []querycheck.QueryResultCheck{
					querycheck.ExpectLength(queryZoneAddr, 0),
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
