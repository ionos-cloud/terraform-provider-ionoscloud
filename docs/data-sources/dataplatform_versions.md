---
subcategory: "Dataplatform"
layout: "ionoscloud"
page_title: "IonosCloud : ionoscloud_dataplatform_versions"
sidebar_current: "docs-dataplatform_versions"
description: |-
  Get information on Managed Dataplatform API versions.
---

# ionoscloud\_dataplatform_versions

⚠️ **Note:** Data Platform is currently in the Early Access (EA) phase.
We recommend keeping usage and testing to non-production critical applications.
Please contact your sales representative or support for more information.
Manages a **Dataplatform Node Pool**.

The **Dataplatform Versions Data Source** can be used to search for and retrieve list of available Managed Dataplatform API versions.


## Example Usage


### Retrieve list of Managed Dataplatform API versions
```hcl
data "ionoscloud_dataplatform_versions" "example" {
}
```

## Attributes Reference

The following attributes are returned by the datasource:

* `versions` - list of Managed Dataplatform API versions.