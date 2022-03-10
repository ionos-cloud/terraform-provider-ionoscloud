//go:build compute || all || s3key

package ionoscloud

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils"
	"testing"
)

func TestAccS3KeyImportBasic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccChecks3KeyDestroyCheck,
		Steps: []resource.TestStep{
			{
				Config: testAccImportS3KeyConfigBasic,
			},
			{
				ResourceName:            S3KeyResource + "." + S3KeyTestResource,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateIdFunc:       testAccS3KeyImportStateID,
				ImportStateVerifyIgnore: []string{},
			},
		},
	})
}

func testAccS3KeyImportStateID(s *terraform.State) (string, error) {
	var importID string = ""

	for _, rs := range s.RootModule().Resources {
		if rs.Type != S3KeyResource {
			continue
		}

		importID = fmt.Sprintf("%s/%s", rs.Primary.Attributes["user_id"], rs.Primary.ID)
	}

	return importID, nil
}

var testAccImportS3KeyConfigBasic = `
resource ` + UserResource + ` "example" {
  first_name = "terraform"
  last_name = "test"
  email = "` + utils.GenerateEmail() + `"
  password = "abc123-321CBA"
  administrator = false
  force_sec_auth= false
}

resource ` + S3KeyResource + ` ` + S3KeyTestResource + ` {
  user_id    = ` + UserResource + `.example.id
  active     = false
}`
