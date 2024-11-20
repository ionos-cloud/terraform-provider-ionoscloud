resource "ionoscloud_k8s_cluster" "example" {
  name        = "k8sClusterExample"
  k8s_version = var.KUBERNETES_VERSION
  location    = ionoscloud_datacenter.example.location
  public      = true

  maintenance_window {
    day_of_the_week = "Monday"
    time            = "10:00:00Z"
  }
}

resource "ionoscloud_k8s_node_pool" "example" {
  name              = "k8sNodePoolExample"
  availability_zone = "AUTO"
  cpu_family        = "INTEL_SKYLAKE"
  datacenter_id     = ionoscloud_datacenter.example.id
  k8s_cluster_id    = ionoscloud_k8s_cluster.example.id
  k8s_version       = var.KUBERNETES_VERSION
  node_count        = 2
  cores_count       = 2
  ram_size          = 2048
  storage_size      = 10
  storage_type      = "SSD"

  lans {
    id   = ionoscloud_lan.example.id
    dhcp = true
  }

  maintenance_window {
    day_of_the_week = "Monday"
    time            = "11:00:00Z"
  }
}

data "ionoscloud_k8s_cluster" "app" {
  name = ionoscloud_k8s_cluster.example.name
}

provider "kubernetes" {
  host                   = data.ionoscloud_k8s_cluster.app.config[0].clusters[0].cluster.server
  token                  = data.ionoscloud_k8s_cluster.app.config[0].users[0].user.token
  cluster_ca_certificate = data.ionoscloud_k8s_cluster.app.config[0].clusters[0].cluster.certificate_authority_data
}

resource "kubernetes_secret" "registry_token" {
  type = "kubernetes.io/dockerconfigjson"
  metadata {
    name      = "docker-registry-token"
    namespace = "default"
  }

  data = {
    ".dockerconfigjson" = base64decode(var.IONOS_REGISTRY_TOKEN)
  }
}

resource "kubernetes_deployment" "app" {

  metadata {
    name      = "app"
    namespace = "default"
    labels    = { app = "app" }
  }
  spec {
    selector {
      match_labels = { app = "app" }
    }
    template {
      metadata {
        labels = { app = "app" }
      }
      spec {
        image_pull_secrets {
          name = kubernetes_secret.registry_token.metadata[0].name
        }
        container {
          name  = "app"
          image = "kicbase/echo-server:1.0"

          env {
            name  = "PORT"
            value = "8080"
          }
          env {
            name  = "LOG_HTTP_HEADERS"
            value = "true"
          }
          env {
            name  = "DB_CONN_STR"
            value = "postgres://${ionoscloud_pg_user.example_pg_user.username}:${ionoscloud_pg_user.example_pg_user.password}@${ionoscloud_pg_cluster.example.dns_name}:5432/${ionoscloud_pg_database.example_pg_database.name}"
          }
        }
      }
    }
  }
}

resource "kubernetes_service" "app" {
  metadata {
    name = "app"
  }
  spec {
    type     = "LoadBalancer"
    selector = kubernetes_deployment.app.spec[0].template[0].metadata[0].labels
    port {
      port        = 3000
      target_port = 8080
    }
  }
}
