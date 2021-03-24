package ionoscloud

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/profitbricks/profitbricks-sdk-go/v5"
)

func resourcePrivateCrossConnect() *schema.Resource {
	return &schema.Resource{
		Create: resourcePrivateCrossConnectCreate,
		Read:   resourcePrivateCrossConnectRead,
		Update: resourcePrivateCrossConnectUpdate,
		Delete: resourcePrivateCrossConnectDelete,
		Importer: &schema.ResourceImporter{
			State: resourcePrivateCrossConnectImport,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Description: "The desired name for the private cross-connect",
				Required:    true,
			},
			"description": {
				Type:        schema.TypeString,
				Description: "The desired description for the private cross-connect",
				Optional:    true,
			},
			"connectable_datacenters": {
				Type:        schema.TypeList,
				Description: "A list containing all the connectable datacenters",
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Description: "The UUID of the connectable datacenter",
							Computed:    true,
						},
						"name": {
							Type:        schema.TypeString,
							Description: "The name of the connectable datacenter",
							Computed:    true,
						},
						"location": {
							Type:        schema.TypeString,
							Description: "The physical location of the connectable datacenter",
							Computed:    true,
						},
					},
				},
			},
			"peers": {
				Type:        schema.TypeList,
				Description: "A list containing the details of all datacenter cross-connected through this private cross-connect",
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"lan_id": {
							Type:        schema.TypeString,
							Description: "The id of the cross-connected LAN",
							Computed:    true,
						},
						"lan_name": {
							Type:        schema.TypeString,
							Description: "The name of the cross-connected LAN",
							Computed:    true,
						},
						"datacenter_id": {
							Type:        schema.TypeString,
							Description: "The id of the cross-connected datacenter",
							Computed:    true,
						},
						"datacenter_name": {
							Type:        schema.TypeString,
							Description: "The name of the cross-connected datacenter",
							Computed:    true,
						},
						"location": {
							Type:        schema.TypeString,
							Description: "The location of the cross-connected datacenter",
							Computed:    true,
						},
					},
				},
			},
		},
		Timeouts: &resourceDefaultTimeouts,
	}
}

func resourcePrivateCrossConnectCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(SdkBundle).LegacyClient

	pcc := profitbricks.PrivateCrossConnect{
		Properties: &profitbricks.PrivateCrossConnectProperties{
			Name: d.Get("name").(string),
		},
	}

	if descVal, descOk := d.GetOk("description"); descOk {
		log.Printf("[INFO] Setting PCC description to : %s", descVal.(string))
		pcc.Properties.Description = descVal.(string)
	}

	createdPCC, err := client.CreatePrivateCrossConnect(pcc)

	if err != nil {
		d.SetId("")
		return fmt.Errorf("Error creating private PCC: %s", err)
	}

	d.SetId(createdPCC.ID)
	log.Printf("[INFO] Created PCC: %s", d.Id())

	for {
		log.Printf("[INFO] Waiting for PCC %s to be ready...", d.Id())
		time.Sleep(5 * time.Second)

		pccReady, rsErr := privateCrossConnectReady(client, d)

		if rsErr != nil {
			return fmt.Errorf("Error while checking readiness status of PCC %s: %s", d.Id(), rsErr)
		}

		if pccReady && rsErr == nil {
			log.Printf("[INFO] PCC ready: %s", d.Id())
			break
		}
	}

	return resourcePrivateCrossConnectRead(d, meta)
}

func resourcePrivateCrossConnectRead(d *schema.ResourceData, meta interface{}) error {

	client := meta.(SdkBundle).LegacyClient
	pcc, err := client.GetPrivateCrossConnect(d.Id())

	if err != nil {
		if apiError, ok := err.(profitbricks.ApiError); ok {
			if apiError.HttpStatusCode() == 404 {
				d.SetId("")
				return nil
			}
		}
		return fmt.Errorf("Error while fetching PCC %s: %s", d.Id(), err)
	}

	log.Printf("[INFO] Successfully retreived PCC %s: %+v", d.Id(), pcc)

	peers := []map[string]string{}

	for _, peer := range *pcc.Properties.Peers {
		peers = append(peers, map[string]string{
			"lan_id":          peer.LANId,
			"lan_name":        peer.LANName,
			"datacenter_id":   peer.DataCenterID,
			"datacenter_name": peer.DataCenterName,
			"location":        peer.Location,
		})
	}

	d.Set("peers", peers)
	log.Printf("[INFO] Setting peers for PCC %s to %+v...", d.Id(), d.Get("peers"))

	connectableDatacenters := []map[string]string{}

	for _, connectableDC := range *pcc.Properties.ConnectableDatacenters {
		connectableDatacenters = append(connectableDatacenters, map[string]string{
			"id":       connectableDC.ID,
			"name":     connectableDC.Name,
			"location": connectableDC.Location,
		})
	}

	d.Set("connectable_datacenters", connectableDatacenters)

	return nil
}

func resourcePrivateCrossConnectUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(SdkBundle).LegacyClient
	request := profitbricks.PrivateCrossConnect{}

	request.Properties = &profitbricks.PrivateCrossConnectProperties{
		Name: d.Get("name").(string),
	}

	if d.HasChange("name") {
		oldName, newName := d.GetChange("name")
		log.Printf("[INFO] PCC name changed from %+v to %+v", oldName, newName)
		request.Properties.Name = newName.(string)
	}

	log.Printf("[INFO] Attempting update PCC %s", d.Id())

	if d.HasChange("description") {
		oldDesc, newDesc := d.GetChange("description")
		log.Printf("[INFO] PCC description changed from %+v to %+v", oldDesc, newDesc)
		if newDesc != nil {
			request.Properties.Description = newDesc.(string)
		}
	}

	_, err := client.UpdatePrivateCrossConnect(d.Id(), request)

	if err != nil {
		if apiError, ok := err.(profitbricks.ApiError); ok {
			if apiError.HttpStatusCode() == 404 {
				d.SetId("")
				return nil
			}
			return fmt.Errorf("Error while updating PCC: %s", err)
		}
		return fmt.Errorf("Error while updating PCC %s: %s", d.Id(), err)
	}

	for {
		log.Printf("[INFO] Waiting for PCC %s to be ready...", d.Id())
		time.Sleep(5 * time.Second)

		pccReady, rsErr := privateCrossConnectReady(client, d)

		if rsErr != nil {
			return fmt.Errorf("Error while checking readiness status of PCC %s: %s", d.Id(), rsErr)
		}

		if pccReady && rsErr == nil {
			log.Printf("[INFO] PCC ready: %s", d.Id())
			break
		}
	}

	return resourcePrivateCrossConnectRead(d, meta)
}

func resourcePrivateCrossConnectDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(SdkBundle).LegacyClient

	_, err := client.DeletePrivateCrossConnect(d.Id())

	if err != nil {
		if apiError, ok := err.(profitbricks.ApiError); ok {
			if apiError.HttpStatusCode() == 404 {
				d.SetId("")
				return nil
			}
			return fmt.Errorf("Error while deleting PCC: %s", err)
		}

		return fmt.Errorf("Error while deleting PCC %s: %s", d.Id(), err)
	}

	for {
		log.Printf("[INFO] Waiting for PCC %s to be deleted...", d.Id())
		time.Sleep(5 * time.Second)

		pccDeleted, dsErr := privateCrossConnectDeleted(client, d)

		if dsErr != nil {
			return fmt.Errorf("Error while checking deletion status of PCC %s: %s", d.Id(), dsErr)
		}

		if pccDeleted && dsErr == nil {
			log.Printf("[INFO] Successfully deleted PCC: %s", d.Id())
			break
		}
	}

	return nil
}

func privateCrossConnectReady(client *profitbricks.Client, d *schema.ResourceData) (bool, error) {
	subjectPCC, err := client.GetPrivateCrossConnect(d.Id())

	if err != nil {
		return true, fmt.Errorf("Error checking PCC status: %s", err)
	}
	return subjectPCC.Metadata.State == "AVAILABLE", nil
}

func privateCrossConnectDeleted(client *profitbricks.Client, d *schema.ResourceData) (bool, error) {
	_, err := client.GetPrivateCrossConnect(d.Id())

	if err != nil {
		if apiError, ok := err.(profitbricks.ApiError); ok {
			if apiError.HttpStatusCode() == 404 {
				return true, nil
			}
			return true, fmt.Errorf("Error checking PCC deletion status: %s", err)
		}
	}
	return false, nil
}
