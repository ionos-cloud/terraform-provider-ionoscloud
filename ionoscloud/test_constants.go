package ionoscloud

const (
	testAccCheckBackupUnitConfigBasic = `
resource ` + BackupUnitResource + ` ` + BackupUnitTestResource + ` {
	name        = "` + BackupUnitTestResource + `"
	password    = "DemoPassword123$"
	email       = "example@ionoscloud.com"
}
`
)

// Datacenter constants
const (
	testAccCheckDatacenterConfigBasic = `
resource ` + DatacenterResource + ` ` + DatacenterTestResource + ` {
	name       = "` + DatacenterTestResource + `"
	location = "us/las"
	description = "Test Datacenter Description"
	sec_auth_protection = false
}`
)

// Lan Constants
const (
	testAccCheckLanConfigBasic = testAccCheckDatacenterConfigBasic + `
resource ` + LanResource + ` ` + LanTestResource + ` {
  datacenter_id = ` + DatacenterResource + `.` + DatacenterTestResource + `.id
  public = true
  name = "` + LanTestResource + `"
}`
)

// Private Crossconnect Constants
const (
	testAccCheckPrivateCrossConnectConfigBasic = `
resource ` + PCCResource + ` ` + PCCTestResource + ` {
  name        = "` + PCCTestResource + `"
  description = "` + PCCTestResource + `"
}`
)

// Server Constants
const (
	testAccCheckServerConfigBasic = `
resource ` + DatacenterResource + ` ` + DatacenterTestResource + ` {
	name       = "server-test"
	location = "us/las"
}
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
  image_name ="debian-10-genericcloud-amd64-20211011-792"
  image_password = "K3tTj8G14a3EgKyNeeiY"
  type = "ENTERPRISE"
  volume {
    name = "system"
    size = 5
    disk_type = "SSD Standard"
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

const testAccDataSourceDatacenterWrongNameError = testAccCheckDatacenterConfigBasic + `
data ` + DatacenterResource + ` ` + DatacenterDataSourceMatching + ` {
    name = "wrong_name"
    location = ` + DatacenterResource + `.` + DatacenterTestResource + `.location
}`

const ImmutableError = "attribute is immutable, therefore not allowed in update requests"
