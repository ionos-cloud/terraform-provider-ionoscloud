---
subcategory: "Compute Engine"
layout: "ionoscloud"
page_title: "IonosCloud : servers"
sidebar_current: "docs-ionoscloud-datasource-servers"
description: |-
  Retrieves a list of Ionos Cloud Servers
---

# ionoscloud\_servers

The **Servers data source** can be used to search for and return existing servers based on filters used.

## Example Usage

### By Name
```hcl
data ionoscloud_servers example {
  datacenter_id = ionoscloud_datacenter.example.id
  filter {
    name = "name"
    value = "server_name_to_look_here"
  }
}
```

### By CPU Family
```hcl
data ionoscloud_servers example {
  datacenter_id = ionoscloud_datacenter.example.id
  filter {
    name = "cpu_family"
    value = "INTEL_XEON"
  }
}
```


### By Name and Cores
```hcl
data ionoscloud_servers example {
  datacenter_id = ionoscloud_datacenter.example.id
  filter {
    name = "name"
    value = "test"
  }
  filter {
    name = "cores"
    value = "1"
  }
}
```

## Argument Reference

* `datacenter_id` - (Required) Name of an existing datacenter that the servers are a part of
* `filter` -  (Required) One or more name/value pairs to filter off of. You can use most base fields in the [server](../resources/server.md) resource. These do **NOT** include nested fields in nics or volume nested fields.


`datacenter_id` and `filter` must be provided. If none either is not provided, the datasource will return an error.

**NOTE:** Lookup by filter is partial. Searching for a server using filter name and value `test`, will find all servers that have `test` in the name. 
For example it will find servers named `test`, `test1`, `testsomething`. 

**NOTE:** You cannot search by image_name by providing an alias like `ubuntu`.

## Attributes Reference

The following attributes are returned by the datasource:

* `servers` - list of servers that matches the filters provided.
For a full reference of all attributes returned, check out [documentation](../resources/server.md)