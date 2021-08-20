package ionoscloud

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	autoscaling "github.com/ionos-cloud/sdk-go-autoscaling"
	"strings"
)

func dataSourceAutoscalingTemplate() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceAutoscalingTemplateRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Required: false,
				Optional: true,
			},
			"name": {
				Type:        schema.TypeString,
				Description: "Name used for VMs.",
				Optional:    true,
			},
			"availability_zone": {
				Type:        schema.TypeString,
				Description: "Zone where the VMs created using this Template.",
				Computed:    true,
			},
			"cores": {
				Type:        schema.TypeInt,
				Description: "The total number of cores for the VMs.",
				Computed:    true,
			},
			"cpu_family": {
				Type:        schema.TypeString,
				Description: "CPU family for the VMs created using this Template. If null, the VM will be created with the default CPU family from the assigned location.",
				Computed:    true,
			},
			"location": {
				Type:        schema.TypeString,
				Description: "Location of the Template.",
				Computed:    true,
			},
			"nics": {
				Type:        schema.TypeList,
				Description: "List of NICs associated with this Template.",
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"lan": {
							Type:        schema.TypeInt,
							Description: "Lan Id for this template Nic.",
							Computed:    true,
						},
						"name": {
							Type:        schema.TypeString,
							Description: "Name for this template Nic.",
							Computed:    true,
						},
					},
				},
			},
			"ram": {
				Type:        schema.TypeInt,
				Description: "The amount of memory for the VMs in MB, e.g. 2048. Size must be specified in multiples of 256 MB with a minimum of 256 MB; however, if you set ramHotPlug to TRUE then you must use a minimum of 1024 MB. If you set the RAM size more than 240GB, then ramHotPlug will be set to FALSE and can not be set to TRUE unless RAM size not set to less than 240GB.",
				Computed:    true,
			},
			"volumes": {
				Type:        schema.TypeList,
				Description: "List of volumes associated with this Template. Only a single volume is currently supported.",
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"image": {
							Type:        schema.TypeString,
							Description: "Image installed on the volume. Only UUID of the image is supported currently.",
							Computed:    true,
						},
						"image_password": {
							Type:        schema.TypeString,
							Description: "Image password for this template volume.",
							Computed:    true,
						},
						"name": {
							Type:        schema.TypeString,
							Description: "Name of the template volume.",
							Computed:    true,
						},
						"size": {
							Type:        schema.TypeInt,
							Description: "User-defined size for this template volume in GB.",
							Computed:    true,
						},
						"ssh_keys": {
							Type:        schema.TypeList,
							Description: "Ssh keys that has access to the volume.",
							Computed:    true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"type": {
							Type:        schema.TypeString,
							Description: "Storage Type for this template volume (SSD or HDD).",
							Computed:    true,
						},
						"user_data": {
							Type:        schema.TypeString,
							Description: "user-data (Cloud Init) for this template volume.",
							Computed:    true,
						},
					},
				},
			},
		},
		Timeouts: &resourceDefaultTimeouts,
	}
}

func dataSourceAutoscalingTemplateRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(SdkBundle).AutoscalingClient

	id, idOk := d.GetOk("id")
	name, nameOk := d.GetOk("name")

	if idOk && nameOk {
		diags := diag.FromErr(fmt.Errorf("id and name cannot be both specified in the same time"))
		return diags
	}
	if !idOk && !nameOk {
		diags := diag.FromErr(fmt.Errorf("please provide either the template id or name"))
		return diags
	}

	var template autoscaling.Template
	var err error

	if idOk {
		/* search by ID */
		template, _, err = client.TemplatesApi.AutoscalingTemplatesFindById(ctx, id.(string)).Execute()
		if err != nil {
			diags := diag.FromErr(fmt.Errorf("an error occurred while fetching template with ID %s: %s", id.(string), err))
			return diags
		}
	} else {
		/* search by name */
		var templates autoscaling.TemplateCollection

		templates, _, err := client.TemplatesApi.AutoscalingTemplatesGet(ctx).Execute()
		if err != nil {
			diags := diag.FromErr(fmt.Errorf("an error occurred while fetching template: %s", err.Error()))
			return diags
		}

		found := false
		if templates.Items != nil {
			for _, t := range *templates.Items {
				tmpTemplate, _, err := client.TemplatesApi.AutoscalingTemplatesFindById(ctx, *t.Id).Execute()
				if err != nil {
					diags := diag.FromErr(fmt.Errorf("an error occurred while fetching group %s: %s", *t.Id, err))
					return diags
				}
				if tmpTemplate.Properties.Name != nil {
					if strings.Contains(*tmpTemplate.Properties.Name, name.(string)) {
						template = tmpTemplate
						found = true
						break
					}
				}
			}
		}

		if !found {
			diags := diag.FromErr(fmt.Errorf("template not found"))
			return diags
		}
	}

	if diags := setAutoscalingTemplateData(d, &template); diags != nil {
		return diags
	}

	return nil
}

func setAutoscalingTemplateData(d *schema.ResourceData, template *autoscaling.Template) diag.Diagnostics {
	d.SetId(*template.Id)
	if err := d.Set("id", *template.Id); err != nil {
		diags := diag.FromErr(err)
		return diags
	}

	if template.Properties != nil {
		if template.Properties.AvailabilityZone != nil {
			if err := d.Set("availability_zone", *template.Properties.AvailabilityZone); err != nil {
				diags := diag.FromErr(fmt.Errorf("error while setting availability_zone property for autoscaling group %s: %s", d.Id(), err))
				return diags
			}
		}

		if template.Properties.Cores != nil {
			if err := d.Set("cores", *template.Properties.Cores); err != nil {
				diags := diag.FromErr(fmt.Errorf("error while setting cores property for autoscaling template %s: %s", d.Id(), err))
				return diags
			}
		}

		if template.Properties.CpuFamily != nil {
			if err := d.Set("cpu_family", *template.Properties.CpuFamily); err != nil {
				diags := diag.FromErr(fmt.Errorf("error while setting cpu_family property for autoscaling template %s: %s", d.Id(), err))
				return diags
			}
		}

		if template.Properties.Location != nil {
			if err := d.Set("location", *template.Properties.Location); err != nil {
				diags := diag.FromErr(fmt.Errorf("error while setting location property for autoscaling template %s: %s", d.Id(), err))
				return diags
			}
		}

		if template.Properties.Name != nil {
			if err := d.Set("name", *template.Properties.Name); err != nil {
				diags := diag.FromErr(fmt.Errorf("error while setting name property for autoscaling template %s: %s", d.Id(), err))
				return diags
			}
		}

		if template.Properties.Nics != nil && len(*template.Properties.Nics) > 0 {
			var nics []interface{}

			for _, nic := range *template.Properties.Nics {
				nicEntry := make(map[string]interface{})
				if nic.Lan != nil {
					nicEntry["lan"] = *nic.Lan
				}
				if nic.Name != nil {
					nicEntry["name"] = *nic.Name
				}
				nics = append(nics, nicEntry)
			}

			if len(nics) > 0 {
				if err := d.Set("nics", nics); err != nil {
					diags := diag.FromErr(fmt.Errorf("error while setting nics property for autoscaling template %s: %s", d.Id(), err))
					return diags
				}
			}
		}

		if template.Properties.Ram != nil {
			if err := d.Set("ram", *template.Properties.Ram); err != nil {
				diags := diag.FromErr(fmt.Errorf("error while setting ram property for autoscaling template %s: %s", d.Id(), err))
				return diags
			}
		}

		if template.Properties.Volumes != nil && len(*template.Properties.Volumes) > 0 {
			var volumes []interface{}
			for _, volume := range *template.Properties.Volumes {
				volumeEntry := make(map[string]interface{})
				if volume.Image != nil {
					volumeEntry["image"] = *volume.Image
				}
				if volume.ImagePassword != nil {
					volumeEntry["image_password"] = *volume.ImagePassword
				}
				if volume.Name != nil {
					volumeEntry["name"] = *volume.Name
				}
				if volume.Size != nil {
					volumeEntry["size"] = *volume.Size
				}
				if volume.SshKeys != nil {
					volumeEntry["ssh_keys"] = *volume.SshKeys
				}
				if volume.Type != nil {
					volumeEntry["type"] = *volume.Type
				}
				if volume.UserData != nil {
					volumeEntry["user_data"] = *volume.UserData
				}
				volumes = append(volumes, volumeEntry)
			}
			if len(volumes) > 0 {
				if err := d.Set("volumes", volumes); err != nil {
					diags := diag.FromErr(fmt.Errorf("error while setting volumes property for autoscaling template %s: %s", d.Id(), err))
					return diags
				}
			}
		}
	}

	return nil

}
