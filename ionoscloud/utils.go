package ionoscloud

import (
	"context"
	"fmt"
	ionoscloud "github.com/ionos-cloud/sdk-go/v5"
	"log"
	"reflect"
	"regexp"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const SleepInterval = 5 * time.Second

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
func resourceShareImporter(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	parts := strings.Split(d.Id(), "/")
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return nil, fmt.Errorf("invalid import id %q. Expecting {group}/{resource}", d.Id())
	}

	grpId := parts[0]
	rscId := parts[1]

	client := meta.(*ionoscloud.APIClient)

	share, apiResponse, err := client.UserManagementApi.UmGroupsSharesFindByResourceId(ctx, grpId, rscId).Execute()
	logApiRequestTime(apiResponse)

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

	lan, apiResponse, err := client.LanApi.DatacentersLansFindById(ctx, dcId, lanId).Execute()
	logApiRequestTime(apiResponse)

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

	loadbalancer, apiResponse, err := client.LoadBalancerApi.DatacentersLoadbalancersFindById(ctx, dcId, lbId).Execute()
	logApiRequestTime(apiResponse)
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

func IsValidUUID(uuid string) bool {
	r := regexp.MustCompile("^[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{12}$")
	return r.MatchString(uuid)
}

func boolAddr(b bool) *bool {
	return &b
}

func GenerateEmail() string {
	email := fmt.Sprintf("terraform_test-%d@mailinator.com", time.Now().UnixNano())
	return email
}

func logApiRequestTime(resp *ionoscloud.APIResponse) {
	if resp != nil {
		log.Printf("[DEBUG] Request time : %s for operation : %s",
			resp.RequestTime, resp.Operation)
	}
}
