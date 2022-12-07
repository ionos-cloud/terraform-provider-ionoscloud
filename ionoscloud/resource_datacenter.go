package ionoscloud

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	ionoscloud "github.com/ionos-cloud/sdk-go/v6"
	"log"
)

func resourceDatacenter() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDatacenterCreate,
		ReadContext:   resourceDatacenterRead,
		UpdateContext: resourceDatacenterUpdate,
		DeleteContext: resourceDatacenterDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceDatacenterImport,
		},
		Schema: map[string]*schema.Schema{

			//Datacenter parameters
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.All(validation.StringIsNotWhiteSpace),
			},
			"location": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.All(validation.StringIsNotWhiteSpace),
			},
			"description": {
				Type:        schema.TypeString,
				Description: "A description for the datacenter, e.g. staging, production",
				Optional:    true,
				Computed:    true,
			},
			"sec_auth_protection": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"version": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"features": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"cpu_architecture": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"cpu_family": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"max_cores": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"max_ram": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"vendor": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
		Timeouts: &resourceDefaultTimeouts,
	}
}

func resourceDatacenterCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	client := meta.(SdkBundle).CloudApiClient

	datacenterName := d.Get("name").(string)
	datacenterLocation := d.Get("location").(string)

	datacenter := ionoscloud.Datacenter{
		Properties: &ionoscloud.DatacenterProperties{
			Name:     &datacenterName,
			Location: &datacenterLocation,
		},
	}

	if attr, ok := d.GetOk("description"); ok {
		attrStr := attr.(string)
		datacenter.Properties.Description = &attrStr
	}

	if attr, ok := d.GetOk("sec_auth_protection"); ok {
		attrStr := attr.(bool)
		datacenter.Properties.SecAuthProtection = &attrStr
	}

	createdDatacenter, apiResponse, err := client.DataCentersApi.DatacentersPost(ctx).Datacenter(datacenter).Execute()
	logApiRequestTime(apiResponse)

	if err != nil {
		diags := diag.FromErr(fmt.Errorf("error creating data center (%s) (%s)", d.Id(), err))
		return diags
	}
	d.SetId(*createdDatacenter.Id)

	log.Printf("[INFO] DataCenter Id: %s", d.Id())

	// Wait, catching any errors
	_, errState := getStateChangeConf(meta, d, apiResponse.Header.Get("Location"), schema.TimeoutCreate).WaitForStateContext(ctx)

	if errState != nil {
		if IsRequestFailed(err) {
			// Request failed, so resource was not created, delete resource from state file
			d.SetId("")
		}
		diags := diag.FromErr(errState)
		return diags
	}

	return resourceDatacenterRead(ctx, d, meta)
}

func resourceDatacenterRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	client := meta.(SdkBundle).CloudApiClient

	datacenter, apiResponse, err := client.DataCentersApi.DatacentersFindById(ctx, d.Id()).Execute()
	logApiRequestTime(apiResponse)

	if err != nil {
		if httpNotFound(apiResponse) {
			d.SetId("")
			return nil
		}
		diags := diag.FromErr(fmt.Errorf("error while fetching a data center ID %s %w", d.Id(), err))
		return diags
	}

	if err := setDatacenterData(d, &datacenter); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceDatacenterUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	client := meta.(SdkBundle).CloudApiClient
	obj := ionoscloud.DatacenterProperties{}

	if d.HasChange("name") {
		_, newName := d.GetChange("name")
		newNameStr := newName.(string)
		obj.Name = &newNameStr
	}

	if d.HasChange("description") {
		_, newDescription := d.GetChange("description")
		newDescriptionStr := newDescription.(string)
		obj.Description = &newDescriptionStr
	}

	if d.HasChange("location") {
		oldLocation, newLocation := d.GetChange("location")
		diags := diag.FromErr(fmt.Errorf("data center is created in %s location. You can not change location of the data center to %s; it requires recreation of the data center", oldLocation, newLocation))
		return diags
	}

	if d.HasChange("sec_auth_protection") {
		_, newSecAuthProtection := d.GetChange("sec_auth_protection")
		newSecAuthProtectionStr := newSecAuthProtection.(bool)
		obj.SecAuthProtection = &newSecAuthProtectionStr
	}
	_, apiResponse, err := client.DataCentersApi.DatacentersPatch(ctx, d.Id()).Datacenter(obj).Execute()
	logApiRequestTime(apiResponse)

	if err != nil {
		diags := diag.FromErr(fmt.Errorf("an error occured while update the data center ID %s %w", d.Id(), err))
		return diags
	}

	// Wait, catching any errors
	_, errState := getStateChangeConf(meta, d, apiResponse.Header.Get("Location"), schema.TimeoutUpdate).WaitForStateContext(ctx)
	if errState != nil {
		diags := diag.FromErr(errState)
		return diags
	}

	return resourceDatacenterRead(ctx, d, meta)
}

func resourceDatacenterDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	client := meta.(SdkBundle).CloudApiClient

	apiResponse, err := client.DataCentersApi.DatacentersDelete(ctx, d.Id()).Execute()
	logApiRequestTime(apiResponse)

	if err != nil {
		diags := diag.FromErr(err)
		return diags
	}

	// Wait, catching any errors
	_, errState := getStateChangeConf(meta, d, apiResponse.Header.Get("Location"), schema.TimeoutDelete).WaitForStateContext(ctx)
	if errState != nil {
		diags := diag.FromErr(errState)
		return diags
	}

	d.SetId("")
	return nil
}

func resourceDatacenterImport(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	client := meta.(SdkBundle).CloudApiClient

	dcId := d.Id()

	datacenter, apiResponse, err := client.DataCentersApi.DatacentersFindById(ctx, d.Id()).Execute()
	logApiRequestTime(apiResponse)

	if err != nil {
		if httpNotFound(apiResponse) {
			d.SetId("")
			return nil, fmt.Errorf("unable to find datacenter %q", dcId)
		}
		return nil, fmt.Errorf("an error occured while retrieving the datacenter %q, %q", dcId, err)
	}

	log.Printf("[INFO] Datacenter found: %+v", datacenter)

	if err := setDatacenterData(d, &datacenter); err != nil {
		return nil, err
	}

	return []*schema.ResourceData{d}, nil
}

func setDatacenterData(d *schema.ResourceData, datacenter *ionoscloud.Datacenter) error {

	if datacenter.Id != nil {
		d.SetId(*datacenter.Id)
	}

	if datacenter.Properties != nil {
		if datacenter.Properties.Location != nil {
			err := d.Set("location", *datacenter.Properties.Location)
			if err != nil {
				return fmt.Errorf("error while setting location property for datacenter %s: %w", d.Id(), err)
			}
		}

		if datacenter.Properties.Description != nil {
			err := d.Set("description", *datacenter.Properties.Description)
			if err != nil {
				return fmt.Errorf("error while setting description property for datacenter %s: %w", d.Id(), err)
			}
		}

		if datacenter.Properties.Name != nil {
			err := d.Set("name", *datacenter.Properties.Name)
			if err != nil {
				return fmt.Errorf("error while setting name property for datacenter %s: %w", d.Id(), err)
			}
		}

		if datacenter.Properties.Version != nil {
			err := d.Set("version", *datacenter.Properties.Version)
			if err != nil {
				return fmt.Errorf("error while setting version property for datacenter %s: %w", d.Id(), err)
			}
		}

		if datacenter.Properties.Features != nil && len(*datacenter.Properties.Features) > 0 {
			err := d.Set("features", *datacenter.Properties.Features)
			if err != nil {
				return fmt.Errorf("error while setting features property for datacenter %s: %w", d.Id(), err)
			}
		}

		if datacenter.Properties.SecAuthProtection != nil {
			err := d.Set("sec_auth_protection", *datacenter.Properties.SecAuthProtection)
			if err != nil {
				return fmt.Errorf("error while setting sec_auth_protection property for datacenter %s: %w", d.Id(), err)
			}
		}

		if datacenter.Properties.CpuArchitecture != nil && len(*datacenter.Properties.CpuArchitecture) > 0 {
			var cpuArchitectures []interface{}
			for _, cpuArchitecture := range *datacenter.Properties.CpuArchitecture {
				architectureEntry := make(map[string]interface{})

				if cpuArchitecture.CpuFamily != nil {
					architectureEntry["cpu_family"] = *cpuArchitecture.CpuFamily
				}

				if cpuArchitecture.MaxCores != nil {
					architectureEntry["max_cores"] = *cpuArchitecture.MaxCores
				}

				if cpuArchitecture.MaxRam != nil {
					architectureEntry["max_ram"] = *cpuArchitecture.MaxRam
				}

				if cpuArchitecture.Vendor != nil {
					architectureEntry["vendor"] = *cpuArchitecture.Vendor
				}

				cpuArchitectures = append(cpuArchitectures, architectureEntry)

				if len(cpuArchitectures) > 0 {
					if err := d.Set("cpu_architecture", cpuArchitectures); err != nil {
						return fmt.Errorf("error while setting cpu_architecture property for datacenter %s: %w", d.Id(), err)
					}
				}
			}
		}
	}

	return nil
}
