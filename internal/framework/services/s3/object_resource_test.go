//go:build all || s3
// +build all s3

package s3_test

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils"
	"io"
	"os"
	"testing"

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
				),
			},
			{
				ResourceName:                         objectResourceName,
				ImportStateId:                        fmt.Sprintf("%s/%s", bucket, key),
				ImportState:                          true,
				ImportStateVerifyIdentifierAttribute: "key",
				ImportStateVerifyIgnore:              []string{"force_destroy", "content"},
				ImportStateCheck: func(s []*terraform.InstanceState) error {
					if len(s) != 1 {
						return fmt.Errorf("expected one state, got %d", len(s))
					}

					state := s[0]
					attributes := map[string]string{
						"bucket": bucket,
						"key":    key,
					}

					for k, v := range attributes {
						if state.Attributes[k] != v {
							return fmt.Errorf("attribute %s expected %s, got %s", k, v, state.Attributes[k])
						}
					}

					return nil
				},
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccObjectResource_ContentType(t *testing.T) {
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
					resource.TestCheckResourceAttr(objectResourceName, "content_type", "text/plain"),
				),
			},
		},
	})
}

func TestAccObjectResource_CacheControl(t *testing.T) {
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
					resource.TestCheckResourceAttr(objectResourceName, "cache_control", "private"),
				),
			},
		},
	})
}

func TestAccObjectResource_Content(t *testing.T) {
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
					resource.TestCheckResourceAttr(objectResourceName, "content_disposition", "attachment"),
					resource.TestCheckResourceAttr(objectResourceName, "content_encoding", "gzip"),
					resource.TestCheckResourceAttr(objectResourceName, "content_language", "en"),
				),
			},
		},
	})
}

func TestAccObjectResource_Expires(t *testing.T) {
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
					resource.TestCheckResourceAttr(objectResourceName, "expires", "2022-01-01T00:00:00Z"),
				),
			},
		},
	})
}

func TestAccObjectResource_WebsiteRedirect(t *testing.T) {
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
					resource.TestCheckResourceAttr(objectResourceName, "website_redirect", "https://example.com"),
				),
			},
		},
	})
}

func TestAccObjectResource_ServerSideEncryption(t *testing.T) {
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
					resource.TestCheckResourceAttr(objectResourceName, "server_side_encryption", "AES256"),
				),
			},
		},
	})
}

// TODO Learn how to test this
func TestAccObjectResource_ServerSideEncryptionCustomer(t *testing.T) {
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
					resource.TestCheckResourceAttr(objectResourceName, "server_side_encryption_customer_algorithm", "AES256"),
					resource.TestCheckResourceAttr(objectResourceName, "server_side_encryption_customer_key", "dpHHYOfjTUlcpotfDSNzkyWUWLtcZkoX1dlua5D1pAM="),
					resource.TestCheckResourceAttr(objectResourceName, "server_side_encryption_customer_key_md5", "56029099e69ec4ea644fb2a34d507e16"),
				),
			},
		},
	})

}

func TestAccObjectResource_Tags(t *testing.T) {
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
					resource.TestCheckResourceAttr(objectResourceName, "tags.key1", "value1"),
					resource.TestCheckResourceAttr(objectResourceName, "tags.key2", "value2"),
				),
			},
		},
	})

}

func TestAccObjectResource_Metadata(t *testing.T) {
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
  cache_control = "private"
}

`, key))
}

func testAccObjectConfig_content(bucketName, key string) string {
	return utils.ConfigCompose(testAccObjectConfig_base(bucketName), fmt.Sprintf(`
resource "ionoscloud_s3_object" "test" {
  bucket = ionoscloud_s3_bucket.test.name
  key = %[1]q
  content = "test"
  content_disposition = "attachment"
  content_encoding = "gzip"
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
  server_side_encryption_customer_algorithm = "AES256"
  server_side_encryption_customer_key = "4ZRNYBCCvL0YZeqo3f2+9qDyIfnLdbg5S99R2XWr0aw="
  server_side_encryption_customer_key_md5 = "ZeDiDFGrdO9ZXpA6TUOo4g=="
}

`, key))
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

		output, _, err := client.ObjectsApi.GetObject(ctx, rs.Primary.Attributes["bucket"], rs.Primary.Attributes["key"]).Execute()
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
