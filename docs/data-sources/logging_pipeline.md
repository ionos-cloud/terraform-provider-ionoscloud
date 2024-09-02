---
subcategory: "Logging Service"
layout: "ionoscloud"
page_title: "IonosCloud: ionoscloud_logging_pipeline"
sidebar_current: "docs-resource-logging_pipeline"
description: |-
  Get information on a Logging pipeline.
---

# ionoscloud_logging_pipeline

The **Logging pipeline** datasource can be used to search for and return an existing Logging pipeline.
If a single match is found, it will be returned. If your search results in multiple matches, an error will be returned.

> ⚠️  Only tokens are accepted for authorization in the **logging_pipeline** data source. Please ensure you are using tokens as other methods will not be valid.

## Example Usage

### By ID
```hcl
data "ionoscloud_logging_pipeline" "example" {
  location = "de/txl"
  id = <pipeline_id>
}
```

### By name
```hcl
data "ionoscloud_logging_pipeline" "example" {
  location = "de/txl"
  name = "pipeline_name"
}
```

## Argument reference
* `location` - (Optional)[string] The location of the Logging pipeline. Default: `de/txl`. One of `de/fra`, `de/txl`, `gb/lhr`, `es/vit`, `fr/par`.
* `id` - (Optional)[string] The ID of the Logging pipeline you want to search for.
* `name` - (Optional)[string] The name of the Logging pipeline you want to search for.

Either `id` or `name` must be provided. If none, or both are provided, the datasource will return an error.

## Attributes Reference

The following attributes are returned by the datasource:

* `id` - The UUID of the Logging pipeline.
* `name` - The name of the Logging pipeline.
* `grafana_address` - The address of the client's grafana instance.
* `log` - [list] Pipeline logs, a list that contains elements with the following structure:
  * `source` - [string] The source parser to be used.
  * `tag` - [string] The tag is used to distinguish different pipelines. Must be unique amongst the pipeline's array items.
  * `protocol` - [string] "Protocol to use as intake. Possible values are: http, tcp."
  * `public` - [bool]
  * `destinations` - [list] The configuration of the logs datastore, a list that contains elements with the following structure:
    * `type` - [string] The internal output stream to send logs to.
    * `retention_in_days` - [int] Defines the number of days a log record should be kept in loki. Works with loki destination type only.