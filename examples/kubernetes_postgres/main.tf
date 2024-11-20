terraform {
  backend "s3" {
    endpoint                    = "s3.eu-central-3.ionoscloud.com"
    region                      = "eu-central-3"
    bucket                      = "example-tf-state"
    key                         = "state/terraform.tfstate"
    force_path_style            = true
    skip_credentials_validation = true
    skip_region_validation      = true
  }

  required_providers {
    ionoscloud = {
      source  = "ionos-cloud/ionoscloud"
      version = "~> 6.6.1"
    }
    kubernetes = {
      source  = "hashicorp/kubernetes"
      version = "~> 2.0"
    }
    local = {
      source  = "hashicorp/local"
      version = "~> 2.0"
    }
  }
}

provider "ionoscloud" {
  s3_access_key = var.IONOS_S3_ACCESS_KEY
  s3_secret_key = var.IONOS_S3_SECRET_KEY
}

resource "ionoscloud_datacenter" "example" {
  name     = "Datacenter"
  location = "de/fra"
}

resource "ionoscloud_lan" "example" {
  datacenter_id = ionoscloud_datacenter.example.id
  public        = false
  name          = "Lan"
  lifecycle {
    create_before_destroy = true
  }
}

resource "random_password" "image_password" {
  length  = 32
  special = false
}

resource "ionoscloud_server" "example" {
  name              = "Jumphost"
  datacenter_id     = ionoscloud_datacenter.example.id
  cores             = 2
  ram               = 2048
  availability_zone = "AUTO"
  cpu_family        = "INTEL_SKYLAKE"
  image_name        = "ubuntu:latest"
  image_password    = random_password.image_password.result
  volume {
    size      = 10
    disk_type = "SSD Standard"
  }
  nic {
    lan  = ionoscloud_lan.example.id
    dhcp = true
  }
}
