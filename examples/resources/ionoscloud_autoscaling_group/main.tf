resource "ionoscloud_datacenter" "datacenter_example" {
    name     = "datacenter_example"
    location = "de/fra"
}

resource "ionoscloud_lan" "lan_example_1" {
    datacenter_id    = ionoscloud_datacenter.datacenter_example.id
    public           = false
    name             = "lan_example_1"
}

resource "ionoscloud_lan" "lan_example_2" {
    datacenter_id    = ionoscloud_datacenter.datacenter_example.id
    public           = false
    name             = "lan_example_2"
}

resource "ionoscloud_target_group" "autoscaling_target_group" {
  name                      = "Target Group Example"
  algorithm                 = "ROUND_ROBIN"
  protocol                  = "HTTP"
  protocol_version          = "HTTP1"
}

resource "ionoscloud_autoscaling_group" "autoscaling_group_example" {
  datacenter_id = ionoscloud_datacenter.datacenter_example.id
  max_replica_count      = 2
  min_replica_count      = 1
  name                   = "autoscaling_group_example"
  policy {
    metric             = "INSTANCE_CPU_UTILIZATION_AVERAGE"
    range              = "PT24H"
    scale_in_action {
      amount                  =  1
      amount_type             = "ABSOLUTE"
      termination_policy_type = "OLDEST_SERVER_FIRST"
      cooldown_period         = "PT5M"
      delete_volumes          = true
    }
    scale_in_threshold = 33
    scale_out_action  {
      amount          =  1
      amount_type     = "ABSOLUTE"
      cooldown_period = "PT5M"
    }
    scale_out_threshold = 77
    unit                = "PER_HOUR"
  }
  replica_configuration {
    availability_zone = "AUTO"
    cores               = "2"
    cpu_family           = "INTEL_SKYLAKE"
    ram                  = 2048
    nic {
      lan   = ionoscloud_lan.lan_example_1.id
      name  = "nic_example_1"
      dhcp  = true
    }
    nic {
      lan   = ionoscloud_lan.lan_example_2.id
      name  = "nic_example_2"
      dhcp  = true
      firewall_active = true
      firewall_type = "INGRESS"
      firewall_rule {
        name = "rule_1"
        protocol = "TCP"
        port_range_start = 1
        port_range_end = 1000
        type = "INGRESS"
      }

      flow_log {
        name="flow_log_1"
        bucket="test-de-bucket"
        action="ALL"
        direction="BIDIRECTIONAL"
      }

      target_group {
        target_group_id = ionoscloud_target_group.autoscaling_target_group.id
        port            = 80
        weight          = 50
      }
    }
    volume    {
      image_alias    = "ubuntu:latest"
      name           = "volume_example"
      size           = 10
      type           = "HDD"
      user_data      = "ZWNobyAiSGVsbG8sIFdvcmxkIgo="
      image_password = random_password.server_image_password.result
      boot_order     = "AUTO"
    }
  }
}

resource "random_password" "server_image_password" {
  length           = 16
  special          = false
}
