---
subcategory: "User Management"
layout: "ionoscloud"
page_title: "IonosCloud: share"
sidebar_current: "docs-resource-share"
description: |-
  Creates and manages share objects.
---

# ionoscloud\_share

Manages **Shares** and list shares permissions granted to the group members for each shared resource.

## Example Usage

```hcl
resource "ionoscloud_datacenter" "example" {
	name                = "Datacenter Example"
	location            = "us/las"
	description         = "Datacenter Description"
	sec_auth_protection = false
}

resource "ionoscloud_group" "example" {
  name                    = "Group Example"
  create_datacenter       = true
  create_snapshot         = true
  reserve_ip              = true
  access_activity_log     = true
  create_pcc              = true
  s3_privilege            = true
  create_backup_unit      = true
  create_internet_access  = true
  create_k8s_cluster      = true
}

resource "ionoscloud_share" "example" {
  group_id                = ionoscloud_group.example.id
  resource_id             = ionoscloud_datacenter.example.id
  edit_privilege          = true
  share_privilege         = false
}
```

## Argument reference

* `edit_privilege` - (Optional)[Boolean] The group has permission to edit privileges on this resource.
* `group_id` - (Required)[string] The ID of the specific group containing the resource to update.
* `resource_id` - (Required)[string] The ID of the specific resource to update.
* `share_privilege` - (Optional)[Boolean] The group has permission to share this resource.

⚠️ **Note:** There is a limitation due to which the creation of several shares at the same time leads
to an error. To avoid this, `parallelism=1` can be used when running `terraform apply` command in order
to create the resources in a sequential manner. Another solution involves the usage of `depends_on`
attributes inside the `ionoscloud_share` resource to enforce the sequential creation of the shares.

## Import

Resource Share can be imported using the `resource id`, e.g.

```shell
terraform import ionoscloud_share.myshare {group uuid}/{resource uuid}
```