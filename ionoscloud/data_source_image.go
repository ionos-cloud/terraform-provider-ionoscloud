package ionoscloud

import (
	"context"
	"fmt"
	ionoscloud "github.com/ionos-cloud/sdk-go/v5"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
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
	client := meta.(SdkBundle).Client

	ctx, cancel := context.WithTimeout(context.Background(), *resourceDefaultTimeouts.Default)

	if cancel != nil {
		defer cancel()
	}

	images, _, err := client.ImageApi.ImagesGet(ctx).Execute()

	if err != nil {
		return fmt.Errorf("An error occured while fetching IonosCloud images %s", err)
	}

	name := d.Get("name").(string)
	imageType, imageTypeOk := d.GetOk("type")
	location, locationOk := d.GetOk("location")
	version, versionOk := d.GetOk("version")

	results := []ionoscloud.Image{}

	// if version value is present then concatenate name - version
	// otherwise search by name or part of the name
	if versionOk {
		name_ver := fmt.Sprintf("%s-%s", name, version.(string))
		if images.Items != nil {
			for _, img := range *images.Items {
				if strings.Contains(strings.ToLower(*img.Properties.Name), strings.ToLower(name_ver)) {
					results = append(results, img)
				}
			}
		}
	} else {
		if images.Items != nil {
			for _, img := range *images.Items {
				if strings.Contains(strings.ToLower(*img.Properties.Name), strings.ToLower(name)) {
					results = append(results, img)
				}
			}
		}
	}

	if imageTypeOk {
		imageTypeResults := []ionoscloud.Image{}
		for _, img := range results {
			if img.Properties.ImageType != nil && *img.Properties.ImageType == imageType.(string) {
				imageTypeResults = append(imageTypeResults, img)
			}

		}
		results = imageTypeResults
	}

	if locationOk {
		locationResults := []ionoscloud.Image{}
		for _, img := range results {
			if img.Properties.Location != nil && *img.Properties.Location == location.(string) {
				locationResults = append(locationResults, img)
			}

		}
		results = locationResults
	}

	if len(results) > 1 {
		return fmt.Errorf("There is more than one image that match the search criteria")
	}

	if len(results) == 0 {
		return fmt.Errorf("There are no images that match the search criteria")
	}

	if results[0].Properties.Name != nil {
		err := d.Set("name", *results[0].Properties.Name)
		if err != nil {
			return fmt.Errorf("Error while setting name property for image %s: %s", d.Id(), err)
		}
	}

	d.SetId(*results[0].Id)

	return nil
}
