//go:build all || dns

package ionoscloud

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	dns "github.com/ionos-cloud/sdk-go-dns"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils/constant"
	"regexp"
	"testing"
)

func TestAccDNSZone(t *testing.T) {
	var Zone dns.ZoneRead

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
	client := testAccProvider.Meta().(services.SdkBundle).DNSClient
	ctx, cancel := context.WithTimeout(context.Background(), *resourceDefaultTimeouts.Default)
	defer cancel()

	for _, rs := range s.RootModule().Resources {
		if rs.Type != constant.DNSZoneResource {
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

func testAccDNSZoneExistenceCheck(path string, zone *dns.ZoneRead) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(services.SdkBundle).DNSClient
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
	name = "` + zoneNameValue[:4] + `"
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
