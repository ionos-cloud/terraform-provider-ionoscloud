package ionoscloud

import (
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/profitbricks/profitbricks-sdk-go/v5"
)

func resourceResourceImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	parts := strings.Split(d.Id(), "/")
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return nil, fmt.Errorf("Invalid import id %q. Expecting {datacenter}/{resource}", d.Id())
	}

	d.Set("datacenter_id", parts[0])
	d.SetId(parts[1])

	return []*schema.ResourceData{d}, nil
}

func resourceServerImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	parts := strings.Split(d.Id(), "/")
	if len(parts) > 4 || len(parts) < 3 || parts[0] == "" || parts[1] == "" {
		return nil, fmt.Errorf("Invalid import id %q. Expecting {datacenter}/{server}/{primary_nic} or {datacenter}/{server}/{primary_nic}/{firewall}", d.Id())
	}

	d.Set("datacenter_id", parts[0])
	d.Set("primary_nic", parts[2])
	if len(parts) > 3 {
		d.Set("firewallrule_id", parts[3])
	}
	d.SetId(parts[1])

	return []*schema.ResourceData{d}, nil
}

func resourceK8sClusterImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	client := meta.(*profitbricks.Client)
	cluster, err := client.GetKubernetesCluster(d.Id())

	if err != nil {
		if apiError, ok := err.(profitbricks.ApiError); ok {
			if apiError.HttpStatusCode() == 404 {
				d.SetId("")
				return nil, fmt.Errorf("Unable to find k8s cluster %q", d.Id())
			}
		}
		return nil, fmt.Errorf("Unable to retreive k8s cluster %q", d.Id())
	}

	log.Printf("[INFO] K8s cluster found: %+v", cluster)
	d.SetId(cluster.ID)
	d.Set("name", cluster.Properties.Name)
	d.Set("k8s_version", cluster.Properties.K8sVersion)

	if cluster.Properties.MaintenanceWindow != nil {
		d.Set("maintenance_window", []map[string]string{
			{
				"day_of_the_week": cluster.Properties.MaintenanceWindow.DayOfTheWeek,
				"time":            cluster.Properties.MaintenanceWindow.Time,
			},
		})
		log.Printf("[INFO] Setting maintenance window for k8s cluster %s to %+v...", d.Id(), cluster.Properties.MaintenanceWindow)
	}

	log.Printf("[INFO] Importing k8s cluster %q...", d.Id())

	return []*schema.ResourceData{d}, nil
}

func resourceK8sNodepoolImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {

	parts := strings.Split(d.Id(), "/")

	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return nil, fmt.Errorf("Invalid import id %q. Expecting {k8sClusterId}/{k8sNodePoolId}", d.Id())
	}

	client := meta.(*profitbricks.Client)
	k8sNodepool, err := client.GetKubernetesNodePool(parts[0], parts[1])

	if err != nil {
		if apiError, ok := err.(profitbricks.ApiError); ok {
			if apiError.HttpStatusCode() == 404 {
				d.SetId("")
				return nil, fmt.Errorf("Unable to find k8s node pool %q", d.Id())
			}
		}
		return nil, fmt.Errorf("Unable to retreive k8s node pool %q", d.Id())
	}

	log.Printf("[INFO] K8s node pool found: %+v", k8sNodepool)

	d.SetId(k8sNodepool.ID)
	d.Set("name", k8sNodepool.Properties.Name)
	d.Set("k8s_version", k8sNodepool.Properties.K8sVersion)
	d.Set("k8s_cluster_id", parts[0])
	d.Set("datacenter_id", k8sNodepool.Properties.DatacenterID)
	d.Set("cpu_family", k8sNodepool.Properties.CPUFamily)
	d.Set("availability_zone", k8sNodepool.Properties.AvailabilityZone)
	d.Set("storage_type", k8sNodepool.Properties.StorageType)
	d.Set("node_count", k8sNodepool.Properties.NodeCount)
	d.Set("cores_count", k8sNodepool.Properties.CoresCount)
	d.Set("ram_size", k8sNodepool.Properties.RAMSize)
	d.Set("storage_size", k8sNodepool.Properties.StorageSize)

	if k8sNodepool.Properties.PublicIPs != nil {
		d.Set("public_ips", k8sNodepool.Properties.PublicIPs)
		log.Printf("[INFO] Setting Public IPs for k8s node pool %s to %+v...", d.Id(), d.Get("public_ips"))
	}

	if k8sNodepool.Properties.AutoScaling != nil && (k8sNodepool.Properties.AutoScaling.MinNodeCount != 0 && k8sNodepool.Properties.AutoScaling.MaxNodeCount != 0) {
		d.Set("auto_scaling", []map[string]uint32{
			{
				"min_node_count": k8sNodepool.Properties.AutoScaling.MinNodeCount,
				"max_node_count": k8sNodepool.Properties.AutoScaling.MaxNodeCount,
			},
		})
		log.Printf("[INFO] Setting AutoScaling for k8s node pool %s to %+v...", d.Id(), k8sNodepool.Properties.AutoScaling)
	}

	if k8sNodepool.Properties.LANs != nil {
		lans := []uint32{}

		for _, lan := range *k8sNodepool.Properties.LANs {
			lans = append(lans, lan.ID)
		}
		d.Set("lans", lans)
		log.Printf("[INFO] Setting LAN's for k8s node pool %s to %+v...", d.Id(), d.Get("lans"))
	}

	if k8sNodepool.Properties.MaintenanceWindow != nil {
		d.Set("maintenance_window", []map[string]string{
			{
				"day_of_the_week": k8sNodepool.Properties.MaintenanceWindow.DayOfTheWeek,
				"time":            k8sNodepool.Properties.MaintenanceWindow.Time,
			},
		})
		log.Printf("[INFO] Setting maintenance window for k8s node pool %s to %+v...", d.Id(), k8sNodepool.Properties.MaintenanceWindow)
	}

	log.Printf("[INFO] Importing k8s node pool %q...", d.Id())

	return []*schema.ResourceData{d}, nil
}

func resourcePrivateCrossConnectImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	client := meta.(*profitbricks.Client)
	pcc, err := client.GetPrivateCrossConnect(d.Id())

	if err != nil {
		if apiError, ok := err.(profitbricks.ApiError); ok {
			if apiError.HttpStatusCode() == 404 {
				d.SetId("")
				return nil, fmt.Errorf("Unable to find PCC %q", d.Id())
			}
		}
		return nil, fmt.Errorf("Unable to retreive PCC %q", d.Id())
	}

	log.Printf("[INFO] PCC found: %+v", pcc)

	d.SetId(pcc.ID)
	d.Set("name", pcc.Properties.Name)
	d.Set("description", pcc.Properties.Description)

	if pcc.Properties.Peers != nil {
		peers := []map[string]interface{}{}

		for _, peer := range *pcc.Properties.Peers {
			peers = append(peers, map[string]interface{}{
				"lan_id":          peer.LANId,
				"lan_name":        peer.LANName,
				"datacenter_id":   peer.DataCenterID,
				"datacenter_name": peer.DataCenterName,
				"location":        peer.Location,
			})
		}

		d.Set("peers", peers)
		log.Printf("[INFO] Setting peers for PCC %s to %+v...", d.Id(), d.Get("peers"))
	}

	if pcc.Properties.ConnectableDatacenters != nil {
		connectableDatacenters := []map[string]interface{}{}

		for _, connectableDatacenter := range *pcc.Properties.ConnectableDatacenters {
			connectableDatacenters = append(connectableDatacenters, map[string]interface{}{
				"id":       connectableDatacenter.ID,
				"name":     connectableDatacenter.Name,
				"location": connectableDatacenter.Location,
			})
		}

		d.Set("connectable_datacenters", connectableDatacenters)
		log.Printf("[INFO] Setting peers for PCC %s to %+v...", d.Id(), d.Get("peers"))
	}

	log.Printf("[INFO] Importing PCC %q...", d.Id())

	return []*schema.ResourceData{d}, nil
}

func resourceBackupUnitImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	client := meta.(*profitbricks.Client)
	backupUnit, err := client.GetBackupUnit(d.Id())

	if err != nil {
		if apiError, ok := err.(profitbricks.ApiError); ok {
			if apiError.HttpStatusCode() == 404 {
				d.SetId("")
				return nil, fmt.Errorf("Unable to find Backup Unit %q", d.Id())
			}
		}
		return nil, fmt.Errorf("Unable to retreive Backup Unit %q", d.Id())
	}

	log.Printf("[INFO] Backup Unit found: %+v", backupUnit)

	d.SetId(backupUnit.ID)

	d.Set("name", backupUnit.Properties.Name)
	d.Set("email", backupUnit.Properties.Email)

	contractResources, cErr := client.GetContractResources()

	if cErr != nil {
		return nil, fmt.Errorf("Error while fetching contract resources for backup unit %q: %s", d.Id(), cErr)
	}

	d.Set("login", fmt.Sprintf("%s-%d", backupUnit.Properties.Name, int64(contractResources.Properties.PBContractNumber)))

	return []*schema.ResourceData{d}, nil
}

func resourceS3KeyImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	parts := strings.Split(d.Id(), "/")

	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return nil, fmt.Errorf("Invalid import id %q. Expecting {userId}/{s3KeyId}", d.Id())
	}

	client := meta.(*profitbricks.Client)
	s3Key, err := client.GetS3Key(parts[0], parts[1])

	if err != nil {
		if apiError, ok := err.(profitbricks.ApiError); ok {
			if apiError.HttpStatusCode() == 404 {
				d.SetId("")
				return nil, fmt.Errorf("Unable to find S3 key %q", d.Id())
			}
		}
		return nil, fmt.Errorf("Unable to retreive S3 key %q", d.Id())
	}

	d.SetId(s3Key.ID)
	d.Set("user_id", parts[0])
	d.Set("secret_key", s3Key.Properties.SecretKey)
	d.Set("active", s3Key.Properties.Active)

	return []*schema.ResourceData{d}, nil
}

func resourceFirewallImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	parts := strings.Split(d.Id(), "/")
	if len(parts) != 4 || parts[0] == "" || parts[1] == "" {
		return nil, fmt.Errorf("Invalid import id %q. Expecting {datacenter}/{server}/{nic}/{firewall}", d.Id())
	}

	d.Set("datacenter_id", parts[0])
	d.Set("server_id", parts[1])
	d.Set("nic_id", parts[2])
	d.SetId(parts[3])

	return []*schema.ResourceData{d}, nil
}

func resourceNicImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	parts := strings.Split(d.Id(), "/")
	if len(parts) != 3 || parts[0] == "" || parts[1] == "" {
		return nil, fmt.Errorf("Invalid import id %q. Expecting {datacenter}/{server}/{nic}", d.Id())
	}

	d.Set("datacenter_id", parts[0])
	d.Set("server_id", parts[1])
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
