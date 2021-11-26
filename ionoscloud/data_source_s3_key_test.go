package ionoscloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccDataSourceS3KeyMatchFields(t *testing.T) {

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceS3KeyConfigBasic,
			},
			{
				Config: testAccDataSourceS3KeyMatchId,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(S3KeyResource+"."+S3KeyTestResource, "id"),
					resource.TestCheckResourceAttrPair(S3KeyResource+"."+S3KeyTestResource, "id", DataSource+"."+S3KeyResource+"."+S3KeyDataSourceById, "id"),
					resource.TestCheckResourceAttrPair(S3KeyResource+"."+S3KeyTestResource, "secret", DataSource+"."+S3KeyResource+"."+S3KeyDataSourceById, "secret"),
					resource.TestCheckResourceAttrPair(S3KeyResource+"."+S3KeyTestResource, "active", DataSource+"."+S3KeyResource+"."+S3KeyDataSourceById, "active"),
				),
			},
		},
	})
}

var testAccDataSourceS3KeyConfigBasic = `
resource ` + UserResource + ` "example" {
  first_name = "terraform"
  last_name = "test"
  email = "` + GenerateEmail() + `"
  password = "abc123-321CBA"
  administrator = false
  force_sec_auth= false
}

resource ` + S3KeyResource + ` ` + S3KeyTestResource + ` {
  user_id    = ` + UserResource + `.example.id
  active     = true
}`

var testAccDataSourceS3KeyMatchId = testAccDataSourceS3KeyConfigBasic + `
data ` + S3KeyResource + ` ` + S3KeyDataSourceById + ` {
 user_id    	= ` + UserResource + `.example.id
 id			= ` + S3KeyResource + `.` + S3KeyTestResource + `.id
}
`
