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
	testAccCheckServerNoPwdOrSSH = `
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
  image_name ="ubuntu:latest"
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

const sshKey = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAACAQC6J7UMVHrx2EztvbnH+xCVOo8i4sg40H4U5NNySxF5ZwmHXHDlOw8BCJCwFAjknDxJPZQgZMPUAvAYZh0gBWcZhqOXTNcDyPCusMBQvEbngiXyAfTJKdSe+lPkpOnoq7RGjdIbrnLzmxtnPNL6pk1Ys+eVBxoOt+FGkfbIhXwEv5zy82Kk2j96fKD6OrfJna7O7xQWDkhIa6GHa9S0LaU6NwWZmaZidbEAbf4/ntjKLtrIJLcc8C5ExquBVg36jdTjsnoW85tY95SScVH5qlk7zEpn9nFLbb3TKNItwewK0pf5jsjbAOXpRWQk+sn2IgayEZ8fOfmQe88mH3ZHrWqAMSvyBl/CXY3wBjHsUiUNy+Z4i3Rx3Gqa+vcUpx8r0ZaryfbrTWkA4WYEsX5Brg6JsgcA/oJ8HNcUY8dexSZMXPV1Ofl+AxkwLMjUjxSKHgfX1EkjdhzVgQraHihSgCbKZCjkEhAzASI/TOQjSPk0/6itX+359fbBE5mahfYzrDFTwDqbgJI295cZxrMH5JU/RHMMq3xzUHO20L02kQgz3By5lDhlLq65qqxbSHncqbWPlbfzqqNaJEfK0tCwuTfMEmKv8PcrF6KrLyaYJTAjYPvOiZUVOp1OlUoArGrsHG2smjgn+juOHPBOWVFSukRTIn869uKWkCWfA1hIjFEhjQ== My nginx key"

const testAccCheckServerSshDirectly = `
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
  image_name ="ubuntu:latest"
  ssh_key_path = ["` + sshKey + `"]
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

const testAccCheckServerSshKeysDirectly = `
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
  image_name ="ubuntu:latest"
  ssh_keys = ["` + sshKey + `"]
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

const testAccCheckServerSshKeysAndKeyPathErr = `
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
  image_name ="ubuntu:latest"
  ssh_keys = ["` + sshKey + `"]
  ssh_key_path = ["` + sshKey + `"]
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

const testAccCheckServerConfigBasic = `
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
  image_name ="ubuntu:latest"
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

const testAccDataSourceDatacenterWrongNameError = testAccCheckDatacenterConfigBasic + `
data ` + DatacenterResource + ` ` + DatacenterDataSourceMatching + ` {
    name = "wrong_name"
    location = ` + DatacenterResource + `.` + DatacenterTestResource + `.location
}`

const ImmutableError = "attribute is immutable, therefore not allowed in update requests"

// Cube Server Constants
const (
	testAccCheckCubeServerConfigBasic = `
data "ionoscloud_template" ` + ServerTestResource + ` {
    name = "CUBES XS"
    cores = 1
    ram   = 1024
    storage_size = 30
}
resource ` + DatacenterResource + ` ` + DatacenterTestResource + ` {
	name       = "server-test"
	location = "de/fra"
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
resource ` + ServerCubeResource + ` ` + ServerTestResource + ` {
  template_uuid     = data.ionoscloud_template.` + ServerTestResource + `.id
  name = "` + ServerTestResource + `"
  datacenter_id = ` + DatacenterResource + `.` + DatacenterTestResource + `.id
  availability_zone = "ZONE_1"
  cpu_family = "INTEL_SKYLAKE"
  image_name ="ubuntu:latest"
  image_password = "K3tTj8G14a3EgKyNeeiY"
  
  volume {
    name = "system"
    licence_type    = "LINUX"
    disk_type = "DAS"
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

const testAccCheckServerCreationWithLabels = `
resource ` + DatacenterResource + ` ` + DatacenterTestResource + ` {
	name       = "server-test"
	location = "us/las"
}
resource ` + LanResource + ` ` + LanTestResource + ` {
  datacenter_id = ` + DatacenterResource + `.` + DatacenterTestResource + `.id
  public = true
  name = "public"
}
resource "random_password" "image_password" {
  length = 16
  special = false
}
resource ` + ServerResource + ` ` + ServerTestResource + ` {
  name = "` + ServerTestResource + `"
  datacenter_id = ` + DatacenterResource + `.` + DatacenterTestResource + `.id
  cores = 1
  ram = 1024
  availability_zone = "ZONE_1"
  cpu_family = "AMD_OPTERON"
  image_name = "ubuntu:latest"
  image_password = random_password.image_password.result
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
    firewall_active = false
  }
  label {
    key = "labelkey0"
    value = "labelvalue0"
  }
  label {
    key = "labelkey1"
    value = "labelvalue1"
  }
}`

const K8sVersion = "1.23.12"
const UpgradedK8sVersion = "1.24.6"
