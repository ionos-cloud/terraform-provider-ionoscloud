package cloudapiserver

import (
	"context"
	"fmt"
	"log"
	"strings"

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

var ErrSuspendCubeLast error

const (
	CUBE_SERVER_TYPE   = "CUBE"
	CUBE_VMSTATE_STOP  = "SUSPENDED"
	CUBE_VMSTATE_START = "RUNNING"

	ENTERPRISE_SERVER_TYPE   = "ENTERPRISE"
	ENTERPRISE_VMSTATE_STOP  = "SHUTOFF"
	ENTERPRISE_VMSTATE_START = "RUNNING"
)

func (ss *Service) FindById(ctx context.Context, datacenterID, serverID string, depth int32) (*ionoscloud.Server, error) {
	server, apiResponse, err := ss.Client.ServersApi.DatacentersServersFindById(ctx, datacenterID, serverID).Depth(depth).Execute()
	apiResponse.LogInfo()
	if err != nil {
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
	// Wait, catching any errors
	_, errState := cloudapi.GetStateChangeConf(ss.Meta, ss.D, apiResponse.Header.Get("Location"), schema.TimeoutDelete).WaitForStateContext(ctx)
	if errState != nil {
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
	// Wait, catching any errors
	_, errState := cloudapi.GetStateChangeConf(ss.Meta, ss.D, apiResponse.Header.Get("Location"), schema.TimeoutCreate).WaitForStateContext(ctx)
	if errState != nil {
		if cloudapi.IsRequestFailed(err) {
			// Request failed, so resource was not created, delete resource from state file
			ss.D.SetId("")
		}
		return nil, apiResponse, fmt.Errorf("an error occured while waiting for server state change on create dcId: %s, Response: (%w)", datacenterID, errState)
	}
	return &server, apiResponse, nil
}

func (ss *Service) Update(ctx context.Context, datacenterID, serverID string, serverProperties ionoscloud.ServerProperties) (*ionoscloud.Server, *ionoscloud.APIResponse, error) {
	updatedServer, apiResponse, err := ss.Client.ServersApi.DatacentersServersPatch(ctx, datacenterID, serverID).Execute()
	apiResponse.LogInfo()
	if err != nil {
		return nil, apiResponse, fmt.Errorf("an error occured while updating server for dcId: %s, server_id: %s, Response: (%w)", datacenterID, serverID, err)
	}
	// Wait, catching any errors
	_, errState := cloudapi.GetStateChangeConf(ss.Meta, ss.D, apiResponse.Header.Get("Location"), schema.TimeoutUpdate).WaitForStateContext(ctx)
	if errState != nil {
		return nil, apiResponse, fmt.Errorf("an error occured while waiting for server state change on update dcId: %s, server_id: %s, Response: (%w)", datacenterID, serverID, errState)
	}
	return &updatedServer, apiResponse, nil
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
	case ENTERPRISE_SERVER_TYPE:
		if strings.EqualFold(newVmState, CUBE_VMSTATE_STOP) {
			return fmt.Errorf("cannot suspend an enterprise server, set to %s instead", ENTERPRISE_VMSTATE_STOP)
		}
		if strings.EqualFold(newVmState, ENTERPRISE_VMSTATE_START) && strings.EqualFold(currentVmState, ENTERPRISE_VMSTATE_STOP) {
			return ss.Start(ctx, datacenterID, serverID, serverType)
		}
		if strings.EqualFold(newVmState, ENTERPRISE_VMSTATE_STOP) && strings.EqualFold(currentVmState, ENTERPRISE_VMSTATE_START) {
			return ss.Stop(ctx, datacenterID, serverID, serverType)
		}

	case CUBE_SERVER_TYPE:
		if strings.EqualFold(newVmState, ENTERPRISE_VMSTATE_STOP) {
			return fmt.Errorf("cannot shut down a cube server, set to %s instead", CUBE_VMSTATE_STOP)
		}
		if newVmState == "" {
			return fmt.Errorf("cannot update a suspended Cube Server, must change the state to %s first", CUBE_VMSTATE_START)
		}
		if strings.EqualFold(newVmState, CUBE_VMSTATE_START) && strings.EqualFold(currentVmState, CUBE_VMSTATE_STOP) {
			return ss.Start(ctx, datacenterID, serverID, serverType)
		}
		if strings.EqualFold(newVmState, CUBE_VMSTATE_STOP) && strings.EqualFold(currentVmState, CUBE_VMSTATE_START) {
			return ErrSuspendCubeLast
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

	case ENTERPRISE_SERVER_TYPE:
		apiResponse, err := ss.Client.ServersApi.DatacentersServersStartPost(ctx, datacenterID, serverID).Execute()
		apiResponse.LogInfo()
		if err != nil {
			return err
		}
		return utils.WaitForResourceToBeReady(ctx, ss.D, ss.checkExpectedVmStateFn(ctx, datacenterID, ENTERPRISE_VMSTATE_START))

	case CUBE_SERVER_TYPE:
		apiResponse, err := ss.Client.ServersApi.DatacentersServersResumePost(ctx, datacenterID, serverID).Execute()
		apiResponse.LogInfo()
		if err != nil {
			return err
		}
		return utils.WaitForResourceToBeReady(ctx, ss.D, ss.checkExpectedVmStateFn(ctx, datacenterID, CUBE_VMSTATE_START))
	}

	return fmt.Errorf("cannot start unknown server type: %s", serverType)

}

func (ss *Service) Stop(ctx context.Context, datacenterID, serverID, serverType string) error {

	switch serverType {

	case ENTERPRISE_SERVER_TYPE:
		apiResponse, err := ss.Client.ServersApi.DatacentersServersStopPost(ctx, datacenterID, serverID).Execute()
		apiResponse.LogInfo()
		if err != nil {
			return err
		}
		return utils.WaitForResourceToBeReady(ctx, ss.D, ss.checkExpectedVmStateFn(ctx, datacenterID, ENTERPRISE_VMSTATE_STOP))

	case CUBE_SERVER_TYPE:
		apiResponse, err := ss.Client.ServersApi.DatacentersServersSuspendPost(ctx, datacenterID, serverID).Execute()
		apiResponse.LogInfo()
		if err != nil {
			return err
		}
		return utils.WaitForResourceToBeReady(ctx, ss.D, ss.checkExpectedVmStateFn(ctx, datacenterID, CUBE_VMSTATE_STOP))
	}

	return fmt.Errorf("cannot stop unknown server type: %s", serverType)

}

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
