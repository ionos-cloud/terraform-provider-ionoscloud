package constant

import (
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// IonosDebug - env variable, set to true to enable debug
const IonosDebug = "IONOS_DEBUG"

// MaxRetries - number of retries in case of rate-limit
const MaxRetries = 999

// MaxWaitTime - waits 4 seconds before retry in case of rate limit
const MaxWaitTime = 4 * time.Second

const SleepInterval = 5 * time.Second

const Available = "AVAILABLE"

// Datetime formats
const (
	DatetimeZLayout        = "2006-01-02 15:04:05Z"
	DatetimeTZOffsetLayout = "2006-01-02 15:04:05 -0700 MST"
)

// VPN Constants

// DefaultTier is the default tier value for Wireguard & IPSec gateways
const DefaultTier = "STANDARD"

// ApplicationLoadBalancer Constants
const (
	ALBResource         = "ionoscloud_application_loadbalancer"
	ALBTestResource     = "test_application_loadbalancer"
	ALBDataSourceById   = "test_application_loadbalancer_id"
	ALBDataSourceByName = "test_application_loadbalancer_name"

	ALBForwardingRuleResource         = "ionoscloud_application_loadbalancer_forwardingrule"
	ALBForwardingRuleTestResource     = "test_application_loadbalancer_forwardingrule"
	ALBForwardingRuleDataSourceById   = "test_application_loadbalancer_forwardingrule_id"
	ALBForwardingRuleDataSourceByName = "test_application_loadbalancer_forwardingrule_name"

	TargetGroupResource         = "ionoscloud_target_group"
	TargetGroupTestResource     = "test_target_group"
	TargetGroupDataSourceById   = "test_target_group_id"
	TargetGroupDataSourceByName = "test_target_group_name"
)

// Image Constants
const (
	ImageResource     = "ionoscloud_image"
	ImageTestResource = "test_image"
)

// Location Constants
const (
	LocationResource     = "ionoscloud_location"
	LocationTestResource = "test_location"
)

// Resource Constants
const (
	ResourceResource     = "ionoscloud_resource"
	ResourceTestResource = "test_resource"
)

// Template Constants
const (
	TemplateResource     = "ionoscloud_template"
	TemplateTestResource = "test_template"
)

// BackupUnit Constants
const (
	BackupUnitResource         = "ionoscloud_backup_unit"
	BackupUnitTestResource     = "testBackupUnit"
	BackupUnitDataSourceById   = "testBackupUnitId"
	BackupUnitDataSourceByName = "testBackupUnitName"
)

// Datacenter Constants
const (
	DatacenterResource           = "ionoscloud_datacenter"
	DatacenterTestResource       = "test_datacenter"
	DatacenterDataSourceById     = "test_datacenter_id"
	DatacenterDataSourceByName   = "test_datacenter_name"
	DatacenterDataSourceMatching = "test_datacenter_matching"
)

// Firewall Constants
const (
	FirewallResource         = "ionoscloud_firewall"
	FirewallTestResource     = "test_firewall"
	FirewallDataSourceById   = "test_firewall_id"
	FirewallDataSourceByName = "test_firewall_name"
)

// FirewallProtocolEnum values for firewall rule protocol
var FirewallProtocolEnum = []string{"TCP", "UDP", "ICMP", "ICMPv6", "GRE", "VRRP", "ESP", "AH", "ANY"}

// NSG Constants
const (
	NSGResource                 = "ionoscloud_nsg"
	NSGTestResource             = "test_nsg"
	NGDataSourceByID            = "test_nsg_id"
	NSGDataSourceByName         = "test_nsg_name"
	NSGFirewallRuleResource     = "ionoscloud_nsg_firewallrule"
	NSGFirewallRuleTestResource = "test_nsg_firewallrule"
	NSGSelectionResource        = "ionoscloud_datacenter_nsg_selection"
)

// Lan Constants
const (
	LanResource                = "ionoscloud_lan"
	LanTestResource            = "test_lan"
	PrivateLANTestResource     = "private_lan"
	LanDataSourceByID          = "test_lan_id"
	LanDataSourceByName        = "test_lan_name"
	PrivateLANDataSourceByName = "private_lan_name"
)

// Group Constants
const (
	GroupResource         = "ionoscloud_group"
	GroupTestResource     = "test_group"
	GroupDataSourceById   = "test_group_id"
	GroupDataSourceByName = "test_group_name"
)

// K8s Constants
const (
	K8sClusterResource                = "ionoscloud_k8s_cluster"
	K8sClustersDataSource             = "ionoscloud_k8s_clusters"
	K8sClusterTestResource            = "test_k8s_cluster"
	PrivateK8sClusterTestResource     = "test_private_k8s_cluster"
	K8sClusterDataSourceByID          = "test_k8s_cluster_id"
	K8sClusterDataSourceByName        = "test_k8s_cluster_name"
	K8sClustersDataSourceFilterName   = "test_k8s_clusters_filter_name"
	K8sClustersDataSourceFilterPublic = "test_k8s_clusters_filter_public"

	K8sNodePoolResource         = "ionoscloud_k8s_node_pool"
	K8sNodePoolNodesResource    = "ionoscloud_k8s_node_pool_nodes"
	K8sNodePoolTestResource     = "test_k8s_node_pool"
	K8sNodePoolDataSourceById   = "test_k8s_node_pool_id"
	K8sNodePoolDataSourceByName = "test_k8s_node_pool_name"

	ResourceNameK8sNodePool      = K8sNodePoolResource + "." + K8sNodePoolTestResource
	DataSourceK8sNodePoolId      = DataSource + "." + K8sNodePoolResource + "." + K8sNodePoolDataSourceById
	DataSourceK8sNodePoolNodesId = DataSource + "." + K8sNodePoolNodesResource + "." + K8sNodePoolDataSourceById
	DataSourceK8sNodePoolName    = DataSource + "." + K8sNodePoolResource + "." + K8sNodePoolDataSourceByName
	K8sNodePoolTimeout           = 3 * time.Hour
)

var ResourceK8sNodePoolTimeout = schema.ResourceTimeout{
	Create:  schema.DefaultTimeout(K8sNodePoolTimeout),
	Update:  schema.DefaultTimeout(K8sNodePoolTimeout),
	Delete:  schema.DefaultTimeout(K8sNodePoolTimeout),
	Default: schema.DefaultTimeout(K8sNodePoolTimeout),
}

// LoadBalancer Constants
const LoadBalancerResource = "ionoscloud_loadbalancer"

// NatGateway Constants
const (
	NatGatewayResource             = "ionoscloud_natgateway"
	NatGatewayTestResource         = "test_nat_gateway"
	NatGatewayDataSourceById       = "test_nat_gateway_id"
	NatGatewayDataSourceByName     = "test_nat_gateway_name"
	NatGatewayRuleResource         = "ionoscloud_natgateway_rule"
	NatGatewayRuleTestResource     = "test_nat_gateway"
	NatGatewayRuleDataSourceById   = "test_nat_gateway_id"
	NatGatewayRuleDataSourceByName = "test_nat_gateway_name"
)

// NetworkLoadBalancer Constants
const (
	NetworkLoadBalancerResource                       = "ionoscloud_networkloadbalancer"
	NetworkLoadBalancerTestResource                   = "test_networkloadbalancer"
	NetworkLoadBalancerDataSourceById                 = "test_networkloadbalancer_id"
	NetworkLoadBalancerDataSourceByName               = "test_networkloadbalancer_name"
	NetworkLoadBalancerForwardingRuleResource         = "ionoscloud_networkloadbalancer_forwardingrule"
	NetworkLoadBalancerForwardingRuleTestResource     = "test_networkloadbalancer_forwardingrule"
	NetworkLoadBalancerForwardingRuleDataSourceById   = "test_networkloadbalancer_forwardingrule_id"
	NetworkLoadBalancerForwardingRuleDataSourceByName = "test_networkloadbalancer_forwardingrule_name"
)

// NFS Constants
const (
	NFSClusterResource = "ionoscloud_nfs_cluster"
	NFSShareResource   = "ionoscloud_nfs_share"
)

// Private Cross Connect Constants
// The resource name was changed from Private Cross Connect to Cross Connect
// But the terraform resources names did not change for backwards compatibility reasons
const (
	PCCResource         = "ionoscloud_private_crossconnect"
	PCCTestResource     = "test_private_crossconnect"
	PCCDataSourceById   = "test_private_crossconnect_id"
	PCCDataSourceByName = "test_private_crossconnect_name"
)

// Server Constants
const (
	ServerResource         = "ionoscloud_server"
	ServerCubeResource     = "ionoscloud_cube_server"
	ServerVCPUResource     = "ionoscloud_vcpu_server"
	ServerTestResource     = "test_server"
	ServerTestHostname     = "myhostname"
	ServerDataSourceById   = "test_server_id"
	ServerDataSourceByName = "test_server_name"
)

// S3Key Constants
const (
	S3KeyResource       = "ionoscloud_s3_key"
	S3KeyTestResource   = "test_s3_key"
	S3KeyDataSourceById = "test_s3_key_id"
)

// Snapshot Constants
const (
	SnapshotResource         = "ionoscloud_snapshot"
	SnapshotTestResource     = "test_snapshot"
	SnapshotDataSourceById   = "test_snapshot_id"
	SnapshotDataSourceByName = "test_snapshot_name"
)

// User Constants
const (
	UserResource = "ionoscloud_user"
	// Used for tests where we need fresh user creation, e.g the tests in which we create the user
	// and also add it to a group in the same time.
	NewUserName          = "new_test_user"
	NewUserResource      = "new_test_user_resource"
	UserTestResource     = "test_user"
	UserDataSourceById   = "test_user_id"
	UserDataSourceByName = "test_user_name"
)

// Ip Block constants
const (
	IpBlockResource           = "ionoscloud_ipblock"
	IpBlockTestResource       = "test_ip_block"
	IpBlockDataSourceById     = "test_ip_block_id"
	IpBlockDataSourceByName   = "test_ip_block_id"
	IpBlockDataSourceMatching = "test_ip_block_id_location"
)

const (
	NicResource         = "ionoscloud_nic"
	FullNicResourceName = NicResource + "." + NicTestResourceName
	NicTestResourceName = "database_nic"
	NicMac              = "00:0a:95:9d:68:16"
)

const (
	VolumeResource         = "ionoscloud_volume"
	VolumeTestResource     = "test_volume"
	VolumeDataSourceById   = "test_volume_id"
	VolumeDataSourceByName = "test_volume_name"
)

// DBaaS Constants
const (
	DBaaSClusterTestDataSourceById   = "test_dbaas_cluster_id"
	DBaaSClusterTestResource         = "test_dbaas_cluster"
	DBaaSClusterTestDataSourceByName = "test_dbaas_cluster_name"
	DBaaSReplicaSetTestResource      = "test_dbaas_replicaset"

	// PgSql constants
	PsqlClusterResource          = "ionoscloud_pg_cluster"
	PsqlDatabaseResource         = "ionoscloud_pg_database"
	PsqlDatabasesResource        = "ionoscloud_pg_databases"
	PsqlDatabaseTestResource     = "test_database"
	PsqlDatabaseDataSourceByName = "test_database_name"
	PsqlDatabasesDataSource      = "test_databases"
	PsqlUserResource             = "ionoscloud_pg_user"
	PsqlBackupsResource          = "ionoscloud_pg_backups"
	PsqlBackupsTest              = "test_dbaas_backups"
	PsqlVersionsResource         = "ionoscloud_pg_versions"
	PsqlVersionsTest             = "test_dbaas_versions"

	// MariaDB constants
	DBaaSMariaDBClusterResource       = "ionoscloud_mariadb_cluster"
	DBaaSMariaDBBackupsDataSource     = "ionoscloud_mariadb_backups"
	DBaasMariaDBBackupsDataSourceName = "test_mariadb_backups"

	DBaasMongoClusterResource        = "ionoscloud_mongo_cluster"
	DBaasMongoUserResource           = "ionoscloud_mongo_user"
	DBaaSMongoTemplateResource       = "ionoscloud_mongo_template"
	DBaaSMongoTemplateTestDataSource = "test_dbaas_mongo_template"

	// InMemoryDB constants
	DBaaSReplicaSetTestDataSourceByID   = "test_dbaas_replicaset_id"
	DBaaSReplicaSetTestDataSourceByName = "test_dbaas_replicaset_name"
	DBaaSInMemoryDBReplicaSetResource   = "ionoscloud_inmemorydb_replicaset"
	DBaaSInMemoryDBSnapshotResource     = "ionoscloud_inmemorydb_snapshot"
)

// Locations slice represents the locations in which services are available.
var Locations = []string{"de/fra", "de/txl", "es/vit", "fr/par", "gb/lhr", "us/ewr", "us/las", "us/mci"}

// Container Registry Constants
const (
	// ContainerRegistryTestResource needs to be with -, do not change
	ContainerRegistryTestResource      = "test-terraform-container-registry"
	ContainerRegistryTokenTestResource = "test-terraform-container-registry-token"

	ContainerRegistryResource                  = "ionoscloud_container_registry"
	ContainerRegistryTestDataSourceById        = "test_container_registry_id"
	ContainerRegistryTestDataSourceByName      = "test_container_registry_name"
	ContainerRegistryTokenResource             = "ionoscloud_container_registry_token"
	ContainerRegistryTokenTestDataSourceById   = "test_container_registry_token_id"
	ContainerRegistryTokenTestDataSourceByName = "test_container_registry_token_name"
	ContainerRegistryLocationsResource         = "ionoscloud_container_registry_locations"
	ContainerRegistryLocationsTest             = "test_container_registry_locations"
)

// Dataplatform Constants
const (
	DataplatformClusterResource              = "ionoscloud_dataplatform_cluster"
	DataplatformClusterTestResource          = "test_dataplatform_cluster"
	DataplatformClusterTestDataSourceById    = "test_dataplatform_cluster_id"
	DataplatformClusterTestDataSourceByName  = "test_dataplatform_cluster_name"
	DataplatformNodePoolResource             = "ionoscloud_dataplatform_node_pool"
	DataplatformNodePoolTestResource         = "test_dataplatform_node_pool"
	DataplatformNodePoolTestDataSourceById   = "test_dataplatform_node_pool_id"
	DataplatformNodePoolTestDataSourceByName = "test_dataplatform_node_pool_name"
	DataplatformNodePoolsDataSource          = "ionoscloud_dataplatform_node_pools"
	DataplatformVersionsDataSource           = "ionoscloud_dataplatform_versions"
	DataplatformNodePoolsTestDataSource      = "test_dataplatform_node_pools"
	DataplatformVersionsTestDataSource       = "test_dataplatform_versions"
	// DataPlatformVersion lowest 'available' version is now 24.3
	DataPlatformVersion              = "24.11"
	DataPlatformNameRegexConstraint  = "^[A-Za-z0-9][-A-Za-z0-9_.]*[A-Za-z0-9]$"
	DataPlatformRegexNameError       = "name should match " + DataPlatformNameRegexConstraint
	ResourceNameDataplatformNodePool = DataplatformNodePoolResource + "." + DataplatformNodePoolTestResource
)

// DNS Constants
const (
	DNSRecordDataSource         = "ionoscloud_dns_record"
	DNSRecordResource           = "ionoscloud_dns_record"
	DNSZoneDataSource           = "ionoscloud_dns_zone"
	DNSZoneResource             = "ionoscloud_dns_zone"
	DNSZoneTestResourceName     = "test_dns_zone"
	DNSZoneTestDataSourceName   = "test_dns_zone_data_source"
	DNSRecordTestResourceName   = "test_dns_record"
	DNSRecordTestDataSourceName = "test_dns_record_data_source"
)
const (
	ShareResource         = "ionoscloud_share"
	ShareResourceFullName = ShareResource + "." + SourceShareName
	SourceShareName       = "share"
)

// CDN Constants
const (
	CDNDistributionResource           = "ionoscloud_cdn_distribution"
	CDNDistributionTestResource       = "test_cdn_distribution"
	CDNDistributionDataSourceByID     = "test_cdn_distribution_id"
	CDNDistributionDataSourceByDomain = "test_cdn_distribution_domain"
	CDNDistributionDataSourceMatching = "test_cdn_distribution_matching"
)

const (
	ResourceIpFailover         = "ionoscloud_ipfailover"
	IpfailoverResourceFullName = ResourceIpFailover + "." + IpfailoverName
	IpfailoverName             = "failover-group"
	SecondIpfailoverName       = "second-failover-group"
)

// General Constants
const (
	DataSource        = "data"
	UpdatedResources  = "test_updated"
	DataSourcePartial = "test"
	RandomPassword    = "random_password"
)

// Logging Service Constants

const (
	LoggingPipelineDataSource         = "ionoscloud_logging_pipeline"
	LoggingPipelineResource           = "ionoscloud_logging_pipeline"
	LoggingPipelineTestResourceName   = "test_logging_pipeline"
	LoggingPipelineTestDataSourceName = "test_logging_pipeline_data_source"
)

const ServersDataSource = "ionoscloud_servers"

// Certificate constants
const (
	TestCertificate = `-----BEGIN CERTIFICATE-----
MIIDazCCAlOgAwIBAgIUOH1cikhurIjCjm5Zxt7sfJmhIVAwDQYJKoZIhvcNAQEL
BQAwRTELMAkGA1UEBhMCQVUxEzARBgNVBAgMClNvbWUtU3RhdGUxITAfBgNVBAoM
GEludGVybmV0IFdpZGdpdHMgUHR5IEx0ZDAeFw0yMjA2MDkxMjM0MzVaFw0zMjA2
MDYxMjM0MzVaMEUxCzAJBgNVBAYTAkFVMRMwEQYDVQQIDApTb21lLVN0YXRlMSEw
HwYDVQQKDBhJbnRlcm5ldCBXaWRnaXRzIFB0eSBMdGQwggEiMA0GCSqGSIb3DQEB
AQUAA4IBDwAwggEKAoIBAQDkVU596LWGR+/nC3r3MndJfimGMDvt9W4SwL0bOa+V
XxgVgUKYTTCPvwZaAQrtJRUjW2bGxwwj8/3uDEY6vwHJ1Yh+OrbrQHPFPKcBbRie
8mqwgjnAveqvlRKxi3VWwG0Bevki54ghwolmZ5GppvzeqLNYFF8nYuSAbseRoPFb
EJMLd5vuEkDytZl42eiZkv/aHEtUGXvcTY29K6G4yGOEr3Pr320ts8tVW4UNlBt4
0mDfBjtXAeSIcQfww/c69Pc3Xrfd3FVf4Qjo3bhMCvbg5shvRHmJrcbOPJO5kUn+
mwPU7DlJM9YeOMQBMgmw3NoKKI4dOU3HUBpBiN3M5tztAgMBAAGjUzBRMB0GA1Ud
DgQWBBRCecVMYml89VvfhBl+DTxzqcwWoTAfBgNVHSMEGDAWgBRCecVMYml89Vvf
hBl+DTxzqcwWoTAPBgNVHRMBAf8EBTADAQH/MA0GCSqGSIb3DQEBCwUAA4IBAQDJ
oF3c5MGbrg5Iu8TF9X8g3tgVANGUVWDHDxx0Fe3zQojW0NSfHPtQ+Qkf6BdYH6oc
OzQBUgWFGnrhPliUW9rD4/8c3BoVvT9ukYPOhwLDd2lPqTTbbfhdkzDSM/BKPP1g
7Ok2m/uk9jnsLQSCQE4zc8+X0M+zG9ZPyC0MJqM3d7gB+LVOE8PKIJz6fXCyoakz
18PV+e4RhL5daTFCdZ1XAL146kIorS4XX5iIyvCt1WBzSS8IUtAIgR/QLxk7ZqrL
BKEkcU1X0yvgyDUkpcJ1BS/++5q/EDEQCYP6gN0cvPFhFvQeNor5SId6EFFlEkMn
MYuea4TP5Gk2UkmDOxuJ
-----END CERTIFICATE-----`
	PrivateKey = `-----BEGIN PRIVATE KEY-----
MIIEvgIBADANBgkqhkiG9w0BAQEFAASCBKgwggSkAgEAAoIBAQDkVU596LWGR+/n
C3r3MndJfimGMDvt9W4SwL0bOa+VXxgVgUKYTTCPvwZaAQrtJRUjW2bGxwwj8/3u
DEY6vwHJ1Yh+OrbrQHPFPKcBbRie8mqwgjnAveqvlRKxi3VWwG0Bevki54ghwolm
Z5GppvzeqLNYFF8nYuSAbseRoPFbEJMLd5vuEkDytZl42eiZkv/aHEtUGXvcTY29
K6G4yGOEr3Pr320ts8tVW4UNlBt40mDfBjtXAeSIcQfww/c69Pc3Xrfd3FVf4Qjo
3bhMCvbg5shvRHmJrcbOPJO5kUn+mwPU7DlJM9YeOMQBMgmw3NoKKI4dOU3HUBpB
iN3M5tztAgMBAAECggEBANJZQFlAA9Kz/O+VpO+L/1amMmzbjKo4evItu0kUiIwM
MezFyurx2XXjnl9WLJGxotqSvokLIEUS5vDhP+Wox2YAIKFhR9hL5RtkN9pZfeAY
JW98WOYWT9j3dWQ4vJ1x4joF5vRf5gpr5BaB/TAUlUoukiHnio2HTkh/Rb0ETrT7
Pvl9hYFO50xmaxwd5Vy+726ZLwOkkraDpXB1jZC9Kp7EfnMi7ekZ8LfBYmEdl87Y
VvBghjSsRL4VdY/WTOpWM1DnOIBrUmM/0UfYW1uaV4upSPScOjFeBllY+lSpyO38
B+L1eQSJghIULOntN5XUGnrTpMSXW8C67qaEFPfa9qECgYEA/UyAQCEXuFiEDpia
CkZ0Ykh6xxY2sA6jMB52RvjpWxqbrVUE6yGMM2UJxNplwdZk9lmpzU9KPfDgulKX
Uq34O94wDSXKJQYI0GsrXs3IgheXKVT/4s3t9oCc2hH0F3/o37jYOkYP571e2sdd
yQd0aTZqG1qp7bZlRKWahKrB/FUCgYEA5sSmc7dIwxgX/kp/4VUdjgFUy/GU5xr9
6xnioGsdnb6rBpicklri/h2E4eLLzgVbuzVQMLIAG/MpwrIxWspUR64yAaPEAVm8
3GIg771JZHl06lYYjAqaSy1qC8v3/3T+masWwa/MNCxXB6YN9ptohAf4M3hlEL/J
jlR8Qp5M9jkCgYAqcPgIRsM3szUlUPJ2iEmV8jkIRLOTGlDDjkcZKznGdxXgnB8/
2pYoQmS5pDJqoSa3lFx8Ny3kZQjyj0Ylp1qxhVAd09gkDffKHDrfHrHbAmLknQZn
FUQrCm+9pkZ07Yyyd8FbOkQN+0/6bm9LcMFTo7dxr+ZLG0Wqk+jpE8d/JQKBgCFf
s0rs6OL+KwolsBTggGO3IZJVH9nEd5B2r+XPV/smRgmwLISmDEn/7uXULPFgqQGM
FkrUk1t3cUStDKI6vLGZKbY+/uvLFJsyvdyuHV0gi54QUYB/UA0rRjjqiLUzMFb8
/U+JoxiwiO2cQEy38QeXN3gKI2OmuPmSkl34Et1RAoGBAIqPAdXfyoMb6stylvC1
N20fcwpG3aiTESteYpnXCNFW8XrMnoBWL6bK6st4eBSUbvOfjTJrSVC/KBLR6awV
i+U582LTWq8y6WA8tdqfeZO+TUl+8DBk6k6aDbA8a3+X/D+sTsRfSavEVyEeV7EO
wkv+4ThHJ677Dpi/P8F8iOJp
-----END PRIVATE KEY-----`
	CertificateResource             = "ionoscloud_certificate" //nolint:golint
	AutoCertificateProviderResource = "ionoscloud_auto_certificate_provider"
	AutoCertificateResource         = "ionoscloud_auto_certificate"
	TestCertName                    = "certTest"
	TestCMProviderName              = "CMProviderTest"
	TestCMAutoCertificateName       = "CMAutoCertificateTest"
)

// Server type constants
const (
	VCPUType       = "VCPU"
	CubeType       = "CUBE"
	EnterpriseType = "ENTERPRISE"
)

// Server power state constants
const (
	CubeVMStateStop = "SUSPENDED"
	VMStateStart    = "RUNNING"
	VMStateStop     = "SHUTOFF"
)

// Server boot devices constants
const (
	BootDeviceTypeCDROM  = "cdrom"
	BootDeviceTypeVolume = "volume"
)

// VolumeBootOrderNone unsets the volume as boot device for the VM it is attached to
const VolumeBootOrderNone = "NONE"

const (
	// FlowlogBucket created on account that runs CI
	FlowlogBucket        = "flowlog-acceptance-test"
	FlowlogBucketUpdated = "flowlog-acceptance-test-updated"
)

const RepoURL = "https://github.com/ionos-cloud/terraform-provider-ionoscloud"

const (
	AutoscalingGroupResource              = "ionoscloud_autoscaling_group"
	AutoscalingGroupTestResource          = "test_autoscaling_group"
	AutoscalingGroupDataSourceById        = "test_autoscaling_group_id"
	AutoscalingGroupDataSourceByName      = "test_autoscaling_group_name"
	AutoscalingGroupServersResource       = "ionoscloud_autoscaling_group_servers"
	AutoscalingGroupServersTestDataSource = "test_autoscaling_servers"
)

const (
	ServerBootDeviceSelectionResource     = "ionoscloud_server_boot_device_selection"
	TestServerBootDeviceSelectionResource = "boot_device_selection_example"
)

const (
	// WireGuardGatewayResource is the full name of the WireGuardGatewayResource
	WireGuardGatewayResource = "ionoscloud_vpn_wireguard_gateway"
	// WireGuardPeerResource is the full name of the WireGuardGatewayResource
	WireGuardPeerResource = "ionoscloud_vpn_wireguard_peer"
	// WireGuardGatewayTestResource is the name used for testing
	WireGuardGatewayTestResource = "test_wireguard_gateway"
	// WireGuardPeerTestResource is the name used for testing
	WireGuardPeerTestResource = "test_wireguard_peer"
	// IPSecGatewayResource is the full name of the IPSecGatewayResource
	IPSecGatewayResource = "ionoscloud_vpn_ipsec_gateway"
	// IPSecTunnelResource is the full name of the IPSecTunnelResource
	IPSecTunnelResource = "ionoscloud_vpn_ipsec_tunnel"
)

var ForwardingRuleAlgorithms = []string{"ROUND_ROBIN", "LEAST_CONNECTION", "RANDOM", "SOURCE_IP"}
var LBTargetProxyProtocolVersions = []string{"none", "v1", "v2", "v2ssl"}

// Maximum limit for various resources
// The limit represents the maximum number of entities that can be fetched in a single 'GET' request
const (
	TargetGroupLimit = 200
	IPBlockLimit     = 1000
)

// Kafka Constants
const (
	KafkaClusterResource      = "ionoscloud_kafka_cluster"
	KafkaClusterTopicResource = "ionoscloud_kafka_cluster_topic"
)

// API Gateway Constants
const (
	APIGatewayResource      = "ionoscloud_apigateway"
	APIGatewayRouteResource = "ionoscloud_apigateway_route"
)
