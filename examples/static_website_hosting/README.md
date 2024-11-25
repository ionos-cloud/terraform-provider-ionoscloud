# Static Website Hosting with Object Storage Bucket and CDN + WAF

This Terraform manifest provisions a static website using IONOS Cloud's Object Storage service. The steps involved are:

- Create S3 Bucket: An S3 bucket is created in the eu-central-3 region.
- Configure Website: Set index and error documents for the website.
- Set Bucket Policy: Grant public read access to all objects in the bucket.
- Configure CDN: Set up a CDN for the bucket to enhance content delivery, with caching and geo-restrictions for access from Germany.
