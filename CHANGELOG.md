## 6.1.0

### Feature:
- Database as a Service: 
  - Resources: 
    - resource_dbaas_pgsql_cluster
  - Data Sources:
    - data_source_dbaas_pgsql_backups
    - data_source_dbaas_pgsql_cluster
    - data_source_dbaas_pgsql_versions

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

