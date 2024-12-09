---
subcategory: "Logging Service"
layout: "ionoscloud"
page_title: "IonosCloud: ionoscloud_logging_pipeline"
sidebar_current: "docs-resource-logging_pipeline"
description: |-
  Creates and manages Logging pipeline objects.
---

# ionoscloud_logging_pipeline

Manages a **Logging pipeline**.

> ⚠️  Only tokens are accepted for authorization in the **logging_pipeline** resource. Please ensure you are using tokens as other methods will not be valid.

## Usage example

```hcl
resource "ionoscloud_logging_pipeline" "example" {
  location = "es/vit"
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

For re-usability, an array of **logs** can be defined in a **tfvars** file or inside the terraform
plan, and used as presented below:

The content inside **vars.tfvars** file:

```hcl
logs = [
  {
    source = "kubernetes"
    tag = "firstlog"
    protocol = "http"
    destinations = {
      type = "loki"
      retention_in_days = 7
    }},
    {
    source = "docker"
    tag = "secondlog"
    protocol = "tcp"
    destinations = {
      type = "loki"
      retention_in_days = 14
    }}]
```

The content inside the tf plan:

```hcl
variable "logs" {
  description = "logs"
  type        = list(object({
    source = string
    tag = string
    protocol = string
    destinations = object({
      type = string
      retention_in_days = number
    } )}))
}

resource "ionoscloud_logging_pipeline" "example" {
  location = "es/vit"
  name = "examplepipeline"
  dynamic "log" {
    for_each = var.logs
    content {
      source = log.value["source"]
      tag = log.value["tag"]
      protocol = log.value["protocol"]
      destinations {
        type = log.value["destinations"]["type"]
        retention_in_days = log.value["destinations"]["retention_in_days"]
      }
    }
  }
}
```
The configuration can then be applied using the following commands:

```shell
terraform plan -var-file="vars.tfvars"
terraform apply -var-file="vars.tfvars"
```

## Argument reference

* `location` - (Optional)[string] The location of the Logging pipeline. Default: `de/txl` One of `de/fra`, `de/txl`, `gb/lhr`, `es/vit`, `fr/par`. If this is not set and if no value is provided for the `IONOS_API_URL` env var, the default `location` will be: `de/fra`.
* `name` - (Required)[string] The name of the Logging pipeline.
* `grafana_address` - (Computed)[string] The address of the client's grafana instance.
* `log` - (Required)[list] Pipeline logs, a list that contains elements with the following structure:
  * `source` - (Required)[string] The source parser to be used.
  * `tag` - (Required)[string] The tag is used to distinguish different pipelines. Must be unique amongst the pipeline's array items.
  * `protocol` - (Required)[string] "Protocol to use as intake. Possible values are: http, tcp."
  * `public` - (Computed)[bool]
  * `destinations` - (Optional)[list] The configuration of the logs datastore, a list that contains elements with the following structure:
    * `type` - (Optional)[string] The internal output stream to send logs to.
    * `retention_in_days` - (Optional)[int] Defines the number of days a log record should be kept in loki. Works with loki destination type only. Can be one of: 7, 14, 30.

## Import

In order to import a Logging pipeline, you can define an empty Logging pipeline resource in the plan:

```hcl
resource "ionoscloud_logging_pipeline" "example" {
}
```

The resource can be imported using the `location` and `pipeline_id`, for example:

```shell
terraform import ionoscloud_logging_pipeline.example {location}:{pipeline_id}
```