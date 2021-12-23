package ionoscloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccDataSourceShareMatchFields(t *testing.T) {

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testacccheckshareconfigBasic,
			},
			{
				Config: testAccDataSourceShareConfigBasic,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(ShareResourceFullName, "id"),
					resource.TestCheckResourceAttrPair(ShareResourceFullName, "id", DataSource+"."+ShareResource+"."+SourceShareName, "id"),
					resource.TestCheckResourceAttrPair(ShareResourceFullName, "edit_privilege",
						DataSource+"."+ShareResource+"."+SourceShareName, "edit_privilege"),
					resource.TestCheckResourceAttrPair(ShareResourceFullName, "share_privilege",
						DataSource+"."+ShareResource+"."+SourceShareName, "share_privilege"),
					resource.TestCheckResourceAttr(DataSource+"."+ShareResource+"."+SourceShareName, "edit_privilege", "true"),
				),
			},
		},
	})
}

var testAccDataSourceShareConfigBasic = testacccheckshareconfigBasic + `
data ` + ShareResource + " " + SourceShareName + `{
  group_id    = "${ionoscloud_group.group.id}"
  resource_id = "${ionoscloud_datacenter.foobar.id}"
  id		  = ` + ShareResourceFullName + `.id
}
`
