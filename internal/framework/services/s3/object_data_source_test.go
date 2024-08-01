//go:build all || s3
// +build all s3

package s3_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/internal/acctest"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils"
)

func TestAccObjectDataSource(t *testing.T) {
	bucketName := acctest.GenerateRandomResourceName(bucketPrefix)
	objKey := acctest.GenerateRandomResourceName(objectPrefix)
	dataSourceName := "data.ionoscloud_s3_object.test"
	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		Steps: []resource.TestStep{
			{
				Config: testAccObjectDataSourceConfig_basic(bucketName, objKey),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "bucket", "ionoscloud_s3_object.test", "bucket"),
					resource.TestCheckResourceAttrPair(dataSourceName, "key", "ionoscloud_s3_object.test", "key"),
					resource.TestCheckResourceAttrPair(dataSourceName, "etag", "ionoscloud_s3_object.test", "etag"),
					resource.TestCheckResourceAttrPair(dataSourceName, "content_type", "ionoscloud_s3_object.test", "content_type"),
					resource.TestCheckResourceAttr(dataSourceName, "body", "test"),
					resource.TestCheckResourceAttr(dataSourceName, "content_length", "4"),
					resource.TestCheckResourceAttr(dataSourceName, "metadata.%", "0"),
					resource.TestCheckResourceAttr(dataSourceName, "tags.%", "0"),
				),
			},
		},
	})
}

func TestAccObjectDataSource_allParams(t *testing.T) {
	bucketName := acctest.GenerateRandomResourceName(bucketPrefix)
	objKey := acctest.GenerateRandomResourceName(objectPrefix)
	dataSourceName := "data.ionoscloud_s3_object.test"
	resourceName := "ionoscloud_s3_object.test"
	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		Steps: []resource.TestStep{
			{
				Config: testAccObjectDataSourceConfig_allParams(bucketName, objKey),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "bucket", resourceName, "bucket"),
					resource.TestCheckResourceAttrPair(dataSourceName, "key", resourceName, "key"),
					resource.TestCheckResourceAttrPair(dataSourceName, "etag", resourceName, "etag"),
					resource.TestCheckResourceAttrPair(dataSourceName, "content_type", resourceName, "content_type"),
					resource.TestCheckResourceAttrPair(dataSourceName, "content_disposition", resourceName, "content_disposition"),
					resource.TestCheckResourceAttrPair(dataSourceName, "content_encoding", resourceName, "content_encoding"),
					resource.TestCheckResourceAttrPair(dataSourceName, "content_language", resourceName, "content_language"),
					resource.TestCheckResourceAttrPair(dataSourceName, "cache_control", resourceName, "cache_control"),
					resource.TestCheckResourceAttrPair(dataSourceName, "expires", resourceName, "expires"),
					resource.TestCheckResourceAttrPair(dataSourceName, "website_redirect", resourceName, "website_redirect"),
					resource.TestCheckResourceAttrPair(dataSourceName, "server_side_encryption", resourceName, "server_side_encryption"),
					resource.TestCheckResourceAttrPair(dataSourceName, "storage_class", resourceName, "storage_class"),
					resource.TestCheckResourceAttrPair(dataSourceName, "version_id", resourceName, "version_id"),
					resource.TestCheckResourceAttr(dataSourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(dataSourceName, "metadata.%", "1"),
					resource.TestCheckResourceAttr(dataSourceName, "metadata.mk", "mv"),
					resource.TestCheckResourceAttr(dataSourceName, "tags.tk", "tv"),
				),
			},
		},
	})

}

func testAccObjectDataSourceConfig_basic(bucketName, objectName string) string {
	return utils.ConfigCompose(testAccObjectConfig_basic(bucketName, objectName), fmt.Sprintf(`
data "ionoscloud_s3_object" "test" {
	bucket = ionoscloud_s3_bucket.test.name
	key =  ionoscloud_s3_object.test.key
}
`))
}

func testAccObjectDataSourceConfig_allParams(bucketName, objectName string) string {
	return utils.ConfigCompose(testAccObjectConfig_base(bucketName), fmt.Sprintf(`
resource "ionoscloud_s3_object" "test" {
	bucket = ionoscloud_s3_bucket.test.name
	key =  %[1]q
	 content             = <<CONTENT
{
  "msg": "Hello World"
}
CONTENT
	content_type        = "application/json"
	cache_control       = "no-cache"
	content_disposition = "attachment"
	content_encoding    = "identity"
	content_language    = "en-GB"
	expires			 = "2022-10-07T12:34:56Z"
	website_redirect = "https://www.ionos.com"
	server_side_encryption = "AES256"

	tags = {
		tk = "tv"
	}

	metadata = {
		"mk" = "mv"
	}

}
data "ionoscloud_s3_object" "test" {
	bucket = ionoscloud_s3_bucket.test.name
	key =  ionoscloud_s3_object.test.key
}
`, objectName))
}
