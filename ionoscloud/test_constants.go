package ionoscloud

const (
	testAccCheckBackupUnitConfigBasic = `
resource ` + BackupUnitResource + ` ` + BackupUnitTestResource + ` {
	name        = "` + BackupUnitTestResource + `"
	password    = ` + RandomPassword + `.backup_unit_password.result
	email       = "example@ionoscloud.com"
}

resource ` + RandomPassword + ` "backup_unit_password" {
  length           = 16
  special          = true
  override_special = "!#$%&*()-_=+[]{}<>:?"
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
  image_password = ` + RandomPassword + `.server_image_password.result
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
}
` + ServerImagePassword

// Solves  #372 crash when ips field in nic resource is a list with an empty string
const testAccCheckServerConfigEmptyNicIps = `
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
  image_password = ` + RandomPassword + `.server_image_password.result
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
    ips            = [""]
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

const testAccDataSourceDatacenterWrongNameError = testAccCheckDatacenterConfigBasic + `
data ` + DatacenterResource + ` ` + DatacenterDataSourceMatching + ` {
    name = "wrong_name"
    location = ` + DatacenterResource + `.` + DatacenterTestResource + `.location
}`

const ImmutableError = "attribute is immutable, therefore not allowed in update requests"

const ServerImagePassword = `
resource ` + RandomPassword + ` "server_image_password" {
  length           = 16
  special          = false
}
`

const ServerImagePasswordUpdated = `
resource ` + RandomPassword + ` "server_image_password_updated" {
  length           = 16
  special          = false
}
`

// Cube Server Constants
const testAccCheckCubeServerConfigBasic = `
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
  image_name ="ubuntu:latest"
  image_password = ` + RandomPassword + `.server_image_password.result
  
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
}
` + ServerImagePassword

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

const K8sVersion = "1.23.12"
const UpgradedK8sVersion = "1.24.6"

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
resource ` + DNSZoneResource + ` ` + DNSZoneTestResourceName + ` {
	` + zoneNameAttribute + ` = "` + zoneNameValue + `"
	` + zoneDescriptionAttribute + ` = "` + zoneDescriptionValue + `"
    ` + zoneEnabledAttribute + ` = ` + zoneEnabledValue + `
}
`

// DNS Records constants
const recordNameAttribute = "name"
const recordNameValue = "testrecord"
const recordTypeAttribute = "type"
const recordTypeValue = "CNAME"
const recordContentAttribute = "content"
const recordContentValue = "1.2.3.4"
const recordUpdatedContentValue = "4.3.2.1"
const recordTtlAttribute = "ttl"
const recordTtlValue = "2000"
const recordUpdatedTtlValue = "3600"
const recordPriorityAttribute = "priority"
const recordPriorityValue = "1024"
const recordEnabledAttribute = "enabled"
const recordEnabledValue = "true"
const recordUpdatedEnabledValue = "false"

const DNSRecordConfig = DNSZoneConfig + `
resource ` + DNSRecordResource + ` ` + DNSRecordTestResourceName + ` {
	zone_id = ` + DNSZoneResource + `.` + DNSZoneTestResourceName + `.id
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
const pipelineNameAttribute = "name"
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

// Standard configuration
const LoggingPipelineConfig = `
resource ` + LoggingPipelineResource + ` ` + LoggingPipelineTestResourceName + ` {
	` + pipelineNameAttribute + ` = "` + pipelineNameValue + `"
	` + pipelineLog + `
}
`

// DBaaS PgSQL constants
// Attributes
const clusterIdAttribute = "cluster_id"
const usernameAttribute = "username"
const passwordAttribute = "password"
const isSystemUserAttribute = "is_system_user"

// Values
const usernameValue = "testusername"
const isSystemUserValue = "false"

// Configurations

const PgSqlUserConfig = `
resource ` + DatacenterResource + ` "datacenter_example" {
  name        = "datacenter_example"
  location    = "de/txl"
  description = "Datacenter for testing DBaaS PgSql user"
}

resource ` + LanResource + ` "lan_example" {
  datacenter_id = ` + DatacenterResource + `.datacenter_example.id 
  public        = false
  name          = "lan_example"
}

resource ` + PsqlClusterResource + ` ` + DBaaSClusterTestResource + ` {
  postgres_version   = 12
  instances          = 1
  cores              = 1
  ram                = 2048
  storage_size       = 2048
  storage_type       = "HDD"
  connections   {
	datacenter_id   =  ` + DatacenterResource + `.datacenter_example.id 
    lan_id          =  ` + LanResource + `.lan_example.id 
    cidr            =  "192.168.1.100/24"
  }
  location = ` + DatacenterResource + `.datacenter_example.location
  backup_location = "de"
  display_name = "` + DBaaSClusterTestResource + `"
  maintenance_window {
    day_of_the_week  = "Sunday"
    time             = "09:00:00"
  }
  credentials {
  	username = "username"
	password = ` + RandomPassword + `.cluster_password.result
  }
  synchronization_mode = "ASYNCHRONOUS"
}

resource ` + PsqlUserResource + ` ` + UserTestResource + ` {
  ` + clusterIdAttribute + ` = ` + PsqlClusterResource + `.` + DBaaSClusterTestResource + `.id 
  ` + usernameAttribute + ` = "` + usernameValue + `"
  ` + passwordAttribute + ` = ` + RandomPassword + `.user_password.result
  ` + isSystemUserAttribute + ` = ` + isSystemUserValue + `
}

resource ` + RandomPassword + ` "cluster_password" {
  length           = 16
  special          = true
  override_special = "!#$%&*()-_=+[]{}<>:?"
}

resource ` + RandomPassword + ` "user_password" {
  length           = 16
  special          = true
  override_special = "!#$%&*()-_=+[]{}<>:?"
}
`
