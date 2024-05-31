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
  name = "kafka-cluster"
  version = "3.5.1"
  size = "S"
}
```

## Argument reference

* `name` - (Required)[string] Name of the Kafka Cluster.
* `version` - (Required)[string] Version of the Kafka Cluster.
* `size` - (Required)[string] Size of the Kafka Cluster. Possible values: `S`

## Import

Kafka Cluster can be imported using the `kafka cluster id`:

```shell
terraform import ionoscloud_kafka_cluster.mycluster {kafka cluster uuid}
```
