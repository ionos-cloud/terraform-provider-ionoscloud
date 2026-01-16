//go:build all || dns

package ionoscloud

import (
	"context"
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"

	dns "github.com/ionos-cloud/sdk-go-bundle/products/dns/v2"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/bundleclient"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils/constant"
)

func TestAccDNSReverseRecord(t *testing.T) {
	var ReverseRecord dns.ReverseRecordRead

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactoriesInternal(t, &testAccProvider),
		CheckDestroy:             testAccDNSReverseRecordDestroyCheck,
		Steps: []resource.TestStep{
			{
				Config: DNSReverseRecordConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccDNSReverseRecordExistenceCheck(constant.DNSReverseRecordResource+"."+constant.DNSReverseRecordTestResourceName, &ReverseRecord),
					resource.TestCheckResourceAttr(constant.DNSReverseRecordResource+"."+constant.DNSReverseRecordTestResourceName, "name", reverseRecordNameValue),
					resource.TestCheckResourceAttr(constant.DNSReverseRecordResource+"."+constant.DNSReverseRecordTestResourceName, "description", reverseRecordDescValue),
					resource.TestCheckResourceAttrSet(constant.DNSReverseRecordResource+"."+constant.DNSReverseRecordTestResourceName, "ip"),
				),
			},
			{
				Config: DNSReverseRecordDataSourceMatchById,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(constant.DataSource+"."+constant.DNSReverseRecordDataSource+"."+constant.DNSReverseRecordTestDataSourceName, "name", reverseRecordNameValue),
					resource.TestCheckResourceAttr(constant.DataSource+"."+constant.DNSReverseRecordDataSource+"."+constant.DNSReverseRecordTestDataSourceName, "description", reverseRecordDescValue),
					resource.TestCheckResourceAttrSet(constant.DataSource+"."+constant.DNSReverseRecordDataSource+"."+constant.DNSReverseRecordTestDataSourceName, "ip"),
				),
			},
			{
				Config: DNSReverseRecordDataSourceMatchByName,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(constant.DataSource+"."+constant.DNSReverseRecordDataSource+"."+constant.DNSReverseRecordTestDataSourceName, "name", reverseRecordNameValue),
					resource.TestCheckResourceAttr(constant.DataSource+"."+constant.DNSReverseRecordDataSource+"."+constant.DNSReverseRecordTestDataSourceName, "description", reverseRecordDescValue),
					resource.TestCheckResourceAttrSet(constant.DataSource+"."+constant.DNSReverseRecordDataSource+"."+constant.DNSReverseRecordTestDataSourceName, "ip"),
				),
			},
			{
				Config: DNSReverseRecordDataSourceMatchByIp,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(constant.DataSource+"."+constant.DNSReverseRecordDataSource+"."+constant.DNSReverseRecordTestDataSourceName, "name", reverseRecordNameValue),
					resource.TestCheckResourceAttr(constant.DataSource+"."+constant.DNSReverseRecordDataSource+"."+constant.DNSReverseRecordTestDataSourceName, "description", reverseRecordDescValue),
					resource.TestCheckResourceAttrSet(constant.DataSource+"."+constant.DNSReverseRecordDataSource+"."+constant.DNSReverseRecordTestDataSourceName, "ip"),
				),
			},
			{
				Config: DNSReverseRecordDataSourceMatchByNamePartialMatch,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(constant.DataSource+"."+constant.DNSReverseRecordDataSource+"."+constant.DNSReverseRecordTestDataSourceName, "name", reverseRecordNameValue),
					resource.TestCheckResourceAttr(constant.DataSource+"."+constant.DNSReverseRecordDataSource+"."+constant.DNSReverseRecordTestDataSourceName, "description", reverseRecordDescValue),
					resource.TestCheckResourceAttrSet(constant.DataSource+"."+constant.DNSReverseRecordDataSource+"."+constant.DNSReverseRecordTestDataSourceName, "ip"),
				),
			},
			{
				Config:      DNSReverseRecordDataSourceInvalidBothIDAndName,
				ExpectError: regexp.MustCompile(`only one of \[Id, name, ip\] can be specified at the same time`),
			},
			{
				Config:      DNSReverseRecordDataSourceInvalidNoIDNoName,
				ExpectError: regexp.MustCompile("please provide either the DNS Record Id, name or IP"),
			},
			{
				Config:      DNSReverseRecordDataSourceInvalidPartialMatchUsedWithID,
				ExpectError: regexp.MustCompile("partial_match can only be used together with the name attribute"),
			},
			{
				Config:      DNSReverseRecordDataSourceWrongNameError,
				ExpectError: regexp.MustCompile("no DNS Reverse Record found with the specified filter"),
			},
			{
				Config:      DNSReverseRecordDataSourceWrongPartialNameError,
				ExpectError: regexp.MustCompile("no DNS Reverse Record found with the specified filter"),
			},
			{
				Config:      DNSReverseRecordDataSourceWrongIpError,
				ExpectError: regexp.MustCompile("no DNS Reverse Record found with the specified filter"),
			},
			{
				Config: DNSReverseRecordsDataSourceMatchByName,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(constant.DataSource+"."+constant.DNSReverseRecordsDataSource+"."+constant.DNSReverseRecordsTestDataSourceName, "reverse_records.#", "1"),
					resource.TestCheckResourceAttr(constant.DataSource+"."+constant.DNSReverseRecordsDataSource+"."+constant.DNSReverseRecordsTestDataSourceName, "reverse_records.0.name", reverseRecordNameValue),
					resource.TestCheckResourceAttr(constant.DataSource+"."+constant.DNSReverseRecordsDataSource+"."+constant.DNSReverseRecordsTestDataSourceName, "reverse_records.0.description", reverseRecordDescValue),
					resource.TestCheckResourceAttrSet(constant.DataSource+"."+constant.DNSReverseRecordsDataSource+"."+constant.DNSReverseRecordsTestDataSourceName, "reverse_records.0.ip"),
				),
			},
			{
				Config: DNSReverseRecordsDataSourceMatchByNamePartial,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(constant.DataSource+"."+constant.DNSReverseRecordsDataSource+"."+constant.DNSReverseRecordsTestDataSourceName, "reverse_records.#", "2"),
				),
			},
			{
				Config: DNSReverseRecordsDataSourceMatchByIp,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(constant.DataSource+"."+constant.DNSReverseRecordsDataSource+"."+constant.DNSReverseRecordsTestDataSourceName, "reverse_records.#", "1"),
					resource.TestCheckResourceAttr(constant.DataSource+"."+constant.DNSReverseRecordsDataSource+"."+constant.DNSReverseRecordsTestDataSourceName, "reverse_records.0.name", reverseRecordNameValue),
					resource.TestCheckResourceAttr(constant.DataSource+"."+constant.DNSReverseRecordsDataSource+"."+constant.DNSReverseRecordsTestDataSourceName, "reverse_records.0.description", reverseRecordDescValue),
					resource.TestCheckResourceAttrSet(constant.DataSource+"."+constant.DNSReverseRecordsDataSource+"."+constant.DNSReverseRecordsTestDataSourceName, "reverse_records.0.ip"),
				),
			},
			{
				Config: DNSReverseRecordsDataSourceMatchByIpAndNameNoMatch,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(constant.DataSource+"."+constant.DNSReverseRecordsDataSource+"."+constant.DNSReverseRecordsTestDataSourceName, "reverse_records.#", "0"),
				),
			},
			{
				Config: DNSReverseRecordsDataSourceMatchByIpAndNamePartialMatch,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(constant.DataSource+"."+constant.DNSReverseRecordsDataSource+"."+constant.DNSReverseRecordsTestDataSourceName, "reverse_records.#", "1"),
					resource.TestCheckResourceAttr(constant.DataSource+"."+constant.DNSReverseRecordsDataSource+"."+constant.DNSReverseRecordsTestDataSourceName, "reverse_records.0.name", reverseRecordNameValueUpdated),
					resource.TestCheckResourceAttr(constant.DataSource+"."+constant.DNSReverseRecordsDataSource+"."+constant.DNSReverseRecordsTestDataSourceName, "reverse_records.0.description", reverseRecordDescValueUpdated),
					resource.TestCheckResourceAttrSet(constant.DataSource+"."+constant.DNSReverseRecordsDataSource+"."+constant.DNSReverseRecordsTestDataSourceName, "reverse_records.0.ip"),
				),
			},
			{
				Config: DNSReverseRecordsDataSourceMatchByIps,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(constant.DataSource+"."+constant.DNSReverseRecordsDataSource+"."+constant.DNSReverseRecordsTestDataSourceName, "reverse_records.#", "2"),
				),
			},
			{
				Config: DNSReverseRecordConfigUpdate,
				Check: resource.ComposeTestCheckFunc(
					testAccDNSReverseRecordExistenceCheck(constant.DNSReverseRecordResource+"."+constant.DNSReverseRecordTestResourceName, &ReverseRecord),
					resource.TestCheckResourceAttr(constant.DNSReverseRecordResource+"."+constant.DNSReverseRecordTestResourceName, "name", reverseRecordNameValueUpdated),
					resource.TestCheckResourceAttr(constant.DNSReverseRecordResource+"."+constant.DNSReverseRecordTestResourceName, "description", reverseRecordDescValueUpdated),
					resource.TestCheckResourceAttrSet(constant.DNSReverseRecordResource+"."+constant.DNSReverseRecordTestResourceName, "ip"),
				),
			},
		},
	})
}

func testAccDNSReverseRecordDestroyCheck(s *terraform.State) error {
	client := testAccProvider.Meta().(bundleclient.SdkBundle).DNSClient
	ctx, cancel := context.WithTimeout(context.Background(), *resourceDefaultTimeouts.Default)
	defer cancel()

	for _, rs := range s.RootModule().Resources {
		if rs.Type != constant.DNSReverseRecordResource {
			continue
		}
		recordId := rs.Primary.ID
		_, apiResponse, err := client.GetReverseRecordById(ctx, recordId)
		if err != nil {
			if !apiResponse.HttpNotFound() {
				return fmt.Errorf("an error occurred while checking the destruction of DNS Reverse Record with ID: %s, error: %w", recordId, err)
			}
		} else {
			return fmt.Errorf("DNS Reverse Record with ID: %s still exists", recordId)
		}
	}
	return nil
}

func testAccDNSReverseRecordExistenceCheck(path string, record *dns.ReverseRecordRead) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(bundleclient.SdkBundle).DNSClient
		rs, ok := s.RootModule().Resources[path]

		if !ok {
			return fmt.Errorf("not found: %s", path)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no ID is set for the DNS Record")
		}
		ctx, cancel := context.WithTimeout(context.Background(), *resourceDefaultTimeouts.Default)
		defer cancel()
		recordId := rs.Primary.ID
		recordResponse, _, err := client.GetReverseRecordById(ctx, recordId)
		if err != nil {
			return fmt.Errorf("an error occurred while fetching DNS Record with ID: %s, error: %w", recordId, err)
		}
		record = &recordResponse
		return nil
	}
}

const DNSReverseRecordDataSourceMatchById = DNSReverseRecordConfig + `
` + constant.DataSource + ` ` + constant.DNSReverseRecordResource + ` ` + constant.DNSReverseRecordTestDataSourceName + ` {
	id = ` + constant.DNSReverseRecordResource + `.` + constant.DNSReverseRecordTestResourceName + `.id
}
`

const DNSReverseRecordDataSourceMatchByName = DNSReverseRecordConfig + `
` + constant.DataSource + ` ` + constant.DNSReverseRecordResource + ` ` + constant.DNSReverseRecordTestDataSourceName + ` {
	name = ` + constant.DNSReverseRecordResource + `.` + constant.DNSReverseRecordTestResourceName + `.name
}
`

const DNSReverseRecordDataSourceMatchByIp = DNSReverseRecordConfig + `
` + constant.DataSource + ` ` + constant.DNSReverseRecordResource + ` ` + constant.DNSReverseRecordTestDataSourceName + ` {
	ip = ` + constant.IpBlockResource + `.` + constant.IpBlockTestResource + `.ips[0]
}
`

var DNSReverseRecordDataSourceMatchByNamePartialMatch = DNSReverseRecordConfig + `
` + constant.DataSource + ` ` + constant.DNSReverseRecordResource + ` ` + constant.DNSReverseRecordTestDataSourceName + ` {
	name = "` + reverseRecordNameValue[:4] + `"
	partial_match = true
}
`

const DNSReverseRecordDataSourceInvalidBothIDAndName = DNSReverseRecordConfig + `
` + constant.DataSource + ` ` + constant.DNSReverseRecordResource + ` ` + constant.DNSReverseRecordTestDataSourceName + ` {
	name = ` + constant.DNSReverseRecordResource + `.` + constant.DNSReverseRecordTestResourceName + `.name
	id = ` + constant.DNSReverseRecordResource + `.` + constant.DNSReverseRecordTestResourceName + `.id
}
`

const DNSReverseRecordDataSourceInvalidNoIDNoName = DNSReverseRecordConfig + `
` + constant.DataSource + ` ` + constant.DNSReverseRecordResource + ` ` + constant.DNSReverseRecordTestDataSourceName + ` {
}
`

const DNSReverseRecordDataSourceInvalidPartialMatchUsedWithID = DNSReverseRecordConfig + `
` + constant.DataSource + ` ` + constant.DNSReverseRecordResource + ` ` + constant.DNSReverseRecordTestDataSourceName + ` {
	id = ` + constant.DNSReverseRecordResource + `.` + constant.DNSReverseRecordTestResourceName + `.id
	partial_match = true
}
`

const DNSReverseRecordDataSourceWrongNameError = DNSReverseRecordConfig + `
` + constant.DataSource + ` ` + constant.DNSReverseRecordResource + ` ` + constant.DNSReverseRecordTestDataSourceName + ` {
	name = "nonexistent"
}
`

const DNSReverseRecordDataSourceWrongPartialNameError = DNSReverseRecordConfig + `
` + constant.DataSource + ` ` + constant.DNSReverseRecordResource + ` ` + constant.DNSReverseRecordTestDataSourceName + ` {
	name = "nonexistent"
	partial_match = true
}
`

const DNSReverseRecordsDataSourceMatchByName = DNSReverseRecordsConfig + `
` + constant.DataSource + ` ` + constant.DNSReverseRecordsDataSource + ` ` + constant.DNSReverseRecordsTestDataSourceName + ` {
	name = ` + constant.DNSReverseRecordResource + `.` + constant.DNSReverseRecordTestResourceName + `.name
}
`

const DNSReverseRecordsDataSourceMatchByNamePartial = DNSReverseRecordsConfig + `
` + constant.DataSource + ` ` + constant.DNSReverseRecordsDataSource + ` ` + constant.DNSReverseRecordsTestDataSourceName + ` {
	name = ` + constant.DNSReverseRecordResource + `.` + constant.DNSReverseRecordTestResourceName + `.name
	partial_match = true
}
`

const DNSReverseRecordsDataSourceMatchByIp = DNSReverseRecordsConfig + `
` + constant.DataSource + ` ` + constant.DNSReverseRecordsDataSource + ` ` + constant.DNSReverseRecordsTestDataSourceName + ` {
	ips = [` + constant.IpBlockResource + `.` + constant.IpBlockTestResource + `.ips[0]]
}
`

const DNSReverseRecordsDataSourceMatchByIpAndNameNoMatch = DNSReverseRecordsConfig + `
` + constant.DataSource + ` ` + constant.DNSReverseRecordsDataSource + ` ` + constant.DNSReverseRecordsTestDataSourceName + ` {
	ips = [` + constant.IpBlockResource + `.` + constant.IpBlockTestResource + `.ips[1]]
	name = ` + constant.DNSReverseRecordResource + `.` + constant.DNSReverseRecordTestResourceName + `.name
}
`

const DNSReverseRecordsDataSourceMatchByIpAndNamePartialMatch = DNSReverseRecordsConfig + `
` + constant.DataSource + ` ` + constant.DNSReverseRecordsDataSource + ` ` + constant.DNSReverseRecordsTestDataSourceName + ` {
	ips = [` + constant.IpBlockResource + `.` + constant.IpBlockTestResource + `.ips[1]]
	name = ` + constant.DNSReverseRecordResource + `.` + constant.DNSReverseRecordTestResourceName + `.name
	partial_match = true
}
`

const DNSReverseRecordsDataSourceMatchByIps = DNSReverseRecordsConfig + `
` + constant.DataSource + ` ` + constant.DNSReverseRecordsDataSource + ` ` + constant.DNSReverseRecordsTestDataSourceName + ` {
	ips = [` + constant.IpBlockResource + `.` + constant.IpBlockTestResource + `.ips[0], ` + constant.IpBlockResource + `.` + constant.IpBlockTestResource + `.ips[1]]
}
`

const DNSReverseRecordDataSourceWrongIpError = DNSReverseRecordConfig + `
` + constant.DataSource + ` ` + constant.DNSReverseRecordResource + ` ` + constant.DNSReverseRecordTestDataSourceName + ` {
	ip = "127.0.0.1"
}
`

const DNSReverseRecordConfigUpdate = `
resource ` + constant.IpBlockResource + ` ` + constant.IpBlockTestResource + ` {
  location = "de/fra"
  size = 1
  name = "` + constant.IpBlockTestResource + `"
}` + `
resource ` + constant.DNSReverseRecordResource + ` ` + constant.DNSReverseRecordTestResourceName + ` {
  name = "` + reverseRecordNameValueUpdated + `"
  description = "` + reverseRecordDescValueUpdated + `"
  ip = ` + constant.IpBlockResource + `.` + constant.IpBlockTestResource + `.ips[0]
}
`
