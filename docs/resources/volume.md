---
layout: "ionoscloud"
page_title: "IonosCloud: server"
sidebar_current: "docs-resource-volume"
description: |-
  Creates and manages IonosCloud Volume objects.
---

# ionoscloud\_volume

Manages a volume on IonosCloud.

## Example Usage

A primary volume will be created with the server. If there is a need for additional volumes, this resource handles it.

```hcl
resource "ionoscloud_volume" "example" {
  datacenter_id = "${ionoscloud_datacenter.example.id}"
  server_id     = "${ionoscloud_server.example.id}"
  image_name    = "${var.ubuntu}"
  size          = 5
  disk_type     = "HDD"
  ssh_key_path  = "${var.private_key_path}"
  bus           = "VIRTIO"
}
```

## Argument reference

* `datacenter_id` - (Required)[string] The ID of a Virtual Data Center.
* `server_id` - (Required)[string] The ID of a server.
* `disk_type` - (Required)[string] The volume type: HDD or SSD.
* `bus` - (Required)[Boolean] The bus type of the volume: VIRTIO or IDE.
* `size` -  (Required)[integer] The size of the volume in GB.
* `ssh_key_path` -  (Required)[list] List of paths to files containing a public SSH key that will be injected into IonosCloud provided Linux images. Required for IonosCloud Linux images. Required if `image_password` is not provided.
* `sshkey` - (Computed) The associated public SSH key.
* `image_password` - [string] Required if `sshkey_path` is not provided.
* `image_name` - [string] The image or snapshot UUID. May also be an image alias. It is required if `licence_type` is not provided.
* `licence_type` - [string] Required if `image_name` is not provided.
* `name` - (Optional)[string] The name of the volume.
* `availability_zone` - (Optional)[string] The storage availability zone assigned to the volume: AUTO, ZONE_1, ZONE_2, or ZONE_3.
