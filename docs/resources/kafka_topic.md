---
subcategory: "Event Streams for Apache Kafka"
layout: "ionoscloud"
page_title: "IonosCloud: ionoscloud_kafka_cluster_topic"
sidebar_current: "docs-resource-kafka-topic"
description: |-
  Creates and manages IonosCloud Kafka Cluster Topic objects.
---

# ionoscloud_kafka_cluster_topic

Manages a **Kafka Cluster Topic** on IonosCloud.

## Example Usage

This resource will create an operational Kafka Cluster Topic. After this section completes, the provisioner can be
called.

```hcl
# Basic example

resource "ionoscloud_datacenter" "example" {
  name     = "example-kafka-datacenter"
  location = "de/fra"
}

resource "ionoscloud_lan" "example" {
  datacenter_id = ionoscloud_datacenter.example.id
  public        = false
  name          = "example-kafka-lan"
}

resource "ionoscloud_kafka_cluster" "example" {
  name     = "example-kafka-cluster"
  location = ionoscloud_datacenter.example.location
  version  = "3.7.0"
  size     = "S"
  connections {
    datacenter_id = ionoscloud_datacenter.example.id
    lan_id = ionoscloud_lan.example.id
    broker_addresses = [
      "192.168.1.101/24",
      "192.168.1.102/24",
      "192.168.1.103/24"
    ]
  }
}

resource "ionoscloud_kafka_cluster_topic" "example" {
  cluster_id           = ionoscloud_kafka_cluster.example.id
  name                 = "kafka-cluster-topic"
  location             = ionoscloud_kafka_cluster.example.location
  replication_factor   = 1
  number_of_partitions = 1
  retention_time       = 86400000
  segment_bytes        = 1073741824
}
```

```hcl
# Complete example

resource "ionoscloud_datacenter" "example" {
  name     = "example-kafka-datacenter"
  location = "de/fra"
}

resource "ionoscloud_lan" "example" {
  datacenter_id = ionoscloud_datacenter.example.id
  public        = false
  name          = "example-kafka-lan"
}

resource "ionoscloud_server" "example" {
  name              = "example-kafka-server"
  datacenter_id     = ionoscloud_datacenter.example.id
  cores             = 1
  ram               = 2 * 1024
  availability_zone = "AUTO"
  image_name = "ubuntu:latest" # alias name
  image_password    = random_password.password.result
  volume {
    name      = "example-kafka-volume"
    size      = 6
    disk_type = "SSD Standard"
  }
  nic {
    lan  = ionoscloud_lan.example.id
    name = "example-kafka-nic"
    dhcp = true
  }
}

resource "random_password" "password" {
  length  = 16
  special = false
}

locals {
  prefix = format("%s/%s", ionoscloud_server.example.nic[0].ips[0], "24")
  server_net_index              = split(".", ionoscloud_server.example.nic[0].ips[0])[3]
  kafka_cluster_broker_ips      = [
    for i in range(local.server_net_index + 1, local.server_net_index + 4) :cidrhost(local.prefix, i)
  ]
  kafka_cluster_broker_ips_cidr = [for ip in local.kafka_cluster_broker_ips : format("%s/%s", ip, "24")]
}

resource "ionoscloud_kafka_cluster" "example" {
  name     = "example-kafka-cluster"
  location = ionoscloud_datacenter.example.location
  version  = "3.7.0"
  size     = "S"
  connections {
    datacenter_id    = ionoscloud_datacenter.example.id
    lan_id           = ionoscloud_lan.example.id
    broker_addresses = local.kafka_cluster_broker_ips_cidr
  }
}

resource "ionoscloud_kafka_cluster_topic" "example" {
  cluster_id           = ionoscloud_kafka_cluster.example.id
  name                 = "kafka-cluster-topic"
  location             = ionoscloud_kafka_cluster.example.location
  replication_factor   = 1
  number_of_partitions = 1
  retention_time       = 86400000
  segment_bytes        = 1073741824
}
```

## Argument reference

* `id` - (Computed)[string] The UUID of the Kafka Cluster Topic.
* `name` - (Required)[string] Name of the Kafka Cluster.
* `location` - (Optional)[string] The location of the Kafka Cluster Topic. Possible values: `de/fra`, `de/txl`. If this is not set and if no value is provided for the `IONOS_API_URL` env var, the default `location` will be: `de/fra`.
* `cluster_id` - (Required)[string] ID of the Kafka Cluster that the topic belongs to.
* `replication_factor` - (Optional)[int] The number of replicas of the topic. The replication factor determines how many
  copies of the topic are stored on different brokers. The replication factor must be less than or equal to the number
  of brokers in the Kafka Cluster. Minimum value: 1. Default value: 3.
* `number_of_partitions` - (Optional)[int] The number of partitions of the topic. Partitions allow for parallel
  processing of messages. The partition count must be greater than or equal to the replication factor. Minimum value: 1.
  Default value: 3.
* `retention_time` - (Optional)[int] This configuration controls the maximum time we will retain a log before we will
  discard old log segments to free up space. This represents an SLA on how soon consumers must read their data. If set
  to -1, no time limit is applied. Default value: 604800000.
* `segment_bytes` - (Optional)[int] This configuration controls the segment file size for the log. Retention and
  cleaning is always done a file at a time so a larger segment size means fewer files but less granular control over
  retention. Default value: 1073741824.

## Import

Kafka Cluster Topic can be imported using the `location`, `kafka cluster id` and the `kafka cluster topic id`:

```shell
terraform import ionoscloud_kafka_cluster_topic.my_topic location:kafka cluster uuid:kafka cluster topic uuid
```
