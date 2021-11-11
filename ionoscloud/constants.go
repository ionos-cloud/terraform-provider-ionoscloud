package ionoscloud

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
	UserResource         = "ionoscloud_user"
	UserTestResource     = "test_user"
	UserDataSourceById   = "test_user_id"
	UserDataSourceByName = "test_user_name"
)

//Ip Block constants
const IpBLockResource = "ionoscloud_ipblock"

const (
	nicResource         = "ionoscloud_nic"
	fullNicResourceName = nicResource + "." + nicTestResourceName
	nicTestResourceName = "database_nic"
)

const (
	VolumeResource         = "ionoscloud_volume"
	VolumeTestResource     = "test_volume"
	VolumeDataSourceById   = "test_volume_id"
	VolumeDataSourceByName = "test_volume_name"
)

const (
	shareResource         = "ionoscloud_share"
	shareResourceFullName = shareResource + "." + sourceShareName
	sourceShareName       = "share"
)

// General Constants
const (
	DataSource       = "data"
	UpdatedResources = "test_updated"
)
