//go:build compute || all || location

package ionoscloud

import (
	"regexp"
	"testing"

	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils/constant"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

const locationTestName = constant.DataSource + "." + constant.LocationResource + "." + constant.LocationTestResource

func TestAccDataSourceLocationBasic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactoriesInternal(t, &testAccProvider),
		Steps: []resource.TestStep{
			{

				Config: testAccDataSourceLocationBasic,
				Check: resource.ComposeTestCheckFunc(resource.TestCheckResourceAttr(locationTestName, "id", "de/fra"),
					resource.TestCheckResourceAttr(locationTestName, "name", "frankfurt"),
				),
			},
			{
				Config:      testAccDataSourceLocationWrongNameError,
				ExpectError: regexp.MustCompile("no location found with the specified criteria"),
			},
			{
				Config:      testAccDataSourceLocationWrongFeature,
				ExpectError: regexp.MustCompile("no location found with the specified criteria"),
			},
			{
				Config:      testAccDataSourceLocationNoFilterError,
				ExpectError: regexp.MustCompile(`either 'name' or 'feature' must be provided`),
			},
			{
				Config: testAccDataSourceLocationByName,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(locationTestName, "id", "de/fra"),
					resource.TestCheckResourceAttr(locationTestName, "name", "frankfurt"),
				),
			},
		},
	})

}

const testAccDataSourceLocationBasic = `
data ` + constant.LocationResource + ` ` + constant.LocationTestResource + ` {
	  name = "frankfurt"
	  feature = "SSD"
}
`
const testAccDataSourceLocationWrongNameError = `
data ` + constant.LocationResource + ` ` + constant.LocationTestResource + ` {
	  name = "wrong_name"
	  feature = "SSD"
}
`
const testAccDataSourceLocationWrongFeature = `
data ` + constant.LocationResource + ` ` + constant.LocationTestResource + ` {
	  name = "frankfurt"
	  feature = "wrong_feature"
}
`

const testAccDataSourceLocationNoFilterError = `
data ` + constant.LocationResource + ` ` + constant.LocationTestResource + ` {
}
`

const testAccDataSourceLocationByName = `
data ` + constant.LocationResource + ` ` + constant.LocationTestResource + ` {
	  name = "frankfurt"
}
`
