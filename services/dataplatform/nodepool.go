package dataplatform

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils/constant"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	dataplatform "github.com/ionos-cloud/sdk-go-dataplatform"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils"
)

var nodePoolResourceName = "Dataplatform Node Pool"

func (c *Client) GetNodePool(ctx context.Context, clusterId, nodePoolId string) (dataplatform.NodePoolResponseData, *dataplatform.APIResponse, error) {
	cluster, apiResponse, err := c.sdkClient.DataPlatformNodePoolApi.ClustersNodepoolsFindById(ctx, clusterId, nodePoolId).Execute()
	apiResponse.LogInfo()
	return cluster, apiResponse, err
}

func (c *Client) IsNodePoolDeleted(ctx context.Context, d *schema.ResourceData) (bool, error) {
	clusterId, ok := d.GetOk("cluster_id")
	if !ok {
		return false, fmt.Errorf("could not get cluster_id from schema")
	}
	clusterIdStr := clusterId.(string)
	nodePoolId := d.Id()
	_, apiResponse, err := c.sdkClient.DataPlatformNodePoolApi.ClustersNodepoolsFindById(ctx, clusterIdStr, nodePoolId).Execute()
	apiResponse.LogInfo()
	return apiResponse.HttpNotFound(), err
}

func (c *Client) ListNodePools(ctx context.Context, clusterId string) (dataplatform.NodePoolListResponseData, *dataplatform.APIResponse, error) {
	nodePool, apiResponse, err := c.sdkClient.DataPlatformNodePoolApi.ClustersNodepoolsGet(ctx, clusterId).Execute()
	apiResponse.LogInfo()
	return nodePool, apiResponse, err
}

func (c *Client) CreateNodePool(ctx context.Context, clusterId string, d *schema.ResourceData) (dataplatform.NodePoolResponseData, *dataplatform.APIResponse, error) {
	dataplatformNodePool := GetDataplatformNodePoolDataCreate(d)
	clusterResponse, apiResponse, err := c.sdkClient.DataPlatformNodePoolApi.ClustersNodepoolsPost(ctx, clusterId).CreateNodePoolRequest(*dataplatformNodePool).Execute()
	apiResponse.LogInfo()
	return clusterResponse, apiResponse, err
}

func (c *Client) UpdateNodePool(ctx context.Context, clusterId, nodePoolId string, d *schema.ResourceData) (dataplatform.NodePoolResponseData, utils.ApiResponseInfo, error) {
	dataplatformNodePool := GetDataplatformNodePoolDataUpdate(d)
	clusterResponse, apiResponse, err := c.sdkClient.DataPlatformNodePoolApi.ClustersNodepoolsPatch(ctx, clusterId, nodePoolId).PatchNodePoolRequest(*dataplatformNodePool).Execute()
	apiResponse.LogInfo()
	return clusterResponse, apiResponse, err
}

func (c *Client) DeleteNodePool(ctx context.Context, clusterId, nodePoolId string) (dataplatform.NodePoolResponseData, *dataplatform.APIResponse, error) {
	clusterResponse, apiResponse, err := c.sdkClient.DataPlatformNodePoolApi.ClustersNodepoolsDelete(ctx, clusterId, nodePoolId).Execute()
	apiResponse.LogInfo()
	return clusterResponse, apiResponse, err
}
func (c *Client) IsNodePoolReady(ctx context.Context, d *schema.ResourceData) (bool, error) {
	clusterId, ok := d.GetOk("cluster_id")
	if !ok {
		return false, fmt.Errorf("could not get cluster_id from schema")
	}
	clusterIdStr := clusterId.(string)
	subjectNodePool, _, err := c.GetNodePool(ctx, clusterIdStr, d.Id())
	if err != nil {
		return false, fmt.Errorf("checking Dataplatform Node Pool status: %w", err)
	}

	if subjectNodePool.Metadata == nil || subjectNodePool.Metadata.State == nil {
		return false, fmt.Errorf("expected nodepool metadata, got empty for id %s", d.Id())
	}
	log.Printf("[DEBUG] dataplatform cluster nodepool state %s", *subjectNodePool.Metadata.State)
	if strings.EqualFold(*subjectNodePool.Metadata.State, "FAILED") {
		return false, fmt.Errorf("nodepool id %s is in failed state", d.Id())
	}

	return strings.EqualFold(*subjectNodePool.Metadata.State, constant.Available), nil
}

func GetDataplatformNodePoolDataCreate(d *schema.ResourceData) *dataplatform.CreateNodePoolRequest {

	dataplatformNodePool := dataplatform.CreateNodePoolRequest{
		Properties: &dataplatform.CreateNodePoolProperties{},
	}

	if nameValue, ok := d.GetOk("name"); ok {
		name := nameValue.(string)
		dataplatformNodePool.Properties.Name = &name
	}

	if nodeCountValue, ok := d.GetOk("node_count"); ok {
		nodeCount := int32(nodeCountValue.(int))
		dataplatformNodePool.Properties.NodeCount = &nodeCount
	}

	if cpuFamilyValue, ok := d.GetOk("cpu_family"); ok {
		cpuFamily := cpuFamilyValue.(string)
		dataplatformNodePool.Properties.CpuFamily = &cpuFamily
	}

	if coresCountValue, ok := d.GetOk("cores_count"); ok {
		coresCount := int32(coresCountValue.(int))
		dataplatformNodePool.Properties.CoresCount = &coresCount
	}

	if ramSizeValue, ok := d.GetOk("ram_size"); ok {
		ramSize := int32(ramSizeValue.(int))
		dataplatformNodePool.Properties.RamSize = &ramSize
	}

	if availabilityZoneValue, ok := d.GetOk("availability_zone"); ok {
		availabilityZone := dataplatform.AvailabilityZone(availabilityZoneValue.(string))
		dataplatformNodePool.Properties.AvailabilityZone = &availabilityZone
	}

	if storageTypeValue, ok := d.GetOk("storage_type"); ok {
		storageType := dataplatform.StorageType(storageTypeValue.(string))
		dataplatformNodePool.Properties.StorageType = &storageType
	}

	if storageSizeValue, ok := d.GetOk("storage_size"); ok {
		storageSize := int32(storageSizeValue.(int))
		dataplatformNodePool.Properties.StorageSize = &storageSize
	}

	if _, ok := d.GetOk("maintenance_window"); ok {
		dataplatformNodePool.Properties.MaintenanceWindow = setMaintenanceWindowData(d)
	}

	if labelsValue, ok := d.GetOk("labels"); ok {
		labels := make(map[string]interface{})
		for k, v := range labelsValue.(map[string]interface{}) {
			labels[k] = v.(string)
		}
		dataplatformNodePool.Properties.Labels = &labels
	}

	if annotationsValue, ok := d.GetOk("annotations"); ok {
		annotations := make(map[string]interface{})
		for k, v := range annotationsValue.(map[string]interface{}) {
			annotations[k] = v.(string)
		}
		dataplatformNodePool.Properties.Annotations = &annotations
	}

	return &dataplatformNodePool
}

func GetDataplatformNodePoolDataUpdate(d *schema.ResourceData) *dataplatform.PatchNodePoolRequest {

	dataplatformNodePool := dataplatform.PatchNodePoolRequest{
		Properties: &dataplatform.PatchNodePoolProperties{},
	}

	if nodeCountValue, ok := d.GetOk("node_count"); ok {
		nodeCount := int32(nodeCountValue.(int))
		dataplatformNodePool.Properties.NodeCount = &nodeCount
	}

	if _, ok := d.GetOk("maintenance_window"); ok {
		dataplatformNodePool.Properties.MaintenanceWindow = setMaintenanceWindowData(d)
	}

	if labelsValue, ok := d.GetOk("labels"); ok {
		labels := make(map[string]interface{})
		for k, v := range labelsValue.(map[string]interface{}) {
			labels[k] = v.(string)
		}
		dataplatformNodePool.Properties.Labels = &labels
	}

	if annotationsValue, ok := d.GetOk("annotations"); ok {
		annotations := make(map[string]interface{})
		for k, v := range annotationsValue.(map[string]interface{}) {
			annotations[k] = v.(string)
		}
		dataplatformNodePool.Properties.Annotations = &annotations
	}

	return &dataplatformNodePool
}

func SetDataplatformNodePoolData(d *schema.ResourceData, nodePool dataplatform.NodePoolResponseData) error {

	if nodePool.Id != nil {
		d.SetId(*nodePool.Id)
	}

	if nodePool.Properties == nil {
		return fmt.Errorf("node pool properties should not be empty for ID %s", *nodePool.Id)
	}

	if nodePool.Properties.Name != nil {
		if err := d.Set("name", *nodePool.Properties.Name); err != nil {
			return utils.GenerateSetError(nodePoolResourceName, "name", err)
		}
	}

	if nodePool.Properties.DataPlatformVersion != nil {
		if err := d.Set("version", *nodePool.Properties.DataPlatformVersion); err != nil {
			return utils.GenerateSetError(nodePoolResourceName, "version", err)
		}
	}

	if nodePool.Properties.DatacenterId != nil {
		if err := d.Set("datacenter_id", *nodePool.Properties.DatacenterId); err != nil {
			return utils.GenerateSetError(nodePoolResourceName, "datacenter_id", err)
		}
	}

	if nodePool.Properties.NodeCount != nil {
		if err := d.Set("node_count", *nodePool.Properties.NodeCount); err != nil {
			return utils.GenerateSetError(nodePoolResourceName, "node_count", err)
		}
	}

	if nodePool.Properties.CpuFamily != nil {
		if err := d.Set("cpu_family", *nodePool.Properties.CpuFamily); err != nil {
			return utils.GenerateSetError(nodePoolResourceName, "cpu_family", err)
		}
	}

	if nodePool.Properties.CoresCount != nil {
		if err := d.Set("cores_count", *nodePool.Properties.CoresCount); err != nil {
			return utils.GenerateSetError(nodePoolResourceName, "cores_count", err)
		}
	}

	if nodePool.Properties.RamSize != nil {
		if err := d.Set("ram_size", *nodePool.Properties.RamSize); err != nil {
			return utils.GenerateSetError(nodePoolResourceName, "ram_size", err)
		}
	}

	if nodePool.Properties.AvailabilityZone != nil {
		if err := d.Set("availability_zone", *nodePool.Properties.AvailabilityZone); err != nil {
			return utils.GenerateSetError(nodePoolResourceName, "availability_zone", err)
		}
	}

	if nodePool.Properties.StorageType != nil {
		if err := d.Set("storage_type", *nodePool.Properties.StorageType); err != nil {
			return utils.GenerateSetError(nodePoolResourceName, "storage_type", err)
		}
	}

	if nodePool.Properties.StorageSize != nil {
		if err := d.Set("storage_size", *nodePool.Properties.StorageSize); err != nil {
			return utils.GenerateSetError(nodePoolResourceName, "storage_size", err)
		}
	}

	if nodePool.Properties.MaintenanceWindow != nil {
		var maintenanceWindow []interface{}
		maintenanceWindowEntry := SetMaintenanceWindowProperties(*nodePool.Properties.MaintenanceWindow)
		maintenanceWindow = append(maintenanceWindow, maintenanceWindowEntry)
		if err := d.Set("maintenance_window", maintenanceWindow); err != nil {
			return utils.GenerateSetError(nodePoolResourceName, "maintenance_window", err)
		}
	}

	if nodePool.Properties.Labels != nil {
		if err := d.Set("labels", *nodePool.Properties.Labels); err != nil {
			return utils.GenerateSetError(nodePoolResourceName, "labels", err)
		}
	}

	if nodePool.Properties.Annotations != nil {
		if err := d.Set("annotations", *nodePool.Properties.Annotations); err != nil {
			return utils.GenerateSetError(nodePoolResourceName, "annotations", err)
		}
	}
	return nil
}

func SetNodePoolsData(d *schema.ResourceData, results []dataplatform.NodePoolResponseData) diag.Diagnostics {

	resourceId := uuid.New()
	d.SetId(resourceId.String())

	if results != nil {
		var nodePools []interface{}
		for _, nodePool := range results {
			nodePoolEntry := make(map[string]interface{})
			if nodePool.Properties != nil {
				utils.SetPropWithNilCheck(nodePoolEntry, "name", nodePool.Properties.Name)
				utils.SetPropWithNilCheck(nodePoolEntry, "version", nodePool.Properties.DataPlatformVersion)
				utils.SetPropWithNilCheck(nodePoolEntry, "datacenter_id", nodePool.Properties.DatacenterId)
				utils.SetPropWithNilCheck(nodePoolEntry, "node_count", nodePool.Properties.NodeCount)
				utils.SetPropWithNilCheck(nodePoolEntry, "cpu_family", nodePool.Properties.CpuFamily)
				utils.SetPropWithNilCheck(nodePoolEntry, "cores_count", nodePool.Properties.CoresCount)
				utils.SetPropWithNilCheck(nodePoolEntry, "ram_size", nodePool.Properties.RamSize)
				utils.SetPropWithNilCheck(nodePoolEntry, "availability_zone", nodePool.Properties.AvailabilityZone)
				utils.SetPropWithNilCheck(nodePoolEntry, "storage_type", nodePool.Properties.StorageType)
				utils.SetPropWithNilCheck(nodePoolEntry, "storage_size", nodePool.Properties.StorageSize)
				if nodePool.Properties.MaintenanceWindow != nil {
					var maintenanceWindow []interface{}
					maintenanceWindowEntry := SetMaintenanceWindowProperties(*nodePool.Properties.MaintenanceWindow)
					maintenanceWindow = append(maintenanceWindow, maintenanceWindowEntry)
					utils.SetPropWithNilCheck(nodePoolEntry, "maintenance_window", maintenanceWindow)
				}
				utils.SetPropWithNilCheck(nodePoolEntry, "labels", nodePool.Properties.Labels)
				utils.SetPropWithNilCheck(nodePoolEntry, "annotations", nodePool.Properties.Annotations)
			}
			nodePools = append(nodePools, nodePoolEntry)
		}

		if nodePools == nil || len(nodePools) == 0 {
			return diag.FromErr(fmt.Errorf("no nodepools found for criteria, please check your filter configuration"))
		}

		err := d.Set("node_pools", nodePools)
		if err != nil {
			diags := diag.FromErr(fmt.Errorf("an error occurred while setting node_pools: %w", err))
			return diags
		}
	}
	return nil
}
