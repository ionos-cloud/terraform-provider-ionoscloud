## 5.2.24 (upcoming release)

### Features:
- add data source for ip failover

### Documentation:
- Improved terraform registry documentation with a more detailed description of environment and terraform variables
- Added badges containing the release and go version in README.md
- Add note on using v6 to Readme.md

### Fixes:
- Creating server with volume from snapshot did not populate volume_boot
- Primary_ip is now set on server creation
- Add versioning to go.mod to allow version import of module
- Immutable k8s node_pool fields should throw error when running plan also, not only on apply

## 5.2.23

### Fixes: 
- Fixed rebuild k8 nodes with the same lan - order of lans is ignored now at diff


## 5.2.22

### Enhancements:
- Update sdk to version v5.1.11. 
- Update go to version 1.17
- Update terraform-plugin-sdk to v2.9.0

### Fixes:
- Password now saved for user on update
- fix sporadic EOF received when making a lot of concurrent https requests to server (fixed in sdk 6.0.0)
- fixed #154: allow url to start with "http" (fixed in sdk v5.1.1)
- fixed #92: user password change (fixed in sdk v5.1.1)
- fix user update and password field is now sensitive
- fix crash when no metadata is received from server

## 5.2.21

### Fixes:
- Add additional fixes to improve code stability and prevent crashes. Revert icmp_type and icmp_code inside server resource and add tests.
- Allow creation of an inner firewall rule for server when updating a terraform plan. 

## 5.2.20

### Fixes:
- fix crash and add additional logs 

## 5.2.19

### Enhancements:
- added http request time log for API calls
- updated sdk-go to v5.1.9
- for `k8s_node_pool`, `nic` and `share`:
  - made tests comprehensive
  - optimized test duration by including both match by id and by name in the same test
  - removed duplicated code from import, data_source, resource and tests files

### Features:
- import for `nic`, data_source for `nic`, `share`

### Fixes 
- k8s_node_pool update node_count didn't work and emptying lans and public_ips. revert icmp_code and icmp_type to string to allow setting to 0

## 5.2.18

### Fixes:
- fixed datacenter datasource

### Enhancements:
- added constants and removed duplicated tests to `backupUnit`, `datacenter`, `lan`, `s3_key`, `firewall`, `server`
- for `pcc`, `group`, `user`,`snapshot`, `volume` and `server`:
  - made tests comprehensive 
  - optimized test duration by including both match by id and by name in the same test
  - removed duplicated code from import, data_source, resource and tests files
  
### Features:
- import for `snapshot`, `ipblock`, data_source for `group`, `user`, `volume`, `ipblock`

## 5.2.17

### Fixes:
- issue #31 - k8s node pool labels and annotations implemented
- fixed issue #112 can't attach existing volume to server after recreating server 
- cannot empty `api_subnet_allow_list` and `s3_buckets`

### Enhancements:
  - made tests comprehensive
  - optimized test duration by including both match by id and by name in the same test
  - removed duplicated code from import, data_source and resource files (set parameters)

## 5.2.16

### Enhancements:
- for `backupUnit`, `datacenter`, `lan`, `s3_key`, and `firewall` resources done the following:
  - made tests comprehensive
  - optimized test duration by including both match by id and by name in the same test
  - removed duplicated code from import, data_source and resource files (set parameters)
  - updated documentation
  - improved import functions

### Features:
-  data_source for `s3_key`

## 5.2.15

### Features:
- updated sdk-go to v5.1.7
- implemented data source for backup unit and firewall rule
- removed public and gateway_ip properties from k8s cluster

## 5.2.14

### Fixes:
- fixed typo in setting user_data and backup_unit_id in the volume entity from server
- test updates

## 5.2.13

### Features:
- added user_data and backup_unit_id in the volume entity from server

### Fixes:
- fix issue #19 - fixed update ssh_key_path although not changed
- issue #93 - updated documentation for image data source

## 5.2.12

### Fixes:

-  correctly saving lans when reading a k8s node pool

## 5.2.11

- documentation updates

## 5.2.10

- fixed set of empty array in terraform state instead of null

## 5.2.9

### Fixes:
- issue #72 - fixed find volume image by name
- error messages for immutable node pool attributes
- issue #84 - fixed build & updated README.md

## 5.2.8

- rollback to the node pool behaviour before the fix of issue #71
- issue #77 - fix import for k8s nodepool

## 5.2.7

- fix: issue #71 - recreate nodepool on change of specifications

## 5.2.6

- issue #66 - detailed kube config attributes implemented 

## 5.2.5

### Fixes:

- fix: fixes #1 - usage example updates
- fix: fixes #13 ignore changes of patch level in k8s cluster & nodepool k8sVersion
- added some missing arguments
- fixed import server import

### Enhancements:
- documentation updates
- set public default to true to remove deprecated GetOkExists function
- API kubernetes security featues implemented (apiSubnetAllowList and S3Buckets) 

## 5.2.4

- issue #47 - corrected nic resource to accept a list of strings for ips parameter

## 5.2.3

- issue #39 - new imports for volume, user, group, share, IPfailover and loadbalancer

## 5.2.2

- issue #36 - correctly setting the value of the active property when creating an s3 resource

## 5.2.1

### Fixes:
- issue #29 - corrected parameter name in volume error message
- issue #30 - creation of volume without password + default value for bus

## 5.2.0

- fixes #17 - documentation updates

## 5.2.0-beta.2

### Fixes:
- fixes #24 - ability to create servers without an image

## 5.2.0-beta.1

### Enhancements:
- terraform sdk upgrade to v2.4.3

## 5.1.7

### Fixes:
- fixes #22 - ability to specify boot_cdrom when creating a server
- fix: respecting resource timeouts when waiting for requests to be fullfiled

### Enhancements:
- ability to debug sdk requests by setting the IONOS_DEBUG=1 env var and TF_LOG=1


## 5.1.6

- fixes #5 - correctly dereferencing possibly nil properties received from the api

## 5.1.5

### Fixes:
- fixes #12 - correctly setting up a custom Ionos Cloud API url

## 5.1.4

### Fixes:
- error handling improvements 
- always displaying the full response body from the API in case of an error

## 5.1.3

### Fixes:
- correctly checking for nil the image volume 

## 5.1.2

### Fixes:
- avoid sending an empty image password to the API if 
  no image password is set

## 5.1.1

- Bug fix: nil check for image password when creating a server 

## 5.1.0

- Using the latest Ionos Cloud GO SDK v5.1.0

## 5.0.4

### Fixes:
- Importing mac info when loading nic information or server information
- Reading PCC info when importing a lan

## 5.0.3

### FEATURES:
- new data sources added: k8s_cluster, k8s_node_pool

## 5.0.2

### Fixes:
- Correctly updating ips on a nic embedded in a server config 

## 5.0.1

### FEATURES:
- new datasources added: lan, server, private cross connect

## 5.0.0

### FEATURES:
- terraform-provider-profitbricks rebranding to terraform-provider-ionoscloud

