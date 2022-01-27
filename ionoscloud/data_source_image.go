package ionoscloud

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	ionoscloud "github.com/ionos-cloud/sdk-go/v6"
	"log"
)

func dataSourceImage() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceImageRead,
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
				Optional: true,
				Computed: true,
			},
		},
		Timeouts: &resourceDefaultTimeouts,
	}
}

func dataSourceImageRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(SdkBundle).CloudApiClient

	name := d.Get("name").(string)
	imageType, imageTypeOk := d.GetOk("type")
	location, locationOk := d.GetOk("location")
	version, versionOk := d.GetOk("version")
	cloudInit, cloudInitOk := d.GetOk("cloud_init")

	var results ionoscloud.Images
	var image ionoscloud.Image

	request := client.ImagesApi.ImagesGet(ctx).Depth(1)

	// if version value is present then concatenate name - version
	// otherwise search by name or part of the name
	if versionOk {
		nameVer := fmt.Sprintf("%s-%s", name, version.(string))
		request = request.Filter("name", nameVer)
	} else {
		request = request.Filter("name", name)
	}

	if imageTypeOk {
		request = request.Filter("imageType", imageType.(string))
	}

	if locationOk {
		request = request.Filter("location", location.(string))
	}

	if cloudInitOk {
		request = request.Filter("cloudInit", cloudInit.(string))
	}

	results, apiResponse, err := request.Execute()
	logApiRequestTime(apiResponse)

	if err != nil {
		return diag.FromErr(fmt.Errorf("an error occurred while fetching Ionos images: %s", err.Error()))
	}

	if results.Items != nil && len(*results.Items) > 0 {
		image = (*results.Items)[len(*results.Items)-1]
		log.Printf("[INFO] %v images found matching the search critiria. Getting the latest image from the list %v", len(*results.Items), *image.Id)
	} else {
		return diag.FromErr(fmt.Errorf("no image found with the specified criteria"))
	}

	if err := setImageData(d, &image); err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func setImageData(d *schema.ResourceData, image *ionoscloud.Image) error {

	if image.Id != nil {
		d.SetId(*image.Id)
	}

	if image.Properties != nil {
		if image.Properties.Name != nil {
			err := d.Set("name", *image.Properties.Name)
			if err != nil {
				return fmt.Errorf("error while setting name property for image %s: %s", d.Id(), err)
			}
		}

		if image.Properties.Description != nil {
			if err := d.Set("description", *image.Properties.Description); err != nil {
				return err
			}
		}

		if image.Properties.Size != nil {
			if err := d.Set("size", *image.Properties.Size); err != nil {
				return err
			}
		}

		if image.Properties.CpuHotPlug != nil {
			if err := d.Set("cpu_hot_plug", *image.Properties.CpuHotPlug); err != nil {
				return err
			}
		}

		if image.Properties.CpuHotUnplug != nil {
			if err := d.Set("cpu_hot_unplug", *image.Properties.CpuHotUnplug); err != nil {
				return err
			}
		}

		if image.Properties.RamHotPlug != nil {
			if err := d.Set("ram_hot_plug", *image.Properties.RamHotPlug); err != nil {
				return err
			}
		}

		if image.Properties.RamHotUnplug != nil {
			if err := d.Set("ram_hot_unplug", *image.Properties.RamHotUnplug); err != nil {
				return err
			}
		}

		if image.Properties.NicHotPlug != nil {
			if err := d.Set("nic_hot_plug", *image.Properties.NicHotPlug); err != nil {
				return err
			}
		}

		if image.Properties.NicHotUnplug != nil {
			if err := d.Set("nic_hot_unplug", *image.Properties.NicHotUnplug); err != nil {
				return err
			}
		}

		if image.Properties.DiscVirtioHotPlug != nil {
			if err := d.Set("disc_virtio_hot_plug", *image.Properties.DiscVirtioHotPlug); err != nil {
				return err
			}
		}

		if image.Properties.DiscVirtioHotUnplug != nil {
			if err := d.Set("disc_virtio_hot_unplug", *image.Properties.DiscVirtioHotUnplug); err != nil {
				return err
			}
		}

		if image.Properties.DiscScsiHotPlug != nil {
			if err := d.Set("disc_scsi_hot_plug", *image.Properties.DiscScsiHotPlug); err != nil {
				return err
			}
		}

		if image.Properties.DiscScsiHotUnplug != nil {
			if err := d.Set("disc_scsi_hot_unplug", *image.Properties.DiscScsiHotUnplug); err != nil {
				return err
			}
		}

		if image.Properties.LicenceType != nil {
			if err := d.Set("license_type", *image.Properties.LicenceType); err != nil {
				return err
			}
		}

		if image.Properties.Public != nil {
			if err := d.Set("public", *image.Properties.Public); err != nil {
				return err
			}
		}

		if image.Properties.ImageAliases != nil && len(*image.Properties.ImageAliases) > 0 {
			if err := d.Set("image_aliases", *image.Properties.ImageAliases); err != nil {
				return err
			}
		}

		if image.Properties.CloudInit != nil {
			if err := d.Set("cloud_init", *image.Properties.CloudInit); err != nil {
				return err
			}
		}

		if image.Properties.ImageType != nil {
			err := d.Set("type", *image.Properties.ImageType)
			if err != nil {
				return err
			}
		}

		if image.Properties.Location != nil {
			err := d.Set("location", *image.Properties.Location)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
