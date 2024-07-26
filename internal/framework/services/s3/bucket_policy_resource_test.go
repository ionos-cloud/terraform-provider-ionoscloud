//go:build all || s3
// +build all s3

package s3_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework-jsontypes/jsontypes"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/internal/framework/services/s3"

	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/internal/acctest"
)

func TestAccBucketPolicyResource(t *testing.T) {
	rName := "tf-acctest-test-bucket-policy"
	name := "ionoscloud_bucket_policy.test"

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
resource "ionoscloud_bucket" "test" {
  bucket = %[1]q
}

resource "ionoscloud_bucket_policy" "test" {
 bucket = ionoscloud_bucket.test.bucket
 policy = %[2]q
}
`, bucketName, PolicyJSON)
}

func testAccBucketPolicyConfig_updated(bucketName string) string {
	return fmt.Sprintf(`
resource "ionoscloud_bucket" "test" {
  bucket = %[1]q
}

resource "ionoscloud_bucket_policy" "test" {
 bucket = ionoscloud_bucket.test.bucket
 policy = %[2]q
}
`, bucketName, policyJSONUpdated)
}

func testAccCheckBucketPolicyDestroy(s *terraform.State) error {
	client, err := acctest.S3Client()
	if err != nil {
		return err
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "ionoscloud_bucket_policy" {
			continue
		}

		if rs.Primary.Attributes["bucket"] != "" {
			_, err = s3.GetBucketPolicy(context.Background(), client, rs.Primary.Attributes["bucket"])
			if err != s3.ErrBucketPolicyNotFound {
				return err
			}
		}
	}

	return nil
}

func testAccCheckBucketPolicyData(policy string) func(s *terraform.State) error {
	return func(s *terraform.State) error {
		for _, rs := range s.RootModule().Resources {
			if rs.Type != "ionoscloud_bucket_policy" {
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
