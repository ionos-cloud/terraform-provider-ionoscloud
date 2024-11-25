terraform {
  required_providers {
    ionoscloud = {
      source  = "ionos-cloud/ionoscloud"
      version = "~> 6.6.1"
    }
  }
}

provider "ionoscloud" {
  s3_access_key = var.IONOS_S3_ACCESS_KEY
  s3_secret_key = var.IONOS_S3_SECRET_KEY
}

resource "ionoscloud_s3_bucket" "example" {
  name   = "example-website"
  region = "eu-central-3"
}

resource "ionoscloud_s3_bucket_website_configuration" "example" {
  bucket = ionoscloud_s3_bucket.example.name

  index_document {
    suffix = "index.html"
  }

  error_document {
    key = "error.html"
  }
}

resource "ionoscloud_s3_bucket_policy" "example" {
  bucket = ionoscloud_s3_bucket.example.name

  policy = jsonencode({
    Version = "2012-10-17",
    Statement = [
      {
        Sid       = "PublicRead",
        Effect    = "Allow",
        Principal = ["*"],
        Action    = ["s3:GetObject"],
        Resource  = ["arn:aws:s3:::${ionoscloud_s3_bucket.example.name}/*"]
      }
    ]
  })
}

resource "ionoscloud_s3_object" "index_document" {
  bucket       = ionoscloud_s3_bucket.example.name
  key          = "index.html"
  source       = "index.html"
  content_type = "text/html"
}

resource "ionoscloud_s3_object" "error_document" {
  bucket       = ionoscloud_s3_bucket.example.name
  key          = "error.html"
  source       = "error.html"
  content_type = "text/html"
}

resource "ionoscloud_cdn_distribution" "example" {
  domain = "${ionoscloud_s3_bucket.example.name}.s3.eu-central-3.ionoscloud.com"
  routing_rules {
    scheme = "http"
    prefix = "/"
    upstream {
      host             = "${ionoscloud_s3_bucket.example.name}.s3.eu-central-3.ionoscloud.com"
      rate_limit_class = "R100"
      sni_mode         = "distribution"
      caching          = true
      waf              = true
      geo_restrictions {
        allow_list = ["DE"]
      }
    }
  }
}
