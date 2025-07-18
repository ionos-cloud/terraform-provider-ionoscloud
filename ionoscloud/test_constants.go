package ionoscloud

import (
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils/constant"
)

const (
	testAccCheckBackupUnitConfigBasic = `
resource ` + constant.BackupUnitResource + ` ` + constant.BackupUnitTestResource + ` {
	name        = "` + constant.BackupUnitTestResource + `"
	password    = ` + constant.RandomPassword + `.backup_unit_password.result
	email       = "example@ionoscloud.com"
}

resource ` + constant.RandomPassword + ` "backup_unit_password" {
  length           = 16
  special          = true
  override_special = "!#$%&*()-_=+[]{}<>:?"
}
`
)

// Datacenter constants
const (
	testAccCheckDatacenterConfigBasic = `
resource ` + constant.DatacenterResource + ` ` + constant.DatacenterTestResource + ` {
	name       = "` + constant.DatacenterTestResource + `"
	location = "us/las"
	description = "Test Datacenter Description"
	sec_auth_protection = false
}`
)

// Lan Constants
const (
	testAccCheckLanConfigBasic = testAccCheckDatacenterConfigBasic + `
resource ` + constant.LanResource + ` ` + constant.LanTestResource + ` {
  datacenter_id = ` + constant.DatacenterResource + `.` + constant.DatacenterTestResource + `.id
  public = true
  name = "` + constant.LanTestResource + `"
}`

	privateLANConfig = testAccCheckDatacenterConfigBasic + `
resource ` + constant.LanResource + ` ` + constant.PrivateLANTestResource + ` {
  datacenter_id = ` + constant.DatacenterResource + `.` + constant.DatacenterTestResource + `.id
  public = false
  name = "` + constant.PrivateLANTestResource + `"
}`
)

// Private Crossconnect Constants
// The resource name was changed from Private Cross Connect to Cross Connect
// But the terraform resources names did not change for backwards compatibility reasons

const (
	testAccCheckPrivateCrossConnectConfigBasic = `
resource ` + constant.PCCResource + ` ` + constant.PCCTestResource + ` {
  name        = "` + constant.PCCTestResource + `"
  description = "` + constant.PCCTestResource + `"
}`
)

// Server Constants

const (
	testAccCheckServerNoPwdOrSSH = `
resource ` + constant.DatacenterResource + ` ` + constant.DatacenterTestResource + ` {
	name       = "server-test"
	location = "us/las"
}
resource "ionoscloud_ipblock" "webserver_ipblock" {
  location = ` + constant.DatacenterResource + `.` + constant.DatacenterTestResource + `.location
  size = 4
  name = "webserver_ipblock"
}
resource ` + constant.LanResource + ` ` + constant.LanTestResource + ` {
  datacenter_id = ` + constant.DatacenterResource + `.` + constant.DatacenterTestResource + `.id
  public = true
  name = "public"
}
resource ` + constant.ServerResource + ` ` + constant.ServerTestResource + ` {
  name = "` + constant.ServerTestResource + `"
  datacenter_id = ` + constant.DatacenterResource + `.` + constant.DatacenterTestResource + `.id
  cores = 1
  ram = 1024
  image_name ="ubuntu:latest"
  type = "ENTERPRISE"
  volume {
    name = "system"
    size = 5
    disk_type = "SSD Standard"
    user_data = "foo"
    bus = "VIRTIO"
  }
  nic {
    lan = ` + constant.LanResource + `.` + constant.LanTestResource + `.id
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
resource ` + constant.DatacenterResource + ` ` + constant.DatacenterTestResource + ` {
	name       = "server-test"
	location = "us/las"
}
resource "ionoscloud_ipblock" "webserver_ipblock" {
  location = ` + constant.DatacenterResource + `.` + constant.DatacenterTestResource + `.location
  size = 4
  name = "webserver_ipblock"
}
resource ` + constant.LanResource + ` ` + constant.LanTestResource + ` {
  datacenter_id = ` + constant.DatacenterResource + `.` + constant.DatacenterTestResource + `.id
  public = true
  name = "public"
}
resource ` + constant.ServerResource + ` ` + constant.ServerTestResource + ` {
  name = "` + constant.ServerTestResource + `"
  datacenter_id = ` + constant.DatacenterResource + `.` + constant.DatacenterTestResource + `.id
  cores = 1
  ram = 1024
  image_name ="ubuntu:latest"
  ssh_key_path = ["` + sshKey + `"]
  type = "ENTERPRISE"
  volume {
    name = "system"
    size = 5
    disk_type = "SSD Standard"
    user_data = "foo"
    bus = "VIRTIO"
  }
  nic {
    lan = ` + constant.LanResource + `.` + constant.LanTestResource + `.id
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
resource ` + constant.DatacenterResource + ` ` + constant.DatacenterTestResource + ` {
	name       = "server-test"
	location = "us/las"
}
resource "ionoscloud_ipblock" "webserver_ipblock" {
  location = ` + constant.DatacenterResource + `.` + constant.DatacenterTestResource + `.location
  size = 4
  name = "webserver_ipblock"
}
resource ` + constant.LanResource + ` ` + constant.LanTestResource + ` {
  datacenter_id = ` + constant.DatacenterResource + `.` + constant.DatacenterTestResource + `.id
  public = true
  name = "public"
}
resource ` + constant.ServerResource + ` ` + constant.ServerTestResource + ` {
  name = "` + constant.ServerTestResource + `"
  datacenter_id = ` + constant.DatacenterResource + `.` + constant.DatacenterTestResource + `.id
  cores = 1
  ram = 1024
  image_name ="ubuntu:latest"
  ssh_keys = ["` + sshKey + `"]
  type = "ENTERPRISE"
  volume {
    name = "system"
    size = 5
    disk_type = "SSD Standard"
    user_data = "foo"
    bus = "VIRTIO"
  }
  nic {
    lan = ` + constant.LanResource + `.` + constant.LanTestResource + `.id
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
resource ` + constant.DatacenterResource + ` ` + constant.DatacenterTestResource + ` {
	name       = "server-test"
	location = "us/las"
}
resource "ionoscloud_ipblock" "webserver_ipblock" {
  location = ` + constant.DatacenterResource + `.` + constant.DatacenterTestResource + `.location
  size = 4
  name = "webserver_ipblock"
}
resource ` + constant.LanResource + ` ` + constant.LanTestResource + ` {
  datacenter_id = ` + constant.DatacenterResource + `.` + constant.DatacenterTestResource + `.id
  public = true
  name = "public"
}
resource ` + constant.ServerResource + ` ` + constant.ServerTestResource + ` {
  name = "` + constant.ServerTestResource + `"
  datacenter_id = ` + constant.DatacenterResource + `.` + constant.DatacenterTestResource + `.id
  cores = 1
  ram = 1024
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
  }
  nic {
    lan = ` + constant.LanResource + `.` + constant.LanTestResource + `.id
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

const testSecurityGroups = `
resource ` + constant.NSGResource + ` example_1 {
  name          = "testing-name-1"
  description   = "testing-description-1"
  datacenter_id = ` + constant.DatacenterResource + `.` + constant.DatacenterTestResource + `.id
}
resource ` + constant.NSGResource + ` example_2 {
  name          = "testing-name-2"
  description   = "testing-description-2"
  datacenter_id = ` + constant.DatacenterResource + `.` + constant.DatacenterTestResource + `.id
}
`

const testAccCheckServerConfigBasic = `
resource ` + constant.DatacenterResource + ` ` + constant.DatacenterTestResource + ` {
	name       = "server-test"
	location = "us/las"
}
resource "ionoscloud_ipblock" "webserver_ipblock" {
  location = ` + constant.DatacenterResource + `.` + constant.DatacenterTestResource + `.location
  size = 4
  name = "webserver_ipblock"
}
resource ` + constant.LanResource + ` ` + constant.LanTestResource + ` {
  datacenter_id = ` + constant.DatacenterResource + `.` + constant.DatacenterTestResource + `.id
  public = true
  name = "public"
}
resource ` + constant.ServerResource + ` ` + constant.ServerTestResource + ` {
  name = "` + constant.ServerTestResource + `"
  hostname = "` + constant.ServerTestHostname + `"
  datacenter_id = ` + constant.DatacenterResource + `.` + constant.DatacenterTestResource + `.id
  cores = 1
  ram = 1024
  image_name ="ubuntu:latest"
  image_password = ` + constant.RandomPassword + `.server_image_password.result
  type = "ENTERPRISE"
  volume {
    name = "system"
    size = 5
    disk_type = "SSD Standard"
    user_data = "foo"
    bus = "VIRTIO"
}
  nic {
    lan = ` + constant.LanResource + `.` + constant.LanTestResource + `.id
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
}
` + ServerImagePassword

const testSnapshotServer = `
resource ` + constant.DatacenterResource + ` ` + constant.DatacenterTestResource + ` {
  name     = "server-test"
  location = "us/las"
}

resource ` + constant.ServerResource + ` ` + constant.ServerTestResource + ` {
  name = "` + constant.ServerTestResource + `"
  datacenter_id = ` + constant.DatacenterResource + `.` + constant.DatacenterTestResource + `.id
  cores = 1
  ram = 1024
  image_name ="ubuntu:latest"
  image_password = ` + constant.RandomPassword + `.server_image_password.result
  type = "ENTERPRISE"
  volume {
    name = "system"
    size = 5
    disk_type = "SSD Standard"
    user_data = "foo"
    bus = "VIRTIO"
  }
}
` + ServerImagePassword

// Solves  #372 crash when ips field in nic resource is a list with an empty string
const testAccCheckServerConfigEmptyNicIps = `
resource ` + constant.DatacenterResource + ` ` + constant.DatacenterTestResource + ` {
	name       = "server-test"
	location = "us/las"
}
resource "ionoscloud_ipblock" "webserver_ipblock" {
  location = ` + constant.DatacenterResource + `.` + constant.DatacenterTestResource + `.location
  size = 4
  name = "webserver_ipblock"
}
resource ` + constant.LanResource + ` ` + constant.LanTestResource + ` {
  datacenter_id = ` + constant.DatacenterResource + `.` + constant.DatacenterTestResource + `.id
  public = true
  name = "public"
}
resource ` + constant.ServerResource + ` ` + constant.ServerTestResource + ` {
  name = "` + constant.ServerTestResource + `"
  datacenter_id = ` + constant.DatacenterResource + `.` + constant.DatacenterTestResource + `.id
  cores = 1
  ram = 1024
  image_name ="ubuntu:latest"
  image_password = ` + constant.RandomPassword + `.server_image_password.result
  type = "ENTERPRISE"
  volume {
    name = "system"
    size = 5
    disk_type = "SSD Standard"
    user_data = "foo"
    bus = "VIRTIO"
}
  nic {
    lan = ` + constant.LanResource + `.` + constant.LanTestResource + `.id
    name = "system"
    dhcp = true
    firewall_active = true
	firewall_type = "BIDIRECTIONAL"
    ips            = [""]
    firewall {
      protocol = "TCP"
      name = "SSH"
      port_range_start = 22
      port_range_end = 23
	  source_mac = "00:0a:95:9d:68:17"
	  source_ip = ionoscloud_ipblock.webserver_ipblock.ips[2]
	  target_ip = ionoscloud_ipblock.webserver_ipblock.ips[3]
	  type = "EGRESS"
    }
  }
}
` + ServerImagePassword

const testAccCheckServerConfigIpv6Enabled = `
resource ` + constant.DatacenterResource + ` ` + constant.DatacenterTestResource + ` {
	name       = "server-test"
	location = "us/las"
}
resource "ionoscloud_ipblock" "webserver_ipblock" {
  location = ` + constant.DatacenterResource + `.` + constant.DatacenterTestResource + `.location
  size = 4
  name = "webserver_ipblock"
}
resource ` + constant.LanResource + ` ` + constant.LanTestResource + ` {
  datacenter_id = ` + constant.DatacenterResource + `.` + constant.DatacenterTestResource + `.id
  public = true
  name = "public"
  ipv6_cidr_block = cidrsubnet(` + constant.DatacenterResource + `.` + constant.DatacenterTestResource + `.ipv6_cidr_block` + `,8,10)
}
resource ` + constant.ServerResource + ` ` + constant.ServerTestResource + ` {
  name = "` + constant.ServerTestResource + `"
  hostname = "updated"
  datacenter_id = ` + constant.DatacenterResource + `.` + constant.DatacenterTestResource + `.id
  cores = 1
  ram = 1024
  image_name ="ubuntu:latest"
  image_password = ` + constant.RandomPassword + `.server_image_password.result
  type = "ENTERPRISE"
  volume {
    name = "system"
    size = 5
    disk_type = "SSD Standard"
    user_data = "foo"
    bus = "VIRTIO"
  }
  nic {
    lan = ` + constant.LanResource + `.` + constant.LanTestResource + `.id
    name = "system"
    dhcp = true
    firewall_active = true
    firewall_type = "BIDIRECTIONAL"
    ips  = [ ionoscloud_ipblock.webserver_ipblock.ips[0], ionoscloud_ipblock.webserver_ipblock.ips[1] ]
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
    dhcpv6 = true
    ipv6_cidr_block = cidrsubnet(` + constant.LanResource + `.` + constant.LanTestResource + `.ipv6_cidr_block,16,10)
    ipv6_ips        = [ 
                        cidrhost(cidrsubnet(` + constant.LanResource + `.` + constant.LanTestResource + `.ipv6_cidr_block,16,10),1),
                        cidrhost(cidrsubnet(` + constant.LanResource + `.` + constant.LanTestResource + `.ipv6_cidr_block,16,10),2),
                        cidrhost(cidrsubnet(` + constant.LanResource + `.` + constant.LanTestResource + `.ipv6_cidr_block,16,10),3)
                      ]
  }
}
` + ServerImagePassword + `

data ` + constant.ServerResource + ` ` + constant.ServerDataSourceById + ` {
  datacenter_id = ` + constant.DatacenterResource + `.` + constant.DatacenterTestResource + `.id
  id = ` + constant.ServerResource + `.` + constant.ServerTestResource + `.id
}
`

const testAccCheckServerConfigIpv6Update = `
resource ` + constant.DatacenterResource + ` ` + constant.DatacenterTestResource + ` {
	name       = "server-test"
	location = "us/las"
}
resource "ionoscloud_ipblock" "webserver_ipblock" {
  location = ` + constant.DatacenterResource + `.` + constant.DatacenterTestResource + `.location
  size = 4
  name = "webserver_ipblock"
}
resource ` + constant.LanResource + ` ` + constant.LanTestResource + ` {
  datacenter_id = ` + constant.DatacenterResource + `.` + constant.DatacenterTestResource + `.id
  public = true
  name = "public"
  ipv6_cidr_block = cidrsubnet(` + constant.DatacenterResource + `.` + constant.DatacenterTestResource + `.ipv6_cidr_block` + `,8,20)
}
resource ` + constant.ServerResource + ` ` + constant.ServerTestResource + ` {
  name = "` + constant.ServerTestResource + `"
  datacenter_id = ` + constant.DatacenterResource + `.` + constant.DatacenterTestResource + `.id
  cores = 1
  ram = 1024
  image_name ="ubuntu:latest"
  image_password = ` + constant.RandomPassword + `.server_image_password.result
  type = "ENTERPRISE"
  volume {
    name = "system"
    size = 5
    disk_type = "SSD Standard"
    user_data = "foo"
    bus = "VIRTIO"
  }
  nic {
    lan = ` + constant.LanResource + `.` + constant.LanTestResource + `.id
    name = "system"
    dhcp = true
    firewall_active = true
    firewall_type = "BIDIRECTIONAL"
    ips  = [ ionoscloud_ipblock.webserver_ipblock.ips[0], ionoscloud_ipblock.webserver_ipblock.ips[1] ]
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
    dhcpv6 = false
    ipv6_cidr_block = cidrsubnet(` + constant.LanResource + `.` + constant.LanTestResource + `.ipv6_cidr_block,16,24)
    ipv6_ips        = [ 
                        cidrhost(cidrsubnet(` + constant.LanResource + `.` + constant.LanTestResource + `.ipv6_cidr_block,16,24),11),
                        cidrhost(cidrsubnet(` + constant.LanResource + `.` + constant.LanTestResource + `.ipv6_cidr_block,16,24),12),
                        cidrhost(cidrsubnet(` + constant.LanResource + `.` + constant.LanTestResource + `.ipv6_cidr_block,16,24),13)
                      ]
  }
}
` + ServerImagePassword + `

data ` + constant.ServerResource + ` ` + constant.ServerDataSourceById + ` {
  datacenter_id = ` + constant.DatacenterResource + `.` + constant.DatacenterTestResource + `.id
  id = ` + constant.ServerResource + `.` + constant.ServerTestResource + `.id
}
`

const testAccCheckServerConfigShutdown = `
resource ` + constant.DatacenterResource + ` ` + constant.DatacenterTestResource + ` {
	name       = "server-test"
	location = "us/las"
}
resource "ionoscloud_ipblock" "webserver_ipblock" {
  location = ` + constant.DatacenterResource + `.` + constant.DatacenterTestResource + `.location
  size = 4
  name = "webserver_ipblock"
}
resource ` + constant.LanResource + ` ` + constant.LanTestResource + ` {
  datacenter_id = ` + constant.DatacenterResource + `.` + constant.DatacenterTestResource + `.id
  public = true
  name = "public"
}
resource ` + constant.ServerResource + ` ` + constant.ServerTestResource + ` {
  name = "` + constant.ServerTestResource + `"
  datacenter_id = ` + constant.DatacenterResource + `.` + constant.DatacenterTestResource + `.id
  vm_state = "` + constant.VMStateStop + `"
  cores = 1
  ram = 1024
  image_name ="ubuntu:latest"
  image_password = ` + constant.RandomPassword + `.server_image_password.result
  type = "ENTERPRISE"
  volume {
    name = "system"
    size = 5
    disk_type = "SSD Standard"
    user_data = "foo"
    bus = "VIRTIO"
  }
  nic {
    lan = ` + constant.LanResource + `.` + constant.LanTestResource + `.id
    name = "system"
    dhcp = true
    firewall_active = true
    firewall_type = "BIDIRECTIONAL"
    ips  = [ ionoscloud_ipblock.webserver_ipblock.ips[0], ionoscloud_ipblock.webserver_ipblock.ips[1] ]
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
}
` + ServerImagePassword + `

data ` + constant.ServerResource + ` ` + constant.ServerDataSourceById + ` {
  datacenter_id = ` + constant.DatacenterResource + `.` + constant.DatacenterTestResource + `.id
  id = ` + constant.ServerResource + `.` + constant.ServerTestResource + `.id
}
`

const testAccCheckServerConfigPowerOn = `
resource ` + constant.DatacenterResource + ` ` + constant.DatacenterTestResource + ` {
	name       = "server-test"
	location = "us/las"
}
resource "ionoscloud_ipblock" "webserver_ipblock" {
  location = ` + constant.DatacenterResource + `.` + constant.DatacenterTestResource + `.location
  size = 4
  name = "webserver_ipblock"
}
resource ` + constant.LanResource + ` ` + constant.LanTestResource + ` {
  datacenter_id = ` + constant.DatacenterResource + `.` + constant.DatacenterTestResource + `.id
  public = true
  name = "public"
}
resource ` + constant.ServerResource + ` ` + constant.ServerTestResource + ` {
  name = "` + constant.ServerTestResource + `"
  datacenter_id = ` + constant.DatacenterResource + `.` + constant.DatacenterTestResource + `.id
  vm_state = "` + constant.VMStateStart + `"
  cores = 1
  ram = 1024
  image_name ="ubuntu:latest"
  image_password = ` + constant.RandomPassword + `.server_image_password.result
  type = "ENTERPRISE"
  volume {
    name = "system"
    size = 5
    disk_type = "SSD Standard"
    user_data = "foo"
    bus = "VIRTIO"
  }
  nic {
    lan = ` + constant.LanResource + `.` + constant.LanTestResource + `.id
    name = "system"
    dhcp = true
    firewall_active = true
    firewall_type = "BIDIRECTIONAL"
    ips  = [ ionoscloud_ipblock.webserver_ipblock.ips[0], ionoscloud_ipblock.webserver_ipblock.ips[1] ]
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
}
` + ServerImagePassword + `

data ` + constant.ServerResource + ` ` + constant.ServerDataSourceById + ` {
  datacenter_id = ` + constant.DatacenterResource + `.` + constant.DatacenterTestResource + `.id
  id = ` + constant.ServerResource + `.` + constant.ServerTestResource + `.id
}
`

const testAccDataSourceDatacenterWrongNameError = testAccCheckDatacenterConfigBasic + `
data ` + constant.DatacenterResource + ` ` + constant.DatacenterDataSourceMatching + ` {
    name = "wrong_name"
    location = ` + constant.DatacenterResource + `.` + constant.DatacenterTestResource + `.location
}`

const ImmutableError = "attribute is immutable, therefore not allowed in update requests"

const ServerImagePassword = `
resource ` + constant.RandomPassword + ` "server_image_password" {
  length           = 16
  special          = false
}
`

const ServerImagePasswordUpdated = `
resource ` + constant.RandomPassword + ` "server_image_password_updated" {
  length           = 16
  special          = false
}
`

// Cube Server Constants
const testAccCheckCubeServerConfigBasic = `
data "ionoscloud_template" ` + constant.ServerTestResource + ` {
    name = "Basic Cube XS"
    cores = 1
    ram   = 2048
    storage_size = 60
}
resource ` + constant.DatacenterResource + ` ` + constant.DatacenterTestResource + ` {
	name       = "server-test"
	location = "de/fra"
}
resource "ionoscloud_ipblock" "webserver_ipblock" {
  location = ` + constant.DatacenterResource + `.` + constant.DatacenterTestResource + `.location
  size = 4
  name = "webserver_ipblock"
}
resource ` + constant.LanResource + ` ` + constant.LanTestResource + ` {
  datacenter_id = ` + constant.DatacenterResource + `.` + constant.DatacenterTestResource + `.id
  public = true
  name = "public"
}
resource ` + constant.ServerCubeResource + ` ` + constant.ServerTestResource + ` {
  template_uuid     = data.ionoscloud_template.` + constant.ServerTestResource + `.id
  name = "` + constant.ServerTestResource + `"
  hostname = "` + constant.ServerTestHostname + `"
  datacenter_id = ` + constant.DatacenterResource + `.` + constant.DatacenterTestResource + `.id
  image_name ="ubuntu:latest"
  image_password = ` + constant.RandomPassword + `.server_image_password.result
  volume {
    name = "system"
    licence_type    = "LINUX"
    disk_type = "DAS"
	expose_serial = true
  }
  nic {
    lan = ` + constant.LanResource + `.` + constant.LanTestResource + `.id
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
}
` + ServerImagePassword

const testAccCheckCubeServerEnableIpv6 = `
data "ionoscloud_template" ` + constant.ServerTestResource + ` {
    name = "Basic Cube XS"
    cores = 1
    ram   = 2048
    storage_size = 60
}
resource ` + constant.DatacenterResource + ` ` + constant.DatacenterTestResource + ` {
  name       = "server-test"
  location = "de/fra"
}
resource "ionoscloud_ipblock" "webserver_ipblock" {
  location = ` + constant.DatacenterResource + `.` + constant.DatacenterTestResource + `.location
  size = 4
  name = "webserver_ipblock"
}
resource ` + constant.LanResource + ` ` + constant.LanTestResource + ` {
  datacenter_id = ` + constant.DatacenterResource + `.` + constant.DatacenterTestResource + `.id
  public = true
  name = "public"
  ipv6_cidr_block = cidrsubnet(` + constant.DatacenterResource + `.` + constant.DatacenterTestResource + `.ipv6_cidr_block` + `,8,10)
}
resource ` + constant.ServerCubeResource + ` ` + constant.ServerTestResource + ` {
  template_uuid     = data.ionoscloud_template.` + constant.ServerTestResource + `.id
  name = "` + constant.ServerTestResource + `"
  datacenter_id = ` + constant.DatacenterResource + `.` + constant.DatacenterTestResource + `.id
  image_name ="ubuntu:latest"
  image_password = ` + constant.RandomPassword + `.server_image_password.result

  volume {
    name = "system"
    licence_type    = "LINUX"
    disk_type = "DAS"
}
  nic {
    lan = ` + constant.LanResource + `.` + constant.LanTestResource + `.id
    name = "system"
    dhcp = true
    firewall_active = true
    firewall_type = "BIDIRECTIONAL"
    ips            = [ ionoscloud_ipblock.webserver_ipblock.ips[0], ionoscloud_ipblock.webserver_ipblock.ips[1] ]

    dhcpv6 = true
    ipv6_cidr_block = cidrsubnet(` + constant.LanResource + `.` + constant.LanTestResource + `.ipv6_cidr_block,16,12)
    ipv6_ips        = [ 
                        cidrhost(cidrsubnet(` + constant.LanResource + `.` + constant.LanTestResource + `.ipv6_cidr_block,16,12),1),
                        cidrhost(cidrsubnet(` + constant.LanResource + `.` + constant.LanTestResource + `.ipv6_cidr_block,16,12),2),
                        cidrhost(cidrsubnet(` + constant.LanResource + `.` + constant.LanTestResource + `.ipv6_cidr_block,16,12),3)
                      ]

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
}
` + ServerImagePassword + `
data ` + constant.ServerCubeResource + ` ` + constant.ServerDataSourceById + ` {
  datacenter_id = ` + constant.DatacenterResource + `.` + constant.DatacenterTestResource + `.id
  id = ` + constant.ServerCubeResource + `.` + constant.ServerTestResource + `.id
}
`
const testAccCheckCubeServerUpdateIpv6 = `
data "ionoscloud_template" ` + constant.ServerTestResource + ` {
    name = "Basic Cube XS"
    cores = 1
    ram   = 2048
    storage_size = 60
}
resource ` + constant.DatacenterResource + ` ` + constant.DatacenterTestResource + ` {
  name       = "server-test"
  location   = "de/fra"
}
resource "ionoscloud_ipblock" "webserver_ipblock" {
  location = ` + constant.DatacenterResource + `.` + constant.DatacenterTestResource + `.location
  size = 4
  name = "webserver_ipblock"
}
resource ` + constant.LanResource + ` ` + constant.LanTestResource + ` {
  datacenter_id = ` + constant.DatacenterResource + `.` + constant.DatacenterTestResource + `.id
  public = true
  name = "public"
  ipv6_cidr_block = cidrsubnet(` + constant.DatacenterResource + `.` + constant.DatacenterTestResource + `.ipv6_cidr_block` + `,8,10)
}
resource ` + constant.ServerCubeResource + ` ` + constant.ServerTestResource + ` {
  template_uuid     = data.ionoscloud_template.` + constant.ServerTestResource + `.id
  name = "` + constant.ServerTestResource + `"
  datacenter_id = ` + constant.DatacenterResource + `.` + constant.DatacenterTestResource + `.id
  image_name ="ubuntu:latest"
  image_password = ` + constant.RandomPassword + `.server_image_password.result
  
  volume {
    name = "system"
    licence_type    = "LINUX"
    disk_type = "DAS"
}
  nic {
    lan = ` + constant.LanResource + `.` + constant.LanTestResource + `.id
    name = "system"
    dhcp = true
    firewall_active = true
    firewall_type = "BIDIRECTIONAL"
    ips            = [ ionoscloud_ipblock.webserver_ipblock.ips[0], ionoscloud_ipblock.webserver_ipblock.ips[1] ]

    dhcpv6 = false
    ipv6_cidr_block = cidrsubnet(` + constant.LanResource + `.` + constant.LanTestResource + `.ipv6_cidr_block,16,16)
    ipv6_ips        = [ 
                        cidrhost(cidrsubnet(` + constant.LanResource + `.` + constant.LanTestResource + `.ipv6_cidr_block,16,16),10),
                        cidrhost(cidrsubnet(` + constant.LanResource + `.` + constant.LanTestResource + `.ipv6_cidr_block,16,16),20),
                        cidrhost(cidrsubnet(` + constant.LanResource + `.` + constant.LanTestResource + `.ipv6_cidr_block,16,16),30)
                      ]

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
}
` + ServerImagePassword + `
data ` + constant.ServerCubeResource + ` ` + constant.ServerDataSourceById + ` {
  datacenter_id = ` + constant.DatacenterResource + `.` + constant.DatacenterTestResource + `.id
  id = ` + constant.ServerCubeResource + `.` + constant.ServerTestResource + `.id
}
`
const testAccCheckCubeServerSuspend = `
data "ionoscloud_template" ` + constant.ServerTestResource + ` {
    name = "Basic Cube XS"
    cores = 1
    ram   = 2048
    storage_size = 60
}
resource ` + constant.DatacenterResource + ` ` + constant.DatacenterTestResource + ` {
  name       = "server-test"
  location   = "de/fra"
}
resource "ionoscloud_ipblock" "webserver_ipblock" {
  location = ` + constant.DatacenterResource + `.` + constant.DatacenterTestResource + `.location
  size = 4
  name = "webserver_ipblock"
}
resource ` + constant.LanResource + ` ` + constant.LanTestResource + ` {
  datacenter_id = ` + constant.DatacenterResource + `.` + constant.DatacenterTestResource + `.id
  public = true
  name = "public"
  ipv6_cidr_block = cidrsubnet(` + constant.DatacenterResource + `.` + constant.DatacenterTestResource + `.ipv6_cidr_block` + `,8,10)
}
resource ` + constant.ServerCubeResource + ` ` + constant.ServerTestResource + ` {
  template_uuid     = data.ionoscloud_template.` + constant.ServerTestResource + `.id
  vm_state = "` + constant.CubeVMStateStop + `"
  name = "` + constant.ServerTestResource + `"
  datacenter_id = ` + constant.DatacenterResource + `.` + constant.DatacenterTestResource + `.id
  image_name ="ubuntu:latest"
  image_password = ` + constant.RandomPassword + `.server_image_password.result
  
  volume {
    name = "system"
    licence_type    = "LINUX"
    disk_type = "DAS"
}
  nic {
    lan = ` + constant.LanResource + `.` + constant.LanTestResource + `.id
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
}
` + ServerImagePassword + `
data ` + constant.ServerCubeResource + ` ` + constant.ServerDataSourceById + ` {
  datacenter_id = ` + constant.DatacenterResource + `.` + constant.DatacenterTestResource + `.id
  id = ` + constant.ServerCubeResource + `.` + constant.ServerTestResource + `.id
}
`
const testAccCheckCubeServerResume = `
data "ionoscloud_template" ` + constant.ServerTestResource + ` {
    name = "Basic Cube XS"
    cores = 1
    ram   = 2048
    storage_size = 60
}
resource ` + constant.DatacenterResource + ` ` + constant.DatacenterTestResource + ` {
  name       = "server-test"
  location   = "de/fra"
}
resource "ionoscloud_ipblock" "webserver_ipblock" {
  location = ` + constant.DatacenterResource + `.` + constant.DatacenterTestResource + `.location
  size = 4
  name = "webserver_ipblock"
}
resource ` + constant.LanResource + ` ` + constant.LanTestResource + ` {
  datacenter_id = ` + constant.DatacenterResource + `.` + constant.DatacenterTestResource + `.id
  public = true
  name = "public"
  ipv6_cidr_block = cidrsubnet(` + constant.DatacenterResource + `.` + constant.DatacenterTestResource + `.ipv6_cidr_block` + `,8,10)
}
resource ` + constant.ServerCubeResource + ` ` + constant.ServerTestResource + ` {
  template_uuid     = data.ionoscloud_template.` + constant.ServerTestResource + `.id
  name = "` + constant.ServerTestResource + `"
  vm_state = "` + constant.VMStateStart + `"
  datacenter_id = ` + constant.DatacenterResource + `.` + constant.DatacenterTestResource + `.id
  image_name ="ubuntu:latest"
  image_password = ` + constant.RandomPassword + `.server_image_password.result
  
  volume {
    name = "system"
    licence_type    = "LINUX"
    disk_type = "DAS"
}
  nic {
    lan = ` + constant.LanResource + `.` + constant.LanTestResource + `.id
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
}
` + ServerImagePassword + `
data ` + constant.ServerCubeResource + ` ` + constant.ServerDataSourceById + ` {
  datacenter_id = ` + constant.DatacenterResource + `.` + constant.DatacenterTestResource + `.id
  id = ` + constant.ServerCubeResource + `.` + constant.ServerTestResource + `.id
}
`

const testAccCheckCubeServerUpdateWhenSuspended = `
data "ionoscloud_template" ` + constant.ServerTestResource + ` {
    name = "Basic Cube XS"
    cores = 1
    ram   = 2048
    storage_size = 60
}
resource ` + constant.DatacenterResource + ` ` + constant.DatacenterTestResource + ` {
  name       = "server-test"
  location   = "de/fra"
}
resource "ionoscloud_ipblock" "webserver_ipblock" {
  location = ` + constant.DatacenterResource + `.` + constant.DatacenterTestResource + `.location
  size = 4
  name = "webserver_ipblock"
}
resource ` + constant.LanResource + ` ` + constant.LanTestResource + ` {
  datacenter_id = ` + constant.DatacenterResource + `.` + constant.DatacenterTestResource + `.id
  public = true
  name = "public"
  ipv6_cidr_block = cidrsubnet(` + constant.DatacenterResource + `.` + constant.DatacenterTestResource + `.ipv6_cidr_block` + `,8,10)
}
resource ` + constant.ServerCubeResource + ` ` + constant.ServerTestResource + ` {
  template_uuid     = data.ionoscloud_template.` + constant.ServerTestResource + `.id
  vm_state = "` + constant.CubeVMStateStop + `"
  name = "` + constant.UpdatedResources + `"
  datacenter_id = ` + constant.DatacenterResource + `.` + constant.DatacenterTestResource + `.id
  image_name ="ubuntu:latest"
  image_password = ` + constant.RandomPassword + `.server_image_password.result
  
  volume {
    name = "system"
    licence_type    = "LINUX"
    disk_type = "DAS"
}
  nic {
    lan = ` + constant.LanResource + `.` + constant.LanTestResource + `.id
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
}
` + ServerImagePassword

const testAccCheckServerCreationWithLabels = `
resource ` + constant.DatacenterResource + ` ` + constant.DatacenterTestResource + ` {
	name       = "server-test"
	location = "us/las"
}
resource ` + constant.LanResource + ` ` + constant.LanTestResource + ` {
  datacenter_id = ` + constant.DatacenterResource + `.` + constant.DatacenterTestResource + `.id
  public = true
  name = "public"
}
resource "random_password" "image_password" {
  length = 16
  special = false
}
resource ` + constant.ServerResource + ` ` + constant.ServerTestResource + ` {
  name = "` + constant.ServerTestResource + `"
  datacenter_id = ` + constant.DatacenterResource + `.` + constant.DatacenterTestResource + `.id
  cores = 1
  ram = 1024
  image_name = "ubuntu:latest"
  image_password = random_password.image_password.result
  type = "ENTERPRISE"
  volume {
    name = "system"
    size = 5
    disk_type = "SSD Standard"
    user_data = "foo"
    bus = "VIRTIO"
}
  nic {
    lan = ` + constant.LanResource + `.` + constant.LanTestResource + `.id
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

const (
	testAccCheckServerNoNic = `
resource ` + constant.DatacenterResource + ` ` + constant.DatacenterTestResource + ` {
  name       = "server-test"
  location = "us/las"
}
resource "ionoscloud_ipblock" "webserver_ipblock" {
  location = ` + constant.DatacenterResource + `.` + constant.DatacenterTestResource + `.location
  size = 4
  name = "webserver_ipblock"
}
resource ` + constant.LanResource + ` ` + constant.LanTestResource + ` {
  datacenter_id = ` + constant.DatacenterResource + `.` + constant.DatacenterTestResource + `.id
  public = true
  name = "public"
}
resource ` + constant.ServerResource + ` ` + constant.ServerTestResource + ` {
  name = "` + constant.ServerTestResource + `"
  hostname = "` + constant.ServerTestHostname + `"
  datacenter_id = ` + constant.DatacenterResource + `.` + constant.DatacenterTestResource + `.id
  cores = 1
  ram = 1024
  image_name ="ubuntu:latest"
  ssh_key_path = ["` + sshKey + `"]
  type = "ENTERPRISE"
  volume {
    name = "system"
    size = 5
    disk_type = "SSD Standard"
    user_data = "foo"
    bus = "VIRTIO"
  }
}`

	testAccCheckServerNoNicUpdate = `
resource ` + constant.DatacenterResource + ` ` + constant.DatacenterTestResource + ` {
	name       = "server-test"
	location = "us/las"
}
resource "ionoscloud_ipblock" "webserver_ipblock" {
  location = ` + constant.DatacenterResource + `.` + constant.DatacenterTestResource + `.location
  size = 4
  name = "webserver_ipblock"
}
resource ` + constant.LanResource + ` ` + constant.LanTestResource + ` {
  datacenter_id = ` + constant.DatacenterResource + `.` + constant.DatacenterTestResource + `.id
  public = true
  name = "public"
}
resource ` + constant.ServerResource + ` ` + constant.ServerTestResource + ` {
  name = "` + constant.ServerTestResource + `"
  hostname = "` + constant.ServerTestHostname + `"
  datacenter_id = ` + constant.DatacenterResource + `.` + constant.DatacenterTestResource + `.id
  cores = 2
  ram = 2048
  image_name ="ubuntu:latest"
  ssh_key_path = ["` + sshKey + `"]
  type = "ENTERPRISE"
  volume {
    name = "system"
    size = 5
    disk_type = "SSD Standard"
    user_data = "foo"
    bus = "VIRTIO"
  }
}`
)

// VCPU Server constants

const (
	testAccCheckServerVCPUNoPwdOrSSH = `
resource ` + constant.DatacenterResource + ` ` + constant.DatacenterTestResource + ` {
	name       = "server-test"
	location   = "de/txl"
}
resource "ionoscloud_ipblock" "webserver_ipblock" {
  location = ` + constant.DatacenterResource + `.` + constant.DatacenterTestResource + `.location
  size = 4
  name = "webserver_ipblock"
}
resource ` + constant.LanResource + ` ` + constant.LanTestResource + ` {
  datacenter_id = ` + constant.DatacenterResource + `.` + constant.DatacenterTestResource + `.id
  public = true
  name = "public"
}
resource ` + constant.ServerVCPUResource + ` ` + constant.ServerTestResource + ` {
  name = "` + constant.ServerTestResource + `"
  datacenter_id = ` + constant.DatacenterResource + `.` + constant.DatacenterTestResource + `.id
  cores = 1
  ram = 1024
  image_name ="ubuntu:latest"
  volume {
    name = "system"
    size = 5
    disk_type = "SSD Standard"
    user_data = "foo"
    bus = "VIRTIO"
    expose_serial = true
  }
  nic {
    lan = ` + constant.LanResource + `.` + constant.LanTestResource + `.id
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

const testAccCheckServerVCPUSshDirectly = `
resource ` + constant.DatacenterResource + ` ` + constant.DatacenterTestResource + ` {
	name       = "server-test"
	location   = "de/txl"
}
resource "ionoscloud_ipblock" "webserver_ipblock" {
  location = ` + constant.DatacenterResource + `.` + constant.DatacenterTestResource + `.location
  size = 4
  name = "webserver_ipblock"
}
resource ` + constant.LanResource + ` ` + constant.LanTestResource + ` {
  datacenter_id = ` + constant.DatacenterResource + `.` + constant.DatacenterTestResource + `.id
  public = true
  name = "public"
}
resource ` + constant.ServerVCPUResource + ` ` + constant.ServerTestResource + ` {
  name = "` + constant.ServerTestResource + `"
  datacenter_id = ` + constant.DatacenterResource + `.` + constant.DatacenterTestResource + `.id
  cores = 1
  ram = 1024
  image_name ="ubuntu:latest"
  ssh_key_path = ["` + sshKey + `"]
  volume {
    name = "system"
    size = 5
    disk_type = "SSD Standard"
    user_data = "foo"
    bus = "VIRTIO"
  }
  nic {
    lan = ` + constant.LanResource + `.` + constant.LanTestResource + `.id
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

const testAccCheckServerVCPUSshKeysDirectly = `
resource ` + constant.DatacenterResource + ` ` + constant.DatacenterTestResource + ` {
	name       = "server-test"
	location   = "de/txl"
}
resource "ionoscloud_ipblock" "webserver_ipblock" {
  location = ` + constant.DatacenterResource + `.` + constant.DatacenterTestResource + `.location
  size = 4
  name = "webserver_ipblock"
}
resource ` + constant.LanResource + ` ` + constant.LanTestResource + ` {
  datacenter_id = ` + constant.DatacenterResource + `.` + constant.DatacenterTestResource + `.id
  public = true
  name = "public"
}
resource ` + constant.ServerVCPUResource + ` ` + constant.ServerTestResource + ` {
  name = "` + constant.ServerTestResource + `"
  datacenter_id = ` + constant.DatacenterResource + `.` + constant.DatacenterTestResource + `.id
  cores = 1
  ram = 1024
  image_name ="ubuntu:latest"
  ssh_keys = ["` + sshKey + `"]
  volume {
    name = "system"
    size = 5
    disk_type = "SSD Standard"
    user_data = "foo"
    bus = "VIRTIO"
	expose_serial = true
  }
  nic {
    lan = ` + constant.LanResource + `.` + constant.LanTestResource + `.id
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

const testAccCheckServerVCPUSshKeysAndKeyPathErr = `
resource ` + constant.DatacenterResource + ` ` + constant.DatacenterTestResource + ` {
	name       = "server-test"
	location   = "de/txl"
}
resource "ionoscloud_ipblock" "webserver_ipblock" {
  location = ` + constant.DatacenterResource + `.` + constant.DatacenterTestResource + `.location
  size = 4
  name = "webserver_ipblock"
}
resource ` + constant.LanResource + ` ` + constant.LanTestResource + ` {
  datacenter_id = ` + constant.DatacenterResource + `.` + constant.DatacenterTestResource + `.id
  public = true
  name = "public"
}
resource ` + constant.ServerVCPUResource + ` ` + constant.ServerTestResource + ` {
  name = "` + constant.ServerTestResource + `"
  datacenter_id = ` + constant.DatacenterResource + `.` + constant.DatacenterTestResource + `.id
  cores = 1
  ram = 1024
  image_name ="ubuntu:latest"
  ssh_keys = ["` + sshKey + `"]
  ssh_key_path = ["` + sshKey + `"]
  volume {
    name = "system"
    size = 5
    disk_type = "SSD Standard"
    user_data = "foo"
    bus = "VIRTIO"
  }
  nic {
    lan = ` + constant.LanResource + `.` + constant.LanTestResource + `.id
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

const testAccCheckServerVCPUConfigBasic = `
resource ` + constant.DatacenterResource + ` ` + constant.DatacenterTestResource + ` {
	name       = "server-test"
	location   = "de/txl"
}
resource "ionoscloud_ipblock" "webserver_ipblock" {
  location = ` + constant.DatacenterResource + `.` + constant.DatacenterTestResource + `.location
  size = 4
  name = "webserver_ipblock"
}
resource ` + constant.LanResource + ` ` + constant.LanTestResource + ` {
  datacenter_id = ` + constant.DatacenterResource + `.` + constant.DatacenterTestResource + `.id
  public = true
  name = "public"
}
resource ` + constant.ServerVCPUResource + ` ` + constant.ServerTestResource + ` {
  name = "` + constant.ServerTestResource + `"
  hostname = "` + constant.ServerTestHostname + `"
  datacenter_id = ` + constant.DatacenterResource + `.` + constant.DatacenterTestResource + `.id
  cores = 1
  ram = 1024
  image_name ="ubuntu:latest"
  image_password = ` + constant.RandomPassword + `.server_image_password.result
  volume {
    name = "system"
    size = 5
    disk_type = "SSD Standard"
    user_data = "foo"
    bus = "VIRTIO"
    expose_serial = true
}
  nic {
    lan = ` + constant.LanResource + `.` + constant.LanTestResource + `.id
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
}
` + ServerImagePassword

// Solves  #372 crash when ips field in nic resource is a list with an empty string
const testAccCheckServerVCPUConfigEmptyNicIps = `
resource ` + constant.DatacenterResource + ` ` + constant.DatacenterTestResource + ` {
	name       = "server-test"
	location   = "de/txl"
}
resource "ionoscloud_ipblock" "webserver_ipblock" {
  location = ` + constant.DatacenterResource + `.` + constant.DatacenterTestResource + `.location
  size = 4
  name = "webserver_ipblock"
}
resource ` + constant.LanResource + ` ` + constant.LanTestResource + ` {
  datacenter_id = ` + constant.DatacenterResource + `.` + constant.DatacenterTestResource + `.id
  public = true
  name = "public"
}
resource ` + constant.ServerVCPUResource + ` ` + constant.ServerTestResource + ` {
  name = "` + constant.ServerTestResource + `"
  datacenter_id = ` + constant.DatacenterResource + `.` + constant.DatacenterTestResource + `.id
  cores = 1
  ram = 1024
  image_name ="ubuntu:latest"
  image_password = ` + constant.RandomPassword + `.server_image_password.result
  volume {
    name = "system"
    size = 5
    disk_type = "SSD Standard"
    user_data = "foo"
    bus = "VIRTIO"
    expose_serial = true
}
  nic {
    lan = ` + constant.LanResource + `.` + constant.LanTestResource + `.id
    name = "system"
    dhcp = true
    firewall_active = true
    mac = "` + constant.NicMac + `"
	  firewall_type = "BIDIRECTIONAL"
    ips            = [""]
    firewall {
      protocol = "TCP"
      name = "SSH"
      port_range_start = 22
      port_range_end = 23
	  source_mac = "00:0a:95:9d:68:17"
	  source_ip = ionoscloud_ipblock.webserver_ipblock.ips[2]
	  target_ip = ionoscloud_ipblock.webserver_ipblock.ips[3]
	  type = "EGRESS"
    }
  }
}
` + ServerImagePassword

const testAccCheckServerVCPUCreationWithLabels = `
resource ` + constant.DatacenterResource + ` ` + constant.DatacenterTestResource + ` {
	name       = "server-test"
	location   = "de/txl"
}
resource ` + constant.LanResource + ` ` + constant.LanTestResource + ` {
  datacenter_id = ` + constant.DatacenterResource + `.` + constant.DatacenterTestResource + `.id
  public = true
  name = "public"
}
resource "random_password" "image_password" {
  length = 16
  special = false
}
resource ` + constant.ServerVCPUResource + ` ` + constant.ServerTestResource + ` {
  name = "` + constant.ServerTestResource + `"
  datacenter_id = ` + constant.DatacenterResource + `.` + constant.DatacenterTestResource + `.id
  cores = 1
  ram = 1024
  image_name = "ubuntu:latest"
  image_password = random_password.image_password.result
  volume {
    name = "system"
    size = 5
    disk_type = "SSD Standard"
    user_data = "foo"
    bus = "VIRTIO"
}
  nic {
    lan = ` + constant.LanResource + `.` + constant.LanTestResource + `.id
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

const (
	testAccCheckServerVCPUNoNic = `
resource ` + constant.DatacenterResource + ` ` + constant.DatacenterTestResource + ` {
  name       = "server-test"
  location   = "de/txl"
}
resource "ionoscloud_ipblock" "webserver_ipblock" {
  location = ` + constant.DatacenterResource + `.` + constant.DatacenterTestResource + `.location
  size = 4
  name = "webserver_ipblock"
}
resource ` + constant.LanResource + ` ` + constant.LanTestResource + ` {
  datacenter_id = ` + constant.DatacenterResource + `.` + constant.DatacenterTestResource + `.id
  public = true
  name = "public"
}
resource ` + constant.ServerVCPUResource + ` ` + constant.ServerTestResource + ` {
  name = "` + constant.ServerTestResource + `"
  hostname = "` + constant.ServerTestHostname + `"
  datacenter_id = ` + constant.DatacenterResource + `.` + constant.DatacenterTestResource + `.id
  cores = 1
  ram = 1024
  image_name ="ubuntu:latest"
  ssh_keys = ["` + sshKey + `"]
  volume {
    name = "system"
    size = 5
    disk_type = "SSD Standard"
    user_data = "foo"
    bus = "VIRTIO"
	expose_serial = true
  }
}`
	testAccCheckServerVCPUNoNicUpdate = `
resource ` + constant.DatacenterResource + ` ` + constant.DatacenterTestResource + ` {
	name       = "server-test"
	location   = "de/txl"
}
resource "ionoscloud_ipblock" "webserver_ipblock" {
  location = ` + constant.DatacenterResource + `.` + constant.DatacenterTestResource + `.location
  size = 4
  name = "webserver_ipblock"
}
resource ` + constant.LanResource + ` ` + constant.LanTestResource + ` {
  datacenter_id = ` + constant.DatacenterResource + `.` + constant.DatacenterTestResource + `.id
  public = true
  name = "public"
}
resource ` + constant.ServerVCPUResource + ` ` + constant.ServerTestResource + ` {
  name = "` + constant.ServerTestResource + `"
  hostname = "updated"
  datacenter_id = ` + constant.DatacenterResource + `.` + constant.DatacenterTestResource + `.id
  cores = 2
  ram = 2048
  image_name ="ubuntu:latest"
  ssh_keys = ["` + sshKey + `"]
  volume {
    name = "system"
    size = 5
    disk_type = "SSD Standard"
    user_data = "foo"
    bus = "VIRTIO"
	expose_serial = true
  }
}`
)

const testAccCheckServerVCPUShutDown = `
resource ` + constant.DatacenterResource + ` ` + constant.DatacenterTestResource + ` {
  name       = "server-test"
  location   = "de/txl"
}
resource "ionoscloud_ipblock" "webserver_ipblock" {
  location = ` + constant.DatacenterResource + `.` + constant.DatacenterTestResource + `.location
  size = 4
  name = "webserver_ipblock"
}
resource ` + constant.LanResource + ` ` + constant.LanTestResource + ` {
  datacenter_id = ` + constant.DatacenterResource + `.` + constant.DatacenterTestResource + `.id
  public = true
  name = "public"
}
resource ` + constant.ServerVCPUResource + ` ` + constant.ServerTestResource + ` {
  name = "` + constant.ServerTestResource + `"
  hostname = "` + constant.ServerTestHostname + `"
  datacenter_id = ` + constant.DatacenterResource + `.` + constant.DatacenterTestResource + `.id
  vm_state = "` + constant.VMStateStop + `"
  cores = 1
  ram = 1024
  image_name ="ubuntu:latest"
  ssh_keys = ["` + sshKey + `"]
  volume {
    name = "system"
    size = 5
    disk_type = "SSD Standard"
    user_data = "foo"
    bus = "VIRTIO"
  }
}`

const testAccCheckServerVCPUPowerOn = `
resource ` + constant.DatacenterResource + ` ` + constant.DatacenterTestResource + ` {
  name       = "server-test"
  location   = "de/txl"
}
resource "ionoscloud_ipblock" "webserver_ipblock" {
  location = ` + constant.DatacenterResource + `.` + constant.DatacenterTestResource + `.location
  size = 4
  name = "webserver_ipblock"
}
resource ` + constant.LanResource + ` ` + constant.LanTestResource + ` {
  datacenter_id = ` + constant.DatacenterResource + `.` + constant.DatacenterTestResource + `.id
  public = true
  name = "public"
}
resource ` + constant.ServerVCPUResource + ` ` + constant.ServerTestResource + ` {
  name = "` + constant.ServerTestResource + `"
  datacenter_id = ` + constant.DatacenterResource + `.` + constant.DatacenterTestResource + `.id
  vm_state = "` + constant.VMStateStart + `"
  cores = 1
  ram = 1024
  image_name ="ubuntu:latest"
  ssh_keys = ["` + sshKey + `"]
  volume {
    name = "system"
    size = 5
    disk_type = "SSD Standard"
    user_data = "foo"
    bus = "VIRTIO"
    expose_serial = true
  }
}`
const testAccCheckCDNDistributionConfigOnlyRequired = `resource ` + constant.CDNDistributionResource + ` ` + constant.CDNDistributionTestResource + ` {
	domain         = "unique.test.example.com"
	routing_rules {
		scheme = "http"
		prefix = "/api"
		upstream {
			host             = "server.example.com"
			caching          = true
			waf              = true
			sni_mode		 = "distribution"
			rate_limit_class = "R100"
		}
	}
}`

const testAccCheckCDNDistributionConfigBasic = `resource ` + constant.CDNDistributionResource + ` ` + constant.CDNDistributionTestResource + ` {
	domain         = "unique.test.example.com"
	routing_rules {
		scheme = "http"
		prefix = "/api"
		upstream {
			host             = "server.example.com"
			caching          = true
			waf              = true
			sni_mode		 = "distribution"
			rate_limit_class = "R100"
			geo_restrictions {
				allow_list = [ "RO"]
			}
		}
	}
}`

const testAccCheckCDNDistributionConfigBasicImport = `resource ` + constant.CDNDistributionResource + ` ` + constant.CDNDistributionTestResource + ` {
	domain         = "ionossdk.terra.example.import"
	routing_rules {
		scheme = "http"
		prefix = "/api"
		upstream {
			host             = "server.example.com"
			caching          = true
			waf              = true
			rate_limit_class = "R100"
			sni_mode		 = "distribution"
			geo_restrictions {
				allow_list = [ "RO"]
			}
		}
	}
}`

const resourceRandomUUID = `
resource "random_uuid" "uuid" {
}
`

const resourceRandomString = `
resource "random_string" "simple_string" {
	length = 16
	special = false
}
`

// K8s values
const (
	K8sVersion                  = "1.29.6"
	UpgradedK8sVersion          = "1.29.7"
	K8sBucket                   = "test_k8s_terraform_v7"
	K8sPrivateClusterNodeSubnet = "192.168.0.0/16"
)

// DNS test constants: configs, attributes and values.

// DNS Zones constants
const zoneNameAttribute = "name"
const zoneNameValue = "test.com"
const zoneDescriptionAttribute = "description"
const zoneDescriptionValue = "test description"
const zoneUpdatedDescriptionValue = "updated description"
const zoneEnabledAttribute = "enabled"
const zoneEnabledValue = "true"
const zoneupdatedEnabledValue = "false"

const DNSZoneConfig = `
resource ` + constant.DNSZoneResource + ` ` + constant.DNSZoneTestResourceName + ` {
	` + zoneNameAttribute + ` = "` + zoneNameValue + `"
	` + zoneDescriptionAttribute + ` = "` + zoneDescriptionValue + `"
    ` + zoneEnabledAttribute + ` = ` + zoneEnabledValue + `
}
`

// DNS Records constants
const recordNameAttribute = "name"
const recordNameValue = "example.com"
const recordTypeAttribute = "type"
const recordTypeValue = "MX"
const recordContentAttribute = "content"
const recordContentValue = "mail.example.com"
const recordUpdatedContentValue = "updated.example.com"
const recordTtlAttribute = "ttl"
const recordTtlValue = "2000"
const recordUpdatedTtlValue = "3600"
const recordPriorityAttribute = "priority"
const recordPriorityValue = "1024"
const recordUpdatedPriorityValue = "2048"
const recordEnabledAttribute = "enabled"
const recordEnabledValue = "true"
const recordUpdatedEnabledValue = "false"

const DNSRecordConfig = DNSZoneConfig + `
resource ` + constant.DNSRecordResource + ` ` + constant.DNSRecordTestResourceName + ` {
	zone_id = ` + constant.DNSZoneResource + `.` + constant.DNSZoneTestResourceName + `.id
	` + recordNameAttribute + ` = "` + recordNameValue + `"
	` + recordTypeAttribute + ` = "` + recordTypeValue + `"
	` + recordContentAttribute + ` = "` + recordContentValue + `"
	` + recordTtlAttribute + ` = ` + recordTtlValue + `
	` + recordPriorityAttribute + ` = ` + recordPriorityValue + `
	` + recordEnabledAttribute + ` = ` + recordEnabledValue + `
}
`

// Logging Pipeline constants
// Attributes
const nameAttribute = "name"
const pipelineGrafanaAddressAttribute = "grafana_address" // nolint:unused
const pipelineLogAttribute = "log"
const pipelineLogSourceAttribute = "source"
const pipelineLogTagAttribute = "tag"
const pipelineLogProtocolAttribute = "protocol"
const pipelineLogDestinationAttribute = "destinations"
const pipelineLogDestinationTypeAttribute = "type"
const pipelineLogDestinationRetentionAttribute = "retention_in_days"

// Values
const pipelineNameValue = "testpipeline"
const pipelineLogSourceValue = "kubernetes"
const pipelineLogTagValue = "testtag"
const pipelineLogProtocolValue = "http"
const pipelineLogDestinationTypeValue = "loki"
const pipelineLogDestinationRetentionValue = "7"
const pipelineLogDestination = pipelineLogDestinationAttribute + `{
	` + pipelineLogDestinationTypeAttribute + ` = "` + pipelineLogDestinationTypeValue + `"
	` + pipelineLogDestinationRetentionAttribute + ` = "` + pipelineLogDestinationRetentionValue + `"
}`
const pipelineLog = pipelineLogAttribute + `{
	` + pipelineLogSourceAttribute + ` = "` + pipelineLogSourceValue + `"
	` + pipelineLogTagAttribute + ` = "` + pipelineLogTagValue + `"
	` + pipelineLogProtocolAttribute + ` = "` + pipelineLogProtocolValue + `"
	` + pipelineLogDestination + `
}`

// Update values
const pipelineNameUpdatedValue = "updatedtestpipeline"
const pipelineLogSourceUpdatedValue = "docker"
const pipelineLogTagUpdatedValue = "updatedtesttag"
const pipelineLogProtocolUpdatedValue = "tcp"
const pipelineLogDestinationRetentionUpdatedValue = "14"
const pipelineLogDestinationUpdated = pipelineLogDestinationAttribute + `{
	` + pipelineLogDestinationTypeAttribute + ` = "` + pipelineLogDestinationTypeValue + `"
	` + pipelineLogDestinationRetentionAttribute + ` = "` + pipelineLogDestinationRetentionUpdatedValue + `"
}`
const pipelineLogUpdated = pipelineLogAttribute + `{
	` + pipelineLogSourceAttribute + ` = "` + pipelineLogSourceUpdatedValue + `"
	` + pipelineLogTagAttribute + ` = "` + pipelineLogTagUpdatedValue + `"
	` + pipelineLogProtocolAttribute + ` = "` + pipelineLogProtocolUpdatedValue + `"
	` + pipelineLogDestinationUpdated + `
}`

// LoggingPipelineConfig Standard configuration
const LoggingPipelineConfig = `
resource ` + constant.LoggingPipelineResource + ` ` + constant.LoggingPipelineTestResourceName + ` {
	` + nameAttribute + ` = "` + pipelineNameValue + `"
	` + pipelineLog + `
	location = "es/vit"
}
`

// DBaaS constants
// Attributes
// These can be used for all clusters, no matter the type
const clusterIdAttribute = "cluster_id"
const clusterInstancesAttribute = "instances"
const clusterCoresAttribute = "cores"
const clusterLocationAttribute = "location"
const clusterRamAttribute = "ram"
const clusterStorageSizeAttribute = "storage_size"
const clusterConnectionsAttribute = "connections"
const clusterConnectionsDatacenterIDAttribute = "datacenter_id"
const clusterConnectionsLanIDAttribute = "lan_id"
const clusterConnectionsCidrAttribute = "cidr"
const clusterDisplayNameAttribute = "display_name"
const clusterMaintenanceWindowAttribute = "maintenance_window"
const clusterMaintenanceWindowDayOfTheWeekAttribute = "day_of_the_week"
const clusterBackupAttribute = "backup"
const clusterBackupLocationAttribute = "location"

const clusterMaintenanceWindowTimeAttribute = "time"
const clusterCredentialsAttribute = "credentials"
const clusterCredentialsUsernameAttribute = "username"
const clusterCredentialsPasswordAttribute = "password"

const datacenterResourceName = "datacenter_example"
const lanResourceName = "lan_example"
const mongoVersion = "7.0"
const mongoVersionOld = "6.0"

// Values
const clusterMaintenanceWindowDayOfTheWeekValue = "Sunday"
const clusterMaintenanceWindowDayOfTheWeekUpdateValue = "Monday"
const clusterMaintenanceWindowTimeValue = "09:00:00"
const clusterMaintenanceWindowTimeUpdateValue = "10:00:00"

const clusterBackupLocationValue = "de"

const wireguardGatewayConfig = `
resource "ionoscloud_datacenter" "datacenter_example" {
  name = "datacenter_example"
  location = "de/fra"
}
resource ` + constant.IpBlockResource + ` ` + constant.IpBlockTestResource + ` {
  location = "de/fra"
  size = 1
  name = "` + constant.IpBlockTestResource + `"
}

resource "ionoscloud_lan" "lan_example" {
  name = "lan_example"
  datacenter_id = ionoscloud_datacenter.datacenter_example.id
}

resource` + ` ` + constant.WireGuardGatewayResource + ` ` + constant.WireGuardGatewayTestResource + `{
  name = "` + constant.WireGuardGatewayTestResource + `"
  location = "de/fra"
  description = "description"
  private_key = "0HpE4BNwGHabeaC4aY/GFxB6fBSc0d49Db0qAzRVSVc="

  gateway_ip = ` + constant.IpBlockResource + `.` + constant.IpBlockTestResource + `.ips[0]
  interface_ipv4_cidr =  "192.168.1.100/24"
  connections   {
    datacenter_id   =  ionoscloud_datacenter.datacenter_example.id
    lan_id          =  ionoscloud_lan.lan_example.id
    ipv4_cidr       =  "192.168.1.108/24"
  }
  maintenance_window {
    day_of_the_week       = "Monday"
    time                  = "09:00:00"
  }
  tier = "STANDARD"
}` // nolint:unused

// This configuration is used because there are some problems with the API and the creation/deletion
// of the setup resources (datacenter, lan, server) is not possible (there are some problems with
// LAN deletion). Because of that, for the moment, only to test the InMemoryDB and NFS functionalities, we
// will use data sources for already existing setup resources.

const temporaryConfigSetup = `
data "ionoscloud_datacenter" "datacenterDS" {
	id = "88eeae0d-515d-44c1-b142-d9293c20e676"
}

data "ionoscloud_lan" "lanDS" {
	id = "1"
	datacenter_id = data.ionoscloud_datacenter.datacenterDS.id
}

data "ionoscloud_server" "serverDS" {
	id = "1f77a37e-2b38-49f2-b2e1-61a47ccf5f15"
	datacenter_id = data.ionoscloud_datacenter.datacenterDS.id
}

locals {
 prefix                   = format("%s/%s", data.ionoscloud_server.serverDS.nics[0].ips[0], "24")
 database_ip              = cidrhost(local.prefix, 1)
 database_ip_cidr         = format("%s/%s", local.database_ip, "24")
}

resource ` + constant.RandomPassword + ` "replicaset_password" {
  length           = 16
  special          = true
  override_special = "!#$%&*()-_=+[]{}<>:?"
}
`
