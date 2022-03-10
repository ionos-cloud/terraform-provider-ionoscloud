package ionoscloud

import (
	"context"
	"fmt"
	ionoscloud "github.com/ionos-cloud/sdk-go/v6"
	"log"
	"net"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const SleepInterval = 5 * time.Second

func resourceIpFailoverImporter(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	parts := strings.Split(d.Id(), "/")
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return nil, fmt.Errorf("invalid import id %q. Expecting {datacenter}/{lan}", d.Id())
	}

	dcId := parts[0]
	lanId := parts[1]

	client := meta.(SdkBundle).CloudApiClient

	lan, apiResponse, err := client.LANsApi.DatacentersLansFindById(ctx, dcId, lanId).Execute()
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

	client := meta.(SdkBundle).CloudApiClient

	loadbalancer, apiResponse, err := client.LoadBalancersApi.DatacentersLoadbalancersFindById(ctx, dcId, lbId).Execute()
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

func responseBody(resp *ionoscloud.APIResponse) string {
	ret := "<nil>"
	if resp != nil {
		ret = string(resp.Payload)
	}

	return ret
}

func IsValidUUID(uuid string) bool {
	r := regexp.MustCompile("^[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{12}$")
	return r.MatchString(uuid)
}

// DiffBasedOnVersion used for k8 node pool and cluster
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

//DiffToLower terraform suppress differences between lower and upper
func DiffToLower(_, old, new string, _ *schema.ResourceData) bool {
	if strings.ToLower(old) == strings.ToLower(new) {
		return true
	}
	return false
}

//DiffCidr terraform suppress differences between ip and cidr
func DiffCidr(_, old, new string, _ *schema.ResourceData) bool {
	oldIp, _, err := net.ParseCIDR(old)
	newIp := net.ParseIP(new)
	// if new is an ip and old is a cidr, suppress
	if err == nil && newIp != nil && oldIp != nil && newIp.Equal(oldIp) {
		return true
	}
	return false
}

// VerifyUnavailableIPs used for DBaaS cluster to check the provided IPs
func VerifyUnavailableIPs(val interface{}, key string) (warns []string, errs []error) {
	v := val.(string)
	unavailableNetworks := []string{"10.233.64.0/18", "10.233.0.0/18", "10.233.114.0/24"}

	ip, _, _ := net.ParseCIDR(v)

	for _, unavailableNetwork := range unavailableNetworks {
		if _, network, _ := net.ParseCIDR(unavailableNetwork); network.Contains(ip) {
			errs = append(errs, fmt.Errorf("for %q the following IP ranges are unavailable: 10.233.64.0/18, 10.233.0.0/18, 10.233.114.0/24; got: %v", key, v))
		}
	}
	return
}

func logApiRequestTime(resp *ionoscloud.APIResponse) {
	if resp != nil {
		log.Printf("[DEBUG] Request time : %s for operation : %s",
			resp.RequestTime, resp.Operation)
	}
}

func httpNotFound(resp *ionoscloud.APIResponse) bool {
	if resp != nil && resp.Response != nil && resp.StatusCode == http.StatusNotFound {
		return true
	}
	return false
}
