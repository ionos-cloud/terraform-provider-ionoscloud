package ionoscloud

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	ionoscloud "github.com/ionos-cloud/sdk-go/v5"
)

func resourceLan() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceLanCreate,
		ReadContext:   resourceLanRead,
		UpdateContext: resourceLanUpdate,
		DeleteContext: resourceLanDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceResourceImport,
		},
		Schema: map[string]*schema.Schema{

			"public": {
				Type:     schema.TypeBool,
				Required: true,
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
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ip": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"nic_uuid": {
							Type:     schema.TypeString,
							Optional: true,
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
	client := meta.(*ionoscloud.APIClient)
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
	rsp, apiResponse, err := client.LanApi.DatacentersLansPost(ctx, dcid).Lan(request).Execute()

	if err != nil {
		d.SetId("")
		diags := diag.FromErr(fmt.Errorf("an error occured while creating LAN: %s", err))
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
		time.Sleep(5 * time.Second)

		clusterReady, rsErr := lanAvailable(ctx, client, d)

		if rsErr != nil {
			diags := diag.FromErr(fmt.Errorf("error while checking readiness status of LAN %s: %s", *rsp.Id, rsErr))
			return diags
		}

		if clusterReady {
			log.Printf("[INFO] LAN ready: %s", d.Id())
			break
		}
	}

	return resourceLanRead(ctx, d, meta)
}

func resourceLanRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*ionoscloud.APIClient)

	dcid := d.Get("datacenter_id").(string)
	rsp, apiResponse, err := client.LanApi.DatacentersLansFindById(ctx, dcid, d.Id()).Execute()

	if err != nil {
		if apiResponse != nil && apiResponse.StatusCode == 404 {
			log.Printf("[INFO] LAN %s not found", d.Id())
			d.SetId("")
			return nil
		}
		diags := diag.FromErr(fmt.Errorf("an error occured while fetching a LAN %s: %s", d.Id(), err))
		return diags
	}

	if rsp.Properties.Public != nil {
		if err := d.Set("public", *rsp.Properties.Public); err != nil {
			diags := diag.FromErr(err)
			return diags
		}
	}

	if rsp.Properties.Name != nil {
		if err := d.Set("name", *rsp.Properties.Name); err != nil {
			diags := diag.FromErr(err)
			return diags
		}
	}

	if rsp.Properties.IpFailover != nil && len(*rsp.Properties.IpFailover) > 0 {
		if err := d.Set("ip_failover", []map[string]string{
			{
				"ip":       *(*rsp.Properties.IpFailover)[0].Ip,
				"nic_uuid": *(*rsp.Properties.IpFailover)[0].NicUuid,
			},
		}); err != nil {
			diags := diag.FromErr(err)
			return diags
		}
	}

	if err := d.Set("datacenter_id", d.Get("datacenter_id").(string)); err != nil {
		diags := diag.FromErr(err)
		return diags
	}

	if rsp.Properties.Pcc != nil {
		if err := d.Set("pcc", *rsp.Properties.Pcc); err != nil {
			diags := diag.FromErr(err)
			return diags
		}
	}
	log.Printf("[INFO] LAN %s found: %+v", d.Id(), rsp)
	return nil
}

func resourceLanUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*ionoscloud.APIClient)
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
	_, apiResponse, err := client.LanApi.DatacentersLansPatch(ctx, dcid, d.Id()).Lan(*properties).Execute()
	if err != nil {
		diags := diag.FromErr(fmt.Errorf("an error occured while patching a lan ID %s %s", d.Id(), err))
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
	client := meta.(*ionoscloud.APIClient)
	dcid := d.Get("datacenter_id").(string)

	_, _, err := client.LanApi.DatacentersLansDelete(ctx, dcid, d.Id()).Execute()

	if err != nil {
		//try again in 120 seconds
		time.Sleep(120 * time.Second)
		_, apiResponse, err := client.LanApi.DatacentersLansDelete(ctx, dcid, d.Id()).Execute()

		if err != nil {
			if apiResponse == nil || apiResponse.StatusCode != 404 {
				diags := diag.FromErr(fmt.Errorf("an error occured while deleting a lan dcId %s ID %s %s", d.Get("datacenter_id").(string), d.Id(), err))
				return diags
			}
		}
	}

	for {
		log.Printf("[INFO] Waiting for LAN %s to be deleted...", d.Id())
		time.Sleep(5 * time.Second)

		lDeleted, dsErr := lanDeleted(ctx, client, d)

		if dsErr != nil {
			diags := diag.FromErr(fmt.Errorf("error while checking deletion status of LAN %s: %s", d.Id(), dsErr))
			return diags
		}

		if lDeleted {
			log.Printf("[INFO] Successfully deleted LAN: %s", d.Id())
			break
		}
	}

	d.SetId("")
	return nil
}

func lanAvailable(ctx context.Context, client *ionoscloud.APIClient, d *schema.ResourceData) (bool, error) {

	dcid := d.Get("datacenter_id").(string)
	rsp, _, err := client.LanApi.DatacentersLansFindById(ctx, dcid, d.Id()).Execute()

	log.Printf("[INFO] Current status for LAN %s: %+v", d.Id(), rsp)

	if err != nil {
		return true, fmt.Errorf("error checking LAN status: %s", err)
	}
	return *rsp.Metadata.State == "AVAILABLE", nil
}

func lanDeleted(ctx context.Context, client *ionoscloud.APIClient, d *schema.ResourceData) (bool, error) {
	dcid := d.Get("datacenter_id").(string)
	rsp, apiResponse, err := client.LanApi.DatacentersLansFindById(ctx, dcid, d.Id()).Execute()

	log.Printf("[INFO] Current deletion status for LAN %s: %+v", d.Id(), rsp)

	if err != nil {
		if apiResponse != nil && apiResponse.StatusCode == 404 {
			return true, nil
		}
		return true, err

	}
	log.Printf("[INFO] LAN %s not deleted yet deleted LAN: %+v", d.Id(), rsp)
	return false, nil
}
