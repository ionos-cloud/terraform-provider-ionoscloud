---
subcategory: "Managed Backup"
layout: "ionoscloud"
page_title: "IonosCloud: backup_unit"
sidebar_current: "docs-datasource-backup-unit"
description: |-
  Get Information on a IonosCloud Backup Unit
---

# ionoscloud\_backup_unit

The **Backup Unit data source** can be used to search for and return an existing Backup Unit.
You can provide a string for either id or name parameters which will be compared with provisioned Backup Units. 
If a single match is found, it will be returned. If your search results in multiple matches, an error will be returned.
When this happens, please refine your search and make sure that your resources have unique names.

## Example Usage

### By ID
```hcl
data "ionoscloud_backup_unit" "example" {
  id			= <backup_unit_id>
}
```

### By Name
```hcl
data "ionoscloud_backup_unit" "example" {
  name			= "Backup Unit Example"
}
```

### By Name with Partial Match
```hcl
data "ionoscloud_backup_unit" "example" {
  name			= "Example"
  partial_match	= true
}
```

## Argument Reference

The following arguments are supported:

* `id` - (Optional) ID of the backup unit you want to search for.
* `name` - (Optional) Name of an existing backup unit that you want to search for. Search by name is case-insensitive. The whole resource name is required if `partial_match` parameter is not set to true.
* `partial_match` - (Optional) Whether partial matching is allowed or not when using name argument. Default value is false.

Either `id` or `name` must be provided. If none, or both are provided, the datasource will return an error.

## Attributes Reference

The following attributes are returned by the datasource:

* `id` - The id of the Backup Unit.
* `name` - The name of the Backup Unit.
* `email` - The e-mail address you want assigned to the backup unit.
* `login` - The login associated with the backup unit. Derived from the contract number.