package ionoscloud

import (
	"context"
	"fmt"
	ionoscloud "github.com/ionos-cloud/sdk-go/v6"
	"log"
	"reflect"
	"regexp"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const SleepInterval = 5 * time.Second

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
		if apiResponse != nil && apiResponse.Response != nil && apiResponse.StatusCode == 404 {
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

		nodePoolLans := getK8sNodePoolLans(*k8sNodepool.Properties.Lans)

		if len(nodePoolLans) > 0 {
			if err := d.Set("lans", nodePoolLans); err != nil {
				return nil, err
			}
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

	labels := make(map[string]interface{})
	if k8sNodepool.Properties.Labels != nil && len(*k8sNodepool.Properties.Labels) > 0 {
		for k, v := range *k8sNodepool.Properties.Labels {
			labels[k] = v
		}
	}

	if err := d.Set("labels", labels); err != nil {
		return nil, fmt.Errorf("error while setting the labels property for k8sNodepool %s: %s", d.Id(), err)
	}

	annotations := make(map[string]interface{})
	if k8sNodepool.Properties.Annotations != nil && len(*k8sNodepool.Properties.Annotations) > 0 {
		for k, v := range *k8sNodepool.Properties.Annotations {
			annotations[k] = v
		}
	}

	if err := d.Set("annotations", annotations); err != nil {
		return nil, fmt.Errorf("error while setting the annotations property for k8sNodepool %s: %s", d.Id(), err)
	}

	log.Printf("[INFO] Importing k8s node pool %q...", d.Id())

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

func resourceGroupImporter(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	client := meta.(*ionoscloud.APIClient)

	grpId := d.Id()

	group, apiResponse, err := client.UserManagementApi.UmGroupsFindById(ctx, grpId).Execute()

	if err != nil {
		if apiResponse != nil && apiResponse.Response != nil && apiResponse.StatusCode == 404 {
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
		if apiResponse != nil && apiResponse.Response != nil && apiResponse.StatusCode == 404 {
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
		if apiResponse != nil && apiResponse.Response != nil && apiResponse.StatusCode == 404 {
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
		if apiResponse != nil && apiResponse.Response != nil && apiResponse.StatusCode == 404 {
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
		if apiResponse != nil && apiResponse.Response != nil && apiResponse.StatusCode == 404 {
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

func IsValidUUID(uuid string) bool {
	r := regexp.MustCompile("^[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{12}$")
	return r.MatchString(uuid)
}

//used for k8 node pool and cluster
func DiffBasedOnVersion(_, old, new string, _ *schema.ResourceData) bool {
	var oldMajor, oldMinor string
	if old != "" {
		oldSplit := strings.Split(old, ".")
		oldMajor = oldSplit[0]
		oldMinor = oldSplit[1]

		newSplit := strings.Split(new, ".")
		newMajor := newSplit[0]
		newMinor := newSplit[1]

		if oldMajor == newMajor && oldMinor == newMinor {
			return true
		}
	}
	return false
}
