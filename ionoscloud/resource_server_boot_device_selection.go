package ionoscloud

import (
	"context"
	"errors"
	"fmt"

	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/cloudapi/cloudapiserver"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils/uuidgen"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"

	ionoscloud "github.com/ionos-cloud/sdk-go/v6"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceServerBootDeviceSelection() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceServerBootDeviceSelectionCreate,
		ReadContext:   resourceServerBootDeviceSelectionRead,
		UpdateContext: resourceServerBootDeviceSelectionUpdate,
		DeleteContext: resourceServerBootDeviceSelectionDelete,
		Schema: map[string]*schema.Schema{

			"datacenter_id": {
				Type:             schema.TypeString,
				Description:      "ID of the Datacenter that holds the server for which the boot volume is selected",
				Required:         true,
				ForceNew:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.IsUUID),
			},
			"server_id": {
				Type:             schema.TypeString,
				Description:      "ID of the Server for which the boot device will be selected.",
				Required:         true,
				ForceNew:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.IsUUID),
			},
			"boot_device_id": {
				Type:             schema.TypeString,
				Description:      "ID of the entity to set as primary boot device. Possible boot devices are CDROM Images and Volumes. If omitted, server will boot from PXE",
				Optional:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.IsUUID),
			},
			"default_boot_volume_id": {
				Type:        schema.TypeString,
				Description: "ID of the first attached volume of the Server, which will be the default boot volume.",
				Computed:    true,
			},
		},

		Timeouts: &resourceDefaultTimeouts,
	}
}

func resourceServerBootDeviceSelectionCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	dcId := d.Get("datacenter_id").(string)
	serverId := d.Get("server_id").(string)

	ss := cloudapiserver.NewUnboundService(serverId, meta)

	// The bootable device to which the server will revert if this resource is destroyed.
	defaultBootVolume, err := ss.GetDefaultBootVolume(ctx, dcId, serverId)
	if err != nil {
		return diag.FromErr(err)
	}

	if err = d.Set("default_boot_volume_id", defaultBootVolume.Id); err != nil {
		return diag.FromErr(fmt.Errorf("error setting a default boot volume for boot selection resource"))
	}

	bootDeviceIDValue, bootDeviceIDOk := d.GetOk("boot_device_id")
	if !bootDeviceIDOk {
		if err = ss.PxeBoot(ctx, dcId, serverId); err != nil {
			return diag.FromErr(fmt.Errorf("error while performing pxe boot for server, serverId: %s, dcId: %s (%w)", serverId, dcId, err))
		}
	} else {
		bootDeviceID := bootDeviceIDValue.(string)
		if err = ss.UpdateBootDevice(ctx, dcId, serverId, bootDeviceID); err != nil {
			return diag.FromErr(err)
		}
	}

	d.SetId(uuidgen.ResourceUuid().String())
	return resourceServerBootDeviceSelectionRead(ctx, d, meta)
}

func resourceServerBootDeviceSelectionRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	ss := cloudapiserver.Service{Client: meta.(services.SdkBundle).CloudAPIClient, Meta: meta, D: d}

	dcId := d.Get("datacenter_id").(string)
	serverId := d.Get("server_id").(string)

	server, err := ss.FindById(ctx, dcId, serverId, 3)
	if err != nil {
		if errors.Is(err, cloudapiserver.ErrServerNotFound) {
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}

	if err = setServerBootDeviceSelectionData(d, server); err != nil {
		return diag.FromErr(fmt.Errorf("error reading boot devices for server, dcId: %s, sId: %s, (%w)", dcId, serverId, err))
	}

	return nil
}

func resourceServerBootDeviceSelectionUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	dcId := d.Get("datacenter_id").(string)
	serverId := d.Get("server_id").(string)

	ss := cloudapiserver.NewUnboundService(serverId, meta)

	if d.HasChange("boot_device_id") {
		bootDeviceIDValue, bootDeviceIDOk := d.GetOk("boot_device_id")
		if !bootDeviceIDOk {
			if err := ss.PxeBoot(ctx, dcId, serverId); err != nil {
				return diag.FromErr(fmt.Errorf("error while performing pxe boot: %s, serverId: %s, dcId: %s", err.Error(), serverId, dcId))
			}
		} else {
			bootDeviceID := bootDeviceIDValue.(string)
			if err := ss.UpdateBootDevice(ctx, dcId, serverId, bootDeviceID); err != nil {
				return diag.FromErr(err)
			}
		}
	}
	return resourceServerBootDeviceSelectionRead(ctx, d, meta)
}

func resourceServerBootDeviceSelectionDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	ss := cloudapiserver.Service{Client: meta.(services.SdkBundle).CloudAPIClient, Meta: meta, D: d}

	dcId := d.Get("datacenter_id").(string)
	serverId := d.Get("server_id").(string)
	defaultBootVolumeId := d.Get("default_boot_volume_id").(string)

	if err := ss.UpdateBootDevice(ctx, dcId, serverId, defaultBootVolumeId); err != nil {
		return diag.FromErr(err)
	}
	d.SetId("")
	return nil
}

func setServerBootDeviceSelectionData(d *schema.ResourceData, server *ionoscloud.Server) error {

	if server.Properties.BootCdrom != nil {
		if err := d.Set("boot_device_id", *server.Properties.BootCdrom.Id); err != nil {
			return err
		}
		return nil
	}

	if server.Properties.BootVolume != nil {
		if err := d.Set("boot_device_id", *server.Properties.BootVolume.Id); err != nil {
			return err
		}
	}

	return nil
}
