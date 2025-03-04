//go:build all || objectstorage
// +build all objectstorage

package objectstorage_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"

	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/internal/acctest"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils"
)

func TestAccObjectCopy_basic(t *testing.T) {
	ctx := context.Background()
	rNameSource := acctest.GenerateRandomResourceName(bucketPrefix)
	rNameTarget := acctest.GenerateRandomResourceName(bucketPrefix)
	resourceName := "ionoscloud_s3_object_copy.test"
	sourceName := "ionoscloud_s3_object.source"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		CheckDestroy: testAccCheckObjectCopyDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccObjectCopyConfig_basic(rNameSource, "source", rNameTarget, "target"),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckObjectCopyExists(ctx, resourceName),
					resource.TestCheckResourceAttr(resourceName, "bucket", rNameTarget),
					resource.TestCheckResourceAttr(resourceName, "content_type", "application/octet-stream"),
					resource.TestCheckNoResourceAttr(resourceName, "copy_if_match"),
					resource.TestCheckNoResourceAttr(resourceName, "copy_if_modified_since"),
					resource.TestCheckNoResourceAttr(resourceName, "copy_if_none_match"),
					resource.TestCheckNoResourceAttr(resourceName, "copy_if_unmodified_since"),
					resource.TestCheckNoResourceAttr(resourceName, "server_side_encryption_algorithm"),
					resource.TestCheckNoResourceAttr(resourceName, "server_side_encryption_key"),
					resource.TestCheckNoResourceAttr(resourceName, "server_side_encryption_key_md5"),
					resource.TestCheckResourceAttrPair(resourceName, "etag", sourceName, "etag"),
					resource.TestCheckNoResourceAttr(resourceName, "expires"),
					resource.TestCheckResourceAttr(resourceName, "force_destroy", "false"),
					resource.TestCheckResourceAttr(resourceName, "grant.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "key", "target"),
					resource.TestCheckResourceAttrSet(resourceName, "last_modified"),
					resource.TestCheckResourceAttr(resourceName, "metadata.%", "0"),
					resource.TestCheckNoResourceAttr(resourceName, "metadata_directive"),
					resource.TestCheckNoResourceAttr(resourceName, "request_payer"),
					resource.TestCheckResourceAttr(resourceName, "source", fmt.Sprintf("%s/%s", rNameSource, "source")),
					resource.TestCheckNoResourceAttr(resourceName, "source_server_side_encryption_algorithm"),
					resource.TestCheckNoResourceAttr(resourceName, "source_server_side_encryption_key"),
					resource.TestCheckNoResourceAttr(resourceName, "source_server_side_encryption_key_md5"),
					resource.TestCheckResourceAttr(resourceName, "storage_class", "STANDARD"),
					resource.TestCheckNoResourceAttr(resourceName, "tagging_directive"),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "0"),
					resource.TestCheckResourceAttr(resourceName, "version_id", ""),
				),
			},
		},
	})
}

func TestAccObjectCopy_metadata(t *testing.T) {
	ctx := context.Background()
	rName1 := acctest.GenerateRandomResourceName(bucketPrefix)
	rName2 := acctest.GenerateRandomResourceName(bucketPrefix)
	resourceName := "ionoscloud_s3_object_copy.test"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		CheckDestroy: testAccCheckObjectCopyDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccObjectCopyConfig_metadata(rName1, "source", rName2, "target"),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckObjectCopyExists(ctx, resourceName),
					resource.TestCheckResourceAttr(resourceName, "metadata_directive", "REPLACE"),
					resource.TestCheckResourceAttr(resourceName, "metadata.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "metadata.mk1", "mv1"),
				),
			},
		},
	})
}

func TestAccObjectCopy_sourceWithSlashes(t *testing.T) {
	ctx := context.Background()
	rName1 := acctest.GenerateRandomResourceName(bucketPrefix)
	rName2 := acctest.GenerateRandomResourceName(bucketPrefix)
	resourceName := "ionoscloud_s3_object_copy.test"
	sourceKey := "dir1/dir2/source"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		CheckDestroy: testAccCheckObjectCopyDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccObjectCopyConfig_baseSourceAndTargetBuckets(rName1, rName2),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckBucketAddObjects(ctx, "ionoscloud_s3_bucket.source", sourceKey),
				),
			},
			{
				Config: testAccObjectCopyConfig_externalSourceObject(rName1, sourceKey, rName2, "target"),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckObjectCopyExists(ctx, resourceName),
					resource.TestCheckResourceAttr(resourceName, "source", fmt.Sprintf("%s/%s", rName1, sourceKey)),
				),
			},
		},
	})
}

func TestAccObjectCopy_objectLockLegalHold(t *testing.T) {
	ctx := context.Background()
	rName1 := acctest.GenerateRandomResourceName(bucketPrefix)
	rName2 := acctest.GenerateRandomResourceName(bucketPrefix)
	resourceName := "ionoscloud_s3_object_copy.test"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		CheckDestroy: testAccCheckObjectCopyDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccObjectCopyConfig_lockLegalHold(rName1, "source", rName2, "target", "ON"),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckObjectCopyExists(ctx, resourceName),
					resource.TestCheckResourceAttr(resourceName, "object_lock_legal_hold", "ON"),
				),
			},
			{
				Config: testAccObjectCopyConfig_lockLegalHold(rName1, "source", rName2, "target", "OFF"),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckObjectCopyExists(ctx, resourceName),
					resource.TestCheckResourceAttr(resourceName, "object_lock_legal_hold", "OFF"),
				),
			},
		},
	})
}

func TestAccObjectCopy_targetWithMultipleSlashes(t *testing.T) {
	ctx := context.Background()
	rName1 := acctest.GenerateRandomResourceName(bucketPrefix)
	rName2 := acctest.GenerateRandomResourceName(bucketPrefix)
	resourceName := "ionoscloud_s3_object_copy.test"
	targetKey := "/dir//target/"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		CheckDestroy: testAccCheckObjectCopyDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccObjectCopyConfig_basic(rName1, "source", rName2, targetKey),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "key", targetKey),
					resource.TestCheckResourceAttr(resourceName, "source", fmt.Sprintf("%s/%s", rName1, "source")),
				),
			},
		},
	})
}

func testAccCheckObjectCopyDestroy(ctx context.Context) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := acctest.NewTestBundleClientFromEnv().S3Client.GetBaseClient()

		for _, rs := range s.RootModule().Resources {
			if rs.Type != "ionoscloud_s3_object_copy" {
				continue
			}

			if rs.Primary.Attributes["bucket"] != "" && rs.Primary.Attributes["key"] != "" {
				_, apiResponse, err := client.ObjectsApi.HeadObject(ctx, rs.Primary.Attributes["bucket"], rs.Primary.Attributes["key"]).Execute()
				if err != nil {
					if !apiResponse.HttpNotFound() {
						return fmt.Errorf("error checking object %s: %w", rs.Primary.ID, err)
					}

					return nil
				}

				return fmt.Errorf("object %s still exists", rs.Primary.ID)
			}
		}

		return nil
	}
}

func testAccCheckObjectCopyExists(ctx context.Context, n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not Found: %s", n)
		}

		client := acctest.NewTestBundleClientFromEnv().S3Client.GetBaseClient()

		if rs.Primary.Attributes["bucket"] != "" && rs.Primary.Attributes["key"] != "" {
			_, _, err := client.ObjectsApi.HeadObject(ctx, rs.Primary.Attributes["bucket"], rs.Primary.Attributes["key"]).Execute()
			if err != nil {
				return fmt.Errorf("error checking object %s: %w", rs.Primary.ID, err)
			}
		}

		return nil
	}
}

func testAccObjectCopyConfig_baseSourceAndTargetBuckets(sourceBucket, targetBucket string) string {
	return fmt.Sprintf(`
resource "ionoscloud_s3_bucket" "source" {
  name = %[1]q

  force_destroy = true
}

resource "ionoscloud_s3_bucket" "target" {
  name = %[2]q
}
`, sourceBucket, targetBucket)
}

func testAccObjectCopyConfig_baseSourceObject(sourceBucket, sourceKey, targetBucket string) string {
	return utils.ConfigCompose(testAccObjectCopyConfig_baseSourceAndTargetBuckets(sourceBucket, targetBucket), fmt.Sprintf(`
resource "ionoscloud_s3_object" "source" {
  bucket  = ionoscloud_s3_bucket.source.name
  key     = %[1]q
  content = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
  content_type = "application/octet-stream"
}
`, sourceKey))
}

func testAccObjectCopyConfig_basic(sourceBucket, sourceKey, targetBucket, targetKey string) string {
	return utils.ConfigCompose(testAccObjectCopyConfig_baseSourceObject(sourceBucket, sourceKey, targetBucket), fmt.Sprintf(`
resource "ionoscloud_s3_object_copy" "test" {
  bucket = ionoscloud_s3_bucket.target.name
  key    = %[1]q
  source = "${ionoscloud_s3_bucket.source.name}/${ionoscloud_s3_object.source.key}"
}
`, targetKey))
}

func testAccObjectCopyConfig_metadata(sourceBucket, sourceKey, targetBucket, targetKey string) string {
	return utils.ConfigCompose(testAccObjectCopyConfig_baseSourceObject(sourceBucket, sourceKey, targetBucket), fmt.Sprintf(`
resource "ionoscloud_s3_object_copy" "test" {
  bucket = ionoscloud_s3_bucket.target.name
  key    = %[1]q
  source = "${ionoscloud_s3_bucket.source.name}/${ionoscloud_s3_object.source.key}"

  metadata_directive = "REPLACE"

  metadata = {
    "mk1" = "mv1"
  }
}
`, targetKey))
}

func testAccObjectCopyConfig_externalSourceObject(sourceBucket, sourceKey, targetBucket, targetKey string) string {
	return utils.ConfigCompose(testAccObjectCopyConfig_baseSourceAndTargetBuckets(sourceBucket, targetBucket), fmt.Sprintf(`
resource "ionoscloud_s3_object_copy" "test" {
  bucket = ionoscloud_s3_bucket.target.name
  key    = %[2]q
  source = "${ionoscloud_s3_bucket.source.name}/%[1]s"
}
`, sourceKey, targetKey))
}

func testAccObjectCopyConfig_lockLegalHold(sourceBucket, sourceKey, targetBucket, targetKey, legalHoldStatus string) string {
	return fmt.Sprintf(`
resource "ionoscloud_s3_bucket" "source" {
  name = %[1]q

  force_destroy = true
}

resource "ionoscloud_s3_bucket" "target" {
  name = %[3]q

  object_lock_enabled = true

  force_destroy = true
}

resource "ionoscloud_s3_bucket_versioning" "target" {
  bucket = ionoscloud_s3_bucket.target.name
  versioning_configuration {
    status = "Enabled"
  }
}

resource "ionoscloud_s3_object" "source" {
  bucket  = ionoscloud_s3_bucket.source.name
  key     = %[2]q
  content = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
}

resource "ionoscloud_s3_object_copy" "test" {
  # Must have bucket versioning enabled first
  bucket = ionoscloud_s3_bucket_versioning.target.bucket
  key    = %[4]q
  source = "${ionoscloud_s3_bucket.source.name}/${ionoscloud_s3_object.source.key}"

  object_lock_legal_hold = %[5]q
  force_destroy          = true
}
`, sourceBucket, sourceKey, targetBucket, targetKey, legalHoldStatus)
}
