---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "synology_vmm_guest Resource - terraform-provider-synology"
subcategory: ""
description: |-
  
---

# synology_vmm_guest (Resource)



## Example Usage

```terraform
terraform {
  required_providers {
    synology = {
      version = "0.1"
      source  = "github.com/arnouthoebreckx/synology"
    }
  }
}

provider "synology" {
  url      = "<SYNOLOGY_ADDRESS>"
  username = "<SYNOLOGY_USERNAME>"
  password = "<SYNOLOGY_PASSWORD>"
  # these variables can be set as env vars in SYNOLOGY_ADDRESS SYNOLOGY_USERNAME and SYNOLOGY_PASSWORD
}

resource "synology_vmm_guest" "my-guest" {
  autorun     = 2
  poweron      = true
  guest_name   = "terraform-guest"
  description  = "Virtual machine setup with terraform"
  storage_name = "synology - VM Storage 1"
  vram_size    = 1024
  vnics {
    network_name = "default"
  }
  vdisks {
    create_type = 0
    vdisk_size  = 10240
  }
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `guest_name` (String) The guest name
- `vdisks` (Block List, Min: 1) (see [below for nested schema](#nestedblock--vdisks))
- `vnics` (Block List, Min: 1) (see [below for nested schema](#nestedblock--vnics))

### Optional

- `autorun` (Number) Optional. 0: off 1: last state 2: on
- `description` (String) Optional. The description of the guest.
- `guest_id` (String) The guest name
- `poweron` (Boolean) Optional. Default VM is not powered on.
- `storage_id` (String) Optional. The id of storage where the guest resides. Note: At least storage_id or storage_name should be given.
- `storage_name` (String) Optional. The name of storage where the guest resides. Note: At least storage_id or storage_name should be given.
- `vcpu_num` (Number) Optional. The vCPU number
- `vram_size` (Number) Optional. The memory size in MB

### Read-Only

- `id` (String) The ID of this resource.
- `status` (String) The guest status. (running/shutdown/inaccessiblen/booting/shutting_down/moving/stor_migrating/creating/importing/preparing/ha_standby/unknown/crashed/undefined

<a id="nestedblock--vdisks"></a>
### Nested Schema for `vdisks`

Required:

- `create_type` (Number) 0: Create an empty vDisk, 1: Clone an existing image

Optional:

- `image_id` (String) Optional. If create_type is 1, at least image_id or image_name should be given. The id of the image that is to be cloned. Note: Image type should be disk.
- `image_name` (String) Optional. If create_type is 1, at least image_id or image_name should be given. The name of the image that is to be cloned. Note: Image type should be disk.
- `vdisk_size` (Number) Optional. If create_type is 0, this field must be set. The created vDisk size in MB.

Read-Only:

- `controller` (Number) 1: VirtIO 2: IDE 3: SATA
- `unmap` (Boolean) Determine whether to enable space reclamation.
- `vdisk_id` (String) The id of this vDisk.


<a id="nestedblock--vnics"></a>
### Nested Schema for `vnics`

Optional:

- `mac` (String) Optional. MAC address. If not specified, a MAC address of this vNIC will be randomly generated.
- `network_id` (String) Optional. Connected network group id. At least network_id or network_name should be given. Note: network_id can be an empty string to represent not being connected.
- `network_name` (String) Optional. Connected network group name. At least network_id or network_name should be given.

Read-Only:

- `model` (Number) 1: VirtIO 2: e1000 3: rtl8139
- `vnic_id` (String) The id of this vNIC.

