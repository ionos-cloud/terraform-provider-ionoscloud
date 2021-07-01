package ionoscloud

import (
	"fmt"
	"math/rand"
	"strconv"
	"time"

	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccUser_ImportBasic(t *testing.T) {
	resourceName := "resource_user"
	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)
	email := strconv.Itoa(r1.Intn(100000)) + "terraform_test" + strconv.Itoa(r1.Intn(100000)) + "@go.com"

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckUserDestroyCheck,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testacccheckuserconfigBasic, email),
			},

			{
				ResourceName:            fmt.Sprintf("ionoscloud_user.%s", resourceName),
				ImportStateIdFunc:       testAccUserImportStateId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"password"},
			},
		},
	})
}

func testAccUserImportStateId(s *terraform.State) (string, error) {
	importID := ""

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "ionoscloud_user" {
			continue
		}

		importID = rs.Primary.Attributes["id"]
	}

	return importID, nil
}
