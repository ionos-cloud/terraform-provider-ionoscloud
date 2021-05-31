package ionoscloud

import (
	"context"
	"fmt"
	ionoscloud "github.com/ionos-cloud/sdk-go/v5"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceResourceImport(_ context.Context, d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData, error) {
	parts := strings.Split(d.Id(), "/")
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return nil, fmt.Errorf("invalid import id %q. Expecting {datacenter}/{resource}", d.Id())
	}

	if err := d.Set("datacenter_id", parts[0]); err != nil {
		return nil, err
	}
	d.SetId(parts[1])

	return []*schema.ResourceData{d}, nil
}

func resourceServerImport(_ context.Context, d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData, error) {
	parts := strings.Split(d.Id(), "/")
	if len(parts) > 4 || len(parts) < 3 || parts[0] == "" || parts[1] == "" {
		return nil, fmt.Errorf("invalid import id %q. Expecting {datacenter}/{server}/{primary_nic} or {datacenter}/{server}/{primary_nic}/{firewall}", d.Id())
	}

	if err := d.Set("datacenter_id", parts[0]); err != nil {
		return nil, err
	}

	if err := d.Set("primary_nic", parts[2]); err != nil {
		return nil, err
	}
	if len(parts) > 3 {
		if err := d.Set("firewallrule_id", parts[3]); err != nil {
			return nil, err
		}
	}
	d.SetId(parts[1])

	return []*schema.ResourceData{d}, nil
}

func resourceK8sClusterImport(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	client := meta.(*ionoscloud.APIClient)

	cluster, apiResponse, err := client.KubernetesApi.K8sFindByClusterId(ctx, d.Id()).Execute()

	if err != nil {
		if apiResponse != nil && apiResponse.StatusCode == 404 {
			d.SetId("")
			return nil, fmt.Errorf("unable to find k8s cluster %q", d.Id())
		}
		return nil, fmt.Errorf("unable to retreive k8s cluster %q", d.Id())
	}

	log.Printf("[INFO] K8s cluster found: %+v", cluster)
	d.SetId(*cluster.Id)
	if err := d.Set("name", *cluster.Properties.Name); err != nil {
		return nil, err
	}

	if err := d.Set("k8s_version", *cluster.Properties.K8sVersion); err != nil {
		return nil, err
	}

	if cluster.Properties.MaintenanceWindow != nil {
		if err := d.Set("maintenance_window", []map[string]string{
			{
				"day_of_the_week": *cluster.Properties.MaintenanceWindow.DayOfTheWeek,
				"time":            *cluster.Properties.MaintenanceWindow.Time,
			},
		}); err != nil {
			return nil, err
		}
		log.Printf("[INFO] Setting maintenance window for k8s cluster %s to %+v...", d.Id(), *cluster.Properties.MaintenanceWindow)
	}

	log.Printf("[INFO] Importing k8s cluster %q...", d.Id())

	return []*schema.ResourceData{d}, nil
}

func resourceK8sNodepoolImport(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {

	parts := strings.Split(d.Id(), "/")

	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return nil, fmt.Errorf("invalid import id %q. Expecting {k8sClusterId}/{k8sNodePoolId}", d.Id())
	}

	client := meta.(*ionoscloud.APIClient)

	k8sNodepool, apiResponse, err := client.KubernetesApi.K8sNodepoolsFindById(ctx, parts[0], parts[1]).Execute()

	if err != nil {
		if _, ok := err.(ionoscloud.GenericOpenAPIError); ok {
			if apiResponse != nil && apiResponse.StatusCode == 404 {
				d.SetId("")
				return nil, fmt.Errorf("unable to find k8s node pool %q", d.Id())
			}
		}
		return nil, fmt.Errorf("unable to retreive k8s node pool %q", d.Id())
	}

	log.Printf("[INFO] K8s node pool found: %+v", k8sNodepool)

	d.SetId(*k8sNodepool.Id)
	if err := d.Set("name", *k8sNodepool.Properties.Name); err != nil {
		return nil, err
	}
	if err := d.Set("k8s_version", *k8sNodepool.Properties.K8sVersion); err != nil {
		return nil, err
	}
	if err := d.Set("k8s_cluster_id", parts[0]); err != nil {
		return nil, err
	}
	if err := d.Set("datacenter_id", *k8sNodepool.Properties.DatacenterId); err != nil {
		return nil, err
	}
	if err := d.Set("cpu_family", *k8sNodepool.Properties.CpuFamily); err != nil {
		return nil, err
	}
	if err := d.Set("availability_zone", *k8sNodepool.Properties.AvailabilityZone); err != nil {
		return nil, err
	}
	if err := d.Set("storage_type", *k8sNodepool.Properties.StorageType); err != nil {
		return nil, err
	}
	if err := d.Set("node_count", *k8sNodepool.Properties.NodeCount); err != nil {
		return nil, err
	}
	if err := d.Set("cores_count", *k8sNodepool.Properties.CoresCount); err != nil {
		return nil, err
	}
	if err := d.Set("ram_size", *k8sNodepool.Properties.RamSize); err != nil {
		return nil, err
	}
	if err := d.Set("storage_size", *k8sNodepool.Properties.StorageSize); err != nil {
		return nil, err
	}

	if k8sNodepool.Properties.PublicIps != nil {
		if err := d.Set("public_ips", k8sNodepool.Properties.PublicIps); err != nil {
			return nil, err
		}
		log.Printf("[INFO] Setting Public IPs for k8s node pool %s to %+v...", d.Id(), d.Get("public_ips"))
	}

	if k8sNodepool.Properties.AutoScaling != nil && (*k8sNodepool.Properties.AutoScaling.MinNodeCount != 0 && *k8sNodepool.Properties.AutoScaling.MaxNodeCount != 0) {
		if err := d.Set("auto_scaling", []map[string]int32{
			{
				"min_node_count": *k8sNodepool.Properties.AutoScaling.MinNodeCount,
				"max_node_count": *k8sNodepool.Properties.AutoScaling.MaxNodeCount,
			},
		}); err != nil {
			return nil, err
		}
		log.Printf("[INFO] Setting AutoScaling for k8s node pool %s to %+v...", d.Id(), k8sNodepool.Properties.AutoScaling)
	}

	if k8sNodepool.Properties.Lans != nil {
		var lans []int32

		for _, lan := range *k8sNodepool.Properties.Lans {
			lans = append(lans, *lan.Id)
		}

		if err := d.Set("lans", lans); err != nil {
			return nil, err
		}

		log.Printf("[INFO] Setting LAN's for k8s node pool %s to %+v...", d.Id(), d.Get("lans"))
	}

	if k8sNodepool.Properties.MaintenanceWindow != nil {
		if err := d.Set("maintenance_window", []map[string]string{
			{
				"day_of_the_week": *k8sNodepool.Properties.MaintenanceWindow.DayOfTheWeek,
				"time":            *k8sNodepool.Properties.MaintenanceWindow.Time,
			},
		}); err != nil {
			return nil, err
		}
		log.Printf("[INFO] Setting maintenance window for k8s node pool %s to %+v...", d.Id(), k8sNodepool.Properties.MaintenanceWindow)
	}

	log.Printf("[INFO] Importing k8s node pool %q...", d.Id())

	return []*schema.ResourceData{d}, nil
}

func resourcePrivateCrossConnectImport(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	client := meta.(*ionoscloud.APIClient)

	pcc, apiResponse, err := client.PrivateCrossConnectApi.PccsFindById(ctx, d.Id()).Execute()

	if err != nil {
		if _, ok := err.(ionoscloud.GenericOpenAPIError); ok {
			if apiResponse != nil && apiResponse.StatusCode == 404 {
				d.SetId("")
				return nil, fmt.Errorf("unable to find PCC %q", d.Id())
			}
		}
		return nil, fmt.Errorf("unable to retreive PCC %q", d.Id())
	}

	log.Printf("[INFO] PCC found: %+v", pcc)

	d.SetId(*pcc.Id)
	if err := d.Set("name", *pcc.Properties.Name); err != nil {
		return nil, err
	}
	if err := d.Set("description", *pcc.Properties.Description); err != nil {
		return nil, err
	}

	if pcc.Properties.Peers != nil {
		var peers []map[string]interface{}

		for _, peer := range *pcc.Properties.Peers {
			peers = append(peers, map[string]interface{}{
				"lan_id":          *peer.Id,
				"lan_name":        *peer.Name,
				"datacenter_id":   *peer.DatacenterId,
				"datacenter_name": *peer.DatacenterName,
				"location":        *peer.Location,
			})
		}

		if err := d.Set("peers", peers); err != nil {
			return nil, err
		}
		log.Printf("[INFO] Setting peers for PCC %s to %+v...", d.Id(), d.Get("peers"))
	}

	if pcc.Properties.ConnectableDatacenters != nil {
		var connectableDatacenters []map[string]interface{}

		for _, connectableDatacenter := range *pcc.Properties.ConnectableDatacenters {
			connectableDatacenters = append(connectableDatacenters, map[string]interface{}{
				"id":       *connectableDatacenter.Id,
				"name":     *connectableDatacenter.Name,
				"location": *connectableDatacenter.Location,
			})
		}

		if err := d.Set("connectable_datacenters", connectableDatacenters); err != nil {
			return nil, err
		}
		log.Printf("[INFO] Setting peers for PCC %s to %+v...", d.Id(), d.Get("peers"))
	}

	log.Printf("[INFO] Importing PCC %q...", d.Id())

	return []*schema.ResourceData{d}, nil
}

func resourceBackupUnitImport(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	client := meta.(*ionoscloud.APIClient)

	backupUnit, apiResponse, err := client.BackupUnitApi.BackupunitsFindById(ctx, d.Id()).Execute()

	if err != nil {
		if _, ok := err.(ionoscloud.GenericOpenAPIError); ok {
			if apiResponse != nil && apiResponse.StatusCode == 404 {
				d.SetId("")
				return nil, fmt.Errorf("unable to find Backup Unit %q", d.Id())
			}
		}
		return nil, fmt.Errorf("unable to retreive Backup Unit %q", d.Id())
	}

	log.Printf("[INFO] Backup Unit found: %+v", backupUnit)

	d.SetId(*backupUnit.Id)

	if err := d.Set("name", *backupUnit.Properties.Name); err != nil {
		return nil, err
	}
	if err := d.Set("email", *backupUnit.Properties.Email); err != nil {
		return nil, err
	}

	contractResources, apiResponse, cErr := client.ContractApi.ContractsGet(ctx).Execute()

	if cErr != nil {
		return nil, fmt.Errorf("error while fetching contract resources for backup unit %q: %s", d.Id(), cErr)
	}

	if err := d.Set("login", fmt.Sprintf("%s-%d", *backupUnit.Properties.Name, *contractResources.Properties.ContractNumber)); err != nil {
		return nil, err
	}

	return []*schema.ResourceData{d}, nil
}

func resourceS3KeyImport(_ context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	parts := strings.Split(d.Id(), "/")

	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return nil, fmt.Errorf("invalid import id %q. Expecting {userId}/{s3KeyId}", d.Id())
	}

	client := meta.(*ionoscloud.APIClient)

	s3Key, apiResponse, err := client.UserManagementApi.UmUsersS3keysFindByKeyId(context.TODO(), parts[0], parts[1]).Execute()

	if err != nil {
		if apiResponse != nil && apiResponse.StatusCode == 404 {
			d.SetId("")
			return nil, fmt.Errorf("unable to find S3 key %q", d.Id())
		}
		return nil, fmt.Errorf("unable to retreive S3 key %q", d.Id())
	}

	d.SetId(*s3Key.Id)
	if err := d.Set("user_id", parts[0]); err != nil {
		return nil, err
	}
	if err := d.Set("secret_key", *s3Key.Properties.SecretKey); err != nil {
		return nil, err
	}
	if err := d.Set("active", *s3Key.Properties.Active); err != nil {
		return nil, err
	}

	return []*schema.ResourceData{d}, nil
}

func resourceFirewallImport(_ context.Context, d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData, error) {
	parts := strings.Split(d.Id(), "/")
	if len(parts) != 4 || parts[0] == "" || parts[1] == "" {
		return nil, fmt.Errorf("invalid import id %q. Expecting {datacenter}/{server}/{nic}/{firewall}", d.Id())
	}

	if err := d.Set("datacenter_id", parts[0]); err != nil {
		return nil, err
	}
	if err := d.Set("server_id", parts[1]); err != nil {
		return nil, err
	}
	if err := d.Set("nic_id", parts[2]); err != nil {
		return nil, err
	}
	d.SetId(parts[3])

	return []*schema.ResourceData{d}, nil
}

func resourceNicImport(_ context.Context, d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData, error) {
	parts := strings.Split(d.Id(), "/")
	if len(parts) != 3 || parts[0] == "" || parts[1] == "" {
		return nil, fmt.Errorf("invalid import id %q. Expecting {datacenter}/{server}/{nic}", d.Id())
	}

	if err := d.Set("datacenter_id", parts[0]); err != nil {
		return nil, err
	}
	if err := d.Set("server_id", parts[1]); err != nil {
		return nil, err
	}
	d.SetId(parts[2])

	return []*schema.ResourceData{d}, nil
}

func convertSlice(slice []interface{}) []string {
	s := make([]string, len(slice))
	for i, v := range slice {
		s[i] = v.(string)
	}
	return s
}

func diffSlice(slice1 []string, slice2 []string) []string {
	var diff []string

	// Loop two times, first to find slice1 strings not in slice2,
	// second loop to find slice2 strings not in slice1
	for i := 0; i < 2; i++ {
		for _, s1 := range slice1 {
			found := false
			for _, s2 := range slice2 {
				if s1 == s2 {
					found = true
					break
				}
			}
			// String not found. We add it to return slice
			if !found {
				diff = append(diff, s1)
			}
		}
		// Swap the slices, only if it was the first loop
		if i == 0 {
			slice1, slice2 = slice2, slice1
		}
	}

	return diff
}
