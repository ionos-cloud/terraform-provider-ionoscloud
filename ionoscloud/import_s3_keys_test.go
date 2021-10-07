package ionoscloud

import (
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccS3KeyImportBasic(t *testing.T) {
	resourceName := "example"
	email := fmt.Sprintf("terraform-s3-import-acc-tester-%d@mailinator.com", time.Now().Unix())
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccChecks3KeyDestroyCheck,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccChecks3KeyImportConfigBasic, email, resourceName),
			},
			{
				ResourceName:            fmt.Sprintf("ionoscloud_s3_key.%s", resourceName),
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateIdFunc:       testAccS3KeyImportStateID,
				ImportStateVerifyIgnore: []string{},
			},
		},
	})
}

func testAccS3KeyImportStateID(s *terraform.State) (string, error) {
	importID := ""

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "ionoscloud_s3_key" {
			continue
		}

		importID = fmt.Sprintf("%s/%s", rs.Primary.Attributes["user_id"], rs.Primary.ID)
	}

	return importID, nil
}

const testAccChecks3KeyImportConfigBasic = `

resource "ionoscloud_user" "example" {
  first_name = "terraform"
  last_name = "test"
  email = "%s"
  password = "abc123-321CBA"
  administrator = false
  force_sec_auth= false
}

resource "ionoscloud_s3_key" "%s" {
  user_id    = ionoscloud_user.example.id
  active     = true
}`
