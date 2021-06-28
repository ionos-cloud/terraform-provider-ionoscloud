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
		if apiResponse != nil && apiResponse.Response.StatusCode == 404 {
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
		if apiResponse != nil && apiResponse.Response.StatusCode == 404 {
			d.SetId("")
			return nil, fmt.Errorf("unable to find k8s node pool %q", d.Id())
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

	if k8sNodepool.Properties.AutoScaling != nil && k8sNodepool.Properties.AutoScaling.MinNodeCount != nil && k8sNodepool.Properties.AutoScaling.MaxNodeCount != nil &&
		(*k8sNodepool.Properties.AutoScaling.MinNodeCount != 0 && *k8sNodepool.Properties.AutoScaling.MaxNodeCount != 0) {
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

	pcc, apiResponse, err := client.PrivateCrossConnectsApi.PccsFindById(ctx, d.Id()).Execute()

	if err != nil {
		if apiResponse != nil && apiResponse.Response.StatusCode == 404 {
			d.SetId("")
			return nil, fmt.Errorf("unable to find PCC %q", d.Id())
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

	backupUnit, apiResponse, err := client.BackupUnitsApi.BackupunitsFindById(ctx, d.Id()).Execute()

	if err != nil {
		if apiResponse != nil && apiResponse.Response.StatusCode == 404 {
			d.SetId("")
			return nil, fmt.Errorf("unable to find Backup Unit %q", d.Id())
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

func resourceS3KeyImport(_ context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	parts := strings.Split(d.Id(), "/")

	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return nil, fmt.Errorf("invalid import id %q. Expecting {userId}/{s3KeyId}", d.Id())
	}

	client := meta.(*ionoscloud.APIClient)

	s3Key, apiResponse, err := client.UserS3KeysApi.UmUsersS3keysFindByKeyId(context.TODO(), parts[0], parts[1]).Execute()

	if err != nil {
		if apiResponse != nil && apiResponse.Response.StatusCode == 404 {
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

func resourceVolumeImporter(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	parts := strings.Split(d.Id(), "/")
	if len(parts) != 3 || parts[0] == "" || parts[1] == "" || parts[2] == "" {
		return nil, fmt.Errorf("invalid import id %q. Expecting {datacenter}/{server}/{volume}", d.Id())
	}

	client := meta.(*ionoscloud.APIClient)

	volume, apiResponse, err := client.VolumesApi.DatacentersVolumesFindById(ctx, parts[0], parts[2]).Execute()

	if err != nil {
		if apiResponse != nil && apiResponse.StatusCode == 404 {
			d.SetId("")
			return nil, fmt.Errorf("volume does not exist %q", d.Id())
		}
		return nil, fmt.Errorf("an error occured while trying to find the volume %q", d.Id())
	}

	log.Printf("[INFO] volume found: %+v", volume)

	d.SetId(*volume.Id)
	if err := d.Set("datacenter_id", parts[0]); err != nil {
		return nil, err
	}

	if err := d.Set("server_id", parts[1]); err != nil {
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

	group, apiResponse, err := client.UserManagementApi.UmGroupsFindById(ctx, d.Id()).Execute()

	if err != nil {
		if apiResponse != nil && apiResponse.StatusCode == 404 {
			d.SetId("")
			return nil, fmt.Errorf("an error occured while trying to fetch the group %q", d.Id())
		}
		return nil, fmt.Errorf("group does not exist%q", d.Id())
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

	usersEntries := make([]interface{}, 0)
	if users.Items != nil && len(*users.Items) > 0 {
		usersEntries = make([]interface{}, len(*users.Items))
		for userIndex, user := range *users.Items {
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

			usersEntries[userIndex] = userEntry
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

	user, apiResponse, err := client.UserManagementApi.UmUsersFindById(ctx, d.Id()).Execute()

	if err != nil {
		if apiResponse != nil && apiResponse.StatusCode == 404 {
			d.SetId("")
			return nil, fmt.Errorf("an error occured while trying to fetch the user %q", d.Id())
		}
		return nil, fmt.Errorf("user does not exist%q", d.Id())
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

	client := meta.(*ionoscloud.APIClient)

	share, apiResponse, err := client.UserManagementApi.UmGroupsSharesFindByResourceId(ctx, parts[0], parts[1]).Execute()

	if err != nil {
		if apiResponse != nil && apiResponse.StatusCode == 404 {
			d.SetId("")
			return nil, fmt.Errorf("an error occured while trying to fetch the share %q", d.Id())
		}
		return nil, fmt.Errorf("share does not exist%q", d.Id())
	}

	log.Printf("[INFO] share found: %+v", share)

	d.SetId(*share.Id)

	if err := d.Set("group_id", parts[0]); err != nil {
		return nil, err
	}

	if err := d.Set("resource_id", parts[1]); err != nil {
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
