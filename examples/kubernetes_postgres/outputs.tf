output "kubernetes_service_ip" {
  value = kubernetes_service.app.status.0.load_balancer.0.ingress.0.ip
}

output "postgres_database_ip_cidr" {
  value = local.postgres_database_ip_cidr
}

output "postgres_database_dns_name" {
  value = ionoscloud_pg_cluster.example.dns_name
}
