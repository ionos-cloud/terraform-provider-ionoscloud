package ionoscloud

import (
	"context"
	"fmt"
	ionoscloud "github.com/ionos-cloud/sdk-go/v5"
	"log"
	"regexp"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
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
				Type:     schema.TypeString,
				Required: true,
			},

			"location": {
				Type:     schema.TypeString,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		},
		Timeouts: &resourceDefaultTimeouts,
	}
}

func resourceDatacenterCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(SdkBundle).Client

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

	ctx, cancel := context.WithTimeout(context.Background(), *resourceDefaultTimeouts.Create)

	if cancel != nil {
		defer cancel()
	}

	createdDatacenter, apiResponse, err := client.DataCenterApi.DatacentersPost(ctx).Datacenter(datacenter).Execute()

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

	client := meta.(SdkBundle).Client

	ctx, cancel := context.WithTimeout(context.Background(), *resourceDefaultTimeouts.Default)

	if cancel != nil {
		defer cancel()
	}

	datacenter, apiResponse, err := client.DataCenterApi.DatacentersFindById(ctx, d.Id()).Execute()

	if err != nil {
		if _, ok := err.(ionoscloud.GenericOpenAPIError); ok {
			if apiResponse.Response.StatusCode == 404 {
				d.SetId("")
				return nil
			}
		}
		return fmt.Errorf("Error while fetching a data center ID %s %s", d.Id(), err)
	}

	if datacenter.Properties.Name != nil{
		err := d.Set("name", *datacenter.Properties.Name)
		if err != nil {
			return fmt.Errorf("Error while setting name property for backup unit %s: %s", d.Id(), err)
		}
	}

	if datacenter.Properties.Location != nil{
		err := d.Set("location", *datacenter.Properties.Location)
		if err != nil {
			return fmt.Errorf("Error while setting location property for backup unit %s: %s", d.Id(), err)
		}
	}

	if datacenter.Properties.Description != nil{
		err := d.Set("description", *datacenter.Properties.Description)
		if err != nil {
			return fmt.Errorf("Error while setting description property for backup unit %s: %s", d.Id(), err)
		}
	}

	return nil
}

func resourceDatacenterUpdate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(SdkBundle).Client
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

	ctx, cancel := context.WithTimeout(context.Background(), *resourceDefaultTimeouts.Update)

	if cancel != nil {
		defer cancel()
	}

	_, apiResponse, err := client.DataCenterApi.DatacentersPatch(ctx, d.Id()).Datacenter(obj).Execute()

	if err != nil {
		return fmt.Errorf("An error occured while update the data center ID %s %s", d.Id(), err)
	}

	// Wait, catching any errors
	_, errState := getStateChangeConf(meta, d, apiResponse.Header.Get("Location"), schema.TimeoutUpdate).WaitForState()
	if errState != nil {
		return errState
	}

	return resourceDatacenterRead(d, meta)
}

func resourceDatacenterDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(SdkBundle).Client

	ctx, cancel := context.WithTimeout(context.Background(), *resourceDefaultTimeouts.Delete)

	if cancel != nil {
		defer cancel()
	}

	_, apiResponse, err := client.DataCenterApi.DatacentersDelete(ctx, d.Id()).Execute()

	if err != nil {
		return fmt.Errorf("An error occured while deleting the data center ID %s %s", d.Id(), err)
	}

	// Wait, catching any errors
	_, errState := getStateChangeConf(meta, d, apiResponse.Header.Get("Location"), schema.TimeoutDelete).WaitForState()
	if errState != nil {
		return errState
	}

	d.SetId("")
	return nil
}

func getImage(client *ionoscloud.APIClient, dcId string, imageName string, imageType string) (*ionoscloud.Image, error) {

	ctx, cancel := context.WithTimeout(context.Background(), *resourceDefaultTimeouts.Default)

	if cancel != nil {
		defer cancel()
	}

	if imageName == "" {
		return nil, fmt.Errorf("imageName not suplied")
	}

	dc, _, err := client.DataCenterApi.DatacentersFindById(ctx, dcId).Execute()

	if err != nil {
		log.Print(fmt.Errorf("Error while fetching a data center ID %s %s", dcId, err))
		return nil, err
	}

	images, _, err := client.ImageApi.ImagesGet(ctx).Execute()

	if err != nil {
		log.Print(fmt.Errorf("Error while fetching the list of images %s", err))
		return nil, err
	}

	if len(*images.Items) > 0 {
		for _, i := range *images.Items {
			imgName := ""
			if i.Properties.Name != nil && *i.Properties.Name != "" {
				imgName = *i.Properties.Name
			}

			if imageType == "SSD" {
				imageType = "HDD"
			}

			if imgName != "" && strings.Contains(strings.ToLower(imgName), strings.ToLower(imageName)) && *i.Properties.ImageType == imageType && *i.Properties.Location == *dc.Properties.Location {
				return &i, err
			}

			if imgName != "" && strings.ToLower(imageName) == strings.ToLower(*i.Id) && *i.Properties.ImageType == imageType && *i.Properties.Location == *dc.Properties.Location {
				return &i, err
			}

		}
	}
	return nil, err
}

func getSnapshotId(client *ionoscloud.APIClient, snapshotName string) string {

	ctx, cancel := context.WithTimeout(context.Background(), *resourceDefaultTimeouts.Default)

	if cancel != nil {
		defer cancel()
	}

	if snapshotName == "" {
		return ""
	}

	snapshots, _, err := client.SnapshotApi.SnapshotsGet(ctx).Execute()

	if err != nil {
		log.Print(fmt.Errorf("Error while fetching the list of snapshots %s", err))
	}

	if len(*snapshots.Items) > 0 {
		for _, i := range *snapshots.Items {
			imgName := ""
			if *i.Properties.Name != "" {
				imgName = *i.Properties.Name
			}

			if imgName != "" && strings.Contains(strings.ToLower(imgName), strings.ToLower(snapshotName)) {
				return *i.Id
			}
		}
	}
	return ""
}


func getImageAlias(client *ionoscloud.APIClient, imageAlias string, location string) string {

	ctx, cancel := context.WithTimeout(context.Background(), *resourceDefaultTimeouts.Default)

	if cancel != nil {
		defer cancel()
	}

	if imageAlias == "" {
		return ""
	}
	parts := strings.SplitN(location, "/", 2)
	if len(parts) != 2 {
		log.Print(fmt.Errorf("Invalid location id %s", location))
	}

	locations, _, err := client.LocationApi.LocationsFindByRegionIdAndId(ctx, parts[0], parts[1]).Execute()

	if err != nil {
		log.Print(fmt.Errorf("Error while fetching the list of snapshots %s", err))
	}

	if len(*locations.Properties.ImageAliases) > 0 {
		for _, i := range *locations.Properties.ImageAliases {
			alias := ""
			if i != "" {
				alias = i
			}

			if alias != "" && strings.ToLower(alias) == strings.ToLower(imageAlias) {
				return i
			}
		}
	}
	return ""
}


func IsValidUUID(uuid string) bool {
	r := regexp.MustCompile("^[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{12}$")
	return r.MatchString(uuid)
}
