package ionoscloud

// BackupUnit Constants
const (
	BackupUnitResource                = "ionoscloud_backup_unit"
	BackupUnitTestResource            = "testBackupUnit"
	BackupUnitDataSourceById          = "testBackupUnitId"
	BackupUnitDataSourceByName        = "testBackupUnitName"
	testAccCheckBackupUnitConfigBasic = `
resource ` + BackupUnitResource + ` ` + BackupUnitTestResource + ` {
	name        = "` + BackupUnitTestResource + `"
	password    = "DemoPassword123$"
	email       = "example@ionoscloud.com"
}
`
)

// Datacenter Constants
const (
	DatacenterResource                = "ionoscloud_datacenter"
	DatacenterTestResource            = "test_datacenter"
	DatacenterDataSourceById          = "test_datacenter_id"
	DatacenterDataSourceByName        = "test_datacenter_name"
	DatacenterDataSourceMatching      = "test_datacenter_matching"
	testAccCheckDatacenterConfigBasic = `
resource ` + DatacenterResource + ` ` + DatacenterTestResource + ` {
	name       = "` + DatacenterTestResource + `"
	location = "us/las"
	description = "Test Datacenter Description"
	sec_auth_protection = false
}`
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
	LanResource                = "ionoscloud_lan"
	LanTestResource            = "test_lan"
	LanDataSourceById          = "test_lan_id"
	LanDataSourceByName        = "test_lan_name"
	testAccCheckLanConfigBasic = testAccCheckDatacenterConfigBasic + `
resource ` + LanResource + ` ` + LanTestResource + ` {
  datacenter_id = ` + DatacenterResource + `.` + DatacenterTestResource + `.id
  public = true
  name = "` + LanTestResource + `"
}`
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
	PCCResource                                = "ionoscloud_private_crossconnect"
	PCCTestResource                            = "test_private_crossconnect"
	PCCDataSourceById                          = "test_private_crossconnect_id"
	PCCDataSourceByName                        = "test_private_crossconnect_name"
	testAccCheckPrivateCrossConnectConfigBasic = `
resource ` + PCCResource + ` ` + PCCTestResource + ` {
  name        = "` + PCCTestResource + `"
  description = "` + PCCTestResource + `"
}`
)

// Server Constants
const (
	ServerResource                = "ionoscloud_server"
	ServerTestResource            = "test_server"
	ServerDataSourceById          = "test_server_id"
	ServerDataSourceByName        = "test_server_name"
	testAccCheckServerConfigBasic = `
resource ` + DatacenterResource + ` ` + DatacenterTestResource + ` {
	name       = "server-test"
	location = "us/las"
}
` + testAccCheckBackupUnitConfigBasic + `

resource "ionoscloud_ipblock" "webserver_ipblock" {
  location = ` + DatacenterResource + `.` + DatacenterTestResource + `.location
  size = 4
  name = "webserver_ipblock"
}
resource ` + LanResource + ` ` + LanTestResource + ` {
  datacenter_id = ` + DatacenterResource + `.` + DatacenterTestResource + `.id
  public = true
  name = "public"
}
resource ` + ServerResource + ` ` + ServerTestResource + ` {
  name = "` + ServerTestResource + `"
  datacenter_id = ` + DatacenterResource + `.` + DatacenterTestResource + `.id
  cores = 1
  ram = 1024
  availability_zone = "ZONE_1"
  cpu_family = "AMD_OPTERON"
  image_name ="Debian-10-cloud-init.qcow2"
  image_password = "K3tTj8G14a3EgKyNeeiY"
  type = "ENTERPRISE"
  volume {
    name = "system"
    size = 5
    disk_type = "SSD Standard"
	backup_unit_id = ` + BackupUnitResource + `.` + BackupUnitTestResource + `.id
    user_data = "foo"
    bus = "VIRTIO"
    availability_zone = "ZONE_1"
}
  nic {
    lan = ` + LanResource + `.` + LanTestResource + `.id
    name = "system"
    dhcp = true
    firewall_active = true
	firewall_type = "BIDIRECTIONAL"
    ips            = [ ionoscloud_ipblock.webserver_ipblock.ips[0], ionoscloud_ipblock.webserver_ipblock.ips[1] ]
    firewall {
      protocol = "TCP"
      name = "SSH"
      port_range_start = 22
      port_range_end = 22
	  source_mac = "00:0a:95:9d:68:17"
	  source_ip = ionoscloud_ipblock.webserver_ipblock.ips[2]
	  target_ip = ionoscloud_ipblock.webserver_ipblock.ips[3]
	  type = "EGRESS"
    }
  }
}`
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
	DBaaSClusterResource             = "ionoscloud_pg_cluster"
	DBaaSClusterTestResource         = "test_dbaas_cluster"
	DBaaSClusterTestDataSourceById   = "test_dbaas_cluster_id"
	DBaaSClusterTestDataSourceByName = "test_dbaas_cluster_name"
	DBaaSBackupsResource             = "ionoscloud_pg_backups"
	DBaaSBackupsTest                 = "test_dbaas_backups"
	DBaaSVersionsResource            = "ionoscloud_pg_versions"
	DBaaSVersionsTest                = "test_dbaas_versions"
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
	DataSource       = "data"
	UpdatedResources = "test_updated"
)
