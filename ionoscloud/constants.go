package ionoscloud

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

// Lan Constants
const (
	LanResource         = "ionoscloud_lan"
	LanTestResource     = "test_lan"
	LanDataSourceById   = "test_lan_id"
	LanDataSourceByName = "test_lan_name"
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
	K8sClusterResource         = "ionoscloud_k8s_cluster"
	K8sClusterTestResource     = "test_k8s_cluster"
	K8sClusterDataSourceById   = "test_k8s_cluster_id"
	K8sClusterDataSourceByName = "test_k8s_cluster_name"

	K8sNodePoolResource         = "ionoscloud_k8s_node_pool"
	K8sNodePoolNodesResource    = "ionoscloud_k8s_node_pool_nodes"
	K8sNodePoolTestResource     = "test_k8s_node_pool"
	K8sNodePoolDataSourceById   = "test_k8s_node_pool_id"
	K8sNodePoolDataSourceByName = "test_k8s_node_pool_name"

	ResourceNameK8sNodePool   = K8sNodePoolResource + "." + K8sNodePoolTestResource
	DataSourceK8sNodePoolId   = DataSource + "." + K8sNodePoolResource + "." + K8sNodePoolDataSourceById
	DataSourceK8sNodePoolName = DataSource + "." + K8sNodePoolResource + "." + K8sNodePoolDataSourceByName
)

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

// Private Crossconnect Constants
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
	ServerTestResource     = "test_server"
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
	fullNicResourceName = NicResource + "." + nicTestResourceName
	nicTestResourceName = "database_nic"
)

const (
	VolumeResource         = "ionoscloud_volume"
	VolumeTestResource     = "test_volume"
	VolumeDataSourceById   = "test_volume_id"
	VolumeDataSourceByName = "test_volume_name"
)

// DBaaS Constants
const (
	DBaaSClusterTestResource         = "test_dbaas_cluster"
	DBaaSClusterTestDataSourceById   = "test_dbaas_cluster_id"
	DBaaSClusterTestDataSourceByName = "test_dbaas_cluster_name"

	PsqlClusterResource  = "ionoscloud_pg_cluster"
	PsqlBackupsResource  = "ionoscloud_pg_backups"
	PsqlBackupsTest      = "test_dbaas_backups"
	PsqlVersionsResource = "ionoscloud_pg_versions"
	PsqlVersionsTest     = "test_dbaas_versions"

	DBaasMongoClusterResource = "ionoscloud_mongo_cluster"
	DBaasMongoUserResource    = "ionoscloud_mongo_user"
)

// Container Registry Constants
const (
	//ContainerRegistryTestResource needs to be with -, do not change
	ContainerRegistryTestResource      = "test-container-registry"
	ContainerRegistryTokenTestResource = "test-container-registry-token"

	ContainerRegistryResource                  = "ionoscloud_container_registry"
	ContainerRegistryTestDataSourceById        = "test_container_registry_id"
	ContainerRegistryTestDataSourceByName      = "test_container_registry_name"
	ContainerRegistryTokenResource             = "ionoscloud_container_registry_token"
	ContainerRegistryTokenTestDataSourceById   = "test_container_registry_token_id"
	ContainerRegistryTokenTestDataSourceByName = "test_container_registry_token_name"
	ContainerRegistryLocationsResource         = "ionoscloud_container_registry_locations"
	ContainerRegistryLocationsTest             = "test_container_registry_locations"
)

const (
	ShareResource         = "ionoscloud_share"
	shareResourceFullName = ShareResource + "." + sourceShareName
	sourceShareName       = "share"
)

const (
	ResourceIpFailover         = "ionoscloud_ipfailover"
	ipfailoverResourceFullName = ResourceIpFailover + "." + ipfailoverName
	ipfailoverName             = "failover-test"
)

// General Constants
const (
	DataSource        = "data"
	UpdatedResources  = "test_updated"
	DataSourcePartial = "test"
)

const ServersDataSource = "ionoscloud_servers"

const (
	CertificateResource = "ionoscloud_certificate"
	TestCertName        = "certTest"
)

type clientType int

const (
	ionosClient clientType = iota
	psqlClient
	certManagerClient
	mongoClient
	containerRegistryClient
)
