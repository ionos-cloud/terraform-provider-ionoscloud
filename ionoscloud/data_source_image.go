package ionoscloud

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	ionoscloud "github.com/ionos-cloud/sdk-go/v6"
)

func dataSourceImage() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceImageRead,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"image_alias": {
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
			"licence_type": {
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
	client := meta.(services.SdkBundle).CloudApiClient

	images, apiResponse, err := client.ImagesApi.ImagesGet(ctx).Depth(1).Execute()
	logApiRequestTime(apiResponse)

	if err != nil {
		return diag.FromErr(fmt.Errorf("an error occurred while fetching IonosCloud images %w", err))
	}

	nameValue, nameOk := d.GetOk("name")
	imageTypeValue, imageTypeOk := d.GetOk("type")
	locationValue, locationOk := d.GetOk("location")
	versionValue, versionOk := d.GetOk("version")
	cloudInitValue, cloudInitOk := d.GetOk("cloud_init")
	imgAliasVal, imgAliasOk := d.GetOk("image_alias")

	name := nameValue.(string)
	imageType := imageTypeValue.(string)
	location := locationValue.(string)
	version := versionValue.(string)
	cloudInit := cloudInitValue.(string)
	imgAlias := imgAliasVal.(string)
	var results []ionoscloud.Image

	// if version value is present then concatenate name - version
	// otherwise search by name or part of the name
	if versionOk && nameOk && version != "" && name != "" {
		nameVer := fmt.Sprintf("%s-%s", name, version)
		if images.Items != nil {
			for _, img := range *images.Items {
				if img.Properties != nil && img.Properties.Name != nil && strings.EqualFold(*img.Properties.Name, nameVer) {
					results = append(results, img)
				}
			}
		}
		if results == nil {
			return diag.FromErr(fmt.Errorf("no image found with the specified criteria: name %s and version %s (%s)", name, version, nameVer))
		}
	} else if nameOk && name != "" {
		var exactMatches []ionoscloud.Image
		if images.Items != nil {
			for _, img := range *images.Items {
				if img.Properties != nil && img.Properties.Name != nil && strings.Contains(strings.ToLower(*img.Properties.Name), strings.ToLower(name)) {
					results = append(results, img)
					// if the image name is an exact match, store it in a separate list of exact matches
					if strings.EqualFold(*img.Properties.Name, name) {
						exactMatches = append(exactMatches, img)
					}
				}
			}
			// if exact matches have been found, only continue filtering with these
			if len(exactMatches) > 0 {
				results = exactMatches
			}
		}
		if results == nil {
			return diag.FromErr(fmt.Errorf("no image found with the specified criteria: name %s", name))
		}
	} else {
		results = *images.Items
	}

	if imageTypeOk && imageType != "" {
		var imageTypeResults []ionoscloud.Image
		for _, img := range results {
			if img.Properties != nil && img.Properties.ImageType != nil && strings.EqualFold(*img.Properties.ImageType, imageType) {
				imageTypeResults = append(imageTypeResults, img)
			}

		}
		results = imageTypeResults
	}
	if imgAliasOk && imgAlias != "" {
		var imageAliasResults []ionoscloud.Image
		for _, img := range results {
			if img.Properties != nil {
				for _, imgAliasIdx := range *img.Properties.ImageAliases {
					if strings.EqualFold(imgAliasIdx, imgAlias) {
						imageAliasResults = append(imageAliasResults, img)
					}
				}
			}
		}
		results = imageAliasResults
	}

	if locationOk && location != "" {
		var locationResults []ionoscloud.Image
		for _, img := range results {
			if img.Properties != nil && img.Properties.Location != nil && strings.EqualFold(*img.Properties.Location, location) {
				locationResults = append(locationResults, img)
			}
		}
		results = locationResults
	}

	if cloudInitOk && cloudInit != "" {
		var cloudInitResults []ionoscloud.Image
		for _, img := range results {
			if img.Properties != nil && img.Properties.CloudInit != nil && strings.EqualFold(*img.Properties.CloudInit, cloudInit) {
				cloudInitResults = append(cloudInitResults, img)
			}
		}
		results = cloudInitResults
	}

	var image ionoscloud.Image

	if results == nil || len(results) == 0 {
		return diag.FromErr(fmt.Errorf("no image found with the specified criteria: name = %s, type = %s, location = %s, version = %s, cloudInit = %s, imageAlias = %s", name, imageType, location, version, cloudInit, imgAlias))
	} else if len(results) > 1 {
		for _, result := range results {
			if result.Properties != nil {
				log.Printf("[DEBUG] found image %s in location %s", *result.Properties.Name, *result.Properties.Location)
			}
		}
		return diag.FromErr(fmt.Errorf("more than one image found, enable debug to learn more. Criteria used name = %s, type = %s, location = %s, version = %s, cloudInit = %s, imageAlias = %s", name, imageType, location, version, cloudInit, imgAlias))
	} else {
		image = results[0]
	}

	if err := ImageSetData(d, &image); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func ImageSetData(d *schema.ResourceData, image *ionoscloud.Image) error {
	if image.Id != nil {
		d.SetId(*image.Id)
	}

	if image.Properties != nil {

		if image.Properties.Name != nil {
			err := d.Set("name", *image.Properties.Name)
			if err != nil {
				return fmt.Errorf("error while setting name property for image %s: %w", d.Id(), err)
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
			if err := d.Set("licence_type", *image.Properties.LicenceType); err != nil {
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
