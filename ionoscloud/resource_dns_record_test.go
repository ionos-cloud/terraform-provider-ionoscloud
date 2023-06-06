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
					testAccDNSRecordExistenceCheck(DNSRecordResource+"."+DNSRecordTestResourceName, &Record),
					resource.TestCheckResourceAttr(DNSRecordResource+"."+DNSRecordTestResourceName, recordNameAttribute, recordNameValue),
					resource.TestCheckResourceAttr(DNSRecordResource+"."+DNSRecordTestResourceName, recordTypeAttribute, recordTypeValue),
					resource.TestCheckResourceAttr(DNSRecordResource+"."+DNSRecordTestResourceName, recordContentAttribute, recordContentValue),
					resource.TestCheckResourceAttr(DNSRecordResource+"."+DNSRecordTestResourceName, recordTtlAttribute, recordTtlValue),
					resource.TestCheckResourceAttr(DNSRecordResource+"."+DNSRecordTestResourceName, recordEnabledAttribute, recordEnabledValue),
				),
			},
			{
				Config: DNSRecordDataSourceMatchById,
				Check: resource.ComposeTestCheckFunc(
					testAccDNSRecordExistenceCheck(DNSRecordResource+"."+DNSRecordTestResourceName, &Record),
					resource.TestCheckResourceAttr(DNSRecordResource+"."+DNSRecordTestResourceName, recordNameAttribute, recordNameValue),
					resource.TestCheckResourceAttr(DNSRecordResource+"."+DNSRecordTestResourceName, recordTypeAttribute, recordTypeValue),
					resource.TestCheckResourceAttr(DNSRecordResource+"."+DNSRecordTestResourceName, recordContentAttribute, recordContentValue),
					resource.TestCheckResourceAttr(DNSRecordResource+"."+DNSRecordTestResourceName, recordTtlAttribute, recordTtlValue),
					resource.TestCheckResourceAttr(DNSRecordResource+"."+DNSRecordTestResourceName, recordEnabledAttribute, recordEnabledValue),
				),
			},
			{
				Config: DNSRecordDataSourceMatchByName,
				Check: resource.ComposeTestCheckFunc(
					testAccDNSRecordExistenceCheck(DNSRecordResource+"."+DNSRecordTestResourceName, &Record),
					resource.TestCheckResourceAttr(DNSRecordResource+"."+DNSRecordTestResourceName, recordNameAttribute, recordNameValue),
					resource.TestCheckResourceAttr(DNSRecordResource+"."+DNSRecordTestResourceName, recordTypeAttribute, recordTypeValue),
					resource.TestCheckResourceAttr(DNSRecordResource+"."+DNSRecordTestResourceName, recordContentAttribute, recordContentValue),
					resource.TestCheckResourceAttr(DNSRecordResource+"."+DNSRecordTestResourceName, recordTtlAttribute, recordTtlValue),
					resource.TestCheckResourceAttr(DNSRecordResource+"."+DNSRecordTestResourceName, recordEnabledAttribute, recordEnabledValue),
				),
			},
			{
				Config: DNSRecordDataSourceMatchByNamePartialMatch,
				Check: resource.ComposeTestCheckFunc(
					testAccDNSRecordExistenceCheck(DNSRecordResource+"."+DNSRecordTestResourceName, &Record),
					resource.TestCheckResourceAttr(DNSRecordResource+"."+DNSRecordTestResourceName, recordNameAttribute, recordNameValue),
					resource.TestCheckResourceAttr(DNSRecordResource+"."+DNSRecordTestResourceName, recordTypeAttribute, recordTypeValue),
					resource.TestCheckResourceAttr(DNSRecordResource+"."+DNSRecordTestResourceName, recordContentAttribute, recordContentValue),
					resource.TestCheckResourceAttr(DNSRecordResource+"."+DNSRecordTestResourceName, recordTtlAttribute, recordTtlValue),
					resource.TestCheckResourceAttr(DNSRecordResource+"."+DNSRecordTestResourceName, recordEnabledAttribute, recordEnabledValue),
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
					testAccDNSRecordExistenceCheck(DNSRecordResource+"."+DNSRecordTestResourceName, &Record),
					resource.TestCheckResourceAttr(DNSRecordResource+"."+DNSRecordTestResourceName, recordNameAttribute, recordNameValue),
					resource.TestCheckResourceAttr(DNSRecordResource+"."+DNSRecordTestResourceName, recordTypeAttribute, recordTypeValue),
					resource.TestCheckResourceAttr(DNSRecordResource+"."+DNSRecordTestResourceName, recordContentAttribute, recordUpdatedContentValue),
					resource.TestCheckResourceAttr(DNSRecordResource+"."+DNSRecordTestResourceName, recordTtlAttribute, recordUpdatedTtlValue),
					resource.TestCheckResourceAttr(DNSRecordResource+"."+DNSRecordTestResourceName, recordEnabledAttribute, recordUpdatedEnabledValue),
				),
			},
		},
	})
}

func testAccDNSRecordDestroyCheck(s *terraform.State) error {
	client := testAccProvider.Meta().(SdkBundle).DNSClient
	ctx, cancel := context.WithTimeout(context.Background(), *resourceDefaultTimeouts.Default)
	defer cancel()

	for _, rs := range s.RootModule().Resources {
		if rs.Type != DNSRecordResource {
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
		client := testAccProvider.Meta().(SdkBundle).DNSClient
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
` + DataSource + ` ` + DNSRecordResource + ` ` + DNSRecordTestDataSourceName + ` {
	zone_id = ` + DNSRecordResource + `.` + DNSRecordTestResourceName + `.zone_id
	id = ` + DNSRecordResource + `.` + DNSRecordTestResourceName + `.id
}
`

const DNSRecordDataSourceMatchByName = DNSRecordConfig + `
` + DataSource + ` ` + DNSRecordResource + ` ` + DNSRecordTestDataSourceName + ` {
	zone_id = ` + DNSRecordResource + `.` + DNSRecordTestResourceName + `.zone_id
	name = ` + DNSRecordResource + `.` + DNSRecordTestResourceName + `.name
}
`

var DNSRecordDataSourceMatchByNamePartialMatch = DNSRecordConfig + `
` + DataSource + ` ` + DNSRecordResource + ` ` + DNSRecordTestDataSourceName + ` {
	zone_id = ` + DNSRecordResource + `.` + DNSRecordTestResourceName + `.zone_id
	name = "` + recordNameValue[:4] + `"
	partial_match = true
}
`

const DNSRecordDataSourceInvalidBothIDAndName = DNSRecordConfig + `
` + DataSource + ` ` + DNSRecordResource + ` ` + DNSRecordTestDataSourceName + ` {
	zone_id = ` + DNSRecordResource + `.` + DNSRecordTestResourceName + `.zone_id
	name = ` + DNSRecordResource + `.` + DNSRecordTestResourceName + `.name
	id = ` + DNSRecordResource + `.` + DNSRecordTestResourceName + `.id
}
`

const DNSRecordDataSourceInvalidNoIDNoName = DNSRecordConfig + `
` + DataSource + ` ` + DNSRecordResource + ` ` + DNSRecordTestDataSourceName + ` {
	zone_id = ` + DNSRecordResource + `.` + DNSRecordTestResourceName + `.zone_id
}
`

const DNSRecordDataSourceInvalidPartialMatchUsedWithID = DNSRecordConfig + `
` + DataSource + ` ` + DNSRecordResource + ` ` + DNSRecordTestDataSourceName + ` {
	zone_id = ` + DNSRecordResource + `.` + DNSRecordTestResourceName + `.zone_id
	id = ` + DNSRecordResource + `.` + DNSRecordTestResourceName + `.id
	partial_match = true
}
`

const DNSRecordDataSourceWrongNameError = DNSRecordConfig + `
` + DataSource + ` ` + DNSRecordResource + ` ` + DNSRecordTestDataSourceName + ` {
	zone_id = ` + DNSRecordResource + `.` + DNSRecordTestResourceName + `.zone_id
	name = "nonexistent"
}
`

const DNSRecordDataSourceWrongPartialNameError = DNSRecordConfig + `
` + DataSource + ` ` + DNSRecordResource + ` ` + DNSRecordTestDataSourceName + ` {
	zone_id = ` + DNSRecordResource + `.` + DNSRecordTestResourceName + `.zone_id
	name = "nonexistent"
	partial_match = true
}
`

const DNSRecordConfigUpdate = DNSZoneConfig + `
resource ` + DNSRecordResource + ` ` + DNSRecordTestResourceName + ` {
	zone_id = ` + DNSZoneResource + `.` + DNSZoneTestResourceName + `.id
	` + recordNameAttribute + ` = "` + recordNameValue + `"
	` + recordTypeAttribute + ` = "` + recordTypeValue + `"
	` + recordContentAttribute + ` = "` + recordUpdatedContentValue + `"
	` + recordTtlAttribute + ` = ` + recordUpdatedTtlValue + `
	` + recordPriorityAttribute + ` = ` + recordPriorityValue + `
	` + recordEnabledAttribute + ` = ` + recordUpdatedEnabledValue + `
}
`
