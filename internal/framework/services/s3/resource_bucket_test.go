//go:build all || s3
// +build all s3

package s3_test

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/internal/framework/services/s3"
	"log"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/internal/acctest"
)

func TestAccBucketResource(t *testing.T) {
	rName := acctest.GenerateRandomResourceName(bucketPrefix)
	name := "ionoscloud_s3_bucket.test"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		CheckDestroy: testAccCheckBucketDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccBucketConfig_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBucketExists(context.Background(), name),
					resource.TestCheckResourceAttr(name, "name", rName),
					resource.TestCheckResourceAttr(name, "region", "eu-central-3"),
					resource.TestCheckResourceAttr(name, "object_lock_enabled", "false"),
				),
			},
			{
				ResourceName:                         name,
				ImportStateId:                        rName,
				ImportState:                          true,
				ImportStateVerifyIdentifierAttribute: "name",
				ImportStateVerifyIgnore:              []string{"force_destroy"},
				ImportStateVerify:                    true,
			},
		},
	})
}

func TestAccBucketResource_ForceDestroy(t *testing.T) {
	rName := acctest.GenerateRandomResourceName(bucketPrefix)
	name := "ionoscloud_s3_bucket.test"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		CheckDestroy: testAccCheckBucketDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccBucketConfig_forceDestroy(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBucketExists(context.Background(), name),
					testAccCheckBucketAddObjects(context.Background(), name, "test1", "test2"),
					resource.TestCheckResourceAttr(name, "force_destroy", "true"),
				),
			},
		},
	})

}

func TestAccBucketResource_ObjectLockEnabled(t *testing.T) {
	rName := acctest.GenerateRandomResourceName(bucketPrefix)
	name := "ionoscloud_s3_bucket.test"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		CheckDestroy: testAccCheckBucketDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccBucketConfig_objectLockEnabled(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(name, "name", rName),
					resource.TestCheckResourceAttr(name, "region", "eu-central-3"),
					resource.TestCheckResourceAttr(name, "object_lock_enabled", "true"),
				),
			},
		},
	})
}

func TestAccBucketResource_ForceDestroyObjectVersions(t *testing.T) {
	ctx := context.Background()
	rName := acctest.GenerateRandomResourceName(bucketPrefix)
	name := "ionoscloud_s3_bucket.test"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		CheckDestroy: testAccCheckBucketDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccBucketConfig_forceDestroyObjectVersions(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBucketExists(ctx, name),
					testAccCheckBucketAddObjects(ctx, name, "test1", "test2"),
					testAccCheckBucketDeleteObjects(ctx, name, "test1"), // Creates Delete marker
					testAccCheckBucketAddObjects(ctx, name, "test1"),
				),
			},
		},
	})
}

func TestAccBucketResource_ForceDestroyObjectLockEnabled(t *testing.T) {
	ctx := context.Background()
	rName := acctest.GenerateRandomResourceName(bucketPrefix)
	name := "ionoscloud_s3_bucket.test"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		CheckDestroy: testAccCheckBucketDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccBucketConfig_forceDestroyObjectLockEnabled(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBucketExists(ctx, name),
					testAccCheckBucketAddObjectsWithLegalHold(ctx, name, "test1", "test2"),
				),
			},
		},
	})
}

func testAccBucketConfig_basic(bucketName string) string {
	return fmt.Sprintf(`
resource "ionoscloud_s3_bucket" "test" {
  name = %[1]q
}
`, bucketName)
}

func testAccBucketConfig_forceDestroy(bucketName string) string {
	return fmt.Sprintf(`
resource "ionoscloud_s3_bucket" "test" {
  name = %[1]q
  force_destroy = true
}
`, bucketName)
}

func testAccBucketConfig_forceDestroyObjectVersions(bucketName string) string {
	return fmt.Sprintf(`
resource "ionoscloud_s3_bucket" "test" {
  name = %[1]q
  force_destroy = true
}

resource "ionoscloud_s3_bucket_versioning" "test" {
  bucket = ionoscloud_s3_bucket.test.name
  versioning_configuration {
    status = "Enabled"
  }
}
`, bucketName)
}

func testAccBucketConfig_objectLockEnabled(bucketName string) string {
	return fmt.Sprintf(`
resource "ionoscloud_s3_bucket" "test" {
  name = %[1]q
  object_lock_enabled = true
}
`, bucketName)
}

func testAccBucketConfig_forceDestroyObjectLockEnabled(bucketName string) string {
	return fmt.Sprintf(`
resource "ionoscloud_s3_bucket" "test" {
  name = %[1]q
  force_destroy = true
  object_lock_enabled = true
}

resource "ionoscloud_s3_bucket_versioning" "bucket" {
  bucket = ionoscloud_s3_bucket.test.name
  versioning_configuration {
    status = "Enabled"
  }
}
`, bucketName)
}

func testAccCheckBucketAddObjectsWithLegalHold(ctx context.Context, n string, keys ...string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs := s.RootModule().Resources[n]
		client, err := acctest.S3Client()
		if err != nil {
			return err
		}

		for _, key := range keys {
			body, err := createTempFile("test", "test")
			if err != nil {
				return fmt.Errorf("createTempFile error: %w", err)
			}
			_, err = client.ObjectsApi.PutObject(ctx, rs.Primary.Attributes["name"], key).Body(body).XAmzObjectLockLegalHold("ON").Execute()
			if err != nil {
				return fmt.Errorf("PutObject error: %w", err)
			}

			err = os.Remove(body.Name())
			if err != nil {
				log.Printf("failed to remove temp file: %s", err.Error())
			}

			if body.Close() != nil {
				log.Printf("failed to close temp file: %s", err.Error())
			}
		}

		return nil
	}
}

func testAccCheckBucketExists(ctx context.Context, n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not Found: %s", n)
		}

		client, err := acctest.S3Client()
		if err != nil {
			return err
		}

		_, err = client.BucketsApi.HeadBucket(ctx, rs.Primary.Attributes["name"]).Execute()
		return err
	}
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

func testAccCheckBucketDeleteObjects(ctx context.Context, n string, keys ...string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs := s.RootModule().Resources[n]
		client, err := acctest.S3Client()
		if err != nil {
			return err
		}

		for _, key := range keys {
			_, _, err = client.ObjectsApi.DeleteObject(ctx, rs.Primary.Attributes["name"], key).Execute()
			if err != nil {
				return fmt.Errorf("DeleteObject error: %w", err)
			}
		}

		return nil
	}
}

func testAccCheckBucketAddObjects(ctx context.Context, n string, keys ...string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs := s.RootModule().Resources[n]

		client, err := acctest.S3Client()
		if err != nil {
			return err
		}

		for _, key := range keys {
			body, err := createTempFile("test", "test")
			if err != nil {
				return fmt.Errorf("createTempFile error: %w", err)
			}
			_, err = client.ObjectsApi.PutObject(ctx, rs.Primary.Attributes["name"], key).Body(body).Execute()
			if err != nil {
				return fmt.Errorf("PutObject error: %w", err)
			}

			err = os.Remove(body.Name())
			if err != nil {
				log.Printf("failed to remove temp file: %s", err.Error())
			}

			if body.Close() != nil {
				log.Printf("failed to close temp file: %s", err.Error())
			}
		}

		return nil
	}
}

func createTempFile(fileName, content string) (*os.File, error) {
	file, err := os.CreateTemp("", fileName)
	if err != nil {
		return nil, err
	}

	_, err = file.WriteString(content)
	if err != nil {
		return nil, err
	}

	f, err := os.Open(file.Name())
	if err != nil {
		return nil, err
	}

	return f, nil
}
