---
layout: "ionoscloud"
page_title: "IonosCloud: user"
sidebar_current: "docs-resource-user"
description: |-
  Creates and manages user objects.
---

# ionoscloud\_user

Manages users and list users and groups associated with that user.

## Example Usage

```hcl
resource "ionoscloud_user" "user" {
  first_name = "terraform"
  last_name = "test"
  email = "%s"
  password = "abc123-321CBA"
  administrator = false
  force_sec_auth= false
}
```

## Argument reference

* `administrator` - (Required)[Boolean] The group has permission to edit privileges on this resource.
* `email` - (Required)[string] An e-mail address for the user.
* `first_name` - (Required)[string] A first name for the user.
* `force_sec_auth` - (Required)[Boolean] Indicates if secure (two-factor) authentication should be enabled for the user (true) or not (false).
* `last_name` - (Required)[string] A last name for the user.
* `password` - (Required)[string] A password for the user.
