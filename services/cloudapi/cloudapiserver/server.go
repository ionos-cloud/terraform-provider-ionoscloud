package cloudapiserver

import (
	"context"
	"fmt"
	"log"

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

func (fs *Service) FindById(ctx context.Context, datacenterID, serverID string, depth int32) (*ionoscloud.Server, error) {
	server, apiResponse, err := fs.Client.ServersApi.DatacentersServersFindById(ctx, datacenterID, serverID).Depth(depth).Execute()
	apiResponse.LogInfo()
	if err != nil {
		return nil, err
	}
	return &server, nil
}

func (fs *Service) Delete(ctx context.Context, datacenterID, serverID, ID string) (*ionoscloud.APIResponse, error) {
	apiResponse, err := fs.Client.ServersApi.DatacentersServersDelete(ctx, datacenterID, serverID).Execute()
	apiResponse.LogInfo()
	if err != nil {
		return apiResponse, err
	}
	// Wait, catching any errors
	_, errState := cloudapi.GetStateChangeConf(fs.Meta, fs.D, apiResponse.Header.Get("Location"), schema.TimeoutDelete).WaitForStateContext(ctx)
	if errState != nil {
		return apiResponse, fmt.Errorf("an error occured while waiting for server state change on delete dcId: %s, server_id: %s, ID: %s, Response: (%w)", datacenterID, serverID, ID, errState)
	}
	return apiResponse, nil
}

func (fs *Service) Create(ctx context.Context, datacenterID string) (*ionoscloud.Server, *ionoscloud.APIResponse, error) {
	server, apiResponse, err := fs.Client.ServersApi.DatacentersServersPost(ctx, datacenterID).Execute()
	apiResponse.LogInfo()
	if err != nil {
		return nil, apiResponse, fmt.Errorf("an error occured while creating server for dcId: %s, Response: (%w)", datacenterID, err)
	}
	// Wait, catching any errors
	_, errState := cloudapi.GetStateChangeConf(fs.Meta, fs.D, apiResponse.Header.Get("Location"), schema.TimeoutCreate).WaitForStateContext(ctx)
	if errState != nil {
		if cloudapi.IsRequestFailed(err) {
			// Request failed, so resource was not created, delete resource from state file
			fs.D.SetId("")
		}
		return nil, apiResponse, fmt.Errorf("an error occured while waiting for server state change on create dcId: %s, Response: (%w)", datacenterID, errState)
	}
	return &server, apiResponse, nil
}

func (fs *Service) Update(ctx context.Context, datacenterID, serverID string, serverProperties ionoscloud.ServerProperties) (*ionoscloud.Server, *ionoscloud.APIResponse, error) {
	updatedServer, apiResponse, err := fs.Client.ServersApi.DatacentersServersPatch(ctx, datacenterID, serverID).Execute()
	apiResponse.LogInfo()
	if err != nil {
		return nil, apiResponse, fmt.Errorf("an error occured while updating server for dcId: %s, server_id: %s, Response: (%w)", datacenterID, serverID, err)
	}
	// Wait, catching any errors
	_, errState := cloudapi.GetStateChangeConf(fs.Meta, fs.D, apiResponse.Header.Get("Location"), schema.TimeoutUpdate).WaitForStateContext(ctx)
	if errState != nil {
		return nil, apiResponse, fmt.Errorf("an error occured while waiting for server state change on update dcId: %s, server_id: %s, Response: (%w)", datacenterID, serverID, errState)
	}
	return &updatedServer, apiResponse, nil
}

func (fs *Service) Start(ctx context.Context, datacenterID, serverID, serverType string) error {

	switch serverType {

	case "ENTERPRISE":
		apiResponse, err := fs.Client.ServersApi.DatacentersServersStartPost(ctx, datacenterID, serverID).Execute()
		apiResponse.LogInfo()
		if err != nil {
			return err
		}
		return utils.WaitForResourceToBeReady(ctx, fs.D, fs.checkExpectedVmStateFn(ctx, datacenterID, "RUNNING"))

	case "CUBE":
		apiResponse, err := fs.Client.ServersApi.DatacentersServersResumePost(ctx, datacenterID, serverID).Execute()
		apiResponse.LogInfo()
		if err != nil {
			return err
		}
		return utils.WaitForResourceToBeReady(ctx, fs.D, fs.checkExpectedVmStateFn(ctx, datacenterID, "RUNNING"))
	}

	return fmt.Errorf("unknown server type: %s", serverType)

}

func (fs *Service) Stop(ctx context.Context, datacenterID, serverID, serverType string) error {

	switch serverType {

	case "ENTERPRISE":
		apiResponse, err := fs.Client.ServersApi.DatacentersServersStopPost(ctx, datacenterID, serverID).Execute()
		apiResponse.LogInfo()
		if err != nil {
			return err
		}
		return utils.WaitForResourceToBeReady(ctx, fs.D, fs.checkExpectedVmStateFn(ctx, datacenterID, "SHUTOFF"))

	case "CUBE":
		apiResponse, err := fs.Client.ServersApi.DatacentersServersSuspendPost(ctx, datacenterID, serverID).Execute()
		apiResponse.LogInfo()
		if err != nil {
			return err
		}
		return utils.WaitForResourceToBeReady(ctx, fs.D, fs.checkExpectedVmStateFn(ctx, datacenterID, "SUSPENDED"))
	}

	return fmt.Errorf("unknown server type: %s", serverType)

}

func (fs *Service) checkExpectedVmStateFn(ctx context.Context, dcId, expectedState string) func(ctx context.Context, d *schema.ResourceData) (bool, error) {

	return func(ctx context.Context, d *schema.ResourceData) (bool, error) {
		ionoscloudServer, _, err := fs.Client.ServersApi.DatacentersServersFindById(ctx, dcId, d.Id()).Execute()
		serverType := *ionoscloudServer.Properties.Type
		if err != nil {
			return false, err
		}
		if *ionoscloudServer.Properties.VmState != expectedState {
			log.Printf("[INFO] Server (type: %s) vmState not yet changed to %s: %s", serverType, expectedState, d.Id())
			return false, nil
		}
		return true, nil
	}
}
