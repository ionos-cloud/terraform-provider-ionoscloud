---
subcategory: "Compute Engine"
layout: "ionoscloud"
page_title: "IonosCloud : ionoscloud_template"
sidebar_current: "docs-ionoscloud-datasource-template"
description: |-
  Get information on a Ionos Cloud Template
---

# ionoscloud_template

The template data source can be used to search for and return existing templates by providing any of template properties (name, cores, ram, storage_size).
If a single match is found, it will be returned. If your search results in multiple matches, an error will be returned.
When this happens, please refine your search string so that it is specific enough to return only one result.

## Example Usage

```hcl
data "ionoscloud_template" "example" {
	name = "BETA CUBES S"
}
```

## Argument Reference

* `name` - (Optional) A name of that resource.
* `cores` - (Optional) The CPU cores count.
* `ram` - (Optional) The RAM size in MB.
* `storage_size` - (Optional) The storage size in GB.

Any of the arguments ca be provided. If none, the datasource will return an error.

## Attributes Reference

The following attributes are returned by the datasource:

* `id` - Id of template
* `name` - Name of template
* `cores`- The CPU cores count
* `ram` - The RAM size in MB
* `storage_size` - The storage size in GB
