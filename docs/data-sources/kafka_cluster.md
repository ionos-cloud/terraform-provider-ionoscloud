---
subcategory: "Event Streams for Apache Kafka"
layout: "ionoscloud"
page_title: "IonosCloud: ionoscloud_kafka_cluster"
sidebar_current: "docs-datasource-kafka-cluster"
description: |-
  Reads IonosCloud Kafka Cluster objects.
---

# ionoscloud_kafka_cluster

The **Kafka Cluster data source** can be used to search for and return an existing Kafka Cluster.
You can provide a string for the name parameter which will be compared with provisioned Kafka Clusters.
If a single match is found, it will be returned. If your search results in multiple matches, an error will be returned.
When this happens, please refine your search string so that it is specific enough to return only one result.

## Example Usage

### By ID

```hcl
data "ionoscloud_kafka_cluster" "example" {
  id = "your_kafka_cluster_id"
  location = "location_of_kafka_cluster"
}
```

### By Name

Needs to have the resource be previously created, or a depends_on clause to ensure that the resource is created before
this data source is called.

```hcl
data "ionoscloud_kafka_cluster" "example" {
  name     = "kafka-cluster"
  location = "location_of_kafka_cluster"
}
```

## Argument Reference

* `id` - (Optional) ID of an existing Kafka Cluster that you want to search for.
* `name` - (Optional) Name of an existing Kafka Cluster that you want to search for.
* `location` - (Required) The location of the Kafka Cluster. Possible values: `de/fra`, `de/fra/2`, `de/txl`, `fr/par`, `es/vit`, `gb/lhr`, `gb/bhx`, `us/las`, `us/mci`, `us/ewr`.

## Attributes Reference

The following attributes are returned by the datasource:

* `id` - UUID of the Kafka Cluster.
* `name` - The name of the Kafka Cluster.
* `version` - The version of the Kafka Cluster.
* `size` - The size of the Kafka Cluster.
* `connections` - Connection information of the Kafka Cluster. Minimum items: 1, maximum items: 1.
    * `datacenter_id` - The datacenter that your instance is connected to.
    * `lan_id` - The numeric LAN ID your instance is connected to.
    * `broker_addresses` - IP addresses and subnet of cluster brokers.
* `broker_addresses` - IP address and port of cluster brokers.
