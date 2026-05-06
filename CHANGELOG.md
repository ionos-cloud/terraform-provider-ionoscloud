## Upcoming

### Documentation
- Improve share docs

### Chore
- Modernize to 1.26 standards `go fix`

### Enhancements
- Add validation error when update-only attributes are set during snapshot creation
- Add more information to error messages
- Use tflog instead of print

### Testing 
- Add more test cases for `compute` data sources.

## 6.7.28
### Fixes
- Fix possible type assertion panic in `location` and `template` data sources.

### Dependencies
- Bump google.golang.org/grpc from 1.72.1 to 1.79.3 to mitigate vulnerability.

## 6.7.27
### Features
- Add `nic_multi_queue` attribute to the `ionoscloud_servers` data source.
- Add DBaaS PostgreSQL v2 support:
  - New resource: `ionoscloud_pg_cluster_v2`
  - New data sources: `ionoscloud_pg_cluster_v2`, `ionoscloud_pg_clusters_v2`, `ionoscloud_pg_backups_v2`, `ionoscloud_pg_versions_v2`, `ionoscloud_pg_backup_location_v2`

### Testing
- Add checks for the `nic_multi_queue` attribute inside VCPU servers tests.
- Modify `ionoscloud_servers` data source test to include a check for the `nic_multi_queue` attribute.

## 6.7.26
### Fixes
- The ionoscloud_s3_bucket_policy resource now correctly handles all standard S3 Principal representations:
  - "Principal": "*" — wildcard string
  - "Principal": ["arn:...", "*"] — flat array
  - "Principal": {"AWS": "arn:..."} — object with single string
  - "Principal": {"AWS": ["arn:...", "arn:..."]} — object with array
- Make DNS Record name immutable, fixes [#953](https://github.com/ionos-cloud/terraform-provider-ionoscloud/pull/953)
- Improved error message when `image_name` matches an image that has a non-HDD type (e.g. CDROM) or is in a different location. The error now reports the found image's type and location.

### Testing
- Add import tests for CUBE servers.

### Docs
- Update documentation for `ionoscloud_nsg_firewallrule`, `ionoscloud_natgateway` and `ionoscloud_ipfailover` resources

### Enhancements
- Add `sbom` to `release` workflow

## 6.7.25
### Features
- File configuration failover for global resources: users, groups, target groups, s3keys, backupunits, contracts, object storage access keys
### Fixes
- Add debug log with error in case config file cannot be read due to whatever reason
### Dependencies
- Use shared v0.1.8

## 6.7.24
### Features
- Added `location` field to several resources and data sources to support regional endpoints via file configuration. The field should only be used if a file configuration is provided (e.g. at `IONOS_CONFIG_FILE`).
    - Updated Resources: `application_loadbalancer_forwardingrule`, `container_registry_token`, `cube_server`, `datacenter_nsg_selection`, `dbaas_mongodb_user`, `dbaas_pgsql_database`, `dbaas_pgsql_user`, `firewall`, `gpu_server`, `k8s_node_pool`, `lan`, `loadbalancer`, `natgateway`, `natgateway_rule`, `networkloadbalancer`, `networkloadbalancer_forwardingrule`, `nic`, `nsg`, `nsg_firewallrule`, `private_crossconnect`, `server`, `server_boot_device_selection`, `snapshot`, `vcpu_server`, `volume`.
    - Updated Data Sources: `application_loadbalancer`, `application_loadbalancer_forwardingrule`, `container_registry_token`, `dbaas_mongo_user`, `dbaas_pgsql_backups`, `dbaas_pgsql_database`, `dbaas_pgsql_databases`, `dbaas_pgsql_user`, `firewall`, `gpu`, `gpus`, `ipfailover`, `k8s_clusters`, `k8s_node_pool`, `k8s_node_pool_nodes`, `lan`, `natgateway`, `natgateway_rule`, `networkloadbalancer`, `networkloadbalancer_forwardingrule`, `nic`, `nsg`, `private_crossconnect`, `server`, `servers`, `vcpu_server`, `volume`.
- Added `ForceNew: true` for the following attributes:
    - `location` inside `datacenter` resource;
    - `datacenter_id` inside `datacenter_nsg_selection` resource;
    - `datacenter_id` inside `k8s_node_pool` resource;
    - `k8s_cluster_id` inside `k8s_node_pool` resource;
    - `location` inside `snapshot` resource;
- Added `Optional: true` for the following attributes:
    - `location` inside `snapshot` resource;
    - `location` inside `k8s_cluster` resource;

## 6.7.23
### Features
- Added Kafka users data source: `ionoscloud_kafka_users`, user access credentials data source: `ionoscloud_kafka_user_credentials` and user access credentials ephemeral resource: `ionoscloud_kafka_user_credentials`.
### Breaking Changes
- **Removed API Gateway resources and data sources**: `ionoscloud_apigateway`, `ionoscloud_apigateway_route` and their corresponding data sources have been removed from the provider. All API Gateway related code, tests, documentation, and examples have been removed.
### Documentation
- Updated documentation for S3 buckets.
- Removed API Gateway documentation and examples.
### Fixes
- Add back missing `id` field in `ionoscloud_server.nic.firewall` sub-resource, which caused an error when trying to create a server with firewall rules defined on a nic.
- Remove unused `password` field from `users` attribute in `ionoscloud_group` resource schema.
- Fixed VPN IPSec Gateway to use correct `fileconfiguration.VPN` constant instead of incorrect `fileconfiguration.APIGateway`.
- Fixed typo in VPN IPSec Gateway data source documentation heading.
- Fixed `logging_format` field in `ionoscloud_networkloadbalancer` resource causing perpetual drift by adding `Computed: true` to the schema ([#918](https://github.com/ionos-cloud/terraform-provider-ionoscloud/issues/918)).

## 6.7.22
### Features
- Added Reverse DNS resource and data sources: `resource_dns_reverse_record`, `data_source_dns_reverse_record` and `data_source_dns_reverse_records`
- Added resource for GPU type servers: `ionoscloud_gpu_server`
- Added datasource for GPU type servers: `ionoscloud_gpu_server`
- Added datasources for GPU / GPUs: `ionoscloud_gpu`, `ionoscloud_gpus`
- Added support for the GPU fields within the existing template datasource.
### Fixes
- Fix an issue related to the IPs order for a NIC that led to continuous changes: [#903](https://github.com/ionos-cloud/terraform-provider-ionoscloud/issues/903)

## 6.7.21
### Features
- Add `require_legacy_bios` option for volumes, snapshots, images and inline volumes (for all server types)
### Fixes
- Fix the way in which `login` value is computed for backup unit resource (fixes [#898](https://github.com/ionos-cloud/terraform-provider-ionoscloud/issues/898))
### Chore
- Update GO version to 1.25.3
### Documentation
- Specify how to work with multiple resources of the same type but in different locations/regions, explain regions logic for S3 buckets

## 6.7.20
### Fixes
- Fix bug introduced by `6.7.19` for `VCPU` type server resource and data-source

## 6.7.19
### Features
- Add new `nic_multi_queue` feature for `ENTERPRISE` servers.
### Fixes
- Fix #881 error on s3 key creation immediately after user and group. "The user needs to be part of a group that has ACCESS_S3_OBJECT_STORAGE privilege"
### Documentation
- Add reference to IONOS Object Storage documentation inside S3 bucket policy doc, update documentation for `IPv4`, `IPv6` addresses for `VPN` gateways, remove extra colons

## 6.7.18
### Fixes
- Fix provider crash on `terraform plan / terraform refresh` after the external deletion of a server in a boot device selection resource

## 6.7.17
### Features
- Add `get_users_data` attribute to the `ionoscloud_group` resource  and data source, this makes fetching user details optional to prevent performance issues in environments with many users or groups.
### Fixes
- Fix [#872](https://github.com/ionos-cloud/terraform-provider-ionoscloud/issues/872)

## 6.7.16
### Fixes
- Fix [#867](https://github.com/ionos-cloud/terraform-provider-ionoscloud/issues/867)
- Fix [#864](https://github.com/ionos-cloud/terraform-provider-ionoscloud/issues/864)
- Fix tests for `ionoscloud_pg_cluster` by modifying the PgSQL version
- Fix functionality for S3 keys (read & delete bug)

## 6.7.15
### Features
- Add support for the new S3 locations: `eu-central-4` and `us-central-1`
## 6.7.14
### Features
- Add support for new location for NFS: `de/fra/2`
### Refactor
- Remove useless checks for some resources
- Modify error messages for some resources for more clarity
- Remove Dataplatform service
### Documentation
- Fix documentation for DBaaS PgSQL cluster resource

## 6.7.13
### Features
- Add support for new location: `de/fra/2`
### Refactor
- Use API Gateway bundle product instead of `sdk-go-api-gateway`

## 6.7.12
### Fixes
- Fix `ionoscloud_server_boot_device_selection` delete error when servers have no initial boot device and a cdrom is used. (https://github.com/ionos-cloud/terraform-provider-ionoscloud/issues/843)
### Documentation
- Update the format of Attaching a NSG to a Datacenter documentation in the resources datacenter.
- Update Kafka documentation to state that version 3.9.0 is supported in all relevant resources and data sources.
### Fixes
- Update go version to 1.24.4

## 6.7.11
### Fixes
- Add `expose_serial` attribute to `ionoscloud_vcpu_server` resource and data source.
- Make `expose_serial` attribute computed for `ionoscloud_server` `ionoscloud_volume` and `ionoscloud_cube_server` resources.
### Documentation
- Fix `ionoscloud_logging_pipeline` resource documentation.
## 6.7.10
### Features
- Add `key` to `ionoscloud_logging_pipeline` resource
- Add `tcp_address` and`http_address` to `ionoscloud_logging_pipeline` resource and data source
- Add `expose_serial` attribute to `ionoscloud_volume`, `ionoscloud_server`, `ionoscloud_cube_server` resources. If set to `true` will expose the serial id of the disk attached to the server.
Some operating systems or software solutions require the serial id to be exposed to work properly. Exposing the serial can influence licensed software (e.g. Windows) behavior

## 6.7.9
### Features
- Add `filter` object with `prefix` field so prefix is actually set for lifecycle rules in `ionoscloud_s3_bucket_lifecycle_configuration` resource. Deprecate `prefix` field as it does nothing..
### Fixes
- Terraform provider not working with latest Terraform version 1.12.2
## 6.7.8
### Chore
- Update `github.com/hashicorp/terraform-plugin-sdk/v2` dependency to latest version
### Documentation
- Remove mentions of `kafka` 3.7.0 version as it is no longer supported.
- Update import method for `ionoscloud_datacenter` and `ionoscloud_server` resources.
### Fixes
- Fix S3 key test, add back update test
- `ionoscloud_logging_pipeline` crashes on read after create
- Retry if `ionoscloud_lan` is delete protected by a managed service
### Features
- Move inmemorydb resources to sdk-go-bundle

## 6.7.7
### Features
- Add backup to mariadb
- Add `ionoscloud_contracts` data source
### Fixes
- Error on `FAILED` for `cdn`, `api_gateway`, `dataplatform`, `inmemorydb`, `mariadb`, `nfs`, `mongo` resources
- Fix [#813](https://github.com/ionos-cloud/terraform-provider-ionoscloud/issues/813) by adding `ForceNew: true` for all attributes for `ionoscloud_pg_database` resource
- Instantiate new http clients for each product (instead of http.Default client) so that insecure mode and certificate can be set independently
- Retrieve accessKey and secretKey from the config file if they are not found
- Check if url is found in locationToURL before overwriting endpoint for certificate manager
- Keep the VM autoscaling group in state if the action fails
- Fix VM Autoscaling tests and examples
### Refactor
- Use VM Autoscaling bundle product instead of sdk-go-vm-autoscaling
- Use MariaDB bundle product instead of sdk-go-dbaas-mariadb
- Use Mongo bundle product instead of sdk-go-dbaas-mongo
- Use Dataplatform bundle product instead of sdk-go-dataplatform
- Use PostgreSQL bundle product instead of sdk-go-dbaas-postgres
- Update `sdk-go-dbaas-postgres` to version `v1.1.4` and use the new sdk methods for PG versions (done before bundle integration)
### Documentation
- Add better documentation for dedicated core for servers
- Add example for `user` field `password_wo` write only field.
## 6.7.6
### Features
- Add `server_type` optional attribute to `ionoscloud_k8s_nodepool` resource and data source
- Add `password_wo` and `password_wo_version` to `ionoscloud_user` resource. Write only field that is not stored in state. Can be used only with Terraform 1.11 or higher.
### Changed
- `cpu_family` is now optional for `ionoscloud_k8s_nodepool` resource
### Refactor
- Use Kafka bundle product instead of sdk-go-kafka
### Chore
- Remove nolint, add comments
- Update golangci-lint to v2
- Update go version to v1.23
- Updates plugin framework and sdkv2 deps

## 6.7.5
## Refactor
- Use Object Storage Management bundle product instead of sdk-go-object-storage-management
### Fixes
- Save `ionoscloud_container_registry_token` password to `password` field in the resource state

## 6.7.4
### Fixes
- Trying to get Ionoscloud provider version for user agent
### Docs
- Added links in resources
### Refactor
- Use DNS bundle product instead of sdk-go-dns
- Use Container Registry bundle product instead of sdk-go-container-registry
- Use Certificate Manager bundle product instead of sdk-go-cert-manager

## 6.7.3
### Fixes
- Remove `cpu_family`, `availability_zone` and` rockylinux-8-GenericCloud-20230518` from docs
- Do not return an error if `ionoscloud_object_storage_acesskey` is not found
- Return early if dataplatform cluster is in `FAILED` state
- `ionoscloud_s3_key` data source should require only `user_id`
## 6.7.2
### Fixes
- Fix provider crashing when `canonical_user_id` is `nil` in the response for object storage access key
- Fail on `FAILED` state for inmemorydb cluster
## 6.7.1
### Fixes
- Remove cpu_family and availability_zone from the tests
### Features
- Add `IONOS_API_URL_OBJECT_STORAGE_MANAGEMENT` to set a custom API URL for the Object Storage Management Product. Setting `endpoint` or `IONOS_API_URL` does not have any effect
- Add the following privileges to the ionoscloud_group Terraform resource and data source to enhance group access control: accessAndManageLogging, accessAndManageCdn, accessAndManageVpn,
  accessAndManageApiGateway, accessAndManageKaas, accessAndManageNetworkFileStorage, accessAndManageAiModelHub, accessAndManageIamResources, createNetworkSecurityGroups, manageDns
  manageRegistry, manageDataPlatform.

## 6.7.0
### Fixes
- Fix [#735](https://github.com/ionos-cloud/terraform-provider-ionoscloud/issues/735) by reading all values for `api_subnet_allow_list`, not only non-nil values.
- Fix [#748](https://github.com/ionos-cloud/terraform-provider-ionoscloud/issues/748) by removing unecessary error check
- S3 key creation fails with 422 if s3 key not found. Add function to check for that specific response from the API.
### Features
- Add new read-only attribute: `ipv4_cidr_block` to `ionoscloud_lan` resource and data source.
- Make `volume` optional for `ionoscloud_server` resource.
- `name` attribute for `ionoscloud_auto_certificate` resource is now required.
- Add `allow_replace` field to `ionoscloud_pg_cluster` resource.
### Docs
- Changed dead link in MongoDB cluster docs.

## 6.6.9
### Features
- Add support for Monitoring Service: `ionoscloud_monitoring_pipeline` resource and data source.
- Add `expose_serial` attribute for image data source
### Docs
- Replace < and > with " in the docs
- Remove { and } from terraform imports
- Replace \_ with _ in resource names

## 6.6.8
### Features
- Add `auto_scaling` attribute to `ionoscloud_dataplatform_node_pool` resource.
### Fixes
- Omitting the `location` attribute for some resources no longer generates an error
- Better check and log for k8s resources polling

## 6.6.7
### Fix
- Remove location validations for `backup_location` in `ionoscloud_pg_cluster` resource
- Remove location validations for `ionoscloud_mongo_cluster` resource
- Remove location validations for `ionoscloud_datacenter` resource
- Remove location validations for `ionoscloud_vpn_ipsec_gateway` resource
- Remove location validations for `ionoscloud_vpn_wireguard_gateway` resource
- Remove location validations for kafa, auto_certificate, inmemorydb and nfs data sources

## 6.6.6
### Features
- Add `maintenance_window`, `tier` and regional endpoints for VPN resources
### Chores
- Replace `paultyng/ghaction-import-gpg` with `crazy-max/ghaction-import-gpg`

## 6.6.5
### Features
- Resource `ionoscloud_mariadb_cluster` now supports updates
### Testing
- Add import tests for MariaDB clusters

## 6.6.4
### New Product - **Object Storage Management**:
- `Resources`:
  - [ionoscloud_object_storage_accesskey](docs/resources/object_storage_accesskey.md)
- `Data Sources`:
  - [ionoscloud_object_storage_accesskey](docs/data-sources/object_storage_accesskey.md)
  - [ionoscloud_object_storage_region](docs/data-sources/object_storage_region.md)
### Enhancement
- make `mac` optional on `ionoscloud_nic`, `ionoscloud_server`, `ionoscloud_cube_server` and `ionoscloud_vcpu_server`
### Fixes
- Refactor `ionoscloud_share` and `ionoscloud_nic` data sources
- Remove sleep and delete from `ionoscloud_share` resource
### Testing
- Fix template test
- Remove cpu_family from server test
- Fix server and vcpu server tests with multiple firewall rules

## 6.6.3
### Documentation
- Add additional infrastructure provisioning examples
- Fix titles for mariadb docs data sources `https://docs.ionos.com/`
- Add Network Security Group to `https://docs.ionos.com/`
- Add bootvolume_selector to `https://docs.ionos.com/
- Add servers to `https://docs.ionos.com/`
- Add cube server and vcpu server to `https://docs.ionos.com/`
### Enhancement
- Add `allow_replace` to `ionoscloud_server` and `ionoscloud_cube_server` resources, which allows the update of immutable server fields by destroying and then re-creating the resource. This field should be used with care, understanding the risks.
### Fixes
- All `id` and `name` fields in data sources need to be computed, so value can be read on first apply.
### Testing
- Add basic NFS tests

## 6.6.2
### Features
- Make `location` optional for `certificate_manager` resources and datasources
- Make `location` optional for `vpn` resources and datasources
- Make `location` optional for `nfs` resources and datasources
- Make `location` optional for `kafka` resources and datasources
- Add `IONOS_API_URL_NFS` to set a custom API URL for the NAS/NFS product. `location` field needs to be empty, otherwise it will override the custom API URL. Setting `endpoint` or `IONOS_API_URL` does not have any effect.
- Add `IONOS_API_URL_VPN` to set a custom API URL for the VPN product. `location` field needs to be empty, otherwise it will override the custom API URL. Setting `endpoint` or `IONOS_API_URL` does not have any effect.
- Add `IONOS_API_URL_CERT` to set a custom API URL for the Certificate Manager product. `location` field needs to be empty, otherwise it will override the custom API URL. Setting `endpoint` or `IONOS_API_URL` does not have any effect.
- Add `IONOS_API_URL_KAFKA` to set a custom API URL for the Event Streams product. `location` field needs to be empty, otherwise it will override the custom API URL. Setting `endpoint` or `IONOS_API_URL` does not have any effect.
- Add `IONOS_API_URL_MARIADB` to set a custom API URL for the MariaDB product. `location` field needs to be empty, otherwise it will override the custom API URL. Setting `endpoint` or `IONOS_API_URL` does not have any effect.
- Add `IONOS_API_URL_INMEMORYDB` to set a custom API URL for InMemoryDB product. `location` field needs to be empty, otherwise it will override the custom API URL. Setting `endpoint` or `IONOS_API_URL` does not have any effect.
- Add `IONOS_API_URL_OBJECT_STORAGE` to set a custom API URL for Object Storage product. `region` field needs to be empty, otherwise it will override the custom API URL. Setting `endpoint` or `IONOS_API_URL` does not have any effect.
- Add `IONOS_ALLOW_INSECURE` env variable and `insecure` field to allow insecure connections to the API. This is useful for testing purposes only.
- Add import tests for VPN Gateway resources
- Add `security_groups_ids` to `ionoscloud_server`, `ionoscloud_cube_server`, `ionoscloud_nic`, `ionoscloud_vcpu_server` resources and data sources

### New Product - **Network Security Groups**:
- `Resources`:
  - [ionoscloud_nsg](docs/resources/nsg.md)
  - [ionoscloud_nsg_firewallrule](docs/resources/nsg_firewallrule.md)
- `Data Sources`:
  - [ionoscloud_nsg](docs/data-sources/nsg.md)
### Documentation
- Update documentation for `s3_region` and `IONOS_S3_REGION` variables

## 6.6.2
### Fixes
- Fix empty `ssh_key` used as variable in `ssh_keys` field in `ionoscloud_server` resource
- `hostname` needs to be computed as it gets the value of the server name if not set. Fix for `resource_server`, `resource_vcpu_server` and `resource_cube_server`
- Add import tests for VPN Gateway resources

## 6.6.1

### Features
- Add `hostname` to `ionoscloud_server` resource and data source
- Add `hostname` to `ionoscloud_vcpu_server` resource and data source
- Add `hostname` to `ionoscloud_cube_server` resource and data source

## 6.6.0
### Refactor
- Rename `S3` occurrences to `Object Storage`

## 6.5.9
### Features
  - Add new, required `sni_mode` attribute for `ionoscloud_cdn_distribution` resource and data source
### Documentation
  - Add `FAQ` section in `README.md`, add information about IP retrieval for `NIC`s

## 6.5.8
### Refactor
  - Remove `image_alias` sets from `ionoscloud_volume` data source and resource
### Documentation
  - Remove `image_alias` from `ionocloud_volume` data source and resource docs
### Fixes
  - Allow empty `prefix` for bucket lifecycle configuration rules

## 6.5.7
### Fixes
  - Fix documentation rendering of `autoscaling_group` resource and data source, `dbaas_mongo_template` data source and `server_boot_device_selection` resource in Terraform registry
  - Fix `application_loadbalancer_forwardingrule` docs typo
  - Fix for [#687](https://github.com/ionos-cloud/terraform-provider-ionoscloud/issues/687) by setting `user_data` and `backupunit_id` in `ionoscloud_cube_server`

## 6.5.6
### Fixes
- Fix `kafka` remove unavailable locations from resources and data sources
- Fix update behavior for container registry property: `apiSubnetAllowList`
- Fix `ionoscloud_certificate` data source
- Fix `DBaaS` tests, change location for clusters creation, mark `connection_pooler` as computed
- `certificate_id` should not be required for API Gateway resource, `custom_domains` field.
- `cdn distribution` add metadata ipv4, ipv6 and resource_urn to resource and data source
- set 'server_side_encryption' as computed for `ionoscloud_s3_object` resource
### Documentation
- Update documentation for `force_destroy` field in `ionoscloud_s3_bucket` resource

## 6.5.5
### Fixes
- Fix for optional blocks in `ionoscloud_s3_bucket_lifecycle_configuration`
  and `ionoscloud_s3_bucket_website_configuration` resources, before were wrongly marked as required
### Features
- Add `connection_pooler` attribute for PostgreSQL clusters

## 6.5.4
### Fixes
- Fixed bucket public access block documentation
- Fixed resources that need generate MD5 header for the API

## 6.5.3
### Fixes
- `ionoslcoud_logging_pipeline` - `location` should be optional with `de/txl` default. Upgrading should not break existing pipelines.
- Fix DBaaS tests
### Enhancements
- Increase GO version to 1.22, update dependencies

## 6.5.2
### Features
- support for all s3 resources configurations
- Add `location` to `logging_pipeline` resource and data source
### Fixes
- Fix nil deref due to `GeoRestrictions` not being checked against nil

## 6.5.1
### Fixes
- Pass timeouts to `WaitForResourceToBeReady` and `WaitForResourceToBeDeleted` methods
- Add configurable timeouts to s3 buckets. Default stays 60 minutes.
- Fix for vpn wireguard peer when there is no endpoint set
- Minor fixes to documentation for api gateway route resource.
- Remove `public_endpoint` from api gateway route resource.
- Add temporary fix for backup units resources.
- Add tags for certificate manager test files.
- Add computed `id` for `ionoscloud_s3_bucket` resource. Same as name, used for crossplane generation.
- New sdk-go-bundle versions to fix default params not being sent when having default values on marshalling
- Fix CDN tests
- Fix small CDN bug that led to an inconsistent state
- Fix k8s tests.

### Documentation
- Update documentation for S3 bucket resource
- Update documentation for `ionoscloud_inmemorydb_replicaset` resource
- Fix error message for `ionoscloud_s3_bucket_policy` data source when bucket or policy does not exist.
- Fix error message for `ionoscloud_s3_bucket_public_access_block` data source when bucket or public access block does not exist.
- Add validation for `persistence_mode` and `eviction_policy` fields of `ionoscloud_inmemorydb_replicaset`
- Add `ForceNew: true` for some attributes in `ionoscloud_inmemorydb_replicaset` resource
- Fixes #632 update docs with `Principal` example for `s3_bucket_policy`
- Update examples for vpn gateway resources
- Minor fixes to documentation for api gateway and api gateway route resources.
- Only valid hcl in resource examples
- `connections` needs to be required for `ionoscloud_vpn_wireguard_gateway` resource
- Minor documentation fix for CDN resource
- Add basic examples for NFS, VPN Gateway and Kafka resources

## 6.5.0
### Features
- Add support for CDN
- Add new attribute `api_subnet_allow_list` to `container_registry` resource and data source
- Add new attribute `api_subnet_allow_list` to `container_registry` resource and data source
- Add new attribute `protocol_version` to `target_group` resource and data source
- Add new attributes `central_logging` and `logging_format` to `networkloadbalancer` resource and data source
- Add new attributes `central_logging` and `logging_format` to `application_loadbalancer` resource and data source
- Add support for Event Streams for Apache Kafka
- Add support for Certificate Manager providers and auto-certificates
- Add support for In-Memory DB
- Add support for API Gateway
- Add support for VPN Gateway

⚠️ **Note:** Upgrading to 6.5.0 also means using a new version for Certificate Manager service. If, after upgrading to 6.5.0, you receive this error: `{"errorCode": "paas-feature-1", "message": "feature is not enabled for the contract"}`, please send an e-mail to one of the addresses listed here: https://docs.ionos.com/support/general-information/contact-information.


## 6.4.19
### Features
- Add Network File Storage API Support
- Add s3 bucket, object, policy resources with base functionality
### Enhancements
- Move to `sdk-go-bundle` for logging sdk
### Fixes
- Fixes #607. Container registry should wait until the resource is ready before returning the ID.
- Move tests from AMD_OPTERON to INTEL_XEON
- Data source `ionoscloud_mongo_template` should have id `computed` and `optional`
- Fail on k8s cluster and nodepool if creation or deletion entered failed state
- K8s, dataplatform and MariaDB tests
### Documentation
- Update documentation for MariaDB cluster


## 6.4.18
### Features
- Add tests for Mongo cluster and user
- Add new fields for NICs in VM Autoscaling group (firewall_active, firewall_type, firewall_rule, flow_log, target_group)
- Refactor VM Autoscaling group
### Fixes
- Wrap missing base error for resource fetching errors
- Properly persist user group ids in state when syncing with remote configuration
- Quick fix for MariaDB State metadata values
### Enhancements
- Add `grafana_address` attribute to `ionoscloud_logging_pipeline` resource and data source
### Misc
- Replace deprecated `--rm-dist` with `--clean` in release workflow
### Documentation
- Updated documentation to specify that `ionoscloud_logging_pipeline`, `ionoscloud_dns_record` and `ionoscloud_dns_zone` only accept tokens for authorization.
- Removed Early Access (EA) warning for `ionoscloud_logging_pipeline`.


## 6.4.17
### Fixes
- Correctly raise immutable error for changes to `template_uuid` when running `terraform plan` for Cube servers

### Documentation
- Update `ionoscloud_user` documentation. Fix `administrator` and add other fields description
- Change to have nested lists show correctly in tf registry docs
- Fix documentation for `ionoscloud_server`, `ionoscloud_volume`, `ionoscloud_lan` resources and `ionoscloud_image` data sources

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
    - ionoscloud_dataplatform_cluster
    - ionoscloud_dataplatform_node_pool
  - `Data Sources`:
    - ionoscloud_dataplatform_cluster
    - ionoscloud_dataplatform_node_pool
    - ionoscloud_dataplatform_node_pools
    - ionoscloud_dataplatform_versions

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
 - Allow server import with nic and firewallId : `terraform import ionoscloud_server.myserver datacenter uuid/server uuid/primary nic id/firewall rule id`
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
    - [ionoscloud_certificate](docs/resources/certificate_manager_certificate.md)
  - `Data Sources`:
    - [ionoscloud_certificate](docs/data-sources/certificate_manager_certificate.md)


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
