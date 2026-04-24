package ionoscloud

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"

	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/bundleclient"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils/constant"
	diagutil "github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils/diags"

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
				Type:        schema.TypeString,
				Description: "The location of the resource. This field should be used only if you are also using a file configuration and should not be configured otherwise.",
				Optional:    true,
				ForceNew:    true,
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
	client, err := meta.(bundleclient.SdkBundle).NewCloudAPIClient(location)
	if err != nil {
		return diag.FromErr(err)
	}

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
		tflog.Info(ctx, "setting PCC for LAN", map[string]interface{}{"lan_id": d.Id(), "pcc_id": pccID})
		request.Properties.Pcc = &pccID
	}

	if d.Get("ipv6_cidr_block") != nil {
		ipv6 := d.Get("ipv6_cidr_block").(string)
		tflog.Info(ctx, "setting ipv6CidrBlock for LAN", map[string]interface{}{"lan_id": d.Id(), "ipv6_cidr_block": ipv6})
		request.Properties.Ipv6CidrBlock = &ipv6
	} else {
		request.Properties.SetIpv6CidrBlockNil()
	}

	dcid := d.Get("datacenter_id").(string)
	rsp, apiResponse, err := client.LANsApi.DatacentersLansPost(ctx, dcid).Lan(request).Execute()
	logApiRequestTime(apiResponse)

	if err != nil {
		d.SetId("")
		requestLocation, _ := apiResponse.SafeLocation()
		return diagutil.ToDiags(d, fmt.Errorf("an error occurred while creating LAN: %w", err), &diagutil.ErrorContext{RequestID: diagutil.ExtractRequestID(requestLocation), StatusCode: apiResponse.SafeStatusCode()})
	}

	d.SetId(*rsp.Id)

	tflog.Info(ctx, "LAN created", map[string]interface{}{"lan_id": d.Id()})

	if errState := bundleclient.WaitForStateChange(ctx, meta, d, apiResponse, schema.TimeoutCreate); errState != nil {
		if errState != nil {
			if bundleclient.IsRequestFailed(errState) {
				d.SetId("")
			}
		}
		requestLocation, _ := apiResponse.SafeLocation()
		return diagutil.ToDiags(d, errState, &diagutil.ErrorContext{Timeout: d.Timeout(schema.TimeoutCreate).String(), RequestID: diagutil.ExtractRequestID(requestLocation)})
	}

	for {
		tflog.Info(ctx, "waiting for LAN to be available", map[string]interface{}{"lan_id": *rsp.Id})

		lanReady, rsErr := lanAvailable(ctx, client, d)

		if rsErr != nil {
			return diagutil.ToDiags(d, fmt.Errorf("error while checking readiness status of LAN %s: %w", *rsp.Id, rsErr), nil)
		}

		if lanReady {
			tflog.Info(ctx, "LAN ready", map[string]interface{}{"lan_id": d.Id()})
			break
		}

		select {
		case <-time.After(constant.SleepInterval):
			tflog.Info(ctx, "LAN not yet ready, retrying")
		case <-ctx.Done():
			tflog.Info(ctx, "LAN creation timed out")
			return diagutil.ToDiags(d, fmt.Errorf("lan creation timed out! WARNING: your lan will still probably be created after some time but the terraform state won't reflect that; check your Ionos Cloud account for updates"), nil)
		}
	}

	return resourceLanRead(ctx, d, meta)
}

func resourceLanRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	location := d.Get("location").(string)
	client, err := meta.(bundleclient.SdkBundle).NewCloudAPIClient(location)
	if err != nil {
		return diag.FromErr(err)
	}

	dcid := d.Get("datacenter_id").(string)

	lan, apiResponse, err := client.LANsApi.DatacentersLansFindById(ctx, dcid, d.Id()).Execute()
	logApiRequestTime(apiResponse)

	if err != nil {
		if httpNotFound(apiResponse) {
			tflog.Info(ctx, "LAN not found", map[string]interface{}{"lan_id": d.Id()})
			d.SetId("")
			return nil
		}

		return diagutil.ToDiags(d, fmt.Errorf("an error occurred while fetching a LAN: %w", err), nil)
	}

	tflog.Info(ctx, "LAN found", map[string]interface{}{"lan_id": d.Id()})

	if err := setLanData(d, &lan); err != nil {
		return diagutil.ToDiags(d, err, nil)
	}

	return nil
}

func resourceLanUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	location := d.Get("location").(string)
	client, err := meta.(bundleclient.SdkBundle).NewCloudAPIClient(location)
	if err != nil {
		return diag.FromErr(err)
	}

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
			tflog.Info(ctx, "setting PCC for LAN", map[string]interface{}{"lan_id": d.Id(), "pcc_id": newPCC.(string)})
			pcc := newPCC.(string)
			properties.Pcc = &pcc
		}
	}

	if d.HasChange("ipv6_cidr_block") {
		_, newIpv6 := d.GetChange("ipv6_cidr_block")
		if newIpv6 != nil && newIpv6.(string) != "" {
			tflog.Info(ctx, "setting ipv6CidrBlock for LAN", map[string]interface{}{"lan_id": d.Id(), "ipv6_cidr_block": newIpv6.(string)})
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
		requestLocation, _ := apiResponse.SafeLocation()
		return diagutil.ToDiags(d, fmt.Errorf("an error occurred while patching a lan: %w", err), &diagutil.ErrorContext{RequestID: diagutil.ExtractRequestID(requestLocation), StatusCode: apiResponse.SafeStatusCode()})
	}

	if errState := bundleclient.WaitForStateChange(ctx, meta, d, apiResponse, schema.TimeoutUpdate); errState != nil {
		requestLocation, _ := apiResponse.SafeLocation()
		return diagutil.ToDiags(d, errState, &diagutil.ErrorContext{Timeout: d.Timeout(schema.TimeoutUpdate).String(), RequestID: diagutil.ExtractRequestID(requestLocation)})
	}

	return resourceLanRead(ctx, d, meta)
}

func resourceLanDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	dcId := d.Get("datacenter_id").(string)
	location := d.Get("location").(string)
	client, err := meta.(bundleclient.SdkBundle).NewCloudAPIClient(location)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := waitForLanNicsDeletion(ctx, client, d); err != nil {
		return diagutil.ToDiags(d, err, nil)
	}

	err = retry.RetryContext(ctx, d.Timeout(schema.TimeoutDelete), func() *retry.RetryError {
		apiResponse, err := client.LANsApi.DatacentersLansDelete(ctx, dcId, d.Id()).Execute()
		if err != nil {
			if isDeleteProtected(apiResponse, err.Error()) {
				tflog.Info(ctx, "LAN is delete-protected, keep trying", map[string]interface{}{"lan_id": d.Id()})
				return retry.RetryableError(fmt.Errorf("lan %s is delete-protected, keep trying %w", d.Id(), err))
			}
			if httpNotFound(apiResponse) {
				tflog.Info(ctx, "LAN already deleted", map[string]interface{}{"lan_id": d.Id()})
				return nil
			}
			return retry.NonRetryableError(fmt.Errorf("an error occurred while deleting lan dcId %s ID %s %w", dcId, d.Id(), err))
		}
		tflog.Debug(ctx, "deletion started for LAN", map[string]interface{}{"lan_id": d.Id()})
		return nil
	})
	if err != nil {
		return diagutil.ToDiags(d, fmt.Errorf("an error occurred while deleting lan dcId %s %w", dcId, err), nil)
	}

	if err := waitForLanDeletion(ctx, client, d); err != nil {
		return diagutil.ToDiags(d, err, nil)
	}

	d.SetId("")
	return nil
}
func isDeleteProtected(apiResponse *ionoscloud.APIResponse, errMessage string) bool {
	if apiResponse.SafeStatusCode() == 403 {
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
		return nil, diagutil.ToError(d, fmt.Errorf("failed validating import identifier %q: %w", importID, err), nil)
	}

	datacenterId := parts[0]
	lanId := parts[1]

	client, err := meta.(bundleclient.SdkBundle).NewCloudAPIClient(location)
	if err != nil {
		return nil, err
	}

	lan, apiResponse, err := client.LANsApi.DatacentersLansFindById(ctx, datacenterId, lanId).Execute()
	logApiRequestTime(apiResponse)

	if err != nil {
		if httpNotFound(apiResponse) {
			d.SetId("")
			return nil, diagutil.ToError(d, fmt.Errorf("unable to find lan %q", lanId), nil)
		}
		return nil, diagutil.ToError(d, fmt.Errorf("an error occurred while retrieving the lan %q, %w", lanId, err), nil)
	}

	tflog.Info(ctx, "LAN imported", map[string]interface{}{"lan_id": d.Id(), "datacenter_id": datacenterId})

	if err := d.Set("datacenter_id", datacenterId); err != nil {
		return nil, diagutil.ToError(d, fmt.Errorf("error while setting datacenter_id property for lan %q: %w", lanId, err), nil)
	}
	if err := d.Set("location", location); err != nil {
		return nil, err
	}

	if err := setLanData(d, &lan); err != nil {
		return nil, diagutil.ToError(d, err, nil)
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
			tflog.Info(ctx, "LAN is delete-protected, keep trying", map[string]interface{}{"lan_id": d.Id()})
			return false, nil
		}
		if httpNotFound(apiResponse) {
			tflog.Info(ctx, "LAN deleted", map[string]interface{}{"lan_id": d.Id()})
			return true, nil
		}
		return false, fmt.Errorf("error checking LAN deletion status: %w", err)
	}

	tflog.Info(ctx, "LAN not yet deleted from the datacenter", map[string]interface{}{"lan_id": d.Id(), "datacenter_id": dcId})

	return false, nil
}

func waitForLanDeletion(ctx context.Context, client *ionoscloud.APIClient, d *schema.ResourceData) error {
	for {
		tflog.Info(ctx, "waiting for LAN to be deleted", map[string]interface{}{"lan_id": d.Id()})

		lDeleted, dsErr := lanDeleted(ctx, client, d)

		if dsErr != nil {
			return fmt.Errorf("error while checking deletion status of LAN %s: %w", d.Id(), dsErr)
		}

		if lDeleted {
			tflog.Info(ctx, "successfully deleted LAN", map[string]interface{}{"lan_id": d.Id()})
			break
		}

		select {
		case <-time.After(constant.SleepInterval):
			tflog.Info(ctx, "LAN not yet deleted, retrying")
		case <-ctx.Done():
			tflog.Info(ctx, "LAN deletion timed out")
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
		tflog.Info(ctx, "nics still present under LAN", map[string]interface{}{"lan_id": d.Id()})
		return false, nil
	}

	return true, nil
}

func waitForLanNicsDeletion(ctx context.Context, client *ionoscloud.APIClient, d *schema.ResourceData) error {
	for {
		tflog.Info(ctx, "waiting for nics under LAN to be deleted", map[string]interface{}{"lan_id": d.Id()})

		nicsDeleted, dsErr := lanNicsDeleted(ctx, client, d)

		if dsErr != nil {
			return fmt.Errorf("error while checking nics under lan %s: %w", d.Id(), dsErr)
		}

		if nicsDeleted {
			tflog.Info(ctx, "no nics under LAN", map[string]interface{}{"lan_id": d.Id()})
			break
		}

		select {
		case <-time.After(constant.SleepInterval):
			tflog.Info(ctx, "nics still present under LAN, retrying")
		case <-ctx.Done():
			tflog.Info(ctx, "nics deletion check timed out")
			return fmt.Errorf("nics deletion check timed out! WARNING: your lan nics may still be deleted; check your Ionos Cloud account for updates and perform again a destroy for remaining resources")
		}
	}
	return nil
}
