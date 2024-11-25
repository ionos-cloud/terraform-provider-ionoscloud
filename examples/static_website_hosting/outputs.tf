output "website_url_subdomain" {
  value = "https://${ionoscloud_s3_bucket.example.name}.s3.eu-central-3.ionoscloud.com/index.html"
}

output "website_url_subpath" {
  value = "https://s3.eu-central-3.ionoscloud.com/${ionoscloud_s3_bucket.example.name}/index.html"
}
