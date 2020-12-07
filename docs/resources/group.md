---
layout: "ionoscloud"
page_title: "IonosCloud: group"
sidebar_current: "docs-resource-group"
description: |-
  Creates and manages group objects.
---

# ionoscloud\_group

Manages groups and group privileges on IonosCloud.

## Example Usage

```hcl
resource "ionoscloud_group" "group" {
  name = "my group"
  create_datacenter = true
  create_snapshot = true
  reserve_ip = true
  access_activity_log = false
  user_id="user_id"
}
```

##Argument reference

* `access_activity_log` - (Required) [Boolean] The group will be allowed to access the activity log.
* `create_datacenter` - (Optional) [Boolean] The group will be allowed to create virtual data centers.
* `create_snapshot` - (Optional) [Boolean] The group will be allowed to create snapshots.
* `name` - (Optional) [string] A name for the group.
* `reserve_ip` - (Optional) [Boolean] The group will be allowed to reserve IP addresses.
* `user_id` - (Optional) [string] The ID of the specific user to add to the group.
