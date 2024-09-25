//go:build all || s3management
// +build all s3management

package s3management_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/internal/acctest"
)

func TestAccACcesskeyResource(t *testing.T) {
	description := acctest.GenerateRandomResourceName("description")
	name := "ionoscloud_s3_accesskey.test"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		Steps: []resource.TestStep{
			{
				Config: testAccAccesskeyConfig_description(description),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(name, "description", description),
				),
			},
		},
	})
}

func testAccAccesskeyConfig_basic() string {
	return `
resource "ionoscloud_s3_accesskey" "test" {
}
`
}

func testAccAccesskeyConfig_description(description string) string {
	return fmt.Sprintf(`
resource "ionoscloud_s3_accesskey" "test" {
  description = %[1]q
}
`, description)
}
