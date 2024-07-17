---
subcategory: "Kafka"
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
resource "ionoscloud_kafka_cluster_topic" "kafka_cluster_topic" {
  cluster_id = <your_kafka_cluster_id>
  name                 = "kafka-cluster-topic"
  location             = <location_of_kafka_cluster>
  replication_factor   = 1
  number_of_partitions = 1
  retention_time       = 86400000
  segment_bytes        = 1073741824
}
```

## Argument reference

* `id` - (Computed)[string] The UUID of the Kafka Cluster Topic.
* `name` - (Required)[string] Name of the Kafka Cluster.
* `location` - (Required)[string] The location of the Kafka Cluster Topic. Possible values: `de/fra`, `de/txl`,
  `es/vit`,`gb/lhr`, `us/ewr`, `us/las`, `us/mci`, `fr/par`
* `cluster_id` - (Required)[string] ID of the Kafka Cluster that the topic belongs to.
* `replication_factor` - (Required)[int] The number of replicas of the topic. The replication factor determines how many
  copies of the topic are stored on different brokers. The replication factor must be less than or equal to the number
  of brokers in the Kafka Cluster. Minimum value: 1.
* `number_of_partitions` - (Required)[int] The number of partitions of the topic. Partitions allow for parallel
  processing of messages. The partition count must be greater than or equal to the replication factor. Minimum value: 1.
* `retention_time` - (Optional)[int] The time in milliseconds that a message is retained in the topic log. Messages
  older than the retention time are deleted. If value is 0, messages are retained indefinitely unless other retention is
  set. Default value: 0.
* `segment_bytes` - (Optional)[int] The maximum size in bytes that the topic log can grow to. When the log reaches this
  size, the oldest messages are deleted. If value is 0, messages are retained indefinitely unless other retention is
  set. Default value: 0.

## Import

Kafka Cluster Topic can be imported using the `kafka cluster topic id` and the `kafka cluster id`:

```shell
terraform import ionoscloud_kafka_cluster_topic.my_topic {kafka cluster uuid}:{kafka cluster topic uuid}
```
