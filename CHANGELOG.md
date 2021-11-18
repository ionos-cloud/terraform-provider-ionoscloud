## 6.0.0-beta.15

- **code enhancements**: added http request time log for api calls
- **new features**: import for `nic`, data_source for `nic`, `share`, `ipfailover`
- **dependency update**: updated sdk-go to v6.0.0-beta.8
- **tests enhancements**: improved tests on natgateway and natgateway_rule
- **code enhancements**: for `k8s_node_pool`, `nic`, `ipfailover`, and `share`:
  - made tests comprehensive
  - optimized test duration by including both match by id and by name in the same test
  - removed duplicated code from import, data_source, resource and tests files
- **bug fixes**: k8s_node_pool update node_count and emptying lans and public_ips didn't work
- **bug fixes**: fixed bug at creating natgateway_rule - target_subnet was not set properly
- **bug fixes**: revert icmp_code and icmp_type to string to allow setting to 0

## 6.0.0-beta.14

- **bug fixes**: fixed datacenter datasource
- **code enhancements**: added constants and removed duplicated tests to `backupUnit`, `datacenter`, `lan`, `s3_key`, `firewall`, `server`
- **code enhancements**: for `pcc`, `group`, `user`, `snapshot`, and `volume` :
  - made tests comprehensive
  - optimized test duration by including both match by id and by name in the same test
  - removed duplicated code from import, data_source, resource and tests files
- **new features**: import for `snapshot`, `ipblock`, data_source for `group`, `user`, `ipblock`, `volume`

## 6.0.0-beta.13
- **bug fixes**: fixed issue #112 can't attach existing volume to server after recreating server
- **bug fixes**: `cube server` could not be deleted
- **functionality enhancements**: improved data_source for template - now `template` can be searched by any of its arguments
- **bug fixes**: cannot empty `api_subnet_allow_list` and `s3_buckets`
- **code enhancements**: for `k8s_cluster`:
  - made tests comprehensive
  - optimized test duration by including both match by id and by name in the same test
  - removed duplicated code from import, data_source and resource files (set parameters)
  
## 6.0.0-beta.12

- **bug fixes**: `server`: can not create cube server, firewall not updated
- **bug fixes**: `firewall`: using type argument throws error
- **code enhancements**: for `backupUnit`, `datacenter`, `lan`, `s3_key`, and `firewall` resources done the following:
  - made tests comprehensive
  - optimized test duration by including both match by id and by name in the same test
  - removed duplicated code from import, data_source and resource files (set parameters)
  - updated documentation
  - improved import functions
- **new features**: data_source for `s3_key`

## 6.0.0-beta.11

- added `image_alias` to volume
- removed `public` and `gateway_ip` properties from `k8s_cluster`
- added `data_sources for `backup_unit` and `firewall_rule`
- added import for `natgateway`, `natgateway_rule`, `networkloadbalancer` and `networkloadbalancer_forwardingrule`
- updated sdk-go to `v6.0.0-beta.7`

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

