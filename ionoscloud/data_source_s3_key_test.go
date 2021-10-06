package ionoscloud

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
	"time"
)

var email = fmt.Sprintf("terraform_test-%d@mailinator.com", time.Now().Unix())

func TestAccDataSourceS3Key_matchFields_expectSuccess(t *testing.T) {

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceS3KeyCreateResource,
			},
			{
				Config: testAccDataSourceS3KeyMatchId,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("ionoscloud_s3_key.example_key", "id"),
					resource.TestCheckResourceAttrPair("ionoscloud_s3_key.example_key", "id", "data.ionoscloud_s3_key.example_key_data", "id"),
					resource.TestCheckResourceAttrPair("ionoscloud_s3_key.example_key", "secret", "data.ionoscloud_s3_key.example_key_data", "secret"),
					resource.TestCheckResourceAttrPair("ionoscloud_s3_key.example_key", "active", "data.ionoscloud_s3_key.example_key_data", "active"),
				),
			},
		},
	})
}

var testAccDataSourceS3KeyCreateResource = `
resource "ionoscloud_user" "example" {
  first_name 	= "terraform"
  last_name  	= "test"
  email 	 	= "` + email + `"
  password   	= "abc123-321CBA"
  administrator = false
  force_sec_auth= false
}

resource "ionoscloud_s3_key" "example_key" {
  user_id    = ionoscloud_user.example.id
  active     = false
}`

var testAccDataSourceS3KeyMatchId = `
resource "ionoscloud_user" "example" {
 first_name 	= "terraform"
 last_name  	= "test"
 email 	 	= "` + email + `"
 password   	= "abc123-321CBA"
 administrator = false
 force_sec_auth= false
}

resource "ionoscloud_s3_key" "example_key" {
  user_id    = ionoscloud_user.example.id
  active     = false
}

data "ionoscloud_s3_key" "example_key_data" {
 user_id    	= ionoscloud_user.example.id
 id			= ionoscloud_s3_key.example_key.id
}
`
