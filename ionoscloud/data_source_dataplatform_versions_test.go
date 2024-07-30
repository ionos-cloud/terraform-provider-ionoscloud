//go:build all || dataplatform
// +build all dataplatform

package ionoscloud

import (
	"testing"

	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils/constant"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDataSourceDataplatformVersions(t *testing.T) {
	t.Skip("problem in the go sdk getting versions")
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactoriesInternal(t, &testAccProvider),
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceDataplatformVersions,
				Check:  utils.TestNotEmptySlice(constant.DataSource+"."+constant.DataplatformVersionsDataSource+"."+constant.DataplatformVersionsTestDataSource, "versions.#"),
			},
		},
	})
}

const testAccDataSourceDataplatformVersions = `
data ` + constant.DataplatformVersionsDataSource + ` ` + constant.DataplatformVersionsTestDataSource + ` {
}
`
