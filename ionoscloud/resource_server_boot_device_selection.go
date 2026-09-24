package ionoscloud

import (
	"context"
	"errors"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	ionoscloud "github.com/ionos-cloud/sdk-go-bundle/products/compute/v2"

	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/bundleclient"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/cloudapi/cloudapiserver"
	diagutil "github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils/diags"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils/uuidgen"
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
			"location": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The location of the resource. This field should be used only if you are also using a file configuration and should not be configured otherwise.",
				ForceNew:    true,
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

func resourceServerBootDeviceSelectionCreate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	dcID := d.Get("datacenter_id").(string)
	serverID := d.Get("server_id").(string)
	location := d.Get("location").(string)

	ss, err := cloudapiserver.NewUnboundService(ctx, serverID, location, meta)
	if err != nil {
		return diag.FromErr(err)
	}

	// The bootable device to which the server will revert if this resource is destroyed.
	defaultBootVolume, err := ss.GetDefaultBootVolume(ctx, dcID, serverID)
	if err != nil {
		return diagutil.ToDiags(d, err, nil)
	}

	if err = d.Set("default_boot_volume_id", defaultBootVolume.Id); err != nil {
		return diagutil.ToDiags(d, fmt.Errorf("error setting a default boot volume for boot selection resource"), nil)
	}

	bootDeviceIDValue, bootDeviceIDOk := d.GetOk("boot_device_id")
	if !bootDeviceIDOk {
		if err = ss.PxeBoot(ctx, dcID, serverID); err != nil {
			return diagutil.ToDiags(d, fmt.Errorf("error while performing pxe boot for server, serverID: %s, dcID: %s (%w)", serverID, dcID, err), nil)
		}
	} else {
		bootDeviceID := bootDeviceIDValue.(string)
		if err = ss.UpdateBootDevice(ctx, dcID, serverID, bootDeviceID); err != nil {
			return diagutil.ToDiags(d, err, nil)
		}
	}

	d.SetId(uuidgen.ResourceUuid().String())
	return resourceServerBootDeviceSelectionRead(ctx, d, meta)
}

func resourceServerBootDeviceSelectionRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	location := d.Get("location").(string)
	client, err := meta.(bundleclient.SdkBundle).NewCloudAPIClient(ctx, location)
	if err != nil {
		return diag.FromErr(err)
	}
	ss := cloudapiserver.Service{Client: client, Meta: meta, D: d}

	dcID := d.Get("datacenter_id").(string)
	serverID := d.Get("server_id").(string)

	server, err := ss.FindById(ctx, dcID, serverID, 3)
	if err != nil {
		if errors.Is(err, cloudapiserver.ErrServerNotFound) {
			d.SetId("")
			return nil
		}
		return diagutil.ToDiags(d, err, nil)
	}

	if err = setServerBootDeviceSelectionData(d, server); err != nil {
		return diagutil.ToDiags(d, fmt.Errorf("error reading boot devices for server, dcID: %s, sID: %s, (%w)", dcID, serverID, err), nil)
	}

	return nil
}

func resourceServerBootDeviceSelectionUpdate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	dcID := d.Get("datacenter_id").(string)
	serverID := d.Get("server_id").(string)
	location := d.Get("location").(string)

	ss, err := cloudapiserver.NewUnboundService(ctx, serverID, location, meta)
	if err != nil {
		return diag.FromErr(err)
	}

	if d.HasChange("boot_device_id") {
		bootDeviceIDValue, bootDeviceIDOk := d.GetOk("boot_device_id")
		if !bootDeviceIDOk {
			if err := ss.PxeBoot(ctx, dcID, serverID); err != nil {
				return diagutil.ToDiags(d, fmt.Errorf("error while performing pxe boot: %w, serverID: %s, dcID: %s", err, serverID, dcID), nil)
			}
		} else {
			bootDeviceID := bootDeviceIDValue.(string)
			if err := ss.UpdateBootDevice(ctx, dcID, serverID, bootDeviceID); err != nil {
				return diagutil.ToDiags(d, err, nil)
			}
		}
	}
	return resourceServerBootDeviceSelectionRead(ctx, d, meta)
}

func resourceServerBootDeviceSelectionDelete(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	location := d.Get("location").(string)
	client, err := meta.(bundleclient.SdkBundle).NewCloudAPIClient(ctx, location)
	if err != nil {
		return diag.FromErr(err)
	}
	ss := cloudapiserver.Service{Client: client, Meta: meta, D: d}

	dcID := d.Get("datacenter_id").(string)
	serverID := d.Get("server_id").(string)
	defaultBootVolumeID := d.Get("default_boot_volume_id").(string)

	if err := ss.UpdateBootDevice(ctx, dcID, serverID, defaultBootVolumeID); err != nil {
		return diagutil.ToDiags(d, err, nil)
	}
	d.SetId("")
	return nil
}

func setServerBootDeviceSelectionData(d *schema.ResourceData, server *ionoscloud.Server) error {

	if server.Properties.BootCdrom != nil && server.Properties.BootCdrom.Id != "" {
		if err := d.Set("boot_device_id", server.Properties.BootCdrom.Id); err != nil {
			return err
		}
		return nil
	}

	if server.Properties.BootVolume != nil && server.Properties.BootVolume.Id != "" {
		if err := d.Set("boot_device_id", server.Properties.BootVolume.Id); err != nil {
			return err
		}
	}

	return nil
}
