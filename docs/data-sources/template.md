---
layout: "ionoscloud"
page_title: "IonosCloud : ionoscloud_template"
sidebar_current: "docs-ionoscloud-datasource-template"
description: |-
Get information on a Ionos Cloud Template
---

# ionoscloud_template

The template data source can be used to search for and return existing templates.

## Example Usage

```hcl
data "ionoscloud_template" "example" {
	name = "BETA CUBES S"
	cores = 1
	ram	  = 2048
	storage_size = 50
}
```

## Argument Reference

* `name` - (Required) A name of that resource.
* `cores` - (Required) The CPU cores count.
* `ram` - (Required) The RAM size in MB.
* `storage_size` - (Required) The storage size in GB.

All arguments must be provided. If none, the datasource will return an error.

## Attributes Reference

The following attributes are returned by the datasource:

* `id`
* `name`
* `cores`
* `ram`
* `storage_size`
