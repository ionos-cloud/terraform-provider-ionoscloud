// Package datacenter holds the SDKv2-compatible schema, identity schema, and
// field-mapping logic for the ionoscloud_datacenter resource, shared between
// the legacy ionoscloud package (which owns the CRUD resource) and the
// framework-based compute package (which builds the ionoscloud_datacenter
// list resource on top of the same schema). Neither of those packages can
// import the other without an import cycle, so this dependency-neutral
// package exists as the single place both can import instead.
package datacenter

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	ionoscloud "github.com/ionos-cloud/sdk-go/v6"
)

// Schema returns the SDKv2 schema for the ionoscloud_datacenter resource.
func Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"name": {
			Type:             schema.TypeString,
			Required:         true,
			ValidateDiagFunc: validation.ToDiagFunc(validation.StringIsNotWhiteSpace),
		},
		"location": {
			Type:     schema.TypeString,
			Required: true,
			ForceNew: true,
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
		"ipv6_cidr_block": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Auto-assigned /56 IPv6 CIDR block, if IPv6 is enabled for the datacenter. Read-only",
		},
	}
}

// IdentitySchema returns the SDKv2 resource identity schema for the
// ionoscloud_datacenter resource.
func IdentitySchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"id": {
			Type:              schema.TypeString,
			RequiredForImport: true,
		},
	}
}

// PopulateResourceData flattens a Cloud API Datacenter into d, using the
// field names defined by Schema.
func PopulateResourceData(d *schema.ResourceData, datacenter *ionoscloud.Datacenter) error {

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
			var cpuArchitectures []any
			for _, cpuArchitecture := range *datacenter.Properties.CpuArchitecture {
				architectureEntry := make(map[string]any)

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

		if datacenter.Properties.Ipv6CidrBlock != nil {
			err := d.Set("ipv6_cidr_block", *datacenter.Properties.Ipv6CidrBlock)
			if err != nil {
				return fmt.Errorf("error while setting ipv6_cidr_block property for datacenter %s: %w", d.Id(), err)
			}
		}

	}

	return nil
}
