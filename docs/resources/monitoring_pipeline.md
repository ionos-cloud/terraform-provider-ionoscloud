---
subcategory: "Monitoring Service"
layout: "ionoscloud"
page_title: "IonosCloud: ionoscloud_monitoring_pipeline"
sidebar_current: "docs-resource-monitoring_pipeline"
description: |-
  Creates and manages Monitoring pipeline objects.
---

# ionoscloud_monitoring_pipeline

Manages a **Monitoring pipeline**.

> ⚠️  Only tokens are accepted for authorization in the **monitoring_pipeline** resource. Please ensure you are using tokens as other methods will not be valid.

## Usage example

```hcl
resource "ionoscloud_monitoring_pipeline" "example" {
  location = "es/vit"
  name = "pipelineExample"
}
```

**NOTE:** The default timeout for all operations is 20 minutes. If you want to change the default value, you can use `timeouts` attribute inside the resource:

```hcl
resource "ionoscloud_monitoring_pipeline" "example" {
  location = "es/vit"
  name = "pipelineExample"
  timeouts {
    create = "10m"
    read = "30s"
    update = "5m"
    delete = "1m"
  }
}
```

## Argument reference

* `name` - (Required)[string] The name of the Monitoring pipeline.
* `location` - (Optional)[string] The location of the Monitoring pipeline. Default is `de/fra`. It can be one of `de/fra`, `de/txl`, `gb/lhr`, `es/vit`, `fr/par`. If this is not set and if no value is provided for the `IONOS_API_URL` env var, the default `location` will be: `de/fra`.
* `grafana_endpoint` - (Computed)[string] The endpoint of the Grafana instance.
* `http_endpoint` - (Computed)[string] The HTTP endpoint of the monitoring instance.

## Import

In order to import a Monitoring pipeline, you can define an empty Monitoring pipeline resource in the plan:

```hcl
resource "ionoscloud_monitoring_pipeline" "example" {
}
```

The resource can be imported using the `location` and `pipeline_id`, for example:

```shell
terraform import ionoscloud_monitoring_pipeline.example {location}:{pipeline_id}
```