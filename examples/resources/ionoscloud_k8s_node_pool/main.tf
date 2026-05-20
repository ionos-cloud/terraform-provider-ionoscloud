resource "ionoscloud_datacenter" "example" {
  name                  = "Datacenter Example"
  location              = "us/las"
  description           = "datacenter description"
  sec_auth_protection   = false
}

resource "ionoscloud_lan" "example" {
  datacenter_id         = ionoscloud_datacenter.example.id
  public                = false
  name                  = "Lan Example"
  lifecycle {
    create_before_destroy = true
  }
}

resource "ionoscloud_ipblock" "example" {
  location              = "us/las"
  size                  = 3
  name                  = "IP Block Example"
}

resource "ionoscloud_k8s_cluster" "example" {
  name                  = "k8sClusterExample"
  k8s_version           = "1.31.2"
  maintenance_window {
    day_of_the_week     = "Sunday"
    time                = "09:00:00Z"
  }
  api_subnet_allow_list = ["1.2.3.4/32"]
  s3_buckets {
     name               = "globally_unique_s3_bucket_name"
  }
}

resource "ionoscloud_k8s_node_pool" "example" {
  datacenter_id         = ionoscloud_datacenter.example.id
  k8s_cluster_id        = ionoscloud_k8s_cluster.example.id
  name                  = "k8sNodePoolExample"
  k8s_version           = ionoscloud_k8s_cluster.example.k8s_version
  maintenance_window {
    day_of_the_week     = "Monday"
    time                = "09:00:00Z"
  }
  auto_scaling {
    min_node_count      = 1
    max_node_count      = 2
  }
  cpu_family            = "INTEL_XEON"
  availability_zone     = "AUTO"
  storage_type          = "SSD"
  node_count            = 1
  cores_count           = 2
  ram_size              = 2048
  storage_size          = 40
  server_type           = "DedicatedCore"
  public_ips            = [ ionoscloud_ipblock.example.ips[0], ionoscloud_ipblock.example.ips[1], ionoscloud_ipblock.example.ips[2] ]
  lans {
    id                  = ionoscloud_lan.example.id
    dhcp                = true
	routes {
       network          = "1.2.3.5/24"
       gateway_ip       = "10.1.5.17"
     }
   }
  labels                = {
    lab1                = "value1"
    lab2                = "value2"
  }
  annotations           = {
    ann1                = "value1"
    ann2                = "value2"
  }
}
