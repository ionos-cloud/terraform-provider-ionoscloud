## 5.2.19

- **code enhancements**: added http request time log for API calls
- **dependency update**: updated sdk-go to v5.1.9
- **code enhancements**: for `k8s_node_pool`, `nic` and `share`:
  - made tests comprehensive
  - optimized test duration by including both match by id and by name in the same test
  - removed duplicated code from import, data_source, resource and tests files
- **new features**: import for `nic`, data_source for `nic`, `share`
- **bug fixes**: k8s_node_pool update node_count didn't work and emptying lans and public_ips. revert icmp_code and icmp_type to string to allow setting to 0

## 5.2.18

- **bug fixes**: fixed datacenter datasource
- **code enhancements**: added constants and removed duplicated tests to `backupUnit`, `datacenter`, `lan`, `s3_key`, `firewall`, `server`
- **code enhancements**: for `pcc`, `group`, `user`,`snapshot`, `volume` and `server`:
  - made tests comprehensive 
  - optimized test duration by including both match by id and by name in the same test
  - removed duplicated code from import, data_source, resource and tests files
- **new features**: import for `snapshot`, `ipblock`, data_source for `group`, `user`, `volume`, `ipblock`

## 5.2.17

- issue #31 - k8s node pool labels and annotations implemented
- **bug fixes**: fixed issue #112 can't attach existing volume to server after recreating server 
- **bug fixes**: cannot empty `api_subnet_allow_list` and `s3_buckets`
- **code enhancements**: for `k8s_cluster`:
  - made tests comprehensive
  - optimized test duration by including both match by id and by name in the same test
  - removed duplicated code from import, data_source and resource files (set parameters)

## 5.2.16

- **code enhancements**: for `backupUnit`, `datacenter`, `lan`, `s3_key`, and `firewall` resources done the following:
  - made tests comprehensive
  - optimized test duration by including both match by id and by name in the same test
  - removed duplicated code from import, data_source and resource files (set parameters)
  - updated documentation
  - improved import functions
- **new features**: data_source for `s3_key`

## 5.2.15

- updated sdk-go to v5.1.7
- implemented data source for backup unit and firewall rule
- removed public and gateway_ip properties from k8s cluster

## 5.2.14

- fixed typo in setting user_data and backup_unit_id in the volume entity from server
- test updates

## 5.2.13

- added user_data and backup_unit_id in the volume entity from server
- fix issue #19 - fixed update ssh_key_path although not changed
- issue #93 - updated documentation for image data source

## 5.2.12

- fix: correctly saving lans when reading a k8s node pool

## 5.2.11

- documentation updates

## 5.2.10

- fixed set of empty array in terraform state instead of null

## 5.2.9

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
- fix: fixes #1 - usage example updates
- documentation updates
- fix: fixes #13 ignore changes of patch level in k8s cluster & nodepool k8sVersion
- set public default to true to remove deprecated GetOkExists function
- added some missing arguments
- API kubernetes security featues implemented (apiSubnetAllowList and S3Buckets) 
- fixed import server import

## 5.2.4

- issue #47 - corrected nic resource to accept a list of strings for ips parameter

## 5.2.3

- issue #39 - new imports for volume, user, group, share, IPfailover and loadbalancer

## 5.2.2

- issue #36 - correctly setting the value of the active property when creating an s3 resource

## 5.2.1

- issue #29 - corrected parameter name in volume error message
- issue #30 - creation of volume without password + default value for bus

## 5.2.0

- fixes #17 - documentation updates

## 5.2.0-beta.2

- fixes #24 - ability to create servers without an image

## 5.2.0-beta.1

- terraform sdk upgrade to v2.4.3

## 5.1.7

- fixes #22 - ability to specify boot_cdrom when creating a server
- fix: respecting resource timeouts when waiting for requests to be fullfiled
- ability to debug sdk requests by setting the IONOS_DEBUG=1 env var and TF_LOG=1


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

