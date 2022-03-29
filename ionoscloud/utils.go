package ionoscloud

import (
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
		if resp.Response != nil {
			log.Printf("[DEBUG] response status code : %d\n", resp.StatusCode)
		}
	}
}

func httpNotFound(resp *ionoscloud.APIResponse) bool {
	if resp != nil && resp.Response != nil && resp.StatusCode == http.StatusNotFound {
		return true
	}
	return false
}
