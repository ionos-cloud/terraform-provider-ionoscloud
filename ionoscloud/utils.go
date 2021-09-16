package ionoscloud

import (
	"context"
	"fmt"
	ionoscloud "github.com/ionos-cloud/sdk-go/v6"
	"log"
	"reflect"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const SleepInterval = 5 * time.Second

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

func resourceServerImport(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	parts := strings.Split(d.Id(), "/")

	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return nil, fmt.Errorf("invalid import id %q. Expecting {datacenter}/{server}", d.Id())
	}

	datacenterId := parts[0]
	serverId := parts[1]

	client := meta.(*ionoscloud.APIClient)

	server, apiResponse, err := client.ServersApi.DatacentersServersFindById(ctx, datacenterId, serverId).Execute()

	if err != nil {
		if apiResponse != nil && apiResponse.Response != nil && apiResponse.StatusCode == 404 {
			d.SetId("")
			return nil, fmt.Errorf("unable to find server %q", serverId)
		}
		return nil, fmt.Errorf("error occured while fetching a server ID %s %s", d.Id(), err)
	}

	d.SetId(*server.Id)

	if err := d.Set("datacenter_id", datacenterId); err != nil {
		return nil, err
	}

	if server.Properties.Name != nil {
		if err := d.Set("name", *server.Properties.Name); err != nil {
			return nil, err
		}
	}

	if server.Properties.Cores != nil {
		if err := d.Set("cores", *server.Properties.Cores); err != nil {
			return nil, err
		}
	}

	if server.Properties.Ram != nil {
		if err := d.Set("ram", *server.Properties.Ram); err != nil {
			return nil, err
		}
	}

	if server.Properties.AvailabilityZone != nil {
		if err := d.Set("availability_zone", *server.Properties.AvailabilityZone); err != nil {
			return nil, err
		}
	}

	if server.Properties.CpuFamily != nil {
		if err := d.Set("cpu_family", *server.Properties.CpuFamily); err != nil {
			return nil, err
		}
	}

	if server.Entities.Volumes != nil &&
		len(*server.Entities.Volumes.Items) > 0 &&
		(*server.Entities.Volumes.Items)[0].Properties.Image != nil {
		if err := d.Set("boot_image", *(*server.Entities.Volumes.Items)[0].Properties.Image); err != nil {
			return nil, err
		}
	}

	if server.Entities.Nics != nil && len(*server.Entities.Nics.Items) > 0 && (*server.Entities.Nics.Items)[0].Id != nil {
		primaryNic := *(*server.Entities.Nics.Items)[0].Id
		if err := d.Set("primary_nic", primaryNic); err != nil {
			return nil, err
		}

		nic, _, err := client.NetworkInterfacesApi.DatacentersServersNicsFindById(ctx, datacenterId, serverId, primaryNic).Execute()
		if err != nil {
			return nil, err
		}

		if len(*nic.Properties.Ips) > 0 {
			if err := d.Set("primary_ip", (*nic.Properties.Ips)[0]); err != nil {
				return nil, err
			}
		}

		network := map[string]interface{}{}

		setPropWithNilCheck(network, "dhcp", nic.Properties.Dhcp)
		setPropWithNilCheck(network, "firewall_active", nic.Properties.FirewallActive)

		setPropWithNilCheck(network, "lan", nic.Properties.Lan)
		setPropWithNilCheck(network, "name", nic.Properties.Name)
		setPropWithNilCheck(network, "ips", nic.Properties.Ips)
		setPropWithNilCheck(network, "mac", nic.Properties.Mac)

		if nic.Properties.Ips != nil && len(*nic.Properties.Ips) > 0 {
			network["ips"] = *nic.Properties.Ips
		}

		firewallRules, _, err := client.FirewallRulesApi.DatacentersServersNicsFirewallrulesGet(ctx, datacenterId, serverId, primaryNic).Execute()

		if err != nil {
			return nil, err
		}

		if firewallRules.Items != nil {
			if len(*firewallRules.Items) > 0 {
				if err := d.Set("firewallrule_id", *(*firewallRules.Items)[0].Id); err != nil {
					return nil, err
				}
			}
		}

		if firewallId, ok := d.GetOk("firewallrule_id"); ok {
			firewall, _, err := client.FirewallRulesApi.DatacentersServersNicsFirewallrulesFindById(ctx, datacenterId, serverId, primaryNic, firewallId.(string)).Execute()
			if err != nil {
				return nil, err
			}

			fw := map[string]interface{}{}
			/*
				"protocol": *firewall.Properties.Protocol,
				"name":     *firewall.Properties.Name,
			*/
			setPropWithNilCheck(fw, "protocol", firewall.Properties.Protocol)
			setPropWithNilCheck(fw, "name", firewall.Properties.Name)
			setPropWithNilCheck(fw, "source_mac", firewall.Properties.SourceMac)
			setPropWithNilCheck(fw, "source_ip", firewall.Properties.SourceIp)
			setPropWithNilCheck(fw, "target_ip", firewall.Properties.TargetIp)
			setPropWithNilCheck(fw, "port_range_start", firewall.Properties.PortRangeStart)
			setPropWithNilCheck(fw, "port_range_end", firewall.Properties.PortRangeEnd)
			setPropWithNilCheck(fw, "icmp_type", firewall.Properties.IcmpType)
			setPropWithNilCheck(fw, "icmp_code", firewall.Properties.IcmpCode)

			network["firewall"] = []map[string]interface{}{fw}
		}

		networks := []map[string]interface{}{network}
		if err := d.Set("nic", networks); err != nil {
			return nil, err
		}
	}

	if server.Properties.BootVolume != nil {
		if server.Properties.BootVolume.Id != nil {
			if err := d.Set("boot_volume", *server.Properties.BootVolume.Id); err != nil {
				return nil, err
			}
		}
		volumeObj, _, err := client.ServersApi.DatacentersServersVolumesFindById(ctx, datacenterId, serverId, *server.Properties.BootVolume.Id).Execute()
		if err == nil {
			volumeItem := map[string]interface{}{}

			setPropWithNilCheck(volumeItem, "name", volumeObj.Properties.Name)
			setPropWithNilCheck(volumeItem, "disk_type", volumeObj.Properties.Type)
			setPropWithNilCheck(volumeItem, "size", volumeObj.Properties.Size)
			setPropWithNilCheck(volumeItem, "licence_type", volumeObj.Properties.LicenceType)
			setPropWithNilCheck(volumeItem, "bus", volumeObj.Properties.Bus)
			setPropWithNilCheck(volumeItem, "availability_zone", volumeObj.Properties.AvailabilityZone)

			volumesList := []map[string]interface{}{volumeItem}
			if err := d.Set("volume", volumesList); err != nil {
				return nil, err
			}
		}
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

	clusterId := d.Id()

	cluster, apiResponse, err := client.KubernetesApi.K8sFindByClusterId(ctx, clusterId).Execute()

	if err != nil {
		if apiResponse != nil && apiResponse.StatusCode == 404 {
			d.SetId("")
			return nil, fmt.Errorf("unable to find k8s cluster %q", clusterId)
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

	clusterId := parts[0]
	npId := parts[1]

	client := meta.(*ionoscloud.APIClient)

	k8sNodepool, apiResponse, err := client.KubernetesApi.K8sNodepoolsFindById(ctx, clusterId, npId).Execute()

	if err != nil {
		if apiResponse != nil && apiResponse.StatusCode == 404 {
			d.SetId("")
			return nil, fmt.Errorf("unable to find k8s node pool %q", npId)
		}
		return nil, fmt.Errorf("unable to retreive k8s node pool %q", npId)
	}

	log.Printf("[INFO] K8s node pool found: %+v", k8sNodepool)

	d.SetId(*k8sNodepool.Id)
	if err := d.Set("name", *k8sNodepool.Properties.Name); err != nil {
		return nil, err
	}
	if err := d.Set("k8s_version", *k8sNodepool.Properties.K8sVersion); err != nil {
		return nil, err
	}
	if err := d.Set("k8s_cluster_id", clusterId); err != nil {
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

	if k8sNodepool.Properties.AutoScaling != nil && k8sNodepool.Properties.AutoScaling.MinNodeCount != nil && k8sNodepool.Properties.AutoScaling.MaxNodeCount != nil && (*k8sNodepool.Properties.AutoScaling.MinNodeCount != 0 && *k8sNodepool.Properties.AutoScaling.MaxNodeCount != 0) {
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

	pccId := d.Id()

	pcc, apiResponse, err := client.PrivateCrossConnectsApi.PccsFindById(ctx, d.Id()).Execute()

	if err != nil {
		if apiResponse != nil && apiResponse.StatusCode == 404 {
			d.SetId("")
			return nil, fmt.Errorf("unable to find PCC %q", pccId)
		}
		return nil, fmt.Errorf("unable to retreive PCC %q", pccId)
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

	buId := d.Id()

	backupUnit, apiResponse, err := client.BackupUnitsApi.BackupunitsFindById(ctx, d.Id()).Execute()

	if err != nil {
		if apiResponse != nil && apiResponse.StatusCode == 404 {
			d.SetId("")
			return nil, fmt.Errorf("unable to find Backup Unit %q", buId)
		}
		return nil, fmt.Errorf("unable to retreive Backup Unit %q", buId)
	}

	log.Printf("[INFO] Backup Unit found: %+v", backupUnit)

	d.SetId(*backupUnit.Id)

	if err := d.Set("name", *backupUnit.Properties.Name); err != nil {
		return nil, err
	}
	if err := d.Set("email", *backupUnit.Properties.Email); err != nil {
		return nil, err
	}

	contractResources, apiResponse, cErr := client.ContractResourcesApi.ContractsGet(ctx).Execute()

	if cErr != nil {
		return nil, fmt.Errorf("error while fetching contract resources for backup unit %q: %s", d.Id(), cErr)
	}

	if contractResources.Items == nil || len(*contractResources.Items) == 0 {
		return nil, fmt.Errorf("no contracts found for user")
	}

	props := (*contractResources.Items)[0].Properties
	if props == nil {
		return nil, fmt.Errorf("could not get first contract properties")
	}

	if props.ContractNumber == nil {
		return nil, fmt.Errorf("contract number not set")
	}

	if err := d.Set("login", fmt.Sprintf("%s-%d", *backupUnit.Properties.Name, *props.ContractNumber)); err != nil {
		return nil, err
	}

	return []*schema.ResourceData{d}, nil
}

func resourceS3KeyImport(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	parts := strings.Split(d.Id(), "/")

	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return nil, fmt.Errorf("invalid import id %q. Expecting {userId}/{s3KeyId}", d.Id())
	}

	userId := parts[0]
	keyId := parts[1]

	client := meta.(*ionoscloud.APIClient)

	s3Key, apiResponse, err := client.UserS3KeysApi.UmUsersS3keysFindByKeyId(ctx, userId, keyId).Execute()

	if err != nil {
		if apiResponse != nil && apiResponse.StatusCode == 404 {
			d.SetId("")
			return nil, fmt.Errorf("unable to find S3 key %q", keyId)
		}
		return nil, fmt.Errorf("unable to retreive S3 key %q", keyId)
	}

	d.SetId(*s3Key.Id)
	if err := d.Set("user_id", userId); err != nil {
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

func resourceVolumeImporter(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	parts := strings.Split(d.Id(), "/")
	if len(parts) != 3 || parts[0] == "" || parts[1] == "" || parts[2] == "" {
		return nil, fmt.Errorf("invalid import id %q. Expecting {datacenter}/{server}/{volume}", d.Id())
	}

	client := meta.(*ionoscloud.APIClient)

	dcId := parts[0]
	srvId := parts[1]
	volumeId := parts[2]

	volume, apiResponse, err := client.VolumesApi.DatacentersVolumesFindById(ctx, dcId, volumeId).Execute()

	if err != nil {
		if apiResponse != nil && apiResponse.StatusCode == 404 {
			d.SetId("")
			return nil, fmt.Errorf("volume does not exist %q", volumeId)
		}
		return nil, fmt.Errorf("an error occured while trying to find the volume %q", volumeId)
	}

	log.Printf("[INFO] volume found: %+v", volume)

	d.SetId(*volume.Id)
	if err := d.Set("datacenter_id", dcId); err != nil {
		return nil, err
	}

	if err := d.Set("server_id", srvId); err != nil {
		return nil, err
	}

	if volume.Properties.Name != nil {
		err := d.Set("name", *volume.Properties.Name)
		if err != nil {
			return nil, err
		}
	}

	if volume.Properties.Type != nil {
		err := d.Set("disk_type", *volume.Properties.Type)
		if err != nil {
			return nil, err
		}
	}

	if volume.Properties.Size != nil {
		err := d.Set("size", *volume.Properties.Size)
		if err != nil {
			return nil, err
		}
	}

	if volume.Properties.Bus != nil {
		err := d.Set("bus", *volume.Properties.Bus)
		if err != nil {
			return nil, err
		}
	}

	if volume.Properties.AvailabilityZone != nil {
		err := d.Set("availability_zone", *volume.Properties.AvailabilityZone)
		if err != nil {
			return nil, err
		}
	}

	if volume.Properties.CpuHotPlug != nil {
		err := d.Set("cpu_hot_plug", *volume.Properties.CpuHotPlug)
		if err != nil {
			return nil, err
		}
	}

	if volume.Properties.RamHotPlug != nil {
		err := d.Set("ram_hot_plug", *volume.Properties.RamHotPlug)
		if err != nil {
			return nil, err
		}
	}

	if volume.Properties.NicHotPlug != nil {
		err := d.Set("nic_hot_plug", *volume.Properties.NicHotPlug)
		if err != nil {
			return nil, err
		}
	}

	if volume.Properties.NicHotUnplug != nil {
		err := d.Set("nic_hot_unplug", *volume.Properties.NicHotUnplug)
		if err != nil {
			return nil, err
		}
	}

	if volume.Properties.DiscVirtioHotPlug != nil {
		err := d.Set("disc_virtio_hot_plug", *volume.Properties.DiscVirtioHotPlug)
		if err != nil {
			return nil, err
		}
	}

	if volume.Properties.DiscVirtioHotUnplug != nil {
		err := d.Set("disc_virtio_hot_unplug", *volume.Properties.DiscVirtioHotUnplug)
		if err != nil {
			return nil, err
		}
	}

	if volume.Properties.BackupunitId != nil {
		err := d.Set("backup_unit_id", *volume.Properties.BackupunitId)
		if err != nil {
			return nil, err
		}
	}

	if volume.Properties.UserData != nil {
		err := d.Set("user_data", *volume.Properties.UserData)
		if err != nil {
			return nil, err
		}
	}
	return []*schema.ResourceData{d}, nil
}

func resourceGroupImporter(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	client := meta.(*ionoscloud.APIClient)

	grpId := d.Id()

	group, apiResponse, err := client.UserManagementApi.UmGroupsFindById(ctx, grpId).Execute()

	if err != nil {
		if apiResponse != nil && apiResponse.StatusCode == 404 {
			d.SetId("")
			return nil, fmt.Errorf("an error occured while trying to fetch the group %q", grpId)
		}
		return nil, fmt.Errorf("group does not exist%q", grpId)
	}

	log.Printf("[INFO] group found: %+v", group)

	d.SetId(*group.Id)

	if group.Properties.Name != nil {
		err := d.Set("name", *group.Properties.Name)
		if err != nil {
			return nil, err
		}
	}

	if group.Properties.CreateDataCenter != nil {
		err := d.Set("create_datacenter", *group.Properties.CreateDataCenter)
		if err != nil {
			return nil, err
		}
	}

	if group.Properties.CreateSnapshot != nil {
		err := d.Set("create_snapshot", *group.Properties.CreateSnapshot)
		if err != nil {
			return nil, err
		}
	}

	if group.Properties.ReserveIp != nil {
		err := d.Set("reserve_ip", *group.Properties.ReserveIp)
		if err != nil {
			return nil, err
		}
	}

	if group.Properties.AccessActivityLog != nil {
		err := d.Set("access_activity_log", *group.Properties.AccessActivityLog)
		if err != nil {
			return nil, err
		}
	}

	if group.Properties.CreatePcc != nil {
		err := d.Set("create_pcc", *group.Properties.CreatePcc)
		if err != nil {
			return nil, err
		}
	}

	if group.Properties.S3Privilege != nil {
		err := d.Set("s3_privilege", *group.Properties.S3Privilege)
		if err != nil {
			return nil, err
		}
	}

	if group.Properties.CreateBackupUnit != nil {
		err := d.Set("create_backup_unit", *group.Properties.CreateBackupUnit)
		if err != nil {
			return nil, err
		}
	}

	if group.Properties.CreateInternetAccess != nil {
		err := d.Set("create_internet_access", *group.Properties.CreateInternetAccess)
		if err != nil {
			return nil, err
		}
	}

	if group.Properties.CreateK8sCluster != nil {
		err := d.Set("create_k8s_cluster", *group.Properties.CreateK8sCluster)
		if err != nil {
			return nil, err
		}
	}

	users, _, err := client.UserManagementApi.UmGroupsUsersGet(ctx, d.Id()).Execute()
	if err != nil {
		return nil, err
	}

	if users.Items != nil && len(*users.Items) > 0 {
		var usersEntries []interface{}
		for _, user := range *users.Items {
			userEntry := make(map[string]interface{})

			if user.Properties.Firstname != nil {
				userEntry["first_name"] = *user.Properties.Firstname
			}

			if user.Properties.Lastname != nil {
				userEntry["last_name"] = *user.Properties.Lastname
			}

			if user.Properties.Email != nil {
				userEntry["email"] = *user.Properties.Email
			}

			if user.Properties.Administrator != nil {
				userEntry["administrator"] = *user.Properties.Administrator
			}

			if user.Properties.ForceSecAuth != nil {
				userEntry["force_sec_auth"] = *user.Properties.ForceSecAuth
			}

			usersEntries = append(usersEntries, userEntry)
		}

		if len(usersEntries) > 0 {
			if err := d.Set("users", usersEntries); err != nil {
				return nil, err
			}
		}
	}

	return []*schema.ResourceData{d}, nil
}

func resourceUserImporter(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	client := meta.(*ionoscloud.APIClient)

	userId := d.Id()

	user, apiResponse, err := client.UserManagementApi.UmUsersFindById(ctx, userId).Execute()

	if err != nil {
		if apiResponse != nil && apiResponse.StatusCode == 404 {
			d.SetId("")
			return nil, fmt.Errorf("an error occured while trying to fetch the user %q", userId)
		}
		return nil, fmt.Errorf("user does not exist%q", userId)
	}

	log.Printf("[INFO] user found: %+v", user)

	d.SetId(*user.Id)

	if user.Properties.Firstname != nil {
		if err := d.Set("first_name", *user.Properties.Firstname); err != nil {
			return nil, err
		}
	}

	if user.Properties.Lastname != nil {
		if err := d.Set("last_name", *user.Properties.Lastname); err != nil {
			return nil, err
		}
	}
	if user.Properties.Email != nil {
		if err := d.Set("email", *user.Properties.Email); err != nil {
			return nil, err
		}
	}
	if user.Properties.Administrator != nil {
		if err := d.Set("administrator", *user.Properties.Administrator); err != nil {
			return nil, err
		}
	}
	if user.Properties.ForceSecAuth != nil {
		if err := d.Set("force_sec_auth", *user.Properties.ForceSecAuth); err != nil {
			return nil, err
		}
	}

	return []*schema.ResourceData{d}, nil
}

func resourceShareImporter(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	parts := strings.Split(d.Id(), "/")
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return nil, fmt.Errorf("invalid import id %q. Expecting {group}/{resource}", d.Id())
	}

	grpId := parts[0]
	rscId := parts[1]

	client := meta.(*ionoscloud.APIClient)

	share, apiResponse, err := client.UserManagementApi.UmGroupsSharesFindByResourceId(ctx, grpId, rscId).Execute()

	if err != nil {
		if apiResponse != nil && apiResponse.StatusCode == 404 {
			d.SetId("")
			return nil, fmt.Errorf("an error occured while trying to fetch the share of resource %q for group %q", rscId, grpId)
		}
		return nil, fmt.Errorf("share does not exist of resource %q for group %q", rscId, grpId)
	}

	log.Printf("[INFO] share found: %+v", share)

	d.SetId(*share.Id)

	if err := d.Set("group_id", grpId); err != nil {
		return nil, err
	}

	if err := d.Set("resource_id", rscId); err != nil {
		return nil, err
	}

	if share.Properties.EditPrivilege != nil {
		if err := d.Set("edit_privilege", *share.Properties.EditPrivilege); err != nil {
			return nil, err
		}
	}

	if share.Properties.SharePrivilege != nil {
		if err := d.Set("share_privilege", *share.Properties.SharePrivilege); err != nil {
			return nil, err
		}
	}

	return []*schema.ResourceData{d}, nil
}

func resourceIpFailoverImporter(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	parts := strings.Split(d.Id(), "/")
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return nil, fmt.Errorf("invalid import id %q. Expecting {datacenter}/{lan}", d.Id())
	}

	dcId := parts[0]
	lanId := parts[1]

	client := meta.(*ionoscloud.APIClient)

	lan, apiResponse, err := client.LansApi.DatacentersLansFindById(ctx, dcId, lanId).Execute()

	if err != nil {
		if apiResponse != nil && apiResponse.StatusCode == 404 {
			d.SetId("")
			return nil, fmt.Errorf("an error occured while trying to fetch the lan %q", lanId)
		}
		return nil, fmt.Errorf("lan does not exist%q", lanId)
	}

	log.Printf("[INFO] lan found: %+v", lan)

	d.SetId(*lan.Id)

	if err := d.Set("datacenter_id", dcId); err != nil {
		return nil, err
	}

	if lan.Properties.IpFailover != nil {
		err := d.Set("ip", *(*lan.Properties.IpFailover)[0].Ip)
		if err != nil {
			return nil, err
		}
	}

	if lan.Properties.IpFailover != nil {
		err := d.Set("nicuuid", *(*lan.Properties.IpFailover)[0].NicUuid)
		if err != nil {
			return nil, err
		}
	}

	if lan.Id != nil {
		err := d.Set("lan_id", *lan.Id)
		if err != nil {
			return nil, err
		}
	}

	return []*schema.ResourceData{d}, nil
}

func resourceLoadbalancerImporter(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	parts := strings.Split(d.Id(), "/")
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return nil, fmt.Errorf("invalid import id %q. Expecting {datacenter}/{loadbalancer}", d.Id())
	}

	dcId := parts[0]
	lbId := parts[1]

	client := meta.(*ionoscloud.APIClient)

	loadbalancer, apiResponse, err := client.LoadBalancersApi.DatacentersLoadbalancersFindById(ctx, dcId, lbId).Execute()

	if err != nil {
		if apiResponse != nil && apiResponse.StatusCode == 404 {
			d.SetId("")
			return nil, fmt.Errorf("an error occured while trying to fetch the loadbalancer %q", lbId)
		}
		return nil, fmt.Errorf("loadbalancer does not exist%q", lbId)
	}

	log.Printf("[INFO] loadbalancer found: %+v", loadbalancer)

	d.SetId(*loadbalancer.Id)

	if err := d.Set("datacenter_id", dcId); err != nil {
		return nil, err
	}

	if loadbalancer.Properties.Name != nil {
		if err := d.Set("name", *loadbalancer.Properties.Name); err != nil {
			return nil, err
		}
	}

	if loadbalancer.Properties.Ip != nil {
		if err := d.Set("ip", *loadbalancer.Properties.Ip); err != nil {
			return nil, err
		}
	}

	if loadbalancer.Properties.Dhcp != nil {
		if err := d.Set("dhcp", *loadbalancer.Properties.Dhcp); err != nil {
			return nil, err
		}
	}

	if loadbalancer.Entities != nil && loadbalancer.Entities.Balancednics != nil &&
		loadbalancer.Entities.Balancednics.Items != nil && len(*loadbalancer.Entities.Balancednics.Items) > 0 {

		var lans []string
		for _, lan := range *loadbalancer.Entities.Balancednics.Items {
			if *lan.Id != "" {
				lans = append(lans, *lan.Id)
			}
		}
		if err := d.Set("nic_ids", lans); err != nil {
			return nil, err
		}
	}

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

func responseBody(resp *ionoscloud.APIResponse) string {
	ret := "<nil>"
	if resp != nil {
		ret = string(resp.Payload)
	}

	return ret
}

func setPropWithNilCheck(m map[string]interface{}, prop string, v interface{}) {

	rVal := reflect.ValueOf(v)
	if rVal.Kind() == reflect.Ptr {
		if !rVal.IsNil() {
			m[prop] = rVal.Elem().Interface()
		}
	} else {
		m[prop] = v
	}
}
