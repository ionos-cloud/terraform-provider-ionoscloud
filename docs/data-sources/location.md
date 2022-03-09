---
subcategory: "Compute Engine"
layout: "ionoscloud"
page_title: "IonosCloud: location"
sidebar_current: "docs-datasource-location"
description: |-
  Get information on a IonosCloud Locations
---

# ionoscloud\_location

The location's data source can be used to search for and return an existing location which can then be used elsewhere in the configuration.
If a single match is found, it will be returned. If your search results in multiple matches, an error will be returned.
When this happens, please refine your search string so that it is specific enough to return only one result.

## Example Usage

```hcl
data "ionoscloud_location" "loc1" {
  name    = "karlsruhe"
  feature = "SSD"
}
```

## Argument Reference

 * `name` - (Required) Name of the location to search for.
 * `feature` - (Optional) A desired feature that the location must be able to provide.

## Attributes Reference

The following attributes are returned by the datasource:

 * `id` - UUID of the location
 * `cpu_architecture` - Array of features and CPU families available in a location
  * `cpu_family` - A valid CPU family name.
  * `max_cores` - The maximum number of cores available.
  * `max_ram` - The maximum number of RAM in MB.
  * `vendor` - A valid CPU vendor name.
* `image_aliases` - List of image aliases available for the location