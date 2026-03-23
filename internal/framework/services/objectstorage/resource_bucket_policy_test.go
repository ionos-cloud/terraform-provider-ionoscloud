//go:build all || objectstorage
// +build all objectstorage

package objectstorage_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework-jsontypes/jsontypes"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"

	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/internal/acctest"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/objectstorage"
)

func TestAccBucketPolicyResource(t *testing.T) {
	rName := acctest.GenerateRandomResourceName(bucketPrefix)
	name := "ionoscloud_s3_bucket_policy.test"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		CheckDestroy: testAccCheckBucketPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccBucketPolicyConfig_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(name, "bucket", rName),
					testAccCheckBucketPolicyData(PolicyJSON),
				),
			},
			{
				Config: testAccBucketPolicyConfig_updated(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(name, "bucket", rName),
					testAccCheckBucketPolicyData(policyJSONUpdated),
				),
			},
			{
				ResourceName:                         name,
				ImportStateId:                        rName,
				ImportState:                          true,
				ImportStateVerifyIdentifierAttribute: "bucket",
				ImportStateVerify:                    false,
			},
		},
	})
}

func testAccBucketPolicyConfig_basic(bucketName string) string {
	return fmt.Sprintf(`
resource "ionoscloud_s3_bucket" "test" {
  name = %[1]q
}

resource "ionoscloud_s3_bucket_policy" "test" {
 bucket = ionoscloud_s3_bucket.test.name
 policy = %[2]q
}
`, bucketName, PolicyJSON)
}

func testAccBucketPolicyConfig_updated(bucketName string) string {
	return fmt.Sprintf(`
resource "ionoscloud_s3_bucket" "test" {
  name = %[1]q
}

resource "ionoscloud_s3_bucket_policy" "test" {
 bucket = ionoscloud_s3_bucket.test.name
 policy = %[2]q
}
`, bucketName, policyJSONUpdated)
}

func testAccCheckBucketPolicyDestroy(s *terraform.State) error {
	client := acctest.NewTestBundleClientFromEnv().S3Client.GetBaseClient()

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "ionoscloud_s3_bucket_policy" {
			continue
		}

		if rs.Primary.Attributes["bucket"] != "" {
			_, apiResponse, err := client.PolicyApi.GetBucketPolicy(context.Background(), rs.Primary.Attributes["bucket"]).Execute()
			if apiResponse.HttpNotFound() {
				return nil
			}

			if err != nil {
				return fmt.Errorf("error checking for bucket policy")
			}

			return fmt.Errorf("bucket policy still exists")
		}
	}

	return nil
}

func testAccCheckBucketPolicyData(policy string) func(s *terraform.State) error {
	return func(s *terraform.State) error {
		for _, rs := range s.RootModule().Resources {
			if rs.Type != "ionoscloud_s3_bucket_policy" {
				continue
			}

			if rs.Primary.Attributes["policy"] != "" {
				n := jsontypes.NewNormalizedValue(rs.Primary.Attributes["policy"])
				v := jsontypes.NewNormalizedValue(policy)
				if eq, _ := n.StringSemanticEquals(context.Background(), v); !eq {
					return fmt.Errorf("Policy attribute not equal")
				}
			}
		}

		return nil
	}
}

// TestAccBucketPolicyResource_PrincipalFormats tests all supported ways to
// specify Principal in a bucket policy: flat array, object with string,
// object with array, and wildcard. Uses a single bucket with sequential steps.
func TestAccBucketPolicyResource_PrincipalFormats(t *testing.T) {
	rName := acctest.GenerateRandomResourceName(bucketPrefix)
	name := "ionoscloud_s3_bucket_policy.test"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acctest.PreCheck(t) },
		CheckDestroy:             testAccCheckBucketPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccBucketPolicyConfig_withPolicy(rName, policyPrincipalFlatArray),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(name, "bucket", rName),
					testAccCheckBucketPolicySemanticEqual(policyPrincipalFlatArray),
				),
			},
			{
				Config: testAccBucketPolicyConfig_withPolicy(rName, policyPrincipalObjectString),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(name, "bucket", rName),
					testAccCheckBucketPolicySemanticEqual(policyPrincipalObjectString),
				),
			},
			{
				Config: testAccBucketPolicyConfig_withPolicy(rName, policyPrincipalObjectArray),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(name, "bucket", rName),
					testAccCheckBucketPolicySemanticEqual(policyPrincipalObjectArray),
				),
			},
			{
				Config: testAccBucketPolicyConfig_withPolicy(rName, policyPrincipalWildcard),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(name, "bucket", rName),
					testAccCheckBucketPolicySemanticEqual(policyPrincipalWildcard),
				),
			},
		},
	})
}

func testAccBucketPolicyConfig_withPolicy(bucketName, policy string) string {
	return fmt.Sprintf(`
resource "ionoscloud_s3_bucket" "test" {
  name = %[1]q
}

resource "ionoscloud_s3_bucket_policy" "test" {
  bucket = ionoscloud_s3_bucket.test.name
  policy = %[2]q
}
`, bucketName, policy)
}

// testAccCheckBucketPolicySemanticEqual checks that the policy in state is
// semantically equal to expected, using PoliciesSemanticEqual which normalizes
// Principal format differences.
func testAccCheckBucketPolicySemanticEqual(expected string) func(s *terraform.State) error {
	return func(s *terraform.State) error {
		for _, rs := range s.RootModule().Resources {
			if rs.Type != "ionoscloud_s3_bucket_policy" {
				continue
			}
			actual := rs.Primary.Attributes["policy"]
			if actual == "" {
				return fmt.Errorf("policy attribute is empty")
			}
			equal, err := objectstorage.PoliciesSemanticEqual(actual, expected)
			if err != nil {
				return fmt.Errorf("failed to compare policies: %w", err)
			}
			if !equal {
				return fmt.Errorf("policy not semantically equal.\nGot:      %s\nExpected: %s", actual, expected)
			}
		}
		return nil
	}
}

// Principal as flat array: ["arn:..."]
const policyPrincipalFlatArray = `{
"Statement": [
{
  "Sid": "FlatArray",
  "Action": ["s3:GetObject"],
  "Effect": "Allow",
  "Principal": [
    "arn:aws:iam:::user/32112648:03a9a933-b732-425d-80e8-36a4da96d5a7"
  ],
  "Resource": [
    "arn:aws:s3:::tf-test-bucket-10189",
    "arn:aws:s3:::tf-test-bucket-10189/*"
  ]
}
],
"Version": "2012-10-17"
}`

// Principal as object with single string: {"AWS": "arn:..."}
const policyPrincipalObjectString = `{
"Statement": [
{
  "Sid": "ObjectString",
  "Action": ["s3:GetObject"],
  "Effect": "Allow",
  "Principal": {
    "AWS": "arn:aws:iam:::user/32112648:03a9a933-b732-425d-80e8-36a4da96d5a7"
  },
  "Resource": [
    "arn:aws:s3:::tf-test-bucket-10189",
    "arn:aws:s3:::tf-test-bucket-10189/*"
  ]
}
],
"Version": "2012-10-17"
}`

// Principal as object with array: {"AWS": ["arn:...", "arn:..."]}
const policyPrincipalObjectArray = `{
"Statement": [
{
  "Sid": "ObjectArray",
  "Action": ["s3:GetObject"],
  "Effect": "Allow",
  "Principal": {
    "AWS": [
      "arn:aws:iam:::user/32112648:03a9a933-b732-425d-80e8-36a4da96d5a7",
      "arn:aws:iam:::user/32112649:03a9a933-b732-425d-80e8-99988196d5a7"
    ]
  },
  "Resource": [
    "arn:aws:s3:::tf-test-bucket-10189",
    "arn:aws:s3:::tf-test-bucket-10189/*"
  ]
}
],
"Version": "2012-10-17"
}`

// Principal as wildcard: ["*"]
const policyPrincipalWildcard = `{
"Statement": [
{
  "Sid": "Wildcard",
  "Action": ["s3:GetObject"],
  "Effect": "Allow",
  "Principal": ["*"],
  "Resource": [
    "arn:aws:s3:::tf-test-bucket-10189",
    "arn:aws:s3:::tf-test-bucket-10189/*"
  ]
}
],
"Version": "2012-10-17"
}`

const PolicyJSON = `{
"Statement": [
{
  "Action": [
	"s3:DeleteObject",
	"s3:DeleteBucketWebsite"
  ],
  "Effect": "Allow",
  "Principal": [
	"arn:aws:iam:::user/32112648:03a9a933-b732-425d-80e8-36a4da96d5a7",
	"arn:aws:iam:::user/32112649:03a9a933-b732-425d-80e8-99988196d5a7"
  ],
  "Resource": [
	"arn:aws:s3:::tf-test-bucket-10189",
	"arn:aws:s3:::tf-test-bucket-10189/*"
  ],
  "Sid": "sid 1"
},
{
  "Action": [
	"s3:GetBucketTagging",
	"s3:GetBucketVersioning"
  ],
  "Condition": {
	  "DateGreaterThan": "2023-01-13 16:27:42Z",
	  "IpAddress": [
		"1.1.1.1",
		"2.2.2.2"
	  ],
	  "NotIpAddress": [
		"3.3.3.3",
		"4.4.4.4"
	  ]
	},
  "Effect": "Allow",
  "Principal": [
	"*"
  ],
  "Resource": [
	"arn:aws:s3:::tf-test-bucket-10189",
	"arn:aws:s3:::tf-test-bucket-10189/*"
  ],
  "Sid": "sid 2"
},
{
  "Action": [
	"s3:GetBucketTagging",
	"s3:GetBucketVersioning"
  ],
  "Condition": {
	  "DateLessThan": "2026-01-13 16:27:42Z",
	  "IpAddress": [
		"6.6.6.6",
		"7.7.7.7"
	  ]
  },
  "Effect": "Deny",
  "Principal": [
	"*"
  ],
  "Resource": [
	"arn:aws:s3:::tf-test-bucket-10189",
	"arn:aws:s3:::tf-test-bucket-10189/*"
  ]
}
],
"Version": "2008-10-17"
}`

const policyJSONUpdated = `{
"Statement": [
{
  "Action": [
	"s3:GetBucketTagging",
	"s3:GetBucketVersioning"
  ],
  "Condition": {
	  "DateLessThan": "2026-01-13 16:27:42Z",
	  "IpAddress": [
		"10.10.10.10",
		"11.11.11.11"
	  ]
  },
  "Effect": "Allow",
  "Principal": [
	"*"
  ],
  "Resource": [
	"arn:aws:s3:::tf-test-bucket-10189",
	"arn:aws:s3:::tf-test-bucket-10189/*"
  ]
}
],
"Version": "2008-10-17"
}`
