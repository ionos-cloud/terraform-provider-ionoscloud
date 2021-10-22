package ionoscloud

import (
	"context"
	"fmt"
	ionoscloud "github.com/ionos-cloud/sdk-go/v6"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccGroup_Basic(t *testing.T) {
	var group ionoscloud.Group

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckGroupDestroyCheck,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckGroupConfigBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGroupExists(GroupResource+"."+GroupTestResource, &group),
					resource.TestCheckResourceAttr(GroupResource+"."+GroupTestResource, "name", GroupTestResource),
					resource.TestCheckResourceAttr(GroupResource+"."+GroupTestResource, "create_datacenter", "true"),
					resource.TestCheckResourceAttr(GroupResource+"."+GroupTestResource, "create_snapshot", "true"),
					resource.TestCheckResourceAttr(GroupResource+"."+GroupTestResource, "reserve_ip", "true"),
					resource.TestCheckResourceAttr(GroupResource+"."+GroupTestResource, "access_activity_log", "true"),
					resource.TestCheckResourceAttr(GroupResource+"."+GroupTestResource, "create_pcc", "true"),
					resource.TestCheckResourceAttr(GroupResource+"."+GroupTestResource, "s3_privilege", "true"),
					resource.TestCheckResourceAttr(GroupResource+"."+GroupTestResource, "create_backup_unit", "true"),
					resource.TestCheckResourceAttr(GroupResource+"."+GroupTestResource, "create_internet_access", "true"),
					resource.TestCheckResourceAttr(GroupResource+"."+GroupTestResource, "create_k8s_cluster", "true"),
					testNotEmptySlice(GroupResource, "users")),
			},
			{
				Config: testAccCheckGroupConfigUpdate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGroupExists(GroupResource+"."+GroupTestResource, &group),
					resource.TestCheckResourceAttr(GroupResource+"."+GroupTestResource, "name", UpdatedResources),
					resource.TestCheckResourceAttr(GroupResource+"."+GroupTestResource, "create_datacenter", "false"),
					resource.TestCheckResourceAttr(GroupResource+"."+GroupTestResource, "create_snapshot", "false"),
					resource.TestCheckResourceAttr(GroupResource+"."+GroupTestResource, "reserve_ip", "false"),
					resource.TestCheckResourceAttr(GroupResource+"."+GroupTestResource, "access_activity_log", "false"),
					resource.TestCheckResourceAttr(GroupResource+"."+GroupTestResource, "create_pcc", "false"),
					resource.TestCheckResourceAttr(GroupResource+"."+GroupTestResource, "s3_privilege", "false"),
					resource.TestCheckResourceAttr(GroupResource+"."+GroupTestResource, "create_backup_unit", "false"),
					resource.TestCheckResourceAttr(GroupResource+"."+GroupTestResource, "create_internet_access", "false"),
					resource.TestCheckResourceAttr(GroupResource+"."+GroupTestResource, "create_k8s_cluster", "false"),
					testNotEmptySlice(GroupResource, "users")),
			},
		},
	})
}

func testAccCheckGroupDestroyCheck(s *terraform.State) error {
	client := testAccProvider.Meta().(*ionoscloud.APIClient)
	ctx, cancel := context.WithTimeout(context.Background(), *resourceDefaultTimeouts.Default)

	if cancel != nil {
		defer cancel()
	}
	for _, rs := range s.RootModule().Resources {

		if rs.Type != GroupResource {
			continue
		}
		_, apiResponse, err := client.UserManagementApi.UmGroupsFindById(ctx, rs.Primary.ID).Execute()

		if err != nil {
			if apiResponse == nil || apiResponse.Response != nil && apiResponse.StatusCode != 404 {
				return fmt.Errorf("an error occurred while checking the destruction of group %s: %s", rs.Primary.ID, err)
			}
		} else {
			return fmt.Errorf("group %s still exists", rs.Primary.ID)
		}

	}

	return nil
}

func testAccCheckGroupExists(n string, group *ionoscloud.Group) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*ionoscloud.APIClient)

		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("testAccCheckGroupExists: Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no Record ID is set")
		}

		ctx, cancel := context.WithTimeout(context.Background(), *resourceDefaultTimeouts.Default)

		if cancel != nil {
			defer cancel()
		}

		foundgroup, _, err := client.UserManagementApi.UmGroupsFindById(ctx, rs.Primary.ID).Execute()

		if err != nil {
			return fmt.Errorf("error occured while fetching Group: %s", rs.Primary.ID)
		}
		if *foundgroup.Id != rs.Primary.ID {
			return fmt.Errorf("record not found")
		}

		group = &foundgroup

		return nil
	}
}

var testAccCheckGroupConfigBasic = `
resource ` + UserResource + ` ` + UserTestResource + ` {
  first_name = "user"
  last_name = "test"
  email = "` + email + `"
  password = "abc123-321CBA"
  administrator = false
  force_sec_auth= false
  active = false
}

resource ` + GroupResource + ` ` + GroupTestResource + ` {
  name = "` + GroupTestResource + `"
  create_datacenter = true
  create_snapshot = true
  reserve_ip = true
  access_activity_log = true
  create_pcc = true
  s3_privilege = true
  create_backup_unit = true
  create_internet_access = true
  create_k8s_cluster = true
  user_id = ` + UserResource + `.` + UserTestResource + `.id
}
`

var testAccCheckGroupConfigUpdate = `
resource ` + UserResource + ` resource_user_updated {
  first_name = "updated"
  last_name = "test"
  email = "updated` + email + `"
  password = "abc123-321CBA"
  administrator = false
  force_sec_auth= false
  active = false
}

resource ` + GroupResource + ` ` + GroupTestResource + ` {
  name = "` + UpdatedResources + `"
  create_datacenter = false
  create_snapshot = false
  reserve_ip = false
  access_activity_log = false
  create_pcc = false
  s3_privilege = false
  create_backup_unit = false
  create_internet_access = false
  create_k8s_cluster = false
  user_id = ` + UserResource + `.resource_user_updated.id
}
`
