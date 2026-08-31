//go:build compute || all || datacenter

package ionoscloud

import (
	"context"
	"fmt"
	"regexp"
	"testing"

	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/bundleclient"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils/constant"

	ionoscloud "github.com/ionos-cloud/sdk-go/v6"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/querycheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAccDataCenterBasic(t *testing.T) {
	var datacenter ionoscloud.Datacenter

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactoriesInternal(t, &testAccProvider),
		CheckDestroy:             testAccCheckDatacenterDestroyCheck,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckDatacenterConfigBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDatacenterExists(constant.DatacenterResource+"."+constant.DatacenterTestResource, &datacenter),
					resource.TestCheckResourceAttr(constant.DatacenterResource+"."+constant.DatacenterTestResource, "name", constant.DatacenterTestResource),
					resource.TestCheckResourceAttr(constant.DatacenterResource+"."+constant.DatacenterTestResource, "location", "us/las"),
					resource.TestCheckResourceAttr(constant.DatacenterResource+"."+constant.DatacenterTestResource, "description", "Test Datacenter Description"),
					resource.TestCheckResourceAttr(constant.DatacenterResource+"."+constant.DatacenterTestResource, "sec_auth_protection", "false"),
				),
			},
			{
				Config: testAccDataSourceDatacenterMatchID,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.DatacenterResource+"."+constant.DatacenterDataSourceById, "name", constant.DatacenterResource+"."+constant.DatacenterTestResource, "name"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.DatacenterResource+"."+constant.DatacenterDataSourceById, "location", constant.DatacenterResource+"."+constant.DatacenterTestResource, "location"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.DatacenterResource+"."+constant.DatacenterDataSourceById, "description", constant.DatacenterResource+"."+constant.DatacenterTestResource, "description"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.DatacenterResource+"."+constant.DatacenterDataSourceById, "version", constant.DatacenterResource+"."+constant.DatacenterTestResource, "version"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.DatacenterResource+"."+constant.DatacenterDataSourceById, "features", constant.DatacenterResource+"."+constant.DatacenterTestResource, "features"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.DatacenterResource+"."+constant.DatacenterDataSourceById, "sec_auth_protection", constant.DatacenterResource+"."+constant.DatacenterTestResource, "sec_auth_protection"),
				),
			},
			{
				Config: testAccDataSourceDatacenterMatchName,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDatacenterExists(constant.DatacenterResource+"."+constant.DatacenterTestResource, &datacenter),
					resource.TestCheckResourceAttrSet(constant.DatacenterResource+"."+constant.DatacenterTestResource, "id"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.DatacenterResource+"."+constant.DatacenterDataSourceByName, "name", constant.DatacenterResource+"."+constant.DatacenterTestResource, "name"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.DatacenterResource+"."+constant.DatacenterDataSourceByName, "location", constant.DatacenterResource+"."+constant.DatacenterTestResource, "location"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.DatacenterResource+"."+constant.DatacenterDataSourceByName, "description", constant.DatacenterResource+"."+constant.DatacenterTestResource, "description"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.DatacenterResource+"."+constant.DatacenterDataSourceByName, "version", constant.DatacenterResource+"."+constant.DatacenterTestResource, "version"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.DatacenterResource+"."+constant.DatacenterDataSourceByName, "features", constant.DatacenterResource+"."+constant.DatacenterTestResource, "features"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.DatacenterResource+"."+constant.DatacenterDataSourceByName, "sec_auth_protection", constant.DatacenterResource+"."+constant.DatacenterTestResource, "sec_auth_protection"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.DatacenterResource+"."+constant.DatacenterDataSourceByName, "ipv6_cidr_block", constant.DatacenterResource+"."+constant.DatacenterTestResource, "ipv6_cidr_block"),
				),
			},
			{
				Config: testAccDataSourceDatacenterMatching,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.DatacenterResource+"."+constant.DatacenterDataSourceMatching, "name", constant.DatacenterResource+"."+constant.DatacenterTestResource, "name"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.DatacenterResource+"."+constant.DatacenterDataSourceMatching, "location", constant.DatacenterResource+"."+constant.DatacenterTestResource, "location"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.DatacenterResource+"."+constant.DatacenterDataSourceMatching, "description", constant.DatacenterResource+"."+constant.DatacenterTestResource, "description"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.DatacenterResource+"."+constant.DatacenterDataSourceMatching, "version", constant.DatacenterResource+"."+constant.DatacenterTestResource, "version"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.DatacenterResource+"."+constant.DatacenterDataSourceMatching, "features", constant.DatacenterResource+"."+constant.DatacenterTestResource, "features"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.DatacenterResource+"."+constant.DatacenterDataSourceMatching, "sec_auth_protection", constant.DatacenterResource+"."+constant.DatacenterTestResource, "sec_auth_protection"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.DatacenterResource+"."+constant.DatacenterDataSourceMatching, "ipv6_cidr_block", constant.DatacenterResource+"."+constant.DatacenterTestResource, "ipv6_cidr_block"),
				),
			},
			{
				Config:      testAccDataSourceDatacenterMultipleResultsError,
				ExpectError: regexp.MustCompile("more than one datacenter found with the specified criteria"),
			},
			{
				Config:      testAccDataSourceDatacenterWrongNameError,
				ExpectError: regexp.MustCompile("no datacenter found with the specified criteria"),
			},
			{
				Config:      testAccDataSourceDatacenterWrongLocationError,
				ExpectError: regexp.MustCompile("no datacenter found with the specified criteria"),
			},
			{
				Config:      testAccDataSourceDatacenterWrongNameAndLocationError,
				ExpectError: regexp.MustCompile("no datacenter found with the specified criteria"),
			},
			{
				Config:      testAccDataSourceDatacenterNoFilterError,
				ExpectError: regexp.MustCompile(`either id, location or name must be set`),
			},
			{
				Config:      testAccDataSourceDatacenterWrongIdError,
				ExpectError: regexp.MustCompile(`error getting datacenter with id`),
			},
			{
				Config:      testAccDataSourceDatacenterIdNameMismatchError,
				ExpectError: regexp.MustCompile(`name of dc \(UUID=.+, name=.+\) does not match expected name`),
			},
			{
				Config:      testAccDataSourceDatacenterIdLocationMismatchError,
				ExpectError: regexp.MustCompile(`location of dc \(UUID=.+, location=.+\) does not match expected location`),
			},
			{
				Config: testAccCheckDatacenterConfigUpdate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDatacenterExists(constant.DatacenterResource+"."+constant.DatacenterTestResource, &datacenter),
					resource.TestCheckResourceAttr(constant.DatacenterResource+"."+constant.DatacenterTestResource, "name", constant.UpdatedResources),
					resource.TestCheckResourceAttr(constant.DatacenterResource+"."+constant.DatacenterTestResource, "location", "us/las"),
					resource.TestCheckResourceAttr(constant.DatacenterResource+"."+constant.DatacenterTestResource, "description", "Test Datacenter Description Updated"),
					resource.TestCheckResourceAttr(constant.DatacenterResource+"."+constant.DatacenterTestResource, "sec_auth_protection", "false"),
				),
			},
		},
	})
}

// TestAccDataCenterQuery exercises the ionoscloud_datacenter list resource and the
// resource identity that listing depends on.
//
// The list resource is served by the plugin-framework half of the provider even though
// the datacenter resource itself is implemented with SDKv2, so this also covers the
// mux serving the two halves under the same type name. See
// internal/framework/services/compute/resource_datacenter_list.go.
func TestAccDataCenterQuery(t *testing.T) {
	const (
		datacenterName = "tf-test-datacenter-query"
		datacenterAddr = constant.DatacenterResource + ".test_datacenter"
		otherLocation  = "de/txl"
	)

	resource.Test(t, resource.TestCase{
		PreCheck: func() { testAccPreCheck(t) },
		// `terraform query` and list blocks were introduced in Terraform 1.14.
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_14_0),
		},
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactoriesInternal(t, &testAccProvider),
		CheckDestroy:             testAccCheckDatacenterDestroyCheck,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
resource %[1]q "test_datacenter" {
  name        = %[2]q
  location    = "us/las"
  description = "Datacenter for the list resource acceptance test"
}`, constant.DatacenterResource, datacenterName),
			},
			// List without filters: the datacenter must show up with its identity.
			{
				Query: true,
				Config: fmt.Sprintf(`list %[1]q "test_datacenter" {
  provider = ionoscloud
}`, constant.DatacenterResource),
				QueryResultChecks: []querycheck.QueryResultCheck{
					querycheck.ExpectIdentity(datacenterAddr, map[string]knownvalue.Check{
						"id":       knownvalue.NotNull(),
						"location": knownvalue.StringExact("us/las"),
					}),
				},
			},
			// Filter by name and location: the unique name guarantees exactly one result.
			{
				Query: true,
				Config: fmt.Sprintf(`list %[1]q "test_datacenter" {
  provider = ionoscloud
  config {
    filters = [
      { field_name = "name",     field_value = %[2]q },
      { field_name = "location", field_value = "us/las" },
    ]
  }
}`, constant.DatacenterResource, datacenterName),
				QueryResultChecks: []querycheck.QueryResultCheck{
					querycheck.ExpectLength(datacenterAddr, 1),
				},
			},
			// Same name, different location: proves the location filter is evaluated.
			{
				Query: true,
				Config: fmt.Sprintf(`list %[1]q "test_datacenter" {
  provider = ionoscloud
  config {
    filters = [
      { field_name = "name",     field_value = %[2]q },
      { field_name = "location", field_value = %[3]q },
    ]
  }
}`, constant.DatacenterResource, datacenterName, otherLocation),
				QueryResultChecks: []querycheck.QueryResultCheck{
					querycheck.ExpectLength(datacenterAddr, 0),
				},
			},
			// Import through the resource identity that the list results carry. This kind
			// already checks that the import succeeds, that the plan it leaves behind is a
			// no-op and that the planned identity matches the one in state; ImportStateVerify
			// cannot be combined with it, only ImportCommandWithID reads that field.
			{
				ResourceName:    datacenterAddr,
				ImportState:     true,
				ImportStateKind: resource.ImportBlockWithResourceIdentity,
			},
		},
	})
}

func testAccCheckDatacenterDestroyCheck(s *terraform.State) error {
	ctx, cancel := context.WithTimeout(context.Background(), *resourceDefaultTimeouts.Default)

	if cancel != nil {
		defer cancel()
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != constant.DatacenterResource {
			continue
		}

		client, err := testAccProvider.Meta().(bundleclient.SdkBundle).NewCloudAPIClient(ctx, rs.Primary.Attributes["location"])
		if err != nil {
			return err
		}
		_, apiResponse, err := client.DataCentersApi.DatacentersFindById(ctx, rs.Primary.ID).Execute()
		logApiRequestTime(apiResponse)

		if err != nil {
			if !httpNotFound(apiResponse) {
				return fmt.Errorf("an error occurred while checking the destruction of datacenter %s: %w", rs.Primary.ID, err)
			}
		} else {
			return fmt.Errorf("datacenter %s still exists", rs.Primary.ID)
		}
	}

	return nil
}

func testAccCheckDatacenterExists(n string, datacenter *ionoscloud.Datacenter) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no Record ID is set")
		}

		ctx, cancel := context.WithTimeout(context.Background(), *resourceDefaultTimeouts.Default)

		if cancel != nil {
			defer cancel()
		}

		client, err := testAccProvider.Meta().(bundleclient.SdkBundle).NewCloudAPIClient(ctx, rs.Primary.Attributes["location"])
		if err != nil {
			return err
		}
		foundDC, apiResponse, err := client.DataCentersApi.DatacentersFindById(ctx, rs.Primary.ID).Execute()
		logApiRequestTime(apiResponse)

		if err != nil {
			return fmt.Errorf("error occurred while fetching DC: %s", rs.Primary.ID)
		}
		if *foundDC.Id != rs.Primary.ID {
			return fmt.Errorf("record not found")
		}
		datacenter = &foundDC

		return nil
	}
}

const testAccCheckDatacenterConfigUpdate = `
resource ` + constant.DatacenterResource + ` ` + constant.DatacenterTestResource + ` {
	name       =  "` + constant.UpdatedResources + `"
	location = "us/las"
	description = "Test Datacenter Description Updated"
	sec_auth_protection = false
}`

const testAccDataSourceDatacenterMatchID = testAccCheckDatacenterConfigBasic + `
data ` + constant.DatacenterResource + ` ` + constant.DatacenterDataSourceById + ` {
  id			= ` + constant.DatacenterResource + `.` + constant.DatacenterTestResource + `.id
}`

const testAccDataSourceDatacenterMatchName = testAccCheckDatacenterConfigBasic + `
data ` + constant.DatacenterResource + ` ` + constant.DatacenterDataSourceByName + ` {
    name = ` + constant.DatacenterResource + `.` + constant.DatacenterTestResource + `.name
}`

const testAccDataSourceDatacenterMatching = testAccCheckDatacenterConfigBasic + `
data ` + constant.DatacenterResource + ` ` + constant.DatacenterDataSourceMatching + ` {
    name = ` + constant.DatacenterResource + `.` + constant.DatacenterTestResource + `.name
    location = ` + constant.DatacenterResource + `.` + constant.DatacenterTestResource + `.location
}`

const testAccDataSourceDatacenterMultipleResultsError = testAccCheckDatacenterConfigBasic + `
resource ` + constant.DatacenterResource + ` ` + constant.DatacenterTestResource + `_multiple_results {
	name       = "` + constant.DatacenterTestResource + `"
	location = "us/las"
	description = "Test Datacenter Description Updated"
	sec_auth_protection = false
}

data ` + constant.DatacenterResource + ` ` + constant.DatacenterDataSourceMatching + ` {
    name = ` + constant.DatacenterResource + `.` + constant.DatacenterTestResource + `.name
    location = ` + constant.DatacenterResource + `.` + constant.DatacenterTestResource + `.location
}`

const testAccDataSourceDatacenterWrongLocationError = testAccCheckDatacenterConfigBasic + `
data ` + constant.DatacenterResource + ` ` + constant.DatacenterDataSourceMatching + ` {
    name = ` + constant.DatacenterResource + `.` + constant.DatacenterTestResource + `.name
    location =  "wrong_location"
}`

const testAccDataSourceDatacenterWrongNameAndLocationError = testAccCheckDatacenterConfigBasic + `
data ` + constant.DatacenterResource + ` ` + constant.DatacenterDataSourceMatching + ` {
    name =  "wrong_name"
    location =  "wrong_location"
}`

const testAccDataSourceDatacenterNoFilterError = testAccCheckDatacenterConfigBasic + `
data ` + constant.DatacenterResource + ` ` + constant.DatacenterDataSourceMatching + ` {
}`

const testAccDataSourceDatacenterWrongIdError = testAccCheckDatacenterConfigBasic + `
data ` + constant.DatacenterResource + ` ` + constant.DatacenterDataSourceById + ` {
    id = "00000000-0000-0000-0000-000000000000"
}`

const testAccDataSourceDatacenterIdNameMismatchError = testAccCheckDatacenterConfigBasic + `
data ` + constant.DatacenterResource + ` ` + constant.DatacenterDataSourceMatching + ` {
    id   = ` + constant.DatacenterResource + `.` + constant.DatacenterTestResource + `.id
    name = "wrong_name"
}`

const testAccDataSourceDatacenterIdLocationMismatchError = testAccCheckDatacenterConfigBasic + `
data ` + constant.DatacenterResource + ` ` + constant.DatacenterDataSourceMatching + ` {
    id       = ` + constant.DatacenterResource + `.` + constant.DatacenterTestResource + `.id
    location = "wrong_location"
}`
