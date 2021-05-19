---
layout: "ionoscloud"
page_title: "IonosCloud : datacenter"
sidebar_current: "docs-datasource-datacenter"
description: |-
  Get information on a IonosCloud Data Centers
---

# ionoscloud\_datacenter

The data centers data source can be used to search for and return an existing Virtual Data Center. You can provide a string for the name and location parameters which will be compared with provisioned Virtual Data Centers. If a single match is found, it will be returned. If your search results in multiple matches, an error will be generated. When this happens, please refine your search string so that it is specific enough to return only one result.

## Example Usage

```hcl
data "ionoscloud_datacenter" "dc_example" {
  name     = "test_dc"
  location = "us"
}
```

## Argument Reference

 * `name` - (Required) Name of an existing Virtual Data Center that you want to search for.
 * `location` - (Optional) Id of the existing Virtual Data Center's location.

## Attributes Reference

 * `id` - UUID of the Virtual Data Center
