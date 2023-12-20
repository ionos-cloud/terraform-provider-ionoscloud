package cloudapiserver

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	ionoscloud "github.com/ionos-cloud/sdk-go/v6"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/cloudapi"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils"
)

type Service struct {
	Client *ionoscloud.APIClient
	Meta   interface{}
	D      *schema.ResourceData
}

// The caller should ignore this error, it only informs that the CUBE server should be suspended after all other updates have been applied.
var ErrSuspendCubeLast error
var ServerNotFound error

const (
	CubeServerType       = "CUBE"
	EnterpriseServerType = "ENTERPRISE"
	VCPUServerType       = "VCPU"

	CubeVMStateStop = "SUSPENDED"

	// These are the vm_state values that are available for VCPU and ENTERPRISE servers
	VMStateStart = "RUNNING"
	VMStateStop  = "SHUTOFF"

	// Types of bootable devices which servers use
	BootDeviceTypeVolume = "volume"
	BootDeviceTypeCDROM  = "cdrom"
)

func (ss *Service) FindById(ctx context.Context, datacenterID, serverID string, depth int32) (*ionoscloud.Server, error) {
	server, apiResponse, err := ss.Client.ServersApi.DatacentersServersFindById(ctx, datacenterID, serverID).Depth(depth).Execute()
	apiResponse.LogInfo()
	if err != nil {
		if apiResponse.HttpNotFound() {
			log.Printf("[DEBUG] cannot find server by id in datacenter dcId: %s, serverId: %s\n", datacenterID, serverID)
			return nil, ServerNotFound
		}
		return nil, err
	}
	return &server, nil
}

func (ss *Service) Delete(ctx context.Context, datacenterID, serverID, ID string) (*ionoscloud.APIResponse, error) {
	apiResponse, err := ss.Client.ServersApi.DatacentersServersDelete(ctx, datacenterID, serverID).Execute()
	apiResponse.LogInfo()
	if err != nil {
		return apiResponse, err
	}
	if errState := cloudapi.WaitForStateChange(ctx, ss.Meta, ss.D, apiResponse, schema.TimeoutDelete); errState != nil {
		return apiResponse, fmt.Errorf("an error occured while waiting for server state change on delete dcId: %s, server_id: %s, ID: %s, Response: (%w)", datacenterID, serverID, ID, errState)
	}
	return apiResponse, nil
}

func (ss *Service) Create(ctx context.Context, datacenterID string) (*ionoscloud.Server, *ionoscloud.APIResponse, error) {
	server, apiResponse, err := ss.Client.ServersApi.DatacentersServersPost(ctx, datacenterID).Execute()
	apiResponse.LogInfo()
	if err != nil {
		return nil, apiResponse, fmt.Errorf("an error occured while creating server for dcId: %s, Response: (%w)", datacenterID, err)
	}
	if errState := cloudapi.WaitForStateChange(ctx, ss.Meta, ss.D, apiResponse, schema.TimeoutCreate); errState != nil {
		if cloudapi.IsRequestFailed(errState) {
			ss.D.SetId("")
		}
		return nil, apiResponse, fmt.Errorf("an error occured while waiting for server state change on create dcId: %s, Response: (%w)", datacenterID, errState)
	}
	return &server, apiResponse, nil
}

func (ss *Service) Update(ctx context.Context, datacenterID, serverID string, serverProperties ionoscloud.ServerProperties) (*ionoscloud.Server, *ionoscloud.APIResponse, error) {
	updatedServer, apiResponse, err := ss.Client.ServersApi.DatacentersServersPatch(ctx, datacenterID, serverID).Server(serverProperties).Execute()
	apiResponse.LogInfo()
	if err != nil {
		return nil, apiResponse, fmt.Errorf("an error occured while updating server for dcId: %s, server_id: %s, Response: (%w)", datacenterID, serverID, err)
	}
	if errState := cloudapi.WaitForStateChange(ctx, ss.Meta, ss.D, apiResponse, schema.TimeoutUpdate); errState != nil {
		return nil, apiResponse, fmt.Errorf("an error occured while waiting for server state change on update dcId: %s, server_id: %s, Response: (%w)", datacenterID, serverID, errState)
	}
	return &updatedServer, apiResponse, nil
}

func (ss *Service) GetAttachedVolumes(ctx context.Context, datacenterID, serverID string) ([]*ionoscloud.Volume, *ionoscloud.APIResponse, error) {

	attachedVolumeIds, apiResponse, err := ss.Client.ServersApi.DatacentersServersVolumesGet(ctx, datacenterID, serverID).Execute()
	apiResponse.LogInfo()
	if err != nil {
		return nil, apiResponse, fmt.Errorf("an error occured while fetching attached volumes for server, dcId: %s, serverId: %s, Response: (%w)", datacenterID, serverID, err)
	}
	attachedVolumes := []*ionoscloud.Volume{}
	for _, v := range *attachedVolumeIds.Items {
		volume, apiResponse, err := ss.Client.ServersApi.DatacentersServersVolumesFindById(ctx, datacenterID, serverID, *v.Id).Execute()
		if err != nil {
			return nil, apiResponse, err
		}
		attachedVolumes = append(attachedVolumes, &volume)
	}

	return attachedVolumes, apiResponse, nil
}

func (ss *Service) GetDefaultBootVolume(ctx context.Context, datacenterId, serverId string) (*ionoscloud.Volume, error) {
	attachedVolumes, _, err := ss.GetAttachedVolumes(ctx, datacenterId, serverId)
	if err != nil {
		return nil, err
	}

	var defaultBootVolume ionoscloud.Volume
	firstCreatedTime := time.Now()
	for _, v := range attachedVolumes {
		if v.Metadata.CreatedDate.Before(firstCreatedTime) {
			defaultBootVolume = *v
			firstCreatedTime = v.Metadata.CreatedDate.Time
		}
	}
	return &defaultBootVolume, nil
}

func (ss *Service) GetCurrentBootDeviceID(ctx context.Context, datacenterId, serverId string) (string, string, error) {
	server, err := ss.FindById(ctx, datacenterId, serverId, 3)
	if err != nil {
		return "", "", err
	}
	if server.Properties == nil {
		return "", "", fmt.Errorf("server has no boot device because properties object was nil")
	}
	if server.Properties.BootCdrom != nil {
		return *server.Properties.BootCdrom.Id, BootDeviceTypeCDROM, nil
	}
	if server.Properties.BootVolume != nil {
		return *server.Properties.BootVolume.Id, BootDeviceTypeVolume, nil
	}
	return "", "", fmt.Errorf("server has no boot device")
}

func (ss *Service) UpdateVmState(ctx context.Context, datacenterID, serverID, newVmState string) error {

	var serverType, currentVmState string
	var err error

	if serverType, err = ss.GetServerType(ctx, datacenterID, serverID); err != nil {
		return err
	}
	if currentVmState, err = ss.GetVmState(ctx, datacenterID, serverID); err != nil {
		return err
	}

	switch serverType {
	case EnterpriseServerType, VCPUServerType:
		if strings.EqualFold(newVmState, CubeVMStateStop) {
			return fmt.Errorf("cannot suspend a %s server, set to %s instead", serverType, VMStateStop)
		}
		if strings.EqualFold(newVmState, VMStateStart) && strings.EqualFold(currentVmState, VMStateStop) {
			return ss.Start(ctx, datacenterID, serverID, serverType)
		}
		if strings.EqualFold(newVmState, VMStateStop) && strings.EqualFold(currentVmState, VMStateStart) {
			return ss.Stop(ctx, datacenterID, serverID, serverType)
		}

	case CubeServerType:
		if strings.EqualFold(newVmState, VMStateStop) {
			return fmt.Errorf("cannot shut down a %s server, set to %s instead", serverType, CubeVMStateStop)
		}
		if strings.EqualFold(newVmState, VMStateStart) && strings.EqualFold(currentVmState, CubeVMStateStop) {
			return ss.Start(ctx, datacenterID, serverID, serverType)
		}
		if strings.EqualFold(newVmState, CubeVMStateStop) && strings.EqualFold(currentVmState, VMStateStart) {
			return ErrSuspendCubeLast
		}

	}

	return nil
}

// UpdateBootDevice will set a new boot device for the server, which can be a volume or bootable image CDROM.
// When the new boot device is a volume, it must be already attached to the server.
// When the new boot device is a CDROM image, it will be attached by default.
// If the current boot device is a CDROM image, it will be detached after it is changed by this operation.
func (ss *Service) UpdateBootDevice(ctx context.Context, datacenterID, serverID, newBootDeviceID string) error {
	oldBootDeviceID, oldBdType, err := ss.GetCurrentBootDeviceID(ctx, datacenterID, serverID)
	if err != nil {
		return err
	}
	if strings.EqualFold(oldBootDeviceID, newBootDeviceID) {
		return nil
	}

	newBdType := BootDeviceTypeCDROM
	_, apiResponse, err := ss.Client.ImagesApi.ImagesFindById(ctx, newBootDeviceID).Execute()
	if err != nil {
		if !apiResponse.HttpNotFound() {
			return err
		}
		log.Printf("[DEBUG] no bootable image found with id : %s\n", newBootDeviceID)
		newBdType = BootDeviceTypeVolume
	}

	switch oldBdType {
	case BootDeviceTypeCDROM:
		if strings.EqualFold(newBdType, BootDeviceTypeVolume) {
			// update to new boot volume
			sp := ionoscloud.ServerProperties{BootVolume: ionoscloud.NewResourceReference(newBootDeviceID)}
			if _, _, err = ss.Update(ctx, datacenterID, serverID, sp); err != nil {
				return err
			}
		} else {
			// attach new cdrom
			img := ionoscloud.Image{Id: &newBootDeviceID}
			_, apiResponse, err := ss.Client.ServersApi.DatacentersServersCdromsPost(ctx, datacenterID, serverID).Cdrom(img).Execute()
			if err != nil {
				return err
			}
			if errState := cloudapi.WaitForStateChange(ctx, ss.Meta, ss.D, apiResponse, schema.TimeoutUpdate); errState != nil {
				return errState
			}
			log.Printf("[DEBUG] attached CDROM image to server: serverId: %s, imageId: %s\n", serverID, newBootDeviceID)
			// update to new boot cdrom
			sp := ionoscloud.ServerProperties{BootCdrom: ionoscloud.NewResourceReference(newBootDeviceID)}
			if _, _, err = ss.Update(ctx, datacenterID, serverID, sp); err != nil {
				return err
			}
		}
		// detach old cdrom
		apiResponse, err = ss.Client.ServersApi.DatacentersServersCdromsDelete(ctx, datacenterID, serverID, oldBootDeviceID).Execute()
		if err != nil {
			return err
		}
		if errState := cloudapi.WaitForStateChange(ctx, ss.Meta, ss.D, apiResponse, schema.TimeoutUpdate); errState != nil {
			return errState
		}
		log.Printf("[DEBUG] detached CDROM image from server: serverId: %s, imageId: %s\n", serverID, oldBootDeviceID)

	case BootDeviceTypeVolume:
		// no cdrom is detached, only update to the new boot device, regardless of type
		sp := ionoscloud.ServerProperties{BootVolume: ionoscloud.NewResourceReference(newBootDeviceID)}
		if strings.EqualFold(newBdType, BootDeviceTypeCDROM) {
			sp = ionoscloud.ServerProperties{BootCdrom: ionoscloud.NewResourceReference(newBootDeviceID)}
		}
		if _, _, err = ss.Update(ctx, datacenterID, serverID, sp); err != nil {
			return err
		}
	}

	return nil
}

func (ss *Service) GetVmState(ctx context.Context, datacenterID, serverID string) (string, error) {
	server, err := ss.FindById(ctx, datacenterID, serverID, 0)
	if err != nil {
		return "", err
	}
	if server.Properties == nil {
		return "", fmt.Errorf("got empty properties for datacenterID %s serverID %s", datacenterID, serverID)
	}
	return *server.Properties.VmState, nil
}

func (ss *Service) GetServerType(ctx context.Context, datacenterID, serverID string) (string, error) {
	server, err := ss.FindById(ctx, datacenterID, serverID, 0)
	if err != nil {
		return "", err
	}
	if server.Properties == nil {
		return "", fmt.Errorf("got empty properties for datacenterID %s serverID %s", datacenterID, serverID)
	}
	return *server.Properties.Type, nil
}

func (ss *Service) Start(ctx context.Context, datacenterID, serverID, serverType string) error {

	switch serverType {

	case EnterpriseServerType, VCPUServerType:
		apiResponse, err := ss.Client.ServersApi.DatacentersServersStartPost(ctx, datacenterID, serverID).Execute()
		apiResponse.LogInfo()
		if err != nil {
			return err
		}
		return utils.WaitForResourceToBeReady(ctx, ss.D, ss.checkExpectedVmStateFn(ctx, datacenterID, VMStateStart))

	case CubeServerType:
		apiResponse, err := ss.Client.ServersApi.DatacentersServersResumePost(ctx, datacenterID, serverID).Execute()
		apiResponse.LogInfo()
		if err != nil {
			return err
		}
		return utils.WaitForResourceToBeReady(ctx, ss.D, ss.checkExpectedVmStateFn(ctx, datacenterID, VMStateStart))
	}

	return fmt.Errorf("cannot start unknown server type: %s", serverType)

}

func (ss *Service) Stop(ctx context.Context, datacenterID, serverID, serverType string) error {

	switch serverType {

	case EnterpriseServerType, VCPUServerType:
		apiResponse, err := ss.Client.ServersApi.DatacentersServersStopPost(ctx, datacenterID, serverID).Execute()
		apiResponse.LogInfo()
		if err != nil {
			return err
		}
		return utils.WaitForResourceToBeReady(ctx, ss.D, ss.checkExpectedVmStateFn(ctx, datacenterID, VMStateStop))

	case CubeServerType:
		apiResponse, err := ss.Client.ServersApi.DatacentersServersSuspendPost(ctx, datacenterID, serverID).Execute()
		apiResponse.LogInfo()
		if err != nil {
			return err
		}
		return utils.WaitForResourceToBeReady(ctx, ss.D, ss.checkExpectedVmStateFn(ctx, datacenterID, CubeVMStateStop))
	}

	return fmt.Errorf("cannot stop unknown server type: %s", serverType)

}

// checkExpectedVmStateFn wraps over the ResourceReadyFunc to allow passing expectedState
// TODO: change ResourceReadyFunc sig to support passing an expectedState param
func (ss *Service) checkExpectedVmStateFn(ctx context.Context, dcId, expectedState string) utils.ResourceReadyFunc {

	return func(ctx context.Context, d *schema.ResourceData) (bool, error) {
		ionoscloudServer, _, err := ss.Client.ServersApi.DatacentersServersFindById(ctx, dcId, d.Id()).Execute()
		if err != nil {
			return false, err
		}

		serverType := *ionoscloudServer.Properties.Type
		if !strings.EqualFold(*ionoscloudServer.Properties.VmState, expectedState) {
			log.Printf("[INFO] Server (type: %s) vmState not yet changed to %s: %s", serverType, expectedState, d.Id())
			return false, nil
		}
		return true, nil
	}
}
