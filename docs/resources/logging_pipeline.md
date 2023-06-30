---
subcategory: "Logging Service"
layout: "ionoscloud"
page_title: "IonosCloud: ionoscloud_logging_pipeline"
sidebar_current: "docs-resource-logging_pipeline"
description: |-
  Creates and manages Logging pipeline objects.
---

# ionoscloud_logging_pipeline

⚠️ **Note:** Logging Service is currently in the Early Access (EA) phase.
We recommend keeping usage and testing to non-production critical applications.
Please contact your sales representative or support for more information.

Manages a **Logging pipeline**.

## Usage example

```hcl
resource "ionoscloud_logging_pipeline" "example" {
  name = "pipelineexample"
  log {
    source = "kubernetes"
    tag = "tagexample"
    protocol = "http"
    destinations {
      type = "loki"
      retention_in_days = 7
    }
  }
  log {
    source = "kubernetes"
    tag = "anothertagexample"
    protocol = "tcp"
    destinations {
      type = "loki"
      retention_in_days = 7
    }
  }
}
```

## Argument reference

* `name` - (Required)[string] The name of the Logging pipeline.
* `log` - (Required)[list] Pipeline logs, a list that contains elements with the following structure:
  * `source` - (Required)[string] The source parser to be used.
  * `tag` - (Required)[string] The tag is used to distinguish different pipelines. Must be unique amongst the pipeline's array items.
  * `protocol` - (Required)[string] "Protocol to use as intake. Possible values are: http, tcp."
  * `public` - (Computed)[bool]
  * `destinations` - (Optional)[list] The configuration of the logs datastore, a list that contains elements with the following structure:
    * `type` - (Optional)[string] The internal output stream to send logs to.
    * `retention_in_days` - (Optional)[int] Defines the number of days a log record should be kept in loki. Works with loki destination type only.

## Import

In order to import a Logging pipeline, you can define an empty Logging pipeline resource in the plan:

```hcl
resource "ionoscloud_logging_pipeline" "example" {
  
}
```

The resource can be imported using the `pipeline_id`, for example:

```shell
terraform import ionoscloud_logging_pipeline.example {pipeline_id}
```