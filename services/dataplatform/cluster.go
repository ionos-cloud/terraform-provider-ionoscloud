package dataplatform

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	dataplatform "github.com/ionos-cloud/sdk-go-autoscaling"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils"
)

var clusterResourceName = "Dataplatform Cluster"

type ClusterService interface {
	GetCluster(ctx context.Context, clusterId string) (dataplatform.ClusterResponseData, *dataplatform.APIResponse, error)
	GetClusterKubeConfig(ctx context.Context, clusterId string) (string, *dataplatform.APIResponse, error)
	ListClusters(ctx context.Context, filterName string) ([]dataplatform.ClusterResponseData, *dataplatform.APIResponse, error)
	CreateCluster(ctx context.Context, cluster dataplatform.CreateClusterRequest) (dataplatform.ClusterResponseData, *dataplatform.APIResponse, error)
	UpdateCluster(ctx context.Context, clusterId string, cluster dataplatform.PatchClusterRequest) (dataplatform.ClusterResponseData, *dataplatform.APIResponse, error)
	DeleteCluster(ctx context.Context, clusterId string) (dataplatform.ClusterResponseData, *dataplatform.APIResponse, error)
}

func (c *Client) GetCluster(ctx context.Context, clusterId string) (dataplatform.ClusterResponseData, *dataplatform.APIResponse, error) {
	cluster, apiResponse, err := c.DataPlatformClusterApi.GetCluster(ctx, clusterId).Execute()
	if apiResponse != nil {
		return cluster, apiResponse, err

	}
	return cluster, nil, err
}

func (c *Client) GetClusterKubeConfig(ctx context.Context, clusterId string) (string, *dataplatform.APIResponse, error) {
	kubeConfig, apiResponse, err := c.DataPlatformClusterApi.GetClusterKubeconfig(ctx, clusterId).Execute()
	if apiResponse != nil {
		return kubeConfig, apiResponse, err
	}
	return kubeConfig, nil, err
}

func (c *Client) ListClusters(ctx context.Context, filterName string) (dataplatform.ClusterListResponseData, *dataplatform.APIResponse, error) {
	request := c.DataPlatformClusterApi.GetClusters(ctx)
	if filterName != "" {
		request = request.Name(filterName)
	}
	clusters, apiResponse, err := c.DataPlatformClusterApi.GetClustersExecute(request)
	if apiResponse != nil {
		return clusters, apiResponse, err
	}
	return clusters, nil, err
}

func (c *Client) CreateCluster(ctx context.Context, cluster dataplatform.CreateClusterRequest) (dataplatform.ClusterResponseData, *dataplatform.APIResponse, error) {
	clusterResponse, apiResponse, err := c.DataPlatformClusterApi.CreateCluster(ctx).CreateClusterRequest(cluster).Execute()
	if apiResponse != nil {
		return clusterResponse, apiResponse, err
	}
	return clusterResponse, nil, err
}

func (c *Client) UpdateCluster(ctx context.Context, clusterId string, cluster dataplatform.PatchClusterRequest) (dataplatform.ClusterResponseData, *dataplatform.APIResponse, error) {
	clusterResponse, apiResponse, err := c.DataPlatformClusterApi.PatchCluster(ctx, clusterId).PatchClusterRequest(cluster).Execute()
	if apiResponse != nil {
		return clusterResponse, apiResponse, err
	}
	return clusterResponse, nil, err
}

func (c *Client) DeleteCluster(ctx context.Context, clusterId string) (dataplatform.ClusterResponseData, *dataplatform.APIResponse, error) {
	clusterResponse, apiResponse, err := c.DataPlatformClusterApi.DeleteCluster(ctx, clusterId).Execute()
	if apiResponse != nil {
		return clusterResponse, apiResponse, err
	}
	return clusterResponse, nil, err
}

func GetDataplatformClusterDataCreate(d *schema.ResourceData) *dataplatform.CreateClusterRequest {

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
		dataplatformCluster.Properties.MaintenanceWindow = GetDataplatformMaintenanceWindowData(d)
	}

	return &dataplatformCluster
}

func GetDataplatformClusterDataUpdate(d *schema.ResourceData) (*dataplatform.PatchClusterRequest, diag.Diagnostics) {

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
		dataplatformCluster.Properties.MaintenanceWindow = GetDataplatformMaintenanceWindowData(d)
	}

	return &dataplatformCluster, nil
}

func GetDataplatformMaintenanceWindowData(d *schema.ResourceData) *dataplatform.MaintenanceWindow {
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
