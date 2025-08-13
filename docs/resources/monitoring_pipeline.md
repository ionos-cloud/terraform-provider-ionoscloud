---
subcategory: "Monitoring Service"
layout: "ionoscloud"
page_title: "IonosCloud: ionoscloud_monitoring_pipeline"
sidebar_current: "docs-resource-monitoring_pipeline"
description: |-
  Creates and manages Monitoring pipeline objects.
---

# ionoscloud_monitoring_pipeline

Manages a [Monitoring pipeline](https://docs.ionos.com/cloud/observability/monitoring-service).

> ⚠️  Only tokens are accepted for authorization in the **monitoring_pipeline** resource. Please ensure you are using tokens as other methods will not be valid.

## Usage example

```hcl
resource "ionoscloud_monitoring_pipeline" "example" {
  location = "es/vit"
  name = "pipelineExample"
}
```

**NOTE:** The default timeout for all operations is 60 minutes. If you want to change the default value, you can use `timeouts` attribute inside the resource:

```hcl
resource "ionoscloud_monitoring_pipeline" "example" {
  location = "es/vit"
  name = "pipelineExample"
  timeouts {
    create = "20m"
    read = "30s"
    update = "10m"
    delete = "10m"
  }
}
```

## Argument reference

* `name` - (Required)[string] The name of the Monitoring pipeline.
* `location` - (Optional)[string] The location of the Monitoring pipeline. Default is `de/fra`. It can be one of `de/fra`, `de/fra/2`, `de/txl`, `es/vit`, `gb/bhx`, `gb/lhr`,`fr/par`, `us/mci`. If this is not set and if no value is provided for the `IONOS_API_URL_MONITORING` env var, the default `location` will be: `de/fra`.
* `grafana_endpoint` - (Computed)[string] The endpoint of the Grafana instance.
* `http_endpoint` - (Computed)[string] The HTTP endpoint of the monitoring instance.
* `key` - (Computed)(Sensitive)[string] The key used to connect to the monitoring pipeline.

> **⚠ NOTE:** `IONOS_API_URL_MONITORING` can be used to set a custom API URL for the resource. `location` field needs to be empty, otherwise it will override the custom API URL.

## Import

In order to import a Monitoring pipeline, you can define an empty Monitoring pipeline resource in the plan:

```hcl
resource "ionoscloud_monitoring_pipeline" "example" {
}
```

The resource can be imported using the `location` and `pipeline_id`, for example:

```shell
terraform import ionoscloud_monitoring_pipeline.example location:pipeline_id
```