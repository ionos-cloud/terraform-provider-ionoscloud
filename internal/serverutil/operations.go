package serverutil

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	ionoscloud "github.com/ionos-cloud/sdk-go/v6"
	
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/bundleclient"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils"
)

func SetServerVolumeProperties(volume ionoscloud.Volume) map[string]interface{} {
	volumeMap := map[string]interface{}{}
	if volume.Properties != nil {
		utils.SetPropWithNilCheck(volumeMap, "name", volume.Properties.Name)
		utils.SetPropWithNilCheck(volumeMap, "disk_type", volume.Properties.Type)
		utils.SetPropWithNilCheck(volumeMap, "pci_slot", volume.Properties.PciSlot)
		utils.SetPropWithNilCheck(volumeMap, "licence_type", volume.Properties.LicenceType)
		utils.SetPropWithNilCheck(volumeMap, "bus", volume.Properties.Bus)
		utils.SetPropWithNilCheck(volumeMap, "availability_zone", volume.Properties.AvailabilityZone)
		utils.SetPropWithNilCheck(volumeMap, "cpu_hot_plug", volume.Properties.CpuHotPlug)
		utils.SetPropWithNilCheck(volumeMap, "ram_hot_plug", volume.Properties.RamHotPlug)
		utils.SetPropWithNilCheck(volumeMap, "nic_hot_plug", volume.Properties.NicHotPlug)
		utils.SetPropWithNilCheck(volumeMap, "nic_hot_unplug", volume.Properties.NicHotUnplug)
		utils.SetPropWithNilCheck(volumeMap, "disc_virtio_hot_plug", volume.Properties.DiscVirtioHotPlug)
		utils.SetPropWithNilCheck(volumeMap, "disc_virtio_hot_unplug", volume.Properties.DiscVirtioHotUnplug)
		utils.SetPropWithNilCheck(volumeMap, "device_number", volume.Properties.DeviceNumber)
		utils.SetPropWithNilCheck(volumeMap, "user_data", volume.Properties.UserData)
		utils.SetPropWithNilCheck(volumeMap, "backup_unit_id", volume.Properties.BackupunitId)
		utils.SetPropWithNilCheck(volumeMap, "boot_server", volume.Properties.BootServer)
		utils.SetPropWithNilCheck(volumeMap, "expose_serial", volume.Properties.ExposeSerial)
		utils.SetPropWithNilCheck(volumeMap, "require_legacy_bios", volume.Properties.RequireLegacyBios)
	}
	return volumeMap
}

func ResourceCommonServerDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(bundleclient.SdkBundle).CloudApiClient
	dcId := d.Get("datacenter_id").(string)

	apiResponse, err := client.ServersApi.DatacentersServersDelete(ctx, dcId, d.Id()).Execute()
	if apiResponse != nil {
		log.Printf("[DEBUG] Request time : %s for operation : %s",
			apiResponse.RequestTime, apiResponse.Operation)
		if apiResponse.Response != nil {
			log.Printf("[DEBUG] response status code : %d\n", apiResponse.StatusCode)
		}
	}
	if err != nil {
		diags := diag.FromErr(fmt.Errorf("an error occurred while deleting a server ID %s %w", d.Id(), err))
		return diags

	}

	if errState := bundleclient.WaitForStateChange(ctx, meta, d, apiResponse, schema.TimeoutDelete); errState != nil {
		return diag.FromErr(fmt.Errorf("error getting state change for server delete %w", errState))
	}

	d.SetId("")
	return nil
}
