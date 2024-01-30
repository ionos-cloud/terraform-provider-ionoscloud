package cloudapiserver

import (
	"context"
	"errors"
	"fmt"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	ionoscloud "github.com/ionos-cloud/sdk-go/v6"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/cloudapi"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils/constant"
)

var (
	ErrSuspendCubeLast error
	ErrServerNotFound  error
	ErrNoBootDevice    = errors.New("server has no boot device")
)

// UnboundService allows usage of a subset of the method set of the Service
// This is useful when using the Service in Provider Context functions that belong to a different resource which needs to interact with an already existing server
// In this case, the Service will not be 'bound' to the state of the Server resource
// The methods of this interface must not write to the state, since the Service will use a placeholder ResourceData reference
type UnboundService interface {
	Update(context.Context, string, string, ionoscloud.ServerProperties) (*ionoscloud.Server, *ionoscloud.APIResponse, error)
	GetDefaultBootVolume(ctx context.Context, datacenterId, serverId string) (*ionoscloud.Volume, error)
	UpdateBootDevice(ctx context.Context, datacenterID, serverID, newBootDeviceID string) error
	PxeBoot(ctx context.Context, datacenterID, serverID string) error
	Reboot(ctx context.Context, datacenterID, serverID string) error
}

type Service struct {
	Client *ionoscloud.APIClient
	Meta   any
	D      *schema.ResourceData
}

// NewUnboundService creates an UnboundService with a subset of the underlying Service methods
// The concrete Service is created with a dummy ResourceData reference which has the ID of the Server this service will interact with
// This ensure state tracking functions such as WaitForResourceToBeReady use the correct ID
func NewUnboundService(serverId string, meta any) UnboundService {
	client := meta.(services.SdkBundle).CloudApiClient
	d := &schema.ResourceData{}
	d.SetId(serverId)
	return &Service{client, meta, d}
}

func (ss *Service) FindById(ctx context.Context, datacenterID, serverID string, depth int32) (*ionoscloud.Server, error) {
	server, apiResponse, err := ss.Client.ServersApi.DatacentersServersFindById(ctx, datacenterID, serverID).Depth(depth).Execute()
	apiResponse.LogInfo()
	if err != nil {
		if apiResponse.HttpNotFound() {
			log.Printf("[DEBUG] cannot find server by id in datacenter dcId: %s, serverId: %s\n", datacenterID, serverID)
			return nil, ErrServerNotFound
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
		return *server.Properties.BootCdrom.Id, constant.BootDeviceTypeCDROM, nil
	}
	if server.Properties.BootVolume != nil {
		return *server.Properties.BootVolume.Id, constant.BootDeviceTypeVolume, nil
	}
	return "", "", ErrNoBootDevice
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
	case constant.EnterpriseType, constant.VCPUType:
		if strings.EqualFold(newVmState, constant.CubeVMStateStop) {
			return fmt.Errorf("cannot suspend a %s server, set to %s instead", serverType, constant.VMStateStop)
		}
		if strings.EqualFold(newVmState, constant.VMStateStart) && strings.EqualFold(currentVmState, constant.VMStateStop) {
			return ss.Start(ctx, datacenterID, serverID, serverType)
		}
		if strings.EqualFold(newVmState, constant.VMStateStop) && strings.EqualFold(currentVmState, constant.VMStateStart) {
			return ss.Stop(ctx, datacenterID, serverID, serverType)
		}

	case constant.CubeType:
		if strings.EqualFold(newVmState, constant.VMStateStop) {
			return fmt.Errorf("cannot shut down a %s server, set to %s instead", serverType, constant.CubeVMStateStop)
		}
		if strings.EqualFold(newVmState, constant.VMStateStart) && strings.EqualFold(currentVmState, constant.CubeVMStateStop) {
			return ss.Start(ctx, datacenterID, serverID, serverType)
		}
		if strings.EqualFold(newVmState, constant.CubeVMStateStop) && strings.EqualFold(currentVmState, constant.VMStateStart) {
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
	var oldBdType string
	oldBootDeviceID, oldBdType, err := ss.GetCurrentBootDeviceID(ctx, datacenterID, serverID)
	if err != nil {
		if !errors.Is(err, ErrNoBootDevice) {
			return err
		}
		oldBdType = constant.BootDeviceTypeVolume
	}
	if strings.EqualFold(oldBootDeviceID, newBootDeviceID) {
		return nil
	}

	newBdType := constant.BootDeviceTypeCDROM
	_, apiResponse, err := ss.Client.ImagesApi.ImagesFindById(ctx, newBootDeviceID).Execute()
	if err != nil {
		if !apiResponse.HttpNotFound() {
			return err
		}
		log.Printf("[DEBUG] no bootable image found with id : %s\n", newBootDeviceID)
		newBdType = constant.BootDeviceTypeVolume
	}

	switch oldBdType {
	case constant.BootDeviceTypeCDROM:
		if strings.EqualFold(newBdType, constant.BootDeviceTypeVolume) {
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

	case constant.BootDeviceTypeVolume:
		// no cdrom is detached, only update to the new boot device, regardless of type
		sp := ionoscloud.ServerProperties{BootVolume: ionoscloud.NewResourceReference(newBootDeviceID)}
		if strings.EqualFold(newBdType, constant.BootDeviceTypeCDROM) {
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

	case constant.EnterpriseType, constant.VCPUType:
		apiResponse, err := ss.Client.ServersApi.DatacentersServersStartPost(ctx, datacenterID, serverID).Execute()
		apiResponse.LogInfo()
		if err != nil {
			return err
		}
		if err = utils.WaitForResourceToBeReady(ctx, ss.D, ss.checkExpectedVmStateFn(datacenterID, constant.VMStateStart)); err != nil {
			return err
		}
		if errState := cloudapi.WaitForStateChange(ctx, ss.Meta, ss.D, apiResponse, schema.TimeoutUpdate); errState != nil {
			return fmt.Errorf("an error occured while waiting for server state change on VM POWER ON dcId: %s, server_id: %s, Response: (%w)", datacenterID, serverID, errState)
		}
		log.Printf("[DEBUG] %s server powered on: serverId: %s \n", serverType, serverID)
		return nil

	case constant.CubeType:
		apiResponse, err := ss.Client.ServersApi.DatacentersServersResumePost(ctx, datacenterID, serverID).Execute()
		apiResponse.LogInfo()
		if err != nil {
			return err
		}
		if err = utils.WaitForResourceToBeReady(ctx, ss.D, ss.checkExpectedVmStateFn(datacenterID, constant.VMStateStart)); err != nil {
			return err
		}
		if errState := cloudapi.WaitForStateChange(ctx, ss.Meta, ss.D, apiResponse, schema.TimeoutUpdate); errState != nil {
			return fmt.Errorf("an error occured while waiting for server state change on VM RESUME dcId: %s, server_id: %s, Response: (%w)", datacenterID, serverID, errState)
		}
		log.Printf("[DEBUG] %s server unsuspended: serverId: %s \n", serverType, serverID)
		return nil

	}

	return fmt.Errorf("cannot start unknown server type: %s", serverType)

}

func (ss *Service) Stop(ctx context.Context, datacenterID, serverID, serverType string) error {

	switch serverType {

	case constant.EnterpriseType, constant.VCPUType:
		apiResponse, err := ss.Client.ServersApi.DatacentersServersStopPost(ctx, datacenterID, serverID).Execute()
		apiResponse.LogInfo()
		if err != nil {
			return err
		}
		if err = utils.WaitForResourceToBeReady(ctx, ss.D, ss.checkExpectedVmStateFn(datacenterID, constant.VMStateStop)); err != nil {
			return err
		}
		if errState := cloudapi.WaitForStateChange(ctx, ss.Meta, ss.D, apiResponse, schema.TimeoutUpdate); errState != nil {
			return fmt.Errorf("an error occured while waiting for server state change on VM SHUTOFF dcId: %s, server_id: %s, Response: (%w)", datacenterID, serverID, errState)
		}
		log.Printf("[DEBUG] %s server powered off: serverId: %s \n", serverType, serverID)
		return nil

	case constant.CubeType:
		apiResponse, err := ss.Client.ServersApi.DatacentersServersSuspendPost(ctx, datacenterID, serverID).Execute()
		apiResponse.LogInfo()
		if err != nil {
			return err
		}
		if err = utils.WaitForResourceToBeReady(ctx, ss.D, ss.checkExpectedVmStateFn(datacenterID, constant.CubeVMStateStop)); err != nil {
			return err
		}
		if errState := cloudapi.WaitForStateChange(ctx, ss.Meta, ss.D, apiResponse, schema.TimeoutUpdate); errState != nil {
			return fmt.Errorf("an error occured while waiting for server state change on VM SUSPEND dcId: %s, server_id: %s, Response: (%w)", datacenterID, serverID, errState)
		}
		log.Printf("[DEBUG] %s server suspended: serverId: %s \n", serverType, serverID)
		return nil
	}

	return fmt.Errorf("cannot stop unknown server type: %s", serverType)

}

func (ss *Service) Reboot(ctx context.Context, datacenterID, serverID string) error {

	apiResponse, err := ss.Client.ServersApi.DatacentersServersRebootPost(ctx, datacenterID, serverID).Execute()
	apiResponse.LogInfo()
	if err != nil {
		return err
	}
	if err = utils.WaitForResourceToBeReady(ctx, ss.D, ss.checkExpectedVmStateFn(datacenterID, constant.VMStateStart)); err != nil {
		return err
	}
	if errState := cloudapi.WaitForStateChange(ctx, ss.Meta, ss.D, apiResponse, schema.TimeoutUpdate); errState != nil {
		return fmt.Errorf("an error occured while waiting for server state change on reboot dcId: %s, server_id: %s, Response: (%w)", datacenterID, serverID, errState)
	}
	log.Printf("[DEBUG] server reboot finished: serverId: %s \n", serverID)
	return nil
}

func (ss *Service) PxeBoot(ctx context.Context, datacenterID, serverID string) error {

	deviceID, deviceType, err := ss.GetCurrentBootDeviceID(ctx, datacenterID, serverID)
	if err != nil {
		if errors.Is(err, ErrNoBootDevice) {
			return ss.Reboot(ctx, datacenterID, serverID)
		}
		return err
	}

	switch deviceType {
	case constant.BootDeviceTypeVolume:
		volumeProperties := ionoscloud.VolumeProperties{}
		volumeProperties.SetBootOrder(constant.VolumeBootOrderNone)
		_, apiResponse, err := ss.Client.VolumesApi.DatacentersVolumesPatch(ctx, datacenterID, deviceID).Volume(volumeProperties).Execute()
		if err != nil {
			return err
		}
		if errState := cloudapi.WaitForStateChange(ctx, ss.Meta, ss.D, apiResponse, schema.TimeoutUpdate); errState != nil {
			return errState
		}
		log.Printf("[DEBUG] enabled PXE boot for server: serverId: %s \n", serverID)

	case constant.BootDeviceTypeCDROM:
		apiResponse, err := ss.Client.ServersApi.DatacentersServersCdromsDelete(ctx, datacenterID, serverID, deviceID).Execute()
		if err != nil {
			return err
		}
		if errState := cloudapi.WaitForStateChange(ctx, ss.Meta, ss.D, apiResponse, schema.TimeoutUpdate); errState != nil {
			return errState
		}
		log.Printf("[DEBUG] detached CDROM image from server: serverId: %s, imageId: %s\n", serverID, deviceID)
	}

	return ss.Reboot(ctx, datacenterID, serverID)
}

// checkExpectedVmStateFn wraps over the ResourceReadyFunc to allow passing expectedState
func (ss *Service) checkExpectedVmStateFn(dcId, expectedState string) utils.ResourceReadyFunc {

	return func(ctx context.Context, d *schema.ResourceData) (bool, error) {
		server, _, err := ss.Client.ServersApi.DatacentersServersFindById(ctx, dcId, d.Id()).Execute()
		if err != nil {
			return false, err
		}

		serverType := *server.Properties.Type
		if !strings.EqualFold(*server.Properties.VmState, expectedState) {
			log.Printf("[INFO] Server (type: %s) vmState not yet changed to %s: %s", serverType, expectedState, d.Id())
			return false, nil
		}
		return true, nil
	}
}
