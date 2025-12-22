package ionoscloud

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	ionoscloud "github.com/ionos-cloud/sdk-go/v6"
)

func responseBody(resp *ionoscloud.APIResponse) string {
	ret := "<nil>"
	if resp != nil {
		ret = string(resp.Payload)
	}

	return ret
}

// DiffBasedOnVersion used for k8 node pool and cluster
// ignores downgrades of the patch versions.
func DiffBasedOnVersion(_, old, new string, _ *schema.ResourceData) bool {
	var oldMajor, oldMinor string
	var newMajor, newMinor string
	var oldPatchInt, newPatchInt int
	if old != "" {
		oldSplit := strings.Split(old, ".")
		if len(oldSplit) > 2 {
			oldMajor = oldSplit[0]
			oldMinor = oldSplit[1]
			oldPatchInt, _ = strconv.Atoi(oldSplit[2])
		}

		newSplit := strings.Split(new, ".")
		if len(newSplit) > 2 {
			newMajor = newSplit[0]
			newMinor = newSplit[1]
			newPatchInt, _ = strconv.Atoi(newSplit[2])
		}

		if oldMajor == newMajor && oldMinor == newMinor {
			// this is a downgrade of the patch version that we will ignore
			// it may happen either manually, or after a maintenance window
			if oldPatchInt > newPatchInt {
				log.Printf("[WARN] Downgrade is not supported on k8s from %d to %d", oldPatchInt, newPatchInt)
				return true
			}
		}
	}
	return false
}

// DiffCidr terraform suppress differences between ip and cidr
func DiffCidr(_, old, new string, _ *schema.ResourceData) bool {
	oldIp, _, err := net.ParseCIDR(old)
	newIp := net.ParseIP(new)
	// if new is an ip and old is a cidr, suppress
	if err == nil && newIp != nil && oldIp != nil && newIp.Equal(oldIp) {
		return true
	}
	return false
}

// DiffExpiryDate terraform suppress differences between layout and default +0000 UTC time format
func DiffExpiryDate(_, old, new string, _ *schema.ResourceData) bool {
	layout := "2006-01-02 15:04:05Z"
	oldTimeString := strings.Split(old, " +")
	oldTime, oldTimeErr := time.Parse(layout, oldTimeString[0]+"Z")
	newTime, newTimeErr := time.Parse(layout, new)
	if oldTimeErr == nil && newTimeErr == nil && newTime.Equal(oldTime) {
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
