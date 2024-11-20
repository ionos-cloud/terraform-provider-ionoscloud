terraform {
  required_providers {
    ionoscloud = {
      source  = "ionos-cloud/ionoscloud"
      version = "~> 6.6.1"
    }
  }
}

locals {
  cloud_init = <<-EOT
    #cloud-config
    package_update: true
    packages:
    - nginx

    write_files:
    - path: /etc/nginx/sites-enabled/backend
      permissions: '0644'
      content: |
        server {
            listen 80;
            server_name example.${data.ionoscloud_dns_zone.example.name};

            location /files {
                # Handle POST requests
                if ($request_method = POST) {
                    return 200 'This is a POST request to /files';
                }

                # Handle GET requests
                if ($request_method = GET) {
                    return 200 'This is a GET request to /files';
                }

                # Catch all for other methods on /files
                return 405 'Method not allowed on /files';
            }

            # Catch all other routes
            location / {
                return 200 'Default message for all other requests';
            }
        }

    runcmd:
    - rm /etc/nginx/sites-enabled/default
    - systemctl reload nginx
    EOT
}

data "ionoscloud_dns_zone" "example" {
  name = var.DOMAIN
}

resource "ionoscloud_auto_certificate_provider" "example" {
  name     = "Let's Encrypt"
  email    = var.CERT_REQUESTER
  location = "de/fra"
  server   = "https://acme-v02.api.letsencrypt.org/directory"
}

resource "ionoscloud_auto_certificate" "example" {
  provider_id   = ionoscloud_auto_certificate_provider.example.id
  name          = "example.${data.ionoscloud_dns_zone.example.name}"
  common_name   = "example.${data.ionoscloud_dns_zone.example.name}"
  location      = ionoscloud_auto_certificate_provider.example.location
  key_algorithm = "rsa4096"
}

resource "ionoscloud_apigateway" "example" {
  name = "example"
  logs = true
}

resource "ionoscloud_apigateway_route" "files_route" {
  gateway_id = ionoscloud_apigateway.example.id
  name       = "files-route"
  type       = "http"
  websocket  = false
  paths = [
    "/files",
  ]
  methods = [
    "GET",
    "POST",
    "HEAD",
    "OPTIONS",
  ]
  upstreams {
    scheme = "http"
    host   = ionoscloud_vcpu_server.example.primary_ip
    port   = 80
  }
}

resource "ionoscloud_apigateway_route" "default_route" {
  gateway_id = ionoscloud_apigateway.example.id
  name       = "default-route"
  type       = "http"
  websocket  = false
  paths = [
    "/",
  ]
  methods = [
    "GET",
    "POST",
    "HEAD",
    "OPTIONS",
  ]
  upstreams {
    scheme = "http"
    host   = ionoscloud_vcpu_server.example.primary_ip
    port   = 80
  }
}

resource "ionoscloud_dns_record" "example_ipv4" {
  zone_id = data.ionoscloud_dns_zone.example.id
  name    = "example"
  type    = "A"
  content = ionoscloud_cdn_distribution.example.public_endpoint_v4
  ttl     = 60
  enabled = true
}

resource "ionoscloud_dns_record" "example_ipv6" {
  zone_id = data.ionoscloud_dns_zone.example.id
  name    = "example"
  type    = "AAAA"
  content = ionoscloud_cdn_distribution.example.public_endpoint_v6
  ttl     = 60
  enabled = true
}

resource "ionoscloud_cdn_distribution" "example" {
  domain         = "example.${data.ionoscloud_dns_zone.example.name}"
  certificate_id = ionoscloud_auto_certificate.example.last_issued_certificate_id
  routing_rules {
    scheme = "https"
    prefix = "/"
    upstream {
      host             = ionoscloud_apigateway.example.public_endpoint
      caching          = true
      waf              = true
      rate_limit_class = "R100"
      sni_mode         = "origin"
      geo_restrictions {
        allow_list = ["DE"]
      }
    }
  }
}

resource "ionoscloud_datacenter" "example" {
  name     = "example"
  location = "de/fra"
}

resource "ionoscloud_lan" "example" {
  datacenter_id = ionoscloud_datacenter.example.id
  public        = true
}

resource "ionoscloud_vcpu_server" "example" {
  name           = "example"
  datacenter_id  = ionoscloud_datacenter.example.id
  cores          = 1
  ram            = 1024 * 2
  image_name     = "ubuntu:latest"
  image_password = random_password.server_image_password.result
  volume {
    size      = 10
    disk_type = "SSD Standard"
    user_data = base64encode(local.cloud_init)
  }
  nic {
    lan  = ionoscloud_lan.example.id
    dhcp = true
  }
}

resource "random_password" "server_image_password" {
  length  = 16
  special = false
}
