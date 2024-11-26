# Table of contents

* [Introduction](README.md)
* [Changelog](../CHANGELOG.md)

## Terraform Registry

* [Terraform Registry](https://registry.terraform.io/providers/ionos-cloud/ionoscloud/latest)
* [Documentation](https://registry.terraform.io/providers/ionos-cloud/ionoscloud/latest/docs)

## Legal

* [Privacy policy](https://www.ionos.com/terms-gtc/terms-privacy/)
* [Imprint](https://www.ionos.de/impressum)


## API 

* Compute Engine
  * Resources
    * [Datacenter](../docs/resources/datacenter.md)
    * [Firewall Rule](../docs/resources/firewall.md)
    * [IpBlock](../docs/resources/ipblock.md)
    * [IP Failover](../docs/resources/ipfailover.md)
    * [Lan](../docs/resources/lan.md)
    * [LoadBalancer](../docs/resources/loadbalancer.md)
    * [Nic](../docs/resources/nic.md)
    * [Cross Connect](../docs/resources/private_crossconnect.md)
    * [Server](../docs/resources/server.md)
    * [VCPU Server](../docs/resources/vcpu_server.md)
    * [Cube Server](../docs/resources/cube_server.md)
    * [Snapshot](../docs/resources/snapshot.md)
    * [Volume](../docs/resources/volume.md)
    * [Server BootVolume Selection](../docs/resources/server_boot_device_selection.md)
  * Data Sources
    * [Datacenter](../docs/data-sources/datacenter.md)
    * [Firewall Rule](../docs/data-sources/firewall.md)
    * [Image](../docs/data-sources/image.md)
    * [IpBlock](../docs/data-sources/ipblock.md)
    * [IP Failover](../docs/data-sources/ipfailover.md)
    * [Lan](../docs/data-sources/lan.md)
    * [Location](../docs/data-sources/location.md)
    * [Nic](../docs/data-sources/nic.md)
    * [Cross Connect](../docs/data-sources/private_crossconnect.md)
    * [Server](../docs/data-sources/server.md)
    * [VCPU Server](../docs/data-sources/vcpu_server.md)
    * [Cube Server](../docs/data-sources/cube_server.md)
    * [Servers](../docs/data-sources/servers.md)
    * [Snapshot](../docs/data-sources/snapshot.md)
    * [Template](../docs/data-sources/template.md)
    * [Volume](../docs/data-sources/volume.md)

* Managed Kubernetes
  * Resources
    * [Kubernetes Cluster](../docs/resources/k8s_cluster.md)
    * [Kubernetes NodePool](../docs/resources/k8s_node_pool.md)
  * Data Sources
    * [Kubernetes Cluster](../docs/data-sources/k8s_cluster.md)
    * [Kubernetes NodePool](../docs/data-sources/k8s_node_pool.md)

* NAT Gateway
  * Resources
    * [NAT Gateway](../docs/resources/natgateway.md)
    * [NAT Gateway Rule](../docs/resources/natgateway_rule.md)
  * Data Sources
    * [NAT Gateway](../docs/data-sources/natgateway.md)
    * [NAT Gateway Rule](../docs/data-sources/natgateway_rule.md)

* Network Load Balancer
  * Resources
    * [Network Load Balancer](../docs/resources/networkloadbalancer.md)
    * [Network Load Balancer Forwarding Rule](../docs/resources/networkloadbalancer_forwardingrule.md)
  * Data Sources
    * [Network Load Balancer](../docs/data-sources/networkloadbalancer.md)
    * [Network Load Balancer Forwarding Rule](../docs/data-sources/networkloadbalancer_forwardingrule.md)

* Network Security Group
  * Resources
    * [Network Security Group](../docs/resources/nsg.md)
    * [Network Security Group Firewall Rule](../docs/resources/nsg_firewallrule.md)
    * [Datacenter Network Security Group Selection](../docs/resources/datacenter_nsg_selection.md)
  * Data Sources
    * [Network Security Group](../docs/data-sources/nsg.md)

* Managed Backup
  * Resources
    * [Backup Unit](../docs/resources/backup_unit.md)
  * Data Sources
    * [Backup Unit](../docs/data-sources/backup_unit.md)

* User Management
  * Resources
    * [Group](../docs/resources/group.md)
    * [Object Storage Key](../docs/resources/s3_key.md)
    * [Share](../docs/resources/share.md)
    * [User](../docs/resources/user.md)
  * Data Sources
    * [Group](../docs/data-sources/group.md)
    * [Resource](../docs/data-sources/resource.md)
    * [Object Storage Key](../docs/data-sources/s3_key.md)
    * [Share](../docs/data-sources/share.md)
    * [User](../docs/data-sources/user.md)

* Database as a Service - Postgres
  * Resources
    * [DBaaS Postgres Cluster](../docs/resources/dbaas_pgsql_cluster.md)
    * [DBaaS Postgres User](../docs/resources/dbaas_pgsql_user.md)
    * [DBaaS Postgres Database](../docs/resources/dbaas_pgsql_database.md)
  * Data Sources
    * [DBaaS Postgres Backup](../docs/data-sources/dbaas_pgsql_backups.md)
    * [DBaaS Postgres Cluster](../docs/data-sources/dbaas_pgsql_cluster.md)
    * [DBaaS Postgres Versions](../docs/data-sources/dbaas_pgsql_versions.md)
    * [DBaaS Postgres User](../docs/data-sources/dbaas_pgsql_user.md)
    * [DBaaS Postgres Database](../docs/data-sources/dbaas_pgsql_database.md)
    * [DBaaS Postgres Databases](../docs/data-sources/dbaas_pgsql_databases.md)

* Database as a Service - MongoDB
  * Resources
    * [DBaaS MongoDB Cluster](../docs/resources/dbaas_mongo_cluster.md)
    * [DBaaS MongoDB User](../docs/resources/dbaas_mongo_user.md)
  * Data Sources
    * [DBaaS MongoDB Cluster](../docs/data-sources/dbaas_mongo_cluster.md)
    * [DBaaS MongoDB Template](../docs/data-sources/dbaas_mongo_template.md)
    * [DBaaS MongoDB User](../docs/data-sources/dbaas_mongo_user.md)

* Database as a Service - MariaDB
  * Resources
    * [DBaaS MariaDB Cluster](../docs/resources/dbaas_mariadb_cluster.md)
  * Data Sources
    * [DBaaS MariaDB Cluster](../docs/data-sources/dbaas_mariadb_cluster.md)
    * [DBaaS MariaDB Backups](../docs/data-sources/dbaas_mariadb_backups.md)

* Database as a Service - InMemory-DB
  * Resources
    * [DBaaS InMemoryDB ReplicaSet](../docs/resources/dbaas_inmemorydb_replica_set.md)
  * Data sources
    * [DBaaS InMemoryDB ReplicaSet](../docs/data-sources/dbaas_inmemorydb_replica_set.md)
    * [DBaaS InMemoryDB Snapshot](../docs/data-sources/dbaas_inmemorydb_snapshot.md)

* Application Load Balancer
  * Resources
    * [Application Load Balancer](../docs/resources/application_loadbalancer.md)
    * [Application Load Balancer Forwarding Rule](../docs/resources/application_loadbalancer_forwardingrule.md)
    * [Target Group](../docs/resources/target_group.md)
  * Data Sources
    * [Application Load Balancer](../docs/data-sources/application_loadbalancer.md)
    * [Application Load Balancer Forwarding Rule](../docs/data-sources/application_loadbalancer_forwardingrule.md)
    * [Target Group](../docs/data-sources/target_group.md)

* Container Registry
  * Resources
    * [Container Registry](../docs/resources/container_registry.md)
    * [Container Registry Token](../docs/resources/container_registry_token.md)
  * Data Sources
    * [Container Registry](../docs/data-sources/container_registry.md)
    * [Container Registry Token](../docs/data-sources/container_registry_token.md)
    * [Container Registry Locations](../docs/data-sources/container_registry_locations.md)

* Data Platform
  * Resources
    * [Data Platform Cluster](../docs/resources/dataplatform_cluster.md)
    * [Data Platform Node Pool](../docs/resources/dataplatform_node_pool.md)
  * Data Sources
    * [Data Platform Cluster](../docs/data-sources/dataplatform_cluster.md)
    * [Data Platform Node Pool](../docs/data-sources/dataplatform_node_pool.md)
    * [Data Platform Node Pools](../docs/data-sources/dataplatform_node_pools.md)
    * [Data Platform Versions](../docs/data-sources/dataplatform_versions.md)

* Certificate Manager
  * Resources
    * [Certificate](../docs/resources/certificate_manager_certificate.md)
    * [Provider](../docs/resources/certificate_manager_provider.md)
    * [Auto-certificate](../docs/resources/certificate_manager_auto_certificate.md)
  * Data Sources
    * [Certificate](../docs/data-sources/certificate_manager_certificate.md)
    * [Provider](../docs/data-sources/certificate_manager_provider.md)
    * [Auto-certificate](../docs/data-sources/certificate_manager_auto_certificate.md)

* Cloud DNS
  * Resources
    * [DNS Record](../docs/resources/dns_record.md)
    * [DNS Zone](../docs/resources/dns_zone.md)
  * Data Sources
    * [DNS Record](../docs/data-sources/dns_record.md)
    * [DNS Zone](../docs/data-sources/dns_zone.md)

* Logging Service
  * Resources
    * [Logging Pipeline](../docs/resources/logging_pipeline.md)
  * Data Sources
    * [Logging Pipeline](../docs/data-sources/logging_pipeline.md)

* Network File Storage
  * Resources
    * [Network File Storage Cluster](../docs/resources/nfs_cluster.md)
    * [Network File Storage Share](../docs/resources/nfs_share.md)
  * Data Sources
    * [Network File Storage Cluster](../docs/data-sources/nfs_cluster.md)
    * [Network File Storage Share](../docs/data-sources/nfs_share.md)

* Object Storage
  * Resources
    * [Bucket](../docs/resources/s3_bucket.md)
    * [Bucket Policy](../docs/resources/s3_bucket_policy.md)
    * [Object](../docs/resources/s3_object.md)
    * [Bucket Public Access Block](../docs/resources/s3_bucket_public_access_block)
    * [Bucket Website Configuration](../docs/resources/s3_bucket_website_configuration.md)
    * [Bucket CORS Configuration](../docs/resources/s3_bucket_cors_configuration.md)
    * [Bucket Lifecycle Configuration](../docs/resources/s3_bucket_lifecycle_configuration.md)
    * [Bucket Object Lock Configuration](../docs/resources/s3_bucket_object_lock_configuration.md)
    * [Bucket Versioning Configuration](../docs/resources/s3_bucket_versioning.md)
    * [Bucket Server Side Encryption Configuration](../docs/resources/s3_bucket_server_side_encryption_configuration.md)
    * [Object Copy](../docs/resources/s3_object_copy.md)
  * Data Sources
    * [Bucket](../docs/data-sources/s3_bucket.md)
    * [Bucket Policy](../docs/data-sources/s3_bucket_policy.md)
    * [Object](../docs/data-sources/s3_object.md)
    * [Objects](../docs/data-sources/s3_objects.md)

* Object Storage Management
  * Resources
    * [Access Key](../docs/resources/object_storage_accesskey.md)
  * Data Sources
    * [Access Key](../docs/data-sources/object_storage_accesskey.md)
    * [Region](../docs/data-sources/object_storage_region.md)
  
* CDN
  * Resources
    * [CDN Distribution](../docs/resources/cdn_distribution.md)
  * Data Sources
    * [CDN Distribution](../docs/data-sources/cdn_distribution.md)

* API Gateway
  * Resources
    * [API Gateway](../docs/resources/apigateway.md)
    * [API Gateway Route](../docs/resources/apigateway_route.md)
  * Data Sources
    * [API Gateway](../docs/data-sources/apigateway.md)
    * [API Gateway Route](../docs/data-sources/apigateway_route.md)

* VPN
  * Resources
    * [VPN IPSEC Gateway](../docs/resources/vpn_ipsec_gateway.md)
    * [VPN IPSEC Tunnel](../docs/resources/vpn_ipsec_tunnel.md)
    * [VPN Wireguard Gateway](../docs/resources/vpn_wireguard_gateway.md)
    * [VPN Wireguard Peer](../docs/resources/vpn_wireguard_peer.md)
  * Data Sources
    * [VPN IPSEC Gateway](../docs/data-sources/vpn_ipsec_gateway.md)
    * [VPN IPSEC Tunnel](../docs/data-sources/vpn_ipsec_tunnel.md)
    * [VPN Wireguard Gateway](../docs/data-sources/vpn_wireguard_gateway.md)
    * [VPN Wireguard Peer](../docs/data-sources/vpn_wireguard_peer.md)

* Event Streams for Apache Kafka
  * Resources
    * [Kafka Cluster](../docs/resources/kafka_cluster.md)
    * [Kafka Topic](../docs/resources/kafka_topic.md)
  * Data Sources
    * [Kafka Cluster](../docs/data-sources/kafka_cluster.md)
    * [Kafka Topic](../docs/data-sources/kafka_topic.md)
