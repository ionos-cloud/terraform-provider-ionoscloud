---
subcategory: "Compute Engine"
layout: "ionoscloud"
page_title: "IonosCloud : ionoscloud_template"
sidebar_current: "docs-ionoscloud-datasource-template"
description: |-
  Get information on a Ionos Cloud Template
---

# ionoscloud_template

The **Template data source** can be used to search for and return existing templates by providing any of template properties (name, cores, ram, storage_size).
If a single match is found, it will be returned. If your search results in multiple matches, an error will be returned.
When this happens, please refine your search string so that it is specific enough to return only one result.

## Example Usage

### By ID
```hcl
data "ionoscloud_template" "example" {
	id = <template_id>
}
```

### By Name
```hcl
data "ionoscloud_template" "example" {
	name = "CUBES S"
}
```

### By Cores
```hcl
data "ionoscloud_template" "example" {
	cores = 6
}
```

### By Ram
```hcl
data "ionoscloud_template" "example" {
	ram = 49152
}
```

### By Storage Size
```hcl
data "ionoscloud_template" "example" {
	storage_size = 80
}
```

## Argument Reference

* `id` - (Optional) ID of the template you want to search for.
* `name` - (Optional) The name of template. Search by name is case-insensitive. The whole resource name is required if `partial_match` parameter is not set to true..
* `partial_match` - (Optional) Whether partial matching is allowed or not when using name argument. Default value is false.
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
