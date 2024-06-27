---
subcategory: "Kafka"
layout: "ionoscloud"
page_title: "IonosCloud: cluster"
sidebar_current: "docs-resource-kafka-cluster"
description: |-
  Creates and manages IonosCloud Kafka Cluster objects.
---

# ionoscloud_kafka_cluster

Manages a **Kafka Cluster** on IonosCloud.

## Example Usage

This resource will create an operational kafka cluster. After this section completes, the provisioner can be called.

```hcl
resource "ionoscloud_kafka_cluster" "kafka_cluster" {
  name    = "kafka-cluster"
  version = "3.7.0"
  size    = "S"
  connections {
    datacenter_id = <your_datacenter_id>
    lan_id = <your_lan_id>
    cidr = "192.168.1.100/24"
    broker_addresses = [
      "192.168.1.101/24",
      "192.168.1.102/24",
      "192.168.1.103/24"
    ]
  }
}
```

## Argument reference

* `name` - (Required)[string] Name of the Kafka Cluster.
* `version` - (Required)[string] Version of the Kafka Cluster. Possible values: `3.7.0`
* `size` - (Required)[string] Size of the Kafka Cluster. Possible values: `S`
* `connections` - (Required) Connection information of the Kafka cluster. Minimum items: 1, maximum items: 1.
    * `datacenter_id` - (Required)[string] The datacenter to connect your instance to.
    * `lan_id` - (Required)[string] The numeric LAN ID to connect your instance to.
    * `cidr` - (Required)[string] The IP and subnet for your instance. **Note** the following unavailable IP range:
      10.244.0.0/11
    * `broker_addresses` - (Required)[list] IP addresses and subnet of cluster brokers. **Note** the following
      unavailable IP range: 10.224.0.0/11

## Import

Kafka Cluster can be imported using the `kafka cluster id`:

```shell
terraform import ionoscloud_kafka_cluster.mycluster {kafka cluster uuid}
```
