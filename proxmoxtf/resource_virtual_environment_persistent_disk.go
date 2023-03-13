/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package proxmoxtf

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/zarichard/terraform-provider-proxmox/proxmox"
)

const (
	dvResourceVirtualEnvironmentPersistentDiskName   = "persistent-1"
	dvResourceVirtualEnvironmentPersistentDiskSize   = "10"
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
				Type:        schema.TypeInt,
				Description: "The disk size in gigabytes",
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
	diskSize := d.Get(mkResourceVirtualEnvironmentPersistentDiskSize).(int)
	diskFormat := d.Get(mkResourceVirtualEnvironmentPersistentDiskFormat).(string)

	var commands []string
	commands = append(
		commands,
		`set -e`,
		fmt.Sprintf(`storage="%s"`, storage),
		fmt.Sprintf(`vm_id="%d"`, vmID),
		fmt.Sprintf(`disk_name="%s"`, diskName),
		fmt.Sprintf(`disk_size="%d"`, diskSize),
		fmt.Sprintf(`disk_format="%s"`, diskFormat),
		`pvesm alloc ${storage} ${vm_id} vm-${vm_id}-${disk_name} ${disk_size}G --format ${disk_format}`,
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
	config := m.(providerConfiguration)
	veClient, err := config.GetVEClient()
	if err != nil {
		return diag.FromErr(err)
	}

	closeOrLogError := proxmox.CloseOrLogError(ctx)

	nodeName := d.Get(mkResourceVirtualEnvironmentPersistentDiskNodeName).(string)
	storage := d.Get(mkResourceVirtualEnvironmentPersistentDiskStorage).(string)
	vmID := d.Get(mkResourceVirtualEnvironmentPersistentDiskVmID).(int)
	diskName := fmt.Sprintf(`vm-%d-%s`, vmID, d.Get(mkResourceVirtualEnvironmentPersistentDiskName).(string))

	sshClient, err := veClient.OpenNodeShell(ctx, nodeName)
	if err != nil {
		return diag.FromErr(err)
	}
	defer closeOrLogError(sshClient)

	sshSession, err := sshClient.NewSession()
	if err != nil {
		return diag.FromErr(err)
	}
	defer closeOrLogError(sshSession)

	var commands []string
	commands = append(
		commands,
		`set -e`,
		fmt.Sprintf(`storage="%s"`, storage),
		fmt.Sprintf(`vm_id="%d"`, vmID),
		`pvesm list ${storage} --vmid ${vm_id} --content images`,
	)

	script := strings.Join(commands, " && \\\n")
	output, err := sshSession.CombinedOutput(
		fmt.Sprintf(
			"/bin/bash -c '%s'",
			strings.ReplaceAll(script, "'", "'\"'\"'"),
		),
	)

	if err != nil {
		return diag.FromErr(errors.New(string(output)))
	}

	outputString := string(output[:])
	lines := strings.Split(outputString, "\n")

	// ignore first line as it's the table headers
	if len(lines) > 1 {
		for _, line := range lines[1:] {
			outputName, outputFormat, outputSize := func() (name, format, size string) {

				values := strings.Fields(line)
				if len(values) >= 4 {
					format = values[1]
					size = values[3]

					storageAndName := strings.Split(values[0], ":")
					if len(storageAndName) >= 2 {
						name = storageAndName[1]
					}
				}

				return
			}()

			if outputName == diskName {
				outputSizeInt, err := strconv.Atoi(outputSize)
				if err != nil {
					return diag.FromErr(err)
				}

				sizeGB := outputSizeInt / 1024 / 1024 / 1024

				d.Set(mkResourceVirtualEnvironmentPersistentDiskSize, sizeGB)
				d.Set(mkResourceVirtualEnvironmentPersistentDiskFormat, outputFormat)

				return nil
			}
		}
	}

	d.SetId("")

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
