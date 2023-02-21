# Terraform Provider for Proxmox

Fork of https://github.com/bpg/terraform-provider-proxmox adding very specific functionality for my use-case:

# Persistent virtual disk resource
Hacky way to preserve VM disks when the VM is destroyed, without requiring twice the storage.
Meant to be used for persistent storage on nodes that should be kept during a reprovisioning.
Only tested on zfs based storage and is not a very generic solution.

# Passthrough physical disk
A field in VM resources to allow for passthrough disks to be specified.

