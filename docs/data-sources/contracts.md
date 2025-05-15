---
subcategory: "Compute Engine"
layout: "ionoscloud"
page_title: "IonosCloud : ionoscloud_contracts"
sidebar_current: "docs-ionoscloud-datasource-contracts"
description: |-
  Get information on the contracts available in your IONOS Cloud account
---

# Contracts Data Source

The `contracts` data source provides information about the contracts available in your IONOS Cloud account, including resource limits and other contract details.

## Example Usage

```hcl
data "ionoscloud_contracts" "example" {}

output "contracts" {
  value = data.ionoscloud_contracts.example.contracts
}
```
The following attributes are returned by the datasource:

Sure! Here's the list of attributes formatted as requested:

* `contracts` 
  * `contract_number` - The contract number.
  * `owner` - The contract owner's user name.
  * `status` - The contract status.
  * `reg_domain` - The registration domain of the contract.
  * `resource_limits`
    * `cores_per_server` - The maximum number of cores per server.
    * `ram_per_server` - The maximum RAM per server in MB.
    * `ram_per_contract` - The maximum RAM per contract in MB.
    * `cores_per_contract` - The maximum number of cores per contract.
    * `cores_provisioned` - The number of cores provisioned.
    * `das_volume_provisioned` - The DAS volume provisioned.
    * `hdd_limit_per_contract` - The HDD limit per contract.
    * `hdd_limit_per_volume` - The HDD limit per volume.
    * `hdd_volume_provisioned` - The HDD volume provisioned.
    * `k8s_cluster_limit_total` - The total Kubernetes cluster limit.
    * `k8s_clusters_provisioned` - The number of Kubernetes clusters provisioned.
    * `nat_gateway_limit_total` - The total NAT gateway limit.
    * `nat_gateway_provisioned` - The number of NAT gateways provisioned.
    * `nlb_limit_total` - The total NLB limit.
    * `nlb_provisioned` - The number of NLBs provisioned.
    * `ram_provisioned` - The RAM provisioned.
    * `reservable_ips` - The number of reservable IPs.
    * `reserved_ips_in_use` - The number of reserved IPs in use.
    * `reserved_ips_on_contract` - The number of reserved IPs on the contract.
    * `ssd_limit_per_contract` - The SSD limit per contract.
    * `ssd_limit_per_volume` - The SSD limit per volume.
    * `ssd_volume_provisioned` - The SSD volume provisioned.
    * `security_groups_per_vdc` - The number of security groups per VDC.
    * `security_groups_per_resource` - The number of security groups per resource.
    * `rules_per_security_group` - The number of rules per security group.  