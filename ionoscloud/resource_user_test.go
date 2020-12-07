package ionoscloud

import (
	"fmt"
	"testing"

	"math/rand"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/profitbricks/profitbricks-sdk-go/v5"
)

func TestAccUser_Basic(t *testing.T) {
	var user profitbricks.User
	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)
	email := strconv.Itoa(r1.Intn(100000)) + "terraform_test" + strconv.Itoa(r1.Intn(100000)) + "@go.com"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckUserDestroyCheck,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccCheckUserConfig_basic, email),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckUserExists("ionoscloud_user.user", &user),
					testAccCheckUserAttributes("ionoscloud_user.user", "terraform"),
					resource.TestCheckResourceAttr("ionoscloud_user.user", "first_name", "terraform"),
				),
			},
			{
				Config: fmt.Sprintf(testAccCheckUserConfig_update, email),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckUserAttributes("ionoscloud_user.user", "updated"),
					resource.TestCheckResourceAttr("ionoscloud_user.user", "first_name", "updated"),
				),
			},
		},
	})
}

func testAccCheckUserDestroyCheck(s *terraform.State) error {
	client := testAccProvider.Meta().(*profitbricks.Client)
	for _, rs := range s.RootModule().Resources {
		_, err := client.GetUser(rs.Primary.ID)

		if err != nil {
			return fmt.Errorf("user still exists %s %s", rs.Primary.ID, err)
		}
	}

	return nil
}

func testAccCheckUserAttributes(n string, name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("testAccCheckUserAttributes: Not found: %s", n)
		}
		if rs.Primary.Attributes["first_name"] != name {
			return fmt.Errorf("Bad first_name: %s", rs.Primary.Attributes["first_name"])
		}

		return nil
	}
}

func testAccCheckUserExists(n string, user *profitbricks.User) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*profitbricks.Client)
		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("testAccCheckUserExists: Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Record ID is set")
		}

		founduser, err := client.GetUser(rs.Primary.ID)

		if err != nil {
			return fmt.Errorf("Error occured while fetching User: %s", rs.Primary.ID)
		}
		if founduser.ID != rs.Primary.ID {
			return fmt.Errorf("Record not found")
		}

		user = founduser

		return nil
	}
}

const testAccCheckUserConfig_basic = `

resource "ionoscloud_group" "group" {
  name = "terraform user group"
  create_datacenter = true
  create_snapshot = true
  reserve_ip = true
  access_activity_log = false
}

resource "ionoscloud_user" "user" {
  first_name = "terraform"
  last_name = "test"
  email = "%s"
  password = "abc123-321CBA"
  administrator = false
  force_sec_auth= false
}`

const testAccCheckUserConfig_update = `


resource "ionoscloud_user" "user" {
  first_name = "updated"
  last_name = "test"
  email = "%s"
  password = "abc123-321CBA"
  administrator = false
  force_sec_auth= false
}

resource "ionoscloud_group" "group" {
  name = "terraform user group"
  create_datacenter = true
  create_snapshot = true
  reserve_ip = true
  access_activity_log = false
  user_id="${ionoscloud_user.user.id}"
}
`
