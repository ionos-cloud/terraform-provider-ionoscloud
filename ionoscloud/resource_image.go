package ionoscloud

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	ionoscloud "github.com/ionos-cloud/sdk-go/v6"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/cloudapi/ftp"
)

func resourceImage() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceImageCreate,
		ReadContext:   resourceImageRead,
		UpdateContext: resourceImageUpdate,
		DeleteContext: resourceImageDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceImageImport,
		},
		Schema: map[string]*schema.Schema{
			"location": {
				Type: schema.TypeString,
				Description: fmt.Sprintf(
					"The location of the FTP server, "+
						"available locations: '%s'", strings.Join(ftp.ValidFTPLocations, ", '"),
				),
				//Optional:         true,
				// TODO -- In the initial phase of the implementation, leave this as required for
				// simplicity. 'location' and 'name' will be required so we can properly retrieve the
				// image using the API, but in practice, 'location' has to be optional since 'ftp_url'
				// can also be defined.
				Required: true,
				ForceNew: true,
				// TODO -- Uncomment this when we also uncomment 'ftp_url'
				// ConflictsWith:    []string{"ftp_url"},
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice(ftp.ValidFTPLocations, false)),
			},
			//"ftp_url": {
			//	Type:          schema.TypeString,
			//	Description:   "Custom URL for the FTP server",
			//	Optional:      true,
			//	ForceNew:      true,
			//	ConflictsWith: []string{"location"},
			//	// TODO -- Check if we need to add a validation function
			//},
			"image_path": {
				Type:        schema.TypeString,
				Description: "The path to the image file to be uploaded",
				Required:    true,
				ForceNew:    true,
			},
			// TODO -- check if we need to attributes like "skip_verify"/"skip_update" as we do have
			// for CLI
			// TODO -- check which of these optional field can also be marked as computed.
			"name": {
				Type:        schema.TypeString,
				Description: "The name of the image",
				// TODO -- leave this as required in the initial phase so we can fetch the image properly
				// using the 'name' and 'location'. Later, these fields will be marked as 'optional'
				// and the code logic will be changed accordingly.
				// Optional: true,
				Required: true,
			},
			"description": {
				Type:        schema.TypeString,
				Description: "The description of the image",
				Optional:    true,
			},
			"licence_type": {
				Type:        schema.TypeString,
				Description: "The OS type of the image",
				Optional:    true,
			},
			"cloud_init": {
				Type:        schema.TypeString,
				Description: "Cloud init compatibility",
				Optional:    true,
			},
			"cpu_hot_plug": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"ram_hot_plug": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"nic_hot_plug": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"disc_virtio_hot_plug": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"disc_scsi_hot_plug": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"cpu_hot_unplug": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"ram_hot_unplug": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"nic_hot_unplug": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"disc_virtio_hot_unplug": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"disc_scsi_hot_unplug": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			// TODO -- add computed attributes: image_alias / image_aliases, expose_serial, etc.
		},
		Timeouts: &resourceDefaultTimeouts,
	}
}

func resourceImageCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	// TODO -- Don't forget to change this when the name will be optional.
	imageName := d.Get("name").(string)
	location := d.Get("location").(string)

	// TODO -- Add the logic to upload the image to the FTP server
	// 1. Upload the image to the FTP server

	// 2. Periodically check the cloud API to see if the image is uploaded
	image, err := ftp.PollImage(ctx, meta, imageName, location)
	if err != nil {
		return diag.FromErr(err)
	}

	return diag.FromErr(setImageData(d, *image))
}

func resourceImageRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	imageID := d.Id()
	client := meta.(services.SdkBundle).CloudApiClient
	image, apiResponse, err := client.ImagesApi.ImagesFindById(ctx, imageID).Execute()
	logApiRequestTime(apiResponse)
	if err != nil {
		if httpNotFound(apiResponse) {
			d.SetId("")
			return nil
		}
		return diag.FromErr(fmt.Errorf("error occured while fetching image with ID: %v, error: %w", imageID, err))
	}
	return diag.FromErr(setImageData(d, image))
}

func resourceImageUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}

func resourceImageDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}

func resourceImageImport(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	return nil, nil
}

func setImageData(d *schema.ResourceData, image ionoscloud.Image) error {
	if image.Id == nil {
		return fmt.Errorf("expected a valid value for the image ID but received 'nil' instead")
	}
	d.SetId(*image.Id)
	if image.Properties != nil {
		// TODO -- Add the logic to set the properties of the image for all the fields.
		if image.Properties.CpuHotPlug != nil {
			if err := d.Set("cpu_hot_plug", *image.Properties.CpuHotPlug); err != nil {
				return err
			}
		}
	}
	return nil
}
