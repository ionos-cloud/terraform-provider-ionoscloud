locals {
  postgres_prefix           = format("%s/%s", ionoscloud_server.example.nic[0].ips[0], "24")
  postgres_database_ip      = cidrhost(local.postgres_prefix, 1)
  postgres_database_ip_cidr = format("%s/%s", local.postgres_database_ip, "24")
}

resource "random_password" "cluster_password" {
  length  = 32
  special = false
}

resource "random_password" "user_password" {
  length  = 32
  special = false
}

resource "ionoscloud_pg_cluster" "example" {
  display_name         = "PostgreSQL Cluster"
  postgres_version     = 15
  location             = ionoscloud_datacenter.example.location
  backup_location      = "eu-central-2"
  instances            = 3
  cores                = 4
  ram                  = 2048
  storage_size         = 2048
  storage_type         = "SSD Standard"
  synchronization_mode = "ASYNCHRONOUS"

  connections {
    datacenter_id = ionoscloud_datacenter.example.id
    lan_id        = ionoscloud_lan.example.id
    cidr          = local.postgres_database_ip_cidr
  }

  maintenance_window {
    day_of_the_week = "Monday"
    time            = "09:00:00"
  }

  credentials {
    username = "username"
    password = random_password.cluster_password.result
  }
}

resource "ionoscloud_pg_user" "example_pg_user" {
  cluster_id = ionoscloud_pg_cluster.example.id
  username   = "exampleuser"
  password   = random_password.user_password.result
}

resource "ionoscloud_pg_database" "example_pg_database" {
  cluster_id = ionoscloud_pg_cluster.example.id
  name       = "exampledatabase"
  owner      = ionoscloud_pg_user.example_pg_user.username
}
