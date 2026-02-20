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
			"location": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
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
	if d.Id() == "" {
		return nil
	}

	if d.HasChange("flowlog") {
		oldVal, newVal := d.GetChange("flowlog")
		oldFLowLogs := oldVal.([]any)
		newFlowLogs := newVal.([]any)

		// "this check is for the scenario in which we have 0 initial flowlogs and we add a new one during a nic update
		if (oldFLowLogs == nil || len(oldFLowLogs) == 0) && (newFlowLogs != nil || len(newFlowLogs) > 0) {
			return nil
		}
		// flowlog deleted from resource
		if (oldFLowLogs != nil || len(oldFLowLogs) > 0) && (newFlowLogs == nil || len(newFlowLogs) == 0) {
			return d.ForceNew("flowlog")
		}
		var oldFlowlogMap map[string]any
		var newFlowLogMap map[string]any
		if oldFLowLogs != nil && len(oldFLowLogs) > 0 {
			oldFlowlogMap = oldFLowLogs[0].(map[string]any)
		}
		if newFlowLogs != nil && len(newFlowLogs) > 0 {
			newFlowLogMap = newFlowLogs[0].(map[string]any)
		}
		// find the diff between the old and new value of the fields.
		// name should not force re-creation
		// all other values should force re-creation, case does not matter
		for k, v := range newFlowLogMap {
			if k != "name" && k != "id" {
				if !strings.EqualFold(strings.ToUpper(v.(string)), strings.ToUpper(oldFlowlogMap[k].(string))) {
					// set ForceNew manually only if a field changes. only set on the field that changes, setting on the entire `flowlog` list does nothing
					err := d.ForceNew("flowlog.0." + k)
					if err != nil {
						return err
					}
				}
			}
		}
		return nil
	}
	return nil
}
func resourceNicCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	location := d.Get("location").(string)
	client := meta.(bundleclient.SdkBundle).NewCloudAPIClient(location)
	ns := cloudapinic.Service{Client: client, Meta: meta, D: d}

	nic, err := cloudapinic.GetNicFromSchemaCreate(d, "")
	if err != nil {
		diags := diag.FromErr(fmt.Errorf("error occurred while getting nic from schema: %w", err))
		return diags
	}

	dcid := d.Get("datacenter_id").(string)
	srvid := d.Get("server_id").(string)
	createdNic, _, err := ns.Create(ctx, dcid, srvid, nic)
	if err != nil {
		diags := diag.FromErr(fmt.Errorf("error occurred while creating a nic: %w", err))
		return diags
	}

	if createdNic.Id != nil {
		d.SetId(*createdNic.Id)

		if v, ok := d.GetOk("security_groups_ids"); ok {
			raw := v.(*schema.Set).List()
			nsgService := nsg.Service{Client: client, Meta: meta, D: d}
			if diagnostic := nsgService.PutNICNSG(ctx, dcid, srvid, d.Id(), raw); diagnostic != nil {
				return diagnostic
			}
		}
	}
	// Sometimes there is an error because the nic is not found after it's created.
	// Probably a read write consistency issue.
	// We're retrying for 5 minutes. 404 - means we keep on trying.
	// Sometimes there is an error because the nic is not found after it's created.
	// Probably a read write consistency issue.
	// We're retrying for 5 minutes. 404 - means we keep on trying.
	var foundNic = &ionoscloud.Nic{}
	var apiResponse = &ionoscloud.APIResponse{}
	err = retry.RetryContext(ctx, 5*time.Minute, func() *retry.RetryError {
		var err error
		foundNic, apiResponse, err = ns.Get(ctx, dcid, srvid, *createdNic.Id, 3)
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
		return diag.FromErr(fmt.Errorf("could not find nic with id %s after creation ", *nic.Id))
	}

	return resourceNicRead(ctx, d, meta)
}

func resourceNicRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	location := d.Get("location").(string)
	client := meta.(bundleclient.SdkBundle).NewCloudAPIClient(location)

	ns := cloudapinic.Service{Client: client, Meta: meta, D: d}
	dcid := d.Get("datacenter_id").(string)
	srvid := d.Get("server_id").(string)
	nicid := d.Id()
	nic, apiResponse, err := ns.Get(ctx, dcid, srvid, nicid, 3)
	if err != nil {
		if apiResponse.HttpNotFound() {
			log.Printf("[INFO] nic resource with id %s not found", nicid)
			d.SetId("")
			return nil
		}
		diags := diag.FromErr(fmt.Errorf("error occurred while fetching a nic ID %s %w", d.Id(), err))
		return diags
	}

	if err := cloudapinic.NicSetData(d, nic); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceNicUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	location := d.Get("location").(string)
	client := meta.(bundleclient.SdkBundle).NewCloudAPIClient(location)

	ns := cloudapinic.Service{Client: client, Meta: meta, D: d}
	dcID := d.Get("datacenter_id").(string)
	srvID := d.Get("server_id").(string)
	nicID := d.Id()
	var err error
	if d.HasChange("flowlog") {
		oldV, newV := d.GetChange("flowlog")
		var firstFlowLogId = ""
		if oldV != nil && len(oldV.([]any)) > 0 {
			firstFlowLogId = oldV.([]any)[0].(map[string]any)["id"].(string)
		}

		if newV != nil && newV.([]any) != nil && len(newV.([]any)) > 0 {
			for _, val := range newV.([]any) {
				if flowLogMap, ok := val.(map[string]any); ok {
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
		}
	}

	nic, err := cloudapinic.GetNicFromSchema(d, "")
	if err != nil {
		diags := diag.FromErr(fmt.Errorf("update error occurred while getting nic from schema: %w", err))
		return diags
	}

	_, _, err = ns.Update(ctx, dcID, srvID, nicID, *nic.Properties)
	if err != nil {
		diags := diag.FromErr(fmt.Errorf("error occurred while updating a nic: %w", err))
		return diags
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
	location := d.Get("location").(string)
	client := meta.(bundleclient.SdkBundle).NewCloudAPIClient(location)

	ns := cloudapinic.Service{Client: client, Meta: meta, D: d}
	dcid := d.Get("datacenter_id").(string)
	srvid := d.Get("server_id").(string)
	nicid := d.Id()
	_, err := ns.Delete(ctx, dcid, srvid, nicid)
	if err != nil {
		diags := diag.FromErr(fmt.Errorf("an error occurred while deleting a nic dcId %s ID %s %s", d.Get("datacenter_id").(string), d.Id(), err))
		return diags
	}
	d.SetId("")
	return nil
}

func resourceNicImport(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	importID := d.Id()

	location, parts := splitImportID(importID, "/")
	if len(parts) != 3 {
		return nil, fmt.Errorf(
			"invalid import identifier: expected one of <location>:<datacenter>/<server>/<nic> or "+
				"<datacenter>/<server>/<nic>, got: %s", importID,
		)
	}

	if err := validateImportIDParts(parts); err != nil {
		return nil, fmt.Errorf("failed validating import identifier %q: %w", importID, err)
	}

	dcId := parts[0]
	sId := parts[1]
	nicId := parts[2]

	client := meta.(bundleclient.SdkBundle).NewCloudAPIClient(location)

	nic, apiResponse, err := client.NetworkInterfacesApi.DatacentersServersNicsFindById(ctx, dcId, sId, nicId).Execute()
	logApiRequestTime(apiResponse)

	if err != nil {
		if apiResponse.HttpNotFound() {
			d.SetId("")
			return nil, fmt.Errorf("lan does not exist%q", nicId)
		}

		return nil, fmt.Errorf("an error occurred while trying to fetch the nic %q, error:%w", nicId, err)

	}

	err = d.Set("datacenter_id", dcId)
	if err != nil {
		return nil, err
	}
	err = d.Set("server_id", sId)
	if err != nil {
		return nil, err
	}

	if err := cloudapinic.NicSetData(d, &nic); err != nil {
		return nil, err
	}

	log.Printf("[INFO] nic found: %+v", nic)

	return []*schema.ResourceData{d}, nil
}
