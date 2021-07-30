package ionoscloud

import (
	"context"
	"fmt"
	autoscaling "github.com/ionos-cloud/sdk-go-autoscaling"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccAutoscalingTemplate_Basic(t *testing.T) {
	var autoscalingTemplate autoscaling.Template
	autoscalingTemplateName := "test_autoscaling_template"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckAutoscalingTemplateDestroyCheck,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccCheckAutoscalingTemplateConfigBasic, autoscalingTemplateName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAutoscalingTemplateExists("ionoscloud_autoscaling_template.autoscaling_template", &autoscalingTemplate),
					resource.TestCheckResourceAttr("ionoscloud_autoscaling_template.autoscaling_template", "name", autoscalingTemplateName),
				),
			},
		},
	})
}

func testAccCheckAutoscalingTemplateDestroyCheck(s *terraform.State) error {
	client := testAccProvider.Meta().(SdkBundle).AutoscalingClient

	ctx, cancel := context.WithTimeout(context.Background(), *resourceDefaultTimeouts.Default)

	if cancel != nil {
		defer cancel()
	}

	for _, rs := range s.RootModule().Resources {

		if rs.Type != "ionoscloud_autoscaling_template" {
			continue
		}

		_, apiResponse, err := client.TemplatesApi.AutoscalingTemplatesFindById(ctx, rs.Primary.ID).Execute()

		if err != nil {
			if apiResponse == nil || apiResponse.StatusCode != 404 {
				return fmt.Errorf("an error occurred while checking for the destruction of autoscaling template %s: %s",
					rs.Primary.ID, err)
			}
		} else {
			return fmt.Errorf("autoscaling template %s still exists", rs.Primary.ID)
		}

	}

	return nil
}

func testAccCheckAutoscalingTemplateExists(n string, autoscalingTemplate *autoscaling.Template) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(SdkBundle).AutoscalingClient

		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no Record ID is set")
		}

		ctx, cancel := context.WithTimeout(context.Background(), *resourceDefaultTimeouts.Default)

		if cancel != nil {
			defer cancel()
		}

		foundTemplate, _, err := client.TemplatesApi.AutoscalingTemplatesFindById(ctx, rs.Primary.ID).Execute()

		if err != nil {
			return fmt.Errorf("error occured while fetching backup unit: %s", rs.Primary.ID)
		}
		if *foundTemplate.Id != rs.Primary.ID {
			return fmt.Errorf("record not found")
		}
		autoscalingTemplate = &foundTemplate

		return nil
	}
}

const testAccCheckAutoscalingTemplateConfigBasic = `
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
    name                 = "%s"
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
