package ionoscloud

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	ionoscloud "github.com/ionos-cloud/sdk-go/v6"

	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/bundleclient"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/cloudapi/cloudapinic"
	cloudapiflowlog "github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/cloudapi/flowlog"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/cloudapi/nsg"
)

func resourceNic() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceNicCreate,
		ReadContext:   resourceNicRead,
		UpdateContext: resourceNicUpdate,
		DeleteContext: resourceNicDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceNicImport,
		},
		Schema: map[string]*schema.Schema{

			"lan": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"dhcp": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"dhcpv6": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Indicates whether this NIC receives an IPv6 address through DHCP.",
			},
			"ipv6_cidr_block": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "IPv6 CIDR block assigned to the NIC.",
			},
			"ips": {
				Type: schema.TypeSet,
				Elem: &schema.Schema{
					Type:             schema.TypeString,
					ValidateDiagFunc: validation.ToDiagFunc(validation.StringIsNotWhiteSpace),
				},
				Computed:    true,
				Optional:    true,
				Description: "Collection of IP addresses assigned to a nic. Explicitly assigned public IPs need to come from reserved IP blocks, Passing value null or empty array will assign an IP address automatically.",
			},
			"ipv6_ips": {
				Type: schema.TypeSet,
				Elem: &schema.Schema{
					Type:             schema.TypeString,
					ValidateDiagFunc: validation.ToDiagFunc(validation.StringIsNotWhiteSpace),
				},
				Optional:    true,
				Computed:    true,
				Description: "Collection for IPv6 addresses assigned to a nic. Explicitly assigned IPv6 addresses need to come from inside the IPv6 CIDR block assigned to the nic.",
			},
			"firewall_active": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"firewall_type": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"server_id": {
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringIsNotWhiteSpace),
			},
			"datacenter_id": {
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringIsNotWhiteSpace),
			},
			"mac": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"device_number": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"pci_slot": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"flowlog": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     cloudapiflowlog.FlowlogSchemaResource,
				MaxItems: 1,
				Description: `Only 1 flow log can be configured. Only the name field can change as part of an update. Flow logs holistically capture network information such as source and destination 
IP addresses, source and destination ports, number of packets, amount of bytes, 
the start and end time of the recording, and the type of protocol – 
and log the extent to which your instances are being accessed.`,
			},
			"security_groups_ids": {
				Type:        schema.TypeSet,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Optional:    true,
				Description: "The list of Security Group IDs",
			},
		},
		Timeouts:      &resourceDefaultTimeouts,
		CustomizeDiff: ForceNewForFlowlogChanges,
	}
}

// ForceNewForFlowlogChanges - sets ForceNew either on `flowlog` if it is being deleted, or
// on the field that changes. This is needed because the API does not support PATCH for all flowlog fields except name.
// The API also does not support DELETE on the flowlog, so the whole resource needs to be re-created.
func ForceNewForFlowlogChanges(_ context.Context, d *schema.ResourceDiff, _ interface{}) error {
	// we do not want to check in case of resource creation
	if d.Id() == "" || !d.HasChange("flowlog") {
		return nil
	}

	oldVal, newVal := d.GetChange("flowlog")
	oldFlowLogs := oldVal.([]any)
	newFlowLogs := newVal.([]any)

	// "this check is for the scenario in which we have 0 initial flowlogs and we add a new one during a nic update
	if len(oldFlowLogs) == 0 && len(newFlowLogs) > 0 {
		return nil
	}
	// flowlog deleted from resource
	if len(oldFlowLogs) > 0 && len(newFlowLogs) == 0 {
		return d.ForceNew("flowlog")
	}

	if len(oldFlowLogs) > 0 && len(newFlowLogs) > 0 {
		oldFlowlogMap := oldFlowLogs[0].(map[string]any)
		newFlowLogMap := newFlowLogs[0].(map[string]any)
		// find the diff between the old and new value of the fields.
		// name should not force re-creation
		// all other values should force re-creation, case does not matter
		for k, v := range newFlowLogMap {
			if k != "name" && k != "id" {
				if !strings.EqualFold(fmt.Sprintf("%v", v), fmt.Sprintf("%v", oldFlowlogMap[k])) {
					// set ForceNew manually only if a field changes. only set on the field that changes, setting on the entire `flowlog` list does nothing
					if err := d.ForceNew("flowlog.0." + k); err != nil {
						return err
					}
				}
			}
		}
	}
	return nil
}
func resourceNicCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(bundleclient.SdkBundle).CloudApiClient
	ns := cloudapinic.Service{Client: client, Meta: meta, D: d}

	nic, err := cloudapinic.GetNicFromSchemaCreate(d, "")
	if err != nil {
		return diag.FromErr(fmt.Errorf("error occurred while getting nic from schema: %w", err))
	}

	dcID := d.Get("datacenter_id").(string)
	srvID := d.Get("server_id").(string)
	createdNic, _, err := ns.Create(ctx, dcID, srvID, nic)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error occurred while creating a nic: %w", err))
	}

	if createdNic.Id != nil {
		d.SetId(*createdNic.Id)

		if v, ok := d.GetOk("security_groups_ids"); ok {
			raw := v.(*schema.Set).List()
			nsgService := nsg.Service{Client: client, Meta: meta, D: d}
			if diagnostic := nsgService.PutNICNSG(ctx, dcID, srvID, d.Id(), raw); diagnostic != nil {
				return diagnostic
			}
		}
	}
	// Sometimes there is an error because the nic is not found after it's created.
	//Probably a read write consistency issue.
	//We're retrying for 5 minutes. 404 - means we keep on trying.
	// Sometimes there is an error because the nic is not found after it's created.
	// Probably a read write consistency issue.
	// We're retrying for 5 minutes. 404 - means we keep on trying.
	var foundNic = &ionoscloud.Nic{}
	var apiResponse = &ionoscloud.APIResponse{}
	err = retry.RetryContext(ctx, 5*time.Minute, func() *retry.RetryError {
		var err error
		foundNic, apiResponse, err = ns.Get(ctx, dcID, srvID, *createdNic.Id, 3)
		if apiResponse.HttpNotFound() {
			log.Printf("[INFO] Could not find nic with Id %s , retrying...", *createdNic.Id)
			return retry.RetryableError(fmt.Errorf("could not find nic, %w", err))
		}
		if err != nil {
			retry.NonRetryableError(err)
		}
		return nil
	})

	if err != nil {
		return diag.FromErr(err)
	}

	if foundNic == nil || *foundNic.Id == "" {
		return diag.FromErr(fmt.Errorf("could not find nic with id %s after creation ", *createdNic.Id))
	}

	return resourceNicRead(ctx, d, meta)
}

func resourceNicRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(bundleclient.SdkBundle).CloudApiClient
	ns := cloudapinic.Service{Client: client, Meta: meta, D: d}
	dcID := d.Get("datacenter_id").(string)
	srvID := d.Get("server_id").(string)
	nicID := d.Id()
	nic, apiResponse, err := ns.Get(ctx, dcID, srvID, nicID, 3)
	if err != nil {
		if apiResponse.HttpNotFound() {
			log.Printf("[INFO] nic resource with id %s not found", nicID)
			d.SetId("")
			return nil
		}
		return diag.FromErr(fmt.Errorf("error occurred while fetching a nic ID %s %w", nicID, err))
	}

	if err := cloudapinic.NicSetData(d, nic); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceNicUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(bundleclient.SdkBundle).CloudApiClient
	ns := cloudapinic.Service{Client: client, Meta: meta, D: d}
	dcID := d.Get("datacenter_id").(string)
	srvID := d.Get("server_id").(string)
	nicID := d.Id()
	var err error
	if d.HasChange("flowlog") {
		oldV, newV := d.GetChange("flowlog")
		var firstFlowLogId = ""
		if oldFlowLogs := oldV.([]any); len(oldFlowLogs) > 0 {
			firstFlowLogId = oldFlowLogs[0].(map[string]any)["id"].(string)
		}

		if newFlowLogs := newV.([]any); len(newFlowLogs) > 0 {
			flowLogMap := newFlowLogs[0].(map[string]any)
			flowLog := cloudapiflowlog.GetFlowlogFromMap(flowLogMap)
			fw := cloudapiflowlog.Service{
				D:      d,
				Client: client,
			}
			err = fw.CreateOrPatchForServer(ctx, dcID, srvID, nicID, firstFlowLogId, flowLog)
			if err != nil {
				// if we have a create that failed, we do not want to save in state
				// saving in state would mean a diff that would force a re-create
				if firstFlowLogId == "" {
					_ = d.Set("flowlog", nil)
				}
				return diag.FromErr(err)
			}
		}
	}

	nic, err := cloudapinic.GetNicFromSchema(d, "")
	if err != nil {
		return diag.FromErr(fmt.Errorf("update error occurred while getting nic from schema: %w", err))
	}

	_, _, err = ns.Update(ctx, dcID, srvID, nicID, *nic.Properties)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error occurred while updating a nic: %w", err))
	}
	if d.HasChange("security_groups_ids") {
		if v, ok := d.GetOk("security_groups_ids"); ok {
			raw := v.(*schema.Set).List()
			nsgService := nsg.Service{Client: client, Meta: meta, D: d}
			if diagnostic := nsgService.PutNICNSG(ctx, dcID, srvID, nicID, raw); diagnostic != nil {
				return diagnostic
			}
		}
	}

	return resourceNicRead(ctx, d, meta)
}

func resourceNicDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(bundleclient.SdkBundle).CloudApiClient
	ns := cloudapinic.Service{Client: client, Meta: meta, D: d}
	dcID := d.Get("datacenter_id").(string)
	srvID := d.Get("server_id").(string)
	nicID := d.Id()
	_, err := ns.Delete(ctx, dcID, srvID, nicID)
	if err != nil {
		return diag.FromErr(fmt.Errorf("an error occurred while deleting a nic dcID %s ID %s %w", dcID, nicID, err))
	}
	d.SetId("")
	return nil
}

func resourceNicImport(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	parts := strings.Split(d.Id(), "/")
	if len(parts) != 3 || parts[0] == "" || parts[1] == "" {
		return nil, fmt.Errorf("invalid import id %q. Expecting {datacenter}/{server}/{nic}", d.Id())
	}
	dcID := parts[0]
	srvID := parts[1]
	nicID := parts[2]

	client := meta.(bundleclient.SdkBundle).CloudApiClient

	nic, apiResponse, err := client.NetworkInterfacesApi.DatacentersServersNicsFindById(ctx, dcID, srvID, nicID).Execute()
	logApiRequestTime(apiResponse)

	if err != nil {
		if apiResponse.HttpNotFound() {
			d.SetId("")
			return nil, fmt.Errorf("nic does not exist %q", nicID)
		}

		return nil, fmt.Errorf("an error occurred while trying to fetch the nic %q, error: %w", nicID, err)

	}

	if err := d.Set("datacenter_id", dcID); err != nil {
		return nil, err
	}
	if err := d.Set("server_id", srvID); err != nil {
		return nil, err
	}

	if err := cloudapinic.NicSetData(d, &nic); err != nil {
		return nil, err
	}

	log.Printf("[INFO] nic found: %+v", nic)

	return []*schema.ResourceData{d}, nil
}
