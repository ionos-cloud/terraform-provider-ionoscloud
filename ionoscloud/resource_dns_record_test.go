//go:build all || dns

package ionoscloud

import (
	"context"
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	dns "github.com/ionos-cloud/sdk-go-dns"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils/constant"
)

func TestAccDNSRecord(t *testing.T) {
	var Record dns.RecordRead

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccDNSRecordDestroyCheck,
		Steps: []resource.TestStep{
			{
				Config: DNSRecordConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccDNSRecordExistenceCheck(constant.DNSRecordResource+"."+constant.DNSRecordTestResourceName, &Record),
					resource.TestCheckResourceAttr(constant.DNSRecordResource+"."+constant.DNSRecordTestResourceName, recordNameAttribute, recordNameValue),
					resource.TestCheckResourceAttr(constant.DNSRecordResource+"."+constant.DNSRecordTestResourceName, recordTypeAttribute, recordTypeValue),
					resource.TestCheckResourceAttr(constant.DNSRecordResource+"."+constant.DNSRecordTestResourceName, recordContentAttribute, recordContentValue),
					resource.TestCheckResourceAttr(constant.DNSRecordResource+"."+constant.DNSRecordTestResourceName, recordTtlAttribute, recordTtlValue),
					resource.TestCheckResourceAttr(constant.DNSRecordResource+"."+constant.DNSRecordTestResourceName, recordEnabledAttribute, recordEnabledValue),
				),
			},
			{
				Config: DNSRecordDataSourceMatchById,
				Check: resource.ComposeTestCheckFunc(
					testAccDNSRecordExistenceCheck(constant.DNSRecordResource+"."+constant.DNSRecordTestResourceName, &Record),
					resource.TestCheckResourceAttr(constant.DNSRecordResource+"."+constant.DNSRecordTestResourceName, recordNameAttribute, recordNameValue),
					resource.TestCheckResourceAttr(constant.DNSRecordResource+"."+constant.DNSRecordTestResourceName, recordTypeAttribute, recordTypeValue),
					resource.TestCheckResourceAttr(constant.DNSRecordResource+"."+constant.DNSRecordTestResourceName, recordContentAttribute, recordContentValue),
					resource.TestCheckResourceAttr(constant.DNSRecordResource+"."+constant.DNSRecordTestResourceName, recordTtlAttribute, recordTtlValue),
					resource.TestCheckResourceAttr(constant.DNSRecordResource+"."+constant.DNSRecordTestResourceName, recordEnabledAttribute, recordEnabledValue),
				),
			},
			{
				Config: DNSRecordDataSourceMatchByName,
				Check: resource.ComposeTestCheckFunc(
					testAccDNSRecordExistenceCheck(constant.DNSRecordResource+"."+constant.DNSRecordTestResourceName, &Record),
					resource.TestCheckResourceAttr(constant.DNSRecordResource+"."+constant.DNSRecordTestResourceName, recordNameAttribute, recordNameValue),
					resource.TestCheckResourceAttr(constant.DNSRecordResource+"."+constant.DNSRecordTestResourceName, recordTypeAttribute, recordTypeValue),
					resource.TestCheckResourceAttr(constant.DNSRecordResource+"."+constant.DNSRecordTestResourceName, recordContentAttribute, recordContentValue),
					resource.TestCheckResourceAttr(constant.DNSRecordResource+"."+constant.DNSRecordTestResourceName, recordTtlAttribute, recordTtlValue),
					resource.TestCheckResourceAttr(constant.DNSRecordResource+"."+constant.DNSRecordTestResourceName, recordEnabledAttribute, recordEnabledValue),
				),
			},
			{
				Config: DNSRecordDataSourceMatchByNamePartialMatch,
				Check: resource.ComposeTestCheckFunc(
					testAccDNSRecordExistenceCheck(constant.DNSRecordResource+"."+constant.DNSRecordTestResourceName, &Record),
					resource.TestCheckResourceAttr(constant.DNSRecordResource+"."+constant.DNSRecordTestResourceName, recordNameAttribute, recordNameValue),
					resource.TestCheckResourceAttr(constant.DNSRecordResource+"."+constant.DNSRecordTestResourceName, recordTypeAttribute, recordTypeValue),
					resource.TestCheckResourceAttr(constant.DNSRecordResource+"."+constant.DNSRecordTestResourceName, recordContentAttribute, recordContentValue),
					resource.TestCheckResourceAttr(constant.DNSRecordResource+"."+constant.DNSRecordTestResourceName, recordTtlAttribute, recordTtlValue),
					resource.TestCheckResourceAttr(constant.DNSRecordResource+"."+constant.DNSRecordTestResourceName, recordEnabledAttribute, recordEnabledValue),
				),
			},
			{
				Config:      DNSRecordDataSourceInvalidBothIDAndName,
				ExpectError: regexp.MustCompile("ID and name cannot be both specified at the same time"),
			},
			{
				Config:      DNSRecordDataSourceInvalidNoIDNoName,
				ExpectError: regexp.MustCompile("please provide either the DNS Record ID or name"),
			},
			{
				Config:      DNSRecordDataSourceInvalidPartialMatchUsedWithID,
				ExpectError: regexp.MustCompile("partial_match can only be used together with the name attribute"),
			},
			{
				Config:      DNSRecordDataSourceWrongNameError,
				ExpectError: regexp.MustCompile("no DNS Record found with the specified name"),
			},
			{
				Config:      DNSRecordDataSourceWrongPartialNameError,
				ExpectError: regexp.MustCompile("no DNS Record found with the specified name"),
			},
			{
				Config: DNSRecordConfigUpdate,
				Check: resource.ComposeTestCheckFunc(
					testAccDNSRecordExistenceCheck(constant.DNSRecordResource+"."+constant.DNSRecordTestResourceName, &Record),
					resource.TestCheckResourceAttr(constant.DNSRecordResource+"."+constant.DNSRecordTestResourceName, recordNameAttribute, recordNameValue),
					resource.TestCheckResourceAttr(constant.DNSRecordResource+"."+constant.DNSRecordTestResourceName, recordTypeAttribute, recordTypeValue),
					resource.TestCheckResourceAttr(constant.DNSRecordResource+"."+constant.DNSRecordTestResourceName, recordContentAttribute, recordUpdatedContentValue),
					resource.TestCheckResourceAttr(constant.DNSRecordResource+"."+constant.DNSRecordTestResourceName, recordTtlAttribute, recordUpdatedTtlValue),
					resource.TestCheckResourceAttr(constant.DNSRecordResource+"."+constant.DNSRecordTestResourceName, recordEnabledAttribute, recordUpdatedEnabledValue),
				),
			},
		},
	})
}

func testAccDNSRecordDestroyCheck(s *terraform.State) error {
	client := testAccProvider.Meta().(services.SdkBundle).DNSClient
	ctx, cancel := context.WithTimeout(context.Background(), *resourceDefaultTimeouts.Default)
	defer cancel()

	for _, rs := range s.RootModule().Resources {
		if rs.Type != constant.DNSRecordResource {
			continue
		}
		zoneId := rs.Primary.Attributes["zone_id"]
		recordId := rs.Primary.ID
		_, apiResponse, err := client.GetRecordById(ctx, zoneId, recordId)
		if err != nil {
			if !apiResponse.HttpNotFound() {
				return fmt.Errorf("an error occured while checking the destruction of DNS Record with ID: %s, zone ID: %s, error: %w", recordId, zoneId, err)
			}
		} else {
			return fmt.Errorf("DNS Record with ID: %s still exists, zone ID: %s", recordId, zoneId)
		}
	}
	return nil
}

func testAccDNSRecordExistenceCheck(path string, record *dns.RecordRead) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(services.SdkBundle).DNSClient
		rs, ok := s.RootModule().Resources[path]

		if !ok {
			return fmt.Errorf("not found: %s", path)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no ID is set for the DNS Record")
		}
		ctx, cancel := context.WithTimeout(context.Background(), *resourceDefaultTimeouts.Default)
		defer cancel()
		zoneId := rs.Primary.Attributes["zone_id"]
		recordId := rs.Primary.ID
		recordResponse, _, err := client.GetRecordById(ctx, zoneId, recordId)
		if err != nil {
			return fmt.Errorf("an error occured while fetching DNS Record with ID: %s, zone ID: %s, error: %w", recordId, zoneId, err)
		}
		record = &recordResponse
		return nil
	}
}

const DNSRecordDataSourceMatchById = DNSRecordConfig + `
` + constant.DataSource + ` ` + constant.DNSRecordResource + ` ` + constant.DNSRecordTestDataSourceName + ` {
	zone_id = ` + constant.DNSRecordResource + `.` + constant.DNSRecordTestResourceName + `.zone_id
	id = ` + constant.DNSRecordResource + `.` + constant.DNSRecordTestResourceName + `.id
}
`

const DNSRecordDataSourceMatchByName = DNSRecordConfig + `
` + constant.DataSource + ` ` + constant.DNSRecordResource + ` ` + constant.DNSRecordTestDataSourceName + ` {
	zone_id = ` + constant.DNSRecordResource + `.` + constant.DNSRecordTestResourceName + `.zone_id
	name = ` + constant.DNSRecordResource + `.` + constant.DNSRecordTestResourceName + `.name
}
`

var DNSRecordDataSourceMatchByNamePartialMatch = DNSRecordConfig + `
` + constant.DataSource + ` ` + constant.DNSRecordResource + ` ` + constant.DNSRecordTestDataSourceName + ` {
	zone_id = ` + constant.DNSRecordResource + `.` + constant.DNSRecordTestResourceName + `.zone_id
	name = "` + recordNameValue[:4] + `"
	partial_match = true
}
`

const DNSRecordDataSourceInvalidBothIDAndName = DNSRecordConfig + `
` + constant.DataSource + ` ` + constant.DNSRecordResource + ` ` + constant.DNSRecordTestDataSourceName + ` {
	zone_id = ` + constant.DNSRecordResource + `.` + constant.DNSRecordTestResourceName + `.zone_id
	name = ` + constant.DNSRecordResource + `.` + constant.DNSRecordTestResourceName + `.name
	id = ` + constant.DNSRecordResource + `.` + constant.DNSRecordTestResourceName + `.id
}
`

const DNSRecordDataSourceInvalidNoIDNoName = DNSRecordConfig + `
` + constant.DataSource + ` ` + constant.DNSRecordResource + ` ` + constant.DNSRecordTestDataSourceName + ` {
	zone_id = ` + constant.DNSRecordResource + `.` + constant.DNSRecordTestResourceName + `.zone_id
}
`

const DNSRecordDataSourceInvalidPartialMatchUsedWithID = DNSRecordConfig + `
` + constant.DataSource + ` ` + constant.DNSRecordResource + ` ` + constant.DNSRecordTestDataSourceName + ` {
	zone_id = ` + constant.DNSRecordResource + `.` + constant.DNSRecordTestResourceName + `.zone_id
	id = ` + constant.DNSRecordResource + `.` + constant.DNSRecordTestResourceName + `.id
	partial_match = true
}
`

const DNSRecordDataSourceWrongNameError = DNSRecordConfig + `
` + constant.DataSource + ` ` + constant.DNSRecordResource + ` ` + constant.DNSRecordTestDataSourceName + ` {
	zone_id = ` + constant.DNSRecordResource + `.` + constant.DNSRecordTestResourceName + `.zone_id
	name = "nonexistent"
}
`

const DNSRecordDataSourceWrongPartialNameError = DNSRecordConfig + `
` + constant.DataSource + ` ` + constant.DNSRecordResource + ` ` + constant.DNSRecordTestDataSourceName + ` {
	zone_id = ` + constant.DNSRecordResource + `.` + constant.DNSRecordTestResourceName + `.zone_id
	name = "nonexistent"
	partial_match = true
}
`

const DNSRecordConfigUpdate = DNSZoneConfig + `
resource ` + constant.DNSRecordResource + ` ` + constant.DNSRecordTestResourceName + ` {
	zone_id = ` + constant.DNSZoneResource + `.` + constant.DNSZoneTestResourceName + `.id
	` + recordNameAttribute + ` = "` + recordNameValue + `"
	` + recordTypeAttribute + ` = "` + recordTypeValue + `"
	` + recordContentAttribute + ` = "` + recordUpdatedContentValue + `"
	` + recordTtlAttribute + ` = ` + recordUpdatedTtlValue + `
	` + recordPriorityAttribute + ` = ` + recordPriorityValue + `
	` + recordEnabledAttribute + ` = ` + recordUpdatedEnabledValue + `
}
`
