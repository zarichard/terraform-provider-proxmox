# Terraform Provider for Proxmox

Fork of https://github.com/zarichard/terraform-provider-proxmox adding very specific functionality for my use-case:

# Persistent virtual disk resource
Hacky way to preserve VM disks when the VM is destroyed, without requiring twice the storage.
Meant to be used for persistent storage on nodes that should be kept during a reprovisioning.
Only tested on zfs based storage and is not a very generic solution.

# Passthrough physical disk
A field in VM resources to allow for passthrough disks to be specified.

## Compatibility Matrix

| Proxmox version | Provider version |
| --------------- | ---------------- |
| 6.x             | \<= 0.4.4        |
| 7.x             | \>= 0.4.5        |

## Requirements

- [Terraform](https://www.terraform.io/downloads.html) 1.2+
- [Go](https://golang.org/doc/install) 1.19+ (to build the provider plugin)
- [GoReleaser](https://goreleaser.com/install/) v1.15+ (to build the provider plugin)

## Table of Contents

- [Building the provider](#building-the-provider)
- [Using the provider](#using-the-provider)
- [Testing the provider](#testing-the-provider)
- [Known issues](#known-issues)

## Building the provider

- Clone the repository to `$GOPATH/src/github.com/zarichard/terraform-provider-proxmox`:

  ```sh
  mkdir -p "${GOPATH}/src/github.com/bpg"
  cd "${GOPATH}/src/github.com/bpg"
  git clone git@github.com:bpg/terraform-provider-proxmox
  ```

- Enter the provider directory and build it:

  ```sh
  cd "${GOPATH}/src/github.com/zarichard/terraform-provider-proxmox"
  make build
  ```

## Using the provider

You can find the latest release and its documentation in the [Terraform Registry](https://registry.terraform.io/providers/bpg/proxmox/latest).

## Testing the provider

In order to test the provider, you can simply run `make test`.

```sh
make test
```

Tests are limited to regression tests, ensuring backwards compatibility.

## Deploying the example resources

There are number of TF examples in the `examples` directory, which can be used
to deploy a Container, VM, or other Proxmox resources on your test Proxmox cluster.
The following assumptions are made about the test Proxmox cluster:

- It has one node named `pve`
- The node has local storages named `local` and `local-lvm`

Create `examples/terraform.tfvars` with the following variables:

```sh
virtual_environment_username = "root@pam"
virtual_environment_password = "put-your-password-here"
virtual_environment_endpoint = "https://<your-cluster-endpoint>:8006/"
```

Then run `make example` to deploy the example resources.

## Known issues

### Disk images cannot be imported by non-PAM accounts

Due to limitations in the Proxmox VE API, certain actions need to be performed
using SSH. This requires the use of a PAM account (standard Linux account).

### Disk images from VMware cannot be uploaded or imported

Proxmox VE is not currently supporting VMware disk images directly. However, you
can still use them as disk images by using this workaround:

```hcl
resource "proxmox_virtual_environment_file" "vmdk_disk_image" {
  content_type = "iso"
  datastore_id = "datastore-id"
  node_name    = "node-name"

  source_file {
    # We must override the file extension to bypass the validation code
    # in the Proxmox VE API.
    file_name = "vmdk-file-name.img"
    path      = "path-to-vmdk-file"
  }
}

resource "proxmox_virtual_environment_vm" "example" {
  //...

  disk {
    datastore_id = "datastore-id"
    # We must tell the provider that the file format is vmdk instead of qcow2.
    file_format  = "vmdk"
    file_id      = "${proxmox_virtual_environment_file.vmdk_disk_image.id}"
  }

  //...
}
```

### Snippets cannot be uploaded by non-PAM accounts

Due to limitations in the Proxmox VE API, certain files need to be uploaded
using SFTP. This requires the use of a PAM account (standard Linux account).
