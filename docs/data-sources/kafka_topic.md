---
subcategory: "Kafka"
layout: "ionoscloud"
page_title: "IonosCloud: ionoscloud_kafka_cluster_topic"
sidebar_current: "docs-datasource-kafka-topic"
description: |-
  Reads IonosCloud Kafka Cluster objects.
---

# ionoscloud_kafka_cluster_topic

The **Kafka topic data source** can be used to search for and return an existing Kafka Cluster Topic.
You can provide a string for the name parameter which will be compared with provisioned Kafka Cluster Topics.
If a single match is found, it will be returned. If your search results in multiple matches, an error will be returned.
When this happens, please refine your search string so that it is specific enough to return only one result.

## Example Usage

### By ID

```hcl
data "ionoscloud_kafka_cluster_topic" "example" {
  id = <your_kafka_cluster_topic_id>
cluster_id = <your_kafka_cluster_id>
location = <your_kafka_cluster_location>
}
```

### By Name

Needs to have the resource be previously created, or a depends_on clause to ensure that the resource is created before
this data source is called.

```hcl
data "ionoscloud_kafka_cluster_topic" "example" {
  name       = "kafka-cluster-topic"
  cluster_id = <your_kafka_cluster_id>
  location = <location_of_kafka_cluster>
}
```

## Argument Reference

* `id` - (Optional) ID of an existing Kafka Cluster Topic that you want to search for.
* `name` - (Optional) Name of an existing Kafka Cluster Topic that you want to search for.
* `cluster_id` - (Required) ID of the Kafka Cluster that the topic belongs to.
* `location` - (Required) The location of the Kafka Cluster Topic. Must be the same as the location of the Kafka
  Cluster. Possible values: `de/fra`, `de/txl`, `es/vit`,`gb/lhr`, `us/ewr`, `us/las`, `us/mci`, `fr/par`

## Attributes Reference

The following attributes are returned by the datasource:

* `id` - UUID of the Kafka Cluster Topic.
* `name` - The name of the Kafka Cluster Topic.
* `cluster_id` - The id of the Kafka Cluster that the topic belongs to.
* `replication_factor` - The number of replicas of the topic. The replication factor determines how many copies of the
  topic are stored on different brokers.
* `number_of_partitions` - The number of partitions of the topic. Partitions allow for parallel processing of messages.
* `retention_time` - The time in milliseconds that a message is retained in the topic log. Messages older than the
  retention time are deleted.
* `segment_bytes` - The maximum size in bytes that the topic log can grow to. When the log reaches this size, the oldest
  messages are deleted.
