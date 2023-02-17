//go:build all || dataplatform
// +build all dataplatform

package ionoscloud

import (
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceDataplatformVersions(t *testing.T) {
	t.Skip("problem in the go sdk getting versions")
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceDataplatformVersions,
				Check:  utils.TestNotEmptySlice(DataSource+"."+DataplatformVersionsDataSource+"."+DataplatformVersionsTestDataSource, "versions.#"),
			},
		},
	})
}

const testAccDataSourceDataplatformVersions = `
data ` + DataplatformVersionsDataSource + ` ` + DataplatformVersionsTestDataSource + ` {
}
`
