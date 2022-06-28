---
subcategory: "Data Stack as a Service"
layout: "ionoscloud"
page_title: "IonosCloud : ionoscloud_dsaas_versions"
sidebar_current: "docs-dsaas_versions"
description: |-
Get information on Managed Data Stack API versions.
---

# ionoscloud\_dsaas_versions

The **DSaaS Versions Data Source** can be used to search for and retrieve list of available Managed Data Stack API versions.


## Example Usage


### Retrieve list of Managed Data Stack API versions
```hcl
data "ionoscloud_dsaas_versions" "example" {
}
```

## Attributes Reference

The following attributes are returned by the datasource:

* `versions` - list of Managed Data Stack API versions.