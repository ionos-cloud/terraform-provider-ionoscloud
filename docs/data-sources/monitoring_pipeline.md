---
subcategory: "Monitoring Service"
layout: "ionoscloud"
page_title: "IonosCloud: ionoscloud_monitoring_pipeline"
sidebar_current: "docs-resource-monitoring_pipeline"
description: |-
  Get information on a Monitoring pipeline.
---

# ionoscloud_monitoring_pipeline

The **Monitoring pipeline** datasource can be used to search for and return an existing Monitoring pipeline.
If a single match is found, it will be returned. If your search results in multiple matches, an error will be returned.

> ⚠️  Only tokens are accepted for authorization in the **monitoring_pipeline** data source. Please ensure you are using tokens as other methods will not be valid.

## Example Usage

### By ID
```hcl
data "ionoscloud_monitoring_pipeline" "example" {
  location = "de/txl"
  id = "pipeline_id"
}
```

### By name
```hcl
data "ionoscloud_monitoring_pipeline" "example" {
  location = "de/txl"
  name = "pipeline_name"
}
```

## Argument reference
* `location` - (Optional)[string] The location of the Monitoring pipeline. Default is `de/fra`. It can be one of `de/fra`, `de/txl`, `gb/lhr`, `es/vit`, `fr/par`. If this is not set and if no value is provided for the `IONOS_API_URL` env var, the default `location` will be: `de/fra`.
* `id` - (Optional)[string] The ID of the Monitoring pipeline you want to search for.
* `name` - (Optional)[string] The name of the Monitoring pipeline you want to search for.

Either `id` or `name` must be provided. If none, or both are provided, the datasource will return an error.

## Attributes Reference

The following attributes are returned by the datasource:

* `id` - The UUID of the Monitoring pipeline.
* `name` - The name of the Monitoring pipeline.
* `grafana_address` - The endpoint of the Grafana instance.
* `http_endpoint`- The HTTP endpoint of the Monitoring instance.