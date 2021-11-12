package ionoscloud

//
//import (
//	"context"
//	"fmt"
//	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
//	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
//	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
//	autoscaling "github.com/ionos-cloud/sdk-go-autoscaling"
//	"log"
//)
//
//func resourceAutoscalingTemplate() *schema.Resource {
//	return &schema.Resource{
//		CreateContext: resourceAutoscalingTemplateCreate,
//		ReadContext:   resourceAutoscalingTemplateRead,
//		UpdateContext: resourceAutoscalingTemplateUpdate,
//		DeleteContext: resourceAutoscalingTemplateDelete,
//		Schema: map[string]*schema.Schema{
//			"availability_zone": {
//				Type:         schema.TypeString,
//				Description:  "Zone where the VMs created using this Template.",
//				Required:     true,
//				ValidateFunc: validation.All(validation.StringIsNotWhiteSpace),
//			},
//			"cores": {
//				Type:        schema.TypeInt,
//				Description: "The total number of cores for the VMs.",
//				Required:    true,
//				ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
//					v := val.(int)
//					if v < 1 {
//						errs = append(errs, fmt.Errorf("%q must be at least 1, got: %d", key, v))
//					}
//					return
//				},
//			},
//			"cpu_family": {
//				Type:         schema.TypeString,
//				Description:  "CPU family for the VMs created using this Template. If null, the VM will be created with the default CPU family from the assigned location.",
//				Optional:     true,
//				ValidateFunc: validation.All(validation.StringIsNotWhiteSpace),
//			},
//			"location": {
//				Type:         schema.TypeString,
//				Description:  "Location of the Template.",
//				Required:     true,
//				ValidateFunc: validation.All(validation.StringIsNotWhiteSpace),
//			},
//			"name": {
//				Type:         schema.TypeString,
//				Description:  "Name used for VMs.",
//				Optional:     true,
//				ValidateFunc: isValidName,
//			},
//			"nics": {
//				Type:        schema.TypeList,
//				Description: "List of NICs associated with this Template.",
//				Optional:    true,
//				Elem: &schema.Resource{
//					Schema: map[string]*schema.Schema{
//						"lan": {
//							Type:        schema.TypeInt,
//							Description: "Lan Id for this template Nic.",
//							Required:    true,
//						},
//						"name": {
//							Type:         schema.TypeString,
//							Description:  "Name for this template Nic.",
//							Required:     true,
//							ValidateFunc: isValidName,
//						},
//					},
//				},
//			},
//			"ram": {
//				Type:        schema.TypeInt,
//				Description: "The amount of memory for the VMs in MB, e.g. 2048. Size must be specified in multiples of 256 MB with a minimum of 256 MB; however, if you set ramHotPlug to TRUE then you must use a minimum of 1024 MB. If you set the RAM size more than 240GB, then ramHotPlug will be set to FALSE and can not be set to TRUE unless RAM size not set to less than 240GB.",
//				Required:    true,
//			},
//			"volumes": {
//				Type:        schema.TypeList,
//				Description: "List of volumes associated with this Template. Only a single volume is currently supported.",
//				MaxItems:    1,
//				Optional:    true,
//				Elem: &schema.Resource{
//					Schema: map[string]*schema.Schema{
//						"image": {
//							Type:        schema.TypeString,
//							Description: "Image installed on the volume. Only UUID of the image is supported currently.",
//							Required:    true,
//							ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
//								v := val.(string)
//								if !IsValidUUID(v) {
//									errs = append(errs, fmt.Errorf("%q must have the pattern \"[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{12}\" and got: %s", key, v))
//								}
//								return
//							},
//						},
//						"image_password": {
//							Type:        schema.TypeString,
//							Description: "Image password for this template volume.",
//							Optional:    true,
//						},
//						"name": {
//							Type:         schema.TypeString,
//							Description:  "Name of the template volume.",
//							Optional:     true,
//							ValidateFunc: isValidName,
//						},
//						"size": {
//							Type:        schema.TypeInt,
//							Description: "User-defined size for this template volume in GB.",
//							Optional:    true,
//							ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
//								v := val.(int)
//								if v < 1 {
//									errs = append(errs, fmt.Errorf("%q must be at least 1, got: %d", key, v))
//								}
//								return
//							},
//						},
//						"ssh_keys": {
//							Type:        schema.TypeList,
//							Description: "Ssh keys that has access to the volume.",
//							Optional:    true,
//							Elem: &schema.Schema{
//								Type: schema.TypeString,
//							},
//						},
//						"type": {
//							Type:         schema.TypeString,
//							Description:  "Storage Type for this template volume (SSD or HDD).",
//							Required:     true,
//							ValidateFunc: validation.All(validation.StringIsNotWhiteSpace),
//						},
//						"user_data": {
//							Type:        schema.TypeString,
//							Description: "user-data (Cloud Init) for this template volume.",
//							Optional:    true,
//						},
//					},
//				},
//			},
//		},
//		Timeouts: &resourceDefaultTimeouts,
//	}
//}
//
//func resourceAutoscalingTemplateCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
//
//	client := meta.(SdkBundle).AutoscalingClient
//
//	template := autoscaling.Template{
//		Properties: &autoscaling.TemplateProperties{},
//	}
//
//	if availabilityZone, availabilityZoneOk := d.GetOk("availability_zone"); availabilityZoneOk {
//		availabilityZone := autoscaling.AvailabilityZone(availabilityZone.(string))
//		template.Properties.AvailabilityZone = &availabilityZone
//	}
//
//	if cores, coresOk := d.GetOk("cores"); coresOk {
//		cores := int32(cores.(int))
//		template.Properties.Cores = &cores
//	}
//
//	if cpuFamily, cpuFamilyOk := d.GetOk("cpu_family"); cpuFamilyOk {
//		cpuFamily := autoscaling.CpuFamily(cpuFamily.(string))
//		template.Properties.CpuFamily = &cpuFamily
//	}
//
//	if location, locationOk := d.GetOk("location"); locationOk {
//		location := location.(string)
//		template.Properties.Location = &location
//	}
//
//	if name, nameOk := d.GetOk("name"); nameOk {
//		name := name.(string)
//		template.Properties.Name = &name
//	}
//
//	if nicsVal, nicsOk := d.GetOk("nics"); nicsOk {
//		nicsVal := nicsVal.([]interface{})
//		if nicsVal != nil {
//			createNics := false
//
//			var nics []autoscaling.TemplateNic
//
//			for nicIndex := range nicsVal {
//				nic := autoscaling.TemplateNic{}
//				addNic := false
//
//				if lan, lanOk := d.GetOk(fmt.Sprintf("nics.%d.lan", nicIndex)); lanOk {
//					lan := int32(lan.(int))
//					nic.Lan = &lan
//					addNic = true
//				}
//				if name, nameOk := d.GetOk(fmt.Sprintf("nics.%d.name", nicIndex)); nameOk {
//					name := name.(string)
//					nic.Name = &name
//				}
//
//				if addNic {
//					nics = append(nics, nic)
//				}
//
//			}
//
//			if len(nics) > 0 {
//				createNics = true
//			}
//
//			if createNics == true {
//				template.Properties.Nics = &nics
//			}
//		}
//	}
//
//	if ram, ramOk := d.GetOk("ram"); ramOk {
//		ram := int32(ram.(int))
//		template.Properties.Ram = &ram
//	}
//
//	if volumesVal, volumeOk := d.GetOk("volumes"); volumeOk {
//		volumesVal := volumesVal.([]interface{})
//		if volumesVal != nil {
//			createVolumes := false
//
//			var volumes []autoscaling.TemplateVolume
//
//			for volumeIndex := range volumesVal {
//				volume := autoscaling.TemplateVolume{}
//				addVolume := false
//
//				if image, imageOk := d.GetOk(fmt.Sprintf("volumes.%d.image", volumeIndex)); imageOk {
//					image := image.(string)
//					volume.Image = &image
//					addVolume = true
//				}
//
//				if imagePassword, imagePasswordOk := d.GetOk(fmt.Sprintf("volumes.%d.image_password", volumeIndex)); imagePasswordOk {
//					imagePassword := imagePassword.(string)
//					volume.ImagePassword = &imagePassword
//				}
//
//				if name, nameOk := d.GetOk(fmt.Sprintf("volumes.%d.name", volumeIndex)); nameOk {
//					name := name.(string)
//					volume.Name = &name
//				}
//
//				if size, sizeOk := d.GetOk(fmt.Sprintf("volumes.%d.size", volumeIndex)); sizeOk {
//					size := int32(size.(int))
//					volume.Size = &size
//				}
//
//				if sshKeys, sshKeysOk := d.GetOk(fmt.Sprintf("volumes.%d.ssh_keys", volumeIndex)); sshKeysOk {
//					sshKeys := sshKeys.([]interface{})
//					if sshKeys != nil {
//						sshKeysEntry := make([]string, 0)
//						for _, sshKey := range sshKeys {
//							sshKeysEntry = append(sshKeysEntry, sshKey.(string))
//						}
//						volume.SshKeys = &sshKeysEntry
//					}
//				}
//
//				if name, nameOk := d.GetOk(fmt.Sprintf("volumes.%d.name", volumeIndex)); nameOk {
//					name := name.(string)
//					volume.Name = &name
//				}
//
//				if name, nameOk := d.GetOk(fmt.Sprintf("volumes.%d.name", volumeIndex)); nameOk {
//					name := name.(string)
//					volume.Name = &name
//				}
//
//				if addVolume {
//					volumes = append(volumes, volume)
//				}
//
//			}
//
//			if len(volumes) > 0 {
//				createVolumes = true
//			}
//
//			if createVolumes {
//				template.Properties.Volumes = &volumes
//			}
//		}
//
//	}
//
//	autoscalingTemplate, _, err := client.TemplatesApi.AutoscalingTemplatesPost(ctx).Template(template).Execute()
//
//	if err != nil {
//		d.SetId("")
//		diags := diag.FromErr(fmt.Errorf("error creating autoscaling template: %s", err))
//		return diags
//	}
//
//	d.SetId(*autoscalingTemplate.Id)
//
//	return resourceAutoscalingTemplateRead(ctx, d, meta)
//}
//
//func resourceAutoscalingTemplateRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
//
//	client := meta.(SdkBundle).AutoscalingClient
//
//	template, apiResponse, err := client.TemplatesApi.AutoscalingTemplatesFindById(ctx, d.Id()).Execute()
//
//	if err != nil {
//		log.Printf("[INFO] Resource %s not found: %+v", d.Id(), err)
//		if apiResponse.StatusCode == 404 {
//			d.SetId("")
//			return nil
//		}
//	}
//
//	log.Printf("[INFO] Successfully retreived autoscaling template %s: %+v", d.Id(), template)
//
//	if template.Properties.AvailabilityZone != nil {
//		if err := d.Set("availability_zone", *template.Properties.AvailabilityZone); err != nil {
//			diags := diag.FromErr(fmt.Errorf("error while setting availability_zone property for autoscaling group %s: %s", d.Id(), err))
//			return diags
//		}
//	}
//
//	if template.Properties.Cores != nil {
//		if err := d.Set("cores", *template.Properties.Cores); err != nil {
//			diags := diag.FromErr(fmt.Errorf("error while setting cores property for autoscaling template %s: %s", d.Id(), err))
//			return diags
//		}
//	}
//
//	if template.Properties.CpuFamily != nil {
//		if err := d.Set("cpu_family", *template.Properties.CpuFamily); err != nil {
//			diags := diag.FromErr(fmt.Errorf("error while setting cpu_family property for autoscaling template %s: %s", d.Id(), err))
//			return diags
//		}
//	}
//
//	if template.Properties.Location != nil {
//		if err := d.Set("location", *template.Properties.Location); err != nil {
//			diags := diag.FromErr(fmt.Errorf("error while setting location property for autoscaling template %s: %s", d.Id(), err))
//			return diags
//		}
//	}
//
//	if template.Properties.Name != nil {
//		if err := d.Set("name", *template.Properties.Name); err != nil {
//			diags := diag.FromErr(fmt.Errorf("error while setting name property for autoscaling template %s: %s", d.Id(), err))
//			return diags
//		}
//	}
//
//	if template.Properties.Nics != nil && len(*template.Properties.Nics) > 0 {
//		var nics []interface{}
//
//		for _, nic := range *template.Properties.Nics {
//			nicEntry := make(map[string]interface{})
//			if nic.Lan != nil {
//				nicEntry["lan"] = *nic.Lan
//			}
//			if nic.Name != nil {
//				nicEntry["name"] = *nic.Name
//			}
//			nics = append(nics, nicEntry)
//		}
//
//		if len(nics) > 0 {
//			if err := d.Set("nics", nics); err != nil {
//				diags := diag.FromErr(fmt.Errorf("error while setting nics property for autoscaling template %s: %s", d.Id(), err))
//				return diags
//			}
//		}
//	}
//
//	if template.Properties.Ram != nil {
//		if err := d.Set("ram", *template.Properties.Ram); err != nil {
//			diags := diag.FromErr(fmt.Errorf("error while setting ram property for autoscaling template %s: %s", d.Id(), err))
//			return diags
//		}
//	}
//
//	if template.Properties.Volumes != nil && len(*template.Properties.Volumes) > 0 {
//		var volumes []interface{}
//		for _, volume := range *template.Properties.Volumes {
//			volumeEntry := make(map[string]interface{})
//			if volume.Image != nil {
//				volumeEntry["image"] = *volume.Image
//			}
//			if volume.ImagePassword != nil {
//				volumeEntry["image_password"] = *volume.ImagePassword
//			}
//			if volume.Name != nil {
//				volumeEntry["name"] = *volume.Name
//			}
//			if volume.Size != nil {
//				volumeEntry["size"] = *volume.Size
//			}
//			if volume.SshKeys != nil {
//				volumeEntry["ssh_keys"] = *volume.SshKeys
//			}
//			if volume.Type != nil {
//				volumeEntry["type"] = *volume.Type
//			}
//			if volume.UserData != nil {
//				volumeEntry["user_data"] = *volume.UserData
//			}
//			volumes = append(volumes, volumeEntry)
//		}
//
//		if len(volumes) > 0 {
//			if err := d.Set("volumes", volumes); err != nil {
//				diags := diag.FromErr(fmt.Errorf("error while setting volumes property for autoscaling template %s: %s", d.Id(), err))
//				return diags
//			}
//		}
//	}
//
//	return nil
//}
//
//func resourceAutoscalingTemplateUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
//
//	return nil
//}
//
//func resourceAutoscalingTemplateDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
//	client := meta.(SdkBundle).AutoscalingClient
//
//	_, err := client.TemplatesApi.AutoscalingTemplatesDelete(ctx, d.Id()).Execute()
//
//	if err != nil {
//		diags := diag.FromErr(fmt.Errorf("an error occured while deleting an atuoscaling template %s %s", d.Id(), err))
//		return diags
//	}
//
//	d.SetId("")
//
//	return nil
//}
