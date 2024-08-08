//go:build all || s3
// +build all s3

package s3_test

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	s3 "github.com/ionos-cloud/sdk-go-s3"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils"
	"io"
	"os"
	"testing"
	"time"

	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/internal/acctest"
)

const (
	objectResourceName = "ionoscloud_s3_object.test"
	objectPrefix       = "acctest-tf-object-"
	bucketPrefix       = "acctest-tf-bucket-"
)

func TestAccObjectResourceBasic(t *testing.T) {
	ctx := context.Background()
	var body string
	bucket := acctest.GenerateRandomResourceName(bucketPrefix)
	key := acctest.GenerateRandomResourceName(objectPrefix)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		CheckDestroy: testAccCheckObjectDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccObjectConfig_basic(bucket, key),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckObjectExists(ctx, objectResourceName, &body),
					testAccCheckObjectBody(&body, "test"),
					resource.TestCheckResourceAttr(objectResourceName, "bucket", bucket),
					resource.TestCheckResourceAttr(objectResourceName, "key", key),
					resource.TestCheckResourceAttr(objectResourceName, "content", "test"),
					resource.TestCheckResourceAttrSet(objectResourceName, "etag"),
					resource.TestCheckResourceAttr(objectResourceName, "content_type", "text/plain"),
					resource.TestCheckResourceAttr(objectResourceName, "storage_class", "STANDARD"),
					resource.TestCheckResourceAttr(objectResourceName, "force_destroy", "false"),
					resource.TestCheckResourceAttr(objectResourceName, "version_id", ""),
				),
			},
			{
				ResourceName:                         objectResourceName,
				ImportStateId:                        fmt.Sprintf("%s/%s", bucket, key),
				ImportState:                          true,
				ImportStateVerifyIdentifierAttribute: "key",
				ImportStateVerifyIgnore:              []string{"force_destroy", "content"},
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					rs, ok := s.RootModule().Resources[objectResourceName]
					if !ok {
						return "", fmt.Errorf("Not Found: %s", objectResourceName)
					}

					return fmt.Sprintf("%s/%s", rs.Primary.Attributes["bucket"], rs.Primary.Attributes["key"]), nil
				},
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccObjectResourceDirectory(t *testing.T) {
	ctx := context.Background()
	var body string
	bucket := acctest.GenerateRandomResourceName(bucketPrefix)
	key := "dir/" + acctest.GenerateRandomResourceName(objectPrefix)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		CheckDestroy: testAccCheckObjectDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccObjectConfig_basic(bucket, key),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckObjectExists(ctx, objectResourceName, &body),
					testAccCheckObjectBody(&body, "test"),
					resource.TestCheckResourceAttr(objectResourceName, "bucket", bucket),
					resource.TestCheckResourceAttr(objectResourceName, "key", key),
					resource.TestCheckResourceAttr(objectResourceName, "content", "test"),
					resource.TestCheckResourceAttrSet(objectResourceName, "etag"),
					resource.TestCheckResourceAttr(objectResourceName, "content_type", "text/plain"),
					resource.TestCheckResourceAttr(objectResourceName, "storage_class", "STANDARD"),
					resource.TestCheckResourceAttr(objectResourceName, "force_destroy", "false"),
				),
			},
		},
	})
}

func TestAccObjectResource_ContentType(t *testing.T) {
	ctx := context.Background()
	var body string
	bucket := acctest.GenerateRandomResourceName(bucketPrefix)
	key := acctest.GenerateRandomResourceName(objectPrefix)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		CheckDestroy: testAccCheckObjectDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccObjectConfig_contentType(bucket, key),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckObjectExists(ctx, objectResourceName, &body),
					testAccCheckObjectBody(&body, "test"),
					resource.TestCheckResourceAttr(objectResourceName, "content_type", "text/plain"),
				),
			},
		},
	})
}

func TestAccObjectResource_CacheControl(t *testing.T) {
	ctx := context.Background()
	var body string
	bucket := acctest.GenerateRandomResourceName(bucketPrefix)
	key := acctest.GenerateRandomResourceName(objectPrefix)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		CheckDestroy: testAccCheckObjectDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccObjectConfig_cacheControl(bucket, key),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckObjectExists(ctx, objectResourceName, &body),
					testAccCheckObjectBody(&body, "test"),
					resource.TestCheckResourceAttr(objectResourceName, "cache_control", "private"),
				),
			},
		},
	})
}

func TestAccObjectResource_Content(t *testing.T) {
	ctx := context.Background()
	var body string
	bucket := acctest.GenerateRandomResourceName(bucketPrefix)
	key := acctest.GenerateRandomResourceName(objectPrefix)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		CheckDestroy: testAccCheckObjectDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccObjectConfig_content(bucket, key),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckObjectExists(ctx, objectResourceName, &body),
					testAccCheckObjectBody(&body, "<data><test>a</test></data>"),
					resource.TestCheckResourceAttr(objectResourceName, "content_disposition", "attachment; filename=\"test.xml\""),
					resource.TestCheckResourceAttr(objectResourceName, "content_encoding", "identity"),
					resource.TestCheckResourceAttr(objectResourceName, "content_language", "en"),
				),
			},
		},
	})
}

func TestAccObjectResource_Expires(t *testing.T) {
	ctx := context.Background()
	var body string
	bucket := acctest.GenerateRandomResourceName(bucketPrefix)
	key := acctest.GenerateRandomResourceName(objectPrefix)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		CheckDestroy: testAccCheckObjectDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccObjectConfig_expires(bucket, key),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckObjectExists(ctx, objectResourceName, &body),
					testAccCheckObjectBody(&body, "test"),
					resource.TestCheckResourceAttr(objectResourceName, "expires", "2022-01-01T00:00:00Z"),
				),
			},
		},
	})
}

func TestAccObjectResource_WebsiteRedirect(t *testing.T) {
	ctx := context.Background()
	var body string
	bucket := acctest.GenerateRandomResourceName(bucketPrefix)
	key := acctest.GenerateRandomResourceName(objectPrefix)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		CheckDestroy: testAccCheckObjectDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccObjectConfig_websiteRedirect(bucket, key),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckObjectExists(ctx, objectResourceName, &body),
					testAccCheckObjectBody(&body, "test"),
					resource.TestCheckResourceAttr(objectResourceName, "website_redirect", "https://example.com"),
				),
			},
		},
	})
}

func TestAccObjectResource_ServerSideEncryption(t *testing.T) {
	ctx := context.Background()
	var body string
	bucket := acctest.GenerateRandomResourceName(bucketPrefix)
	key := acctest.GenerateRandomResourceName(objectPrefix)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		CheckDestroy: testAccCheckObjectDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccObjectConfig_serverSideEncryption(bucket, key),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckObjectExists(ctx, objectResourceName, &body),
					testAccCheckObjectBody(&body, "test"),
					resource.TestCheckResourceAttr(objectResourceName, "server_side_encryption", "AES256"),
				),
			},
		},
	})
}

func TestAccObjectResource_ServerSideEncryptionCustomer(t *testing.T) {
	ctx := context.Background()
	var body string
	bucket := acctest.GenerateRandomResourceName(bucketPrefix)
	key := acctest.GenerateRandomResourceName(objectPrefix)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		CheckDestroy: testAccCheckObjectDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccObjectConfig_serverSideEncryptionCustomer(bucket, key),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckObjectExists(ctx, objectResourceName, &body),
					testAccCheckObjectBody(&body, "test"),
					resource.TestCheckResourceAttr(objectResourceName, "server_side_encryption_customer_algorithm", "AES256"),
					resource.TestCheckResourceAttr(objectResourceName, "server_side_encryption_customer_key", "MzItYnl0ZS1lbmNyeXB0aW9uLWtleS0xMjM0NTY3ODk="),
					resource.TestCheckResourceAttr(objectResourceName, "server_side_encryption_customer_key_md5", "mUyItJiR9XqV9ARcO72seQ=="),
				),
			},
		},
	})

}

func TestAccObjectResource_ObjectLock(t *testing.T) {
	ctx := context.Background()
	var body string
	bucket := acctest.GenerateRandomResourceName(bucketPrefix)
	key := acctest.GenerateRandomResourceName(objectPrefix)

	retainUntilDate := time.Now().UTC().AddDate(0, 0, 1).Format(time.RFC3339)
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		CheckDestroy: testAccCheckObjectDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccObjectConfig_objectLockNoRetention(bucket, key, "stuff"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckObjectExists(ctx, objectResourceName, &body),
					testAccCheckObjectBody(&body, "stuff"),
					resource.TestCheckNoResourceAttr(objectResourceName, "object_lock_mode"),
					resource.TestCheckNoResourceAttr(objectResourceName, "object_legal_hold"),
					resource.TestCheckNoResourceAttr(objectResourceName, "object_lock_retain_until_date"),
				),
			},
			{
				Config: testAccObjectConfig_objectLockRetention(bucket, key, retainUntilDate, "stuff"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckObjectExists(ctx, objectResourceName, &body),
					testAccCheckObjectBody(&body, "stuff"),
					resource.TestCheckResourceAttr(objectResourceName, "object_lock_mode", "GOVERNANCE"),
					resource.TestCheckResourceAttr(objectResourceName, "object_lock_retain_until_date", retainUntilDate),
				),
			},
			{
				Config: testAccObjectConfig_objectLockNoRetention(bucket, key, "changed stuff"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckObjectExists(ctx, objectResourceName, &body),
					testAccCheckObjectBody(&body, "changed stuff"),
					resource.TestCheckNoResourceAttr(objectResourceName, "object_lock_mode"),
					resource.TestCheckNoResourceAttr(objectResourceName, "object_legal_hold"),
					resource.TestCheckNoResourceAttr(objectResourceName, "object_lock_retain_until_date"),
				),
			},
		},
	})
}

func TestAccObjectResource_Tags(t *testing.T) {
	ctx := context.Background()
	var body string
	bucket := acctest.GenerateRandomResourceName(bucketPrefix)
	key := acctest.GenerateRandomResourceName(objectPrefix)
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		CheckDestroy: testAccCheckObjectDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccObjectConfig_tags(bucket, key),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckObjectExists(ctx, objectResourceName, &body),
					testAccCheckObjectBody(&body, "test"),
					resource.TestCheckResourceAttr(objectResourceName, "tags.key1", "value1"),
					resource.TestCheckResourceAttr(objectResourceName, "tags.key2", "value2"),
				),
			},
		},
	})

}

func TestAccObjectResource_Metadata(t *testing.T) {
	ctx := context.Background()
	var body string
	bucket := acctest.GenerateRandomResourceName(bucketPrefix)
	key := acctest.GenerateRandomResourceName(objectPrefix)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		CheckDestroy: testAccCheckObjectDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccObjectConfig_metadata(bucket, key),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckObjectExists(ctx, objectResourceName, &body),
					testAccCheckObjectBody(&body, "test"),
					resource.TestCheckResourceAttr(objectResourceName, "metadata.key1", "value1"),
					resource.TestCheckResourceAttr(objectResourceName, "metadata.key2", "value2"),
				),
			},
		},
	})
}

func TestAccObjectResource_SourceFile(t *testing.T) {
	ctx := context.Background()
	bucket := acctest.GenerateRandomResourceName(bucketPrefix)
	key := acctest.GenerateRandomResourceName(objectPrefix)
	source := testAccObjectCreateTempFile(t, "{anything will do }")
	var body string
	defer os.Remove(source)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		CheckDestroy: testAccCheckObjectDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccObjectConfig_sourceFile(bucket, key, source),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckObjectExists(ctx, objectResourceName, &body),
					testAccCheckObjectBody(&body, "{anything will do }"),
					resource.TestCheckResourceAttr(objectResourceName, "source", source),
				),
			},
		},
	})

}

func testAccObjectConfig_base(bucketName string) string {
	return fmt.Sprintf(`
resource "ionoscloud_s3_bucket" "test" {
  name = %[1]q
}
`, bucketName)
}

func testAccObjectConfig_basic(bucketName, key string) string {
	return utils.ConfigCompose(testAccObjectConfig_base(bucketName), fmt.Sprintf(`
resource "ionoscloud_s3_object" "test" {
  bucket = ionoscloud_s3_bucket.test.name
  key = %[1]q
  content = "test"
  content_type = "text/plain"
}

`, key))
}

func testAccObjectConfig_contentType(bucketName, key string) string {
	return utils.ConfigCompose(testAccObjectConfig_base(bucketName), fmt.Sprintf(`
resource "ionoscloud_s3_object" "test" {
  bucket = ionoscloud_s3_bucket.test.name
  key = %[1]q
  content = "test"
  content_type = "text/plain"
}

`, key))
}

func testAccObjectConfig_cacheControl(bucketName, key string) string {
	return utils.ConfigCompose(testAccObjectConfig_base(bucketName), fmt.Sprintf(`
resource "ionoscloud_s3_object" "test" {
  bucket = ionoscloud_s3_bucket.test.name
  key = %[1]q
  content = "test"
  content_type = "text/plain"
  cache_control = "private"
}

`, key))
}

func testAccObjectConfig_content(bucketName, key string) string {
	return utils.ConfigCompose(testAccObjectConfig_base(bucketName), fmt.Sprintf(`
resource "ionoscloud_s3_object" "test" {
  bucket = ionoscloud_s3_bucket.test.name
  key = %[1]q
  content = "<data><test>a</test></data>"
  content_type = "application/xml"
  content_disposition = "attachment; filename=\"test.xml\""
  content_encoding = "identity"
  content_language = "en"
}

`, key))
}

func testAccObjectConfig_expires(bucketName, key string) string {
	return utils.ConfigCompose(testAccObjectConfig_base(bucketName), fmt.Sprintf(`
resource "ionoscloud_s3_object" "test" {
  bucket = ionoscloud_s3_bucket.test.name
  key = %[1]q
  content = "test"
  expires = "2022-01-01T00:00:00Z"
}

`, key))
}

func testAccObjectConfig_websiteRedirect(bucketName, key string) string {
	return utils.ConfigCompose(testAccObjectConfig_base(bucketName), fmt.Sprintf(`
resource "ionoscloud_s3_object" "test" {
  bucket = ionoscloud_s3_bucket.test.name
  key = %[1]q
  content = "test"
  website_redirect = "https://example.com"
}

`, key))
}

func testAccObjectConfig_serverSideEncryption(bucketName, key string) string {
	return utils.ConfigCompose(testAccObjectConfig_base(bucketName), fmt.Sprintf(`
resource "ionoscloud_s3_object" "test" {
  bucket = ionoscloud_s3_bucket.test.name
  key = %[1]q
  content = "test"
  server_side_encryption = "AES256"
}

`, key))
}

func testAccObjectConfig_serverSideEncryptionCustomer(bucketName, key string) string {
	return utils.ConfigCompose(testAccObjectConfig_base(bucketName), fmt.Sprintf(`
resource "ionoscloud_s3_object" "test" {
  bucket = ionoscloud_s3_bucket.test.name
  key = %[1]q
  content = "test"
  content_type = "text/plain"
  server_side_encryption_customer_algorithm = "AES256"
  server_side_encryption_customer_key = "MzItYnl0ZS1lbmNyeXB0aW9uLWtleS0xMjM0NTY3ODk="
  server_side_encryption_customer_key_md5 = "mUyItJiR9XqV9ARcO72seQ=="
}

`, key))
}

func testAccObjectConfig_objectLockNoRetention(bucketName, key, content string) string {
	return fmt.Sprintf(`
resource "ionoscloud_s3_bucket" "test" {
  name = %[1]q
  object_lock_enabled = true
}

resource "ionoscloud_s3_object" "test" {
  bucket = ionoscloud_s3_bucket.test.name
  key = %[2]q
  content = %[3]q
  content_type = "text/plain"
  force_destroy = true
}

`, bucketName, key, content)
}

func testAccObjectConfig_objectLockRetention(bucketName, key, retention, content string) string {
	return fmt.Sprintf(`
resource "ionoscloud_s3_bucket" "test" {
  name = %[1]q
  object_lock_enabled = true
}

resource "ionoscloud_s3_object" "test" {
  bucket = ionoscloud_s3_bucket.test.name
  key = %[2]q
  content = %[4]q
  content_type = "text/plain"
  object_lock_mode = "GOVERNANCE"
  object_lock_retain_until_date = %[3]q
  force_destroy = true
}

`, bucketName, key, retention, content)
}

func testAccObjectConfig_tags(bucketName, key string) string {
	return utils.ConfigCompose(testAccObjectConfig_base(bucketName), fmt.Sprintf(`
resource "ionoscloud_s3_object" "test" {
  bucket = ionoscloud_s3_bucket.test.name
  key = %[1]q
  content = "test"
  tags = {
	key1 = "value1"
	key2 = "value2"
  }
}

`, key))
}

func testAccObjectConfig_metadata(bucketName, key string) string {
	return utils.ConfigCompose(testAccObjectConfig_base(bucketName), fmt.Sprintf(`
resource "ionoscloud_s3_object" "test" {
  bucket = ionoscloud_s3_bucket.test.name
  key = %[1]q
  content = "test"
  metadata = {
	key1 = "value1"
    key2 = "value2"
  }
}

`, key))
}

func testAccObjectConfig_sourceFile(bucketName, key, source string) string {
	return utils.ConfigCompose(testAccObjectConfig_base(bucketName), fmt.Sprintf(`
resource "ionoscloud_s3_object" "test" {
  bucket = ionoscloud_s3_bucket.test.name
  key = %[1]q
  source = %[2]q
}

`, key, source))
}

func testAccCheckObjectExists(ctx context.Context, n string, body *string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not Found: %s", n)
		}

		client, err := acctest.S3Client()
		if err != nil {
			return err
		}

		output, _, err := buildGetObjectRequest(ctx, client, rs.Primary.Attributes).Execute()
		if err != nil {
			return fmt.Errorf("error fetching object %s: %w", rs.Primary.ID, err)
		}

		data, err := io.ReadAll(output)
		if err != nil {
			return fmt.Errorf("error reading object body: %w", err)
		}

		*body = string(data)

		return nil
	}
}

func buildGetObjectRequest(ctx context.Context, client *s3.APIClient, attributes map[string]string) s3.ApiGetObjectRequest {
	req := client.ObjectsApi.GetObject(ctx, attributes["bucket"], attributes["key"])
	if attributes["version_id"] != "" {
		req = req.VersionId(attributes["version_id"])
	}

	if attributes["etag"] != "" {
		req = req.IfMatch(attributes["etag"])
	}

	if attributes["server_side_encryption_customer_algorithm"] != "" {
		req = req.XAmzServerSideEncryptionCustomerAlgorithm(attributes["server_side_encryption_customer_algorithm"])
	}

	if attributes["server_side_encryption_customer_key"] != "" {
		req = req.XAmzServerSideEncryptionCustomerKey(attributes["server_side_encryption_customer_key"])
	}

	if attributes["server_side_encryption_customer_key_md5"] != "" {
		req = req.XAmzServerSideEncryptionCustomerKeyMD5(attributes["server_side_encryption_customer_key_md5"])
	}

	return req
}

func testAccCheckObjectBody(got *string, want string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if *got != want {
			return fmt.Errorf("S3 Object body = %v, want %v", got, want)
		}

		return nil
	}
}

func testAccObjectCreateTempFile(t *testing.T, data string) string {
	tmpFile, err := os.CreateTemp("", "tf-acc-object")
	if err != nil {
		t.Fatal(err)
	}
	filename := tmpFile.Name()

	err = os.WriteFile(filename, []byte(data), 0644)
	if err != nil {
		os.Remove(filename)
		t.Fatal(err)
	}

	return filename
}

func testAccCheckObjectDestroy(s *terraform.State) error {
	client, err := acctest.S3Client()
	if err != nil {
		return err
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "ionoscloud_s3_object" {
			continue
		}

		if rs.Primary.Attributes["bucket"] != "" && rs.Primary.Attributes["key"] != "" {
			_, apiResponse, err := client.ObjectsApi.HeadObject(context.Background(), rs.Primary.Attributes["bucket"], rs.Primary.Attributes["key"]).Execute()
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
