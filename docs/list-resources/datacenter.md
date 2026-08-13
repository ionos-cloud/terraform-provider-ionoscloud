---
subcategory: "Compute Engine"
layout: "ionoscloud"
page_title: "IONOS CLOUD: datacenter"
description: |-
  Lists IONOS CLOUD Virtual Data Centers.
---

# List Resource: ionoscloud_datacenter

⚠️ **Note:** List Resources require HashiCorp Terraform version 1.14 or later and are queried using `terraform query`.

Lists IONOS CLOUD Virtual Data Centers.

## Example Usage

⚠️ **Note:** `list` blocks must be placed in a dedicated `tfquery.hcl` file, separate from your main Terraform configuration.

### List all datacenters

```hcl
list "ionoscloud_datacenter" "all" {
  provider         = ionoscloud
  include_resource = true
}
```

### Filter datacenters by name

```hcl
list "ionoscloud_datacenter" "prod" {
  provider         = ionoscloud
  include_resource = true
  config {
    filters = [
      { field_name = "name", field_value = "prod" }
    ]
  }
}
```

### Filter datacenters by name and location

```hcl
list "ionoscloud_datacenter" "prod_txl" {
  provider         = ionoscloud
  include_resource = true
  config {
    filters = [
      { field_name = "name",     field_value = "prod" },
      { field_name = "location", field_value = "de/txl" },
    ]
  }
}
```

### Generate resource configuration from existing datacenters

Use `terraform query` with `-generate-config-out` to produce ready-to-use `ionoscloud_datacenter` resource blocks for all existing datacenters:

```shell
terraform query -generate-config-out=imported.tf
```

Terraform will write an `ionoscloud_datacenter` resource block for each discovered datacenter into `imported.tf`, which can then be used directly in your configuration.

## Argument Reference

The `config` block supports the following arguments:

- `filters` - (Optional) List of filters to apply. All filters must match (AND logic). Each filter supports:
  - `field_name` - (Required) The field to filter on. Supported values: `name`, `location`.
  - `field_value` - (Required) The value to match against. The Cloud API performs a substring ("contains") match.

## Identity Attributes

Each result exposes the following identity attribute, usable for import:

| Attribute | Description                 |
|-----------|------------------------------|
| `id`      | The UUID of the datacenter. |

## Attributes Reference

Each result exposes the following attributes when `include_resource = true`, matching the `ionoscloud_datacenter` resource schema:

- `id` - The UUID of the datacenter.
- `name` - The name of the datacenter.
- `location` - The physical location where the datacenter is provisioned.
- `description` - A description for the datacenter.
- `version` - The version of the datacenter.
- `features` - A list of features supported by the datacenter.
- `sec_auth_protection` - Whether secure authentication protection is enabled.
- `cpu_architecture` - CPU architecture block:
  - `cpu_family` - The CPU family.
  - `max_cores` - Maximum number of cores.
  - `max_ram` - Maximum amount of RAM.
  - `vendor` - The CPU vendor.
- `ipv6_cidr_block` - Auto-assigned /56 IPv6 CIDR block, if IPv6 is enabled for the datacenter.
