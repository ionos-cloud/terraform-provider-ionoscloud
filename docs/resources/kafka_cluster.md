---
subcategory: "Event Streams for Apache Kafka"
layout: "ionoscloud"
page_title: "IonosCloud: ionoscloud_kafka_cluster"
sidebar_current: "docs-resource-kafka-cluster"
description: |-
  Creates and manages IonosCloud Kafka Cluster objects.
---

# ionoscloud_kafka_cluster

Manages a [Kafka Cluster](https://docs.ionos.com/cloud/data-analytics/kafka/overview) on IonosCloud.

## Example Usage

This resource will create an operational Kafka Cluster. After this section completes, the provisioner can be called.

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
  location = "de/fra"
  version  = "3.8.0"
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
  version  = "3.8.0"
  size     = "S"
  connections {
    datacenter_id = ionoscloud_datacenter.example.id
    lan_id = ionoscloud_lan.example.id
    broker_addresses = local.kafka_cluster_broker_ips_cidr
  }
}
```

## Argument reference

* `id` - (Computed)[string] The UUID of the Kafka Cluster.
* `name` - (Required)[string] Name of the Kafka Cluster.
* `location` - (Optional)[string] The location of the Kafka Cluster. Possible values: `de/fra`, `de/txl`. If this is not set and if no value is provided for the `IONOS_API_URL` env var, the default `location` will be: `de/fra`.
* `version` - (Required)[string] Version of the Kafka Cluster. Possible values: `3.8.0`
* `size` - (Required)[string] Size of the Kafka Cluster. Possible values: `XS`, `S`
* `connections` - (Required) Connection information of the Kafka Cluster. Minimum items: 1, maximum items: 1.
    * `datacenter_id` - (Required)[string] The datacenter to connect your instance to.
    * `lan_id` - (Required)[string] The numeric LAN ID to connect your instance to.
    * `broker_addresses` - (Required)[list] IP addresses and subnet of cluster brokers. **Note** the following
      unavailable IP range: 10.224.0.0/11
* `broker_addresses` - (Computed)[list] IP address and port of cluster brokers.

> **⚠ NOTE:** `IONOS_API_URL_KAFKA` can be used to set a custom API URL for the kafka resource. `location` field needs to be empty, otherwise it will override the custom API URL. Setting `endpoint` or `IONOS_API_URL` does not have any effect.


## Import

Kafka Cluster can be imported using the `location` and `kafka cluster id`:

```shell
terraform import ionoscloud_kafka_cluster.mycluster location:kafka cluster uuid
```
