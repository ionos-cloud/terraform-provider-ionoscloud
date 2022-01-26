//go:build compute || all || user

package ionoscloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccDataSourceUser(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceUserConfigBasic,
			},
			{
				Config: testAccDataSourceUserMatchId,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(DataSource+"."+UserResource+"."+UserDataSourceById, "first_name", UserResource+"."+UserTestResource, "first_name"),
					resource.TestCheckResourceAttrPair(DataSource+"."+UserResource+"."+UserDataSourceById, "last_name", UserResource+"."+UserTestResource, "last_name"),
					resource.TestCheckResourceAttrPair(DataSource+"."+UserResource+"."+UserDataSourceById, "email", UserResource+"."+UserTestResource, "email"),
					resource.TestCheckResourceAttrPair(DataSource+"."+UserResource+"."+UserDataSourceById, "administrator", UserResource+"."+UserTestResource, "administrator"),
					resource.TestCheckResourceAttrPair(DataSource+"."+UserResource+"."+UserDataSourceById, "force_sec_auth", UserResource+"."+UserTestResource, "force_sec_auth"),
					resource.TestCheckResourceAttrPair(DataSource+"."+UserResource+"."+UserDataSourceById, "sec_auth_active", UserResource+"."+UserTestResource, "sec_auth_active"),
					resource.TestCheckResourceAttrPair(DataSource+"."+UserResource+"."+UserDataSourceById, "s3_canonical_user_id", UserResource+"."+UserTestResource, "s3_canonical_user_id"),
					resource.TestCheckResourceAttrPair(DataSource+"."+UserResource+"."+UserDataSourceById, "active", UserResource+"."+UserTestResource, "active"),
				),
			},
			{
				Config: testAccDataSourceUserMatchName,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(DataSource+"."+UserResource+"."+UserDataSourceByName, "first_name", UserResource+"."+UserTestResource, "first_name"),
					resource.TestCheckResourceAttrPair(DataSource+"."+UserResource+"."+UserDataSourceByName, "last_name", UserResource+"."+UserTestResource, "last_name"),
					resource.TestCheckResourceAttrPair(DataSource+"."+UserResource+"."+UserDataSourceByName, "email", UserResource+"."+UserTestResource, "email"),
					resource.TestCheckResourceAttrPair(DataSource+"."+UserResource+"."+UserDataSourceByName, "administrator", UserResource+"."+UserTestResource, "administrator"),
					resource.TestCheckResourceAttrPair(DataSource+"."+UserResource+"."+UserDataSourceByName, "force_sec_auth", UserResource+"."+UserTestResource, "force_sec_auth"),
					resource.TestCheckResourceAttrPair(DataSource+"."+UserResource+"."+UserDataSourceByName, "sec_auth_active", UserResource+"."+UserTestResource, "sec_auth_active"),
					resource.TestCheckResourceAttrPair(DataSource+"."+UserResource+"."+UserDataSourceByName, "s3_canonical_user_id", UserResource+"."+UserTestResource, "s3_canonical_user_id"),
					resource.TestCheckResourceAttrPair(DataSource+"."+UserResource+"."+UserDataSourceByName, "active", UserResource+"."+UserTestResource, "active"),
				),
			},
		},
	})
}

var testAccDataSourceUserConfigBasic = `
resource ` + UserResource + ` ` + UserTestResource + ` {
  first_name = "` + UserTestResource + `"
  last_name = "` + UserTestResource + `"
  email = "` + GenerateEmail() + `"
  password = "abc123-321CBA"
  administrator = true
  force_sec_auth= true
  active  = true
}`

var testAccDataSourceUserMatchId = testAccDataSourceUserConfigBasic + `
data ` + UserResource + ` ` + UserDataSourceById + ` {
  id			= ` + UserResource + `.` + UserTestResource + `.id
}
`

var testAccDataSourceUserMatchName = testAccDataSourceUserConfigBasic + `
data ` + UserResource + ` ` + UserDataSourceByName + ` {
  email			= ` + UserResource + `.` + UserTestResource + `.email
}
`
