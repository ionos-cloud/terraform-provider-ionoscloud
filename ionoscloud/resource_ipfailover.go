package ionoscloud

import (
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/profitbricks/profitbricks-sdk-go/v5"
)

func resourceLanIPFailover() *schema.Resource {
	return &schema.Resource{
		Create: resourceLanIPFailoverCreate,
		Read:   resourceLanIPFailoverRead,
		Update: resourceLanIPFailoverUpdate,
		Delete: resourceLanIPFailoverDelete,
		Schema: map[string]*schema.Schema{
			"ip": {
				Type:     schema.TypeString,
				Required: true,
			},
			"nicuuid": {
				Type:     schema.TypeString,
				Required: true,
			},
			"lan_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"datacenter_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
		Timeouts: &resourceDefaultTimeouts,
	}
}

func resourceLanIPFailoverCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(SdkBundle).LegacyClient
	dcid := d.Get("datacenter_id").(string)
	lanid := d.Get("lan_id").(string)
	if lanid == "" {
		return fmt.Errorf("'lan_id' is missing, please provide a valid lan ID ")
	}
	ip := d.Get("ip").(string)
	nicUuid := d.Get("nicuuid").(string)
	properties := &profitbricks.LanProperties{}

	properties.IPFailover = &[]profitbricks.IPFailover{
		{
			IP:      ip,
			NicUUID: nicUuid,
		}}

	if properties != nil {
		lan, err := client.UpdateLan(dcid, lanid, *properties)
		if err != nil {
			return fmt.Errorf("An error occured while patching a lans failover group  %s %s", lanid, err)
		}

		// Wait, catching any errors
		_, errState := getStateChangeConf(meta, d, lan.Headers.Get("Location"), schema.TimeoutCreate).WaitForState()
		if errState != nil {
			return errState
		}

		d.SetId(lan.ID)
	}
	return resourceLanIPFailoverRead(d, meta)
}

func resourceLanIPFailoverRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(SdkBundle).LegacyClient
	lan, err := client.GetLan(d.Get("datacenter_id").(string), d.Id())

	if err != nil {
		if apiError, ok := err.(profitbricks.ApiError); ok {
			if apiError.HttpStatusCode() == 404 {
				d.SetId("")
				return nil
			}
		}
		return fmt.Errorf("An error occured while fetching a lan ID %s %s", d.Id(), err)
	}

	d.Set("public", lan.Properties.Public)
	d.Set("name", lan.Properties.Name)
	d.Set("ip_failover", lan.Properties.IPFailover)
	d.Set("datacenter_id", d.Get("datacenter_id").(string))
	return nil
}

func resourceLanIPFailoverUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(SdkBundle).LegacyClient
	properties := &profitbricks.LanProperties{}
	dcid := d.Get("datacenter_id").(string)
	lanid := d.Get("lan_id").(string)
	ip := d.Get("ip").(string)
	nicUuid := d.Get("nicuuid").(string)

	properties.IPFailover = &[]profitbricks.IPFailover{
		{
			IP:      ip,
			NicUUID: nicUuid,
		}}

	if properties != nil {
		lan, err := client.UpdateLan(dcid, lanid, *properties)
		if err != nil {
			return fmt.Errorf("An error occured while patching a lan ID %s %s", d.Id(), err)
		}

		// Wait, catching any errors
		_, errState := getStateChangeConf(meta, d, lan.Headers.Get("Location"), schema.TimeoutUpdate).WaitForState()
		if errState != nil {
			return errState
		}
	}
	return resourceLanIPFailoverRead(d, meta)
}

func resourceLanIPFailoverDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(SdkBundle).LegacyClient
	dcid := d.Get("datacenter_id").(string)
	lanid := d.Get("lan_id").(string)

	//remove the failover group
	properties := &profitbricks.LanProperties{
		IPFailover: &[]profitbricks.IPFailover{},
	}

	ipfailover, err := client.UpdateLan(dcid, lanid, *properties)
	if err != nil {
		//try again in 90 seconds
		time.Sleep(90 * time.Second)
		ipfailover, err = client.UpdateLan(dcid, lanid, *properties)

		if err != nil {
			if apiError, ok := err.(profitbricks.ApiError); ok {
				if apiError.HttpStatusCode() != 404 {
					return fmt.Errorf("An error occured while removing a lans ipfailover groups dcId %s ID %s %s", d.Get("datacenter_id").(string), d.Id(), err)
				}
			}
		}
	}

	// Wait, catching any errors
	_, errState := getStateChangeConf(meta, d, ipfailover.Headers.Get("Location"), schema.TimeoutDelete).WaitForState()
	if errState != nil {
		return errState
	}

	d.SetId("")
	return nil
}
