---
layout: "ionoscloud"
page_title: "IonosCloud: share"
sidebar_current: "docs-resource-share"
description: |-
  Creates and manages share objects.
---

# ionoscloud\_share

Manages shares and list shares permissions granted to the group members for each shared resource.

## Example Usage

```hcl
resource "ionoscloud_share" "share" {
  group_id = "groupId"
  resource_id = "resourceId"
  edit_privilege = true
  share_privilege = false
}
```

## Argument reference

* `edit_privilege` - (Optional)[Boolean] The group has permission to edit privileges on this resource.
* `group_id` - (Required)[string] The ID of the specific group containing the resource to update.
* `resource_id` - (Required)[string] The ID of the specific resource to update.
* `share_privilege` - (Optional)[Boolean] The group has permission to share this resource.

## Import

Resource Share can be imported using the `resource id`, e.g.

```shell
terraform import ionoscloud_share.myshare {group uuid}/{resource uuid}
```