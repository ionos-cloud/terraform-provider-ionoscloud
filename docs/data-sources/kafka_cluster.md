---
subcategory: "Kafka"
layout: "ionoscloud"
page_title: "IonosCloud: cluster"
sidebar_current: "docs-resource-kafka-cluster"
description: |-
  Reads IonosCloud Kafka Cluster objects.
---

# ionoscloud_kafka_cluster

The **Kafka cluster data source** can be used to search for and return an existing Kafka Cluster.
You can provide a string for the name parameter which will be compared with provisioned Kafka clusters.
If a single match is found, it will be returned. If your search results in multiple matches, an error will be returned.
When this happens, please refine your search string so that it is specific enough to return only one result.

## Example Usage

### By ID
```hcl
data "ionoscloud_kafka_cluster" "example" {
  id       = <kafka_cluster_id>
}
```

### By Name

Needs to have the resource be previously created, or a depends_on clause to ensure that the resource is created before this data source is called.

```hcl
data "ionoscloud_kafka_cluster" "example" {
  name     = "kafka-cluster"
}
```

## Argument Reference
* `id` - (Optional) Id of an existing Kafka cluster that you want to search for.
* `name` - (Optional) Name of an existing Kafka cluster that you want to search for.

## Attributes Reference

The following attributes are returned by the datasource:

* `id` - UUID of the Kafka cluster
* `name` - The name of the Kafka cluster
* `version` - The version of the Kafka cluster
* `size` - The size of the Kafka cluster
* `broker_addresses` - List of Kafka broker addresses
* `server_address` - The server address of the Kafka cluster