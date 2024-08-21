//go:build all || s3
// +build all s3

package s3_test

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/internal/acctest"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils"
	"testing"
)

func TestAccS3ObjectsDataSource_basic(t *testing.T) {
	rName := acctest.GenerateRandomResourceName(bucketPrefix)
	dataSourceName := "data.ionoscloud_s3_objects.test"

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		PreventPostDestroyRefresh: true,
		Steps: []resource.TestStep{
			{
				Config: testAccObjectsDataSourceConfig_basic(rName, 1),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceName, "common_prefixes.#", "0"),
					resource.TestCheckResourceAttr(dataSourceName, "keys.#", "3"),
					resource.TestCheckResourceAttr(dataSourceName, "owners.#", "0"),
				),
			},
		},
	})
}

func TestAccS3ObjectsDataSource_prefixes(t *testing.T) {
	rName := acctest.GenerateRandomResourceName(bucketPrefix)
	dataSourceName := "data.ionoscloud_s3_objects.test"

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		PreventPostDestroyRefresh: true,
		Steps: []resource.TestStep{
			{
				Config: testAccObjectsDataSourceConfig_prefixes(rName, 1),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceName, "common_prefixes.#", "2"),
					resource.TestCheckResourceAttr(dataSourceName, "keys.#", "0"),
					resource.TestCheckResourceAttr(dataSourceName, "owners.#", "0"),
				),
			},
		},
	})
}

func TestAccS3ObjectsDataSource_encoded(t *testing.T) {
	rName := acctest.GenerateRandomResourceName(bucketPrefix)
	dataSourceName := "data.ionoscloud_s3_objects.test"

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		PreventPostDestroyRefresh: true,
		Steps: []resource.TestStep{
			{
				Config: testAccObjectsDataSourceConfig_encoded(rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceName, "common_prefixes.#", "0"),
					resource.TestCheckResourceAttr(dataSourceName, "keys.#", "1"),
					resource.TestCheckResourceAttr(dataSourceName, "keys.0", "prefix%2Fa%2Cb"),
					resource.TestCheckResourceAttr(dataSourceName, "owners.#", "0"),
				),
			},
		},
	})
}

func TestAccS3ObjectsDataSource_maxKeysSmall(t *testing.T) {
	rName := acctest.GenerateRandomResourceName(bucketPrefix)
	dataSourceName := "data.ionoscloud_s3_objects.test"

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		PreventPostDestroyRefresh: true,
		Steps: []resource.TestStep{
			{
				Config: testAccObjectsDataSourceConfig_maxKeysSmall(rName, 1, 5),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceName, "common_prefixes.#", "0"),
					resource.TestCheckResourceAttr(dataSourceName, "keys.#", "3"),
					resource.TestCheckResourceAttr(dataSourceName, "owners.#", "0"),
				),
			},
			{
				Config: testAccObjectsDataSourceConfig_maxKeysSmall(rName, 2, 5),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceName, "common_prefixes.#", "0"),
					resource.TestCheckResourceAttr(dataSourceName, "keys.#", "5"),
					resource.TestCheckResourceAttr(dataSourceName, "owners.#", "0"),
				),
			},
		},
	})
}

func TestAccS3ObjectsDataSource_maxKeysLarge(t *testing.T) {
	ctx := context.Background()
	rName := acctest.GenerateRandomResourceName(bucketPrefix)
	dataSourceName := "data.ionoscloud_s3_objects.test"
	var keys []string
	for i := 0; i < 1500; i++ {
		keys = append(keys, fmt.Sprintf("data%d", i))
	}

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		PreventPostDestroyRefresh: true,
		Steps: []resource.TestStep{
			{
				Config: testAccObjectsDataSourceConfig_maxKeysLarge(rName, 1002),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceName, "common_prefixes.#", "0"),
					resource.TestCheckResourceAttr(dataSourceName, "keys.#", "0"),
					resource.TestCheckResourceAttr(dataSourceName, "owners.#", "0"),
					testAccCheckBucketAddObjects(ctx, "ionoscloud_s3_bucket.test", keys...),
				),
			},
			{
				Config: testAccObjectsDataSourceConfig_maxKeysLarge(rName, 1002),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceName, "common_prefixes.#", "0"),
					resource.TestCheckResourceAttr(dataSourceName, "keys.#", "1002"),
					resource.TestCheckResourceAttr(dataSourceName, "owners.#", "0"),
				),
			},
		},
	})
}

func TestAccS3ObjectsDataSource_startAfter(t *testing.T) {
	rName := acctest.GenerateRandomResourceName(bucketPrefix)
	dataSourceName := "data.ionoscloud_s3_objects.test"

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		PreventPostDestroyRefresh: true,
		Steps: []resource.TestStep{
			{
				Config: testAccObjectsDataSourceConfig_startAfter(rName, 1, "prefix1/sub2/0"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceName, "common_prefixes.#", "0"),
					resource.TestCheckResourceAttr(dataSourceName, "keys.#", "1"),
					resource.TestCheckResourceAttr(dataSourceName, "owners.#", "0"),
				),
			},
		},
	})
}

func TestAccS3ObjectsDataSource_fetchOwner(t *testing.T) {
	rName := acctest.GenerateRandomResourceName(bucketPrefix)
	dataSourceName := "data.ionoscloud_s3_objects.test"

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		PreventPostDestroyRefresh: true,
		Steps: []resource.TestStep{
			{
				Config: testAccObjectsDataSourceConfig_owners(rName, 1),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceName, "common_prefixes.#", "0"),
					resource.TestCheckResourceAttr(dataSourceName, "keys.#", "3"),
					resource.TestCheckResourceAttr(dataSourceName, "owners.#", "3"),
				),
			},
		},
	})
}

func testAccObjectsDataSourceConfig_base(rName string, n int) string {
	return fmt.Sprintf(`
resource "ionoscloud_s3_bucket" "test" {
  name = %[1]q
}

resource "ionoscloud_s3_object" "test1" {
  count = %[2]d

  bucket  = ionoscloud_s3_bucket.test.name
  key     = "prefix1/sub1/${count.index}"
  content = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
}

resource "ionoscloud_s3_object" "test2" {
  count = %[2]d

  bucket  = ionoscloud_s3_bucket.test.name
  key     = "prefix1/sub2/${count.index}"
  content = "0123456789"
}

resource "ionoscloud_s3_object" "test3" {
  count = %[2]d

  bucket  = ionoscloud_s3_bucket.test.name
  key     = "prefix2/${count.index}"
  content = "abcdefghijklmnopqrstuvwxyz"
}
`, rName, n)
}

func testAccObjectsDataSourceConfig_basic(rName string, n int) string {
	return utils.ConfigCompose(testAccObjectsDataSourceConfig_base(rName, n), `
data "ionoscloud_s3_objects" "test" {
  bucket = ionoscloud_s3_bucket.test.name

  depends_on = [ionoscloud_s3_object.test1, ionoscloud_s3_object.test2, ionoscloud_s3_object.test3]
}
`)
}

func testAccObjectsDataSourceConfig_prefixes(rName string, n int) string {
	return utils.ConfigCompose(testAccObjectsDataSourceConfig_base(rName, n), `
data "ionoscloud_s3_objects" "test" {
  bucket    = ionoscloud_s3_bucket.test.name
  prefix    = "prefix1/"
  delimiter = "/"

  depends_on = [ionoscloud_s3_object.test1, ionoscloud_s3_object.test2, ionoscloud_s3_object.test3]
}
`)
}

func testAccObjectsDataSourceConfig_encoded(rName string) string {
	return fmt.Sprintf(`
resource "ionoscloud_s3_bucket" "test" {
  name = %[1]q
}

resource "ionoscloud_s3_object" "test" {
  bucket  = ionoscloud_s3_bucket.test.name
  key     = "prefix/a,b"
  content = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
}

data "ionoscloud_s3_objects" "test" {
  bucket        = ionoscloud_s3_bucket.test.name
  encoding_type = "url"

  depends_on = [ionoscloud_s3_object.test]
}
`, rName)
}

func testAccObjectsDataSourceConfig_maxKeysSmall(rName string, n, maxKeys int) string {
	return utils.ConfigCompose(testAccObjectsDataSourceConfig_base(rName, n), fmt.Sprintf(`
data "ionoscloud_s3_objects" "test" {
  bucket   = ionoscloud_s3_bucket.test.name
  max_keys = %[1]d

  depends_on = [ionoscloud_s3_object.test1, ionoscloud_s3_object.test2, ionoscloud_s3_object.test3]
}
`, maxKeys))
}

// Objects are added to the bucket outside this configuration.
func testAccObjectsDataSourceConfig_maxKeysLarge(rName string, maxKeys int) string {
	return fmt.Sprintf(`
resource "ionoscloud_s3_bucket" "test" {
  name        = %[1]q
  force_destroy = true
}

data "ionoscloud_s3_objects" "test" {
  bucket   = ionoscloud_s3_bucket.test.name
  max_keys = %[2]d
}
`, rName, maxKeys)
}

func testAccObjectsDataSourceConfig_startAfter(rName string, n int, startAfter string) string {
	return utils.ConfigCompose(testAccObjectsDataSourceConfig_base(rName, n), fmt.Sprintf(`
data "ionoscloud_s3_objects" "test" {
  bucket      = ionoscloud_s3_bucket.test.name
  start_after = %[1]q

  depends_on = [ionoscloud_s3_object.test1, ionoscloud_s3_object.test2, ionoscloud_s3_object.test3]
}
`, startAfter))
}

func testAccObjectsDataSourceConfig_owners(rName string, n int) string {
	return utils.ConfigCompose(testAccObjectsDataSourceConfig_base(rName, n), `
data "ionoscloud_s3_objects" "test" {
  bucket      = ionoscloud_s3_bucket.test.name
  fetch_owner = true

  depends_on = [ionoscloud_s3_object.test1, ionoscloud_s3_object.test2, ionoscloud_s3_object.test3]
}
`)
}
