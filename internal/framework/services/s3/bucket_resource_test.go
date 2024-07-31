//go:build all || s3
// +build all s3

package s3_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"

	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/internal/acctest"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/internal/framework/services/s3"
)

func TestAccBucketResource(t *testing.T) {
	rName := "acctest-tf-bucket"
	name := "ionoscloud_s3_bucket.test"

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		CheckDestroy: testAccCheckBucketDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccBucketConfig_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(name, "name", rName),
					resource.TestCheckResourceAttr(name, "region", "eu-central-3"),
				),
			},
			{
				ResourceName:                         name,
				ImportStateId:                        rName,
				ImportState:                          true,
				ImportStateVerifyIdentifierAttribute: "name",
				ImportStateVerify:                    true,
			},
		},
	})
}

func testAccBucketConfig_basic(bucketName string) string {
	return fmt.Sprintf(`
resource "ionoscloud_s3_bucket" "test" {
  name = %[1]q
  region = "eu-central-3"
}
`, bucketName)
}

func testAccCheckBucketDestroy(s *terraform.State) error {
	client, err := acctest.S3Client()
	if err != nil {
		return err
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "ionoscloud_s3_bucket" {
			continue
		}

		if rs.Primary.Attributes["name"] != "" {
			err = s3.IsBucketDeleted(context.Background(), client, rs.Primary.Attributes["name"])
			if err != nil {
				return err
			}
		}
	}

	return nil
}
