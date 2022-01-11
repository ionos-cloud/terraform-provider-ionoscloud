---
layout: "ionoscloud"
page_title: "IonosCloud: share"
sidebar_current: "docs-datasource-share"
description: |-
Get Information on share permission objects
---

# ionoscloud\_share

The share data source can be used to search for and return an existing share object.
You need to provide the group_id and resource_id to get the group resources for the shared resource.

## Example Usage

```hcl
data "ionoscloud_share" "share" {
  group_id = "groupId"
  resource_id = "resourceId"
  id			= "share uuid"
}
```

## Argument Reference

The following arguments are supported:

* `group_id` - (Required)The ID of the specific group containing the resource to update.
* `resource_id` - (Required)The ID of the specific resource to update.
* `id` - (Required)The uuid of the share object


`id`, `resource_id` and `group_id` must be provided. If any of them are missing, the datasource will return an error.

## Attributes Reference

The following attributes are returned by the datasource:

* `id` - The id of the share resource.
* `group_id` - The ID of the specific group containing the resource to update.
* `resource_id` - The ID of the specific resource to update.
* `edit_privilege` - The flag that specifies if the group has permission to edit privileges on this resource.
* `share_privilege` - The group has permission to share this resource.