package ionoscloud

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/profitbricks/profitbricks-sdk-go/v5"
)

func dataSourceImage() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceImageRead,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"type": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"location": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"version": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
		Timeouts: &resourceDefaultTimeouts,
	}
}

func dataSourceImageRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*profitbricks.Client)

	images, err := client.ListImages()

	if err != nil {
		return fmt.Errorf("An error occured while fetching IonosCloud images %s", err)
	}

	name := d.Get("name").(string)
	imageType, imageTypeOk := d.GetOk("type")
	location, locationOk := d.GetOk("location")
	version, versionOk := d.GetOk("version")

	var results []profitbricks.Image

	// if version value is present then concatenate name - version
	// otherwise search by name or part of the name
	if versionOk {
		nameVer := fmt.Sprintf("%s-%s", name, version.(string))
		for _, img := range images.Items {
			if strings.ToLower(img.Properties.Name) == strings.ToLower(nameVer) {
				results = append(results, img)
			}
		}
		if results == nil {
			return fmt.Errorf("could not find an image with name %s and version %s (%s)", name, version.(string), nameVer)
		}
	} else {
		for _, img := range images.Items {
			if strings.ToLower(img.Properties.Name) == strings.ToLower(name) {
				results = append(results, img)
				break
			}
		}
		if results == nil {
			return fmt.Errorf("could not find an image with name %s", name)
		}
	}

	if imageTypeOk {
		imageTypeResults := []profitbricks.Image{}
		for _, img := range results {
			if img.Properties.ImageType == imageType.(string) {
				imageTypeResults = append(imageTypeResults, img)
			}

		}
		results = imageTypeResults
	}

	if locationOk {
		locationResults := []profitbricks.Image{}
		for _, img := range results {
			if img.Properties.Location == location.(string) {
				locationResults = append(locationResults, img)
			}

		}
		results = locationResults
	}

	if len(results) == 0 {
		return fmt.Errorf("There are no images that match the search criteria")
	}

	err = d.Set("name", results[0].Properties.Name)
	if err != nil {
		return err
	}

	err = d.Set("type", results[0].Properties.ImageType)
	if err != nil {
		return err
	}

	err = d.Set("location", results[0].Properties.Location)
	if err != nil {
		return err
	}

	d.SetId(results[0].ID)

	return nil
}
