package ionoscloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceAutoscalingTemplate_matchId(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAutoscalingTemplateCreateResources,
			},
			{
				Config: testAccDataSourceAutoscalingTemplateMatchId,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.ionoscloud_autoscaling_template.autoscaling_template", "name", "test_ds_autoscaling_template"),
				),
			},
		},
	})
}

func TestAccDataSourceAutoscalingTemplate_matchName(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAutoscalingTemplateCreateResources,
			},
			{
				Config: testAccDataSourceAutoscalingTemplateMatchName,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.ionoscloud_autoscaling_template.autoscaling_template", "name", "test_ds_autoscaling_template"),
				),
			},
		},
	})

}

const testAccDataSourceAutoscalingTemplateCreateResources = `
resource "ionoscloud_datacenter" "autoscaling_template" {
   name     = "test_autoscaling_template"
   location = "de/txl"
}

resource "ionoscloud_lan" "autoscaling_template" {
	datacenter_id    = ionoscloud_datacenter.autoscaling_template.id
    public           = false
    name             = "test_autoscaling_template"
}

resource "ionoscloud_autoscaling_template" "autoscaling_template" {
	availability_zone    = "AUTO"
    cores				 = 2
	cpu_family           = "INTEL_SKYLAKE"
	location			 = "de/txl"
    name                 = "test_ds_autoscaling_template"
    nics    {
		lan              = ionoscloud_lan.autoscaling_template.id
        name             = "test_autoscaling_template"
    }
    ram                  = 1024
	volumes  {
    	image            = "e309f108-b48d-11eb-b9b3-d2869b2d44d9"
		image_password   = "test12345678"
        name             = "test_autoscaling_template"
		size             = 50
    	type             = "HDD"
	}
}
`

const testAccDataSourceAutoscalingTemplateMatchId = `
resource "ionoscloud_datacenter" "autoscaling_template" {
   name     = "test_autoscaling_template"
   location = "de/txl"
}

resource "ionoscloud_lan" "autoscaling_template" {
	datacenter_id    = ionoscloud_datacenter.autoscaling_template.id
    public           = false
    name             = "test_autoscaling_template"
}

resource "ionoscloud_autoscaling_template" "autoscaling_template" {
	availability_zone    = "AUTO"
    cores				 = 2
	cpu_family           = "INTEL_SKYLAKE"
	location			 = "de/txl"
    name                 = "test_ds_autoscaling_template"
    nics    {
		lan              = ionoscloud_lan.autoscaling_template.id
        name             = "test_autoscaling_template"
    }
    ram                  = 1024
	volumes  {
    	image            = "e309f108-b48d-11eb-b9b3-d2869b2d44d9"
		image_password   = "test12345678"
        name             = "test_autoscaling_template"
		size             = 50
    	type             = "HDD"
	}
}

data "ionoscloud_autoscaling_template" "autoscaling_template" {
  id			= ionoscloud_autoscaling_template.autoscaling_template.id
}
`

const testAccDataSourceAutoscalingTemplateMatchName = `
resource "ionoscloud_datacenter" "autoscaling_template" {
   name     = "test_autoscaling_template"
   location = "de/txl"
}

resource "ionoscloud_lan" "autoscaling_template" {
	datacenter_id    = ionoscloud_datacenter.autoscaling_template.id
    public           = false
    name             = "test_autoscaling_template"
}

resource "ionoscloud_autoscaling_template" "autoscaling_template" {
	availability_zone    = "AUTO"
    cores				 = 2
	cpu_family           = "INTEL_SKYLAKE"
	location			 = "de/txl"
    name                 = "test_ds_autoscaling_template"
    nics    {
		lan              = ionoscloud_lan.autoscaling_template.id
        name             = "test_autoscaling_template"
    }
    ram                  = 1024
	volumes  {
    	image            = "e309f108-b48d-11eb-b9b3-d2869b2d44d9"
		image_password   = "test12345678"
        name             = "test_autoscaling_template"
		size             = 50
    	type             = "HDD"
	}
}

data "ionoscloud_autoscaling_template" "autoscaling_template" {
  name			= "test_ds"
}
`
