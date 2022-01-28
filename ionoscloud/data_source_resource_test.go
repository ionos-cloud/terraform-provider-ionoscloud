package ionoscloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

const resourceName = DataSource + "." + ResourceResource + "." + ResourceTestResource

func TestAccResourceBasic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceResourceBasic,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "resource_type", "datacenter"),
				),
			},
		},
	})

}

const testAccDataSourceResourceBasic = `
resource ` + DatacenterResource + ` "foobar" {
  name       = "test_name"
  location = "us/las"
}

data ` + ResourceResource + ` ` + ResourceTestResource + ` {
  resource_type = "datacenter"
  resource_id= ` + DatacenterResource + `.foobar.id
}`
