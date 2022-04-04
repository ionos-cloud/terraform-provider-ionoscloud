## 6.2.1

### Documentation
- improved all the examples to be ready to use 
- added units where missing
- added example for adding a secondary NIC to an IP Failover
- updated provider version to the latest release in main registry page
- added details in [README.md](README.md) about testing locally

### Enhancement
- add `allow_replace` to node pool resource, which allows the update of immutable node_pool fields will first
  destroy and then re-create the resource. This field should be used with care, understanding the risks.
- update sdk-go dependency to v6.0.2
- update sdk-go-dbaas-postgres dependency to v1.0.2
- update terraform-plugin-sdk to v2.12.0
- token and username+password does not conflict anymore, all three can be set, the token having priority

### Features 
- added `backup_location` property for `ionoscloud_pg_cluster`. For more details refer to the [documentation](/docs/resources/dbaas_pgsql_cluster.md)

### Fixes
- fixed image data-source bug when `name` not provided - data-source returned 0 results
- when you try to change an immutable field, you get an error, but before that the tf state is changed. 
Before applying a real change you need to `apply` it back with an error again. 
To fix, when you try to change immutable fields they will throw an error in the plan phase.
- reintroduced in group resource the `user_id` argument, as deprecated, to provide a period of transition
- check slice length to prevent crash
- fixed k8s_cluster data_source bug when searching by name 
- fix lan deletion error, when trying to delete it immediately after the deletion of the DBaaS cluster that contained it

## 6.2.0

### Enhancement
- modified group_resource to accept multiple users. **Note: Modify your plan according to the documentation**

## 6.1.6

### Fixes
- fixed data sources to provide an exact match (roll-back to 6.1.3 + errors in case of multiple results)

### Documentation
- updated k8s cluster and node pool version from examples

## 6.1.5

### Fixes
- Limit max concurrent connections to the same host to 3.
- Set max retries in case of rate-limit(429) to 999.
- Set backoff time to 4s.

### Documentation:
- updated gitbook documentation with `legal` subheading

## 6.1.4

### Enhancements:
- improved lookup in data_sources by using filters
- improved tests duration by moving steps from data_source test files in the corresponding resource test files 
- added workflow to run tests from GitHub actions 
- split tests with build tags
- improve http client performance and timeouts

### Documentations: 
- a more accurate example on how can the cidr be set automatically on a DBaaS Cluster
- update doc of how to dump kube_config into a file in yaml format.

### Fixes: 
- fix on creating a DBaaS Cluster without specifying the maintenance window
- solve #204 - targets in nlb forwarding rule(switched to Set instead of List), lb_private_ips(set to computed), features in datacenter resources(switched to Set instead of List)
- fix of plugin crash when updating k8s_node_pool node_count
- fix of diff when creating a k8s_node_pool without maintenance_window

## 6.1.3

### Features:
- added **public** parameter for k8s_cluster (creation of private clusters is possible now)
- added **gateway_ip** parameter for k8s_nodepool
- added **boot_server** read-only property for volume

### Fixes:
- do not diff on gateway ips set as normal ips instead of cidr

### Enhancements:
- terraform plugin sdk upgrade to v2.10.1
- use depth explicitly on api calls to improve performance
- sdk-go updated to v6.0.1


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
- improved tests for networkloadbalancer and networkloadbalancer_forwardingrule

### Fixes:
- fixed bug regarding updating listener_lan and target_lan on networkloadbalancer
- fixed diff on availableUpgradeVersions for k8s cluster and nodepool
- fixed lan deletion - wait for completion of nic deletion

### Documentation:
- restructured documentation by adding subcategories

## 6.0.2

### Fixes:
- fixes #168: Add versioning to allow module import.
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
- added http request time log for api calls
- updated to go version 1.17, updated to sdk version 6.0.0
- for `k8s_node_pool`, `nic`, `ipfailover`, and `share`:
  - made tests comprehensive
  - optimized test duration by including both match by id and by name in the same test
  - removed duplicated code from import, data_source, resource and tests files
- improved tests on natgateway and natgateway_rule

### Features:
- import for `nic`, data_source for `nic`, `share`, `ipfailover`

### Fixes:
- k8s_node_pool update node_count and emptying lans and public_ips didn't work
- fixed bug at creating natgateway_rule - target_subnet was not set properly
- revert icmp_code and icmp_type to string to allow setting to 0
- Add additional fixes to improve code stability and prevent crashes. Revert icmp_type and icmp_code inside server resource and add tests.
- Allow creation of an inner firewall rule for server when updating a terraform plan.
- fixed issue #155: added stateUpgrader for handling change of lan field structure
- fix sporadic EOF received when making a lot of https requests to server (fixed in sdk)
- fixed #154: allow url to start with "http" (fixed in sdk)
- fixed #92: fix user update, user password change and password field is now sensitive
- fix crash when no metadata is received from server

## 6.0.0-beta.14

### Fixes:
- fixed datacenter datasource

### Enhancements:
- added constants and removed duplicated tests to `backupUnit`, `datacenter`, `lan`, `s3_key`, `firewall`, `server`
- for `pcc`, `group`, `user`, `snapshot`, and `volume` :
  - made tests comprehensive
  - optimized test duration by including both match by id and by name in the same test
  - removed duplicated code from import, data_source, resource and tests files

### Features:
- added import for `snapshot`, `ipblock`, data_source for `group`, `user`, `ipblock`, `volume`

## 6.0.0-beta.13

### Fixes:
- fixed issue #112 can't attach existing volume to server after recreating server
- `cube server` could not be deleted
- cannot empty `api_subnet_allow_list` and `s3_buckets`

### Enhancements:
- improved data_source for template - now `template` can be searched by any of its arguments
- **code enhancements**: for `k8s_cluster`:
  - made tests comprehensive
  - optimized test duration by including both match by id and by name in the same test
  - removed duplicated code from import, data_source and resource files (set parameters)
  
## 6.0.0-beta.12
### Fixes:
- `server`: can not create cube server, firewall not updated
- `firewall`: using type argument throws error

### Enhancements:
- for `backupUnit`, `datacenter`, `lan`, `s3_key`, and `firewall` resources done the following:
  - made tests comprehensive
  - optimized test duration by including both match by id and by name in the same test
  - removed duplicated code from import, data_source and resource files (set parameters)
  - updated documentation
  - improved import functions

### Features:
- data_source for `s3_key`

## 6.0.0-beta.11
### Fixes:
- added `image_alias` to volume
- removed `public` and `gateway_ip` properties from `k8s_cluster`

### Enhancements:
- updated sdk-go to `v6.0.0-beta.7`

### Features:
- added `data_sources for `backup_unit` and `firewall_rule`
- added import for `natgateway`, `natgateway_rule`, `networkloadbalancer` and `networkloadbalancer_forwardingrule`

## 6.0.0-beta.10

- issue #19 - fixed update `ssh_key_path` although not changed
- issue #93 - updated `documentation` for image data source
- made `backup_unit_id` configurable for volume
- fixed `server import`

## 6.0.0-beta.9

- issue #31 - k8s node pool labels and annotations implemented
- ipblock `k8s_nodepool_uuid` attribute fixed
- correctly importing private lans from k8s node pools

## 6.0.0-beta.8

- fixed set of empty array in terraform state instead of null

## 6.0.0-beta.7

- k8s security features implemented

## 6.0.0-beta.6

- updated arguments for datacenter, ipblock, location and user
- issue #72 - fixed find volume image by name
- error message for immutable node pool attributes
- issue #84 - fixed build & updated README.md

## 6.0.0-beta.5

- rollback to the node pool behaviour before the fix of issue #71
- issue #77 - fix import for k8s nodepool

## 6.0.0-beta.4

- fix: issue #71 - recreate nodepool on change of specifications

## 6.0.0-beta.3

- issue #66 - detailed kube config attributes implemented

## 6.0.0-beta.2

- updated dependencies 
- updated server, nic and volume resources with the missing arguments

## 6.0.0-beta.1

- documentation updates
- fix: fixes #13 ignore changes of patch level for k8s

## 6.0.0-alpha.4

- documentation updates

## Enhancements:
- terraform plugin sdk upgrade to v2.4.3
- fix: create volume without password
- fix: ability to create server without image
- fix: fixes #25 correctly set of dhcp + nil check + added firewall_type field in server resource
- fix: fixes #39 - new imports for volume, user, group, share, IPfailover and loadbalancer
- fix: fixes #47 - corrected nic resource to accept a list of strings for ips parameter
- fix: fixes #36 - correctly setting the active property of the s3 key upon creation

## 6.0.0-alpha.3

- documentation updates

## 6.0.0-alpha.2

- IONOS_DEBUG env var support for debugging sdk/api request payloads
- fix: contract number correctly computed when generating backup-unit names
- fix: segfault avoided on missing volume image
- test suite improvements

## 6.0.0-alpha.1

- initial v6 version supporting Ionos Cloud API v6

## 5.1.6

- fixes #5 - correctly dereferencing possibly nil properties received from the api

## 5.1.5

- fixes #12 - correctly setting up a custom Ionos Cloud API url

## 5.1.4

- error handling improvements 
- always displaying the full response body from the API in case of an error

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
- new data sources added: k8s_cluster, k8s_node_pool

## 5.0.2

BUG FIXES:

- Correctly updating ips on a nic embedded in a server config 

## 5.0.1

FEATURES:
- new datasources added: lan, server, private cross connect

## 5.0.0

FEATURES:
- terraform-provider-profitbricks rebranding to terraform-provider-ionoscloud

