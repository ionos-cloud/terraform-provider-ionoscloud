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

