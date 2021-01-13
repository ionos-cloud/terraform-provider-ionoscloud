package ionoscloud

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/profitbricks/profitbricks-sdk-go/v5"
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
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"pcc": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
		Timeouts: &resourceDefaultTimeouts,
	}
}

func resourceLanCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*profitbricks.Client)
	request := profitbricks.Lan{
		Properties: profitbricks.LanProperties{
			Public: d.Get("public").(bool),
		},
	}

	log.Printf("[DEBUG] NAME %s", d.Get("name"))
	if d.Get("name") != nil {
		request.Properties.Name = d.Get("name").(string)
	}

	if d.Get("pcc") != nil && d.Get("pcc").(string) != "" {
		pccID := d.Get("pcc").(string)
		log.Printf("[INFO] Setting PCC for LAN %s to %s...", d.Id(), pccID)
		request.Properties.PCC = pccID
	}

	lan, err := client.CreateLan(d.Get("datacenter_id").(string), request)

	if err != nil {
		d.SetId("")
		return fmt.Errorf("An error occured while creating LAN: %s", err)
	}

	log.Printf("[DEBUG] LAN ID: %s", lan.ID)
	log.Printf("[DEBUG] LAN RESPONSE: %s", lan.Response)

	d.SetId(lan.ID)

	log.Printf("[INFO] LAN ID: %s", d.Id())

	// Wait, catching any errors
	_, errState := getStateChangeConf(meta, d, lan.Headers.Get("Location"), schema.TimeoutCreate).WaitForState()
	if errState != nil {
		if IsRequestFailed(err) {
			// Request failed, so resource was not created, delete resource from state file
			d.SetId("")
		}
		return errState
	}

	for {
		log.Printf("[INFO] Waiting for LAN %s to be available...", lan.ID)
		time.Sleep(5 * time.Second)

		clusterReady, rsErr := lanAvailable(client, d)

		if rsErr != nil {
			return fmt.Errorf("Error while checking readiness status of LAN %s: %s", lan.ID, rsErr)
		}

		if clusterReady && rsErr == nil {
			log.Printf("[INFO] LAN ready: %s", d.Id())
			break
		}
	}

	return resourceLanRead(d, meta)
}

func resourceLanRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*profitbricks.Client)
	lan, err := client.GetLan(d.Get("datacenter_id").(string), d.Id())

	if err != nil {
		if apiError, ok := err.(profitbricks.ApiError); ok {
			if apiError.HttpStatusCode() == 404 {
				log.Printf("[INFO] LAN %s not found", d.Id())
				d.SetId("")
				return nil
			}
		}

		return fmt.Errorf("An error occured while fetching a LAN %s: %s", d.Id(), err)
	}

	d.Set("public", lan.Properties.Public)
	d.Set("name", lan.Properties.Name)
	d.Set("ip_failover", lan.Properties.IPFailover)
	d.Set("datacenter_id", d.Get("datacenter_id").(string))
	d.Set("pcc", lan.Properties.PCC)
	log.Printf("[INFO] LAN %s found: %+v", d.Id(), lan)
	return nil
}

func resourceLanUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*profitbricks.Client)
	properties := &profitbricks.LanProperties{}
	newValue := d.Get("public")
	properties.Public = newValue.(bool)

	if d.HasChange("name") {
		_, newValue := d.GetChange("name")
		properties.Name = newValue.(string)
	}

	if d.HasChange("pcc") {
		_, newPCC := d.GetChange("pcc")

		if newPCC != nil && newPCC.(string) != "" {
			log.Printf("[INFO] Setting PCC for LAN %s to %s...", d.Id(), newPCC.(string))
			properties.PCC = newPCC.(string)
		}
	}

	if properties != nil {
		updatedLAN, err := client.UpdateLan(d.Get("datacenter_id").(string), d.Id(), *properties)
		if err != nil {
			return fmt.Errorf("An error occured while patching a lan ID %s %s", d.Id(), err)
		}

		for {
			log.Printf("[INFO] Waiting for LAN %s to be available...", d.Id())
			time.Sleep(5 * time.Second)

			clusterReady, rsErr := lanAvailable(client, d)

			if rsErr != nil {
				return fmt.Errorf("Error while checking readiness status of LAN %s: %s", d.Id(), rsErr)
			}

			if clusterReady && rsErr == nil {
				log.Printf("[INFO] LAN %s ready: %+v", d.Id(), updatedLAN)
				break
			}
		}

	}

	return resourceLanRead(d, meta)
}

func resourceLanDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*profitbricks.Client)
	dcID := d.Get("datacenter_id").(string)

	_, err := client.DeleteLan(dcID, d.Id())

	if err != nil {
		//try again in 120 seconds
		time.Sleep(120 * time.Second)
		_, err = client.DeleteLan(dcID, d.Id())

		if err != nil {
			if apiError, ok := err.(profitbricks.ApiError); ok {
				if apiError.HttpStatusCode() != 404 {
					return fmt.Errorf("An error occured while deleting a lan dcId %s ID %s %s", d.Get("datacenter_id").(string), d.Id(), err)
				}
			}
		}
	}

	for {
		log.Printf("[INFO] Waiting for LAN %s to be deleted...", d.Id())
		time.Sleep(5 * time.Second)

		lDeleted, dsErr := lanDeleted(client, d)

		if dsErr != nil {
			return fmt.Errorf("Error while checking deletion status of LAN %s: %s", d.Id(), dsErr)
		}

		if lDeleted && dsErr == nil {
			log.Printf("[INFO] Successfully deleted LAN: %s", d.Id())
			break
		}
	}

	d.SetId("")
	return nil
}

func lanAvailable(client *profitbricks.Client, d *schema.ResourceData) (bool, error) {
	subjectLAN, err := client.GetLan(d.Get("datacenter_id").(string), d.Id())

	log.Printf("[INFO] Current status for LAN %s: %+v", d.Id(), subjectLAN)

	if err != nil {
		return true, fmt.Errorf("Error checking LAN status: %s", err)
	}
	return subjectLAN.Metadata.State == "AVAILABLE", nil
}

func lanDeleted(client *profitbricks.Client, d *schema.ResourceData) (bool, error) {
	subjectLAN, err := client.GetLan(d.Get("datacenter_id").(string), d.Id())

	log.Printf("[INFO] Current deletion status for LAN %s: %+v", d.Id(), subjectLAN)

	if err != nil {
		if apiError, ok := err.(profitbricks.ApiError); ok {
			if apiError.HttpStatusCode() == 404 {
				return true, nil
			}
			return true, fmt.Errorf("Error checking LAN deletion status: %s", err)
		}
	}
	log.Printf("[INFO] LAN %s not deleted yet deleted LAN: %+v", d.Id(), subjectLAN)
	return false, nil
}
