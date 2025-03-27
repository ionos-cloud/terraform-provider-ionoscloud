---
subcategory: "User Management"
layout: "ionoscloud"
page_title: "IonosCloud: share"
sidebar_current: "docs-datasource-share"
description: |-
  Get Information on share permission objects
---

# ionoscloud_share

The **Share data source** can be used to search for and return an existing share object.
You need to provide the group_id and resource_id to get the group resources for the shared resource.
If a single match is found, it will be returned. If your search results in multiple matches, an error will be returned.
When this happens, please refine your search string so that it is specific enough to return only one result.

## Example Usage

```hcl
data "ionoscloud_share" "example" {
  group_id      = "group_id"
  resource_id   = "resource_id"
}
```

## Argument Reference

The following arguments are supported:

* `group_id` - (Required)The ID of the specific group containing the resource to update.
* `resource_id` - (Required)The ID of the specific resource to update.


`resource_id` and `group_id` must be provided. If any of them are missing, the datasource will return an error.

## Attributes Reference

The following attributes are returned by the datasource:

* `id` - The id of the share resource.
* `group_id` - The ID of the specific group containing the resource to update.
* `resource_id` - The ID of the specific resource to update.
* `edit_privilege` - The flag that specifies if the group has permission to edit privileges on this resource.
* `share_privilege` - The group has permission to share this resource.