---
subcategory: "Event Streams for Apache Kafka"
layout: "ionoscloud"
page_title: "IonosCloud: ionoscloud_kafka_user_credentials"
sidebar_current: "docs-ionoscloud-ephemeral-kafka-user-credentials"
description: |-
  Gets information about Kafka users access credentials.
---

# ionoscloud_kafka_user_credentials

The **Kafka user credentials** ephemeral can be used to retrieve access credentials for a specific user without storing sensitive data into the state.

## Example Usage

### By ID
```hcl
ephemeral "ionoscloud_kafka_user_credentials" "kafka_user_credentials_ephemeral" {
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
ephemeral "ionoscloud_kafka_user_credentials" "kafka_user_credentials_ephemeral" {
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

The following information is returned by the ephemeral resource:

* `id` - the ID of the user;
* `username` - the name of the user;
* `certificate_authority` - PEM for the certificate authority;
* `private_key` - PEM for the private key;
* `certificate` - PEM for the certificate;