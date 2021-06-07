package ionoscloud

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	ionoscloud "github.com/ionos-cloud/sdk-go/v6"
)

func resourceLan() *schema.Resource {
	return &schema.Resource{
		Create: resourceLanCreate,
		Read:   resourceLanRead,
		Update: resourceLanUpdate,
		Delete: resourceLanDelete,
		Importer: &schema.ResourceImporter{
			State: resourceResourceImport,
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

func resourceLanCreate(d *schema.ResourceData, meta interface{}) error {
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

	ctx, cancel := context.WithTimeout(context.Background(), *resourceDefaultTimeouts.Create)
	if cancel != nil {
		defer cancel()
	}
	dcid := d.Get("datacenter_id").(string)
	rsp, apiResponse, err := client.LansApi.DatacentersLansPost(ctx, dcid).Lan(request).Execute()

	if err != nil {
		d.SetId("")
		return fmt.Errorf("an error occured while creating LAN: %s", err)
	}

	log.Printf("[DEBUG] LAN ID: %s", *rsp.Id)
	log.Printf("[DEBUG] LAN RESPONSE: %s", string(apiResponse.Payload))

	d.SetId(*rsp.Id)

	log.Printf("[INFO] LAN ID: %s", d.Id())

	// Wait, catching any errors
	_, errState := getStateChangeConf(meta, d, apiResponse.Header.Get("Location"), schema.TimeoutCreate).WaitForState()
	if errState != nil {
		if IsRequestFailed(err) {
			// Request failed, so resource was not created, delete resource from state file
			d.SetId("")
		}
		return errState
	}

	for {
		log.Printf("[INFO] Waiting for LAN %s to be available...", *rsp.Id)

		clusterReady, rsErr := lanAvailable(client, d)

		if rsErr != nil {
			return fmt.Errorf("error while checking readiness status of LAN %s: %s", *rsp.Id, rsErr)
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
			return fmt.Errorf("lan creation timed out! WARNING: your lan will still probably be created after some time but the terraform state won't reflect that; check your Ionos Cloud account for updates")
		}
	}

	return resourceLanRead(d, meta)
}

func resourceLanRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ionoscloud.APIClient)
	ctx, cancel := context.WithTimeout(context.Background(), *resourceDefaultTimeouts.Default)
	if cancel != nil {
		defer cancel()
	}
	dcid := d.Get("datacenter_id").(string)
	rsp, apiResponse, err := client.LansApi.DatacentersLansFindById(ctx, dcid, d.Id()).Execute()

	if err != nil {
		if _, ok := err.(ionoscloud.GenericOpenAPIError); ok {
			if apiResponse != nil && apiResponse.Response.StatusCode == 404 {
				log.Printf("[INFO] LAN %s not found", d.Id())
				d.SetId("")
				return nil
			}
		}

		return fmt.Errorf("an error occured while fetching a LAN %s: %s", d.Id(), err)
	}

	d.Set("public", *rsp.Properties.Public)
	d.Set("name", *rsp.Properties.Name)
	d.Set("ip_failover", *rsp.Properties.IpFailover)
	d.Set("datacenter_id", d.Get("datacenter_id").(string))
	if rsp.Properties.Pcc != nil {
		d.Set("pcc", *rsp.Properties.Pcc)
	}
	log.Printf("[INFO] LAN %s found: %+v", d.Id(), rsp)
	return nil
}

func resourceLanUpdate(d *schema.ResourceData, meta interface{}) error {
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

	ctx, cancel := context.WithTimeout(context.Background(), *resourceDefaultTimeouts.Update)
	if cancel != nil {
		defer cancel()
	}
	dcid := d.Get("datacenter_id").(string)
	rsp, _, err := client.LansApi.DatacentersLansPatch(ctx, dcid, d.Id()).Lan(*properties).Execute()
	if err != nil {
		return fmt.Errorf("an error occured while patching a lan ID %s %s", d.Id(), err)
	}

	for {
		log.Printf("[INFO] Waiting for LAN %s to be available...", d.Id())

		clusterReady, rsErr := lanAvailable(client, d)

		if rsErr != nil {
			return fmt.Errorf("error while checking readiness status of LAN %s: %s", d.Id(), rsErr)
		}

		if clusterReady {
			log.Printf("[INFO] LAN %s ready: %+v", d.Id(), rsp)
			break
		}

		select {
		case <-time.After(SleepInterval):
			log.Printf("[INFO] trying again ...")
		case <-ctx.Done():
			log.Printf("[INFO] update timed out")
			return fmt.Errorf("lan update timed out! WARNING: your lan will still probably be created after some time but the terraform state won't reflect that; check your Ionos Cloud account for updates")
		}
	}

	return resourceLanRead(d, meta)
}

func resourceLanDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ionoscloud.APIClient)
	dcid := d.Get("datacenter_id").(string)

	ctx, cancel := context.WithTimeout(context.Background(), *resourceDefaultTimeouts.Delete)
	if cancel != nil {
		defer cancel()
	}
	apiResponse, err := client.LansApi.DatacentersLansDelete(ctx, dcid, d.Id()).Execute()

	if err != nil {
		var payload string = "<nil>"
		if apiResponse != nil {
			payload = string(apiResponse.Payload)
		}
		return fmt.Errorf("an error occured while deleting lan dcId %s ID %s %s: %s", dcid, d.Id(), err, payload)
	}

	for {
		log.Printf("[INFO] Waiting for LAN %s to be deleted...", d.Id())

		lDeleted, dsErr := lanDeleted(client, d)

		if dsErr != nil {
			return fmt.Errorf("error while checking deletion status of LAN %s: %s", d.Id(), dsErr)
		}

		if lDeleted {
			log.Printf("[INFO] Successfully deleted LAN: %s", d.Id())
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

	d.SetId("")
	return nil
}

func lanAvailable(client *ionoscloud.APIClient, d *schema.ResourceData) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), *resourceDefaultTimeouts.Default)
	if cancel != nil {
		defer cancel()
	}
	dcid := d.Get("datacenter_id").(string)
	rsp, _, err := client.LansApi.DatacentersLansFindById(ctx, dcid, d.Id()).Execute()

	log.Printf("[INFO] Current status for LAN %s: %+v", d.Id(), rsp)

	if err != nil {
		return true, fmt.Errorf("error checking LAN status: %s", err)
	}
	return *rsp.Metadata.State == "AVAILABLE", nil
}

func lanDeleted(client *ionoscloud.APIClient, d *schema.ResourceData) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), *resourceDefaultTimeouts.Default)
	if cancel != nil {
		defer cancel()
	}
	dcid := d.Get("datacenter_id").(string)
	rsp, apiResponse, err := client.LansApi.DatacentersLansFindById(ctx, dcid, d.Id()).Execute()

	log.Printf("[INFO] Current deletion status for LAN %s: %+v", d.Id(), rsp)

	if err != nil {
		if _, ok := err.(ionoscloud.GenericOpenAPIError); ok {
			if apiResponse != nil && apiResponse.Response.StatusCode == 404 {
				return true, nil
			}
			return true, fmt.Errorf("error checking LAN deletion status: %s", err)
		}
	}
	log.Printf("[INFO] LAN %s not deleted yet deleted LAN: %+v", d.Id(), rsp)
	return false, nil
}
