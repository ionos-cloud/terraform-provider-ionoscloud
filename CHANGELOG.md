## 5.2.3

- issue #43 - new imports for volume, user, group, share, IPfailover and loadbalancer
- issue #47 - corrected nic resource to accept a list of strings for ips parameter 

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

