---
subcategory: "Autoscaling"
layout: "ionoscloud"
page_title: "IonosCloud : ionoscloud_autoscaling_group_servers"
sidebar_current: "docs-ionoscloud_autoscaling_group_servers"
description: |-
  Get information on servers generated as part of the autoscaling group.
---

# ionoscloud\_autoscaling_group_servers

The autoscaling group servers data source can be used to search for and return existing servers that are part of a specific autoscaling group.

## Example Usage

```hcl
data "ionoscloud_autoscaling_group_servers" "autoscaling_group_servers" {
	group_id = "autoscaling_group_uuid"
}
```

## Argument Reference

* `group_id` - (Required) The unique ID of the autoscaling group.

`group_id` must be provided. If it is not provided, the datasource will return an error.

## Attributes Reference

The following attributes are returned by the datasource:

* `group_id` - Id of the autoscaling group.
* `servers` - List of servers.
    * `id` - The unique ID of the server.