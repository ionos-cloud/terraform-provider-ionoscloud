package dataplatform

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	dataplatform "github.com/ionos-cloud/sdk-go-dataplatform"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils"
	"log"
	"strings"
	"time"
)

var clusterResourceName = "Dataplatform Cluster"

type ClusterService interface {
	GetById(ctx context.Context, id string) (dataplatform.ClusterResponseData, *dataplatform.APIResponse, error)
	GetClusterKubeConfig(ctx context.Context, clusterId string) (string, *dataplatform.APIResponse, error)
	ListClusters(ctx context.Context, filterName string) ([]dataplatform.ClusterResponseData, *dataplatform.APIResponse, error)
}

func (c *Client) GetClusterById(ctx context.Context, id string) (dataplatform.ClusterResponseData, *dataplatform.APIResponse, error) {
	cluster, apiResponse, err := c.DataPlatformClusterApi.GetCluster(ctx, id).Execute()
	apiResponse.LogInfo()
	return cluster, apiResponse, err
}

// DoesResourceExist - returns apiResponse to check if resource still exists. To be used with WaitForResourceToBeDeleted
func (c *Client) DoesResourceExist(ctx context.Context, id string) (utils.ApiResponseInfo, error) {
	_, apiResponse, err := c.DataPlatformClusterApi.GetCluster(ctx, id).Execute()
	apiResponse.LogInfo()
	return apiResponse, err
}

func (c *Client) GetClusterKubeConfig(ctx context.Context, clusterId string) (string, *dataplatform.APIResponse, error) {
	kubeConfig, apiResponse, err := c.DataPlatformClusterApi.GetClusterKubeconfig(ctx, clusterId).Execute()
	apiResponse.LogInfo()
	return kubeConfig, apiResponse, err
}

func (c *Client) ListClusters(ctx context.Context, filterName string) (dataplatform.ClusterListResponseData, *dataplatform.APIResponse, error) {
	request := c.DataPlatformClusterApi.GetClusters(ctx)
	if filterName != "" {
		request = request.Name(filterName)
	}
	clusters, apiResponse, err := c.DataPlatformClusterApi.GetClustersExecute(request)
	apiResponse.LogInfo()
	return clusters, apiResponse, err
}

// CreateResource - creates the request from the schema and sends it to the API. returns the id of the create resource,
// the apiResponse, or an error if an error occurred
func (c *Client) CreateResource(ctx context.Context, d *schema.ResourceData) (id string, responseInfo utils.ApiResponseInfo, err error) {
	createRequest := setCreateClusterRequestProperties(d)
	clusterResponse, apiResponse, err := c.DataPlatformClusterApi.CreateCluster(ctx).CreateClusterRequest(*createRequest).Execute()
	apiResponse.LogInfo()
	if err != nil {
		return "", apiResponse, err
	}
	return *clusterResponse.Id, apiResponse, err
}

func (c *Client) UpdateResource(ctx context.Context, id string, d *schema.ResourceData) (utils.ApiResponseInfo, error) {
	cluster := setPatchClusterRequestProperties(d)
	_, apiResponse, err := c.DataPlatformClusterApi.PatchCluster(ctx, id).PatchClusterRequest(*cluster).Execute()
	apiResponse.LogInfo()
	return apiResponse, err
}

func (c *Client) DeleteResource(ctx context.Context, id string) (utils.ApiResponseInfo, error) {
	_, apiResponse, err := c.DataPlatformClusterApi.DeleteCluster(ctx, id).Execute()
	apiResponse.LogInfo()
	return apiResponse, err
}

// WaitForClusterToBeReady - keeps retrying until cluster is in 'available' state, or context deadline is reached
func (c *Client) WaitForClusterToBeReady(ctx context.Context, clusterId string) error {
	var clusterRequest = dataplatform.NewClusterResponseDataWithDefaults()
	err := resource.RetryContext(ctx, 60*time.Minute, func() *resource.RetryError {
		var err error
		var apiResponse *dataplatform.APIResponse
		*clusterRequest, apiResponse, err = c.GetClusterById(ctx, clusterId)
		if apiResponse.HttpNotFound() {
			log.Printf("[INFO] Could not find cluster %s retrying...", clusterId)
			return resource.RetryableError(fmt.Errorf("could not find cluster %s, %w, retrying", clusterId, err))
		}
		if err != nil {
			resource.NonRetryableError(err)
		}

		if clusterRequest != nil && clusterRequest.Metadata != nil && !strings.EqualFold(*clusterRequest.Metadata.State, utils.Available) {
			log.Printf("[INFO] dataplatform cluster %s is still in state %s", clusterId, *clusterRequest.Metadata.State)
			return resource.RetryableError(fmt.Errorf(" dataplatform cluster is still in state %s", *clusterRequest.Metadata.State))
		}
		return nil
	})
	if clusterRequest == nil || clusterRequest.Properties == nil || *clusterRequest.Properties.DatacenterId == "" {
		return fmt.Errorf("could not find dataplatform cluster %s", clusterId)
	}
	return err
}

func setCreateClusterRequestProperties(d *schema.ResourceData) *dataplatform.CreateClusterRequest {

	dataplatformCluster := dataplatform.CreateClusterRequest{
		Properties: &dataplatform.CreateClusterProperties{},
	}

	if nameValue, ok := d.GetOk("name"); ok {
		name := nameValue.(string)
		dataplatformCluster.Properties.Name = &name
	}

	if datacenterIdValue, ok := d.GetOk("datacenter_id"); ok {
		datacenterId := datacenterIdValue.(string)
		dataplatformCluster.Properties.DatacenterId = &datacenterId
	}

	if dataPlatformVersionValue, ok := d.GetOk("data_platform_version"); ok {
		dataPlatformVersion := dataPlatformVersionValue.(string)
		dataplatformCluster.Properties.DataPlatformVersion = &dataPlatformVersion
	}

	if _, ok := d.GetOk("maintenance_window"); ok {
		dataplatformCluster.Properties.MaintenanceWindow = setMaintenanceWindowData(d)
	}

	return &dataplatformCluster
}

func setPatchClusterRequestProperties(d *schema.ResourceData) *dataplatform.PatchClusterRequest {

	dataplatformCluster := dataplatform.PatchClusterRequest{
		Properties: &dataplatform.PatchClusterProperties{},
	}

	if nameValue, ok := d.GetOk("name"); ok {
		name := nameValue.(string)
		dataplatformCluster.Properties.Name = &name
	}

	if dataPlatformVersionValue, ok := d.GetOk("data_platform_version"); ok {
		dataPlatformVersion := dataPlatformVersionValue.(string)
		dataplatformCluster.Properties.DataPlatformVersion = &dataPlatformVersion
	}

	if _, ok := d.GetOk("maintenance_window"); ok {
		dataplatformCluster.Properties.MaintenanceWindow = setMaintenanceWindowData(d)
	}

	return &dataplatformCluster
}

func setMaintenanceWindowData(d *schema.ResourceData) *dataplatform.MaintenanceWindow {
	var maintenanceWindow dataplatform.MaintenanceWindow

	if timeV, ok := d.GetOk("maintenance_window.0.time"); ok {
		timeV := timeV.(string)
		maintenanceWindow.Time = &timeV
	}

	if dayOfTheWeek, ok := d.GetOk("maintenance_window.0.day_of_the_week"); ok {
		dayOfTheWeek := dayOfTheWeek.(string)
		maintenanceWindow.DayOfTheWeek = &dayOfTheWeek
	}

	return &maintenanceWindow
}

func SetDataplatformClusterData(d *schema.ResourceData, cluster dataplatform.ClusterResponseData) error {

	if cluster.Id != nil {
		d.SetId(*cluster.Id)
	}

	if cluster.Properties.Name != nil {
		if err := d.Set("name", *cluster.Properties.Name); err != nil {
			return utils.GenerateSetError(clusterResourceName, "name", err)
		}
	}

	if cluster.Properties.DataPlatformVersion != nil {
		if err := d.Set("data_platform_version", *cluster.Properties.DataPlatformVersion); err != nil {
			return utils.GenerateSetError(clusterResourceName, "data_platform_version", err)
		}
	}

	if cluster.Properties.DatacenterId != nil {
		if err := d.Set("datacenter_id", *cluster.Properties.DatacenterId); err != nil {
			return utils.GenerateSetError(clusterResourceName, "datacenter_id", err)
		}
	}

	if cluster.Properties.MaintenanceWindow != nil {
		var maintenanceWindow []interface{}
		maintenanceWindowEntry := SetMaintenanceWindowProperties(*cluster.Properties.MaintenanceWindow)
		maintenanceWindow = append(maintenanceWindow, maintenanceWindowEntry)
		if err := d.Set("maintenance_window", maintenanceWindow); err != nil {
			return utils.GenerateSetError(clusterResourceName, "maintenance_window", err)
		}
	}

	return nil
}

func SetMaintenanceWindowProperties(maintenanceWindow dataplatform.MaintenanceWindow) map[string]interface{} {

	maintenance := map[string]interface{}{}

	utils.SetPropWithNilCheck(maintenance, "time", maintenanceWindow.Time)
	utils.SetPropWithNilCheck(maintenance, "day_of_the_week", maintenanceWindow.DayOfTheWeek)

	return maintenance
}
