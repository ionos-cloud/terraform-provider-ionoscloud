//go:build all || dns

package ionoscloud

import (
	"context"
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
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

// TestAccDNSZoneQuery exercises the ionoscloud_dns_zone list resource and the
// resource identity that listing depends on.
//
// The list resource is served by the plugin-framework half of the provider even though
// the DNS Zone resource itself is implemented with SDKv2, so this also covers the
// mux serving the two halves under the same type name. See
// resource_dns_zone_list.go in this package.
func TestAccDNSZoneQuery(t *testing.T) {
	const (
		dnsZoneName               = "tf-test-dns-zone-query.com"
		dnsZoneDescription        = "DNS Zone for the list resource acceptance test"
		dnsZoneUpdatedDescription = "DNS Zone for the list resource acceptance test, updated"
		dnsZoneAddr               = constant.DNSZoneResource + ".test_dns_zone_query"
		otherDescription          = "no DNS Zone has this description" // the zero-result filter value: no step sets it
	)

	// Own name, own label, own create step - not DNSZoneConfig:
	// ExpectLength asserts a contract-wide total, so a same-named resource made by
	// another test in the package makes ExpectLength(1) flap.

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
}`, constant.DNSZoneResource, dnsZoneName, dnsZoneDescription),
			},
			// List without filters: the DNS Zone must show up with its identity.
			{
				Query: true,
				Config: fmt.Sprintf(`list %[1]q "test_dns_zone_query" {
  provider = ionoscloud
}`, constant.DNSZoneResource),
				QueryResultChecks: []querycheck.QueryResultCheck{
					querycheck.ExpectIdentity(dnsZoneAddr, map[string]knownvalue.Check{
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
}`, constant.DNSZoneResource, dnsZoneName, dnsZoneDescription),
				QueryResultChecks: []querycheck.QueryResultCheck{
					querycheck.ExpectLength(dnsZoneAddr, 1),
				},
			},
			// Same name, different description: proves the description filter is evaluated.
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
}`, constant.DNSZoneResource, dnsZoneName, otherDescription),
				QueryResultChecks: []querycheck.QueryResultCheck{
					querycheck.ExpectLength(dnsZoneAddr, 0),
				},
			},
			// Change the description in place: zoneUpdate does not end in a read but writes
			// the identity itself, which the SDK compares with the planned one on this apply.
			{
				Config: fmt.Sprintf(`
resource %[1]q "test_dns_zone_query" {
  name        = %[2]q
  description = %[3]q
}`, constant.DNSZoneResource, dnsZoneName, dnsZoneUpdatedDescription),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(dnsZoneAddr, plancheck.ResourceActionUpdate),
					},
				},
				Check: resource.TestCheckResourceAttr(dnsZoneAddr, "description", dnsZoneUpdatedDescription),
			},
			// Import through the resource identity that the list results carry. This kind
			// already checks that the import succeeds, that the plan it leaves behind is a
			// no-op and that the planned identity matches the one in state; ImportStateVerify
			// cannot be combined with it, only ImportCommandWithID reads that field.
			{
				ResourceName:    dnsZoneAddr,
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
