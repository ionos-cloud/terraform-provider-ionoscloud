package ionoscloud

import (
	"fmt"
	"log"
	"regexp"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/profitbricks/profitbricks-sdk-go/v5"
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
	client := meta.(SdkBundle).LegacyClient
	datacenter := profitbricks.Datacenter{
		Properties: profitbricks.DatacenterProperties{
			Name:     d.Get("name").(string),
			Location: d.Get("location").(string),
		},
	}

	if attr, ok := d.GetOk("description"); ok {
		datacenter.Properties.Description = attr.(string)
	}
	dc, err := client.CreateDatacenter(datacenter)

	if err != nil {
		return fmt.Errorf(
			"Error creating data center (%s) (%s)", d.Id(), err)
	}
	d.SetId(dc.ID)

	log.Printf("[INFO] DataCenter Id: %s", d.Id())

	// Wait, catching any errors
	_, errState := getStateChangeConf(meta, d, dc.Headers.Get("Location"), schema.TimeoutCreate).WaitForState()
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
	client := meta.(SdkBundle).LegacyClient
	datacenter, err := client.GetDatacenter(d.Id())

	if err != nil {
		if apiError, ok := err.(profitbricks.ApiError); ok {
			if apiError.HttpStatusCode() == 404 {
				d.SetId("")
				return nil
			}
		}
		return fmt.Errorf("Error while fetching a data center ID %s %s", d.Id(), err)
	}

	d.Set("name", datacenter.Properties.Name)
	d.Set("location", datacenter.Properties.Location)
	d.Set("description", datacenter.Properties.Description)
	return nil
}

func resourceDatacenterUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(SdkBundle).LegacyClient
	obj := profitbricks.DatacenterProperties{}

	if d.HasChange("name") {
		_, newName := d.GetChange("name")

		obj.Name = newName.(string)
	}

	if d.HasChange("description") {
		_, newDescription := d.GetChange("description")
		obj.Description = newDescription.(string)
	}

	if d.HasChange("location") {
		oldLocation, newLocation := d.GetChange("location")
		return fmt.Errorf("Data center is created in %s location. You can not change location of the data center to %s. It requires recreation of the data center.", oldLocation, newLocation)
	}

	dc, err := client.UpdateDataCenter(d.Id(), obj)

	if err != nil {
		return fmt.Errorf("An error occured while update the data center ID %s %s", d.Id(), err)
	}

	// Wait, catching any errors
	_, errState := getStateChangeConf(meta, d, dc.Headers.Get("Location"), schema.TimeoutUpdate).WaitForState()
	if errState != nil {
		return errState
	}

	return resourceDatacenterRead(d, meta)
}

func resourceDatacenterDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(SdkBundle).LegacyClient
	dcid := d.Id()
	resp, err := client.DeleteDatacenter(dcid)

	if err != nil {
		return fmt.Errorf("An error occured while deleting the data center ID %s %s", d.Id(), err)
	}

	// Wait, catching any errors
	_, errState := getStateChangeConf(meta, d, resp.Get("Location"), schema.TimeoutDelete).WaitForState()
	if errState != nil {
		return errState
	}

	d.SetId("")
	return nil
}

func getImage(client *profitbricks.Client, dcId string, imageName string, imageType string) (*profitbricks.Image, error) {
	if imageName == "" {
		return nil, fmt.Errorf("imageName not suplied")
	}
	dc, err := client.GetDatacenter(dcId)
	if err != nil {
		log.Print(fmt.Errorf("Error while fetching a data center ID %s %s", dcId, err))
		return nil, err
	}

	images, err := client.ListImages()
	if err != nil {
		log.Print(fmt.Errorf("Error while fetching the list of images %s", err))
		return nil, err
	}

	if len(images.Items) > 0 {
		for _, i := range images.Items {
			imgName := ""
			if i.Properties.Name != "" {
				imgName = i.Properties.Name
			}

			if imageType == "SSD" {
				imageType = "HDD"
			}

			if imgName != "" && strings.Contains(strings.ToLower(imgName), strings.ToLower(imageName)) && i.Properties.ImageType == imageType && i.Properties.Location == dc.Properties.Location {
				return &i, err
			}

			if imgName != "" && strings.ToLower(imageName) == strings.ToLower(i.ID) && i.Properties.ImageType == imageType && i.Properties.Location == dc.Properties.Location {
				return &i, err
			}

		}
	}
	return nil, err
}

func getSnapshotId(client *profitbricks.Client, snapshotName string) string {
	if snapshotName == "" {
		return ""
	}
	snapshots, err := client.ListSnapshots()
	if err != nil {
		log.Print(fmt.Errorf("Error while fetching the list of snapshots %s", err))
	}

	if len(snapshots.Items) > 0 {
		for _, i := range snapshots.Items {
			imgName := ""
			if i.Properties.Name != "" {
				imgName = i.Properties.Name
			}

			if imgName != "" && strings.Contains(strings.ToLower(imgName), strings.ToLower(snapshotName)) {
				return i.ID
			}
		}
	}
	return ""
}

func getImageAlias(client *profitbricks.Client, imageAlias string, location string) string {
	if imageAlias == "" {
		return ""
	}
	locations, err := client.GetLocation(location)
	if err != nil {
		log.Print(fmt.Errorf("Error while fetching the list of snapshots %s", err))
	}

	if len(locations.Properties.ImageAliases) > 0 {
		for _, i := range locations.Properties.ImageAliases {
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
