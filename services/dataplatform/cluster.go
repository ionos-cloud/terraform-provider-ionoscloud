package dataplatform

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/ionos-cloud/sdk-go-bundle/products/dataplatform/v2"
	"github.com/ionos-cloud/sdk-go-bundle/shared"

	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils/constant"
)

var clusterResourceName = "Dataplatform Cluster"

func (c *Client) IsClusterReady(ctx context.Context, d *schema.ResourceData) (bool, error) {
	cluster, _, err := c.GetClusterById(ctx, d.Id())
	if err != nil {
		return false, err
	}
	if cluster.Metadata.State == nil {
		return false, fmt.Errorf("expected metadata, got empty for cluster id %s", d.Id())
	}
	if utils.IsStateFailed(*cluster.Metadata.State) {
		return false, fmt.Errorf("cluster %s is in a failed state", d.Id())
	}
	log.Printf("[DEBUG] dataplatform cluster state %s", *cluster.Metadata.State)
	return strings.EqualFold(*cluster.Metadata.State, constant.Available), nil
}

func (c *Client) GetClusterById(ctx context.Context, id string) (dataplatform.ClusterResponseData, *shared.APIResponse, error) {
	cluster, apiResponse, err := c.sdkClient.DataPlatformClusterApi.ClustersFindById(ctx, id).Execute()
	apiResponse.LogInfo()
	return cluster, apiResponse, err
}

// IsClusterDeleted - checks if resource still exists. To be used with WaitForResourceToBeDeleted
func (c *Client) IsClusterDeleted(ctx context.Context, d *schema.ResourceData) (bool, error) {
	_, apiResponse, err := c.sdkClient.DataPlatformClusterApi.ClustersFindById(ctx, d.Id()).Execute()
	apiResponse.LogInfo()
	return apiResponse.HttpNotFound(), err
}

// GetClusterKubeConfig - gets the kube config for the cluster
func (c *Client) GetClusterKubeConfig(ctx context.Context, clusterID string) (map[string]interface{}, *shared.APIResponse, error) {
	kubeConfig, apiResponse, err := c.sdkClient.DataPlatformClusterApi.ClustersKubeconfigFindByClusterId(ctx, clusterID).Execute()
	apiResponse.LogInfo()
	return kubeConfig, apiResponse, err
}

func (c *Client) ListClusters(ctx context.Context, filterName string) (dataplatform.ClusterListResponseData, *shared.APIResponse, error) {
	request := c.sdkClient.DataPlatformClusterApi.ClustersGet(ctx)
	if filterName != "" {
		request = request.Name(filterName)
	}
	clusters, apiResponse, err := c.sdkClient.DataPlatformClusterApi.ClustersGetExecute(request)
	apiResponse.LogInfo()
	return clusters, apiResponse, err
}

// CreateCluster - creates the request from the schema and sends it to the API. returns the id of the created resource,
// the apiResponse, or an error if an error occurred
func (c *Client) CreateCluster(ctx context.Context, d *schema.ResourceData) (id string, responseInfo utils.ApiResponseInfo, err error) {
	createRequest := setCreateClusterRequestProperties(d)
	clusterResponse, apiResponse, err := c.sdkClient.DataPlatformClusterApi.ClustersPost(ctx).CreateClusterRequest(*createRequest).Execute()
	apiResponse.LogInfo()
	if err != nil {
		return "", apiResponse, err
	}
	return clusterResponse.Id, apiResponse, err
}

func (c *Client) UpdateCluster(ctx context.Context, id string, d *schema.ResourceData) (utils.ApiResponseInfo, error) {
	cluster := setPatchClusterRequestProperties(d)
	_, apiResponse, err := c.sdkClient.DataPlatformClusterApi.ClustersPatch(ctx, id).PatchClusterRequest(*cluster).Execute()
	apiResponse.LogInfo()
	return apiResponse, err
}

func (c *Client) DeleteCluster(ctx context.Context, id string) (utils.ApiResponseInfo, error) {
	_, apiResponse, err := c.sdkClient.DataPlatformClusterApi.ClustersDelete(ctx, id).Execute()
	apiResponse.LogInfo()
	return apiResponse, err
}

func setCreateClusterRequestProperties(d *schema.ResourceData) *dataplatform.CreateClusterRequest {

	dataplatformCluster := dataplatform.CreateClusterRequest{
		Properties: dataplatform.CreateClusterProperties{},
	}

	if nameValue, ok := d.GetOk("name"); ok {
		name := nameValue.(string)
		dataplatformCluster.Properties.Name = name
	}

	if datacenterIdValue, ok := d.GetOk("datacenter_id"); ok {
		datacenterId := datacenterIdValue.(string)
		dataplatformCluster.Properties.DatacenterId = datacenterId
	}

	if dataPlatformVersionValue, ok := d.GetOk("version"); ok {
		dataPlatformVersion := dataPlatformVersionValue.(string)
		dataplatformCluster.Properties.DataPlatformVersion = &dataPlatformVersion
	}

	if _, ok := d.GetOk("maintenance_window"); ok {
		dataplatformCluster.Properties.MaintenanceWindow = setMaintenanceWindowData(d)
	}

	dataplatformCluster.Properties.Lans = *setLansData(d)

	return &dataplatformCluster
}

func setPatchClusterRequestProperties(d *schema.ResourceData) *dataplatform.PatchClusterRequest {

	dataplatformCluster := dataplatform.PatchClusterRequest{
		Properties: dataplatform.PatchClusterProperties{},
	}

	if nameValue, ok := d.GetOk("name"); ok {
		name := nameValue.(string)
		dataplatformCluster.Properties.Name = &name
	}

	if dataPlatformVersionValue, ok := d.GetOk("version"); ok {
		dataPlatformVersion := dataPlatformVersionValue.(string)
		dataplatformCluster.Properties.DataPlatformVersion = &dataPlatformVersion
	}

	if _, ok := d.GetOk("maintenance_window"); ok {
		dataplatformCluster.Properties.MaintenanceWindow = setMaintenanceWindowData(d)
	}

	dataplatformCluster.Properties.Lans = *setLansData(d)

	return &dataplatformCluster
}

func setLansData(d *schema.ResourceData) *[]dataplatform.Lan {
	lansBody := make([]dataplatform.Lan, 0)
	if lansData, ok := d.GetOk("lans"); ok {
		if lansData, ok := lansData.(*schema.Set); ok {
			for _, lanData := range lansData.List() {
				if lan, ok := lanData.(map[string]interface{}); ok {
					lanID := lan["lan_id"].(string)
					dhcp := lan["dhcp"].(bool)
					var lanBody dataplatform.Lan
					lanBody.LanId = lanID
					lanBody.Dhcp = &dhcp
					lanBody.Routes = *setRoutesData(lan)
					lansBody = append(lansBody, lanBody)
				}
			}
		}
	}
	return &lansBody
}

func setRoutesData(lan map[string]interface{}) *[]dataplatform.Route {
	routesBody := make([]dataplatform.Route, 0)
	if routesData, ok := lan["routes"].(*schema.Set); ok {
		for _, routeData := range routesData.List() {
			if route, ok := routeData.(map[string]interface{}); ok {
				network := route["network"].(string)
				gateway := route["gateway"].(string)
				var routeBody dataplatform.Route
				routeBody.Network = network
				routeBody.Gateway = gateway
				routesBody = append(routesBody, routeBody)
			}
		}
	}
	return &routesBody
}

func setMaintenanceWindowData(d *schema.ResourceData) *dataplatform.MaintenanceWindow {
	var maintenanceWindow dataplatform.MaintenanceWindow

	if timeV, ok := d.GetOk("maintenance_window.0.time"); ok {
		timeV := timeV.(string)
		maintenanceWindow.Time = timeV
	}

	if dayOfTheWeek, ok := d.GetOk("maintenance_window.0.day_of_the_week"); ok {
		dayOfTheWeek := dayOfTheWeek.(string)
		maintenanceWindow.DayOfTheWeek = dayOfTheWeek
	}

	return &maintenanceWindow
}

func SetDataplatformClusterData(d *schema.ResourceData, cluster dataplatform.ClusterResponseData) error {

	d.SetId(cluster.Id)

	if cluster.Properties.Name != nil {
		if err := d.Set("name", *cluster.Properties.Name); err != nil {
			return utils.GenerateSetError(clusterResourceName, "name", err)
		}
	}

	if cluster.Properties.DataPlatformVersion != nil {
		if err := d.Set("version", *cluster.Properties.DataPlatformVersion); err != nil {
			return utils.GenerateSetError(clusterResourceName, "version", err)
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

	if cluster.Properties.Lans != nil {
		if err := d.Set("lans", setLansProperties(cluster.Properties.Lans)); err != nil {
			return utils.GenerateSetError(clusterResourceName, "lans", err)
		}
	}

	return nil
}

func setLansProperties(retrievedLans []dataplatform.Lan) []interface{} {
	lans := make([]interface{}, len(retrievedLans))
	for i, lanElem := range retrievedLans {
		lanEntry := make(map[string]interface{})
		utils.SetPropWithNilCheck(lanEntry, "lan_id", lanElem.LanId)
		utils.SetPropWithNilCheck(lanEntry, "dhcp", lanElem.Dhcp)
		if lanElem.Routes != nil {
			routes := make([]interface{}, len(lanElem.Routes))
			for i, routeElem := range lanElem.Routes {
				routeEntry := make(map[string]interface{})
				utils.SetPropWithNilCheck(routeEntry, "network", routeElem.Network)
				utils.SetPropWithNilCheck(routeEntry, "gateway", routeElem.Gateway)
				routes[i] = routeEntry
			}
			utils.SetPropWithNilCheck(lanEntry, "routes", routes)
		}
		lans[i] = lanEntry
	}
	return lans
}

func SetMaintenanceWindowProperties(maintenanceWindow dataplatform.MaintenanceWindow) map[string]interface{} {

	maintenance := map[string]interface{}{}

	utils.SetPropWithNilCheck(maintenance, "time", maintenanceWindow.Time)
	utils.SetPropWithNilCheck(maintenance, "day_of_the_week", maintenanceWindow.DayOfTheWeek)

	return maintenance
}
