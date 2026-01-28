---
subcategory: "Event Streams for Apache Kafka"
layout: "ionoscloud"
page_title: "IonosCloud: ionoscloud_kafka_user_credentials"
sidebar_current: "docs-ionoscloud-datasource-kafka-user-credentials"
description: |-
  Gets information about Kafka users access credentials.
---

# ionoscloud_kafka_user_credentials

The **Kafka user credentials** data source can be used to retrieve access credentials for a specific user. 

> ⚠️  In order to avoid storing sensitive data in the state, the analogue [ephemeral resource](./../ephemerals/kafka_user_credentials.md) can be used.

## Example Usage

### By ID
```hcl
data "ionoscloud_kafka_user_credentials" "kafka_user_credentials_ds" {
  cluster_id = "kafka_cluster_id"
  id = "kafka_user_id"
  location = "kafka_cluster_location"
  timeouts = {
    read = "1m"
  }
}
```

### By name
```hcl
data "ionoscloud_kafka_user_credentials" "kafka_user_credentials_ds" {
  cluster_id = "kafka_cluster_id"
  username = "kafka_username"
  location = "kafka_cluster_location"
  timeouts = {
    read = "1m"
  }
}
```

## Argument reference
* `cluster_id` - (Required)[string] the ID of the Kafka cluster;
* `id` - (Optional)[string] the ID of the Kafka user, can be retrieved using [`ionoscloud_kafka_users` data source](kafka_users.md);
* `username` - (Optional)[string] the name of the Kafka user, can be retrieved using [`ionoscloud_kafka_users` data source](kafka_users.md);
* `location` - (Optional)[string] the location of the Kafka cluster, can be one of: `de/fra`, `de/fra/2`, `de/txl`, `fr/par`, `es/vit`, `gb/lhr`, `gb/bhx`, `us/las`, `us/mci`, `us/ewr`. If omitted, the default location will be used: `de/fra`;

## Attributes reference

The following attributes are returned by the data source:

* `id` - the ID of the user;
* `username` - the name of the user;
* `certificate_authority` - PEM for the certificate authority;
* `private_key` - PEM for the private key;
* `certificate` - PEM for the certificate;