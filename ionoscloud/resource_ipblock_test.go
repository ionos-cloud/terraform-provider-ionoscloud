//go:build compute || all || ipblock

package ionoscloud

import (
	"context"
	"fmt"
	"regexp"
	"testing"

	ionoscloud "github.com/ionos-cloud/sdk-go/v6"

	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/bundleclient"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils/constant"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/querycheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

const fullIpBlockResourceName = constant.IpBlockResource + "." + constant.IpBlockTestResource

var dataSourceIpBlockNameByID = fmt.Sprintf("%s.%s.%s", constant.DataSource, constant.IpBlockResource, constant.IpBlockDataSourceById)
var dataSourceIpBlockNameMatching = fmt.Sprintf("%s.%s.%s", constant.DataSource, constant.IpBlockResource, constant.IpBlockDataSourceMatching)
var dataSourceIpBlockNameMatchName = fmt.Sprintf("%s.%s.%s", constant.DataSource, constant.IpBlockResource, constant.IpBlockDataSourceByName)

const location = "us/las"

func TestAccIPBlockBasic(t *testing.T) {
	var ipblock ionoscloud.IpBlock

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactoriesInternal(t, &testAccProvider),
		CheckDestroy:             testAccCheckIPBlockDestroyCheck,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckIPBlockConfigBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIPBlockExists(fullIpBlockResourceName, &ipblock),
					resource.TestCheckResourceAttr(fullIpBlockResourceName, "location", location),
					resource.TestCheckResourceAttr(fullIpBlockResourceName, "name", constant.IpBlockTestResource),
					resource.TestCheckResourceAttr(fullIpBlockResourceName, "size", "1"),
				),
			}, {
				Config: testAccDataSourceIpBlockMatchID,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceIpBlockNameByID, "name", fullIpBlockResourceName, "name"),
					resource.TestCheckResourceAttrPair(dataSourceIpBlockNameByID, "location", fullIpBlockResourceName, "location"),
					resource.TestCheckResourceAttrPair(dataSourceIpBlockNameByID, "size", fullIpBlockResourceName, "size"),
					resource.TestCheckResourceAttrPair(dataSourceIpBlockNameByID, "ips", fullIpBlockResourceName, "ips"),
					resource.TestCheckResourceAttrPair(dataSourceIpBlockNameByID, "ip_consumers", fullIpBlockResourceName, "ip_consumers"),
				),
			},
			{
				Config: testAccDataSourceIpBlockMatching,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceIpBlockNameMatching, "name", fullIpBlockResourceName, "name"),
					resource.TestCheckResourceAttrPair(dataSourceIpBlockNameMatching, "location", fullIpBlockResourceName, "location"),
					resource.TestCheckResourceAttrPair(dataSourceIpBlockNameMatching, "size", fullIpBlockResourceName, "size"),
					resource.TestCheckResourceAttrPair(dataSourceIpBlockNameMatching, "ips", fullIpBlockResourceName, "ips"),
					resource.TestCheckResourceAttrPair(dataSourceIpBlockNameMatching, "ip_consumers", fullIpBlockResourceName, "ip_consumers"),
				),
			},
			{
				Config: testAccDataSourceIpBlockMatchName,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceIpBlockNameMatchName, "name", fullIpBlockResourceName, "name"),
					resource.TestCheckResourceAttrPair(dataSourceIpBlockNameMatchName, "location", fullIpBlockResourceName, "location"),
					resource.TestCheckResourceAttrPair(dataSourceIpBlockNameMatchName, "size", fullIpBlockResourceName, "size"),
					resource.TestCheckResourceAttrPair(dataSourceIpBlockNameMatchName, "ips", fullIpBlockResourceName, "ips"),
					resource.TestCheckResourceAttrPair(dataSourceIpBlockNameMatchName, "ip_consumers", fullIpBlockResourceName, "ip_consumers"),
				),
			},
			{
				Config:      testAccDataSourceIpBlockMultipleResultsError,
				ExpectError: regexp.MustCompile(`more than one ip block found with the specified criteria`),
			},
			{
				Config:      testAccDataSourceIpBlockNameError,
				ExpectError: regexp.MustCompile(`no ip block found with the specified criteria`),
			},
			{
				Config:      testAccDataSourceIpBlockMatchNameLocationError,
				ExpectError: regexp.MustCompile(`no ip block found with the specified criteria`),
			},
			{
				Config:      testAccDataSourceIpBlockLocationError,
				ExpectError: regexp.MustCompile(`no ip block found with the specified criteria`),
			},
			{
				Config:      testIpBlockGoodIDLocationError,
				ExpectError: regexp.MustCompile(`location of ip block`),
			},
			{
				Config:      testAccDataSourceIpBlockNoFilterError,
				ExpectError: regexp.MustCompile(`either id, location or name must be set`),
			},
			{
				Config:      testAccDataSourceIpBlockWrongIdError,
				ExpectError: regexp.MustCompile(`error getting ip block with id`),
			},
			{
				Config:      testAccDataSourceIpBlockGoodIdNameError,
				ExpectError: regexp.MustCompile(`name of ip block \(UUID=.+, name=.+\) does not match expected name`),
			},
			// name is the only non-ForceNew attribute, so this step is the only one that
			// actually enters resourceIPBlockUpdate - and therefore the only coverage of
			// the identity write on the update path. The size change below is a
			// destroy-and-create, which goes through Create and Read instead.
			{
				Config: testAccCheckIPBlockConfigUpdateNameOnly,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIPBlockExists(fullIpBlockResourceName, &ipblock),
					testAccCheckIPBlockAttributes(fullIpBlockResourceName, location),
					resource.TestCheckResourceAttr(fullIpBlockResourceName, "name", constant.UpdatedResources),
					resource.TestCheckResourceAttr(fullIpBlockResourceName, "size", "1"),
				),
			},
			{
				Config: testAccCheckIPBlockConfigUpdate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIPBlockExists(fullIpBlockResourceName, &ipblock),
					testAccCheckIPBlockAttributes(fullIpBlockResourceName, location),
					resource.TestCheckResourceAttr(fullIpBlockResourceName, "name", constant.UpdatedResources),
					resource.TestCheckResourceAttr(fullIpBlockResourceName, "size", "2"),
				),
			},
		},
	})
}

// TestAccIPBlockQuery exercises the ionoscloud_ipblock list resource and the resource
// identity that listing depends on.
//
// The list resource is served by the plugin-framework half of the provider even though
// the ipblock resource itself is implemented with SDKv2, so this also covers the mux
// serving the two halves under the same type name. See resource_ipblock_list.go in this
// package.
func TestAccIPBlockQuery(t *testing.T) {
	const (
		ipBlockName   = "tf-test-ipblock-query"
		ipBlockAddr   = constant.IpBlockResource + ".test_ipblock"
		otherLocation = "de/fra"
	)

	resource.Test(t, resource.TestCase{
		PreCheck: func() { testAccPreCheck(t) },
		// `terraform query` and list blocks were introduced in Terraform 1.14.
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_14_0),
		},
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactoriesInternal(t, &testAccProvider),
		CheckDestroy:             testAccCheckIPBlockDestroyCheck,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
resource %[1]q "test_ipblock" {
  location = %[2]q
  size     = 1
  name     = %[3]q
}`, constant.IpBlockResource, location, ipBlockName),
			},
			// List without filters: the ip block must show up with its identity.
			{
				Query: true,
				Config: fmt.Sprintf(`list %[1]q "test_ipblock" {
  provider = ionoscloud
}`, constant.IpBlockResource),
				QueryResultChecks: []querycheck.QueryResultCheck{
					querycheck.ExpectIdentity(ipBlockAddr, map[string]knownvalue.Check{
						"id":       knownvalue.NotNull(),
						"location": knownvalue.StringExact(location),
					}),
				},
			},
			// Filter by name and location: the unique name guarantees exactly one result.
			{
				Query: true,
				Config: fmt.Sprintf(`list %[1]q "test_ipblock" {
  provider = ionoscloud
  config {
    filters = [
      { field_name = "name",     field_value = %[2]q },
      { field_name = "location", field_value = %[3]q },
    ]
  }
}`, constant.IpBlockResource, ipBlockName, location),
				QueryResultChecks: []querycheck.QueryResultCheck{
					querycheck.ExpectLength(ipBlockAddr, 1),
				},
			},
			// Same name, different location: proves the location filter is evaluated.
			{
				Query: true,
				Config: fmt.Sprintf(`list %[1]q "test_ipblock" {
  provider = ionoscloud
  config {
    filters = [
      { field_name = "name",     field_value = %[2]q },
      { field_name = "location", field_value = %[3]q },
    ]
  }
}`, constant.IpBlockResource, ipBlockName, otherLocation),
				QueryResultChecks: []querycheck.QueryResultCheck{
					querycheck.ExpectLength(ipBlockAddr, 0),
				},
			},
			// Import through the resource identity that the list results carry. This kind
			// already checks that the import succeeds, that the plan it leaves behind is a
			// no-op and that the planned identity matches the one in state; ImportStateVerify
			// cannot be combined with it, only ImportCommandWithID reads that field.
			{
				ResourceName:    ipBlockAddr,
				ImportState:     true,
				ImportStateKind: resource.ImportBlockWithResourceIdentity,
			},
		},
	})
}

func testAccCheckIPBlockDestroyCheck(s *terraform.State) error {
	ctx, cancel := context.WithTimeout(context.Background(), *resourceDefaultTimeouts.Delete)
	if cancel != nil {
		defer cancel()
	}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != constant.IpBlockResource {
			continue
		}

		client, err := testAccProvider.Meta().(bundleclient.SdkBundle).NewCloudAPIClient(ctx, rs.Primary.Attributes["location"])
		if err != nil {
			return err
		}
		_, apiResponse, err := client.IPBlocksApi.IpblocksFindById(ctx, rs.Primary.ID).Execute()
		logApiRequestTime(apiResponse)

		if err != nil {
			if !httpNotFound(apiResponse) {
				return fmt.Errorf("an error occurred while checking deletion of IPBlock %s %s", rs.Primary.ID, err)
			}
		} else {
			return fmt.Errorf("IPBlock still exists %s %s", rs.Primary.ID, err)
		}
	}

	return nil
}

func testAccCheckIPBlockAttributes(n string, location string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("testAccCheckLanAttributes: Not found: %s", n)
		}
		if rs.Primary.Attributes["location"] != location {
			return fmt.Errorf("bad name: %s", rs.Primary.Attributes["location"])
		}

		return nil
	}
}

func testAccCheckIPBlockExists(n string, ipblock *ionoscloud.IpBlock) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("testAccCheckIPBlockExists: Not found: %s", n)
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
		foundIP, apiResponse, err := client.IPBlocksApi.IpblocksFindById(ctx, rs.Primary.ID).Execute()
		logApiRequestTime(apiResponse)

		if err != nil {
			return fmt.Errorf("error occurred while fetching IP Block: %s", rs.Primary.ID)
		}
		if *foundIP.Id != rs.Primary.ID {
			return fmt.Errorf("record not found")
		}

		ipblock = &foundIP

		return nil
	}
}

const testAccCheckIPBlockConfigBasic = `
resource ` + constant.IpBlockResource + ` ` + constant.IpBlockTestResource + ` {
  location = "` + location + `"
  size = 1
  name = "` + constant.IpBlockTestResource + `"
}`

const testAccCheckIPBlockConfigUpdateNameOnly = `
resource ` + constant.IpBlockResource + ` ` + constant.IpBlockTestResource + ` {
  location = "` + location + `"
  size = 1
  name = "` + constant.UpdatedResources + `"
}`

const testAccCheckIPBlockConfigUpdate = `
resource ` + constant.IpBlockResource + ` ` + constant.IpBlockTestResource + `{
  location = "` + location + `"
  size = 2
  name = "` + constant.UpdatedResources + `"
}`

const testAccDataSourceIpBlockMatchID = testAccCheckIPBlockConfigBasic + `
data ` + constant.IpBlockResource + `  ` + constant.IpBlockDataSourceById + ` {
	id = ` + fullIpBlockResourceName + `.id 
}
`

const testAccDataSourceIpBlockMatching = testAccCheckIPBlockConfigBasic + `
data ` + constant.IpBlockResource + ` ` + constant.IpBlockDataSourceMatching + ` { 
	name = ` + fullIpBlockResourceName + `.name
	location = ` + fullIpBlockResourceName + `.location 
}`

const testAccDataSourceIpBlockMatchName = testAccCheckIPBlockConfigBasic + `
data ` + constant.IpBlockResource + ` ` + constant.IpBlockDataSourceByName + ` { 
	name = ` + fullIpBlockResourceName + `.name
}`

const testAccDataSourceIpBlockMultipleResultsError = testAccCheckIPBlockConfigBasic + `
resource ` + constant.IpBlockResource + ` ` + constant.IpBlockTestResource + `_same_name{
  location = "` + location + `"
  size = 2
  name = ` + fullIpBlockResourceName + `.name
}

data ` + constant.IpBlockResource + ` ` + constant.IpBlockDataSourceByName + ` { 
	name = ` + fullIpBlockResourceName + `.name
}`

const testAccDataSourceIpBlockNameError = testAccCheckIPBlockConfigBasic + `
data ` + constant.IpBlockResource + ` ` + constant.IpBlockDataSourceByName + ` { 
	name = ` + fullIpBlockResourceName + `.size
}`
const testAccDataSourceIpBlockMatchNameLocationError = testAccCheckIPBlockConfigBasic + `
data ` + constant.IpBlockResource + ` ` + constant.IpBlockDataSourceByName + ` { 
	name = ` + fullIpBlockResourceName + `.name
	location = "none"
}`
const testAccDataSourceIpBlockLocationError = testAccCheckIPBlockConfigBasic + `
data ` + constant.IpBlockResource + ` ` + constant.IpBlockDataSourceByName + ` {
	location = "none"
}`

const testIpBlockGoodIDLocationError = testAccCheckIPBlockConfigBasic + `
data ` + constant.IpBlockResource + ` ` + constant.IpBlockDataSourceByName + ` {
    id = ` + fullIpBlockResourceName + `.id
	location = "none"
}`

const testAccDataSourceIpBlockNoFilterError = testAccCheckIPBlockConfigBasic + `
data ` + constant.IpBlockResource + ` ` + constant.IpBlockDataSourceByName + ` {
}`

const testAccDataSourceIpBlockWrongIdError = testAccCheckIPBlockConfigBasic + `
data ` + constant.IpBlockResource + ` ` + constant.IpBlockDataSourceById + ` {
    id = "00000000-0000-0000-0000-000000000000"
}`

const testAccDataSourceIpBlockGoodIdNameError = testAccCheckIPBlockConfigBasic + `
data ` + constant.IpBlockResource + ` ` + constant.IpBlockDataSourceByName + ` {
    id   = ` + fullIpBlockResourceName + `.id
    name = "wrong_name"
}`
