/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package proxmoxtf

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	dvResourceVirtualEnvironmentPersistentDiskName   = "persistent-1"
	dvResourceVirtualEnvironmentPersistentDiskSize   = "10G"
	dvResourceVirtualEnvironmentPersistentDiskFormat = "raw"

	mkResourceVirtualEnvironmentPersistentDiskNodeName = "node_name"
	mkResourceVirtualEnvironmentPersistentDiskStorage  = "storage_pool"
	mkResourceVirtualEnvironmentPersistentDiskVmID     = "vm_id"
	mkResourceVirtualEnvironmentPersistentDiskName     = "name"
	mkResourceVirtualEnvironmentPersistentDiskSize     = "size"
	mkResourceVirtualEnvironmentPersistentDiskFormat   = "format"
)

func resourceVirtualEnvironmenPersistentDisk() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			mkResourceVirtualEnvironmentPersistentDiskNodeName: {
				Type:        schema.TypeString,
				Description: "The node name",
				Required:    true,
				ForceNew:    true,
			},
			mkResourceVirtualEnvironmentPersistentDiskStorage: {
				Type:        schema.TypeString,
				Description: "Storage pool",
				Required:    true,
				ForceNew:    true,
			},
			mkResourceVirtualEnvironmentPersistentDiskVmID: {
				Type:        schema.TypeInt,
				Description: "Id of a non-existing VM which will be the owner. IMPORTANT: Make sure a VM with this Id never exists, as it will delete the disk when the vm is removed.",
				Required:    true,
				ForceNew:    true,
			},
			mkResourceVirtualEnvironmentPersistentDiskName: {
				Type:        schema.TypeString,
				Description: "Name of disk",
				Optional:    true,
				Default:     dvResourceVirtualEnvironmentPersistentDiskName,
				ForceNew:    true,
			},
			mkResourceVirtualEnvironmentPersistentDiskSize: {
				Type:        schema.TypeString,
				Description: "Size of disk. Default unit is kilobytes, supports suffixes M (1024K) and G (1024M)",
				Optional:    true,
				Default:     dvResourceVirtualEnvironmentPersistentDiskSize,
				ForceNew:    true,
			},
			mkResourceVirtualEnvironmentPersistentDiskFormat: {
				Type:        schema.TypeString,
				Description: "Disk format (raw | qcow2 | subvol)",
				Optional:    true,
				Default:     dvResourceVirtualEnvironmentPersistentDiskFormat,
				ForceNew:    true,
			},
		},
		CreateContext: resourceVirtualEnvironmentPersistentDiskCreate,
		ReadContext:   resourceVirtualEnvironmentPersistentDiskRead,
		//UpdateContext: resourceVirtualEnvironmentPersistentDiskUpdate,
		DeleteContext: resourceVirtualEnvironmentPersistentDiskDelete,
	}
}

func resourceVirtualEnvironmentPersistentDiskCreate(
	ctx context.Context,
	d *schema.ResourceData,
	m interface{},
) diag.Diagnostics {
	config := m.(providerConfiguration)
	veClient, err := config.GetVEClient()
	if err != nil {
		return diag.FromErr(err)
	}

	nodeName := d.Get(mkResourceVirtualEnvironmentPersistentDiskNodeName).(string)
	storage := d.Get(mkResourceVirtualEnvironmentPersistentDiskStorage).(string)
	vmID := d.Get(mkResourceVirtualEnvironmentPersistentDiskVmID).(int)
	diskName := d.Get(mkResourceVirtualEnvironmentPersistentDiskName).(string)
	diskSize := d.Get(mkResourceVirtualEnvironmentPersistentDiskSize).(string)
	diskFormat := d.Get(mkResourceVirtualEnvironmentPersistentDiskFormat).(string)

	var commands []string
	commands = append(
		commands,
		`set -e`,
		fmt.Sprintf(`storage="%s"`, storage),
		fmt.Sprintf(`vm_id="%d"`, vmID),
		fmt.Sprintf(`disk_name="%s"`, diskName),
		fmt.Sprintf(`disk_size="%s"`, diskSize),
		fmt.Sprintf(`disk_format="%s"`, diskFormat),
		`pvesm alloc ${storage} ${vm_id} vm-${vm_id}-${disk_name} ${disk_size} --format ${disk_format}`,
	)

	err = veClient.ExecuteNodeCommands(ctx, nodeName, commands)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(fmt.Sprintf(`%s_%s_%d-%s`, nodeName, storage, vmID, diskName))

	return nil
}

func resourceVirtualEnvironmentPersistentDiskRead(
	ctx context.Context,
	d *schema.ResourceData,
	m interface{},
) diag.Diagnostics {
	return nil
}

// func resourceVirtualEnvironmentPersistentDiskUpdate(
// 	ctx context.Context,
// 	d *schema.ResourceData,
// 	m interface{},
// ) diag.Diagnostics {
// 	return nil
// }

func resourceVirtualEnvironmentPersistentDiskDelete(
	ctx context.Context,
	d *schema.ResourceData,
	m interface{},
) diag.Diagnostics {
	config := m.(providerConfiguration)
	veClient, err := config.GetVEClient()
	if err != nil {
		return diag.FromErr(err)
	}

	nodeName := d.Get(mkResourceVirtualEnvironmentPersistentDiskNodeName).(string)
	storage := d.Get(mkResourceVirtualEnvironmentPersistentDiskStorage).(string)
	vmID := d.Get(mkResourceVirtualEnvironmentPersistentDiskVmID).(int)
	diskName := d.Get(mkResourceVirtualEnvironmentPersistentDiskName).(string)

	var commands []string

	commands = append(
		commands,
		`set -e`,
		fmt.Sprintf(`storage="%s"`, storage),
		fmt.Sprintf(`vm_id="%d"`, vmID),
		fmt.Sprintf(`disk_name="%s"`, diskName),
		`pvesm free vm-${vm_id}-${disk_name} --storage ${storage}`,
	)

	err = veClient.ExecuteNodeCommands(ctx, nodeName, commands)

	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")

	return nil
}
