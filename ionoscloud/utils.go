package ionoscloud

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/ionos-cloud/sdk-go-bundle/shared"
	"log"
	"net"
	"net/http"
	"strings"
	"time"
)

func convertSlice(slice []interface{}) []string {
	s := make([]string, len(slice))
	for i, v := range slice {
		s[i] = v.(string)
	}
	return s
}

func responseBody(resp *shared.APIResponse) string {
	ret := "<nil>"
	if resp != nil {
		ret = string(resp.Payload)
	}

	return ret
}

// DiffBasedOnVersion used for k8 node pool and cluster
func DiffBasedOnVersion(_, old, new string, _ *schema.ResourceData) bool {
	var oldMajor, oldMinor string
	var newMajor, newMinor string
	if old != "" {
		oldSplit := strings.Split(old, ".")
		if len(oldSplit) > 1 {
			oldMajor = oldSplit[0]
			oldMinor = oldSplit[1]
		}

		newSplit := strings.Split(new, ".")
		if len(newSplit) > 1 {
			newMajor = newSplit[0]
			newMinor = newSplit[1]
		}

		if oldMajor == newMajor && oldMinor == newMinor {
			return true
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

func logApiRequestTime(resp *shared.APIResponse) {
	if resp != nil {
		log.Printf("[DEBUG] Request time : %s for operation : %s",
			resp.RequestTime, resp.Operation)
		if resp.Response != nil {
			log.Printf("[DEBUG] response status code : %d\n", resp.StatusCode)
		}
	}
}

func httpNotFound(resp *shared.APIResponse) bool {
	if resp != nil && resp.Response != nil && resp.StatusCode == http.StatusNotFound {
		return true
	}
	return false
}

// used for the datasource, when the nic is a member of the server object
var nicServerDSResource = &schema.Resource{
	Schema: map[string]*schema.Schema{
		"id": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"name": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"mac": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"ips": {
			Type:     schema.TypeList,
			Computed: true,
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
		"dhcp": {
			Type:     schema.TypeBool,
			Computed: true,
		},
		"lan": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"firewall_active": {
			Type:     schema.TypeBool,
			Computed: true,
		},
		"firewall_type": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"device_number": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"pci_slot": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"firewall_rules": {
			Type:     schema.TypeList,
			Computed: true,
			Elem:     firewallServerDSResource,
		},
	},
}

// used for the datasource, when the firewall is a member of the server object
var firewallServerDSResource = &schema.Resource{
	Schema: map[string]*schema.Schema{
		"id": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"name": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"protocol": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"source_mac": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"source_ip": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"target_ip": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"icmp_code": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"icmp_type": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"port_range_start": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"port_range_end": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"type": {
			Type:     schema.TypeString,
			Computed: true,
		},
	},
}

// used for the datasource, when the cdrom is a member of the server object
var cdromsServerDSResource = &schema.Resource{
	Schema: map[string]*schema.Schema{
		"id": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"name": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"description": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"location": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"size": {
			Type:     schema.TypeFloat,
			Computed: true,
		},
		"cpu_hot_plug": {
			Type:     schema.TypeBool,
			Computed: true,
		},
		"cpu_hot_unplug": {
			Type:     schema.TypeBool,
			Computed: true,
		},
		"ram_hot_plug": {
			Type:     schema.TypeBool,
			Computed: true,
		},
		"ram_hot_unplug": {
			Type:     schema.TypeBool,
			Computed: true,
		},
		"nic_hot_plug": {
			Type:     schema.TypeBool,
			Computed: true,
		},
		"nic_hot_unplug": {
			Type:     schema.TypeBool,
			Computed: true,
		},
		"disc_virtio_hot_plug": {
			Type:     schema.TypeBool,
			Computed: true,
		},
		"disc_virtio_hot_unplug": {
			Type:     schema.TypeBool,
			Computed: true,
		},
		"disc_scsi_hot_plug": {
			Type:     schema.TypeBool,
			Computed: true,
		},
		"disc_scsi_hot_unplug": {
			Type:     schema.TypeBool,
			Computed: true,
		},
		"licence_type": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"image_type": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"public": {
			Type:     schema.TypeBool,
			Computed: true,
		},
		"image_aliases": {
			Type:     schema.TypeList,
			Computed: true,
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
		"cloud_init": {
			Type:     schema.TypeString,
			Computed: true,
		},
	},
}
