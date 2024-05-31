## 6.4.17
### Fixes
- Correctly raise immutable error for changes to `template_uuid` when running `terraform plan` for Cube servers

### Documentation
- Update `ionoscloud_user` documentation. Fix `administrator` and add other fields description
- Change to have nested lists show correctly in tf registry docs

### Enhancements
- Add configurable fields to `ionoscloud_share` resource. Fields that can be set on creation: `description`, `sec_auth_protection`, `licence_type`. 
Updatable fields: `description`, `licence_type`, `nic_hot_plug`, `cpu_hot_plug`, `nic_hot_unplug`, `disc_virtio_hot_plug`, `disc_virtio_hot_unplug`, `ram_hot_plug`.
- Allow MariaDB cluster creation in other zones than `de/txl` by adding `location` parameter to resources and data sources


## 6.4.16
### Enhancements
- Modify DBaaS workflow to run tests in multiple stages for every service (Mongo, MariaDB, PgSQL) rather than running all tests in one stage
### Fixes
- Fix MongoDB user import
- Fix k8s cluster tests
- Fix #552 in order to allow Dataplatform cluster creation without `lans` or `routes`

## 6.4.15
### Fixes
- Increase max result limit of data sources for target groups (200) and IP blocks (1000), as a workaround for pagination issues.
- Add email filter for user data source, fixes pagination issues for users.
- Fix name validation for Dataplatform resources.

### Enhancements
- Change location for MongoDB tests to improve running time.
- Change location for PgSQL tests to improve running time.

### Features
- Add new attribute for Dataplatform clusters: `lans`.

## 6.4.14
### Features
- Add MariaDB cluster resource, data source and backups data source.
### Fixes
- #524 `filters` is now optional for `ionoscloud_servers` data source. If not provided, all servers in the configured datacenter will be returned.
- `filters` is now optional for `ionoscloud_clusters` data source. If not provided, all k8s clusters will be returned.
- Populate `server` field for k8s cluster data sources

### Documentation
- Update documentation for pgsql cluster and mongo cluster

## 6.4.13
### Features
- Added ability to boot from network for `ionoscloud_server`, `ionoscloud_vcpu_server`, `ionoscloud_cube_server`
- Add `ionoscloud_k8s_clusters` data source
- Add `vulnerability_scanning` parameter to `ionoscloud_container_registry` resource.
### Refactor
- Remove duplicate functions for image retrieval (`checkImage`, `resolveImageName`) in `resource_volume.go`
### Fixes
- Remove `credentials` field from `ionoscloud_mongo_cluster` resource
- Remove `credentials` field from `ionoscloud_pg_cluster` data source
### Documentation
- Update documentation for K8s cluster, nodepools and shares

## 6.4.12
### Features
- Add `ionoscloud_server_boot_device_selection` resource for selecting the boot device of `ionoscloud_server`, `ionoscloud_vcpu_server` and `ionoscloud_cube_server` resources
- Increase the timeout period for `ionoscloud_node_pool` resource to 3 hours

### Fixes
- `is_system_user` is actually read-only. You cannot set it.
- Add `priority` in state for DNS Records

### Features
- Add support for creating private k8s clusters to `ionoscloud_k8s_cluster` : `public`, `location`, `nat_gateway_ip`, `node_subnet`

## 6.4.11
### Documentation
- Refactor readme files to better explain the usage of the provider

### Features
- Add `flowlog` to `ionoscloud_nic` resource
- Add `flowlog` to `ionoscloud_networkloadbalancer` resource
- Add `flowlog` to `ionoscloud_application_loadbalancer` resource
- Update dependency for terraform-plugin-sdk v2.30.0
- Use v6.4.10 of cloudapiv6 sdk
- #494 add `proxy-protocol` to `ionoscloud_networkloadbalancer_forwarding_rule` resource

### New Product - **Autoscaling**:
  - `Resources`:
    - [ionoscloud_autoscaling_group](docs/resources/autoscaling_group.md)
  - `Data Sources`:
    - [ionoscloud_autoscaling_group](docs/data-sources/autoscaling_group.md)
    - [ionoscloud_autoscaling_group_servers](docs/data-sources/autoscaling_group_servers.md)

### Fixes
- #487. Crash on server import without inline `nic`
- #503. Use `Location` func for state tracking request instead of getting Location header directly and minor refactor.
- #497. allow to set empty `name` for `ionoscloud_dns_record`
- Refactor validation to use `validation.AllDiag` instead of `validation.All`, remove unnecessary usage of `validation.All`

## 6.4.10
### Refactor
- Add `nic` service
- Use `error.As` for `requestFailed` err
- `%w` instead of `%s` for some printed errors
- Use `serve` for debug mode
### Features
- #460 add `contract_number` to provider configuration
- #412 add support to set power state of Enterprise and Cube servers, by adding the new field `vm_state` in `ionoscloud_server`, `ionoscloud_cube_server` and `ionoscloud_vcpu_server `resources
### Fixes
- #467 removing an inline `nic` of the `server` resource from dcd should not throw 404 when running plan or apply after
- #432 Now it is possible to create and delete multiple `ionoscloud_ipfailover` resources at the same time. The UUID is generated based on the IP of the 
failover group. The resources that are created using Terraform cannot be modified/deleted outside Terraform.
- Fix `nil` deref error on list for nic datasource
- #470 fix image name searching in `ionoscloud_image` and `ionoscloud_volume`. Exact matches are returned correctly now if they exist.

## 6.4.9

### Features
- Cloud DNS is now Generally Available
- Data Platform is now in Generally Available
- #451 update go sdk, allow `IONOS_CONTRACT_NUMBER` to be used to run terraform on different contract numbers for reseller accounts
- Update dependency for terraform-plugin-sdk. Stop using deprecated functions from `resource` package


## 6.4.8
### Fixes
- `primary_ip` in `ionoscloud_server` should be set on creation
- `ssh_keys` was no longer being set if server was not vcpu.
- `ssh_keys` will no longer be computed on any type of server
- `ssh_key_path` will now be set to schema on creation
- Setting explicit `ipv6_cidr_block` on `nic` resource.
- Ipv6 fields `dhcpv6`, `ipv6_cidr_block`, `ipv6_ips` not updating correctly on `ionoscloud_server` and `ionoscloud_cube_server`
- Issue caused by `dhcpv6` field for plans which do not enable the IPv6 feature
- #449. Increase `NotFoundChecks` to 9999.
### Documentation
- Example IPv6 usage for `ionoscloud_server` and `ionoscloud_cube_server`

## 6.4.7
### Features
- Add support for mongo clusters enterprise edition

## 6.4.6
### Fixes
- Fix `ipv6_ips` should not request a re-apply of the plan if `ipv6_cidr_block` is not set on the lan
- Fix `dhcpv6` should not be set on server nic if IPv6 is not enabled on the lan
- Fix `boot_cdrom` should not crash even if not set to an UUID
### Documentation
- Fix `inoscloud_image` docs to get cdrom image
- Fix `boot_cdrom` - add description and examples
### Features
- Add support for `VCPU` servers

## 6.4.5
### Features
- Make `nic` list in `ionoscloud_server` resource optional
- Make `firewall` list in `ionoscloud_server` resource optional and allow multiple inline firewall rules in the list
- Add ipv6 functionality for `ionoscloud_datacenter`, `ionoscloud_lan` and `ionoscloud_nic` resources

### Refactor
- Separate `cloudapi` code from `ionoscloud` folder, to be able to write services easier.
- Refactor validation to use `ValidateDiagFunc` instead of `ValidateFunc`, remove unnecessary usage of `validation.All()`
## 6.4.4

### Features
- Add support for PgSQL User & Databases

### Dependency update
- Update `sdk-go-dbaas-postgres` to version 1.1.2

## 6.4.3

### Documentation
- Improve example for `ionoscloud_private_crossconnect`

### Fixes
- Remove unpopulated `credentials` field from mongodb cluster data source.
- Add `ram` and `cores` fields to cube server data source.

## 6.4.2
### Fixes
- Fix `ssh_keys` field upgrade `ionoscloud_server` from `6.3.3` to higher versions should not replace server. `ssh_keys` and `ssh_key_path` fields no longer forceNew. 
`ssh_keys` is no longer computed.
- Fix `ssh_keys` suppress diff on upgrade for `ionoscloud_server` when having `volume.0.ssh_keys`
- Add validation to `label` `key` and `value` fields for `ionoscloud_server` resource
- Fix gitbook references

### Docs
- Add new products to Gitbook docs

## 6.4.1
### Fixes
- Fix `inline_volume_ids` field upgrade for `ionoscloud_server`

### Docs
- Fix `ionoscloud_image` examples
- Improve docs for `ssh_keys` and `ssh_key_path`

## 6.4.0
### Enhancement:
- Increase go version to 1.20
### Features
- Add `inline_volume_ids` computed field.
- New Product: **DNS**:
  - `Resources`:
    - [ionoscloud_dns_zone](docs/resources/dns_zone.md)
    - [ionoscloud_dns_record](docs/resources/dns_record.md)
  - `Data Sources`:
    - [ionoscloud_dns_zone](docs/data-sources/dns_zone.md)
    - [ionoscloud_dns_record](docs/data-sources/dns_record.md)
### Dependency update
- Update `sdk-go-dbaas-mongo` to [v1.0.6](https://github.com/ionos-cloud/sdk-go-dbaas-mongo/releases/tag/v1.2.2)
- Update `sdk-go-container-registry` to [v1.0.1](https://github.com/ionos-cloud/sdk-go-container-registry/releases/tag/v1.0.1)
- Update `sdk-go` to [v6.1.7](https://github.com/ionos-cloud/sdk-go/releases/tag/v6.1.7)
- Update `sdk-go-cert-manager` to [v1.0.1](https://github.com/ionos-cloud/sdk-go-cert-manager/releases/tag/v1.0.1)
- Update `terraform-plugin-sdk` to [v2.26.1](https://github.com/hashicorp/terraform-plugin-sdk/releases/tag/v2.26.1)
### Fixes
- Log levels need to be shown and filtered correctly when set with `TF_LOG`. Also change `WARNING` log levels to `WARN`.
- Update code to work with new mongo version
- Ignore downgrades of `k8s_version` patch level.
- Allow upgrades of `k8s_version` patch level.

## 6.3.6
### Feature
- Rewrite a part of the psql service to use new functionality.
- Add `dns_name` to `ionoscloud_pg_cluster` datasource and resource
- Add option to search for images in the `ionoscloud_image` data source using `image_alias`. Search will be performed with exact match.
- New Product: **DataPlatform**:
  - `Resources`:
    - [ionoscloud_dataplatform_cluster](docs/resources/dataplatform_cluster.md)
    - [ionoscloud_dataplatform_node_pool](docs/resources/dataplatform_node_pool.md)
  - `Data Sources`:
    - [ionoscloud_dataplatform_cluster](docs/data-sources/dataplatform_cluster.md)
    - [ionoscloud_dataplatform_node_pool](docs/data-sources/dataplatform_node_pool.md)
    - [ionoscloud_dataplatform_node_pools](docs/data-sources/dataplatform_node_pools.md):
    - [ionoscloud_dataplatform_versions](docs/data-sources/dataplatform_versions.md):

⚠️ **Note:** Data Platform is currently in the Early Access (EA) phase.
We recommend keeping usage and testing to non-production critical applications.
Please contact your sales representative or support for more information.

### Dependency update
- Update `go` to version 1.19
- Update `sdk-go-dbaas-postgres` to [v1.0.6](https://github.com/ionos-cloud/sdk-go-dbaas-postgres/releases/tag/v1.0.6)
### Documentation
- Update documentation for K8s node pools
- Update documentation for `ionoscloud_image` to clarify what type of search is done.
- Improve documentation for `endpoint` field
## Fixes
- Throw error on `404` for mongo cluster creation.
- Solves  #372 crash when ips field in nic resource is a list with an empty string

## 6.3.5
### Feature:
- Removed EA note for container registry and dbaas mongo docs

## 6.3.4
### Feature:
 - Add update for mongo database resources
 - Add update for mongo cluster and user
 - Add labels for servers
 - Add data source for DBaaS Mongo Templates
 - Update mongo sdk to v1.2.0
 - Added server ssh_keys tests

### Refactor:
- Refactor services, add generic `WaitForResourceToBeReady` and `WaitForResourceToBeDeleted` methods
- Removed hard coded passwords from docs and tests and replaced with dynamically generated passwords
- Remove useless checks from services

## Fixes
 - Fix mongo user tests to check for cluster state instead of user state which was removed
 - Defining a separate firewall rule for server should not set firewall_id inside server resource, as it moves the firewall resource inside the server on re-apply
 - Fixes creating share resource edit and share privileges mix up
 - `viable_node_pool_versions`  in k8s cluster is no longer optional, is only computed
 - Allow server import with nic and firewallId : `terraform import ionoscloud_server.myserver {datacenter uuid}/{server uuid}/{primary nic id}/{firewall rule id}`
 - Mongo tests update mongo version
 - Change the way in which we set the NIC data
 - Allow server import with nic and firewall ids
 - Typo in group resource
 - Readme fix link to test suite, dbaas test use correct checking function
 - Make viable_node_pool_versions only computed
 - K8s nodepool test
 - Mix up share and edit privileges on create

## 6.3.3
### Feature
- New Product: **ContainerRegistry**:
  - `Resources`:
    - [ionoscloud_container_registry](docs/resources/container_registry.md)
    - [ionoscloud_container_registry_token](docs/resources/container_registry_token.md)
  - `Data Sources`:
    - [ionoscloud_container_registry](docs/data-sources/container_registry.md)
    - [ionoscloud_container_registry_token](docs/data-sources/container_registry_token.md)
    - [ionoscloud_container_registry_locations](docs/data-sources/container_registry_locations.md)
     
⚠️ **Note:** Container Registry is currently in the Early Access (EA) phase. We recommend keeping usage and testing to non-production critical applications.
Please contact your sales representative or support for more information.
### Fixes
- Fixes #326, now removing a s3_buckets block from an ionoscloud_k8s_cluster resource triggers a change in the terraform plan.
- Fixes user creation bug, now user creation works properly when `group_ids` is specified in the plan.

## 6.3.2
### Feature
- New Product: **MongoDB**:
  - `Resources`:
    - [ionoscloud_mongo_cluster](docs/resources/dbaas_mongo_cluster.md)
    - [ionoscloud_mongo_user](docs/resources/dbaas_mongo_user.md)
  - `Data Sources`:
    - [ionoscloud_mongo_cluster](docs/data-sources/dbaas_mongo_cluster.md)
    - [ionoscloud_mongo_user](docs/data-sources/dbaas_mongo_user.md)

⚠️ **Note:** DBaaS - MongoDB is currently in the Early Access (EA) phase. We recommend keeping usage and testing to non-production critical applications.
Please contact your sales representative or support for more information.

- New Product: **Certificate Manager**:
  - `Resources`:
    - [ionoscloud_certificate](docs/resources/certificate.md)
  - `Data Sources`:
    - [ionoscloud_certificate](docs/data-sources/certificate.md)


### Enhancement:
- Increase go version to 1.18
- Update dependencies to latest versions
- Update Ionos Cloud GO SDK v6.1.3. Release notes here [v6.1.3](https://github.com/ionos-cloud/sdk-go/releases/tag/v6.1.3)
- Update SDK GO DBaaS Postgres to v1.0.4. Release notes here [v1.0.4](https://github.com/ionos-cloud/sdk-go-dbaas-postgres/releases/tag/v1.0.4)
- `ssh_key_path` will now allow the keys to be passed directly also. In the future, will be renamed to `ssh_keys`.

### Fixes
- Reproduces rarely: sometimes the `nic` resource is not found after creation. As a fix we added a retry for 5 minutes to be able to get the NIC. The retry will keep trying if the response 
is `not found`(404)
- Fix cube server creation. Some attributes were not populated - name, boot_cdrom, availability_zone
- Crash on update of k8s version when we have a value without `.`

### Documentation
- Add links to documentation for `cube` and `enterprise` fields

## 6.3.1

### Feature
- When no argument is provided for user data source, try to get the email from the client configuration
- Update Ionos Cloud GO SDK v6.1.2. Release notes here [v6.1.2](https://github.com/ionos-cloud/sdk-go/releases/tag/v6.1.2)
- Refactor server and volume creation code
- Make maintenance_window computed


## 6.3.0

### Feature
- Adds ionoscloud_servers data source that returns a list of servers based on filters set. The filters do partial matching. See docs [here](docs/data-sources/servers.md)
- New Product: **Application Load Balancer**:
  - `Resources`:
    - [ionoscloud_application_loadbalancer](docs/resources/application_loadbalancer.md)
    - [ionoscloud_application_loadbalancer_forwarding_rule](docs/resources/application_loadbalancer_forwardingrule.md)
    - [ionoscloud_target_group](docs/resources/target_group.md)
  - `Data Sources`:
    - [ionoscloud_application_loadbalancer](docs/data-sources/application_loadbalancer.md)
    - [ionoscloud_application_loadbalancer_forwarding_rule](docs/data-sources/application_loadbalancer_forwardingrule.md)
    - [ionoscloud_target_group](docs/data-sources/target_group.md)

### Dependency-update
  - updated sdk-go version from [6.0.3](https://github.com/ionos-cloud/sdk-go/releases/tag/v6.0.3) to [v6.1.0](https://github.com/ionos-cloud/sdk-go/releases/tag/v6.1.0)

## 6.2.5
### Enhancement
- Update sdk-go dependency to v6.0.3. 
  * enable certificate pinning, by setting IONOS_PINNED_CERT env variable
- Temporarily removed `gateway_ip` and `public` fields for k8s
- Introduced error when trying to set `max_node_count` equal to `min_node_count` in `k8s_node_pool`

### Fixes 
- Crash when trying to disable `autoscaling` on `k8s_node_pool`

## 6.2.4
### Fixes
- Bug when upgrading from a v6.0.0-beta.X version to a stable one (_number is required_ error)
- Reintroduced error for image data source when finding multiple results with data source
### Enhancement
- Update sdk-go-dbaas-postgres dependency to v1.0.3

### Documentation
- Updated multiple nics under the same IP Failover example, with a [one_step example](examples/ip_failover/README.md)

## 6.2.3

### Fixes
- Do not allow empty string AvailabilityZone. Only allow "AUTO", "ZONE_1", "ZONE_2", "ZONE_3"
- Type field in server resource should be case-insensitive
- Remove deprecated image_name field on volume level from server resource
- Solve #266 crash on resource_volume when using image_alias with no image_password, or ssh_key_path

### Features
- Added `group_ids` property for `ionoscloud_user` resource. For more details refer to the [documentation](docs/resources/user.md)
- Added `groups` property for `ionoscloud_user` data source. For more details refer to the [documentation](docs/data-sources/user.md)

## 6.2.2

### Fix
- Fixed error from upgrading from 6.2.0 to 6.2.1 (version compatibility issue)

## 6.2.1

### Documentation
- Improved all the examples to be ready to use 
- Added units where missing
- Added example for adding a secondary NIC to an IP Failover
- Updated provider version to the latest release in main registry page
- Added details in [README.md](README.md) about testing locally

### Enhancement
- Add `allow_replace` to node pool resource, which allows the update of immutable node_pool fields will first
  destroy and then re-create the resource. This field should be used with care, understanding the risks.
- Update sdk-go dependency to v6.0.2
- Update sdk-go-dbaas-postgres dependency to v1.0.2
- Update terraform-plugin-sdk to v2.12.0
- Token and username+password does not conflict anymore, all three can be set, the token having priority

### Features 
- Added `backup_location` property for `ionoscloud_pg_cluster`. For more details refer to the [documentation](docs/resources/dbaas_pgsql_cluster.md)

### Fixes
- Fixed image data-source bug when `name` not provided - data-source returned 0 results
- When you try to change an immutable field, you get an error, but before that the tf state is changed. 
Before applying a real change you need to `apply` it back with an error again. 
To fix, when you try to change immutable fields they will throw an error in the plan phase.
- Reintroduced in group resource the `user_id` argument, as deprecated, to provide a period of transition
- Check slice length to prevent crash
- Fixed k8s_cluster data_source bug when searching by name 
- Fix lan deletion error, when trying to delete it immediately after the deletion of the DBaaS cluster that contained it

## 6.2.0

### Enhancement
- Modified group_resource to accept multiple users. **Note: Modify your plan according to the documentation**

## 6.1.6

### Fixes
- Fixed data sources to provide an exact match (roll-back to 6.1.3 + errors in case of multiple results)

### Documentation
- Updated k8s cluster and node pool version from examples

## 6.1.5

### Fixes
- Limit max concurrent connections to the same host to 3.
- Set max retries in case of rate-limit(429) to 999.
- Set backoff time to 4s.

### Documentation:
- Updated gitbook documentation with `legal` subheading

## 6.1.4

### Enhancements:
- Improved lookup in data_sources by using filters
- Improved tests duration by moving steps from data_source test files in the corresponding resource test files 
- Added workflow to run tests from GitHub actions 
- Split tests with build tags
- Improve http client performance and timeouts

### Documentations: 
- A more accurate example on how can the cidr be set automatically on a DBaaS Cluster
- Update doc of how to dump kube_config into a file in yaml format.

### Fixes: 
- Fix on creating a DBaaS Cluster without specifying the maintenance window
- Solve #204 - targets in nlb forwarding rule(switched to Set instead of List), lb_private_ips(set to computed), features in datacenter resources(switched to Set instead of List)
- Fix of plugin crash when updating k8s_node_pool node_count
- Fix of diff when creating a k8s_node_pool without maintenance_window

## 6.1.3

### Features:
- Added **public** parameter for k8s_cluster (creation of private clusters is possible now)
- Added **gateway_ip** parameter for k8s_nodepool
- Added **boot_server** read-only property for volume

### Fixes:
- Do not diff on gateway ips set as normal ips instead of cidr

### Enhancements:
- Terraform plugin sdk upgrade to v2.10.1
- Use depth explicitly on api calls to improve performance
- Sdk-go updated to v6.0.1


## 6.1.2

### Docs:
- Fix documentation in terraform registry

## 6.1.1

### Docs:
- Fix documentation in terraform registry 

## 6.1.0

### Features:
- New Product: **Database as a Service**: 
  - Resources: 
    - resource_dbaas_pgsql_cluster
  - Data Sources:
    - data_source_dbaas_pgsql_backups
    - data_source_dbaas_pgsql_cluster
    - data_source_dbaas_pgsql_versions
  - Dependency-update: added SDK-Go-DBaaS Postgres version [v1.0.0](https://github.com/ionos-cloud/sdk-go-dbaas-postgres/releases/tag/v1.0.0)

## 6.0.3

### Enhancements:
- Improved tests for networkloadbalancer and networkloadbalancer_forwardingrule

### Fixes:
- Fixed bug regarding updating listener_lan and target_lan on networkloadbalancer
- Fixed diff on availableUpgradeVersions for k8s cluster and nodepool
- Fixed lan deletion - wait for completion of nic deletion

### Documentation:
- Restructured documentation by adding subcategories

## 6.0.2

### Fixes:
- Fixes #168: Add versioning to allow module import.
- Modify UserAgent string

### Documentation:
- Improved terraform registry documentation with a more detailed description of environment and terraform variables
- Added badges containing the release and go version in README.md

### Fixes:
- Immutable k8s node_pool fields should throw error when running plan also, not only on apply

## 6.0.1

### Fixes: 
- Fixed rebuild k8 nodes with the same lan - order of lans is ignored now at diff
- Fixed conversion coming from a v5 state - added nil check in lans interface conversion

## 6.0.0

### Enhancements:
- Added http request time log for api calls
- Updated to go version 1.17, updated to sdk version 6.0.0
- For `k8s_node_pool`, `nic`, `ipfailover`, and `share`:
  - made tests comprehensive
  - optimized test duration by including both match by id and by name in the same test
  - removed duplicated code from import, data_source, resource and tests files
- Improved tests on natgateway and natgateway_rule

### Features:
- Import for `nic`, data_source for `nic`, `share`, `ipfailover`

### Fixes:
- K8s_node_pool update node_count and emptying lans and public_ips didn't work
- Fixed bug at creating natgateway_rule - target_subnet was not set properly
- Revert icmp_code and icmp_type to string to allow setting to 0
- Add additional fixes to improve code stability and prevent crashes. Revert icmp_type and icmp_code inside server resource and add tests.
- Allow creation of an inner firewall rule for server when updating a terraform plan.
- Fixed issue #155: added stateUpgrader for handling change of lan field structure
- Fix sporadic EOF received when making a lot of https requests to server (fixed in sdk)
- Fixed #154: allow url to start with "http" (fixed in sdk)
- Fixed #92: fix user update, user password change and password field is now sensitive
- Fix crash when no metadata is received from server

## 6.0.0-beta.14

### Fixes:
- Fixed datacenter datasource

### Enhancements:
- Added constants and removed duplicated tests to `backupUnit`, `datacenter`, `lan`, `s3_key`, `firewall`, `server`
- For `pcc`, `group`, `user`, `snapshot`, and `volume` :
  - made tests comprehensive
  - optimized test duration by including both match by id and by name in the same test
  - removed duplicated code from import, data_source, resource and tests files

### Features:
- Added import for `snapshot`, `ipblock`, data_source for `group`, `user`, `ipblock`, `volume`

## 6.0.0-beta.13

### Fixes:
- Fixed issue #112 can't attach existing volume to server after recreating server
- `cube server` could not be deleted
- Cannot empty `api_subnet_allow_list` and `s3_buckets`

### Enhancements:
- Improved data_source for template - now `template` can be searched by any of its arguments
- **code enhancements**: for `k8s_cluster`:
  - made tests comprehensive
  - optimized test duration by including both match by id and by name in the same test
  - removed duplicated code from import, data_source and resource files (set parameters)
  
## 6.0.0-beta.12
### Fixes:
- `server`: can not create cube server, firewall not updated
- `firewall`: using type argument throws error

### Enhancements:
- For `backupUnit`, `datacenter`, `lan`, `s3_key`, and `firewall` resources done the following:
  - made tests comprehensive
  - optimized test duration by including both match by id and by name in the same test
  - removed duplicated code from import, data_source and resource files (set parameters)
  - updated documentation
  - improved import functions

### Features:
- Data_source for `s3_key`

## 6.0.0-beta.11
### Fixes:
- Added `image_alias` to volume
- Removed `public` and `gateway_ip` properties from `k8s_cluster`

### Enhancements:
- Updated sdk-go to `v6.0.0-beta.7`

### Features:
- Added `data_sources for `backup_unit` and `firewall_rule`
- Added import for `natgateway`, `natgateway_rule`, `networkloadbalancer` and `networkloadbalancer_forwardingrule`

## 6.0.0-beta.10

- Issue #19 - fixed update `ssh_key_path` although not changed
- Issue #93 - updated `documentation` for image data source
- Made `backup_unit_id` configurable for volume
- Fixed `server import`

## 6.0.0-beta.9

- Issue #31 - k8s node pool labels and annotations implemented
- Ipblock `k8s_nodepool_uuid` attribute fixed
- Correctly importing private lans from k8s node pools

## 6.0.0-beta.8

- Fixed set of empty array in terraform state instead of null

## 6.0.0-beta.7

- K8s security features implemented

## 6.0.0-beta.6

- Updated arguments for datacenter, ipblock, location and user
- Issue #72 - fixed find volume image by name
- Error message for immutable node pool attributes
- Issue #84 - fixed build & updated README.md

## 6.0.0-beta.5

- Rollback to the node pool behaviour before the fix of issue #71
- Issue #77 - fix import for k8s nodepool

## 6.0.0-beta.4

- Fix: issue #71 - recreate nodepool on change of specifications

## 6.0.0-beta.3

- Issue #66 - detailed kube config attributes implemented

## 6.0.0-beta.2

- Updated dependencies 
- Updated server, nic and volume resources with the missing arguments

## 6.0.0-beta.1

- Documentation updates
- Fix: fixes #13 ignore changes of patch level for k8s

## 6.0.0-alpha.4

- Documentation updates

## Enhancements:
- Terraform plugin sdk upgrade to v2.4.3
- Fix: create volume without password
- Fix: ability to create server without image
- Fix: fixes #25 correctly set of dhcp + nil check + added firewall_type field in server resource
- Fix: fixes #39 - new imports for volume, user, group, share, IPfailover and loadbalancer
- Fix: fixes #47 - corrected nic resource to accept a list of strings for ips parameter
- Fix: fixes #36 - correctly setting the active property of the s3 key upon creation

## 6.0.0-alpha.3

- Documentation updates

## 6.0.0-alpha.2

- IONOS_DEBUG env var support for debugging sdk/api request payloads
- Fix: contract number correctly computed when generating backup-unit names
- Fix: segfault avoided on missing volume image
- Test suite improvements

## 6.0.0-alpha.1

- Initial v6 version supporting Ionos Cloud API v6

## 5.1.6

- Fixes #5 - correctly dereferencing possibly nil properties received from the api

## 5.1.5

- Fixes #12 - correctly setting up a custom Ionos Cloud API url

## 5.1.4

- Error handling improvements 
- Always displaying the full response body from the API in case of an error

## 5.1.3

- Bug fix: correctly checking for nil the image volume 

## 5.1.2

- Bug fix: avoid sending an empty image password to the API if 
  no image password is set

## 5.1.1

- Bug fix: nil check for image password when creating a server 

## 5.1.0

- Using the latest Ionos Cloud GO SDK v5.1.0

## 5.0.4

BUG FIXES:
- Importing mac info when loading nic information or server information
- Reading PCC info when importing a lan

## 5.0.3

FEATURES:
- New data sources added: k8s_cluster, k8s_node_pool

## 5.0.2

BUG FIXES:

- Correctly updating ips on a nic embedded in a server config 

## 5.0.1

FEATURES:
- New datasources added: lan, server, private cross connect

## 5.0.0

FEATURES:
- Terraform-provider-profitbricks rebranding to terraform-provider-ionoscloud

