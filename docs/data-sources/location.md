---
layout: "ionoscloud"
page_title: "IonosCloud: location"
sidebar_current: "docs-datasource-location"
description: |-
  Get information on a IonosCloud Locations
---

# ionoscloud\_location

The locations data source can be used to search for and return an existing location which can then be used elsewhere in the configuration.

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