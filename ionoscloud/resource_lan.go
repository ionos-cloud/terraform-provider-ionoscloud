package ionoscloud

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
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
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.All(validation.StringIsNotWhiteSpace),
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
		},
		Timeouts: &resourceDefaultTimeouts,
	}
}

func resourceLanCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(SdkBundle).CloudApiClient
	public := d.Get("public").(bool)
	request := ionoscloud.LanPost{
		Properties: &ionoscloud.LanPropertiesPost{
			Public: &public,
		},
	}

	name := d.Get("name").(string)
	log.Printf("[DEBUG] NAME %s", d.Get("name"))
	if d.Get("name") != nil {
		request.Properties.Name = &name
	}

	if d.Get("pcc") != nil && d.Get("pcc").(string) != "" {
		pccID := d.Get("pcc").(string)
		log.Printf("[INFO] Setting PCC for LAN %s to %s...", d.Id(), pccID)
		request.Properties.Pcc = &pccID
	}

	dcid := d.Get("datacenter_id").(string)
	rsp, apiResponse, err := client.LANsApi.DatacentersLansPost(ctx, dcid).Lan(request).Execute()
	logApiRequestTime(apiResponse)

	if err != nil {
		d.SetId("")
		diags := diag.FromErr(fmt.Errorf("an error occured while creating LAN: %w", err))
		return diags
	}

	d.SetId(*rsp.Id)

	log.Printf("[INFO] LAN ID: %s", d.Id())

	// Wait, catching any errors
	_, errState := getStateChangeConf(meta, d, apiResponse.Header.Get("Location"), schema.TimeoutCreate).WaitForStateContext(ctx)
	if errState != nil {
		if IsRequestFailed(err) {
			// Request failed, so resource was not created, delete resource from state file
			d.SetId("")
		}
		diags := diag.FromErr(errState)
		return diags
	}

	for {
		log.Printf("[INFO] Waiting for LAN %s to be available...", *rsp.Id)

		clusterReady, rsErr := lanAvailable(ctx, client, d)

		if rsErr != nil {
			diags := diag.FromErr(fmt.Errorf("error while checking readiness status of LAN %s: %w", *rsp.Id, rsErr))
			return diags
		}

		if clusterReady {
			log.Printf("[INFO] LAN ready: %s", d.Id())
			break
		}

		select {
		case <-time.After(SleepInterval):
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
	client := meta.(SdkBundle).CloudApiClient

	dcid := d.Get("datacenter_id").(string)

	lan, apiResponse, err := client.LANsApi.DatacentersLansFindById(ctx, dcid, d.Id()).Execute()
	logApiRequestTime(apiResponse)

	if err != nil {
		if httpNotFound(apiResponse) {
			log.Printf("[INFO] LAN %s not found", d.Id())
			d.SetId("")
			return nil
		}

		diags := diag.FromErr(fmt.Errorf("an error occured while fetching a LAN %s: %w", d.Id(), err))
		return diags
	}

	log.Printf("[INFO] LAN %s found: %+v", d.Id(), lan)

	if err := setLanData(d, &lan); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceLanUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(SdkBundle).CloudApiClient
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

	dcid := d.Get("datacenter_id").(string)

	_, apiResponse, err := client.LANsApi.DatacentersLansPatch(ctx, dcid, d.Id()).Lan(*properties).Execute()
	logApiRequestTime(apiResponse)

	if err != nil {
		diags := diag.FromErr(fmt.Errorf("an error occured while patching a lan ID %s %w", d.Id(), err))
		return diags
	}

	_, errState := getStateChangeConf(meta, d, apiResponse.Header.Get("Location"), schema.TimeoutUpdate).WaitForStateContext(ctx)
	if errState != nil {
		diags := diag.FromErr(errState)
		return diags
	}

	return resourceLanRead(ctx, d, meta)
}

func resourceLanDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(SdkBundle).CloudApiClient
	dcId := d.Get("datacenter_id").(string)

	if err := waitForLanNicsDeletion(ctx, client, d); err != nil {
		return diag.FromErr(err)
	}

	apiResponse, err := client.LANsApi.DatacentersLansDelete(ctx, dcId, d.Id()).Execute()
	logApiRequestTime(apiResponse)

	if err != nil {
		diags := diag.FromErr(fmt.Errorf("an error occured while deleting lan dcId %s ID %s %w", dcId, d.Id(), err))
		return diags

	}

	if err := waitForLanDeletion(ctx, client, d); err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")
	return nil
}

func resourceLanImport(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	client := meta.(SdkBundle).CloudApiClient

	parts := strings.Split(d.Id(), "/")

	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return nil, fmt.Errorf("invalid import id %q. Expecting {datacenter}/{lan}", d.Id())
	}

	datacenterId := parts[0]
	lanId := parts[1]

	lan, apiResponse, err := client.LANsApi.DatacentersLansFindById(ctx, datacenterId, lanId).Execute()
	logApiRequestTime(apiResponse)

	if err != nil {
		if httpNotFound(apiResponse) {
			d.SetId("")
			return nil, fmt.Errorf("unable to find lan %q", lanId)
		}
		return nil, fmt.Errorf("an error occured while retrieving the lan %q, %w", lanId, err)
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

	return *rsp.Metadata.State == "AVAILABLE", nil
}

func lanDeleted(ctx context.Context, client *ionoscloud.APIClient, d *schema.ResourceData) (bool, error) {
	dcid := d.Get("datacenter_id").(string)

	rsp, apiResponse, err := client.LANsApi.DatacentersLansFindById(ctx, dcid, d.Id()).Execute()
	logApiRequestTime(apiResponse)

	if err != nil {
		if httpNotFound(apiResponse) {
			log.Printf("[INFO] LAN deleted %s", d.Id())
			return true, nil
		}
		return true, fmt.Errorf("error checking LAN deletion status: %w", err)
	}
	if rsp.Metadata != nil && rsp.Metadata.State != nil {
		log.Printf("[INFO] LAN %s not deleted yet; LAN status: %+v", d.Id(), *rsp.Metadata.State)
	}

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
		case <-time.After(SleepInterval):
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
		return false, fmt.Errorf("an error occured while searching for nics in datacenter with id: %s for lan with: id %s %w", dcId, d.Id(), err)
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
		case <-time.After(SleepInterval):
			log.Printf("[INFO] trying again ...")
		case <-ctx.Done():
			log.Printf("[INFO] nics deletion check timed out")
			return fmt.Errorf("nics deletion check timed out! WARNING: your lan nics may still be deleted; check your Ionos Cloud account for updates and perform again a destroy for remaining resources")
		}
	}
	return nil
}
