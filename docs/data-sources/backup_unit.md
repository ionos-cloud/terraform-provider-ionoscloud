---
subcategory: "Managed Backup"
layout: "ionoscloud"
page_title: "IonosCloud: backup_unit"
sidebar_current: "docs-datasource-backup-unit"
description: |-
  Get Information on a IonosCloud Backup Unit
---

# ionoscloud_backup_unit

The **Backup Unit data source** can be used to search for and return an existing Backup Unit.
You can provide a string for either id or name parameters which will be compared with provisioned Backup Units. 
If a single match is found, it will be returned. If your search results in multiple matches, an error will be returned. 
When this happens, please refine your search string so that it is specific enough to return only one result.

## Example Usage

### By ID
```hcl
data "ionoscloud_backup_unit" "example" {
  id			= "backup_unit_id"
}
```

### By Name
```hcl
data "ionoscloud_backup_unit" "example" {
  name			= "Backup Unit Example"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Optional) Name of an existing backup unit that you want to search for.
* `id` - (Optional) ID of the backup unit you want to search for.

Either `name` or `id` must be provided. If none, or both are provided, the datasource will return an error.

## Attributes Reference

The following attributes are returned by the datasource:

* `id` - The id of the Backup Unit.
* `name` - The name of the Backup Unit.
* `email` - The e-mail address you want assigned to the backup unit.
* `login` - The login associated with the backup unit. Derived from the contract number.