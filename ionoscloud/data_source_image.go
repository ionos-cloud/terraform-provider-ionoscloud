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
			"description": {
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
			"size": {
				Type:     schema.TypeFloat,
				Computed: true,
			},
			"cpu_hot_plug": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"cpu_hot_unplug": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"ram_hot_plug": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"ram_hot_unplug": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"nic_hot_plug": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"nic_hot_unplug": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"disc_virtio_hot_plug": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"disc_virtio_hot_unplug": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"disc_scsi_hot_plug": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"disc_scsi_hot_unplug": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"license_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"public": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"image_aliases": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"cloud_init": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
		Timeouts: &resourceDefaultTimeouts,
	}
}

func dataSourceImageRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ionoscloud.APIClient)

	ctx, cancel := context.WithTimeout(context.Background(), *resourceDefaultTimeouts.Default)

	if cancel != nil {
		defer cancel()
	}

	images, _, err := client.ImageApi.ImagesGet(ctx).Execute()

	if err != nil {
		return fmt.Errorf("an error occured while fetching IonosCloud images %s", err)
	}

	name := d.Get("name").(string)
	imageType, imageTypeOk := d.GetOk("type")
	location, locationOk := d.GetOk("location")
	version, versionOk := d.GetOk("version")
	cloudInit, cloudInitOk := d.GetOk("cloud_init")

	var results []ionoscloud.Image

	// if version value is present then concatenate name - version
	// otherwise search by name or part of the name
	if versionOk {
		nameVer := fmt.Sprintf("%s-%s", name, version.(string))
		if images.Items != nil {
			for _, img := range *images.Items {
				if img.Properties.Name != nil && strings.ToLower(*img.Properties.Name) == strings.ToLower(nameVer) {
					results = append(results, img)
				}
			}
		}
		if results == nil {
			return fmt.Errorf("could not find an image with name %s and version %s (%s)", name, version.(string), nameVer)
		}
	} else {
		if images.Items != nil {
		for _, img := range *images.Items {
			if img.Properties.Name != nil && strings.ToLower(*img.Properties.Name) == strings.ToLower(name) {
					results = append(results, img)
					break
				}
			}
		}
		if results == nil {
			return fmt.Errorf("could not find an image with name %s", name)
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

	if cloudInitOk {
		cloudInitResults := []ionoscloud.Image{}
		for _, img := range results {
			if img.Properties.CloudInit != nil && *img.Properties.CloudInit == cloudInit.(string) {
				cloudInitResults = append(cloudInitResults, img)
			}
		}
		results = cloudInitResults
	}

	if len(results) == 0 {
		return fmt.Errorf("there are no images that match the search criteria")
	}

	if results[0].Properties.Name != nil {
		err := d.Set("name", *results[0].Properties.Name)
		if err != nil {
			return fmt.Errorf("error while setting name property for image %s: %s", d.Id(), err)
		}
	}

	if results[0].Properties.Description != nil {
		if err := d.Set("description", *results[0].Properties.Description); err != nil {
			return err
		}
	}

	if results[0].Properties.Size != nil {
		if err := d.Set("size", *results[0].Properties.Size); err != nil {
			return err
		}
	}

	if results[0].Properties.CpuHotPlug != nil {
		if err := d.Set("cpu_hot_plug", *results[0].Properties.CpuHotPlug); err != nil {
			return err
		}
	}

	if results[0].Properties.CpuHotUnplug != nil {
		if err := d.Set("cpu_hot_unplug", *results[0].Properties.CpuHotUnplug); err != nil {
			return err
		}
	}

	if results[0].Properties.RamHotPlug != nil {
		if err := d.Set("ram_hot_plug", *results[0].Properties.RamHotPlug); err != nil {
			return err
		}
	}

	if results[0].Properties.RamHotUnplug != nil {
		if err := d.Set("ram_hot_unplug", *results[0].Properties.RamHotUnplug); err != nil {
			return err
		}
	}

	if results[0].Properties.NicHotPlug != nil {
		if err := d.Set("nic_hot_plug", *results[0].Properties.NicHotPlug); err != nil {
			return err
		}
	}

	if results[0].Properties.NicHotUnplug != nil {
		if err := d.Set("nic_hot_unplug", *results[0].Properties.NicHotUnplug); err != nil {
			return err
		}
	}

	if results[0].Properties.DiscVirtioHotPlug != nil {
		if err := d.Set("disc_virtio_hot_plug", *results[0].Properties.DiscVirtioHotPlug); err != nil {
			return err
		}
	}

	if results[0].Properties.DiscVirtioHotUnplug != nil {
		if err := d.Set("disc_virtio_hot_unplug", *results[0].Properties.DiscVirtioHotUnplug); err != nil {
			return err
		}
	}

	if results[0].Properties.DiscScsiHotPlug != nil {
		if err := d.Set("disc_scsi_hot_plug", *results[0].Properties.DiscScsiHotPlug); err != nil {
			return err
		}
	}

	if results[0].Properties.DiscScsiHotUnplug != nil {
		if err := d.Set("disc_scsi_hot_unplug", *results[0].Properties.DiscScsiHotUnplug); err != nil {
			return err
		}
	}

	if results[0].Properties.LicenceType != nil {
		if err := d.Set("license_type", *results[0].Properties.LicenceType); err != nil {
			return err
		}
	}

	if results[0].Properties.Public != nil {
		if err := d.Set("public", *results[0].Properties.Public); err != nil {
			return err
		}
	}

	if results[0].Properties.ImageAliases != nil {
		if err := d.Set("image_aliases", *results[0].Properties.ImageAliases); err != nil {
			return err
		}
	}

	if results[0].Properties.CloudInit != nil {
		if err := d.Set("cloud_init", *results[0].Properties.CloudInit); err != nil {
			return err
		}
	}

	if results[0].Properties.ImageType != nil {
		err = d.Set("type", *results[0].Properties.ImageType)
		if err != nil {
			return err
		}
	}

	if results[0].Properties.Location != nil {
		err = d.Set("location", *results[0].Properties.Location)
		if err != nil {
			return err
		}
	}
	d.SetId(*results[0].Id)

	return nil
}
