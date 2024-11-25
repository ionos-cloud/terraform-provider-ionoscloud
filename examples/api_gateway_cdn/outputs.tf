output "backend_server_ip_address" {
  value = ionoscloud_vcpu_server.example.primary_ip
}

output "apigw_public_endpoint" {
  value = ionoscloud_apigateway.example.public_endpoint
}

output "cdn_distribution_ipv4_address" {
  value = ionoscloud_cdn_distribution.example.public_endpoint_v4
}

output "cdn_distribution_ipv6_address" {
  value = ionoscloud_cdn_distribution.example.public_endpoint_v6
}

output "cdn_domain" {
  value = ionoscloud_cdn_distribution.example.domain
}
