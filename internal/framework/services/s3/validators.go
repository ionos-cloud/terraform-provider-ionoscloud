package s3

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

func bucketPolicyStatementActionValidators() []validator.List {
	actions := []string{"s3:*", "s3:AbortMultipartUpload", "s3:CreateBucket", "s3:DeleteBucketPolicy", "s3:DeleteBucket",
		"s3:DeleteBucketWebsite", "s3:DeleteObject", "s3:DeleteObjectVersion", "s3:DeleteBucketPublicAccessBlock", "s3:DeleteReplicationConfiguration",
		"s3:GetAccelerateConfiguration", "s3:GetBucketAcl", "s3:GetBucketCORS", "s3:GetBucketLocation", "s3:GetBucketLogging",
		"s3:GetBucketNotification", "s3:GetBucketPolicy", "s3:GetBucketRequestPayment", "s3:GetBucketTagging", "s3:GetBucketVersioning",
		"s3:GetBucketWebsite", "s3:GetLifecycleConfiguration", "s3:GetObjectAcl", "s3:GetObject", "s3:GetObjectTorrent",
		"s3:GetObjectVersionAcl", "s3:GetObjectVersion", "s3:GetObjectVersionTorrent", "s3:GetBucketPublicAccessBlock", "s3:GetReplicationConfiguration",
		"s3:ListAllMyBuckets", "s3:ListBucketMultipartUploads", "s3:ListBucket", "s3:ListBucketVersions", "s3:ListMultipartUploadParts",
		"s3:PutAccelerateConfiguration", "s3:PutBucketAcl", "s3:PutBucketCORS", "s3:PutBucketLogging", "s3:PutBucketNotification", "s3:PutBucketPolicy",
		"s3:PutBucketRequestPayment", "s3:PutBucketTagging", "s3:PutBucketVersioning", "s3:PutBucketWebsite", "s3:PutLifecycleConfiguration",
		"s3:PutBucketPublicAccessBlock", "s3:PutObjectAcl", "s3:PutObject", "s3:PutObjectVersionAcl", "s3:PutReplicationConfiguration", "s3:RestoreObject"}

	return []validator.List{listvalidator.ValueStringsAre(stringvalidator.OneOf(actions...))}
}
