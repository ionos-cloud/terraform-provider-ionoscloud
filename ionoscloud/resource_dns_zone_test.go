//go:build all || dns

package ionoscloud

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	dns "github.com/ionos-cloud/sdk-go-dns"
	"regexp"
	"testing"
)

func TestAccDNSZone(t *testing.T) {
	var Zone dns.ZoneResponse

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccDNSZoneDestroyCheck,
		Steps: []resource.TestStep{
			{
				Config: DNSZoneConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccDNSZoneExistenceCheck(DNSZoneResource+"."+DNSZoneTestResourceName, &Zone),
					resource.TestCheckResourceAttr(DNSZoneResource+"."+DNSZoneTestResourceName, zoneNameAttribute, zoneNameValue),
					resource.TestCheckResourceAttr(DNSZoneResource+"."+DNSZoneTestResourceName, zoneDescriptionAttribute, zoneDescriptionValue),
					resource.TestCheckResourceAttr(DNSZoneResource+"."+DNSZoneTestResourceName, zoneEnabledAttribute, zoneEnabledValue),
				),
			},
			{
				Config: DNSZoneDataSourceMatchById,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(DataSource+"."+DNSZoneResource+"."+DNSZoneTestDataSourceName, zoneNameAttribute, DNSZoneResource+"."+DNSZoneTestResourceName, zoneNameAttribute),
					resource.TestCheckResourceAttrPair(DataSource+"."+DNSZoneResource+"."+DNSZoneTestDataSourceName, zoneDescriptionAttribute, DNSZoneResource+"."+DNSZoneTestResourceName, zoneDescriptionAttribute),
					resource.TestCheckResourceAttrPair(DataSource+"."+DNSZoneResource+"."+DNSZoneTestDataSourceName, zoneEnabledAttribute, DNSZoneResource+"."+DNSZoneTestResourceName, zoneEnabledAttribute),
				),
			},
			{
				Config: DNSZoneDataSourceMatchByName,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(DataSource+"."+DNSZoneResource+"."+DNSZoneTestDataSourceName, zoneNameAttribute, DNSZoneResource+"."+DNSZoneTestResourceName, zoneNameAttribute),
					resource.TestCheckResourceAttrPair(DataSource+"."+DNSZoneResource+"."+DNSZoneTestDataSourceName, zoneDescriptionAttribute, DNSZoneResource+"."+DNSZoneTestResourceName, zoneDescriptionAttribute),
					resource.TestCheckResourceAttrPair(DataSource+"."+DNSZoneResource+"."+DNSZoneTestDataSourceName, zoneEnabledAttribute, DNSZoneResource+"."+DNSZoneTestResourceName, zoneEnabledAttribute),
				),
			},
			{
				Config: DNSZoneDataSourceMatchByNamePartialMatch,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(DataSource+"."+DNSZoneResource+"."+DNSZoneTestDataSourceName, zoneNameAttribute, DNSZoneResource+"."+DNSZoneTestResourceName, zoneNameAttribute),
					resource.TestCheckResourceAttrPair(DataSource+"."+DNSZoneResource+"."+DNSZoneTestDataSourceName, zoneDescriptionAttribute, DNSZoneResource+"."+DNSZoneTestResourceName, zoneDescriptionAttribute),
					resource.TestCheckResourceAttrPair(DataSource+"."+DNSZoneResource+"."+DNSZoneTestDataSourceName, zoneEnabledAttribute, DNSZoneResource+"."+DNSZoneTestResourceName, zoneEnabledAttribute),
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
					testAccDNSZoneExistenceCheck(DNSZoneResource+"."+DNSZoneTestResourceName, &Zone),
					resource.TestCheckResourceAttr(DNSZoneResource+"."+DNSZoneTestResourceName, zoneNameAttribute, zoneNameValue),
					resource.TestCheckResourceAttr(DNSZoneResource+"."+DNSZoneTestResourceName, zoneDescriptionAttribute, zoneUpdatedDescriptionValue),
					resource.TestCheckResourceAttr(DNSZoneResource+"."+DNSZoneTestResourceName, zoneEnabledAttribute, zoneupdatedEnabledValue),
				),
			},
		},
	})
}

func testAccDNSZoneDestroyCheck(s *terraform.State) error {
	client := testAccProvider.Meta().(SdkBundle).DNSClient
	ctx, cancel := context.WithTimeout(context.Background(), *resourceDefaultTimeouts.Default)
	defer cancel()

	for _, rs := range s.RootModule().Resources {
		if rs.Type != DNSZoneResource {
			continue
		}
		zoneId := rs.Primary.ID
		_, apiResponse, err := client.GetZoneById(ctx, zoneId)
		if err != nil {
			if !apiResponse.HttpNotFound() {
				return fmt.Errorf("an error occured while checking the destruction of DNS Zone with ID: %s, error: %w", zoneId, err)
			}
		} else {
			return fmt.Errorf("DNS Zone with ID: %s still exists", zoneId)
		}
	}
	return nil
}

func testAccDNSZoneExistenceCheck(path string, zone *dns.ZoneResponse) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(SdkBundle).DNSClient
		rs, ok := s.RootModule().Resources[path]

		if !ok {
			return fmt.Errorf("not found: %s", path)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no ID is set for the DNS Zone")
		}
		ctx, cancel := context.WithTimeout(context.Background(), *resourceDefaultTimeouts.Default)
		defer cancel()
		zoneId := rs.Primary.ID
		zoneResponse, _, err := client.GetZoneById(ctx, zoneId)

		if err != nil {
			return fmt.Errorf("an error occured while fetching DNS Zone with ID: %s, error: %w", zoneId, err)
		}
		zone = &zoneResponse
		return nil
	}
}

const DNSZoneDataSourceMatchById = DNSZoneConfig + `
` + DataSource + ` ` + DNSZoneResource + ` ` + DNSZoneTestDataSourceName + ` {
	id = ` + DNSZoneResource + `.` + DNSZoneTestResourceName + `.id
}
`

const DNSZoneDataSourceMatchByName = DNSZoneConfig + `
` + DataSource + ` ` + DNSZoneResource + ` ` + DNSZoneTestDataSourceName + ` {
	name = ` + DNSZoneResource + `.` + DNSZoneTestResourceName + `.name
}
`

var DNSZoneDataSourceMatchByNamePartialMatch = DNSZoneConfig + `
` + DataSource + ` ` + DNSZoneResource + ` ` + DNSZoneTestDataSourceName + ` {
	name = "` + zoneNameValue[:4] + `"
	partial_match = true
}
`

const DNSZoneDataSourceInvalidBothIDAndName = DNSZoneConfig + `
` + DataSource + ` ` + DNSZoneResource + ` ` + DNSZoneTestDataSourceName + ` {
	name = ` + DNSZoneResource + `.` + DNSZoneTestResourceName + `.name
	id = ` + DNSZoneResource + `.` + DNSZoneTestResourceName + `.id
}
`

const DNSZoneDataSourceInvalidNoIDNoName = `
` + DataSource + ` ` + DNSZoneResource + ` ` + DNSZoneTestDataSourceName + ` {
}
`

const DNSZoneDataSourceInvalidPartialMatchUsedWithID = DNSZoneConfig + `
` + DataSource + ` ` + DNSZoneResource + ` ` + DNSZoneTestDataSourceName + ` {
	id = ` + DNSZoneResource + `.` + DNSZoneTestResourceName + `.id
	partial_match = true
}
`

const DNSZoneDataSourceWrongNameError = `
` + DataSource + ` ` + DNSZoneResource + ` ` + DNSZoneTestDataSourceName + ` {
	name = "nonexistent"
}
`

const DNSZoneDataSourceWrongPartialNameError = `
` + DataSource + ` ` + DNSZoneResource + ` ` + DNSZoneTestDataSourceName + ` {
	name = "nonexistent"
	partial_match = true
}
`

const DNSZoneConfigUpdate = `
resource ` + DNSZoneResource + ` ` + DNSZoneTestResourceName + ` {
	` + zoneNameAttribute + ` = "` + zoneNameValue + `"
	` + zoneDescriptionAttribute + ` = "` + zoneUpdatedDescriptionValue + `"
    ` + zoneEnabledAttribute + ` = ` + zoneupdatedEnabledValue + `
}
`
