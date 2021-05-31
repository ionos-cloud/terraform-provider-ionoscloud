package ionoscloud

import (
	"context"
	"fmt"
<<<<<<< HEAD
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
=======
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
>>>>>>> master
	ionoscloud "github.com/ionos-cloud/sdk-go/v6"
	"log"
	"regexp"
	"time"
)

func resourceDatacenter() *schema.Resource {
	return &schema.Resource{
		Create: resourceDatacenterCreate,
		Read:   resourceDatacenterRead,
		Update: resourceDatacenterUpdate,
		Delete: resourceDatacenterDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{

			//Datacenter parameters
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.All(validation.StringIsNotWhiteSpace),
			},

			"location": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.All(validation.StringIsNotWhiteSpace),
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"sec_auth_protection": {
				Type:     schema.TypeBool,
				Optional: true,
			},
		},
		Timeouts: &resourceDefaultTimeouts,
	}
}

func resourceDatacenterCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*ionoscloud.APIClient)

	datacenterName := d.Get("name").(string)
	datacenterLocation := d.Get("location").(string)

	datacenter := ionoscloud.Datacenter{
		Properties: &ionoscloud.DatacenterProperties{
			Name:     &datacenterName,
			Location: &datacenterLocation,
		},
	}

	if attr, ok := d.GetOk("description"); ok {
		attrStr := attr.(string)
		datacenter.Properties.Description = &attrStr
	}

	if attr, ok := d.GetOk("sec_auth_protection"); ok {
		attrStr := attr.(bool)
		datacenter.Properties.SecAuthProtection = &attrStr
	}

	ctx, cancel := context.WithTimeout(context.Background(), *resourceDefaultTimeouts.Create)

	if cancel != nil {
		defer cancel()
	}

	createdDatacenter, apiResponse, err := client.DataCentersApi.DatacentersPost(ctx).Datacenter(datacenter).Execute()

	if err != nil {
		return fmt.Errorf(
			"Error creating data center (%s) (%s)", d.Id(), err)
	}
	d.SetId(*createdDatacenter.Id)

	log.Printf("[INFO] DataCenter Id: %s", d.Id())

	// Wait, catching any errors
	_, errState := getStateChangeConf(meta, d, apiResponse.Header.Get("Location"), schema.TimeoutCreate).WaitForState()

	if errState != nil {
		if IsRequestFailed(err) {
			// Request failed, so resource was not created, delete resource from state file
			d.SetId("")
		}
		return errState
	}

	return resourceDatacenterRead(d, meta)
}

func resourceDatacenterRead(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*ionoscloud.APIClient)

	ctx, cancel := context.WithTimeout(context.Background(), *resourceDefaultTimeouts.Default)

	if cancel != nil {
		defer cancel()
	}

	datacenter, apiResponse, err := client.DataCentersApi.DatacentersFindById(ctx, d.Id()).Execute()

	if err != nil {
		if apiResponse != nil && apiResponse.Response.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("error while fetching a data center ID %s %s", d.Id(), err)
	}

	if datacenter.Properties.Name != nil {
		err := d.Set("name", *datacenter.Properties.Name)
		if err != nil {
			return fmt.Errorf("error while setting name property for datacenter %s: %s", d.Id(), err)
		}
	}

	if datacenter.Properties.Location != nil {
		err := d.Set("location", *datacenter.Properties.Location)
		if err != nil {
			return fmt.Errorf("error while setting location property for datacenter %s: %s", d.Id(), err)
		}
	}

	if datacenter.Properties.Description != nil {
		err := d.Set("description", *datacenter.Properties.Description)
		if err != nil {
			return fmt.Errorf("error while setting description property for datacenter %s: %s", d.Id(), err)
		}
	}

	if datacenter.Properties.SecAuthProtection != nil {
		err := d.Set("sec_auth_protection", *datacenter.Properties.SecAuthProtection)
		if err != nil {
			return fmt.Errorf("error while setting sec_auth_protection property for datacenter %s: %s", d.Id(), err)
		}
	}

	return nil
}

func resourceDatacenterUpdate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*ionoscloud.APIClient)
	obj := ionoscloud.DatacenterProperties{}

	if d.HasChange("name") {
		_, newName := d.GetChange("name")
		newNameStr := newName.(string)
		obj.Name = &newNameStr
	}

	if d.HasChange("description") {
		_, newDescription := d.GetChange("description")
		newDescriptionStr := newDescription.(string)
		obj.Description = &newDescriptionStr
	}

	if d.HasChange("location") {
		oldLocation, newLocation := d.GetChange("location")
		return fmt.Errorf("Data center is created in %s location. You can not change location of the data center to %s. It requires recreation of the data center.", oldLocation, newLocation)
	}

	if d.HasChange("sec_auth_protection") {
		_, newSecAuthProtection := d.GetChange("sec_auth_protection")
		newSecAuthProtectionStr := newSecAuthProtection.(bool)
		obj.SecAuthProtection = &newSecAuthProtectionStr
	}

	ctx, cancel := context.WithTimeout(context.Background(), *resourceDefaultTimeouts.Update)

	if cancel != nil {
		defer cancel()
	}

	_, apiResponse, err := client.DataCentersApi.DatacentersPatch(ctx, d.Id()).Datacenter(obj).Execute()

	if err != nil {
		return fmt.Errorf("an error occured while update the data center ID %s %s", d.Id(), err)
	}

	// Wait, catching any errors
	_, errState := getStateChangeConf(meta, d, apiResponse.Header.Get("Location"), schema.TimeoutUpdate).WaitForState()
	if errState != nil {
		return errState
	}

	return resourceDatacenterRead(d, meta)
}

func resourceDatacenterDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*ionoscloud.APIClient)

	ctx, cancel := context.WithTimeout(context.Background(), *resourceDefaultTimeouts.Delete)

	if cancel != nil {
		defer cancel()
	}

	_, apiResponse, err := client.DataCentersApi.DatacentersDelete(ctx, d.Id()).Execute()

	/*todo: error Can not check a state when path is empty - getStateChangeConf without path - datacenter deletion retried after 20 seconds due to the fact that in k8s_node_pool tests datacenter deletion is not permitted from the first try*/
	if err != nil {
		if apiResponse.Response.StatusCode == 404 {
			return err
		}
		//try again in 20 seconds
		time.Sleep(20 * time.Second)

		_, apiResponse, err := client.DataCentersApi.DatacentersDelete(ctx, d.Id()).Execute()

		if err != nil {
			return fmt.Errorf("an error occured while deleting the data center ID %s %s \n ApiError: %s", d.Id(), err, string(apiResponse.Payload))
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

func IsValidUUID(uuid string) bool {
	r := regexp.MustCompile("^[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{12}$")
	return r.MatchString(uuid)
}
