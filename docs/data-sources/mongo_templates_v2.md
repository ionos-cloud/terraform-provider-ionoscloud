---
subcategory: "Database as a Service - MongoDB v2"
layout: "ionoscloud"
page_title: "IONOS CLOUD: ionoscloud_mongo_templates_v2"
sidebar_current: "docs-datasource-mongo_templates_v2"
description: |-
  Lists IONOS CLOUD MongoDB V2 Instance Templates.
---

# ionoscloud_mongo_templates_v2

The `ionoscloud_mongo_templates_v2` data source can be used to retrieve the available MongoDB V2 instance templates in a given location, with an optional filter by edition. A template ID is used for the `instance_template_id` attribute of a `BUSINESS` edition cluster.

## Example Usage

```hcl
data "ionoscloud_mongo_templates_v2" "all" {
  location = "gb/bhx"
}

data "ionoscloud_mongo_templates_v2" "business" {
  location = "gb/bhx"
  edition  = "BUSINESS"
}
```

## Argument Reference

* `location` - (Required)[string] The location to query. Available locations: `gb/bhx`, `us/ewr`.
* `edition` - (Optional)[string] Filter templates by edition (`ENTERPRISE` or `BUSINESS`).

## Attributes Reference

The following attributes are returned by the datasource:

* `items` - List of instance templates. Each item exposes:
  * `id` - The UUID of the instance template.
  * `name` - The name of the instance template.
  * `edition` - The edition the template applies to.
  * `cores` - The number of CPU cores.
  * `ram` - The amount of memory in gigabytes (GB).
  * `storage_size` - The amount of storage in gigabytes (GB).
