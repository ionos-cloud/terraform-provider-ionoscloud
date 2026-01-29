---
subcategory: "Event Streams for Apache Kafka"
layout: "ionoscloud"
page_title: "IonosCloud: ionoscloud_kafka_users"
sidebar_current: "docs-ionoscloud-datasource-kafka-users"
description: |-
  Gets information about Kafka users.
---

# ionoscloud_kafka_users

The **Kafka users** data source can be used to retrieve information about users.

## Example Usage

```hcl
data "ionoscloud_kafka_users" "kafka_users_ds" {
  cluster_id = "kafka_cluster_id"
  location = "kafka_cluster_location"
  timeouts = {
    read = "1s"
  }
}
```

## Argument reference
* `cluster_id` - (Required)[string] the ID of the Kafka cluster;
* `location` - (Optional)[string] the location of the Kafka cluster, can be one of: `de/fra`, `de/fra/2`, `de/txl`, `fr/par`, `es/vit`, `gb/lhr`, `gb/bhx`, `us/las`, `us/mci`, `us/ewr`. If omitted, the default location will be used: `de/fra`;

## Attributes Reference

The following attributes are returned by the data source:

* `users` - the list of users, for each user inside the list, the following information is retrieved:
  * `id` - the ID of the user;
  * `username` - the name of the user;