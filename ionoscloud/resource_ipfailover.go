package ionoscloud

import (
	"context"
	"fmt"
	ionoscloud "github.com/ionos-cloud/sdk-go/v5"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
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
	client := meta.(*ionoscloud.APIClient)
	dcid := d.Get("datacenter_id").(string)
	lanid := d.Get("lan_id").(string)
	if lanid == "" {
		return fmt.Errorf("'lan_id' is missing, please provide a valid lan ID ")
	}
	ip := d.Get("ip").(string)
	nicUuid := d.Get("nicuuid").(string)

	properties := &ionoscloud.LanProperties{}

	properties.IpFailover = &[]ionoscloud.IPFailover{
		{
			Ip:      &ip,
			NicUuid: &nicUuid,
		}}

	ctx, cancel := context.WithTimeout(context.Background(), *resourceDefaultTimeouts.Create)

	if cancel != nil {
		defer cancel()
	}

	if properties != nil {
		lan, apiResponse, err := client.LanApi.DatacentersLansPatch(ctx, dcid, lanid).Lan(*properties).Execute()
		if err != nil {
			return fmt.Errorf("An error occured while patching a lans failover group  %s %s", lanid, err)
		}

		// Wait, catching any errors
		_, errState := getStateChangeConf(meta, d, apiResponse.Header.Get("Location"), schema.TimeoutCreate).WaitForState()
		if errState != nil {
			return errState
		}

		d.SetId(*lan.Id)
	}
	return resourceLanIPFailoverRead(d, meta)
}

func resourceLanIPFailoverRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ionoscloud.APIClient)

	ctx, cancel := context.WithTimeout(context.Background(), *resourceDefaultTimeouts.Default)
	if cancel != nil {
		defer cancel()
	}

	lan, apiResponse, err := client.LanApi.DatacentersLansFindById(ctx, d.Get("datacenter_id").(string), d.Id()).Execute()

	if err != nil {
		if apiResponse != nil && apiResponse.Response.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("an error occured while fetching a lan ID %s %s", d.Id(), err)
	}

	if lan.Properties.IpFailover != nil {
		err := d.Set("ip", *(*lan.Properties.IpFailover)[0].Ip)
		if err != nil {
			return fmt.Errorf("Error while setting ip property for IpFailover %s: %s", d.Id(), err)
		}
	}

	if lan.Properties.IpFailover != nil {
		err := d.Set("nicuuid", *(*lan.Properties.IpFailover)[0].NicUuid)
		if err != nil {
			return fmt.Errorf("Error while setting nicuuid property for IpFailover %s: %s", d.Id(), err)
		}
	}

	if lan.Id != nil {
		err := d.Set("lan_id", *lan.Id)
		if err != nil {
			return fmt.Errorf("Error while setting lan_id property for IpFailover %s: %s", d.Id(), err)
		}
	}

	d.Set("datacenter_id", d.Get("datacenter_id").(string))

	return nil
}

func resourceLanIPFailoverUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ionoscloud.APIClient)

	properties := &ionoscloud.LanProperties{}
	dcid := d.Get("datacenter_id").(string)
	lanid := d.Get("lan_id").(string)
	ip := d.Get("ip").(string)
	nicUuid := d.Get("nicuuid").(string)

	properties.IpFailover = &[]ionoscloud.IPFailover{
		{
			Ip:      &ip,
			NicUuid: &nicUuid,
		}}

	if properties != nil {
		ctx, cancel := context.WithTimeout(context.Background(), *resourceDefaultTimeouts.Update)

		if cancel != nil {
			defer cancel()
		}

		_, apiResponse, err := client.LanApi.DatacentersLansPatch(ctx, dcid, lanid).Lan(*properties).Execute()
		if err != nil {
			return fmt.Errorf("An error occured while patching a lan ID %s %s", d.Id(), err)
		}

		// Wait, catching any errors
		_, errState := getStateChangeConf(meta, d, apiResponse.Header.Get("Location"), schema.TimeoutUpdate).WaitForState()
		if errState != nil {
			return errState
		}
	}
	return resourceLanIPFailoverRead(d, meta)
}

func resourceLanIPFailoverDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ionoscloud.APIClient)

	dcid := d.Get("datacenter_id").(string)
	lanid := d.Get("lan_id").(string)

	//remove the failover group
	properties := &ionoscloud.LanProperties{
		IpFailover: &[]ionoscloud.IPFailover{},
	}

	ctx, cancel := context.WithTimeout(context.Background(), *resourceDefaultTimeouts.Update)

	if cancel != nil {
		defer cancel()
	}

	_, apiResponse, err := client.LanApi.DatacentersLansPatch(ctx, dcid, lanid).Lan(*properties).Execute()
	if err != nil {
		//try again in 90 seconds
		time.Sleep(90 * time.Second)
		_, apiResponse, err = client.LanApi.DatacentersLansPatch(ctx, dcid, lanid).Lan(*properties).Execute()

		if err != nil && (apiResponse == nil || apiResponse.Response.StatusCode != 404) {
			return fmt.Errorf("an error occured while removing a lans ipfailover groups dcId %s ID %s %s", d.Get("datacenter_id").(string), d.Id(), err)
		}
	}

	// Wait, catching any errors
	_, errState := getStateChangeConf(meta, d, apiResponse.Header.Get("Location"), schema.TimeoutDelete).WaitForState()
	if errState != nil {
		return errState
	}

	d.SetId("")
	return nil
}
