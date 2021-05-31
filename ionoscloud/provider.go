package ionoscloud

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/httpclient"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/ionos-cloud/sdk-go/v6"
)

// Provider returns a schema.Provider for ionoscloud.
func Provider() terraform.ResourceProvider {
	provider := &schema.Provider{
		Schema: map[string]*schema.Schema{
			"username": {
				Type:          schema.TypeString,
				Optional:      true,
				DefaultFunc:   schema.EnvDefaultFunc("IONOS_USERNAME", nil),
				Description:   "IonosCloud username for API operations. If token is provided, token is preferred",
				ConflictsWith: []string{"token"},
			},
			"password": {
				Type:          schema.TypeString,
				Optional:      true,
				DefaultFunc:   schema.EnvDefaultFunc("IONOS_PASSWORD", nil),
				Description:   "IonosCloud password for API operations. If token is provided, token is preferred",
				ConflictsWith: []string{"token"},
			},
			"token": {
				Type:          schema.TypeString,
				Optional:      true,
				DefaultFunc:   schema.EnvDefaultFunc("IONOS_TOKEN", nil),
				Description:   "IonosCloud bearer token for API operations.",
				ConflictsWith: []string{"username", "password"},
			},
			"endpoint": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("IONOS_API_URL", ""),
				Description: "IonosCloud REST API URL.",
			},
			"retries": {
				Type:       schema.TypeInt,
				Optional:   true,
				Default:    50,
				Deprecated: "Timeout is used instead of this functionality",
			},
		},

		ResourcesMap: map[string]*schema.Resource{
			"ionoscloud_datacenter":                         resourceDatacenter(),
			"ionoscloud_ipblock":                            resourceIPBlock(),
			"ionoscloud_firewall":                           resourceFirewall(),
			"ionoscloud_lan":                                resourceLan(),
			"ionoscloud_loadbalancer":                       resourceLoadbalancer(),
			"ionoscloud_nic":                                resourceNic(),
			"ionoscloud_server":                             resourceServer(),
			"ionoscloud_volume":                             resourceVolume(),
			"ionoscloud_group":                              resourceGroup(),
			"ionoscloud_share":                              resourceShare(),
			"ionoscloud_user":                               resourceUser(),
			"ionoscloud_snapshot":                           resourceSnapshot(),
			"ionoscloud_ipfailover":                         resourceLanIPFailover(),
			"ionoscloud_k8s_cluster":                        resourcek8sCluster(),
			"ionoscloud_k8s_node_pool":                      resourcek8sNodePool(),
			"ionoscloud_private_crossconnect":               resourcePrivateCrossConnect(),
			"ionoscloud_backup_unit":                        resourceBackupUnit(),
			"ionoscloud_s3_key":                             resourceS3Key(),
			"ionoscloud_natgateway":                         resourceNatGateway(),
			"ionoscloud_natgateway_rule":                    resourceNatGatewayRule(),
			"ionoscloud_networkloadbalancer":                resourceNetworkLoadBalancer(),
			"ionoscloud_networkloadbalancer_forwardingrule": resourceNetworkLoadBalancerForwardingRule(),
		},
		DataSourcesMap: map[string]*schema.Resource{
			"ionoscloud_datacenter":                         dataSourceDataCenter(),
			"ionoscloud_location":                           dataSourceLocation(),
			"ionoscloud_image":                              dataSourceImage(),
			"ionoscloud_resource":                           dataSourceResource(),
			"ionoscloud_snapshot":                           dataSourceSnapshot(),
			"ionoscloud_lan":                                dataSourceLan(),
			"ionoscloud_private_crossconnect":               dataSourcePcc(),
			"ionoscloud_server":                             dataSourceServer(),
			"ionoscloud_k8s_cluster":                        dataSourceK8sCluster(),
			"ionoscloud_k8s_node_pool":                      dataSourceK8sNodePool(),
			"ionoscloud_natgateway":                         dataSourceNatGateway(),
			"ionoscloud_natgateway_rule":                    dataSourceNatGatewayRule(),
			"ionoscloud_networkloadbalancer":                dataSourceNetworkLoadBalancer(),
			"ionoscloud_networkloadbalancer_forwardingrule": dataSourceNetworkLoadBalancerForwardingRule(),
			"ionoscloud_template":                           dataSourceTemplate(),
		},
	}

	provider.ConfigureFunc = func(d *schema.ResourceData) (interface{}, error) {

		terraformVersion := provider.TerraformVersion

		if terraformVersion == "" {
			// Terraform 0.12 introduced this field to the protocol
			// We can therefore assume that if it's missing it's 0.10 or 0.11
			terraformVersion = "0.11+compatible"
		}

		log.Printf("[DEBUG] Setting terraformVersion to %s", terraformVersion)

		return providerConfigure(d, terraformVersion)
	}

	return provider
}

func providerConfigure(d *schema.ResourceData, terraformVersion string) (interface{}, error) {

	username, usernameOk := d.GetOk("username")
	password, passwordOk := d.GetOk("password")
	token, tokenOk := d.GetOk("token")

	if !tokenOk {
		if !usernameOk {
			return nil, fmt.Errorf("neither IonosCloud token, nor IonosCloud username has been provided")
		}

		if !passwordOk {
			return nil, fmt.Errorf("neither IonosCloud token, nor IonosCloud password has been provided")
		}
	} else {
		if usernameOk || passwordOk {
			return nil, fmt.Errorf("only provide IonosCloud token OR IonosCloud username/password")
		}
	}

	cleanedUrl := cleanURL(d.Get("endpoint").(string))

	newConfig := ionoscloud.NewConfiguration(username.(string), password.(string), token.(string))
	if len(cleanedUrl) > 0 {
		parts := strings.Split(cleanedUrl, "://")
		if len(parts) == 1 {
			newConfig.Host = cleanedUrl
		} else {
			newConfig.Scheme = parts[1]
			newConfig.Host = parts[2]
		}
	}
	newConfig.UserAgent = httpclient.TerraformUserAgent(terraformVersion)
	newClient := ionoscloud.NewAPIClient(newConfig)

	return newClient, nil
}

// cleanURL makes sure trailing slash does not corrupt the state
func cleanURL(url string) string {
	length := len(url)
	if length > 1 && url[length-1] == '/' {
		url = url[:length-1]
	}

	return url
}

// getStateChangeConf gets the default configuration for tracking a request progress
func getStateChangeConf(meta interface{}, d *schema.ResourceData, location string, timeoutType string) *resource.StateChangeConf {
	stateConf := &resource.StateChangeConf{
		Pending:        resourcePendingStates,
		Target:         resourceTargetStates,
		Refresh:        resourceStateRefreshFunc(meta, location),
		Timeout:        d.Timeout(timeoutType),
		MinTimeout:     10 * time.Second,
		Delay:          10 * time.Second, // Wait 10 secs before starting
		NotFoundChecks: 600,              //Setting high number, to support long timeouts
	}

	return stateConf
}

type RequestFailedError struct {
	msg string
}

func (e RequestFailedError) Error() string {
	return e.msg
}

func IsRequestFailed(err error) bool {
	_, ok := err.(RequestFailedError)
	return ok
}

// resourceStateRefreshFunc tracks progress of a request
func resourceStateRefreshFunc(meta interface{}, path string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		client := meta.(*ionoscloud.APIClient)

		fmt.Printf("[INFO] Checking PATH %s\n", path)
		if path == "" {
			return nil, "", fmt.Errorf("can not check a state when path is empty")
		}

		request, _, err := client.GetRequestStatus(context.Background(), path)

		if err != nil {
			return nil, "", fmt.Errorf("request failed with following error: %s", err)
		}

		if *request.Metadata.Status == "FAILED" {
			var msg string
			if request.Metadata.Message != nil {
				msg = fmt.Sprintf("Request failed with following error: %s", *request.Metadata.Message)
			} else {
				msg = "Request failed with an unknown error"
			}
			return nil, "", RequestFailedError{msg}
		}

		if *request.Metadata.Status == "DONE" {
			return request, "DONE", nil
		}

		return nil, *request.Metadata.Status, nil
	}
}

// resourcePendingStates defines states of working in progress
var resourcePendingStates = []string{
	"RUNNING",
	"QUEUED",
}

// resourceTargetStates defines states of completion
var resourceTargetStates = []string{
	"DONE",
}

// resourceDefaultTimeouts sets default value for each Timeout type
var resourceDefaultTimeouts = schema.ResourceTimeout{
	Create:  schema.DefaultTimeout(60 * time.Minute),
	Update:  schema.DefaultTimeout(60 * time.Minute),
	Delete:  schema.DefaultTimeout(60 * time.Minute),
	Default: schema.DefaultTimeout(60 * time.Minute),
}
