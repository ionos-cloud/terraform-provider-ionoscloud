package ionoscloud

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"

	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/bundleclient"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils/constant"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	ionoscloud "github.com/ionos-cloud/sdk-go/v6"
)

func resourceLan() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceLanCreate,
		ReadContext:   resourceLanRead,
		UpdateContext: resourceLanUpdate,
		DeleteContext: resourceLanDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceLanImport,
		},
		Schema: map[string]*schema.Schema{
			"public": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
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
			"pcc": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"ip_failover": {
				Type:     schema.TypeList,
				Computed: true,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ip": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"nic_uuid": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"ipv4_cidr_block": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "For public LANs this property is null, for private LANs it contains the private IPv4 CIDR range. This property is a read only property.",
			},
			"ipv6_cidr_block": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "IPv6 CIDR block assigned to the LAN. Can be set to 'AUTO' for an automatically assigned address or the address can be explicitly supplied.",
				// If a value has already been assigned by the backend, avoids reassignment if the field is set to AUTO.
				DiffSuppressFunc: func(_, old, new string, _ *schema.ResourceData) bool {
					if old != "" && new == "AUTO" {
						return true
					}
					return false
				},
			},
		},
		Timeouts: &resourceDefaultTimeouts,
	}
}

func resourceLanCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	location := d.Get("location").(string)
	client := meta.(bundleclient.SdkBundle).NewCloudAPIClient(location)

	public := d.Get("public").(bool)
	request := ionoscloud.Lan{
		Properties: &ionoscloud.LanProperties{
			Public: &public,
		},
	}

	if d.Get("name") != nil {
		name := d.Get("name").(string)
		request.Properties.Name = &name
	}

	if d.Get("pcc") != nil && d.Get("pcc").(string) != "" {
		pccID := d.Get("pcc").(string)
		log.Printf("[INFO] Setting PCC for LAN %s to %s...", d.Id(), pccID)
		request.Properties.Pcc = &pccID
	}

	if d.Get("ipv6_cidr_block") != nil {
		ipv6 := d.Get("ipv6_cidr_block").(string)
		log.Printf("[INFO] Setting ipv6CidrBlock for LAN %s to %s...", d.Id(), ipv6)
		request.Properties.Ipv6CidrBlock = &ipv6
	} else {
		request.Properties.SetIpv6CidrBlockNil()
	}

	dcid := d.Get("datacenter_id").(string)
	rsp, apiResponse, err := client.LANsApi.DatacentersLansPost(ctx, dcid).Lan(request).Execute()
	logApiRequestTime(apiResponse)

	if err != nil {
		d.SetId("")
		diags := diag.FromErr(fmt.Errorf("an error occurred while creating LAN: %w", err))
		return diags
	}

	d.SetId(*rsp.Id)

	log.Printf("[INFO] LAN ID: %s", d.Id())

	if errState := bundleclient.WaitForStateChange(ctx, meta, d, apiResponse, schema.TimeoutCreate); errState != nil {
		if errState != nil {
			if bundleclient.IsRequestFailed(errState) {
				d.SetId("")
			}
		}
		return diag.FromErr(errState)
	}

	for {
		log.Printf("[INFO] Waiting for LAN %s to be available...", *rsp.Id)

		lanReady, rsErr := lanAvailable(ctx, client, d)

		if rsErr != nil {
			diags := diag.FromErr(fmt.Errorf("error while checking readiness status of LAN %s: %w", *rsp.Id, rsErr))
			return diags
		}

		if lanReady {
			log.Printf("[INFO] LAN ready: %s", d.Id())
			break
		}

		select {
		case <-time.After(constant.SleepInterval):
			log.Printf("[INFO] trying again ...")
		case <-ctx.Done():
			log.Printf("[INFO] lan creation timed out")
			diags := diag.FromErr(fmt.Errorf("lan creation timed out! WARNING: your lan will still probably be created after some time but the terraform state won't reflect that; check your Ionos Cloud account for updates"))
			return diags
		}
	}

	return resourceLanRead(ctx, d, meta)
}

func resourceLanRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	location := d.Get("location").(string)
	client := meta.(bundleclient.SdkBundle).NewCloudAPIClient(location)

	dcid := d.Get("datacenter_id").(string)

	lan, apiResponse, err := client.LANsApi.DatacentersLansFindById(ctx, dcid, d.Id()).Execute()
	logApiRequestTime(apiResponse)

	if err != nil {
		if httpNotFound(apiResponse) {
			log.Printf("[INFO] LAN %s not found", d.Id())
			d.SetId("")
			return nil
		}

		diags := diag.FromErr(fmt.Errorf("an error occurred while fetching a LAN %s: %w", d.Id(), err))
		return diags
	}

	log.Printf("[INFO] LAN %s found: %+v", d.Id(), lan)

	if err := setLanData(d, &lan); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceLanUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	location := d.Get("location").(string)
	client := meta.(bundleclient.SdkBundle).NewCloudAPIClient(location)

	properties := &ionoscloud.LanProperties{}
	newValue := d.Get("public")
	public := newValue.(bool)
	properties.Public = &public

	if d.HasChange("name") {
		_, newValue := d.GetChange("name")
		name := newValue.(string)
		properties.Name = &name
	}

	if d.HasChange("pcc") {
		_, newPCC := d.GetChange("pcc")

		if newPCC != nil && newPCC.(string) != "" {
			log.Printf("[INFO] Setting PCC for LAN %s to %s...", d.Id(), newPCC.(string))
			pcc := newPCC.(string)
			properties.Pcc = &pcc
		}
	}

	if d.HasChange("ipv6_cidr_block") {
		_, newIpv6 := d.GetChange("ipv6_cidr_block")
		if newIpv6 != nil && newIpv6.(string) != "" {
			log.Printf("[INFO] Setting ipv6CidrBlock for LAN %s to %s...", d.Id(), newIpv6.(string))
			ipv6 := newIpv6.(string)
			properties.Ipv6CidrBlock = &ipv6
		} else {
			properties.SetIpv6CidrBlockNil()
		}
	}

	dcid := d.Get("datacenter_id").(string)

	_, apiResponse, err := client.LANsApi.DatacentersLansPatch(ctx, dcid, d.Id()).Lan(*properties).Execute()
	logApiRequestTime(apiResponse)

	if err != nil {
		diags := diag.FromErr(fmt.Errorf("an error occurred while patching a lan ID %s %w", d.Id(), err))
		return diags
	}

	if errState := bundleclient.WaitForStateChange(ctx, meta, d, apiResponse, schema.TimeoutUpdate); errState != nil {
		return diag.FromErr(errState)
	}

	return resourceLanRead(ctx, d, meta)
}

func resourceLanDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	dcId := d.Get("datacenter_id").(string)
	location := d.Get("location").(string)
	client := meta.(bundleclient.SdkBundle).NewCloudAPIClient(location)

	if err := waitForLanNicsDeletion(ctx, client, d); err != nil {
		return diag.FromErr(err)
	}

	err := retry.RetryContext(ctx, d.Timeout(schema.TimeoutDelete), func() *retry.RetryError {
		apiResponse, err := client.LANsApi.DatacentersLansDelete(ctx, dcId, d.Id()).Execute()
		if err != nil {
			if isDeleteProtected(apiResponse, err.Error()) {
				log.Printf("[INFO] LAN %s is delete-protected, keep trying", d.Id())
				return retry.RetryableError(fmt.Errorf("lan %s is delete-protected, keep trying %w", d.Id(), err))
			}
			if httpNotFound(apiResponse) {
				log.Printf("[INFO] LAN already deleted %s", d.Id())
				return nil
			}
			return retry.NonRetryableError(fmt.Errorf("an error occurred while deleting lan dcId %s ID %s %w", dcId, d.Id(), err))
		}
		log.Printf("[DEBUG] deletion started for LAN with ID: %v", d.Id())
		return nil
	})
	if err != nil {
		return diag.FromErr(fmt.Errorf("an error occurred while deleting lan dcId %s ID %s %w", dcId, d.Id(), err))
	}

	if err := waitForLanDeletion(ctx, client, d); err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")
	return nil
}
func isDeleteProtected(apiResponse *ionoscloud.APIResponse, errMessage string) bool {
	if apiResponse != nil && apiResponse.Response != nil && apiResponse.StatusCode == 403 {
		if strings.Contains(errMessage, "is delete-protected by") {
			return true
		}
	}
	return false
}

func resourceLanImport(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	importID := d.Id()

	location, parts := splitImportID(importID, "/")
	if len(parts) != 2 {
		return nil, fmt.Errorf(
			"invalid import identifier: expected one of <location>:<datacenter-id>/<lan-id> "+
				"or <datacenter-id>/<lan-id>, got: %s", importID,
		)
	}

	if err := validateImportIDParts(parts); err != nil {
		return nil, fmt.Errorf("failed validating import identifier %q: %w", importID, err)
	}

	datacenterId := parts[0]
	lanId := parts[1]

	client := meta.(bundleclient.SdkBundle).NewCloudAPIClient(location)

	lan, apiResponse, err := client.LANsApi.DatacentersLansFindById(ctx, datacenterId, lanId).Execute()
	logApiRequestTime(apiResponse)

	if err != nil {
		if httpNotFound(apiResponse) {
			d.SetId("")
			return nil, fmt.Errorf("unable to find lan %q", lanId)
		}
		return nil, fmt.Errorf("an error occurred while retrieving the lan %q, %w", lanId, err)
	}

	log.Printf("[INFO] LAN %s found: %+v", d.Id(), lan)

	if err := d.Set("datacenter_id", datacenterId); err != nil {
		return nil, fmt.Errorf("error while setting datacenter_id property for lan %q: %w", lanId, err)
	}

	if err := setLanData(d, &lan); err != nil {
		return nil, err
	}
	return []*schema.ResourceData{d}, nil
}

func setLanData(d *schema.ResourceData, lan *ionoscloud.Lan) error {
	d.SetId(*lan.Id)

	if lan.Properties != nil {
		if lan.Properties.Name != nil {
			if err := d.Set("name", *lan.Properties.Name); err != nil {
				return err
			}
		}
		if lan.Properties.IpFailover != nil && len(*lan.Properties.IpFailover) > 0 {
			if err := d.Set("ip_failover", convertIpFailoverList(lan.Properties.IpFailover)); err != nil {
				return err
			}
		}
		if lan.Properties.Pcc != nil {
			if err := d.Set("pcc", *lan.Properties.Pcc); err != nil {
				return err
			}
		}
		if lan.Properties.Public != nil {
			if err := d.Set("public", *lan.Properties.Public); err != nil {
				return err
			}
		}

		if lan.Properties.Ipv4CidrBlock != nil {
			if err := d.Set("ipv4_cidr_block", *lan.Properties.Ipv4CidrBlock); err != nil {
				return utils.GenerateSetError("lan", "ipv4_cidr_block", err)
			}
		}

		if lan.Properties.Ipv6CidrBlock != nil {
			if err := d.Set("ipv6_cidr_block", *lan.Properties.Ipv6CidrBlock); err != nil {
				return utils.GenerateSetError("lan", "ipv6_cidr_block", err)
			}
		}
	}

	return nil
}

func lanAvailable(ctx context.Context, client *ionoscloud.APIClient, d *schema.ResourceData) (bool, error) {
	dcId := d.Get("datacenter_id").(string)
	rsp, apiResponse, err := client.LANsApi.DatacentersLansFindById(ctx, dcId, d.Id()).Execute()
	logApiRequestTime(apiResponse)

	if err != nil {
		return true, fmt.Errorf("error checking LAN status: %w", err)
	}

	if rsp.Metadata == nil || rsp.Metadata.State == nil {
		return false, fmt.Errorf("could not retrieve state of lan %s", d.Id())
	}

	return strings.EqualFold(*rsp.Metadata.State, constant.Available), nil
}

func lanDeleted(ctx context.Context, client *ionoscloud.APIClient, d *schema.ResourceData) (bool, error) {
	dcId := d.Get("datacenter_id").(string)

	_, apiResponse, err := client.LANsApi.DatacentersLansFindById(ctx, dcId, d.Id()).Execute()
	logApiRequestTime(apiResponse)

	if err != nil {
		if isDeleteProtected(apiResponse, err.Error()) {
			log.Printf("[INFO] LAN %s is delete-protected, keep trying", d.Id())
			return false, nil
		}
		if httpNotFound(apiResponse) {
			log.Printf("[INFO] LAN deleted %s", d.Id())
			return true, nil
		}
		return false, fmt.Errorf("error checking LAN deletion status: %w", err)
	}

	log.Printf("[INFO] LAN %s not yet deleted from the datacenter %s", d.Id(), dcId)

	return false, nil
}

func waitForLanDeletion(ctx context.Context, client *ionoscloud.APIClient, d *schema.ResourceData) error {
	for {
		log.Printf("[INFO] waiting for LAN %s to be deleted...", d.Id())

		lDeleted, dsErr := lanDeleted(ctx, client, d)

		if dsErr != nil {
			return fmt.Errorf("error while checking deletion status of LAN %s: %w", d.Id(), dsErr)
		}

		if lDeleted {
			log.Printf("[INFO] successfully deleted LAN: %s", d.Id())
			break
		}

		select {
		case <-time.After(constant.SleepInterval):
			log.Printf("[INFO] trying again ...")
		case <-ctx.Done():
			log.Printf("[INFO] lan deletion timed out")
			return fmt.Errorf("lan deletion timed out! WARNING: your lan will still probably be deleted after some time but the terraform state won't reflect that; check your Ionos Cloud account for updates")
		}
	}
	return nil
}

func lanNicsDeleted(ctx context.Context, client *ionoscloud.APIClient, d *schema.ResourceData) (bool, error) {
	dcId := d.Get("datacenter_id").(string)

	nics, apiResponse, err := client.LANsApi.DatacentersLansNicsGet(ctx, dcId, d.Id()).Execute()
	logApiRequestTime(apiResponse)

	if err != nil {
		return false, fmt.Errorf("an error occurred while searching for nics in datacenter with id: %s for lan with: id %s %w", dcId, d.Id(), err)
	}

	if nics.Items != nil && len(*nics.Items) > 0 {
		log.Printf("[INFO] there are still nics under LAN  with id %s", d.Id())
		return false, nil
	}

	return true, nil
}

func waitForLanNicsDeletion(ctx context.Context, client *ionoscloud.APIClient, d *schema.ResourceData) error {
	for {
		log.Printf("[INFO] waiting for nics under LAN %s to be deleted...", d.Id())

		nicsDeleted, dsErr := lanNicsDeleted(ctx, client, d)

		if dsErr != nil {
			return fmt.Errorf("error while checking nics under lan %s: %w", d.Id(), dsErr)
		}

		if nicsDeleted {
			log.Printf("[INFO] no nics under LAN: %s", d.Id())
			break
		}

		select {
		case <-time.After(constant.SleepInterval):
			log.Printf("[INFO] trying again ...")
		case <-ctx.Done():
			log.Printf("[INFO] nics deletion check timed out")
			return fmt.Errorf("nics deletion check timed out! WARNING: your lan nics may still be deleted; check your Ionos Cloud account for updates and perform again a destroy for remaining resources")
		}
	}
	return nil
}
