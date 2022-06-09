package dsaas

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	dsaas "github.com/ionos-cloud/sdk-go-autoscaling"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils"
)

var clusterResourceName = "DSaaS Cluster"

type ClusterService interface {
	GetCluster(ctx context.Context, clusterId string) (dsaas.ClusterResponseData, *dsaas.APIResponse, error)
	ListClusters(ctx context.Context, filterName string) ([]dsaas.ClusterResponseData, *dsaas.APIResponse, error)
	CreateCluster(ctx context.Context, cluster dsaas.CreateClusterRequest) (dsaas.ClusterResponseData, *dsaas.APIResponse, error)
	UpdateCluster(ctx context.Context, clusterId string, cluster dsaas.PatchClusterRequest) (dsaas.ClusterResponseData, *dsaas.APIResponse, error)
	DeleteCluster(ctx context.Context, clusterId string) (dsaas.ClusterResponseData, *dsaas.APIResponse, error)
}

func (c *Client) GetCluster(ctx context.Context, clusterId string) (dsaas.ClusterResponseData, *dsaas.APIResponse, error) {
	cluster, apiResponse, err := c.DataPlatformClusterApi.GetCluster(ctx, clusterId).Execute()
	if apiResponse != nil {
		return cluster, apiResponse, err

	}
	return cluster, nil, err
}

func (c *Client) ListClusters(ctx context.Context, filterName string) (dsaas.ClusterListResponseData, *dsaas.APIResponse, error) {
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

func (c *Client) CreateCluster(ctx context.Context, cluster dsaas.CreateClusterRequest) (dsaas.ClusterResponseData, *dsaas.APIResponse, error) {
	clusterResponse, apiResponse, err := c.DataPlatformClusterApi.CreateCluster(ctx).CreateClusterRequest(cluster).Execute()
	if apiResponse != nil {
		return clusterResponse, apiResponse, err
	}
	return clusterResponse, nil, err
}

func (c *Client) UpdateCluster(ctx context.Context, clusterId string, cluster dsaas.PatchClusterRequest) (dsaas.ClusterResponseData, *dsaas.APIResponse, error) {
	clusterResponse, apiResponse, err := c.DataPlatformClusterApi.PatchCluster(ctx, clusterId).PatchClusterRequest(cluster).Execute()
	if apiResponse != nil {
		return clusterResponse, apiResponse, err
	}
	return clusterResponse, nil, err
}

func (c *Client) DeleteCluster(ctx context.Context, clusterId string) (dsaas.ClusterResponseData, *dsaas.APIResponse, error) {
	clusterResponse, apiResponse, err := c.DataPlatformClusterApi.DeleteCluster(ctx, clusterId).Execute()
	if apiResponse != nil {
		return clusterResponse, apiResponse, err
	}
	return clusterResponse, nil, err
}

func GetDSaaSClusterDataCreate(d *schema.ResourceData) *dsaas.CreateClusterRequest {

	dsaasCluster := dsaas.CreateClusterRequest{
		Properties: &dsaas.CreateClusterProperties{},
	}

	if nameValue, ok := d.GetOk("name"); ok {
		name := nameValue.(string)
		dsaasCluster.Properties.Name = &name
	}

	if datacenterIdValue, ok := d.GetOk("datacenter_id"); ok {
		datacenterId := datacenterIdValue.(string)
		dsaasCluster.Properties.DatacenterId = &datacenterId
	}

	if dataPlatformVersionValue, ok := d.GetOk("data_platform_version"); ok {
		dataPlatformVersion := dataPlatformVersionValue.(string)
		dsaasCluster.Properties.DataPlatformVersion = &dataPlatformVersion
	}

	if _, ok := d.GetOk("maintenance_window"); ok {
		dsaasCluster.Properties.MaintenanceWindow = GetDSaaSMaintenanceWindowData(d)
	}

	return &dsaasCluster
}

func GetDSaaSClusterDataUpdate(d *schema.ResourceData) (*dsaas.PatchClusterRequest, diag.Diagnostics) {

	dsaasCluster := dsaas.PatchClusterRequest{
		Properties: &dsaas.PatchClusterProperties{},
	}

	if nameValue, ok := d.GetOk("name"); ok {
		name := nameValue.(string)
		dsaasCluster.Properties.Name = &name
	}

	if dataPlatformVersionValue, ok := d.GetOk("data_platform_version"); ok {
		dataPlatformVersion := dataPlatformVersionValue.(string)
		dsaasCluster.Properties.DataPlatformVersion = &dataPlatformVersion
	}

	if _, ok := d.GetOk("maintenance_window"); ok {
		dsaasCluster.Properties.MaintenanceWindow = GetDSaaSMaintenanceWindowData(d)
	}

	if _, ok := d.GetOk("datacenter_id"); ok {
		return nil, diag.FromErr(utils.GenerateImmutableError(clusterResourceName, "datacenter_id"))
	}

	return &dsaasCluster, nil
}

func GetDSaaSMaintenanceWindowData(d *schema.ResourceData) *dsaas.MaintenanceWindow {
	var maintenanceWindow dsaas.MaintenanceWindow

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

func SetDSaaSClusterData(d *schema.ResourceData, cluster dsaas.ClusterResponseData) error {

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

func SetMaintenanceWindowProperties(maintenanceWindow dsaas.MaintenanceWindow) map[string]interface{} {

	maintenance := map[string]interface{}{}

	utils.SetPropWithNilCheck(maintenance, "time", maintenanceWindow.Time)
	utils.SetPropWithNilCheck(maintenance, "day_of_the_week", maintenanceWindow.DayOfTheWeek)

	return maintenance
}
