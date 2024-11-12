package ionoscloud

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/cloudapi/ftp"
)

func resourceImage() *schema.Resource {
	return &schema.Resource{
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
				Optional:         true,
				ForceNew:         true,
				ConflictsWith:    []string{"ftp_url"},
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice(ftp.ValidFTPLocations, false)),
			},
			"ftp_url": {
				Type:          schema.TypeString,
				Description:   "Custom URL for the FTP server",
				Optional:      true,
				ForceNew:      true,
				ConflictsWith: []string{"location"},
				// TODO -- Check if we need to add a validation function
			},
			"image_path": {
				Type:        schema.TypeString,
				Description: "The path to the image file to be uploaded",
				Required:    true,
				ForceNew:    true,
			},
			// TODO -- check if we need to attributes like "skip_verify"/"skip_update" as we do have
			// for CLI
			"name": {
				Type:        schema.TypeString,
				Description: "The name of the image",
				Optional:    true,
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
			},
			"ram_hot_plug": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"nic_hot_plug": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"disc_virtio_hot_plug": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"disc_scsi_hot_plug": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"cpu_hot_unplug": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"ram_hot_unplug": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"nic_hot_unplug": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"disc_virtio_hot_unplug": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"disc_scsi_hot_unplug": {
				Type:     schema.TypeBool,
				Optional: true,
			},
		},
	}
}

func resourceImageImport(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	return nil, nil
}
