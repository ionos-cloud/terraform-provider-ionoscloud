---
layout: "ionoscloud"
page_title: "IonosCloud: backup_unit"
sidebar_current: "docs-resource-backup-unit"
description: |-
  Creates and manages IonosCloud Backup Units.
---

# ionoscloud_backup_unit

Manages a Backup Unit on IonosCloud.

## Example Usage

```hcl
resource "ionoscloud_backup_unit" "example" {
  name        = "example"
  password    = "<example-password>"
  email       = "example@example-domain.com"
}
```

## Argument Reference

The following arguments are supported:

- `name` - (Required)[string] The name of the Backup Unit. This argument is immutable.
- `password` - (Required)[string] The desired password for the Backup Unit
- `email` - (Required)[string] The email address assigned to the backup unit
- `login`- (Computed) The login associated with the backup unit. Derived from the contract number

## Import

A Backup Unit resource can be imported using its `resource id`, e.g.

```shell
terraform import ionoscloud_backup_unit.demo {backup_unit_uuid}
```

This can be helpful when you want to import backup units which you have already created manually or using other means, outside of terraform. Please note that you need to manually specify the password when first declaring the resource in terraform, as there is no way to retrieve the password from the Cloud API.

## Important Notes

- Please note that at the moment, Backup Units cannot be renamed
- Please note that the password attribute is write-only, and it cannot be retrieved from the API when importing a ionoscloud_backup_unit. The only way to keep track of it in Terraform is to specify it on the resource to be imported, thus, making it a required attribute.