//go:build compute || all || location

package ionoscloud

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

const locationTestName = DataSource + "." + LocationResource + "." + LocationTestResource

func TestAccDataSourceLocationBasic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: testAccProviderFactories,
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
		},
	})

}

const testAccDataSourceLocationBasic = `
data ` + LocationResource + ` ` + LocationTestResource + ` {
	  name = "frankfurt"
	  feature = "SSD"
}
`
const testAccDataSourceLocationWrongNameError = `
data ` + LocationResource + ` ` + LocationTestResource + ` {
	  name = "wrong_name"
	  feature = "SSD"
}
`
const testAccDataSourceLocationWrongFeature = `
data ` + LocationResource + ` ` + LocationTestResource + ` {
	  name = "frankfurt"
	  feature = "wrong_feature"
}
`
